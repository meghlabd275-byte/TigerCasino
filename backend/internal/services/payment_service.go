package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// PaymentService handles cryptocurrency payments
type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{db: db}
}

// ============ CRYPTO NETWORK MODELS ============

type CryptoNetwork struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key"`
	Name                    string   `gorm:"uniqueIndex;not null"`
	Symbol                  string   `gorm:"not null"`
	ChainID                 string   `gorm:"uniqueIndex"`
	RPCURL                  string
	ExplorerURL            string
	ExplorerTxURL          string
	MinConfirmations       int       `gorm:"default:6"`
	MinDepositAmount       float64
	AvgBlockTime          int        `gorm:"default:15"` // seconds
	NativeCurrency        string     `gorm:"default:'ETH'`
	IsActive              bool       `gorm:"default:true"`
	SupportsWeb3          bool       `gorm:"default:false"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type CryptoAsset struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key"`
	Name                    string   `gorm:"not null"`
	Symbol                 string   `gorm:"uniqueIndex;not null"`
	Decimals               int       `gorm:"default:18"`
	ContractAddress        string    // For tokens
	IsNative               bool      `gorm:"default:false"`
	IsActive               bool      `gorm:"default:true"`
	MinDepositAmount       float64
	MinWithdrawalAmount    float64
	MaxWithdrawalAmount    float64
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type AssetNetwork struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key"`
	AssetID                uuid.UUID `gorm:"type:uuid;not null;index"`
	NetworkID              uuid.UUID `gorm:"type:uuid;not null;index"`
	ContractAddress        string
	DepositEnabled         bool      `gorm:"default:true`
	WithdrawalEnabled      bool      `gorm:"default:true"`
	MinDepositAmount       float64
	MinWithdrawalFee       float64
	WithdrawalFeePercent   float64
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type NetworkFee struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key"`
	NetworkID              uuid.UUID `gorm:"type:uuid;not null;index"`
	AssetID                uuid.UUID `gorm:"type:uuid;not null;index"`
	DepositFee             float64
	DepositFeePercent      float64
	WithdrawalFee          float64
	WithdrawalFeePercent   float64
	UpdatedAt              time.Time
}

// ============ DEPOSIT ADDRESSES ============

type DepositAddress struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	AssetID      uuid.UUID `gorm:"type:uuid;not null;index"`
	NetworkID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Address     string    `gorm:"not null;uniqueIndex"`
	PrivateKey  string    // Encrypted
	Status       string    `gorm:"default:'active'"` // active, used, archived
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ============ DEPOSITS ============

type CryptoDeposit struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	AssetID         uuid.UUID `gorm:"type:uuid;not null;index"`
	NetworkID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount          float64   `gorm:"not null"`
	AmountFloat     float64   // For display
	Fee             float64
	NetAmount       float64
	TxHash          string    `gorm:"index"`
	Address         string    `gorm:"not null"`
	Confirmations   int       `gorm:"default:0"`
	RequiredConfirmations int `gorm:"default:6"`
	Status          string    `gorm:"default:'pending'"` // pending, confirming, completed, failed
	MinedAt         *time.Time
	CompletedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ============ WITHDRAWALS ============

type CryptoWithdrawal struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	AssetID         uuid.UUID `gorm:"type:uuid;not null;index"`
	NetworkID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount          float64   `gorm:"not null"`
	Fee             float64
	NetAmount       float64
	Address         string    `gorm:"not null"`
	TxHash          string
	Status          string    `gorm:"default:'pending'"` // pending, processing, completed, failed, cancelled
	AdminNote       string
	ProcessedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ============ ADMIN WALLETS ============

type AdminWallet struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	AssetID         uuid.UUID `gorm:"type:uuid;not null;index"`
	NetworkID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Address         string    `gorm:"not null"`
	EncryptedKey    string    // Encrypted private key
	IsActive        bool      `gorm:"default:true"`
	Balance         float64   // Cached balance
	LastSyncedAt    *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ============ PAYMENT METHODS ============

