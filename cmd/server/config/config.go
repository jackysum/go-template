package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Port string
}

func New(log zerolog.Logger) Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("config#New - error loading .env")
	}

	return Config{
		Port: os.Getenv("PORT"),
	}
}
