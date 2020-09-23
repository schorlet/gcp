package http

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func NewServer(handler http.Handler) (*http.Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8011"
	}

	return &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}, nil
}

func NewRouter(version string, helloHandler http.HandlerFunc) (http.Handler, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/favicon.png", http.NotFound)
	mux.HandleFunc("/liveness", Liveness)
	mux.HandleFunc("/readiness", Readiness)

	mux.HandleFunc("/v1/hello/", Logger(Hello(helloHandler)))

	mux.HandleFunc("/", Logger(Index(version)))

	return mux, nil
}

func Index(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "version: %s", version)
	}
}

func Hello(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// simulate latency
		random := rand.New(rand.NewSource(time.Now().Unix()))
		time.Sleep(time.Duration(random.Intn(1000)) * time.Millisecond)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	}
}

func Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
