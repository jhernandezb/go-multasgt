package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	encode(w, getTickets())
}
func main() {
	m := http.NewServeMux()
	m.HandleFunc("/api/request", handleRequest)
	s := &http.Server{
		Addr:    ":8080",
		Handler: m,
	}
	s.ListenAndServe()
}

// Request represents the json request submitted to an endpoint.
type Request struct {
	Type  string `json:"type"`
	Plate string `json:"plate"`
}

func decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func encode(w io.Writer, v interface{}) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}
