package services

import (
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/audit"
	"goticka/pkg/events"
	"log"
)

type AuditService struct {
	auditRepository repositories.AuditRepositoryInterface
}

func NewAuditService() AuditService {
	return AuditService{
		auditRepository: dependencies.DI().AuditRepository,
	}
}

func (as AuditService) Create(a audit.Audit) (audit.Audit, error) {
	// Validation here

	createdAudit, err := as.auditRepository.Save(a)
	if err != nil {
		return audit.Audit{}, err
	}
	log.Printf("created audit record %d\n", createdAudit.ID)

	return createdAudit, err
}

func (as AuditService) GetByID(id int64) (audit.Audit, error) {
	return as.auditRepository.GetByID(id)
}

// Register to local events
func init() {
	auditService := NewAuditService()
	eventHandler := events.Handler()

	// -------------------------------------
	// User created Audit
	// -------------------------------------
	eventHandler.RegisterSyncCallBack(
		events.USER_CREATED,
		func(event events.LocalEvent) error {
			auditService.Create(audit.Audit{
				UserID:  event.UserID,
				Message: "New user created",
			})
			return nil
		})

	// -------------------------------------
	// Ticket created Audit
	// -------------------------------------
	eventHandler.RegisterSyncCallBack(
		events.TICKET_CREATED,
		func(event events.LocalEvent) error {
			auditService.Create(audit.Audit{
				TicketID: event.TicketID,
				Message:  "New ticket created",
			})
			return nil
		})
}
