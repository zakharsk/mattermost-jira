package main

import (
	"net/http"
	"encoding/json"
)

func index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Accept-Encoding", "utf-8")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "utf-8")

	json.NewEncoder(w).Encode("test")
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":5000", nil)
}
