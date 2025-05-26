package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackysum/go-template/cmd/server/config"
	"github.com/jackysum/go-template/src/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	log := getLogger()
	cfg := config.New(log)

	s := server.Setup(
		server.WithLogger(log),
	)

	svr := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: s,
	}

	log.Info().Str("port", cfg.Port).Msg("main#main - server starting")

	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("main#main - server failed")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := svr.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("main#main - server shutdown error")
	}

	log.Info().Msg("main#main - server shutdown complete")
}

func getLogger() zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	out := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	return zerolog.New(out).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}
