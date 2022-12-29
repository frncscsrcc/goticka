package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/ticket"
	"log"
	"time"
)

type TicketRepositoryInterface interface {
	GetByID(ID int64) (ticket.Ticket, error)
	StoreTicket(t ticket.Ticket) (ticket.Ticket, error)
	Delete(t ticket.Ticket) error
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

func (tr TicketRepositorySQL) fetchTicketRow(rows *sql.Rows) ([]ticket.Ticket, error) {
	tickets := make([]ticket.Ticket, 0)
	for rows.Next() {
		var t ticket.Ticket
		var deleted sql.NullTime

		errScan := rows.Scan(
			&t.ID,
			&t.Queue.ID,
			&t.Subject,
			&t.Created,
			&t.Changed,
			&deleted,
		)
		if errScan != nil {
			return []ticket.Ticket{}, errScan
		}

		if deleted.Valid {
			t.Deleted = deleted.Time
		}

		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (tr TicketRepositorySQL) GetByID(ID int64) (ticket.Ticket, error) {
	rows, err := tr.db.Query(`
		SELECT
			t.ID,
			t.queueID,
			t.subject,
			t.created,
			t.changed,
			t.deleted
		FROM tickets t
		WHERE t.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return ticket.Ticket{}, err
	}

	defer rows.Close()

	tickets, err := tr.fetchTicketRow(rows)
	if err != nil {
		return ticket.Ticket{}, err
	}
	if len(tickets) == 0 {
		return ticket.Ticket{}, errors.New("ticket not found")
	}

	return tickets[0], nil
}

func (tr TicketRepositorySQL) StoreTicket(t ticket.Ticket) (ticket.Ticket, error) {
	log.Print("Storing a ticket")

	now := time.Now()
	res, err := tr.db.Exec(`
		INSERT INTO Tickets (
			subject,
			queueId,
			created,
			changed
		)
		VALUES (?, ?, ?, ?);`,

		t.Subject,
		t.Queue.ID,
		now,
		now,
	)

	if err != nil {
		return ticket.Ticket{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return ticket.Ticket{}, err
	}

	createdTicket := t
	createdTicket.ID = id
	createdTicket.Created = now

	articles := make([]article.Article, 0)
	for _, article := range t.Articles {
		if storedArticle, err := tr.articleRepository.StoreArticle(article, id); err != nil {
			return ticket.Ticket{}, err
		} else {
			articles = append(articles, storedArticle)
		}
	}
	createdTicket.Articles = articles

	return createdTicket, nil
}

func (tr TicketRepositorySQL) Delete(t ticket.Ticket) error {
	log.Print("Deleting a ticket")

	now := time.Now()
	_, err := tr.db.Exec(`
		UPDATE Tickets
		SET deleted = ?
		WHERE tickets.id = ?`,

		now,
		t.ID,
	)

	return err
}
