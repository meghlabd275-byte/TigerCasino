package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SportsbookMessage types for real-time sports updates
type SportsbookMessageType string

const (
	// Client messages
	MsgSubscribeMatch   SportsbookMessageType = "subscribe_match"
	MsgUnsubscribeMatch SportsbookMessageType = "unsubscribe_match"
	MsgSubscribeSport   SportsbookMessageType = "subscribe_sport"
	MsgPlaceBet         SportsbookMessageType = "place_bet"

	// Server messages
	MsgMatchUpdate      SportsbookMessageType = "match_update"
	MsgOddsChange       SportsbookMessageType = "odds_change"
	MsgScoreUpdate      SportsbookMessageType = "score_update"
	MsgGameStart        SportsbookMessageType = "game_start"
	MsgGameEnd          SportsbookMessageType = "game_end"
	MsgBetResult        SportsbookMessageType = "bet_result"
	MsgLiveScore        SportsbookMessageType = "live_score"
	MsgMarketUpdate     SportsbookMessageType = "market_update"
)

// SportsbookMessage represents a WebSocket message for sportsbook
type SportsbookMessage struct {
	Type    SportsbookMessageType `json:"type"`
	Payload interface{}          `json:"payload"`
	Timestamp int64              `json:"timestamp"`
}

// MatchUpdatePayload contains updated match data
type MatchUpdatePayload struct {
	MatchID     string    `json:"match_id"`
	Sport       string    `json:"sport"`
	League      string    `json:"league"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	HomeScore   int       `json:"home_score"`
	AwayScore   int       `json:"away_score"`
	Minute      int       `json:"minute"`
	Status      string    `json:"status"`
	HomeOdds    float64   `json:"home_odds"`
	DrawOdds    float64   `json:"draw_odds"`
	AwayOdds    float64   `json:"away_odds"`
	OverOdds    float64   `json:"over_odds"`
	UnderOdds   float64   `json:"under_odds"`
}

// OddsChangePayload contains odds change information
type OddsChangePayload struct {
	MatchID   string  `json:"match_id"`
	Market    string  `json:"market"`
	Selection string  `json:"selection"`
	OldOdds   float64 `json:"old_odds"`
	NewOdds   float64 `json:"new_odds"`
	Timestamp int64   `json:"timestamp"`
}

// LiveScorePayload contains live score updates
type LiveScorePayload struct {
	MatchID   string `json:"match_id"`
	Sport     string `json:"sport"`
	Event     string `json:"event"` // goal, foul, timeout, etc.
	Team      string `json:"team"`  // home or away
	Score     string `json:"score"`
	Minute    int    `json:"minute"`
	Timestamp int64  `json:"timestamp"`
}

// BetPlacementPayload contains bet placement data
type BetPlacementPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	MatchID   string    `json:"match_id"`
	BetType   string    `json:"bet_type"`
	Stake     float64   `json:"stake"`
	Odds      float64   `json:"odds"`
}

// BetResultPayload contains bet settlement results
type BetResultPayload struct {
	BetID     string    `json:"bet_id"`
	UserID    uuid.UUID `json:"user_id"`
	MatchID   string    `json:"match_id"`
	Result    string    `json:"result"` // won, lost, void
	Payout    float64   `json:"payout"`
	Profit    float64   `json:"profit"`
	Timestamp int64    `json:"timestamp"`
}

// SportsbookHub manages WebSocket connections for sportsbook
type SportsbookHub struct {
	// Registered clients
	clients map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Message broadcast
	broadcast chan []byte

	// Sport-specific channels for targeted updates
	sportChannels map[string]chan []byte

	// Match-specific channels
	matchChannels map[string]map[*Client]bool

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Statistics
	stats struct {
		totalConnections   int64
		activeConnections  int64
		messagesSent       int64
		messagesReceived   int64
	}
}

// NewSportsbookHub creates a new sportsbook WebSocket hub
func NewSportsbookHub() *SportsbookHub {
	hub := &SportsbookHub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte, 256),
		sportChannels:  make(map[string]chan []byte),
		matchChannels:  make(map[string]map[*Client]bool),
	}

	// Initialize sport channels
	sports := []string{"football", "basketball", "tennis", "esports", "hockey", "baseball", "mma", "cricket"}
	for _, sport := range sports {
		hub.sportChannels[sport] = make(chan []byte, 64)
	}

	return hub
}

// Run starts the sportsbook hub
func (h *SportsbookHub) Run() {
	log.Println("Sportsbook WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.stats.totalConnections++
			h.stats.activeConnections = int64(len(h.clients))
			h.mu.Unlock()
			log.Printf("Client connected to sportsbook. Total: %d", h.stats.activeConnections)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.stats.activeConnections = int64(len(h.clients))

				// Remove from match channels
				for matchID, clients := range h.matchChannels {
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.matchChannels, matchID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client disconnected from sportsbook. Total: %d", h.stats.activeConnections)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.stats.messagesSent += int64(len(h.clients))
			h.mu.RUnlock()
		}
	}
}

// SubscribeToMatch adds a client to a match's update channel
func (h *SportsbookHub) SubscribeToMatch(client *Client, matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.matchChannels[matchID]; !ok {
		h.matchChannels[matchID] = make(map[*Client]bool)
	}
	h.matchChannels[matchID][client] = true

	log.Printf("Client subscribed to match %s", matchID)
}

// UnsubscribeFromMatch removes a client from a match's update channel
func (h *SportsbookHub) UnsubscribeFromMatch(client *Client, matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.matchChannels[matchID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.matchChannels, matchID)
		}
	}
}

// BroadcastMatchUpdate broadcasts a match update to all subscribed clients
func (h *SportsbookHub) BroadcastMatchUpdate(update MatchUpdatePayload) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := SportsbookMessage{
		Type:      MsgMatchUpdate,
		Payload:   update,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling match update: %v", err)
		return
	}

	// Send to all clients subscribed to this match
	if clients, ok := h.matchChannels[update.MatchID]; ok {
		for client := range clients {
			select {
			case client.send <- data:
			default:
				delete(clients, client)
			}
		}
	}

	// Also broadcast to all clients
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			delete(h.clients, client)
		}
	}

	h.stats.messagesSent++
}

// BroadcastOddsChange broadcasts odds changes to relevant clients
func (h *SportsbookHub) BroadcastOddsChange(change OddsChangePayload) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := SportsbookMessage{
		Type:      MsgOddsChange,
		Payload:   change,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling odds change: %v", err)
		return
	}

	// Send to clients subscribed to this match
	if clients, ok := h.matchChannels[change.MatchID]; ok {
		for client := range clients {
			select {
			case client.send <- data:
			default:
				delete(clients, client)
			}
		}
	}

	// Also broadcast to sport-specific channel
	if sportChan, ok := h.sportChannels["all"]; ok {
		select {
		case sportChan <- data:
		default:
		}
	}
}

// BroadcastLiveScore broadcasts live score updates
func (h *SportsbookHub) BroadcastLiveScore(score LiveScorePayload) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := SportsbookMessage{
		Type:      MsgLiveScore,
		Payload:   score,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling live score: %v", err)
		return
	}

	// Send to all connected clients
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			delete(h.clients, client)
		}
	}

	h.stats.messagesSent++
}

// BroadcastBetResult broadcasts bet settlement results
func (h *SportsbookHub) BroadcastBetResult(result BetResultPayload) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg := SportsbookMessage{
		Type:      MsgBetResult,
		Payload:   result,
		Timestamp: time.Now().UnixMilli(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling bet result: %v", err)
		return
	}

	// Send to the specific user (in production, this would be targeted)
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			delete(h.clients, client)
		}
	}

	h.stats.messagesSent++
}

// HandleSportsbookMessage handles incoming sportsbook WebSocket messages
func (h *SportsbookHub) HandleSportsbookMessage(client *Client, message []byte) {
	h.stats.messagesReceived++

	var msg SportsbookMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling sportsbook message: %v", err)
		return
	}

	switch msg.Type {
	case MsgSubscribeMatch:
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if matchID, ok := payload["match_id"].(string); ok {
				h.SubscribeToMatch(client, matchID)
			}
		}

	case MsgUnsubscribeMatch:
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if matchID, ok := payload["match_id"].(string); ok {
				h.UnsubscribeFromMatch(client, matchID)
			}
		}

	case MsgSubscribeSport:
		// Handle sport subscription
		log.Printf("Client subscribed to sport updates")

	case MsgPlaceBet:
		// Handle real-time bet placement
		log.Printf("Received bet placement request")
	}
}

// GetStats returns current hub statistics
func (h *SportsbookHub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"total_connections":  h.stats.totalConnections,
		"active_connections": h.stats.activeConnections,
		"messages_sent":      h.stats.messagesSent,
		"messages_received":  h.stats.messagesReceived,
		"subscribed_matches": len(h.matchChannels),
	}
}

// GameHub manages game-specific WebSocket connections
type GameHub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	gameRooms  map[string]map[*Client]bool
	mu         sync.RWMutex
}

// NewGameHub creates a new game WebSocket hub
func NewGameHub() *GameHub {
	return &GameHub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
		gameRooms:  make(map[string]map[*Client]bool),
	}
}

// Run starts the game hub
func (h *GameHub) Run() {
	log.Println("Game WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove from all game rooms
				for _, room := range h.gameRooms {
					delete(room, client)
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// JoinGameRoom adds a client to a game room
func (h *GameHub) JoinGameRoom(client *Client, gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.gameRooms[gameID]; !ok {
		h.gameRooms[gameID] = make(map[*Client]bool)
	}
	h.gameRooms[gameID][client] = true

	log.Printf("Client joined game room %s", gameID)
}

// LeaveGameRoom removes a client from a game room
func (h *GameHub) LeaveGameRoom(client *Client, gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.gameRooms[gameID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.gameRooms, gameID)
		}
	}
}

// BroadcastToGameRoom sends a message to all clients in a game room
func (h *GameHub) BroadcastToGameRoom(gameID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.gameRooms[gameID]; ok {
		for client := range room {
			select {
			case client.send <- message:
			default:
				delete(room, client)
			}
		}
	}
}

// ChatHub manages real-time chat functionality
type ChatHub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan *ChatMessage
	rooms       map[string]map[*Client]bool
	mu          sync.RWMutex
	maxMessages int
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Timestamp int64     `json:"timestamp"`
	Type      string    `json:"type"` // message, system, bet
}

// NewChatHub creates a new chat WebSocket hub
func NewChatHub() *ChatHub {
	return &ChatHub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *ChatMessage, 256),
		rooms:       make(map[string]map[*Client]bool),
		maxMessages: 100,
	}
}

// Run starts the chat hub
func (h *ChatHub) Run() {
	log.Println("Chat WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove from all rooms
				for _, room := range h.rooms {
					delete(room, client)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			if room, ok := h.rooms[msg.RoomID]; ok {
				data, _ := json.Marshal(msg)
				for client := range room {
					select {
					case client.send <- data:
					default:
						delete(room, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// JoinRoom adds a client to a chat room
func (h *ChatHub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
}

// LeaveRoom removes a client from a chat room
func (h *ChatHub) LeaveRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// SendMessage broadcasts a chat message to a room
func (h *ChatHub) SendMessage(msg *ChatMessage) {
	h.broadcast <- msg
}

// GlobalHub holds all WebSocket hubs
type GlobalHub struct {
	Sportsbook *SportsbookHub
	Game       *GameHub
	Chat       *ChatHub
}

// NewGlobalHub creates and initializes all WebSocket hubs
func NewGlobalHub() *GlobalHub {
	return &GlobalHub{
		Sportsbook: NewSportsbookHub(),
		Game:       NewGameHub(),
		Chat:       NewChatHub(),
	}
}

// StartAll starts all WebSocket hubs
func (h *GlobalHub) StartAll() {
	go h.Sportsbook.Run()
	go h.Game.Run()
	go h.Chat.Run()
	log.Println("All WebSocket hubs started")
}

// GetHubStats returns statistics for all hubs
func (h *GlobalHub) GetHubStats() map[string]interface{} {
	return map[string]interface{}{
		"sportsbook": h.Sportsbook.GetStats(),
		"game":       h.Game.GetStats(),
		"chat":       h.Chat.GetStats(),
	}
}

// GetStats returns game hub statistics
func (h *GameHub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"active_connections": len(h.clients),
		"active_rooms":       len(h.gameRooms),
	}
}
