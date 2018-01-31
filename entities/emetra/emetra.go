package emetra

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	multasgt "github.com/jhernandezb/go-multasgt"
)

const (
	emetraURL = "http://consulta.muniguate.com/emetra/despliega.php"
	photoURL  = "http://consultas.muniguate.com/consultas/remisiones/remision_pic.jsp?r=%v"
	detailURL = "http://consultas.muniguate.com/consultas/remisiones/detalle.jsp?s=%v&r=%v&id=%v"
	entity    = "EMETRA"
)

func fetch(plateType, plateNumber string, cli *http.Client) (*goquery.Document, error) {
	resp, err := cli.PostForm(emetraURL, url.Values{
		"tplaca": {plateType},
		"nplaca": {plateNumber},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request not successful: status %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func parse(doc *goquery.Document, cli *http.Client) ([]multasgt.Ticket, error) {
	var wg sync.WaitGroup
	rows := doc.Find("tr.row")
	tickets := make([]multasgt.Ticket, len(rows.Nodes))
	rows.Each(func(idx int, sel *goquery.Selection) {
		wg.Add(1)
		go processEmetraTicket(&wg, tickets, idx, sel, cli)
	})
	wg.Wait()
	return tickets, nil
}

// Check fetch and parse a request.
func Check(plateType, plateNumber string, cli *http.Client) ([]multasgt.Ticket, error) {
	doc, err := fetch(plateType, plateNumber, cli)
	if err != nil {
		return nil, err
	}
	return parse(doc, cli)
}

func processEmetraTicket(wg *sync.WaitGroup, tickets []multasgt.Ticket, idx int, sel *goquery.Selection, cli *http.Client) {
	defer wg.Done()
	info := sel.Find(".texto")
	if len(info.Nodes) < 3 {
		return
	}
	rURL, err := url.Parse(multasgt.GetAttribute("href", info.Get(1)))
	if err != nil {
		return
	}
	date := info.Get(0).FirstChild.Data
	loc := info.Get(1).FirstChild.Data
	ammount := info.Get(2).FirstChild.Data
	id := rURL.Query().Get("r")
	rtype := rURL.Query().Get("id")
	serie := rURL.Query().Get("s")
	tickets[idx] = multasgt.Ticket{
		Entity:   entity,
		ID:       id,
		Date:     date,
		Amount:   ammount,
		Location: loc,
	}
	if id == "" {
		return
	}
	parsedDetailURL := fmt.Sprintf(detailURL, serie, id, rtype)
	res, err := cli.Get(parsedDetailURL)
	if err != nil {
		return
	}
	detailsDoc, err := goquery.NewDocumentFromResponse(res)
	info = detailsDoc.Find(`.head_blue2:contains("Descripci")`).Parent().Find(".texto")
	tickets[idx].Info = multasgt.CleanStrings(info.Text())
	// only series L have pictures
	if multasgt.CleanStrings(serie) != "L" {
		return
	}
	parsedPhotoURL := fmt.Sprintf(photoURL, id)
	res, err = cli.Get(parsedPhotoURL)
	if err != nil {
		return
	}
	photoDoc, err := goquery.NewDocumentFromResponse(res)
	m := photoDoc.Find(`strong:contains("MOTIVO DE REMISION:")`).Parent().Text()
	tickets[idx].Info = multasgt.CleanStrings(strings.Replace(m, "MOTIVO DE REMISION:", "", -1))
	regex := `[src#=(^http:\/\/consultas.muniguate.com\/consultas\/fotos)]`
	// URL comes with an extra \n.
	tickets[idx].Photo = multasgt.CleanStrings(multasgt.GetAttribute("src", photoDoc.Find(regex).Get(0)))
}
