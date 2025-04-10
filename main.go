package main

import (
	"consumer/services"
	websocket "consumer/websoket"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	services.StartConsumer(hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWS(hub, w, r)
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	http.Handle("/", corsHandler.Handler(http.DefaultServeMux))

	log.Println("Servidor WebSocket escuchando en :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
