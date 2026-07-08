package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// SportsbookService handles sports betting operations
type SportsbookService struct {
	db *gorm.DB
}

func NewSportsbookService(db *gorm.DB) *SportsbookService {
	return &SportsbookService{db: db}
}

// ============ SPORTS MODELS ============

type SportsEvent struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	ExternalID      string    `gorm:"uniqueIndex"`
	Sport           string    `gorm:"not null"` // football, basketball, tennis, etc.
	League          string    `gorm:"not null"`
	HomeTeam        string    `gorm:"not null"`
	AwayTeam        string    `gorm:"not null"`
	StartTime       time.Time `gorm:"not null"`
	Status          string    `gorm:"default:'upcoming'"` // upcoming, live, finished, cancelled
	HomeScore       int       `gorm:"default:0"`
	AwayScore       int       `gorm:"default:0"`
	HomeOdds        float64   `gorm:"not null"`
	DrawOdds        float64   `gorm:"default:0"`
	AwayOdds        float64   `gorm:"not null"`
	OverUnderLine   float64   `gorm:"default:2.5"`
	OverOdds        float64   `gorm:"default:1.95"`
	UnderOdds       float64   `gorm:"default:1.95"`
	Handicap        float64   `gorm:"default:0"`
	HomeHandicapOdds float64   `gorm:"default:1.95"`
	AwayHandicapOdds float64  `gorm:"default:1.95"`
	Featured        bool      `gorm:"default:false"`
	LiveEnabled     bool      `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type SportsBet struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	EventID       uuid.UUID `gorm:"type:uuid;not null"`
	Selection      string    `gorm:"not null"` // home_win, away_win, draw, over, under, handicap_home, handicap_away
	Odds          float64   `gorm:"not null"`
	Stake         float64   `gorm:"not null"`
	PotentialWin   float64   `gorm:"not null"`
	ActualWin     float64   `gorm:"default:0"`
	Status        string    `gorm:"default:'pending'"` // pending, won, lost, cancelled, refunded
	Result         string    // Actual outcome
	SettledAt     *time.Time
	CreatedAt      time.Time
}

type ParlayBet struct {
	ID            uuid.UUID        `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	Stake         float64          `gorm:"not null"`
	Odds          float64          `gorm:"not null"`
	PotentialWin  float64          `gorm:"not null"`
	ActualWin     float64          `gorm:"default:0"`
	Status        string           `gorm:"default:'pending'"`
	Selections    string           `gorm:"type:jsonb"` // JSON array of selections
	CreatedAt     time.Time
	SettledAt     *time.Time
}

// ============ ODDS ENGINE ============

// OddsEngine generates realistic odds for sports events
type OddsEngine struct{}

func NewOddsEngine() *OddsEngine {
	return &OddsEngine{}
}

// CalculateMatchOdds calculates 1X2 odds for a match
func (e *OddsEngine) CalculateMatchOdds(homeStrength, awayStrength float64) (home, draw, away float64) {
	// Home advantage factor
	homeAdvantage := 1.15

	// Calculate implied probabilities
	totalStrength := homeStrength + awayStrength + 0.1 // 0.1 represents draw tendency

	homeProb := (homeStrength * homeAdvantage) / totalStrength
	awayProb := awayStrength / totalStrength
	drawProb := 0.35 - (homeProb + awayProb)/3 // Normalize draw probability

	if drawProb < 0.15 {
		drawProb = 0.15
	}

	// Apply bookmaker margin (typically 5-10%)
	margin := 1.05

	homeOdds = 1.0 / (homeProb / margin)
	drawOdds = 1.0 / (drawProb / margin)
	awayOdds = 1.0 / (awayProb / margin)

	return
}

// CalculateOverUnderOdds calculates over/under odds
func (e *OddsEngine) CalculateOverUnderOdds(line float64, expectedGoals float64) (over, under float64) {
	// Line adjustment based on expected goals
	adjustedLine := expectedGoals / 2

	if adjustedLine > line {
		over = 1.85 + (adjustedLine-line)*0.1
		under = 2.00 - (adjustedLine-line)*0.1
	} else {
		over = 1.95 - (line-adjustedLine)*0.1
		under = 1.90 + (line-adjustedLine)*0.1
	}

	// Clamp values
	if over > 2.10 {
		over = 2.10
	}
	if under < 1.80 {
		under = 1.80
	}

	return
}

