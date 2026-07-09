package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// TournamentService manages casino tournaments
type TournamentService struct {
	db           *gorm.DB
	redis        *redis.Client
	config       *TournamentConfig
	activeTournaments map[uuid.UUID]*Tournament
	tournamentMu sync.RWMutex
}

type TournamentConfig struct {
	// Prize pool settings
	MinPrizePool    float64
	DefaultPrizePool float64
	
	// Registration settings
	MinParticipants int
	MaxParticipants int
	RegistrationEndDuration time.Duration
	
	// Scoring settings
	PointsPerWin    float64
	PointsPerBet    float64  // Points per currency unit wagered
	MultiplierBonus float64  // Bonus for winning streaks
	
	// Time settings
	TournamentCheckInterval time.Duration
	LeaderboardUpdateInterval time.Duration
	ResultsFinalizeDelay   time.Duration
}

// Tournament represents a casino tournament
type Tournament struct {
	ID              uuid.UUID
	Name            string
	Description     string
	Type            string // slots, table_games, live_casino, all_games, sports
	Status          string // upcoming, registration, active, completed, cancelled
	GameFilter      []string // allowed game IDs or categories
	MinBet          float64
	StartTime       time.Time
	EndTime         time.Time
	RegistrationEnd time.Time
	
	// Prize settings
	PrizePool        float64
	Currency         string
	PrizeDistribution []PrizeBreakdown
	MinWagerToQualify float64
	
	// Scoring
	ScoringType     string // wager, wins, profit
	PointsMultiplier float64
	
	// Participants
	RegisteredUsers map[uuid.UUID]bool
	ParticipantCount int
	
	// State
	CurrentLeaderboard []LeaderboardEntry
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

type PrizeBreakdown struct {
	Position   int
	Percent    float64
	Amount     float64
	MinWager  float64
}

type LeaderboardEntry struct {
	Rank        int
	UserID      uuid.UUID
	Username    string
	Score       float64
	Wagered     float64
	Wins        int
	WinStreak   int
	UpdatedAt   time.Time
}

type TournamentParticipant struct {
	UserID       uuid.UUID
	TournamentID uuid.UUID
	JoinedAt    time.Time
	Score        float64
	Wagered      float64
	Wins         int
	CurrentStreak int
	BestStreak   int
	Position     int
}

// NewTournamentService creates a new tournament service
func NewTournamentService(db *gorm.DB, redisClient *redis.Client, config *TournamentConfig) *TournamentService {
	if config == nil {
		config = &TournamentConfig{
			MinPrizePool:         100,
			DefaultPrizePool:     1000,
			MinParticipants:      10,
			MaxParticipants:      10000,
			RegistrationEndDuration: 1 * time.Hour,
			PointsPerWin:         100,
			PointsPerBet:         1,
			MultiplierBonus:     0.1,
			TournamentCheckInterval: 1 * time.Minute,
			LeaderboardUpdateInterval: 30 * time.Second,
			ResultsFinalizeDelay:   5 * time.Minute,
		}
	}
	
	return &TournamentService{
		db:            db,
		redis:         redisClient,
		config:        config,
		activeTournaments: make(map[uuid.UUID]*Tournament),
	}
}

// ============ TOURNAMENT MANAGEMENT ============

// CreateTournament creates a new tournament
func (s *TournamentService) CreateTournament(ctx context.Context, req *CreateTournamentRequest) (*Tournament, error) {
	if req.PrizePool < s.config.MinPrizePool {
		return nil, fmt.Errorf("minimum prize pool is %.2f", s.config.MinPrizePool)
	}
	
	tournament := &Tournament{
		ID:                  uuid.New(),
		Name:                req.Name,
		Description:         req.Description,
		Type:               req.Type,
		Status:             "upcoming",
		GameFilter:         req.GameFilter,
		MinBet:             req.MinBet,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		RegistrationEnd:    req.StartTime.Add(-time.Hour),
		PrizePool:          req.PrizePool,
		Currency:           "USD",
		ScoringType:        req.ScoringType,
		PointsMultiplier:   req.PointsMultiplier,
		MinWagerToQualify:  req.MinWagerToQualify,
		RegisteredUsers:    make(map[uuid.UUID]bool),
		CurrentLeaderboard: make([]LeaderboardEntry, 0),
		CreatedAt:          time.Now(),
	}
	
	// Calculate prize distribution
	tournament.PrizeDistribution = s.calculatePrizeDistribution(tournament.PrizePool)
	
	// Save to database
	tournamentModel := models.Tournament{
		ID:               tournament.ID,
		Name:             tournament.Name,
		Description:      tournament.Description,
		Type:             tournament.Type,
		Status:           tournament.Status,
		GameFilter:       strings.Join(tournament.GameFilter, ","),
		MinBet:           tournament.MinBet,
		StartTime:        tournament.StartTime,
		EndTime:          tournament.EndTime,
		RegistrationEnd:  tournament.RegistrationEnd,
		PrizePool:        tournament.PrizePool,
		Currency:         tournament.Currency,
		ScoringType:      tournament.ScoringType,
		PointsMultiplier: tournament.PointsMultiplier,
		MinWagerToQualify: tournament.MinWagerToQualify,
		CreatedAt:        tournament.CreatedAt,
	}
	
	if err := s.db.Create(&tournamentModel).Error; err != nil {
		return nil, err
	}
	
	// Add to active tournaments
	s.tournamentMu.Lock()
	s.activeTournaments[tournament.ID] = tournament
	s.tournamentMu.Unlock()
	
	return tournament, nil
}

type CreateTournamentRequest struct {
	Name              string
	Description       string
	Type              string
	GameFilter        []string
	MinBet            float64
	StartTime         time.Time
	EndTime           time.Time
	PrizePool         float64
	ScoringType       string
	PointsMultiplier  float64
	MinWagerToQualify float64
}

func (s *TournamentService) calculatePrizeDistribution(prizePool float64) []PrizeBreakdown {
	// Standard prize distribution: Top 10 get prizes
	distribution := []PrizeBreakdown{
		{Position: 1, Percent: 30},
		{Position: 2, Percent: 20},
		{Position: 3, Percent: 12},
		{Position: 4, Percent: 8},
		{Position: 5, Percent: 6},
		{Position: 6, Percent: 5},
		{Position: 7, Percent: 4},
		{Position: 8, Percent: 3},
		{Position: 9, Percent: 2},
		{Position: 10, Percent: 2},
	}
	
	// Adjust for smaller prize pools
	if prizePool < 1000 {
		distribution = []PrizeBreakdown{
			{Position: 1, Percent: 40},
			{Position: 2, Percent: 25},
			{Position: 3, Percent: 15},
			{Position: 4, Percent: 10},
			{Position: 5, Percent: 10},
		}
	}
	
	// Calculate amounts
	result := make([]PrizeBreakdown, len(distribution))
	for i, d := range distribution {
		result[i] = PrizeBreakdown{
			Position: d.Position,
			Percent:  d.Percent,
			Amount:   prizePool * (d.Percent / 100),
		}
	}
	
	return result
}

// StartTournament starts a tournament
func (s *TournamentService) StartTournament(ctx context.Context, tournamentID uuid.UUID) error {
	s.tournamentMu.Lock()
	defer s.tournamentMu.Unlock()
	
	tournament, ok := s.activeTournaments[tournamentID]
	if !ok {
		return fmt.Errorf("tournament not found")
	}
	
	if tournament.Status != "upcoming" && tournament.Status != "registration" {
		return fmt.Errorf("tournament cannot be started")
	}
	
	tournament.Status = "active"
	tournament.UpdatedAt = time.Now()
	
	// Update database
	return s.db.Model(&models.Tournament{}).
		Where("id = ?", tournamentID).
		Update("status", "active").Error
}

// EndTournament ends a tournament and distributes prizes
func (s *TournamentService) EndTournament(ctx context.Context, tournamentID uuid.UUID) error {
	s.tournamentMu.Lock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.tournamentMu.Unlock()
	
	if !ok {
		return fmt.Errorf("tournament not found")
	}
	
	if tournament.Status != "active" {
		return fmt.Errorf("tournament is not active")
	}
	
	tournament.Status = "completed"
	tournament.UpdatedAt = time.Now()
	
	// Get final leaderboard
	leaderboard, err := s.GetLeaderboard(ctx, tournamentID, 100)
	if err != nil {
		return err
	}
	
	// Distribute prizes
	if err := s.distributePrizes(ctx, tournament, leaderboard); err != nil {
		return err
	}
	
	// Update database
	s.db.Model(&models.Tournament{}).
		Where("id = ?", tournamentID).
		Update("status", "completed")
	
	return nil
}

func (s *TournamentService) distributePrizes(ctx context.Context, tournament *Tournament, leaderboard []LeaderboardEntry) error {
	if len(leaderboard) < s.config.MinParticipants {
		// Not enough participants - return prizes
		return fmt.Errorf("not enough participants")
	}
	
	for i, entry := range leaderboard {
		if i >= len(tournament.PrizeDistribution) {
			break
		}
		
		prize := tournament.PrizeDistribution[i]
		
		// Credit prize to user
		var wallet models.Wallet
		if err := s.db.Where("user_id = ? AND currency = ?", entry.UserID, tournament.Currency).First(&wallet).Error; err != nil {
			continue
		}
		
		wallet.Balance += prize.Amount
		wallet.UpdatedAt = time.Now()
		s.db.Save(&wallet)
		
		// Create prize record
		prizeRecord := models.TournamentPrize{
			ID:           uuid.New(),
			TournamentID: tournament.ID,
			UserID:       entry.UserID,
			Position:     prize.Position,
			PrizeAmount: prize.Amount,
			Currency:    tournament.Currency,
			CreatedAt:    time.Now(),
		}
		s.db.Create(&prizeRecord)
		
		// Create notification
		notification := models.Notification{
			ID:        uuid.New(),
			UserID:    entry.UserID,
			Type:      "tournament_prize",
			Title:     "Tournament Prize Won!",
			Message:   fmt.Sprintf("Congratulations! You won %d place in %s and received %.2f %s", 
				prize.Position, tournament.Name, prize.Amount, tournament.Currency),
			CreatedAt: time.Now(),
		}
		s.db.Create(&notification)
	}
	
	return nil
}

// RegisterUser registers a user for a tournament
func (s *TournamentService) RegisterUser(ctx context.Context, tournamentID, userID uuid.UUID) error {
	s.tournamentMu.Lock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.tournamentMu.Unlock()
	
	if !ok {
		return fmt.Errorf("tournament not found")
	}
	
	if tournament.Status != "upcoming" && tournament.Status != "registration" {
		return fmt.Errorf("tournament registration is closed")
	}
	
	if time.Now().After(tournament.RegistrationEnd) {
		return fmt.Errorf("registration period has ended")
	}
	
	if tournament.RegisteredUsers[userID] {
		return fmt.Errorf("already registered")
	}
	
	if tournament.MaxParticipants > 0 && tournament.ParticipantCount >= tournament.MaxParticipants {
		return fmt.Errorf("tournament is full")
	}
	
	// Add to tournament
	tournament.RegisteredUsers[userID] = true
	tournament.ParticipantCount++
	tournament.UpdatedAt = time.Now()
	
	// Save to database
	participant := models.TournamentParticipant{
		ID:           uuid.New(),
		TournamentID: tournamentID,
		UserID:       userID,
		JoinedAt:    time.Now(),
	}
	s.db.Create(&participant)
	
	return nil
}

// UnregisterUser removes a user from a tournament
func (s *TournamentService) UnregisterUser(ctx context.Context, tournamentID, userID uuid.UUID) error {
	s.tournamentMu.Lock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.tournamentMu.Unlock()
	
	if !ok {
		return fmt.Errorf("tournament not found")
	}
	
	if tournament.Status == "active" {
		return fmt.Errorf("cannot unregister from active tournament")
	}
	
	if !tournament.RegisteredUsers[userID] {
		return fmt.Errorf("not registered")
	}
	
	delete(tournament.RegisteredUsers, userID)
	tournament.ParticipantCount--
	tournament.UpdatedAt = time.Now()
	
	// Remove from database
	s.db.Where("tournament_id = ? AND user_id = ?", tournamentID, userID).
		Delete(&models.TournamentParticipant{})
	
	return nil
}

// ============ SCORING ============

// RecordScore records a score for a user in a tournament
func (s *TournamentService) RecordScore(ctx context.Context, userID, gameID uuid.UUID, betAmount, winAmount float64, isWin bool) error {
	s.tournamentMu.RLock()
	defer s.tournamentMu.RUnlock()
	
	for _, tournament := range s.activeTournaments {
		if tournament.Status != "active" {
			continue
		}
		
		// Check if user is registered
		if !tournament.RegisteredUsers[userID] {
			continue
		}
		
		// Check if game is eligible
		if !s.isGameEligible(tournament, gameID.String()) {
			continue
		}
		
		// Check minimum bet
		if betAmount < tournament.MinBet {
			continue
		}
		
		// Calculate score
		score := s.calculateScore(tournament, betAmount, winAmount, isWin)
		
		// Update participant score
		s.updateParticipantScore(ctx, tournament.ID, userID, score, betAmount, isWin)
	}
	
	return nil
}

func (s *TournamentService) isGameEligible(tournament *Tournament, gameID string) bool {
	if len(tournament.GameFilter) == 0 {
		return true
	}
	
	for _, filter := range tournament.GameFilter {
		if filter == gameID || filter == "*" {
			return true
		}
	}
	
	return false
}

func (s *TournamentService) calculateScore(tournament *Tournament, betAmount, winAmount float64, isWin bool) float64 {
	multiplier := tournament.PointsMultiplier
	if multiplier == 0 {
		multiplier = 1
	}
	
	switch tournament.ScoringType {
	case "wager":
		return betAmount * s.config.PointsPerBet * multiplier
	case "wins":
		if isWin {
			return s.config.PointsPerWin * multiplier
		}
		return 0
	case "profit":
		profit := winAmount - betAmount
		if profit > 0 {
			return profit * 10 * multiplier
		}
		return 0
	default:
		return betAmount * s.config.PointsPerBet * multiplier
	}
}

func (s *TournamentService) updateParticipantScore(ctx context.Context, tournamentID, userID uuid.UUID, score float64, betAmount float64, isWin bool) {
	var participant models.TournamentParticipant
	result := s.db.Where("tournament_id = ? AND user_id = ?", tournamentID, userID).First(&participant)
	
	if result.Error != nil {
		// Create new participant
		participant = models.TournamentParticipant{
			ID:           uuid.New(),
			TournamentID: tournamentID,
			UserID:       userID,
			JoinedAt:    time.Now(),
			Score:        score,
			Wagered:      betAmount,
		}
		if isWin {
			participant.Wins = 1
			participant.CurrentStreak = 1
			participant.BestStreak = 1
		}
		s.db.Create(&participant)
	} else {
		// Update existing participant
		participant.Score += score
		participant.Wagered += betAmount
		if isWin {
			participant.Wins++
			participant.CurrentStreak++
			if participant.CurrentStreak > participant.BestStreak {
				participant.BestStreak = participant.CurrentStreak
			}
		} else {
			participant.CurrentStreak = 0
		}
		s.db.Save(&participant)
	}
}

// ============ LEADERBOARD ============

// GetLeaderboard returns the tournament leaderboard
func (s *TournamentService) GetLeaderboard(ctx context.Context, tournamentID uuid.UUID, limit int) ([]LeaderboardEntry, error) {
	s.tournamentMu.RLock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.tournamentMu.RUnlock()
	
	if !ok {
		return nil, fmt.Errorf("tournament not found")
	}
	
	// Try cache first
	cacheKey := fmt.Sprintf("tournament:%s:leaderboard:%d", tournamentID.String(), limit)
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var leaderboard []LeaderboardEntry
			if json.Unmarshal([]byte(cached), &leaderboard) == nil {
				return leaderboard, nil
			}
		}
	}
	
	// Query database
	var participants []models.TournamentParticipant
	s.db.Where("tournament_id = ?", tournamentID).
		Order("score DESC").
		Limit(limit).
		Find(&participants)
	
	leaderboard := make([]LeaderboardEntry, len(participants))
	
	for i, p := range participants {
		var user models.User
		s.db.First(&user, p.UserID)
		
		leaderboard[i] = LeaderboardEntry{
			Rank:      i + 1,
			UserID:    p.UserID,
			Username:  user.Username,
			Score:     p.Score,
			Wagered:   p.Wagered,
			Wins:      p.Wins,
			WinStreak: p.BestStreak,
			UpdatedAt: p.JoinedAt,
		}
	}
	
	// Update tournament leaderboard
	s.tournamentMu.Lock()
	if tournament != nil {
		tournament.CurrentLeaderboard = leaderboard
	}
	s.tournamentMu.Unlock()
	
	// Cache leaderboard
	if s.redis != nil && len(leaderboard) > 0 {
		if data, err := json.Marshal(leaderboard); err == nil {
			s.redis.Set(ctx, cacheKey, data, s.config.LeaderboardUpdateInterval)
		}
	}
	
	return leaderboard, nil
}

