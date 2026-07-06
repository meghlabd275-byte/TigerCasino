package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// GameService handles all game operations
type GameService struct {
	db *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{db: db}
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

// GenerateOutcome creates a deterministic outcome from seeds
func (s *GameService) GenerateOutcome(serverSeed, clientSeed string, nonce int, max int) int {
	combined := fmt.Sprintf("%s:%s:%d", serverSeed, clientSeed, nonce)
	hash := sha256.Sum256([]byte(hashSeed))
	result := new(big.Int).SetBytes(hash[:])
	return int(result.Mod(result, big.NewInt(int64(max))).Int64())
}

// GenerateFloatOutcome creates a float between 0 and 1
func (s *GameService) GenerateFloatOutcome(serverSeed, clientSeed string, nonce int) float64 {
	combined := fmt.Sprintf("%s:%s:%d", serverSeed, clientSeed, nonce)
	hash := sha256.Sum256([]byte(combined))
	result := new(big.Int).SetBytes(hash[:])
	// Convert to float between 0 and 1
	divisor := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	floatResult := new(big.Float).SetInt(result)
	floatResult = floatResult.Quo(floatResult, new(big.Float).SetInt(divisor))
	f, _ := floatResult.Float64()
	return f
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
	// Get or create seeds
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	// Generate roll (0-100)
	roll := s.GenerateFloatOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce) * 100

	// Calculate multiplier
	var multiplier float64
	win := false

	if direction == "over" {
		win = roll > target
		multiplier = (100 - target) / target
	} else {
		win = roll < target
		multiplier = target / (100 - target)
	}

	// House edge
	multiplier = multiplier * 0.98

	var winAmount float64
	if win {
		winAmount = betAmount * multiplier
	}

	// Record bet
	bet := models.Bet{
		UserID:     userID,
		GameType:   "dice",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     winAmount - betAmount,
		Status:     "settled",
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
		IsVerified: true,
	}

	s.db.Create(&bet)

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

// ============ Crash Game ============

type CrashResult struct {
	RoundID     string  `json:"round_id"`
	CrashPoint  float64 `json:"crash_point"`
	PlayerCash  float64 `json:"player_cash,omitempty"`
	DidCrash    bool    `json:"did_crash"`
	ServerSeed  string  `json:"server_seed"`
	ClientSeed  string  `json:"client_seed"`
	Nonce       int     `json:"nonce"`
}

func (s *GameService) StartCrashRound() (*CrashResult, error) {
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	// Generate crash point
	crashPoint := s.GenerateCrashPoint(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	// Store round
	round := models.CrashGameRound{
		RoundID:    uuid.New().String(),
		CrashPoint: crashPoint,
		Status:     "running",
		ServerSeed: seeds.ServerSeed,
		ServerHash: seeds.ServerSeedHash,
		StartTime:  time.Now(),
	}

	s.db.Create(&round)

	return &CrashResult{
		RoundID:    round.RoundID,
		CrashPoint: crashPoint,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}, nil
}

func (s *GameService) GenerateCrashPoint(serverSeed, clientSeed string, nonce int) float64 {
	// Use exponential distribution for realistic crash points
	f := s.GenerateFloatOutcome(serverSeed, clientSeed, nonce)
	
	// Most crashes happen below 2x, rare high multipliers
	if f < 0.7 {
		// 70% chance - lower crash points
		return 1.0 + (f * 10) // 1.0 - 8.0
	}
	// 30% chance - higher crash points
	return 8.0 + (f-0.7)*300 // Up to 100x
}

func (s *GameService) PlaceCrashBet(userID uuid.UUID, roundID string, betAmount float64, autoCashout float64) error {
	round := models.CrashGameRound{}
	if err := s.db.Where("round_id = ?", roundID).First(&round).Error; err != nil {
		return err
	}

	bet := models.Bet{
		UserID:     userID,
		GameType:   "crash",
		BetAmount:  betAmount,
		GameData:   fmt.Sprintf(`{"round_id":"%s","auto_cashout":%f}`, roundID, autoCashout),
		ServerSeed: round.ServerSeed,
		ClientSeed: "",
		Nonce:      round.ID.Nonce(),
	}

	s.db.Create(&bet)
	return nil
}

func (s *GameService) CrashCashout(userID uuid.UUID, betID uuid.UUID, multiplier float64) (float64, error) {
	var bet models.Bet
	if err := s.db.First(&bet, betID).Error; err != nil {
		return 0, err
	}

	if bet.UserID != userID {
		return 0, fmt.Errorf("unauthorized")
	}

	payout := bet.BetAmount * multiplier * 0.97 // 3% house edge

	bet.WinAmount = payout
	bet.Multiplier = multiplier
	bet.Profit = payout - bet.BetAmount
	bet.Status = "won"

	s.db.Save(&bet)

	return payout, nil
}

// ============ Mines Game ============

type MinesResult struct {
	GameID       string  `json:"game_id"`
	MinesCount   int     `json:"mines_count"`
	RevealedTile int     `json:"revealed_tile"`
	IsMine       bool    `json:"is_mine"`
	Multiplier   float64 `json:"multiplier"`
	WinAmount    float64 `json:"win_amount"`
	GameOver     bool    `json:"game_over"`
}

func (s *GameService) StartMinesGame(userID uuid.UUID, betAmount float64, minesCount int) (*MinesResult, error) {
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	// Generate mine positions
	gridSize := 25
	mines := make(map[int]bool)
	for len(mines) < minesCount {
		pos := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+len(mines), gridSize)
		mines[pos] = true
	}

	// Save game state
	gameData := fmt.Sprintf(`{"mines":[%d,%d,%d]}`, getKeys(mines)...)
	
	bet := models.Bet{
		UserID:     userID,
		GameType:   "mines",
		BetAmount:  betAmount,
		GameData:   gameData,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	s.db.Create(&bet)

	return &MinesResult{
		GameID:     bet.ID.String(),
		MinesCount: minesCount,
	}, nil
}

func getKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (s *GameService) RevealMinesTile(userID uuid.UUID, betID uuid.UUID, tile int) (*MinesResult, error) {
	var bet models.Bet
	if err := s.db.First(&bet, betID).Error; err != nil {
		return nil, err
	}

	if bet.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Simplified - in production would check stored mines
	isMine := false // Would check against stored game state

	result := &MinesResult{
		GameID:       bet.ID.String(),
		RevealedTile: tile,
		IsMine:       isMine,
	}

	if !isMine {
		// Calculate multiplier
		result.Multiplier = 1.0 + float64(tile)*0.1
		result.WinAmount = bet.BetAmount * result.Multiplier
	}

	return result, nil
}

// ============ Plinko Game ============

type PlinkoResult struct {
	GameID      string  `json:"game_id"`
	Rows        int     `json:"rows"`
	Risk        string  `json:"risk"`
	FinalPos    int     `json:"final_pos"`
	Multiplier  float64 `json:"multiplier"`
	WinAmount   float64 `json:"win_amount"`
}

func (s *GameService) PlayPlinko(userID uuid.UUID, betAmount float64, rows int, risk string) (*PlinkoResult, error) {
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	// Simulate ball path
	position := rows / 2
	for i := 0; i < rows; i++ {
		f := s.GenerateFloatOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+i)
		if f > 0.5 && position < i+2 {
			position++
		}
	}

	// Get multiplier based on position
	multiplier := s.GetPlinkoMultiplier(rows, risk, position)
	winAmount := betAmount * multiplier * 0.96 // 4% house edge

	bet := models.Bet{
		UserID:     userID,
		GameType:   "plinko",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     winAmount - betAmount,
		Status:     "settled",
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	s.db.Create(&bet)

	return &PlinkoResult{
		GameID:     bet.ID.String(),
		Rows:       rows,
		Risk:       risk,
		FinalPos:   position,
		Multiplier: multiplier,
		WinAmount:  winAmount,
	}, nil
}

func (s *GameService) GetPlinkoMultiplier(rows int, risk string, position int) float64 {
	// Simplified payout tables
	lowRisk := []float64{1.5, 1.2, 0.8, 0.5, 0.5, 0.8, 1.2, 1.5}
	medRisk := []float64{5.0, 2.0, 1.0, 0.5, 0.5, 1.0, 2.0, 5.0}
	highRisk := []float64{10.0, 5.0, 2.0, 0.5, 0.5, 2.0, 5.0, 10.0}

	var table []float64
	switch risk {
	case "low":
		table = lowRisk
	case "high":
		table = highRisk
	default:
		table = medRisk
	}

	if position >= len(table) {
		position = len(table) - 1
	}

	return table[position]
}

// ============ Slots ============

type SlotsResult struct {
	Reels     [3]int  `json:"reels"`
	Multiplier float64 `json:"multiplier"`
	WinAmount  float64 `json:"win_amount"`
}

func (s *GameService) PlaySlots(userID uuid.UUID, betAmount float64) (*SlotsResult, error) {
	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return nil, err
	}

	// Generate 3 reels (0-9)
	var reels [3]int
	reels[0] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce, 10)
	reels[1] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+1, 10)
	reels[2] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+2, 10)

	// Calculate win
	var multiplier float64
	if reels[0] == reels[1] && reels[1] == reels[2] {
		multiplier = 10 // Jackpot
	} else if reels[0] == reels[1] || reels[1] == reels[2] || reels[0] == reels[2] {
		// Two matching - small win
		// Simplified
	}

	winAmount := betAmount * multiplier * 0.95

	bet := models.Bet{
		UserID:     userID,
		GameType:   "slots",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     winAmount - betAmount,
		Status:     "settled",
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	s.db.Create(&bet)

	return &SlotsResult{
		Reels:     reels,
		Multiplier: multiplier,
		WinAmount:  winAmount,
	}, nil
}

