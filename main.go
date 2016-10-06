package main

import (
	"net/http"
	"fmt"
	"os"
	"log"
)

func index(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Accept", "application/json")
	//w.Header().Set("Accept-Encoding", "utf-8")
	//w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Encoding", "utf-8")
	//
	//json.NewEncoder(w).Encode("test")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println(port)
	http.HandleFunc("/", index)
	http.ListenAndServe(":" + port, nil)
}
