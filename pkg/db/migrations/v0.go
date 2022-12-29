package migrations

import (
	"database/sql"
)

func migrate_v0(db *sql.DB) error {
	err := exec_statement(db,
		`CREATE TABLE _migrations (
			Version int PRIMARY KEY
		);`,

		`INSERT INTO _migrations VALUES (0)`,
	)
	return err
}
