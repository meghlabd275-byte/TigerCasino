package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/google/uuid"
)

// ============ Caribbean Stud Poker ============

type CaribbeanStudService struct{}

func NewCaribbeanStudService() *CaribbeanStudService {
	return &CaribbeanStudService{}
}

type CaribbeanStudHand struct {
	PlayerCards []Card    `json:"player_cards"`
	DealerCards []Card    `json:"dealer_cards"`
	PlayerRank  string    `json:"player_rank"`
	DealerRank  string    `json:"dealer_rank"`
	PlayerScore int       `json:"player_score"`
	DealerScore int      `json:"dealer_score"`
	Won         bool     `json:"won"`
	Payout      float64  `json:"payout"`
}

var caribbeanPayouts = map[string]float64{
	"royalFlush":     100.0,
	"straightFlush":  50.0,
	"fourOfAKind":    20.0,
	"fullHouse":      7.0,
	"flush":          5.0,
	"straight":       4.0,
	"threeOfAKind":   3.0,
	"twoPair":        2.0,
	"pair":           1.0,
	"highCard":       0.0,
}

func (s *CaribbeanStudService) Play(userID uuid.UUID, betAmount float64, clientSeed string, fold bool) (*CaribbeanStudHand, error) {
	deck := createDeck(1)
	deck.Shuffle()

	// Deal 5 cards to player and dealer
	playerCards := make([]Card, 5)
	dealerCards := make([]Card, 5)

	for i := 0; i < 5; i++ {
		card, ok := deck.Draw()
		if !ok {
			return nil, fmt.Errorf("deck exhausted")
		}
		playerCards[i] = card
	}

	for i := 0; i < 5; i++ {
		card, ok := deck.Draw()
		if !ok {
			return nil, fmt.Errorf("deck exhausted")
		}
		dealerCards[i] = card
	}

	result := &CaribbeanStudHand{
		PlayerCards: playerCards,
		DealerCards: dealerCards,
	}

	// Evaluate hands
	playerRank, playerScore := evaluateCaribbeanPokerHand(playerCards)
	dealerRank, dealerScore := evaluateCaribbeanPokerHand(dealerCards)

	result.PlayerRank = playerRank
	result.DealerRank = dealerRank
	result.PlayerScore = playerScore
	result.DealerScore = dealerScore

	if fold {
		// Fold - lose the bet
		result.Won = false
		result.Payout = 0
		return result, nil
	}

	// Check if dealer qualifies (Ace/King or better)
	dealerQualifies := dealerScore >= evaluateHand(map[string]int{"A": 14, "K": 13})

	if !dealerQualifies {
		// Dealer doesn't qualify - ante wins 1:1
		result.Won = true
		result.Payout = betAmount
		return result, nil
	}

	// Compare hands
	if playerScore > dealerScore {
		result.Won = true
		payout, ok := caribbeanPayouts[playerRank]
		if !ok {
			payout = 1.0
		}
		result.Payout = betAmount * (1 + payout)
	} else if playerScore == dealerScore {
		// Tie - push
		result.Won = false
		result.Payout = 0
	} else {
		result.Won = false
		result.Payout = 0
	}

	return result, nil
}

