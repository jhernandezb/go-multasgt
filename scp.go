package multasgt

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &SCP{}

const (
	scpURL = "http://ws.scp.gob.gt/transito/"
)

// SCP implementation.
type SCP struct {
}

// Check retrieves all tickes and aditional information.
func (e SCP) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := http.PostForm(scpURL, url.Values{
		"cbTipoPlaca":          {plateType},
		"txtNumeroPlaca":       {plateNumber},
		"btnObtener":           {"Consultar"},
		"__VIEWSTATEGENERATOR": {""},
		"__VIEWSTATE":          {""},
	})

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)

	var tickets []Ticket
	doc.Find(`#pMultas table tbody tr`).Each(func(idx int, sel *goquery.Selection) {
		if idx%2 == 0 {
			return
		}
		ticket := Ticket{Entity: "SCP"}
		sel.Children().EachWithBreak(func(cIdx int, cSel *goquery.Selection) bool {
			if cIdx > 3 {
				return false
			}
			switch cIdx {
			case 0:
				ticket.Date = cleanStrings(cSel.First().Text())
			case 1:
				// TODO: https://golang.org/src/image/decode_example_test.go
				ticket.Photo = cSel.Find("img").AttrOr("src", "")
			case 2:
				// id := cSel.Find("#ctl03_lbBoleta").Text()
				ticket.Info = cleanStrings(cSel.Find(`[id*="lbInfraccion"]`).Text())
				ticket.Location = cleanStrings(cSel.Find(`[id*="_lbLugar"]`).Text())
			case 3:
				ticket.Ammount = cSel.Find(`[id*="lbCostoMulta"]`).Text()
				tickets = append(tickets, ticket)
			}

			return true
		})

	})

	return tickets, nil
}
