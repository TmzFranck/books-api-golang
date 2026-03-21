package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	viper := viper.New()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./../")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return viper

}