// CalculateHandicapOdds calculates Asian handicap odds
func (e *OddsEngine) CalculateHandicapOdds(handicap float64, homeStrength, awayStrength float64) (home, away float64) {
	homeProb := homeStrength / (homeStrength + awayStrength)
	awayProb := 1 - homeProb

	// Adjust for handicap
	if handicap > 0 {
		homeProb += handicap * 0.05
	} else if handicap < 0 {
		awayProb -= handicap * 0.05
	}

	margin := 1.08

	homeOdds = 1.0 / (homeProb / margin)
	awayOdds = 1.0 / (awayProb / margin)

	return
}

// ============ EVENT MANAGEMENT ============

func (s *SportsbookService) CreateEvent(event *SportsEvent) error {
	event.ID = uuid.New()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	return s.db.Create(event).Error
}

func (s *SportsbookService) GetEvent(id uuid.UUID) (*SportsEvent, error) {
	var event SportsEvent
	err := s.db.First(&event, id).Error
	return &event, err
}

func (s *SportsbookService) GetEventsBySport(sport string, status string) ([]SportsEvent, error) {
	var events []SportsEvent
	query := s.db.Where("sport = ?", sport)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("start_time ASC").Find(&events).Error
	return events, err
}

func (s *SportsbookService) GetLiveEvents() ([]SportsEvent, error) {
	var events []SportsEvent
	err := s.db.Where("status = ? AND live_enabled = ?", "live", true).
		Order("start_time ASC").
		Find(&events).Error
	return events, err
}

