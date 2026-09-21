# MiniSocial API

Минималистичная социальная сеть, реализованная в виде REST API на Go.

## Возможности

- Регистрация пользователей
- Создание постов
- Подписка на других пользователей
- Лайки на посты
- Получение ленты новостей

## Стек технологий

- **Go** (chi, pgx)
- **PostgreSQL**
- **Docker, Docker Compose**

## Структура проекта

- `cmd/api` - точка входа и настройка сервера
- `internal/handler` - HTTP обработчики и роутинг
- `internal/service` - бизнес-логика приложения
- `internal/repository` - работа с базой данных
- `internal/domain` - сущности предметной области

## Запуск проекта

### Требования
- Docker & Docker Compose
- Make (опционально)

### 1. Запуск через Docker Compose

```bash
make docker-up
# или вручную
docker-compose up -d --build
```

API: http://localhost:8080

Swagger: http://localhost:8080/swagger/index.html

### 2. Локальный запуск

При локальной установке Go и Postgres:

```bash
make migrate-up
make run
```

## API Endpoints

Для защищенных эндпоинтов нужно передавать заголовок X-User-ID с UUID пользователя.

### Пользователи

- `POST /api/v1/users` - зарегистрировать пользователя

```json
{ "username": "john_doe" }
```

### Посты и Лента

- `POST /api/v1/posts` - создать пост (требует X-User-ID)

```json
{ "content": "Hello world!" }
```

- `GET /api/v1/feed` - получить ленту новостей (требует X-User-ID)

### Социальные взаимодействия

- `POST /api/v1/users/{id}/follow` - подписаться на пользователя (требует X-User-ID)
- `POST /api/v1/posts/{id}/like` - поставить лайк (требует X-User-ID)

## Примеры запросов curl

1. Создать пользователя:

```bash
curl -X POST http://localhost:8080/api/v1/users \
-H "Content-Type: application/json" \
-d '{"username": "alice"}'
```

2. Создать пост:

```bash
curl -X POST http://localhost:8080/api/v1/posts \
-H "Content-Type: application/json" \
-H "X-User-ID: <alice_id>" \
-d '{"content": "My first post!"}'
```

3. Получить ленту:

```bash
curl -X GET http://localhost:8080/api/v1/feed \
-H "X-User-ID: <alice_id>"
```

## Краткая инструкция

1. Запусти `make docker-up`.
2. Тестируй API через Postman или cURL!
3. `docker-compose stop` или `docker-compose down`
