package multasgt

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var cleanStringsRegex = regexp.MustCompile(`[\t\n\r]`)

// Ticket represents the information related to the ticket.
type Ticket struct {
	ID       string `json:"id"`
	Entity   string `json:"entity"`
	Date     string `json:"date"`
	Ammount  string `json:"ammount"`
	Discount string `json:"discount"`
	Total    string `json:"total"`
	Location string `json:"location"`
	Info     string `json:"info"`
	Photo    string `json:"photo"`
}

// TicketChecker is the interface that all checkers must implement.
type TicketChecker interface {
	Check(plateType, plateNumber string) ([]Ticket, error)
}

func GetAttribute(attrName string, n *html.Node) string {
	if n == nil {
		return ""
	}
	for i, a := range n.Attr {
		if a.Key == attrName {
			return n.Attr[i].Val
		}
	}
	return ""
}

// CleanStrings removes any white unnecessary whitespace
func CleanStrings(s string) string {
	return cleanStringsRegex.ReplaceAllString(strings.TrimSpace(s), "")
}
