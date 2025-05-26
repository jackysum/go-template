package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jackysum/go-template/cmd/server/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	log := getLogger()
	cfg := config.New(log)

	fmt.Println(cfg.Port)
}

func getLogger() zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	out := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	return zerolog.New(out).Level(zerolog.InfoLevel).With().Timestamp().Logger()
}
