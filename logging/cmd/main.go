package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/odlev/animal-delivery/lib/logger"
	"github.com/odlev/animal-delivery/logging/internal/app"
	"github.com/odlev/animal-delivery/logging/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log := logger.SetupZerolog()
	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))

	app, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Err(err).Send()
	}
	log.Info().Msg("app starting")
	go app.Start()
	<-ctx.Done()

	log.Info().Msg("app waiting for consumers and processors")
	app.Wait()
}
