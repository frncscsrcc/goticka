package repositories

import (
	"database/sql"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/ticket"
	"log"
	"time"
)

type TicketRepositoryInterface interface {
	StoreTicket(t ticket.Ticket) (ticket.Ticket, error)
}

type TicketRepositorySQL struct {
	db                *sql.DB
	articleRepository ArticleRepositoryInterface
}

func NewTicketRepositorySQL(db *sql.DB, articleRepository ArticleRepositoryInterface) *TicketRepositorySQL {
	return &TicketRepositorySQL{
		db:                db,
		articleRepository: articleRepository,
	}
}

func (tr TicketRepositorySQL) StoreTicket(t ticket.Ticket) (ticket.Ticket, error) {
	log.Print("Storing a ticket")

	res, err := tr.db.Exec(`
		INSERT INTO Tickets
			(subject, created)
		VALUES (?, ?);`, t.Subject, time.Now())

	if err != nil {
		return ticket.Ticket{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return ticket.Ticket{}, err
	}

	createdTicket := t
	createdTicket.ID = id

	var a article.Article
	if a, err = tr.articleRepository.StoreArticle(t.Articles[0], id); err != nil {
		return ticket.Ticket{}, err
	}

	createdTicket.Articles = []article.Article{a}

	return createdTicket, nil
}
