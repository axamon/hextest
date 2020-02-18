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

// Repository is thte interface to contact to interact with the core.
type Repository interface {
	Create(ticket *Ticket) error
	DeleteByID(id string) error
	FindByID(id string) (*Ticket, error)
	FindAll() ([]*Ticket, error)
	CloseByID(id string) (*Ticket, error)
}

// Service is the interface to the business rules logic.
type Service interface {
	CreateTicket(ticket *Ticket) (string, error)
	DeleteTicketByID(id string) error
	FindTicketByID(id string) (*Ticket, error)
	CloseTicketByID(id string) (*Ticket, error)
	FindAllTickets() ([]*Ticket, error)
}

type service struct {
	repo Repository
}

// Handler interface is the way to deal with the repository
// via http routes.
type Handler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	CloseByID(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ticketService Service
}
