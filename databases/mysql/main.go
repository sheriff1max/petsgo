package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func openMySQL() (*sql.DB, error) {
    src := "user:password@tcp(localhost:3307)/testdb?parseTime=true"

    db, err := sql.Open("mysql", src)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

func initOrders(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS orders")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255),
			price INT
		)`)
	return err
}

func addOrderTx(db *sql.DB, orders []Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO orders (title, price) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, o := range orders {
		if o.price < 0 {
			return fmt.Errorf("invalid price for order %s: %d", o.title, o.price)
		}
		_, err = stmt.Exec(o.title, o.price)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func countOrders(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	return count, err
}

type Order struct {
	id int
	title string
	price int
}

func main() {
	db, err := openMySQL()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e0 := initOrders(db)
	e1 := addOrderTx(db, []Order{{0, "book", 10}, {0, "pen", 5}})
	n1, _ := countOrders(db)
	e2 := addOrderTx(db, []Order{{0, "ok", 1}, {0, "bad", -5}})
	n2, _ := countOrders(db)

	fmt.Println(e0 == nil && e1 == nil && n1 == 2 && e2 != nil && n2 == 2)
}