// GetGlobalLeaderboard returns the global leaderboard across all active tournaments
func (s *TournamentService) GetGlobalLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	s.tournamentMu.RLock()
	defer s.tournamentMu.RUnlock()
	
	type aggregateScore struct {
		UserID    uuid.UUID
		TotalScore float64
		TotalWager float64
		TotalWins  int
	}
	
	var scores []aggregateScore
	
	s.db.Raw(`
		SELECT user_id, SUM(score) as total_score, SUM(wagered) as total_wager, SUM(wins) as total_wins
		FROM tournament_participants tp
		JOIN tournaments t ON tp.tournament_id = t.id
		WHERE t.status = 'active'
		GROUP BY user_id
		ORDER BY total_score DESC
		LIMIT ?
	`, limit).Scan(&scores)
	
	leaderboard := make([]LeaderboardEntry, len(scores))
	
	for i, score := range scores {
		var user models.User
		s.db.First(&user, score.UserID)
		
		leaderboard[i] = LeaderboardEntry{
			Rank:      i + 1,
			UserID:    score.UserID,
			Username:  user.Username,
			Score:     score.TotalScore,
			Wagered:   score.TotalWager,
			Wins:      score.TotalWins,
		}
	}
	
	return leaderboard, nil
}

// ============ TOURNAMENT QUERIES ============

