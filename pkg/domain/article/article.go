package article

import (
	"errors"
	"goticka/pkg/domain/attachment"
	"goticka/pkg/domain/user"
)

type Article struct {
	ID          int64
	From        user.User
	To          user.User
	Body        string
	Attachments []attachment.Attachment
}

func (a Article) Validate() error {
	if a.From.ID == 0 {
		return errors.New("missing 'from' in article")
	}
	if a.To.ID == 0 {
		return errors.New("missing 'to' in article")
	}
	if a.From.ID == a.To.ID {
		return errors.New("sender and receiver can not be the same user")
	}
	if a.Body == "" {
		return errors.New("missing 'body' in article")
	}
	for _, attachment := range a.Attachments {
		if err := attachment.Validate(); err != nil {
			return err
		}
	}
	return nil
}
