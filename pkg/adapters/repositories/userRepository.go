package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"goticka/pkg/domain/user"
	"log"
	"time"
)

type PasswordHasherInterface interface {
	Hash(string) string
}

type PlainTextPasswordHasher struct{}

func NewPlainTextPasswordHasher() PlainTextPasswordHasher {
	return PlainTextPasswordHasher{}
}

func (h PlainTextPasswordHasher) Hash(plain string) string {
	return plain
}

// -------

type UserRepositoryInterface interface {
	CreateUser(u user.User) (user.User, error)
	GetByID(ID int64) (user.User, error)
	GetByUserName(userName string) (user.User, error)
	GetByUserNameAndPassword(userName string, password string) (user.User, error)
}

type UserRepositorySQL struct {
	db     *sql.DB
	hasher PasswordHasherInterface
}

func NewUserRepositorySQL(db *sql.DB, hasher PasswordHasherInterface) *UserRepositorySQL {
	return &UserRepositorySQL{
		db:     db,
		hasher: hasher,
	}
}

func (ur UserRepositorySQL) CreateUser(u user.User) (user.User, error) {
	log.Print("Creating an user")

	hashedPassword := ur.hasher.Hash(u.Password)

	res, err := ur.db.Exec(`
		INSERT INTO Users
			(username, password, created)
		VALUES (?, ?, ?);`, u.UserName, hashedPassword, time.Now())

	if err != nil {
		return user.User{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return user.User{}, err
	}
	u.ID = id

	return u, nil
}

func (ur UserRepositorySQL) fetchUserRow(rows *sql.Rows) ([]user.User, error) {
	users := make([]user.User, 0)
	for rows.Next() {
		var u user.User
		errScan := rows.Scan(
			&u.ID,
			&u.UserName,
		)
		if errScan != nil {
			return []user.User{}, errScan
		}
		users = append(users, u)
	}
	return users, nil
}

func (ur UserRepositorySQL) GetByID(ID int64) (user.User, error) {
	rows, err := ur.db.Query(`
		SELECT
			ID,
			username
		FROM users 
		WHERE users.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return user.User{}, err
	}

	defer rows.Close()

	users, err := ur.fetchUserRow(rows)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, errors.New("user not found")
	}

	return users[0], nil
}

func (ur UserRepositorySQL) GetByUserName(userName string) (user.User, error) {
	rows, err := ur.db.Query(`
		SELECT
			ID,
			username
		FROM users 
		WHERE users.username = ?
		LIMIT 1`,

		userName,
	)

	if err != nil {
		return user.User{}, err
	}

	defer rows.Close()

	users, err := ur.fetchUserRow(rows)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, errors.New("user not found")
	}

	return users[0], nil
}

func (ur UserRepositorySQL) GetByUserNameAndPassword(userName string, password string) (user.User, error) {
	hashedPassword := ur.hasher.Hash(password)
	rows, err := ur.db.Query(`
		SELECT
			ID,
			username
		FROM users
		WHERE users.username = ? AND users.password = ?
		LIMIT 1`,

		userName, hashedPassword,
	)

	if err != nil {
		return user.User{}, err
	}

	defer rows.Close()

	users, err := ur.fetchUserRow(rows)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, errors.New("user not found")
	}

	fmt.Print(users)

	return users[0], nil
}
