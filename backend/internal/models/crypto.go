package models

import (
	"time"

	"github.com/google/uuid"
)

// CryptoNetwork represents a blockchain network
// GORM model stub - actual migration from schema.sql
type CryptoNetwork struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name                string    `json:"name" gorm:"type:varchar(50);not null"`
	ChainID             string    `json:"chain_id" gorm:"type:varchar(20)"`
	Symbol              string    `json:"symbol" gorm:"type:varchar(20);not null"`
	ExplorerURL         string    `json:"explorer_url" gorm:"type:varchar(255)"`
	RPCURL              string    `json:"rpc_url" gorm:"type:text"`
	IsWithdrawalEnabled bool      `json:"is_withdrawal_enabled" gorm:"default:true"`
	IsDepositEnabled    bool      `json:"is_deposit_enabled" gorm:"default:true"`
	IsActive            bool      `json:"is_active" gorm:"default:true"`
	MinConfirmationBlocks int     `json:"min_confirmation_blocks" gorm:"default:6"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CryptoAsset represents a cryptocurrency
type CryptoAsset struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name                string    `json:"name" gorm:"type:varchar(100);not null"`
	Symbol              string    `json:"symbol" gorm:"type:varchar(20);not null;uniqueIndex"`
	Decimals            int       `json:"decimals" gorm:"default:18"`
	ContractAddress     string    `json:"contract_address" gorm:"type:varchar(100)"`
	IsActive            bool      `json:"is_active" gorm:"default:true"`
	MinDepositAmount    float64   `json:"min_deposit_amount" gorm:"type:decimal(20,8)"`
	MinWithdrawalAmount float64   `json:"min_withdrawal_amount" gorm:"type:decimal(20,8)"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CryptoAssetNetwork represents the many-to-many relationship between assets and networks
type CryptoAssetNetwork struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	AssetID           uuid.UUID `json:"asset_id" gorm:"type:uuid;not null"`
	NetworkID         uuid.UUID `json:"network_id" gorm:"type:uuid;not null"`
	DepositEnabled    bool      `json:"deposit_enabled" gorm:"default:true"`
	WithdrawalEnabled bool      `json:"withdrawal_enabled" gorm:"default:true"`
	ContractAddress   string    `json:"contract_address" gorm:"type:varchar(100)"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Asset   CryptoAsset   `json:"asset,omitempty" gorm:"foreignKey:AssetID"`
	Network CryptoNetwork `json:"network,omitempty" gorm:"foreignKey:NetworkID"`
}

// NetworkFee represents fees for deposit/withdrawal per asset per network
type NetworkFee struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	AssetID             uuid.UUID `json:"asset_id" gorm:"type:uuid;not null"`
	NetworkID           uuid.UUID `json:"network_id" gorm:"type:uuid;not null"`
	DepositFee          float64   `json:"deposit_fee" gorm:"type:decimal(20,8);default:0"`
	WithdrawalFee       float64   `json:"withdrawal_fee" gorm:"type:decimal(20,8);default:0"`
	DepositFeePercent   float64   `json:"deposit_fee_percent" gorm:"type:decimal(5,4);default:0"`
	WithdrawalFeePercent float64  `json:"withdrawal_fee_percent" gorm:"type:decimal(5,4);default:0"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	Asset   CryptoAsset   `json:"asset,omitempty" gorm:"foreignKey:AssetID"`
	Network CryptoNetwork `json:"network,omitempty" gorm:"foreignKey:NetworkID"`
}

