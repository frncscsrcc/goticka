package user

import "time"

type User struct {
	ID       int64
	UserName string
	Password string

	Created time.Time
	Deleted time.Time
}
