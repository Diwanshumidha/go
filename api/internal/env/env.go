package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Environment Variable %s not an integer Using fallback value: %d. Error: %v", key, fallback, err)
		return fallback
	}

	return intValue
}

func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment Variable %s not a boolean Using fallback value: %t. Error: %v", key, fallback, err)
		return fallback
	}

	return boolValue
}
