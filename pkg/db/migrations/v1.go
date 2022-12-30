package migrations

import (
	"database/sql"
)

func migrate_v1(db *sql.DB) error {
	err := exec_statement(db,
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			external BOOL,
			username varchar(255) UNIQUE,
			email varchar(255) UNIQUE,
			password varchar(255),
			created datetime,
			changed datetime,
			deleted datetime
		);`,

		`CREATE TABLE queues (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name varchar(255), 
			description text,
			created datetime,
			changed datetime,
			deleted datetime
		);`,

		`CREATE TABLE tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			status INTEGER,
			subject text,
			queueID INTEGER,
			created datetime,
			changed datetime, 
			deleted datetime
		);`,

		`CREATE TABLE articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ticketID INTEGER, 
			external BOOL,
			fromUserID INTEGER, 
			toUserID INTEGER, 
			body text, 
			created datetime,
			changed datetime,
			deleted datetime 
		);`,

		`CREATE TABLE attachments (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			articleID integer,
			URI         text,
			FileName    varchar(255),
			ContentType varchar(50),
			Size        integer,
			created datetime,
			changed datetime,
			deleted datetime
		);`,

		`CREATE TABLE audits (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			ticketID INTEGER, 
			articleID INTEGER, 
			attachmentID INTEGER,
			userID INTEGER, 
			message TEXT,
			extra TEXT,
			created datetime
		);`,

		`INSERT INTO _migrations VALUES (1)`,
	)
	return err
}
