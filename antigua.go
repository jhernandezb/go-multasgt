package multasgt

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &Antigua{}

const (
	antiguaURL    = "http://muniantigua.com/sql/pmt/buscar.php"
	antiguaEntity = "ANTIGUA"
)

// Antigua implementation.
type Antigua struct {
	Client *http.Client
}

// Check retrieves all tickes and aditional information.
func (a *Antigua) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := a.Client.PostForm(antiguaURL, url.Values{
		"placa_total": {fmt.Sprintf("%s-%s", plateType, plateNumber)},
	})

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	var tickets []Ticket

	doc.Find(`.main table:last-child tbody tr`).Each(func(idx int, sel *goquery.Selection) {
		if idx == 0 {
			return
		}
		if sel.Children().Length() == 0 {
			return
		}
		ticket := Ticket{Entity: antiguaEntity}
		sel.Children().Each(func(cIdx int, cSel *goquery.Selection) {
			switch cIdx {
			case 0:
				ticket.ID = cleanStrings(cSel.First().Text())
			case 1:
				ticket.Date = cleanStrings(cSel.First().Text())
			case 2:
				ticket.Location = cleanStrings(cSel.Text())
			case 3:
				ticket.Info = cleanStrings(cSel.Text())
			case 5:
				ticket.Ammount = cleanStrings(cSel.Text())
			case 7:
				ticket.Total = cleanStrings(cSel.Text())
				tickets = append(tickets, ticket)
			}
		})

	})

	return tickets, nil
}
