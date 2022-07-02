package server

import (
	"fmt"
	"net"

	"github.com/nahidhasan98/socket-chat-go/client"
)

func StartNewServer() {
	fmt.Println("Starting server...")

	listener, error := net.Listen("tcp", ":12345")
	if error != nil {
		fmt.Println(error)
	}

	cm := ClientManager{
		Clients:    make(map[*client.Client]bool),
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
	}
	go cm.manage()

	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}

		c := &client.Client{Socket: connection, Data: make(chan []byte)}
		cm.Register <- c

		dm := DataManager{c, []byte("joined!")}
		go dm.broadcast(&cm)

		go dm.Receive(&cm)
		go dm.send(&cm)
	}
}
