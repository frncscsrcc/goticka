package services

import (
	"errors"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/ticket"
	"log"
)

type TicketService struct {
	ticketRepository repositories.TicketRepositoryInterface
}

func NewTicketService() TicketService {
	return TicketService{}
}

func (ts TicketService) Create(t ticket.Ticket) (ticket.Ticket, error) {
	if validationError := t.Validate(); validationError != nil {
		return ticket.Ticket{}, validationError
	}
	createdTicket, err := dependencies.DI().TicketRepository.StoreTicket(t)
	if err != nil {
		return ticket.Ticket{}, err
	}
	log.Printf("created ticket %d\n", createdTicket.ID)
	return createdTicket, err
}

func (ts TicketService) GetByID(id int64) (ticket.Ticket, error) {
	return dependencies.DI().TicketRepository.GetByID(id)
}

func (ts TicketService) EnrichTicketInfo(t ticket.Ticket) (ticket.Ticket, error) {
	enrichedTicket := t

	queueRepo := dependencies.DI().QueueRepository
	articleRepo := dependencies.DI().ArticleRepository
	attachmentRepo := dependencies.DI().AttachmentRepository
	userRepo := dependencies.DI().UserRepository

	// Add the queue
	queue, queueError := queueRepo.GetByID(enrichedTicket.Queue.ID)
	if queueError != nil {
		return ticket.Ticket{}, queueError
	}
	enrichedTicket.Queue = queue

	// Add the articles
	articles, articlesFetchError := articleRepo.GetByTicketID(t.ID)
	if articlesFetchError != nil {
		return ticket.Ticket{}, articlesFetchError
	}

	for i, article := range t.Articles {
		// Add the users
		from, errorFrom := userRepo.GetByID(article.From.ID)
		if errorFrom != nil {
			return ticket.Ticket{}, errorFrom
		}
		articles[i].From = from

		to, errorTo := userRepo.GetByID(article.To.ID)
		if errorFrom != nil {
			return ticket.Ticket{}, errorTo
		}
		articles[i].To = to

		// Add the attachments
		attachments, attachmentsError := attachmentRepo.GetByArticleID(article.ID)
		if attachmentsError != nil {
			return ticket.Ticket{}, attachmentsError
		}
		articles[i].Attachments = attachments
	}

	enrichedTicket.Articles = articles

	// Add the attachments
	for _, article := range enrichedTicket.Articles {
		Skip(article)
	}

	Skip(queueRepo, attachmentRepo)
	return enrichedTicket, nil
}

func Skip(i ...interface{}) {

}

func (ts TicketService) Delete(t ticket.Ticket) error {
	if t.ID == 0 {
		return errors.New("can not delete an invalid ticket")
	}
	return dependencies.DI().TicketRepository.Delete(t)
}
