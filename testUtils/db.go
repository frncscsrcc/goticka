package testUtils

import (
	"database/sql"
	"goticka/pkg/db/migrations"

	_ "github.com/mattn/go-sqlite3"
)

func NewTestDB() *sql.DB {
	inMemoryConnection, err := sql.Open("sqlite3", ":memory:?parseTime=true")
	if err != nil {
		panic(err)
	}

	migrations.Migrate(inMemoryConnection)

	return inMemoryConnection
}
