package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// VIPService handles VIP levels, rakeback, and loyalty rewards
type VIPService struct {
	db *gorm.DB
}

func NewVIPService(db *gorm.DB) *VIPService {
	return &VIPService{db: db}
}

// ============ VIP LEVELS ============

type VIPLevelInfo struct {
	Level                    int     `json:"level"`
	Name                     string  `json:"name"`
	RequiredPoints          float64 `json:"required_points"`
	DepositBonusPercent     float64 `json:"deposit_bonus_percent"`
	WithdrawalBonusPercent float64 `json:"withdrawal_bonus_percent"`
	RakebackPercent         float64 `json:"rakeback_percent"`
	MaxDailyWithdrawal     float64 `json:"max_daily_withdrawal"`
	PrioritySupport        bool    `json:"priority_support"`
	MonthlyBonus          float64 `json:"monthly_bonus"`
	WeeklyCashback         float64 `json:"weekly_cashback"`
}

func (s *VIPService) GetAllVIPLevels() ([]VIPLevelInfo, error) {
	var levels []models.VIPLevel
	err := s.db.Order("level ASC").Find(&levels).Error
	if err != nil {
		return nil, err
	}

	result := make([]VIPLevelInfo, len(levels))
	for i, l := range levels {
		result[i] = VIPLevelInfo{
			Level:                    l.Level,
			Name:                     l.Name,
			RequiredPoints:          l.RequiredPoints,
			DepositBonusPercent:     l.DepositBonusPercent,
			WithdrawalBonusPercent: l.WithdrawalBonusPercent,
			RakebackPercent:         l.RakebackPercent,
			MaxDailyWithdrawal:     l.MaxDailyWithdrawal,
			PrioritySupport:        l.PrioritySupport,
		}
	}

	return result, nil
}

func (s *VIPService) GetUserVIPLevel(userID uuid.UUID) (*VIPLevelInfo, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var level models.VIPLevel
	err := s.db.Where("level = ?", user.VIPLevel).First(&level).Error
	if err != nil {
		return nil, err
	}

	return &VIPLevelInfo{
		Level:                    level.Level,
		Name:                     level.Name,
		RequiredPoints:          level.RequiredPoints,
		DepositBonusPercent:     level.DepositBonusPercent,
		WithdrawalBonusPercent: level.WithdrawalBonusPercent,
		RakebackPercent:         level.RakebackPercent,
		MaxDailyWithdrawal:     level.MaxDailyWithdrawal,
		PrioritySupport:        level.PrioritySupport,
	}, nil
}

// ============ RAKEBACK ============

// RakebackEarned tracks rakeback that users have earned
type RakebackEarned struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount      float64  `gorm:"not null"`
	PeriodStart time.Time
	PeriodEnd   time.Time
	Claimed     bool     `gorm:"default:false"`
	ClaimedAt   *time.Time
	CreatedAt   time.Time
}

// CalculateRakeback calculates the rakeback a user has earned based on their wagers
func (s *VIPService) CalculateRakeback(userID uuid.UUID) (float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	// Get user's VIP level rakeback percentage
	var vipLevel models.VIPLevel
	err := s.db.Where("level = ?", user.VIPLevel).First(&vipLevel).Error
	if err != nil {
		return 0, err
	}

	rakebackPercent := vipLevel.RakebackPercent / 100.0

	// Calculate wagers since last rakeback calculation
	var totalWagered float64
	weekAgo := time.Now().AddDate(0, 0, -7)

	// Sum up bets from the last week
	rows, err := s.db.Table("bets").
		Where("user_id = ? AND created_at >= ? AND status IN ('won', 'lost')", userID, weekAgo).
		Select("COALESCE(SUM(bet_amount), 0)").Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&totalWagered)
	}

	// Calculate rakeback (typically 0.1% to 0.5% of wagers)
	rakeback := totalWagered * rakebackPercent

	return rakeback, nil
}

