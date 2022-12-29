package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	connection, err := ConnectSQLite3()
	if err != nil {
		panic(err)
	}
	db = connection
}

func ConnectSQLite3() (*sql.DB, error) {
	return sql.Open("sqlite3", "./test.db?parseTime=true")
}

func GetDB() *sql.DB {
	return db
}
