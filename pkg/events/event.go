package events

import (
	"errors"
	"time"
)

type EventType int64

const (
	UNDEFINED EventType = iota

	TICKET_CREATED
	TICKET_DELETED

	USER_CREATED
	USER_DELETED

	QUEUE_CREATING
	QUEUE_CREATED
)

type LocalEvent struct {
	EventType    EventType
	TicketID     int64
	QueueID      int64
	ArticleID    int64
	AttachmentID int64
	UserID       int64
	Extra        string
	Time         time.Time
}

type LocalEventHandler struct {
	syncEventMap  map[EventType][]func(LocalEvent) error
	asyncEventMap map[EventType][]func(LocalEvent)
}

var localEventHandler *LocalEventHandler

func init() {
	localEventHandler = &LocalEventHandler{
		syncEventMap:  make(map[EventType][]func(LocalEvent) error),
		asyncEventMap: make(map[EventType][]func(LocalEvent)),
	}
}

func Handler() *LocalEventHandler {
	return localEventHandler
}

func (handler *LocalEventHandler) RegisterSyncCallBack(
	eventType EventType,
	cb func(LocalEvent) error,
) {
	_, exists := handler.syncEventMap[eventType]
	if !exists {
		handler.syncEventMap[eventType] = make([]func(LocalEvent) error, 0)
	}
	handler.syncEventMap[eventType] = append(handler.syncEventMap[eventType], cb)
}

func (handler *LocalEventHandler) RegisterAsyncCallBack(
	eventType EventType,
	cb func(LocalEvent),
) {
	_, exists := handler.asyncEventMap[eventType]
	if !exists {
		handler.asyncEventMap[eventType] = make([]func(LocalEvent), 0)
	}
	handler.asyncEventMap[eventType] = append(handler.asyncEventMap[eventType], cb)
}

func (handler *LocalEventHandler) SendLocalEvent(e LocalEvent) error {
	if e.EventType == UNDEFINED {
		return errors.New("missing EventType, skipping")
	}

	// Trigger async callbacks and forget
	asyncCallbacks, _ := handler.asyncEventMap[e.EventType]
	for _, callback := range asyncCallbacks {
		go callback(e)
	}

	// Trigger sync
	syncCallbacks, _ := handler.syncEventMap[e.EventType]
	for _, callback := range syncCallbacks {
		err := callback(e)
		if err != nil {
			return err
		}
	}
	return nil
}
