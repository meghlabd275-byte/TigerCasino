package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// VirtualSportsService handles all virtual sports betting
type VirtualSportsService struct {
	db          *gorm.DB
	mu          sync.RWMutex
	events      map[string]*VirtualEvent
	currentRound int
}

type VirtualEvent struct {
	ID           string              `json:"id"`
	Sport        string              `json:"sport"`
	League       string              `json:"league"`
	HomeTeam    string              `json:"home_team"`
	AwayTeam    string              `json:"away_team"`
	StartTime   time.Time          `json:"start_time"`
	Status       string              `json:"status"` // scheduled, running, finished
	HomeScore   int               `json:"home_score"`
	AwayScore   int               `json:"away_score"`
	HomeOdds    float64            `json:"home_odds"`
	DrawOdds    float64            `json:"draw_odds"`
	AwayOdds    float64            `json:"away_odds"`
	Markets     map[string]map[string]float64 `json:"markets"`
	Winner      string              `json:"winner"`
	Minute      int                `json:"minute"`
}

type VirtualBet struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	EventID       string    `json:"event_id"`
	Sport         string    `json:"sport"`
	BetType       string    `json:"bet_type"`
	Stake         float64   `json:"stake"`
	Odds          float64   `json:"odds"`
	Selection     string    `json:"selection"`
	PotentialWin  float64   `json:"potential_win"`
	Status        string    `json:"status"` // pending, won, lost
	Result        string    `json:"result"`
	PlacedAt     time.Time `json:"placed_at"`
	SettledAt     *time.Time `json:"settled_at"`
}

// Virtual sports teams
var footballTeams = []string{
	"Tigers", "Eagles", "Wolves", "Sharks", "Panthers", "Dragons",
	"Lions", "Bears", "Falcons", "Hawks", "Cobras", "Vipers",
	"Storm", "Thunder", "Lightning", "Fire", "Ice", "Rock",
}

var basketballTeams = []string{
	"Warriors", "Lakers", "Celtics", "Heat", "Bulls", "Knicks",
	"Suns", "Nets", "Clippers", "Mavericks", "Bucks", "Raptors",
}

var tennisPlayers = []string{
	"Novak", "Rafael", "Roger", "Serena", "Stefanos", "Daniil",
	"Alexander", "Andrey", "Jannik", "Cameron", "Matteo", "Hubert",
}

var horseRacingHorses = []string{
	"Thunder Bolt", "Lightning Flash", "Star Dancer", "Midnight Run",
	"Golden Hoof", "Silver Streak", "Diamond Dust", "Emerald Eyes",
	"Ruby Rose", "Sapphire Blue", "Amber Glow", "Crystal Clear",
}

var greyhoundDogs = []string{
	"Swift Runner", "Quick Silver", "Lightning Paw", "Thunder Paw",
	"Star Blazer", "Moon Dancer", "Sun Chaser", "Wind Whisper",
}

func NewVirtualSportsService(db *gorm.DB) *VirtualSportsService {
	s := &VirtualSportsService{
		db:      db,
		events:  make(map[string]*VirtualEvent),
	}

	// Initialize virtual events
	s.initializeVirtualEvents()

	// Start simulation loop
	go s.simulationLoop()

	return s
}

