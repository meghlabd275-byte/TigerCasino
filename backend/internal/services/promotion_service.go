package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// PromotionService handles daily races, raffles, and promotions
type PromotionService struct {
	db *gorm.DB
}

func NewPromotionService(db *gorm.DB) *PromotionService {
	return &PromotionService{db: db}
}

// ============ Daily Race ============

type DailyRace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	PrizePool float64   `json:"prize_pool"`
	Status     string    `json:"status"` // active, completed
	Entries    []RaceEntry
}

type RaceEntry struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Wagered  float64 `json:"wagered"`
	Profit   float64 `json:"profit"`
	Rank     int     `json:"rank"`
}

func (s *PromotionService) CreateDailyRace(name string, prizePool float64, durationHours int) (*DailyRace, error) {
	race := DailyRace{
		ID:         uuid.New().String(),
		Name:       name,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(time.Duration(durationHours) * time.Hour),
		PrizePool:  prizePool,
		Status:     "active",
		Entries:    []RaceEntry{},
	}
	
	// Store in database
	raceModel := models.DailyRace{
		ID:         race.ID,
		Name:       race.Name,
		StartTime:  race.StartTime,
		EndTime:    race.EndTime,
		PrizePool:  race.PrizePool,
		Status:     "active",
	}
	
	if err := s.db.Create(&raceModel).Error; err != nil {
		return nil, err
	}
	
	return &race, nil
}

func (s *PromotionService) UpdateRaceProgress(raceID string, userID, username string, wagered, profit float64) error {
	// Update race entry
	var entry models.DailyRaceEntry
	result := s.db.Where("race_id = ? AND user_id = ?", raceID, userID).First(&entry)
	
	if result.Error == gorm.ErrRecordNotFound {
		entry = models.DailyRaceEntry{
			ID:       uuid.New().String(),
			RaceID:   raceID,
			UserID:   userID,
			Username: username,
			Wagered:  wagered,
			Profit:   profit,
		}
		s.db.Create(&entry)
	} else {
		entry.Wagered += wagered
		entry.Profit += profit
		s.db.Save(&entry)
	}
	
	return nil
}

func (s *PromotionService) GetDailyRaceLeaderboard(raceID string, limit int) ([]RaceEntry, error) {
	var entries []models.DailyRaceEntry
	err := s.db.Where("race_id = ?", raceID).
		Order("wagered DESC").
		Limit(limit).
		Find(&entries).Error
	
	if err != nil {
		return nil, err
	}
	
	var leaderboard []RaceEntry
	for i, e := range entries {
		leaderboard = append(leaderboard, RaceEntry{
			UserID:   e.UserID,
			Username: e.Username,
			Wagered:  e.Wagered,
			Profit:   e.Profit,
			Rank:     i + 1,
		})
	}
	
	return leaderboard, nil
}

func (s *PromotionService) SettleDailyRace(raceID string) error {
	// Get top entries
	entries, err := s.GetDailyRaceLeaderboard(raceID, 10)
	if err != nil {
		return err
	}
	
	// Calculate prizes (prize pool distribution)
	prizeDistribution := []float64{0.25, 0.15, 0.10, 0.08, 0.07, 0.06, 0.05, 0.05, 0.04, 0.03}
	
	// Get race info
	var race models.DailyRace
	if err := s.db.First(&race, "id = ?", raceID).Error; err != nil {
		return err
	}
	
	// Award prizes
	for i, entry := range entries {
		if i >= len(prizeDistribution) {
			break
		}
		
		prize := race.PrizePool * prizeDistribution[i]
		
		// Credit prize to user
		s.db.Model(&models.User{}).
			Where("id = ?", entry.UserID).
			Update("balance", gorm.Expr("balance + ?", prize))
		
		// Create transaction
		tx := models.Transaction{
			ID:          uuid.New().String(),
			UserID:      uuid.MustParse(entry.UserID),
			Type:        "race_prize",
			Amount:      prize,
			Currency:    "USD",
			Status:      "confirmed",
			Description: fmt.Sprintf("Daily Race #%d prize - Rank %d", i+1, i+1),
		}
		s.db.Create(&tx)
	}
	
	// Update race status
	s.db.Model(&models.DailyRace{}).
		Where("id = ?", raceID).
		Update("status", "completed")
	
	return nil
}

// ============ Weekly Raffle ============

type Raffle struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	TicketCost float64   `json:"ticket_cost"`
	PrizePool float64   `json:"prize_pool"`
	Status     string    `json:"status"`
	Tickets    []RaffleTicket
}

type RaffleTicket struct {
	UserID   string `json:"user_id"`
	TicketID string `json:"ticket_id"`
	IsWinner bool   `json:"is_winner"`
}

