package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Payment types
type PaymentType string

const (
	PaymentTypeDeposit    PaymentType = "deposit"
	PaymentTypeWithdrawal PaymentType = "withdrawal"
)

// Payment status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Payment method types
type PaymentMethod string

const (
	// Crypto
	MethodCrypto    PaymentMethod = "crypto"
	
	// Fiat
	MethodCreditCard PaymentMethod = "credit_card"
	MethodDebitCard  PaymentMethod = "debit_card"
	MethodSkrill    PaymentMethod = "skrill"
	MethodNeteller  PaymentMethod = "neteller"
	MethodPayPal    PaymentMethod = "paypal"
	MethodApplePay  PaymentMethod = "apple_pay"
	MethodGooglePay PaymentMethod = "google_pay"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodInstadebit PaymentMethod = "instadebit"
	MethodiDebit    PaymentMethod = "idebit"
)

// Payment request
type PaymentRequest struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	Type            PaymentType       `json:"type"`
	Method          PaymentMethod     `json:"method"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Status          PaymentStatus     `json:"status"`
	TransactionID   string            `json:"transaction_id,omitempty"`
	ExternalID      string            `json:"external_id,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ProcessedAt     *time.Time       `json:"processed_at,omitempty"`
}

// Card info (for payment processing)
type CardInfo struct {
	Number          string `json:"number"`
	CVV             string `json:"cvv"`
	ExpiryMonth     int    `json:"expiry_month"`
	ExpiryYear      int    `json:"expiry_year"`
	CardholderName string `json:"cardholder_name"`
	BillingAddress  Address `json:"billing_address"`
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// Bank info for transfers
type BankInfo struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	SWIFTBIC      string `json:"swift_bic"`
	BankName      string `json:"bank_name"`
	BankAddress   string `json:"bank_address"`
	Country       string `json:"country"`
}

// Payment processor interface
type PaymentProcessor interface {
	ProcessDeposit(req *PaymentRequest, cardInfo *CardInfo) (*PaymentResult, error)
	ProcessWithdrawal(req *PaymentRequest, bankInfo *BankInfo) (*PaymentResult, error)
	GetBalance() (float64, error)
	GetSupportedCurrencies() []string
}

