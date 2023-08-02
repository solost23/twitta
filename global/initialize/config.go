package initialize

import (
	"fmt"
	"os"

	"twitta/global"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	provider = "consul"

	consulDefaultHost = "127.0.0.1"
	consulDefaultPort = 8500
)

func InitConfig(configFilePath string) {
	v := viper.New()
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = consulDefaultHost
	}
	v.SetDefault("consul.host", consulHost)
	v.SetDefault("consul.port", consulDefaultPort)

	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	// 从配置中心读取配置
	err := v.AddRemoteProvider(provider,
		fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port),
		global.ServerConfig.ConfigPath)
	if err != nil {
		panic(err)
	}

	v.SetConfigType("YAML")

	if err = v.ReadRemoteConfig(); err != nil {
		panic(err)
	}

	if err = v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
}
