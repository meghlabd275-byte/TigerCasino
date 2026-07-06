package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SportsbookService handles sports betting with advanced features
type SportsbookService struct {
	mu           sync.RWMutex
	bets         map[string]*SportsBet
	events       map[string]*SportsEvent
	odds         map[string]map[string]float64
	cashouts     map[string]*Cashout
}

// SportsBet represents a sports bet
type SportsBet struct {
	ID            string
	UserID        string
	EventID       string
	BetType       string // single, parlay, system
	Stake         float64
	Odds          float64
	PotentialWin  float64
	Selection    string // home, draw, away, over, under, etc.
	Status        string // pending, won, lost, cancelled, cashed_out
	PlacedAt      time.Time
	SettledAt     *time.Time
	CashoutAmount float64
	CashoutAt     *time.Time
}

// SportsEvent represents a sports event
type SportsEvent struct {
	ID           string
	Sport        string
	League       string
	HomeTeam     string
	AwayTeam     string
	StartTime    time.Time
	Status       string // scheduled, live, finished
	HomeScore    int
	AwayScore    int
	Period       string
	Minute       int
	Markets      map[string]map[string]float64 // market -> selection -> odds
}

// Cashout represents a cashout request
type Cashout struct {
	ID          string
	BetID       string
	UserID      string
	Amount      float64
	OriginalOdds float64
	CurrentOdds float64
	Status      string // pending, completed, failed
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// NewSportsbookService creates a new sportsbook service
func NewSportsbookService() *SportsbookService {
	s := &SportsbookService{
		bets:     make(map[string]*SportsBet),
		events:   make(map[string]*SportsEvent),
		odds:     make(map[string]map[string]float64),
		cashouts: make(map[string]*Cashout),
	}
	s.initializeEvents()
	return s
}

func (s *SportsbookService) initializeEvents() {
	// Sample football events
	s.events["fb_001"] = &SportsEvent{
		ID: "fb_001", Sport: "Football", League: "Premier League",
		HomeTeam: "Arsenal", AwayTeam: "Liverpool",
		StartTime: time.Now().Add(2 * time.Hour),
		Status: "scheduled",
		Markets: map[string]map[string]float64{
			"1x2": {"home": 2.50, "draw": 3.25, "away": 2.75},
			"over_2.5": {"over": 1.85, "under": 2.00},
			"btts": {"yes": 1.70, "no": 2.10},
		},
	}
	s.events["fb_002"] = &SportsEvent{
		ID: "fb_002", Sport: "Football", League: "La Liga",
		HomeTeam: "Real Madrid", AwayTeam: "Barcelona",
		StartTime: time.Now().Add(4 * time.Hour),
		Status: "scheduled",
		Markets: map[string]map[string]float64{
			"1x2": {"home": 2.20, "draw": 3.50, "away": 3.00},
			"over_2.5": {"over": 1.75, "under": 2.10},
			"btts": {"yes": 1.55, "no": 2.40},
		},
	}
	// Basketball
	s.events["bb_001"] = &SportsEvent{
		ID: "bb_001", Sport: "Basketball", League: "NBA",
		HomeTeam: "Lakers", AwayTeam: "Celtics",
		StartTime: time.Now().Add(3 * time.Hour),
		Status: "scheduled",
		Markets: map[string]map[string]float64{
			"1x2": {"home": 1.90, "away": 1.95},
			"over_210.5": {"over": 1.90, "under": 1.90},
		},
	}
	// Tennis
	s.events["tn_001"] = &SportsEvent{
		ID: "tn_001", Sport: "Tennis", League: "ATP",
		HomeTeam: "Djokovic", AwayTeam: "Nadal",
		StartTime: time.Now().Add(5 * time.Hour),
		Status: "scheduled",
		Markets: map[string]map[string]float64{
			"1x2": {"home": 1.65, "away": 2.30},
			"over_3.5_sets": {"over": 1.80, "under": 2.00},
		},
	}
}

// PlaceBet places a new bet
func (s *SportsbookService) PlaceBet(userID, eventID, betType, selection string, stake float64) (*SportsBet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	// Find odds for selection
	var odds float64
	found := false
	for market, selections := range event.Markets {
		if selOdds, ok := selections[selection]; ok {
			odds = selOdds
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("selection not available")
	}

	potentialWin := stake * odds

	bet := &SportsBet{
		ID:           uuid.New().String(),
		UserID:       userID,
		EventID:      eventID,
		BetType:      betType,
		Stake:        stake,
		Odds:         odds,
		PotentialWin: potentialWin,
		Selection:    selection,
		Status:       "pending",
		PlacedAt:     time.Now(),
	}

	s.bets[bet.ID] = bet
	return bet, nil
}

// CalculateCashout calculates cashout amount for a bet
func (s *SportsbookService) CalculateCashout(betID string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bet, ok := s.bets[betID]
	if !ok {
		return 0, fmt.Errorf("bet not found")
	}

	if bet.Status != "pending" {
		return 0, fmt.Errorf("bet cannot be cashed out")
	}

	// Get current odds (simplified - real implementation would fetch live odds)
	event, ok := s.events[bet.EventID]
	if !ok {
		return 0, fmt.Errorf("event not found")
	}

	currentOdds := event.Markets["1x2"][bet.Selection]
	if currentOdds == 0 {
		currentOdds = bet.Odds
	}

	// Cashout formula: stake * (current_odds / original_odds) * cashout_factor
	elapsedMinutes := time.Since(bet.PlacedAt).Minutes()
	cashoutFactor := 1.0 - (elapsedMinutes * 0.001) // Decrease slightly over time
	if cashoutFactor < 0.5 {
		cashoutFactor = 0.5
	}

	cashoutAmount := bet.Stake * (currentOdds / bet.Odds) * cashoutFactor
	return cashoutAmount, nil
}

// RequestCashout requests a cashout
func (s *SportsbookService) RequestCashout(betID string) (*Cashout, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bet, ok := s.bets[betID]
	if !ok {
		return nil, fmt.Errorf("bet not found")
	}

	cashoutAmount, err := s.CalculateCashout(betID)
	if err != nil {
		return nil, err
	}

	cashout := &Cashout{
		ID:           uuid.New().String(),
		BetID:        betID,
		UserID:       bet.UserID,
		Amount:       cashoutAmount,
		OriginalOdds: bet.Odds,
		CurrentOdds:  cashoutAmount / bet.Stake,
		Status:       "completed",
		CreatedAt:    time.Now(),
	}

	now := time.Now()
	cashout.CompletedAt = &now

	// Update bet status
	bet.Status = "cashed_out"
	bet.CashoutAmount = cashoutAmount
	bet.CashoutAt = &now

	s.cashouts[cashout.ID] = cashout

	return cashout, nil
}

// GetUserBets returns all bets for a user
func (s *SportsbookService) GetUserBets(userID string) []SportsBet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bets []SportsBet
	for _, bet := range s.bets {
		if bet.UserID == userID {
			bets = append(bets, *bet)
		}
	}
	return bets
}

// GetAvailableEvents returns all available events
func (s *SportsbookService) GetAvailableEvents() []SportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []SportsEvent
	for _, event := range s.events {
		if event.Status == "scheduled" {
			events = append(events, *event)
		}
	}
	return events
}

