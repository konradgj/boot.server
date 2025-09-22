package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	port := "8080"
	rootPath := "."
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	handlerFileServer := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(rootPath)))
	mux.Handle("/app/", http.StripPrefix("/app", handlerFileServer))

	mux.HandleFunc("GET /api/healthz", handlerReady)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
