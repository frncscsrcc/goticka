package audit

import (
	"goticka/pkg/events"
	"time"
)

type Audit struct {
	ID           int64
	EventType    events.EventType
	TicketID     int64
	ArticleID    int64
	AttachmentID int64
	UserID       int64
	Created      time.Time
	Message      string
	Extra        string
}
