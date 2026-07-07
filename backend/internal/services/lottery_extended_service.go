package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// ============ Lottery Games Service ============

type LotteryService struct{}

func NewLotteryService() *LotteryService {
	return &LotteryService{}
}

type LotteryGame struct {
	ID           string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // keno, bingo, scratch, lottery
	MinPlayers  int     `json:"min_players"`
	MaxPlayers  int     `json:"max_players"`
	DrawFrequency string `json:"draw_frequency"`
	Jackpot     float64 `json:"jackpot"`
}

type LotteryTicket struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	GameID     string    `json:"game_id"`
	Numbers    []int     `json:"numbers"`
	Stake      float64   `json:"stake"`
	Prize     float64   `json:"prize"`
	DrawTime   time.Time `json:"draw_time"`
	Status     string    `json:"status"` // pending, won, lost
	CreatedAt  time.Time `json:"created_at"`
}

type KenoResult struct {
	DrawNumber   int     `json:"draw_number"`
	DrawTime    time.Time `json:"draw_time"`
	DrawnNumbers []int   `json:"drawn_numbers"`
	PickedNumbers []int  `json:"picked_numbers"`
	Matches     int     `json:"matches"`
	Prize       float64 `json:"prize"`
	Payout      float64 `json:"payout"`
}

// Keno paytable (picks matched vs prize multiplier)
var kenoPaytable = map[int]float64{
	1:  2.0,
	2:  5.0,
	3:  10.0,
	4:  50.0,
	5:  100.0,
	6:  500.0,
	7:  1000.0,
	8:  5000.0,
	9:  10000.0,
	10: 50000.0,
}

// Scratch card game
type ScratchCard struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Cost        float64  `json:"cost"`
	PrizeTiers []PrizeTier `json:"prize_tiers"`
	ImageURL    string   `json:"image_url"`
}

type PrizeTier struct {
	Prize   float64 `json:"prize"`
	Count   int     `json:"count"`
	Odds    float64 `json:"odds"`
}

type ScratchResult struct {
	CardID      string   `json:"card_id"`
	CardName    string   `json:"card_name"`
	Symbols     []string `json:"symbols"`
	MatchCount  int      `json:"match_count"`
	PrizeWon   float64  `json:"prize_won"`
	IsWinner   bool     `json:"is_winner"`
}

// Bingo game
type BingoGame struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // 75-ball, 90-ball
	CardPrice   float64   `json:"card_price"`
	MinCards   int       `json:"min_cards"`
	MaxCards   int       `json:"max_cards"`
	Jackpot     float64   `json:"jackpot"`
	PrizePool  float64   `json:"prize_pool"`
	Players    int       `json:"players"`
	Status     string    `json:"status"` // waiting, playing, finished
	StartTime  time.Time `json:"start_time"`
}

type BingoCard struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	GameID    string `json:"game_id"`
	Numbers   [][]int `json:"numbers"` // For 90-ball: 3x9 grid
	Matched   []int   `json:"matched"`
	Pattern   string  `json:"pattern"` // line, two-lines, full-house
	Prize     float64 `json:"prize"`
	Won       bool    `json:"won"`
}

type BingoCall struct {
	Number   int     `json:"number"`
	Letter   string  `json:"letter"` // B, I, N, G, O
	Time     time.Time `json:"time"`
}

// Generate Keno numbers
func (s *LotteryService) GenerateKenoNumbers(count int) []int {
	numbers := make([]int, 0)
	used := make(map[int]bool)
	
	for len(numbers) < count {
		num := rand.Intn(80) + 1
		if !used[num] {
			used[num] = true
			numbers = append(numbers, num)
		}
	}
	return numbers
}

// Play Keno
func (s *LotteryService) PlayKeno(userID string, picks int, stake float64, numbers []int) (*KenoResult, error) {
	if picks < 1 || picks > 10 {
		return nil, fmt.Errorf("picks must be between 1 and 10")
	}
	
	if len(numbers) != picks {
		numbers = s.GenerateKenoNumbers(picks)
	}
	
	// Draw winning numbers
	winningNumbers := s.GenerateKenoNumbers(20)
	
	// Count matches
	matches := 0
	for _, num := range numbers {
		for _, win := range winningNumbers {
			if num == win {
				matches++
				break
			}
		}
	}
	
	// Calculate prize
	multiplier := kenoPaytable[matches]
	prize := stake * multiplier
	
	return &KenoResult{
		DrawNumber:   rand.Intn(1000000),
		DrawTime:    time.Now(),
		DrawnNumbers: winningNumbers,
		PickedNumbers: numbers,
		Matches:     matches,
		Prize:       prize,
		Payout:      multiplier,
	}, nil
}

