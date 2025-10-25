// Package processor is a package who represents transport layer + service layer on one layer!
package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/odlev/animal-delivery/contracts/gen/go/loggingpb"
	"github.com/odlev/animal-delivery/logging/internal/storage"
	"google.golang.org/protobuf/proto"
)

type LogWriter interface {
	InsertLogs(ctx context.Context, logs []storage.LogRecord) error
}

type Processor struct {
	logWriter LogWriter
	batch     []storage.LogRecord
	batchSize int
}

func NewProcessor(logWriter LogWriter, batchSize int) *Processor {
	return &Processor{
		logWriter: logWriter,
		batch:     make([]storage.LogRecord, batchSize),
		batchSize: batchSize,
	}
}

func (p *Processor) HandleLogEvent(ctx context.Context, msg []byte, topic kafka.TopicPartition) error {
	var event loggingpb.LogEvent
	if err := proto.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("failed to unmarshal log event: %w", err)
	}

	fieldsJSON, err := json.Marshal(event.Fields)
	if err != nil {
		return fmt.Errorf("failed to marshal event fields: %w", err)
	}

	log := storage.LogRecord{
		KafkaTopic: *topic.Topic,
		Timestamp:  time.Unix(event.Timestamp, 0),
		Level:      event.Level,
		Message:    event.Message,
		Service:    event.Service,
		TraceID:    event.TraceId,
		SpanID:     event.SpanId,
		Fields:     string(fieldsJSON),
	}

	p.batch = append(p.batch, log)
	if len(p.batch) >= p.batchSize {
		if err := p.logWriter.InsertLogs(ctx, p.batch); err != nil {
			return fmt.Errorf("failed to insert logs: %w", err)
		}
		p.batch = p.batch[:0]
	}

	return nil
}

func (p *Processor) FLush(ctx context.Context) error {
	if len(p.batch) > 0 {
		return p.logWriter.InsertLogs(ctx, p.batch)
	}
	return nil
}
