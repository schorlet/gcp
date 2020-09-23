package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWrapper struct {
	http.ResponseWriter
	code int
}

func (w *responseWrapper) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.code = code
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := responseWrapper{ResponseWriter: w}

		defer func() {
			if rcv := recover(); rcv != nil {
				http.Error(w, fmt.Sprintf("%v", rcv), http.StatusInternalServerError)
				wrapper.code = http.StatusInternalServerError
				log.Printf("%s\n", debug.Stack())
			}
			if wrapper.code == 0 {
				wrapper.code = http.StatusOK
			}
			log.Printf("%d %s %s (%s)\n", wrapper.code, r.Method, r.URL, time.Since(start))
		}()

		next(&wrapper, r)
	}
}
