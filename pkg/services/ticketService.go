package services

import (
	"errors"
	"goticka/pkg/adapters/cache"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/config"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/ticket"
	"goticka/pkg/events"
	"log"
	"strconv"
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

	events.Handler().SendLocalEvent(events.LocalEvent{
		EventType: events.TICKET_CREATED,
		TicketID:  createdTicket.ID,
	})

	return createdTicket, err
}

func (ts TicketService) GetByID(id int64) (ticket.Ticket, error) {
	// Check in the cache
	cached := dependencies.DI().Cache.Get(cache.Item{
		Type: "ticket",
		Key:  strconv.FormatInt(id, 10),
	})
	if cached.IsValid() {
		if value, ok := cached.Value.(ticket.Ticket); ok {
			return value, nil
		}
	}

	t, err := dependencies.DI().TicketRepository.GetByID(id)
	if err != nil {
		return ticket.Ticket{}, err
	}

	// Save in cache
	dependencies.DI().Cache.Set(cache.Item{
		Type:  "ticket",
		Key:   strconv.FormatInt(id, 10),
		Value: t,
		TTL:   config.GetConfig().Cache.TicketTTL,
	})

	return t, nil
}

func (ts TicketService) EnrichTicketInfo(t ticket.Ticket) (ticket.Ticket, error) {
	// Check in the cache
	cached := dependencies.DI().Cache.Get(cache.Item{
		Type: "ticket",
		Key:  strconv.FormatInt(t.ID, 10) + "_FULL",
	})
	if cached.IsValid() {
		if value, ok := cached.Value.(ticket.Ticket); ok {
			return value, nil
		}
	}

	enrichedTicket := t

	queueService := NewQueueService()
	userService := NewUserService()

	articleRepo := dependencies.DI().ArticleRepository
	attachmentRepo := dependencies.DI().AttachmentRepository

	// Add the queue
	queue, queueError := queueService.GetByID(enrichedTicket.Queue.ID)
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
		from, errorFrom := userService.GetByID(article.From.ID)
		if errorFrom != nil {
			return ticket.Ticket{}, errorFrom
		}
		articles[i].From = from

		to, errorTo := userService.GetByID(article.To.ID)
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

	// Save in cache
	dependencies.DI().Cache.Set(cache.Item{
		Type:  "ticket",
		Key:   strconv.FormatInt(enrichedTicket.ID, 10) + "_FULL",
		Value: t,
		TTL:   config.GetConfig().Cache.TicketTTL,
	})

	return enrichedTicket, nil
}

func Skip(i ...interface{}) {

}

func (ts TicketService) Delete(t ticket.Ticket) error {
	if t.ID == 0 {
		return errors.New("can not delete an invalid ticket")
	}

	// Check in the cache
	dependencies.DI().Cache.Delete(cache.Item{
		Type: "ticket",
		Key:  strconv.FormatInt(t.ID, 10),
	})
	dependencies.DI().Cache.Delete(cache.Item{
		Type: "ticket",
		Key:  strconv.FormatInt(t.ID, 10) + "_FULL",
	})

	err := dependencies.DI().TicketRepository.Delete(t)

	if err == nil {
		events.Handler().SendLocalEvent(events.LocalEvent{
			EventType: events.TICKET_DELETED,
			TicketID:  t.ID,
		})
	}

	return err
}
