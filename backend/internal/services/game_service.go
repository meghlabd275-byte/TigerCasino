package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

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

// UserSeed stores user's provably fair seed pair
type UserSeed struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;index"`
	ServerSeed       string
	ServerSeedHash   string
	ClientSeed       string
	Active           bool
	NextRevealNonce  int
	CreatedAt        time.Time
}

func (s *GameService) InitializeUserSeeds(userID uuid.UUID) error {
	var existing UserSeed
	err := s.db.Where("user_id = ? AND active = ?", userID, true).First(&existing).Error
	if err == nil {
		return nil // Already has active seeds
	}

	seeds, err := s.GenerateSeeds("")
	if err != nil {
		return err
	}

	userSeed := UserSeed{
		ID:             uuid.New(),
		UserID:         userID,
		ServerSeed:     seeds.ServerSeed,
		ServerSeedHash: seeds.ServerSeedHash,
		ClientSeed:     seeds.ClientSeed,
		Active:         true,
		NextRevealNonce: 0,
	}

	return s.db.Create(&userSeed).Error
}

func (s *GameService) GetUserSeeds(userID uuid.UUID) (*UserSeed, error) {
	var seeds UserSeed
	err := s.db.Where("user_id = ? AND active = ?", userID, true).First(&seeds).Error
	if err != nil {
		return nil, err
	}
	return &seeds, nil
}

func (s *GameService) RegenerateSeeds(userID uuid.UUID, clientSeed string) error {
	// Deactivate current seeds
	s.db.Model(&UserSeed{}).Where("user_id = ?", userID).Update("active", false)

	// Generate new seeds
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return err
	}

	userSeed := UserSeed{
		ID:             uuid.New(),
		UserID:         userID,
		ServerSeed:     seeds.ServerSeed,
		ServerSeedHash: seeds.ServerSeedHash,
		ClientSeed:     seeds.ClientSeed,
		Active:         true,
		NextRevealNonce: 0,
	}

	return s.db.Create(&userSeed).Error
}

// GenerateSeeds creates new provably fair seeds
func (s *GameService) GenerateSeeds(clientSeed string) (*Seeds, error) {
	serverSeedBytes := make([]byte, 32)
	if _, err := rand.Read(serverSeedBytes); err != nil {
		return nil, fmt.Errorf("failed to generate server seed: %w", err)
	}

	serverSeed := hex.EncodeToString(serverSeedBytes)
	serverSeedHash := s.HashSeed(serverSeed)

	if clientSeed == "" {
		clientBytes := make([]byte, 16)
		rand.Read(clientBytes)
		clientSeed = hex.EncodeToString(clientBytes)
	}

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

// VerifyResult verifies a game result using the original seeds
func (s *GameService) VerifyResult(serverSeed, clientSeed string, nonce int, result float64, max float64) bool {
	generated := s.security.GenerateOutcome(serverSeed, clientSeed, nonce)
	normalized := generated * max
	return math.Abs(normalized-result) < 0.0001
}

// ============ Game Models ============

type GameHistory struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	GameType    string                 `json:"game_type"`
	BetAmount   float64                `json:"bet_amount"`
	WinAmount   float64                `json:"win_amount"`
	Multiplier  float64                `json:"multiplier"`
	Profit      float64                `json:"profit"`
	Status      string                 `json:"status"`
	GameData    map[string]interface{} `json:"game_data"`
	Seeds       *Seeds                 `json:"seeds,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

type CrashGameState struct {
	RoundID     string    `json:"round_id"`
	CrashPoint  float64   `json:"crash_point"`
	Status      string    `json:"status"` // waiting, running, crashed
	StartTime   time.Time `json:"start_time"`
	Entries     []CrashEntry `json:"entries"`
}

type CrashEntry struct {
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	BetAmount   float64   `json:"bet_amount"`
	CashoutAt   float64   `json:"cashout_at,omitempty"`
	Won         bool      `json:"won"`
}

type PlinkoGameState struct {
	Rows   int       `json:"rows"`
	Risk   string    `json:"risk"` // low, medium, high
	Path   []int     `json:"path"`
	Bucket int       `json:"bucket"`
}

// ============ DICE GAME ============

type DiceResult struct {
	ID           uuid.UUID `json:"id"`
	Roll         float64   `json:"roll"`
	Target       float64   `json:"target"`
	Multiplier   float64   `json:"multiplier"`
	WinAmount    float64   `json:"win_amount"`
	Direction    string    `json:"direction"` // "over" or "under"
	Profit       float64   `json:"profit"`
	ServerSeed   string    `json:"server_seed"`
	ClientSeed   string    `json:"client_seed"`
	Nonce        int       `json:"nonce"`
	Verified     bool      `json:"verified"`
}

func (s *GameService) PlayDice(userID uuid.UUID, betAmount float64, target float64, direction string, clientSeed string) (*DiceResult, error) {
	// Get or initialize user seeds
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds(clientSeed)
	}

	nonce := seeds.NextRevealNonce

	// Generate roll using provably fair system
	roll := s.GenerateFloatOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce) * 100.0

	// Calculate win
	var multiplier float64
	win := false

	if direction == "over" {
		win = roll > target
		if win {
			multiplier = (100 - target) / target
		}
	} else {
		win = roll < target
		if win {
			multiplier = target / (100 - target)
		}
	}

	// Apply house edge (1%)
	multiplier = multiplier * 0.99

	var winAmount float64
	profit := -betAmount

	if win {
		winAmount = betAmount * multiplier
		profit = winAmount - betAmount
	}

	// Update balance
	userService := NewUserService(s.db)
	if win {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "dice",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     func() string { if win { return "won" }; return "lost" }(),
		GameData:   fmt.Sprintf(`{"roll":%f,"target":%f,"direction":"%s"}`, roll, target, direction),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+1)

	// Update user stats
	userService.UpdateWagered(userID, betAmount)

	return &DiceResult{
		ID:          bet.ID,
		Roll:        roll,
		Target:      target,
		Multiplier:  multiplier,
		WinAmount:   winAmount,
		Direction:   direction,
		Profit:      profit,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed:  seeds.ClientSeed,
		Nonce:       nonce,
		Verified:    true,
	}, nil
}

// ============ CRASH GAME ============

type CrashResult struct {
	ID           uuid.UUID `json:"id"`
	RoundID      string    `json:"round_id"`
	CrashPoint   float64   `json:"crash_point"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Status      string    `json:"status"` // waiting, running, crashed
	Profit       float64   `json:"profit"`
	WinAmount    float64   `json:"win_amount"`
	AutoCashout float64   `json:"auto_cashout,omitempty"`
}

