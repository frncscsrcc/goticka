package queue

import (
	"errors"
	"time"
)

type Queue struct {
	ID          int64
	Name        string
	Description string

	Created time.Time
	Deleted time.Time
}

func (q Queue) Validate() error {
	if q.Name == "" {
		return errors.New("missing 'name' in queue")
	}
	return nil
}
