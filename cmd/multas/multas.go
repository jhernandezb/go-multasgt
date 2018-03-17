package main

import (
	"flag"
	"fmt"

	"net/http"

	"time"

	"github.com/jhernandezb/go-multasgt/entities/emetra"
)

func main() {
	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	var pType = flag.String("type", "P", "Plate Type")
	var pNumber = flag.String("number", "123ABC", "Plate Number")
	flag.Parse()
	plates, _ := emetra.Check(*pType, *pNumber, client)
	fmt.Printf("%+v", plates)
	// ts, _ := multasgt.GetAllTickets(*pType, *pNumber, client)
	// for _, t := range ts {
	// 	fmt.Printf("%#v \n", t)
	// }

}
