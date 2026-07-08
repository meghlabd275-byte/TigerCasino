package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// TournamentService handles casino tournaments
type TournamentService struct {
	db *gorm.DB
}

func NewTournamentService(db *gorm.DB) *TournamentService {
	return &TournamentService{db: db}
}

// ============ TOURNAMENT MODELS ============

type Tournament struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key"`
	Name            string     `gorm:"not null"`
	Description    string
	GameType       string     `gorm:"not null"` // crash, slots, dice, etc.
	Status          string     `gorm:"default:'upcoming'"` // upcoming, registration, active, completed, cancelled
	StartTime      time.Time
	EndTime        time.Time
	MinBet         float64    `gorm:"default:0"`
	MaxBet         float64    `gorm:"default:0"`
	EntryFee       float64    `gorm:"default:0"`
	TotalPrizePool float64    `gorm:"not null"`
	1stPrize       float64
	2ndPrize       float64
	3rdPrize       float64
	4thPrize       float64
	5thPrize       float64
	MaxParticipants int        `gorm:"default:0"` // 0 = unlimited
	CreatedAt      time.Time
}

type TournamentParticipant struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key"`
	TournamentID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	Username     string
	Score        float64    `gorm:"default:0"`
	BestMultiplier float64   `gorm:"default:0"`
	TotalBets    int        `gorm:"default:0"`
	Wins         int        `gorm:"default:0"`
	Rank         int        `gorm:"default:0"`
	IsActive     bool       `gorm:"default:true"`
	JoinedAt     time.Time
	UpdatedAt    time.Time
}

type TournamentBet struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key"`
	TournamentID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	BetID        uuid.UUID  `gorm:"not null"`
	Amount       float64    `gorm:"not null"`
	Multiplier   float64   `gorm:"not null"`
	Score        float64    // Score contributed to tournament
	CreatedAt    time.Time
}

// ============ TOURNAMENT CRUD ============

func (s *TournamentService) CreateTournament(t *Tournament) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	return s.db.Create(t).Error
}

func (s *TournamentService) GetTournament(id uuid.UUID) (*Tournament, error) {
	var t Tournament
	err := s.db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TournamentService) GetAllTournaments(status string) ([]Tournament, error) {
	var tournaments []Tournament
	query := s.db.Order("start_time DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&tournaments).Error
	return tournaments, err
}

func (s *TournamentService) GetActiveTournaments() ([]Tournament, error) {
	var tournaments []Tournament
	now := time.Now()

	err := s.db.Where("status IN ('upcoming', 'registration', 'active') AND end_time > ?", now).
		Order("start_time ASC").
		Find(&tournaments).Error

	return tournaments, err
}

func (s *TournamentService) UpdateTournament(id uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&Tournament{}).Where("id = ?", id).Updates(updates).Error
}

func (s *TournamentService) CancelTournament(id uuid.UUID) error {
	// Refund entry fees to participants
	var participants []TournamentParticipant
	s.db.Where("tournament_id = ?", id).Find(&participants)

	userService := NewUserService(s.db)

	var tournament Tournament
	s.db.First(&tournament, id)

	if tournament.EntryFee > 0 {
		for _, p := range participants {
			userService.UpdateBalance(p.UserID, tournament.EntryFee)
		}
	}

	return s.db.Model(&Tournament{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "cancelled"}).Error
}

// ============ PARTICIPATION ============

func (s *TournamentService) JoinTournament(tournamentID, userID uuid.UUID, username string) (*TournamentParticipant, error) {
	var tournament Tournament
	if err := s.db.First(&tournament, tournamentID).Error; err != nil {
		return nil, fmt.Errorf("tournament not found")
	}

	// Check if tournament allows joining
	if tournament.Status != "registration" && tournament.Status != "upcoming" {
		return nil, fmt.Errorf("tournament is not accepting participants")
	}

	// Check if already joined
	var existing TournamentParticipant
	err := s.db.Where("tournament_id = ? AND user_id = ?", tournamentID, userID).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("already joined this tournament")
	}

	// Check max participants
	if tournament.MaxParticipants > 0 {
		var count int64
		s.db.Model(&TournamentParticipant{}).Where("tournament_id = ?", tournamentID).Count(&count)
		if count >= int64(tournament.MaxParticipants) {
			return nil, fmt.Errorf("tournament is full")
		}
	}

	// Deduct entry fee if applicable
	if tournament.EntryFee > 0 {
		userService := NewUserService(s.db)
		var user models.User
		s.db.First(&user, userID)
		if user.Balance < tournament.EntryFee {
			return nil, fmt.Errorf("insufficient balance for entry fee")
		}
		userService.UpdateBalance(userID, -tournament.EntryFee)
	}

	// Create participant
	participant := TournamentParticipant{
		ID:            uuid.New(),
		TournamentID:  tournamentID,
		UserID:       userID,
		Username:     username,
		Score:        0,
		BestMultiplier: 0,
		TotalBets:    0,
		Wins:         0,
		IsActive:     true,
		JoinedAt:     time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(&participant).Error; err != nil {
		return nil, err
	}

	return &participant, nil
}

