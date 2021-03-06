package redis

import (
	"encoding/json"

	"github.com/axamon/hextest/ticket"
	"github.com/go-redis/redis"
)

const table = "tickets"

type ticketRepository struct {
	connection *redis.Client
}

// NewRedisTicketRepository creates a repository on Redis.
func NewRedisTicketRepository(connection *redis.Client) ticket.Repository {
	return &ticketRepository{
		connection,
	}
}

// Create creates a ticket in the repository.
func (r *ticketRepository) Create(ticket *ticket.Ticket) error {
	encoded, err := json.Marshal(ticket)

	if err != nil {
		return err
	}

	r.connection.HSet(table, ticket.ID, encoded) // Does not expire
	return nil
}

// DeleteByID method returns the ticket with id passed as argument.
func (r *ticketRepository) DeleteByID(id string) error {
	_ = r.connection.HDel(table, id)

	return nil
}

// CloseByID method changes the status of tt to cloesed.
func (r *ticketRepository) CloseByID(id string) (*ticket.Ticket, error) {
	b, err := r.connection.HGet(table, id).Bytes()

	if err != nil {
		return nil, err
	}

	t := new(ticket.Ticket)
	err = json.Unmarshal(b, t)

	if err != nil {
		return nil, err
	}

	// cancel old tt version.
	_ = r.connection.HDel(table, id)

	t.Status = "closed"

	encoded, err := json.Marshal(t)

	_ = r.connection.HSet(table, id, encoded) // Does not expire

	return t, nil
}

// FindByID method returns the ticket with id passed as argument.
func (r *ticketRepository) FindByID(id string) (*ticket.Ticket, error) {
	b, err := r.connection.HGet(table, id).Bytes()

	if err != nil {
		return nil, err
	}

	t := new(ticket.Ticket)
	err = json.Unmarshal(b, t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

// FindAll method returns all tickets present in the redis repository.
func (r *ticketRepository) FindAll() (tickets []*ticket.Ticket, err error) {
	ts := r.connection.HGetAll(table).Val()
	for key, value := range ts {
		t := new(ticket.Ticket)
		err = json.Unmarshal([]byte(value), t)

		if err != nil {
			return nil, err
		}

		t.ID = key
		tickets = append(tickets, t)
	}
	return tickets, nil
}
