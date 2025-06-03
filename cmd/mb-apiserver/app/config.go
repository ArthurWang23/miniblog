package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 加载配置文件方法
const (
	defaultHomeDir = ".miniblog"

	defaultConfigName = "mb-apiserver.yaml"
)

// onInitialize读取配置文件名、环境变量，并将其内容读取到viper中
// 优先读configFile参数指定的文件，若空则加载默认
func onInitialize() {
	if configFile != "" {
		// 从命令行选项指定的配置文件中提取
		viper.SetConfigFile(configFile)
	} else {
		// 使用默认的配置文件路径和名称
		for _, dir := range searchDirs() {
			viper.AddConfigPath(dir)
		}
		// 配置文件类型
		viper.SetConfigType("yaml")
		// 配置文件名称
		viper.SetConfigName(defaultConfigName)
	}

	setupEnvironmentVariables()
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read viper configuration file,err:%v", err)
	}
	log.Printf("Using config file :%s", viper.ConfigFileUsed())
}

func setupEnvironmentVariables() {
	// 允许自动匹配环境变量
	viper.AutomaticEnv()
	// 设置环境变量前缀，避免同一机器多个服务环境变量名称冲突
	viper.SetEnvPrefix("MINIBLOG")
	// 替换key中的分隔符
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
}

func searchDirs() []string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return []string{filepath.Join(homeDir, defaultHomeDir)}
}

func filePath() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return filepath.Join(home, defaultHomeDir, defaultConfigName)
}
