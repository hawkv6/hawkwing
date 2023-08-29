package config

import (
	"log"

	"github.com/spf13/viper"
)

var viperInstance = viper.New()

var Params Config

type Config struct {
	HawkEye struct {
		Hostname string
		Port     int
	}
	Services struct {
	}
}

func ReadConfig(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viperInstance.SetConfigFile(cfgFile)
	} else {
		// Look for config in the working directory with name "config" (without extension).
		viperInstance.AddConfigPath(".")
		viperInstance.SetConfigType("yaml")
		viperInstance.SetConfigName("config")
	}

	viperInstance.SetEnvPrefix("HAWKING")
	viperInstance.AutomaticEnv()

	err := viperInstance.ReadInConfig()
	if err != nil {
		return err
	} else {
		log.Println("Using config file:", viperInstance.ConfigFileUsed())
	}

	err = viperInstance.Unmarshal(&Params)
	if err != nil {
		return err
	}

	return err
}

func GetInstance() *viper.Viper {
	return viperInstance
}
