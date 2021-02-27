package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richard-to/go-poker/pkg/server"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	hub := server.NewHub()
	gameState := server.StartGame()

	go hub.Run()

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		server.ServeWs(hub, gameState, c.Writer, c.Request)
	})

	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.Run("localhost:8000")
}
