package services

import (
	"context"
	"fmt"
	"math"
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

// VIPService handles VIP programs, rakeback, and loyalty rewards
type VIPService struct {
	db           *gorm.DB
	redis        *redis.Client
	config       *VIPConfig
	levels       map[int]*VIPLevel
	bonusCache   *BonusCache
}

type VIPConfig struct {
	// Rakeback settings
	BaseRakeback     float64   // Base rakeback percentage (e.g., 3.5)
	MaxRakeback      float64   // Maximum rakeback (e.g., 25.0)
	RakeContribution float64   // How much rake = 1 point (e.g., 1 USD)
	
	// Level thresholds (total wagered)
	LevelThresholds []float64
	
	// Level benefits
	LevelBenefits map[int]*LevelBenefits
	
	// Bonus settings
	WelcomeBonusAmount    float64
	WelcomeBonusWager    float64  // Multiplier (e.g., 10x)
	DepositBonusPercent   float64
	DepositBonusMax      float64
	
	// Cashback settings
	DailyCashbackPercent float64
	WeeklyCashbackPercent float64
	MonthlyCashbackPercent float64
	
	// Points settings
	PointsPerWager      float64   // Points earned per $1 wagered
	PointsRedemptionRate float64   // $ value per 100 points
	
	// Time settings
	PromotionCheckInterval time.Duration
	LevelUpdateInterval   time.Duration
}

type VIPLevel struct {
	Level             int
	Name              string
	MinWagered        float64
	RakebackPercent   float64
	WeeklyCashback    float64
	MonthlyCashback   float64
	MaxWithdrawal     float64
	WithdrawalFee     float64
	PersonalHost      bool
	ExclusiveEvents   bool
	CustomLimits      bool
}

type LevelBenefits struct {
	MaxBet           float64
	MaxWin           float64
	WithdrawalLimit  float64
	WithdrawalFee    float64
	DailyWithdrawalLimit int
	MonthlyWithdrawalLimit int
	RakebackPercent  float64
	CashbackPercent  float64
	PointsMultiplier float64
	InstantWithdraw  bool
	PrioritySupport  bool
	PersonalHost     bool
	ExclusiveBonuses bool
	TournamentAccess bool
	CustomPromoCodes bool
}

type BonusCache struct {
	mu      sync.RWMutex
	bonuses map[string]*ActiveBonus // userID + bonusID -> bonus
}

type ActiveBonus struct {
	UserID      uuid.UUID
	BonusID     uuid.UUID
	BonusType   string // welcome, deposit, cashback, free_spins, loyalty
	Amount      float64
	Wagered     float64
	WagerReq    float64
	ExpiresAt   time.Time
	Status      string // active, completed, expired, cancelled
	CreatedAt   time.Time
}

// DefaultVIPConfig returns default VIP configuration
func DefaultVIPConfig() *VIPConfig {
	return &VIPConfig{
		BaseRakeback:     3.5,
		MaxRakeback:      25.0,
		RakeContribution: 1.0, // $1 rake = 1 point
		
		LevelThresholds: []float64{
			0,      // Bronze
			1000,   // Silver
			5000,   // Gold
			25000,  // Platinum
			100000, // Diamond
			500000, // VIP
		},
		
		LevelBenefits: map[int]*LevelBenefits{
			0: { // Bronze
				MaxBet:            1000,
				MaxWin:            5000,
				WithdrawalLimit:   10000,
				WithdrawalFee:     2.0,
				DailyWithdrawalLimit: 3,
				RakebackPercent:   3.5,
				CashbackPercent:   0,
				PointsMultiplier:  1.0,
			},
			1: { // Silver
				MaxBet:            2500,
				MaxWin:            12500,
				WithdrawalLimit:   25000,
				WithdrawalFee:     1.5,
				DailyWithdrawalLimit: 5,
				RakebackPercent:   5.0,
				CashbackPercent:   2.0,
				PointsMultiplier:  1.25,
			},
			2: { // Gold
				MaxBet:            5000,
				MaxWin:            25000,
				WithdrawalLimit:   50000,
				WithdrawalFee:     1.0,
				DailyWithdrawalLimit: 10,
				RakebackPercent:   7.5,
				CashbackPercent:   3.0,
				PointsMultiplier:  1.5,
				PrioritySupport:   true,
			},
			3: { // Platinum
				MaxBet:            10000,
				MaxWin:            50000,
				WithdrawalLimit:   100000,
				WithdrawalFee:     0.5,
				DailyWithdrawalLimit: 20,
				RakebackPercent:   10.0,
				CashbackPercent:   5.0,
				PointsMultiplier:  2.0,
				PrioritySupport:   true,
				InstantWithdraw:   true,
			},
			4: { // Diamond
				MaxBet:            25000,
				MaxWin:            125000,
				WithdrawalLimit:   250000,
				WithdrawalFee:     0,
				DailyWithdrawalLimit: 50,
				RakebackPercent:   15.0,
				CashbackPercent:   7.5,
				PointsMultiplier:  2.5,
				PrioritySupport:   true,
				InstantWithdraw:   true,
				ExclusiveBonuses:  true,
				TournamentAccess: true,
			},
			5: { // VIP
				MaxBet:            100000,
				MaxWin:            500000,
				WithdrawalLimit:   1000000,
				WithdrawalFee:     0,
				DailyWithdrawalLimit: 100,
				RakebackPercent:   25.0,
				CashbackPercent:   10.0,
				PointsMultiplier:  3.0,
				PrioritySupport:   true,
				InstantWithdraw:   true,
				ExclusiveBonuses:  true,
				TournamentAccess: true,
				PersonalHost:      true,
				CustomPromoCodes: true,
			},
		},
		
		WelcomeBonusAmount:   100,
		WelcomeBonusWager:    10,
		DepositBonusPercent:   100,
		DepositBonusMax:      1000,
		
		DailyCashbackPercent:   0,
		WeeklyCashbackPercent:   5.0,
		MonthlyCashbackPercent: 10.0,
		
		PointsPerWager:       1.0,
		PointsRedemptionRate: 0.01, // $0.01 per point
		
		PromotionCheckInterval: time.Hour,
		LevelUpdateInterval:   24 * time.Hour,
	}
}

func NewVIPLevel(level int) *VIPLevel {
	names := []string{"Bronze", "Silver", "Gold", "Platinum", "Diamond", "VIP"}
	name := names[level]
	if level >= len(names) {
		name = "VIP"
	}
	
	return &VIPLevel{
		Level:           level,
		Name:            name,
		MinWagered:      float64(level) * 100000,
		RakebackPercent: 3.5 + float64(level)*2.5,
		WeeklyCashback:  float64(level) * 1.0,
		MonthlyCashback: float64(level) * 2.0,
	}
}

// NewVIPService creates a new VIP service
func NewVIPService(db *gorm.DB, redisClient *redis.Client, config *VIPConfig) *VIPService {
	if config == nil {
		config = DefaultVIPConfig()
	}
	
	service := &VIPService{
		db:         db,
		redis:      redisClient,
		config:     config,
		bonusCache: &BonusCache{bonuses: make(map[string]*ActiveBonus)},
	}
	
	// Initialize VIP levels
	service.levels = make(map[int]*VIPLevel)
	for i := 0; i <= 5; i++ {
		service.levels[i] = NewVIPLevel(i)
	}
	
	return service
}

// ============ USER VIP STATUS ============

// GetUserVIPStatus returns user's current VIP status
func (s *VIPService) GetUserVIPStatus(ctx context.Context, userID uuid.UUID) (*VIPUserStatus, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	
	// Calculate total wagered
	var totalWagered float64
	s.db.Model(&models.GameHistory{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(bet_amount), 0)").
		Scan(&totalWagered)
	
	// Get current level
	currentLevel := s.calculateLevel(totalWagered)
	level := s.levels[currentLevel]
	
	// Get benefits
	benefits := s.config.LevelBenefits[currentLevel]
	
	// Get points
	var points int64
	s.db.Model(&models.LoyaltyPoints{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points), 0)").
		Scan(&points)
	
	// Get next level info
	nextLevel := currentLevel + 1
	var nextLevelInfo *VIPLevel
	var progressToNext float64
	
	if nextLevel <= 5 {
		nextLevelInfo = s.levels[nextLevel]
		if currentLevel < 5 {
			threshold := s.config.LevelThresholds[currentLevel]
			nextThreshold := s.config.LevelThresholds[nextLevel]
			progressToNext = (totalWagered - threshold) / (nextThreshold - threshold) * 100
		}
	}
	
	// Get rakeback balance
	var rakebackBalance float64
	s.db.Model(&models.RakebackBalance{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(available), 0)").
		Scan(&rakebackBalance)
	
	return &VIPUserStatus{
		UserID:          userID,
		Level:           currentLevel,
		LevelName:       level.Name,
		TotalWagered:    totalWagered,
		Points:          points,
		RakebackPercent: level.RakebackPercent,
		RakebackBalance: rakebackBalance,
		Benefits:        benefits,
		NextLevel:       nextLevelInfo,
		ProgressToNext:  math.Min(progressToNext, 100),
	}, nil
}

type VIPUserStatus struct {
	UserID          uuid.UUID
	Level           int
	LevelName       string
	TotalWagered    float64
	Points          int64
	RakebackPercent float64
	RakebackBalance float64
	Benefits        *LevelBenefits
	NextLevel       *VIPLevel
	ProgressToNext  float64
}

func (s *VIPService) calculateLevel(totalWagered float64) int {
	for i := len(s.config.LevelThresholds) - 1; i >= 0; i-- {
		if totalWagered >= s.config.LevelThresholds[i] {
			return i
		}
	}
	return 0
}

// ============ RAKEBACK SYSTEM ============

// CalculateRake calculates rake from a bet
func (s *VIPService) CalculateRake(betAmount float64, gameType string) float64 {
	// Different game types have different rake percentages
	rakePercent := map[string]float64{
		"slots":       0.05,  // 5% rake
		"table_games": 0.02,  // 2% rake
		"live_casino": 0.03,  // 3% rake
		"sports":      0.04,  // 4% rake
		"poker":       0.05,  // 5% rake
	}
	
	percent := rakePercent[gameType]
	if percent == 0 {
		percent = 0.05 // Default 5%
	}
	
	return betAmount * percent
}

// ProcessRakeback processes rakeback for a user after a bet
func (s *VIPService) ProcessRakeback(ctx context.Context, userID uuid.UUID, betAmount float64, gameType string) error {
	rake := s.CalculateRake(betAmount, gameType)
	if rake <= 0 {
		return nil
	}
	
	// Get user's VIP level
	status, err := s.GetUserVIPStatus(ctx, userID)
	if err != nil {
		return err
	}
	
	// Calculate rakeback based on level
	rakebackAmount := rake * (status.RakebackPercent / 100)
	
	// Create rakeback balance record
	rakeback := models.RakebackBalance{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      rakebackAmount,
		Available:   rakebackAmount,
		Locked:     0,
		Source:     gameType,
		ExternalRef: uuid.New().String(),
		CreatedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour), // 30 days expiry
	}
	
	if err := s.db.Create(&rakeback).Error; err != nil {
		return err
	}
	
	// Award loyalty points
	points := int64(rake * s.config.PointsPerWager * status.Benefits.PointsMultiplier)
	if points > 0 {
		s.AwardPoints(ctx, userID, points, "rakeback")
	}
	
	return nil
}

// ClaimRakeback allows user to claim their available rakeback
func (s *VIPService) ClaimRakeback(ctx context.Context, userID uuid.UUID) (float64, error) {
	var available float64
	s.db.Model(&models.RakebackBalance{}).
		Where("user_id = ? AND available > 0 AND expires_at > ?", userID, time.Now()).
		Select("COALESCE(SUM(available), 0)").
		Scan(&available)
	
	if available <= 0 {
		return 0, fmt.Errorf("no rakeback available to claim")
	}
	
	// Credit user wallet
	var wallet models.Wallet
	if err := s.db.Where("user_id = ? AND currency = ?", userID, "USD").First(&wallet).Error; err != nil {
		return 0, err
	}
	
	tx := s.db.Begin()
	
	wallet.Balance += available
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	// Mark rakeback as claimed
	if err := tx.Model(&models.RakebackBalance{}).
		Where("user_id = ? AND available > 0", userID).
		Update("available", 0).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	tx.Commit()
	
	// Create transaction record
	transaction := models.Transaction{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      "rakeback",
		Amount:    available,
		Status:    "completed",
		CreatedAt: time.Now(),
	}
	s.db.Create(&transaction)
	
	return available, nil
}

// ============ LOYALTY POINTS ============

// AwardPoints awards loyalty points to a user
func (s *VIPService) AwardPoints(ctx context.Context, userID uuid.UUID, points int64, source string) error {
	pointsRecord := models.LoyaltyPoints{
		ID:        uuid.New(),
		UserID:    userID,
		Points:    points,
		Source:    source,
		CreatedAt: time.Now(),
	}
	
	return s.db.Create(&pointsRecord).Error
}

// RedeemPoints allows user to redeem points for bonus money
func (s *VIPService) RedeemPoints(ctx context.Context, userID uuid.UUID, points int64) (float64, error) {
	// Check available points
	var availablePoints int64
	s.db.Model(&models.LoyaltyPoints{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points), 0)").
		Scan(&availablePoints)
	
	if availablePoints < points {
		return 0, fmt.Errorf("insufficient points")
	}
	
	// Calculate redemption value
	redemptionValue := float64(points) * s.config.PointsRedemptionRate
	
	// Deduct points
	tx := s.db.Begin()
	
	// Create negative points record
	pointsRecord := models.LoyaltyPoints{
		ID:        uuid.New(),
		UserID:    userID,
		Points:    -points,
		Source:    "redemption",
		CreatedAt: time.Now(),
	}
	
	if err := tx.Create(&pointsRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	// Credit bonus to user wallet
	var wallet models.Wallet
	if err := tx.Where("user_id = ? AND currency = ?", userID, "USD").First(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	wallet.Balance += redemptionValue
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	tx.Commit()
	
	return redemptionValue, nil
}

// ============ BONUS SYSTEM ============

// ClaimWelcomeBonus claims the welcome bonus for new users
func (s *VIPService) ClaimWelcomeBonus(ctx context.Context, userID uuid.UUID) (*BonusResult, error) {
	// Check if already claimed
	var existingBonus int64
	s.db.Model(&models.BonusClaim{}).
		Where("user_id = ? AND bonus_type = ?", userID, "welcome").
		Count(&existingBonus)
	
	if existingBonus > 0 {
		return nil, fmt.Errorf("welcome bonus already claimed")
	}
	
	bonusAmount := s.config.WelcomeBonusAmount
	wagerReq := bonusAmount * s.config.WelcomeBonusWager
	
	// Credit bonus
	var wallet models.Wallet
	if err := s.db.Where("user_id = ? AND currency = ?", userID, "USD").First(&wallet).Error; err != nil {
		return nil, err
	}
	
	tx := s.db.Begin()
	
	wallet.Balance += bonusAmount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	// Create bonus tracking record
	bonusClaim := models.BonusClaim{
		ID:            uuid.New(),
		UserID:        userID,
		BonusType:     "welcome",
		Amount:        bonusAmount,
		WagerRequired: wagerReq,
		Wagered:       0,
		Status:        "active",
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(30 * 24 * time.Hour),
	}
	
	if err := tx.Create(&bonusClaim).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	tx.Commit()
	
	return &BonusResult{
		BonusID:      bonusClaim.ID,
		Amount:       bonusAmount,
		WagerReq:     wagerReq,
		ExpiresAt:    bonusClaim.ExpiresAt,
	}, nil
}

type BonusResult struct {
	BonusID     uuid.UUID
	Amount      float64
	WagerReq    float64
	ExpiresAt   time.Time
}

// ClaimDepositBonus claims a deposit bonus
func (s *VIPService) ClaimDepositBonus(ctx context.Context, userID uuid.UUID, depositAmount float64) (*BonusResult, error) {
	bonusAmount := depositAmount * (s.config.DepositBonusPercent / 100)
	if bonusAmount > s.config.DepositBonusMax {
		bonusAmount = s.config.DepositBonusMax
	}
	
	wagerReq := bonusAmount * 5 // 5x wager requirement for deposit bonus
	
	// Credit bonus
	var wallet models.Wallet
	if err := s.db.Where("user_id = ? AND currency = ?", userID, "USD").First(&wallet).Error; err != nil {
		return nil, err
	}
	
	tx := s.db.Begin()
	
	wallet.Balance += bonusAmount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	// Create bonus tracking record
	bonusClaim := models.BonusClaim{
		ID:            uuid.New(),
		UserID:        userID,
		BonusType:     "deposit",
		Amount:        bonusAmount,
		WagerRequired: wagerReq,
		Wagered:       0,
		Status:        "active",
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(7 * 24 * time.Hour),
	}
	
	if err := tx.Create(&bonusClaim).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	tx.Commit()
	
	return &BonusResult{
		BonusID:   bonusClaim.ID,
		Amount:    bonusAmount,
		WagerReq:  wagerReq,
		ExpiresAt: bonusClaim.ExpiresAt,
	}, nil
}

// ProcessWagerProgress updates wager progress for a bonus
func (s *VIPService) ProcessWagerProgress(ctx context.Context, userID uuid.UUID, betAmount float64) error {
	// Find active bonuses
	var activeBonuses []models.BonusClaim
	s.db.Where("user_id = ? AND status = ? AND expires_at > ?", userID, "active", time.Now()).
		Find(&activeBonuses)
	
	if len(activeBonuses) == 0 {
		return nil
	}
	
	tx := s.db.Begin()
	
	for i := range activeBonuses {
		bonus := &activeBonuses[i]
		
		// Add wager progress
		bonus.Wagered += betAmount
		
		// Check if wager requirement met
		if bonus.Wagered >= bonus.WagerRequired {
			bonus.Status = "completed"
			bonus.CompletedAt = time.Now()
		}
		
		if err := tx.Save(bonus).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
}

// ============ CASHBACK SYSTEM ============

// ProcessDailyCashback processes daily cashback for all users
func (s *VIPService) ProcessDailyCashback(ctx context.Context) error {
	if s.config.DailyCashbackPercent <= 0 {
		return nil
	}
	
	return s.processCashback(ctx, "daily")
}

// ProcessWeeklyCashback processes weekly cashback
func (s *VIPService) ProcessWeeklyCashback(ctx context.Context) error {
	if s.config.WeeklyCashbackPercent <= 0 {
		return nil
	}
	
	return s.processCashback(ctx, "weekly")
}

func (s *VIPService) processCashback(ctx context.Context, period string) error {
	// Get all users who made activity in the period
	var users []models.User
	s.db.Find(&users)
	
	percentMap := map[string]float64{
		"daily":   s.config.DailyCashbackPercent,
		"weekly":  s.config.WeeklyCashbackPercent,
		"monthly": s.config.MonthlyCashbackPercent,
	}
	
	percent := percentMap[period]
	if percent <= 0 {
		return nil
	}
	
	for _, user := range users {
		// Calculate net loss for period
		var totalWagered, totalWon float64
		
		since := time.Now().AddDate(0, 0, -1)
		if period == "weekly" {
			since = time.Now().AddDate(0, 0, -7)
		} else if period == "monthly" {
			since = time.Now().AddDate(0, -1, 0)
		}
		
		s.db.Model(&models.GameHistory{}).
			Where("user_id = ? AND timestamp > ?", user.ID, since).
			Select("COALESCE(SUM(bet_amount), 0)", "COALESCE(SUM(win_amount), 0)").
			Scan(&totalWagered, &totalWon)
		
		netLoss := totalWagered - totalWon
		if netLoss <= 0 {
			continue // No loss, no cashback
		}
		
		// Get user's cashback percentage based on level
		status, err := s.GetUserVIPStatus(ctx, user.ID)
		if err != nil {
			continue
		}
		
		cashbackPercent := percent
		if status.Benefits != nil && status.Benefits.CashbackPercent > 0 {
			cashbackPercent = status.Benefits.CashbackPercent
		}
		
		cashbackAmount := netLoss * (cashbackPercent / 100)
		
		// Credit cashback
		var wallet models.Wallet
		if err := s.db.Where("user_id = ? AND currency = ?", user.ID, "USD").First(&wallet).Error; err != nil {
			continue
		}
		
		wallet.Balance += cashbackAmount
		wallet.UpdatedAt = time.Now()
		s.db.Save(&wallet)
		
		// Create cashback record
		cashback := models.CashbackRecord{
			ID:        uuid.New(),
			UserID:    user.ID,
			Period:    period,
			NetLoss:   netLoss,
			Percent:   cashbackPercent,
			Amount:    cashbackAmount,
			CreatedAt: time.Now(),
		}
		s.db.Create(&cashback)
	}
	
	return nil
}

// ============ LEADERBOARD ============

// GetLeaderboard returns top users by wager volume
func (s *VIPService) GetLeaderboard(ctx context.Context, period string, limit int) ([]LeaderboardEntry, error) {
	var since time.Time
	now := time.Now()
	
	switch period {
	case "daily":
		since = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "weekly":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		since = now.AddDate(0, 0, -(weekday - 1))
		since = time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, since.Location())
	case "monthly":
		since = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "all":
		since = time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
	default:
		since = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	
	type leaderboardRow struct {
		UserID      uuid.UUID
		Username    string
		TotalWager  float64
		TotalWins   float64
		BetCount    int64
	}
	
	var rows []leaderboardRow
	
	s.db.Model(&models.GameHistory{}).
		Select("user_id, SUM(bet_amount) as total_wager, SUM(win_amount) as total_wins, COUNT(*) as bet_count").
		Where("timestamp > ?", since).
		Group("user_id").
		Order("total_wager DESC").
		Limit(limit).
		Scan(&rows)
	
	entries := make([]LeaderboardEntry, len(rows))
	
	for i, row := range rows {
		var user models.User
		s.db.First(&user, row.UserID)
		
		entries[i] = LeaderboardEntry{
			Rank:        i + 1,
			UserID:      row.UserID,
			Username:    user.Username,
			TotalWager:  row.TotalWager,
			TotalWins:   row.TotalWins,
			NetProfit:   row.TotalWins - row.TotalWager,
			BetCount:    row.BetCount,
		}
	}
	
	return entries, nil
}

type LeaderboardEntry struct {
	Rank       int
	UserID     uuid.UUID
	Username   string
	TotalWager float64
	TotalWins  float64
	NetProfit  float64
	BetCount   int64
}

// ============ PROMOTIONS ============

// GetActivePromotions returns currently active promotions
func (s *VIPService) GetActivePromotions(ctx context.Context) ([]Promotion, error) {
	var promotions []models.Promotion
	s.db.Where("status = ? AND start_date <= ? AND end_date >= ?", "active", time.Now(), time.Now()).
		Find(&promotions)
	
	result := make([]Promotion, len(promotions))
	for i, p := range promotions {
		result[i] = Promotion{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Type:        p.Type,
			BonusAmount: p.BonusAmount,
			WagerReq:    p.WagerReq,
			StartDate:   p.StartDate,
			EndDate:    p.EndDate,
			Terms:       p.Terms,
		}
	}
	
	return result, nil
}

type Promotion struct {
	ID          uuid.UUID
	Name        string
	Description string
	Type        string
	BonusAmount float64
	WagerReq    float64
	StartDate   time.Time
	EndDate     time.Time
	Terms       string
}

// ============ LEVEL UP BONUS ============

// ProcessLevelUp processes level up rewards
func (s *VIPService) ProcessLevelUp(ctx context.Context, userID uuid.UUID, oldLevel, newLevel int) error {
	if newLevel <= oldLevel {
		return nil
	}
	
	// Calculate level up bonus
	levelBonus := float64(newLevel-oldLevel) * 50 // $50 per level
	
	// Credit bonus
	var wallet models.Wallet
	if err := s.db.Where("user_id = ? AND currency = ?", userID, "USD").First(&wallet).Error; err != nil {
		return err
	}
	
	tx := s.db.Begin()
	
	wallet.Balance += levelBonus
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Create level up record
	levelUp := models.LevelUpReward{
		ID:        uuid.New(),
		UserID:    userID,
		OldLevel:  oldLevel,
		NewLevel:  newLevel,
		Bonus:     levelBonus,
		CreatedAt: time.Now(),
	}
	
	if err := tx.Create(&levelUp).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

// ============ VIP STATUS SUMMARY ============

// GetAllVIPLevels returns all VIP levels with their benefits
func (s *VIPService) GetAllVIPLevels(ctx context.Context) []VIPLevelInfo {
	levels := make([]VIPLevelInfo, 6)
	
	for i := 0; i <= 5; i++ {
		level := s.levels[i]
		benefits := s.config.LevelBenefits[i]
		
		levels[i] = VIPLevelInfo{
			Level:           i,
			Name:             level.Name,
			MinWagered:       s.config.LevelThresholds[i],
			RakebackPercent:  level.RakebackPercent,
			WeeklyCashback:   level.WeeklyCashback,
			MonthlyCashback:  level.MonthlyCashback,
			Benefits:         benefits,
		}
	}
	
	return levels
}

type VIPLevelInfo struct {
	Level           int
	Name            string
	MinWagered      float64
	RakebackPercent float64
	WeeklyCashback  float64
	MonthlyCashback float64
	Benefits        *LevelBenefits
}

// ============ REFERENCES TO EXISTING MODELS ============

// Ensure models exist - these are references to models in models package
var _ = func() {}(
	models.User{},
	models.GameHistory{},
	models.Wallet{},
	models.RakebackBalance{},
	models.LoyaltyPoints{},
	models.BonusClaim{},
	models.CashbackRecord{},
	models.LevelUpReward{},
	models.Promotion{},
)
