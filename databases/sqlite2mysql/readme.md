# SQLite в MySQL

## Описание

migrateUsers(src, dst *sql.DB) (int, error): читает всех пользователей из SQLite (Query + цикл Scan), вставляет их в MySQL-таблицу users одной транзакцией и возвращает число перенесённых.

## Запуск

1.

```bash
docker compose up -d
```

2.

```bash
go run main.go
```

3. Поверка добавленной записи, или выполнив:

```bash
docker exec -it <имя_контейнера> mysql -u user -ppassword testdb -e "SELECT * FROM users;"
```

4. Останровка:

```bash
docker compose down
```
