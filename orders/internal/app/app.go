// Package app is a nice package
package app

import (
	"context"
	"fmt"
	"net"

	"github.com/odlev/animal-delivery/orders/internal/config"
	"github.com/odlev/animal-delivery/orders/internal/repo/postgres"
	"github.com/odlev/animal-delivery/orders/internal/transport"
	"github.com/odlev/animal-delivery/orders/internal/usecase"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	log     *logrus.Logger
	server  *grpc.Server
	storage usecase.OrdersManager
	cfg     *config.Config
}

func New(ctx context.Context, cfg *config.Config, dsn string, log *logrus.Logger) (*App, error) {
	storage, err := postgres.Init(ctx, dsn)
	if err != nil {
		return nil, err
	}

	ordersUS := usecase.New(storage)
	srv := grpc.NewServer()
	transport.Register(srv, ordersUS, log)

	return &App{
		log:     log,
		server:  srv,
		storage: storage,
		cfg:     cfg,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		return err
	}
	defer lis.Close()

	a.log.Infof("orders gRPC listening on %s", lis.Addr())
	return a.server.Serve(lis)
}

func (a *App) Shutdown(ctx context.Context) {
	a.log.Info("shutting down gRPC orders server")
	a.server.GracefulStop()

	a.log.Info("shutting down storage")
	a.storage.Close()
}