func evaluateCaribbeanPokerHand(cards []Card) (string, int) {
	// Sort by value descending
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value > cards[j].Value
	})

	// Check for hands
	suits := make(map[string]int)
	values := make(map[int]int)
	for _, c := range cards {
		suits[c.Suit]++
		values[c.Value]++
	}

	isFlush := false
	for _, count := range suits {
		if count >= 5 {
			isFlush = true
			break
		}
	}

	isStraight := true
	for i := 0; i < 4; i++ {
		if cards[i].Value-cards[i+1].Value != 1 {
			isStraight = false
			break
		}
	}
	// Check for Ace low straight
	if !isStraight && cards[0].Value == 14 && cards[1].Value == 5 && cards[2].Value == 4 && cards[3].Value == 3 && cards[4].Value == 2 {
		isStraight = true
	}

	// Count duplicates
	trips := 0
	pairs := 0
	fourKind := 0
	for _, count := range values {
		if count == 4 {
			fourKind++
		} else if count == 3 {
			trips++
		} else if count == 2 {
			pairs++
		}
	}

	// Determine hand rank
	if isFlush && isStraight {
		if cards[0].Value == 14 && cards[4].Value == 10 {
			return "royalFlush", 1000
		}
		return "straightFlush", 900
	}
	if fourKind == 1 {
		return "fourOfAKind", 800
	}
	if trips == 1 && pairs == 1 {
		return "fullHouse", 700
	}
	if isFlush {
		return "flush", 600
	}
	if isStraight {
		return "straight", 500
	}
	if trips == 1 {
		return "threeOfAKind", 400
	}
	if pairs == 2 {
		return "twoPair", 300
	}
	if pairs == 1 {
		return "pair", 200
	}

	// High card - use top 5 cards
	score := 0
	for i := 0; i < 5; i++ {
		score = score*15 + cards[i].Value
	}
	return "highCard", score
}

// ============ Three Card Poker ============

type ThreeCardPokerService struct{}

func NewThreeCardPokerService() *ThreeCardPokerService {
	return &ThreeCardPokerService{}
}

type ThreeCardPokerResult struct {
	PlayerCards []Card `json:"player_cards"`
	DealerCard  []Card `json:"dealer_cards"`
	PlayerRank  string `json:"player_rank"`
	Payout      float64 `json:"payout"`
	Won         bool   `json:"won"`
}

var threeCardPayouts = map[string]float64{
	"straightFlush": 40.0,
	"threeOfAKind":  30.0,
	"straight":      6.0,
	"flush":         4.0,
	"pair":          1.0,
	"highCard":      0.0,
}

func (s *ThreeCardPokerService) Play(userID uuid.UUID, betAmount float64, clientSeed string) (*ThreeCardPokerResult, error) {
	deck := createDeck(1)
	deck.Shuffle()

	playerCards := make([]Card, 3)
	dealerCards := make([]Card, 3)

	for i := 0; i < 3; i++ {
		card, ok := deck.Draw()
		if !ok {
			return nil, fmt.Errorf("deck exhausted")
		}
		playerCards[i] = card
	}

	for i := 0; i < 3; i++ {
		card, ok := deck.Draw()
		if !ok {
			return nil, fmt.Errorf("deck exhausted")
		}
		dealerCards[i] = card
	}

	result := &ThreeCardPokerResult{
		PlayerCards: playerCards,
		DealerCard:  dealerCards,
	}

	// Sort player cards by value for evaluation
	sort.Slice(playerCards, func(i, j int) bool {
		return playerCards[i].Value > playerCards[i].Value
	})
	sort.Slice(dealerCards, func(i, j int) bool {
		return dealerCards[i].Value > dealerCards[i].Value
	})

	playerRank := evaluateThreeCardHand(playerCards)
	dealerRank := evaluateThreeCardHand(dealerCards)

	result.PlayerRank = playerRank

	// Check if dealer qualifies (Queen high or better)
	dealerQualifies := dealerCards[0].Value >= 12 // Queen or higher

	if !dealerQualifies {
		// Dealer doesn't qualify - ante pays 1:1
		result.Won = true
		result.Payout = betAmount * 2
		return result, nil
	}

	// Compare hands
	playerScore := handScoreThreeCard(playerRank, playerCards)
	dealerScore := handScoreThreeCard(dealerRank, dealerCards)

	if playerScore > dealerScore {
		result.Won = true
		if playerRank == "pair" {
			result.Payout = betAmount * 2 // 1:1 for ante + pair plus bet
		} else {
			payout, _ := threeCardPayouts[playerRank]
			result.Payout = betAmount * (1 + payout)
		}
	} else if playerScore == dealerScore {
		// Tie - push
		result.Won = false
		result.Payout = 0
	} else {
		result.Won = false
		result.Payout = 0
	}

	return result, nil
}

