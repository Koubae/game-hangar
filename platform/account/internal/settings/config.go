package settings

import (
	"fmt"
	"github.com/koubae/game-hangar/account/pkg/utils"
	"os"
	"slices"
	"strconv"
)

type Config struct {
	Port           uint16
	Environment    string
	TrustedProxies []string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		panic("Config is not initialized, call NewConfig() first!")
	}
	return config
}
func NewConfig() *Config {
	port := utils.GetEnvInt("APP_PORT", 8001)

	errTemp := os.Setenv("PORT", strconv.Itoa(port)) // For gin-gonic
	if errTemp != nil {
		panic(errTemp.Error())
	}

	environment := utils.GetEnvString("APP_ENVIRONMENT", "development")
	if !slices.Contains(Environments[:], environment) {
		panic(fmt.Sprintf("Invalid environment: %s, supported envs are %v", environment, Environments))
	}
	trustedProxies := utils.GetEnvStringSlice("APP_NETWORKING_PROXIES", []string{})

	config = &Config{
		Port:           uint16(port),
		Environment:    environment,
		TrustedProxies: trustedProxies,
	}
	return config
}
