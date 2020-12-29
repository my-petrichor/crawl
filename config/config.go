package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	StateDepartmentNews     string `mapstructure: "stateDepartmentNews"`
	StateDepartmentPolicies string `mapstructure: "stateDepartmentPolicies`
	Economist               string `mapstructure: "economist`
	Academician             string `mapstructure: "academician`
}

var Conf = new(Config)

func Set() {
	viper.SetConfigFile("./config/conf/test.yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("config file could found, other error: %v\n", err)
		} else {
			fmt.Printf("config file not found: %v\n", err)
		}
	}

	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Sprintf("unmarshal error %v\n", err))
	}

	// viper.WatchConfig()

	// viper.OnConfigChange(func(in fsnotify.Event) {
	// 	if err := viper.Unmarshal(Conf); err != nil {
	// 		panic(fmt.Sprintf("unmarshal error %v\n", err))
	// 	}
	// })
}