type CrashBet struct {
	UserID      uuid.UUID `json:"user_id"`
	BetAmount   float64   `json:"bet_amount"`
	CashoutAt   float64   `json:"cashout_at,omitempty"`
	AutoCashout float64   `json:"auto_cashout,omitempty"`
	PlacedAt    time.Time `json:"placed_at"`
}

var crashGames = make(map[string]*CrashGameState)
var crashHistory []float64

func (s *GameService) StartCrashRound() (*CrashGameState, error) {
	seeds, _ := s.GenerateSeeds("")
	roundID := uuid.New().String()

	// Generate crash point using exponential distribution
	f := s.security.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, 0)
	var crashPoint float64

	if f < 0.70 {
		// 70% crash between 1.0x and 8.0x
		crashPoint = 1.0 + (f * 10.0)
	} else {
		// 30% can go higher
		crashPoint = 8.0 + ((f - 0.70) * 300.0)
	}

	// Apply house edge
	crashPoint = crashPoint * 0.97

	// Cap at reasonable maximum
	if crashPoint > 1000.0 {
		crashPoint = 1000.0
	}

	state := &CrashGameState{
		RoundID:    roundID,
		CrashPoint: crashPoint,
		Status:     "waiting",
		StartTime:  time.Now(),
		Entries:    make([]CrashEntry, 0),
	}

	crashGames[roundID] = state
	return state, nil
}

func (s *GameService) PlaceCrashBet(userID uuid.UUID, username string, roundID string, betAmount float64, autoCashout float64) error {
	game, exists := crashGames[roundID]
	if !exists {
		return fmt.Errorf("game not found")
	}

	entry := CrashEntry{
		UserID:      userID,
		Username:    username,
		BetAmount:   betAmount,
		AutoCashout: autoCashout,
		PlacedAt:    time.Now(),
	}

	game.Entries = append(game.Entries, entry)

	// Deduct balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, -betAmount)

	return nil
}

func (s *GameService) CashoutCrash(userID uuid.UUID, roundID string, multiplier float64) (*CrashResult, error) {
	game, exists := crashGames[roundID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	// Find user's bet
	var entry *CrashEntry
	var entryIdx int
	for i := range game.Entries {
		if game.Entries[i].UserID == userID && game.Entries[i].Won == false && game.Entries[i].CashoutAt == 0 {
			entry = &game.Entries[i]
			entryIdx = i
			break
		}
	}

	if entry == nil {
		return nil, fmt.Errorf("no active bet found")
	}

	// Check if already crashed
	if game.Status == "crashed" {
		return nil, fmt.Errorf("game already crashed")
	}

	// Check if cashout is valid
	if multiplier >= game.CrashPoint {
		return nil, fmt.Errorf("invalid cashout multiplier")
	}

	// Calculate winnings
	winAmount := entry.BetAmount * multiplier
	profit := winAmount - entry.BetAmount

	// Update entry
	game.Entries[entryIdx].CashoutAt = multiplier
	game.Entries[entryIdx].Won = true

	// Update balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, winAmount)
	userService.UpdateWagered(userID, entry.BetAmount)

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "crash",
		BetAmount:  entry.BetAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     "won",
		GameData:   fmt.Sprintf(`{"round_id":"%s","cashout":%f}`, roundID, multiplier),
	}
	s.db.Create(&bet)

	return &CrashResult{
		ID:         bet.ID,
		RoundID:    roundID,
		CrashPoint: game.CrashPoint,
		Profit:     profit,
		WinAmount:  winAmount,
		Status:     "cashouted",
	}, nil
}

func (s *GameService) GetCurrentCrashState() *CrashGameState {
	for _, game := range crashGames {
		if game.Status != "crashed" {
			return game
		}
	}
	// Start new game if none exists
	state, _ := s.StartCrashRound()
	return state
}

func (s *GameService) GetCrashHistory() []float64 {
	return crashHistory
}

// ============ SLOTS ============

type SlotsResult struct {
	ID           uuid.UUID       `json:"id"`
	Reels        [5]int          `json:"reels"`
	Symbols      [5]string       `json:"symbols"`
	Multiplier    float64         `json:"multiplier"`
	WinAmount     float64         `json:"win_amount"`
	Profit        float64         `json:"profit"`
	WinLine       string          `json:"win_line"`
	IsJackpot     bool            `json:"is_jackpot"`
	IsFreeSpin    bool            `json:"is_free_spin"`
	ServerSeed   string          `json:"server_seed"`
	ClientSeed   string          `json:"client_seed"`
	Nonce        int             `json:"nonce"`
}

var slotSymbols = []string{"Cherry", "Lemon", "Orange", "Grape", "Watermelon", "Bell", "Diamond", "Seven", "Star"}