func (s *PromotionService) CreateWeeklyRaffle(name string, ticketCost, prizePool float64) (*Raffle, error) {
	raffle := Raffle{
		ID:         uuid.New().String(),
		Name:       name,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(7 * 24 * time.Hour),
		TicketCost: ticketCost,
		PrizePool:  prizePool,
		Status:     "active",
		Tickets:    []RaffleTicket{},
	}
	
	raffleModel := models.WeeklyRaffle{
		ID:         raffle.ID,
		Name:       raffle.Name,
		StartTime:  raffle.StartTime,
		EndTime:    raffle.EndTime,
		TicketCost: raffle.TicketCost,
		PrizePool:  raffle.PrizePool,
		Status:     "active",
	}
	
	if err := s.db.Create(&raffleModel).Error; err != nil {
		return nil, err
	}
	
	return &raffle, nil
}

func (s *PromotionService) BuyRaffleTicket(raffleID, userID string, quantity int) ([]string, error) {
	var ticketIDs []string
	
	for i := 0; i < quantity; i++ {
		ticketID := uuid.New().String()
		ticketIDs = append(ticketIDs, ticketID)
		
		ticket := models.RaffleTicket{
			ID:       ticketID,
			RaffleID: raffleID,
			UserID:   userID,
			TicketNumber: i + 1,
		}
		
		if err := s.db.Create(&ticket).Error; err != nil {
			return nil, err
		}
	}
	
	return ticketIDs, nil
}

func (s *PromotionService) DrawRaffle(raffleID string, winnerCount int) ([]string, error) {
	// Get all tickets
	var tickets []models.RaffleTicket
	s.db.Where("raffle_id = ? AND is_winner = ?", raffleID, false).Find(&tickets)
	
	if len(tickets) < winnerCount {
		winnerCount = len(tickets)
	}
	
	// Random selection (simplified)
	var winners []string
	selected := make(map[int]bool)
	
	for i := 0; i < winnerCount; i++ {
		// Simple random (in production, use crypto random)
		for j := 0; j < len(tickets); j++ {
			if !selected[j] {
				selected[j] = true
				winners = append(winners, tickets[j].UserID)
				
				// Mark as winner
				s.db.Model(&models.RaffleTicket{}).
					Where("id = ?", tickets[j].ID).
					Update("is_winner", true)
				break
			}
		}
	}
	
	// Get prize amount
	var raffle models.WeeklyRaffle
	if err := s.db.First(&raffle, "id = ?", raffleID).Error; err != nil {
		return winners, err
	}
	
	prizePerWinner := raffle.PrizePool / float64(winnerCount)
	
	// Credit prizes
	for _, winnerID := range winners {
		s.db.Model(&models.User{}).
			Where("id = ?", winnerID).
			Update("balance", gorm.Expr("balance + ?", prizePerWinner))
		
		tx := models.Transaction{
			ID:          uuid.New().String(),
			UserID:      uuid.MustParse(winnerID),
			Type:        "raffle_prize",
			Amount:      prizePerWinner,
			Currency:    "USD",
			Status:      "confirmed",
			Description: "Weekly Raffle prize",
		}
		s.db.Create(&tx)
	}
	
	// Update raffle status
	s.db.Model(&models.WeeklyRaffle{}).
		Where("id = ?", raffleID).
		Update("status", "completed")
	
	return winners, nil
}

// ============ Quests/Challenges ============

type Quest struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	QuestType   string    `json:"quest_type"` // daily, weekly, milestone
	Target      float64   `json:"target"`
	Current     float64   `json:"current"`
	Reward      float64   `json:"reward"`
	ExpiresAt   time.Time `json:"expires_at"`
	Status      string    `json:"status"` // active, completed, expired
}

func (s *PromotionService) CreateQuest(name, description, questType string, target, reward float64, durationHours int) (*Quest, error) {
	quest := Quest{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		QuestType:   questType,
		Target:      target,
		Current:     0,
		Reward:      reward,
		ExpiresAt:   time.Now().Add(time.Duration(durationHours) * time.Hour),
		Status:      "active",
	}
	
	questModel := models.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		QuestType:   quest.QuestType,
		Target:      quest.Target,
		Reward:      quest.Reward,
		ExpiresAt:   quest.ExpiresAt,
		Status:      "active",
	}
	
	if err := s.db.Create(&questModel).Error; err != nil {
		return nil, err
	}
	
	return &quest, nil
}

func (s *PromotionService) UpdateQuestProgress(questID, userID string, progress float64) error {
	var userQuest models.UserQuest
	result := s.db.Where("quest_id = ? AND user_id = ?", questID, userID).First(&userQuest)
	
	if result.Error == gorm.ErrRecordNotFound {
		userQuest = models.UserQuest{
			ID:        uuid.New().String(),
			QuestID:   questID,
			UserID:    userID,
			Progress:  progress,
			Completed: false,
		}
		s.db.Create(&userQuest)
	} else {
		userQuest.Progress += progress
		
		// Check if completed
		var quest models.Quest
		if err := s.db.First(&quest, "id = ?", questID).Error; err == nil {
			if userQuest.Progress >= quest.Target {
				userQuest.Completed = true
				userQuest.CompletedAt = &[]time.Time{time.Now()}[0]
			}
		}
		
		s.db.Save(&userQuest)
	}
	
	return nil
}

