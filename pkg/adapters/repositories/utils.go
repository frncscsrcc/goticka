package repositories

import (
	"database/sql"
)

func exec_statement(db *sql.DB, queries ...string) error {
	for _, query := range queries {
		stmn, err := db.Prepare(query)
		if err != nil {
			return err
		}
		defer stmn.Close()
		if _, err := stmn.Exec(); err != nil {
			return err
		}
	}
	return nil
}
