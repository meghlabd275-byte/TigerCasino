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

	// Chainlink
	s.networks["LINK"] = &NetworkConfig{
		Name:            "Chainlink",
		Symbol:          "LINK",
		ChainID:         1,
		RPCURL:          "https://eth-rpc.example.com",
		ExplorerURL:     "https://etherscan.io",
		MinConfirmations: 12,
		MinDeposit:      5.0,
		MinWithdrawal:   5.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Uniswap
	s.networks["UNI"] = &NetworkConfig{
		Name:            "Uniswap",
		Symbol:          "UNI",
		ChainID:         1,
		RPCURL:          "https://eth-rpc.example.com",
		ExplorerURL:     "https://etherscan.io",
		MinConfirmations: 12,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.5,
		IsActive:        true,
	}

	// Aave
	s.networks["AAVE"] = &NetworkConfig{
		Name:            "Aave",
		Symbol:          "AAVE",
		ChainID:         1,
		RPCURL:          "https://eth-rpc.example.com",
		ExplorerURL:     "https://etherscan.io",
		MinConfirmations: 12,
		MinDeposit:      0.5,
		MinWithdrawal:   0.5,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Maker
	s.networks["MKR"] = &NetworkConfig{
		Name:            "Maker",
		Symbol:          "MKR",
		ChainID:         1,
		RPCURL:          "https://eth-rpc.example.com",
		ExplorerURL:     "https://etherscan.io",
		MinConfirmations: 12,
		MinDeposit:      0.1,
		MinWithdrawal:   0.1,
		WithdrawalFee:   0.05,
		IsActive:        true,
	}

	// Cosmos
	s.networks["ATOM"] = &NetworkConfig{
		Name:            "Cosmos",
		Symbol:          "ATOM",
		ChainID:         0,
		RPCURL:          "https://cosmos-rpc.example.com",
		ExplorerURL:     "https://mintscan.io",
		MinConfirmations: 10,
		MinDeposit:      0.5,
		MinWithdrawal:   0.5,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Algorand
	s.networks["ALGO"] = &NetworkConfig{
		Name:            "Algorand",
		Symbol:          "ALGO",
		ChainID:         0,
		RPCURL:          "https://algo-rpc.example.com",
		ExplorerURL:     "https://algoexplorer.io",
		MinConfirmations: 4,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// VeChain
	s.networks["VET"] = &NetworkConfig{
		Name:            "VeChain",
		Symbol:          "VET",
		ChainID:         0,
		RPCURL:          "https://vet-rpc.example.com",
		ExplorerURL:     "https://vechainstats.com",
		MinConfirmations: 20,
		MinDeposit:      100.0,
		MinWithdrawal:   100.0,
		WithdrawalFee:   10.0,
		IsActive:        true,
	}

	// Hedera
	s.networks["HBAR"] = &NetworkConfig{
		Name:            "Hedera Hashgraph",
		Symbol:          "HBAR",
		ChainID:         0,
		RPCURL:          "https://hbar-rpc.example.com",
		ExplorerURL:     "https://hashscan.io",
		MinConfirmations: 10,
		MinDeposit:      50.0,
		MinWithdrawal:   50.0,
		WithdrawalFee:   5.0,
		IsActive:        true,
	}

	// Fantom
	s.networks["FTM"] = &NetworkConfig{
		Name:            "Fantom",
		Symbol:          "FTM",
		ChainID:         250,
		RPCURL:          "https://fantom-rpc.example.com",
		ExplorerURL:     "https://ftmscan.com",
		MinConfirmations: 15,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Near
	s.networks["NEAR"] = &NetworkConfig{
		Name:            "NEAR Protocol",
		Symbol:          "NEAR",
		ChainID:         0,
		RPCURL:          "https://near-rpc.example.com",
		ExplorerURL:     "https://explorer.near.org",
		MinConfirmations: 5,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Aptos
	s.networks["APT"] = &NetworkConfig{
		Name:            "Aptos",
		Symbol:          "APT",
		ChainID:         0,
		RPCURL:          "https://aptos-rpc.example.com",
		ExplorerURL:     "https://aptoscan.com",
		MinConfirmations: 1,
		MinDeposit:      0.5,
		MinWithdrawal:   0.5,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Arbitrum
	s.networks["ARB"] = &NetworkConfig{
		Name:            "Arbitrum One",
		Symbol:          "ARB",
		ChainID:         42161,
		RPCURL:          "https://arb-rpc.example.com",
		ExplorerURL:     "https://arbiscan.io",
		MinConfirmations: 15,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Optimism
	s.networks["OP"] = &NetworkConfig{
		Name:            "Optimism",
		Symbol:          "OP",
		ChainID:         10,
		RPCURL:          "https://opt-rpc.example.com",
		ExplorerURL:     "https://optimistic.etherscan.io",
		MinConfirmations: 15,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Polygon zkEVM
	s.networks["POL"] = &NetworkConfig{
		Name:            "Polygon zkEVM",
		Symbol:          "POL",
		ChainID:         1101,
		RPCURL:          "https://zkevm-rpc.example.com",
		ExplorerURL:     "https://zkevm.polygonscan.com",
		MinConfirmations: 10,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Base
	s.networks["BASE"] = &NetworkConfig{
		Name:            "Base",
		Symbol:          "ETH",
		ChainID:         8453,
		RPCURL:          "https://base-rpc.example.com",
		ExplorerURL:     "https://basescan.org",
		MinConfirmations: 15,
		MinDeposit:      0.001,
		MinWithdrawal:   0.001,
		WithdrawalFee:   0.0005,
		IsActive:        true,
	}

	// zkSync Era
	s.networks["ZK"] = &NetworkConfig{
		Name:            "zkSync Era",
		Symbol:          "ETH",
		ChainID:         324,
		RPCURL:          "https://zksync-rpc.example.com",
		ExplorerURL:     "https://explorer.zksync.io",
		MinConfirmations: 1,
		MinDeposit:      0.001,
		MinWithdrawal:   0.001,
		WithdrawalFee:   0.0005,
		IsActive:        true,
	}

	// Stellar
	s.networks["XLM"] = &NetworkConfig{
		Name:            "Stellar",
		Symbol:          "XLM",
		ChainID:         0,
		RPCURL:          "https://stellar-rpc.example.com",
		ExplorerURL:     "https://stellarscan.io",
		MinConfirmations: 1,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Monero
	s.networks["XMR"] = &NetworkConfig{
		Name:            "Monero",
		Symbol:          "XMR",
		ChainID:         0,
		RPCURL:          "https://monero-rpc.example.com",
		ExplorerURL:     "https://xmrchain.net",
		MinConfirmations: 10,
		MinDeposit:      0.1,
		MinWithdrawal:   0.1,
		WithdrawalFee:   0.01,
		IsActive:        true,
	}

	// Zcash
	s.networks["ZEC"] = &NetworkConfig{
		Name:            "Zcash",
		Symbol:          "ZEC",
		ChainID:         0,
		RPCURL:          "https://zcash-rpc.example.com",
		ExplorerURL:     "https://zcashblockexplorer.com",
		MinConfirmations: 10,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Dash
	s.networks["DASH"] = &NetworkConfig{
		Name:            "Dash",
		Symbol:          "DASH",
		ChainID:         0,
		RPCURL:          "https://dash-rpc.example.com",
		ExplorerURL:     "https://dashblockexplorer.com",
		MinConfirmations: 6,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Neo
	s.networks["NEO"] = &NetworkConfig{
		Name:            "Neo",
		Symbol:          "NEO",
		ChainID:         0,
		RPCURL:          "https://neo-rpc.example.com",
		ExplorerURL:     "https://neotube.io",
		MinConfirmations: 2,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// EOS
	s.networks["EOS"] = &NetworkConfig{
		Name:            "EOS",
		Symbol:          "EOS",
		ChainID:         0,
		RPCURL:          "https://eos-rpc.example.com",
		ExplorerURL:     "https://bloks.io",
		MinConfirmations: 1,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Tezos
	s.networks["XTZ"] = &NetworkConfig{
		Name:            "Tezos",
		Symbol:          "XTZ",
		ChainID:         0,
		RPCURL:          "https://tezos-rpc.example.com",
		ExplorerURL:     "https://tzstats.com",
		MinConfirmations: 1,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Flow
	s.networks["FLOW"] = &NetworkConfig{
		Name:            "Flow",
		Symbol:          "FLOW",
		ChainID:         0,
		RPCURL:          "https://flow-rpc.example.com",
		ExplorerURL:     "https://flowscan.io",
		MinConfirmations: 2,
		MinDeposit:      1.0,
		MinWithdrawal:   1.0,
		WithdrawalFee:   0.1,
		IsActive:        true,
	}

	// Internet Computer
	s.networks["ICP"] = &NetworkConfig{
		Name:            "Internet Computer",
		Symbol:          "ICP",
		ChainID:         0,
		RPCURL:          "https://icp-rpc.example.com",
		ExplorerURL:     "https://icpcash.com",
		MinConfirmations: 1,
		MinDeposit:      0.1,
		MinWithdrawal:   0.1,
		WithdrawalFee:   0.01,
		IsActive:        true,
	}

	// Optimism Bedrock
	s.networks["ETH_OPT"] = &NetworkConfig{
		Name:            "Optimism (ETH)",
		Symbol:          "ETH",
		ChainID:         10,
		RPCURL:          "https://opt-rpc.example.com",
		ExplorerURL:     "https://optimistic.etherscan.io",
		MinConfirmations: 15,
		MinDeposit:      0.001,
		MinWithdrawal:   0.001,
		WithdrawalFee:   0.0005,
		IsActive:        true,
	}

	// Gnosis
	s.networks["GNO"] = &NetworkConfig{
		Name:            "Gnosis Chain",
		Symbol:          "GNO",
		ChainID:         100,
		RPCURL:          "https://gnosis-rpc.example.com",
		ExplorerURL:     "https://gnosisscan.io",
		MinConfirmations: 15,
		MinDeposit:      0.01,
		MinWithdrawal:   0.01,
		WithdrawalFee:   0.001,
		IsActive:        true,
	}

	// Cronos
	s.networks["CRO"] = &NetworkConfig{
		Name:            "Cronos",
		Symbol:          "CRO",
		ChainID:         25,
		RPCURL:          "https://cronos-rpc.example.com",
		ExplorerURL:     "https://cronoscan.com",
		MinConfirmations: 20,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// Klaytn
	s.networks["KLAY"] = &NetworkConfig{
		Name:            "Klaytn",
		Symbol:          "KLAY",
		ChainID:         8217,
		RPCURL:          "https://klaytn-rpc.example.com",
		ExplorerURL:     "https://klaytnscope.com",
		MinConfirmations: 1,
		MinDeposit:      10.0,
		MinWithdrawal:   10.0,
		WithdrawalFee:   1.0,
		IsActive:        true,
	}

	// THORChain
	s.networks["RUNE"] = &NetworkConfig{
		Name:            "THORChain",
		Symbol:          "RUNE",
		ChainID:         0,
		RPCURL:          "https://thor-rpc.example.com",
		ExplorerURL:     "https://runescan.io",
		MinConfirmations: 1,
		MinDeposit:      0.1,
		MinWithdrawal:   0.1,
		WithdrawalFee:   0.01,
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
