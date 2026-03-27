package main

import (
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	head := w.Header()
	head.Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(200)

	message := "OK"
	w.Write([]byte(message))

}
