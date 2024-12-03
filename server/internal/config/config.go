package config

import "github.com/joho/godotenv"

type DBConfig interface {
	DSN() string
}

func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