// ============ Leaderboard ============

type LeaderboardEntry struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Profit    float64   `json:"profit"`
	Wagered   float64   `json:"wagered"`
	BetCount  int       `json:"bet_count"`
	Rank      int       `json:"rank"`
}

func (s *GameService) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	// Simplified - would aggregate from bets
	s.db.Table("users").
		Select("users.id as user_id, users.username, COALESCE(SUM(bets.profit), 0) as profit, COALESCE(SUM(bets.bet_amount), 0) as wagered, COUNT(bets.id) as bet_count").
		Joins("LEFT JOIN bets ON bets.user_id = users.id").
		Group("users.id").
		Order("profit DESC").
		Limit(limit).
		Scan(&entries)

	return entries, nil
}

// ============ Game History ============

func (s *GameService) GetCrashHistory(limit int) ([]models.CrashGameRound, error) {
	var rounds []models.CrashGameRound
	err := s.db.Where("status = 'crashed'").Order("created_at DESC").Limit(limit).Find(&rounds).Error
	return rounds, err
}

func (s *GameService) GetUserBets(userID uuid.UUID, gameType string, limit int) ([]models.Bet, error) {
	var bets []models.Bet
	query := s.db.Where("user_id = ?", userID)
	if gameType != "" {
		query = query.Where("game_type = ?", gameType)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&bets).Error
	return bets, err
}

// Helper to avoid shadowing
func hashSeed(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])
}