// GetActiveTournaments returns all active tournaments
func (s *TournamentService) GetActiveTournaments(ctx context.Context) ([]Tournament, error) {
	s.tournamentMu.RLock()
	defer s.tournamentMu.RUnlock()
	
	result := make([]Tournament, 0)
	for _, t := range s.activeTournaments {
		if t.Status == "active" || t.Status == "registration" || t.Status == "upcoming" {
			result = append(result, *t)
		}
	}
	
	// Sort by start time
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartTime.Before(result[j].StartTime)
	})
	
	return result, nil
}

// GetUpcomingTournaments returns upcoming tournaments
func (s *TournamentService) GetUpcomingTournaments(ctx context.Context, limit int) ([]Tournament, error) {
	var tournaments []models.Tournament
	s.db.Where("status IN (?, ?) AND start_time > ?", "upcoming", "registration", time.Now()).
		Order("start_time ASC").
		Limit(limit).
		Find(&tournaments)
	
	result := make([]Tournament, len(tournaments))
	for i, t := range tournaments {
		result[i] = Tournament{
			ID:              t.ID,
			Name:            t.Name,
			Description:     t.Description,
			Type:            t.Type,
			Status:          t.Status,
			GameFilter:      strings.Split(t.GameFilter, ","),
			MinBet:          t.MinBet,
			StartTime:       t.StartTime,
			EndTime:         t.EndTime,
			RegistrationEnd: t.RegistrationEnd,
			PrizePool:       t.PrizePool,
			Currency:        t.Currency,
			CreatedAt:       t.CreatedAt,
		}
	}
	
	return result, nil
}

