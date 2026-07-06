package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message types
const (
	MsgTypeChat       = "chat"
	MsgTypeBet        = "bet"
	MsgTypeGameState  = "game_state"
	MsgTypeLeaderboard = "leaderboard"
	MsgTypeSystem     = "system"
	MsgTypeUserJoin   = "user_join"
	MsgTypeUserLeave  = "user_leave"
)

// upgrader config
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Hub maintains active connections
type Hub struct {
	// Registered clients
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mutex      sync.RWMutex
}

// Client represents a WebSocket connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   string
	username string
}

// Message represents a WebSocket message
type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
	Sender    string          `json:"sender,omitempty"`
	Username  string          `json:"username,omitempty"`
}

// ChatPayload represents chat message payload
type ChatPayload struct {
	Message   string `json:"message"`
	GameID   string `json:"game_id,omitempty"`
	Room     string `json:"room"`
}

// GameStatePayload represents game state update
type GameStatePayload struct {
	GameType string          `json:"game_type"`
	State    json.RawMessage `json:"state"`
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("Client connected: %s", client.username)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected: %s", client.username)

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message.encode():
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// Broadcast broadcasts a message to all clients
func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- msg
}

// BroadcastToRoom broadcasts to clients in a specific room
func (h *Hub) BroadcastToRoom(room string, msg *Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		// In production, would check room membership
		select {
		case client.send <- msg.encode():
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// HandleWS handles WebSocket connections
func HandleWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		username := c.Query("username")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := &Client{
			hub:      hub,
			conn:     conn,
			send:     make(chan []byte, 256),
			userID:   userID,
			username: username,
		}

		hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}

// readPump reads messages from the WebSocket
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msg.Sender = c.userID
		msg.Username = c.username
		msg.Timestamp = time.Now().Unix()

		// Handle message based on type
		switch msg.Type {
		case MsgTypeChat:
			c.handleChat(&msg)
		case MsgTypeBet:
			// Broadcast bet to all clients
			c.hub.Broadcast(&msg)
		}
	}
}

// writePump writes messages to the WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleChat handles chat messages
func (c *Client) handleChat(msg *Message) {
	var payload ChatPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	// Broadcast to room
	msg.Payload, _ = json.Marshal(payload)
	c.hub.BroadcastToRoom(payload.Room, msg)
}

// encode encodes the message to JSON bytes
func (m *Message) encode() []byte {
	data, _ := json.Marshal(m)
	return data
}

// ============ Real-time Game Updates ============

// BroadcastCrashUpdate broadcasts crash game update
func (h *Hub) BroadcastCrashUpdate(crashPoint float64, status string) {
	payload := map[string]interface{}{
		"crash_point": crashPoint,
		"status":      status,
		"timestamp":   time.Now().UnixMilli(),
	}

	data, _ := json.Marshal(payload)

	msg := &Message{
		Type:      MsgTypeGameState,
		Payload:   data,
		Timestamp: time.Now().Unix(),
	}

	h.Broadcast(msg)
}

// BroadcastLeaderboardUpdate broadcasts leaderboard update
func (h *Hub) BroadcastLeaderboardUpdate(entries []map[string]interface{}) {
	data, _ := json.Marshal(entries)

	msg := &Message{
		Type:      MsgTypeLeaderboard,
		Payload:   data,
		Timestamp: time.Now().Unix(),
	}

	h.Broadcast(msg)
}

// BroadcastChatMessage broadcasts a chat message
func (h *Hub) BroadcastChatMessage(room, username, message string) {
	payload := ChatPayload{
		Message: message,
		Room:    room,
	}

	data, _ := json.Marshal(payload)

	msg := &Message{
		Type:     MsgTypeChat,
		Payload:  data,
		Sender:   username,
		Username: username,
		Timestamp: time.Now().Unix(),
	}

	h.BroadcastToRoom(room, msg)
}
