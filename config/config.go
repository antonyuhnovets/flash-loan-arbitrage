package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Log
	HttpServer
	Storage
	Blockchain
}

type Log struct {
	Level string `env:"LOG_LEVEL" env-default:"debug"`
}
type HttpServer struct {
	Port string `env:"HTTP_PORT" env-default:"8080"`
	Host string `env:"HTTP_HOST" env-default:"0.0.0.0"`
}

type Storage struct {
	Type string `env:"STORAGE_TYPE" env-default:"localfile"`
	Localstorage
	Database
}

type Localstorage struct {
	Path string `env:"STORAGE_PATH" env-default:"./storage_test/test.json"`
}

type Database struct {
	Driver   string `env:"DATABASE_DRIVER"`
	Host     string `env:"DATABASE_HOST"`
	Port     string `env:"DATABASE_PORT"`
	Username string `env:"DATABASE_USERNAME"`
	Password string `env:"DATABASE_PASSWORD"`
	Name     string `env:"DATABASE_NAME"`
}

type Blockchain struct {
	Name string `env:"BLOCKCHAIN_NAME" env-default:"goerli"`
	Url  string `env:"BLOCKCHAIN_RPC_URL"`
	Account
	Contract
}

type Account struct {
	Address string `env:"ACCOUNT_ADDRESS"`
	// pk      string `env:"ACCOUNT_PRIVATE_KEY"`
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