func (s *SportsbookService) GetUpcomingEvents(limit int) ([]SportsEvent, error) {
	if limit <= 0 {
		limit = 50
	}

	var events []SportsEvent
	err := s.db.Where("status = ? AND start_time > ?", "upcoming", time.Now()).
		Order("start_time ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (s *SportsbookService) GetFeaturedEvents() ([]SportsEvent, error) {
	var events []SportsEvent
	err := s.db.Where("featured = ? AND status IN ?", true, []string{"upcoming", "live"}).
		Order("start_time ASC").
		Find(&events).Error
	return events, err
}

// ============ BETTING ============

func (s *SportsbookService) PlaceBet(userID, eventID uuid.UUID, selection string, stake float64) (*SportsBet, error) {
	var event SportsEvent
	if err := s.db.First(&event, eventID).Error; err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Validate selection
	var odds float64
	switch selection {
	case "home_win":
		odds = event.HomeOdds
	case "draw":
		odds = event.DrawOdds
	case "away_win":
		odds = event.AwayOdds
	case "over":
		odds = event.OverOdds
	case "under":
		odds = event.UnderOdds
	case "handicap_home":
		odds = event.HomeHandicapOdds
	case "handicap_away":
		odds = event.AwayHandicapOdds
	default:
		return nil, fmt.Errorf("invalid selection")
	}

	if odds <= 0 {
		return nil, fmt.Errorf("odds not available for this selection")
	}

	// Check minimum bet
	if stake < 0.10 {
		return nil, fmt.Errorf("minimum bet is 0.10")
	}

	// Deduct stake from balance
	userService := NewUserService(s.db)
	var user models.User
	s.db.First(&user, userID)
	if user.Balance < stake {
		return nil, fmt.Errorf("insufficient balance")
	}
	userService.UpdateBalance(userID, -stake)

	// Create bet
	bet := SportsBet{
		ID:           uuid.New(),
		UserID:       userID,
		EventID:      eventID,
		Selection:    selection,
		Odds:         odds,
		Stake:        stake,
		PotentialWin: stake * odds,
		Status:       "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.db.Create(&bet).Error; err != nil {
		// Refund on error
		userService.UpdateBalance(userID, stake)
		return nil, err
	}

	return &bet, nil
}

func (s *SportsbookService) PlaceParlayBet(userID uuid.UUID, selections []struct {
	EventID   uuid.UUID
	Selection string
}, stake float64) (*ParlayBet, error) {
	if len(selections) < 2 {
		return nil, fmt.Errorf("parlay must have at least 2 selections")
	}

	if len(selections) > 10 {
		return nil, fmt.Errorf("parlay cannot have more than 10 selections")
	}

	// Calculate combined odds
	totalOdds := 1.0
	selectionsJSON, _ := json.Marshal(selections)

	for _, sel := range selections {
		var event SportsEvent
		s.db.First(&event, sel.EventID)

		var odds float64
		switch sel.Selection {
		case "home_win":
			odds = event.HomeOdds
		case "draw":
			odds = event.DrawOdds
		case "away_win":
			odds = event.AwayOdds
		default:
			odds = 1.0
		}

		totalOdds *= odds
	}

	// Check minimum bet
	if stake < 0.50 {
		return nil, fmt.Errorf("minimum parlay bet is 0.50")
	}

	// Deduct stake
	userService := NewUserService(s.db)
	var user models.User
	s.db.First(&user, userID)
	if user.Balance < stake {
		return nil, fmt.Errorf("insufficient balance")
	}
	userService.UpdateBalance(userID, -stake)

	// Create parlay bet
	bet := ParlayBet{
		ID:           uuid.New(),
		UserID:       userID,
		Stake:        stake,
		Odds:         totalOdds,
		PotentialWin: stake * totalOdds,
		Selections:   string(selectionsJSON),
		CreatedAt:   time.Now(),
	}

	if err := s.db.Create(&bet).Error; err != nil {
		userService.UpdateBalance(userID, stake)
		return nil, err
	}

	return &bet, nil
}

// ============ SETTLEMENT ============

func (s *SportsbookService) SettleEvent(eventID uuid.UUID, result string) error {
	var event SportsEvent
	if err := s.db.First(&event, eventID).Error; err != nil {
		return err
	}

	// Update event status
	updates := map[string]interface{}{
		"status": "finished",
	}
	if result == "home_win" {
		event.HomeScore++
	} else if result == "away_win" {
		event.AwayScore++
	}
	s.db.Save(&event)

	// Settle all bets
	var bets []SportsBet
	s.db.Where("event_id = ? AND status = ?", eventID, "pending").Find(&bets)

	userService := NewUserService(s.db)

	for _, bet := range bets {
		won := bet.Selection == result

		if won {
			bet.Status = "won"
			bet.ActualWin = bet.Stake * bet.Odds
			bet.Result = result
			userService.UpdateBalance(bet.UserID, bet.ActualWin)
		} else {
			bet.Status = "lost"
			bet.Result = result
		}

		now := time.Now()
		bet.SettledAt = &now
		s.db.Save(&bet)
	}

	return nil
}

func (s *SportsbookService) SettleParlayBet(betID uuid.UUID) error {
	var bet ParlayBet
	if err := s.db.First(&bet, betID).Error; err != nil {
		return err
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet already settled")
	}

	var selections []struct {
		EventID   uuid.UUID
		Selection string
	}
	json.Unmarshal([]byte(bet.Selections), &selections)

	allWon := true
	for _, sel := range selections {
		var event SportsEvent
		s.db.First(&event, sel.EventID)

		if event.Status != "finished" {
			allWon = false
			break
		}

		// Determine event result (simplified)
		eventResult := "draw"
		if event.HomeScore > event.AwayScore {
			eventResult = "home_win"
		} else if event.AwayScore > event.HomeScore {
			eventResult = "away_win"
		}

		if sel.Selection != eventResult {
			allWon = false
			break
		}
	}

	userService := NewUserService(s.db)
	now := time.Now()

	if allWon {
		bet.Status = "won"
		bet.ActualWin = bet.Stake * bet.Odds
		userService.UpdateBalance(bet.UserID, bet.ActualWin)
	} else {
		bet.Status = "lost"
	}

	bet.SettledAt = &now
	s.db.Save(&bet)

	return nil
}

// ============ LIVE BETTING ============

func (s *SportsbookService) UpdateLiveScore(eventID uuid.UUID, homeScore, awayScore int) error {
	var event SportsEvent
	if err := s.db.First(&event, eventID).Error; err != nil {
		return err
	}

	event.HomeScore = homeScore
	event.AwayScore = awayScore
	event.UpdatedAt = time.Now()

	// Adjust odds based on score (simplified)
	// In production, this would be much more sophisticated
	margin := 0.02 // 2% shift per goal
	goalDiff := homeScore - awayScore

	if goalDiff > 0 {
		event.HomeOdds = event.HomeOdds * (1 - margin*float64(goalDiff))
		event.AwayOdds = event.AwayOdds * (1 + margin*float64(goalDiff))
	} else if goalDiff < 0 {
		event.HomeOdds = event.HomeOdds * (1 + margin*float64(-goalDiff))
		event.AwayOdds = event.AwayOdds * (1 - margin*float64(-goalDiff))
	}

	return s.db.Save(&event).Error
}

// ============ BET HISTORY ============

func (s *SportsbookService) GetUserBets(userID uuid.UUID, status string, limit int) ([]SportsBet, error) {
	if limit <= 0 {
		limit = 20
	}

	var bets []SportsBet
	query := s.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&bets).Error
	return bets, err
}

func (s *SportsbookService) GetUserParlayBets(userID uuid.UUID, limit int) ([]ParlayBet, error) {
	if limit <= 0 {
		limit = 20
	}

	var bets []ParlayBet
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&bets).Error
	return bets, err
}

