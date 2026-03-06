package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Result string `json:"result"`
}

func main() {
	http.HandleFunc("/", uppercaseHandler)
	log.Println("Service 1 (uppercase) running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func uppercaseHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp := Response{
		Result: strings.ToUpper(req.Text),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
