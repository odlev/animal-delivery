FROM golang:1.25 AS builder
WORKDIR /app

COPY go.work go.work.sum ./
COPY api-gateway api-gateway
COPY contracts contracts
COPY delivery delivery
COPY lib lib
COPY orders orders
COPY orders/migrations orders/migrations

WORKDIR /app/orders
RUN go mod download
RUN mkdir -p /app/bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/orders ./cmd

FROM alpine:3.20 AS runner
WORKDIR /srv/orders
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bin/orders ./orders
COPY --from=builder /app/orders/configs ./configs
COPY --from=builder /app/orders/migrations ./migrations
COPY orders/.env .env
EXPOSE 50051
ENTRYPOINT [ "./orders" ]
