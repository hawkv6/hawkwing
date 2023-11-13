package test

import (
	"log"
	"testing"

	"github.com/hawkv6/hawkwing/internal/config"
)

const testConfigPath = "../../test_assets/test_config.yaml"

func SetupTestConfig(tb testing.TB) {
	config.GetInstance().SetConfigFile(testConfigPath)
	if err := config.Parse(); err != nil {
		log.Fatalln(err)
	}
}
