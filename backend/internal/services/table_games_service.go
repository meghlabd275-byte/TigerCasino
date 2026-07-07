package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// TableGamesService handles all table game operations
type TableGamesService struct {
	db *gorm.DB
}

func NewTableGamesService(db *gorm.DB) *TableGamesService {
	return &TableGamesService{db: db}
}

// ============ Card Management ============

// Card represents a playing card
type Card struct {
	Suit  string `json:"suit"`
	Rank  string `json:"rank"`
	Value int    `json:"value"`
}

// Deck represents a shoe of cards
type Deck struct {
	Cards   []Card `json:"cards"`
	Decks   int    `json:"decks"`
	CutCard int    `json:"cut_card"`
}

var suits = []string{"♠", "♥", "♦", "♣"}
var ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func createDeck(decks int) *Deck {
	deck := &Deck{
		Cards: make([]Card, 0, 52*decks),
		Decks: decks,
	}

	for d := 0; d < decks; d++ {
		for _, suit := range suits {
			for i, rank := range ranks {
				value := i + 1
				if value > 10 {
					value = 10
				}
				deck.Cards = append(deck.Cards, Card{
					Suit:  suit,
					Rank:  rank,
					Value: value,
				})
			}
		}
	}

	// Shuffle
	rand.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})

	// Set cut card (typically around 75% through)
	deck.CutCard = len(deck.Cards) * 3 / 4

	return deck
}

