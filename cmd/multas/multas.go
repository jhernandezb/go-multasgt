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
		Timeout: time.Duration(15 * time.Second),
	}
	var pType = flag.String("type", "P", "Plate Type")
	var pNumber = flag.String("number", "123ABC", "Plate Number")
	flag.Parse()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(7)
	var ts []multasgt.Ticket
	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		e := &multasgt.Emetra{}
		e.Client = client
		tickets, _ = e.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		em := &multasgt.Emixtra{}
		em.Client = client
		tickets, _ = em.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		scp := &multasgt.SCP{}
		scp.Client = client
		tickets, _ = scp.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		f := &multasgt.Fraijanes{}
		f.Client = client
		tickets, _ = f.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		v := &multasgt.VillaNueva{}
		v.Client = client
		tickets, _ = v.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		p := &multasgt.PNC{}
		p.Client = client
		tickets, _ = p.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		var tickets []multasgt.Ticket
		a := &multasgt.Antigua{}
		a.Client = client
		tickets, _ = a.Check(*pType, *pNumber)
		mutex.Lock()
		ts = append(ts, tickets...)
		mutex.Unlock()
	}()
	wg.Wait()
	for _, t := range ts {
		fmt.Printf("%#v \n", t)
	}

}
