package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// JackpotService handles progressive jackpots
type JackpotService struct {
	mu            sync.RWMutex
	jacpots      map[string]*Jackpot
	contributions map[string][]*JackpotContribution
	 winners      map[string]*JackpotWinner
}

// Jackpot represents a progressive jackpot
type Jackpot struct {
	ID            string
	Name          string
	Provider      string
	GameID        string
	CurrentAmount float64
	SeedAmount    float64
	MinBet        float64
	Increment     float64 // Percentage of each bet that goes to jackpot
	TriggerAmount float64
	LastWinAmount float64
	LastWinTime   *time.Time
	LastWinnerID  string
	Status        string // active, triggered, reset
	Type          string // mini, minor, major, grand
}

// JackpotContribution represents a bet contribution to jackpot
type JackpotContribution struct {
	ID          string
	JackpotID  string
	UserID     string
	GameID     string
	BetAmount  float64
	ContribAmt float64
	Timestamp  time.Time
}

// JackpotWinner represents a jackpot win
type JackpotWinner struct {
	ID          string
	JackpotID   string
	UserID      string
	Amount      float64
	GameID      string
	Timestamp   time.Time
}

// NewJackpotService creates a new jackpot service
func NewJackpotService() *JackpotService {
	s := &JackpotService{
		jacpots:      make(map[string]*Jackpot),
		contributions: make(map[string][]*JackpotContribution),
		winners:      make(map[string]*JackpotWinner),
	}
	s.initializeJackpots()
	return s
}

func (s *JackpotService) initializeJackpots() {
	// Create different jackpot levels
	jackpots := []*Jackpot{
		{
			ID: "jp_mini_001", Name: "Mini Jackpot", Provider: "Tiger",
			GameID: "tiger_slots", CurrentAmount: 500, SeedAmount: 100,
			MinBet: 0.5, Increment: 0.02, TriggerAmount: 1000,
			Status: "active", Type: "mini",
		},
		{
			ID: "jp_minor_001", Name: "Minor Jackpot", Provider: "Tiger",
			GameID: "tiger_slots", CurrentAmount: 5000, SeedAmount: 1000,
			MinBet: 1.0, Increment: 0.03, TriggerAmount: 10000,
			Status: "active", Type: "minor",
		},
		{
			ID: "jp_major_001", Name: "Major Jackpot", Provider: "Tiger",
			GameID: "tiger_slots", CurrentAmount: 50000, SeedAmount: 10000,
			MinBet: 5.0, Increment: 0.04, TriggerAmount: 100000,
			Status: "active", Type: "major",
		},
		{
			ID: "jp_grand_001", Name: "Grand Jackpot", Provider: "Tiger",
			GameID: "tiger_slots", CurrentAmount: 500000, SeedAmount: 100000,
			MinBet: 10.0, Increment: 0.05, TriggerAmount: 1000000,
			Status: "active", Type: "grand",
		},
		{
			ID: "jp_prog_001", Name: "Progressive Slots Jackpot", Provider: "Pragmatic",
			GameID: "pragmatic_jackpot", CurrentAmount: 25000, SeedAmount: 5000,
			MinBet: 1.0, Increment: 0.03, TriggerAmount: 50000,
			Status: "active", Type: "major",
		},
		{
			ID: "jp_prog_002", Name: "Network Progressive", Provider: "Relax",
			GameID: "relax_jackpot", CurrentAmount: 100000, SeedAmount: 20000,
			MinBet: 2.0, Increment: 0.04, TriggerAmount: 200000,
			Status: "active", Type: "grand",
		},
	}

	for _, jp := range jackpots {
		s.jacpots[jp.ID] = jp
	}
}

