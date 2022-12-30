package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/role"
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
	Delete(u user.User) error

	AddRole(user user.User, role role.Role) error
	RemoveRole(user user.User, role role.Role) error
}

type UserRepositorySQL struct {
	db             *sql.DB
	hasher         PasswordHasherInterface
	roleRepository RoleRepositoryInterface
}

func NewUserRepositorySQL(
	db *sql.DB,
	hasher PasswordHasherInterface,
	roleRepository RoleRepositoryInterface,
) *UserRepositorySQL {
	return &UserRepositorySQL{
		db:             db,
		hasher:         hasher,
		roleRepository: roleRepository,
	}
}

func (ur UserRepositorySQL) CreateUser(u user.User) (user.User, error) {
	log.Print("Creating an user")

	hashedPassword := ur.hasher.Hash(u.Password)

	now := time.Now()
	res, err := ur.db.Exec(`
		INSERT INTO Users
			(
				external,
				email,
				username,
				password, 
				created, 
				changed
			)
		VALUES (?,?, ?, ?, ?, ?);`,

		u.External,
		u.Email,
		u.UserName,
		hashedPassword,
		now,
		now,
	)

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
		var deleted sql.NullTime
		errScan := rows.Scan(
			&u.ID,
			&u.External,
			&u.Email,
			&u.UserName,
			&u.Created,
			&u.Changed,
			&deleted,
		)
		if deleted.Valid {
			u.Deleted = deleted.Time
		}
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
			external,
			email,
			username,
			created,
			changed,
			deleted
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

	// Add the roles
	roles, rolesError := ur.roleRepository.GetByUserID(ID)
	if rolesError != nil {
		return user.User{}, rolesError
	}
	users[0].Roles = roles

	return users[0], nil
}

func (ur UserRepositorySQL) GetByUserName(userName string) (user.User, error) {
	rows, err := ur.db.Query(`
		SELECT
			ID,
			external,
			email,
			username,
			created,
			changed,
			deleted
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

	// Add the roles
	roles, rolesError := ur.roleRepository.GetByUserID(users[0].ID)
	if rolesError != nil {
		return user.User{}, rolesError
	}
	users[0].Roles = roles

	return users[0], nil
}

func (ur UserRepositorySQL) GetByUserNameAndPassword(userName string, password string) (user.User, error) {
	hashedPassword := ur.hasher.Hash(password)
	rows, err := ur.db.Query(`
		SELECT
			ID,
			external,
			email,
			username,
			created,
			changed,
			deleted
		FROM users
		WHERE users.username = ? AND users.password = ? AND deleted IS NULL
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

	return users[0], nil
}

func (ur UserRepositorySQL) Delete(u user.User) error {
	log.Print("Deleting an user")

	_, err := ur.db.Exec(`
		UPDATE users
		SET deleted = ?
		WHERE users.id = ?`,

		time.Now(), u.ID)

	return err
}

func (ur UserRepositorySQL) AddRole(user user.User, role role.Role) error {
	return ur.roleRepository.AddRoleToUser(user, role)
}

func (ur UserRepositorySQL) RemoveRole(user user.User, role role.Role) error {
	return ur.roleRepository.RemoveRoleFromUser(user, role)
}
