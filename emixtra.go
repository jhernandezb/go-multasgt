package multasgt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &Emixtra{}

const (
	emixtraURL      = "http://consultas.munimixco.gob.gt/vista/emixtra.php"
	emixtraPhotoURL = "http://consultas.munimixco.gob.gt/vista/views/foto.php?rem=%v&T=%v&P=%v&s=F"
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

	b, _ := ioutil.ReadAll(resp.Body)
	// Nasty hack since malformed html with unclosed `b` tags generates a deep tree.
	s := strings.Replace(strings.Replace(string(b), "<b>", "", -1), "</b>", "", -1)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	status := 0
	ticket := Ticket{}
	var tickets []Ticket
	doc.Find("#foo").Children().Each(func(idx int, sel *goquery.Selection) {
		if sel.Is("hr") {
			status = 0
			return
		}
		if !sel.Is("div") {
			return
		}

		switch status {
		case 0:
			ticket.Entity = "EMIXTRA"
			ticket.Location = cleanStrings(sel.Children().Last().Text())
		case 1:
			ticket.Date = cleanStrings(sel.Children().Last().Text())
		case 3:
			var pID string
			sel.Find("h6").EachWithBreak(func(i int, s *goquery.Selection) bool {
				switch i {
				case 0:
					pID = cleanStrings(s.Text())
					ticket.Photo = fmt.Sprintf(emixtraPhotoURL, pID, plateType, plateNumber)
				case 1:
					ticket.Ammount = cleanStrings(s.Text())
					return false
				}
				return true
			})
		case 5:
			sel.Find("h6").EachWithBreak(func(i int, s *goquery.Selection) bool {
				switch i {
				case 0:
					ticket.Info = cleanStrings(s.Text())
					return false
				}
				return true
			})
			tickets = append(tickets, ticket)
			ticket = Ticket{}
		}
		status++
	})
	return tickets, nil
}
