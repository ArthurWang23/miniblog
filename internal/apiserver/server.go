package apiserver

// 运行时代码保存在internal/apiserver
// 初始化配置正确加载并读取后，基于初始化配置创建运行时配置，并基于运行时配置创建服务器实例
// 采用面向对象风格UnionServer结构体封装服务相关功能

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/spf13/viper"
)

type Config struct {
	ServerMode string
	JWTKey     string
	Expiration time.Duration
}

type UnionServer struct {
	cfg *Config
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	return &UnionServer{cfg: cfg}, nil
}

func (s *UnionServer) Run() error {
	log.Infow("ServerMode from ServerOptions", "jwt-key", s.cfg.JWTKey)
	log.Infow("ServerMode from Viper", "jwt-key", viper.GetString("jwt-key"))
	jsonData, _ := json.MarshalIndent(s.cfg, "", " ")
	fmt.Println(string(jsonData))
	select {}
	return nil
}
