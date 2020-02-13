package ticket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewTicketHandler handles creations of new tickets via http.
func NewTicketHandler(ticketService Service) Handler {
	return &handler{
		ticketService,
	}
}

// GetAll method returns all tickets via http.
func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	tickets, _ := h.ticketService.FindAllTickets()

	response, _ := json.Marshal(tickets)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// DeleteByID method deletes one ticket by id via http.
func (h *handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.ticketService.DeleteTicketByID(id)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Deleted ticket with id: " + id))
}

// GetByID method returns one ticket by id via http.
func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ticket, _ := h.ticketService.FindTicketByID(id)

	response, _ := json.Marshal(ticket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// Create method creates a new ticket in the repository.
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&ticket)
	id, _ := h.ticketService.CreateTicket(&ticket)

	response, _ := json.Marshal(ticket)
	w.Header().Set("Content-Type", "application/json")
	doneOk := fmt.Sprintf("Created ticket with id: %s\n", id)
	w.Write([]byte(doneOk))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