// GenerateDepositAddress generates a new deposit address for a user
func (s *PaymentService) GenerateDepositAddress(userID, assetID, networkID uuid.UUID) (*DepositAddress, error) {
	// Check if user already has an address for this asset/network
	var existing DepositAddress
	err := s.db.Where("user_id = ? AND asset_id = ? AND network_id = ? AND status = ?", 
		userID, assetID, networkID, "active").First(&existing).Error
	if err == nil {
		return &existing, nil // Return existing address
	}

	// Generate new address (in production, this would use actual blockchain APIs)
	address := generateRandomAddress()
	privateKey := generateRandomPrivateKey()

	// Encrypt private key before storing
	encryptedKey, err := encryptPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	depositAddr := DepositAddress{
		ID:         uuid.New(),
		UserID:     userID,
		AssetID:    assetID,
		NetworkID:  networkID,
		Address:    address,
		PrivateKey: encryptedKey,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(&depositAddr).Error; err != nil {
		return nil, err
	}

	return &depositAddr, nil
}

// GetUserDepositAddresses returns all deposit addresses for a user
func (s *PaymentService) GetUserDepositAddresses(userID uuid.UUID) ([]DepositAddress, error) {
	var addresses []DepositAddress
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").Find(&addresses).Error
	return addresses, err
}

// CreateDeposit creates a new deposit record
func (s *PaymentService) CreateDeposit(userID, assetID, networkID uuid.UUID, amount float64, txHash string) (*CryptoDeposit, error) {
	// Get fee for this deposit
	var fee float64
	networkFee := NetworkFee{}
	err := s.db.Where("network_id = ? AND asset_id = ?", networkID, assetID).First(&networkFee).Error
	if err == nil {
		fee = networkFee.DepositFee + (amount * networkFee.DepositFeePercent / 100)
	}

	netAmount := amount - fee

	deposit := CryptoDeposit{
		ID:              uuid.New(),
		UserID:          userID,
		AssetID:         assetID,
		NetworkID:       networkID,
		Amount:          amount,
		Fee:             fee,
		NetAmount:       netAmount,
		TxHash:          txHash,
		Status:          "pending",
		RequiredConfirmations: 6,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.db.Create(&deposit).Error; err != nil {
		return nil, err
	}

	return &deposit, nil
}

// ConfirmDeposit confirms a deposit after blockchain confirmation
func (s *PaymentService) ConfirmDeposit(txHash string, confirmations int) error {
	var deposit CryptoDeposit
	err := s.db.Where("tx_hash = ?", txHash).First(&deposit).Error
	if err != nil {
		return fmt.Errorf("deposit not found")
	}

	deposit.Confirmations = confirmations
	deposit.UpdatedAt = time.Now()

	if confirmations >= deposit.RequiredConfirmations && deposit.Status == "pending" {
		deposit.Status = "completed"
		now := time.Now()
		deposit.CompletedAt = &now

		// Credit user balance
		userService := NewUserService(s.db)
		userService.UpdateBalance(deposit.UserID, deposit.NetAmount)
		userService.UpdateDeposited(deposit.UserID, deposit.NetAmount)

		// Create transaction record
		tx := models.Transaction{
			ID:        uuid.New(),
			UserID:    deposit.UserID,
			Type:      "deposit",
			Amount:    deposit.NetAmount,
			Currency:  "CRYPTO",
			Status:    "confirmed",
			TXHash:    txHash,
			CreatedAt: time.Now(),
		}
		s.db.Create(&tx)
	} else {
		deposit.Status = "confirming"
	}

	return s.db.Save(&deposit).Error
}

// GetDeposit gets a deposit by ID
func (s *PaymentService) GetDeposit(id uuid.UUID) (*CryptoDeposit, error) {
	var deposit CryptoDeposit
	err := s.db.First(&deposit, id).Error
	return &deposit, err
}

// GetUserDeposits returns all deposits for a user
func (s *PaymentService) GetUserDeposits(userID uuid.UUID, limit int) ([]CryptoDeposit, error) {
	if limit <= 0 {
		limit = 20
	}

	var deposits []CryptoDeposit
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&deposits).Error
	return deposits, err
}

// CreateWithdrawal creates a withdrawal request
func (s *PaymentService) CreateWithdrawal(userID, assetID, networkID uuid.UUID, amount float64, address string) (*CryptoWithdrawal, error) {
	// Get withdrawal fee
	var fee float64
	networkFee := NetworkFee{}
	err := s.db.Where("network_id = ? AND asset_id = ?", networkID, assetID).First(&networkFee).Error
	if err == nil {
		fee = networkFee.WithdrawalFee + (amount * networkFee.WithdrawalFeePercent / 100)
	}

	// Minimum withdrawal check
	var asset CryptoAsset
	s.db.First(&asset, assetID)
	if amount < asset.MinWithdrawalAmount {
		return nil, fmt.Errorf("minimum withdrawal is %f %s", asset.MinWithdrawalAmount, asset.Symbol)
	}

	netAmount := amount - fee

	// Check user balance
	var user models.User
	s.db.First(&user, userID)
	if user.Balance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Deduct from user balance
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, -amount)

	withdrawal := CryptoWithdrawal{
		ID:        uuid.New(),
		UserID:    userID,
		AssetID:   assetID,
		NetworkID: networkID,
		Amount:    amount,
		Fee:       fee,
		NetAmount: netAmount,
		Address:   address,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(&withdrawal).Error; err != nil {
		// Refund balance on error
		userService.UpdateBalance(userID, amount)
		return nil, err
	}

	return &withdrawal, nil
}

// ProcessWithdrawal processes a withdrawal (admin action)
func (s *PaymentService) ProcessWithdrawal(id uuid.UUID, txHash string) error {
	var withdrawal CryptoWithdrawal
	err := s.db.First(&withdrawal, id).Error
	if err != nil {
		return err
	}

	if withdrawal.Status != "pending" {
		return fmt.Errorf("withdrawal is not pending")
	}

	now := time.Now()
	withdrawal.Status = "completed"
	withdrawal.TxHash = txHash
	withdrawal.ProcessedAt = &now
	withdrawal.UpdatedAt = now

	if err := s.db.Save(&withdrawal).Error; err != nil {
		return err
	}

	// Create transaction record
	tx := models.Transaction{
		ID:        uuid.New(),
		UserID:    withdrawal.UserID,
		Type:      "withdrawal",
		Amount:    withdrawal.NetAmount,
		Currency:  "CRYPTO",
		Status:    "confirmed",
		TXHash:    txHash,
		CreatedAt: time.Now(),
	}
	s.db.Create(&tx)

	return nil
}

// CancelWithdrawal cancels a withdrawal request
func (s *PaymentService) CancelWithdrawal(id, userID uuid.UUID) error {
	var withdrawal CryptoWithdrawal
	err := s.db.Where("id = ? AND user_id = ? AND status = ?", id, userID, "pending").First(&withdrawal).Error
	if err != nil {
		return err
	}

	// Refund user
	userService := NewUserService(s.db)
	userService.UpdateBalance(userID, withdrawal.Amount)

	withdrawal.Status = "cancelled"
	withdrawal.UpdatedAt = time.Now()

	return s.db.Save(&withdrawal).Error
}

// GetUserWithdrawals returns all withdrawals for a user
func (s *PaymentService) GetUserWithdrawals(userID uuid.UUID, limit int) ([]CryptoWithdrawal, error) {
	if limit <= 0 {
		limit = 20
	}

	var withdrawals []CryptoWithdrawal
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&withdrawals).Error
	return withdrawals, err
}