// GetUserTournaments returns tournaments a user is participating in
func (s *TournamentService) GetUserTournaments(ctx context.Context, userID uuid.UUID) ([]Tournament, error) {
	var participantIDs []uuid.UUID
	s.db.Model(&models.TournamentParticipant{}).
		Where("user_id = ?", userID).
		Pluck("tournament_id", &participantIDs)
	
	if len(participantIDs) == 0 {
		return []Tournament{}, nil
	}
	
	var tournaments []models.Tournament
	s.db.Where("id IN ?", participantIDs).Find(&tournaments)
	
	result := make([]Tournament, len(tournaments))
	for i, t := range tournaments {
		result[i] = Tournament{
			ID:              t.ID,
			Name:            t.Name,
			Description:     t.Description,
			Type:            t.Type,
			Status:          t.Status,
			PrizePool:       t.PrizePool,
			StartTime:       t.StartTime,
			EndTime:         t.EndTime,
		}
	}
	
	return result, nil
}

// GetTournamentDetails returns detailed tournament information
func (s *TournamentService) GetTournamentDetails(ctx context.Context, tournamentID uuid.UUID) (*Tournament, []LeaderboardEntry, error) {
	s.tournamentMu.RLock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.tournamentMu.RUnlock()
	
	if !ok {
		// Try database
		var t models.Tournament
		if err := s.db.First(&t, tournamentID).Error; err != nil {
			return nil, nil, err
		}
		
		tournament = &Tournament{
			ID:              t.ID,
			Name:            t.Name,
			Description:     t.Description,
			Type:            t.Type,
			Status:          t.Status,
			GameFilter:      strings.Split(t.GameFilter, ","),
			MinBet:          t.MinBet,
			StartTime:       t.StartTime,
			EndTime:         t.EndTime,
			RegistrationEnd: t.RegistrationEnd,
			PrizePool:       t.PrizePool,
			Currency:        t.Currency,
			CreatedAt:       t.CreatedAt,
		}
	}
	
	// Get leaderboard
	leaderboard, err := s.GetLeaderboard(ctx, tournamentID, 100)
	if err != nil {
		return tournament, nil, err
	}
	
	return tournament, leaderboard, nil
}

