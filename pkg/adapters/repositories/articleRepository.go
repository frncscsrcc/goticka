package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/article"
	"log"
	"time"
)

type ArticleRepositoryInterface interface {
	GetByID(ID int64) (article.Article, error)
	GetByTicketID(ticketID int64) ([]article.Article, error)
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

func (ar ArticleRepositorySQL) fetchArticleRow(rows *sql.Rows) ([]article.Article, error) {
	articles := make([]article.Article, 0)
	for rows.Next() {
		var a article.Article
		var deleted sql.NullTime

		errScan := rows.Scan(
			&a.ID,
			&a.Body,
			&a.From.ID,
			&a.To.ID,
			&a.Created,
			&deleted,
		)
		if errScan != nil {
			return []article.Article{}, errScan
		}

		if deleted.Valid {
			a.Deleted = deleted.Time
		}

		articles = append(articles, a)
	}
	return articles, nil
}

func (ar ArticleRepositorySQL) GetByID(ID int64) (article.Article, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.ID,
			a.Body,
			a.FromUserID,
			a.ToUserID,
			a.created,
			a.deleted
		FROM articles a
		WHERE a.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return article.Article{}, err
	}

	defer rows.Close()

	articles, err := ar.fetchArticleRow(rows)
	if err != nil {
		return article.Article{}, err
	}
	if len(articles) == 0 {
		return article.Article{}, errors.New("article not found")
	}

	return articles[0], nil
}

func (ar ArticleRepositorySQL) GetByTicketID(ticketID int64) ([]article.Article, error) {
	rows, err := ar.db.Query(`
		SELECT
			a.ID,
			a.Body,
			a.FromUserID,
			a.ToUserID,
			a.created,
			a.deleted
		FROM articles a
		WHERE a.ticketId = ?`,

		ticketID,
	)

	if err != nil {
		return []article.Article{}, err
	}

	defer rows.Close()

	articles, err := ar.fetchArticleRow(rows)
	if err != nil {
		return []article.Article{}, err
	}
	if len(articles) == 0 {
		return []article.Article{}, errors.New("articles not found")
	}

	return articles, nil
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
