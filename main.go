package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nahidhasan98/socket-chat-go/client"
	"github.com/nahidhasan98/socket-chat-go/server"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func main() {
	// flagMode := flag.String("mode", "server", "start in client or server mode")
	// flag.Parse()
	// if strings.ToLower(*flagMode) == "server" {
	// 	server.StartNewServer()
	// } else {
	// 	client.StartNewClient()
	// }

	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	r.LoadHTMLGlob("view/*")
	r.Static("/assets", "./assets")

	r.GET("/", index)
	r.GET("/ws", webSocket)

	r.Run(":6001")
	fmt.Println("Server running on port 6001...")
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Public Chat Room",
	})
}

func webSocket(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	dm := &server.DataManager{
		Client: &client.Client{
			WS: ws,
		},
		Message: nil,
	}
	dm.Receive(&server.ClientManager{})
}
