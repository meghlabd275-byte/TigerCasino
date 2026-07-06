package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EsportsService handles esports betting
type EsportsService struct {
	mu         sync.RWMutex
	events     map[string]*EsportsEvent
	bets       map[string][]*EsportsBet
	providers  []EsportsProvider
}

// EsportsEvent represents an esports match/event
type EsportsEvent struct {
	ID             string
	Sport          string // CS2, LoL, Dota2, Valorant, etc.
	League         string
	Team1          string
	Team2          string
	StartTime      time.Time
	Status         string // upcoming, live, completed
	Odds1          float64
	Odds2          float64
	DrawOdds       float64
	Winner         string // team1, team2, draw
	Scores         map[string]int
	Markets        []BettingMarket
}

// BettingMarket represents a betting market
type BettingMarket struct {
	ID      string
	Name    string
	Odds    map[string]float64
	Status  string
}

// EsportsBet represents a bet on an esports event
type EsportsBet struct {
	ID        string
	UserID   string
	EventID  string
	MarketID string
	Selection string
	Odds     float64
	Amount   float64
	PotentialWin float64
	Status   string // pending, won, lost
	CreatedAt time.Time
}

// EsportsProvider interface for external esports data
type EsportsProvider interface {
	GetName() string
	FetchEvents() ([]EsportsEvent, error)
	GetOdds(eventID string) (map[string]float64, error)
}

// NewEsportsService creates a new esports service
func NewEsportsService() *EsportsService {
	s := &EsportsService{
		events: make(map[string]*EsportsEvent),
		bets:   make(map[string][]*EsportsBet),
	}
	s.initializeMockEvents()
	return s
}

func (s *EsportsService) initializeMockEvents() {
	// CS2 Events
	s.events["cs2_001"] = &EsportsEvent{
		ID:         "cs2_001",
		Sport:      "CS2",
		League:     "ESL Pro League",
		Team1:      "Natus Vincere",
		Team2:      "FaZe Clan",
		StartTime:  time.Now().Add(2 * time.Hour),
		Status:     "upcoming",
		Odds1:      1.85,
		Odds2:      2.10,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("cs2_001"),
	}

	s.events["cs2_002"] = &EsportsEvent{
		ID:         "cs2_002",
		Sport:      "CS2",
		League:     "BLAST Premier",
		Team1:      "Team Vitality",
		Team2:      "G2 Esports",
		StartTime:  time.Now().Add(4 * time.Hour),
		Status:     "upcoming",
		Odds1:      1.75,
		Odds2:      2.25,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("cs2_002"),
	}

	// League of Legends Events
	s.events["lol_001"] = &EsportsEvent{
		ID:         "lol_001",
		Sport:      "League of Legends",
		League:     "LCK",
		Team1:      "T1",
		Team2:      "Gen.G",
		StartTime:  time.Now().Add(1 * time.Hour),
		Status:     "upcoming",
		Odds1:      1.90,
		Odds2:      2.00,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("lol_001"),
	}

	s.events["lol_002"] = &EsportsEvent{
		ID:         "lol_002",
		Sport:      "League of Legends",
		League:     "LEC",
		Team1:      "G2 Esports",
		Team2:      "Fnatic",
		StartTime:  time.Now().Add(3 * time.Hour),
		Status:     "upcoming",
		Odds1:      1.65,
		Odds2:      2.40,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("lol_002"),
	}

	// Dota 2 Events
	s.events["dota2_001"] = &EsportsEvent{
		ID:         "dota2_001",
		Sport:      "Dota 2",
		League:     "The International",
		Team1:      "Team Spirit",
		Team2:      "OG",
		StartTime:  time.Now().Add(5 * time.Hour),
		Status:     "upcoming",
		Odds1:      1.80,
		Odds2:      2.15,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("dota2_001"),
	}

	// Valorant Events
	s.events["val_001"] = &EsportsEvent{
		ID:         "val_001",
		Sport:      "Valorant",
		League:     "VCT Masters",
		Team1:      "Sentinels",
		Team2:      "LOUD",
		StartTime:  time.Now().Add(6 * time.Hour),
		Status:     "upcoming",
		Odds1:      2.00,
		Odds2:      1.90,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("val_001"),
	}

	// Rocket League
	s.events["rl_001"] = &EsportsEvent{
		ID:         "rl_001",
		Sport:      "Rocket League",
		League:     "RLCS World",
		Team1:      "Team BDS",
		Team2:      "NRG",
		StartTime:  time.Now().Add(8 * time.Hour),
		Status:     "upcoming",
		Odds1:      1.70,
		Odds2:      2.30,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("rl_001"),
	}

	// FIFA
	s.events["fifa_001"] = &EsportsEvent{
		ID:         "fifa_001",
		Sport:      "FIFA",
		League:     "ePremier League",
		Team1:      "Man City",
		Team2:      "Liverpool",
		StartTime:  time.Now().Add(30 * time.Minute),
		Status:     "upcoming",
		Odds1:      2.10,
		Odds2:      1.80,
		DrawOdds:   3.50,
		Scores:     map[string]int{"team1": 0, "team2": 0},
		Markets:    s.createDefaultMarkets("fifa_001"),
	}
}

