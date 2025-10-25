// Package app is a nice package
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/odlev/animal-delivery/api-gateway/internal/clients"
	"github.com/odlev/animal-delivery/api-gateway/internal/config"
	"github.com/odlev/animal-delivery/api-gateway/internal/infrastructure"
	"github.com/odlev/animal-delivery/api-gateway/internal/repo/kafka"
	"github.com/odlev/animal-delivery/api-gateway/internal/transport"
	"github.com/odlev/animal-delivery/api-gateway/internal/usecase/orders"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	httpserver *http.Server
	log        *zerolog.Logger
}

const (
	kafkaOrdersCreatedTopic string = "orders_created"
)
var (
	address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}
)

func New(ctx context.Context, cfg *config.Config, log *zerolog.Logger) (*App, error) {
	shutdownTracer, err := infrastructure.InitTracer("http://jaeger:14268/api/traces") 
	if err != nil {
		return nil, err
	}
	defer shutdownTracer(ctx)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Clients.Orders.GRPCHost, cfg.Clients.Orders.GRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	log.Info().Str("target", conn.CanonicalTarget()).Msg("grpc client is ready")

	k, error:= kafka.NewProducer(address)
	if error != nil {
		return nil, err
	}
	ordersClient := clients.NewOrdersClient(conn)
	usecase := orders.NewOrdersUsecase(ordersClient, k, kafkaOrdersCreatedTopic)
	handler := transport.NewHandler(usecase, log)

	mux := transport.InitRoutes(handler)


	httpserver := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler: otelhttp.NewHandler(transport.Auth(mux), "http-server"),
	}

	return &App{httpserver: httpserver, log: log}, nil
}

func (a *App) Run() error {
	a.log.Info().Str("httpserver address", a.httpserver.Addr).Msg("listen")
	return a.httpserver.ListenAndServe()
}

func (a *App) Shoutdown(ctx context.Context) error {
	if err := a.httpserver.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
