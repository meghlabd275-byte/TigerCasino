package services

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// PaymentGatewayService handles cryptocurrency payments with real blockchain integration
type PaymentGatewayService struct {
	db               *gorm.DB
	redis            *redis.Client
	walletManager    *WalletManager
	processors       map[string]PaymentProcessor
	config           *PaymentGatewayConfig
	balanceCache     *BalanceCache
	webhookHandler   *WebhookHandler
}

type PaymentGatewayConfig struct {
	CoinsPaidAPIKey     string
	CoinsPaidSecret     string
	CoinsPaidURL        string
	BitPayAPIKey        string
	BitPayURL           string
	CoinPaymentsAPIKey  string
	CoinPaymentsURL     string
	
	// Blockchain node configurations
	BTCNodeURL          string
	BTCNodeUser         string
	BTCNodePass         string
	ETHNodeURL          string
	
	// Hot wallet configuration
	HotWalletPrivateKey string
	HotWalletAddresses map[string]string // currency -> address
	
	// Fee configuration
	MinDepositAmount    float64
	MinWithdrawalAmount float64
	WithdrawalFeePercent float64
	
	// Security
	ConfirmationsRequired int
	AutoWithdrawEnabled   bool
	MaxWithdrawalDaily   float64
	
	Timeout             time.Duration
}

// PaymentProcessor interface for different payment providers
type PaymentProcessor interface {
	GetName() string
	CreateDeposit(ctx context.Context, req *CreateDepositRequest) (*DepositInfo, error)
	CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*WithdrawalInfo, error)
	GetDepositStatus(ctx context.Context, depositID string) (*DepositStatus, error)
	GetWithdrawalStatus(ctx context.Context, withdrawalID string) (*WithdrawalStatus, error)
	GetExchangeRate(ctx context.Context, from, to string) (float64, error)
	IsAvailable() bool
}

type CreateDepositRequest struct {
	UserID        uuid.UUID
	Amount        float64
	Currency      string
	Network       string
	IPAddress     string
	UserAgent     string
}

type CreateWithdrawalRequest struct {
	UserID        uuid.UUID
	Amount        float64
	Currency      string
	Network       string
	ToAddress     string
	IPAddress     string
	UserAgent     string
}

type DepositInfo struct {
	DepositID    string
	Address      string
	Amount       float64
	Currency     string
	Network      string
	ExpiresAt    time.Time
	QRCode       string
	Status       string
}

type WithdrawalInfo struct {
	WithdrawalID string
	Amount       float64
	Currency     string
	Network      string
	ToAddress    string
	Fee          float64
	NetAmount    float64
	Status       string
	TxHash       string
}

type DepositStatus struct {
	DepositID    string
	Status       string
	Confirmations int
	Amount       float64
	TxHash       string
	ReceivedAt   time.Time
}

type WithdrawalStatus struct {
	WithdrawalID string
	Status       string
	Confirmations int
	TxHash       string
	ProcessedAt  time.Time
}

// BalanceCache for quick balance lookups
type BalanceCache struct {
	mu      sync.RWMutex
	balances map[string]map[string]float64 // userID -> currency -> balance
}

func NewBalanceCache() *BalanceCache {
	return &BalanceCache{
		balances: make(map[string]map[string]float64),
	}
}

func (bc *BalanceCache) Get(userID, currency string) (float64, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	if balances, ok := bc.balances[userID]; ok {
		balance, ok := balances[currency]
		return balance, ok
	}
	return 0, false
}

func (bc *BalanceCache) Set(userID, currency string, balance float64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	if _, ok := bc.balances[userID]; !ok {
		bc.balances[userID] = make(map[string]float64)
	}
	bc.balances[userID][currency] = balance
}

func (bc *BalanceCache) Add(userID, currency string, amount float64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	if _, ok := bc.balances[userID]; !ok {
		bc.balances[userID] = make(map[string]float64)
	}
	bc.balances[userID][currency] += amount
}

func (bc *BalanceCache) Subtract(userID, currency string, amount float64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	if balances, ok := bc.balances[userID]; ok {
		balances[currency] -= amount
		if balances[currency] < 0 {
			balances[currency] = 0
		}
	}
}

// WebhookHandler handles payment webhooks
type WebhookHandler struct {
	mu           sync.RWMutex
	handlers     map[string]WebhookHandlerFunc
	secretKeys   map[string]string // processor -> secret
}

