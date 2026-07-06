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
