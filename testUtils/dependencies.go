package testUtils

import (
	"goticka/pkg/adapters/cache"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
)

func ResetTestDependencies() {
	dbConn := NewTestDB()

	cache := cache.GetInMemoryCache()
	passwordHasher := repositories.NewPlainTextPasswordHasher()
	userRepository := repositories.NewUserRepositorySQL(dbConn, passwordHasher)
	binaryStorer := repositories.NewAttachmentBinaryStorerFS("./")
	attachmentRepository := repositories.NewAttachmentRepositorySQL(dbConn, binaryStorer)
	articleRepository := repositories.NewArticleRepositorySQL(dbConn, attachmentRepository)
	queueRepository := repositories.NewQueueRepositorySQL(dbConn)
	auditRepository := repositories.NewAuditRepositorySQL(dbConn)
	ticketRepository := repositories.NewTicketRepositorySQL(dbConn, articleRepository)

	fakeDependencies := dependencies.Dependencies{
		Testing:              true,
		Cache:                cache,
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
