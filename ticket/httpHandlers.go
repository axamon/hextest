package ticket

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// NewTicketHandler handles creations of new tickets via http.
func NewTicketHandler(ticketService Service) Handler {
	return &handler{
		ticketService,
	}
}

// Get handles requests of tickets via http.
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	tickets, _ := h.ticketService.FindAllTickets()

	response, _ := json.Marshal(tickets)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ticket, _ := h.ticketService.FindTicketByID(id)

	response, _ := json.Marshal(ticket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&ticket)
	_ = h.ticketService.CreateTicket(&ticket)

	response, _ := json.Marshal(ticket)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}