func evaluateThreeCardHand(cards []Card) string {
	if len(cards) != 3 {
		return "highCard"
	}

	// Check for pairs
	values := make(map[int]int)
	for _, c := range cards {
		values[c.Value]++
	}

	for _, count := range values {
		if count == 3 {
			return "threeOfAKind"
		}
		if count == 2 {
			return "pair"
		}
	}

	// Check for flush
	suits := make(map[string]int)
	for _, c := range cards {
		suits[c.Suit]++
	}
	isFlush := false
	for _, count := range suits {
		if count == 3 {
			isFlush = true
			break
		}
	}

	// Check for straight
	isStraight := true
	sortedVals := []int{cards[0].Value, cards[1].Value, cards[2].Value}
	sort.Ints(sortedVals)
	for i := 0; i < 2; i++ {
		if sortedVals[i+1]-sortedVals[i] != 1 {
			isStraight = false
			break
		}
	}
	// Check for Ace low straight (A-2-3)
	if !isStraight && sortedVals[2] == 2 && sortedVals[1] == 3 && sortedVals[0] == 14 {
		isStraight = true
	}

	if isStraight && isFlush {
		return "straightFlush"
	}
	if isStraight {
		return "straight"
	}
	if isFlush {
		return "flush"
	}

	return "highCard"
}

func handScoreThreeCard(rank string, cards []Card) int {
	rankScores := map[string]int{
		"straightFlush": 6,
		"threeOfAKind":  5,
		"straight":      4,
		"flush":         3,
		"pair":          2,
		"highCard":      1,
	}

	score := rankScores[rank] * 1000000
	for i, c := range cards {
		score += c.Value * powInt(100, 2-i)
	}
	return score
}

// ============ Casino Hold'em ============

type CasinoHoldemService struct{}

func NewCasinoHoldemService() *CasinoHoldemService {
	return &CasinoHoldemService{}
}

type CasinoHoldemResult struct {
	PlayerCards    []Card `json:"player_cards"`
	CommunityCards []Card `json:"community_cards"`
	DealerCards    []Card `json:"dealer_cards"`
	PlayerBestHand string `json:"player_best_hand"`
	DealerBestHand string `json:"dealer_best_hand"`
	Payout         float64 `json:"payout"`
	Won            bool   `json:"won"`
}

var casinoHoldemPayouts = map[string]float64{
	"royalFlush":     100.0,
	"straightFlush":  50.0,
	"fourOfAKind":    20.0,
	"fullHouse":      7.0,
	"flush":          5.0,
	"straight":       4.0,
	"threeOfAKind":   3.0,
	"twoPair":        2.0,
	"pair":           1.0,
	"highCard":       0.0,
}

func (s *CasinoHoldemService) Play(userID uuid.UUID, betAmount float64, clientSeed string, fold bool) (*CasinoHoldemResult, error) {
	deck := createDeck(1)
	deck.Shuffle()

	// Deal player cards (2)
	playerCards := make([]Card, 2)
	for i := 0; i < 2; i++ {
		card, _ := deck.Draw()
		playerCards[i] = card
	}

	// Deal dealer cards (2)
	dealerCards := make([]Card, 2)
	for i := 0; i < 2; i++ {
		card, _ := deck.Draw()
		dealerCards[i] = card
	}

	// Deal community cards (5)
	communityCards := make([]Card, 5)
	for i := 0; i < 5; i++ {
		card, _ := deck.Draw()
		communityCards[i] = card
	}

	result := &CasinoHoldemResult{
		PlayerCards:    playerCards,
		CommunityCards: communityCards,
		DealerCards:    dealerCards,
	}

	// Make best 5-card hands
	playerBest, playerRank := makeBestPokerHand(append(playerCards, communityCards...))
	dealerBest, dealerRank := makeBestPokerHand(append(dealerCards, communityCards...))

	result.PlayerBestHand = playerRank
	result.DealerBestHand = dealerRank

	if fold {
		result.Won = false
		result.Payout = 0
		return result, nil
	}

	// Evaluate
	playerScore := pokerHandScore(playerRank, playerBest)
	dealerScore := pokerHandScore(dealerRank, dealerBest)

	// Dealer must qualify with pair of 4s or better
	dealerQualifies := dealerRankScores[dealerRank] >= dealerRankScores["fourOfAKind"]

	if !dealerQualifies {
		result.Won = true
		result.Payout = betAmount * 2
		return result, nil
	}

	if playerScore > dealerScore {
		result.Won = true
		payout, _ := casinoHoldemPayouts[playerRank]
		result.Payout = betAmount * (1 + payout)
	} else {
		result.Won = false
		result.Payout = 0
	}

	return result, nil
}

