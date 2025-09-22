package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
)

func main() {
	// 1) 创建 etcd 客户端（按需替换 endpoints / 认证 / TLS）
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		// Username: "root",
		// Password: "password",
		// TLS: tlsConfig, // 如启用 TLS
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 2) 构建 discovery
	dis := etcd.New(cli)

	// 3) 通过服务名拨号（使用 Kratos gRPC 客户端）
	ctx := context.Background()
	conn, err := kratosgrpc.DialInsecure(
		ctx,
		kratosgrpc.WithDiscovery(dis),
		kratosgrpc.WithEndpoint("discovery:///miniblog"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 4) 生成客户端并调用 Healthz
	cliV1 := apiv1.NewMiniBlogClient(conn)
	resp, err := cliV1.Healthz(ctx, &apiv1.HealthzRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("healthz:", resp)
}