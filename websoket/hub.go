package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan interface{}
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan interface{}),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.Clients[conn] = true
			log.Println("Nuevo cliente conectado")
		case conn := <-h.Unregister:
			if _, ok := h.Clients[conn]; ok {
				delete(h.Clients, conn)
				conn.Close()
			}
		case msg := <-h.Broadcast:
			data, _ := json.Marshal(msg)
			for conn := range h.Clients {
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("Error enviando mensaje: %v", err)
					h.Unregister <- conn
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al hacer upgrade a websocket: %v", err)
		return
	}
	hub.Register <- conn

	go func() {
		defer func() { hub.Unregister <- conn }()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
