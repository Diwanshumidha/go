package key

import (
	"errors"
	"fmt"
	"tmdb/internal/tmdb"

	"github.com/zalando/go-keyring"
)

const (
	service = "tmdb-cli"
	keyName = "tmdb-api-key"
)

func SaveAPIKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	return keyring.Set(service, keyName, key)
}

func ValidateAPIKey(key string) error {
	if key == "" {
		return errors.New("API key cannot be empty")
	}

	if err := tmdb.ValidateKey(key); err != nil {
		return err
	}

	return nil
}

func GetAPIKey() (string, error) {
	return keyring.Get(service, keyName)
}

func DeleteAPIKey() error {
	return keyring.Delete(service, keyName)
}
