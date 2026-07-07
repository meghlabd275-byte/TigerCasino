package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// LotteryService handles all lottery and instant win games
type LotteryService struct {
	db *gorm.DB
}

func NewLotteryService(db *gorm.DB) *LotteryService {
	return &LotteryService{db: db}
}

// ============ KENO ============

type KenoResult struct {
	GameID       string  `json:"game_id"`
	Spots        int     `json:"spots"`
	Selected     []int   `json:"selected"`
	Drawn        []int   `json:"drawn"`
	Matches      int     `json:"matches"`
	BetAmount    float64 `json:"bet_amount"`
	WinAmount    float64 `json:"win_amount"`
	Multiplier   float64 `json:"multiplier"`
	ServerSeed   string  `json:"server_seed"`
	ClientSeed   string  `json:"client_seed"`
	Nonce        int     `json:"nonce"`
}

// Keno paytable (multiplier based on matches and spots)
var kenoPaytable = map[int]map[int]float64{
	1:  {1: 3.0},
	2:  {1: 1.0, 2: 9.0},
	3:  {1: 1.0, 2: 2.0, 3: 27.0},
	4:  {1: 0.5, 2: 1.0, 3: 4.0, 4: 80.0},
	5:  {0.5: 1.0, 2: 2.0, 3: 5.0, 4: 10.0, 5: 400.0},
	6:  {0.5: 1.0, 2: 1.5, 3: 3.0, 4: 5.0, 5: 20.0, 6: 1000.0},
	7:  {0.5: 0.5, 2: 1.0, 3: 2.0, 4: 5.0, 5: 20.0, 6: 100.0, 7: 3000.0},
	8:  {0.5: 0.5, 2: 1.0, 3: 2.0, 4: 4.0, 5: 10.0, 6: 50.0, 7: 500.0, 8: 10000.0},
	9:  {0.5: 0.5, 2: 0.5, 3: 1.5, 4: 3.0, 5: 5.0, 6: 20.0, 7: 200.0, 8: 2000.0, 9: 20000.0},
	10: {0.5: 0.5, 2: 0.5, 3: 1.0, 4: 2.0, 5: 5.0, 6: 10.0, 7: 50.0, 8: 500.0, 9: 5000.0, 10: 50000.0},
}

