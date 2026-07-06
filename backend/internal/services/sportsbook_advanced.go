package services

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

// Cashout types
type CashoutType string

const (
	CashoutTypeManual   CashoutType = "manual"
	CashoutTypeAuto     CashoutType = "auto"
	CashoutTypePartial  CashoutType = "partial"
)

// Cashout status
type CashoutStatus string

const (
	CashoutStatusPending  CashoutStatus = "pending"
	CashoutStatusApproved CashoutStatus = "approved"
	CashoutStatusRejected CashoutStatus = "rejected"
)

// Cashout request
type CashoutRequest struct {
	ID              string        `json:"id"`
	BetID           string        `json:"bet_id"`
	UserID          string        `json:"user_id"`
	OriginalStake   decimal.Decimal `json:"original_stake"`
	OriginalOdds    decimal.Decimal `json:"original_odds"`
	CurrentOdds     decimal.Decimal `json:"current_odds"`
	CashoutAmount   decimal.Decimal `json:"cashout_amount"`
	Type            CashoutType   `json:"type"`
	Status          CashoutStatus `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	ProcessedAt     *time.Time   `json:"processed_at,omitempty"`
	Reason          string        `json:"reason,omitempty"`
}

// Sportsbook with advanced features
type SportsbookService struct {
	oddsService    *OddsService
	cashoutService *CashoutService
	analytics      *BettingAnalytics
}

// NewSportsbookService creates a new sportsbook service
func NewSportsbookService() *SportsbookService {
	return &SportsbookService{
		oddsService:    NewOddsService(),
		cashoutService: NewCashoutService(),
		analytics:      NewBettingAnalytics(),
	}
}

// PlaceBet places a sports bet
func (s *SportsbookService) PlaceBet(userID string, selections []BetSelection, stake decimal.Decimal) (*Bet, error) {
	// Validate selections
	if len(selections) == 0 {
		return nil, fmt.Errorf("no selections")
	}

	if stake.LessThan(decimal.NewFromInt(1)) {
		return nil, fmt.Errorf("minimum stake is 1")
	}

	// Calculate total odds
	totalOdds := decimal.NewFromInt(1)
	for _, sel := range selections {
		odds, err := s.oddsService.GetOdds(sel.EventID, sel.MarketType, sel.Selection)
		if err != nil {
			return nil, fmt.Errorf("invalid selection: %v", err)
		}
		totalOdds = totalOdds.Mul(odds)
	}

	// Validate maximum odds
	maxOdds := decimal.NewFromInt(1000)
	if totalOdds.GreaterThan(maxOdds) {
		return nil, fmt.Errorf("maximum combined odds is 1000")
	}

	// Create bet
	bet := &Bet{
		ID:            generateUUID(),
		UserID:        userID,
		Type:          BetTypeSports,
		Stake:         stake,
		Odds:          totalOdds,
		PotentialWin:  stake.Mul(totalOdds),
		Status:        BetStatusPending,
		Selections:    selections,
		CreatedAt:     time.Now(),
	}

	// Initialize cashout
	s.cashoutService.InitCashout(bet)

	return bet, nil
}

// CalculateCashout calculates the cashout value
func (s *SportsbookService) CalculateCashout(bet *Bet) (decimal.Decimal, error) {
	return s.cashoutService.Calculate(bet)
}

// RequestCashout requests a cashout
func (s *SportsbookService) RequestCashout(betID, userID string, cashoutType CashoutType) (*CashoutRequest, error) {
	return s.cashoutService.Request(betID, userID, cashoutType)
}

// ProcessCashout processes a cashout request
func (s *SportsbookService) ProcessCashout(requestID string, approved bool, reason string) error {
	return s.cashoutService.Process(requestID, approved, reason)
}

// GetCashoutHistory gets cashout history for a user
func (s *SportsbookService) GetCashoutHistory(userID string) ([]CashoutRequest, error) {
	return s.cashoutService.GetHistory(userID)
}

// Odds Service
type OddsService struct {
	oddsCache map[string]map[string]decimal.Decimal
}

func NewOddsService() *OddsService {
	return &OddsService{
		oddsCache: make(map[string]map[string]decimal.Decimal),
	}
}

func (s *OddsService) GetOdds(eventID, marketType, selection string) (decimal.Decimal, error) {
	key := fmt.Sprintf("%s:%s", eventID, marketType)
	
	if odds, ok := s.oddsCache[key][selection]; ok {
		return odds, nil
	}

	// Default odds (would come from odds feed in production)
	return decimal.NewFromFloat(1.95), nil
}

func (s *OddsService) UpdateOdds(eventID, marketType string, odds map[string]decimal.Decimal) {
	key := fmt.Sprintf("%s:%s", eventID, marketType)
	s.oddsCache[key] = odds
}

// Cashout Service
type CashoutService struct {
	pendingRequests map[string]*CashoutRequest
	userHistory    map[string][]CashoutRequest
}

func NewCashoutService() *CashoutService {
	return &CashoutService{
		pendingRequests: make(map[string]*CashoutRequest),
		userHistory:    make(map[string][]CashoutRequest),
	}
}

func (s *CashoutService) InitCashout(bet *Bet) {
	// Initialize cashout availability
	bet.CashoutAvailable = true
	bet.CashoutAmount = bet.PotentialWin
}

func (s *CashoutService) Calculate(bet *Bet) (decimal.Decimal, error) {
	if !bet.CashoutAvailable {
		return decimal.Zero, fmt.Errorf("cashout not available")
	}

	// Cashout formula: (current_odds / original_odds) * stake
	ratio := bet.CurrentOdds.Div(bet.Odds)
	cashout := ratio.Mul(bet.Stake)

	// Apply cashout margin (e.g., 5%)
	margin := decimal.NewFromFloat(0.95)
	cashout = cashout.Mul(margin)

	// Ensure minimum cashout
	minCashout := bet.Stake.Mul(decimal.NewFromFloat(0.5))
	if cashout.LessThan(minCashout) {
		cashout = minCashout
	}

	return cashout, nil
}

func (s *CashoutService) Request(betID, userID string, cashoutType CashoutType) (*CashoutRequest, error) {
	// Calculate cashout amount (would fetch bet from DB)
	// For now, return a mock request
	request := &CashoutRequest{
		ID:            generateUUID(),
		BetID:         betID,
		UserID:        userID,
		OriginalStake: decimal.NewFromInt(100),
		OriginalOdds:  decimal.NewFromFloat(2.5),
		CurrentOdds:   decimal.NewFromFloat(1.8),
		CashoutAmount: decimal.NewFromFloat(140),
		Type:          cashoutType,
		Status:        CashoutStatusPending,
		CreatedAt:     time.Now(),
	}

	s.pendingRequests[request.ID] = request
	return request, nil
}

func (s *CashoutService) Process(requestID string, approved bool, reason string) error {
	request, ok := s.pendingRequests[requestID]
	if !ok {
		return fmt.Errorf("request not found")
	}

	if approved {
		request.Status = CashoutStatusApproved
		now := time.Now()
		request.ProcessedAt = &now
	} else {
		request.Status = CashoutStatusRejected
		request.Reason = reason
		now := time.Now()
		request.ProcessedAt = &now
	}

	// Add to history
	s.userHistory[request.UserID] = append(s.userHistory[request.UserID], *request)
	delete(s.pendingRequests, requestID)

	return nil
}

func (s *CashoutService) GetHistory(userID string) ([]CashoutRequest, error) {
	return s.userHistory[userID], nil
}

// Betting Analytics
type BettingAnalytics struct {
	betHistory map[string][]Bet
}

func NewBettingAnalytics() *BettingAnalytics {
	return &BettingAnalytics{
		betHistory: make(map[string][]Bet),
	}
}

func (a *BettingAnalytics) RecordBet(bet Bet) {
	a.betHistory[bet.UserID] = append(a.betHistory[bet.UserID], bet)
}

func (a *BettingAnalytics) GetUserStats(userID string) UserBettingStats {
	bets := a.betHistory[userID]
	
	totalBets := len(bets)
	totalStaked := decimal.Zero
	totalWon := decimal.Zero
	wins := 0

	for _, bet := range bets {
		totalStaked = totalStaked.Add(bet.Stake)
		if bet.Status == BetStatusWon {
			wins++
			totalWon = totalWon.Add(bet.WinAmount)
		}
	}

	winRate := float64(wins) / float64(totalBets)
	roi := decimal.Zero
	if totalStaked.GreaterThan(decimal.Zero) {
		roi = totalWon.Sub(totalStaked).Div(totalStaked).Mul(decimal.NewFromInt(100))
	}

	return UserBettingStats{
		UserID:       userID,
		TotalBets:    totalBets,
		TotalStaked:  totalStaked,
		TotalWon:     totalWon,
		WinRate:      winRate,
		ROI:          roi,
		AverageStake: totalStaked.Div(decimal.NewFromInt(totalBets)),
	}
}

// UserBettingStats represents user betting statistics
type UserBettingStats struct {
	UserID       string          `json:"user_id"`
	TotalBets    int             `json:"total_bets"`
	TotalStaked  decimal.Decimal `json:"total_staked"`
	TotalWon     decimal.Decimal `json:"total_won"`
	WinRate      float64         `json:"win_rate"`
	ROI          decimal.Decimal `json:"roi"`
	AverageStake decimal.Decimal `json:"average_stake"`
}

// Betting data structures
type BetType string

const (
	BetTypeSingle   BetType = "single"
	BetTypeMultiple BetType = "multiple"
	BetTypeSystem   BetType = "system"
	BetTypeSports   BetType = "sports"
)

type BetStatus string

const (
	BetStatusPending  BetStatus = "pending"
	BetStatusWon      BetStatus = "won"
	BetStatusLost    BetStatus = "lost"
	BetStatusVoid    BetStatus = "void"
	BetStatusCashedOut BetStatus = "cashed_out"
)

type Bet struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	Type            BetType           `json:"type"`
	Stake           decimal.Decimal   `json:"stake"`
	Odds            decimal.Decimal   `json:"odds"`
	PotentialWin    decimal.Decimal   `json:"potential_win"`
	WinAmount       decimal.Decimal   `json:"win_amount"`
	Status          BetStatus         `json:"status"`
	Selections      []BetSelection    `json:"selections"`
	CashoutAvailable bool             `json:"cashout_available"`
	CashoutAmount   decimal.Decimal   `json:"cashout_amount"`
	CurrentOdds     decimal.Decimal   `json:"current_odds"`
	CreatedAt       time.Time         `json:"created_at"`
	SettledAt       *time.Time       `json:"settled_at,omitempty"`
}

type BetSelection struct {
	EventID     string          `json:"event_id"`
	MarketType  string          `json:"market_type"`
	Selection   string          `json:"selection"`
	Odds        decimal.Decimal `json:"odds"`
}

// Parlay/Accumulator with advanced options
type ParlayBet struct {
	Bet           *Bet
	AutoCashout   bool
	AutoCashoutOdds decimal.Decimal
	PartialCashout bool
	PartialAmount  decimal.Decimal
}

func (s *SportsbookService) CreateParlay(selections []BetSelection, stake decimal.Decimal, options *ParlayOptions) (*ParlayBet, error) {
	// Create base bet
	bet, err := s.PlaceBet("user", selections, stake)
	if err != nil {
		return nil, err
	}

	parlay := &ParlayBet{
		Bet:          bet,
		AutoCashout:  false,
		PartialCashout: false,
	}

	if options != nil {
		parlay.AutoCashout = options.AutoCashout
		parlay.AutoCashoutOdds = options.AutoCashoutOdds
		parlay.PartialCashout = options.PartialCashout
		parlay.PartialAmount = options.PartialAmount
	}

	return parlay, nil
}

type ParlayOptions struct {
	AutoCashout       bool            `json:"auto_cashout"`
	AutoCashoutOdds   decimal.Decimal `json:"auto_cashout_odds"`
	PartialCashout    bool            `json:"partial_cashout"`
	PartialAmount     decimal.Decimal `json:"partial_amount"`
}

// Helper function to generate UUID (would use proper UUID library in production)
func generateUUID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), int64(math.rand.Float64()*1000000))
}