var dealerRankScores = map[string]int{
	"royalFlush":     10,
	"straightFlush":  9,
	"fourOfAKind":    8,
	"fullHouse":      7,
	"flush":          6,
	"straight":       5,
	"threeOfAKind":   4,
	"twoPair":        3,
	"pair":           2,
	"highCard":       1,
}

// ============ Ultimate Texas Hold'em ============

type UltimateTexasHoldemService struct{}

func NewUltimateTexasHoldemService() *UltimateTexasHoldemService {
	return &UltimateTexasHoldemService{}
}

type UltimateTexasHoldemResult struct {
	PlayerCards      []Card `json:"player_cards"`
	CommunityCards   []Card `json:"community_cards"`
	DealerCards      []Card `json:"dealer_cards"`
	PlayerBestHand  string `json:"player_best_hand"`
	DealerBestHand  string `json:"dealer_best_hand"`
	AnteBet         float64 `json:"ante_bet"`
	BlindBet        float64 `json:"blind_bet"`
	PlayBet         float64 `json:"play_bet"`
	TripsBet        float64 `json:"trips_bet"`
	Payout          float64 `json:"payout"`
	Won             bool   `json:"won"`
}

func (s *UltimateTexasHoldemService) Play(userID uuid.UUID, anteAmount float64, clientSeed string, action string) (*UltimateTexasHoldemResult, error) {
	deck := createDeck(1)
	deck.Shuffle()

	// Deal
	playerCards := make([]Card, 2)
	dealerCards := make([]Card, 2)
	communityCards := make([]Card, 5)

	for i := 0; i < 2; i++ {
		playerCards[i], _ = deck.Draw()
		dealerCards[i], _ = deck.Draw()
	}
	for i := 0; i < 5; i++ {
		communityCards[i], _ = deck.Draw()
	}

	result := &UltimateTexasHoldemResult{
		PlayerCards:    playerCards,
		CommunityCards: communityCards,
		DealerCards:    dealerCards,
		AnteBet:       anteAmount,
		BlindBet:      anteAmount,
	}

	// Evaluate hands
	_, playerRank := makeBestPokerHand(append(playerCards, communityCards...))
	_, dealerRank := makeBestPokerHand(append(dealerCards, communityCards...))

	result.PlayerBestHand = playerRank
	result.DealerBestHand = dealerRank

	// Check dealer qualification
	dealerQualifies := dealerRankScores[dealerRank] >= dealerRankScores["pair"]

	playerScore := pokerHandScore(playerRank, nil)
	dealerScore := pokerHandScore(dealerRank, nil)

	if dealerQualifies && playerScore > dealerScore {
		result.Won = true
		result.Payout = anteAmount * 2 // Ante 1:1
		// Blind bet pays according to paytable
		if blindPayout, ok := casinoHoldemPayouts[playerRank]; ok {
			result.Payout += anteAmount * blindPayout
		}
	} else if !dealerQualifies {
		result.Won = true
		result.Payout = anteAmount // Push ante
	} else {
		result.Won = false
		result.Payout = 0
	}

	return result, nil
}

// ============ Teen Patti (Indian Poker) ============

type TeenPattiService struct{}

func NewTeenPattiService() *TeenPattiService {
	return &TeenPattiService{}
}

type TeenPattiResult struct {
	PlayerCards []Card `json:"player_cards"`
	DealerCards []Card `json:"dealer_cards"`
	PlayerRank  string `json:"player_rank"`
	DealerRank  string `json:"dealer_rank"`
	Payout      float64 `json:"payout"`
	Won         bool   `json:"won"`
}

