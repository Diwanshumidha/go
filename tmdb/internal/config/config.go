package config

import (
	"fmt"
	"slices"
	"strings"
)

var ValidTypes = []string{"popular", "top", "upcoming", "playing"}

type Config struct {
	Type string
}

func (c *Config) Validate(folderType string) error {
	if !slices.Contains(ValidTypes, folderType) {
		return fmt.Errorf(
			"invalid type: '%s'\nValid types are: %s",
			folderType,
			strings.Join(ValidTypes, ", "),
		)
	}
	return nil
}

func (c *Config) SetFolderType(folderType string) error {
	if err := c.Validate(folderType); err != nil {
		return err
	}
	c.Type = folderType
	return nil
}
