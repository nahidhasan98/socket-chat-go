package manager

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
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
		case connection := <-cm.Register:
			cm.Clients[connection] = true
			fmt.Println(connection.WebSocket.RemoteAddr(), "joined!")
			// broadcast to be done
		case connection := <-cm.Unregister:
			if _, ok := cm.Clients[connection]; ok {
				delete(cm.Clients, connection)
				fmt.Println(connection.WebSocket.RemoteAddr(), "left!")
				// broadcast to be done
			}
		}
	}
}
