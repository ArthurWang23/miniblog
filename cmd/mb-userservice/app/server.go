package app

import (
	"github.com/ArthurWang23/miniblog/cmd/mb-userservice/app/options"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	kratoslog "github.com/ArthurWang23/miniblog/pkg/log"
	"github.com/ArthurWang23/miniblog/pkg/registry"
	"github.com/ArthurWang23/miniblog/pkg/token"
	"github.com/ArthurWang23/miniblog/pkg/version"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/userservice"
	"github.com/go-kratos/kratos/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func NewUserServiceCommand() *cobra.Command {
	opts := options.NewServerOptions()
	cmd := &cobra.Command{
		Use:   "mb-userservice",
		Short: "User service of miniblog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKratos(opts)
		},
		Args: cobra.NoArgs,
	}
	// 在命令初始化阶段加载配置与环境变量
	cobra.OnInitialize(onInitialize)

	// 增加 --config/-c 标志，支持自定义配置文件路径
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the userservice configuration file.")

	opts.AddFlags(cmd.PersistentFlags())
	version.AddFlags(cmd.PersistentFlags())

	return cmd
}

func logOptions() *log.Options {
	opts := log.NewOptions()
	// 可按需通过 viper 读取日志相关配置
	return opts
}

func runKratos(opts *options.ServerOptions) error {
	version.PrintAndExitIfRequested()
	log.Init(logOptions())
	defer log.Sync()

	if err := viper.Unmarshal(opts); err != nil {
		return err
	}
	if err := opts.Validate(); err != nil {
		return err
	}
	cfg := opts.Config()

	// 初始化 JWT（为后续接入认证中间件做准备）
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)

	// 构建 ServerConfig
	sc := userservice.NewServerConfig(cfg)

	// 构建 Kratos 传输服务器（当前仅 HTTP，可按需扩展 gRPC/gateway）
	kservers, err := sc.NewKratosServers()
	if err != nil {
		return err
	}

	// etcd 注册（可选）
	var registrar kratos.Registrar
	if cfg.EtcdOptions != nil && len(cfg.EtcdOptions.Endpoints) > 0 {
		reg, _, err := registry.NewEtcdRegistryWithOptions(cfg.EtcdOptions)
		if err != nil {
			return err
		}
		registrar = reg
	}

	appOpts := []kratos.Option{
		kratos.Name("userservice"),
		kratos.Version(version.Get().GitVersion),
		kratos.Logger(kratoslog.Kratos()),
	}
	for _, s := range kservers {
		appOpts = append(appOpts, kratos.Server(s))
	}
	if registrar != nil {
		appOpts = append(appOpts, kratos.Registrar(registrar))
	}

	app := kratos.New(appOpts...)
	return app.Run()
}