// Buy scratch card
func (s *LotteryService) BuyScratchCard(userID, cardID string, cost float64) (*ScratchResult, error) {
	// Define scratch cards
	cards := map[string]ScratchCard{
		"lucky7s": {
			ID:   "lucky7s",
			Name: "Lucky 7s",
			Cost: 1.0,
			PrizeTiers: []PrizeTier{
				{Prize: 10000, Count: 1, Odds: 0.0001},
				{Prize: 1000, Count: 10, Odds: 0.001},
				{Prize: 100, Count: 100, Odds: 0.01},
				{Prize: 10, Count: 1000, Odds: 0.1},
				{Prize: 2, Count: 10000, Odds: 1.0},
			},
		},
		"gold_rush": {
			ID:   "gold_rush",
			Name: "Gold Rush",
			Cost: 0.5,
			PrizeTiers: []PrizeTier{
				{Prize: 5000, Count: 1, Odds: 0.0001},
				{Prize: 500, Count: 10, Odds: 0.001},
				{Prize: 50, Count: 100, Odds: 0.01},
				{Prize: 5, Count: 1000, Odds: 0.1},
				{Prize: 0.5, Count: 10000, Odds: 1.0},
			},
		},
		"diamond_dazzle": {
			ID:   "diamond_dazzle",
			Name: "Diamond Dazzle",
			Cost: 2.0,
			PrizeTiers: []PrizeTier{
				{Prize: 50000, Count: 1, Odds: 0.00005},
				{Prize: 5000, Count: 5, Odds: 0.0005},
				{Prize: 500, Count: 50, Odds: 0.005},
				{Prize: 50, Count: 500, Odds: 0.05},
				{Prize: 5, Count: 5000, Odds: 0.5},
			},
		},
	}
	
	card, ok := cards[cardID]
	if !ok {
		return nil, fmt.Errorf("invalid card type")
	}
	
	// Generate symbols (3x3 grid = 9 symbols)
	symbols := []string{}
	symbolPool := []string{"🍒", "🍋", "🍇", "💎", "⭐", "🔔", "7️⃣", "💰"}
	
	for i := 0; i < 9; i++ {
		symbols = append(symbols, symbolPool[rand.Intn(len(symbolPool))])
	}
	
	// Check for win (center symbol matches any other)
	matchCount := 0
	centerSymbol := symbols[4]
	for i, sym := range symbols {
		if i != 4 && sym == centerSymbol {
			matchCount++
		}
	}
	
	// Calculate prize based on matches
	prizeWon := 0.0
	isWinner := false
	
	if matchCount >= 2 {
		isWinner = true
		if matchCount >= 4 {
			prizeWon = card.Cost * 50
		} else if matchCount >= 3 {
			prizeWon = card.Cost * 10
		} else {
			prizeWon = card.Cost * 2
		}
	}
	
	return &ScratchResult{
		CardID:     cardID,
		CardName:   card.Name,
		Symbols:    symbols,
		MatchCount: matchCount,
		PrizeWon:   prizeWon,
		IsWinner:   isWinner,
	}, nil
}

// Create Bingo game
func (s *LotteryService) CreateBingoGame(name string, gameType string, cardPrice float64) *BingoGame {
	return &BingoGame{
		ID:          uuid.New().String(),
		Name:        name,
		Type:        gameType,
		CardPrice:   cardPrice,
		MinCards:    1,
		MaxCards:    10,
		Jackpot:     cardPrice * 1000,
		PrizePool:   0,
		Players:     0,
		Status:      "waiting",
		StartTime:   time.Now().Add(5 * time.Minute),
	}
}

