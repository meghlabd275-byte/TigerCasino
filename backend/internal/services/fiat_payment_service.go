package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FiatPaymentService handles fiat currency payments
type FiatPaymentService struct {
	mu            sync.RWMutex
	transactions map[string]*FiatTransaction
	balances     map[string]float64
	methods      map[string]*PaymentMethod
}

// FiatTransaction represents a fiat payment
type FiatTransaction struct {
	ID            string
	UserID        string
	Type          string
	Amount        float64
	Currency      string
	Method        string
	Status        string
	Fee           float64
	NetAmount     float64
	ExternalRef   string
	CreatedAt     time.Time
	CompletedAt   *time.Time
	FailureReason string
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID           string
	Name         string
	Type         string
	MinDeposit   float64
	MaxDeposit   float64
	MinWithdrawal float64
	MaxWithdrawal float64
	Fee          float64
	FeePercent   float64
	ProcessingTime string
	IsActive     bool
}

// NewFiatPaymentService creates a new fiat payment service
func NewFiatPaymentService() *FiatPaymentService {
	s := &FiatPaymentService{
		transactions: make(map[string]*FiatTransaction),
		balances:     make(map[string]float64),
		methods:      make(map[string]*PaymentMethod),
	}
	s.initializeMethods()
	return s
}

func (s *FiatPaymentService) initializeMethods() {
	// Credit/Debit Cards
	s.methods["visa"] = &PaymentMethod{ID: "visa", Name: "Visa", Type: "card", MinDeposit: 20, MaxDeposit: 10000, MinWithdrawal: 20, MaxWithdrawal: 50000, FeePercent: 2.5, ProcessingTime: "Instant", IsActive: true}
	s.methods["mastercard"] = &PaymentMethod{ID: "mastercard", Name: "Mastercard", Type: "card", MinDeposit: 20, MaxDeposit: 10000, MinWithdrawal: 20, MaxWithdrawal: 50000, FeePercent: 2.5, ProcessingTime: "Instant", IsActive: true}
	s.methods["amex"] = &PaymentMethod{ID: "amex", Name: "American Express", Type: "card", MinDeposit: 50, MaxDeposit: 25000, MinWithdrawal: 50, MaxWithdrawal: 50000, FeePercent: 3.0, ProcessingTime: "Instant", IsActive: true}

	// E-Wallets
	s.methods["skrill"] = &PaymentMethod{ID: "skrill", Name: "Skrill", Type: "ewallet", MinDeposit: 10, MaxDeposit: 50000, MinWithdrawal: 10, MaxWithdrawal: 50000, FeePercent: 2.5, ProcessingTime: "Instant", IsActive: true}
	s.methods["neteller"] = &PaymentMethod{ID: "neteller", Name: "Neteller", Type: "ewallet", MinDeposit: 10, MaxDeposit: 50000, MinWithdrawal: 10, MaxWithdrawal: 50000, FeePercent: 2.5, ProcessingTime: "Instant", IsActive: true}
	s.methods["paypal"] = &PaymentMethod{ID: "paypal", Name: "PayPal", Type: "ewallet", MinDeposit: 10, MaxDeposit: 10000, MinWithdrawal: 10, MaxWithdrawal: 10000, FeePercent: 3.0, ProcessingTime: "Instant", IsActive: true}
	s.methods["applepay"] = &PaymentMethod{ID: "applepay", Name: "Apple Pay", Type: "ewallet", MinDeposit: 10, MaxDeposit: 10000, MinWithdrawal: 10, MaxWithdrawal: 10000, FeePercent: 1.5, ProcessingTime: "Instant", IsActive: true}
	s.methods["googlepay"] = &PaymentMethod{ID: "googlepay", Name: "Google Pay", Type: "ewallet", MinDeposit: 10, MaxDeposit: 10000, MinWithdrawal: 10, MaxWithdrawal: 10000, FeePercent: 1.5, ProcessingTime: "Instant", IsActive: true}

	// Bank Transfer
	s.methods["bank_transfer"] = &PaymentMethod{ID: "bank_transfer", Name: "Bank Transfer", Type: "bank", MinDeposit: 100, MaxDeposit: 100000, MinWithdrawal: 100, MaxWithdrawal: 100000, Fee: 15, ProcessingTime: "1-3 days", IsActive: true}
	s.methods["swift"] = &PaymentMethod{ID: "swift", Name: "SWIFT Transfer", Type: "bank", MinDeposit: 500, MaxDeposit: 500000, MinWithdrawal: 500, MaxWithdrawal: 500000, Fee: 25, ProcessingTime: "2-5 days", IsActive: true}
	s.methods["sepa"] = &PaymentMethod{ID: "sepa", Name: "SEPA Transfer", Type: "bank", MinDeposit: 50, MaxDeposit: 50000, MinWithdrawal: 50, MaxWithdrawal: 50000, ProcessingTime: "1-2 days", IsActive: true}

	// Crypto on/off ramps
	s.methods["binance_pay"] = &PaymentMethod{ID: "binance_pay", Name: "Binance Pay", Type: "crypto", MinDeposit: 10, MaxDeposit: 100000, MinWithdrawal: 10, MaxWithdrawal: 100000, FeePercent: 1.0, ProcessingTime: "Instant", IsActive: true}
}

