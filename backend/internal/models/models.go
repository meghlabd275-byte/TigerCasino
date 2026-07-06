package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email         string        `gorm:"uniqueIndex;size:255" json:"email"`
	Username      string        `gorm:"uniqueIndex;size:50" json:"username"`
	PasswordHash  string        `gorm:"size:255" json:"-"`
	WalletAddress string        `gorm:"size:100" json:"walletAddress,omitempty"`
	WalletType    string        `gorm:"size:20" json:"walletType,omitempty"`
	Balance       float64       `gorm:"default:0" json:"balance"`
	BonusBalance  float64       `gorm:"default:0" json:"bonusBalance"`
	VIPLevel      int           `gorm:"default:0" json:"vipLevel"`
	KYCStatus     string        `gorm:"size:20;default:pending" json:"kycStatus"`
	IsVerified    bool          `gorm:"default:false" json:"isVerified"`
	IsAdmin       bool          `gorm:"default:false" json:"isAdmin"`
	IsSuperAdmin  bool          `gorm:"default:false" json:"isSuperAdmin"`
	IsBanned      bool          `gorm:"default:false" json:"isBanned"`
	BanReason     string        `json:"banReason,omitempty"`
	TwoFASecret  string        `gorm:"size:255" json:"-"`
	Is2FAEnabled  bool          `gorm:"default:false" json:"is2FAEnabled"`
	EmailVerified bool          `gorm:"default:false" json:"emailVerified"`
	PhoneVerified bool          `gorm:"default:false" json:"phoneVerified"`
	Phone         string        `gorm:"size:20" json:"phone,omitempty"`
	// White Label fields
	WhiteLabelID   *uuid.UUID     `gorm:"type:uuid" json:"whiteLabelId,omitempty"`
	WhiteLabel    *WhiteLabel   `gorm:"foreignKey:WhiteLabelID" json:"whiteLabel,omitempty"`
	CommissionRate float64       `gorm:"default:80" json:"commissionRate"` // 80% to brand, 20% to TigerCasino
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	LastLogin     *time.Time    `json:"lastLogin,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate generates UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"userId"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Type        string         `gorm:"size:20;not null" json:"type"` // deposit, withdrawal, bet, win, bonus
	Amount      float64        `gorm:"not null" json:"amount"`
	Currency    string         `gorm:"size:20;not null" json:"currency"`
	Status      string         `gorm:"size:20;default:pending" json:"status"` // pending, completed, rejected
	TxHash      string         `gorm:"size:100" json:"txHash,omitempty"`
	Address     string         `gorm:"size:100" json:"address,omitempty"`
	Fee         float64        `json:"fee"`
	Chain       string         `gorm:"size:20" json:"chain"` // btc, eth, trx, etc.
	PlatformFee float64        `json:"platformFee"` // 20% fee for white label
	BrandRevenue float64       `json:"brandRevenue"` // 80% to white label brand
	CreatedAt   time.Time      `json:"createdAt"`
	ProcessedAt *time.Time    `json:"processedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// WhiteLabel represents a white label brand
