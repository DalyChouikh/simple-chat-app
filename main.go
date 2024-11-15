package main

import (
	"flag"
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		log.Printf("Status: %d | Method: %s | Path: %s", rw.status, r.Method, r.URL.Path)
	})
}

func main() {

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))

	})

	handler := loggingMiddleware(mux)

	log.Printf("Server is starting on port %s", *addr)

	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