// CreateDeposit creates a new deposit request
func (s *FiatPaymentService) CreateDeposit(userID, methodID string, amount float64, currency string) (*FiatTransaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	method, ok := s.methods[methodID]
	if !ok {
		return nil, fmt.Errorf("payment method not found")
	}

	if amount < method.MinDeposit || amount > method.MaxDeposit {
		return nil, fmt.Errorf("amount outside allowed range")
	}

	fee := amount * method.FeePercent / 100
	netAmount := amount - fee

	tx := &FiatTransaction{
		ID: uuid.New().String(), UserID: userID, Type: "deposit",
		Amount: amount, Currency: currency, Method: methodID,
		Status: "pending", Fee: fee, NetAmount: netAmount,
		ExternalRef: fmt.Sprintf("TX%d%s", time.Now().Unix(), uuid.New().String()[:8]),
		CreatedAt: time.Now(),
	}

	s.transactions[tx.ID] = tx
	s.balances[userID] += netAmount

	return tx, nil
}

// CreateWithdrawal creates a new withdrawal request
func (s *FiatPaymentService) CreateWithdrawal(userID, methodID string, amount float64, currency string) (*FiatTransaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	method, ok := s.methods[methodID]
	if !ok {
		return nil, fmt.Errorf("payment method not found")
	}

	currentBalance := s.balances[userID]
	fee := amount*method.FeePercent/100 + method.Fee
	totalAmount := amount + fee

	if currentBalance < totalAmount {
		return nil, fmt.Errorf("insufficient balance")
	}

	tx := &FiatTransaction{
		ID: uuid.New().String(), UserID: userID, Type: "withdrawal",
		Amount: amount, Currency: currency, Method: methodID,
		Status: "pending", Fee: fee, NetAmount: amount - fee,
		ExternalRef: fmt.Sprintf("TX%d%s", time.Now().Unix(), uuid.New().String()[:8]),
		CreatedAt: time.Now(),
	}

	s.transactions[tx.ID] = tx
	s.balances[userID] -= totalAmount

	return tx, nil
}

// GetTransaction returns a transaction by ID
func (s *FiatPaymentService) GetTransaction(txID string) (*FiatTransaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.transactions[txID]
	if !ok {
		return nil, fmt.Errorf("transaction not found")
	}
	return tx, nil
}

// GetUserTransactions returns all transactions for a user
func (s *FiatPaymentService) GetUserTransactions(userID string) []FiatTransaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var txs []FiatTransaction
	for _, tx := range s.transactions {
		if tx.UserID == userID {
			txs = append(txs, *tx)
		}
	}
	return txs
}

// GetAvailableMethods returns all available payment methods
func (s *FiatPaymentService) GetAvailableMethods(txType string) []PaymentMethod {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var methods []PaymentMethod
	for _, m := range s.methods {
		if m.IsActive {
			methods = append(methods, *m)
		}
	}
	return methods
}

// GetBalance returns user balance
func (s *FiatPaymentService) GetBalance(userID string) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.balances[userID]
}
