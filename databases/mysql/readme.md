# MySQL: транзакции и откат

## Описание

Реализовать:
- openMySQL()
- initOrders(db) создаёт orders(id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255), price INT)
- addOrderTx(db, orders []Order) error: внутри одной транзакции вставляет все заказы; если у какого-то price < 0 — возвращает ошибку и откатывает всё (проверить, что «хороший» заказ из неудачной пачки тоже не остался)

## Драйвера:

```bash
go get github.com/go-sql-driver/mysql
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
docker exec -it <имя_контейнера> mysql -u user -ppassword testdb -e "SELECT * FROM orders;"
```

4. Останровка:

```bash
docker compose down
```
