package wsocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Configuración del upgrader para WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Permitir todas las conexiones para desarrollo
		return true
	},
}

// Cliente representa una conexión WebSocket
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	userID string
}

// Hub mantiene el conjunto de clientes activos y transmite mensajes
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// WSHub es la instancia global del hub
var WSHub = &Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Run ejecuta el hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Cliente WebSocket conectado. Total: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Cliente WebSocket desconectado. Total: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// readPump maneja los mensajes recibidos del cliente
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error WebSocket: %v", err)
			}
			break
		}
	}
}

// writePump maneja el envío de mensajes al cliente
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error enviando mensaje WebSocket: %v", err)
				return
			}
		}
	}
}

// HandleWebSocket maneja las conexiones WebSocket
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al actualizar conexión WebSocket: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  WSHub,
	}

	client.hub.register <- client

	// Iniciar goroutines para leer y escribir
	go client.writePump()
	go client.readPump()
}

// BroadcastMessage envía un mensaje a todos los clientes conectados
func BroadcastMessage(message []byte) {
	WSHub.broadcast <- message
}
