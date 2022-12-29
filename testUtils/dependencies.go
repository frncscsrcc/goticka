package testUtils

import (
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
)

func ResetTestDependencies() {
	dbConn := NewTestDB()

	passwordHasher := repositories.NewPlainTextPasswordHasher()
	userRepository := repositories.NewUserRepositorySQL(dbConn, passwordHasher)
	binaryStorer := repositories.NewAttachmentBinaryStorerFS("./")
	attachmentRepository := repositories.NewAttachmentRepositorySQL(dbConn, binaryStorer)
	articleRepository := repositories.NewArticleRepositorySQL(dbConn, attachmentRepository)
	queueRepository := repositories.NewQueueRepositorySQL(dbConn)
	ticketRepository := repositories.NewTicketRepositorySQL(dbConn, articleRepository)
	auditRepository := repositories.NewAuditRepositorySQL(dbConn)

	fakeDependencies := dependencies.Dependencies{
		Testing:              true,
		PasswordHasher:       passwordHasher,
		UserRepository:       userRepository,
		QueueRepository:      queueRepository,
		TicketRepository:     ticketRepository,
		ArticleRepository:    articleRepository,
		AttachmentRepository: attachmentRepository,
		BinaryStorer:         binaryStorer,
		AuditRepository:      auditRepository,
	}

	dependencies.OverwriteDependencies(fakeDependencies)
}
