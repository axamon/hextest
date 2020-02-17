// Package ticket manages tickets in repository.
package ticket

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	ticket.Creator = getInfo("Creator")
	ticket.Description = getInfo("Description")
	ticket.Title = getInfo("Title")
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

func getInfo(s string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter " + s + ":")
	text, _ := reader.ReadString('\n')
	return text
}