// Payment result
type PaymentResult struct {
	Success       bool              `json:"success"`
	TransactionID string            `json:"transaction_id"`
	Message       string            `json:"message"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// Payment service
type PaymentService struct {
	processors     map[PaymentMethod]PaymentProcessor
	pendingPayments map[string]*PaymentRequest
}

// NewPaymentService creates a new payment service
func NewPaymentService() *PaymentService {
	return &PaymentService{
		processors:     make(map[PaymentMethod]PaymentProcessor),
		pendingPayments: make(map[string]*PaymentRequest),
	}
}

// RegisterProcessor registers a payment processor
func (s *PaymentService) RegisterProcessor(method PaymentMethod, processor PaymentProcessor) {
	s.processors[method] = processor
}

// ProcessPayment processes a payment
func (s *PaymentService) ProcessPayment(req *PaymentRequest, data interface{}) (*PaymentResult, error) {
	processor, ok := s.processors[req.Method]
	if !ok {
		return nil, fmt.Errorf("unsupported payment method: %s", req.Method)
	}

	var result *PaymentResult
	var err error

	switch req.Type {
	case PaymentTypeDeposit:
		cardInfo, ok := data.(*CardInfo)
		if !ok {
			return nil, fmt.Errorf("invalid card info")
		}
		result, err = processor.ProcessDeposit(req, cardInfo)
	case PaymentTypeWithdrawal:
		bankInfo, ok := data.(*BankInfo)
		if !ok {
			return nil, fmt.Errorf("invalid bank info")
		}
		result, err = processor.ProcessWithdrawal(req, bankInfo)
	default:
		return nil, fmt.Errorf("unknown payment type: %s", req.Type)
	}

	if err != nil {
		req.Status = PaymentStatusFailed
		return nil, err
	}

	req.TransactionID = result.TransactionID
	req.Status = PaymentStatusCompleted
	now := time.Now()
	req.ProcessedAt = &now

	return result, nil
}

// ============== Stripe-like Credit Card Processor ==============

type CreditCardProcessor struct {
	apiKey        string
	merchantID    string
 webhookSecret string
}

func NewCreditCardProcessor(apiKey, merchantID, webhookSecret string) *CreditCardProcessor {
	return &CreditCardProcessor{
		apiKey:        apiKey,
		merchantID:    merchantID,
		webhookSecret: webhookSecret,
	}
}

func (p *CreditCardProcessor) ProcessDeposit(req *PaymentRequest, cardInfo *CardInfo) (*PaymentResult, error) {
	// Validate card
	if err := p.validateCard(cardInfo); err != nil {
		return nil, err
	}

	// In production, this would call Stripe API
	// For now, simulate processing
	transactionID := generateTransactionID()

	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Payment processed successfully",
		Metadata:      map[string]string{"last4": cardInfo.Number[len(cardInfo.Number)-4:]},
	}, nil
}

func (p *CreditCardProcessor) ProcessWithdrawal(req *PaymentRequest, bankInfo *BankInfo) (*PaymentResult, error) {
	// Process bank transfer
	transactionID := generateTransactionID()

	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Withdrawal initiated",
	}, nil
}

func (p *CreditCardProcessor) GetBalance() (float64, error) {
	// Would check processor balance
	return 10000.00, nil
}

func (p *CreditCardProcessor) GetSupportedCurrencies() []string {
	return []string{"USD", "EUR", "GBP", "CAD", "AUD"}
}

func (p *CreditCardProcessor) validateCard(card *CardInfo) error {
	if len(card.Number) < 13 || len(card.Number) > 19 {
		return fmt.Errorf("invalid card number")
	}

	if card.CVV == "" || len(card.CVV) < 3 {
		return fmt.Errorf("invalid CVV")
	}

	now := time.Now()
	expiry := time.Date(card.ExpiryYear, time.Month(card.ExpiryMonth), 1, 0, 0, 0, 0, time.UTC)
	if expiry.Before(now) {
		return fmt.Errorf("card expired")
	}

	return nil
}

// ============== E-Wallet Processors ==============

type EWalletProcessor struct {
	apiKey     string
	merchantID string
	ewallet    PaymentMethod
}

func NewSkrillProcessor(apiKey, merchantID string) *EWalletProcessor {
	return &EWalletProcessor{
		apiKey:     apiKey,
		merchantID: merchantID,
		ewallet:    MethodSkrill,
	}
}

func (p *EWalletProcessor) ProcessDeposit(req *PaymentRequest, cardInfo *CardInfo) (*PaymentResult, error) {
	transactionID := generateTransactionID()
	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Skrill deposit processed",
	}, nil
}

func (p *EWalletProcessor) ProcessWithdrawal(req *PaymentRequest, bankInfo *BankInfo) (*PaymentResult, error) {
	transactionID := generateTransactionID()
	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Skrill withdrawal processed",
	}, nil
}

func (p *EWalletProcessor) GetBalance() (float64, error) {
	return 50000.00, nil
}

func (p *EWalletProcessor) GetSupportedCurrencies() []string {
	return []string{"USD", "EUR", "GBP"}
}

type NetellerProcessor struct {
	EWalletProcessor
}

func NewNetellerProcessor(apiKey, merchantID string) *NetellerProcessor {
	return &NetellerProcessor{
		EWalletProcessor: EWalletProcessor{
			apiKey:     apiKey,
			merchantID: merchantID,
			ewallet:    MethodNeteller,
		},
	}
}

// ============== Crypto Processor ==============

type CryptoPaymentProcessor struct {
	confirmations map[string]int
	supportedCoins []string
}

func NewCryptoPaymentProcessor() *CryptoPaymentProcessor {
	return &CryptoPaymentProcessor{
		confirmations: map[string]int{
			"BTC": 3,
			"ETH": 12,
			"USDT": 15,
			"LTC": 6,
		},
		supportedCoins: []string{"BTC", "ETH", "USDT", "LTC", "XRP", "TRX"},
	}
}

func (p *CryptoPaymentProcessor) ProcessDeposit(req *PaymentRequest, _ *CardInfo) (*PaymentResult, error) {
	// Generate deposit address (would use actual crypto nodes)
	transactionID := generateTransactionID()
	address := generateCryptoAddress(req.Currency)

	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Crypto deposit address generated",
		Metadata: map[string]string{
			"address":      address,
			"confirmations": fmt.Sprintf("%d", p.confirmations[req.Currency]),
		},
	}, nil
}

func (p *CryptoPaymentProcessor) ProcessWithdrawal(req *PaymentRequest, _ *BankInfo) (*PaymentResult, error) {
	transactionID := generateTransactionID()
	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Crypto withdrawal initiated",
	}, nil
}

func (p *CryptoPaymentProcessor) GetBalance() (float64, error) {
	return 1000000.00, nil
}

func (p *CryptoPaymentProcessor) GetSupportedCurrencies() []string {
	return p.supportedCoins
}

func (p *CryptoPaymentProcessor) GenerateDepositAddress(currency string) (string, error) {
	return generateCryptoAddress(currency), nil
}

// ============== Apple Pay / Google Pay ==============

type DigitalWalletProcessor struct {
	processorType PaymentMethod
	merchantID    string
}

func NewApplePayProcessor(merchantID string) *DigitalWalletProcessor {
	return &DigitalWalletProcessor{
		processorType: MethodApplePay,
		merchantID:    merchantID,
	}
}

func (p *DigitalWalletProcessor) ProcessDeposit(req *PaymentRequest, cardInfo *CardInfo) (*PaymentResult, error) {
	transactionID := generateTransactionID()
	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Apple Pay payment processed",
	}, nil
}

func (p *DigitalWalletProcessor) ProcessWithdrawal(req *PaymentRequest, bankInfo *BankInfo) (*PaymentResult, error) {
	return nil, fmt.Errorf("withdrawal not supported for Apple Pay")
}

func (p *DigitalWalletProcessor) GetBalance() (float64, error) {
	return 100000.00, nil
}

func (p *DigitalWalletProcessor) GetSupportedCurrencies() []string {
	return []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD"}
}

type GooglePayProcessor struct {
	DigitalWalletProcessor
}

func NewGooglePayProcessor(merchantID string) *GooglePayProcessor {
	return &GooglePayProcessor{
		DigitalWalletProcessor: DigitalWalletProcessor{
			processorType: MethodGooglePay,
			merchantID:    merchantID,
		},
	}
}

// ============== Bank Transfer Processor ==============

type BankTransferProcessor struct {
	bankCode    string
	correspondentBank string
}

func NewBankTransferProcessor(bankCode, correspondentBank string) *BankTransferProcessor {
	return &BankTransferProcessor{
		bankCode:        bankCode,
		correspondentBank: correspondentBank,
	}
}

func (p *BankTransferProcessor) ProcessDeposit(req *PaymentRequest, bankInfo *BankInfo) (*PaymentResult, error) {
	transactionID := generateTransactionID()
	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Bank transfer deposit initiated",
		Metadata: map[string]string{
			"reference": generateReference(),
		},
	}, nil
}

func (p *BankTransferProcessor) ProcessWithdrawal(req *PaymentRequest, bankInfo *BankInfo) (*PaymentResult, error) {
	transactionID := generateTransactionID()
	return &PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Message:       "Bank transfer withdrawal initiated",
		Metadata: map[string]string{
			"processing_time": "3-5 business days",
		},
	}, nil
}

func (p *BankTransferProcessor) GetBalance() (float64, error) {
	return 500000.00, nil
}

func (p *BankTransferProcessor) GetSupportedCurrencies() []string {
	return []string{"USD", "EUR", "GBP", "CHF", "CAD", "AUD", "JPY"}
}

// ============== Helpers ==============

func generateTransactionID() string {
	return fmt.Sprintf("TXN%d%d", time.Now().Unix(), time.Now().Nanosecond()%10000)
}

func generateCryptoAddress(currency string) string {
	switch currency {
	case "BTC":
		return "bc1q" + generateRandomString(38)
	case "ETH":
		return "0x" + generateRandomString(40)
	case "USDT":
		return "0x" + generateRandomString(40)
	default:
		return generateRandomString(32)
	}
}

func generateRandomString(length int) string {
	chars := "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[i%len(chars)]
	}
	return string(result)
}

func generateReference() string {
	return fmt.Sprintf("TC%d%d", time.Now().Unix(), time.Now().Nanosecond()%100000)
}

// VerifyWebhook verifies a webhook signature
func VerifyWebhook(payload []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

// Webhook event types
type WebhookEvent struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

func ParseWebhookEvent(payload []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
