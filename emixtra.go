package multasgt

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &Emixtra{}

const (
	emixtraURL      = "http://consultas.munimixco.gob.gt/vista/emixtra.php"
	emixtraPhotoURL = "http://consultas.munimixco.gob.gt/vista/views/foto.php?rem=%v&T=%v&P=%v&s=%v"
	exmitraEntity   = "EMIXTRA"
)

// Emixtra implementation.
type Emixtra struct {
	Client *http.Client
}

// Check retrieves all tickes and aditional information.
func (e *Emixtra) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := e.Client.PostForm(emixtraURL, url.Values{
		"tPlaca": {plateType},
		"placa":  {plateNumber},
		"estado": {"0"},
	})

	if err != nil {
		return nil, err
	}
	var tickets []Ticket
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	rows := doc.Find(".panel-body .col-xs-12 > .panel > .panel-body > .row")
	currentIndex := 0
	var ticket Ticket
	rows.Each(func(idx int, sel *goquery.Selection) {
		if idx == 0 {
			return
		}
		currentIndex++
		switch currentIndex {
		case 1:
			if idx+1 == rows.Length() {
				return
			}
			ticket = Ticket{Entity: exmitraEntity}

			sel.Find(".col-xs-2").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					ticket.ID = cleanStrings(s.Text())
				case 1:
					ticket.Date = cleanStrings(s.Text())
				case 2:
					ticket.Location = cleanStrings(s.Text())
				case 3:
					ticket.Ammount = cleanStrings(s.Text())
				case 4:
					ticket.Discount = cleanStrings(s.Text())
				case 5:
					ticket.Total = cleanStrings(s.Text())
				}
			})
		case 2:
			form := sel.Find("form")
			formS := form.Find(`[name="s"]`).AttrOr("value", "F")
			formRem := form.Find(`[name="rem"]`).AttrOr("value", "")
			formT := form.Find(`[name="T"]`).AttrOr("value", "P")
			formP := form.Find(`[name="P"]`).AttrOr("value", "")
			ticket.Photo = fmt.Sprintf(emixtraPhotoURL, formRem, formT, formP, formS)
		case 3:
			ticket.Info = cleanStrings(sel.Find(".col-xs-7 .row:nth-child(2) > .col-xs-6").Text())
		case 4:
			tickets = append(tickets, ticket)
			currentIndex = 0
		}
	})
	return tickets, nil
}
