# PostgreSQL

## Описание

Реализовать:
- openPG(), initProducts(db) (с DROP TABLE IF EXISTS, чтобы id начинались с 1)
- addProduct(db, name, price) (int, error) через INSERT ... RETURNING id
- productByName(db, name) (int, int, error) (id и цена)
- avgPrice(db) (float64, error)

## Драйверы

```bash
go get github.com/jackc/pgx/v5
```

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
docker exec -it <имя_контейнера> psql -U user -d testdb -c "SELECT * FROM products;"
```

4. Останровка:

```bash
docker compose down
```
