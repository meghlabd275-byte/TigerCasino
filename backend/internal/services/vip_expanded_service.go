package services

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// VIPService handles the expanded VIP program (50+ levels)
type VIPService struct {
	db *gorm.DB
}

func NewVIPService(db *gorm.DB) *VIPService {
	return &VIPService{db: db}
}

// VIP Level definitions (50 levels)
var VIPLevels = []map[string]interface{}{
	// Bronze Tier (0-4)
	{ "level": 0, "name": "Bronze I", "points": 0, "rakeback": 0.0, "deposit_bonus": 0 },
	{ "level": 1, "name": "Bronze II", "points": 100, "rakeback": 1.0, "deposit_bonus": 1 },
	{ "level": 2, "name": "Bronze III", "points": 500, "rakeback": 2.0, "deposit_bonus": 2 },
	{ "level": 3, "name": "Bronze IV", "points": 1000, "rakeback": 3.0, "deposit_bonus": 3 },
	{ "level": 4, "name": "Bronze V", "points": 2500, "rakeback": 4.0, "deposit_bonus": 4 },

	// Silver Tier (5-9)
	{ "level": 5, "name": "Silver I", "points": 5000, "rakeback": 5.0, "deposit_bonus": 5 },
	{ "level": 6, "name": "Silver II", "points": 10000, "rakeback": 6.0, "deposit_bonus": 6 },
	{ "level": 7, "name": "Silver III", "points": 20000, "rakeback": 7.0, "deposit_bonus": 7 },
	{ "level": 8, "name": "Silver IV", "points": 35000, "rakeback": 8.0, "deposit_bonus": 8 },
	{ "level": 9, "name": "Silver V", "points": 50000, "rakeback": 9.0, "deposit_bonus": 9 },

	// Gold Tier (10-14)
	{ "level": 10, "name": "Gold I", "points": 75000, "rakeback": 10.0, "deposit_bonus": 10 },
	{ "level": 11, "name": "Gold II", "points": 100000, "rakeback": 11.0, "deposit_bonus": 11 },
	{ "level": 12, "name": "Gold III", "points": 150000, "rakeback": 12.0, "deposit_bonus": 12 },
	{ "level": 13, "name": "Gold IV", "points": 200000, "rakeback": 13.0, "deposit_bonus": 13 },
	{ "level": 14, "name": "Gold V", "points": 300000, "rakeback": 14.0, "deposit_bonus": 14 },

	// Platinum Tier (15-19)
	{ "level": 15, "name": "Platinum I", "points": 400000, "rakeback": 15.0, "deposit_bonus": 15 },
	{ "level": 16, "name": "Platinum II", "points": 500000, "rakeback": 16.0, "deposit_bonus": 16 },
	{ "level": 17, "name": "Platinum III", "points": 650000, "rakeback": 17.0, "deposit_bonus": 17 },
	{ "level": 18, "name": "Platinum IV", "points": 800000, "rakeback": 18.0, "deposit_bonus": 18 },
	{ "level": 19, "name": "Platinum V", "points": 1000000, "rakeback": 19.0, "deposit_bonus": 19 },

	// Diamond Tier (20-24)
	{ "level": 20, "name": "Diamond I", "points": 1250000, "rakeback": 20.0, "deposit_bonus": 20 },
	{ "level": 21, "name": "Diamond II", "points": 1500000, "rakeback": 21.0, "deposit_bonus": 21 },
	{ "level": 22, "name": "Diamond III", "points": 1750000, "rakeback": 22.0, "deposit_bonus": 22 },
	{ "level": 23, "name": "Diamond IV", "points": 2000000, "rakeback": 23.0, "deposit_bonus": 23 },
	{ "level": 24, "name": "Diamond V", "points": 2500000, "rakeback": 24.0, "deposit_bonus": 24 },

	// Elite Tier (25-29)
	{ "level": 25, "name": "Elite I", "points": 3000000, "rakeback": 25.0, "deposit_bonus": 25, "priority_support": true },
	{ "level": 26, "name": "Elite II", "points": 3500000, "rakeback": 26.0, "deposit_bonus": 26, "priority_support": true },
	{ "level": 27, "name": "Elite III", "points": 4000000, "rakeback": 27.0, "deposit_bonus": 27, "priority_support": true },
	{ "level": 28, "name": "Elite IV", "points": 4500000, "rakeback": 28.0, "deposit_bonus": 28, "priority_support": true },
	{ "level": 29, "name": "Elite V", "points": 5000000, "rakeback": 29.0, "deposit_bonus": 29, "priority_support": true },

	// Champion Tier (30-34)
	{ "level": 30, "name": "Champion I", "points": 6000000, "rakeback": 30.0, "deposit_bonus": 30, "priority_support": true, "vip_host": true },
	{ "level": 31, "name": "Champion II", "points": 7000000, "rakeback": 31.0, "deposit_bonus": 31, "priority_support": true, "vip_host": true },
	{ "level": 32, "name": "Champion III", "points": 8000000, "rakeback": 32.0, "deposit_bonus": 32, "priority_support": true, "vip_host": true },
	{ "level": 33, "name": "Champion IV", "points": 9000000, "rakeback": 33.0, "deposit_bonus": 33, "priority_support": true, "vip_host": true },
	{ "level": 34, "name": "Champion V", "points": 10000000, "rakeback": 34.0, "deposit_bonus": 34, "priority_support": true, "vip_host": true },

	// Legend Tier (35-39)
	{ "level": 35, "name": "Legend I", "points": 12000000, "rakeback": 35.0, "deposit_bonus": 35, "all_bonus": true, "vip_host": true },
	{ "level": 36, "name": "Legend II", "points": 14000000, "rakeback": 36.0, "deposit_bonus": 36, "all_bonus": true, "vip_host": true },
	{ "level": 37, "name": "Legend III", "points": 16000000, "rakeback": 37.0, "deposit_bonus": 37, "all_bonus": true, "vip_host": true },
	{ "level": 38, "name": "Legend IV", "points": 18000000, "rakeback": 38.0, "deposit_bonus": 38, "all_bonus": true, "vip_host": true },
	{ "level": 39, "name": "Legend V", "points": 20000000, "rakeback": 39.0, "deposit_bonus": 39, "all_bonus": true, "vip_host": true },

	// Tiger Tier (40-44)
	{ "level": 40, "name": "Tiger I", "points": 25000000, "rakeback": 40.0, "deposit_bonus": 40, "all_bonus": true, "exclusive_tournaments": true },
	{ "level": 41, "name": "Tiger II", "points": 30000000, "rakeback": 41.0, "deposit_bonus": 41, "all_bonus": true, "exclusive_tournaments": true },
	{ "level": 42, "name": "Tiger III", "points": 35000000, "rakeback": 42.0, "deposit_bonus": 42, "all_bonus": true, "exclusive_tournaments": true },
	{ "level": 43, "name": "Tiger IV", "points": 40000000, "rakeback": 43.0, "deposit_bonus": 43, "all_bonus": true, "exclusive_tournaments": true },
	{ "level": 44, "name": "Tiger V", "points": 50000000, "rakeback": 44.0, "deposit_bonus": 44, "all_bonus": true, "exclusive_tournaments": true },

	// Eternal Tier (45-49)
	{ "level": 45, "name": "Eternal I", "points": 60000000, "rakeback": 45.0, "deposit_bonus": 45, "custom_bonus": true },
	{ "level": 46, "name": "Eternal II", "points": 70000000, "rakeback": 46.0, "deposit_bonus": 46, "custom_bonus": true },
	{ "level": 47, "name": "Eternal III", "points": 80000000, "rakeback": 47.0, "deposit_bonus": 47, "custom_bonus": true },
	{ "level": 48, "name": "Eternal IV", "points": 90000000, "rakeback": 48.0, "deposit_bonus": 48, "custom_bonus": true },
	{ "level": 49, "name": "Eternal V", "points": 100000000, "rakeback": 49.0, "deposit_bonus": 49, "custom_bonus": true },

	// Ultimate (50)
	{ "level": 50, "name": "Ultimate", "points": 200000000, "rakeback": 50.0, "deposit_bonus": 50, "custom_bonus": true, "lifetime": true },
}

