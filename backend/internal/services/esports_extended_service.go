package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// EsportsService handles all esports betting operations
type EsportsService struct {
	db      *gorm.DB
	mu     sync.RWMutex
	events map[string]*EsportsEvent
}

type EsportsEvent struct {
	ID            string                   `json:"id"`
	Game          string                   `json:"game"` // cs2, lol, dota2, valorant, overwatch
	League        string                   `json:"league"`
	Team1         *EsportsTeam            `json:"team1"`
	Team2         *EsportsTeam            `json:"team2"`
	Format        string                   `json:"format"` // bo1, bo3, bo5
	Status        string                   `json:"status"` // scheduled, live, finished
	StartTime     time.Time               `json:"start_time"`
	Winner        string                   `json:"winner"`
	_maps         []string                `json:"maps"`
	CurrentMap    int                     `json:"current_map"`
	Scores        map[string]int           `json:"scores"`
	Markets       map[string]map[string]float64 `json:"markets"`
	LiveOdds      map[string]float64      `json:"live_odds"`
	Streams       []string                `json:"streams"`
}

type EsportsTeam struct {
	ID      string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Seed   int    `json:"seed"`
	Odds   float64 `json:"odds"`
}

// Esports betting markets
var esportsMarkets = map[string]map[string]bool{
	"cs2": {
		"map_winner": true, "map1_winner": true, "map2_winner": true, "map3_winner": true,
		"total_maps": true, "exact_score": true, "first_blood": true,
		"first_tower": true, "first_pistol": true,
	},
	"lol": {
		"map_winner": true, "first_blood": true, "first_tower": true,
		"first_dragon": true, "first_baron": true, "total_maps": true,
		"correct_score": true,
	},
	"dota2": {
		"map_winner": true, "first_blood": true, "first_tower": true,
		"first_roshan": true, "total_maps": true, "correct_score": true,
	},
	"valorant": {
		"map_winner": true, "first_blood": true, "first_spike": true,
		"total_maps": true, "correct_score": true,
	},
}

// Team databases
var cs2Teams = []*EsportsTeam{
	{ID: "navi", Name: "Natus Vincere", Region: "EU", Seed: 1, Odds: 2.5},
	{ID: "faze", Name: "FaZe Clan", Region: "EU", Seed: 2, Odds: 3.0},
	{ID: "g2", Name: "G2 Esports", Region: "EU", Seed: 3, Odds: 3.5},
	{ID: "liquid", Name: "Team Liquid", Region: "NA", Seed: 4, Odds: 4.0},
	{ID: "cloud9", Name: "Cloud9", Region: "NA", Seed: 5, Odds: 4.5},
	{ID: "vitality", Name: "Team Vitality", Region: "EU", Seed: 6, Odds: 5.0},
	{ID: "astralis", Name: "Astralis", Region: "EU", Seed: 7, Odds: 5.5},
	{ID: "heroic", Name: "Heroic", Region: "EU", Seed: 8, Odds: 6.0},
}

var lolTeams = []*EsportsTeam{
	{ID: "t1", Name: "T1", Region: "KR", Seed: 1, Odds: 2.0},
	{ID: "gen", Name: "Gen.G", Region: "KR", Seed: 2, Odds: 2.5},
	{ID: "jdg", Name: "JD Gaming", Region: "CN", Seed: 3, Odds: 3.0},
	{ID: "bilibili", Name: "Bilibili Gaming", Region: "CN", Seed: 4, Odds: 3.5},
	{ID: "geng", Name: "Gen.G", Region: "KR", Seed: 5, Odds: 4.0},
	{ID: "fnc", Name: "Fnatic", Region: "EU", Seed: 6, Odds: 5.0},
	{ID: "g2_lol", Name: "G2 Esports", Region: "EU", Seed: 7, Odds: 6.0},
	{ID: "c9_lol", Name: "Cloud9", Region: "NA", Seed: 8, Odds: 7.0},
}

var dota2Teams = []*EsportsTeam{
	{ID: "og", Name: "OG", Region: "EU", Seed: 1, Odds: 2.5},
	{ID: "spirit", Name: "Team Spirit", Region: "EU", Seed: 2, Odds: 3.0},
	{ID: "lgd", Name: "PSG.LGD", Region: "CN", Seed: 3, Odds: 3.5},
	{ID: "xtreme", Name: "Xtreme Gaming", Region: "CN", Seed: 4, Odds: 4.0},
	{ID: "talon", Name: "Talon Esports", Region: "SEA", Seed: 5, Odds: 5.0},
	{ID: "betera", Name: "Betera", Region: "CIS", Seed: 6, Odds: 5.5},
	{ID: "nigma", Name: "Nigma Galaxy", Region: "EU", Seed: 7, Odds: 6.0},
	{ID: "sacre", Name: "Team Secret", Region: "EU", Seed: 8, Odds: 7.0},
}

