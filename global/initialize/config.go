package initialize

import (
	"twitta/global"

	"github.com/spf13/viper"
)

func InitConfig(configFilePath string) {
	v := viper.New()
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
}
