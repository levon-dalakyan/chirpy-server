package main

import (
	"log"
	"net/http"
)

func main() {
	filepathRoot := "."
	port := "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
