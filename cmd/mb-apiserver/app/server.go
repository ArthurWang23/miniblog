package app

import (
	"encoding/json"
	"fmt"

	"github.com/ArthurWang23/miniblog/cmd/mb-apiserver/app/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string // 配置文件路径
// 创建 *cobra.Command 对象用于启动应用程序
// 将ServerOptions结构体中特定字段绑定到 *cobra.Command
func NewMiniBlogCommand() *cobra.Command {
	opts := options.NewServerOptions()
	cmd := &cobra.Command{
		Use:   "mb-apiserver",
		Short: "A mini blog show best practices for develop a full-featured Go project",
		Long: `A mini blog show best practices for develop a full-featured Go project.
The project features include:
• Utilization of a clean architecture;
• Use of many commonly used Go packages: gorm, casbin, govalidator, jwt, gin, 
  cobra, viper, pflag, zap, pprof, grpc, protobuf, grpc-gateway, etc.;
• A standardized directory structure following the project-layout convention;
• Authentication (JWT) and authorization features (casbin);
• Independently designed log and error packages;
• Management of the project using a high-quality Makefile;
• Static code analysis;
• Includes unit tests, performance tests, fuzz tests, and mock tests;
• Rich web functionalities (tracing, graceful shutdown, middleware, CORS, 
  recovery from panics, etc.);
• Implementation of HTTP, HTTPS, and gRPC servers;
• Implementation of JSON and Protobuf data exchange formats;
• The project adheres to numerous development standards: 
  code standards, versioning standards, API standards, logging standards, 
  error handling standards, commit standards, etc.;
• Access to MySQL with programming implementation;
• Implemented business functionalities: user management and blog management;
• RESTful API design standards;
• OpenAPI 3.0/Swagger 2.0 API documentation;
• High-quality code.`,
		// 命令出错时，不打印帮助信息
		SilenceUsage: true,
		// 指定调用cmd.Execute()时，执行的Run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			// Unmarshal方法可以将viper加载的配置项解析到opt中
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			fmt.Printf("JWTKey from ServerOptions : %s\n", opts.JWTKey)
			fmt.Printf("JWTKey from Viper : %s\n\n", viper.GetString("jwt-key"))
			// 无前缀 空格分隔
			jsonData, _ := json.MarshalIndent(opts, "", " ")
			fmt.Println(string(jsonData))
			return nil
		},
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(onInitialize)

	// 持久性标志 可用于它所分配的命令以及命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the miniblog configuration file.")
	opts.AddFlags(cmd.PersistentFlags())
	return cmd
}
