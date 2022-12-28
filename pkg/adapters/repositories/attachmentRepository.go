package repositories

import (
	"database/sql"
	"goticka/pkg/domain/attachment"
	"log"
	"time"
)

type AttachmentBinaryStorerInterface interface {
	StoreBinary(attachment.Attachment) (attachment.Attachment, error)
}

type AttachmentBinaryStorerFS struct {
	basePath string
}

func NewAttachmentBinaryStorerFS(basePath string) AttachmentBinaryStorerFS {
	return AttachmentBinaryStorerFS{
		basePath: basePath,
	}
}

func (bs AttachmentBinaryStorerFS) StoreBinary(a attachment.Attachment) (attachment.Attachment, error) {
	return a, nil
}

// -----

type AttachmentRepositoryInterface interface {
	StoreAttachment(attachment.Attachment, int64) (attachment.Attachment, error)
}

type AttachmentRepositorySQL struct {
	db           *sql.DB
	binaryStorer AttachmentBinaryStorerInterface
}

func NewAttachmentRepositorySQL(db *sql.DB, binaryStorer AttachmentBinaryStorerInterface) *AttachmentRepositorySQL {
	return &AttachmentRepositorySQL{
		db:           db,
		binaryStorer: binaryStorer,
	}
}

func (ar AttachmentRepositorySQL) StoreAttachment(a attachment.Attachment, articleID int64) (attachment.Attachment, error) {
	log.Print("Storing an attachment")

	a, storedBinaryError := ar.binaryStorer.StoreBinary(a)
	if storedBinaryError != nil {
		return attachment.Attachment{}, storedBinaryError
	}

	res, err := ar.db.Exec(`
		INSERT INTO Attachments
			(articleID, URI, filename, contentType, size)
		VALUES (?, ?, ?, ?, ?);`,
		articleID, a.URI, a.FileName, a.ContentType, time.Now())

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
