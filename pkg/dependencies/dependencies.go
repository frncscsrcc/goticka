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
	RoleRepository       repositories.RoleRepositoryInterface
	UserRepository       repositories.UserRepositoryInterface
	QueueRepository      repositories.QueueRepositoryInterface
	TicketRepository     repositories.TicketRepositoryInterface
	ArticleRepository    repositories.ArticleRepositoryInterface
	AttachmentRepository repositories.AttachmentRepositoryInterface
	BinaryRepository     repositories.BinaryRepositoryInterface
	AuditRepository      repositories.AuditRepositoryInterface
}

var dependencies Dependencies

func init() {
	dbConn := db.GetDB()

	cache := cache.GetInMemoryCache()
	roleRepository := repositories.NewRoleRepositorySQL(dbConn)
	passwordHasher := repositories.NewPlainTextPasswordHasher()
	userRepository := repositories.NewUserRepositorySQL(dbConn, passwordHasher, roleRepository)
	binaryRepository := repositories.NewBinaryRepositoryFS(GetConfig().Storage.BasePath)
	attachmentRepository := repositories.NewAttachmentRepositorySQL(dbConn, binaryRepository)
	articleRepository := repositories.NewArticleRepositorySQL(dbConn, attachmentRepository)
	queueRepository := repositories.NewQueueRepositorySQL(dbConn)
	auditRepository := repositories.NewAuditRepositorySQL(dbConn)
	ticketRepository := repositories.NewTicketRepositorySQL(dbConn, articleRepository)

	dependencies = Dependencies{
		Testing: false,

		Cache:                cache,
		PasswordHasher:       passwordHasher,
		RoleRepository:       roleRepository,
		UserRepository:       userRepository,
		QueueRepository:      queueRepository,
		TicketRepository:     ticketRepository,
		ArticleRepository:    articleRepository,
		AttachmentRepository: attachmentRepository,
		BinaryRepository:     binaryRepository,
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
