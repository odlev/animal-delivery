// Package storage is a nice package
package storage

import "time"

type LogRecord struct {
	KafkaTopic string
	Timestamp time.Time
	Level string
	Message string
	Service string
	TraceID string
	SpanID string
	Fields string
}
