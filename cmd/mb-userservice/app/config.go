package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 默认配置位置与文件名
const (
	defaultHomeDir     = ".miniblog"
	defaultConfigName  = "mb-userservice.yaml"
)

func onInitialize() {
	if configFile != "" {
		// 使用命令行指定的配置文件
		viper.SetConfigFile(configFile)
	} else {
		// 默认搜索路径
		for _, dir := range searchDirs() {
			viper.AddConfigPath(dir)
		}
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	setupEnvironmentVariables()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read viper configuration file, err: %v", err)
	}
	log.Printf("Using config file: %s", viper.ConfigFileUsed())
}

func setupEnvironmentVariables() {
	// 自动读取环境变量
	viper.AutomaticEnv()
	// 为用户服务设置单独前缀，避免多服务冲突
	viper.SetEnvPrefix("USERSERVICE")
	// 替换 key 分隔符
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