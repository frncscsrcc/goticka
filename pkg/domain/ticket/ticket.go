package ticket

import (
	"errors"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/queue"
)

type Ticket struct {
	ID       int64
	Queue    queue.Queue
	Subject  string
	Articles []article.Article
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
