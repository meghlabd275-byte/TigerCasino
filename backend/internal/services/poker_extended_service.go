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

// PokerService handles all poker game operations
type PokerService struct {
	db           *gorm.DB
	mu           sync.RWMutex
	tables       map[string]*PokerTable
	players      map[string]*PokerPlayer
	tournaments  map[string]*PokerTournament
}

type Card struct {
	Suit  string `json:"suit"`
	Rank  string `json:"rank"`
	Value int    `json:"value"`
}

type PokerTable struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	GameType      string              `json:"game_type"` // texas_holdem, omaha, seven_card
	StakeLevel    string              `json:"stake_level"` // $0.01/$0.02, $0.05/$0.10, etc.
	MinPlayers    int                 `json:"min_players"`
	MaxPlayers    int                 `json:"max_players"`
	Players       []*PokerPlayer     `json:"players"`
	Dealer        int                 `json:"dealer"`
	Pot           float64             `json:"pot"`
	SidePots      []float64           `json:"side_pots"`
	CommunityCards []Card              `json:"community_cards"`
	CurrentPlayer int                 `json:"current_player"`
	Stage         string              `json:"stage"` // preflop, flop, turn, river, showdown
	Bets          map[string]float64  `json:"bets"`
	CurrentBet    float64             `json:"current_bet"`
	MaxBuyIn      float64             `json:"max_buy_in"`
	MinBuyIn      float64             `json:"min_buy_in"`
	Rake          float64             `json:"rake"` // 5% typically
	Flop          []Card               `json:"-"`
	Turn          Card                `json:"-"`
	River         Card                `json:"-"`
	Started       bool                `json:"started"`
	CreatedAt     time.Time           `json:"created_at"`
}

type PokerPlayer struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Chips        float64   `json:"chips"`
	CurrentBet   float64   `json:"current_bet"`
	Hand         []Card    `json:"hand"`
	Folded       bool      `json:"folded"`
	AllIn        bool      `json:"all_in"`
	LastAction   string    `json:"last_action"`
	Seat         int       `json:"seat"`
	Connected    bool      `json:"connected"`
}

type PokerTournament struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	GameType       string            `json:"game_type"`
	BuyIn          float64           `json:"buy_in"`
	StartingChips  int               `json:"starting_chips"`
	MaxPlayers     int               `json:"max_players"`
	Registered     int               `json:"registered"`
	BlindLevel     int               `json:"blind_level"`
	BlindTimer     time.Duration     `json:"blind_timer"`
	LevelDuration  time.Duration     `json:"level_duration"`
	Status         string            `json:"status"` // registered, running, finished
	Players        []*TournamentPlayer `json:"players"`
	Tables         []*PokerTable      `json:"tables"`
	PrizePool      float64           `json:"prize_pool"`
	StartedAt      *time.Time        `json:"started_at"`
	EndsAt         *time.Time        `json:"ends_at"`
}

type TournamentPlayer struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Chips         int       `json:"chips"`
	Position      int       `json:"position"`
	Eliminated    bool      `json:"eliminated"`
	EliminatedAt   *time.Time `json:"eliminated_at"`
	PrizeWon      float64   `json:"prize_won"`
}

type PokerHandResult struct {
	TableID      string           `json:"table_id"`
	Winners      []PokerPlayer    `json:"winners"`
	HandType     string           `json:"hand_type"`
	HandCards    []Card           `json:"hand_cards"`
	Pot          float64          `json:"pot"`
	Rake         float64          `json:"rake"`
	CommunityCards []Card         `json:"community_cards"`
	ServerSeed   string           `json:"server_seed"`
	ClientSeed   string           `json:"client_seed"`
}

