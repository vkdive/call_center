package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"call_center/agents"
	priorityqueue "call_center/queue"
	"call_center/ticket"
)

// create a handler struct
type TicketsHandler struct {
	pq     *priorityqueue.HierarchicalQueue
	agents []*agents.Agents
}

func NewTicketHandler(queue *priorityqueue.HierarchicalQueue, agents []*agents.Agents) TicketsHandler {
	return TicketsHandler{
		queue, agents,
	}
}

// implement `ServeHTTP` method on `HttpHandler` struct
func (h TicketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ticketDetails ticket.Ticket
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		json.Unmarshal(body, &ticketDetails)

		if ticketDetails.Platform == "" {
			http.Error(w, "Invalid request method", http.StatusNotAcceptable)
			return
		}
	}

	h.queueTickets(ticketDetails)

	data := []byte("you requested has been taken, we are trying to connect to agent now...tring...tring...tring...tring...")
	w.Write(data)

}

func (h TicketsHandler) queueTickets(detail ticket.Ticket) {
	if detail.Platform == "call" {
		h.pq.Enqueue(detail, 0)
		return
	} else {
		h.pq.Enqueue(detail, 1)
		return
	}
}

func (h TicketsHandler) CheckAndAssignAgentsForPriorityTicketsIfFree() {
	for h.pq.LenPriority(0) > 0 {
		agent := agents.GetFreePriorityAgents(h.agents)
		if agent != nil {
			ticket, _ := h.pq.GetHighestPriorityTicketsFromQueue(0)
			agent.NumberOfTickets++
			go agent.HandleTickets(ticket)
			continue
		}
		fmt.Println("no agents are free at this time will solve the issue once agents are ready")
		return
	}
}

func (h TicketsHandler) CheckAndAssignAgentsForNonPriorityTicketsIfFree() {
	for h.pq.LenPriority(1) > 0 {
		agent := agents.GetFreeNonPriorityAgents(h.agents)
		if agent != nil {
			ticket, _ := h.pq.GetHighestPriorityTicketsFromQueue(1)
			agent.NumberOfTickets++
			go agent.HandleTickets(ticket)
			continue
		}
		fmt.Println("no agents are free at this time will solve the issue once agents are ready")
		return
	}
}