func (s *GameService) PlaySlots(userID uuid.UUID, betAmount float64, lines int) (*SlotsResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Generate reel outcomes
	var reels [5]int
	var symbols [5]string

	for i := 0; i < 5; i++ {
		reels[i] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i, len(slotSymbols))
		symbols[i] = slotSymbols[reels[i]]
	}

	// Check for wins
	isJackpot := reels[0] == 7 && reels[1] == 7 && reels[2] == 7 && reels[3] == 7 && reels[4] == 7
	isThreeOfKind := (reels[0] == reels[1] && reels[1] == reels[2]) ||
		(reels[1] == reels[2] && reels[2] == reels[3]) ||
		(reels[2] == reels[3] && reels[3] == reels[4])
	isTwoOfKind := reels[0] == reels[1] || reels[1] == reels[2] || reels[2] == reels[3] || reels[3] == reels[4]

	var multiplier float64
	var winLine string

	if isJackpot {
		multiplier = 100.0
		winLine = "JACKPOT!"
	} else if reels[0] == reels[1] && reels[1] == reels[2] && reels[2] == reels[3] {
		multiplier = 20.0
		winLine = "Four of a Kind!"
	} else if isThreeOfKind {
		multiplier = 5.0
		winLine = "Three of a Kind!"
	} else if isTwoOfKind {
		multiplier = 1.0
		winLine = "Two of a Kind!"
	}

	// Apply house edge (5%)
	multiplier = multiplier * 0.95

	winAmount := betAmount * float64(lines) * multiplier
	profit := winAmount - (betAmount * float64(lines))

	// Update balance
	userService := NewUserService(s.db)
	if winAmount > 0 {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "slots",
		BetAmount:  betAmount * float64(lines),
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     func() string { if winAmount > 0 { return "won" }; return "lost" }(),
		GameData:   fmt.Sprintf(`{"reels":[%d,%d,%d,%d,%d],"symbols":["%s","%s","%s","%s","%s"],"lines":%d}`, reels[0], reels[1], reels[2], reels[3], reels[4], symbols[0], symbols[1], symbols[2], symbols[3], symbols[4], lines),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+5)

	// Update user stats
	userService.UpdateWagered(userID, betAmount*float64(lines))

	return &SlotsResult{
		ID:          bet.ID,
		Reels:       reels,
		Symbols:     symbols,
		Multiplier:  multiplier,
		WinAmount:   winAmount,
		Profit:      profit,
		WinLine:     winLine,
		IsJackpot:   isJackpot,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:       nonce,
	}, nil
}

// ============ ROULETTE ============

type RouletteResult struct {
	ID          uuid.UUID `json:"id"`
	Number      int       `json:"number"`
	Color       string    `json:"color"`
	Parity      string    `json:"parity"`
	Range       string    `json:"range"`
	WinAmount   float64   `json:"win_amount"`
	Profit       float64   `json:"profit"`
	BetType     string    `json:"bet_type"`
	BetValue    interface{} `json:"bet_value"`
	ServerSeed  string    `json:"server_seed"`
	ClientSeed  string    `json:"client_seed"`
	Nonce       int       `json:"nonce"`
}