var suits = []string{"♠", "♥", "♦", "♣"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var pokerHandRanks = map[string]int{
	"royal_flush":     10,
	"straight_flush":  9,
	"four_of_a_kind":  8,
	"full_house":      7,
	"flush":           6,
	"straight":        5,
	"three_of_a_kind": 4,
	"two_pair":        3,
	"one_pair":        2,
	"high_card":       1,
}

func NewPokerService(db *gorm.DB) *PokerService {
	s := &PokerService{
		db:          db,
		tables:      make(map[string]*PokerTable),
		players:     make(map[string]*PokerPlayer),
		tournaments: make(map[string]*PokerTournament),
	}
	s.initializeDefaultTables()
	return s
}

func (s *PokerService) initializeDefaultTables() {
	// Texas Hold'em tables
	tableConfigs := []struct {
		name       string
		gameType   string
		stakes     string
		maxPlayers int
		minBuyIn   float64
		maxBuyIn   float64
	}{
		{"Texas Hold'em NL $0.01/$0.02", "texas_holdem", "$0.01/$0.02", 6, 2, 200},
		{"Texas Hold'em NL $0.05/$0.10", "texas_holdem", "$0.05/$0.10", 6, 10, 1000},
		{"Texas Hold'em NL $0.25/$0.50", "texas_holdem", "$0.25/$0.50", 6, 50, 5000},
		{"Texas Hold'em NL $1/$2", "texas_holdem", "$1/$2", 6, 200, 20000},
		{"Texas Hold'em NL $2/$5", "texas_holdem", "$2/$5", 6, 500, 50000},
		{"Omaha PL $0.05/$0.10", "omaha", "$0.05/$0.10", 6, 10, 1000},
		{"Omaha PL $0.25/$0.50", "omaha", "$0.25/$0.50", 6, 50, 5000},
	}

	for _, cfg := range tableConfigs {
		table := &PokerTable{
			ID:           uuid.New().String(),
			Name:         cfg.name,
			GameType:     cfg.gameType,
			StakeLevel:   cfg.stakes,
			MaxPlayers:   cfg.maxPlayers,
			MinPlayers:   2,
			Players:      make([]*PokerPlayer, 0, cfg.maxPlayers),
			Bets:         make(map[string]float64),
			MinBuyIn:     cfg.minBuyIn,
			MaxBuyIn:     cfg.maxBuyIn,
			Rake:         0.05,
			Stage:        "waiting",
			CreatedAt:    time.Now(),
		}
		s.tables[table.ID] = table
	}
}

// ============ Table Management ============

// GetTables returns all available tables
func (s *PokerService) GetTables(gameType string, stakeLevel string) []*PokerTable {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*PokerTable
	for _, table := range s.tables {
		if gameType != "" && table.GameType != gameType {
			continue
		}
		if stakeLevel != "" && table.StakeLevel != stakeLevel {
			continue
		}
		result = append(result, table)
	}

	return result
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

// JoinTable joins a player to a table
func (s *PokerService) JoinTable(tableID string, userID string, username string, buyIn float64) (*PokerTable, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.tables[tableID]
	if !ok {
		return nil, fmt.Errorf("table not found")
	}

	if len(table.Players) >= table.MaxPlayers {
		return nil, fmt.Errorf("table is full")
	}

	if buyIn < table.MinBuyIn || buyIn > table.MaxBuyIn {
		return nil, fmt.Errorf("buy-in must be between %.2f and %.2f", table.MinBuyIn, table.MaxBuyIn)
	}

	// Check if player already at table
	for _, p := range table.Players {
		if p.UserID == userID {
			return nil, fmt.Errorf("already seated at this table")
		}
	}

	// Find empty seat
	seatsTaken := make(map[int]bool)
	for _, p := range table.Players {
		seatsTaken[p.Seat] = true
	}

	seat := 1
	for {
		if !seatsTaken[seat] {
			break
		}
		seat++
	}

	player := &PokerPlayer{
		ID:        uuid.New().String(),
		UserID:    userID,
		Username:  username,
		Chips:     buyIn,
		Seat:      seat,
		Connected: true,
	}

	table.Players = append(table.Players, player)
	s.players[player.ID] = player

	// Start game if enough players and not started
	if len(table.Players) >= table.MinPlayers && !table.Started {
		s.startHand(table)
	}

	return table, nil
}

// LeaveTable removes a player from a table
func (s *PokerService) LeaveTable(tableID string, userID string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.tables[tableID]
	if !ok {
		return 0, fmt.Errorf("table not found")
	}

	for i, p := range table.Players {
		if p.UserID == userID {
			if p.Folded {
				// Player already folded, can leave immediately
				refund := p.Chips + p.CurrentBet
				table.Players = append(table.Players[:i], table.Players[i+1:]...)
				return refund, nil
			}
			// Auto-fold if in hand
			p.Folded = true
			refund := p.Chips
			table.Players = append(table.Players[:i], table.Players[i+1:]...)
			return refund, nil
		}
	}

	return 0, fmt.Errorf("player not found at table")
}

// ============ Hand Management ============

func (s *PokerService) startHand(table *PokerTable) {
	// Reset for new hand
	table.CommunityCards = nil
	table.Pot = 0
	table.SidePots = nil
	table.CurrentBet = 0
	table.Bets = make(map[string]float64)
	table.Stage = "preflop"
	table.Started = true

	// Rotate dealer
	if table.Dealer >= len(table.Players)-1 {
		table.Dealer = 0
	} else {
		table.Dealer++
	}

	// Post blinds
	smallBlindSeat := (table.Dealer + 1) % len(table.Players)
	bigBlindSeat := (table.Dealer + 2) % len(table.Players)

	// Determine blind amounts from stake level
	parts := strings.Split(table.StakeLevel, "/")
	var sb, bb float64
	if len(parts) == 2 {
		sbStr := strings.Replace(parts[0], "$", "", 1)
		bbStr := strings.Replace(parts[1], "$", "", 1)
		sb, _ = strconv.ParseFloat(sbStr, 64)
		bb, _ = strconv.ParseFloat(bbStr, 64)
	} else {
		sb = 0.01
		bb = 0.02
	}

	// Post small blind
	if len(table.Players) > smallBlindSeat {
		table.Players[smallBlindSeat].CurrentBet = sb
		table.Players[smallBlindSeat].Chips -= sb
		table.Pot += sb
	}

	// Post big blind
	if len(table.Players) > bigBlindSeat {
		table.Players[bigBlindSeat].CurrentBet = bb
		table.Players[bigBlindSeat].Chips -= bb
		table.Pot += bb
		table.CurrentBet = bb
	}

	// Deal hole cards
	deck := createPokerDeck()
	for _, p := range table.Players {
		if !p.Folded {
			card1, ok := deck.Draw()
			if ok {
				p.Hand = append(p.Hand, card1)
			}
			card2, ok := deck.Draw()
			if ok {
				p.Hand = append(p.Hand, card2)
			}
		}
	}

	// Store deck for later streets
	table.Flop = make([]Card, 0)
}

func (s *PokerService) dealFlop(table *PokerTable) {
	deck := createPokerDeck()
	// Burn one card
	deck.Draw()
	// Deal three cards
	for i := 0; i < 3; i++ {
		card, ok := deck.Draw()
		if ok {
			table.CommunityCards = append(table.CommunityCards, card)
			table.Flop = append(table.Flop, card)
		}
	}
}

func (s *PokerService) dealTurn(table *PokerTable) {
	deck := createPokerDeck()
	card, ok := deck.Draw()
	if ok {
		table.CommunityCards = append(table.CommunityCards, card)
		table.Turn = card
	}
}

func (s *PokerService) dealRiver(table *PokerTable) {
	deck := createPokerDeck()
	card, ok := deck.Draw()
	if ok {
		table.CommunityCards = append(table.CommunityCards, card)
		table.River = card
	}
}

// ============ Player Actions ============

type PlayerAction struct {
	TableID  string  `json:"table_id"`
	PlayerID string  `json:"player_id"`
	Action   string  `json:"action"` // fold, check, call, bet, raise, all_in
	Amount   float64 `json:"amount"`
}

// PerformAction processes a player's action
func (s *PokerService) PerformAction(action PlayerAction) (*PokerTable, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.tables[action.TableID]
	if !ok {
		return nil, fmt.Errorf("table not found")
	}

	var player *PokerPlayer
	for _, p := range table.Players {
		if p.ID == action.PlayerID {
			player = p
			break
		}
	}

	if player == nil {
		return nil, fmt.Errorf("player not found")
	}

	if player.Folded || player.AllIn {
		return nil, fmt.Errorf("player cannot act")
	}

	switch action.Action {
	case "fold":
		player.Folded = true
		player.LastAction = "fold"

	case "check":
		player.LastAction = "check"

	case "call":
		callAmount := table.CurrentBet - player.CurrentBet
		if callAmount > player.Chips {
			// All-in
			player.AllIn = true
			player.CurrentBet += player.Chips
			table.Pot += player.Chips
			player.Chips = 0
		} else {
			player.Chips -= callAmount
			player.CurrentBet += callAmount
			table.Pot += callAmount
		}
		player.LastAction = "call"

	case "bet":
		if action.Amount <= 0 {
			return nil, fmt.Errorf("bet amount must be positive")
		}
		if action.Amount > player.Chips {
			return nil, fmt.Errorf("insufficient chips")
		}
		player.Chips -= action.Amount
		player.CurrentBet += action.Amount
		table.CurrentBet = player.CurrentBet
		table.Pot += action.Amount
		player.LastAction = "bet"

	case "raise":
		raiseAmount := action.Amount
		if raiseAmount < table.CurrentBet*2 {
			raiseAmount = table.CurrentBet * 2
		}
		if raiseAmount > player.Chips+player.CurrentBet {
			return nil, fmt.Errorf("insufficient chips")
		}
		totalBet := raiseAmount
		player.Chips -= (totalBet - player.CurrentBet)
		player.CurrentBet = totalBet
		table.CurrentBet = totalBet
		table.Pot += (totalBet - player.CurrentBet)
		player.LastAction = "raise"

	case "all_in":
		player.AllIn = true
		table.Pot += player.Chips
		player.CurrentBet += player.Chips
		if player.CurrentBet > table.CurrentBet {
			table.CurrentBet = player.CurrentBet
		}
		player.Chips = 0
		player.LastAction = "all_in"
	}

	// Move to next player
	s.advancePlayer(table)

	// Check if we can proceed to next street
	if s.canProceedToNextStreet(table) {
		s.proceedToNextStreet(table)
	}

	return table, nil
}

func (s *PokerService) canProceedToNextStreet(table *PokerTable) bool {
	activePlayers := 0
	var lastBetter string
	betAmount := table.CurrentBet

	for _, p := range table.Players {
		if !p.Folded && !p.AllIn {
			activePlayers++
		}
		if !p.Folded && p.CurrentBet > 0 {
			lastBetter = p.ID
		}
	}

	if activePlayers <= 1 {
		return true
	}

	// All bets matched
	allMatched := true
	for _, p := range table.Players {
		if !p.Folded && !p.AllIn && p.CurrentBet != betAmount {
			allMatched = false
			break
		}
	}

	return allMatched && lastBetter != ""
}

func (s *PokerService) proceedToNextStreet(table *PokerTable) {
	switch table.Stage {
	case "preflop":
		table.Stage = "flop"
		s.dealFlop(table)
	case "flop":
		table.Stage = "turn"
		s.dealTurn(table)
	case "turn":
		table.Stage = "river"
		s.dealRiver(table)
	case "river":
		table.Stage = "showdown"
		s.resolveHand(table)
	case "showdown":
		// Start new hand
		s.startHand(table)
	}

	// Reset bets for new street
	table.CurrentBet = 0
	for _, p := range table.Players {
		p.CurrentBet = 0
	}
}

func (s *PokerService) advancePlayer(table *PokerTable) {
	// Find next active player
	startSeat := table.CurrentPlayer
	for {
		table.CurrentPlayer++
		if table.CurrentPlayer >= len(table.Players) {
			table.CurrentPlayer = 0
		}

		p := table.Players[table.CurrentPlayer]
		if !p.Folded && !p.AllIn {
			break
		}

		if table.CurrentPlayer == startSeat {
			break
		}
	}
}

func (s *PokerService) resolveHand(table *PokerTable) {
	// Find best hands
	bestHands := s.evaluateHands(table)

	// Calculate pot and side pots
	totalPot := table.Pot

	// Award pot to winner(s)
	winners := bestHands[0].players
	winnings := totalPot * (1 - table.Rake) / float64(len(winners))

	for _, p := range winners {
		p.Chips += winnings
	}

	// Reset for next hand
	table.Started = false
	table.Stage = "waiting"
}

// ============ Hand Evaluation ============

type evaluatedHand struct {
	players    []*PokerPlayer
	handType   string
	handCards  []Card
	mainValue  int
	highCards  []int
}

func (s *PokerService) evaluateHands(table *PokerTable) []evaluatedHand {
	var results []evaluatedHand

	for _, p := range table.Players {
		if p.Folded {
			continue
		}

		allCards := append(p.Hand, table.CommunityCards...)
		handType, handCards := evaluatePokerHand(allCards)

		results = append(results, evaluatedHand{
			players:    []*PokerPlayer{p},
			handType:   handType,
			handCards:  handCards,
			mainValue:  pokerHandRanks[handType],
			highCards:  extractHighCards(handCards),
		})
	}

	// Sort by hand strength
	sort.Slice(results, func(i, j int) bool {
		if results[i].mainValue != results[j].mainValue {
			return results[i].mainValue > results[j].mainValue
		}
		return compareHighCards(results[i].highCards, results[j].highCards) > 0
	})

	return results
}

func evaluatePokerHand(cards []Card) (string, []Card) {
	// Check for flush
	suitsCount := make(map[string]int)
	for _, c := range cards {
		suitsCount[c.Suit]++
	}

	var flushCards []Card
	var flushSuit string
	for suit, count := range suitsCount {
		if count >= 5 {
			flushSuit = suit
			for _, c := range cards {
				if c.Suit == suit {
					flushCards = append(flushCards, c)
				}
			}
			break
		}
	}

	// Check for straight
	sortedCards := make([]Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Value > sortedCards[j].Value
	})

	var straightCards []Card
	uniqueValues := make(map[int]bool)
	for _, c := range sortedCards {
		if !uniqueValues[c.Value] {
			uniqueValues[c.Value] = true
			straightCards = append(straightCards, c)
		}
	}

	// Check straight flush
	if len(flushCards) >= 5 {
		straightFlush := checkStraight(flushCards)
		if straightFlush != nil {
			if straightFlush[0].Value == 14 {
				return "royal_flush", straightFlush
			}
			return "straight_flush", straightFlush
		}
	}

	// Check four of a kind
	valueCount := make(map[int]int)
	for _, c := range cards {
		valueCount[c.Value]++
	}

	for value, count := range valueCount {
		if count == 4 {
			// Find kicker
			var kickers []Card
			for _, c := range cards {
				if c.Value != value {
					kickers = append(kickers, c)
				}
			}
			sort.Slice(kickers, func(i, j int) bool {
				return kickers[i].Value > kickers[j].Value
			})
			quads := []Card{{Value: value}, {Value: value}, {Value: value}, {Value: value}, kickers[0]}
			return "four_of_a_kind", quads
		}
	}

	// Check full house
	var trips, pair []int
	for value, count := range valueCount {
		if count == 3 {
			trips = append(trips, value)
		}
		if count >= 2 {
			pair = append(pair, value)
		}
	}

	if len(trips) >= 1 {
		sort.Ints(trips)
		if len(trips) >= 2 || len(pair) >= 1 {
			var fullHouse []Card
			// Add trips
			for i := 0; i < 3; i++ {
				fullHouse = append(fullHouse, Card{Value: trips[len(trips)-1]})
			}
			// Add best pair
			if len(trips) >= 2 {
				for i := 0; i < 2; i++ {
					fullHouse = append(fullHouse, Card{Value: trips[len(trips)-2]})
				}
			} else {
				sort.Ints(pair)
				for i := 0; i < 2; i++ {
					fullHouse = append(fullHouse, Card{Value: pair[len(pair)-1]})
				}
			}
			return "full_house", fullHouse
		}
	}

	// Check flush
	if len(flushCards) >= 5 {
		sort.Slice(flushCards, func(i, j int) bool {
			return flushCards[i].Value > flushCards[j].Value
		})
		return "flush", flushCards[:5]
	}

	// Check straight
	if len(straightCards) >= 5 {
		straight := checkStraight(straightCards)
		if straight != nil {
			return "straight", straight
		}
	}

	// Check three of a kind
	if len(trips) >= 1 {
		var kickers []Card
		for _, c := range cards {
			if c.Value != trips[len(trips)-1] {
				kickers = append(kickers, c)
			}
		}
		sort.Slice(kickers, func(i, j int) bool {
			return kickers[i].Value > kickers[j].Value
		})
		tripsHand := []Card{
			{Value: trips[len(trips)-1]}, {Value: trips[len(trips)-1]}, {Value: trips[len(trips)-1]},
			kickers[0], kickers[1],
		}
		return "three_of_a_kind", tripsHand
	}

	// Check two pair
	if len(pair) >= 2 {
		sort.Ints(pair)
		var kicker Card
		for _, c := range cards {
			if c.Value != pair[len(pair)-1] && c.Value != pair[len(pair)-2] {
				kicker = c
				break
			}
		}
		twoPair := []Card{
			{Value: pair[len(pair)-1]}, {Value: pair[len(pair)-1]},
			{Value: pair[len(pair)-2]}, {Value: pair[len(pair)-2]},
			kicker,
		}
		return "two_pair", twoPair
	}

	// Check one pair
	if len(pair) >= 1 {
		var kickers []Card
		for _, c := range cards {
			if c.Value != pair[len(pair)-1] {
				kickers = append(kickers, c)
			}
		}
		sort.Slice(kickers, func(i, j int) bool {
			return kickers[i].Value > kickers[j].Value
		})
		pairHand := []Card{
			{Value: pair[len(pair)-1]}, {Value: pair[len(pair)-1]},
			kickers[0], kickers[1], kickers[2],
		}
		return "one_pair", pairHand
	}

	// High card
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Value > sortedCards[j].Value
	})
	return "high_card", sortedCards[:5]
}

