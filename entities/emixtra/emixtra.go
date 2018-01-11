package emixtra

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	multasgt "github.com/jhernandezb/go-multasgt"
)

const (
	emixtraURL    = "https://consultas.munimixco.gob.gt/pvisa/emixtra"
	exmitraEntity = "EMIXTRA"
)

// Fetch fetches the ticket from the remote endpoint.
func Fetch(plateType, plateNumber string, cli *http.Client) (*goquery.Document, error) {
	resp, err := cli.PostForm(emixtraURL, url.Values{
		"tPlaca": {plateType},
		"placa":  {plateNumber},
		"estado": {"0"},
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request not successful: status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Parse parses the retrieved document and returns and array of tickets.
func Parse(doc *goquery.Document) ([]multasgt.Ticket, error) {
	rows := doc.Find(".row > .col-md-8.col-xs-10> .panel.panel-primary")
	var tickets []multasgt.Ticket
	var ticket multasgt.Ticket
	rows.Each(func(idx int, sel *goquery.Selection) {
		// Ommit header info
		if idx == 0 {
			return
		}
		ticket = multasgt.Ticket{Entity: exmitraEntity}
		ticket.ID = multasgt.CleanStrings(sel.Find(".panel-heading > .row > .col-md-5.col-xs-5").Text())
		ticket.Date = multasgt.CleanStrings(sel.Find(".panel-heading > .row > .col-md-6.col-xs-5").Text())
		ticket.Location = multasgt.CleanStrings(sel.Find(".panel-body > .row:nth-of-type(1) > .col-md-3.col-xs-9").Text())
		ticket.Info = multasgt.CleanStrings(sel.Find(".panel-body > .row:nth-of-type(3) > .col-md-5.col-xs-9").Text())
		ticket.Ammount = multasgt.CleanStrings(sel.Find(".panel-body > .row:nth-of-type(6) > div:nth-of-type(1)").Text())
		ticket.Discount = multasgt.CleanStrings(sel.Find(".panel-body > .row:nth-of-type(6) > div:nth-of-type(2)").Text())
		ticket.Total = multasgt.CleanStrings(sel.Find(".panel-body > .row:nth-of-type(6) > div:nth-of-type(3)").Text())
		_photo := sel.Find(".panel-body > .row:nth-of-type(4) > .col-md-5 >  img")
		ticket.Photo = multasgt.GetAttribute("src", _photo.Get(0))
		tickets = append(tickets, ticket)
	})
	return tickets, nil
}
