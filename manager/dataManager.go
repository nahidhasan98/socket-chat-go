package manager

import (
	"bufio"
	"fmt"
	"os"
)

type DataManager struct {
	Client  *Client
	Message []byte
}

func (dm *DataManager) Receive(cm *ClientManager) {
	for {
		messageType, p, err := dm.Client.WebSocket.ReadMessage()
		if err != nil {
			cm.Unregister <- dm.Client
			dm.Message = []byte("left!")
			dm.Client.WebSocket.Close()
		}
		// print out that message for clarity
		fmt.Println(28, messageType, string(p))
		dm.Broadcast(messageType, cm)
	}
}

func (dm *DataManager) Send(cm *ClientManager) {
	defer dm.Client.WebSocket.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		if len(cm.Clients) == 0 {
			fmt.Println("No participants present here!")
		}

		for connection := range cm.Clients {
			msg := fmt.Sprintf("%s: %s", connection.WebSocket.LocalAddr(), message)
			connection.WebSocket.WriteMessage(1, []byte(msg))
		}
	}
}

func (dm *DataManager) Broadcast(messageType int, cm *ClientManager) {
	fmt.Println(len(cm.Clients))
	for connection := range cm.Clients {
		if connection != dm.Client {
			msg := fmt.Sprintf("%s: %s\n", dm.Client.WebSocket.RemoteAddr(), string(dm.Message))
			connection.WebSocket.WriteMessage(1, []byte(msg))
		}
	}
}
