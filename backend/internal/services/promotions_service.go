package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ============ Promotions Service ============

type PromotionsService struct{}

func NewPromotionsService() *PromotionsService {
	return &PromotionsService{}
}

// Promotion types
type Promotion struct {
	ID            string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"` // welcome, deposit, reload, cashback, free_spins, tournament, vip
	BonusAmount  float64   `json:"bonus_amount"`
	BonusPercent float64   `json:"bonus_percent"`
	MaxBonus    float64   `json:"max_bonus"`
	MinDeposit  float64   `json:"min_deposit"`
	MinWager   float64   `json:"min_wager"`
	WagerRequirement float64 `json:"wager_requirement"` // e.g., 30x bonus
	QualifyingGames []string `json:"qualifying_games"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"` // active, expired, upcoming
	MaxClaimants int      `json:"max_claimants"`
	ClaimedCount int     `json:"claimed_count"`
	Code        string    `json:"code"`
	IsFeatured  bool     `json:"is_featured"`
	Tier        string    `json:"tier"` // all, bronze, silver, gold, platinum, diamond
}

// User's claimed bonus
type ClaimedBonus struct {
	ID              string    `json:"id"`
	UserID         string    `json:"user_id"`
	PromotionID    string    `json:"promotion_id"`
	BonusAmount    float64   `json:"bonus_amount"`
	OriginalAmount float64   `json:"original_amount"`
	WageredAmount  float64   `json:"wagered_amount"`
	WagerRequired float64   `json:"wager_required"`
	ClaimedAt      time.Time `json:"claimed_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	Status        string    `json:"status"` // active, completed, expired, cancelled
}

// Free spins promotion
type FreeSpinsPromotion struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	GameID     string    `json:"game_id"`
	GameName   string    `json:"game_name"`
	SpinsCount int       `json:"spins_count"`
	MinDeposit float64   `json:"min_deposit"`
	WagerRequirement float64 `json:"wager_requirement"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	MaxClaims  int      `json:"max_claims"`
	Claimed    int      `json:"claimed"`
	Status    string    `json:"status"`
}

// Tournament promotion
type TournamentPromotion struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // slots, table, live, sports
	MinBet     float64   `json:"min_bet"`
	MinWager  float64   `json:"min_wager"`
	PrizePool  float64   `json:"prize_pool"`
	Prizes     []TournamentPrize `json:"prizes"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status    string    `json:"status"`
}

type TournamentPrize struct {
	Position int     `json:"position"`
	MinRank  int     `json:"min_rank"`
	MaxRank  int     `json:"max_rank"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"` // bonus, cash, free_spins
}

// Get all active promotions
func (s *PromotionsService) GetActivePromotions() []Promotion {
	promotions := []Promotion{
		{
			ID:            "welcome_100",
			Name:         "100% Welcome Bonus up to $1,000",
			Description:  "Get a 100% match on your first deposit",
			Type:         "welcome",
			BonusPercent: 100,
			MaxBonus:    1000,
			MinDeposit:  20,
			WagerRequirement: 30,
			StartDate:   time.Now().Add(-30 * 24 * time.Hour),
			EndDate:     time.Now().Add(60 * 24 * time.Hour),
			Status:      "active",
			MaxClaimants: 10000,
			ClaimedCount: 5432,
			Code:        "WELCOME100",
			IsFeatured:  true,
			Tier:       "all",
		},
		{
			ID:            "reload_50",
			Name:         "50% Weekend Reload Bonus",
			Description:  "Get 50% extra on weekend deposits",
			Type:         "reload",
			BonusPercent: 50,
			MaxBonus:    500,
			MinDeposit:  50,
			WagerRequirement: 25,
			StartDate:   time.Now().Add(-7 * 24 * time.Hour),
			EndDate:     time.Now().Add(7 * 24 * time.Hour),
			Status:      "active",
			Code:        "WEEKEND50",
			Tier:       "all",
		},
		{
			ID:            "cashback_10",
			Name:         "10% Daily Cashback",
			Description:  "Get 10% back on daily losses",
			Type:         "cashback",
			BonusAmount:  10,
			MaxBonus:    100,
			MinWager:   100,
			WagerRequirement: 5,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(30 * 24 * time.Hour),
			Status:      "active",
			Tier:        "all",
		},
		{
			ID:            "vip_bonus",
			Name:         "VIP Exclusive 25% Bonus",
			Description:  "Exclusive bonus for Gold+ members",
			Type:         "deposit",
			BonusPercent: 25,
			MaxBonus:    5000,
			MinDeposit:  100,
			WagerRequirement: 20,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(90 * 24 * time.Hour),
			Status:      "active",
			IsFeatured:  true,
			Tier:        "gold",
		},
		{
			ID:            "crypto_deposit",
			Name:         "Crypto Bonus - Extra 5%",
			Description:  "Get extra 5% on crypto deposits",
			Type:         "deposit",
			BonusPercent: 5,
			MaxBonus:    500,
			MinDeposit:  50,
			WagerRequirement: 15,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(180 * 24 * time.Hour),
			Status:      "active",
			Code:        "CRYPTO5",
			Tier:        "all",
		},
	}
	return promotions
}