// ClaimRakeback claims the user's accumulated rakeback
func (s *VIPService) ClaimRakeback(userID uuid.UUID) (float64, error) {
	rakeback, err := s.CalculateRakeback(userID)
	if err != nil {
		return 0, err
	}

	if rakeback <= 0 {
		return 0, nil
	}

	// Create rakeback record
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)

	record := RakebackEarned{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      rakeback,
		PeriodStart: weekAgo,
		PeriodEnd:   now,
		Claimed:     true,
		ClaimedAt:   &now,
	}

	if err := s.db.Create(&record).Error; err != nil {
		return 0, err
	}

	// Credit to user balance
	userService := NewUserService(s.db)
	if err := userService.UpdateBalance(userID, rakeback); err != nil {
		return 0, err
	}

	// Create transaction record
	tx := models.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        "rakeback",
		Amount:      rakeback,
		Currency:    "USD",
		Status:      "confirmed",
		Description: "Weekly rakeback bonus",
	}
	s.db.Create(&tx)

	return rakeback, nil
}

// GetUnclaimedRakeback returns the amount of unclaimed rakeback
func (s *VIPService) GetUnclaimedRakeback(userID uuid.UUID) (float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	// Get VIP level
	var vipLevel models.VIPLevel
	err := s.db.Where("level = ?", user.VIPLevel).First(&vipLevel).Error
	if err != nil {
		return 0, err
	}

	rakebackPercent := vipLevel.RakebackPercent / 100.0

	// Calculate unclaimed wagered amount
	var unclaimedWagered float64
	weekAgo := time.Now().AddDate(0, 0, -7)

	s.db.Table("bets").
		Where("user_id = ? AND created_at >= ? AND status IN ('won', 'lost')", userID, weekAgo).
		Select("COALESCE(SUM(bet_amount), 0)").Scan(&unclaimedWagered)

	return unclaimedWagered * rakebackPercent, nil
}

// ============ DEPOSIT BONUSES ============

type DepositBonus struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount        float64   `gorm:"not null"`
	BonusAmount  float64   `gorm:"not null"`
	WagerRequired float64  `gorm:"not null"`
	Wagered      float64   `gorm:"default:0"`
	Status       string    `gorm:"default:'active'"` // active, completed, expired
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// CalculateDepositBonus calculates bonus for a deposit based on VIP level
func (s *VIPService) CalculateDepositBonus(userID uuid.UUID, depositAmount float64) (float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	var vipLevel models.VIPLevel
	err := s.db.Where("level = ?", user.VIPLevel).First(&vipLevel).Error
	if err != nil {
		return 0, err
	}

	bonusAmount := depositAmount * (vipLevel.DepositBonusPercent / 100.0)
	return bonusAmount, nil
}

// CreateDepositBonus creates a new deposit bonus for a user
func (s *VIPService) CreateDepositBonus(userID uuid.UUID, depositAmount, bonusAmount float64) (*DepositBonus, error) {
	// Calculate wager requirement (typically 3x the bonus amount)
	wagerRequired := bonusAmount * 3.0
	expiresAt := time.Now().AddDate(0, 1, 0) // Expires in 1 month

	bonus := DepositBonus{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        depositAmount,
		BonusAmount:  bonusAmount,
		WagerRequired: wagerRequired,
		ExpiresAt:    expiresAt,
	}

	if err := s.db.Create(&bonus).Error; err != nil {
		return nil, err
	}

	// Credit the bonus
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, bonusAmount)

	return &bonus, nil
}

// UpdateWagerProgress updates the wager progress for active deposit bonuses
func (s *VIPService) UpdateWagerProgress(userID uuid.UUID, wagerAmount float64) error {
	var bonuses []DepositBonus
	s.db.Where("user_id = ? AND status = ?", userID, "active").Find(&bonuses)

	for i := range bonuses {
		bonuses[i].Wagered += wagerAmount
		if bonuses[i].Wagered >= bonuses[i].WagerRequired {
			bonuses[i].Status = "completed"
		}
		s.db.Save(&bonuses[i])
	}

	return nil
}

// GetActiveBonuses returns all active bonuses for a user
func (s *VIPService) GetActiveBonuses(userID uuid.UUID) ([]DepositBonus, error) {
	var bonuses []DepositBonus
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").Find(&bonuses).Error
	return bonuses, err
}

// ============ WEEKLY/MONTHLY CASHBACK ============

type CashbackConfig struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	Level           int       `gorm:"uniqueIndex;not null"`
	Name            string    `gorm:"not null"`
	CashbackPercent float64   `gorm:"not null"`
	MinWagerRequired float64 `gorm:"default:0"`
	MaxCashback     float64   `gorm:"default:0"` // 0 = unlimited
	CreatedAt       time.Time
}

