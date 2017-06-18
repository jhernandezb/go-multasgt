package multasgt

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &Emixtra{}

const (
	emixtraURL    = "https://consultas.munimixco.gob.gt/pvisa/emixtra"
	exmitraEntity = "EMIXTRA"
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
	rows := doc.Find(".row > .col-md-8.col-xs-10> .panel.panel-primary")
	var ticket Ticket
	rows.Each(func(idx int, sel *goquery.Selection) {
		// Ommit header info
		if idx == 0 {
			return
		}
		ticket = Ticket{Entity: exmitraEntity}
		ticket.ID = cleanStrings(sel.Find(".panel-heading > .row > .col-md-5.col-xs-5").Text())
		ticket.Date = cleanStrings(sel.Find(".panel-heading > .row > .col-md-6.col-xs-5").Text())
		ticket.Location = cleanStrings(sel.Find(".panel-body > .row:nth-of-type(1) > .col-md-3.col-xs-9").Text())
		ticket.Info = cleanStrings(sel.Find(".panel-body > .row:nth-of-type(3) > .col-md-5.col-xs-9").Text())
		ticket.Ammount = cleanStrings(sel.Find(".panel-body > .row:nth-of-type(6) > div:nth-of-type(1)").Text())
		ticket.Discount = cleanStrings(sel.Find(".panel-body > .row:nth-of-type(6) > div:nth-of-type(2)").Text())
		ticket.Total = cleanStrings(sel.Find(".panel-body > .row:nth-of-type(6) > div:nth-of-type(3)").Text())
		_photo := sel.Find(".panel-body > .row:nth-of-type(4) > .col-md-5 >  img")
		ticket.Photo = getAttribute("src", _photo.Get(0))
		tickets = append(tickets, ticket)
	})
	return tickets, nil
}
