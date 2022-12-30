package role

import (
	"errors"
	"time"
)

type Role struct {
	ID          int64
	Name        string
	Description string
	Created     time.Time
	Changed     time.Time
	Deleted     time.Time
}

func (r Role) Validate() error {
	if r.Name == "" {
		return errors.New("missing 'name' in role")
	}
	return nil
}
