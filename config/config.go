package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/ortin779/private_theatre_api/db"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Postgres db.PostgresConfig
}

func LoadConfigFromEnv() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("load env config: %w", err)
	}
	return &Config{
		Server: struct {
			Host string
			Port string
		}{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Postgres: db.PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_DBNAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}, nil
}
