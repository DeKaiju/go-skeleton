package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var AppRoot string

func SetAppRoot(_ string) {
	AppRoot, _ = os.Getwd()
}

func LoadEnv(configFile string) (string, error) {
	configPath := configFile
	if !filepath.IsAbs(configFile) {
		if AppRoot == "" {
			AppRoot, _ = os.Getwd()
		}

		configPath = filepath.Join(AppRoot, configFile)
	}

	if err := godotenv.Load(configPath); err != nil {
		return configPath, err
	}

	return configPath, nil
}

func GetEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func GetIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

func GetBoolEnv(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "":
		return fallback
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
