package settings

import (
	"fmt"
	"github.com/koubae/game-hangar/account/pkg/utils"
	"os"
	"slices"
	"strconv"
)

type Config struct {
	port           uint8
	Environment    string
	TrustedProxies []string
}

var config *Config

func GetConfig() *Config {
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

	return &Config{
		port:           uint8(port),
		Environment:    environment,
		TrustedProxies: trustedProxies,
	}
}