func (s *LotteryService) PlayKeno(userID uuid.UUID, betAmount float64, spots int, selected []int, clientSeed string) (*KenoResult, error) {
	if spots < 1 || spots > 10 {
		return nil, fmt.Errorf("spots must be between 1 and 10")
	}

	if len(selected) != spots {
		return nil, fmt.Errorf("must select exactly %d spots", spots)
	}

	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	// Draw 20 numbers
	drawn := make([]int, 0, 20)
	used := make(map[int]bool)

	for len(drawn) < 20 {
		num := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+len(drawn), 80) + 1
		if !used[num] {
			used[num] = true
			drawn = append(drawn, num)
		}
	}

	// Count matches
	matches := 0
	for _, sel := range selected {
		for _, d := range drawn {
			if sel == d {
				matches++
				break
			}
		}
	}

	// Calculate payout
	multiplier := 0.0
	paytable := kenoPaytable[spots]
	
	// Find best payout
	for matchCount, mult := range paytable {
		if matchCount == matches {
			multiplier = mult
			break
		}
		// Check for "catch" (fewer matches but still wins)
		if matchCount == 0 && matches > 0 && matches >= spots-1 {
			if mult > multiplier {
				multiplier = mult
			}
		}
	}

	winAmount := betAmount * multiplier

	result := &KenoResult{
		GameID:     uuid.New().String(),
		Spots:      spots,
		Selected:   selected,
		Drawn:      drawn,
		Matches:    matches,
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: multiplier,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	// Record bet
	s.recordLotteryBet(userID, "keno", betAmount, winAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ BINGO ============

type BingoResult struct {
	GameID       string  `json:"game_id"`
	Card         [][]int `json:"card"`
	Drawn        []int   `json:"drawn"`
	Pattern      string  `json:"pattern"`
	Complete     bool    `json:"complete"`
	PatternName  string  `json:"pattern_name"`
	BetAmount    float64 `json:"bet_amount"`
	WinAmount    float64 `json:"win_amount"`
	ServerSeed   string  `json:"server_seed"`
	ClientSeed   string  `json:"client_seed"`
	Nonce        int     `json:"nonce"`
}

var bingoPatterns = map[string][][][2]int{
	"four_corners": {{{0, 0}, {0, 4}, {4, 0}, {4, 4}}},
	"blackout":     {{{-1, -1}}, // All cells
	"line_row":     {{{0, -1}}, {{1, -1}}, {{2, -1}}, {{3, -1}}, {{4, -1}}}, // Any row
	"line_col":     {{{-1, 0}}, {{-1, 1}}, {{-1, 2}}, {{-1, 3}}, {{-1, 4}}}, // Any column
	"x_pattern":    {{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {0, 4}, {1, 3}, {2, 2}, {3, 1}, {4, 0}}},
}

func (s *LotteryService) PlayBingo(userID uuid.UUID, betAmount float64, pattern string, clientSeed string) (*BingoResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	// Generate 5x5 bingo card
	card := make([][]int, 5)
	for i := range card {
		card[i] = make([]int, 5)
	}

	// Fill with random unique numbers (1-75)
	used := make(map[int]bool)
	num := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			// Center is free
			if i == 2 && j == 2 {
				card[i][j] = 0
				continue
			}
			for {
				num = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+num, 75) + 1
				if !used[num] {
					used[num] = true
					card[i][j] = num
					break
				}
				num++
			}
		}
	}

	// Draw numbers
	drawn := make([]int, 0, 75)
	for len(drawn) < 75 {
		num := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+len(drawn)+100, 75) + 1
		if !contains(drawn, num) {
			drawn = append(drawn, num)
		}
	}

	result := &BingoResult{
		GameID:      uuid.New().String(),
		Card:        card,
		Drawn:       drawn[:30], // First 30 draws
		Pattern:     pattern,
		Complete:    false,
		BetAmount:   betAmount,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed:  seeds.ClientSeed,
		Nonce:       seeds.Nonce,
	}

	// Check if pattern is complete (simplified - just check if any line can be formed)
	result.Complete = checkBingoPattern(card, pattern)

	if result.Complete {
		result.WinAmount = betAmount * 1000
		result.PatternName = getPatternName(pattern)
	}

	s.recordLotteryBet(userID, "bingo", betAmount, result.WinAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

func checkBingoPattern(card [][]int, pattern string) bool {
	switch pattern {
	case "four_corners":
		return card[0][0] > 0 && card[0][4] > 0 && card[4][0] > 0 && card[4][4] > 0
	case "line_row", "line_col", "line_any":
		// Check rows
		for i := 0; i < 5; i++ {
			complete := true
			for j := 0; j < 5; j++ {
				if i == 2 && j == 2 {
					continue // Free space
				}
				if card[i][j] == 0 {
					complete = false
					break
				}
			}
			if complete {
				return true
			}
		}
		// Check columns
		for j := 0; j < 5; j++ {
			complete := true
			for i := 0; i < 5; i++ {
				if i == 2 && j == 2 {
					continue
				}
				if card[i][j] == 0 {
					complete = false
					break
				}
			}
			if complete {
				return true
			}
		}
	}
	return false
}

func getPatternName(pattern string) string {
	names := map[string]string{
		"four_corners": "Four Corners",
		"blackout":     "Blackout",
		"line_row":     "Any Row",
		"line_col":     "Any Column",
		"line_any":     "Any Line",
		"x_pattern":   "X Pattern",
	}
	return names[pattern]
}

// ============ SCRATCH CARDS ============

type ScratchCardResult struct {
	GameID      string   `json:"game_id"`
	GameType    string   `json:"game_type"`
	Symbols     [][]string `json:"symbols"`
	MatchCount  int      `json:"match_count"`
	MatchSymbol string    `json:"match_symbol"`
	BetAmount   float64  `json:"bet_amount"`
	WinAmount   float64  `json:"win_amount"`
	ServerSeed  string   `json:"server_seed"`
	ClientSeed  string   `json:"client_seed"`
	Nonce       int      `json:"nonce"`
}

var scratchCardSymbols = []string{"🍒", "🍋", "🍊", "🍇", "💎", "⭐", "7️⃣", "🎁"}
var scratchCardPayouts = map[string]float64{
	"🍒":  2.0,
	"🍋":  2.0,
	"🍊":  3.0,
	"🍇":  5.0,
	"💎":  10.0,
	"⭐":  20.0,
	"7️⃣":  50.0,
	"🎁":  100.0,
}

func (s *LotteryService) PlayScratchCard(userID uuid.UUID, gameType string, betAmount float64, clientSeed string) (*ScratchCardResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	// Generate 3x3 grid
	symbols := make([][]string, 3)
	for i := 0; i < 3; i++ {
		symbols[i] = make([]string, 3)
		for j := 0; j < 3; j++ {
			idx := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+i*3+j, len(scratchCardSymbols))
			symbols[i][j] = scratchCardSymbols[idx]
		}
	}

	// Find matches
	matchCount := 0
	matchSymbol := ""

	// Check rows
	for i := 0; i < 3; i++ {
		if symbols[i][0] == symbols[i][1] && symbols[i][1] == symbols[i][2] {
			matchCount += 3
			matchSymbol = symbols[i][0]
		}
	}

	// Check columns
	for j := 0; j < 3; j++ {
		if symbols[0][j] == symbols[1][j] && symbols[1][j] == symbols[2][j] {
			matchCount += 3
			if matchSymbol == "" {
				matchSymbol = symbols[0][j]
			}
		}
	}

	// Calculate win
	winAmount := 0.0
	if matchSymbol != "" {
		winAmount = betAmount * scratchCardPayouts[matchSymbol]
	}

	result := &ScratchCardResult{
		GameID:      uuid.New().String(),
		GameType:    gameType,
		Symbols:     symbols,
		MatchCount:  matchCount,
		MatchSymbol: matchSymbol,
		BetAmount:   betAmount,
		WinAmount:   winAmount,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed:  seeds.ClientSeed,
		Nonce:       seeds.Nonce,
	}

	s.recordLotteryBet(userID, "scratch_"+gameType, betAmount, winAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ MEGA BALL ============

type MegaBallResult struct {
	GameID        string `json:"game_id"`
	MainNumbers   []int  `json:"main_numbers"`
	MegaBall     int    `json:"mega_ball"`
	BetAmount    float64 `json:"bet_amount"`
	MainMatches  int     `json:"main_matches"`
	MegaMatch   bool    `json:"mega_match"`
	WinAmount    float64 `json:"win_amount"`
	ServerSeed   string  `json:"server_seed"`
	ClientSeed   string  `json:"client_seed"`
	Nonce        int     `json:"nonce"`
}

var megaBallPayouts = map[string]float64{
	"0+0":    2.0,
	"1+0":    4.0,
	"2+0":    10.0,
	"3+0":    100.0,
	"4+0":    10000.0,
	"5+0":    1000000.0,
	"0+1":    4.0,
	"1+1":    10.0,
	"2+1":    50.0,
	"3+1":    500.0,
	"4+1":    50000.0,
	"5+1":    10000000.0,
}

func (s *LotteryService) PlayMegaBall(userID uuid.UUID, betAmount float64, selected []int, megaBall int, clientSeed string) (*MegaBallResult, error) {
	if len(selected) != 5 {
		return nil, fmt.Errorf("must select exactly 5 main numbers")
	}

	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	// Draw main numbers
	mainNumbers := make([]int, 0, 5)
	used := make(map[int]bool)

	for len(mainNumbers) < 5 {
		num := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+len(mainNumbers), 70) + 1
		if !used[num] {
			used[num] = true
			mainNumbers = append(mainNumbers, num)
		}
	}

	// Draw mega ball
	mb := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, 100, 25) + 1

	// Sort for comparison
	sort.Ints(mainNumbers)
	sort.Ints(selected)

	// Count matches
	mainMatches := 0
	for _, sel := range selected {
		for _, num := range mainNumbers {
			if sel == num {
				mainMatches++
				break
			}
		}
	}

	megaMatch := (mb == megaBall)

	// Calculate win
	key := fmt.Sprintf("%d+%t", mainMatches, megaMatch)
	multiplier := megaBallPayouts[key]
	winAmount := betAmount * multiplier

	result := &MegaBallResult{
		GameID:       uuid.New().String(),
		MainNumbers:  mainNumbers,
		MegaBall:     mb,
		BetAmount:    betAmount,
		MainMatches:  mainMatches,
		MegaMatch:    megaMatch,
		WinAmount:    winAmount,
		ServerSeed:   seeds.ServerSeed,
		ClientSeed:   seeds.ClientSeed,
		Nonce:        seeds.Nonce,
	}

	s.recordLotteryBet(userID, "megaball", betAmount, winAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ INSTANT WIN ============

type InstantWinResult struct {
	GameID     string  `json:"game_id"`
	GameType   string  `json:"game_type"`
	Numbers    []int   `json:"numbers"`
	TargetSum  int     `json:"target_sum"`
	ActualSum  int     `json:"actual_sum"`
	Win        bool    `json:"win"`
	BetAmount  float64 `json:"bet_amount"`
	WinAmount  float64 `json:"win_amount"`
	ServerSeed string  `json:"server_seed"`
	ClientSeed string  `json:"client_seed"`
	Nonce      int     `json:"nonce"`
}

func (s *LotteryService) PlayInstantWin(userID uuid.UUID, gameType string, betAmount float64, targetSum int, clientSeed string) (*InstantWinResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	// Generate 20 numbers (1-100)
	numbers := make([]int, 20)
	sum := 0

	for i := range numbers {
		numbers[i] = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+i, 100) + 1
		sum += numbers[i]
	}

	win := false
	if gameType == "over" {
		win = sum > targetSum
	} else if gameType == "under" {
		win = sum < targetSum
	} else {
		win = sum == targetSum
	}

	multiplier := 0.0
	if win {
		// Calculate based on difficulty
		diff := float64(targetSum - sum)
		if diff < 0 {
			diff = -diff
		}
		// Closer to target = higher multiplier
		if diff <= 10 {
			multiplier = 10.0
		} else if diff <= 25 {
			multiplier = 5.0
		} else if diff <= 50 {
			multiplier = 2.0
		} else {
			multiplier = 1.5
		}
	}

	winAmount := betAmount * multiplier

	result := &InstantWinResult{
		GameID:     uuid.New().String(),
		GameType:   gameType,
		Numbers:    numbers,
		TargetSum:  targetSum,
		ActualSum:  sum,
		Win:        win,
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	s.recordLotteryBet(userID, "instant_"+gameType, betAmount, winAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ HELPER FUNCTIONS ============

func (s *LotteryService) GenerateSeeds(clientSeed string) (*Seeds, error) {
	serverSeedBytes := make([]byte, 32)
	if _, err := rand.Read(serverSeedBytes); err != nil {
		return nil, fmt.Errorf("failed to generate server seed: %w", err)
	}

	serverSeed := hex.EncodeToString(serverSeedBytes)
	serverSeedHash := s.HashSeed(serverSeed)

	if clientSeed == "" {
		clientSeedBytes := make([]byte, 16)
		rand.Read(clientSeedBytes)
		clientSeed = hex.EncodeToString(clientSeedBytes)
	}

	return &Seeds{
		ServerSeed:     serverSeed,
		ServerSeedHash: serverSeedHash,
		ClientSeed:     clientSeed,
		Nonce:          0,
	}, nil
}

func (s *LotteryService) HashSeed(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])
}

func (s *LotteryService) GenerateOutcome(serverSeed, clientSeed string, nonce int, max int) int {
	combined := fmt.Sprintf("%s:%s:%d", serverSeed, clientSeed, nonce)
	hash := sha256.Sum256([]byte(combined))
	result := new(big.Int).SetBytes(hash[:])
	return int(result.Mod(result, big.NewInt(int64(max))).Int64())
}

func (s *LotteryService) recordLotteryBet(userID uuid.UUID, gameType string, betAmount float64, winAmount float64, serverSeed string, clientSeed string, nonce int) {
	profit := winAmount - betAmount
	if winAmount > 0 {
		profit = winAmount - betAmount
	} else {
		profit = -betAmount
	}

	bet := models.Bet{
		UserID:     userID,
		GameType:   gameType,
		BetAmount:  betAmount,
		WinAmount:  winAmount,
		Multiplier: winAmount / betAmount,
		Profit:     profit,
		Status:     "settled",
		ServerSeed: serverSeed,
		ClientSeed: clientSeed,
		Nonce:      nonce,
		IsVerified: true,
	}

	s.db.Create(&bet)
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// GetLotteryHistory returns user's lottery game history
func (s *LotteryService) GetLotteryHistory(userID uuid.UUID, gameType string, limit int) ([]models.Bet, error) {
	var bets []models.Bet
	query := s.db.Where("user_id = ? AND game_type LIKE ?", userID, gameType+"%")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&bets).Error
	return bets, err
}

// GetLotteryStats returns lottery game statistics
func (s *LotteryService) GetLotteryStats(userID uuid.UUID) (map[string]interface{}, error) {
	var totalBets int64
	var totalWagered float64
	var totalWon float64

	s.db.Model(&models.Bet{}).Where("user_id = ? AND game_type IN ?", userID, 
		[]string{"keno", "bingo", "scratch", "megaball", "instant"}).Count(&totalBets)

	s.db.Model(&models.Bet{}).Where("user_id = ? AND game_type IN ?", userID,
		[]string{"keno", "bingo", "scratch", "megaball", "instant"}).Select("COALESCE(SUM(bet_amount), 0)").Scan(&totalWagered)

	s.db.Model(&models.Bet{}).Where("user_id = ? AND game_type IN ? AND win_amount > 0", userID,
		[]string{"keno", "bingo", "scratch", "megaball", "instant"}).Select("COALESCE(SUM(win_amount), 0)").Scan(&totalWon)

	return map[string]interface{}{
		"total_bets":      totalBets,
		"total_wagered":   totalWagered,
		"total_won":       totalWon,
		"net_profit":      totalWon - totalWagered,
	}, nil
}
