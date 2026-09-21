package main

import (
	"fmt"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)


func openPG() (*sql.DB, error) {
    src := "postgres://user:password@localhost:5432/testdb?sslmode=disable"

    db, err := sql.Open("pgx", src)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

func initProducts(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS products;
	CREATE TABLE products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		price INT
	)`)
	return err
}

func addProduct(db *sql.DB, name string, price int) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id", name, price).Scan(&id)
	return id, err
}

func productByName(db *sql.DB, name string) (int, int, error) {
	var id, price int
	err := db.QueryRow("SELECT id, price FROM products WHERE name = $1", name).Scan(&id, &price)
	return id, price, err
}

func avgPrice(db *sql.DB) (float64, error) {
	var avg float64
	err := db.QueryRow("SELECT AVG(price)::float8 FROM products").Scan(&avg)
	return avg, err
}

func main() {
	db, err := openPG()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e0 := initProducts(db)
	id1, e1 := addProduct(db, "chair", 100)
	id2, e2 := addProduct(db, "table", 200)
	_, p, e3 := productByName(db, "chair")
	avg, e4 := avgPrice(db)

	fmt.Println(e0 == nil && e1 == nil && e2 == nil && id1 == 1 && id2 == 2 &&
		e3 == nil && p == 100 && e4 == nil && avg == 150)
}
