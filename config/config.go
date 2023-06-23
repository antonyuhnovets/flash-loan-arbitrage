package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Blockchain
	Account
	Contract
	Storage
}

type Storage struct {
	Type string `env:"STORAGE_TYPE" env-default:"localfile"`
	Path string `env:"STORAGE_PATH" env-default:"./storage_test/test.json"`
}
type Blockchain struct {
	Name string `env:"NETWORK_CHAIN_NAME" env-default:"goerli"`
	Url  string `env:"NETWORK_CHAIN_URL"`
}

type Account struct {
	Address string `env:"ACCOUNT_ADDRESS"`
	pk      string `env:"ACCOUNT_PRIVATE_KEY"`
}

type Contract struct {
	Name    string `env:"CONTRACT_NAME" env-default:"FlashLoanArbitrage"`
	Address string `env:"CONTRACT_ADDRESS"`
	Input   string `env:"CONTRACT_INPUT"`
}

func LoadConfig() (*Config, error) {
	var conf Config

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