var teenPattiRankOrder = map[string]int{
	"trail":       13, // Three of a kind
	"pureSequence": 12, // Straight flush
	"sequence":     11, // Straight
	"color":       10, // Flush
	"pair":         9, // One pair
	"highCard":     0, // High card
}

func (s *TeenPattiService) Play(userID uuid.UUID, betAmount float64, clientSeed string) (*TeenPattiResult, error) {
	deck := createDeck(1)
	deck.Shuffle()

	playerCards := make([]Card, 3)
	dealerCards := make([]Card, 3)

	for i := 0; i < 3; i++ {
		playerCards[i], _ = deck.Draw()
		dealerCards[i], _ = deck.Draw()
	}

	result := &TeenPattiResult{
		PlayerCards: playerCards,
		DealerCards: dealerCards,
	}

	playerRank := evaluateTeenPattiHand(playerCards)
	dealerRank := evaluateTeenPattiHand(dealerCards)

	result.PlayerRank = playerRank
	result.DealerRank = dealerRank

	playerScore := teenPattiScore(playerRank, playerCards)
	dealerScore := teenPattiScore(dealerRank, dealerCards)

	if playerScore > dealerScore {
		result.Won = true
		result.Payout = betAmount * 2
	} else if playerScore == dealerScore {
		// Compare high cards
		playerHigh := playerCards[0].Value
		dealerHigh := dealerCards[0].Value
		if playerHigh > dealerHigh {
			result.Won = true
			result.Payout = betAmount * 2
		} else if playerHigh == dealerHigh {
			result.Won = false
			result.Payout = 0 // Tie - split
		} else {
			result.Won = false
			result.Payout = 0
		}
	} else {
		result.Won = false
		result.Payout = 0
	}

	return result, nil
}

func evaluateTeenPattiHand(cards []Card) string {
	// Sort by value
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value > cards[j].Value
	})

	// Check for trail (three of a kind)
	values := make(map[int]int)
	for _, c := range cards {
		values[c.Value]++
	}
	for _, count := range values {
		if count == 3 {
			return "trail"
		}
	}

	// Check for pure sequence (straight flush)
	isStraight := true
	for i := 0; i < 2; i++ {
		if cards[i].Value-cards[i+1].Value != 1 {
			isStraight = false
			break
		}
	}
	// Also check A-2-3
	if !isStraight && cards[0].Value == 14 && cards[1].Value == 3 && cards[2].Value == 2 {
		isStraight = true
	}

	isFlush := true
	firstSuit := cards[0].Suit
	for i := 1; i < 3; i++ {
		if cards[i].Suit != firstSuit {
			isFlush = false
			break
		}
	}

	if isStraight && isFlush {
		return "pureSequence"
	}

	if isStraight {
		return "sequence"
	}

	if isFlush {
		return "color"
	}

	// Check for pair
	for _, count := range values {
		if count == 2 {
			return "pair"
		}
	}

	return "highCard"
}

func teenPattiScore(rank string, cards []Card) int {
	baseScore := teenPattiRankOrder[rank] * 1000000
	for i, c := range cards {
		baseScore += c.Value * powInt(100, 2-i)
	}
	return baseScore
}

// ============ Andar Bahar ============

type AndarBaharService struct{}

func NewAndarBaharService() *AndarBaharService {
	return &AndarBaharService{}
}

type AndarBaharResult struct {
	JokerCard     Card    `json:"joker_card"`
	AndarCards    []Card  `json:"andar_cards"`
	BaharCards    []Card  `json:"bahar_cards"`
	AndarCount    int     `json:"andar_count"`
	BaharCount    int     `json:"bahar_count"`
	Winner        string  `json:"winner"`
	BetOn         string  `json:"bet_on"`
	Payout        float64 `json:"payout"`
}

