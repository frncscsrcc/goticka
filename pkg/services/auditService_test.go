package services

import (
	"goticka/pkg/domain/audit"
	"goticka/testUtils"
	"testing"
)

func TestCreateAndRetrieveAudit(t *testing.T) {
	testUtils.ResetTestDependencies()

	ticketID := int64(1)
	articleID := int64(1)
	attachmentID := int64(1)
	message := "MESSAGE"
	extra := "EXTRA"

	as := NewAuditService()
	createdAudit, err := as.Create(
		audit.Audit{
			TicketID:     ticketID,
			ArticleID:    articleID,
			AttachmentID: attachmentID,
			Message:      message,
			Extra:        extra,
		},
	)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if createdAudit.ID == 0 {
		t.Error("audit ID should be initialized")
	}

	retrivedAudit, retrivedAuditError := as.GetByID(createdAudit.ID)
	if retrivedAuditError != nil {
		t.Errorf("unexpected error %s", retrivedAuditError)
	}
	if retrivedAudit.TicketID != ticketID {
		t.Errorf("wrong ticketID, expected %d, got %d", ticketID, retrivedAudit.TicketID)
	}
	if retrivedAudit.ArticleID != articleID {
		t.Errorf("wrong articleID, expected %d, got %d", articleID, retrivedAudit.ArticleID)
	}
	if retrivedAudit.AttachmentID != attachmentID {
		t.Errorf("wrong attachmentID, expected %d, got %d", attachmentID, retrivedAudit.AttachmentID)
	}
	if retrivedAudit.UserID != 0 {
		t.Errorf("wrong userID, expected %d, got %d", 0, retrivedAudit.UserID)
	}
	if retrivedAudit.Message != message {
		t.Errorf("wrong message, expected %s, got %s", message, retrivedAudit.Message)
	}
	if retrivedAudit.Extra != extra {
		t.Errorf("wrong extra, expected %s, got %s", extra, retrivedAudit.Extra)
	}
	if retrivedAudit.Created.IsZero() {
		t.Error("wrong created field (empty)")
	}
}