// ============ SPORTS STATISTICS ============

type SportsStats struct {
	TotalBets      int64   `json:"total_bets"`
	WonBets       int64   `json:"won_bets"`
	LostBets      int64   `json:"lost_bets"`
	PendingBets   int64   `json:"pending_bets"`
	TotalStaked   float64 `json:"total_staked"`
	TotalWon      float64 `json:"total_won"`
	NetProfit     float64 `json:"net_profit"`
	WinRate       float64 `json:"win_rate"`
}

func (s *SportsbookService) GetUserStats(userID uuid.UUID) (*SportsStats, error) {
	var stats SportsStats

	s.db.Model(&SportsBet{}).Where("user_id = ?", userID).Count(&stats.TotalBets)
	s.db.Model(&SportsBet{}).Where("user_id = ? AND status = ?", userID, "won").Count(&stats.WonBets)
	s.db.Model(&SportsBet{}).Where("user_id = ? AND status = ?", userID, "lost").Count(&stats.LostBets)
	s.db.Model(&SportsBet{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&stats.PendingBets)

	s.db.Model(&SportsBet{}).Where("user_id = ?", userID).Select("COALESCE(SUM(stake), 0)").Scan(&stats.TotalStaked)
	s.db.Model(&SportsBet{}).Where("user_id = ? AND status = ?", userID, "won").Select("COALESCE(SUM(actual_win), 0)").Scan(&stats.TotalWon)

	stats.NetProfit = stats.TotalWon - stats.TotalStaked
	if stats.TotalBets > 0 {
		stats.WinRate = float64(stats.WonBets) / float64(stats.TotalBets) * 100
	}

	return &stats, nil
}

// ============ EVENT GENERATION (FOR DEMO/TESTING) ============

// GenerateEvents creates sample events for testing
func (s *SportsbookService) GenerateEvents() error {
	sports := []struct {
		name    string
		leagues []string
		teams  [][]string
	}{
		{
			"football", []string{"Premier League", "La Liga", "Champions League", "Bundesliga"},
			[][]string{
				{"Manchester United", "Liverpool", "Chelsea", "Arsenal", "Manchester City"},
				{"Real Madrid", "Barcelona", "Atletico Madrid", "Sevilla", "Valencia"},
				{"Bayern Munich", "Dortmund", "RB Leipzig", "Leverkusen", "Wolfsburg"},
			},
		},
		{
			"basketball", []string{"NBA", "EuroLeague"},
			[][]string{
				{"Lakers", "Warriors", "Celtics", "Heat", "Bucks"},
				{"Real Madrid", "Barcelona", "Fenerbahce", "CSKA Moscow"},
			},
		},
		{
			"tennis", []string{"Wimbledon", "US Open", "ATP Finals"},
			[][]string{
				{"Djokovic", "Nadal", "Alcaraz", "Medvedev", "Sinner"},
				{"Swiatek", "Sabalenka", "Rybakina", "Gauff", "Djokovic"},
			},
		},
	}

	engine := NewOddsEngine()

	for _, sport := range sports {
		for leagueIdx, league := range sport.leagues {
			for teamIdx := 0; teamIdx < len(sport.teams[leagueIdx])-1; teamIdx += 2 {
				homeTeam := sport.teams[leagueIdx][teamIdx]
				awayTeam := sport.teams[leagueIdx][teamIdx+1]

				// Random team strengths
				homeStrength := 0.5 + rand.Float64()*0.4
				awayStrength := 0.5 + rand.Float64()*0.4

				homeOdds, drawOdds, awayOdds := engine.CalculateMatchOdds(homeStrength, awayStrength)
				overOdds, underOdds := engine.CalculateOverUnderOdds(2.5, (homeStrength+awayStrength)*2.5)

				startTime := time.Now().Add(time.Duration(rand.Intn(72)) * time.Hour)

				event := SportsEvent{
					ID:            uuid.New(),
					ExternalID:    fmt.Sprintf("%s_%d", sport.name, time.Now().UnixNano()),
					Sport:         sport.name,
					League:        league,
					HomeTeam:      homeTeam,
					AwayTeam:      awayTeam,
					StartTime:     startTime,
					Status:        "upcoming",
					HomeOdds:      homeOdds,
					DrawOdds:      drawOdds,
					AwayOdds:      awayOdds,
					OverUnderLine: 2.5,
					OverOdds:      overOdds,
					UnderOdds:     underOdds,
					Featured:      rand.Float64() > 0.7,
				}

				s.db.Create(&event)
			}
		}
	}

	return nil
}
