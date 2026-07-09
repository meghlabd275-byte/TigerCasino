package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a player in the casino
type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	Username        string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash    string    `gorm:"not null" json:"-"`
	WalletAddress   string    `json:"wallet_address,omitempty"`
	Balance         float64   `gorm:"default:0" json:"balance"`
	BonusBalance    float64   `gorm:"default:0" json:"bonus_balance"`
	VIPLevel        int       `gorm:"default:0" json:"vip_level"`
	VIPPoints       float64   `gorm:"default:0" json:"vip_points"`
	KYCStatus       string    `gorm:"default:'pending'" json:"kyc_status"`
	IsVerified      bool      `gorm:"default:false" json:"is_verified"`
	IsAdmin         bool      `gorm:"default:false" json:"is_admin"`
	IsBanned        bool      `gorm:"default:false" json:"is_banned"`
	BanReason       string    `json:"ban_reason,omitempty"`
	TwoFASecret     string    `json:"-"`
	Is2FAEnabled    bool      `gorm:"default:false" json:"is_2fa_enabled"`
	EmailVerified   bool      `gorm:"default:false" json:"email_verified"`
	ReferralCode    string    `gorm:"uniqueIndex" json:"referral_code"`
	ReferredBy      *uuid.UUID `json:"referred_by,omitempty"`
	TotalWagered   float64   `gorm:"default:0" json:"total_wagered"`
	TotalDeposited  float64   `gorm:"default:0" json:"total_deposited"`
	TotalWithdrawn float64   `gorm:"default:0" json:"total_withdrawn"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
}

// Wallet represents a user's wallet
type Wallet struct {
	ID        string          `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string          `gorm:"type:uuid;not null;index" json:"user_id"`
	Balances  map[string]float64 `gorm:"-" json:"balances"`
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DepositAddress represents a deposit address
type DepositAddress struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Currency  string    `gorm:"not null" json:"currency"`
	Address   string    `gorm:"not null" json:"address"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// WithdrawalRequest represents a withdrawal request
type WithdrawalRequest struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Currency  string    `gorm:"not null" json:"currency"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Address   string    `gorm:"not null" json:"address"`
	Status    string    `gorm:"default:'pending'" json:"status"`
	TxHash    string    `json:"tx_hash,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

// VIPLevel represents VIP tier configuration
type VIPLevel struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Level                  int       `gorm:"uniqueIndex;not null" json:"level"`
	Name                   string    `gorm:"not null" json:"name"`
	DepositBonusPercent    float64   `gorm:"default:0" json:"deposit_bonus_percent"`
	WithdrawalBonusPercent float64   `gorm:"default:0" json:"withdrawal_bonus_percent"`
	RakebackPercent        float64   `gorm:"default:0" json:"rakeback_percent"`
	RequiredPoints         float64   `gorm:"default:0" json:"required_points"`
	MaxDailyWithdrawal     float64   `gorm:"default:0" json:"max_daily_withdrawal"`
	PrioritySupport       bool      `gorm:"default:false" json:"priority_support"`
	IsActive              bool      `gorm:"default:true" json:"is_active"`
	CreatedAt             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Referral represents a referral relationship
type Referral struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ReferrerID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"referrer_id"`
	RefereeID    uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"referee_id"`
	ReferralCode string    `gorm:"not null" json:"referral_code"`
	Commission    float64   `gorm:"default:0" json:"commission"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        string    `gorm:"not null" json:"type"` // deposit, withdrawal, bet, win, refund
	Amount      float64   `gorm:"not null" json:"amount"`
	Currency    string    `gorm:"not null" json:"currency"`
	Status      string    `gorm:"default:'pending'" json:"status"` // pending, confirmed, rejected
	TXHash      string    `json:"tx_hash,omitempty"`
	Address     string    `json:"address,omitempty"`
	Fee         float64   `gorm:"default:0" json:"fee"`
	GameID      *uuid.UUID `json:"game_id,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

// Game represents a casino game
type Game struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Type         string    `gorm:"not null" json:"type"` // slots, crash, mines, plinko, dice, etc.
	Provider     string    `json:"provider,omitempty"`
	RTP          float64   `json:"rtp,omitempty"` // Return to player percentage
	MinBet       float64   `json:"min_bet,omitempty"`
	MaxBet       float64   `json:"max_bet,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	IsHot        bool      `gorm:"default:false" json:"is_hot"`
	IsNew        bool      `gorm:"default:false" json:"is_new"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Category     string    `json:"category,omitempty"`
	HouseEdge    float64   `gorm:"default:0.03" json:"house_edge"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Bet represents a game bet
