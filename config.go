package main

import (
	"github.com/BurntSushi/toml"
)

type TelegramConfig struct {
	Token string
}
type Config struct {
	Telegram TelegramConfig
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}