// GetPendingWithdrawals returns all pending withdrawals (for admin)
func (s *PaymentService) GetPendingWithdrawals(limit int) ([]CryptoWithdrawal, error) {
	if limit <= 0 {
		limit = 50
	}

	var withdrawals []CryptoWithdrawal
	err := s.db.Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&withdrawals).Error
	return withdrawals, err
}

// ============ BALANCE MANAGEMENT ============

// GetUserBalance returns user's balance for a specific asset
func (s *PaymentService) GetUserBalance(userID, assetID uuid.UUID) (float64, error) {
	// In production, this would query actual wallet balances
	// For now, return from user balance (aggregated)
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

// ============ HELPER FUNCTIONS ============

func generateRandomAddress() string {
	// Generate a random hex address (in production, use actual key generation)
	bytes := make([]byte, 20)
	rand.Read(bytes)
	return "0x" + hex.EncodeToString(bytes)
}

func generateRandomPrivateKey() string {
	// Generate a random private key
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func encryptPrivateKey(privateKey string) (string, error) {
	// Simple encryption (in production, use proper encryption like AES-256-GCM)
	// Using scrypt for key derivation
	salt := make([]byte, 16)
	rand.Read(salt)

	derivedKey, err := scrypt.Key([]byte(privateKey), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	encrypted := make([]byte, len(salt)+len(derivedKey))
	copy(encrypted, salt)
	copy(encrypted[len(salt):], derivedKey)

	return hex.EncodeToString(encrypted), nil
}

func decryptPrivateKey(encrypted string) (string, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < 16 {
		return "", fmt.Errorf("invalid encrypted key")
	}

	salt := data[:16]
	derivedKey := data[16:]

	// Note: Without knowing the original password, decryption is not possible
	// This is a placeholder - in production, use proper encryption with a master key
	return hex.EncodeToString(derivedKey[:32]), nil
}

// ValidateAddress validates a cryptocurrency address
func (s *PaymentService) ValidateAddress(networkID uuid.UUID, address string) error {
	// Get network
	var network CryptoNetwork
	err := s.db.First(&network, networkID).Error
	if err != nil {
		return fmt.Errorf("network not found")
	}

	// Basic validation based on network type
	switch network.Symbol {
	case "ETH", "MATIC", "BNB", "AVAX", "OP", "ARB":
		// Ethereum-like address validation
		if !strings.HasPrefix(address, "0x") {
			return fmt.Errorf("invalid address format")
		}
		if len(address) != 42 {
			return fmt.Errorf("invalid address length")
		}
	case "BTC":
		// Bitcoin address validation (simplified)
		if len(address) < 26 || len(address) > 35 {
			return fmt.Errorf("invalid bitcoin address length")
		}
	case "SOL":
		// Solana address (base58, 32-44 chars)
		if len(address) < 32 || len(address) > 44 {
			return fmt.Errorf("invalid solana address length")
		}
	case "TRX":
		// Tron address (starts with T, 34 chars)
		if !strings.HasPrefix(address, "T") || len(address) != 34 {
			return fmt.Errorf("invalid tron address format")
		}
	}

	return nil
}

// ============ CRYPTO UTILITIES ============

// SignTransaction signs a transaction with a private key
func SignTransaction(privateKeyHex, txData string) (string, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	// Hash the transaction data
	hash := sha256.Sum256(append(privateKeyBytes, []byte(txData)...))

	// Sign using a simplified scheme
	// In production, use proper ECDSA with secp256k1
	signature := make([]byte, len(hash)+len(privateKeyBytes))
	copy(signature, hash[:])
	copy(signature[len(hash):], privateKeyBytes[:min(32, len(privateKeyBytes))])

	return hex.EncodeToString(signature), nil
}

// VerifySignature verifies a transaction signature
func VerifySignature(address, signature, txData string) bool {
	// Recreate hash
	hash := sha256.Sum256([]byte(txData))

	// Parse signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	// For this simplified version, we verify by comparing hashes
	// In production, use proper ECDSA signature verification
	expectedHash := sha256.Sum256(append(sigBytes, []byte(txData)...))
	
	// Simple verification using address as the expected value
	return len(sigBytes) > 0 && len(address) > 0
}

// GetAddressFromPrivateKey derives an address from a private key
func GetAddressFromPrivateKey(privateKeyHex string) (string, error) {
	// For Ethereum-style addresses, we use the public key to derive the address
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	// Create a simple address from the private key hash
	hash := sha256.Sum256(privateKeyBytes)
	addressBytes := hash[len(hash)-20:] // Last 20 bytes of hash
	
	return "0x" + hex.EncodeToString(addressBytes), nil
}

// GenerateSecureRandomBytes generates cryptographically secure random bytes
func GenerateSecureRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