var valorantTeams = []*EsportsTeam{
	{ID: "sen", Name: "Sentinels", Region: "NA", Seed: 1, Odds: 2.0},
	{ID: "fnc_val", Name: "Fnatic", Region: "EU", Seed: 2, Odds: 2.5},
	{ID: "prx", Name: "Paper Rex", Region: "SEA", Seed: 3, Odds: 3.0},
	{ID: "drx", Name: "DRX", Region: "KR", Seed: 4, Odds: 3.5},
	{ID: "lev", Name: "LOUD", Region: "BR", Seed: 5, Odds: 4.0},
	{ID: "kru", Name: "KRU Esports", Region: "LATAM", Seed: 6, Odds: 5.0},
	{ID: "navi_val", Name: "Natus Vincere", Region: "EU", Seed: 7, Odds: 6.0},
	{ID: "optic", Name: "OpTic Gaming", Region: "NA", Seed: 8, Odds: 7.0},
}

var cs2Maps = []string{"Inferno", "Mirage", "Nuke", "Overpass", "Ancient", "Vertigo", "Anubis"}

func NewEsportsService(db *gorm.DB) *EsportsService {
	s := &EsportsService{
		db:      db,
		events:  make(map[string]*EsportsEvent),
	}

	s.initializeEvents()

	// Start simulation loop
	go s.simulationLoop()

	return s
}

func (s *EsportsService) initializeEvents() {
	// CS2 events
	for i := 0; i < 5; i++ {
		team1 := cs2Teams[i%len(cs2Teams)]
		team2 := cs2Teams[(i+2)%len(cs2Teams)]

		event := s.createEsportsEvent("cs2", fmt.Sprintf("CS2 Match %d", i+1), team1, team2, "bo3")
		s.events[event.ID] = event
	}

	// League of Legends events
	for i := 0; i < 4; i++ {
		team1 := lolTeams[i%len(lolTeams)]
		team2 := lolTeams[(i+2)%len(lolTeams)]

		event := s.createEsportsEvent("lol", fmt.Sprintf("LEC Match %d", i+1), team1, team2, "bo1")
		s.events[event.ID] = event
	}

	// Dota 2 events
	for i := 0; i < 4; i++ {
		team1 := dota2Teams[i%len(dota2Teams)]
		team2 := dota2Teams[(i+2)%len(dota2Teams)]

		event := s.createEsportsEvent("dota2", fmt.Sprintf("DPC Match %d", i+1), team1, team2, "bo3")
		s.events[event.ID] = event
	}

	// Valorant events
	for i := 0; i < 4; i++ {
		team1 := valorantTeams[i%len(valorantTeams)]
		team2 := valorantTeams[(i+2)%len(valorantTeams)]

		event := s.createEsportsEvent("valorant", fmt.Sprintf("VCT Match %d", i+1), team1, team2, "bo3")
		s.events[event.ID] = event
	}
}

func (s *EsportsService) createEsportsEvent(game, league string, team1, team2 *EsportsTeam, format string) *EsportsEvent {
	// Calculate odds based on seeds
	team1Odds := 1.5 + float64(team2.Seed)*0.3 + rand.Float64()*0.5
	team2Odds := 1.5 + float64(team1.Seed)*0.3 + rand.Float64()*0.5

	event := &EsportsEvent{
		ID:         fmt.Sprintf("es_%s_%d", game, time.Now().UnixNano()),
		Game:       game,
		League:     league,
		Team1:      team1,
		Team2:      team2,
		Format:     format,
		Status:     "scheduled",
		StartTime:  time.Now().Add(time.Duration(5+i*3) * time.Minute),
		Scores:     make(map[string]int),
		Markets:    s.generateMarkets(game, team1Odds, team2Odds),
		LiveOdds:   map[string]float64{"team1": team1Odds, "team2": team2Odds},
		CurrentMap:  0,
		_maps:       s.getMaps(game),
	}

	// Add map count
	event.Scores["team1_maps"] = 0
	event.Scores["team2_maps"] = 0

	return event
}

func (s *EsportsService) getMaps(game string) []string {
	switch game {
	case "cs2":
		// Random 3-5 maps from pool
		maps := make([]string, len(cs2Maps))
		copy(maps, cs2Maps)
		rand.Shuffle(len(maps), func(i, j int) {
			maps[i], maps[j] = maps[j], maps[i]
		})
		numMaps := 3 + rand.Intn(3) // 3-5 maps
		return maps[:numMaps]
	case "lol", "dota2":
		return []string{"Map 1", "Map 2", "Map 3", "Map 4", "Map 5"}
	case "valorant":
		return []string{"Ascent", "Haven", "Split", "Bind", "Icebox"}
	default:
		return []string{"Map 1", "Map 2", "Map 3"}
	}
}

