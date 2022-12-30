package user

import (
	"goticka/pkg/domain/role"
	"time"
)

type User struct {
	ID       int64
	External bool
	UserName string
	Password string
	Email    string
	Created  time.Time
	Changed  time.Time
	Deleted  time.Time

	Roles []role.Role
}
