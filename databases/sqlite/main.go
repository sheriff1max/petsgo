package main

import (
	"fmt"
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)


func initUsersDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE users(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, age INT)")
	if err != nil {
		return nil, err
	}

	return db, err
}

func addUser(db *sql.DB, name string, age int) (int64, error) {

	res, err := db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", name, age)
	if err != nil {
		return 0, err
	}

	last_id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return last_id, nil
}

func userAge(db *sql.DB, name string) (int, error) {
	row := db.QueryRow("SELECT age FROM users WHERE name = $1", name)

	var age int
	err := row.Scan(&age)
	if err != nil {
		return 0, err
	}
	return age, nil
}

func countUsers(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func renameUser(db *sql.DB, old, new string) (int64, error) {
	res, err := db.Exec("UPDATE users SET name = $1 WHERE name = $2", new, old)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func main() {
	err_rem := os.Remove("shop.db")
	if err_rem != nil {
		panic(err_rem)
	}

	db, err := initUsersDB("shop.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	id1, e1 := addUser(db, "Ann", 30)
	id2, e2 := addUser(db, "Bob", 25)
	age, e3 := userAge(db, "Ann")
	cnt, e4 := countUsers(db)
	ra, e5 := renameUser(db, "Bob", "Robert")
	fmt.Println(e1 == nil && e2 == nil && id1 == 1 && id2 == 2 &&
		e3 == nil && age == 30 && e4 == nil && cnt == 2 && e5 == nil && ra == 1)
}