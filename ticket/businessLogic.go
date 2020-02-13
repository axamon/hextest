// Package ticket creates tickets in repository.
package ticket

import (
	"errors"
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
func (s *service) CreateTicket(ticket *Ticket) (string, error) {
	ticket.ID = uuid.New().String()
	ticket.Created = time.Now()
	ticket.Updated = time.Now()
	ticket.Status = "open"
	if s.repo.Create(ticket) != nil {
		return "", errors.New("ticket creation impossible")
	}
	return ticket.ID, nil
}

// DeleteTicketByID method deletes ticket with id passed as argument
// from the repository.
func (s *service) DeleteTicketByID(id string) error {
	return s.repo.DeleteByID(id)
}

// FindTicketByID method returns ticket with id passed as argument
// from the repository.
func (s *service) FindTicketByID(id string) (*Ticket, error) {
	return s.repo.FindByID(id)
}

// FindAllTickets method returns all tickets in the repository.
func (s *service) FindAllTickets() ([]*Ticket, error) {
	return s.repo.FindAll()
}