func checkStraight(cards []Card) []Card {
	if len(cards) < 5 {
		return nil
	}

	// Check for Ace-low straight (A-2-3-4-5)
	hasAce := false
	hasFive := false
	for _, c := range cards {
		if c.Value == 14 {
			hasAce = true
		}
		if c.Value == 5 {
			hasFive = true
		}
	}

	var values []int
	if hasAce && hasFive {
		// Create A-2-3-4-5 straight
		return []Card{{Value: 5}, {Value: 4}, {Value: 3}, {Value: 2}, {Value: 1}}
	}

	// Check for regular straight
	for i := 0; i <= len(cards)-5; i++ {
		straight := true
		for j := 0; j < 4; j++ {
			if cards[i+j].Value != cards[i+j+1].Value+1 {
				straight = false
				break
			}
		}
		if straight {
			return cards[i : i+5]
		}
	}

	return nil
}

func extractHighCards(cards []Card) []int {
	var values []int
	for _, c := range cards {
		values = append(values, c.Value)
	}
	return values
}

func compareHighCards(a, b []int) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] > b[i] {
			return 1
		}
		if a[i] < b[i] {
			return -1
		}
	}
	return 0
}

// ============ Deck Management ============

type PokerDeck struct {
	Cards []Card
}

func createPokerDeck() *PokerDeck {
	deck := &PokerDeck{
		Cards: make([]Card, 0, 52),
	}

	for _, suit := range suits {
		for i, rank := range ranks {
			value := i + 2
			if rank == "A" {
				value = 14
			} else if rank == "K" {
				value = 13
			} else if rank == "Q" {
				value = 12
			} else if rank == "J" {
				value = 11
			}

			deck.Cards = append(deck.Cards, Card{
				Suit:  suit,
				Rank:  rank,
				Value: value,
			})
		}
	}

	// Shuffle
	rand.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})

	return deck
}

