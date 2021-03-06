package multasgt

import (
	"net/http"
	"net/url"

	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &VillaNueva{}

const (
	villaNuevaURL      = "http://www.villanueva.gob.gt/consultas/consulta_pmt.php"
	villaNuevaPhotoURL = "http://www.villanueva.gob.gt/consultas/fotos/%v%v%v%v.JPG"
	villaNuevaEntity   = "VILLANUEVA"
)

// VillaNueva implementation.
type VillaNueva struct {
	Client *http.Client
}

// Check retrieves all tickes and aditional information.
func (f *VillaNueva) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := f.Client.PostForm(villaNuevaURL, url.Values{
		"tplaca":        {plateType},
		"nplaca":        {plateNumber},
		"op":            {"Consultar"},
		"page_num":      {"1"},
		"page_count":    {"page_count"},
		"form_build_id": {"form-A3iFsUaS2b03xqxM6nVJGkuZYva6dcrT_apE2x90lWg"},
		"form_id":       {"webform_client_form_66"},
	})

	if err != nil {
		return nil, err
	}

	var tickets []Ticket
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	rows := doc.Find("table.consulta.remision-info tbody tr")
	rl := rows.Length()
	rows.EachWithBreak(func(idx int, s *goquery.Selection) bool {
		if idx < 3 {
			return true
		}

		if rl < 3 || idx >= rl-1 {
			return false
		}
		ticket := Ticket{Entity: villaNuevaEntity}
		ticket.Info = ""
		var serie string

		s.Find("td").Each(func(i int, sel *goquery.Selection) {
			switch i {
			case 1:
				serie = CleanStrings(sel.Text())
			case 2:
				ticket.ID = CleanStrings(sel.Text())
			case 3:
				ticket.Date = CleanStrings(sel.Text())
			case 4:
				ticket.Location = CleanStrings(sel.Text())
			case 5:
				ticket.Info = CleanStrings(sel.Text())
			case 6:
				ticket.Amount = CleanStrings(sel.Text())
			case 7:
				ticket.Discount = CleanStrings(sel.Text())
			case 8:
				ticket.Total = CleanStrings(sel.Text())
			}
		})
		ticket.Photo = fmt.Sprintf(villaNuevaPhotoURL, plateType, plateNumber, serie, ticket.ID)
		tickets = append(tickets, ticket)
		return true
	})
	// TODO: Visit http://www.villanueva.gob.gt/consultas/foto.php?serie=T&boleta=%s&tipo=%s&numero=%s for each ticket in order to generate images.
	return tickets, nil
}
