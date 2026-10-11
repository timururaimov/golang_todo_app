package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfiq() (Config, error) {
	var confiq Config

	if err := envconfig.Process("LOGGER", &confiq); err != nil {
		return Config{}, fmt.Errorf("process envconfig:%w", err)
	}
	return confiq, nil
}

func NewConfigMust() Config {
	config, err := NewConfiq()
	if err != nil {
		panic(fmt.Errorf("get Logger config:%w", err))
	}
	return config
}
