package services

import (
	"goticka/pkg/adapters/cache"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/queue"
	"goticka/pkg/events"
	"log"
	"strconv"
	"time"
)

type QueueService struct {
	queueRepository repositories.QueueRepositoryInterface
}

func NewQueueService() QueueService {
	return QueueService{
		queueRepository: dependencies.DI().QueueRepository,
	}
}

func (qs QueueService) Create(q queue.Queue) (queue.Queue, error) {
	if validationError := q.Validate(); validationError != nil {
		return queue.Queue{}, validationError
	}

	createdQueue, err := qs.queueRepository.Create(q)
	if err != nil {
		return queue.Queue{}, err
	}
	log.Printf("created queue %d\n", createdQueue.ID)

	events.Handler().SendLocalEvent(events.LocalEvent{
		EventType: events.QUEUE_CREATED,
		QueueID:   createdQueue.ID,
	})

	return createdQueue, err
}

func (qs QueueService) GetByID(id int64) (queue.Queue, error) {
	// Check in the cache
	cached := dependencies.DI().Cache.Get(cache.Item{
		Type: "queue",
		Key:  strconv.FormatInt(id, 10),
	})
	if cached.IsValid() {
		if value, ok := cached.Value.(queue.Queue); ok {
			return value, nil
		}
	}

	q, err := qs.queueRepository.GetByID(id)
	if err != nil {
		return queue.Queue{}, err
	}

	// Save in cache
	dependencies.DI().Cache.Set(cache.Item{
		Type:  "queue",
		Key:   strconv.FormatInt(id, 10),
		Value: q,
		TTL:   10 * time.Minute,
	})

	return q, nil
}
