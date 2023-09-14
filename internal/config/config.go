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
	Hostname string
	Port     int
}

type ServiceConfig struct {
	Intent string
	Port   int
	Sid    []string
}

type Config struct {
	HawkEye  HawkEyeConfig
	Services map[string][]ServiceConfig
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
