// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package app

import (
	"github.com/ArthurWang23/miniblog/cmd/mb-apiserver/app/options"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/pkg/version"
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
			return run(opts)
		},
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(onInitialize)

	// 持久性标志 可用于它所分配的命令以及命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the miniblog configuration file.")
	opts.AddFlags(cmd.PersistentFlags())

	version.AddFlags(cmd.PersistentFlags())
	return cmd
}

func run(opts *options.ServerOptions) error {
	version.PrintAndExitIfRequested()
	log.Init(logOptions())
	defer log.Sync()
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	if err := opts.Validate(); err != nil {
		return err
	}
	// 获取应用配置
	// 将命令行选项和应用配置分开，可以更加灵活地处理2中不同类型配置
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	server, err := cfg.NewUnionServer()
	if err != nil {
		return err
	}
	return server.Run()
}

func logOptions() *log.Options {
	opts := log.NewOptions()

	if viper.IsSet("log.disable-caller") {
		opts.DisableCaller = viper.GetBool("log.disable-caller")
	}

	if viper.IsSet("log.disable-stacktrace") {
		opts.DisableStacktrace = viper.GetBool("log.disable-stacktrace")
	}

	if viper.IsSet("log.level") {
		opts.Level = viper.GetString("log.level")
	}

	if viper.IsSet("log.format") {
		opts.Format = viper.GetString("log.format")
	}

	if viper.IsSet("log.output-paths") {
		opts.OutputPaths = viper.GetStringSlice("log.output-paths")
	}

	return opts
}