// ============ TOURNAMENT RESULTS ============

// GetUserTournamentResults returns user's tournament results
func (s *TournamentService) GetUserTournamentResults(ctx context.Context, userID uuid.UUID) ([]TournamentResult, error) {
	var prizes []models.TournamentPrize
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&prizes)
	
	results := make([]TournamentResult, len(prizes))
	
	for i, p := range prizes {
		var tournament models.Tournament
		s.db.First(&tournament, p.TournamentID)
		
		results[i] = TournamentResult{
			TournamentID:   p.TournamentID,
			TournamentName: tournament.Name,
			Position:      p.Position,
			PrizeAmount:   p.PrizeAmount,
			Currency:      p.Currency,
			Date:          p.CreatedAt,
		}
	}
	
	return results, nil
}

type TournamentResult struct {
	TournamentID   uuid.UUID
	TournamentName string
	Position       int
	PrizeAmount    float64
	Currency       string
	Date           time.Time
}

// ============ AUTO-MANAGEMENT ============

// StartTournamentScheduler starts the tournament scheduler
func (s *TournamentService) StartTournamentScheduler(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.config.TournamentCheckInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.processTournaments(ctx)
			}
		}
	}()
}

func (s *TournamentService) processTournaments(ctx context.Context) {
	s.tournamentMu.Lock()
	defer s.tournamentMu.Unlock()
	
	now := time.Now()
	
	for _, tournament := range s.activeTournaments {
		// Start registration
		if tournament.Status == "upcoming" && now.After(tournament.StartTime.Add(-tournament.RegistrationEndDuration)) {
			tournament.Status = "registration"
			s.db.Model(&models.Tournament{}).Where("id = ?", tournament.ID).Update("status", "registration")
		}
		
		// Start tournament
		if tournament.Status == "registration" && now.After(tournament.StartTime) {
			tournament.Status = "active"
			s.db.Model(&models.Tournament{}).Where("id = ?", tournament.ID).Update("status", "active")
		}
		
		// End tournament
		if tournament.Status == "active" && now.After(tournament.EndTime) {
			tournament.Status = "completed"
			s.db.Model(&models.Tournament{}).Where("id = ?", tournament.ID).Update("status", "completed")
			
			// Distribute prizes
			s.distributePrizesFromDB(ctx, tournament.ID)
		}
	}
}

