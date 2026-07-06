package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// LiveChatService handles live chat support
type LiveChatService struct {
	mu           sync.RWMutex
	conversations map[string]*ChatConversation
	agents        map[string]*ChatAgent
	messages      map[string][]*ChatMessage
	queue         []string // UserIDs waiting
}

// ChatConversation represents a chat conversation
type ChatConversation struct {
	ID           string
	UserID       string
	AgentID      string
	Status       string // waiting, active, closed
	Priority     string // low, normal, high
	Department   string // general, payments, technical, VIP
	CreatedAt    time.Time
	ClosedAt     *time.Time
	Messages     int
	Rating       int
	Feedback     string
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID           string
	ConversationID string
	SenderType   string // user, agent, system
	SenderID     string
	Content      string
	Timestamp    time.Time
	ReadAt       *time.Time
}

// ChatAgent represents a support agent
type ChatAgent struct {
	ID           string
	Name         string
	Email        string
	Department   string
	Status       string // online, busy, offline
	MaxChats     int
	CurrentChats int
	Skills       []string
	JoinedAt     time.Time
}

// NewLiveChatService creates a new live chat service
func NewLiveChatService() *LiveChatService {
	s := &LiveChatService{
		conversations: make(map[string]*ChatConversation),
		agents:        make(map[string]*ChatAgent),
		messages:      make(map[string][]*ChatMessage),
		queue:         []string{},
	}
	s.initializeAgents()
	return s
}

func (s *LiveChatService) initializeAgents() {
	// Add default agents
	agents := []struct {
		id, name, email, dept string
	}{
		{"agent_1", "John Smith", "john@tigercasino.com", "general"},
		{"agent_2", "Sarah Johnson", "sarah@tigercasino.com", "payments"},
		{"agent_3", "Mike Wilson", "mike@tigercasino.com", "technical"},
		{"agent_4", "Emily Brown", "emily@tigercasino.com", "VIP"},
		{"agent_5", "David Lee", "david@tigercasino.com", "general"},
	}

	for _, a := range agents {
		s.agents[a.id] = &ChatAgent{
			ID:         a.id,
			Name:       a.name,
			Email:      a.email,
			Department: a.dept,
			Status:     "online",
			MaxChats:   5,
			Skills:     []string{a.dept, "general"},
			JoinedAt:  time.Now(),
		}
	}
}

// StartConversation starts a new chat conversation
func (s *LiveChatService) StartConversation(userID, department, priority string) (*ChatConversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation := &ChatConversation{
		ID:         uuid.New().String(),
		UserID:     userID,
		Status:     "waiting",
		Priority:   priority,
		Department: department,
		CreatedAt:  time.Now(),
	}

	s.conversations[conversation.ID] = conversation
	s.messages[conversation.ID] = []*ChatMessage{}
	s.queue = append(s.queue, conversation.ID)

	// Try to assign agent immediately
	s.assignAgent(conversation.ID)

	return conversation, nil
}

// SendMessage sends a message in a conversation
func (s *LiveChatService) SendMessage(conversationID, senderType, senderID, content string) (*ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation, ok := s.conversations[conversationID]
	if !ok {
		return nil, fmt.Errorf("conversation not found")
	}

	message := &ChatMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		SenderType:     senderType,
		SenderID:       senderID,
		Content:        content,
		Timestamp:      time.Now(),
	}

	s.messages[conversationID] = append(s.messages[conversationID], message)
	conversation.Messages++

	// If user sends message and no agent assigned, try to assign
	if senderType == "user" && conversation.AgentID == "" {
		s.assignAgent(conversationID)
	}

	return message, nil
}

// GetMessages returns all messages in a conversation
func (s *LiveChatService) GetMessages(conversationID string) ([]ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conversation, ok := s.conversations[conversationID]
	if !ok {
		return nil, fmt.Errorf("conversation not found")
	}

	msgs, ok := s.messages[conversationID]
	if !ok {
		return []ChatMessage{}, nil
	}

	// Mark messages as read
	now := time.Now()
	for _, m := range msgs {
		if m.SenderType == "agent" && m.ReadAt == nil {
			m.ReadAt = &now
		}
	}

	result := make([]ChatMessage, len(msgs))
	for i, m := range msgs {
		result[i] = *m
	}
	return result, nil
}

// AssignAgent assigns an agent to a conversation
func (s *LiveChatService) assignAgent(conversationID string) {
	conversation, ok := s.conversations[conversationID]
	if !ok {
		return
	}

	// Find available agent
	for _, agent := range s.agents {
		if agent.Status == "online" && agent.CurrentChats < agent.MaxChats {
			// Check department match
			if agent.Department == conversation.Department || agent.Department == "general" {
				conversation.AgentID = agent.ID
				conversation.Status = "active"
				agent.CurrentChats++

				// Add system message
				systemMsg := &ChatMessage{
					ID:             uuid.New().String(),
					ConversationID: conversationID,
					SenderType:     "system",
					SenderID:       "system",
					Content:        fmt.Sprintf("You are now connected with %s", agent.Name),
					Timestamp:      time.Now(),
				}
				s.messages[conversationID] = append(s.messages[conversationID], systemMsg)
				return
			}
		}
	}
}

// CloseConversation closes a conversation
func (s *LiveChatService) CloseConversation(conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversation, ok := s.conversations[conversationID]
	if !ok {
		return fmt.Errorf("conversation not found")
	}

	conversation.Status = "closed"
	now := time.Now()
	conversation.ClosedAt = &now

	// Free up agent
	if conversation.AgentID != "" {
		if agent, ok := s.agents[conversation.AgentID]; ok {
			agent.CurrentChats--
		}
	}

	return nil
}

// GetAgentStats returns agent statistics
func (s *LiveChatService) GetAgentStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})
	var totalActive, totalWaiting int

	for _, conv := range s.conversations {
		if conv.Status == "active" {
			totalActive++
		} else if conv.Status == "waiting" {
			totalWaiting++
		}
	}

	stats["active_conversations"] = totalActive
	stats["waiting_in_queue"] = totalWaiting
	stats["total_agents"] = len(s.agents)

	var onlineAgents int
	for _, agent := range s.agents {
		if agent.Status == "online" {
			onlineAgents++
		}
	}
	stats["online_agents"] = onlineAgents

	return stats
}

// SetAgentStatus sets agent status
func (s *LiveChatService) SetAgentStatus(agentID, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, ok := s.agents[agentID]
	if !ok {
		return fmt.Errorf("agent not found")
	}

	agent.Status = status
	return nil
}

// GetConversation returns a conversation by ID
func (s *LiveChatService) GetConversation(conversationID string) (*ChatConversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return nil, fmt.Errorf("conversation not found")
	}
	return conv, nil
}

// GetUserConversations returns all conversations for a user
func (s *LiveChatService) GetUserConversations(userID string) []ChatConversation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []ChatConversation
	for _, conv := range s.conversations {
		if conv.UserID == userID {
			result = append(result, *conv)
		}
	}
	return result
}

// RateConversation rates a conversation
func (s *LiveChatService) RateConversation(conversationID string, rating int, feedback string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return fmt.Errorf("conversation not found")
	}

	conv.Rating = rating
	conv.Feedback = feedback
	return nil
}
