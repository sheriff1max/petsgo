# Проект: укорачиватель ссылок

Это сервис, который предоставляет API по созданию сокращенных ссылок.

Ссылка должна:
- быть уникальной и на один оригинальный URL должна ссылаться только одна сокращенная ссылка;
- быть длинной 10 символов;
- состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание).

Сервис написан на Go и принимает следующие http-запросы:
1. метод `POST`, который сохраняет оригинальный URL в базе и возвращает сокращённый;
2. метод `GET`, который принимает сокращённый URL и возвращает оригинальный URL.

Что сделано:
- сервис распространён в виде Docker-образа;
- в качестве хранилища реализовано две реализации (какое хранилище использовать, указывается параметром при запуске сервиса):
    - PostgreSQL;
    - хранение ссылок в памяти;
- реализованный функционал покрыт Unit-тестами.

## Структура проекта

- cmd/short-urls/main.go - точка входа

- internal/config/config.go - env/flags

- internal/handler - HTTP и тесты

- internal/service - логика генерации

- internal/storage/storage.go - тут интерфейс хранилища
- internal/storage/memory.go - реализация хранения в памяти
- interna/storage/postgres.go - реализация хранения в postgresql

## Инструкция к API

В текущей реализации один GET и один POST запросы.

Пример `POST`-запроса:
- http://localhost:8080/generate
- в `Body` кладётся текст с ссылкой, например, **https://github.com**
- Возвращает текст с сокращённой ссылкой, например, **http://localhost:8080/dqyN6V3ZTs**

Пример `GET`-запроса:
- http://localhost:8080/dqyN6V3ZTs - ссылка, полученная из ответа `POST`-запроса
- Возвращает текст с сокращённой ссылкой, например, **https://github.com**

## Инструкция запуска

Сервер работает тут: http://localhost:8080

### Запуск докера

Выбрать можно 2 типа хранилища: `memory` и `postgres`.

Во время запуска автоматически запускаются тесты, поднимается PostgreSQL, запускается сервер.

#### Хранилище - память

Для Linux / macOS:

```bash
STORAGE_TYPE=memory docker-compose up --build
```

Для Windows (PowerShell):

```bash
$env:STORAGE_TYPE="memory"; docker-compose up --build
```

#### Хранилище - PostgreSQL

Для Linux / macOS:

```bash
STORAGE_TYPE=postgres docker-compose up --build
```

Для Windows (PowerShell):

```bash
$env:STORAGE_TYPE="postgres"; docker-compose up --build
```

### Запуск локально

```bash
go mod tidy
go test ./...
```

#### Хранилище - память

```bash
go run ./cmd/short-urls -s memory
```

#### Хранилище - PostgreSQL

```bash
go run ./cmd/short-urls -s postgres
```