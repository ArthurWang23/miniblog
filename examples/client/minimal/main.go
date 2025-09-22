package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:6666", "gRPC server address")
	useTLS := flag.Bool("use-tls", false, "Enable TLS")
	caCert := flag.String("cacert", "", "CA certificate file for TLS (PEM)")
	serverName := flag.String("servername", "", "Override TLS server name (SNI) if needed")

	username := flag.String("username", "", "Login username (optional)")
	password := flag.String("password", "", "Login password (optional)")
	timeout := flag.Duration("timeout", 5*time.Second, "Per-RPC timeout")

	flag.Parse()

	// 1) 建立连接
	var opts []grpc.DialOption
	if *useTLS {
		// 若为自签发证书，请提供 -cacert 指向 server.crt
		if *caCert == "" {
			log.Fatal("TLS enabled but -cacert not provided")
		}
		creds, err := credentials.NewClientTLSFromFile(*caCert, *serverName)
		if err != nil {
			log.Fatalf("Failed to load CA cert: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial %s: %v", *addr, err)
	}
	defer conn.Close()

	client := apiv1.NewMiniBlogClient(conn)

	// 2) 调用 Healthz（无需鉴权）
	{
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		resp, err := client.Healthz(ctx, nil)
		if err != nil {
			log.Fatalf("[Healthz] call failed: %v", err)
		}
		fmt.Printf("[Healthz] ok: %+v\n", resp)
	}

	// 3) 如果提供了用户名/密码，则尝试登录，获取 JWT，并演示如何在后续请求中加上 Authorization
	if *username != "" && *password != "" {
		ctx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()

		loginResp, err := client.Login(ctx, &apiv1.LoginRequest{
			Username: *username,
			Password: *password,
		})
		if err != nil {
			log.Fatalf("[Login] failed: %v", err)
		}
		if loginResp.GetToken() == "" {
			log.Fatalf("[Login] no token returned")
		}
		fmt.Printf("[Login] token: %s, expireAt: %s\n", loginResp.GetToken(), loginResp.GetExpireAt())

		// 将 token 放入 metadata，后续调用受保护接口时使用
		authCtx := metadata.NewOutgoingContext(context.Background(),
			metadata.Pairs("Authorization", "Bearer "+loginResp.GetToken()),
		)

		// 示例：带上 Authorization 再次调用 Healthz（Healthz 本身不需要鉴权，仅作演示）
		{
			ctx2, cancel2 := context.WithTimeout(authCtx, *timeout)
			defer cancel2()
			resp2, err := client.Healthz(ctx2, nil)
			if err != nil {
				log.Fatalf("[Healthz with token] failed: %v", err)
			}
			fmt.Printf("[Healthz with token] ok: %+v\n", resp2)
		}

		// 你可以在这里继续调用受保护的接口，例如：
		// ctx3 := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Authorization", "Bearer "+loginResp.GetToken()))
		// _, err = client.ListUser(ctx3, &apiv1.ListUserRequest{ ... })
		// if err != nil { ... }
	}
}