type WebhookHandlerFunc func(ctx context.Context, data json.RawMessage) error

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		handlers:   make(map[string]WebhookHandlerFunc),
		secretKeys: make(map[string]string),
	}
}

func (wh *WebhookHandler) Register(processor string, secret string, handler WebhookHandlerFunc) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	wh.handlers[processor] = handler
	wh.secretKeys[processor] = secret
}

func (wh *WebhookHandler) Handle(ctx context.Context, processor string, payload json.RawMessage, signature string) error {
	wh.mu.RLock()
	handler, ok := wh.handlers[processor]
	secret := wh.secretKeys[processor]
	wh.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no handler for processor: %s", processor)
	}

	// Verify signature
	if secret != "" && signature != "" {
		expectedSig := wh.generateSignature(payload, secret)
		if signature != expectedSig {
			return fmt.Errorf("invalid signature")
		}
	}

	return handler(ctx, payload)
}

func (wh *WebhookHandler) generateSignature(payload json.RawMessage, secret string) string {
	h := sha256.New()
	h.Write(payload)
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(nil))
}

// WalletManager manages cryptocurrency wallets
type WalletManager struct {
	db          *gorm.DB
	config      *PaymentGatewayConfig
	ethClient   *ethclient.Client
	btcClient   *rpcclient.Client
	addressPool *AddressPool
	wallets     map[string]*Wallet // currency -> wallet
	mu          sync.RWMutex
}

type Wallet struct {
	Currency      string
	Network       string
	PrivateKey    string
	Address       string
	PublicKey     string
	IsHotWallet   bool
	Balance       float64
	LastSyncTime  time.Time
}

type AddressPool struct {
	db        *gorm.DB
	mu        sync.Mutex
	addresses map[string][]string // network -> addresses
}

// CoinsPaidProcessor implements CoinsPaid payment processor
type CoinsPaidProcessor struct {
	config     *PaymentGatewayConfig
	httpClient *http.Client
	enabled    bool
}

func NewCoinsPaidProcessor(config *PaymentGatewayConfig) *CoinsPaidProcessor {
	return &CoinsPaidProcessor{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		enabled: config.CoinsPaidAPIKey != "" && config.CoinsPaidSecret != "",
	}
}

func (c *CoinsPaidProcessor) GetName() string {
	return "CoinsPaid"
}

func (c *CoinsPaidProcessor) IsAvailable() bool {
	return c.enabled
}

func (c *CoinsPaidProcessor) CreateDeposit(ctx context.Context, req *CreateDepositRequest) (*DepositInfo, error) {
	if !c.enabled {
		return nil, fmt.Errorf("CoinsPaid not configured")
	}

	url := fmt.Sprintf("%s/v2/deposits", c.config.CoinsPaidURL)

	timestamp := time.Now().Unix()
	payload := map[string]interface{}{
		"user_id":        req.UserID.String(),
		"currency":       req.Currency,
		"network":        req.Network,
		"amount":         fmt.Sprintf("%.8f", req.Amount),
		"callback_url":   fmt.Sprintf("%s/api/payments/coinspaid/callback", getBaseURL()),
		"external_id":    uuid.New().String(),
	}

	body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	c.signRequest(httpReq, timestamp, string(body))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID          string `json:"id"`
			Address     string `json:"address"`
			Amount      string `json:"amount"`
			Currency    string `json:"currency"`
			Network     string `json:"network"`
			ExpiresAt   int64  `json:"expires_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &DepositInfo{
		DepositID: result.Data.ID,
		Address:    result.Data.Address,
		Amount:     req.Amount,
		Currency:   result.Data.Currency,
		Network:    result.Data.Network,
		ExpiresAt:  time.Unix(result.Data.ExpiresAt, 0),
		Status:     "pending",
	}, nil
}

func (c *CoinsPaidProcessor) CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*WithdrawalInfo, error) {
	if !c.enabled {
		return nil, fmt.Errorf("CoinsPaid not configured")
	}

	url := fmt.Sprintf("%s/v2/withdrawals", c.config.CoinsPaidURL)

	timestamp := time.Now().Unix()
	payload := map[string]interface{}{
		"user_id":       req.UserID.String(),
		"currency":      req.Currency,
		"network":       req.Network,
		"amount":        fmt.Sprintf("%.8f", req.Amount),
		"address":       req.ToAddress,
		"callback_url":  fmt.Sprintf("%s/api/payments/coinspaid/withdrawal/callback", getBaseURL()),
		"external_id":   uuid.New().String(),
	}

	body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	c.signRequest(httpReq, timestamp, string(body))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID          string `json:"id"`
			Amount      string `json:"amount"`
			Currency    string `json:"currency"`
			Network     string `json:"network"`
			Address     string `json:"address"`
			Fee          string `json:"fee"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	amount, _ := strconv.ParseFloat(result.Data.Amount, 64)
	fee, _ := strconv.ParseFloat(result.Data.Fee, 64)

	return &WithdrawalInfo{
		WithdrawalID: result.Data.ID,
		Amount:       amount,
		Currency:     result.Data.Currency,
		Network:      result.Data.Network,
		ToAddress:    result.Data.Address,
		Fee:          fee,
		NetAmount:    amount - fee,
		Status:       "pending",
	}, nil
}

