package manager

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID        string
	WebSocket *websocket.Conn
}

type ClientManager struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}

func (cm *ClientManager) Manage() {
	for {
		select {
		case cl := <-cm.Register:
			for connection := range cm.Clients {
				msg := fmt.Sprintf("Server: [%v] joined!\n", cl.WebSocket.RemoteAddr())
				connection.WebSocket.WriteMessage(1, []byte(msg))
			}
			cm.Clients[cl] = true
		case cl := <-cm.Unregister:
			delete(cm.Clients, cl)
			for connection := range cm.Clients {
				msg := fmt.Sprintf("Server: [%v] left!\n", cl.WebSocket.RemoteAddr())
				connection.WebSocket.WriteMessage(1, []byte(msg))
			}
		}
	}
}

func NewManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}
