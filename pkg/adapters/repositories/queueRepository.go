package repositories

import (
	"database/sql"
	"goticka/pkg/domain/queue"
	"log"
	"time"
)

type QueueRepositoryInterface interface {
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
