package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// ComprehensiveSportsbookService handles all sports betting operations
type ComprehensiveSportsbookService struct {
	db           *gorm.DB
	mu           sync.RWMutex
	events       map[string]*SportsEvent
	bets         map[string]*SportsBet
	oddsCache    map[string]map[string]float64
	cashouts     map[string]*Cashout
	settledBets  []*SportsBet
}

type SportsEvent struct {
	ID              string              `json:"id"`
	Sport           string              `json:"sport"`
	League          string              `json:"league"`
	HomeTeam        string              `json:"home_team"`
	AwayTeam        string              `json:"away_team"`
	StartTime       time.Time           `json:"start_time"`
	Status          string              `json:"status"` // scheduled, live, finished, cancelled
	HomeScore       int                `json:"home_score"`
	AwayScore       int                `json:"away_score"`
	Period          string              `json:"period"`
	Minute          int                `json:"minute"`
	Markets         map[string]map[string]float64 `json:"markets"`
	HomeOdds        float64            `json:"home_odds"`
	DrawOdds        float64            `json:"draw_odds"`
	AwayOdds        float64            `json:"away_odds"`
	LastUpdate      time.Time          `json:"last_update"`
}

type SportsBet struct {
	ID             string              `json:"id"`
	UserID         string              `json:"user_id"`
	EventID        string              `json:"event_id"`
	BetType        string              `json:"bet_type"` // single, parlay, system, teaser
	Stake          float64             `json:"stake"`
	Odds           float64             `json:"odds"`
	PotentialWin   float64             `json:"potential_win"`
	Selection      string              `json:"selection"` // home, draw, away, over, under, etc.
	Status         string              `json:"status"` // pending, won, lost, cancelled, cashed_out
	PlacedAt       time.Time          `json:"placed_at"`
	SettledAt      *time.Time         `json:"settled_at"`
	CashoutAmount  float64             `json:"cashout_amount"`
	CashoutAt      *time.Time         `json:"cashout_at"`
	Metadata       map[string]interface{} `json:"metadata"`
	// Parlay fields
	MultiBetID     string              `json:"multi_bet_id,omitempty"`
	MultiBetOdds   float64             `json:"multi_bet_odds,omitempty"`
}

