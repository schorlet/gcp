package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "localhost:8080", "bind address")

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

//go:generate minify -o assets/app.min.js assets/upload.js assets/media.js assets/app.js
//go:generate minify -o assets/style.min.css assets/style.css

func main() {
	flag.Parse()
	if port := os.Getenv("PORT"); port != "" {
		*addr = ":" + port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/favicon.png", http.NotFound)
	mux.HandleFunc("/opensearch.xml", http.NotFound)

	mux.HandleFunc("/uploadURL", signedURL)
	mux.HandleFunc("/detectText/", detectText)
	mux.Handle("/", cacheControl(http.FileServer(http.Dir("assets"))))

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func cacheControl(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, max-age=3600")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
