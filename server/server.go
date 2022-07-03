package server

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/nahidhasan98/socket-chat-go/manager"
)

func StartNewServer(ws *websocket.Conn, cm *manager.ClientManager) {
	fmt.Println("Starting server...")

	cl := &manager.Client{WebSocket: ws}
	cm.Register <- cl

	dm := manager.DataManager{
		Client:  cl,
		Message: []byte("joined!"),
	}
	go dm.Broadcast(1, cm)

	go dm.Receive(cm)
	go dm.Send(cm)
}
