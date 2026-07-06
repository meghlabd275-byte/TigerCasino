package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// PokerService handles poker room operations
type PokerService struct {
	mu           sync.RWMutex
	tables       map[string]*PokerTable
	players      map[string]*PokerPlayer
	tournaments  map[string]*PokerTournament
}

// PokerTable represents a poker table
type PokerTable struct {
	ID           string
	Name         string
	GameType     string // Texas Hold'em, Omaha, 7-Card Stud
	StakeLevel   string // NL, PL, FL
	MinBuyIn     float64
	MaxBuyIn     float64
	SmallBlind   float64
	BigBlind     float64
	MaxPlayers   int
	CurrentPlayers int
	Pot          float64
	CommunityCards []string
	Dealer       int
	Status       string // waiting, playing, paused
	Players      []*PokerPlayer
	CurrentTurn  int
	Street       string // preflop, flop, turn, river, showdown
}

// PokerPlayer represents a player at the table
type PokerPlayer struct {
	ID        string
	Username  string
	Balance   float64
	CurrentBet float64
	Cards     []string
	Folded    bool
	AllIn     bool
	Seat      int
}

// PokerTournament represents a poker tournament
type PokerTournament struct {
	ID          string
	Name        string
	GameType    string
	BuyIn       float64
	StartingChips float64
	StartTime   time.Time
	MaxPlayers  int
	Registered  int
	Status      string // upcoming, running, completed
	BlindLevel  int
	Levels      []BlindLevel
	Players     []*PokerPlayer
	Tables      []*PokerTable
}

// BlindLevel represents a tournament blind level
type BlindLevel struct {
	Level   int
	SmallBlind float64
	BigBlind  float64
	Duration  int // minutes
}

// NewPokerService creates a new poker service
func NewPokerService() *PokerService {
	s := &PokerService{
		tables:      make(map[string]*PokerTable),
		players:     make(map[string]*PokerPlayer),
		tournaments: make(map[string]*PokerTournament),
	}
	s.initializeDefaultTables()
	return s
}

func (s *PokerService) initializeDefaultTables() {
	// Texas Hold'em Cash Tables
	s.tables["th_nl_001"] = &PokerTable{
		ID: "th_nl_001", Name: "Texas Hold'em NL $1/$2", GameType: "Texas Hold'em",
		StakeLevel: "NL", MinBuyIn: 200, MaxBuyIn: 2000, SmallBlind: 1, BigBlind: 2,
		MaxPlayers: 6, CurrentPlayers: 0, Status: "waiting",
	}

	s.tables["th_nl_002"] = &PokerTable{
		ID: "th_nl_002", Name: "Texas Hold'em NL $2/$5", GameType: "Texas Hold'em",
		StakeLevel: "NL", MinBuyIn: 500, MaxBuyIn: 5000, SmallBlind: 2, BigBlind: 5,
		MaxPlayers: 6, CurrentPlayers: 0, Status: "waiting",
	}

	s.tables["th_nl_003"] = &PokerTable{
		ID: "th_nl_003", Name: "Texas Hold'em NL $5/$10", GameType: "Texas Hold'em",
		StakeLevel: "NL", MinBuyIn: 1000, MaxBuyIn: 10000, SmallBlind: 5, BigBlind: 10,
		MaxPlayers: 6, CurrentPlayers: 0, Status: "waiting",
	}

	s.tables["th_pl_001"] = &PokerTable{
		ID: "th_pl_001", Name: "Texas Hold'em PL $0.50/$1", GameType: "Texas Hold'em",
		StakeLevel: "PL", MinBuyIn: 100, MaxBuyIn: 1000, SmallBlind: 0.5, BigBlind: 1,
		MaxPlayers: 6, CurrentPlayers: 0, Status: "waiting",
	}

	// Omaha Tables
	s.tables["omaha_001"] = &PokerTable{
		ID: "omaha_001", Name: "Omaha PL $1/$2", GameType: "Omaha",
		StakeLevel: "PL", MinBuyIn: 200, MaxBuyIn: 2000, SmallBlind: 1, BigBlind: 2,
		MaxPlayers: 6, CurrentPlayers: 0, Status: "waiting",
	}

	// 7-Card Stud
	s.tables["stud_001"] = &PokerTable{
		ID: "stud_001", Name: "7-Card Stud $1/$2", GameType: "7-Card Stud",
		StakeLevel: "FL", MinBuyIn: 200, MaxBuyIn: 2000, SmallBlind: 1, BigBlind: 2,
		MaxPlayers: 7, CurrentPlayers: 0, Status: "waiting",
	}
}

// GetTables returns all available poker tables
func (s *PokerService) GetTables() []PokerTable {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tables []PokerTable
	for _, t := range s.tables {
		tables = append(tables, *t)
	}
	return tables
}

