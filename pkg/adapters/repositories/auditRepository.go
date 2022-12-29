package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/audit"
	"log"
	"time"
)

type AuditRepositoryInterface interface {
	GetByID(ID int64) (audit.Audit, error)
	Save(a audit.Audit) (audit.Audit, error)
}

type AuditRepositorySQL struct {
	db *sql.DB
}

func NewAuditRepositorySQL(db *sql.DB) *AuditRepositorySQL {
	return &AuditRepositorySQL{
		db: db,
	}
}

func (ar AuditRepositorySQL) Save(a audit.Audit) (audit.Audit, error) {
	log.Print("Creating an audit record")

	createdAudit := a

	now := time.Now()
	res, err := ar.db.Exec(`
		INSERT INTO Audits
			(
				ticketID, 
				articleID, 
				attachmentID, 
				userID,
				message,
				extra,
				created
			)
		VALUES (?, ?, ?, ?, ?, ?, ?);`,

		a.TicketID,
		a.ArticleID,
		a.AttachmentID,
		a.UserID,
		a.Message,
		a.Extra,
		now,
	)

	if err != nil {
		return audit.Audit{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return audit.Audit{}, err
	}
	createdAudit.ID = id

	return createdAudit, nil
}

func (ar AuditRepositorySQL) fetchAuditRow(rows *sql.Rows) ([]audit.Audit, error) {
	audits := make([]audit.Audit, 0)
	for rows.Next() {
		var a audit.Audit
		var ticketID sql.NullInt64
		var articleID sql.NullInt64
		var attachmentID sql.NullInt64
		var userID sql.NullInt64

		errScan := rows.Scan(
			&a.ID,
			&ticketID,
			&articleID,
			&attachmentID,
			&userID,
			&a.Message,
			&a.Extra,
			&a.Created,
		)

		if errScan != nil {
			return []audit.Audit{}, errScan
		}

		if ticketID.Valid {
			a.TicketID = ticketID.Int64
		}
		if articleID.Valid {
			a.ArticleID = articleID.Int64
		}
		if attachmentID.Valid {
			a.AttachmentID = attachmentID.Int64
		}
		if userID.Valid {
			a.UserID = userID.Int64
		}

		audits = append(audits, a)
	}
	return audits, nil
}

func (ar AuditRepositorySQL) GetByID(ID int64) (audit.Audit, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.ID,
			a.ticketID,
			a.articleID,
			a.attachmentID,
			a.userID,
			a.Message,
			a.Extra,
			a.Created
		FROM Audits a
		WHERE a.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return audit.Audit{}, err
	}

	defer rows.Close()

	audits, err := ar.fetchAuditRow(rows)
	if err != nil {
		return audit.Audit{}, err
	}
	if len(audits) == 0 {
		return audit.Audit{}, errors.New("audit record not found")
	}

	return audits[0], nil
}
