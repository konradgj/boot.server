package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	rootPath := "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(rootPath)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
