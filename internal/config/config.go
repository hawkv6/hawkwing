package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	// Use \ as key delimiter because we have . in the key
	viperInstance = viper.NewWithOptions(viper.KeyDelimiter("\\"))
	Params        Config
)

type HawkEyeConfig struct {
	Hostname string `mapstructure:"hostname"`
	Port     int    `mapstructure:"port"`
}

type Intent struct {
	Intent   string   `mapstructure:"intent"`
	Port     int      `mapstructure:"port"`
	MinValue int      `mapstructure:"min_value"`
	MaxValue int      `mapstructure:"max_value"`
	Sfc      []string `mapstructure:"sfc"`
	FlexAlgo int      `mapstructure:"flex_algo"`
	Sid      []string `mapstructure:"sid"`
}

type Application struct {
	Port    int      `mapstructure:"port"`
	Sid     []string `mapstructure:"sid"`
	Intents []Intent `mapstructure:"intents"`
}

type ServiceConfig struct {
	DomainName    string        `mapstructure:"domain_name"`
	Ipv6Addresses []string      `mapstructure:"ipv6_addresses"`
	Applications  []Application `mapstructure:"applications"`
}

type Config struct {
	HawkEye  HawkEyeConfig
	Services map[string]ServiceConfig
}

func init() {
	viperInstance.SetEnvPrefix("HAWKWING")
	viperInstance.AutomaticEnv()
}

func Parse() error {
	if len(viperInstance.ConfigFileUsed()) != 0 {
		if err := viperInstance.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to load config file %s: %v", viperInstance.ConfigFileUsed(), err)
		}
	}

	if err := viperInstance.UnmarshalExact(&Params); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	return nil
}

func GetInstance() *viper.Viper {
	return viperInstance
}
