package dependencies

import (
	"goticka/pkg/adapters/cache"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/config"
	"goticka/pkg/db"
)

func GetConfig() config.Config {
	return config.GetConfig()
}

type Dependencies struct {
	Testing              bool
	Cache                cache.CacheInterface
	PasswordHasher       repositories.PasswordHasherInterface
	UserRepository       repositories.UserRepositoryInterface
	QueueRepository      repositories.QueueRepositoryInterface
	TicketRepository     repositories.TicketRepositoryInterface
	ArticleRepository    repositories.ArticleRepositoryInterface
	AttachmentRepository repositories.AttachmentRepositoryInterface
	BinaryStorer         repositories.AttachmentBinaryStorerFS
	AuditRepository      repositories.AuditRepositoryInterface
}

var dependencies Dependencies

func init() {
	dbConn := db.GetDB()

	cache := cache.GetInMemoryCache()
	passwordHasher := repositories.NewPlainTextPasswordHasher()
	userRepository := repositories.NewUserRepositorySQL(dbConn, passwordHasher)
	binaryStorer := repositories.NewAttachmentBinaryStorerFS("./")
	attachmentRepository := repositories.NewAttachmentRepositorySQL(dbConn, binaryStorer)
	articleRepository := repositories.NewArticleRepositorySQL(dbConn, attachmentRepository)
	queueRepository := repositories.NewQueueRepositorySQL(dbConn)
	auditRepository := repositories.NewAuditRepositorySQL(dbConn)
	ticketRepository := repositories.NewTicketRepositorySQL(dbConn, articleRepository)

	dependencies = Dependencies{
		Testing: false,

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
}

func OverwriteDependencies(newDependencies Dependencies) Dependencies {
	dependencies = newDependencies
	return dependencies
}

func DI() Dependencies {
	return dependencies
}