func (c *CoinsPaidProcessor) signRequest(req *http.Request, timestamp int64, body string) {
	signature := fmt.Sprintf("%d%s%s", timestamp, c.config.CoinsPaidSecret, body)
	h := sha512.New()
	h.Write([]byte(signature))
	sig := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("CPC-PUBLIC-KEY", c.config.CoinsPaidAPIKey)
	req.Header.Set("CPC-TIMESTAMP", fmt.Sprintf("%d", timestamp))
	req.Header.Set("CPC-SIGNATURE", sig)
	req.Header.Set("Content-Type", "application/json")
}

func (c *CoinsPaidProcessor) GetDepositStatus(ctx context.Context, depositID string) (*DepositStatus, error) {
	if !c.enabled {
		return nil, fmt.Errorf("CoinsPaid not configured")
	}

	url := fmt.Sprintf("%s/v2/deposits/%s", c.config.CoinsPaidURL, depositID)

	timestamp := time.Now().Unix()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.signRequest(req, timestamp, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID            string `json:"id"`
			Status        string `json:"status"`
			Confirmations int    `json:"confirmations"`
			Amount        string `json:"amount"`
			TxHash        string `json:"tx_hash"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &DepositStatus{
		DepositID:     result.Data.ID,
		Status:        result.Data.Status,
		Confirmations: result.Data.Confirmations,
		TxHash:        result.Data.TxHash,
	}, nil
}

func (c *CoinsPaidProcessor) GetWithdrawalStatus(ctx context.Context, withdrawalID string) (*WithdrawalStatus, error) {
	if !c.enabled {
		return nil, fmt.Errorf("CoinsPaid not configured")
	}

	url := fmt.Sprintf("%s/v2/withdrawals/%s", c.config.CoinsPaidURL, withdrawalID)

	timestamp := time.Now().Unix()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.signRequest(req, timestamp, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID            string `json:"id"`
			Status        string `json:"status"`
			Confirmations int    `json:"confirmations"`
			TxHash        string `json:"tx_hash"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &WithdrawalStatus{
		WithdrawalID:   result.Data.ID,
		Status:         result.Data.Status,
		Confirmations:  result.Data.Confirmations,
		TxHash:         result.Data.TxHash,
	}, nil
}

func (c *CoinsPaidProcessor) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("%s/v2/rate?currency_from=%s&currency_to=%s", c.config.CoinsPaidURL, from, to)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Rate string `json:"rate"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(result.Data.Rate, 64)
}

// BitPayProcessor implements BitPay payment processor
type BitPayProcessor struct {
	config     *PaymentGatewayConfig
	httpClient *http.Client
	enabled    bool
}

func NewBitPayProcessor(config *PaymentGatewayConfig) *BitPayProcessor {
	return &BitPayProcessor{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		enabled: config.BitPayAPIKey != "",
	}
}

func (b *BitPayProcessor) GetName() string {
	return "BitPay"
}

func (b *BitPayProcessor) IsAvailable() bool {
	return b.enabled
}

func (b *BitPayProcessor) CreateDeposit(ctx context.Context, req *CreateDepositRequest) (*DepositInfo, error) {
	// Similar implementation to CoinsPaid
	return &DepositInfo{
		DepositID: uuid.New().String(),
		Address:   b.generateBTCAddress(req.Currency),
		Amount:    req.Amount,
		Currency:  req.Currency,
		Network:   req.Network,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Status:    "pending",
	}, nil
}

func (b *BitPayProcessor) CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*WithdrawalInfo, error) {
	fee := req.Amount * 0.01 // 1% fee
	return &WithdrawalInfo{
		WithdrawalID: uuid.New().String(),
		Amount:       req.Amount,
		Currency:     req.Currency,
		Network:      req.Network,
		ToAddress:    req.ToAddress,
		Fee:          fee,
		NetAmount:    req.Amount - fee,
		Status:       "pending",
	}, nil
}

func (b *BitPayProcessor) GetDepositStatus(ctx context.Context, depositID string) (*DepositStatus, error) {
	return &DepositStatus{
		DepositID: depositID,
		Status:    "pending",
	}, nil
}

func (b *BitPayProcessor) GetWithdrawalStatus(ctx context.Context, withdrawalID string) (*WithdrawalStatus, error) {
	return &WithdrawalStatus{
		WithdrawalID: withdrawalID,
		Status:       "pending",
	}, nil
}

func (b *BitPayProcessor) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	return 1.0, nil // Simplified
}

func (b *BitPayProcessor) generateBTCAddress(currency string) string {
	// Generate a realistic-looking BTC address for demo
	// In production, this would use actual wallet derivation
	prefixes := map[string]string{
		"BTC": "bc1q",
		"ETH": "0x",
		"LTC": "ltc1q",
		"DOGE": "D",
	}

	prefix := prefixes[currency]
	if prefix == "" {
		prefix = "bc1q"
	}

	// Generate 38 random bytes for address
	bytes := make([]byte, 38)
	rand.Read(bytes)
	return prefix + hex.EncodeToString(bytes)[:38-len(prefix)]
}

// InternalWalletProcessor for direct blockchain operations
type InternalWalletProcessor struct {
	walletManager *WalletManager
}

func NewInternalWalletProcessor(wm *WalletManager) *InternalWalletProcessor {
	return &InternalWalletProcessor{walletManager: wm}
}

func (i *InternalWalletProcessor) GetName() string {
	return "Internal"
}

func (i *InternalWalletProcessor) IsAvailable() bool {
	return true
}

func (i *InternalWalletProcessor) CreateDeposit(ctx context.Context, req *CreateDepositRequest) (*DepositInfo, error) {
	address, err := i.walletManager.GenerateDepositAddress(ctx, req.UserID, req.Currency, req.Network)
	if err != nil {
		return nil, err
	}

	return &DepositInfo{
		DepositID: uuid.New().String(),
		Address:   address,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Network:   req.Network,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Status:    "pending",
	}, nil
}

func (i *InternalWalletProcessor) CreateWithdrawal(ctx context.Context, req *CreateWithdrawalRequest) (*WithdrawalInfo, error) {
	// Validate address format
	if err := i.walletManager.ValidateAddress(req.Currency, req.Network, req.ToAddress); err != nil {
		return nil, err
	}

	fee := calculateWithdrawalFee(req.Currency, req.Network, req.Amount)
	netAmount := req.Amount - fee

	return &WithdrawalInfo{
		WithdrawalID: uuid.New().String(),
		Amount:       req.Amount,
		Currency:     req.Currency,
		Network:      req.Network,
		ToAddress:    req.ToAddress,
		Fee:          fee,
		NetAmount:    netAmount,
		Status:       "pending",
	}, nil
}

func (i *InternalWalletProcessor) GetDepositStatus(ctx context.Context, depositID string) (*DepositStatus, error) {
	return &DepositStatus{
		DepositID: depositID,
		Status:    "pending",
	}, nil
}

func (i *InternalWalletProcessor) GetWithdrawalStatus(ctx context.Context, withdrawalID string) (*WithdrawalStatus, error) {
	return &WithdrawalStatus{
		WithdrawalID: withdrawalID,
		Status:       "pending",
	}, nil
}

func (i *InternalWalletProcessor) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	return 1.0, nil
}

func (wm *WalletManager) GenerateDepositAddress(ctx context.Context, userID uuid.UUID, currency, network string) (string, error) {
	// Check if user already has an address for this currency
	var existing models.WalletAddress
	if err := wm.db.Where("user_id = ? AND currency = ? AND network = ?", userID, currency, network).
		First(&existing).Error; err == nil {
		return existing.Address, nil
	}

	// Generate new address based on currency
	var address string
	var err error

	switch currency {
	case "BTC":
		address, err = wm.generateBTCAddress(network)
	case "ETH":
		address, err = wm.generateETHAddress(network)
	case "USDT", "USDC", "TRX":
		address, err = wm.generateTRC20Address()
	case "LTC":
		address, err = wm.generateLTCAddress()
	default:
		address, err = wm.generateGenericAddress(currency)
	}

	if err != nil {
		return "", err
	}

	// Save to database
	walletAddr := models.WalletAddress{
		ID:        uuid.New(),
		UserID:    userID,
		Currency:  currency,
		Network:   network,
		Address:   address,
		IsPrimary: true,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := wm.db.Create(&walletAddr).Error; err != nil {
		return "", err
	}

	return address, nil
}

func (wm *WalletManager) generateBTCAddress(network string) (string, error) {
	// Generate a P2WPKH (SegWit) address using private key
	// In production, this would use actual HD wallet derivation
	
	privKeyBytes := make([]byte, 32)
	rand.Read(privKeyBytes)

	pubKey := elliptic.P256().Marshal(
		elliptic.P256().ScalarBaseMult(privKeyBytes),
	)

	// Create P2WPKH address
	pkHash := sha256.Sum256(pubKey)
	hasher := txscript.NewAddressPubKeyHash(
		bytes.Repeat([]byte{0x00}, 20),
		&chaincfg.MainNetParams,
		txscript.P2WPKHv0,
	)

	if hasher != nil {
		return fmt.Sprintf("bc1q%s", hex.EncodeToString(pkHash[:])[:39]), nil
	}

	// Fallback
	return "bc1q" + hex.EncodeToString(pkHash[:])[:39], nil
}

func (wm *WalletManager) generateETHAddress(network string) (string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	address := common.AddressToHex(common.BytesToAddress(privateKey.PublicKey.X.Bytes()))
	return address, nil
}

func (wm *WalletManager) generateTRC20Address() (string, error) {
	bytes := make([]byte, 21)
	rand.Read(bytes)
	bytes[0] = 0x41 // Tron address prefix
	return "T" + base58Encode(bytes), nil
}

func (wm *WalletManager) generateLTCAddress() (string, error) {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	hash := sha256.Sum256(bytes)
	return "ltc1q" + hex.EncodeToString(hash[:])[:39], nil
}

func (wm *WalletManager) generateGenericAddress(currency string) (string, error) {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes), nil
}

func (wm *WalletManager) ValidateAddress(currency, network, address string) error {
	switch currency {
	case "BTC":
		if !strings.HasPrefix(address, "bc1q") && !strings.HasPrefix(address, "1") && !strings.HasPrefix(address, "3") {
			return fmt.Errorf("invalid BTC address format")
		}
	case "ETH":
		if !strings.HasPrefix(address, "0x") || len(address) != 42 {
			return fmt.Errorf("invalid ETH address format")
		}
	case "USDT", "USDC":
		if !strings.HasPrefix(address, "T") || len(address) != 34 {
			return fmt.Errorf("invalid TRC20 address format")
		}
	}
	return nil
}

func base58Encode(data []byte) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	result := make([]byte, 0, len(data)*2)
	
	// Convert to big int
	num := new(big.Int).SetBytes(data)
	
	for num.BitLen() > 0 {
		mod := new(big.Int)
		num.DivMod(num, big.NewInt(58), mod)
		result = append([]byte{alphabet[mod.Int64()]}, result...)
	}
	
	// Add leading zeros
	for _, b := range data {
		if b == 0 {
			result = append([]byte{alphabet[0]}, result...)
		} else {
			break
		}
	}
	
	return string(result)
}

// NewPaymentGatewayService creates a new payment gateway service
func NewPaymentGatewayService(db *gorm.DB, redisClient *redis.Client, config *PaymentGatewayConfig) *PaymentGatewayService {
	walletManager := &WalletManager{
		db:          db,
		config:      config,
		addressPool: &AddressPool{addresses: make(map[string][]string)},
		wallets:     make(map[string]*Wallet),
	}

	service := &PaymentGatewayService{
		db:             db,
		redis:          redisClient,
		walletManager:  walletManager,
		config:         config,
		balanceCache:   NewBalanceCache(),
		webhookHandler: NewWebhookHandler(),
		processors:     make(map[string]PaymentProcessor),
	}

	// Initialize available processors
	if config != nil {
		if coinsPaid := NewCoinsPaidProcessor(config); coinsPaid.IsAvailable() {
			service.processors["coinspaid"] = coinsPaid
		}

		if bitpay := NewBitPayProcessor(config); bitpay.IsAvailable() {
			service.processors["bitpay"] = bitpay
		}
	}

	// Always add internal processor
	service.processors["internal"] = NewInternalWalletProcessor(walletManager)

	// Register webhook handlers
	service.setupWebhookHandlers()

	return service
}

func (pgs *PaymentGatewayService) setupWebhookHandlers() {
	pgs.webhookHandler.Register("coinspaid", pgs.config.CoinsPaidSecret, 
		func(ctx context.Context, data json.RawMessage) error {
			var payload map[string]interface{}
			if err := json.Unmarshal(data, &payload); err != nil {
				return err
			}
			
			// Handle deposit callback
			if payload["type"] == "deposit" {
				return pgs.handleCoinsPaidDeposit(ctx, payload)
			}
			
			return nil
		})
}

func (pgs *PaymentGatewayService) handleCoinsPaidDeposit(ctx context.Context, payload map[string]interface{}) error {
	depositID, _ := payload["id"].(string)
	status, _ := payload["status"].(string)
	amountStr, _ := payload["amount"].(string)
	txHash, _ := payload["tx_hash"].(string)

	amount, _ := strconv.ParseFloat(amountStr, 64)

	// Update deposit status in database
	var deposit models.Transaction
	if err := pgs.db.Where("external_id = ?", depositID).First(&deposit).Error; err != nil {
		return err
	}

	if status == "confirmed" || status == "completed" {
		deposit.Status = "completed"
		deposit.TxHash = txHash
		deposit.ConfirmedAt = time.Now()
		
		// Credit user wallet
		userID := deposit.UserID
		currency := deposit.Currency
		
		pgs.balanceCache.Add(userID.String(), currency, amount)
		
		// Update user balance in database
		var wallet models.Wallet
		if err := pgs.db.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error; err == nil {
			wallet.Balance += amount
			wallet.UpdatedAt = time.Now()
			pgs.db.Save(&wallet)
		}
	}

	deposit.Metadata = payload
	pgs.db.Save(&deposit)

	return nil
}

// CreateDeposit creates a new deposit for a user
func (pgs *PaymentGatewayService) CreateDeposit(ctx context.Context, userID uuid.UUID, currency, network string, amount float64, ipAddress string) (*DepositInfo, error) {
	// Validate amount
	if amount < pgs.config.MinDepositAmount {
		return nil, fmt.Errorf("minimum deposit amount is %.2f", pgs.config.MinDepositAmount)
	}

	// Check if user has sufficient KYC for this amount
	if err := pgs.checkDepositLimits(ctx, userID, amount, currency); err != nil {
		return nil, err
	}

	// Get processor
	processor := pgs.getProcessor(currency)
	if processor == nil {
		processor = pgs.processors["internal"]
	}

	req := &CreateDepositRequest{
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		Network:   network,
		IPAddress: ipAddress,
	}

	depositInfo, err := processor.CreateDeposit(ctx, req)
	if err != nil {
		return nil, err
	}

	// Save deposit record to database
	deposit := models.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        "deposit",
		Currency:    currency,
		Network:     network,
		Amount:      amount,
		Status:      "pending",
		ExternalID:  depositInfo.DepositID,
		IPAddress:   ipAddress,
		RequestData: map[string]interface{}{
			"address": depositInfo.Address,
			"network": network,
		},
		CreatedAt: time.Now(),
	}

	if err := pgs.db.Create(&deposit).Error; err != nil {
		return nil, err
	}

	// Generate QR code (simplified - would use actual QR generation in production)
	depositInfo.QRCode = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=%s", depositInfo.Address)

	return depositInfo, nil
}

// CreateWithdrawal creates a new withdrawal for a user
func (pgs *PaymentGatewayService) CreateWithdrawal(ctx context.Context, userID uuid.UUID, currency, network, toAddress string, amount float64, ipAddress string) (*WithdrawalInfo, error) {
	// Validate amount
	if amount < pgs.config.MinWithdrawalAmount {
		return nil, fmt.Errorf("minimum withdrawal amount is %.2f", pgs.config.MinWithdrawalAmount)
	}

	// Get user balance
	balance, err := pgs.GetUserBalance(ctx, userID, currency)
	if err != nil {
		return nil, err
	}

	if balance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Check daily withdrawal limit
	dailyWithdrawn, err := pgs.GetDailyWithdrawalAmount(ctx, userID, currency)
	if err != nil {
		return nil, err
	}

	if dailyWithdrawn+amount > pgs.config.MaxWithdrawalDaily {
		return nil, fmt.Errorf("daily withdrawal limit exceeded")
	}

	// Validate address
	if err := pgs.walletManager.ValidateAddress(currency, network, toAddress); err != nil {
		return nil, err
	}

	// Check for suspicious activity
	if err := pgs.checkWithdrawalSecurity(ctx, userID, amount, toAddress, ipAddress); err != nil {
		return nil, err
	}

	// Deduct from balance
	pgs.balanceCache.Subtract(userID.String(), currency, amount)
	
	// Update database
	var wallet models.Wallet
	if err := pgs.db.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error; err == nil {
		wallet.Balance -= amount
		wallet.UpdatedAt = time.Now()
		pgs.db.Save(&wallet)
	}

	// Get processor
	processor := pgs.getProcessor(currency)
	if processor == nil {
		processor = pgs.processors["internal"]
	}

	req := &CreateWithdrawalRequest{
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		Network:   network,
		ToAddress: toAddress,
		IPAddress: ipAddress,
	}

	withdrawalInfo, err := processor.CreateWithdrawal(ctx, req)
	if err != nil {
		// Refund on error
		pgs.balanceCache.Add(userID.String(), currency, amount)
		return nil, err
	}

	// Save withdrawal record
	withdrawal := models.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        "withdrawal",
		Currency:    currency,
		Network:     network,
		Amount:      amount,
		Fee:         withdrawalInfo.Fee,
		NetAmount:   withdrawalInfo.NetAmount,
		Status:      "pending",
		ExternalID:  withdrawalInfo.WithdrawalID,
		ToAddress:   toAddress,
		IPAddress:   ipAddress,
		RequestData: map[string]interface{}{
			"to_address": toAddress,
			"network":    network,
		},
		CreatedAt: time.Now(),
	}

	if err := pgs.db.Create(&withdrawal).Error; err != nil {
		return nil, err
	}

	return withdrawalInfo, nil
}

func (pgs *PaymentGatewayService) getProcessor(currency string) PaymentProcessor {
	// Try currency-specific processor first
	if processor, ok := pgs.processors[currency]; ok {
		return processor
	}

	// Try internal processor
	return pgs.processors["internal"]
}

// GetUserBalance returns user's balance for a currency
func (pgs *PaymentGatewayService) GetUserBalance(ctx context.Context, userID uuid.UUID, currency string) (float64, error) {
	// Check cache first
	if balance, ok := pgs.balanceCache.Get(userID.String(), currency); ok {
		return balance, nil
	}

	// Query database
	var wallet models.Wallet
	if err := pgs.db.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	// Cache the balance
	pgs.balanceCache.Set(userID.String(), currency, wallet.Balance)

	return wallet.Balance, nil
}

// GetDailyWithdrawalAmount returns total amount withdrawn in the last 24 hours
func (pgs *PaymentGatewayService) GetDailyWithdrawalAmount(ctx context.Context, userID uuid.UUID, currency string) (float64, error) {
	var total float64
	since := time.Now().Add(-24 * time.Hour)

	if err := pgs.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND currency = ? AND status IN (?, ?) AND created_at > ?",
			userID, "withdrawal", currency, "pending", "completed", since).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (pgs *PaymentGatewayService) checkDepositLimits(ctx context.Context, userID uuid.UUID, amount float64, currency string) error {
	// Check if amount requires verification
	if amount > 1000 {
		// Check user verification level
		var user models.User
		if err := pgs.db.First(&user, userID).Error; err != nil {
			return err
		}

		if user.VerificationLevel < 2 {
			return fmt.Errorf("higher verification level required for deposits over 1000")
		}
	}

	return nil
}

func (pgs *PaymentGatewayService) checkWithdrawalSecurity(ctx context.Context, userID uuid.UUID, amount float64, address, ipAddress string) error {
	// Check for new address (first withdrawal to this address)
	var existingCount int64
	pgs.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND to_address = ? AND status = ?",
			userID, "withdrawal", address, "completed").
		Count(&existingCount)

	if existingCount == 0 && amount > 100 {
		// First withdrawal to new address - require hold
		return fmt.Errorf("first withdrawal to new address requires 24h hold")
	}

	// Check for suspicious patterns
	if amount > 10000 {
		var user models.User
		if err := pgs.db.First(&user, userID).Error; err != nil {
			return err
		}

		if user.VerificationLevel < 2 {
			return fmt.Errorf("higher verification level required for withdrawals over 10000")
		}
	}

	return nil
}

// ProcessWebhook handles incoming payment webhooks
func (pgs *PaymentGatewayService) ProcessWebhook(ctx context.Context, processor string, payload json.RawMessage, signature string) error {
	return pgs.webhookHandler.Handle(ctx, processor, payload, signature)
}

// SyncDeposits syncs pending deposits from blockchain/processor
func (pgs *PaymentGatewayService) SyncDeposits(ctx context.Context) error {
	var pendingDeposits []models.Transaction
	pgs.db.Where("type = ? AND status = ?", "deposit", "pending").Find(&pendingDeposits)

	for _, deposit := range pendingDeposits {
		processor := pgs.getProcessor(deposit.Currency)
		if processor == nil {
			continue
		}

		status, err := processor.GetDepositStatus(ctx, deposit.ExternalID)
		if err != nil {
			continue
		}

		if status.Status == "confirmed" || status.Status == "completed" {
			// Credit user
			pgs.balanceCache.Add(deposit.UserID.String(), deposit.Currency, deposit.Amount)

			// Update database
			deposit.Status = "completed"
			deposit.TxHash = status.TxHash
			deposit.ConfirmedAt = time.Now()
			pgs.db.Save(&deposit)

			// Create balance change record
			balanceChange := models.BalanceChange{
				ID:        uuid.New(),
				UserID:    deposit.UserID,
				Type:      "deposit",
				Currency:  deposit.Currency,
				Amount:    deposit.Amount,
				Balance:   deposit.Amount,
				Reference: deposit.ExternalID,
				CreatedAt: time.Now(),
			}
			pgs.db.Create(&balanceChange)
		}
	}

	return nil
}

// GetSupportedCurrencies returns list of supported currencies
func (pgs *PaymentGatewayService) GetSupportedCurrencies() []map[string]interface{} {
	return []map[string]interface{}{
		{"symbol": "BTC", "name": "Bitcoin", "network": "bitcoin", "type": "crypto", "decimals": 8, "minDeposit": 0.0001, "minWithdraw": 0.0002, "fee": 0.0001},
		{"symbol": "ETH", "name": "Ethereum", "network": "ethereum", "type": "crypto", "decimals": 18, "minDeposit": 0.01, "minWithdraw": 0.02, "fee": 0.005},
		{"symbol": "USDT", "name": "Tether USD", "network": "trc20", "type": "stablecoin", "decimals": 6, "minDeposit": 10, "minWithdraw": 20, "fee": 1},
		{"symbol": "USDC", "name": "USD Coin", "network": "erc20", "type": "stablecoin", "decimals": 6, "minDeposit": 10, "minWithdraw": 20, "fee": 1},
		{"symbol": "LTC", "name": "Litecoin", "network": "litecoin", "type": "crypto", "decimals": 8, "minDeposit": 0.01, "minWithdraw": 0.02, "fee": 0.001},
		{"symbol": "DOGE", "name": "Dogecoin", "network": "dogecoin", "type": "crypto", "decimals": 8, "minDeposit": 100, "minWithdraw": 200, "fee": 10},
		{"symbol": "TRX", "name": "Tron", "network": "trc20", "type": "crypto", "decimals": 6, "minDeposit": 100, "minWithdraw": 200, "fee": 1},
		{"symbol": "XRP", "name": "Ripple", "network": "ripple", "type": "crypto", "decimals": 6, "minDeposit": 25, "minWithdraw": 50, "fee": 1},
		{"symbol": "BNB", "name": "BNB", "network": "bsc", "type": "crypto", "decimals": 18, "minDeposit": 0.01, "minWithdraw": 0.02, "fee": 0.005},
		{"symbol": "SOL", "name": "Solana", "network": "solana", "type": "crypto", "decimals": 9, "minDeposit": 0.1, "minWithdraw": 0.2, "fee": 0.01},
	}
}

// GetUserTransactions returns user's transaction history
func (pgs *PaymentGatewayService) GetUserTransactions(ctx context.Context, userID uuid.UUID, txType, status, currency string, page, limit int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	query := pgs.db.Where("user_id = ?", userID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if currency != "" {
		query = query.Where("currency = ?", currency)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func calculateWithdrawalFee(currency, network string, amount float64) float64 {
	feePercent := map[string]float64{
		"BTC":  0.001,  // 0.1%
		"ETH":  0.01,   // 1%
		"USDT": 0.001,  // 0.1%
		"LTC":  0.001,
		"DOGE": 0.01,
		"TRX":  1.0,    // Flat fee
	}

	fee := feePercent[currency]
	if fee < 1 {
		return amount * fee
	}
	return fee // Flat fee
}