// GetVIPLevel returns VIP level info for a user
func (s *VIPService) GetVIPLevel(userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	level := int(user.VIPPoints / 100000) // Simplified calculation
	if level > 50 {
		level = 50
	}

	return VIPLevels[level], nil
}

// GetAllVIPLevels returns all VIP levels
func (s *VIPService) GetAllVIPLevels() []map[string]interface{} {
	return VIPLevels
}

// CalculatePoints calculates VIP points from wager
func (s *VIPService) CalculatePoints(wagered float64) float64 {
	return wagered * 10 // 10 points per $1 wagered
}

// AwardPoints awards VIP points to a user
func (s *VIPService) AwardPoints(userID uuid.UUID, points float64) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("vip_points", gorm.Expr("vip_points + ?", points)).Error
}

// UpdateLevel checks and updates user's VIP level
func (s *VIPService) UpdateLevel(userID uuid.UUID) (int, bool, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, false, err
	}

	oldLevel := user.VIPLevel
	newLevel := calculateLevelFromPoints(user.VIPPoints)

	if newLevel != oldLevel {
		err := s.db.Model(&models.User{}).
			Where("id = ?", userID).
			Update("vip_level", newLevel).Error
		
		// Award level-up bonus
		levelUpBonus := calculateLevelUpBonus(oldLevel, newLevel)
		if levelUpBonus > 0 {
			s.db.Model(&models.User{}).
				Where("id = ?", userID).
				Update("balance", gorm.Expr("balance + ?", levelUpBonus))
		}
		
		return newLevel, true, err
	}

	return newLevel, false, nil
}

