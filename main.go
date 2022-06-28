package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type Broadcast struct {
	client  *Client
	message []byte
}

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan Broadcast
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println(connection.socket.RemoteAddr(), "joined!")
		case connection := <-manager.unregister:
			if ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println(connection.socket.RemoteAddr(), "left!")
			}
		case broadcast := <-manager.broadcast:
			for connection := range manager.clients {
				if connection != broadcast.client {
					msg := fmt.Sprintf("%s: %s\n", broadcast.client.socket.RemoteAddr(), string(broadcast.message))
					connection.socket.Write([]byte(msg))
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			manager.broadcast <- Broadcast{client, []byte("left!")}
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println(client.socket.RemoteAddr(), ": "+string(message))
			manager.broadcast <- Broadcast{client, message}
			// manager.sendAll(client, message)
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		if len(manager.clients) == 0 {
			fmt.Println("No participants present here!")
		}
		for connection := range manager.clients {
			msg := fmt.Sprintf("%s: %s", connection.socket.LocalAddr(), message)
			connection.socket.Write([]byte(msg))
		}
	}
}

func (manager *ClientManager) sendAll(client *Client, message []byte) {
	for connection := range manager.clients {
		if connection != client {
			msg := fmt.Sprintf("%s: %s", client.socket.RemoteAddr(), string(message))
			connection.socket.Write([]byte(msg))
		}
	}
}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Print(string(message))
		}
	}
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":12345")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Broadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		manager.broadcast <- Broadcast{client, []byte("joined!")}
		go manager.receive(client)
		go manager.send(client)
	}
}

func startClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	client := Client{socket: connection}
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func main() {
	flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	}
}
