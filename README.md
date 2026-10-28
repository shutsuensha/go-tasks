# Go Tasks

Go Tasks — production-oriented сервис для управления задачами. Проект показывает, как построить HTTP API на Go с PostgreSQL, Redis-кэшем, асинхронными заданиями, миграциями, интеграционными тестами и observability-инструментами.

## Возможности

- создание задачи;
- получение задачи по идентификатору;
- пагинированный список задач;
- запись события `task_created` в отдельную таблицу;
- асинхронная постановка фоновой задачи после создания задачи;
- Redis-кэширование списка задач;
- защита от cache stampede через `singleflight`;
- распределённый rate limit через Redis — 10 запросов в секунду на IP;
- health- и readiness-checks;
- Prometheus-метрики HTTP, PostgreSQL и Redis;
- distributed tracing через OpenTelemetry и Jaeger;
- Swagger-документация API;
- graceful shutdown HTTP-сервера;
- интеграционные тесты PostgreSQL через Testcontainers.

## Архитектура

В проекте два исполняемых приложения:

- `cmd/app` — HTTP API;
- `cmd/worker` — обработчик фоновых задач и cron-задач.

Основной запрос проходит через middleware Chi, handler, service layer и слой доступа к данным. SQL-запросы находятся в `internal/db/query.sql`, а Go-код для них генерируется SQLC. PostgreSQL используется как основное хранилище, Redis — для кэша, rate limiting и очереди Asynq.

Создание задачи выполняется в транзакции: одновременно создаются запись в `tasks` и событие в `task_events`. После успешного commit задача `task:created` отправляется в Redis-backed очередь. Благодаря этому бизнес-логика отделена от HTTP-транспорта и фоновой обработки.

## API

| Метод | Endpoint | Назначение |
| --- | --- | --- |
| `POST` | `/tasks` | Создать задачу |
| `GET` | `/tasks?limit=10&offset=0` | Получить список задач |
| `GET` | `/tasks/{id}` | Получить задачу по ID |
| `GET` | `/health` | Liveness check |
| `GET` | `/ready` | Readiness check PostgreSQL |
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/swagger/index.html` | Swagger UI |
| `GET` | `/debug/pprof/` | Go pprof |

Пример создания задачи:

```bash
curl -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Prepare README","description":"Document the project"}'
```

## Технологии

### Backend

- Go 1.25;
- стандартный `net/http`;
- Chi — HTTP router и middleware;
- PostgreSQL 16;
- `pgx/v5` и `pgxpool` — PostgreSQL driver и connection pool;
- SQLC — типобезопасная генерация Go-кода из SQL;
- golang-migrate — versioned database migrations;
- Viper — конфигурация из environment variables и `.env`;
- go-playground/validator — валидация входных DTO;
- Zap — structured logging;
- UUID — request correlation ID.

### Cache, queues and reliability

- Redis 7;
- `go-redis` — Redis client;
- Asynq — фоновые задачи, retries, очереди с приоритетами и cron scheduler;
- Redis rate limiter — распределённое ограничение запросов;
- `singleflight` — объединение параллельных cache-miss запросов;
- транзакции PostgreSQL и connection pooling;
- graceful shutdown и таймауты HTTP-сервера;
- recovery middleware для обработки panic.

### Observability and performance

- Prometheus — сбор метрик;
- Grafana — визуализация метрик;
- OpenTelemetry — трассировка HTTP, PostgreSQL и Redis;
- Jaeger — просмотр distributed traces;
- Go `pprof` — профилирование CPU, памяти, goroutines и блокировок;
- `hey` — нагрузочное тестирование (предусмотрено подходом проекта).

### Testing and tooling

- Testcontainers Go — PostgreSQL в контейнере для интеграционных тестов;
- `httptest` — HTTP-тестирование;
- Testify — assertions;
- Docker и Docker Compose;
- Makefile;
- golangci-lint;
- Swagger/OpenAPI через Swaggo.

## Структура проекта

```text
cmd/app              HTTP API
cmd/worker           background worker and scheduler
internal/handler     HTTP handlers and DTOs
internal/service     business logic
internal/db          SQLC-generated database layer
internal/cache       Redis client
internal/queue       Asynq client, worker and scheduler
internal/middleware  request ID, logging, metrics, recovery, rate limit
internal/metrics     application, database and Redis metrics
internal/observability OpenTelemetry setup
migrations           PostgreSQL migrations
tests/integration    integration tests with Testcontainers
docs                 generated Swagger documentation
```

## Быстрый запуск через Docker Compose

Требования: Docker и Docker Compose.

```bash
docker compose up -d --build
```

Сервисы будут доступны по адресам:

- API: `http://localhost:8080`;
- Swagger: `http://localhost:8080/swagger/index.html`;
- Prometheus: `http://localhost:9090`;
- Grafana: `http://localhost:3000`;
- Jaeger: `http://localhost:16686`.

Миграции запускаются отдельно:

```bash
make migrate-up
```

Для локального запуска worker после старта Redis:

```bash
go run ./cmd/worker
```

Основные переменные приложения:

| Переменная | Назначение | Пример |
| --- | --- | --- |
| `HTTP_PORT` | HTTP-порт API | `8080` |
| `DB_URL` | PostgreSQL DSN | `postgres://postgres:postgres@localhost:3333/tasks_db?sslmode=disable` |
| `REDIS_ADDR` | Адрес Redis | `localhost:7744` |

## Команды Makefile

```bash
make build         # собрать API
make run           # запустить API
make tests         # запустить тесты
make lint          # запустить golangci-lint
make migrate-up    # применить миграции
make migrate-down  # откатить последнюю миграцию
make docker-up     # запустить Docker Compose
make docker-down   # остановить Docker Compose
```

## Текущее состояние

Проект является демонстрационным примером архитектуры backend-сервиса. В текущей реализации API не содержит update/delete операций. Обработчик `task:created` намеренно возвращает ошибку `test retry`, чтобы демонстрировать механизм повторных попыток Asynq; cron-задача пока только пишет сообщение в лог. Worker не оформлен отдельным сервисом в `docker-compose.yml` и запускается отдельно.