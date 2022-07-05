package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nahidhasan98/socket-chat-go/manager"
)

var cm *manager.ClientManager

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	cm = manager.NewManager()
	go cm.Manage()

	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.LoadHTMLGlob("view/*")
	router.Static("/assets", "./assets")

	router.GET("/", index)
	router.GET("/ws", webSocket)

	router.Run(":6001")
	fmt.Println("Server running on port 6001...")
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Public Chat Room",
	})
}

func webSocket(c *gin.Context) {
	fmt.Println("46#", c.Request.RemoteAddr)
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	checkError(err)

	cl := &manager.Client{WebSocket: ws}
	cm.Register <- cl

	dm := manager.DataManager{Client: cl}

	go dm.Receive(cm)
	go dm.Send(cm)
	fmt.Println("66#", dm.Client.WebSocket.LocalAddr(), "#", dm.Client.WebSocket.RemoteAddr())
}