func (s *GameService) PlayRoulette(userID uuid.UUID, betAmount float64, betType string, betValue interface{}) (*RouletteResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Generate winning number (0-36)
	winningNumber := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce, 37)

	// Determine color
	color := "green"
	if winningNumber != 0 {
		if isRedNumber(winningNumber) {
			color = "red"
		} else {
			color = "black"
		}
	}

	// Determine parity
	parity := "even"
	if winningNumber%2 != 0 {
		parity = "odd"
	}

	// Determine range
	range_ := "low"
	if winningNumber > 18 {
		range_ = "high"
	} else if winningNumber == 0 {
		range_ = "zero"
	}

	// Check win
	win := false
	multiplier := 0.0

	switch betType {
	case "straight":
		if winningNumber == betValue.(int) {
			win = true
			multiplier = 35.0
		}
	case "red":
		if color == "red" {
			win = true
			multiplier = 2.0
		}
	case "black":
		if color == "black" {
			win = true
			multiplier = 2.0
		}
	case "odd":
		if parity == "odd" && winningNumber != 0 {
			win = true
			multiplier = 2.0
		}
	case "even":
		if parity == "even" && winningNumber != 0 {
			win = true
			multiplier = 2.0
		}
	case "low":
		if range_ == "low" && winningNumber != 0 {
			win = true
			multiplier = 2.0
		}
	case "high":
		if range_ == "high" {
			win = true
			multiplier = 2.0
		}
	case "dozen":
		dozen := betValue.(int)
		start := (dozen - 1) * 12
		if winningNumber > start && winningNumber <= start+12 {
			win = true
			multiplier = 3.0
		}
	case "column":
		col := betValue.(int)
		colNums := []int{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34}
		for _, n := range colNums {
			if winningNumber == n+(col-1)*3 {
				win = true
				break
			}
		}
		if win {
			multiplier = 3.0
		}
	}

	// Apply house edge (2.7% for single zero)
	if winningNumber == 0 {
		win = false
	}
	multiplier = multiplier * 0.973 // House edge

	winAmount := 0.0
	if win {
		winAmount = betAmount * multiplier
	}
	profit := winAmount - betAmount

	// Update balance
	userService := NewUserService(s.db)
	if winAmount > 0 {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	betValueJSON, _ := json.Marshal(betValue)
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "roulette",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     func() string { if win { return "won" }; return "lost" }(),
		GameData:   fmt.Sprintf(`{"winning_number":%d,"color":"%s","bet_type":"%s","bet_value":%s}`, winningNumber, color, betType, string(betValueJSON)),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+1)

	// Update user stats
	userService.UpdateWagered(userID, betAmount)

	return &RouletteResult{
		ID:         bet.ID,
		Number:     winningNumber,
		Color:      color,
		Parity:     parity,
		Range:      range_,
		WinAmount:  winAmount,
		Profit:     profit,
		BetType:    betType,
		BetValue:   betValue,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
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

// ============ BLACKJACK ============

type BlackjackResult struct {
	ID           uuid.UUID `json:"id"`
	PlayerHand   []string  `json:"player_hand"`
	DealerHand  []string  `json:"dealer_hand"`
	PlayerScore  int       `json:"player_score"`
	DealerScore  int       `json:"dealer_score"`
	WinAmount    float64   `json:"win_amount"`
	Profit        float64   `json:"profit"`
	Status       string    `json:"status"` // win, lose, push, blackjack
	ServerSeed   string    `json:"server_seed"`
	ClientSeed   string    `json:"client_seed"`
	Nonce        int       `json:"nonce"`
}

var cardSuits = []string{"♠", "♥", "♦", "♣"}
var cardValues = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func (s *GameService) PlayBlackjack(userID uuid.UUID, betAmount float64) (*BlackjackResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Generate cards
	playerHand := make([]string, 2)
	dealerHand := make([]string, 2)

	for i := 0; i < 2; i++ {
		cardIdx := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i, 52)
		suit := cardSuits[cardIdx/13]
		value := cardValues[cardIdx%13]
		playerHand[i] = fmt.Sprintf("%s%s", value, suit)

		cardIdx = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i+2, 52)
		suit = cardSuits[cardIdx/13]
		value = cardValues[cardIdx%13]
		dealerHand[i] = fmt.Sprintf("%s%s", value, suit)
	}

	// Calculate scores
	playerScore := s.calculateBlackjackScore(playerHand)
	dealerScore := s.calculateBlackjackScore(dealerHand)

	// Determine outcome
	status := "lose"
	multiplier := 0.0

	if playerScore > 21 {
		status = "lose"
	} else if dealerScore > 21 {
		status = "win"
		multiplier = 2.0
	} else if playerScore == 21 && len(playerHand) == 2 {
		if dealerScore == 21 && len(dealerHand) == 2 {
			status = "push"
		} else {
			status = "blackjack"
			multiplier = 2.5
		}
	} else if playerScore > dealerScore {
		status = "win"
		multiplier = 2.0
	} else if playerScore < dealerScore {
		status = "lose"
	} else {
		status = "push"
	}

	winAmount := 0.0
	if status == "blackjack" {
		winAmount = betAmount * multiplier
	} else if status == "win" {
		winAmount = betAmount * multiplier
	}
	profit := winAmount - betAmount

	// Update balance
	userService := NewUserService(s.db)
	if winAmount > 0 {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "blackjack",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     status,
		GameData:   fmt.Sprintf(`{"player_hand":["%s","%s"],"dealer_hand":["%s","%s"],"player_score":%d,"dealer_score":%d}`, playerHand[0], playerHand[1], dealerHand[0], dealerHand[1], playerScore, dealerScore),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+4)

	// Update user stats
	userService.UpdateWagered(userID, betAmount)

	return &BlackjackResult{
		ID:          bet.ID,
		PlayerHand:  playerHand,
		DealerHand:  dealerHand,
		PlayerScore: playerScore,
		DealerScore: dealerScore,
		WinAmount:   winAmount,
		Profit:      profit,
		Status:      status,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:       nonce,
	}, nil
}

func (s *GameService) calculateBlackjackScore(hand []string) int {
	score := 0
	aces := 0

	for _, card := range hand {
		value := card[:len(card)-1] // Remove suit
		switch value {
		case "J", "Q", "K":
			score += 10
		case "A":
			aces++
			score += 11
		default:
			v, _ := strconv.Atoi(value)
			score += v
		}
	}

	// Adjust for aces
	for aces > 0 && score > 21 {
		score -= 10
		aces--
	}

	return score
}

// ============ VIDEO POKER ============

type VideoPokerResult struct {
	ID          uuid.UUID `json:"id"`
	Hand        []string  `json:"hand"`
	FinalHand   []string  `json:"final_hand"`
	WinAmount    float64   `json:"win_amount"`
	Profit       float64   `json:"profit"`
	HandRank    string    `json:"hand_rank"`
	ServerSeed  string    `json:"server_seed"`
	ClientSeed  string    `json:"client_seed"`
	Nonce       int       `json:"nonce"`
}

