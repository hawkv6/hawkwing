package config

import (
	"fmt"

	"github.com/hawkv6/hawkwing/pkg/logging"
	"github.com/spf13/viper"
)

const Subsystem = "go-config"

var (
	viperInstance = viper.NewWithOptions(viper.KeyDelimiter("\\"))
	Params        Config
	log           = logging.DefaultLogger.WithField("subsystem", Subsystem)
)

type HawkEyeConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address" validate:"required,ipv6"`
	Port    int    `mapstructure:"port" validate:"required,gt=0,lt=65535"`
}

type Intent struct {
	Intent     string   `mapstructure:"intent"`
	MinValue   int      `mapstructure:"min_value"`
	MaxValue   int      `mapstructure:"max_value"`
	Functions  []string `mapstructure:"functions"`
	FlexAlgoNr int      `mapstructure:"flex_algo_number"`
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
	ClientIpv6Address string                   `mapstructure:"client_ipv6_address" validate:"ipv6"`
	HawkEye           HawkEyeConfig            `validate:"required,dive,required"`
	Services          map[string]ServiceConfig `validate:"required,dive,required"`
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

	Validate()

	return nil
}

func GetInstance() *viper.Viper {
	return viperInstance
}