func (d *Deck) Draw() (Card, bool) {
	if len(d.Cards) == 0 {
		return Card{}, false
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card, true
}

func (d *Deck) NeedsShuffle() bool {
	return len(d.Cards) <= d.CutCard
}

// ============ Blackjack ============

// BlackjackHand represents a player's hand
type BlackjackHand struct {
	Cards     []Card `json:"cards"`
	Bet       float64 `json:"bet"`
	IsDealer  bool    `json:"is_dealer"`
	Finished  bool    `json:"finished"`
	Blackjack bool    `json:"blackjack"`
	Split     bool    `json:"split"`
}

func (h *BlackjackHand) Score() int {
	score := 0
	aces := 0

	for _, card := range h.Cards {
		score += card.Value
		if card.Rank == "A" {
			aces++
		}
	}

	for aces > 0 && score > 21 {
		score -= 10
		aces--
	}

	return score
}

func (h *BlackjackHand) CanSplit() bool {
	return len(h.Cards) == 2 && h.Cards[0].Value == h.Cards[1].Value && !h.Split
}

func (h *BlackjackHand) CanDouble() bool {
	return len(h.Cards) == 2 && !h.Split
}

// BlackjackResult represents the result of a blackjack game
type BlackjackResult struct {
	PlayerHand   BlackjackHand   `json:"player_hand"`
	DealerHand   BlackjackHand   `json:"dealer_hand"`
	PlayerScore  int             `json:"player_score"`
	DealerScore  int             `json:"dealer_score"`
	WinAmount    float64         `json:"win_amount"`
	Result       string          `json:"result"` // "win", "lose", "push", "blackjack"
	Payout       float64         `json:"payout"` // 1.0, 1.5, 2.0
	ServerSeed   string          `json:"server_seed"`
	ClientSeed   string          `json:"client_seed"`
	Nonce        int             `json:"nonce"`
}

// PlayBlackjack plays a hand of blackjack
func (s *TableGamesService) PlayBlackjack(userID uuid.UUID, betAmount float64, clientSeed string, action string) (*BlackjackResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	deck := createDeck(6) // 6-deck shoe

	// Deal initial hands
	playerHand := BlackjackHand{Bet: betAmount}
	dealerHand := BlackjackHand{IsDealer: true}

	// Player gets two cards
	card1, ok := deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	playerHand.Cards = append(playerHand.Cards, card1)

	card2, ok := deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	dealerHand.Cards = append(dealerHand.Cards, card2)

	card3, ok := deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	playerHand.Cards = append(playerHand.Cards, card3)

	card4, ok := deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	dealerHand.Cards = append(dealerHand.Cards, card4)

	// Check for blackjack
	playerBlackjack := playerHand.Score() == 21
	dealerBlackjack := dealerHand.Score() == 21

	if playerBlackjack && dealerBlackjack {
		// Both have blackjack - push
		playerHand.Blackjack = true
		dealerHand.Blackjack = true
		playerHand.Finished = true
		dealerHand.Finished = true

		return &BlackjackResult{
			PlayerHand:   playerHand,
			DealerHand:   dealerHand,
			PlayerScore:  playerHand.Score(),
			DealerScore:  dealerHand.Score(),
			WinAmount:    0,
			Result:       "push",
			Payout:       1.0,
			ServerSeed:   seeds.ServerSeed,
			ClientSeed:   seeds.ClientSeed,
			Nonce:        seeds.Nonce,
		}, nil
	}

	if playerBlackjack {
		playerHand.Blackjack = true
		playerHand.Finished = true

		// Dealer reveals hole card and plays
		dealerHand.Finished = true
		for dealerHand.Score() < 17 {
			card, ok := deck.Draw()
			if !ok {
				break
			}
			dealerHand.Cards = append(dealerHand.Cards, card)
		}

		dealerScore := dealerHand.Score()

		var winAmount float64
		var result string
		var payout float64

		if dealerScore != 21 {
			winAmount = betAmount * 1.5
			result = "blackjack"
			payout = 1.5
		} else {
			winAmount = 0
			result = "push"
			payout = 1.0
		}

		// Record bet
		s.recordBet(userID, "blackjack", betAmount, winAmount, winAmount-betAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

		return &BlackjackResult{
			PlayerHand:   playerHand,
			DealerHand:   dealerHand,
			PlayerScore:  playerHand.Score(),
			DealerScore:  dealerScore,
			WinAmount:    winAmount,
			Result:       result,
			Payout:       payout,
			ServerSeed:   seeds.ServerSeed,
			ClientSeed:   seeds.ClientSeed,
			Nonce:        seeds.Nonce,
		}, nil
	}

	// Player stands (simplified - no hit/split/double for now)
	playerHand.Finished = true

	// Dealer plays
	dealerHand.Finished = true
	for dealerHand.Score() < 17 {
		card, ok := deck.Draw()
		if !ok {
			break
		}
		dealerHand.Cards = append(dealerHand.Cards, card)
	}

	playerScore := playerHand.Score()
	dealerScore := dealerHand.Score()

	var winAmount float64
	var result string
	var payout float64

	if dealerScore > 21 || playerScore > dealerScore {
		winAmount = betAmount * 2
		result = "win"
		payout = 2.0
	} else if playerScore < dealerScore {
		winAmount = 0
		result = "lose"
		payout = 0
	} else {
		winAmount = betAmount
		result = "push"
		payout = 1.0
	}

	// Record bet
	s.recordBet(userID, "blackjack", betAmount, winAmount, winAmount-betAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return &BlackjackResult{
		PlayerHand:   playerHand,
		DealerHand:   dealerHand,
		PlayerScore:  playerScore,
		DealerScore:  dealerScore,
		WinAmount:    winAmount,
		Result:       result,
		Payout:       payout,
		ServerSeed:   seeds.ServerSeed,
		ClientSeed:   seeds.ClientSeed,
		Nonce:        seeds.Nonce,
	}, nil
}

// ============ Roulette ============

// RouletteBet represents a bet on the roulette table
type RouletteBet struct {
	Type     string   `json:"type"` // "straight", "split", "street", "corner", "line", "column", "dozen", "red", "black", "even", "odd", "1-18", "19-36"
	Numbers  []int    `json:"numbers"`
	Amount   float64  `json:"amount"`
	Multiplier float64 `json:"multiplier"`
}

// RouletteResult represents the result of a roulette spin
type RouletteResult struct {
	Number      int           `json:"number"`
	Color      string        `json:"color"`
	Parity     string        `json:"parity"`
	Range      string        `json:"range"`
	Column     int           `json:"column"`
	Dozen      int           `json:"dozen"`
	Bets       []RouletteBet `json:"bets"`
	WinAmount  float64       `json:"win_amount"`
	TotalWager float64       `json:"total_wager"`
	ServerSeed string        `json:"server_seed"`
	ClientSeed string        `json:"client_seed"`
	Nonce      int           `json:"nonce"`
}

var rouletteNumbers = []int{0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10, 5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26}
var redNumbers = []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
var blackNumbers = []int{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35}

func isRed(n int) bool {
	for _, r := range redNumbers {
		if r == n {
			return true
		}
	}
	return false
}

func isBlack(n int) bool {
	for _, b := range blackNumbers {
		if b == n {
			return true
		}
	}
	return false
}

// SpinRoulette spins the roulette wheel
func (s *TableGamesService) SpinRoulette(userID uuid.UUID, bets []RouletteBet, clientSeed string) (*RouletteResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	// Generate winning number
	winningNumber := s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce, 38)

	result := &RouletteResult{
		Number:      rouletteNumbers[winningNumber],
		Bets:        bets,
		ServerSeed:  seeds.ServerSeed,
		ClientSeed:  seeds.ClientSeed,
		Nonce:       seeds.Nonce,
	}

	// Determine result properties
	if result.Number == 0 {
		result.Color = "green"
		result.Parity = "neither"
		result.Range = "neither"
		result.Column = 0
		result.Dozen = 0
	} else {
		if isRed(result.Number) {
			result.Color = "red"
		} else {
			result.Color = "black"
		}

		if result.Number%2 == 0 {
			result.Parity = "even"
		} else {
			result.Parity = "odd"
		}

		if result.Number <= 18 {
			result.Range = "1-18"
		} else {
			result.Range = "19-36"
		}

		result.Column = (result.Number - 1) % 3
		if result.Column == 0 {
			result.Column = 3
		}

		result.Dozen = (result.Number - 1) / 12
		if result.Dozen == 0 {
			result.Dozen = 1
		} else if result.Dozen == 1 {
			result.Dozen = 2
		} else {
			result.Dozen = 3
		}
	}

	// Calculate winnings
	totalWager := 0.0
	totalWin := 0.0

	for i := range bets {
		bet := &bets[i]
		totalWager += bet.Amount

		for _, num := range bet.Numbers {
			if num == result.Number {
				totalWin += bet.Amount * bet.Multiplier
				break
			}
		}
	}

	result.TotalWager = totalWager
	result.WinAmount = totalWin - totalWager

	// Record bet
	s.recordBet(userID, "roulette", totalWager, totalWin, totalWin-totalWager, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ Baccarat ============

// BaccaratResult represents the result of a baccarat hand
type BaccaratResult struct {
	PlayerCards   []Card     `json:"player_cards"`
	BankerCards   []Card     `json:"banker_cards"`
	PlayerScore   int        `json:"player_score"`
	BankerScore   int        `json:"banker_score"`
	Winner       string     `json:"winner"` // "player", "banker", "tie"
	WinAmount    float64    `json:"win_amount"`
	Payout       float64    `json:"payout"`
	ServerSeed   string     `json:"server_seed"`
	ClientSeed   string     `json:"client_seed"`
	Nonce        int        `json:"nonce"`
}

func baccaratScore(cards []Card) int {
	score := 0
	for _, card := range cards {
		score += card.Value
	}
	return score % 10
}

// PlayBaccarat plays a round of baccarat
func (s *TableGamesService) PlayBaccarat(userID uuid.UUID, betAmount float64, betOn string, clientSeed string) (*BaccaratResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	deck := createDeck(8) // 8-deck shoe

	result := &BaccaratResult{
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	// Deal initial cards
	card, ok := deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	result.PlayerCards = append(result.PlayerCards, card)

	card, ok = deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	result.BankerCards = append(result.BankerCards, card)

	card, ok = deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	result.PlayerCards = append(result.PlayerCards, card)

	card, ok = deck.Draw()
	if !ok {
		return nil, fmt.Errorf("deck exhausted")
	}
	result.BankerCards = append(result.BankerCards, card)

	result.PlayerScore = baccaratScore(result.PlayerCards)
	result.BankerScore = baccaratScore(result.BankerCards)

	// Check for natural
	playerNatural := result.PlayerScore >= 8
	bankerNatural := result.BankerScore >= 8

	if playerNatural || bankerNatural {
		// No more cards
	} else {
		// Player draws third card rules
		playerDraw := false
		if result.PlayerScore <= 5 {
			playerDraw = true
			card, ok = deck.Draw()
			if ok {
				result.PlayerCards = append(result.PlayerCards, card)
				result.PlayerScore = baccaratScore(result.PlayerCards)
			}
		}

		// Banker draws third card rules
		if playerDraw {
			bankerScore := baccaratScore(result.BankerCards)
			if bankerScore <= 2 {
				card, ok = deck.Draw()
				if ok {
					result.BankerCards = append(result.BankerCards, card)
					result.BankerScore = baccaratScore(result.BankerCards)
				}
			} else if bankerScore == 3 && result.PlayerCards[2].Value != 8 {
				card, ok = deck.Draw()
				if ok {
					result.BankerCards = append(result.BankerCards, card)
					result.BankerScore = baccaratScore(result.BankerCards)
				}
			} else if bankerScore == 4 && result.PlayerCards[2].Value >= 2 && result.PlayerCards[2].Value <= 7 {
				card, ok = deck.Draw()
				if ok {
					result.BankerCards = append(result.BankerCards, card)
					result.BankerScore = baccaratScore(result.BankerCards)
				}
			} else if bankerScore == 5 && result.PlayerCards[2].Value >= 4 && result.PlayerCards[2].Value <= 7 {
				card, ok = deck.Draw()
				if ok {
					result.BankerCards = append(result.BankerCards, card)
					result.BankerScore = baccaratScore(result.BankerCards)
				}
			} else if bankerScore == 6 && result.PlayerCards[2].Value >= 6 && result.PlayerCards[2].Value <= 7 {
				card, ok = deck.Draw()
				if ok {
					result.BankerCards = append(result.BankerCards, card)
					result.BankerScore = baccaratScore(result.BankerCards)
				}
			}
		}
	}

	// Determine winner
	if result.PlayerScore > result.BankerScore {
		result.Winner = "player"
	} else if result.BankerScore > result.PlayerScore {
		result.Winner = "banker"
	} else {
		result.Winner = "tie"
	}

	// Calculate payout
	var winAmount float64
	var payout float64

	switch betOn {
	case "player":
		payout = 2.0
		if result.Winner == "player" {
			winAmount = betAmount * payout
		}
	case "banker":
		payout = 1.95 // Banker wins have 5% commission
		if result.Winner == "banker" {
			winAmount = betAmount * payout
		}
	case "tie":
		payout = 9.0
		if result.Winner == "tie" {
			winAmount = betAmount * payout
		}
	}

	result.WinAmount = winAmount
	result.Payout = payout

	// Record bet
	s.recordBet(userID, "baccarat", betAmount, winAmount, winAmount-betAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ Video Poker ============

// VideoPokerHand represents a video poker hand
type VideoPokerHand struct {
	Cards      []Card  `json:"cards"`
	Held       []bool  `json:"held"`
	FinalCards []Card  `json:"final_cards"`
	Payout     float64 `json:"payout"`
}

// VideoPokerResult represents the result of a video poker game
type VideoPokerResult struct {
	InitialHand VideoPokerHand `json:"initial_hand"`
	FinalHand   VideoPokerHand `json:"final_hand"`
	WinAmount   float64       `json:"win_amount"`
	HandRank    string        `json:"hand_rank"` // "royal_flush", "straight_flush", etc.
	ServerSeed  string        `json:"server_seed"`
	ClientSeed  string        `json:"client_seed"`
	Nonce       int           `json:"nonce"`
}

var videoPokerPaytable = map[string]float64{
	"royal_flush":     800,
	"straight_flush":  50,
	"four_of_a_kind":  25,
	"full_house":      9,
	"flush":           6,
	"straight":        4,
	"three_of_a_kind": 3,
	"two_pair":        2,
	"jacks_or_better": 1,
}

func evaluateVideoPoker(cards []Card) (string, float64) {
	if len(cards) != 5 {
		return "nothing", 0
	}

	// Sort cards by value for easier evaluation
	sorted := make([]Card, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	// Check flush
	isFlush := true
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Suit != sorted[0].Suit {
			isFlush = false
			break
		}
	}

	// Check straight
	isStraight := true
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Value != sorted[i-1].Value+1 {
			isStraight = false
			break
		}
	}
	// Check for Ace-low straight
	if !isStraight && sorted[0].Value == 1 && sorted[1].Value == 10 && sorted[2].Value == 11 && sorted[3].Value == 12 && sorted[4].Value == 13 {
		isStraight = true
	}

	// Count values for pairs/trips/quads
	valueCount := make(map[int]int)
	for _, card := range sorted {
		valueCount[card.Value]++
	}

	// Evaluate hand
	if isFlush && isStraight {
		if sorted[0].Value == 1 && sorted[4].Value == 13 {
			return "royal_flush", videoPokerPaytable["royal_flush"]
		}
		return "straight_flush", videoPokerPaytable["straight_flush"]
	}

	hasFourOfAKind := false
	hasThreeOfAKind := false
	hasTwoPair := false
	hasPair := false

	for _, count := range valueCount {
		if count == 4 {
			hasFourOfAKind = true
		} else if count == 3 {
			hasThreeOfAKind = true
		} else if count == 2 {
			if hasPair {
				hasTwoPair = true
			}
			hasPair = true
		}
	}

	if hasFourOfAKind {
		return "four_of_a_kind", videoPokerPaytable["four_of_a_kind"]
	}

	if hasThreeOfAKind && hasPair {
		return "full_house", videoPokerPaytable["full_house"]
	}

	if isFlush {
		return "flush", videoPokerPaytable["flush"]
	}

	if isStraight {
		return "straight", videoPokerPaytable["straight"]
	}

	if hasThreeOfAKind {
		return "three_of_a_kind", videoPokerPaytable["three_of_a_kind"]
	}

	if hasTwoPair {
		return "two_pair", videoPokerPaytable["two_pair"]
	}

	if hasPair {
		// Check for jacks or better
		for val, count := range valueCount {
			if count == 2 && val >= 11 || val == 1 { // J, Q, K, A
				return "jacks_or_better", videoPokerPaytable["jacks_or_better"]
			}
		}
	}

	return "nothing", 0
}

// PlayVideoPoker plays a game of video poker
func (s *TableGamesService) PlayVideoPoker(userID uuid.UUID, betAmount float64, held []bool, clientSeed string) (*VideoPokerResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	deck := createDeck(1)

	result := &VideoPokerResult{
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	// Draw initial 5 cards
	initialHand := VideoPokerHand{
		Cards: make([]Card, 5),
		Held:  held,
	}

	for i := 0; i < 5; i++ {
		card, ok := deck.Draw()
		if !ok {
			return nil, fmt.Errorf("deck exhausted")
		}
		initialHand.Cards[i] = card
	}

	result.InitialHand = initialHand

	// Draw replacement cards for non-held cards
	finalHand := VideoPokerHand{
		Cards: make([]Card, 5),
		Held:  held,
	}

	for i := 0; i < 5; i++ {
		if held[i] {
			finalHand.Cards[i] = initialHand.Cards[i]
		} else {
			card, ok := deck.Draw()
			if !ok {
				return nil, fmt.Errorf("deck exhausted")
			}
			finalHand.Cards[i] = card
		}
	}

	// Evaluate hand
	handRank, multiplier := evaluateVideoPoker(finalHand.Cards)
	result.HandRank = handRank
	result.FinalHand = finalHand
	result.FinalHand.Payout = multiplier
	result.WinAmount = betAmount * multiplier

	// Record bet
	s.recordBet(userID, "video_poker", betAmount, result.WinAmount, result.WinAmount-betAmount, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ Sic Bo ============

// SicBoResult represents the result of a sic bo game
type SicBoResult struct {
	Dice1      int     `json:"dice1"`
	Dice2      int     `json:"dice2"`
	Dice3      int     `json:"dice3"`
	Total      int     `json:"total"`
	Bets       map[string]float64 `json:"bets"`
	WinAmount  float64 `json:"win_amount"`
	ServerSeed string  `json:"server_seed"`
	ClientSeed string  `json:"client_seed"`
	Nonce      int     `json:"nonce"`
}

var sicBoPayouts = map[string]float64{
	"small":    1.0,  // 4-10 except triples
	"big":      1.0,   // 11-17 except triples
	"odd":      1.0,
	"even":     1.0,
	"triple":   180,
	"any_triple": 30,
	"double":   11,
	"4":        62,
	"5":        31,
	"6":        18,
	"7":        12,
	"8":        8,
	"9":        7,
	"10":       6,
	"11":       6,
	"12":       7,
	"13":       8,
	"14":       12,
	"15":       18,
	"16":       31,
	"17":       62,
}

// RollSicBo rolls the sic bo dice
func (s *TableGamesService) RollSicBo(userID uuid.UUID, bets map[string]float64, clientSeed string) (*SicBoResult, error) {
	seeds, err := s.GenerateSeeds(clientSeed)
	if err != nil {
		return nil, err
	}

	result := &SicBoResult{
		Bets:       bets,
		ServerSeed: seeds.ServerSeed,
		ClientSeed: seeds.ClientSeed,
		Nonce:      seeds.Nonce,
	}

	// Roll three dice
	result.Dice1 = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce, 6) + 1
	result.Dice2 = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+1, 6) + 1
	result.Dice3 = s.GenerateOutcome(seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce+2, 6) + 1

	result.Total = result.Dice1 + result.Dice2 + result.Dice3

	// Calculate winnings
	totalWager := 0.0
	totalWin := 0.0

	for betType, amount := range bets {
		totalWager += amount
		win := false

		switch betType {
		case "small":
			// 4-10, not triple
			if result.Total >= 4 && result.Total <= 10 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				win = true
			}
		case "big":
			// 11-17, not triple
			if result.Total >= 11 && result.Total <= 17 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				win = true
			}
		case "odd":
			if result.Total%2 == 1 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				win = true
			}
		case "even":
			if result.Total%2 == 0 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				win = true
			}
		case "triple":
			if result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3 {
				win = true
			}
		case "any_triple":
			if result.Dice1 == result.Dice2 || result.Dice2 == result.Dice3 || result.Dice1 == result.Dice3 {
				win = true
			}
		case "double":
			if result.Dice1 == result.Dice2 || result.Dice2 == result.Dice3 || result.Dice1 == result.Dice3 {
				win = true
			}
		default:
			// Number bets (4-17)
			if fmt.Sprintf("%d", result.Total) == betType {
				win = true
			}
		}

		if win {
			totalWin += amount * sicBoPayouts[betType]
		}
	}

	result.WinAmount = totalWin - totalWager

	// Record bet
	s.recordBet(userID, "sic_bo", totalWager, totalWin, totalWin-totalWager, seeds.ServerSeed, seeds.ClientSeed, seeds.Nonce)

	return result, nil
}

// ============ Helper Methods ============

// recordBet records a bet in the database
func (s *TableGamesService) recordBet(userID uuid.UUID, gameType string, betAmount float64, winAmount float64, profit float64, serverSeed string, clientSeed string, nonce int) {
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

// GenerateSeeds creates new provably fair seeds
func (s *TableGamesService) GenerateSeeds(clientSeed string) (*Seeds, error) {
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

// HashSeed creates SHA-256 hash of a seed
func (s *TableGamesService) HashSeed(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])
}

// GenerateOutcome creates a deterministic outcome from seeds
func (s *TableGamesService) GenerateOutcome(serverSeed, clientSeed string, nonce int, max int) int {
	combined := fmt.Sprintf("%s:%s:%d", serverSeed, clientSeed, nonce)
	hash := sha256.Sum256([]byte(combined))
	result := new(big.Int).SetBytes(hash[:])
	return int(result.Mod(result, big.NewInt(int64(max))).Int64())
}