func (s *GameService) PlayVideoPoker(userID uuid.UUID, betAmount float64) (*VideoPokerResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Generate initial hand
	hand := make([]string, 5)
	deck := make([]int, 52)
	for i := range deck {
		deck[i] = i
	}

	// Shuffle using Fisher-Yates
	for i := len(deck) - 1; i > 0; i-- {
		j := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i, i+1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	for i := 0; i < 5; i++ {
		cardIdx := deck[i]
		suit := cardSuits[cardIdx/13]
		value := cardValues[cardIdx%13]
		hand[i] = fmt.Sprintf("%s%s", value, suit)
	}

	// Evaluate hand
	handRank, multiplier := s.evaluateVideoPokerHand(hand)

	// Apply house edge (5%)
	multiplier = multiplier * 0.95

	winAmount := betAmount * multiplier
	profit := winAmount - betAmount

	// Update balance
	userService := NewUserService(s.db)
	if winAmount > 0 {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "video_poker",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     func() string { if winAmount > 0 { return "won" }; return "lost" }(),
		GameData:   fmt.Sprintf(`{"hand":["%s","%s","%s","%s","%s"],"hand_rank":"%s"}`, hand[0], hand[1], hand[2], hand[3], hand[4], handRank),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+52)

	// Update user stats
	userService.UpdateWagered(userID, betAmount)

	return &VideoPokerResult{
		ID:         bet.ID,
		Hand:       hand,
		FinalHand:  hand,
		WinAmount:  winAmount,
		Profit:     profit,
		HandRank:   handRank,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
	}, nil
}

func (s *GameService) evaluateVideoPokerHand(hand []string) (string, float64) {
	// Extract values and suits
	values := make([]int, 5)
	suits := make([]int, 5)

	for i, card := range hand {
		value := card[:len(card)-1]
		suit := card[len(card)-1:]

		switch value {
		case "A":
			values[i] = 14
		case "K":
			values[i] = 13
		case "Q":
			values[i] = 12
		case "J":
			values[i] = 11
		default:
			v, _ := strconv.Atoi(value)
			values[i] = v
		}

		switch suit {
		case "♠":
			suits[i] = 0
		case "♥":
			suits[i] = 1
		case "♦":
			suits[i] = 2
		case "♣":
			suits[i] = 3
		}
	}

	// Check flush
	isFlush := true
	for i := 1; i < 5; i++ {
		if suits[i] != suits[0] {
			isFlush = false
			break
		}
	}

	// Sort values
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 5; j++ {
			if values[i] > values[j] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}

	// Check straight
	isStraight := true
	for i := 1; i < 5; i++ {
		if values[i] != values[i-1]+1 {
			isStraight = false
			break
		}
	}

	// Check for Royal Flush
	if isFlush && isStraight && values[0] == 10 {
		return "Royal Flush", 250.0
	}

	// Straight Flush
	if isFlush && isStraight {
		return "Straight Flush", 50.0
	}

	// Four of a Kind
	for i := 0; i < 2; i++ {
		if values[i] == values[i+1] && values[i+1] == values[i+2] && values[i+2] == values[i+3] {
			return "Four of a Kind", 25.0
		}
	}

	// Full House
	if (values[0] == values[1] && values[2] == values[3] && values[3] == values[4] &&
		(values[1] == values[2] || values[3] == values[4])) ||
		(values[0] == values[1] && values[1] == values[2] && values[3] == values[4] &&
			(values[2] == values[3])) {
		return "Full House", 9.0
	}

	// Flush
	if isFlush {
		return "Flush", 6.0
	}

	// Straight
	if isStraight {
		return "Straight", 4.0
	}

	// Three of a Kind
	for i := 0; i < 3; i++ {
		if values[i] == values[i+1] && values[i+1] == values[i+2] {
			return "Three of a Kind", 3.0
		}
	}

	// Two Pair
	if (values[0] == values[1] && values[2] == values[3]) ||
		(values[1] == values[2] && values[3] == values[4]) {
		return "Two Pair", 2.0
	}

	// Jacks or Better
	for i := 0; i < 4; i++ {
		if values[i] >= 11 && values[i] == values[i+1] {
			return "Jacks or Better", 1.0
		}
	}

	return "Nothing", 0.0
}

// ============ MINES GAME ============

type MinesResult struct {
	ID               uuid.UUID `json:"id"`
	GridSize         int       `json:"grid_size"`
	MinesCount       int       `json:"mines_count"`
	Revealed         []int     `json:"revealed"`
	MineLocations    []int     `json:"mine_locations"`
	CurrentMultiplier float64   `json:"current_multiplier"`
	WinAmount        float64   `json:"win_amount"`
	Profit           float64    `json:"profit"`
	Status           string     `json:"status"`
	ServerSeed       string    `json:"server_seed"`
	ClientSeed       string    `json:"client_seed"`
	Nonce            int       `json:"nonce"`
}

var minesGames = make(map[string]*MinesResult)

func (s *GameService) StartMinesGame(userID uuid.UUID, betAmount float64, minesCount int, gridSize int) (*MinesResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce
	gameID := uuid.New().String()

	// Generate mine positions
	totalTiles := gridSize * gridSize
 minePositions := make([]int, minesCount)
 used := make(map[int]bool)
 
 for i := 0; i < minesCount; i++ {
     pos := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i, totalTiles)
     for used[pos] {
         pos = (pos + 1) % totalTiles
     }
     used[pos] = true
     minePositions[i] = pos
 }

	result := &MinesResult{
		ID:                uuid.New(),
		GridSize:          gridSize,
		MinesCount:        minesCount,
		Revealed:          make([]int, 0),
		MineLocations:      minePositions,
		CurrentMultiplier: 1.0,
		Status:            "playing",
		ServerSeed:        seeds.ServerSeed,
		ClientSeed:        seeds.ClientSeed,
		Nonce:             nonce,
	}

	minesGames[gameID] = result

	// Deduct balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, -betAmount)

	return result, nil
}

