package server

import (
	"fmt"

	"github.com/nahidhasan98/socket-chat-go/client"
)

type ClientManager struct {
	Clients    map[*client.Client]bool
	Register   chan *client.Client
	Unregister chan *client.Client
}

func (cm *ClientManager) manage() {
	for {
		select {
		case connection := <-cm.Register:
			cm.Clients[connection] = true
			fmt.Println(connection.Socket.RemoteAddr(), "joined!")
		case connection := <-cm.Unregister:
			if _, ok := cm.Clients[connection]; ok {
				close(connection.Data)
				delete(cm.Clients, connection)
				fmt.Println(connection.Socket.RemoteAddr(), "left!")
			}
		}
	}
}