func (s *TournamentService) distributePrizesFromDB(ctx context.Context, tournamentID uuid.UUID) {
	var tournament models.Tournament
	if err := s.db.First(&tournament, tournamentID).Error; err != nil {
		return
	}
	
	var participants []models.TournamentParticipant
	s.db.Where("tournament_id = ?", tournamentID).
		Order("score DESC").
		Limit(10).
		Find(&participants)
	
	if len(participants) < s.config.MinParticipants {
		return
	}
	
	distribution := s.calculatePrizeDistribution(tournament.PrizePool)
	
	for i, p := range participants {
		if i >= len(distribution) {
			break
		}
		
		prize := distribution[i]
		
		var wallet models.Wallet
		if err := s.db.Where("user_id = ? AND currency = ?", p.UserID, tournament.Currency).First(&wallet).Error; err != nil {
			continue
		}
		
		wallet.Balance += prize.Amount
		s.db.Save(&wallet)
		
		prizeRecord := models.TournamentPrize{
			ID:           uuid.New(),
			TournamentID: tournamentID,
			UserID:       p.UserID,
			Position:     prize.Position,
			PrizeAmount: prize.Amount,
			Currency:    tournament.Currency,
			CreatedAt:   time.Now(),
		}
		s.db.Create(&prizeRecord)
	}
}