func (s *GameService) RevealMinesTile(userID uuid.UUID, gameID string, tileIndex int) (*MinesResult, error) {
	result, exists := minesGames[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	// Check if already revealed
	for _, revealed := range result.Revealed {
		if revealed == tileIndex {
			return nil, fmt.Errorf("tile already revealed")
		}
	}

	// Check if it's a mine
	for _, mine := range result.MineLocations {
		if mine == tileIndex {
			result.Status = "lost"
			return result, nil
		}
	}

	// Safe tile - reveal it
	result.Revealed = append(result.Revealed, tileIndex)

	// Calculate new multiplier based on remaining safe tiles
	totalTiles := result.GridSize * result.GridSize
	safeTiles := totalTiles - result.MinesCount
	revealedSafe := len(result.Revealed)
	remainingSafe := safeTiles - revealedSafe
	totalRemaining := totalTiles - revealedSafe

	// Multiplier increases as more tiles are revealed
	result.CurrentMultiplier = float64(totalRemaining) / float64(remainingSafe)

	return result, nil
}

func (s *GameService) CashoutMines(userID uuid.UUID, gameID string, betAmount float64) (*MinesResult, error) {
	result, exists := minesGames[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if result.Status != "playing" {
		return nil, fmt.Errorf("game already finished")
	}

	result.WinAmount = betAmount * result.CurrentMultiplier * 0.98 // House edge
	result.Profit = result.WinAmount - betAmount
	result.Status = "won"

	// Update balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, result.WinAmount)
	userService.UpdateWagered(userID, betAmount)

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "mines",
		BetAmount:  betAmount,
		WinAmount:  result.WinAmount,
		Multiplier: result.CurrentMultiplier,
		Profit:     result.Profit,
		Status:     "won",
		GameData:   fmt.Sprintf(`{"mines_count":%d,"revealed":%v}`, result.MinesCount, result.Revealed),
	}
	s.db.Create(&bet)

	return result, nil
}

// ============ PLINKO GAME ============

type PlinkoResult struct {
	ID           uuid.UUID `json:"id"`
	Rows         int       `json:"rows"`
	Risk         string    `json:"risk"`
	Path         []int     `json:"path"`
	Bucket       int       `json:"bucket"`
	Multiplier   float64   `json:"multiplier"`
	WinAmount    float64   `json:"win_amount"`
	Profit       float64   `json:"profit"`
	ServerSeed   string    `json:"server_seed"`
	ClientSeed   string    `json:"client_seed"`
	Nonce        int       `json:"nonce"`
}

func (s *GameService) PlayPlinko(userID uuid.UUID, betAmount float64, rows int, risk string) (*PlinkoResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Simulate ball path
	path := make([]int, rows)
	bucket := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce, rows+1)
	
	position := rows / 2
	for i := 0; i < rows; i++ {
		dir := s.GenerateFloatOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i+100)
		if dir < 0.5 {
			position = max(0, position-1)
		} else {
			position = min(i+1, position+1)
		}
		path[i] = position
	}

	// Calculate multiplier based on bucket and risk
	multiplier := s.calculatePlinkoMultiplier(bucket, rows, risk)

	// Apply house edge
	multiplier = multiplier * 0.97

	winAmount := betAmount * multiplier
	profit := winAmount - betAmount

	// Update balance
	userService := NewUserService(s.db)
	if winAmount > 0 {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "plinko",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     func() string { if winAmount > 0 { return "won" }; return "lost" }(),
		GameData:   fmt.Sprintf(`{"rows":%d,"risk":"%s","bucket":%d,"path":%v}`, rows, risk, bucket, path),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+rows+100)

	// Update user stats
	userService.UpdateWagered(userID, betAmount)

	return &PlinkoResult{
		ID:         bet.ID,
		Rows:       rows,
		Risk:       risk,
		Path:       path,
		Bucket:     bucket,
		Multiplier: multiplier,
		WinAmount:  winAmount,
		Profit:     profit,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
	}, nil
}

func (s *GameService) calculatePlinkoMultiplier(bucket, rows int, risk string) float64 {
	// Simplified plinko payout table
	lowRiskMultipliers := []float64{0.5, 1.0, 1.5, 2.0, 3.0, 5.0, 3.0, 2.0, 1.5, 1.0, 0.5}
	medRiskMultipliers := []float64{0.5, 1.0, 2.0, 3.0, 5.0, 10.0, 5.0, 3.0, 2.0, 1.0, 0.5}
	highRiskMultipliers := []float64{0.5, 2.0, 3.0, 5.0, 10.0, 25.0, 10.0, 5.0, 3.0, 2.0, 0.5}

	var multipliers []float64
	switch risk {
	case "high":
		multipliers = highRiskMultipliers
	case "medium":
		multipliers = medRiskMultipliers
	default:
		multipliers = lowRiskMultipliers
	}

	// Adjust array size based on rows
	maxBuckets := len(lowRiskMultipliers)
	bucketIdx := bucket * maxBuckets / (rows + 1)
	if bucketIdx >= maxBuckets {
		bucketIdx = maxBuckets - 1
	}

	return multipliers[bucketIdx]
}

// ============ LIMBO GAME ============

type LimboResult struct {
	ID                  uuid.UUID `json:"id"`
	TargetMultiplier     float64   `json:"target_multiplier"`
	ResultMultiplier     float64   `json:"result_multiplier"`
	WinAmount            float64   `json:"win_amount"`
	Profit               float64   `json:"profit"`
	Status               string    `json:"status"`
	ServerSeed           string    `json:"server_seed"`
	ClientSeed           string    `json:"client_seed"`
	Nonce                int       `json:"nonce"`
}

func (s *GameService) PlayLimbo(userID uuid.UUID, betAmount float64, targetMultiplier float64) (*LimboResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Generate result multiplier using exponential-like distribution
	f := s.security.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce)
	resultMultiplier := 1.0 / (1.0 - f)

	// Cap at reasonable maximum
	if resultMultiplier > 1000.0 {
		resultMultiplier = 1000.0
	}

	// Determine win
	win := resultMultiplier >= targetMultiplier

	// Calculate multiplier (payout = result / target)
	var multiplier float64
	if win {
		multiplier = resultMultiplier / targetMultiplier
	}

	// Apply house edge (2%)
	multiplier = multiplier * 0.98

	winAmount := 0.0
	if win {
		winAmount = betAmount * multiplier
	}
	profit := winAmount - betAmount

	// Update balance
	userService := NewUserService(s.db)
	if winAmount > 0 {
		userService.UpdateBalance(userID, winAmount)
	}

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "limbo",
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     func() string { if win { return "won" }; return "lost" }(),
		GameData:   fmt.Sprintf(`{"target":%f,"result":%f}`, targetMultiplier, resultMultiplier),
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}
	s.db.Create(&bet)

	// Update seed nonce
	s.db.Model(&UserSeed{}).Where("user_id = ? AND active = ?", userID, true).
		Update("next_reveal_nonce", nonce+1)

	// Update user stats
	userService.UpdateWagered(userID, betAmount)

	return &LimboResult{
		ID:                  bet.ID,
		TargetMultiplier:    targetMultiplier,
		ResultMultiplier:    resultMultiplier,
		WinAmount:           winAmount,
		Profit:              profit,
		Status:              func() string { if win { return "won" }; return "lost" }(),
		ServerSeed:          seeds.ServerSeed,
		ClientSeed:          seeds.ClientSeed,
		Nonce:               nonce,
	}, nil
}

