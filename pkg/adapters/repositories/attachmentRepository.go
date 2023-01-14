package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/attachment"
	"log"
	"time"
)

type AttachmentRepositoryInterface interface {
	GetByID(ID int64) (attachment.Attachment, error)
	GetByArticleID(articleId int64) ([]attachment.Attachment, error)
	StoreAttachment(attachment.Attachment, int64) (attachment.Attachment, error)
}

type AttachmentRepositorySQL struct {
	db               *sql.DB
	binaryRepository BinaryRepositoryInterface
}

func NewAttachmentRepositorySQL(db *sql.DB, binaryRepository BinaryRepositoryInterface) *AttachmentRepositorySQL {
	return &AttachmentRepositorySQL{
		db:               db,
		binaryRepository: binaryRepository,
	}
}

func (ar AttachmentRepositorySQL) fetchAttachmentRow(rows *sql.Rows) ([]attachment.Attachment, error) {
	attachments := make([]attachment.Attachment, 0)
	for rows.Next() {
		var a attachment.Attachment
		var deleted sql.NullTime

		errScan := rows.Scan(
			&a.ID,
			&a.URI,
			&a.FileName,
			&a.ContentType,
			&a.Size,
			&a.Created,
			&a.Changed,
			&deleted,
		)
		if errScan != nil {
			return []attachment.Attachment{}, errScan
		}

		if deleted.Valid {
			a.Deleted = deleted.Time
		}

		attachments = append(attachments, a)
	}
	return attachments, nil
}

func (ar AttachmentRepositorySQL) GetByID(ID int64) (attachment.Attachment, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.ID,
			a.URI,
			a.FileName,
			a.ContentType,
			a.Size,
			a.created,
			a.changed,
			a.deleted
		FROM attachments a
		WHERE a.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return attachment.Attachment{}, err
	}

	defer rows.Close()

	attachments, err := ar.fetchAttachmentRow(rows)
	if err != nil {
		return attachment.Attachment{}, err
	}
	if len(attachments) == 0 {
		return attachment.Attachment{}, errors.New("queue not found")
	}

	return attachments[0], nil
}

func (ar AttachmentRepositorySQL) GetByArticleID(articleId int64) ([]attachment.Attachment, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.ID,
			a.URI,
			a.FileName,
			a.ContentType,
			a.Size,
			a.created,
			a.changed,
			a.deleted
		FROM attachments a
		WHERE a.articleId = ?`,

		articleId,
	)

	if err != nil {
		return []attachment.Attachment{}, err
	}

	defer rows.Close()

	attachments, err := ar.fetchAttachmentRow(rows)
	if err != nil {
		return []attachment.Attachment{}, err
	}

	return attachments, nil
}

func (ar AttachmentRepositorySQL) StoreAttachment(a attachment.Attachment, articleID int64) (attachment.Attachment, error) {
	log.Print("Storing an attachment")

	a, storedBinaryError := ar.binaryRepository.StoreBinary(a)
	if storedBinaryError != nil {
		return attachment.Attachment{}, storedBinaryError
	}

	now := time.Now()
	res, err := ar.db.Exec(`
		INSERT INTO Attachments (
			articleID,
			URI,
			filename,
			contentType,
			size,
			created,
			changed
		)
		VALUES (?, ?, ?, ?, ?, ?, ?);`,
		articleID, a.URI, a.FileName, a.ContentType, a.Size, now, now)

	if err != nil {
		return attachment.Attachment{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return attachment.Attachment{}, err
	}

	a.ID = id

	return a, nil
}