// GetLiveEvents returns live events
func (s *SportsbookService) GetLiveEvents() []SportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []SportsEvent
	for _, event := range s.events {
		if event.Status == "live" {
			events = append(events, *event)
		}
	}
	return events
}

// UpdateLiveOdds updates odds for live events (simulated)
func (s *SportsbookService) UpdateLiveOdds(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return fmt.Errorf("event not found")
	}

	// Simulate odds changes
	for market := range event.Markets {
		for selection := range event.Markets[market] {
			// Random odds adjustment
			change := 0.05
			event.Markets[market][selection] *= (1 + change)
			if event.Markets[market][selection] > 10 {
				event.Markets[market][selection] /= (1 + 2*change)
			}
		}
	}

	return nil
}

// SettleBet settles a bet
func (s *SportsbookService) SettleBet(betID string, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bet, ok := s.bets[betID]
	if !ok {
		return fmt.Errorf("bet not found")
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet already settled")
	}

	// Simplified settlement logic
	if result == bet.Selection {
		bet.Status = "won"
	} else {
		bet.Status = "lost"
	}

	now := time.Now()
	bet.SettledAt = &now

	return nil
}

// GetBetHistory returns bet history with filters
func (s *SportsbookService) GetBetHistory(userID string, status string) []SportsBet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bets []SportsBet
	for _, bet := range s.bets {
		if bet.UserID == userID {
			if status == "" || bet.Status == status {
				bets = append(bets, *bet)
			}
		}
	}
	return bets
}

// GetEvent returns event by ID
func (s *SportsbookService) GetEvent(eventID string) (*SportsEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

// CalculateParlayOdds calculates parlay odds
func (s *SportsbookService) CalculateParlayOdds(bets []SportsBet) float64 {
	multiplier := 1.0
	for _, bet := range bets {
		multiplier *= bet.Odds
	}
	// Apply parlay bonus
	if len(bets) >= 3 {
		multiplier *= 1.10 // 10% bonus for 3+ selections
	}
	return multiplier
}