// Contribute adds a bet contribution to a jackpot
func (s *JackpotService) Contribute(jackpotID, userID, gameID string, betAmount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jp, ok := s.jacpots[jackpotID]
	if !ok {
		return fmt.Errorf("jackpot not found")
	}

	if betAmount < jp.MinBet {
		return fmt.Errorf("bet amount below minimum for this jackpot")
	}

	contribAmt := betAmount * jp.Increment
	jp.CurrentAmount += contribAmt

	contribution := &JackpotContribution{
		ID:         uuid.New().String(),
		JackpotID:  jackpotID,
		UserID:     userID,
		GameID:     gameID,
		BetAmount:  betAmount,
		ContribAmt: contribAmt,
		Timestamp:  time.Now(),
	}

	s.contributions[jackpotID] = append(s.contributions[jackpotID], contribution)

	// Check if jackpot should trigger
	if jp.CurrentAmount >= jp.TriggerAmount && jp.Status == "active" {
		jp.Status = "triggered"
	}

	return nil
}

// GetJackpot returns a specific jackpot
func (s *JackpotService) GetJackpot(jackpotID string) (*Jackpot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jp, ok := s.jacpots[jackpotID]
	if !ok {
		return nil, fmt.Errorf("jackpot not found")
	}
	return jp, nil
}

// GetAllJackpots returns all jackpots
func (s *JackpotService) GetAllJackpots() []Jackpot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var jps []Jackpot
	for _, jp := range s.jacpots {
		jps = append(jps, *jp)
	}
	return jps
}

// GetJackpotsByType returns jackpots of a specific type
func (s *JackpotService) GetJackpotsByType(jackpotType string) []Jackpot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var jps []Jackpot
	for _, jp := range s.jacpots {
		if jp.Type == jackpotType {
			jps = append(jps, *jp)
		}
	}
	return jps
}

// TriggerJackpot triggers a jackpot win
func (s *JackpotService) TriggerJackpot(jackpotID, userID string) (*JackpotWinner, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jp, ok := s.jacpots[jackpotID]
	if !ok {
		return nil, fmt.Errorf("jackpot not found")
	}

	if jp.Status == "triggered" {
		return nil, fmt.Errorf("jackpot already triggered")
	}

	winAmount := jp.CurrentAmount
	jp.Status = "triggered"
	jp.LastWinAmount = winAmount
	now := time.Now()
	jp.LastWinTime = &now
	jp.LastWinnerID = userID

	// Record winner
	winner := &JackpotWinner{
		ID:        uuid.New().String(),
		JackpotID: jackpotID,
		UserID:    userID,
		Amount:    winAmount,
		GameID:    jp.GameID,
		Timestamp: now,
	}
	s.winners[winner.ID] = winner

	// Reset jackpot
	jp.CurrentAmount = jp.SeedAmount
	jp.Status = "active"

	return winner, nil
}

// GetWinners returns recent jackpot winners
func (s *JackpotService) GetWinners(limit int) []JackpotWinner {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ws []JackpotWinner
	for _, w := range s.winners {
		ws = append(ws, *w)
	}

	// Sort by timestamp (most recent first)
	// In production, use proper sorting

	if limit > 0 && len(ws) > limit {
		ws = ws[:limit]
	}

	return ws
}

// GetJackpotStats returns jackpot statistics
func (s *JackpotService) GetJackpotStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})
	var totalAmount float64
	var activeCount, triggeredCount int

	for _, jp := range s.jacpots {
		totalAmount += jp.CurrentAmount
		if jp.Status == "active" {
			activeCount++
		} else if jp.Status == "triggered" {
			triggeredCount++
		}
	}

	stats["total_jackpots"] = len(s.jacpots)
	stats["active_jackpots"] = activeCount
	stats["triggered_jackpots"] = triggeredCount
	stats["total_amount"] = totalAmount
	stats["total_winners"] = len(s.winners)

	return stats
}

// GetTotalJackpotAmount returns the sum of all jackpot amounts
func (s *JackpotService) GetTotalJackpotAmount() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total float64
	for _, jp := range s.jacpots {
		total += jp.CurrentAmount
	}
	return total
}

// GetTopWinners returns top winners by amount
func (s *JackpotService) GetTopWinners(limit int) []JackpotWinner {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ws []JackpotWinner
	for _, w := range s.winners {
		ws = append(ws, *w)
	}

	// Sort by amount descending (simplified)
	// In production, use proper sorting

	if limit > 0 && len(ws) > limit {
		ws = ws[:limit]
	}

	return ws
}
