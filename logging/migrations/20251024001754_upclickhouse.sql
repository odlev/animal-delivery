-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logs (
    timestamp Datetime,
    level String,
    message String,
    service String,
    trace_id String,
    span_id String,
    fields String --json строка
) ENGINE = MergeTree()
ORDER BY (timestamp, level);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE logs;
-- +goose StatementEnd
