package models

import (
	"github.com/Alexandr-Fox/fox-core/internal/database"
)

type Config struct {
	database.Entity
	Name  string `gorm:"uniqueIndex" json:"name"`
	Value string `json:"value"`
}

func GetConfig(name string, value string) (c *Config, err error) {
	c = &Config{}
	err = c.First(name, value)

	return c, err
}
