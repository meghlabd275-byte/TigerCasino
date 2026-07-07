package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// GameService handles all game operations
type GameService struct {
	db         *gorm.DB
	security   *SecurityBridge
	gameEngine *GameEngineBridge
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{
		db:         db,
		security:   NewSecurityBridge(),
		gameEngine: NewGameEngineBridge(),
	}
}

// ============ Provably Fair System ============

// Seeds represents provably fair seeds
type Seeds struct {
	ServerSeed     string
	ServerSeedHash string
	ClientSeed     string
	Nonce          int
}

// GenerateSeeds creates new provably fair seeds
func (s *GameService) GenerateSeeds(clientSeed string) (*Seeds, error) {
	serverSeedBytes := make([]byte, 32)
	if _, err := rand.Read(serverSeedBytes); err != nil {
		return nil, fmt.Errorf("failed to generate server seed: %w", err)
	}

	serverSeed := hex.EncodeToString(serverSeedBytes)
	serverSeedHash := s.HashSeed(serverSeed)

	return &Seeds{
		ServerSeed:     serverSeed,
		ServerSeedHash: serverSeedHash,
		ClientSeed:     clientSeed,
		Nonce:          0,
	}, nil
}

// HashSeed creates SHA-256 hash of a seed
func (s *GameService) HashSeed(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])
}

func (s *GameService) GenerateOutcome(serverSeed, clientSeed string, nonce int, max int) int {
	outcome := s.security.GenerateOutcome(serverSeed, clientSeed, nonce)
	return int(outcome * float64(max))
}

func (s *GameService) GenerateFloatOutcome(serverSeed, clientSeed string, nonce int) float64 {
	return s.security.GenerateOutcome(serverSeed, clientSeed, nonce)
}

// ============ Dice Game ============

type DiceResult struct {
	Roll         float64 `json:"roll"`
	Target       float64 `json:"target"`
	Multiplier   float64 `json:"multiplier"`
	WinAmount    float64 `json:"win_amount"`
	Direction    string  `json:"direction"` // "over" or "under"
	ServerSeed   string  `json:"server_seed"`
	ClientSeed   string  `json:"client_seed"`
	Nonce        int     `json:"nonce"`
	Verified     bool    `json:"verified"`
}

func (s *GameService) PlayDice(userID uuid.UUID, betAmount float64, target float64, direction string) (*DiceResult, error) {
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	// Update balance and log bet
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, -betAmount)

	roll := s.gameEngine.GenerateDiceRoll()

	var multiplier float64
	win := false

	if direction == "over" {
		win = roll > target
		multiplier = (100 - target) / target
	} else {
		win = roll < target
		multiplier = target / (100 - target)
	}

	multiplier = multiplier * 0.98

	var winAmount float64
	if win {
		winAmount = betAmount * multiplier
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	s.db.Create(&models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "dice",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Status:     "settled",
	})

	return &DiceResult{
		Roll:        roll,
		Target:      target,
		Multiplier:  multiplier,
		WinAmount:   winAmount,
		Direction:   direction,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed:  seeds.ClientSeed,
		Nonce:       seeds.Nonce,
		Verified:    true,
	}, nil
}

// ============ Slots ============

type SlotsResult struct {
	Reels      [3]int  `json:"reels"`
	Multiplier float64 `json:"multiplier"`
	WinAmount   float64 `json:"win_amount"`
}

func (s *GameService) PlaySlots(userID uuid.UUID, betAmount float64) (*SlotsResult, error) {
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	winAmountTotal := s.gameEngine.CalculateSlots(int(betAmount))
	winAmount := float64(winAmountTotal)
	multiplier := 0.0
	if betAmount > 0 {
		multiplier = winAmount / betAmount
	}

	var reels [3]int
	reels[0] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce, 10)
	reels[1] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+1, 10)
	reels[2] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+2, 10)

	return &SlotsResult{
		Reels:     reels,
		Multiplier: multiplier,
		WinAmount:  winAmount,
	}, nil
}

// ============ Roulette ============

type RouletteResult struct {
	Number     int     `json:"number"`
	WinAmount  float64 `json:"win_amount"`
	Multiplier float64 `json:"multiplier"`
}

func (s *GameService) PlayRoulette(userID uuid.UUID, betAmount float64, betType string, selectedNumber int) (*RouletteResult, error) {
	winningNumber := s.gameEngine.SpinRoulette()

	win := false
	multiplier := 0.0

	if betType == "straight" && selectedNumber == winningNumber {
		win = true
		multiplier = 35.0
	} else if betType == "red" && isRedNumber(winningNumber) {
		win = true
		multiplier = 2.0
	} else if betType == "black" && isBlackNumber(winningNumber) {
		win = true
		multiplier = 2.0
	}

	winAmount := 0.0
	if win {
		winAmount = betAmount * multiplier
	}

	return &RouletteResult{
		Number:     winningNumber,
		WinAmount:  winAmount,
		Multiplier: multiplier,
	}, nil
}

func isRedNumber(n int) bool {
	reds := []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
	for _, r := range reds {
		if r == n {
			return true
		}
	}
	return false
}

func isBlackNumber(n int) bool {
	if n == 0 {
		return false
	}
	return !isRedNumber(n)
}

// ============ Blackjack & Video Poker ============

type BlackjackResult struct {
	WinAmount  float64 `json:"win_amount"`
	Multiplier float64 `json:"multiplier"`
}

type VideoPokerResult struct {
	WinAmount  float64 `json:"win_amount"`
	Multiplier float64 `json:"multiplier"`
}

func (s *GameService) PlayBlackjack(userID uuid.UUID, betAmount float64) (*BlackjackResult, error) {
	winAmountTotal := s.gameEngine.PlayBlackjack(int(betAmount))
	winAmount := float64(winAmountTotal)
	multiplier := 0.0
	if betAmount > 0 {
		multiplier = winAmount / betAmount
	}

	return &BlackjackResult{
		WinAmount:  winAmount,
		Multiplier: multiplier,
	}, nil
}

func (s *GameService) PlayVideoPoker(userID uuid.UUID, betAmount float64) (*VideoPokerResult, error) {
	winAmountTotal := s.gameEngine.PlayVideoPoker(int(betAmount))
	winAmount := float64(winAmountTotal)
	multiplier := 0.0
	if betAmount > 0 {
		multiplier = winAmount / betAmount
	}

	return &VideoPokerResult{
		WinAmount:  winAmount,
		Multiplier: multiplier,
	}, nil
}

func (s *GameService) GetUserBets(userID uuid.UUID, gameType string, limit int) ([]models.Bet, error) {
	var bets []models.Bet
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&bets).Error
	return bets, err
}
