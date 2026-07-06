package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// ReferralService handles referral and affiliate operations
type ReferralService struct {
	db *gorm.DB
}

func NewReferralService(db *gorm.DB) *ReferralService {
	return &ReferralService{db: db}
}

// GenerateReferralCode generates a unique referral code for a user
func (s *ReferralService) GenerateReferralCode(userID uuid.UUID) (string, error) {
	// Generate random 8-character code
	bytes := make([]byte, 4)
	rand.Read(bytes)
	code := hex.EncodeToString(bytes)[:8]

	// Check if code exists
	var existing models.User
	count := 0
	s.db.Model(&models.User{}).Where("referral_code = ?", code).Count(&count)
	for count > 0 {
		rand.Read(bytes)
		code = hex.EncodeToString(bytes)[:8]
		s.db.Model(&models.User{}).Where("referral_code = ?", code).Count(&count)
	}

	// Save to user
	err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("referral_code", code).Error
	return code, err
}

// GetReferralCode returns user's referral code
func (s *ReferralService) GetReferralCode(userID uuid.UUID) (string, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		return "", err
	}

	if user.ReferralCode == "" {
		return s.GenerateReferralCode(userID)
	}

	return user.ReferralCode, nil
}

// ProcessReferral processes a new user registration with a referral code
func (s *ReferralService) ProcessReferral(referrerID, refereeID uuid.UUID, code string) error {
	// Verify the code belongs to referrer
	var referrer models.User
	if err := s.db.First(&referrer, referrerID).Error; err != nil {
		return fmt.Errorf("referrer not found")
	}

	if referrer.ReferralCode != code {
		return fmt.Errorf("invalid referral code")
	}

	// Create referral relationship
	referral := models.Referral{
		ReferrerID:    referrerID,
		RefereeID:    refereeID,
		ReferralCode: code,
		Commission:    0, // Will be calculated on deposits
	}

	return s.db.Create(&referral).Error
}

// GetReferrals returns all users referred by a user
func (s *ReferralService) GetReferrals(userID uuid.UUID) ([]models.User, error) {
	var referrals []models.Referral
	err := s.db.Where("referrer_id = ?", userID).Find(&referrals).Error
	if err != nil {
		return nil, err
	}

	var refereeIDs []uuid.UUID
	for _, r := range referrals {
		refereeIDs = append(refereeIDs, r.RefereeID)
	}

	var users []models.User
	if len(refereeIDs) > 0 {
		err = s.db.Where("id IN ?", refereeIDs).Find(&users).Error
	}

	return users, err
}

// GetReferralStats returns statistics about a user's referrals
func (s *ReferralService) GetReferralStats(userID uuid.UUID) (map[string]interface{}, error) {
	var referrals []models.Referral
	err := s.db.Where("referrer_id = ?", userID).Find(&referrals).Error
	if err != nil {
		return nil, err
	}

	refereeIDs := make([]uuid.UUID, len(referrals))
	for i, r := range referrals {
		refereeIDs[i] = r.RefereeID
	}

	// Count active referees (who have made deposits)
	var activeCount int64
	if len(refereeIDs) > 0 {
		s.db.Model(&models.Transaction{}).
			Where("user_id IN ? AND type = ?", refereeIDs, "deposit").
			Distinct("user_id").
			Count(&activeCount)
	}

	// Calculate total commission
	var totalCommission float64
	s.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "referral_commission").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCommission)

	// Calculate total deposits from referees
	var totalDeposits float64
	if len(refereeIDs) > 0 {
		s.db.Model(&models.Transaction{}).
			Where("user_id IN ? AND type = ?", refereeIDs, "deposit").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalDeposits)
	}

	stats := map[string]interface{}{
		"total_referrals":   len(referrals),
		"active_referrals":  activeCount,
		"total_commission":  totalCommission,
		"total_deposits":    totalDeposits,
		"commission_rate":   0.20, // 20% commission
	}

	return stats, nil
}

// CalculateCommission calculates referral commission for a deposit
func (s *ReferralService) CalculateCommission(depositAmount float64) float64 {
	// 20% commission on deposits made by referred users
	return depositAmount * 0.20
}

// AwardCommission awards commission to referrer for referee's deposit
func (s *ReferralService) AwardCommission(referrerID, refereeID uuid.UUID, depositAmount float64) error {
	commission := s.CalculateCommission(depositAmount)

	// Create commission transaction for referrer
	commissionTx := models.Transaction{
		UserID:      referrerID,
		Type:        "referral_commission",
		Amount:      commission,
		Currency:    "USD",
		Status:       "confirmed",
		Description: fmt.Sprintf("Referral commission for deposit by referee %s", refereeID),
	}

	if err := s.db.Create(&commissionTx).Error; err != nil {
		return err
	}

	// Update referral record
	return s.db.Model(&models.Referral{}).
		Where("referrer_id = ? AND referee_id = ?", referrerID, refereeID).
		Update("commission", gorm.Expr("commission + ?", commission)).Error
}

// GetEarnings returns user's referral earnings
func (s *ReferralService) GetEarnings(userID uuid.UUID) (map[string]interface{}, error) {
	var totalCommission float64
	s.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "referral_commission").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCommission)

	var pendingCommission float64
	s.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, "referral_commission", "pending").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&pendingCommission)

	return map[string]interface{}{
		"total_earned":    totalCommission,
		"pending":         pendingCommission,
		"available":       totalCommission - pendingCommission,
		"commission_rate": "20%",
	}, nil
}