type WhiteLabel struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name            string        `gorm:"size:100;not null" json:"name"`
	BrandDomain     string        `gorm:"size:100;uniqueIndex" json:"brandDomain"`
	LogoURL         string        `gorm:"size:255" json:"logoUrl"`
	PrimaryColor   string        `gorm:"size:20" json:"primaryColor"`
	SecondaryColor string        `gorm:"size:20" json:"secondaryColor"`
	Status         string        `gorm:"size:20;default:pending" json:"status"` // pending, active, suspended, rejected
	OwnerID        uuid.UUID     `gorm:"type:uuid;index" json:"ownerId"`
	Owner          User          `gorm:"foreignKey:OwnerID" json:"-"`
	CommissionRate float64       `gorm:"default:80" json:"commissionRate"` // 80% to brand, 20% to TigerCasino
	FeeWallet      string        `gorm:"size:100" json:"feeWallet"` // TigerCasino fee wallet address
	TotalRevenue   float64       `gorm:"default:0" json:"totalRevenue"`
	TotalFees     float64        `gorm:"default:0" json:"totalFees"`
	ApprovedBy    *uuid.UUID    `gorm:"type:uuid" json:"approvedBy,omitempty"`
	ApprovedAt    *time.Time   `json:"approvedAt,omitempty"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

func (w *WhiteLabel) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// CryptoWallet represents multi-chain wallet
type CryptoWallet struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID     `gorm:"type:uuid;index" json:"userId"`
	User          User          `gorm:"foreignKey:UserID" json:"-"`
	Chain         string        `gorm:"size:20;not null" json:"chain"` // btc, eth, bnb, sol, trx, etc.
	Address       string        `gorm:"size:100;not null" json:"address"`
	PrivateKeyEnc string        `gorm:"size:255" json:"-"` // Encrypted private key
	IsPrimary     bool          `gorm:"default:false" json:"isPrimary"`
	Balance       float64       `gorm:"default:0" json:"balance"`
	Network       string        `gorm:"size:50" json:"network"` // mainnet, testnet
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

func (c *CryptoWallet) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// SupportedChains lists all supported blockchain networks
var SupportedChains = []string{"BTC", "ETH", "BNB", "SOL", "TRX", "USDT", "USDC", "DAI", "MATIC", "AVAX", "FTM"}

// ChainInfo contains details about each supported chain
var ChainInfo = map[string]struct {
	Name       string
	Symbol    string
	Decimals  int
	Explorer  string
	MinDeposit float64
}{
	"BTC":   {Name: "Bitcoin", Symbol: "BTC", Decimals: 8, Explorer: "https://blockstream.info", MinDeposit: 0.0001},
	"ETH":   {Name: "Ethereum", Symbol: "ETH", Decimals: 18, Explorer: "https://etherscan.io", MinDeposit: 0.001},
	"BNB":   {Name: "BNB Chain", Symbol: "BNB", Decimals: 18, Explorer: "https://bscscan.com", MinDeposit: 0.01},
	"SOL":   {Name: "Solana", Symbol: "SOL", Decimals: 9, Explorer: "https://solscan.io", MinDeposit: 0.01},
	"TRX":   {Name: "Tron", Symbol: "TRX", Decimals: 6, Explorer: "https://tronscan.org", MinDeposit: 10},
	"USDT":  {Name: "Tether USDT", Symbol: "USDT", Decimals: 6, Explorer: "https://tether.to", MinDeposit: 10},
	"USDC":  {Name: "USD Coin", Symbol: "USDC", Decimals: 6, Explorer: "https://usdc.co", MinDeposit: 10},
	"DAI":   {Name: "Dai", Symbol: "DAI", Decimals: 18, Explorer: "https://dai Stablecoin.com", MinDeposit: 10},
	"MATIC": {Name: "Polygon", Symbol: "MATIC", Decimals: 18, Explorer: "https://polygonscan.com", MinDeposit: 1},
	"AVAX":  {Name: "Avalanche", Symbol: "AVAX", Decimals: 18, Explorer: "https://snowtrace.io", MinDeposit: 0.1},
	"FTM":   {Name: "Fantom", Symbol: "FTM", Decimals: 18, Explorer: "https://ftmscan.com", MinDeposit: 1},
}

// Game represents a casino game
type Game struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Type        string         `gorm:"size:50;not null" json:"type"` // slots, dice, roulette, blackjack, baccarat
	Provider    string         `gorm:"size:50" json:"provider"`
	RTP         float64        `json:"rtp"` // Return to Player percentage
	MinBet      float64        `json:"minBet"`
	MaxBet      float64        `json:"maxBet"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	ThumbnailURL string       `json:"thumbnailUrl,omitempty"`
	GameData    string         `gorm:"type:text" json:"gameData,omitempty"` // JSON configuration
	CreatedAt   time.Time      `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// Bet represents a bet placed by a user
type Bet struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"userId"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	GameID      uuid.UUID      `gorm:"type:uuid;index" json:"gameId"`
	Game        Game           `gorm:"foreignKey:GameID" json:"-"`
	BetAmount   float64        `gorm:"not null" json:"betAmount"`
	WinAmount   float64        `gorm:"default:0" json:"winAmount"`
	Multiplier  float64        `gorm:"default:0" json:"multiplier"`
	GameData    string         `gorm:"type:text" json:"gameData,omitempty"` // Game-specific data (deck, dice roll, etc.)
	Status      string         `gorm:"size:20;default:pending" json:"status"` // pending, won, lost
	SettledAt   *time.Time    `json:"settledAt,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Bet) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// Session represents an active user session
type Session struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;index" json:"userId"`
	User         User           `gorm:"foreignKey:UserID" json:"-"`
	Token        string         `gorm:"size:500;not null;uniqueIndex" json:"token"`
	IPAddress    string         `gorm:"size:45" json:"ipAddress"`
	UserAgent    string         `gorm:"size:500" json:"userAgent"`
	ExpiresAt    time.Time      `json:"expiresAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    *uuid.UUID    `gorm:"type:uuid" json:"userId,omitempty"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Action    string         `gorm:"size:100;not null" json:"action"`
	Details   string         `gorm:"type:text" json:"details,omitempty"` // JSON
	IPAddress string         `gorm:"size:45" json:"ipAddress"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
