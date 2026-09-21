# SQLite: CRUD через database/sql (реляционка без Docker)

## Описание

initUsersDB(path) (*sql.DB, error): открыть файл, создать users(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, age INT).

Реализовать:
- addUser(db, name, age) (int64, error) (LastInsertId)
- userAge(db, name) (int, error) (QueryRow)
- countUsers(db) (int, error)
- renameUser(db, old, new) (int64, error) (RowsAffected)

## Драйвера на выбор

- go get github.com/mattn/go-sqlite3
- go get modernc.org/sqlite