// Get promotion by ID
func (s *PromotionsService) GetPromotion(id string) (*Promotion, error) {
	promos := s.GetActivePromotions()
	for _, p := range promos {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("promotion not found")
}

// Get featured promotions
func (s *PromotionsService) GetFeaturedPromotions() []Promotion {
	var featured []Promotion
	for _, p := range s.GetActivePromotions() {
		if p.IsFeatured {
			featured = append(featured, p)
		}
	}
	return featured
}

// Get promotions by tier
func (s *PromotionsService) GetPromotionsByTier(tier string) []Promotion {
	var tierPromos []Promotion
	for _, p := range s.GetActivePromotions() {
		if p.Tier == "all" || p.Tier == tier {
			tierPromos = append(tierPromos, p)
		}
	}
	return tierPromos
}

// Claim a promotion
func (s *PromotionsService) ClaimPromotion(userID, promoID string) (*ClaimedBonus, error) {
	promo, err := s.GetPromotion(promoID)
	if err != nil {
		return nil, err
	}

	if promo.ClaimedCount >= promo.MaxClaimants {
		return nil, fmt.Errorf("promotion fully claimed")
	}

	bonus := &ClaimedBonus{
		ID:              uuid.New().String(),
		UserID:         userID,
		PromotionID:    promoID,
		BonusAmount:    promo.MaxBonus,
		OriginalAmount: promo.MaxBonus,
		WageredAmount:  0,
		WagerRequired:  promo.WagerRequirement,
		ClaimedAt:      time.Now(),
		ExpiresAt:     time.Now().Add(7 * 24 * time.Hour),
		Status:        "active",
	}

	return bonus, nil
}

// Calculate bonus after wagering
func (s *PromotionsService) CalculateBonusWager(bonus *ClaimedBonus, newWager float64) {
	bonus.WageredAmount += newWager
	
	// Check if wager requirement met
	if bonus.WageredAmount >= bonus.WagerRequired {
		// Bonus becomes withdrawable
		bonus.Status = "completed"
	}
}

// Get free spins promotions
func (s *PromotionsService) GetFreeSpinsPromotions() []FreeSpinsPromotion {
	return []FreeSpinsPromotion{
		{
			ID:            "fs_100_welcome",
			Name:        "100 Free Spins Welcome Package",
			GameID:      "pp_gates_of_olympus",
			GameName:    "Gates of Olympus",
			SpinsCount:  100,
			MinDeposit:  50,
			WagerRequirement: 35,
			StartDate:  time.Now().Add(-30 * 24 * time.Hour),
			EndDate:    time.Now().Add(60 * 24 * time.Hour),
			MaxClaims:  5000,
			Claimed:    2341,
			Status:    "active",
		},
		{
			ID:            "fs_50_deposit",
			Name:        "50 Free Spins on Sweet Bonanza",
			GameID:      "pp_sweet_bonanza",
			GameName:    "Sweet Bonanza",
			SpinsCount:  50,
			MinDeposit:  30,
			WagerRequirement: 30,
			StartDate:  time.Now(),
			EndDate:    time.Now().Add(14 * 24 * time.Hour),
			MaxClaims:  2000,
			Claimed:    1234,
			Status:    "active",
		},
	}
}

// Get tournament promotions
func (s *PromotionsService) GetTournaments() []TournamentPromotion {
	return []TournamentPromotion{
		{
			ID:       "slots_weekly_10k",
			Name:     "Weekly Slots Tournament - $10,000 Prize Pool",
			Type:     "slots",
			MinBet:   0.50,
			PrizePool: 10000,
			Prizes: []TournamentPrize{
				{Position: 1, MinRank: 1, MaxRank: 1, Amount: 2500, Type: "cash"},
				{Position: 2, MinRank: 2, MaxRank: 2, Amount: 1500, Type: "cash"},
				{Position: 3, MinRank: 3, MaxRank: 3, Amount: 1000, Type: "cash"},
				{Position: 4, MinRank: 4, MaxRank: 10, Amount: 500, Type: "bonus"},
				{Position: 5, MinRank: 11, MaxRank: 50, Amount: 100, Type: "bonus"},
			},
			StartDate: time.Now().Add(-3 * 24 * time.Hour),
			EndDate:   time.Now().Add(4 * 24 * time.Hour),
			Status:   "active",
		},
		{
			ID:       "live_daily_5k",
			Name:     "Daily Live Casino Challenge - $5,000",
			Type:     "live",
			MinBet:   1.00,
			PrizePool: 5000,
			Prizes: []TournamentPrize{
				{Position: 1, MinRank: 1, MaxRank: 1, Amount: 1500, Type: "cash"},
				{Position: 2, MinRank: 2, MaxRank: 2, Amount: 1000, Type: "cash"},
				{Position: 3, MinRank: 3, MaxRank: 3, Amount: 500, Type: "cash"},
				{Position: 4, MinRank: 4, MaxRank: 10, Amount: 250, Type: "bonus"},
			},
			StartDate: time.Now().Add(-12 * time.Hour),
			EndDate:   time.Now().Add(12 * time.Hour),
			Status:   "active",
		},
		{
			ID:       "rakeback_monthly",
			Name:     "Monthly Rakeback - Up to 25%",
			Type:     "table",
			MinWager: 100,
			PrizePool: 50000,
			Prizes: []TournamentPrize{
				{Position: 1, MinRank: 1, MaxRank: 1, Amount: 10000, Type: "cash"},
				{Position: 2, MinRank: 2, MaxRank: 5, Amount: 5000, Type: "cash"},
				{Position: 3, MinRank: 6, MaxRank: 20, Amount: 2500, Type: "bonus"},
			},
			StartDate: time.Now().Add(-15 * 24 * time.Hour),
			EndDate:   time.Now().Add(15 * 24 * time.Hour),
			Status:   "active",
		},
	}
}

// Referral promotion
type ReferralPromotion struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	ReferrerBonus     float64   `json:"referrer_bonus"`
	ReferralBonus     float64   `json:"referral_bonus"`
	MinDeposit        float64   `json:"min_deposit"`
	MaxReferrals     int       `json:"max_referrals"`
	WagerRequirement float64   `json:"wager_requirement"`
	Status           string    `json:"status"`
}

func (s *PromotionsService) GetReferralProgram() *ReferralPromotion {
	return &ReferralPromotion{
		ID:                "referral_program",
		Name:              "Refer a Friend - Earn $50",
		ReferrerBonus:     50,
		ReferralBonus:     20,
		MinDeposit:        100,
		MaxReferrals:     20,
		WagerRequirement: 30,
		Status:           "active",
	}
}

// Promotion stats
type PromotionStats struct {
	TotalPromos      int     `json:"total_promos"`
	ActivePromos     int     `json:"active_promos"`
	TotalBonusPaid   float64 `json:"total_bonus_paid"`
	TotalWagered    float64 `json:"total_wagered"`
	AvgBonusSize    float64 `json:"avg_bonus_size"`
	PopularPromo    string  `json:"popular_promo"`
}

func (s *PromotionsService) GetStats() *PromotionStats {
	return &PromotionStats{
		TotalPromos:    15,
		ActivePromos:   8,
		TotalBonusPaid: 543210.50,
		TotalWagered:   2345678.90,
		AvgBonusSize:   125.50,
		PopularPromo:   "WELCOME100",
	}
}