func (s *EsportsService) generateMarkets(game string, team1Odds, team2Odds float64) map[string]map[string]float64 {
	markets := map[string]map[string]float64{
		"map_winner": {
			"team1": team1Odds,
			"team2": team2Odds,
		},
		"total_maps": {
			"over_2.5": 1.85,
			"under_2.5": 1.95,
		},
		"first_blood": {
			"team1": 1.90,
			"team2": 1.90,
		},
	}

	// Add game-specific markets
	switch game {
	case "cs2":
		markets["first_pistol"] = map[string]float64{
			"team1": 1.90, "team2": 1.90,
		}
		markets["first_tower"] = map[string]float64{
			"team1": 1.90, "team2": 1.90,
		}
	case "lol", "dota2":
		markets["first_dragon"] = map[string]float64{
			"team1": 1.90, "team2": 1.90,
		}
		markets["first_baron"] = map[string]float64{
			"team1": 1.90, "team2": 1.90,
		}
	case "valorant":
		markets["first_spike"] = map[string]float64{
			"team1": 1.90, "team2": 1.90,
		}
	}

	return markets
}

// simulationLoop runs the esports simulation
func (s *EsportsService) simulationLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		
		now := time.Now()
		for _, event := range s.events {
			// Start scheduled events
			if event.Status == "scheduled" && now.After(event.StartTime) {
				event.Status = "live"
				event.CurrentMap = 1
				go s.simulateEvent(event.ID)
			}
		}

		s.mu.Unlock()
	}
}

func (s *EsportsService) simulateEvent(eventID string) {
	s.mu.Lock()
	event, ok := s.events[eventID]
	s.mu.Unlock()

	if !ok {
		return
	}

	// Simulate until finished
	for {
		s.mu.Lock()
		if event.Status == "finished" {
			s.mu.Unlock()
			break
		}

		// Random scoring (simplified)
		if rand.Float64() < 0.3 {
			if rand.Float64() < 0.5 {
				event.Scores["team1"]++
			} else {
				event.Scores["team2"]++
			}
		}

		// Update live odds
		totalScore := event.Scores["team1"] + event.Scores["team2"]
		if totalScore > 0 {
			team1Prob := float64(event.Scores["team1"]) / float64(totalScore)
			event.LiveOdds["team1"] = 1.2 + (1-team1Prob)*3
			event.LiveOdds["team2"] = 1.2 + team1Prob*3
		}

		// Check if map is complete (first to win based on format)
		mapsToWin := 2
		if event.Format == "bo1" {
			mapsToWin = 1
		} else if event.Format == "bo5" {
			mapsToWin = 3
		}

		if event.Scores["team1"] >= mapsToWin {
			event.Scores["team1_maps"] = event.Scores["team1_maps"].(int) + 1
			if event.Scores["team1_maps"].(int) >= mapsToWin {
				event.Winner = "team1"
				event.Status = "finished"
				s.settleEsportsBets(eventID, "team1")
			} else {
				// Start next map
				event.CurrentMap++
				event.Scores["team1"] = 0
				event.Scores["team2"] = 0
			}
		} else if event.Scores["team2"] >= mapsToWin {
			event.Scores["team2_maps"] = event.Scores["team2_maps"].(int) + 1
			if event.Scores["team2_maps"].(int) >= mapsToWin {
				event.Winner = "team2"
				event.Status = "finished"
				s.settleEsportsBets(eventID, "team2")
			} else {
				event.CurrentMap++
				event.Scores["team1"] = 0
				event.Scores["team2"] = 0
			}
		}

		s.mu.Unlock()

		time.Sleep(5 * time.Second)
	}

	// Create replacement event
	s.mu.Lock()
	s.createReplacementEvent(event)
	s.mu.Unlock()
}

func (s *EsportsService) createReplacementEvent(oldEvent *EsportsEvent) {
	var teams1, teams2 []*EsportsTeam

	switch oldEvent.Game {
	case "cs2":
		teams1 = cs2Teams
		teams2 = cs2Teams
	case "lol":
		teams1 = lolTeams
		teams2 = lolTeams
	case "dota2":
		teams1 = dota2Teams
		teams2 = dota2Teams
	case "valorant":
		teams1 = valorantTeams
		teams2 = valorantTeams
	}

	team1 := teams1[rand.Intn(len(teams1))]
	team2 := teams2[rand.Intn(len(teams2))]

	newEvent := s.createEsportsEvent(oldEvent.Game, oldEvent.League, team1, team2, oldEvent.Format)
	s.events[newEvent.ID] = newEvent
}

