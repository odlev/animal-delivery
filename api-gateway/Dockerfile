FROM golang:1.25 AS builder
WORKDIR /app

COPY go.work go.work.sum ./
COPY orders orders
COPY delivery delivery
COPY contracts contracts
COPY lib lib
COPY api-gateway api-gateway

WORKDIR /app/api-gateway
RUN go mod download
RUN mkdir -p /app/bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api-gateway ./cmd

FROM alpine:3.20 AS runner
WORKDIR /srv/orders
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/bin/api-gateway ./api-gateway
COPY --from=builder /app/api-gateway/configs ./configs
COPY api-gateway/.env .env
EXPOSE 8080
ENTRYPOINT [ "./api-gateway" ]
