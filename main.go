package main

import (
	"consumer/services"
	websocket "consumer/websoket"
	"log"
	"net/http"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	services.StartConsumer(hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWS(hub, w, r)
	})

	log.Println("Servidor WebSocket escuchando en :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