type UserCashback struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount      float64   `gorm:"not null"`
	Period      string    `gorm:"not null"` // weekly, monthly
	PeriodStart time.Time
	PeriodEnd   time.Time
	Claimed     bool      `gorm:"default:false"`
	ClaimedAt   *time.Time
	CreatedAt   time.Time
}

// CalculateWeeklyCashback calculates cashback for the week
func (s *VIPService) CalculateWeeklyCashback(userID uuid.UUID) (float64, error) {
	weekAgo := time.Now().AddDate(0, 0, -7)

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	// Get cashback config for user's level
	var config CashbackConfig
	err := s.db.Where("level = ?", user.VIPLevel).First(&config).Error
	if err != nil {
		// Default cashback
		config.CashbackPercent = 0.1 // 0.1%
	}

	// Calculate net loss (wagers - wins) for the week
	var totalWagered, totalWon float64

	s.db.Table("bets").
		Where("user_id = ? AND created_at >= ?", userID, weekAgo).
		Select("COALESCE(SUM(bet_amount), 0)").Scan(&totalWagered)

	s.db.Table("bets").
		Where("user_id = ? AND created_at >= ? AND status = ?", userID, weekAgo, "won").
		Select("COALESCE(SUM(win_amount), 0)").Scan(&totalWon)

	netLoss := totalWagered - totalWon
	if netLoss < 0 || config.CashbackPercent == 0 {
		return 0, nil
	}

	cashback := netLoss * (config.CashbackPercent / 100.0)

	// Apply max cap if set
	if config.MaxCashback > 0 && cashback > config.MaxCashback {
		cashback = config.MaxCashback
	}

	return cashback, nil
}

// ClaimWeeklyCashback claims weekly cashback
func (s *VIPService) ClaimWeeklyCashback(userID uuid.UUID) (float64, error) {
	cashback, err := s.CalculateWeeklyCashback(userID)
	if err != nil || cashback <= 0 {
		return cashback, err
	}

	// Check if already claimed
	weekAgo := time.Now().AddDate(0, 0, -7)
	var existing UserCashback
	err = s.db.Where("user_id = ? AND period = ? AND period_start >= ?", userID, "weekly", weekAgo).First(&existing).Error
	if err == nil && existing.Claimed {
		return 0, fmt.Errorf("cashback already claimed for this period")
	}

	// Create cashback record
	now := time.Now()
	userCashback := UserCashback{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      cashback,
		Period:      "weekly",
		PeriodStart: weekAgo,
		PeriodEnd:   now,
		Claimed:     true,
		ClaimedAt:   &now,
	}

	s.db.Create(&userCashback)

	// Credit user
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, cashback)

	// Create transaction
	tx := models.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        "cashback",
		Amount:      cashback,
		Currency:    "USD",
		Status:      "confirmed",
		Description: "Weekly cashback",
	}
	s.db.Create(&tx)

	return cashback, nil
}

// ============ VIP POINT TRACKING ============

// VIPPointsEvent tracks VIP point earning events
type VIPPointsEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Points    float64   `gorm:"not null"`
	EventType string    `gorm:"not null"` // wager, deposit, bonus, referral
	EventID   string    // Reference to bet/deposit ID
	CreatedAt time.Time
}

// AddVIPPoints adds VIP points for a user
func (s *VIPService) AddVIPPoints(userID uuid.UUID, points float64, eventType string, eventID string) error {
	// Create points event
	event := VIPPointsEvent{
		ID:        uuid.New(),
		UserID:    userID,
		Points:    points,
		EventType: eventType,
		EventID:   eventID,
	}

	if err := s.db.Create(&event).Error; err != nil {
		return err
	}

	// Update user's VIP points
	userService := NewUserService(s.db)
	var user models.User
	s.db.First(&user, userID)
	newPoints := user.VIPPoints + points
	s.db.Model(&user).Update("vip_points", newPoints)

	// Check for level upgrade
	s.CheckLevelUpgrade(userID)

	return nil
}