// ============ HI-LO GAME ============

type HiloResult struct {
	ID            uuid.UUID `json:"id"`
	CurrentCard   string    `json:"current_card"`
	NextCard      string    `json:"next_card"`
	Choice        string    `json:"choice"`
	IsCorrect     bool      `json:"is_correct"`
	Streak        int       `json:"streak"`
	Multiplier    float64   `json:"multiplier"`
	WinAmount     float64   `json:"win_amount"`
	Profit        float64   `json:"profit"`
	Status        string    `json:"status"`
	ServerSeed    string    `json:"server_seed"`
	ClientSeed    string    `json:"client_seed"`
	Nonce         int       `json:"nonce"`
}

type HiloGameState struct {
	UserID       uuid.UUID
	Cards        []string
	CurrentIndex int
	BetAmount    float64
	Streak       int
	ServerSeed   string
	ClientSeed   string
	Nonce        int
}

var hiloGames = make(map[uuid.UUID]*HiloGameState)

func (s *GameService) StartHiloGame(userID uuid.UUID, betAmount float64) (*HiloResult, error) {
	seeds, err := s.GetUserSeeds(userID)
	if err != nil {
		s.InitializeUserSeeds(userID)
		seeds, _ = s.GetUserSeeds(userID)
	}
	if seeds == nil {
		seeds, _ = s.GenerateSeeds("")
	}

	nonce := seeds.NextRevealNonce

	// Generate deck
	cards := make([]string, 52)
	for i := 0; i < 52; i++ {
		suit := cardSuits[i/13]
		value := cardValues[i%13]
		cards[i] = fmt.Sprintf("%s%s", value, suit)
	}

	// Shuffle
	for i := len(cards) - 1; i > 0; i-- {
		j := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, nonce+i, i+1)
		cards[i], cards[j] = cards[j], cards[i]
	}

	state := &HiloGameState{
		UserID:       userID,
		Cards:        cards,
		CurrentIndex: 0,
		BetAmount:    betAmount,
		Streak:       0,
		ServerSeed:   seeds.ServerSeed,
		ClientSeed:   seeds.ClientSeed,
		Nonce:        nonce,
	}
	hiloGames[userID] = state

	// Deduct balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, -betAmount)

	return &HiloResult{
		CurrentCard: cards[0],
		Streak:      0,
		Multiplier:  1.0,
		WinAmount:   0,
		Status:      "playing",
		ServerSeed:  seeds.ServerSeed,
		ClientSeed:  seeds.ClientSeed,
		Nonce:       nonce,
	}, nil
}

func (s *GameService) PlayHiloChoice(userID uuid.UUID, choice string) (*HiloResult, error) {
	state, exists := hiloGames[userID]
	if !exists {
		return nil, fmt.Errorf("no active game")
	}

	// Draw next card
	state.CurrentIndex++
	if state.CurrentIndex >= len(state.Cards) {
		// Reshuffle
		state.CurrentIndex = 0
	}

	currentCard := state.Cards[state.CurrentIndex-1]
	nextCard := state.Cards[state.CurrentIndex]

	currentValue := s.getCardValue(state.Cards[state.CurrentIndex-1])
	nextValue := s.getCardValue(nextCard)

	// Check if choice is correct
	isCorrect := false
	switch choice {
	case "higher":
		isCorrect = nextValue > currentValue
	case "lower":
		isCorrect = nextValue < currentValue
	case "equal":
		isCorrect = nextValue == currentValue
	}

	if isCorrect {
		state.Streak++
	} else {
		state.Status = "lost"
		delete(hiloGames, userID)

		// Record bet
		bet := models.Bet{
			ID:         uuid.New(),
			UserID:     userID,
			GameType:   "hilo",
			BetAmount:  state.BetAmount,
			WinAmount:  0,
			Multiplier: 0,
			Profit:     -state.BetAmount,
			Status:     "lost",
		}
		s.db.Create(&bet)

		return &HiloResult{
			CurrentCard: currentCard,
			NextCard:    nextCard,
			Choice:      choice,
			IsCorrect:   false,
			Streak:      state.Streak,
			Status:      "lost",
			ServerSeed:  state.ServerSeed,
			ClientSeed:  state.ClientSeed,
			Nonce:       state.Nonce,
		}, nil
	}

	// Calculate multiplier based on streak
	multiplier := math.Pow(1.5, float64(state.Streak))
	if multiplier > 50.0 {
		multiplier = 50.0
	}

	return &HiloResult{
		CurrentCard: currentCard,
		NextCard:    nextCard,
		Choice:      choice,
		IsCorrect:   true,
		Streak:      state.Streak,
		Multiplier:  multiplier,
		Status:      "playing",
		ServerSeed:  state.ServerSeed,
		ClientSeed:  state.ClientSeed,
		Nonce:       state.Nonce,
	}, nil
}

