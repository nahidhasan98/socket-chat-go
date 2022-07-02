package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	Socket net.Conn
	WS     *websocket.Conn
	Data   chan []byte
}

func (c *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := c.Socket.Read(message)
		if err != nil {
			c.Socket.Close()
			break
		}
		if length > 0 {
			fmt.Print(string(message))
		}
	}
}

func StartNewClient(conn *websocket.Conn) {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	c := Client{Socket: connection}
	go c.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}