func (s *VirtualSportsService) initializeVirtualEvents() {
	// Create upcoming football matches
	for i := 0; i < 10; i++ {
		homeTeam := footballTeams[i%len(footballTeams)]
		awayTeam := footballTeams[(i+3)%len(footballTeams)]
		
		homeOdds := 1.5 + float64(i)*0.2 + rand.Float64()*0.5
		drawOdds := 3.2 + rand.Float64()*0.5
		awayOdds := 2.5 + float64(i)*0.15 + rand.Float64()*0.5

		event := &VirtualEvent{
			ID:        fmt.Sprintf("vf_%03d", i),
			Sport:     "football",
			League:    getRandomLeague("football"),
			HomeTeam: homeTeam,
			AwayTeam: awayTeam,
			StartTime: time.Now().Add(time.Duration(i*5) * time.Minute),
			Status:    "scheduled",
			HomeOdds:  homeOdds,
			DrawOdds:  drawOdds,
			AwayOdds:  awayOdds,
			Markets: map[string]map[string]float64{
				"1x2":        {"home": homeOdds, "draw": drawOdds, "away": awayOdds},
				"over_2.5":  {"over": 1.85 + rand.Float64()*0.2},
				"under_2.5": {"under": 1.85 + rand.Float64()*0.2},
				"btts":      {"yes": 1.75 + rand.Float64()*0.2, "no": 2.0 + rand.Float64()*0.3},
			},
		}
		s.events[event.ID] = event
	}

	// Create upcoming basketball matches
	for i := 0; i < 6; i++ {
		homeTeam := basketballTeams[i%len(basketballTeams)]
		awayTeam := basketballTeams[(i+2)%len(basketballTeams)]

		homeOdds := 1.7 + rand.Float64()*0.4
		awayOdds := 2.0 + rand.Float64()*0.4

		event := &VirtualEvent{
			ID:        fmt.Sprintf("vb_%03d", i),
			Sport:     "basketball",
			League:    "Virtual NBA",
			HomeTeam: homeTeam,
			AwayTeam: awayTeam,
			StartTime: time.Now().Add(time.Duration(i*8) * time.Minute),
			Status:    "scheduled",
			HomeOdds:  homeOdds,
			AwayOdds:  awayOdds,
			Markets: map[string]map[string]float64{
				"1x2":     {"home": homeOdds, "away": awayOdds},
				"spread":  {"home_-5.5": 1.90, "away_+5.5": 1.90},
				"total":   {"over_210.5": 1.90, "under_210.5": 1.90},
			},
		}
		s.events[event.ID] = event
	}

	// Create tennis matches
	for i := 0; i < 6; i++ {
		player1 := tennisPlayers[i%len(tennisPlayers)]
		player2 := tennisPlayers[(i+2)%len(tennisPlayers)]

		player1Odds := 1.6 + rand.Float64()*0.6
		player2Odds := 2.0 + rand.Float64()*0.6

		event := &VirtualEvent{
			ID:        fmt.Sprintf("vt_%03d", i),
			Sport:     "tennis",
			League:    "Virtual ATP",
			HomeTeam: player1,
			AwayTeam: player2,
			StartTime: time.Now().Add(time.Duration(i*10) * time.Minute),
			Status:    "scheduled",
			HomeOdds:  player1Odds,
			AwayOdds:  player2Odds,
			Markets: map[string]map[string]float64{
				"1x2":   {"home": player1Odds, "away": player2Odds},
				"total": {"over_20.5": 1.90, "under_20.5": 1.90},
			},
		}
		s.events[event.ID] = event
	}

	// Create horse racing
	for i := 0; i < 3; i++ {
		raceName := fmt.Sprintf("Virtual Race %d", i+1)
		
		event := &VirtualEvent{
			ID:        fmt.Sprintf("vh_%03d", i),
			Sport:     "horse_racing",
			League:    "Virtual Horse Racing",
			HomeTeam: raceName,
			AwayTeam: "",
			StartTime: time.Now().Add(time.Duration(i*15) * time.Minute),
			Status:    "scheduled",
			Markets:   make(map[string]map[string]float64),
		}

		// Add win markets for each horse
		for j, horse := range horseRacingHorses {
			odds := 3.0 + float64(j)*0.5 + rand.Float64()*2.0
			if event.Markets["winner"] == nil {
				event.Markets["winner"] = make(map[string]float64)
			}
			event.Markets["winner"][fmt.Sprintf("horse_%d", j)] = odds
		}
		
		s.events[event.ID] = event
	}

	// Create greyhound racing
	for i := 0; i < 3; i++ {
		raceName := fmt.Sprintf("Greyhound Race %d", i+1)
		
		event := &VirtualEvent{
			ID:        fmt.Sprintf("vg_%03d", i),
			Sport:     "greyhound_racing",
			League:    "Virtual Greyhound Racing",
			HomeTeam: raceName,
			AwayTeam: "",
			StartTime: time.Now().Add(time.Duration(i*12) * time.Minute),
			Status:    "scheduled",
			Markets:   make(map[string]map[string]float64),
		}

		for j, dog := range greyhoundDogs {
			odds := 2.5 + float64(j)*0.4 + rand.Float64()*1.5
			if event.Markets["winner"] == nil {
				event.Markets["winner"] = make(map[string]float64)
			}
			event.Markets["winner"][fmt.Sprintf("dog_%d", j)] = odds
		}
		
		s.events[event.ID] = event
	}
}

func getRandomLeague(sport string) string {
	leagues := map[string][]string{
		"football": {"Virtual Premier League", "Virtual Champions League", "Virtual Cup"},
	}
	
	if sportLeagues, ok := leagues[sport]; ok {
		return sportLeagues[rand.Intn(len(sportLeagues))]
	}
	return "Virtual League"
}

