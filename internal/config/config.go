package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPAddress string
	DatabaseURL string
	JWTKey      string
}

func Load() (Config, error) {
	JWTKey := os.Getenv("JWT_KEY")
	if JWTKey == "" {
		return Config{}, errors.New("JWT_KEY environment variable is not set")
	}

	return Config{
		HTTPAddress: os.Getenv("HTTP_ADDRESS"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTKey:      JWTKey,
	}, nil
}
