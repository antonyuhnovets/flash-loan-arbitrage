package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Blockchain
}

type Blockchain struct {
	Network  string `env:"NETWORK_NAME" env-default:"goerli"`
	Url      string `env:"RPC_URL"`
	WalletPK string `env:"PRIVAT_KEY"`
}

func LoadConfig() (*Config, error) {
	var conf Config

	err := cleanenv.ReadEnv(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
