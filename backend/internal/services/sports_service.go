package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// SportsService handles sports betting operations
type SportsService struct {
	db *gorm.DB
}

func NewSportsService(db *gorm.DB) *SportsService {
	return &SportsService{db: db}
}

// ListEvents returns upcoming sports events
func (s *SportsService) ListEvents(sport string, status string, limit int) ([]models.SportsEvent, error) {
	var events []models.SportsEvent
	query := s.db.Model(&models.SportsEvent{})

	if sport != "" {
		query = query.Where("sport = ?", sport)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Where("start_time > ?", time.Now()).
		Order("start_time ASC").
		Limit(limit).
		Find(&events).Error

	return events, err
}

// GetEvent returns a single sports event
func (s *SportsService) GetEvent(eventID uuid.UUID) (*models.SportsEvent, error) {
	var event models.SportsEvent
	err := s.db.First(&event, eventID).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetOdds returns current odds for an event
func (s *SportsService) GetOdds(eventID uuid.UUID) (map[string]interface{}, error) {
	event, err := s.GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	// Calculate dynamic odds (simplified)
	// In production, would integrate with odds providers
	baseOdds := 1.85

	odds := map[string]interface{}{
		"event_id":  event.ID,
		"home_win":  baseOdds,
		"draw":      3.20,
		"away_win":  baseOdds,
		"over_25":   1.90,
		"under_25":  1.90,
		"btts_yes":  1.75,
		"btts_no":   2.00,
		"updated_at": time.Now(),
	}

	return odds, nil
}

// PlaceBet places a sports bet
func (s *SportsService) PlaceBet(userID uuid.UUID, eventID uuid.UUID, betType string, stake float64, odds float64) (*models.SportsBet, error) {
	// Get user balance
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	if user.Balance < stake {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Check event is still open for betting
	var event models.SportsEvent
	if err := s.db.First(&event, eventID).Error; err != nil {
		return nil, err
	}

	if event.Status == "finished" || event.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("betting closed for this event")
	}

	// Calculate potential win
	potentialWin := stake * odds

	// Create bet
	bet := models.SportsBet{
		UserID:        userID,
		EventID:       eventID,
		BetType:       betType,
		Odds:          odds,
		Stake:         stake,
		PotentialWin:  potentialWin,
		Status:        "pending",
	}

	if err := s.db.Create(&bet).Error; err != nil {
		return nil, err
	}

	// Deduct stake from balance
	user.Balance -= stake
	s.db.Save(&user)

	// Create transaction
	tx := models.Transaction{
		UserID:      userID,
		Type:        "sports_bet",
		Amount:      stake,
		Currency:    "USD",
		Status:      "confirmed",
		GameID:      &eventID,
		Description: fmt.Sprintf("Sports bet: %s - %s", event.HomeTeam, event.AwayTeam),
	}
	s.db.Create(&tx)

	return &bet, nil
}

// SettleBet settles a sports bet based on result
func (s *SportsService) SettleBet(betID uuid.UUID, result string) error {
	var bet models.SportsBet
	if err := s.db.First(&bet, betID).Error; err != nil {
		return err
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet already settled")
	}

	bet.Result = result
	bet.Status = "settled"
	s.db.Save(&bet)

	// Check if won
	isWin := bet.BetType == result

	if isWin {
		// Credit winnings
		var user models.User
		s.db.First(&user, bet.UserID)

		user.Balance += bet.PotentialWin
		s.db.Save(&user)

		// Create win transaction
		tx := models.Transaction{
			UserID:      bet.UserID,
			Type:        "sports_win",
			Amount:      bet.PotentialWin,
			Currency:    "USD",
			Status:      "confirmed",
			GameID:      &bet.EventID,
			Description: fmt.Sprintf("Sports bet win: %s", result),
		}
		s.db.Create(&tx)

		bet.ActualWin = bet.PotentialWin
	} else {
		bet.ActualWin = 0
		bet.Status = "lost"
	}

	now := time.Now()
	bet.SettledAt = &now
	s.db.Save(&bet)

	return nil
}

// GetMyBets returns user's sports bets
func (s *SportsService) GetMyBets(userID uuid.UUID, status string, limit int) ([]models.SportsBet, error) {
	var bets []models.SportsBet
	query := s.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&bets).Error
	return bets, err
}

// SeedSportsEvents seeds some sample sports events
func (s *SportsService) SeedSportsEvents() error {
	var count int64
	s.db.Model(&models.SportsEvent{}).Count(&count)

	if count > 0 {
		return nil
	}

	events := []models.SportsEvent{
		{
			Sport:     "Football",
			League:    "Premier League",
			HomeTeam:  "Manchester City",
			AwayTeam:  "Liverpool",
			StartTime: time.Now().Add(2 * time.Hour),
			Status:    "upcoming",
		},
		{
			Sport:     "Football",
			League:    "La Liga",
			HomeTeam:  "Real Madrid",
			AwayTeam:  "Barcelona",
			StartTime: time.Now().Add(24 * time.Hour),
			Status:    "upcoming",
		},
		{
			Sport:     "Basketball",
			League:    "NBA",
			HomeTeam:  "Lakers",
			AwayTeam:  "Warriors",
			StartTime: time.Now().Add(3 * time.Hour),
			Status:    "upcoming",
		},
		{
			Sport:     "Tennis",
			League:    "ATP",
			HomeTeam:  "Djokovic",
			AwayTeam:  "Alcaraz",
			StartTime: time.Now().Add(5 * time.Hour),
			Status:    "upcoming",
		},
	}

	for i := range events {
		s.db.Create(&events[i])
	}

	return nil
}