// BrandLevel represents VIP/brand levels with fee discounts
type BrandLevel struct {
	ID                          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name                        string    `json:"name" gorm:"type:varchar(50);not null"`
	Level                       int       `json:"level" gorm:"not null;uniqueIndex"`
	DepositFeeDiscountPercent   float64   `json:"deposit_fee_discount_percent" gorm:"type:decimal(5,4);default:0"`
	WithdrawalFeeDiscountPercent float64  `json:"withdrawal_fee_discount_percent" gorm:"type:decimal(5,4);default:0"`
	IsActive                    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

// CryptoDeposit represents a crypto deposit transaction
type CryptoDeposit struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	AssetID     uuid.UUID `json:"asset_id" gorm:"type:uuid;not null"`
	NetworkID   uuid.UUID `json:"network_id" gorm:"type:uuid;not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(20,8);not null"`
	Fee         float64   `json:"fee" gorm:"type:decimal(20,8);default:0"`
	NetAmount   float64   `json:"net_amount" gorm:"type:decimal(20,8);not null"`
	Address     string    `json:"address" gorm:"type:varchar(100);not null"`
	TxHash      string    `json:"tx_hash" gorm:"type:varchar(100)"`
	Confirmations int     `json:"confirmations" gorm:"default:0"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`

	User    User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Asset   CryptoAsset   `json:"asset,omitempty" gorm:"foreignKey:AssetID"`
	Network CryptoNetwork `json:"network,omitempty" gorm:"foreignKey:NetworkID"`
}

// CryptoWithdrawal represents a crypto withdrawal transaction
type CryptoWithdrawal struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	AssetID     uuid.UUID `json:"asset_id" gorm:"type:uuid;not null"`
	NetworkID   uuid.UUID `json:"network_id" gorm:"type:uuid;not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(20,8);not null"`
	Fee         float64   `json:"fee" gorm:"type:decimal(20,8);default:0"`
	NetAmount   float64   `json:"net_amount" gorm:"type:decimal(20,8);not null"`
	Address     string    `json:"address" gorm:"type:varchar(100);not null"`
	TxHash      string    `json:"tx_hash" gorm:"type:varchar(100)"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`

	User    User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Asset   CryptoAsset   `json:"asset,omitempty" gorm:"foreignKey:AssetID"`
	Network CryptoNetwork `json:"network,omitempty" gorm:"foreignKey:NetworkID"`
}

// AdminWalletAddress represents admin wallet addresses for each asset/network
type AdminWalletAddress struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	AssetID           uuid.UUID `json:"asset_id" gorm:"type:uuid;not null"`
	NetworkID         uuid.UUID `json:"network_id" gorm:"type:uuid;not null"`
	Address           string    `json:"address" gorm:"type:varchar(100);not null"`
	PrivateKeyEncrypted string  `json:"private_key_encrypted" gorm:"type:text"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Asset   CryptoAsset   `json:"asset,omitempty" gorm:"foreignKey:AssetID"`
	Network CryptoNetwork `json:"network,omitempty" gorm:"foreignKey:NetworkID"`
}

// DepositRequest represents a user's deposit request
type DepositRequest struct {
	AssetID   uuid.UUID `json:"asset_id" binding:"required"`
	NetworkID uuid.UUID `json:"network_id" binding:"required"`
}

// WithdrawalRequest represents a user's withdrawal request
type WithdrawalRequest struct {
	AssetID   uuid.UUID `json:"asset_id" binding:"required"`
	NetworkID uuid.UUID `json:"network_id" binding:"required"`
	Amount    float64   `json:"amount" binding:"required,gt=0"`
	Address   string    `json:"address" binding:"required"`
}

// NetworkWithAsset represents a network with its associated assets
type NetworkWithAsset struct {
	CryptoNetwork
	Assets []CryptoAssetWithNetwork `json:"assets"`
}

// CryptoAssetWithNetwork represents an asset with its network info
type CryptoAssetWithNetwork struct {
	CryptoAsset
	NetworkID         uuid.UUID `json:"network_id"`
	DepositEnabled    bool      `json:"deposit_enabled"`
	WithdrawalEnabled bool      `json:"withdrawal_enabled"`
	ContractAddress   string    `json:"contract_address"`
	DepositFee        float64   `json:"deposit_fee"`
	WithdrawalFee     float64   `json:"withdrawal_fee"`
	DepositFeePercent   float64 `json:"deposit_fee_percent"`
	WithdrawalFeePercent float64 `json:"withdrawal_fee_percent"`
	MinDepositAmount  float64   `json:"min_deposit_amount"`
	MinWithdrawalAmount float64 `json:"min_withdrawal_amount"`
}
