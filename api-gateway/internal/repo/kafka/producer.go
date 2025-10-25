// Package kafka is a nice package
package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/odlev/animal-delivery/contracts/gen/go/loggingpb"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	producer *kafka.Producer
}

const (
	flushTimeout = 5000 // ms
)

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	prod, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new producer: %w", err)
	}
	return &Producer{producer: prod}, nil
}

func (p *Producer) Produce(msg string, topic string) error {

	/* payload, err := proto.Marshal(event)
	if err != nil {
		return err
	} */
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(kafka.PartitionAny),
		},
		Value: []byte(msg),
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("failed to produce msg: %w", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			return fmt.Errorf("Produce: kafka delivery failed: %w", ev.TopicPartition.Error)
		}
		return nil
	case kafka.Error:
		return fmt.Errorf("kafka return error event: %w", ev)
	default:
		return errUnknownType
	}
}

func (p *Producer) ProduceLog(event *loggingpb.LogEvent, topic string) error {

	payload, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(kafka.PartitionAny),
		},
		Value: payload,
		Key:   []byte(event.GetTraceId()),
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("failed to produce msg: %w", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			return fmt.Errorf("Produce: kafka delivery failed: %w", err)
		}
		return nil
	case kafka.Error:
		return fmt.Errorf("kafka return error event: %w", ev)
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout) // 5s: прежде чем закрыть продюсер, хочу дождаться получения наших отправленных, но еще не полученных сообщений
	p.producer.Close()
}