func (s *EsportsService) createDefaultMarkets(eventID string) []BettingMarket {
	return []BettingMarket{
		{
			ID:     eventID + "_winner",
			Name:   "Match Winner",
			Odds:   map[string]float64{},
			Status: "open",
		},
		{
			ID:     eventID + "_map1",
			Name:   "First Map Winner",
			Odds:   map[string]float64{},
			Status: "open",
		},
		{
			ID:     eventID + "_total",
			Name:   "Total Maps Over/Under 2.5",
			Odds:   map[string]float64{"over": 1.90, "under": 1.90},
			Status: "open",
		},
	}
}

// GetEvents returns all available esports events
func (s *EsportsService) GetEvents(sport string, status string) []EsportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []EsportsEvent
	for _, e := range s.events {
		if sport != "" && e.Sport != sport {
			continue
		}
		if status != "" && e.Status != status {
			continue
		}
		events = append(events, *e)
	}
	return events
}

// GetEvent returns a specific event
func (s *EsportsService) GetEvent(eventID string) (*EsportsEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

// PlaceBet places a bet on an esports event
func (s *EsportsService) PlaceBet(userID, eventID, marketID, selection string, amount float64) (*EsportsBet, error) {
	s.mu.RLock()
	event, ok := s.events[eventID]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	if event.Status != "upcoming" {
		return nil, fmt.Errorf("cannot bet on live or completed events")
	}

	// Get odds
	var odds float64
	if selection == event.Team1 {
		odds = event.Odds1
	} else if selection == event.Team2 {
		odds = event.Odds2
	} else if selection == "draw" {
		odds = event.DrawOdds
	} else {
		return nil, fmt.Errorf("invalid selection")
	}

	bet := &EsportsBet{
		ID:            uuid.New().String(),
		UserID:        userID,
		EventID:       eventID,
		MarketID:      marketID,
		Selection:     selection,
		Odds:          odds,
		Amount:        amount,
		PotentialWin:  amount * odds,
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	s.mu.Lock()
	s.bets[eventID] = append(s.bets[eventID], bet)
	s.mu.Unlock()

	return bet, nil
}

// GetUserBets returns all bets for a user
func (s *EsportsService) GetUserBets(userID string) []EsportsBet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userBets []EsportsBet
	for _, eventBets := range s.bets {
		for _, bet := range eventBets {
			if bet.UserID == userID {
				userBets = append(userBets, *bet)
			}
		}
	}
	return userBets
}

// SettleEvent settles an event with the winner
func (s *EsportsService) SettleEvent(eventID, winner string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return fmt.Errorf("event not found")
	}

	event.Status = "completed"
	event.Winner = winner

	// Settle all bets
	for _, bet := range s.bets[eventID] {
		if bet.Status != "pending" {
			continue
		}

		if bet.Selection == winner {
			bet.Status = "won"
		} else {
			bet.Status = "lost"
		}
	}

	return nil
}

// GetLiveEvents returns all live events
func (s *EsportsService) GetLiveEvents() []EsportsEvent {
	return s.GetEvents("", "live")
}

// GetUpcomingEvents returns all upcoming events
func (s *EsportsService) GetUpcomingEvents() []EsportsEvent {
	return s.GetEvents("", "upcoming")
}

// GetEventsBySport returns events filtered by sport
func (s *EsportsService) GetEventsBySport(sport string) []EsportsEvent {
	return s.GetEvents(sport, "")
}

// GetPopularSports returns list of popular esports
func (s *EsportsService) GetPopularSports() []string {
	return []string{
		"CS2",
		"League of Legends",
		"Dota 2",
		"Valorant",
		"Rocket League",
		"FIFA",
	}
}
