package server

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Fix this later; only using for testing
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	conn      *websocket.Conn
	hub       *Hub
	id        string
	gameState *GameState
	// Buffered channel of outbound messages.
	send     chan Event
	username string
	seatID   string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	// Send unregister event to hub using defer is a good idea since it will run
	// after the function finishes
	defer func() {
		DisconnectPlayer(c)
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// Set max message size
	c.conn.SetReadLimit(maxMessageSize)
	// Set timeout for read
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// Unsure what this means. But heartbeat, I think
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// Read message
		var e Event
		err := c.conn.ReadJSON(&e)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		ProcessEvent(c, e)

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	// Use defer for clean up
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		// On send event, write message to to client
		// This is done via broadcast event in hub
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(event)
			if err != nil {
				return
			}
		case <-ticker.C:
			// Set write timeout
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Write ping message. Basically heartbeat
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, gameState *GameState, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn:      conn,
		gameState: gameState,
		hub:       hub,
		id:        uuid.New().String(),
		send:      make(chan Event, 256),
	}

	// So when the websocket is activated, add/register client to hub
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	// Start read and write operations in goroutines
	go client.writePump()
	go client.readPump()
}
