package dependencies

import (
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/config"
	"goticka/pkg/db"
)

func GetConfig() config.Config {
	return config.GetConfig()
}

type Dependencies struct {
	Testing              bool
	PasswordHasher       repositories.PasswordHasherInterface
	UserRepository       repositories.UserRepositoryInterface
	QueueRepository      repositories.QueueRepositoryInterface
	TicketRepository     repositories.TicketRepositoryInterface
	ArticleRepository    repositories.ArticleRepositoryInterface
	AttachmentRepository repositories.AttachmentRepositoryInterface
	BinaryStorer         repositories.AttachmentBinaryStorerFS
}

var dependencies Dependencies

func init() {
	dbConn := db.GetDB()

	passwordHasher := repositories.NewPlainTextPasswordHasher()
	userRepository := repositories.NewUserRepositorySQL(dbConn, passwordHasher)
	binaryStorer := repositories.NewAttachmentBinaryStorerFS("./")
	attachmentRepository := repositories.NewAttachmentRepositorySQL(dbConn, binaryStorer)
	articleRepository := repositories.NewArticleRepositorySQL(dbConn, attachmentRepository)
	queueRepository := repositories.NewQueueRepositorySQL(dbConn)
	ticketRepository := repositories.NewTicketRepositorySQL(dbConn, articleRepository)

	dependencies = Dependencies{
		Testing: false,

		PasswordHasher:       passwordHasher,
		UserRepository:       userRepository,
		QueueRepository:      queueRepository,
		TicketRepository:     ticketRepository,
		ArticleRepository:    articleRepository,
		AttachmentRepository: attachmentRepository,
		BinaryStorer:         binaryStorer,
	}
}

func OverwriteDependencies(newDependencies Dependencies) Dependencies {
	dependencies = newDependencies
	return dependencies
}

func DI() Dependencies {
	return dependencies
}
