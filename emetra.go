package multasgt

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &Emetra{}

const (
	emetraURL = "http://consulta.muniguate.com/emetra/despliega.php"
)

// Emetra implementation.
type Emetra struct {
	Client *http.Client
}

func (e *Emetra) processEmetraTicket(wg *sync.WaitGroup, tickets []Ticket, idx int, sel *goquery.Selection) {
	defer wg.Done()
	info := sel.Find(".texto")
	if len(info.Nodes) < 3 {
		return
	}
	photoURL := "http://consultas.muniguate.com/consultas/remisiones/remision_pic.jsp?r=%v"

	rURL, err := url.Parse(GetAttribute("href", info.Get(1)))
	if err != nil {
		return
	}

	date := info.Get(0).FirstChild.Data
	loc := info.Get(1).FirstChild.Data
	ammount := info.Get(2).FirstChild.Data
	id := rURL.Query().Get("r")
	parsedPhotoURL := fmt.Sprintf(photoURL, id)
	var photo string
	tickets[idx] = Ticket{
		Entity:   "EMETRA",
		ID:       id,
		Date:     date,
		Amount:   ammount,
		Photo:    photo,
		Location: loc,
	}

	res, err := e.Client.Get(parsedPhotoURL)
	if err != nil {
		return
	}

	photoDoc, err := goquery.NewDocumentFromResponse(res)
	m := photoDoc.Find(`strong:contains("MOTIVO DE REMISION:")`).Parent().Text()
	tickets[idx].Info = CleanStrings(strings.Replace(m, "MOTIVO DE REMISION:", "", -1))
	regex := `[src#=(^http:\/\/consultas.muniguate.com\/consultas\/fotos)]`
	// URL comes with an extra \n.
	tickets[idx].Photo = CleanStrings(GetAttribute("src", photoDoc.Find(regex).Get(0)))

}

// Check retrieves all tickes and aditional information.
func (e *Emetra) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := e.Client.PostForm(emetraURL, url.Values{
		"tplaca": {plateType},
		"nplaca": {plateNumber},
	})

	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)

	var wg sync.WaitGroup
	if err != nil {
		return nil, err
	}
	rows := doc.Find("tr.row")
	tickets := make([]Ticket, len(rows.Nodes))
	rows.Each(func(idx int, sel *goquery.Selection) {
		wg.Add(1)
		e.processEmetraTicket(&wg, tickets, idx, sel)
	})
	wg.Wait()
	return tickets, nil
}