type Cashout struct {
	ID            string    `json:"id"`
	BetID         string    `json:"bet_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	OriginalOdds  float64   `json:"original_odds"`
	CurrentOdds   float64   `json:"current_odds"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// Betting Markets
var sportsMarkets = map[string]map[string]bool{
	"football": {
		"1x2": true, "double_chance": true, "draw_no_bet": true,
		"over_0.5": true, "over_1.5": true, "over_2.5": true, "over_3.5": true, "over_4.5": true,
		"under_0.5": true, "under_1.5": true, "under_2.5": true, "under_3.5": true, "under_4.5": true,
		"btts": true, "btts_and_over_2.5": true, "win_to_nil": true,
		"first_half_1x2": true, "first_half_over_0.5": true, "first_half_under_1.5": true,
		"corners_over_7.5": true, "corners_over_8.5": true, "corners_over_9.5": true,
		"cards_over_3.5": true, "cards_over_4.5": true,
		"both_teams_to_score_1st_half": true,
	},
	"basketball": {
		"1x2": true, "spread": true, "total": true,
		"over_150.5": true, "over_160.5": true, "over_170.5": true, "over_180.5": true,
		"under_150.5": true, "under_160.5": true, "under_170.5": true, "under_180.5": true,
		"race_to_10": true, "first_halfWinner": true, "second_halfWinner": true,
		"highest_scoring_quarter": true,
	},
	"tennis": {
		"1x2": true, "handicap": true, "total_games": true,
		"over_18.5": true, "over_19.5": true, "over_20.5": true, "over_21.5": true,
		"under_18.5": true, "under_19.5": true, "under_20.5": true, "under_21.5": true,
		"first_set_winner": true, "set_winner": true,
		"correct_score": true,
	},
	"esports": {
		"1x2": true, "map_winner": true, "correct_score": true,
		"over_2.5": true, "under_2.5": true,
		"first_blood": true, "first_tower": true, "first_dragon": true,
	},
}

// NewComprehensiveSportsbookService creates a new sportsbook service
func NewComprehensiveSportsbookService(db *gorm.DB) *ComprehensiveSportsbookService {
	s := &ComprehensiveSportsbookService{
		db:         db,
		events:     make(map[string]*SportsEvent),
		bets:       make(map[string]*SportsBet),
		oddsCache:  make(map[string]map[string]float64),
		cashouts:   make(map[string]*Cashout),
	}
	s.initializeEvents()
	return s
}

func (s *ComprehensiveSportsbookService) initializeEvents() {
	// Football - Premier League
	s.events["fb_001"] = &SportsEvent{
		ID:          "fb_001",
		Sport:       "football",
		League:       "Premier League",
		HomeTeam:    "Manchester City",
		AwayTeam:    "Arsenal",
		StartTime:   time.Now().Add(2 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.75,
		DrawOdds:    3.50,
		AwayOdds:    4.50,
		Markets:     s.generateFootballMarkets(1.75, 3.50, 4.50),
	}

	s.events["fb_002"] = &SportsEvent{
		ID:          "fb_002",
		Sport:       "football",
		League:       "Premier League",
		HomeTeam:    "Liverpool",
		AwayTeam:    "Chelsea",
		StartTime:   time.Now().Add(4 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.85,
		DrawOdds:    3.40,
		AwayOdds:    4.00,
		Markets:     s.generateFootballMarkets(1.85, 3.40, 4.00),
	}

	s.events["fb_003"] = &SportsEvent{
		ID:          "fb_003",
		Sport:       "football",
		League:       "La Liga",
		HomeTeam:    "Real Madrid",
		AwayTeam:    "Barcelona",
		StartTime:   time.Now().Add(6 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    2.30,
		DrawOdds:    3.60,
		AwayOdds:    2.90,
		Markets:     s.generateFootballMarkets(2.30, 3.60, 2.90),
	}

	// Football - Champions League
	s.events["fb_004"] = &SportsEvent{
		ID:          "fb_004",
		Sport:       "football",
		League:       "Champions League",
		HomeTeam:    "Bayern Munich",
		AwayTeam:    "PSG",
		StartTime:   time.Now().Add(8 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    2.00,
		DrawOdds:    3.75,
		AwayOdds:    3.40,
		Markets:     s.generateFootballMarkets(2.00, 3.75, 3.40),
	}

	// Basketball - NBA
	s.events["bb_001"] = &SportsEvent{
		ID:          "bb_001",
		Sport:       "basketball",
		League:       "NBA",
		HomeTeam:    "Lakers",
		AwayTeam:    "Celtics",
		StartTime:   time.Now().Add(3 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.95,
		AwayOdds:    1.90,
		Markets:     s.generateBasketballMarkets(1.95, 1.90),
	}

	s.events["bb_002"] = &SportsEvent{
		ID:          "bb_002",
		Sport:       "basketball",
		League:       "NBA",
		HomeTeam:    "Warriors",
		AwayTeam:    "Heat",
		StartTime:   time.Now().Add(5 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.85,
		AwayOdds:    2.00,
		Markets:     s.generateBasketballMarkets(1.85, 2.00),
	}

	// Tennis - ATP
	s.events["tn_001"] = &SportsEvent{
		ID:          "tn_001",
		Sport:       "tennis",
		League:       "ATP",
		HomeTeam:    "Djokovic",
		AwayTeam:    "Alcaraz",
		StartTime:   time.Now().Add(5 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.70,
		AwayOdds:    2.20,
		Markets:     s.generateTennisMarkets(1.70, 2.20),
	}

	s.events["tn_002"] = &SportsEvent{
		ID:          "tn_002",
		Sport:       "tennis",
		League:       "ATP",
		HomeTeam:    "Sinner",
		AwayTeam:    "Medvedev",
		StartTime:   time.Now().Add(7 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.85,
		AwayOdds:    2.00,
		Markets:     s.generateTennisMarkets(1.85, 2.00),
	}

	// Esports - CS:GO
	s.events["es_001"] = &SportsEvent{
		ID:          "es_001",
		Sport:       "esports",
		League:       "ESL Pro League",
		HomeTeam:    "FaZe Clan",
		AwayTeam:    "Navi",
		StartTime:   time.Now().Add(1 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.80,
		AwayOdds:    2.05,
		Markets:     s.generateEsportsMarkets(1.80, 2.05),
	}

	s.events["es_002"] = &SportsEvent{
		ID:          "es_002",
		Sport:       "esports",
		League:       "LCK",
		HomeTeam:    "T1",
		AwayTeam:    "Gen.G",
		StartTime:   time.Now().Add(4 * time.Hour),
		Status:       "scheduled",
		HomeOdds:    1.75,
		AwayOdds:    2.15,
		Markets:     s.generateEsportsMarkets(1.75, 2.15),
	}

	// More football events for variety
	for i := 5; i <= 20; i++ {
		homeTeam := []string{"Man Utd", "Spurs", "Newcastle", "Villa", "Brighton"}[i%5]
		awayTeam := []string{"Fulham", "Wolves", "Everton", "Forest", "Palace"}[i%5]
		league := []string{"Premier League", "Champions League", "Europa League"}[i%3]

		homeOdds := 1.50 + float64(i)*0.1
		drawOdds := 3.50 + float64(i)*0.05
		awayOdds := 3.00 + float64(i)*0.15

		s.events[fmt.Sprintf("fb_%03d", i)] = &SportsEvent{
			ID:          fmt.Sprintf("fb_%03d", i),
			Sport:       "football",
			League:       league,
			HomeTeam:    homeTeam,
			AwayTeam:    awayTeam,
			StartTime:   time.Now().Add(time.Duration(i) * time.Hour),
			Status:       "scheduled",
			HomeOdds:    homeOdds,
			DrawOdds:    drawOdds,
			AwayOdds:    awayOdds,
			Markets:     s.generateFootballMarkets(homeOdds, drawOdds, awayOdds),
		}
	}
}

func (s *ComprehensiveSportsbookService) generateFootballMarkets(home, draw, away float64) map[string]map[string]float64 {
	return map[string]map[string]float64{
		"1x2":                {"home": home, "draw": draw, "away": away},
		"double_chance":      {"home_draw": 1.20, "home_away": 1.25, "draw_away": 1.40},
		"draw_no_bet":       {"home": home * 0.85, "away": away * 0.85},
		"over_0.5":          {"over": 1.15},
		"over_1.5":          {"over": 1.40},
		"over_2.5":          {"over": 1.90},
		"over_3.5":          {"over": 2.75},
		"over_4.5":          {"over": 4.00},
		"under_0.5":         {"under": 2.50},
		"under_1.5":         {"under": 1.65},
		"under_2.5":         {"under": 1.95},
		"under_3.5":         {"under": 2.85},
		"btts":              {"yes": 1.70, "no": 2.10},
		"btts_and_over_2.5":  {"yes": 2.50, "no": 1.50},
		"first_half_1x2":    {"home": home + 0.15, "draw": 2.20, "away": away + 0.15},
		"first_half_over_0.5": {"over": 1.40},
		"corners_over_7.5":  {"over": 1.75},
		"corners_over_8.5":  {"over": 2.00},
		"cards_over_3.5":    {"over": 1.85},
	}
}

func (s *ComprehensiveSportsbookService) generateBasketballMarkets(home, away float64) map[string]map[string]float64 {
	return map[string]map[string]float64{
		"1x2":           {"home": home, "away": away},
		"spread":        {"home_-5.5": 1.90, "away_+5.5": 1.90},
		"total":         {"over_210.5": 1.90, "under_210.5": 1.90},
		"over_150.5":    {"over": 1.75},
		"over_160.5":    {"over": 1.85},
		"over_170.5":    {"over": 1.95},
		"under_150.5":   {"under": 1.75},
		"under_160.5":   {"under": 1.85},
		"under_170.5":   {"under": 1.95},
	}
}

func (s *ComprehensiveSportsbookService) generateTennisMarkets(home, away float64) map[string]map[string]float64 {
	return map[string]map[string]float64{
		"1x2":              {"home": home, "away": away},
		"handicap":         {"home_-3.5": 1.90, "away_+3.5": 1.90},
		"total_games":      {"over_20.5": 1.90, "under_20.5": 1.90},
		"over_18.5":       {"over": 1.75},
		"over_19.5":       {"over": 1.85},
		"under_18.5":      {"under": 1.75},
		"under_19.5":      {"under": 1.85},
		"first_set_winner": {"home": home + 0.10, "away": away + 0.10},
	}
}

func (s *ComprehensiveSportsbookService) generateEsportsMarkets(home, away float64) map[string]map[string]float64 {
	return map[string]map[string]float64{
		"1x2":          {"home": home, "away": away},
		"map_winner":   {"home": home, "away": away},
		"over_2.5":     {"over": 1.85},
		"under_2.5":    {"under": 1.95},
		"first_blood":  {"home": 1.90, "away": 1.90},
		"first_tower":  {"home": 1.90, "away": 1.90},
		"first_dragon": {"home": 1.90, "away": 1.90},
	}
}

// ============ Betting Operations ============

// PlaceSingleBet places a single bet
func (s *ComprehensiveSportsbookService) PlaceSingleBet(userID, eventID, selection string, stake float64) (*SportsBet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found: %s", eventID)
	}

	if event.Status != "scheduled" {
		return nil, fmt.Errorf("event is not available for betting")
	}

	// Find odds for selection
	var odds float64
	found := false

	// Try to find in markets first
	for marketName, market := range event.Markets {
		if selOdds, ok := market[selection]; ok {
			odds = selOdds
			found = true
			break
		}
	}

	// Fallback to main odds
	if !found {
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
	}

	if !found {
		return nil, fmt.Errorf("selection not available: %s", selection)
	}

	potentialWin := stake * odds

	bet := &SportsBet{
		ID:           uuid.New().String(),
		UserID:       userID,
		EventID:      eventID,
		BetType:      "single",
		Stake:        stake,
		Odds:         odds,
		PotentialWin: potentialWin,
		Selection:    selection,
		Status:       "pending",
		PlacedAt:     time.Now(),
		Metadata: map[string]interface{}{
			"event":    fmt.Sprintf("%s vs %s", event.HomeTeam, event.AwayTeam),
			"league":   event.League,
			"sport":    event.Sport,
		},
	}

	s.bets[bet.ID] = bet

	return bet, nil
}

// PlaceParlayBet places a multi-selection bet (parlay/accumulator)
func (s *ComprehensiveSportsbookService) PlaceParlayBet(userID string, selections []map[string]string, stake float64) (*SportsBet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(selections) < 2 {
		return nil, fmt.Errorf("parlay must have at least 2 selections")
	}

	if len(selections) > 10 {
		return nil, fmt.Errorf("parlay cannot have more than 10 selections")
	}

	multiBetID := uuid.New().String()
	totalOdds := 1.0
	var multiBets []*SportsBet

	for _, sel := range selections {
		eventID := sel["event_id"]
		selection := sel["selection"]

		event, ok := s.events[eventID]
		if !ok {
			return nil, fmt.Errorf("event not found: %s", eventID)
		}

		// Find odds
		var odds float64
		found := false

		for marketName, market := range event.Markets {
			if selOdds, ok := market[selection]; ok {
				odds = selOdds
				found = true
				break
			}
		}

		if !found {
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
		}

		if !found {
			return nil, fmt.Errorf("selection not available: %s for event %s", selection, eventID)
		}

		totalOdds *= odds

		bet := &SportsBet{
			ID:            uuid.New().String(),
			UserID:        userID,
			EventID:       eventID,
			BetType:       "parlay",
			Stake:         stake / float64(len(selections)), // Divide stake equally
			Odds:          odds,
			PotentialWin:  0,
			Selection:     selection,
			Status:        "pending",
			PlacedAt:      time.Now(),
			MultiBetID:    multiBetID,
			MultiBetOdds:  totalOdds,
			Metadata: map[string]interface{}{
				"event":    fmt.Sprintf("%s vs %s", event.HomeTeam, event.AwayTeam),
				"league":   event.League,
				"sport":    event.Sport,
			},
		}

		s.bets[bet.ID] = bet
		multiBets = append(multiBets, bet)
	}

	// Apply parlay bonus
	if len(selections) >= 3 {
		totalOdds *= 1.05 // 5% bonus for 3+
	}
	if len(selections) >= 5 {
		totalOdds *= 1.10 // 10% bonus for 5+
	}

	potentialWin := stake * totalOdds

	// Update all bets with total odds and potential win
	for _, bet := range multiBets {
		bet.MultiBetOdds = totalOdds
		bet.PotentialWin = potentialWin
	}

	return multiBets[0], nil
}

// PlaceSystemBet places a system bet (e.g., 2/3, 3/4)
func (s *ComprehensiveSportsbookService) PlaceSystemBet(userID string, selections []map[string]string, stake float64, systemType string) ([]*SportsBet, error) {
	// System bet breakdown:
	// 2/3 = 3 combinations
	// 3/4 = 4 combinations
	// 4/5 = 5 combinations

	s.mu.Lock()
	defer s.mu.Unlock()

	parts := strings.Split(systemType, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid system type: %s", systemType)
	}

	required, _ := strconv.Atoi(parts[0])
	total := len(selections)

	if required > total || required < 2 {
		return nil, fmt.Errorf("invalid system: %d out of %d", required, total)
	}

	// Calculate combinations
	combinations := factorial(total) / (factorial(required) * factorial(total-required))
	stakePerCombo := stake / float64(combinations)

	var systemBets []*SportsBet
	multiBetID := uuid.New().String()
	totalOdds := 1.0

	// This is a simplified version - in production, generate all combinations
	for i := 0; i < len(selections); i++ {
		eventID := selections[i]["event_id"]
		selection := selections[i]["selection"]

		event, ok := s.events[eventID]
		if !ok {
			continue
		}

		var odds float64
		found := false

		for marketName, market := range event.Markets {
			if selOdds, ok := market[selection]; ok {
				odds = selOdds
				found = true
				break
			}
		}

		if !found {
			switch selection {
			case "home":
				odds = event.HomeOdds
			case "away":
				odds = event.AwayOdds
			case "draw":
				odds = event.DrawOdds
			}
		}

		totalOdds *= odds
	}

	// Place single bet as representative
	bet := &SportsBet{
		ID:            uuid.New().String(),
		UserID:        userID,
		EventID:       selections[0]["event_id"],
		BetType:       "system",
		Stake:         stake,
		Odds:          totalOdds,
		PotentialWin:  stake * totalOdds,
		Selection:     "system",
		Status:        "pending",
		PlacedAt:      time.Now(),
		MultiBetID:    multiBetID,
		Metadata: map[string]interface{}{
			"system_type":   systemType,
			"combinations": combinations,
			"selections":   len(selections),
		},
	}

	s.bets[bet.ID] = bet
	systemBets = append(systemBets, bet)

	return systemBets, nil
}

// ============ Cashout ============

// CalculateCashout calculates the cashout amount for a bet
func (s *ComprehensiveSportsbookService) CalculateCashout(betID string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bet, ok := s.bets[betID]
	if !ok {
		return 0, fmt.Errorf("bet not found")
	}

	if bet.Status != "pending" {
		return 0, fmt.Errorf("bet cannot be cashed out")
	}

	// Get current odds
	event, ok := s.events[bet.EventID]
	if !ok {
		return 0, fmt.Errorf("event not found")
	}

	var currentOdds float64
	switch bet.Selection {
	case "home":
		currentOdds = event.HomeOdds
	case "away":
		currentOdds = event.AwayOdds
	case "draw":
		currentOdds = event.DrawOdds
	default:
		currentOdds = bet.Odds
	}

	// Cashout formula: adjusted by time remaining and odds movement
	elapsedMinutes := time.Since(bet.PlacedAt).Minutes()
	totalMinutes := event.StartTime.Sub(bet.PlacedAt).Minutes()
	remainingRatio := 1.0

	if totalMinutes > 0 {
		remainingRatio = 1.0 - (elapsedMinutes / totalMinutes)
		if remainingRatio < 0.1 {
			remainingRatio = 0.1
		}
	}

	// Cashout factor decreases as time passes
	cashoutFactor := 0.95 - (elapsedMinutes * 0.001)
	if cashoutFactor < 0.5 {
		cashoutFactor = 0.5
	}

	cashoutAmount := bet.Stake * (currentOdds / bet.Odds) * cashoutFactor * remainingRatio

	return cashoutAmount, nil
}

// RequestCashout requests a cashout for a bet
func (s *ComprehensiveSportsbookService) RequestCashout(betID string) (*Cashout, error) {
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

	now := time.Now()

	cashout := &Cashout{
		ID:           uuid.New().String(),
		BetID:        betID,
		UserID:       bet.UserID,
		Amount:       cashoutAmount,
		OriginalOdds: bet.Odds,
		CurrentOdds:  cashoutAmount / bet.Stake,
		Status:       "completed",
		CreatedAt:    now,
		CompletedAt:  &now,
	}

	bet.Status = "cashed_out"
	bet.CashoutAmount = cashoutAmount
	bet.CashoutAt = &now

	s.cashouts[cashout.ID] = cashout

	return cashout, nil
}

// ============ Bet Settlement ============

// SettleBet settles a single bet
func (s *ComprehensiveSportsbookService) SettleBet(betID string, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bet, ok := s.bets[betID]
	if !ok {
		return fmt.Errorf("bet not found")
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet already settled")
	}

	event, ok := s.events[bet.EventID]
	if !ok {
		return fmt.Errorf("event not found")
	}

	// Simplified settlement
	win := false
	switch bet.Selection {
	case "home":
		win = result == "home"
	case "away":
		win = result == "away"
	case "draw":
		win = result == "draw"
	}

	if win {
		bet.Status = "won"
	} else {
		bet.Status = "lost"
	}

	now := time.Now()
	bet.SettledAt = &now

	// Record in database
	s.recordSportsBet(bet)

	return nil
}

// SettleAllBetsForEvent settles all bets for an event
func (s *ComprehensiveSportsbookService) SettleAllBetsForEvent(eventID string, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return fmt.Errorf("event not found")
	}

	event.Status = "finished"

	for _, bet := range s.bets {
		if bet.EventID == eventID && bet.Status == "pending" {
			win := false
			switch bet.Selection {
			case "home":
				win = result == "home"
			case "away":
				win = result == "away"
			case "draw":
				win = result == "draw"
			}

			if win {
				bet.Status = "won"
			} else {
				bet.Status = "lost"
			}

			now := time.Now()
			bet.SettledAt = &now
		}
	}

	return nil
}

// ============ Query Operations ============

// GetEvent returns event by ID
func (s *ComprehensiveSportsbookService) GetEvent(eventID string) (*SportsEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

// GetLiveEvents returns all live events
func (s *ComprehensiveSportsbookService) GetLiveEvents() []SportsEvent {
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

// GetUpcomingEvents returns upcoming events
func (s *ComprehensiveSportsbookService) GetUpcomingEvents(sport string, limit int) []SportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []SportsEvent
	now := time.Now()

	for _, event := range s.events {
		if event.Status != "scheduled" {
			continue
		}
		if sport != "" && event.Sport != sport {
			continue
		}
		if event.StartTime.After(now) {
			events = append(events, *event)
		}
	}

	// Sort by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	if limit > 0 && len(events) > limit {
		events = events[:limit]
	}

	return events
}

// GetEventsBySport returns events grouped by sport
func (s *ComprehensiveSportsbookService) GetEventsBySport() map[string][]SportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	eventsBySport := make(map[string][]SportsEvent)
	now := time.Now()

	for _, event := range s.events {
		if event.Status == "scheduled" && event.StartTime.After(now) {
			eventsBySport[event.Sport] = append(eventsBySport[event.Sport], *event)
		}
	}

	return eventsBySport
}

// GetUserBets returns all bets for a user
func (s *ComprehensiveSportsbookService) GetUserBets(userID string) []SportsBet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bets []SportsBet
	for _, bet := range s.bets {
		if bet.UserID == userID {
			bets = append(bets, *bet)
		}
	}

	// Sort by date
	sort.Slice(bets, func(i, j int) bool {
		return bets[i].PlacedAt.After(bets[j].PlacedAt)
	})

	return bets
}

// GetUserPendingBets returns pending bets for a user
func (s *ComprehensiveSportsbookService) GetUserPendingBets(userID string) []SportsBet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bets []SportsBet
	for _, bet := range s.bets {
		if bet.UserID == userID && bet.Status == "pending" {
			bets = append(bets, *bet)
		}
	}

	return bets
}

// ============ Live Odds Updates ============

// UpdateLiveOdds simulates live odds changes
func (s *ComprehensiveSportsbookService) UpdateLiveOdds(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return fmt.Errorf("event not found")
	}

	// Simulate odds changes
	change := 0.02
	event.HomeOdds *= (1 + change)
	event.AwayOdds *= (1 - change)

	// Keep within reasonable bounds
	if event.HomeOdds > 5.0 {
		event.HomeOdds = 5.0
	}
	if event.AwayOdds > 5.0 {
		event.AwayOdds = 5.0
	}

	event.LastUpdate = time.Now()

	return nil
}

// ============ Helper Functions ============

func (s *ComprehensiveSportsbookService) recordSportsBet(bet *SportsBet) {
	dbBet := models.Bet{
		UserID:       uuid.MustParse(bet.UserID),
		GameType:     "sportsbook",
		BetAmount:    bet.Stake,
		WinAmount:    bet.PotentialWin,
		Multiplier:   bet.Odds,
		Status:       bet.Status,
		GameData:     fmt.Sprintf(`{"event_id":"%s","selection":"%s"}`, bet.EventID, bet.Selection),
		ServerSeed:   "",
		ClientSeed:   "",
	}

	if bet.Status == "won" {
		dbBet.Profit = bet.PotentialWin - bet.Stake
	} else {
		dbBet.Profit = -bet.Stake
	}

	s.db.Create(&dbBet)
}

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// GetAvailableSports returns list of available sports
func (s *ComprehensiveSportsbookService) GetAvailableSports() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sportsSet := make(map[string]bool)
	for _, event := range s.events {
		if event.Status == "scheduled" {
			sportsSet[event.Sport] = true
		}
	}

	sports := make([]string, 0, len(sportsSet))
	for sport := range sportsSet {
		sports = append(sports, sport)
	}

	sort.Strings(sports)
	return sports
}

// GetEventsByLeague returns events for a specific league
func (s *ComprehensiveSportsbookService) GetEventsByLeague(league string) []SportsEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []SportsEvent
	for _, event := range s.events {
		if event.League == league && event.Status == "scheduled" {
			events = append(events, *event)
		}
	}

	return events
}

// PlaceTeaserBet places a teaser bet (basketball/football only)
func (s *ComprehensiveSportsbookService) PlaceTeaserBet(userID string, selections []map[string]string, stake float64, teaserPoints float64) (*SportsBet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Teaser bet: adjust point spread in user's favor
	// 6-point teaser for football, 4-6 for basketball

	if len(selections) < 2 {
		return nil, fmt.Errorf("teaser must have at least 2 selections")
	}

	if len(selections) > 6 {
		return nil, fmt.Errorf("teaser cannot have more than 6 selections")
	}

	multiBetID := uuid.New().String()
	totalOdds := 1.0

	// Teaser odds are fixed based on number of teams
	teaserOdds := map[int]float64{
		2: 1.50,
		3: 1.75,
		4: 2.00,
		5: 2.50,
		6: 3.00,
	}

	odds := teaserOdds[len(selections)]
	totalOdds = odds
	potentialWin := stake * totalOdds

	bet := &SportsBet{
		ID:            uuid.New().String(),
		UserID:        userID,
		EventID:       selections[0]["event_id"],
		BetType:       "teaser",
		Stake:         stake,
		Odds:          totalOdds,
		PotentialWin:  potentialWin,
		Selection:     "teaser",
		Status:        "pending",
		PlacedAt:      time.Now(),
		MultiBetID:    multiBetID,
		Metadata: map[string]interface{}{
			"teaser_points": teaserPoints,
			"selections":    len(selections),
		},
	}

	s.bets[bet.ID] = bet
	return bet, nil
}

// GetBetHistory returns bet history with filters
func (s *ComprehensiveSportsbookService) GetBetHistory(userID string, status string, betType string, limit int) []SportsBet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bets []SportsBet
	for _, bet := range s.bets {
		if bet.UserID != userID {
			continue
		}
		if status != "" && bet.Status != status {
			continue
		}
		if betType != "" && bet.BetType != betType {
			continue
		}
		bets = append(bets, *bet)
	}

	// Sort by date descending
	sort.Slice(bets, func(i, j int) bool {
		return bets[i].PlacedAt.After(bets[j].PlacedAt)
	})

	if limit > 0 && len(bets) > limit {
		bets = bets[:limit]
	}

	return bets
}

// CancelBet cancels a pending bet (if within time limit)
func (s *ComprehensiveSportsbookService) CancelBet(betID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bet, ok := s.bets[betID]
	if !ok {
		return fmt.Errorf("bet not found")
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet cannot be cancelled")
	}

	// Check if within time limit (5 minutes before event)
	event, ok := s.events[bet.EventID]
	if !ok {
		return fmt.Errorf("event not found")
	}

	if time.Now().Add(5 * time.Minute).After(event.StartTime) {
		return fmt.Errorf("too close to event start to cancel")
	}

	bet.Status = "cancelled"

	return nil
}
