package core

import (
	"fmt"
	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	config := "config.toml"
	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", config))
	}
	//监听配置文件变化
	v.WatchConfig()

	return v
}
