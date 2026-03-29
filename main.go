package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/assets", http.FileServer(http.Dir("/assets/logo.png")))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
