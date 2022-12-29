package services

import (
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/audit"
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
