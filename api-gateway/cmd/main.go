package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/odlev/animal-delivery/api-gateway/internal/app"
	"github.com/odlev/animal-delivery/api-gateway/internal/config"
	"github.com/odlev/animal-delivery/lib/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))
	log := logger.SetupZerolog()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:2112", nil)

	app, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("error initialization app")
	}
	go func() {
		log.Info().Msg("http server running")
		if err := app.Run(); err != nil {
			log.Fatal().Err(err).Msg("http server error")
		}
	}()
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shoutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("http server shutdown error")
	}
}
