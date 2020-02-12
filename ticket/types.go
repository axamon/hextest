package ticket

import (
	"net/http"
	"time"
)

// Ticket is the struct for tickets.
type Ticket struct {
	ID          string    `json:"id" db:"id"`
	Creator     string    `json:"creator" db:"creator"`
	Assigned    string    `json:"assigned" db:"assigned"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	Points      int       `json:"points" db:"points"`
	Created     time.Time `json:"created" db:"created"`
	Updated     time.Time `json:"updated" db:"updated"`
	Deleted     time.Time `json:"deleted" db:"deleted"`
}

// Repository is thte inferface to contact to interact with the core.
type Repository interface {
	Create(ticket *Ticket) error
	FindById(id string) (*Ticket, error)
	FindAll() ([]*Ticket, error)
}

type Service interface {
	CreateTicket(ticket *Ticket) error
	FindTicketById(id string) (*Ticket, error)
	FindAllTickets() ([]*Ticket, error)
}

type service struct {
	repo Repository
}

type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ticketService Service
}
