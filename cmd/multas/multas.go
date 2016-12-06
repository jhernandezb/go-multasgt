package main

import (
	"flag"
	"fmt"
	"sync"

	"net/http"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jhernandezme/go-multasgt"
)

func downloadImage(url string) error {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// Session configuration
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-east-1")},
		Profile: "jhernandez",
	})
	if err != nil {
		return err
	}
	svc := s3manager.NewUploader(sess)

	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String("img.el-infractor.jhernandez.me"),
		//The location inside the bucket
		Key:         aws.String("images/test.jpg"),
		ContentType: aws.String(resp.Header.Get("Content-Type")),
		Body:        resp.Body,
	})
	if err != nil {
		fmt.Println("error to upload file: ", err)
		return err
	}
	fmt.Printf("Successfully uploaded %s to \n", result.Location)
	fmt.Println(result)
	return nil
}
func main() {
	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	var pType = flag.String("type", "P", "Plate Type")
	var pNumber = flag.String("number", "123ABC", "Plate Number")
	flag.Parse()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	// err := downloadImage("http://consultas.muniguate.com/consultas/fotos/Periferico_Sur_19_calle_Z11/27-05-2016-07-09-29-0.jpg")
	// fmt.Println(err)
	// if err != nil {
	// 	return
	// }
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
