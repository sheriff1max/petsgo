package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)


func initUsersDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			age INTEGER
		)`,
	)
	return db, err
}

func addUser(db *sql.DB, name string, age int) error {
	_, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", name, age)
	return err
}

func countUsers(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func openMySQL() (*sql.DB, error) {
	dsn := "user:password@tcp(localhost:3307)/testdb?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initMySQLUsers(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		age INT
	)`)
	return err
}

func migrateUsers(src, dst *sql.DB) (int, error) {
	tx, err := dst.Begin()
	if err != nil {
		return 0, err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	rows, err := src.Query("SELECT name, age FROM users")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	stmt, err := tx.Prepare("INSERT INTO users (name, age) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var name string
		var age int
		if err := rows.Scan(&name, &age); err != nil {
			return 0, err
		}
		
		if _, err := stmt.Exec(name, age); err != nil {
			return 0, err
		}
		count++
	}
	
	if err = rows.Err(); err != nil {
		return 0, err
	}

	// 5. Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}


func main() {
	dbFile := "migr.db"
	os.Remove(dbFile)

	sdb, _ := initUsersDB(dbFile)
	defer sdb.Close()
	
	addUser(sdb, "Ann", 30)
	addUser(sdb, "Bob", 25)

	mdb, err := openMySQL()
	if err != nil {
		panic(err)
	}
	defer mdb.Close()

	initMySQLUsers(mdb)

	n, err_migr := migrateUsers(sdb, mdb)
	
	cs, _ := countUsers(sdb)
	cm, _ := countUsers(mdb)

	fmt.Println(n == 2 && err_migr == nil && cs == 2 && cm == 2)
}