type Bet struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	GameID      uuid.UUID  `gorm:"type:uuid;not null" json:"game_id"`
	GameType    string    `gorm:"not null" json:"game_type"` // crash, mines, plinko, dice
	BetAmount   float64   `gorm:"not null" json:"bet_amount"`
	WinAmount   float64   `gorm:"default:0" json:"win_amount"`
	Multiplier  float64   `gorm:"default:0" json:"multiplier"`
	Profit      float64   `gorm:"default:0" json:"profit"`
	Status      string    `gorm:"default:'pending'" json:"status"` // pending, won, lost
	GameData    string    `gorm:"type:jsonb" json:"game_data,omitempty"` // JSON game-specific data
	ServerSeed  string    `json:"server_seed,omitempty"`
	ClientSeed  string    `json:"client_seed,omitempty"`
	Nonce       int       `json:"nonce"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	SettledAt   *time.Time `json:"settled_at,omitempty"`
}

// CrashGameRound represents a crash game round
type CrashGameRound struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	RoundID     string    `gorm:"uniqueIndex;not null" json:"round_id"`
	CrashPoint  float64   `gorm:"not null" json:"crash_point"`
	Status      string    `gorm:"default:'waiting'" json:"status"` // waiting, running, crashed
	ServerSeed  string    `gorm:"not null" json:"server_seed"`
	ServerHash  string    `gorm:"not null" json:"server_hash"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// SportsEvent represents a sports betting event
type SportsEvent struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ExternalID  string    `gorm:"uniqueIndex" json:"external_id"`
	Sport       string    `gorm:"not null" json:"sport"`
	League      string    `json:"league,omitempty"`
	HomeTeam   string    `gorm:"not null" json:"home_team"`
	AwayTeam   string    `gorm:"not null" json:"away_team"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	Status     string    `gorm:"default:'upcoming'" json:"status"` // upcoming, live, finished
	HomeScore  int       `gorm:"default:0" json:"home_score"`
	AwayScore  int       `gorm:"default:0" json:"away_score"`
	IsLive     bool      `gorm:"default:false" json:"is_live"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// SportsBet represents a sports betting wager
type SportsBet struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID      uuid.UUID  `gorm:"type:uuid;not null" json:"event_id"`
	BetType      string    `gorm:"not null" json:"bet_type"` // home_win, away_win, draw, over_under, etc.
	Odds         float64   `gorm:"not null" json:"odds"`
	Stake        float64   `gorm:"not null" json:"stake"`
	PotentialWin float64   `gorm:"not null" json:"potential_win"`
	ActualWin    float64   `gorm:"default:0" json:"actual_win"`
	Status       string    `gorm:"default:'pending'" json:"status"` // pending, won, lost, cancelled
	Result       string    `json:"result,omitempty"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	SettledAt    *time.Time `json:"settled_at,omitempty"`
}

// LeaderboardEntry represents a user on the leaderboard
type LeaderboardEntry struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Username     string    `gorm:"not null" json:"username"`
	TotalProfit  float64   `gorm:"default:0" json:"total_profit"`
	TotalWagered float64   `gorm:"default:0" json:"total_wagered"`
	BetCount     int       `gorm:"default:0" json:"bet_count"`
	Rank         int       `json:"rank"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// FraudAlert represents a suspicious activity alert
type FraudAlert struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AlertType     string    `gorm:"not null" json:"alert_type"` // rapid_betting, unusual_win, bonus_abuse, etc.
	Severity      string    `gorm:"not null" json:"severity"` // low, medium, high, critical
	Description   string    `json:"description"`
	Evidence      string    `gorm:"type:jsonb" json:"evidence"`
	Status        string    `gorm:"default:'open'" json:"status"` // open, investigating, resolved, dismissed
	ResolvedBy    *uuid.UUID `json:"resolved_by,omitempty"`
	ResolutionNote string   `json:"resolution_note,omitempty"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
}

// AuditLog represents system audit logs
type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Action     string    `gorm:"not null" json:"action"`
	EntityType string    `json:"entity_type,omitempty"`
	EntityID   string    `json:"entity_id,omitempty"`
	Details    string    `gorm:"type:jsonb" json:"details"`
	IPAddress  string    `json:"ip_address,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Session represents a user session
type Session struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// ============ VIP & LOYALTY MODELS ============

// RakebackBalance represents user's rakeback balance
type RakebackBalance struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Available   float64   `gorm:"not null" json:"available"`
	Locked      float64   `gorm:"default:0" json:"locked"`
	Source      string    `json:"source"`
	ExternalRef string    `json:"external_ref"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// LoyaltyPoints represents user's loyalty points
type LoyaltyPoints struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Points    int64     `gorm:"not null" json:"points"`
	Source    string    `json:"source"` // rakeback, bonus, promotion
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// BonusClaim represents a claimed bonus
type BonusClaim struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	BonusType     string    `gorm:"not null" json:"bonus_type"` // welcome, deposit, cashback, free_spins, loyalty
	Amount        float64   `gorm:"not null" json:"amount"`
	WagerRequired float64   `gorm:"not null" json:"wager_required"`
	Wagered       float64   `gorm:"default:0" json:"wagered"`
	Status        string    `gorm:"default:'active'" json:"status"` // active, completed, expired, cancelled
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// CashbackRecord represents a cashback award
type CashbackRecord struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Period    string    `gorm:"not null" json:"period"` // daily, weekly, monthly
	NetLoss   float64   `gorm:"not null" json:"net_loss"`
	Percent   float64   `gorm:"not null" json:"percent"`
	Amount    float64   `gorm:"not null" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// LevelUpReward represents a level up reward
