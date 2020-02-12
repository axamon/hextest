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

// Create creates a ticket.
func (r *ticketRepository) Create(ticket *ticket.Ticket) error {
	encoded, err := json.Marshal(ticket)

	if err != nil {
		return err
	}

	r.connection.HSet(table, ticket.ID, encoded) //Don't expire
	return nil
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
