// Package app is a nice package
package app

import (
	"context"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/odlev/animal-delivery/logging/internal/config"
	"github.com/odlev/animal-delivery/logging/internal/consumer"
	"github.com/odlev/animal-delivery/logging/internal/processor"
	"github.com/odlev/animal-delivery/logging/internal/storage/clickhouse"
	"github.com/rs/zerolog"
)

type App struct {
	log *zerolog.Logger
	wg *sync.WaitGroup
	storage processor.LogWriter
	cfg *config.Config
	ctx context.Context
}

// const (
// 	logBatchSize = 1000
// 	numConsumers = 3
// )

func New(ctx context.Context, cfg *config.Config, log *zerolog.Logger) (*App, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	storage, err := clickhouse.Init(os.Getenv("CLICKHOUSE_DSN"), os.Getenv("MIGRATIONS_DIR"))
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup

	return &App{ctx: ctx, storage: storage, log: log, wg: &wg, cfg: cfg}, nil
}

func (a *App) Start() {
	for i := range a.cfg.Kafka.NumConsumers {
		idx := i
		a.wg.Go(func() {
			defer a.wg.Done()
			p := processor.NewProcessor(a.storage, a.cfg.Logging.LogBatchSize)
			c, err := consumer.NewConsumer(p, a.cfg.Kafka.BrokerAdresses, a.cfg.Kafka.Topic, a.cfg.Kafka.ConsumerGroup, a.log)
			if err != nil {
				a.log.Err(err).Msg("failed to start consumer")
				return
			}
			a.log.Info().Int("n", idx).Msg("Consumer started")
			go func() {
				<-a.ctx.Done()
				if err := c.Stop(); err != nil {
					a.log.Err(err).Msg("failed to stop consumer")
				}
				if err := p.FLush(a.ctx); err != nil {
					a.log.Err(err).Msg("failed to flush processor")
				}
			}()
			c.Start(a.ctx)
		})
	}
}

func (a *App) Wait() {
	a.wg.Wait()
}
