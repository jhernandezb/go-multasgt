package multasgt

import (
	"net/http"
	"sync"
)

// Tickets represents all collected tickets
type Tickets struct {
	Tickets []Ticket
	M       sync.Mutex
}

func getTickets(wg *sync.WaitGroup,
	plateType, plateNumber string,
	tckts *Tickets,
	checker TicketChecker) {
	defer wg.Done()
	var tickets []Ticket
	tickets, _ = checker.Check(plateType, plateNumber)
	tckts.M.Lock()
	tckts.Tickets = append(tckts.Tickets, tickets...)
	tckts.M.Unlock()
}

// GetAllTickets retrieves tickets from all entities.
func GetAllTickets(plateType, plateNumber string, c *http.Client) ([]Ticket, error) {
	var wg sync.WaitGroup
	wg.Add(7)
	var tickets Tickets
	// Emetra
	go getTickets(&wg, plateType, plateNumber, &tickets, &Emetra{Client: c})
	// Mixco
	go getTickets(&wg, plateType, plateNumber, &tickets, &Emixtra{Client: c})
	// SCP
	go getTickets(&wg, plateType, plateNumber, &tickets, &SCP{Client: c})
	// Fraijanes
	go getTickets(&wg, plateType, plateNumber, &tickets, &Fraijanes{Client: c})
	// Antigua
	go getTickets(&wg, plateType, plateNumber, &tickets, &Antigua{Client: c})
	// PNC
	go getTickets(&wg, plateType, plateNumber, &tickets, &PNC{Client: c})
	// VillaNueva
	go getTickets(&wg, plateType, plateNumber, &tickets, &VillaNueva{Client: c})
	wg.Wait()
	return tickets.Tickets, nil
}
