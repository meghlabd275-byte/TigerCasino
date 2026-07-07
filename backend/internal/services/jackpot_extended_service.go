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

// JackpotService handles all jackpot operations
type JackpotService struct {
	db          *gorm.DB
	mu          sync.RWMutex
	jackpots    map[string]*Jackpot
	history     []JackpotWin
	currentRound uint64
}

type Jackpot struct {
	ID          string     `json:"id"`
	Name       string     `json:"name"`
	GameType   string     `json:"game_type"`
	MinBet     float64    `json:"min_bet"`
	Current    float64    `json:"current"`
	SeedAmount float64    `json:"seed_amount"`
	Increment  float64    `json:"increment"` // How much each bet adds
	MaxWin     float64    `json:"max_win"`
	Multipliers []float64 `json:"multipliers"` // For different win tiers
	WinChance  float64   `json:"win_chance"` // 0.001 = 0.1%
	Active     bool       `json:"active"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type JackpotWin struct {
	ID          string    `json:"id"`
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	JackpotID  string    `json:"jackpot_id"`
	JackpotName string   `json:"jackpot_name"`
	WinAmount  float64   `json:"win_amount"`
	Tier       int       `json:"tier"`
	GameType   string    `json:"game_type"`
	Timestamp  time.Time `json:"timestamp"`
}

type JackpotBet struct {
	JackpotID  string  `json:"jackpot_id"`
	UserID     string  `json:"user_id"`
	GameType   string  `json:"game_type"`
	BetAmount  float64 `json:"bet_amount"`
	Eligible   bool    `json:"eligible"`
}

// Jackpot tiers (from highest to lowest)
var jackpotTiers = []string{"mega", "major", "minor", "mini"}

// Default multipliers for each tier
var defaultMultipliers = map[string][]float64{
	"daily":   {1000, 500, 100, 10},
	"hourly":  {500, 250, 50, 5},
	"mini":    {100, 50, 10, 2},
}

func NewJackpotService(db *gorm.DB) *JackpotService {
	s := &JackpotService{
		db:       db,
		jackpots: make(map[string]*Jackpot),
		history:  make([]JackpotWin, 0),
	}

	s.initializeJackpots()

	// Start jackpot rollover timer
	go s.rolloverLoop()

	return s
}

func (s *JackpotService) initializeJackpots() {
	// Daily Jackpot - resets every 24 hours
	s.jackpots["daily"] = &Jackpot{
		ID:          "daily",
		Name:        "Daily Jackpot",
		GameType:    "all",
		MinBet:      0.50,
		Current:     5000.0,
		SeedAmount:  5000.0,
		Increment:   0.01, // 1% of eligible bets
		MaxWin:      100000.0,
		Multipliers: defaultMultipliers["daily"],
		WinChance:   0.001, // 0.1%
		Active:      true,
		UpdatedAt:   time.Now(),
	}

	// Hourly Jackpot - resets every hour
	s.jackpots["hourly"] = &Jackpot{
		ID:          "hourly",
		Name:        "Hourly Jackpot",
		GameType:    "all",
		MinBet:      0.10,
		Current:     1000.0,
		SeedAmount:  1000.0,
		Increment:   0.005, // 0.5% of eligible bets
		MaxWin:      10000.0,
		Multipliers: defaultMultipliers["hourly"],
		WinChance:   0.005, // 0.5%
		Active:      true,
		UpdatedAt:   time.Now(),
	}

	// Mini Jackpot - drops frequently
	s.jackpots["mini"] = &Jackpot{
		ID:          "mini",
		Name:        "Mini Jackpot",
		GameType:    "all",
		MinBet:      0.05,
		Current:     100.0,
		SeedAmount: 100.0,
		Increment:   0.002, // 0.2% of eligible bets
		MaxWin:      1000.0,
		Multipliers: defaultMultipliers["mini"],
		WinChance:   0.01, // 1%
		Active:      true,
		UpdatedAt:   time.Now(),
	}

	// Game-specific jackpots
	s.jackpots["slots_mega"] = &Jackpot{
		ID:          "slots_mega",
		Name:        "Mega Slots Jackpot",
		GameType:    "slots",
		MinBet:      1.0,
		Current:     50000.0,
		SeedAmount:  50000.0,
		Increment:   0.02, // 2% of slot bets
		MaxWin:      500000.0,
		Multipliers: []float64{1000, 500, 200, 50},
		WinChance:   0.0005, // 0.05%
		Active:      true,
		UpdatedAt:   time.Now(),
	}

	s.jackpots["table_jackpot"] = &Jackpot{
		ID:          "table_jackpot",
		Name:        "Table Games Jackpot",
		GameType:    "table",
		MinBet:      5.0,
		Current:     10000.0,
		SeedAmount:  10000.0,
		Increment:   0.01,
		MaxWin:      50000.0,
		Multipliers: []float64{500, 200, 50, 10},
		WinChance:   0.001,
		Active:      true,
		UpdatedAt:   time.Now(),
	}
}

// rolloverLoop handles automatic jackpot rollovers
func (s *JackpotService) rolloverLoop() {
	// Check for rollovers every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()

		for _, jackpot := range s.jackpots {
			// Check if it's time to roll over (simplified - in production, use scheduled times)
			if now.Sub(jackpot.UpdatedAt) > 24*time.Hour && jackpot.ID == "daily" {
				// Trigger roll over - seed the jackpot
				jackpot.Current = jackpot.SeedAmount
			} else if now.Sub(jackpot.UpdatedAt) > time.Hour && jackpot.ID == "hourly" {
				jackpot.Current = jackpot.SeedAmount
			}
		}

		s.mu.Unlock()
	}
}

// ProcessBet processes a bet and contributes to jackpots
func (s *JackpotService) ProcessBet(userID, gameType string, betAmount float64) []JackpotBet {
	s.mu.Lock()
	defer s.mu.Unlock()

	var eligibleBets []JackpotBet

	for _, jackpot := range s.jackpots {
		if !jackpot.Active {
			continue
		}

		// Check if game type is eligible
		if jackpot.GameType != "all" && jackpot.GameType != gameType {
			continue
		}

		// Check minimum bet
		if betAmount < jackpot.MinBet {
			continue
		}

		// Contribute to jackpot
		contribution := betAmount * jackpot.Increment
		jackpot.Current += contribution

		// Check for win
		eligibleBets = append(eligibleBets, JackpotBet{
			JackpotID: jackpot.ID,
			UserID:    userID,
			GameType:  gameType,
			BetAmount: betAmount,
			Eligible:  true,
		})

		// Try to trigger win
		if s.triggerWin(jackpot, userID, gameType) {
			// Reset jackpot after win
			jackpot.Current = jackpot.SeedAmount
		}
	}

	jackpot.UpdatedAt = time.Now()

	return eligibleBets
}

func (s *JackpotService) triggerWin(jackpot *Jackpot, userID, gameType string) bool {
	// Generate random number for win check
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randNum := new(big.Int).SetBytes(randBytes)
	randProb := float64(randNum.Int64()) / float64(1<<63)

	if randProb < jackpot.WinChance {
		// WIN! Calculate tier and amount
		tier := s.determineTier()
		multiplier := jackpot.Multipliers[tier]
		
		// Cap at max win
		winAmount := jackpot.Current * multiplier
		if winAmount > jackpot.MaxWin {
			winAmount = jackpot.MaxWin
		}

		// Record win
		win := JackpotWin{
			ID:           uuid.New().String(),
			UserID:       userID,
			Username:     "Player", // Would lookup from user service
			JackpotID:    jackpot.ID,
			JackpotName:  jackpot.Name,
			WinAmount:    winAmount,
			Tier:         tier,
			GameType:     gameType,
			Timestamp:    time.Now(),
		}

		s.history = append(s.history, win)

		// Keep only last 100 wins
		if len(s.history) > 100 {
			s.history = s.history[len(s.history)-100:]
		}

		return true
	}

	return false
}

func (s *JackpotService) determineTier() int {
	// Random tier determination based on probabilities
	// Mega: 1%, Major: 4%, Minor: 20%, Mini: 75%
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randNum := new(big.Int).SetBytes(randBytes)
	randVal := float64(randNum.Int64()) / float64(1<<63)

	if randVal < 0.01 {
		return 0 // Mega
	} else if randVal < 0.05 {
		return 1 // Major
	} else if randVal < 0.25 {
		return 2 // Minor
	}
	return 3 // Mini
}

// GetJackpots returns all active jackpots
func (s *JackpotService) GetJackpots() []*Jackpot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Jackpot
	for _, jackpot := range s.jackpots {
		if jackpot.Active {
			result = append(result, jackpot)
		}
	}

	// Sort by current amount descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Current > result[j].Current
	})

	return result
}

// GetJackpot returns a specific jackpot
func (s *JackpotService) GetJackpot(jackpotID string) (*Jackpot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jackpot, ok := s.jackpots[jackpotID]
	if !ok {
		return nil, fmt.Errorf("jackpot not found")
	}
	return jackpot, nil
}

// GetHistory returns recent jackpot wins
func (s *JackpotService) GetHistory(limit int) []JackpotWin {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit > 0 && len(s.history) > limit {
		return s.history[len(s.history)-limit:]
	}
	return s.history
}

// GetUserWins returns wins for a specific user
func (s *JackpotService) GetUserWins(userID string) []JackpotWin {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var wins []JackpotWin
	for _, win := range s.history {
		if win.UserID == userID {
			wins = append(wins, win)
		}
	}
	return wins
}

// GetStats returns jackpot statistics
func (s *JackpotService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalWon := 0.0
	winCount := len(s.history)

	for _, win := range s.history {
		totalWon += win.WinAmount
	}

	totalJackpot := 0.0
	for _, jackpot := range s.jackpots {
		totalJackpot += jackpot.Current
	}

	return map[string]interface{}{
		"total_jackpot":   totalJackpot,
		"total_won":      totalWon,
		"win_count":       winCount,
		"average_win":    totalWon / float64(winCount),
		"largest_win":    s.getLargestWin(),
		"jackpot_count":  len(s.jackpots),
	}
}

func (s *JackpotService) getLargestWin() float64 {
	largest := 0.0
	for _, win := range s.history {
		if win.WinAmount > largest {
			largest = win.WinAmount
		}
	}
	return largest
}

// ManualTrigger manually triggers a jackpot win (for testing/admin)
func (s *JackpotService) ManualTrigger(jackpotID, userID string) (*JackpotWin, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jackpot, ok := s.jackpots[jackpotID]
	if !ok {
		return nil, fmt.Errorf("jackpot not found")
	}

	// Calculate win
	tier := s.determineTier()
	multiplier := jackpot.Multipliers[tier]
	winAmount := jackpot.Current * multiplier
	if winAmount > jackpot.MaxWin {
		winAmount = jackpot.MaxWin
	}

	win := JackpotWin{
		ID:           uuid.New().String(),
		UserID:       userID,
		Username:     "Player",
		JackpotID:    jackpotID,
		JackpotName:  jackpot.Name,
		WinAmount:    winAmount,
		Tier:         tier,
		GameType:     "manual",
		Timestamp:    time.Now(),
	}

	s.history = append(s.history, win)
	jackpot.Current = jackpot.SeedAmount

	return &win, nil
}

// SeedAllJackpots seeds all jackpots to their starting amounts
func (s *JackpotService) SeedAllJackpots() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, jackpot := range s.jackpots {
		jackpot.Current = jackpot.SeedAmount
	}
}

// ============ PROGRESSIVE JACKPOTS ============

type ProgressiveJackpot struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	GameType      string    `json:"game_type"`
	SeedAmount    float64   `json:"seed_amount"`
	Current       float64   `json:"current"`
	Increment     float64   `json:"increment"`
	MinBet        float64   `json:"min_bet"`
	MaxBet        float64   `json:"max_bet"`
	LastWin       float64   `json:"last_win"`
	LastWinTime   time.Time `json:"last_win_time"`
	AverageWin    float64   `json:"average_win"`
	WinCount      int       `json:"win_count"`
	Contributions float64   `json:"contributions"`
}

var progressiveJackpots = map[string]*ProgressiveJackpot{
	"megajackpot": {
		ID:         "megajackpot",
		Name:       "Mega Progressive",
		GameType:   "slots",
		SeedAmount: 100000,
		Current:    150000,
		Increment:  0.03,
		MinBet:     2.0,
		MaxBet:     100.0,
	},
	"majorjackpot": {
		ID:         "majorjackpot",
		Name:       "Major Progressive",
		GameType:   "slots",
		SeedAmount: 10000,
		Current:    15000,
		Increment:  0.025,
		MinBet:     1.0,
		MaxBet:     50.0,
	},
	"minorjackpot": {
		ID:         "minorjackpot",
		Name:       "Minor Progressive",
		GameType:   "slots",
		SeedAmount: 1000,
		Current:    1500,
		Increment:  0.02,
		MinBet:     0.5,
		MaxBet:     25.0,
	},
	"minijackpot": {
		ID:         "minijackpot",
		Name:       "Mini Progressive",
		GameType:   "slots",
		SeedAmount: 100,
		Current:    200,
		Increment:  0.015,
		MinBet:     0.1,
		MaxBet:     10.0,
	},
}

// ContributeToProgressive adds to a progressive jackpot
func (s *JackpotService) ContributeToProgressive(jackpotID string, betAmount float64) bool {
	jackpot, ok := progressiveJackpots[jackpotID]
	if !ok {
		return false
	}

	if betAmount < jackpot.MinBet || betAmount > jackpot.MaxBet {
		return false
	}

	// Add contribution
	contribution := betAmount * jackpot.Increment
	jackpot.Current += contribution
	jackpot.Contributions += contribution

	// Check for win (simplified)
	if rand.Float64() < 0.0001 { // Very rare
		return true
	}

	return false
}

// GetProgressiveJackpots returns all progressive jackpots
func (s *JackpotService) GetProgressiveJackpots() map[string]*ProgressiveJackpot {
	return progressiveJackpots
}

// ============ DAILY DROPS ============

type DailyDrop struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	GameType    string    `json:"game_type"`
	WinAmount   float64   `json:"win_amount"`
	LastDrop   time.Time `json:"last_drop"`
	AvgDropTime float64 `json:"avg_drop_time"`
	MinDrop    float64   `json:"min_drop"`
	MaxDrop    float64   `json:"max_drop"`
}

var dailyDrops = map[string]*DailyDrop{
	" Pragmatic Drops": {
		ID:          "pragmatic_drops",
		Name:        "Pragmatic Play Drops",
		GameType:    "slots",
		WinAmount:   0,
		LastDrop:    time.Now().Add(-2 * time.Hour),
		AvgDropTime: 3.0, // hours
		MinDrop:    50,
		MaxDrop:    5000,
	},
	" Evolution Drops": {
		ID:          "evolution_drops",
		Name:        "Evolution Live Drops",
		GameType:    "live",
		WinAmount:   0,
		LastDrop:    time.Now().Add(-1 * time.Hour),
		AvgDropTime: 2.0,
		MinDrop:     100,
		MaxDrop:     10000,
	},
}

// GetDailyDrops returns all daily drop campaigns
func (s *JackpotService) GetDailyDrops() map[string]*DailyDrop {
	return dailyDrops
}

// CheckDropWin checks if a drop should occur
func (s *JackpotService) CheckDropWin(dropID string) (float64, bool) {
	drop, ok := dailyDrops[dropID]
	if !ok {
		return 0, false
	}

	// Check time since last drop
	hoursSince := time.Since(drop.LastDrop).Hours()

	// If enough time has passed, chance increases
	if hoursSince >= drop.AvgDropTime*0.5 {
		// Calculate drop chance (higher if overdue)
		chance := (hoursSince / drop.AvgDropTime) * 0.1
		if chance > 0.5 {
			chance = 0.5
		}

		if rand.Float64() < chance {
			// Trigger drop
			drop.WinAmount = drop.MinDrop + rand.Float64()*(drop.MaxDrop-drop.MinDrop)
			drop.LastDrop = time.Now()
			return drop.WinAmount, true
		}
	}

	return 0, false
}

// GenerateSeeds for provably fair
func (s *JackpotService) GenerateSeeds() (string, string, error) {
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
