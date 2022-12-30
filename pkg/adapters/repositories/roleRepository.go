package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/role"
	"goticka/pkg/domain/user"
	"log"
	"time"
)

type RoleRepositoryInterface interface {
	GetByID(ID int64) (role.Role, error)
	GetByName(roleName string) (role.Role, error)
	GetByUserID(userID int64) ([]role.Role, error)
	Create(r role.Role) (role.Role, error)

	AddRoleToUser(user user.User, role role.Role) error
	RemoveRoleFromUser(user user.User, role role.Role) error
}

type RoleRepositorySQL struct {
	db *sql.DB
}

func NewRoleRepositorySQL(db *sql.DB) *RoleRepositorySQL {
	return &RoleRepositorySQL{
		db: db,
	}
}

func (rr RoleRepositorySQL) fetchQueueRow(rows *sql.Rows) ([]role.Role, error) {
	roles := make([]role.Role, 0)
	for rows.Next() {
		var r role.Role
		var deleted sql.NullTime

		errScan := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.Created,
			&r.Changed,
			&deleted,
		)
		if errScan != nil {
			return []role.Role{}, errScan
		}

		if deleted.Valid {
			r.Deleted = deleted.Time
		}

		roles = append(roles, r)
	}
	return roles, nil
}

func (rr RoleRepositorySQL) GetByID(ID int64) (role.Role, error) {
	rows, err := rr.db.Query(`
		SELECT
			r.ID,
			r.name,
			r.description,
			r.created,
			r.changed,
			r.deleted
		FROM roles r
		WHERE r.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return role.Role{}, err
	}

	defer rows.Close()

	roles, err := rr.fetchQueueRow(rows)
	if err != nil {
		return role.Role{}, err
	}
	if len(roles) == 0 {
		return role.Role{}, errors.New("role not found")
	}

	return roles[0], nil
}

func (rr RoleRepositorySQL) GetByName(roleName string) (role.Role, error) {
	rows, err := rr.db.Query(`
		SELECT
			r.ID,
			r.name,
			r.description,
			r.created,
			r.changed,
			r.deleted
		FROM roles r
		WHERE r.name = ?
		LIMIT 1`,

		roleName,
	)

	if err != nil {
		return role.Role{}, err
	}

	defer rows.Close()

	roles, err := rr.fetchQueueRow(rows)
	if err != nil {
		return role.Role{}, err
	}
	if len(roles) == 0 {
		return role.Role{}, errors.New("role not found")
	}

	return roles[0], nil
}

func (rr RoleRepositorySQL) GetByUserID(userID int64) ([]role.Role, error) {
	rows, err := rr.db.Query(`
		SELECT
			r.ID,
			r.name,
			r.description,
			r.created,
			r.changed,
			r.deleted
		FROM roles r
		JOIN users_roles ur on ur.roleID = r.id
		WHERE ur.userId = ?`,

		userID,
	)

	if err != nil {
		return []role.Role{}, err
	}

	defer rows.Close()

	roles, err := rr.fetchQueueRow(rows)
	if err != nil {
		return []role.Role{}, err
	}

	return roles, nil
}

func (ar RoleRepositorySQL) Create(r role.Role) (role.Role, error) {
	log.Print("Creating a role")

	now := time.Now()
	res, err := ar.db.Exec(`
		INSERT INTO Roles 
			(name, description, created, changed)
		VALUES (?, ?, ?, ?);`,
		r.Name, r.Description, now, now,
	)

	if err != nil {
		return role.Role{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return role.Role{}, err
	}

	r.ID = id
	return r, nil
}

func (rr RoleRepositorySQL) AddRoleToUser(user user.User, role role.Role) error {
	now := time.Now()
	_, err := rr.db.Exec(`
		INSERT INTO Users_Roles
			(
				userID,
				roleID,
				created 
			)
		VALUES (?, ?, ?);`,

		user.ID,
		role.ID,
		now,
	)

	return err
}

func (rr RoleRepositorySQL) RemoveRoleFromUser(user user.User, role role.Role) error {
	_, err := rr.db.Exec(`
		DELETE FROM Users_Roles
		WHERE userID = ? AND roleID = ?;`,

		user.ID,
		role.ID,
	)

	return err
}
