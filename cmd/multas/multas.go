package main

import (
	"flag"
	"fmt"

	"net/http"

	"time"

	"github.com/jhernandezb/go-multasgt"
)

func main() {
	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	var pType = flag.String("type", "P", "Plate Type")
	var pNumber = flag.String("number", "123ABC", "Plate Number")
	flag.Parse()
	ts, _ := multasgt.GetAllTickets(*pType, *pNumber, client)
	for _, t := range ts {
		fmt.Printf("%#v \n", t)
	}
}
