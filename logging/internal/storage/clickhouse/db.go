// Package clickhouse is a nice package
package clickhouse

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/odlev/animal-delivery/logging/internal/storage"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	conn driver.Conn
}

func Init(dsn, migrationsDir string) (*Storage, error) {
	applyMigrations(dsn, migrationsDir)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return nil, err
	}
	// log.Log().Any("contributors", conn.Contributors()).Send()
	return &Storage{conn: conn}, nil
}

func applyMigrations(dsn, migrationsDir string) error {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return fmt.Errorf("failed to open connections: %w", err)
	}
	defer db.Close()
	if err := goose.Up(db, migrationsDir); err != nil {
		if err != goose.ErrNoNextVersion {
			return err
		}
	}
	return nil
}

func (s *Storage) InsertLogs(ctx context.Context, logs []storage.LogRecord) error {
	batch, err := s.conn.PrepareBatch(ctx,
		`INSERT INTO logs (timestamp, level, message, service, trace_id, span_id, fields)
	VALUES (?, ?, ?, ?, ?, ?, ?`)
	if err != nil {
		return fmt.Errorf("faled to prepare batch for insert logs")
	}

	for _, log := range logs {
		err = batch.Append(
			log.KafkaTopic,
			log.Timestamp,
			log.Level,
			log.Message,
			log.Service,
			log.TraceID,
			log.SpanID,
			log.Fields,
		)
		if err != nil {
			return fmt.Errorf("failed to append logs for batch insert")
		}
	}
	return batch.Send()
}
