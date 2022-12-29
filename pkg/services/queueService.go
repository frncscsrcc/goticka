package services

import (
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/queue"
	"goticka/pkg/events"
	"log"
)

type QueueService struct {
	queueRepository repositories.QueueRepositoryInterface
}

func NewQueueService() QueueService {
	return QueueService{}
}

func (qs QueueService) Create(q queue.Queue) (queue.Queue, error) {
	if validationError := q.Validate(); validationError != nil {
		return queue.Queue{}, validationError
	}

	createdQueue, err := dependencies.DI().QueueRepository.Create(q)
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
