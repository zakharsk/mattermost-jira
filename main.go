package main

import (
	"net/http"
	"os"
	"encoding/json"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Accept-Encoding", "utf-8")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "utf-8")

	json.NewEncoder(w).Encode("test")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/", index)
	http.ListenAndServe(":" + port, nil)
}
