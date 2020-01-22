package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/schorlet/gcp/dump"
)

var addr = flag.String("addr", "localhost:8001", "bind address")

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/favicon.png", http.NotFound)
	mux.HandleFunc("/", dump.Dump)

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf(err)
	}
}
