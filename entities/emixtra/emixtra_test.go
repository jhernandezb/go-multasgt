package emixtra

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/h2non/gock"
)

func emixtraMatcher(req *http.Request, ereq *gock.Request) (bool, error) {
	req.ParseForm()
	estado := req.Form.Get("estado")
	if estado != "0" {
		return false, errors.New("Should have `estado` prop")
	}
	return true, nil
}
func TestFetchAndCheck(t *testing.T) {
	defer gock.Off()
	gock.New("https://consultas.munimixco.gob.gt").
		Post("/pvisa/emixtra").
		AddMatcher(emixtraMatcher).
		Reply(200).
		File("../../test_data/emixtra.html")
	doc, err := fetch("P", "308FZS", http.DefaultClient)
	if err != nil {
		t.Fatal("Should not return an error")
	}
	if doc == nil {
		t.Fatal("Document should not be nil")
	}
}

func TestFetchInvalidStatus(t *testing.T) {
	defer gock.Off()
	gock.New("https://consultas.munimixco.gob.gt").
		Post("/pvisa/emixtra").
		Reply(500)
	_, err := fetch("P", "308FZS", http.DefaultClient)
	if err == nil {
		t.Fatal("Should return an error")
	}
}

func TestFetchError(t *testing.T) {
	defer gock.Off()
	gock.New("https://consultas.munimixco.gob.gt").
		Post("/pvisa/emixtra").
		ReplyError(http.ErrServerClosed)
	_, err := fetch("P", "308FZS", http.DefaultClient)
	if err == nil {
		t.Fatal("Should return an error")
	}
}

func TestParse(t *testing.T) {
	testFile, err := os.Open("../../test_data/emixtra.html")
	if err != nil {
		t.Fatalf("Error opening test data")
	}
	defer testFile.Close()

	doc, err := goquery.NewDocumentFromReader(testFile)
	if err != nil {
		t.Fatalf("Error generating document")
	}
	tickets, err := parse(doc)
	if err != nil {
		t.Fatalf("Error should be nil")
	}
	if len(tickets) != 21 {
		t.Fatalf("Should have 21 tickets")
	}
}