func (s *TournamentService) LeaveTournament(tournamentID, userID uuid.UUID) error {
	var tournament Tournament
	s.db.First(&tournament, tournamentID)

	// Only allow leaving if tournament hasn't started
	if tournament.Status == "active" || tournament.Status == "completed" {
		return fmt.Errorf("cannot leave after tournament has started")
	}

	// Refund entry fee
	if tournament.EntryFee > 0 {
		userService := NewUserService(s.db)
		userService.UpdateBalance(userID, tournament.EntryFee)
	}

	return s.db.Where("tournament_id = ? AND user_id = ?", tournamentID, userID).
		Delete(&TournamentParticipant{}).Error
}

func (s *TournamentService) GetParticipant(tournamentID, userID uuid.UUID) (*TournamentParticipant, error) {
	var p TournamentParticipant
	err := s.db.Where("tournament_id = ? AND user_id = ?", tournamentID, userID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *TournamentService) GetTournamentParticipants(tournamentID uuid.UUID) ([]TournamentParticipant, error) {
	var participants []TournamentParticipant
	err := s.db.Where("tournament_id = ? AND is_active = ?", tournamentID, true).
		Order("score DESC").
		Find(&participants).Error
	return participants, err
}

// ============ SCORING ============

// ScoreCalculation determines how points are calculated for each game type
func calculateScore(gameType string, betAmount, multiplier float64) float64 {
	switch gameType {
	case "crash":
		// Score = profit from crash game
		return betAmount * (multiplier - 1)
	case "slots":
		// Score = profit from slots
		return betAmount * (multiplier - 1)
	case "dice":
		// Score = profit from dice
		return betAmount * (multiplier - 1)
	case "limbo":
		// Score = profit from limbo
		return betAmount * (multiplier - 1)
	case "plinko":
		// Score = profit from plinko
		return betAmount * (multiplier - 1)
	case "mines":
		// Score = profit from mines
		return betAmount * (multiplier - 1)
	default:
		// Default: score is the profit
		return betAmount * (multiplier - 1)
	}
}

// RecordTournamentBet records a bet for tournament scoring
func (s *TournamentService) RecordTournamentBet(tournamentID, userID, betID uuid.UUID, betAmount, multiplier float64) error {
	var tournament Tournament
	if err := s.db.First(&tournament, tournamentID).Error; err != nil {
		return fmt.Errorf("tournament not found")
	}

	// Only active tournaments accept bets
	if tournament.Status != "active" {
		return nil
	}

	// Check minimum bet
	if tournament.MinBet > 0 && betAmount < tournament.MinBet {
		return nil
	}

	// Check maximum bet
	if tournament.MaxBet > 0 && betAmount > tournament.MaxBet {
		betAmount = tournament.MaxBet
	}

	// Check game type matches
	if tournament.GameType != "all" && tournament.GameType != gameType {
		return nil
	}

	// Calculate score
	score := calculateScore(tournament.GameType, betAmount, multiplier)

	// Create tournament bet record
	tBet := TournamentBet{
		ID:           uuid.New(),
		TournamentID: tournamentID,
		UserID:      userID,
		BetID:       betID,
		Amount:      betAmount,
		Multiplier:  multiplier,
		Score:       score,
		CreatedAt:   time.Now(),
	}

	if err := s.db.Create(&tBet).Error; err != nil {
		return err
	}

	// Update participant stats
	if score > 0 {
		s.db.Model(&TournamentParticipant{}).
			Where("tournament_id = ? AND user_id = ?", tournamentID, userID).
			Updates(map[string]interface{}{
				"score":           gorm.Expr("score + ?", score),
				"best_multiplier": gorm.Expr("GREATEST(best_multiplier, ?)", multiplier),
				"total_bets":      gorm.Expr("total_bets + 1"),
				"wins":           gorm.Expr("wins + 1"),
				"updated_at":      time.Now(),
			})
	} else {
		s.db.Model(&TournamentParticipant{}).
			Where("tournament_id = ? AND user_id = ?", tournamentID, userID).
			Updates(map[string]interface{}{
				"total_bets": gorm.Expr("total_bets + 1"),
				"updated_at": time.Now(),
			})
	}

	return nil
}

// RecalculateRankings recalculates all rankings for a tournament
func (s *TournamentService) RecalculateRankings(tournamentID uuid.UUID) error {
	// Get all participants ordered by score
	var participants []TournamentParticipant
	s.db.Where("tournament_id = ?", tournamentID).Order("score DESC").Find(&participants)

	rank := 1
	for i, p := range participants {
		if i > 0 && participants[i].Score == participants[i-1].Score {
			// Same score, same rank (handle ties)
		} else {
			rank = i + 1
		}
		s.db.Model(&p).Update("rank", rank)
	}

	return nil
}

// GetLeaderboard returns the tournament leaderboard
func (s *TournamentService) GetLeaderboard(tournamentID uuid.UUID, limit int) ([]TournamentParticipant, error) {
	if limit <= 0 {
		limit = 100
	}

	var participants []TournamentParticipant
	err := s.db.Where("tournament_id = ? AND is_active = ?", tournamentID, true).
		Order("rank ASC").
		Limit(limit).
		Find(&participants).Error

	return participants, err
}

// GetUserRanking returns a user's ranking in a tournament
func (s *TournamentService) GetUserRanking(tournamentID, userID uuid.UUID) (*TournamentParticipant, int, error) {
	var p TournamentParticipant
	err := s.db.Where("tournament_id = ? AND user_id = ?", tournamentID, userID).First(&p).Error
	if err != nil {
		return nil, 0, err
	}

	// Count users with higher score
	var higherCount int64
	s.db.Model(&TournamentParticipant{}).
		Where("tournament_id = ? AND score > ?", tournamentID, p.Score).
		Count(&higherCount)

	return &p, int(higherCount) + 1, nil
}

// ============ PRIZE DISTRIBUTION ============

// DistributePrizes distributes prizes to winners
func (s *TournamentService) DistributePrizes(tournamentID uuid.UUID) error {
	var tournament Tournament
	if err := s.db.First(&tournament, tournamentID).Error; err != nil {
		return err
	}

	if tournament.Status != "active" && tournament.Status != "completed" {
		return fmt.Errorf("tournament is not active or completed")
	}

	// Get top participants
	var winners []TournamentParticipant
	s.db.Where("tournament_id = ?", tournamentID).
		Order("rank ASC").
		Limit(5).
		Find(&winners)

	userService := NewUserService(s.db)

	prizes := []float64{tournament.1stPrize, tournament.2ndPrize, tournament.3rdPrize, tournament.4thPrize, tournament.5thPrize}

	for i, winner := range winners {
		if i < len(prizes) && prizes[i] > 0 {
			userService.UpdateBalance(winner.UserID, prizes[i])

			// Create transaction
			tx := models.Transaction{
				ID:          uuid.New(),
				UserID:      winner.UserID,
				Type:        "tournament_win",
				Amount:      prizes[i],
				Currency:    "USD",
				Status:      "confirmed",
				Description: fmt.Sprintf("Tournament prize - %s place", positionName(i+1)),
			}
			s.db.Create(&tx)
		}
	}

	// Mark tournament as completed
	s.db.Model(&tournament).Update("status", "completed")

	return nil
}

func positionName(pos int) string {
	switch pos {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	case 4:
		return "4th"
	case 5:
		return "5th"
	default:
		return fmt.Sprintf("%dth", pos)
	}
}

// ============ TOURNAMENT AUTOMATION ============

// StartTournament starts a tournament
func (s *TournamentService) StartTournament(id uuid.UUID) error {
	var tournament Tournament
	if err := s.db.First(&tournament, id).Error; err != nil {
		return err
	}

	if tournament.Status != "upcoming" && tournament.Status != "registration" {
		return fmt.Errorf("tournament cannot be started")
	}

	return s.db.Model(&tournament).Update("status", "active").Error
}

// EndTournament ends a tournament and distributes prizes
func (s *TournamentService) EndTournament(id uuid.UUID) error {
	var tournament Tournament
	if err := s.db.First(&tournament, id).Error; err != nil {
		return err
	}

	if tournament.Status != "active" {
		return fmt.Errorf("tournament is not active")
	}

	// Finalize rankings
	s.RecalculateRankings(id)

	// Distribute prizes
	s.DistributePrizes(id)

	return nil
}

// ============ SCHEDULED TOURNAMENTS ============

// CreateScheduledTournaments creates tournaments for the week
func (s *TournamentService) CreateScheduledTournaments() error {
	// Daily tournaments
	dailyTournaments := []struct {
		name        string
		gameType    string
		prizePool   float64
		startHour   int
		durationHrs int
	}{
		{"Daily Crash Championship", "crash", 1000, 0, 24},
		{"Daily Slots Marathon", "slots", 1000, 6, 18},
		{"Daily Dice Duel", "dice", 500, 12, 12},
	}

	now := time.Now()

	for _, t := range dailyTournaments {
		startTime := time.Date(now.Year(), now.Month(), now.Day(), t.startHour, 0, 0, 0, now.Location())
		endTime := startTime.Add(time.Duration(t.durationHrs) * time.Hour)

		// Check if already exists
		var existing Tournament
		err := s.db.Where("name = ? AND start_time = ?", t.name, startTime).First(&existing).Error
		if err == nil {
			continue // Already exists
		}

		tournament := Tournament{
			Name:            t.name,
			Description:    fmt.Sprintf("Daily %s tournament with $%.0f prize pool", t.gameType, t.prizePool),
			GameType:        t.gameType,
			Status:          "upcoming",
			StartTime:       startTime,
			EndTime:         endTime,
			MinBet:          1.0,
			MaxBet:          100.0,
			EntryFee:        0,
			TotalPrizePool:  t.prizePool,
			1stPrize:        t.prizePool * 0.5,
			2ndPrize:        t.prizePool * 0.25,
			3rdPrize:        t.prizePool * 0.12,
			4thPrize:        t.prizePool * 0.08,
			5thPrize:        t.prizePool * 0.05,
		}

		s.CreateTournament(&tournament)
	}

	// Weekly tournament
	weekday := int(now.Weekday())
	if weekday == 1 { // Monday
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime := startTime.AddDate(0, 0, 7)

		var existing Tournament
		err := s.db.Where("name = ? AND start_time >= ?", "Weekly Grand Prix", startTime).First(&existing).Error
		if err != nil {
			tournament := Tournament{
				Name:            "Weekly Grand Prix",
				Description:    "The biggest weekly tournament with $10,000 prize pool!",
				GameType:        "all",
				Status:          "upcoming",
				StartTime:       startTime,
				EndTime:         endTime,
				MinBet:          5.0,
				MaxBet:          500.0,
				EntryFee:        0,
				TotalPrizePool:  10000,
				1stPrize:        5000,
				2ndPrize:        2500,
				3rdPrize:        1200,
				4thPrize:        800,
				5thPrize:        500,
			}
			s.CreateTournament(&tournament)
		}
	}

	return nil
}

// ============ TOURNAMENT HISTORY ============

type TournamentHistory struct {
	TournamentID   uuid.UUID
	Name           string
	GameType       string
	StartTime      time.Time
	EndTime        time.Time
	TotalPrizePool float64
	Rank           int
	PrizeWon       float64
	Score          float64
}

func (s *TournamentService) GetUserTournamentHistory(userID uuid.UUID, limit int) ([]TournamentHistory, error) {
	if limit <= 0 {
		limit = 20
	}

	var participants []TournamentParticipant
	err := s.db.Where("user_id = ?", userID).
		Order("joined_at DESC").
		Limit(limit).
		Find(&participants).Error

	if err != nil {
		return nil, err
	}

	history := make([]TournamentHistory, len(participants))
	for i, p := range participants {
		var t Tournament
		s.db.First(&t, p.TournamentID)

		prize := 0.0
		if p.Rank == 1 {
			prize = t.1stPrize
		} else if p.Rank == 2 {
			prize = t.2ndPrize
		} else if p.Rank == 3 {
			prize = t.3rdPrize
		} else if p.Rank == 4 {
			prize = t.4thPrize
		} else if p.Rank == 5 {
			prize = t.5thPrize
		}

		history[i] = TournamentHistory{
			TournamentID:   p.TournamentID,
			Name:           t.Name,
			GameType:       t.GameType,
			StartTime:      t.StartTime,
			EndTime:        t.EndTime,
			TotalPrizePool: t.TotalPrizePool,
			Rank:           p.Rank,
			PrizeWon:       prize,
			Score:          p.Score,
		}
	}

	return history, nil
}