type LevelUpReward struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	OldLevel  int       `gorm:"not null" json:"old_level"`
	NewLevel  int       `gorm:"not null" json:"new_level"`
	Bonus     float64   `gorm:"not null" json:"bonus"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Promotion represents an active promotion
type Promotion struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        string    `gorm:"not null" json:"type"` // bonus, tournament, cashback, free_spins
	BonusAmount float64   `json:"bonus_amount"`
	WagerReq    float64   `json:"wager_req"`
	StartDate  time.Time `gorm:"not null" json:"start_date"`
	EndDate    time.Time `gorm:"not null" json:"end_date"`
	Status     string    `gorm:"default:'active'" json:"status"` // active, scheduled, expired
	Terms      string    `gorm:"type:text" json:"terms"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// ============ TOURNAMENT MODELS ============

// Tournament represents a casino tournament
type Tournament struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name             string    `gorm:"not null" json:"name"`
	Description      string    `json:"description"`
	Type             string    `gorm:"not null" json:"type"` // slots, table_games, live_casino, all_games, sports
	Status           string    `gorm:"default:'upcoming'" json:"status"` // upcoming, registration, active, completed, cancelled
	GameFilter      string    `json:"game_filter"` // comma-separated game IDs or categories
	MinBet          float64   `json:"min_bet"`
	StartTime       time.Time `gorm:"not null" json:"start_time"`
	EndTime         time.Time `gorm:"not null" json:"end_time"`
	RegistrationEnd time.Time `json:"registration_end"`
	PrizePool       float64   `gorm:"not null" json:"prize_pool"`
	Currency        string    `gorm:"default:'USD'" json:"currency"`
	ScoringType     string    `gorm:"default:'wager'" json:"scoring_type"` // wager, wins, profit
	PointsMultiplier float64  `gorm:"default:1" json:"points_multiplier"`
	MinWagerToQualify float64 `json:"min_wager_to_qualify"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TournamentParticipant represents a user in a tournament
type TournamentParticipant struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TournamentID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"tournament_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	JoinedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
	Score          float64    `gorm:"default:0" json:"score"`
	Wagered        float64    `gorm:"default:0" json:"wagered"`
	Wins           int        `gorm:"default:0" json:"wins"`
	CurrentStreak int        `gorm:"default:0" json:"current_streak"`
	BestStreak    int        `gorm:"default:0" json:"best_streak"`
}

// TournamentPrize represents a prize won in a tournament
type TournamentPrize struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TournamentID uuid.UUID  `gorm:"type:uuid;not null;index" json:"tournament_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Position     int        `gorm:"not null" json:"position"`
	PrizeAmount float64    `gorm:"not null" json:"prize_amount"`
	Currency     string     `gorm:"default:'USD'" json:"currency"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// ============ GAME AGGREGATOR MODELS ============

// GameProvider represents a game provider
type GameProvider struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"not null;uniqueIndex" json:"name"`
	Code        string    `gorm:"uniqueIndex" json:"code"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website"`
	Status      string    `gorm:"default:'active'" json:"status"` // active, inactive, maintenance
	IsAggregator bool     `gorm:"default:false" json:"is_aggregator"`
	GameCount  int       `gorm:"default:0" json:"game_count"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// SportsMarket represents a betting market
type SportsMarket struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	EventID     string    `gorm:"not null;index" json:"event_id"`
	Name        string    `gorm:"not null" json:"name"`
	MarketType  string    `gorm:"not null" json:"market_type"` // moneyline, spread, over_under, etc.
	Outcomes    string    `gorm:"type:jsonb" json:"outcomes"` // JSON array of outcomes
	Suspended   bool      `gorm:"default:false" json:"suspended"`
	Status      string    `gorm:"default:'active'" json:"status"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// WalletAddress represents a user's wallet address
type WalletAddress struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Currency    string    `gorm:"not null" json:"currency"`
	Network     string    `gorm:"not null" json:"network"`
	Address     string    `gorm:"not null" json:"address"`
	IsPrimary   bool      `gorm:"default:false" json:"is_primary"`
	Status      string    `gorm:"default:'active'" json:"status"` // active, inactive
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// BalanceChange represents a balance change history
type BalanceChange struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Type       string    `gorm:"not null" json:"type"` // deposit, withdrawal, bet, win, bonus, rakeback, cashback
	Currency   string    `gorm:"not null" json:"currency"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Balance    float64   `gorm:"not null" json:"balance"`
	Reference  string    `json:"reference"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Notification represents a user notification
type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      string    `gorm:"not null" json:"type"` // info, bonus, tournament, withdrawal, deposit
	Title     string    `gorm:"not null" json:"title"`
	Message   string    `json:"message"`
	Read      bool      `gorm:"default:false" json:"read"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
