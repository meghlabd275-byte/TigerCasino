package services

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// VIPService handles VIP and loyalty program operations
type VIPService struct {
	db *gorm.DB
}

func NewVIPService(db *gorm.DB) *VIPService {
	return &VIPService{db: db}
}

// GetUserVIPStatus returns the user's current VIP status
func (s *VIPService) GetUserVIPStatus(userID uuid.UUID) (*models.VIPLevel, float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, 0, err
	}

	var vipLevel models.VIPLevel
	err := s.db.Where("level = ?", user.VIPLevel).First(&vipLevel).Error
	if err != nil {
		return nil, 0, err
	}

	return &vipLevel, user.VIPPoints, nil
}

// CalculateRakeback calculates the rakeback for a user based on their VIP level
func (s *VIPService) CalculateRakeback(userID uuid.UUID, betAmount float64) (float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	var vipLevel models.VIPLevel
	if err := s.db.Where("level = ?", user.VIPLevel).First(&vipLevel).Error; err != nil {
		return 0, err
	}

	// Calculate rakeback (typically 0.5% - 1% of house edge)
	rakebackPercent := vipLevel.RakebackPercent
	rakeback := betAmount * (rakebackPercent / 100)

	return rakeback, nil
}

// AwardPoints awards VIP points to a user
func (s *VIPService) AwardPoints(userID uuid.UUID, points float64) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("vip_points", gorm.Expr("vip_points + ?", points)).Error
}

// UpdateVIPLevel checks and updates user's VIP level based on points
func (s *VIPService) UpdateVIPLevel(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Get all VIP levels ordered by required points
	var levels []models.VIPLevel
	s.db.Where("is_active = ?", true).Order("required_points DESC").Find(&levels)

	// Find the appropriate level
	newLevel := 0
	for _, level := range levels {
		if user.VIPPoints >= level.RequiredPoints {
			newLevel = level.Level
			break
		}

		if newLevel == 0 {
			newLevel = level.Level
		}
	}

	// Update if changed
	if newLevel != user.VIPLevel {
		return s.db.Model(&models.User{}).
			Where("id = ?", userID).
			Update("vip_level", newLevel).Error
	}

	return nil
}

// GetRakebackHistory returns the user's rakeback history
func (s *VIPService) GetRakebackHistory(userID uuid.UUID, days int) ([]map[string]interface{}, error) {
	var transactions []models.Transaction
	startDate := time.Now().AddDate(0, 0, -days)

	err := s.db.Where("user_id = ? AND type = ? AND created_at > ?", userID, "rakeback", startDate).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, tx := range transactions {
		results = append(results, map[string]interface{}{
			"amount":      tx.Amount,
			"created_at": tx.CreatedAt,
		})
	}

	return results, nil
}

// GetVIPBonuses returns available bonuses for user's VIP level
func (s *VIPService) GetVIPBonuses(userID uuid.UUID) ([]map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Return available bonuses based on VIP level
	bonuses := []map[string]interface{}{
		{
			"type":        "deposit",
			"name":        "Weekly Deposit Bonus",
			"percentage":  float64(user.VIPLevel) * 5,
			"max_amount":  1000.0 * float64(user.VIPLevel+1),
			"wager_req":   10,
		},
		{
			"type":        "cashback",
			"name":        "Weekly Cashback",
			"percentage":  float64(user.VIPLevel) * 2,
			"max_amount":  500.0 * float64(user.VIPLevel+1),
			"wager_req":   0,
		},
	}

	return bonuses, nil
}

// CalculateLevelUpBonus calculates bonus for leveling up
func (s *VIPService) CalculateLevelUpBonus(oldLevel, newLevel int) (float64, error) {
	bonusAmounts := map[int]float64{
		1: 50,
		2: 100,
		3: 250,
		4: 500,
		5: 1000,
	}

	bonus, ok := bonusAmounts[newLevel]
	if !ok {
		return 0, nil
	}

	return bonus, nil
}
