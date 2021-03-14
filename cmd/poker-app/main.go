package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/richard-to/go-poker/pkg/server"
)

func main() {
	// Set a random seed to get random card shuffle
	rand.Seed(time.Now().UnixNano())

	// Load environment variables
	env := os.Getenv("POKER_APP_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()

	// Start websocket hub and game state manager
	hub := server.NewHub()
	gameState := server.NewGameState()

	go hub.Run()

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Websocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		server.ServeWs(hub, gameState, c.Writer, c.Request)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": os.Getenv("POKER_APP_ENV"),
		})
	})

	// Serve static react build directory
	buildDir := os.Getenv("REACT_CLIENT_BUILD_DIR")
	r.StaticFile("/", buildDir+"/index.html")
	r.StaticFile("/robots.txt", buildDir+"/robots.txt")
	r.Static("/static", buildDir+"/static")
	r.Static("/images", buildDir+"/images")

	// Google App Engine will set the port automatically
	port := os.Getenv("PORT")

	r.Run(":" + port)
}
