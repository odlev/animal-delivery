// Package consumer is a nice package
package consumer

import (
	"context"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
)

type Consumer struct {
	consumer       *kafka.Consumer
	log            *zerolog.Logger
	processor      Processor
}

type Processor interface {
	// HandleMessage(message []byte, topic kafka.TopicPartition, cn int) error
	HandleLogEvent(ctx context.Context, event []byte, topic kafka.TopicPartition) error
}

const (
	sessionTimeOutMs = 7000
	noTimeout        = -1
)

func NewConsumer(processor Processor, address []string, topic, consumerGroup string, log *zerolog.Logger) (*Consumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeOutMs,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subcribe on topic %s: %w", topic, err)
	}

	return &Consumer{consumer: c, processor: processor, log: log}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for { // event loop))
		select {
		case <-ctx.Done():
			c.log.Info().Msg("context cancelled, stopping consumer loop")
			return
		default:
			kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
			if err != nil {
				c.log.Error().Err(err).Msg("error reading message")
			}
			if kafkaMsg == nil {
				continue
			}
			if err := c.processor.HandleLogEvent(ctx, kafkaMsg.Value, kafkaMsg.TopicPartition); err != nil {
				c.log.Err(err).Msg("error handling message")
				continue
			}
			if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
				c.log.Err(err).Msg("error storing message (store offset partition)")
				continue
			}
		}
	}
}

func (c *Consumer) Stop() error {
	if _, err := c.consumer.Commit(); err != nil {
		return fmt.Errorf("error commiting: %w", err)
	}
	c.log.Info().Msg("commited offset")
	return c.consumer.Close()
}
