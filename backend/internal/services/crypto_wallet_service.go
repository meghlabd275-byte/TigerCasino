package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CryptoWalletService handles multi-cryptocurrency wallet operations
type CryptoWalletService struct {
	mu            sync.RWMutex
	wallets      map[string]*CryptoWallet
	transactions map[string]*CryptoTransaction
	networks     map[string]*NetworkConfig
}

// CryptoWallet represents a user's cryptocurrency wallet
type CryptoWallet struct {
	ID            string
	UserID        string
	Currency     string
	Address       string
	Balance       float64
	PendingBalance float64
	Network       string
	IsInternal    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CryptoTransaction represents a crypto transaction
type CryptoTransaction struct {
	ID            string
	WalletID     string
	Hash         string
	Type         string // deposit, withdrawal, transfer
	Amount       float64
	Fee          float64
	Status       string // pending, confirmed, failed
	Confirmations int
	Network      string
	FromAddress  string
	ToAddress    string
	CreatedAt    time.Time
	ConfirmedAt  *time.Time
}

// NetworkConfig represents blockchain network configuration
type NetworkConfig struct {
	Name            string
	Symbol          string
	ChainID         int
	RPCURL          string
	ExplorerURL     string
	MinConfirmations int
	MinDeposit      float64
	MinWithdrawal   float64
	WithdrawalFee   float64
	IsActive        bool
}

// NewCryptoWalletService creates a new crypto wallet service
func NewCryptoWalletService() *CryptoWalletService {
	s := &CryptoWalletService{
		wallets:      make(map[string]*CryptoWallet),
		transactions: make(map[string]*CryptoTransaction),
		networks:    make(map[string]*NetworkConfig),
	}
	s.initializeNetworks()
	return s
}

func (s *CryptoWalletService) initializeNetworks() {
	// Bitcoin Network
	s.networks["BTC"] = &NetworkConfig{
		Name:            "Bitcoin",
		Symbol:          "BTC",
		ChainID:         1,
		RPCURL:          "https://btc-rpc.example.com",
		ExplorerURL:     "https://blockstream.info",
		MinConfirmations: 3,
		MinDeposit:      0.0001,
		MinWithdrawal:   0.0001,
		WithdrawalFee:   0.0001,
		IsActive:        true,
	}

	// Ethereum Network
	s.networks["ETH"] = &NetworkConfig{
		Name:            "Ethereum",
		Symbol:          "ETH",
		ChainID:         1,
		RPCURL:          "https://eth-rpc.example.com",
		ExplorerURL:     "https://etherscan.io",
		MinConfirmations: 12,
		MinDeposit:      0.001,
		MinWithdrawal:   0.001,
		WithdrawalFee:   0.005,
		IsActive:        true,
	}

	// USDT (ERC-20)
	s.networks["USDT"] = &NetworkConfig{
		Name:            "Tether USD (Ethereum)",
		Symbol:          "USDT",
		ChainID:         1,
		RPCURL:          "https://eth-rpc.example.com",
		ExplorerURL:     "https://etherscan.io",
		MinConfirmations: 12,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   5.0,
		IsActive:        true,
	}

	// Litecoin
	s.networks["LTC"] = &NetworkConfig{
		Name:            "Litecoin",
		Symbol:          "LTC",
		ChainID:         2,
		RPCURL:          "https://ltc-rpc.example.com",
		ExplorerURL:     "https://blockchair.com/litecoin",
		MinConfirmations: 6,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Dogecoin
	s.networks["DOGE"] = &NetworkConfig{
		Name:            "Dogecoin",
		Symbol:          "DOGE",
		ChainID:         3,
		RPCURL:          "https://doge-rpc.example.com",
		ExplorerURL:     "https://dogecoin.com",
		MinConfirmations: 6,
		MinDeposit:      100.0,
		MinWithdrawal:   100.0,
		WithdrawalFee:   10.0,
		IsActive:        true,
	}

	// Bitcoin Cash
	s.networks["BCH"] = &NetworkConfig{
		Name:            "Bitcoin Cash",
		Symbol:          "BCH",
		ChainID:         4,
		RPCURL:          "https://bch-rpc.example.com",
		ExplorerURL:     "https://blockchair.com/bitcoin-cash",
		MinConfirmations: 6,
		MinDeposit:      0.001,
		MinWithdrawal:   0.001,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Binance Smart Chain
	s.networks["BNB"] = &NetworkConfig{
		Name:            "BNB Smart Chain",
		Symbol:          "BNB",
		ChainID:         56,
		RPCURL:          "https://bsc-rpc.example.com",
		ExplorerURL:     "https://bscscan.com",
		MinConfirmations: 15,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.005,
		IsActive:        true,
	}

	// USDT (BEP-20)
	s.netWORKS["USDT_BSC"] = &NetworkConfig{
		Name:            "Tether USD (BSC)",
		Symbol:          "USDT",
		ChainID:         56,
		RPCURL:          "https://bsc-rpc.example.com",
		ExplorerURL:     "https://bscscan.com",
		MinConfirmations: 15,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Polygon
	s.networks["MATIC"] = &NetworkConfig{
		Name:            "Polygon",
		Symbol:          "MATIC",
		ChainID:         137,
		RPCURL:          "https://polygon-rpc.example.com",
		ExplorerURL:     "https://polygonscan.com",
		MinConfirmations: 15,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Avalanche
	s.networks["AVAX"] = &NetworkConfig{
		Name:            "Avalanche",
		Symbol:          "AVAX",
		ChainID:         43114,
		RPCURL:          "https://avax-rpc.example.com",
		ExplorerURL:     "https://snowtrace.io",
		MinConfirmations: 15,
		MinDeposit:      0.5,
		MinWithdrawal:   0.5,
		WithdrawalFee:   0.01,
		IsActive:        true,
	}

	// Solana
	s.networks["SOL"] = &NetworkConfig{
		Name:            "Solana",
		Symbol:          "SOL",
		ChainID:         101,
		RPCURL:          "https://solana-rpc.example.com",
		ExplorerURL:     "https://solscan.io",
		MinConfirmations: 32,
		MinDeposit:      0.1,
		MinWithdrawal:   0.1,
		WithdrawalFee:   0.01,
		IsActive:        true,
	}

	// Ripple
	s.networks["XRP"] = &NetworkConfig{
		Name:            "Ripple",
		Symbol:          "XRP",
		ChainID:         0,
		RPCURL:          "https://xrp-rpc.example.com",
		ExplorerURL:     "https://xrpscan.io",
		MinConfirmations: 1,
		MinDeposit:      20.0,
		MinWithdrawal:   20.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Tron
	s.networks["TRX"] = &NetworkConfig{
		Name:            "Tron",
		Symbol:          "TRX",
		ChainID:         0,
		RPCURL:          "https://trx-rpc.example.com",
		ExplorerURL:     "https://tronscan.io",
		MinConfirmations: 19,
		MinDeposit:      100.0,
		MinWithdrawal:   100.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Cardano
	s.networks["ADA"] = &NetworkConfig{
		Name:            "Cardano",
		Symbol:          "ADA",
		ChainID:         0,
		RPCURL:          "https://cardano-rpc.example.com",
		ExplorerURL:     "https://cardanoscan.io",
		MinConfirmations: 15,
		MinDeposit:      2.0,
		MinWithdrawal:   2.0,
		WithdrawalFee:   0.2,
		IsActive:        true,
	}

	// Polkadot
	s.networks["DOT"] = &NetworkConfig{
		Name:            "Polkadot",
		Symbol:          "DOT",
		ChainID:         0,
		RPCURL:          "https://polkadot-rpc.example.com",
		ExplorerURL:     "https://polkadot.subscan.io",
		MinConfirmations: 10,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}
}

// CreateWallet creates a new wallet for a user
func (s *CryptoWalletService) CreateWallet(userID, currency, network string) (*CryptoWallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if network exists
	net, ok := s.networks[network]
	if !ok {
		return nil, fmt.Errorf("network not supported")
	}

	wallet := &CryptoWallet{
		ID:            uuid.New().String(),
		UserID:        userID,
		Currency:      currency,
		Address:       generateAddress(currency, network), // Generate address
		Balance:       0.0,
		PendingBalance: 0.0,
		Network:       network,
		IsInternal:    true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.wallets[wallet.ID] = wallet
	return wallet, nil
}

// GetWallet returns a wallet by ID
func (s *CryptoWalletService) GetWallet(walletID string) (*CryptoWallet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wallet, ok := s.wallets[walletID]
	if !ok {
		return nil, fmt.Errorf("wallet not found")
	}
	return wallet, nil
}

// GetUserWallets returns all wallets for a user
func (s *CryptoWalletService) GetUserWallets(userID string) []CryptoWallet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var wallets []CryptoWallet
	for _, w := range s.wallets {
		if w.UserID == userID {
			wallets = append(wallets, *w)
		}
	}
	return wallets
}

// Deposit creates a deposit transaction
func (s *CryptoWalletService) Deposit(walletID string, amount float64, txHash, fromAddress string) (*CryptoTransaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, ok := s.wallets[walletID]
	if !ok {
		return nil, fmt.Errorf("wallet not found")
	}

	tx := &CryptoTransaction{
		ID:           uuid.New().String(),
		WalletID:     walletID,
		Hash:         txHash,
		Type:         "deposit",
		Amount:       amount,
		Fee:          0.0,
		Status:       "pending",
		Network:      wallet.Network,
		FromAddress: fromAddress,
		ToAddress:    wallet.Address,
		CreatedAt:   time.Now(),
	}

	s.transactions[tx.ID] = tx
	wallet.PendingBalance += amount

	return tx, nil
}

// ConfirmDeposit confirms a deposit after required confirmations
func (s *CryptoWalletService) ConfirmDeposit(txID string, confirmations int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, ok := s.transactions[txID]
	if !ok {
		return fmt.Errorf("transaction not found")
	}

	wallet, ok := s.wallets[tx.WalletID]
	if !ok {
		return fmt.Errorf("wallet not found")
	}

	network := s.networks[tx.Network]
	if confirmations >= network.MinConfirmations {
		tx.Status = "confirmed"
		tx.Confirmations = confirmations
		now := time.Now()
		tx.ConfirmedAt = &now
		wallet.Balance += tx.Amount
		wallet.PendingBalance -= tx.Amount
	}

	return nil
}

// Withdrawal creates a withdrawal request
func (s *CryptoWalletService) Withdrawal(walletID string, amount float64, toAddress string) (*CryptoTransaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, ok := s.wallets[walletID]
	if !ok {
		return nil, fmt.Errorf("wallet not found")
	}

	network := s.networks[wallet.Network]
	if amount < network.MinWithdrawal {
		return nil, fmt.Errorf("amount below minimum withdrawal: %f", network.MinWithdrawal)
	}

	totalAmount := amount + network.WithdrawalFee
	if wallet.Balance < totalAmount {
		return nil, fmt.Errorf("insufficient balance")
	}

	tx := &CryptoTransaction{
		ID:           uuid.New().String(),
		WalletID:     walletID,
		Hash:         "",
		Type:         "withdrawal",
		Amount:       amount,
		Fee:          network.WithdrawalFee,
		Status:       "pending",
		Network:      wallet.Network,
		FromAddress:  wallet.Address,
		ToAddress:    toAddress,
		CreatedAt:    time.Now(),
	}

	s.transactions[tx.ID] = tx
	wallet.Balance -= totalAmount

	return tx, nil
}

// GetSupportedNetworks returns all supported networks
func (s *CryptoWalletService) GetSupportedNetworks() []NetworkConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var configs []NetworkConfig
	for _, n := range s.networks {
		if n.IsActive {
			configs = append(configs, *n)
		}
	}
	return configs
}

// GetTransaction returns a transaction by ID
func (s *CryptoWalletService) GetTransaction(txID string) (*CryptoTransaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.transactions[txID]
	if !ok {
		return nil, fmt.Errorf("transaction not found")
	}
	return tx, nil
}

// Helper function to generate address (in production, use actual crypto libraries)
func generateAddress(currency, network string) string {
	return fmt.Sprintf("0x%s_%s", uuid.New().String()[:16], network)
}
