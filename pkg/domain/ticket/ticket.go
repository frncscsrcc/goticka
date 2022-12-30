package ticket

import (
	"errors"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/queue"
	"time"
)

type TicketStatus int64

const NEW TicketStatus = 1
const OPEN TicketStatus = 2
const PENDING_INTERNAL TicketStatus = 4
const PENDING_EXTERNAL TicketStatus = 5
const WAITING_APPROVAL TicketStatus = 6
const CLOSED TicketStatus = 7

type Ticket struct {
	ID       int64
	Status   TicketStatus
	Queue    queue.Queue
	Subject  string
	Articles []article.Article

	Created time.Time
	Changed time.Time
	Deleted time.Time
}

func (t Ticket) Validate() error {
	if t.Queue.ID == 0 {
		return errors.New("missing 'queue' in ticket")
	}
	if t.Subject == "" {
		return errors.New("missing 'subject' in ticket")
	}
	if len(t.Articles) == 0 {
		return errors.New("missing 'articles' in ticket")
	}
	for _, article := range t.Articles {
		if err := article.Validate(); err != nil {
			return err
		}
	}
	return nil
}
