package main

import (
	"flag"
	"fmt"
	"sync"

	"net/http"

	"time"

	"github.com/jhernandezme/go-multasgt"
)

func main() {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	var pType = flag.String("type", "P", "Plate Type")
	var pNumber = flag.String("number", "123ABC", "Plate Number")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		e := &multasgt.Emetra{}
		e.Client = client
		tickets, _ = e.Check(*pType, *pNumber)
		for _, r := range tickets {
			fmt.Printf("%#v \n", r)
		}
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		em := &multasgt.Emixtra{}
		em.Client = client
		tickets, _ = em.Check(*pType, *pNumber)
		for _, r := range tickets {
			fmt.Printf("%#v \n", r)
		}
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		scp := &multasgt.SCP{}
		scp.Client = client
		tickets, _ = scp.Check(*pType, *pNumber)
		for _, r := range tickets {
			fmt.Printf("%#v \n", r)
		}
	}()
	wg.Wait()

}