func (s *PromotionService) ClaimQuestReward(questID, userID string) (float64, error) {
	var userQuest models.UserQuest
	if err := s.db.Where("quest_id = ? AND user_id = ? AND completed = ?", questID, userID, true).First(&userQuest).Error; err != nil {
		return 0, err
	}
	
	if userQuest.RewardClaimed {
		return 0, fmt.Errorf("reward already claimed")
	}
	
	// Get reward amount
	var quest models.Quest
	if err := s.db.First(&quest, "id = ?", questID).Error; err != nil {
		return 0, err
	}
	
	// Credit reward
	s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", quest.Reward))
	
	// Mark as claimed
	s.db.Model(&models.UserQuest{}).
		Where("id = ?", userQuest.ID).
		Update("reward_claimed", true)
	
	// Create transaction
	tx := models.Transaction{
		ID:          uuid.New().String(),
		UserID:      uuid.MustParse(userID),
		Type:        "quest_reward",
		Amount:      quest.Reward,
		Currency:    "USD",
		Status:      "confirmed",
		Description: fmt.Sprintf("Quest reward: %s", quest.Name),
	}
	s.db.Create(&tx)
	
	return quest.Reward, nil
}

// ============ Promo Code ============

type PromoCode struct {
	Code        string    `json:"code"`
	BonusType   string    `json:"bonus_type"` // deposit, free_spins, cashback
	BonusAmount float64   `json:"bonus_amount"`
	MinDeposit  float64   `json:"min_deposit"`
	MaxBonus    float64   `json:"max_bonus"`
	WagerReq    int       `json:"wager_req"` // times
	ExpiresAt   time.Time `json:"expires_at"`
	UsageLimit  int       `json:"usage_limit"`
	UsedCount   int       `json:"used_count"`
}

func (s *PromotionService) CreatePromoCode(code, bonusType string, bonusAmount, minDeposit, maxBonus float64, wagerReq int, durationHours int, usageLimit int) (*PromoCode, error) {
	promo := PromoCode{
		Code:        code,
		BonusType:   bonusType,
		BonusAmount: bonusAmount,
		MinDeposit:  minDeposit,
		MaxBonus:    maxBonus,
		WagerReq:    wagerReq,
		ExpiresAt:   time.Now().Add(time.Duration(durationHours) * time.Hour),
		UsageLimit:  usageLimit,
		UsedCount:   0,
	}
	
	promoModel := models.PromoCode{
		Code:        promo.Code,
		BonusType:   promo.BonusType,
		BonusAmount: promo.BonusAmount,
		MinDeposit:  promo.MinDeposit,
		MaxBonus:    promo.MaxBonus,
		WagerReq:    promo.WagerReq,
		ExpiresAt:   promo.ExpiresAt,
		UsageLimit:  promo.UsageLimit,
		UsedCount:   0,
	}
	
	if err := s.db.Create(&promoModel).Error; err != nil {
		return nil, err
	}
	
	return &promo, nil
}

func (s *PromotionService) ClaimPromoCode(code, userID string, depositAmount float64) (float64, error) {
	var promo models.PromoCode
	if err := s.db.Where("code = ?", code).First(&promo).Error; err != nil {
		return 0, fmt.Errorf("invalid promo code")
	}
	
	// Check expiration
	if time.Now().After(promo.ExpiresAt) {
		return 0, fmt.Errorf("promo code expired")
	}
	
	// Check usage limit
	if promo.UsedCount >= promo.UsageLimit {
		return 0, fmt.Errorf("promo code usage limit reached")
	}
	
	// Check min deposit
	if depositAmount < promo.MinDeposit {
		return 0, fmt.Errorf("minimum deposit not met")
	}
	
	// Calculate bonus
	bonus := depositAmount * (promo.BonusAmount / 100)
	if bonus > promo.MaxBonus {
		bonus = promo.MaxBonus
	}
	
	// Credit bonus
	s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("bonus_balance", gorm.Expr("bonus_balance + ?", bonus))
	
	// Update usage
	s.db.Model(&models.PromoCode{}).
		Where("code = ?", code).
		Update("used_count", gorm.Expr("used_count + 1"))
	
	// Record usage
	usage := models.PromoCodeUsage{
		ID:        uuid.New().String(),
		Code:      code,
		UserID:    userID,
		Bonus:     bonus,
	}
	s.db.Create(&usage)
	
	return bonus, nil
}
