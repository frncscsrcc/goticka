package repositories

import (
	"database/sql"
	"goticka/pkg/domain/article"
	"log"
	"time"
)

type ArticleRepositoryInterface interface {
	StoreArticle(a article.Article, ticketID int64) (article.Article, error)
}

type ArticleRepositorySQL struct {
	db                   *sql.DB
	attachmentRepository AttachmentRepositoryInterface
}

func NewArticleRepositorySQL(db *sql.DB, attachmentRepository AttachmentRepositoryInterface) *ArticleRepositorySQL {
	return &ArticleRepositorySQL{
		db:                   db,
		attachmentRepository: attachmentRepository,
	}
}

func (ar ArticleRepositorySQL) StoreArticle(a article.Article, ticketID int64) (article.Article, error) {
	log.Print("Storing an article")

	res, err := ar.db.Exec(`
		INSERT INTO Articles 
			(ticketID, body, fromUserID, toUserID, created)
		VALUES (?, ?, ?, ?, ?);`,
		ticketID, a.Body, a.From.ID, a.To.ID, time.Now(),
	)

	if err != nil {
		return article.Article{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return article.Article{}, err
	}

	createdArticle := a
	createdArticle.ID = id

	if len(a.Attachments) > 0 {
		for _, attachment := range a.Attachments {
			if att, err := ar.attachmentRepository.StoreAttachment(attachment, createdArticle.ID); err != nil {
				return article.Article{}, err
			} else {
				createdArticle.Attachments = append(createdArticle.Attachments, att)
			}
		}
	}

	return createdArticle, nil
}
