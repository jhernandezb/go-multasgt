package multasgt

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Emetra implementation.
type Emetra struct {
}

func processEmetraTicket(wg *sync.WaitGroup, tickets []Ticket, idx int, sel *goquery.Selection) {
	defer wg.Done()
	info := sel.Find(".texto")
	if len(info.Nodes) < 3 {
		return
	}
	photoURL := "http://consultas.muniguate.com/consultas/remisiones/remision_pic.jsp?r=%v"

	rURL, err := url.Parse(getAttribute("href", info.Get(1)))
	if err != nil {
		return
	}

	date := info.Get(0).FirstChild.Data
	loc := info.Get(1).FirstChild.Data
	ammount := info.Get(2).FirstChild.Data
	parsedPhotoURL := fmt.Sprintf(photoURL, rURL.Query().Get("r"))
	var photo string
	tickets[idx] = Ticket{Date: date,
		Ammount:  ammount,
		Photo:    photo,
		Location: loc,
	}
	if err != nil {
		return
	}
	photoDoc, err := goquery.NewDocument(parsedPhotoURL)
	m := photoDoc.Find(`strong:contains("MOTIVO DE REMISION:")`).Parent().Text()
	tickets[idx].Info = cleanStrings(strings.Replace(m, "MOTIVO DE REMISION:", "", -1))
	regex := `[src#=(^http:\/\/consultas.muniguate.com\/consultas\/fotos)]`
	// URL comes with an extra \n.
	tickets[idx].Photo = strings.Replace(getAttribute("src", photoDoc.Find(regex).Get(0)), "\n", "", -1)

}

// Check retrieves all tickes and aditional information.
func (e Emetra) Check(plateType, plateNumber string) ([]Ticket, error) {
	resp, err := http.PostForm(emetraURL, url.Values{
		"tplaca": {plateType},
		"nplaca": {plateNumber},
	})

	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)

	var wg sync.WaitGroup

	rows := doc.Find("tr.row")
	tickets := make([]Ticket, len(rows.Nodes))
	rows.Each(func(idx int, sel *goquery.Selection) {
		wg.Add(1)
		processEmetraTicket(&wg, tickets, idx, sel)
	})
	wg.Wait()
	return tickets, nil
}