// simulationLoop runs the virtual sports simulation
func (s *VirtualSportsService) simulationLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		
		// Check for events that need to start
		now := time.Now()
		for _, event := range s.events {
			if event.Status == "scheduled" && now.After(event.StartTime) {
				event.Status = "running"
				event.Minute = 0
				go s.simulateEvent(event.ID)
			}
		}

		s.mu.Unlock()
	}
}

func (s *VirtualSportsService) simulateEvent(eventID string) {
	s.mu.Lock()
	event, ok := s.events[eventID]
	s.mu.Unlock()

	if !ok {
		return
	}

	duration := 90 // 90 minutes for football
	if event.Sport == "basketball" {
		duration = 48 // 48 minutes
	} else if event.Sport == "tennis" {
		duration = 180 // Best of 3 sets
	} else if event.Sport == "horse_racing" || event.Sport == "greyhound_racing" {
		duration = 120 // Race duration
	}

	// Simulate match/race
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for event.Minute < duration {
		<-ticker.C
		
		s.mu.Lock()
		
		event.Minute += 5
		
		// Random scoring
		if event.Sport == "football" {
			scoreChange := rand.Float64()
			if scoreChange < 0.05 { // 5% chance per 5 min
				if rand.Float64() < 0.6 {
					event.HomeScore++
				} else {
					event.AwayScore++
				}
			}
		} else if event.Sport == "basketball" {
			scoreChange := rand.Float64()
			if scoreChange < 0.4 { // Higher scoring
				if rand.Float64() < 0.5 {
					event.HomeScore += 2 + rand.Intn(3)
				} else {
					event.AwayScore += 2 + rand.Intn(3)
				}
			}
		}

		// Update odds dynamically
		homeProb := float64(event.HomeScore+1) / float64(event.HomeScore+event.AwayScore+2)
		awayProb := float64(event.AwayScore+1) / float64(event.HomeScore+event.AwayScore+2)
		event.HomeOdds = 1.3 + (1-homeProb)*3
		event.AwayOdds = 1.3 + (1-awayProb)*3
		event.DrawOdds = 3.0 + rand.Float64()*1.0

		s.mu.Unlock()
	}

	// Finish event
	s.mu.Lock()
	event.Status = "finished"
	
	// Determine winner
	if event.HomeScore > event.AwayScore {
		event.Winner = "home"
	} else if event.AwayScore > event.HomeScore {
		event.Winner = "away"
	} else {
		event.Winner = "draw"
	}

	// Settle bets
	s.settleBets(eventID, event.Winner)

	// Create replacement event
	s.createReplacementEvent(event)

	s.mu.Unlock()
}

func (s *VirtualSportsService) createReplacementEvent(oldEvent *VirtualEvent) {
	// Create a new event to replace the finished one
	newEventID := fmt.Sprintf("%s_%d", oldEvent.Sport[:1], time.Now().UnixNano())
	
	newEvent := &VirtualEvent{
		ID:        newEventID,
		Sport:     oldEvent.Sport,
		League:    oldEvent.League,
		HomeTeam:  getRandomTeam(oldEvent.Sport),
		AwayTeam:  getRandomTeam(oldEvent.Sport),
		StartTime: time.Now().Add(5 * time.Minute),
		Status:    "scheduled",
		HomeOdds:  1.5 + rand.Float64()*1.0,
		DrawOdds:  3.0 + rand.Float64()*1.0,
		AwayOdds:  2.0 + rand.Float64()*1.0,
		Markets:   oldEvent.Markets,
	}

	s.events[newEventID] = newEvent
}

func getRandomTeam(sport string) string {
	switch sport {
	case "football":
		return footballTeams[rand.Intn(len(footballTeams))]
	case "basketball":
		return basketballTeams[rand.Intn(len(basketballTeams))]
	default:
		return tennisPlayers[rand.Intn(len(tennisPlayers))]
	}
}

func (s *VirtualSportsService) settleBets(eventID string, result string) {
	// In a full implementation, this would settle user bets in the database
	// For now, it's a placeholder for the betting logic
}

// ============ BETTING OPERATIONS ============

