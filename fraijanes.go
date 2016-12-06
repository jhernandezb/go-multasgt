package multasgt

import (
	"net/http"
	"net/url"

	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &Fraijanes{}

const (
	fraijanesURL    = "http://190.149.238.69/consultapmt/daemonpt.php"
	fraijanesEntity = "FRAIJANES"
)

// Fraijanes implementation.
type Fraijanes struct {
	Client *http.Client
}

// Check retrieves all tickes and aditional information.
func (f *Fraijanes) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := f.Client.PostForm(fraijanesURL, url.Values{
		"cbxtipo":    {plateType},
		"txtidplaca": {plateNumber},
		"cmdsmpmt":   {"Consultar"},
	})

	if err != nil {
		return nil, err
	}
	var tickets []Ticket
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	rows := doc.Find("table tbody tr")
	rl := rows.Length()
	rows.EachWithBreak(func(idx int, s *goquery.Selection) bool {
		if idx == 0 {
			return true
		}

		if rl < 3 || idx >= rl-2 {
			return false
		}
		ticket := Ticket{Entity: fraijanesEntity}
		ticket.Info = ""

		s.Find("td").Each(func(i int, sel *goquery.Selection) {
			switch i {
			case 0:
				ticket.ID = cleanStrings(sel.Text())
			case 2:
				ticket.Info = cleanStrings(sel.Text())
			case 3:
				ticket.Info = fmt.Sprintf("%v, %v", ticket.Info, cleanStrings(sel.Text()))
			case 4:
				val := cleanStrings(sel.Text())
				ticket.Ammount = val
				ticket.Total = val
			}
		})
		tickets = append(tickets, ticket)
		return true
	})
	return tickets, nil
}
