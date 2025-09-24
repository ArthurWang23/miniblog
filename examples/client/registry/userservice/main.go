package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	pb "github.com/ArthurWang23/miniblog/pkg/api/userservice/v1"
)

func main() {
	// 1) 构建 etcd discovery
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	dis := etcd.New(cli)

	// 2) 通过服务名 "userservice" 进行 gRPC 拨号
	ctx := context.Background()
	conn, err := kratosgrpc.DialInsecure(
		ctx,
		kratosgrpc.WithDiscovery(dis),
		kratosgrpc.WithEndpoint("discovery:///userservice"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 3) 生成客户端并调用
	cliV1 := pb.NewUserServiceClient(conn)

	// 3.1 Healthz
	health, err := cliV1.Healthz(ctx, &pb.HealthzRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("healthz:", health.GetStatus())

	// 3.2 CreateUser
	_, _ = cliV1.CreateUser(ctx, &pb.CreateUserRequest{
		Username: "alice",
		Password: "alice@12345",
	})

	// 3.3 Login
	login, err := cliV1.Login(ctx, &pb.LoginRequest{
		Username: "alice",
		Password: "alice@12345",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", login.GetAccessToken(), "expires_in:", login.GetExpiresIn())
}