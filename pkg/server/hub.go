package server

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan BroadcastEvent

	// Send message to a specific client.
	send chan Event

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub creates a new hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

// Run game hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case e := <-h.broadcast:
			for id, client := range h.clients {
				if _, ok := e.ExcludeClients[id]; ok {
					continue
				}
				select {
				case client.send <- e.Event:
				default:
					delete(h.clients, id)
					close(client.send)
				}
			}
		}
	}
}
