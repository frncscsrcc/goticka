package repositories

import (
	"database/sql"
	"errors"
	"goticka/pkg/domain/queue"
	"log"
	"time"
)

type QueueRepositoryInterface interface {
	GetByID(ID int64) (queue.Queue, error)
	Create(q queue.Queue) (queue.Queue, error)
}

type QueueRepositorySQL struct {
	db *sql.DB
}

func NewQueueRepositorySQL(db *sql.DB) *QueueRepositorySQL {
	return &QueueRepositorySQL{
		db: db,
	}
}

func (qr QueueRepositorySQL) fetchQueueRow(rows *sql.Rows) ([]queue.Queue, error) {
	queues := make([]queue.Queue, 0)
	for rows.Next() {
		var q queue.Queue
		var deleted sql.NullTime

		errScan := rows.Scan(
			&q.ID,
			&q.Name,
			&q.Description,
			&q.Created,
			&deleted,
		)
		if errScan != nil {
			return []queue.Queue{}, errScan
		}

		if deleted.Valid {
			q.Deleted = deleted.Time
		}

		queues = append(queues, q)
	}
	return queues, nil
}

func (qr QueueRepositorySQL) GetByID(ID int64) (queue.Queue, error) {
	rows, err := qr.db.Query(`
		SELECT
			q.ID,
			q.name,
			q.description,
			q.created,
			q.deleted
		FROM queues q
		WHERE q.id = ?
		LIMIT 1`,

		ID,
	)

	if err != nil {
		return queue.Queue{}, err
	}

	defer rows.Close()

	queues, err := qr.fetchQueueRow(rows)
	if err != nil {
		return queue.Queue{}, err
	}
	if len(queues) == 0 {
		return queue.Queue{}, errors.New("queue not found")
	}

	return queues[0], nil
}

func (ar QueueRepositorySQL) Create(q queue.Queue) (queue.Queue, error) {
	log.Print("Creating a queue")

	res, err := ar.db.Exec(`
		INSERT INTO Queues 
			(name, description, created)
		VALUES (?, ?, ?);`,
		q.Name, q.Description, time.Now(),
	)

	if err != nil {
		return queue.Queue{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return queue.Queue{}, err
	}

	q.ID = id
	return q, nil
}
