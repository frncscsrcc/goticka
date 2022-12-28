package migrations

import (
	"database/sql"
	"log"
)

var migrations []func(*sql.DB) error

func init() {
	migrations = make([]func(*sql.DB) error, 0)

	// -----------------------------------------------------------
	// Reference your migrations here
	// -----------------------------------------------------------
	migrations = append(migrations, migrate_v0)
	migrations = append(migrations, migrate_v1)
	migrations = append(migrations, migrate_v2)
	// -----------------------------------------------------------
}

func Migrate(db *sql.DB) {
	lastMigration := getLastMigration(db)
	originalLastMigration := lastMigration

	log.Printf("Starting migrations. Current version = %d", lastMigration)

	for version, migration := range migrations {
		if version > lastMigration {
			err := migration(db)
			if err != nil {
				panic(err)
			}
			lastMigration = version
		}
	}

	if originalLastMigration == lastMigration {
		log.Printf("Nothing to migrate\n")
	} else {
		log.Printf("Current version %d\n", lastMigration)
	}

}

func getLastMigration(db *sql.DB) int {
	rows, err := db.Query("SELECT MAX(Version) FROM _migrations LIMIT 1")
	if err != nil {
		return -1
	}

	defer rows.Close()

	for rows.Next() {
		var version int
		if err = rows.Scan(&version); err != nil {
			return -1
		}
		return version
	}
	return -1
}

func exec_statement(db *sql.DB, queries ...string) error {
	for _, query := range queries {
		log.Printf("Executing query: %s\n", query)
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
