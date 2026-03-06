package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Text string `json:"text"`
}

type ServiceResponse struct {
	Result string `json:"result"`
}

type AggregatedResponse struct {
	Original  string `json:"original"`
	Uppercase string `json:"uppercase"`
	Reverse   string `json:"reverse"`
}

func main() {
	http.HandleFunc("/", aggregateHandler)
	log.Println("Main server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func aggregateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// creating an HTTP client with timeout
	client := &http.Client{Timeout: 5 * time.Second}

	// calling both services in parallel
	uppercaseCh := make(chan string)
	reverseCh := make(chan string)
	errorCh := make(chan error, 2)

	go callService(client, "http://localhost:8081", req.Text, uppercaseCh, errorCh)
	go callService(client, "http://localhost:8082", req.Text, reverseCh, errorCh)

	// get results
	var uppercase, reverse string
	for i := 0; i < 2; i++ {
		select {
		case res := <-uppercaseCh:
			uppercase = res
		case res := <-reverseCh:
			reverse = res
		case err := <-errorCh:
			log.Printf("Service error: %v", err)
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
	}

	response := AggregatedResponse{
		Original:  req.Text,
		Uppercase: uppercase,
		Reverse:   reverse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func callService(client *http.Client, url, text string, resultCh chan<- string, errorCh chan<- error) {
	reqBody := Request{Text: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		errorCh <- err
		return
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		errorCh <- err
		return
	}
	defer resp.Body.Close()

	var serviceResp ServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		errorCh <- err
		return
	}

	resultCh <- serviceResp.Result
}
