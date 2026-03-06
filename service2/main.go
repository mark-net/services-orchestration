package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Result string `json:"result"`
}

func main() {
	http.HandleFunc("/", reverseHandler)
	log.Println("Service 2 (reverse) running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func reverseHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp := Response{
		Result: reverseString(req.Text),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
