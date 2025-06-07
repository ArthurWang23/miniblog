package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to. ")
	limit = flag.Int64("limit", 10, "Limit to list users.")
)

func main() {
	flag.Parse()

	// grpc建立连接
	// grpc.Dial用于建立客户端与grpc服务端的连接
	// grpc.WithTransportCredentials(insecure.NewCredentials()) 使用不安全的传输凭证(不实用TLS)
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	client := apiv1.NewMiniBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// 发起grpc请求
	resp, err := client.Healthz(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get healthz: %v", err)
	}
	jsonData, _ := json.Marshal(resp)
	fmt.Println(string(jsonData))
}
