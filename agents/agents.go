package agents

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Agents struct {
	Name                string
	PriorityType        int
	TicketsHandlingType string
	Languages           []string
	NumberOfTickets     int64
	TicketDetails       []interface{}
}

func NewAgents() []*Agents {
	return getAgents()
}
func getAgents() []*Agents {
	var agents []*Agents
	for i := 1; i <= 5; i++ {
		ticketHandlingType := ""
		priorityType:= 1
		if i%2 == 0 {
			ticketHandlingType = "call"
			priorityType =0
		}
		agents = append(agents, &Agents{
			Name:                fmt.Sprint("AGENT", i),
			PriorityType:        priorityType,
			TicketsHandlingType: ticketHandlingType,
			Languages:           nil,
			NumberOfTickets:     0,
			TicketDetails:       []interface{}{},
		})
	}
	return agents
}

func (a *Agents) HandleTickets(ticket interface{}) {
	atomic.AddInt64(&a.NumberOfTickets, 1)
	fmt.Println(a.Name, "is handling ticket now", ticket)
	time.Sleep(10 * time.Minute)
	fmt.Println("ticket handling completed")
	a.NumberOfTickets--
	return
}

func GetFreePriorityAgents(agents []*Agents) *Agents {
	for _,agent := range agents{
		if agent.TicketsHandlingType=="call" && agent.NumberOfTickets <=3 {
			return agent
		}
	}
	return nil
}

func GetFreeNonPriorityAgents(agents []*Agents) *Agents {
	for _,agent := range agents{
		if agent.TicketsHandlingType!="call" && agent.NumberOfTickets <=4 {
			return agent
		}
	}
	return nil
}