// Generate Bingo card
func (s *LotteryService) GenerateBingoCard(userID, gameID string, gameType string) *BingoCard {
	card := &BingoCard{
		ID:        uuid.New().String(),
		UserID:    userID,
		GameID:    gameID,
		Numbers:   make([][]int, 0),
		Matched:   make([]int, 0),
		Pattern:   "",
		Prize:     0,
		Won:       false,
	}
	
	if gameType == "75-ball" {
		// 5x5 grid with free space in center
		grid := make([]int, 25)
		used := make(map[int]bool)
		
		for i := 0; i < 25; i++ {
			if i == 12 { // Center is free
				grid[i] = 0 // Free space
				continue
			}
			
			num := rand.Intn(75) + 1
			for used[num] {
				num = rand.Intn(75) + 1
			}
			used[num] = true
			grid[i] = num
		}
		
		// Convert to 2D
		for row := 0; row < 5; row++ {
			card.Numbers = append(card.Numbers, grid[row*5:(row+1)*5])
		}
	} else {
		// 90-ball: 3x9 grid with 2 empty per row
		for row := 0; row < 3; row++ {
			rowNums := make([]int, 9)
			used := make(map[int]bool)
			emptyCount := 0
			
			for col := 0; col < 9; col++ {
				// Randomly leave 2 spots empty per row
				if rand.Float32() < 0.22 && emptyCount < 2 {
					rowNums[col] = 0
					emptyCount++
					continue
				}
				
				// Generate number in column range
				minNum := col*10 + 1
				if col == 8 {
					minNum = 80
				}
				maxNum := minNum + 9
				
				num := minNum + rand.Intn(maxNum-minNum)
				for used[num] || num > 90 {
					num = minNum + rand.Intn(maxNum-minNum)
				}
				used[num] = true
				rowNums[col] = num
			}
			card.Numbers = append(card.Numbers, rowNums)
		}
	}
	
	return card
}

// Call next number in Bingo
func (s *LotteryService) CallBingoNumber(gameType string) *BingoCall {
	number := rand.Intn(75) + 1
	letter := ""
	
	if gameType == "75-ball" {
		if number <= 15 {
			letter = "B"
		} else if number <= 30 {
			letter = "I"
		} else if number <= 45 {
			letter = "N"
		} else if number <= 60 {
			letter = "G"
		} else {
			letter = "O"
		}
	} else {
		letter = "" // 90-ball uses just numbers
	}
	
	return &BingoCall{
		Number: number,
		Letter: letter,
		Time:   time.Now(),
	}
}

// Check Bingo card for win
func (s *LotteryService) CheckBingoWin(card *BingoCard, calledNumbers []int, gameType string) (bool, string) {
	calledMap := make(map[int]bool)
	for _, n := range calledNumbers {
		calledMap[n] = true
	}
	
	if gameType == "75-ball" {
		// Check for different patterns
		patterns := []string{"full-house", "four-corners", "line"}
		
		// Check full house
		fullHouse := true
		for row := 0; row < 5; row++ {
			for col := 0; col < 5; col++ {
				if row == 2 && col == 2 { // Skip free space
					continue
				}
				if !calledMap[card.Numbers[row][col]] {
					fullHouse = false
					break
				}
			}
		}
		if fullHouse {
			return true, "full-house"
		}
		
		// Check four corners
		if calledMap[card.Numbers[0][0]] &&
			calledMap[card.Numbers[0][4]] &&
			calledMap[card.Numbers[4][0]] &&
			calledMap[card.Numbers[4][4]] {
			return true, "four-corners"
		}
		
		// Check for any complete line
		for row := 0; row < 5; row++ {
			lineComplete := true
			for col := 0; col < 5; col++ {
				if row == 2 && col == 2 { // Skip free space
					continue
				}
				if !calledMap[card.Numbers[row][col]] {
					lineComplete = false
					break
				}
			}
			if lineComplete {
				return true, "line"
			}
		}
	}
	
	return false, ""
}

// Lottery statistics
type LotteryStats struct {
	TotalPlayers int     `json:"total_players"`
	TotalWinnings float64 `json:"total_winnings"`
	BiggestWin   float64 `json:"biggest_win"`
	RecentWins   []RecentWin `json:"recent_wins"`
}

type RecentWin struct {
	UserID   string  `json:"user_id"`
	Game     string  `json:"game"`
	Prize    float64 `json:"prize"`
	TimeAgo  string  `json:"time_ago"`
}

func (s *LotteryService) GetStats() *LotteryStats {
	return &LotteryStats{
		TotalPlayers: 12543,
		TotalWinnings: 543210.50,
		BiggestWin: 50000.00,
		RecentWins: []RecentWin{
			{UserID: "user***123", Game: "Lucky 7s", Prize: 5000, TimeAgo: "2 min ago"},
			{UserID: "user***456", Game: "Keno", Prize: 2500, TimeAgo: "5 min ago"},
			{UserID: "user***789", Game: "Bingo", Prize: 10000, TimeAgo: "10 min ago"},
			{UserID: "user***012", Game: "Diamond Dazzle", Prize: 500, TimeAgo: "15 min ago"},
			{UserID: "user***345", Game: "Gold Rush", Prize: 100, TimeAgo: "20 min ago"},
		},
	}
}
