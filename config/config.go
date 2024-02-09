package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ID                       string
	NetWorkGroupName         string
	ListenPort               int
	Roles                    []string
	CommitteeSelectionPeriod int
	MinNodes                 int
}

func NewConfig() (Config, error) {
	viper.SetConfigFile("./config/config.yml")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("error mapping config file: %s", err)
	}

	return config, nil
}
