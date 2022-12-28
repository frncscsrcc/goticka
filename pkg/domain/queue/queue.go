package queue

import (
	"errors"
)

type Queue struct {
	ID          int64
	Name        string
	Description string
}

func (q Queue) Validate() error {
	if q.Name == "" {
		return errors.New("missing 'name' in queue")
	}
	return nil
}
