package services

import (
	websocket "consumer/websoket"
	"encoding/json"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

type Pedido struct {
	PedidoID int    `json:"pedido_id"`
	Estado   string `json:"estado"`
}

func StartConsumer(hub *websocket.Hub) {
	conn, err := amqp.Dial("amqp://xk27:mando18D@44.215.213.150:5672/")
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}
	q, err := ch.QueueDeclare(
		"notificacion_pedidos", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar cola: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al consumir: %v", err)
	}

	go func() {
		for msg := range msgs {
			contentType := msg.ContentType
			log.Printf("Mensaje recibido con Content-Type: %s", contentType)

			switch contentType {
			case "application/json":
				var pedido Pedido
				if err := json.Unmarshal(msg.Body, &pedido); err != nil {
					log.Printf("Error deserializando JSON: %v", err)
					continue
				}
				log.Printf("Pedido recibido: %+v", pedido)
				hub.Broadcast <- pedido

			case "text/plain", "":
				texto := strings.TrimSpace(string(msg.Body))
				log.Printf("Mensaje de texto recibido: %s", texto)
				hub.Broadcast <- texto

			default:
				log.Printf("Tipo de mensaje no soportado: %s", contentType)
			}
		}
	}()
}
