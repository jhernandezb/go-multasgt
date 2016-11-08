package multasgt

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Ensure interface implementation.
var _ TicketChecker = &PNC{}

const (
	pncURL    = "http://sistemas.transito.gob.gt/consultaweb/ConsultaRemisiones.aspx"
	pncEntity = "PNC"
)

// PNC implementation.
type PNC struct {
	Client *http.Client
}

// Check retrieves all tickes and aditional information.
func (p *PNC) Check(plateType, plateNumber string) ([]Ticket, error) {
	data := url.Values{
		"ctl00$MainContent$RadScriptManager1":             {"ctl00$MainContent$ctl00$MainContent$btnConsultarPanel|ctl00$MainContent$btnConsultar"},
		"ctl00_MainContent_RadScriptManager1_HiddenField": {";;System.Web.Extensions, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35:es-ES:3de828f0-5e0d-4c7d-a36b-56a9773c0def:ea597d4b:b25378d2;Telerik.Web.UI, Version=2009.2.826.35, Culture=neutral, PublicKeyToken=121fae78165ba3d4:es-ES:d2d891f5-3533-469c-b9a2-ac7d16eb23ff:16e4e7cd:ed16cbdc:58366029;"},
		"ctl00$MainContent$cmbTipoPlaca":                  {"1"},
		"ctl00$MainContent$TxtNoPlaca":                    {plateNumber},
		"__EVENTTARGET":                                   {"ctl00$MainContent$btnConsultar"},
		"RadAJAXControlID":                                {"ctl00_MainContent_RadAjaxManager1"},
		"ctl00_MainContent_grdResul_ClientState":          {""},
		"__EVENTVALIDATION":                               {"/wEWEQKyxp+rBAKY5fiKCQKgrb3ADQKhrb3ADQKirb3ADQKjrb3ADQKkrb3ADQKlrb3ADQKmrb3ADQK3rb3ADQK4rb3ADQKgrf3DDQKgrfHDDQKgrfXDDQKgrcnDDQKlrfHDDQKY3JWZD4e0MUr+3flZhXzEyTqBoCAxFmpw"},
		"__EVENTARGUMENT":                                 {""},
		"__VIEWSTATE":                                     {"/wEPDwUKLTY0ODI1MjQzNA9kFgJmD2QWAgIDD2QWAgIBD2QWCAIFDw8WAh4XRW5hYmxlQWpheFNraW5SZW5kZXJpbmdoZGQCCQ8QDxYGHg1EYXRhVGV4dEZpZWxkBQlUSVBQTEFJTkkeDkRhdGFWYWx1ZUZpZWxkBQlUSVBQTEFDT0QeC18hRGF0YUJvdW5kZ2QQFQ4BUAFDAUEBTQJDRAJDQwNESVMCVEMBTwFVAk1JA1NUUANUUkMDUE5DFQ4BMQEyATMBNAE1ATYBNwE4ATkCMTACMTECMTICMTMCNjEUKwMOZ2dnZ2dnZ2dnZ2dnZ2dkZAINDw8WAh4RVXNlU3VibWl0QmVoYXZpb3JoZGQCDw88KwANAgAUKwACDxYIHgtfIUl0ZW1Db3VudAIBHwNnHwBoHgtFZGl0SW5kZXhlcxYAZBcBBQ9TZWxlY3RlZEluZGV4ZXMWAAEWAhYKDwIHFCsABxQrAAUWBB4IRGF0YVR5cGUZKwIeBG9pbmQCAmRkZAUJT1JHUkVNREVTFCsABRYEHwcZKwIfCAIDZGRkBQVGRUNIQRQrAAUWBB8HGSsCHwgCBGRkZAUJTlVNU0VSQk9MFCsABRYEHwcZKwIfCAIFZGRkBQZOVU1CT0wUKwAFFgQfBxkrAh8IAgZkZGQFBVBMQUNBFCsABRYEHwcZKwIfCAIHZGRkBQZMVUdJTkYUKwAFFgQfBxkpWlN5c3RlbS5Eb3VibGUsIG1zY29ybGliLCBWZXJzaW9uPTIuMC4wLjAsIEN1bHR1cmU9bmV1dHJhbCwgUHVibGljS2V5VG9rZW49Yjc3YTVjNTYxOTM0ZTA4OR8IAghkZGQFBlRPVFJFTWRlFCsAAAspeVRlbGVyaWsuV2ViLlVJLkdyaWRDaGlsZExvYWRNb2RlLCBUZWxlcmlrLldlYi5VSSwgVmVyc2lvbj0yMDA5LjIuODI2LjM1LCBDdWx0dXJlPW5ldXRyYWwsIFB1YmxpY0tleVRva2VuPTEyMWZhZTc4MTY1YmEzZDQBPCsABwALKXRUZWxlcmlrLldlYi5VSS5HcmlkRWRpdE1vZGUsIFRlbGVyaWsuV2ViLlVJLCBWZXJzaW9uPTIwMDkuMi44MjYuMzUsIEN1bHR1cmU9bmV1dHJhbCwgUHVibGljS2V5VG9rZW49MTIxZmFlNzgxNjViYTNkNAEWAh4EX2Vmc2RkFgweCkRhdGFNZW1iZXJlHgRfaGxtCysFAR8DZx4IRGF0YUtleXMWAB4FXyFDSVMXAB8FAgFmFgRmDxQrAAMPZBYCHgVzdHlsZQULd2lkdGg6MTAwJTtkZGQCAQ8WBBQrAAIPFgwfCmUfCwsrBQEfA2cfDBYAHw0XAB8FAgFkFwMFBl8hRFNJQwIBBQtfIUl0ZW1Db3VudAIBBQhfIVBDb3VudGQWAh4DX3NlFgIeAl9jZmQWB2RkZGRkZGQWAmYPZBYIZg9kFgJmD2QWEmYPDxYEHgRUZXh0BQYmbmJzcDseB1Zpc2libGVoZGQCAQ8PFgQfEQUGJm5ic3A7HxJoZGQCAg8PFgIfEQUGT3JpZ2VuZGQCAw8PFgIfEQUFRmVjaGFkZAIEDw8WAh8RBQVTZXJpZWRkAgUPDxYCHxEFCk5vLiBCb2xldGFkZAIGDw8WAh8RBQVQbGFjYWRkAgcPDxYCHxEFBUx1Z2FyZGQCCA8PFgIfEQUFVmFsb3JkZAIBDw8WAh8SaGQWAmYPZBYSZg8PFgIfEQUGJm5ic3A7ZGQCAQ8PFgIfEQUGJm5ic3A7ZGQCAg8PFgIfEQUGJm5ic3A7ZGQCAw8PFgIfEQUGJm5ic3A7ZGQCBA8PFgIfEQUGJm5ic3A7ZGQCBQ8PFgIfEQUGJm5ic3A7ZGQCBg8PFgIfEQUGJm5ic3A7ZGQCBw8PFgIfEQUGJm5ic3A7ZGQCCA8PFgIfEQUGJm5ic3A7ZGQCAg8PFgIeBF9paWgFATBkFhJmDw8WAh8SaGQWAmYPDxYCHwRoZGQCAQ8PFgQfEQUGJm5ic3A7HxJoZGQCAg8PFgIfEQUDUE5DZGQCAw8PFgIfEQUIMi80LzIwMTZkZAIEDw8WAh8RBQFBZGQCBQ8PFgIfEQUFNzMxMDBkZAIGDw8WAh8RBQhQLTg3MkZXWWRkAgcPDxYCHxEFE0tNLiA0NyBDQS0xIE9SSUVOVEVkZAIIDw8WAh8RBQY0MDAuMDBkZAIDD2QWAmYPDxYCHxJoZGQYAQUeX19Db250cm9sc1JlcXVpcmVQb3N0QmFja0tleV9fFgEFGmN0bDAwJE1haW5Db250ZW50JGdyZFJlc3Vs/gd6va++k9UMqwhIJfCtA2WUy2k="},
		"__ASYNCPOST":                                     {"true"},
	}
	req, err := http.NewRequest(http.MethodPost, pncURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36")
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var tickets []Ticket
	doc, err := goquery.NewDocumentFromResponse(resp)
	rows := doc.Find("table tbody tr.rgRow")

	rows.Each(func(idx int, sel *goquery.Selection) {
		ticket := Ticket{Entity: pncEntity}
		sel.Children().Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1:
				ticket.Date = cleanStrings(s.Text())
			case 2:
				ticket.ID = cleanStrings(s.Text())
			case 3:
				ticket.ID = fmt.Sprintf("%s%s", ticket.ID, cleanStrings(s.Text()))
			case 5:
				ticket.Location = cleanStrings(s.Text())
			case 6:
				val := cleanStrings(s.Text())
				ticket.Ammount = val
				ticket.Total = val
			}
		})
		tickets = append(tickets, ticket)
	})
	return tickets, nil
}

func (PNC) getPlateType(plateType string) string {
	switch plateType {
	case "P":
		return "1"
	case "C":
		return "2"
	case "A":
		return "3"
	case "M":
		return "4"
	case "CD":
		return "5"
	case "CC":
		return "6"
	case "DIS":
		return "7"
	case "TC":
		return "8"
	case "O":
		return "9"
	case "U":
		return "10"
	case "MI":
		return "11"
	case "STP":
		return "12"
	case "TRC":
		return "13"
	case "PNC":
		return "61"
	}
	return "1"
}
