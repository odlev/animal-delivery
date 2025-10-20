package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/odlev/animal-delivery/lib/logger"
	"github.com/odlev/animal-delivery/orders/internal/app"
	"github.com/odlev/animal-delivery/orders/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))
	log := logger.SetupLogrus()

	app, err := app.New(ctx, cfg, os.Getenv("POSTGRES_DSN"), log)
	if err != nil {
		panic(err)
	}
	errCh := make(chan error, 1)
	go func() {
		if err := app.Run(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		app.Shutdown(ctx)
	case err := <-errCh:
		log.Fatal("error:", err)
	}

}
