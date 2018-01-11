package emixtra_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/jhernandezb/go-multasgt/entities/emixtra"
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
	doc, err := emixtra.Fetch("P", "308FZS", http.DefaultClient)
	if err != nil {
		t.Fatal("Should not return an error")
	}
	rows := doc.Find(`img[src*="http://consultas\.munimixco\.gob\.gt/img/view/"]`)
	if len(rows.Nodes) != 21 {
		t.Fatal("Expected 21 tickets")
	}
}

func TestFetchInvalidStatus(t *testing.T) {
	defer gock.Off()
	gock.New("https://consultas.munimixco.gob.gt").
		Post("/pvisa/emixtra").
		Reply(500)
	_, err := emixtra.Fetch("P", "308FZS", http.DefaultClient)
	if err == nil {
		t.Fatal("Should return an error")
	}
}

func TestFetchError(t *testing.T) {
	defer gock.Off()
	gock.New("https://consultas.munimixco.gob.gt").
		Post("/pvisa/emixtra").
		ReplyError(http.ErrServerClosed)
	_, err := emixtra.Fetch("P", "308FZS", http.DefaultClient)
	if err == nil {
		t.Fatal("Should return an error")
	}
}
