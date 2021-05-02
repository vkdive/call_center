package main

import (
	"log"
	"net/http"
	"time"

	agents "call_center/agents"
	"call_center/handler"
	priorityqueue "call_center/queue"
)

func main() {
	pq := priorityqueue.NewHierarchicalQueue(1, false)
	agents := agents.NewAgents()

	ticketHandler := handler.NewTicketHandler(pq, agents)
	go startProcessingTickets(ticketHandler)//can be implemented using middleware
	http.HandleFunc("/ticket", ticketHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func startProcessingTickets(ticketHandler handler.TicketsHandler) {
	for {
		time.Sleep(2 * time.Second)
		ticketHandler.CheckAndAssignAgentsForNonPriorityTicketsIfFree()
		ticketHandler.CheckAndAssignAgentsForPriorityTicketsIfFree()
	}
}
