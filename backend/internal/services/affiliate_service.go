package services

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// AffiliateService handles all affiliate operations
type AffiliateService struct {
	db *gorm.DB
}

// Affiliate represents an affiliate partner
type Affiliate struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	AffiliateCode  string    `gorm:"uniqueIndex"`
	ParentID       *uuid.UUID `gorm:"type:uuid"`
	Level          int       `gorm:"default:1"`
	Tier           string    `gorm:"default:'bronze'"` // bronze, silver, gold, platinum, diamond
	CommissionRate float64   `gorm:"default:0.20"` // 20% default
	TotalReferrals int       `gorm:"default:0"`
	ActiveReferrals int      `gorm:"default:0"`
	TotalEarnings  float64  `gorm:"default:0"`
	PendingEarnings float64  `gorm:"default:0"`
	PaidEarnings   float64  `gorm:"default:0"`
	Clicks         int64    `gorm:"default:0"`
	Conversions    int      `gorm:"default:0"`
	ConversionRate float64  `gorm:"default:0"`
	Status         string   `gorm:"default:'active'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// AffiliateCommission tracks commission earnings
type AffiliateCommission struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	AffiliateID   uuid.UUID `gorm:"type:uuid;index"`
	UserID        uuid.UUID `gorm:"type:uuid;index"`
	ReferralID    uuid.UUID `gorm:"type:uuid"`
	BetID         uuid.UUID `gorm:"type:uuid"`
	GameID        string
	NetRevenue    float64   // Revenue after bonuses and fees
	Commission    float64   // Commission amount
	Tier          string
	Status        string   `gorm:"default:'pending'"` // pending, approved, paid, cancelled
	PaymentID     *uuid.UUID
	CreatedAt     time.Time
	PaidAt        *time.Time
}

// AffiliateClick tracks link clicks
type AffiliateClick struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	AffiliateID  uuid.UUID `gorm:"type:uuid;index"`
	IPAddress    string
	UserAgent    string
	Country      string
	Referrer     string
	Campaign     string
	Timestamp    time.Time
}

// Tier configuration
var TierConfig = map[string]struct {
	MinReferrals   int
	MinEarnings    float64
	CommissionRate float64
}{
	"bronze":   {0, 0, 0.20},
	"silver":   {10, 1000, 0.25},
	"gold":     {25, 5000, 0.30},
	"platinum": {50, 15000, 0.35},
	"diamond":  {100, 50000, 0.40},
}

func NewAffiliateService(db *gorm.DB) *AffiliateService {
	return &AffiliateService{db: db}
}

// CreateAffiliate creates a new affiliate account
func (s *AffiliateService) CreateAffiliate(userID uuid.UUID, parentCode string) (*Affiliate, error) {
	// Generate unique affiliate code
	affiliateCode := generateAffiliateCode()

	// Check parent affiliate if provided
	var parentID *uuid.UUID
	if parentCode != "" {
		var parent Affiliate
		if err := s.db.Where("affiliate_code = ?", parentCode).First(&parent).Error; err == nil {
			parentID = &parent.ID
		}
	}

	affiliate := &Affiliate{
		ID:             uuid.New(),
		UserID:         userID,
		AffiliateCode:  affiliateCode,
		ParentID:       parentID,
		Level:          1,
		Tier:           "bronze",
		CommissionRate: 0.20,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.Create(affiliate).Error; err != nil {
		return nil, err
	}

	// Update parent's referral count if there's a parent
	if parentID != nil {
		s.db.Model(&Affiliate{}).Where("id = ?", parentID).Updates(map[string]interface{}{
			"total_referrals": gorm.Expr("total_referrals + 1"),
			"active_referrals": gorm.Expr("active_referrals + 1"),
		})
	}

	return affiliate, nil
}

// GetAffiliateByUserID retrieves affiliate by user ID
func (s *AffiliateService) GetAffiliateByUserID(userID uuid.UUID) (*Affiliate, error) {
	var affiliate Affiliate
	if err := s.db.Where("user_id = ?", userID).First(&affiliate).Error; err != nil {
		return nil, err
	}
	return &affiliate, nil
}

// GetAffiliateByCode retrieves affiliate by code
func (s *AffiliateService) GetAffiliateByCode(code string) (*Affiliate, error) {
	var affiliate Affiliate
	if err := s.db.Where("affiliate_code = ?", code).First(&affiliate).Error; err != nil {
		return nil, err
	}
	return &affiliate, nil
}

// TrackClick records an affiliate link click
func (s *AffiliateService) TrackClick(affiliateID uuid.UUID, ip, userAgent, country, referrer, campaign string) error {
	click := &AffiliateClick{
		ID:          uuid.New(),
		AffiliateID: affiliateID,
		IPAddress:   ip,
		UserAgent:   userAgent,
		Country:     country,
		Referrer:    referrer,
		Campaign:    campaign,
		Timestamp:   time.Now(),
	}

	if err := s.db.Create(click).Error; err != nil {
		return err
	}

	// Update affiliate click count
	s.db.Model(&Affiliate{}).Where("id = ?", affiliateID).Update("clicks", gorm.Expr("clicks + 1"))

	return nil
}

// RecordBetCommission records commission for a bet
func (s *AffiliateService) RecordBetCommission(affiliateID, userID, referralID, betID uuid.UUID, gameID string, netRevenue float64) error {
	affiliate, err := s.GetAffiliateByUserID(affiliateID)
	if err != nil {
		return err
	}

	// Calculate commission based on tier
	commission := netRevenue * affiliate.CommissionRate

	// Create commission record
	comm := &AffiliateCommission{
		ID:           uuid.New(),
		AffiliateID:  affiliate.ID,
		UserID:       userID,
		ReferralID:   referralID,
		BetID:        betID,
		GameID:       gameID,
		NetRevenue:   netRevenue,
		Commission:   commission,
		Tier:         affiliate.Tier,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(comm).Error; err != nil {
		return err
	}

	// Update affiliate earnings
	s.db.Model(&Affiliate{}).Where("id = ?", affiliate.ID).Updates(map[string]interface{}{
		"pending_arnings": gorm.Expr("pending_arnings + ?", commission),
		"total_arnings":   gorm.Expr("total_arnings + ?", commission),
	})

	// Update conversion stats
	s.db.Model(&Affiliate{}).Where("id = ?", affiliate.ID).Updates(map[string]interface{}{
		"conversions":    gorm.Expr("conversions + 1"),
		"conversion_rate": gorm.Expr("clicks > 0.0 THEN (conversions::float / clicks::float) ELSE 0.0"),
	})

	return nil
}

// ApproveCommission approves pending commissions
func (s *AffiliateService) ApproveCommission(commissionID uuid.UUID) error {
	var comm AffiliateCommission
	if err := s.db.First(&comm, commissionID).Error; err != nil {
		return err
	}

	// Update commission status
	if err := s.db.Model(&comm).Update("status", "approved").Error; err != nil {
		return err
	}

	// Move from pending to approved earnings
	s.db.Model(&Affiliate{}).Where("id = ?", comm.AffiliateID).Updates(map[string]interface{}{
		"pending_arnings": gorm.Expr("pending_arnings - ?", comm.Commission),
	})

	return nil
}

// PayCommission marks commission as paid
func (s *AffiliateService) PayCommission(commissionID uuid.UUID, paymentID uuid.UUID) error {
	var comm AffiliateCommission
	if err := s.db.First(&comm, commissionID).Error; err != nil {
		return err
	}

	now := time.Now()
	if err := s.db.Model(&comm).Updates(map[string]interface{}{
		"status":   "paid",
		"payment_id": paymentID,
		"paid_at":   now,
	}).Error; err != nil {
		return err
	}

	// Update affiliate paid earnings
	s.db.Model(&Affiliate{}).Where("id = ?", comm.AffiliateID).Updates(map[string]interface{}{
		"paid_arnings": gorm.Expr("paid_arnings + ?", comm.Commission),
	})

	return nil
}

// GetAffiliateStats returns affiliate statistics
func (s *AffiliateService) GetAffiliateStats(affiliateID uuid.UUID) (map[string]interface{}, error) {
	var affiliate Affiliate
	if err := s.db.First(&affiliate, affiliateID).Error; err != nil {
		return nil, err
	}

	// Get commissions by status
	var pending, approved, paid float64
	s.db.Model(&AffiliateCommission{}).Where("affiliate_id = ? AND status = ?", affiliateID, "pending").Select("COALESCE(SUM(commission), 0)").Scan(&pending)
	s.db.Model(&AffiliateCommission{}).Where("affiliate_id = ? AND status = ?", affiliateID, "approved").Select("COALESCE(SUM(commission), 0)").Scan(&approved)
	s.db.Model(&AffiliateCommission{}).Where("affiliate_id = ? AND status = ?", affiliateID, "paid").Select("COALESCE(SUM(commission), 0)").Scan(&paid)

	// Get recent clicks
	var recentClicks int64
	s.db.Model(&AffiliateClick{}).Where("affiliate_id = ? AND timestamp > ?", affiliateID, time.Now().AddDate(0, 0, -7)).Count(&recentClicks)

	// Get top countries
	type CountryStats struct {
		Country string
		Clicks  int64
	}
	var countries []CountryStats
	s.db.Model(&AffiliateClick{}).Where("affiliate_id = ?", affiliateID).
		Select("country, COUNT(*) as clicks").
		Group("country").
		Order("clicks DESC").
		Limit(5).
		Scan(&countries)

	stats := map[string]interface{}{
		"affiliate_code":   affiliate.AffiliateCode,
		"tier":             affiliate.Tier,
		"commission_rate":   affiliate.CommissionRate,
		"total_referrals":  affiliate.TotalReferrals,
		"active_referrals": affiliate.ActiveReferrals,
		"total_earnings":   affiliate.TotalEarnings,
		"pending_earnings": pending,
		"approved_earnings": approved,
		"paid_earnings":    paid,
		"total_clicks":     affiliate.Clicks,
		"recent_clicks":    recentClicks,
		"conversions":      affiliate.Conversions,
		"conversion_rate":  affiliate.ConversionRate,
		"countries":       countries,
	}

	return stats, nil
}

// GetLeaderboard returns top affiliates
func (s *AffiliateService) GetLeaderboard(limit int) ([]Affiliate, error) {
	var affiliates []Affiliate
	if err := s.db.Order("total_earnings DESC").Limit(limit).Find(&affiliates).Error; err != nil {
		return nil, err
	}
	return affiliates, nil
}

// UpdateAffiliateTier updates affiliate tier based on performance
func (s *AffiliateService) UpdateAffiliateTier(affiliateID uuid.UUID) error {
	var affiliate Affiliate
	if err := s.db.First(&affiliate, affiliateID).Error; err != nil {
		return err
	}

	// Determine new tier
	newTier := "bronze"
	for tier, config := range TierConfig {
		if affiliate.TotalReferrals >= config.MinReferrals && affiliate.TotalEarnings >= config.MinEarnings {
			newTier = tier
		}
	}

	// Update if tier changed
	if newTier != affiliate.Tier {
		newRate := TierConfig[newTier].CommissionRate
		s.db.Model(&affiliate).Updates(map[string]interface{}{
			"tier":            newTier,
			"commission_rate": newRate,
			"level":           getTierLevel(newTier),
		})
	}

	return nil
}

// GetSubAffiliates returns all sub-affiliates
func (s *AffiliateService) GetSubAffiliates(affiliateID uuid.UUID) ([]Affiliate, error) {
	var affiliates []Affiliate
	if err := s.db.Where("parent_id = ?", affiliateID).Find(&affiliates).Error; err != nil {
		return nil, err
	}
	return affiliates, nil
}

// GenerateAffiliateReport generates a comprehensive affiliate report
func (s *AffiliateService) GenerateAffiliateReport(affiliateID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	var commissions []AffiliateCommission
	s.db.Where("affiliate_id = ? AND created_at BETWEEN ? AND ?", affiliateID, startDate, endDate).Find(&commissions)

	var totalRevenue, totalCommission float64
	var betCount int
	gameStats := make(map[string]float64)

	for _, comm := range commissions {
		totalRevenue += comm.NetRevenue
		totalCommission += comm.Commission
		betCount++
		gameStats[comm.GameID] += comm.NetRevenue
	}

	// Get click stats for period
	var clicks, uniqueVisitors int64
	s.db.Model(&AffiliateClick{}).Where("affiliate_id = ? AND timestamp BETWEEN ? AND ?", affiliateID, startDate, endDate).Count(&clicks)
	s.db.Model(&AffiliateClick{}).Where("affiliate_id = ? AND timestamp BETWEEN ? AND ?", affiliateID, startDate, endDate).Distinct("ip_address").Count(&uniqueVisitors)

	report := map[string]interface{}{
		"period": map[string]string{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
		"revenue": map[string]interface{}{
			"total_net_revenue":  totalRevenue,
			"total_commission":  totalCommission,
			"bet_count":          betCount,
			"average_bet_value":  func() float64 { if betCount > 0 { return totalRevenue / float64(betCount) }; return 0 }(),
			"conversion_rate":    func() float64 { if uniqueVisitors > 0 { return float64(betCount) / float64(uniqueVisitors) }; return 0 }(),
		},
		"traffic": map[string]interface{}{
			"total_clicks":       clicks,
			"unique_visitors":    uniqueVisitors,
			"click_through_rate": func() float64 { if clicks > 0 { return float64(uniqueVisitors) / float64(clicks) }; return 0 }(),
		},
		"by_game": gameStats,
	}

	return report, nil
}

// Helper functions
func generateAffiliateCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(time.Nanosecond)
	}
	return string(code)
}

func getTierLevel(tier string) int {
	levels := map[string]int{"bronze": 1, "silver": 2, "gold": 3, "platinum": 4, "diamond": 5}
	if level, ok := levels[tier]; ok {
		return level
	}
	return 1
}

// CalculateMLCProfit calculates multi-level commission
func (s *AffiliateService) CalculateMLCProfit(affiliateID uuid.UUID, betAmount float64) error {
	// Get affiliate chain
	var affiliate Affiliate
	if err := s.db.First(&affiliate, affiliateID).Error; err != nil {
		return err
	}

	// Calculate commission for each level up the chain
	commissionRate := affiliate.CommissionRate
	currentAffiliate := &affiliate

	for currentAffiliate.ParentID != nil {
		parentID := *currentAffiliate.ParentID
		var parent Affiliate
		if err := s.db.First(&parent, parentID).Error; err != nil {
			break
		}

		// Parent gets percentage of child's commission
		parentCommission := betAmount * commissionRate * 0.10 // 10% of child's commission

		comm := &AffiliateCommission{
			ID:          uuid.New(),
			AffiliateID: parent.ID,
			UserID:      parent.UserID,
			Commission:  parentCommission,
			Tier:        parent.Tier,
			Status:      "pending",
			CreatedAt:   time.Now(),
		}
		s.db.Create(comm)

		// Update parent earnings
		s.db.Model(&Affiliate{}).Where("id = ?", parent.ID).Updates(map[string]interface{}{
			"pending_arnings": gorm.Expr("pending_arnings + ?", parentCommission),
			"total_arnings":   gorm.Expr("total_arnings + ?", parentCommission),
		})

		currentAffiliate = &parent
		commissionRate *= 0.5 // Reduce rate for each level
	}

	return nil
}

// GetTopPerformers returns top performing affiliates
func (s *AffiliateService) GetTopPerformers(limit int, timeframe string) ([]map[string]interface{}, error) {
	var startDate time.Time
	switch timeframe {
	case "today":
		startDate = time.Now().AddDate(0, 0, -1)
	case "week":
		startDate = time.Now().AddDate(0, 0, -7)
	case "month":
		startDate = time.Now().AddDate(0, -1, 0)
	default:
		startDate = time.Now().AddDate(0, -1, 0)
	}

	type Result struct {
		AffiliateID    uuid.UUID
		TotalEarnings   float64
		TotalReferrals  int
		TotalClicks    int64
		ConversionRate  float64
	}
	var results []Result

	s.db.Model(&AffiliateCommission{}).
		Select("affiliate_id, SUM(commission) as total_earnings").
		Where("created_at > ?", startDate).
		Group("affiliate_id").
		Order("total_earnings DESC").
		Limit(limit).
		Scan(&results)

	var topPerformers []map[string]interface{}
	for _, r := range results {
		var aff Affiliate
		if err := s.db.First(&aff, r.AffiliateID).Error; err == nil {
			topPerformers = append(topPerformers, map[string]interface{}{
				"affiliate_code":  aff.AffiliateCode,
				"total_earnings":  r.TotalEarnings,
				"referrals":       aff.TotalReferrals,
				"clicks":          aff.Clicks,
				"conversion_rate": math.Round(aff.ConversionRate*100) / 100,
				"tier":            aff.Tier,
			})
		}
	}

	return topPerformers, nil
}
