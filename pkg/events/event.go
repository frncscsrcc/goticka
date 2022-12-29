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
	eventMap map[EventType][]func(LocalEvent) error
}

var localEventHandler *LocalEventHandler

func init() {
	localEventHandler = &LocalEventHandler{
		eventMap: make(map[EventType][]func(LocalEvent) error),
	}
}

func Handler() *LocalEventHandler {
	return localEventHandler
}

func (handler *LocalEventHandler) RegisterCallBack(
	eventType EventType,
	cb func(LocalEvent) error,
) {
	_, exists := handler.eventMap[eventType]
	if !exists {
		handler.eventMap[eventType] = make([]func(LocalEvent) error, 0)
	}
	handler.eventMap[eventType] = append(handler.eventMap[eventType], cb)
}

func (handler *LocalEventHandler) SendSyncLocalEvent(e LocalEvent) error {
	if e.EventType == UNDEFINED {
		return errors.New("missing EventType, skipping")
	}

	callbacks, _ := handler.eventMap[e.EventType]
	for _, callback := range callbacks {
		err := callback(e)
		if err != nil {
			return err
		}
	}
	return nil
}
