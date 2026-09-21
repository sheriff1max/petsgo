# MongoDB: документный CRUD

## Описание

Реализовать:
- connectMongo() (*mongo.Client, error) + Ping.

На коллекции shop.users:
- insertUser(ctx, c, name, age)
- userAgeMongo(ctx, c, name) (int, error) (FindOne + bson-теги)
- countMongo(ctx, c) (int64, error)
- setAgeMongo(ctx, c, name, age) error через $set.

Перед тестом сделать c.Drop(ctx)

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
docker exec -it <имя_контейнера> mysql -u user -ppassword testdb -e "SELECT * FROM orders;"
```

4. Останровка:

```bash
docker compose down
```
