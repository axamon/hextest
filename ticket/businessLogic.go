package ticket

import (
	"time"

	"github.com/google/uuid"
)

// NewService creates a new service.
func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

// CreateTicket creates a new Ticket.
func (s *service) CreateTicket(ticket *Ticket) error {
	ticket.ID = uuid.New().String()
	ticket.Created = time.Now()
	ticket.Updated = time.Now()
	ticket.Status = "open"
	return s.repo.Create(ticket)
}

func (s *service) FindTicketByID(id string) (*Ticket, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindAllTickets() ([]*Ticket, error) {
	return s.repo.FindAll()
}