func (s *EsportsService) settleEsportsBets(eventID, winner string) {
	// In full implementation, settle user bets
}

// ============ BETTING OPERATIONS ============

// PlaceEsportsBet places a bet on an esports event
func (s *EsportsService) PlaceEsportsBet(userID, eventID, selection string, stake float64) (*EsportsBet, error) {
	s.mu.RLock()
	event, ok := s.events[eventID]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	if event.Status != "scheduled" && event.Status != "live" {
		return nil, fmt.Errorf("event is not available for betting")
	}

	// Get odds
	var odds float64
	found := false

	// Check live odds first for live events
	if event.Status == "live" {
		if selection == "team1" {
			odds = event.LiveOdds["team1"]
			found = true
		} else if selection == "team2" {
			odds = event.LiveOdds["team2"]
			found = true
		}
	}

	// Fallback to market odds
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

	bet := &EsportsBet{
		ID:           uuid.New().String(),
		UserID:       userID,
		EventID:      eventID,
		Game:         event.Game,
		BetType:      "map_winner",
		Stake:        stake,
		Odds:         odds,
		Selection:    selection,
		PotentialWin:  stake * odds,
		Status:       "pending",
		PlacedAt:     time.Now(),
	}

	s.recordEsportsBet(bet)

	return bet, nil
}

// GetEvents returns esports events filtered by game and status
func (s *EsportsService) GetEvents(game, status string) []*EsportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*EsportsEvent
	for _, event := range s.events {
		if game != "" && event.Game != game {
			continue
		}
		if status != "" && event.Status != status {
			continue
		}
		result = append(result, event)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime.Before(result[j].StartTime)
	})

	return result
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

// GetLiveEvents returns live events
func (s *EsportsService) GetLiveEvents() []*EsportsEvent {
	return GetEvents("", "live")
}

// GetGames returns available games
func (s *EsportsService) GetGames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	gamesSet := make(map[string]bool)
	for _, event := range s.events {
		gamesSet[event.Game] = true
	}

	games := make([]string, 0, len(gamesSet))
	for game := range gamesSet {
		games = append(games, game)
	}

	return games
}

// GetMarkets returns available markets for a game
func (s *EsportsService) GetMarkets(game string) []string {
	if markets, ok := esportsMarkets[game]; ok {
		result := make([]string, 0, len(markets))
		for market := range markets {
			result = append(result, market)
		}
		return result
	}
	return []string{}
}

func (s *EsportsService) recordEsportsBet(bet *EsportsBet) {
	dbBet := models.Bet{
		UserID:      uuid.MustParse(bet.UserID),
		GameType:    "esports_" + bet.Game,
		BetAmount:   bet.Stake,
		WinAmount:   bet.PotentialWin,
		Multiplier:  bet.Odds,
		Status:      bet.Status,
		GameData:    fmt.Sprintf(`{"event_id":"%s","selection":"%s"}`, bet.EventID, bet.Selection),
	}

	s.db.Create(&dbBet)
}

// PlaceSpecialtyBet places a bet on specialty markets (first blood, etc.)
func (s *EsportsService) PlaceSpecialtyBet(userID, eventID, market, selection string, stake float64) (*EsportsBet, error) {
	s.mu.RLock()
	event, ok := s.events[eventID]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	if event.Status != "live" {
		return nil, fmt.Errorf("event must be live for specialty bets")
	}

	marketData, ok := event.Markets[market]
	if !ok {
		return nil, fmt.Errorf("market not available")
	}

	odds, ok := marketData[selection]
	if !ok {
		return nil, fmt.Errorf("selection not available")
	}

	bet := &EsportsBet{
		ID:           uuid.New().String(),
		UserID:       userID,
		EventID:      eventID,
		Game:         event.Game,
		BetType:      market,
		Stake:        stake,
		Odds:         odds,
		Selection:    selection,
		PotentialWin:  stake * odds,
		Status:       "pending",
		PlacedAt:     time.Now(),
	}

	s.recordEsportsBet(bet)

	return bet, nil
}

type EsportsBet struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	EventID      string    `json:"event_id"`
	Game         string    `json:"game"`
	BetType      string    `json:"bet_type"`
	Stake        float64   `json:"stake"`
	Odds         float64   `json:"odds"`
	Selection    string    `json:"selection"`
	PotentialWin float64   `json:"potential_win"`
	Status       string    `json:"status"`
	PlacedAt     time.Time `json:"placed_at"`
}

// GenerateSeeds for provably fair
func (s *EsportsService) GenerateSeeds() (string, string, error) {
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
