package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 创建 *cobra.Command 对象用于启动应用程序
func NewMiniBlogCommand() *cobra.Command {
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
			fmt.Println("Hello MiniBlog!")
			return nil
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数
		Args: cobra.NoArgs,
	}
	return cmd
}