// CalculateRakeback calculates rakeback for a user
func (s *VIPService) CalculateRakeback(userID uuid.UUID, betAmount float64) (float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	level := int(user.VIPPoints / 100000)
	if level > 50 {
		level = 50
	}

	rakebackPercent := VIPLevels[level]["rakeback"].(float64)
	return betAmount * (rakebackPercent / 100), nil
}

// GetVIPBenefits returns all benefits for a user's VIP level
func (s *VIPService) GetVIPBenefits(userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	level := int(user.VIPPoints / 100000)
	if level > 50 {
		level = 50
	}

	return VIPLevels[level], nil
}

// GetWeeklyBonus returns weekly bonus for VIP level
func (s *VIPService) GetWeeklyBonus(userID uuid.UUID) (float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	level := int(user.VIPPoints / 100000)
	if level > 50 {
		level = 50
	}

	// Weekly bonus based on level
	bonusMap := map[int]float64{
		0: 0, 5: 10, 10: 25, 15: 50, 20: 100, 
		25: 200, 30: 500, 35: 1000, 40: 2500, 45: 5000, 50: 10000,
	}

	bonus := bonusMap[level]
	if bonus == 0 {
		// Linear scaling for intermediate levels
		bonus = float64(level) * 2
	}

	return bonus, nil
}

// GetBirthdayBonus returns birthday bonus for VIP level
func (s *VIPService) GetBirthdayBonus(userID uuid.UUID) (float64, error) {
	level := 20 // Birthday bonus starts at Diamond I
	
	bonusMap := map[int]float64{
		20: 500, 25: 1000, 30: 2500, 35: 5000, 40: 10000, 45: 25000, 50: 50000,
	}

	return bonusMap[level], nil
}

// Helper functions
func calculateLevelFromPoints(points float64) int {
	for i := len(VIPLevels) - 1; i >= 0; i-- {
		requiredPoints := VIPLevels[i]["points"].(int)
		if int(points) >= requiredPoints {
			return i
		}
	}
	return 0
}

func calculateLevelUpBonus(oldLevel, newLevel int) float64 {
	bonusMap := map[int]float64{
		5: 50, 10: 100, 15: 250, 20: 500, 25: 1000, 
		30: 2500, 35: 5000, 40: 10000, 45: 25000, 50: 50000,
	}

	return bonusMap[newLevel]
}
