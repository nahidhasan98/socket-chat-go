package server

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/nahidhasan98/socket-chat-go/client"
)

type DataManager struct {
	Client  *client.Client
	Message []byte
}

func (dm *DataManager) Receive(cm *ClientManager) {
	for {
		messageType, p, err := dm.Client.WS.ReadMessage()
		if err != nil {
			cm.Unregister <- dm.Client
			dm.Message = []byte("left!")
			dm.broadcast(cm)
			dm.Client.Socket.Close()
			break
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := dm.Client.WS.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		// if length > 0 {
		// 	fmt.Println(dm.Client.Socket.RemoteAddr(), ": "+string(message))
		// 	dm.Message = []byte(message)
		// 	dm.broadcast(cm)
		// }
	}
}

func (dm *DataManager) send(cm *ClientManager) {
	defer dm.Client.Socket.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		if len(cm.Clients) == 0 {
			fmt.Println("No participants present here!")
		}

		for connection := range cm.Clients {
			msg := fmt.Sprintf("%s: %s", connection.Socket.LocalAddr(), message)
			connection.Socket.Write([]byte(msg))
		}
	}
}

func (dm *DataManager) broadcast(cm *ClientManager) {
	for connection := range cm.Clients {
		if connection != dm.Client {
			msg := fmt.Sprintf("%s: %s\n", dm.Client.Socket.RemoteAddr(), string(dm.Message))
			connection.Socket.Write([]byte(msg))
		}
	}
}
