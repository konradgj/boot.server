package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/konradgj/boot.server/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	database       *database.Queries
	fileserverHits atomic.Int32
	platform       string
}

func main() {
	port := "8080"
	rootPath := "."
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		database:       dbQueries,
		fileserverHits: atomic.Int32{},
		platform:       os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	handlerFileServer := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(rootPath)))
	mux.Handle("/app/", http.StripPrefix("/app", handlerFileServer))

	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirpy)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

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