func (s *GameService) CashoutHilo(userID uuid.UUID) (*HiloResult, error) {
	state, exists := hiloGames[userID]
	if !exists {
		return nil, fmt.Errorf("no active game")
	}

	multiplier := math.Pow(1.5, float64(state.Streak))
	if multiplier > 50.0 {
		multiplier = 50.0
	}
	multiplier = multiplier * 0.98 // House edge

	winAmount := state.BetAmount * multiplier
	profit := winAmount - state.BetAmount

	// Update balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, winAmount)
	userService.UpdateWagered(userID, state.BetAmount)

	// Record bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID,
		GameType:   "hilo",
		BetAmount:  state.BetAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		Profit:     profit,
		Status:     "won",
	}
	s.db.Create(&bet)

	delete(hiloGames, userID)

	return &HiloResult{
		CurrentCard: state.Cards[state.CurrentIndex],
		Streak:      state.Streak,
		Multiplier:  multiplier,
		WinAmount:   winAmount,
		Profit:      profit,
		Status:      "won",
		ServerSeed:  state.ServerSeed,
		ClientSeed:  state.ClientSeed,
		Nonce:       state.Nonce,
	}, nil
}

func (s *GameService) getCardValue(card string) int {
	value := card[:len(card)-1]
	switch value {
	case "A":
		return 1
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	default:
		v, _ := strconv.Atoi(value)
		return v
	}
}

// ============ BETS & HISTORY ============

func (s *GameService) GetUserBets(userID uuid.UUID, gameType string, limit int) ([]models.Bet, error) {
	var bets []models.Bet
	query := s.db.Where("user_id = ?", userID)
	if gameType != "" {
		query = query.Where("game_type = ?", gameType)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&bets).Error
	return bets, err
}

func (s *GameService) GetBetByID(betID uuid.UUID) (*models.Bet, error) {
	var bet models.Bet
	err := s.db.First(&bet, betID).Error
	if err != nil {
		return nil, err
	}
	return &bet, nil
}

// ============ LEADERBOARD ============

type LeaderboardEntry struct {
	Rank       int       `json:"rank"`
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	TotalProfit float64   `json:"total_profit"`
	TotalWagered float64 `json:"total_wagered"`
	BetCount   int       `json:"bet_count"`
}

func (s *GameService) GetLeaderboard(gameType string, timeframe string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	// Calculate timeframe
	var startTime time.Time
	now := time.Now()
	switch timeframe {
	case "daily":
		startTime = now.AddDate(0, 0, -1)
	case "weekly":
		startTime = now.AddDate(0, 0, -7)
	case "monthly":
		startTime = now.AddDate(0, -1, 0)
	default:
		startTime = now.AddDate(0, 0, -30)
	}

	query := `
		SELECT user_id, 
			   SUM(profit) as total_profit, 
			   SUM(bet_amount) as total_wagered, 
			   COUNT(*) as bet_count
		FROM bets
		WHERE created_at >= ? AND status IN ('won', 'lost')
	`
	args := []interface{}{startTime}

	if gameType != "" {
		query += " AND game_type = ?"
		args = append(args, gameType)
	}

	query += " GROUP BY user_id ORDER BY total_profit DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rank := 1
	for rows.Next() {
		var entry LeaderboardEntry
		var user models.User
		rows.Scan(&entry.UserID, &entry.TotalProfit, &entry.TotalWagered, &entry.BetCount)
		s.db.First(&user, entry.UserID)
		entry.Username = user.Username
		entry.Rank = rank
		entries = append(entries, entry)
		rank++
	}

	return entries, nil
}

// ============ STATISTICS ============

type GameStats struct {
	TotalBets      int64   `json:"total_bets"`
	TotalWagered   float64 `json:"total_wagered"`
	TotalWins      float64 `json:"total_wins"`
	TotalLosses    float64 `json:"total_losses"`
	NetProfit      float64 `json:"net_profit"`
	AverageBet     float64 `json:"average_bet"`
	WinRate        float64 `json:"win_rate"`
}

func (s *GameService) GetUserStats(userID uuid.UUID) (*GameStats, error) {
	var stats GameStats

	row := s.db.Model(&models.Bet{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_bets, COALESCE(SUM(bet_amount), 0) as total_wagered, COALESCE(SUM(win_amount), 0) as total_wins, COALESCE(SUM(CASE WHEN status = 'lost' THEN bet_amount ELSE 0 END), 0) as total_losses").
		Row()

	row.Scan(&stats.TotalBets, &stats.TotalWagered, &stats.TotalWins, &stats.TotalLosses)
	stats.NetProfit = stats.TotalWins - stats.TotalWagered
	if stats.TotalBets > 0 {
		stats.AverageBet = stats.TotalWagered / float64(stats.TotalBets)
	}
	if stats.TotalWagered > 0 {
		stats.WinRate = stats.TotalWins / stats.TotalWagered * 100
	}

	return &stats, nil
}

func (s *GameService) GetGameStats(gameType string) (*GameStats, error) {
	var stats GameStats

	query := s.db.Model(&models.Bet{}).Where("game_type = ?", gameType)

	row := query.
		Select("COUNT(*) as total_bets, COALESCE(SUM(bet_amount), 0) as total_wagered, COALESCE(SUM(win_amount), 0) as total_wins, COALESCE(SUM(CASE WHEN status = 'lost' THEN bet_amount ELSE 0 END), 0) as total_losses").
		Row()

	row.Scan(&stats.TotalBets, &stats.TotalWagered, &stats.TotalWins, &stats.TotalLosses)
	stats.NetProfit = stats.TotalWins - stats.TotalWagered
	if stats.TotalBets > 0 {
		stats.AverageBet = stats.TotalWagered / float64(stats.TotalBets)
	}
	if stats.TotalWagered > 0 {
		stats.WinRate = stats.TotalWins / stats.TotalWagered * 100
	}

	return &stats, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