// CheckLevelUpgrade checks if user qualifies for a level upgrade
func (s *VIPService) CheckLevelUpgrade(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Get all VIP levels ordered by level
	var levels []models.VIPLevel
	s.db.Order("level DESC").Find(&levels)

	for _, level := range levels {
		if user.VIPPoints >= level.RequiredPoints && user.VIPLevel < level.Level {
			// Upgrade user
			s.db.Model(&user).Update("vip_level", level.Level)

			// Create notification
			// In production, this would trigger a notification system
			return nil
		}
	}

	return nil
}

// GetVIPProgress returns the user's progress towards the next VIP level
func (s *VIPService) GetVIPProgress(userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var currentLevel, nextLevel models.VIPLevel
	s.db.Where("level = ?", user.VIPLevel).First(&currentLevel)

	// Get next level
	s.db.Where("level = ?", user.VIPLevel+1).First(&nextLevel)

	progress := map[string]interface{}{
		"current_level":     currentLevel.Level,
		"current_level_name": currentLevel.Name,
		"current_points":     user.VIPPoints,
		"rakeback_percent":   currentLevel.RakebackPercent,
	}

	if nextLevel.Level > 0 {
		pointsNeeded := nextLevel.RequiredPoints - user.VIPPoints
		progress["next_level"] = nextLevel.Level
		progress["next_level_name"] = nextLevel.Name
		progress["points_needed"] = pointsNeeded
		progress["progress_percent"] = (user.VIPPoints / nextLevel.RequiredPoints) * 100
	} else {
		progress["next_level"] = nil
		progress["is_max_level"] = true
	}

	return progress, nil
}

// ============ REFERRAL SYSTEM ============

type ReferralBonus struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	ReferrerID      uuid.UUID `gorm:"type:uuid;not null;index"`
	RefereeID      uuid.UUID `gorm:"type:uuid;not null;index"`
	CommissionPercent float64 `gorm:"not null"`
	TotalEarned    float64   `gorm:"default:0"`
	RefereeDeposited float64 `gorm:"default:0"`
	Status         string    `gorm:"default:'active'"` // active, completed
	CreatedAt      time.Time
}

// GetReferralBonus calculates referral bonus for a user
func (s *VIPService) GetReferralBonus(referrerID uuid.UUID, refereeDeposit float64) (float64, error) {
	// Default referral commission is 10%
	commissionPercent := 10.0

	// Check if referrer has elevated referral bonus from VIP level
	var referrer models.User
	if err := s.db.First(&referrer, referrerID).Error; err != nil {
		return 0, err
	}

	// VIP Diamond+ might have higher commission
	if referrer.VIPLevel >= 4 {
		commissionPercent = 15.0
	}

	bonus := refereeDeposit * (commissionPercent / 100.0)
	return bonus, nil
}

// CreateReferralBonus creates a referral bonus record
func (s *VIPService) CreateReferralBonus(referrerID, refereeID uuid.UUID, deposit float64) (*ReferralBonus, error) {
	bonus, err := s.GetReferralBonus(referrerID, deposit)
	if err != nil {
		return nil, err
	}

	referral := ReferralBonus{
		ID:              uuid.New(),
		ReferrerID:      referrerID,
		RefereeID:      refereeID,
		CommissionPercent: bonus / deposit * 100,
		TotalEarned:    bonus,
		RefereeDeposited: deposit,
	}

	if err := s.db.Create(&referral).Error; err != nil {
		return nil, err
	}

	// Credit referrer
	userService := NewUserService(s.db)
	userService.UpdateBalance(referrerID, bonus)

	// Add VIP points to both
	s.AddVIPPoints(referrerID, bonus, "referral", referral.ID.String())
	s.AddVIPPoints(refereeID, bonus/2, "referral", referral.ID.String())

	return &referral, nil
}

// GetReferralStats returns referral statistics for a user
func (s *VIPService) GetReferralStats(userID uuid.UUID) (map[string]interface{}, error) {
	var referrals []ReferralBonus
	s.db.Where("referrer_id = ?", userID).Find(&referrals)

	totalEarned := 0.0
	activeReferees := 0

	for _, r := range referrals {
		totalEarned += r.TotalEarned
		if r.Status == "active" {
			activeReferees++
		}
	}

	return map[string]interface{}{
		"total_referrals":   len(referrals),
		"active_referees": activeReferees,
		"total_earned":    totalEarned,
		"referrals":       referrals,
	}, nil
}