func (s *AndarBaharService) Play(userID uuid.UUID, betAmount float64, betOn string, clientSeed string) (*AndarBaharResult, error) {
	deck := createDeck(1)
	deck.Shuffle()

	// Draw Joker card
	jokerCard, _ := deck.Draw()

	result := &AndarBaharResult{
		JokerCard: jokerCard,
		BetOn:     betOn,
	}

	// Alternate dealing: Andar, Bahar, Andar, Bahar...
	andarCards := []Card{}
	baharCards := []Card{}

	// Deal cards alternating
	for i := 0; i < 26; i++ {
		card, _ := deck.Draw()
		if i%2 == 0 {
			andarCards = append(andarCards, card)
			// Check if this card matches joker
			if card.Value == jokerCard.Value {
				result.AndarCount = len(andarCards)
				result.BaharCount = 0
				break
			}
		} else {
			baharCards = append(baharCards, card)
			// Check if this card matches joker
			if card.Value == jokerCard.Value {
				result.BaharCount = len(baharCards)
				result.AndarCount = 0
				break
			}
		}
	}

	result.AndarCards =andarCards
	result.BaharCards = baharCards

	// Determine winner
	if result.AndarCount > 0 && (result.BaharCount == 0 || result.AndarCount < result.BaharCount) {
		result.Winner = "andar"
	} else {
		result.Winner = "bahar"
	}

	// Calculate payout
	if result.Winner == betOn {
		// Payout based on position
		var position int
		if result.Winner == "andar" {
			position = result.AndarCount
		} else {
			position = result.BaharCount
		}

		// Payout table: closer to 1 = higher payout
		payout := 0.85
		if position <= 5 {
			payout = 1.5
		} else if position <= 10 {
			payout = 1.2
		} else if position <= 15 {
			payout = 1.0
		} else if position <= 20 {
			payout = 0.9
		}

		result.Payout = betAmount * (1 + payout)
	} else {
		result.Payout = 0
	}

	return result, nil
}

// ============ Sic Bo ============

type SicBoService struct{}

func NewSicBoService() *SicBoService {
	return &SicBoService{}
}

type SicBoResult struct {
	Dice1       int     `json:"dice_1"`
	Dice2       int     `json:"dice_2"`
	Dice3       int     `json:"dice_3"`
	Total       int     `json:"total"`
	Bets        []SicBoBet `json:"bets"`
	TotalWin    float64 `json:"total_win"`
}

type SicBoBet struct {
	BetType  string  `json:"bet_type"`
	Amount   float64 `json:"amount"`
	Selection string `json:"selection"`
	Paid     float64 `json:"paid"`
	Won      bool    `json:"won"`
}

var sicBoPayouts = map[string]float64{
	"small":      1.0,   // 4-10 except triples
	"big":        1.0,   // 11-17 except triples
	"odd":        1.0,   // Odd
	"even":       1.0,   // Even
	"triple":     30.0,  // Any triple
	"anyTriple":  30.0,  // Any three of a kind
	"double":     10.0,  // Any double
	"threeDice":  180.0, // Specific triple
	"threeTotal4": 60.0,
	"threeTotal5": 30.0,
	"threeTotal6": 18.0,
	"threeTotal7": 12.0,
	"threeTotal8": 8.0,
	"threeTotal9": 7.0,
	"threeTotal10":6.0,
	"threeTotal11":6.0,
	"threeTotal12":7.0,
	"threeTotal13":8.0,
	"threeTotal14":12.0,
	"threeTotal15":18.0,
	"threeTotal16":30.0,
	"threeTotal17":60.0,
	"domino":     5.0,   // Two specific numbers
	"single":     1.0,   // Single number (1:1 for one die, 2:1 for two, 3:1 for three)
}

