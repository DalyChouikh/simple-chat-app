package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/DalyChouikh/simple-chat-app/internal/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
		if websocket.IsWebSocketUpgrade(r) {
			next.ServeHTTP(w, r)
			log.Printf("WebSocket | Method: %s | Path: %s", r.Method, r.URL.Path)
			return
		}
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		log.Printf("Status: %d | Method: %s | Path: %s", rw.status, r.Method, r.URL.Path)
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSockets(manager *types.ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !websocket.IsWebSocketUpgrade(r) {
			http.Error(w, "Expected WebSocket protocol", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.NotFound(w, r)
			log.Printf("Error upgrading connection: %v", err)
			return
		}
		id, _ := uuid.NewRandom()
		client := types.NewClient(id.String(), conn)
		manager.Register <- client
		log.Printf("Client connected: %v", client.ID)
		go client.Read(manager)
		go client.Write()
	}
}

func main() {

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()

	manager := types.NewClientManager()

	go manager.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", handleWebSockets(&manager))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websocket.html")
	})

	handler := loggingMiddleware(mux)

	log.Printf("Server is starting on port %s", *addr)

	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
