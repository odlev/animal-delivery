что внедрить
oidc - https://github.com/coreos/go-oidc
rbac - https://github.com/casbin/casbin
кафка
редис
кликхаус
метрики прометеус и графана и трассировка жаегер (опентелеметри)
кейклок подрубить (рассмотреть это)


Project includes:
    misroservices: 
        - api-gateway
        - orders
        - logging 
- REST
- GRPC, Protobuf (separate go module for contracts)
- Kafka cluster
- Prometheus, Grafana
- Jaeger, OpenTelemetry
- Docker, Docker-compose
- go work))

# Структура проекта:

animal-delivery/
├─ contracts/
│  ├─ go.mod (module github.com/odlev/animal-delivery/contracts)
│  ├─ go.sum
│  ├─ buf.gen.yaml
│  ├─ buf.yaml
│  ├─ Taskfile.yaml
│  ├─ proto/
│  │  ├─ delivery.proto
│  │  └─ orders.proto
│  ├─ gen/
│  │  └─ go/
│  │      └─ anidelive/
│  │          ├─ delivery.pb.go
│  │          ├─ delivery_grpc.pb.go
│  │          ├─ orders.pb.go
│  │          └─ orders_grpc.pb.go
│  └─ .task/…                     # служебные файлы Task
│
├─ api-gateway/
│  ├─ go.mod (module github.com/odlev/animal-delivery/api-gateway)
│  ├─ cmd/gateway/main.go          # точка входа, wiring зависимостей
│  ├─ internal/http/
│  │  ├─ router.go                 # построение chi/gin маршрутов
│  │  ├─ handlers/                 # REST-эндпоинты
│  │  └─ middleware/               # логирование, auth, метрики
│  ├─ internal/service/            # бизнес-оркестрация, вызовы gRPC
│  ├─ pkg/clients/processor/       # обёртка над gRPC-клиентом Processor’а
│  ├─ configs/                     # yaml/json/env шаблоны
│  ├─ deploy/                      # Dockerfile, Helm, k8s
│  └─ test/                        # e2e/contract тесты
│
├─ orders/
│  ├─ go.mod (module github.com/odlev/animal-delivery/processor)
│  ├─ cmd/processor/main.go        # запуск gRPC сервера, DI
│  ├─ internal/server/grpc.go      # регистрация сервисов, interceptors
│  ├─ internal/service/            # доменные use-cases (доставка, заказы)
│  ├─ internal/repo/               # БД, очереди, внешние ресурсы
│  ├─ internal/worker/             # фоновые джобы, обработчики событий
│  ├─ configs/                     # конфигурация сервиса
│  ├─ deploy/                      # Dockerfile, Helm, k8s
│  └─ test/                        # unit/integration тесты
│
─ delivery/
│  ├─ go.mod (module github.com/odlev/animal-delivery/processor)
│  ├─ cmd/processor/main.go        # запуск gRPC сервера, DI
│  ├─ internal/server/grpc.go      # регистрация сервисов, interceptors
│  ├─ internal/service/            # доменные use-cases (доставка, заказы)
│  ├─ internal/repo/               # БД, очереди, внешние ресурсы
│  ├─ internal/worker/             # фоновые джобы, обработчики событий
│  ├─ configs/                     # конфигурация сервиса
│  ├─ deploy/                      # Dockerfile, Helm, k8s
│  └─ test/                        # unit/integration тесты
└─ README.md                         # unit/integration тесты

## Contracts: 
    write ```task newgen``` for delete old and generate new .pb.go files

WHAT WILL BE FUTURE
Подтяни stubs в оба сервисных модулей: в go.mod и go.mod сделай go get github.com/odlev/animal-delivery/contracts/gen/go@latest (можно пока через replace ../contracts).
Разверни каркас processor: создавай cmd/processor/main.go, internal/server/grpc.go, зарегистрируй DeliveryServiceServer и OrderServiceServer (пока с заглушками). Настрой логгирование, чтение конфигов, graceful shutdown.
Продумай доменную логику и storage для processor: что хранится (memory map, Postgres), какие статусы заказов, как обновляются доставки. Реализуй сервисные методы, покрывай unit-тестами.
В api-gateway подними HTTP маршрутизатор (chi/gin/echo), заведи хендлеры POST /orders, POST /deliveries, GET /deliveries/{id} и пр. Внутри — валидация запросов, вызов gRPC клиента processor, трансформация ответов в REST.
Добавь конфиги (env) для адреса gRPC, уровни логов, таймауты. Подумай про middleware: логгирование, метрики, auth, rate-limit (по необходимости).
Настрой наблюдаемость: Logrus формат, Prometheus/otel interceptors на gRPC и HTTP. Добавь health-check endpoints.
В contracts зафиксируй процесс генерации (Taskfile уже есть) и CI-проверку (buf lint, buf breaking, buf generate).
Собери Dockerfile/compose, чтобы локально одновременно поднять gateway+processor и прогнать e2e тесты.
Дальше — нагрузочный клиент, сценарии для интеграционных тестов и документирование REST API (OpenAPI).