func (s *SicBoService) Roll(clientSeed string, bets []SicBoBet) (*SicBoResult, error) {
	// Generate random dice
	dice1, _ := rand.Int(rand.Reader, big.NewInt(6))
	dice2, _ := rand.Int(rand.Reader, big.NewInt(6))
	dice3, _ := rand.Int(rand.Reader, big.NewInt(6))

	result := &SicBoResult{
		Dice1: int(dice1.Int64()) + 1,
		Dice2: int(dice2.Int64()) + 1,
		Dice3: int(dice3.Int64()) + 1,
		Total: int(dice1.Int64()+dice2.Int64()+dice3.Int64()) + 3,
		Bets:  bets,
	}

	totalWin := 0.0

	for i := range bets {
		bet := &bets[i]
		bet.Won = false
		bet.Paid = 0

		switch bet.BetType {
		case "small":
			// 4-10, not triple
			if result.Total >= 4 && result.Total <= 10 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + sicBoPayouts["small"])
			}
		case "big":
			// 11-17, not triple
			if result.Total >= 11 && result.Total <= 17 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + sicBoPayouts["big"])
			}
		case "odd":
			if result.Total%2 == 1 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + sicBoPayouts["odd"])
			}
		case "even":
			if result.Total%2 == 0 && !(result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3) {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + sicBoPayouts["even"])
			}
		case "triple":
			if bet.Selection != "any" {
				expected, _ := strconv.Atoi(bet.Selection)
				if result.Dice1 == expected && result.Dice2 == expected && result.Dice3 == expected {
					bet.Won = true
					bet.Paid = bet.Amount * (1 + sicBoPayouts["threeDice"])
				}
			}
		case "anyTriple":
			if result.Dice1 == result.Dice2 && result.Dice2 == result.Dice3 {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + sicBoPayouts["anyTriple"])
			}
		case "double":
			if result.Dice1 == result.Dice2 || result.Dice2 == result.Dice3 || result.Dice1 == result.Dice3 {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + sicBoPayouts["double"])
			}
		case "total":
			totalBet, _ := strconv.Atoi(bet.Selection)
			if result.Total == totalBet {
				payout := sicBoPayouts["threeTotal"+bet.Selection]
				if payout == 0 {
					payout = 6.0 // Default for 9-12
				}
				bet.Won = true
				bet.Paid = bet.Amount * (1 + payout)
			}
		case "single":
			 num, _ := strconv.Atoi(bet.Selection)
			 matches := 0
			 if result.Dice1 == num { matches++ }
			 if result.Dice2 == num { matches++ }
			 if result.Dice3 == num { matches++ }
			 if matches > 0 {
				 bet.Won = true
				 bet.Paid = bet.Amount * float64(matches) // 1:1 per matching die
			 }
		}

		totalWin += bet.Paid
	}

	result.TotalWin = totalWin
	result.Bets = bets
	return result, nil
}

// ============ Dragon Tiger ============

type DragonTigerService struct{}

func NewDragonTigerService() *DragonTigerService {
	return &DragonTigerService{}
}

type DragonTigerResult struct {
	DragonCard Card    `json:"dragon_card"`
	TigerCard  Card    `json:"tiger_card"`
	Winner     string  `json:"winner"`
	BetOn      string  `json:"bet_on"`
	Payout     float64 `json:"payout"`
}

func (s *DragonTigerService) Play(userID uuid.UUID, betAmount float64, betOn string, clientSeed string) (*DragonTigerResult, error) {
	deck := createDeck(1)
	deck.Shuffle()

	dragonCard, _ := deck.Draw()
	tigerCard, _ := deck.Draw()

	result := &DragonTigerResult{
		DragonCard: dragonCard,
		TigerCard:  tigerCard,
		BetOn:      betOn,
	}

	if dragonCard.Value > tigerCard.Value {
		result.Winner = "dragon"
	} else if tigerCard.Value > dragonCard.Value {
		result.Winner = "tiger"
	} else {
		result.Winner = "tie"
	}

	// Calculate payout
	if result.Winner == betOn {
		if result.Winner == "tie" {
			result.Payout = betAmount * 8 // Tie pays 8:1
		} else {
			result.Payout = betAmount * 2
		}
	} else if result.Winner == "tie" && (betOn == "dragon" || betOn == "tiger") {
		// Dragon/Tiger bet loses on tie, but usually returns half
		result.Payout = betAmount * 0.5
	}

	return result, nil
}

// ============ Craps ============

type CrapsService struct{}

func NewCrapsService() *CrapsService {
	return &CrapsService{}
}

type CrapsResult struct {
	Dice1       int     `json:"dice_1"`
	Dice2       int     `json:"dice_2"`
	Total       int     `json:"total"`
	Phase       string  `json:"phase"` // comeOut, point, resolve
	Point       int     `json:"point"`
	Bets        []CrapsBet `json:"bets"`
	TotalWin    float64 `json:"total_win"`
}

