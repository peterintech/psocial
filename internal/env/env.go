package env

import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetEnvAsInt(name string, fallback int) int {
	val, exists := os.LookupEnv(name)
	if !exists {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetEnvAsBool(name string, fallback bool) bool {
	val, exists := os.LookupEnv(name)
	if !exists {
		return fallback
	}
	valAsBool, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return valAsBool
}