// PlaceVirtualBet places a bet on a virtual sports event
func (s *VirtualSportsService) PlaceVirtualBet(userID, eventID, selection string, stake float64) (*VirtualBet, error) {
	s.mu.RLock()
	event, ok := s.events[eventID]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	if event.Status != "scheduled" && event.Status != "running" {
		return nil, fmt.Errorf("event is not available for betting")
	}

	// Find odds
	var odds float64
	found := false

	// Try main markets first
	switch selection {
	case "home":
		odds = event.HomeOdds
		found = true
	case "away":
		odds = event.AwayOdds
		found = true
	case "draw":
		odds = event.DrawOdds
		found = true
	}

	// Try other markets
	if !found {
		for _, market := range event.Markets {
			if selOdds, ok := market[selection]; ok {
				odds = selOdds
				found = true
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("selection not available")
	}

	bet := &VirtualBet{
		ID:           uuid.New().String(),
		UserID:       userID,
		EventID:      eventID,
		Sport:        event.Sport,
		BetType:      "single",
		Stake:        stake,
		Odds:         odds,
		Selection:    selection,
		PotentialWin:  stake * odds,
		Status:       "pending",
		PlacedAt:     time.Now(),
	}

	// Record bet
	s.recordVirtualBet(bet)

	return bet, nil
}

// PlaceParlayBet places a multi-event virtual parlay
func (s *VirtualSportsService) PlaceParlayBet(userID string, selections []map[string]string, stake float64) ([]*VirtualBet, error) {
	if len(selections) < 2 || len(selections) > 6 {
		return nil, fmt.Errorf("parlay must have 2-6 selections")
	}

	var bets []*VirtualBet
	totalOdds := 1.0

	for _, sel := range selections {
		bet, err := s.PlaceVirtualBet(userID, sel["event_id"], sel["selection"], stake/float64(len(selections)))
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)
		totalOdds *= bet.Odds
	}

	// Update parlay odds
	for _, bet := range bets {
		bet.Odds = totalOdds
		bet.PotentialWin = stake * totalOdds
	}

	return bets, nil
}

// GetEvents returns all virtual events filtered by sport
func (s *VirtualSportsService) GetEvents(sport string, status string) []*VirtualEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*VirtualEvent
	for _, event := range s.events {
		if sport != "" && event.Sport != sport {
			continue
		}
		if status != "" && event.Status != status {
			continue
		}
		result = append(result, event)
	}

	// Sort by start time
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime.Before(result[j].StartTime)
	})

	return result
}

// GetEvent returns a specific event
func (s *VirtualSportsService) GetEvent(eventID string) (*VirtualEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

// GetLiveEvents returns currently running events
func (s *VirtualSportsService) GetLiveEvents() []*VirtualEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*VirtualEvent
	for _, event := range s.events {
		if event.Status == "running" {
			result = append(result, event)
		}
	}
	return result
}

// GetUpcomingEvents returns scheduled events
func (s *VirtualSportsService) GetUpcomingEvents(sport string, limit int) []*VirtualEvent {
	events := s.GetEvents(sport, "scheduled")
	if limit > 0 && len(events) > limit {
		events = events[:limit]
	}
	return events
}

// GetVirtualSports returns available virtual sports
func (s *VirtualSportsService) GetVirtualSports() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sportsSet := make(map[string]bool)
	for _, event := range s.events {
		sportsSet[event.Sport] = true
	}

	sports := make([]string, 0, len(sportsSet))
	for sport := range sportsSet {
		sports = append(sports, sport)
	}

	return sports
}

// GetHorseRacingRunners returns horse racing runners with odds
func (s *VirtualSportsService) GetHorseRacingRunners(eventID string) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	if event.Sport != "horse_racing" {
		return nil, fmt.Errorf("not a horse racing event")
	}

	var runners []map[string]interface{}
	for i, horse := range horseRacingHorses {
		odds := 3.0 + float64(i)*0.5 + rand.Float64()*2.0
		runners = append(runners, map[string]interface{}{
			"number": i + 1,
			"name":   horse,
			"odds":   odds,
		})
	}

	return runners, nil
}

// ============ HELPER FUNCTIONS ============

func (s *VirtualSportsService) recordVirtualBet(bet *VirtualBet) {
	dbBet := models.Bet{
		UserID:      uuid.MustParse(bet.UserID),
		GameType:    "virtual_" + bet.Sport,
		BetAmount:   bet.Stake,
		WinAmount:   bet.PotentialWin,
		Multiplier:  bet.Odds,
		Status:      bet.Status,
		GameData:    fmt.Sprintf(`{"event_id":"%s","selection":"%s"}`, bet.EventID, bet.Selection),
	}

	s.db.Create(&dbBet)
}

// GenerateSeeds creates seeds for provably fair games
func (s *VirtualSportsService) GenerateSeeds() (string, string, error) {
	serverSeedBytes := make([]byte, 32)
	if _, err := rand.Read(serverSeedBytes); err != nil {
		return "", "", err
	}
	serverSeed := hex.EncodeToString(serverSeedBytes)
	
	clientSeedBytes := make([]byte, 16)
	rand.Read(clientSeedBytes)
	clientSeed := hex.EncodeToString(clientSeedBytes)

	return serverSeed, clientSeed, nil
}
