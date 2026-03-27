package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(""))))
	mux.Handle("/assets", http.FileServer(http.Dir("/assets/logo.png")))

	mux.HandleFunc("/healthz", handler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to start server")
		fmt.Println(err)
	}
}