func (d *PokerDeck) Draw() (Card, bool) {
	if len(d.Cards) == 0 {
		return Card{}, false
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card, true
}

// ============ Tournament Management ============

// CreateTournament creates a new poker tournament
func (s *PokerService) CreateTournament(name string, gameType string, buyIn float64, maxPlayers int, startingChips int) *PokerTournament {
	tournament := &PokerTournament{
		ID:              uuid.New().String(),
		Name:            name,
		GameType:        gameType,
		BuyIn:           buyIn,
		StartingChips:   startingChips,
		MaxPlayers:      maxPlayers,
		Registered:      0,
		BlindLevel:      1,
		BlindTimer:      10 * time.Minute,
		LevelDuration:   10 * time.Minute,
		Status:          "registered",
		Players:         make([]*TournamentPlayer, 0),
		PrizePool:       0,
	}

	s.tournaments[tournament.ID] = tournament
	return tournament
}

// RegisterForTournament registers a player for a tournament
func (s *PokerService) RegisterForTournament(tournamentID string, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tournament, ok := s.tournaments[tournamentID]
	if !ok {
		return fmt.Errorf("tournament not found")
	}

	if tournament.Status != "registered" {
		return fmt.Errorf("tournament is not accepting registrations")
	}

	if tournament.Registered >= tournament.MaxPlayers {
		return fmt.Errorf("tournament is full")
	}

	player := &TournamentPlayer{
		ID:           uuid.New().String(),
		UserID:       userID,
		Chips:        tournament.StartingChips,
		Position:     tournament.Registered + 1,
		Eliminated:   false,
	}

	tournament.Players = append(tournament.Players, player)
	tournament.Registered++
	tournament.PrizePool += tournament.BuyIn

	return nil
}

// GetTournaments returns available tournaments
func (s *PokerService) GetTournaments(status string) []*PokerTournament {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*PokerTournament
	for _, t := range s.tournaments {
		if status != "" && t.Status != status {
			continue
		}
		result = append(result, t)
	}

	return result
}

// Helper for parsing
func strconv.ParseFloat(s string, bitSize int) (float64, error) {
	parsed, err := fmt.Sscanf(s, "%f", &parsed)
	if err != nil {
		return 0, err
	}
	if parsed != 1 {
		return 0, fmt.Errorf("ParseFloat: %s", s)
	}
	return parsed, nil
}
