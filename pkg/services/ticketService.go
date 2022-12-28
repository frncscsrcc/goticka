package services

import (
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