// GetTable returns a specific table
func (s *PokerService) GetTable(tableID string) (*PokerTable, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	table, ok := s.tables[tableID]
	if !ok {
		return nil, fmt.Errorf("table not found")
	}
	return table, nil
}

// JoinTable allows a player to join a table
func (s *PokerService) JoinTable(tableID, userID, username string, buyIn float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.tables[tableID]
	if !ok {
		return fmt.Errorf("table not found")
	}

	if table.CurrentPlayers >= table.MaxPlayers {
		return fmt.Errorf("table is full")
	}

	if buyIn < table.MinBuyIn || buyIn > table.MaxBuyIn {
		return fmt.Errorf("buy-in must be between %f and %f", table.MinBuyIn, table.MaxBuyIn)
	}

	player := &PokerPlayer{
		ID:        userID,
		Username:  username,
		Balance:   buyIn,
		CurrentBet: 0,
		Folded:    false,
		AllIn:     false,
		Seat:      table.CurrentPlayers,
	}

	table.Players = append(table.Players, player)
	table.CurrentPlayers++
	table.Status = "playing"

	return nil
}

// LeaveTable allows a player to leave a table
func (s *PokerService) LeaveTable(tableID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.tables[tableID]
	if !ok {
		return fmt.Errorf("table not found")
	}

	for i, p := range table.Players {
		if p.ID == userID {
			table.Players = append(table.Players[:i], table.Players[i+1:]...)
			table.CurrentPlayers--
			if table.CurrentPlayers == 0 {
				table.Status = "waiting"
			}
			return nil
		}
	}

	return fmt.Errorf("player not found at table")
}

// CreateTournament creates a new poker tournament
func (s *PokerService) CreateTournament(name, gameType string, buyIn, startingChips float64, maxPlayers int, startTime time.Time) *PokerTournament {
	tournament := &PokerTournament{
		ID:             uuid.New().String(),
		Name:           name,
		GameType:       gameType,
		BuyIn:          buyIn,
		StartingChips:  startingChips,
		StartTime:      startTime,
		MaxPlayers:     maxPlayers,
		Status:         "upcoming",
		BlindLevel:     1,
		Players:        make([]*PokerPlayer, 0),
		Tables:         make([]*PokerTable, 0),
		Levels: []BlindLevel{
			{Level: 1, SmallBlind: 10, BigBlind: 20, Duration: 20},
			{Level: 2, SmallBlind: 20, BigBlind: 40, Duration: 20},
			{Level: 3, SmallBlind: 30, BigBlind: 60, Duration: 20},
			{Level: 4, SmallBlind: 50, BigBlind: 100, Duration: 20},
			{Level: 5, SmallBlind: 100, BigBlind: 200, Duration: 20},
		},
	}

	s.mu.Lock()
	s.tournaments[tournament.ID] = tournament
	s.mu.Unlock()

	return tournament
}

// RegisterForTournament registers a player for a tournament
func (s *PokerService) RegisterForTournament(tournamentID, userID, username string, buyIn float64) error {
	s.mu.RLock()
	tournament, ok := s.tournaments[tournamentID]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("tournament not found")
	}

	if tournament.Registered >= tournament.MaxPlayers {
		return fmt.Errorf("tournament is full")
	}

	player := &PokerPlayer{
		ID:         userID,
		Username:   username,
		Balance:    tournament.StartingChips,
		CurrentBet: 0,
		Folded:     false,
		AllIn:      false,
		Seat:       tournament.Registered,
	}

	s.mu.Lock()
	tournament.Players = append(tournament.Players, player)
	tournament.Registered++
	s.mu.Unlock()

	return nil
}

// GetTournaments returns all tournaments
func (s *PokerService) GetTournaments(status string) []PokerTournament {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tournaments []PokerTournament
	for _, t := range s.tournaments {
		if status != "" && t.Status != status {
			continue
		}
		tournaments = append(tournaments, *t)
	}
	return tournaments
}

// GetTournament returns a specific tournament
func (s *PokerService) GetTournament(tournamentID string) (*PokerTournament, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tournament, ok := s.tournaments[tournamentID]
	if !ok {
		return nil, fmt.Errorf("tournament not found")
	}
	return tournament, nil
}

// Action represents a player action at the table
func (s *PokerService) Action(tableID, userID, action string, amount float64) error {
	s.mu.RLock()
	table, ok := s.tables[tableID]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("table not found")
	}

	// Find player
	var player *PokerPlayer
	for _, p := range table.Players {
		if p.ID == userID {
			player = p
			break
		}
	}

	if player == nil {
		return fmt.Errorf("player not found at table")
	}

	switch action {
	case "fold":
		player.Folded = true
	case "check":
		// Player checks
	case "call":
		player.CurrentBet = table.BigBlind
	case "raise":
		player.CurrentBet = amount
	case "all-in":
		player.AllIn = true
		player.CurrentBet = player.Balance
	default:
		return fmt.Errorf("invalid action")
	}

	return nil
}