type CrapsBet struct {
	BetType  string  `json:"bet_type"`
	Amount   float64 `json:"amount"`
	Paid     float64 `json:"paid"`
	Won      bool    `json:"won"`
}

var crapsPayouts = map[string]float64{
	"passLine":      1.0,
	"dontPass":      1.0,
	"come":          1.0,
	"dontCome":      1.0,
	"field":         1.0,
	"place6":        7.0/6.0,
	"place8":        7.0/6.0,
	"place5":        7.0/5.0,
	"place9":       7.0/5.0,
	"place10":       9.0/5.0,
	"place4":        9.0/5.0,
	"hard4":         7.0,
	"hard6":         9.0,
	"hard8":         9.0,
	"hard10":        7.0,
	"any7":          4.0,
	"anyCraps":      7.0,
	"craps2":        30.0,
	"craps3":       15.0,
	"craps12":      30.0,
	"horn2":        30.0,
	"horn3":       15.0,
	"horn12":      30.0,
	"hop2":         30.0,
	"hop3":        15.0,
	"hop12":       30.0,
	"hopAny":       14.0,
}

func (s *CrapsService) Roll(clientSeed string, bets []CrapsBet, phase string, point int) (*CrapsResult, error) {
	dice1, _ := rand.Int(rand.Reader, big.NewInt(6))
	dice2, _ := rand.Int(rand.Reader, big.NewInt(6))

	total := int(dice1.Int64()+dice2.Int64()) + 2

	result := &CrapsResult{
		Dice1:   int(dice1.Int64()) + 1,
		Dice2:   int(dice2.Int64()) + 1,
		Total:   total,
		Phase:   phase,
		Point:   point,
		Bets:    bets,
	}

	totalWin := 0.0

	// Determine outcome
	var won bool
	newPhase := phase

	if phase == "comeOut" {
		if total == 7 || total == 11 {
			// Natural - pass line wins
			won = true
			newPhase = "comeOut"
		} else if total == 2 || total == 3 || total == 12 {
			// Craps - pass line loses
			won = false
			newPhase = "comeOut"
		} else {
			// Point established
			won = false
			newPhase = "point"
			result.Point = total
		}
	} else {
		// Point phase
		if total == result.Point {
			// Point hit - pass line wins
			won = true
			newPhase = "comeOut"
		} else if total == 7 {
			// Seven out - pass line loses
			won = false
			newPhase = "comeOut"
		} else {
			// Continue point
			won = false
			newPhase = "point"
		}
	}

	// Evaluate bets
	for i := range bets {
		bet := &bets[i]
		bet.Won = false
		bet.Paid = 0

		switch bet.BetType {
		case "passLine":
			bet.Won = won
			if bet.Won {
				bet.Paid = bet.Amount * (1 + crapsPayouts["passLine"])
			}
		case "dontPass":
			bet.Won = !won
			if bet.Won {
				bet.Paid = bet.Amount * (1 + crapsPayouts["dontPass"])
			}
		case "field":
			if total == 2 || total == 3 || total == 4 || total == 9 || total == 10 || total == 11 || total == 12 {
				bet.Won = true
				if total == 2 || total == 12 {
					bet.Paid = bet.Amount * 3 // Double on 2 and 12
				} else {
					bet.Paid = bet.Amount * 2
				}
			}
		case "place6", "place8":
			if total == 6 || total == 8 {
				bet.Won = true
				bet.Paid = bet.Amount * crapsPayouts[bet.BetType]
			}
		case "any7":
			if total == 7 {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + crapsPayouts["any7"])
			}
		case "anyCraps":
			if total == 2 || total == 3 || total == 12 {
				bet.Won = true
				bet.Paid = bet.Amount * (1 + crapsPayouts["anyCraps"])
			}
		}

		totalWin += bet.Paid
	}

	result.TotalWin = totalWin
	result.Bets = bets
	return result, nil
}

// Helper function
func powInt(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