// LoadTournaments loads active tournaments from database
func (s *TournamentService) LoadTournaments(ctx context.Context) error {
	var tournaments []models.Tournament
	s.db.Where("status IN (?, ?, ?)", "upcoming", "registration", "active").Find(&tournaments)
	
	s.tournamentMu.Lock()
	defer s.tournamentMu.Unlock()
	
	for _, t := range tournaments {
		participants := make(map[uuid.UUID]bool)
		
		var participantList []models.TournamentParticipant
		s.db.Where("tournament_id = ?", t.ID).Find(&participantList)
		
		for _, p := range participantList {
			participants[p.UserID] = true
		}
		
		tournament := &Tournament{
			ID:                  t.ID,
			Name:                t.Name,
			Description:         t.Description,
			Type:                t.Type,
			Status:              t.Status,
			GameFilter:          strings.Split(t.GameFilter, ","),
			MinBet:              t.MinBet,
			StartTime:           t.StartTime,
			EndTime:             t.EndTime,
			RegistrationEnd:    t.RegistrationEnd,
			PrizePool:           t.PrizePool,
			Currency:            t.Currency,
			ScoringType:         t.ScoringType,
			PointsMultiplier:   t.PointsMultiplier,
			MinWagerToQualify:  t.MinWagerToQualify,
			RegisteredUsers:     participants,
			ParticipantCount:    len(participants),
			CreatedAt:           t.CreatedAt,
		}
		
		s.activeTournaments[t.ID] = tournament
	}
	
	return nil
}
