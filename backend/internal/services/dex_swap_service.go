package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DEXQuote represents a swap quote from DEX
type DEXQuote struct {
	FromToken string
	ToToken   string
	FromAmount float64
	ToAmount   float64
	Rate      float64
	Slippage  float64
	GasEstimate float64
	Expiry    time.Time
}

// SwapOrder represents a swap order
type SwapOrder struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	FromToken    string    `json:"from_token"`
	ToToken      string    `json:"to_token"`
	FromAmount   float64   `json:"from_amount"`
	ToAmount     float64   `json:"to_amount"`
	Status       string    `json:"status"` // pending, processing, completed, failed
	TxHash       string    `json:"tx_hash"`
	FromAddress  string    `json:"from_address"`
	ToAddress    string    `json:"to_address"`
	CreatedAt    time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// DEXSwapService handles in-wallet crypto swaps
type DEXSwapService struct {
	httpClient *http.Client
	apiKeys    map[string]string
}

// NewDEXSwapService creates a new DEX swap service
func NewDEXSwapService() *DEXSwapService {
	return &DEXSwapService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKeys:    make(map[string]string),
	}
}

// GetSupportedTokens returns supported tokens for swapping
func (s *DEXSwapService) GetSupportedTokens() []map[string]string {
	return []map[string]string{
		{"symbol": "BTC", "name": "Bitcoin", "address": "", "network": "bitcoin"},
		{"symbol": "ETH", "name": "Ethereum", "address": "0x0000000000000000000000000000000000000000", "network": "ethereum"},
		{"symbol": "USDT", "name": "Tether USD", "address": "0xdAC17F958D2ee523a2206206994597C13D831ec7", "network": "ethereum"},
		{"symbol": "USDC", "name": "USD Coin", "address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "network": "ethereum"},
		{"symbol": "BNB", "name": "BNB", "address": "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bd095", "network": "bsc"},
		{"symbol": "TRX", "name": "Tron", "address": "", "network": "tron"},
		{"symbol": "SOL", "name": "Solana", "address": "", "network": "solana"},
		{"symbol": "MATIC", "name": "Polygon", "address": "0x0000000000000000000000000000000000000000", "network": "polygon"},
	}
}

// GetSwapQuote gets a swap quote
func (s *DEXSwapService) GetSwapQuote(fromToken, toToken string, amount float64) (*DEXQuote, error) {
	// In production, would call actual DEX APIs (Uniswap, 1inch, etc.)
	// For demo, return mock quote
	
	// Simulate API call delay
	time.Sleep(100 * time.Millisecond)
	
	// Mock exchange rates
	rates := map[string]float64{
		"ETH-BTC": 0.02,
		"BTC-ETH": 50.0,
		"USDT-ETH": 0.0004,
		"ETH-USDT": 2500.0,
		"USDT-BTC": 0.000016,
		"BTC-USDT": 62500.0,
		"BNB-ETH": 0.3,
		"ETH-BNB": 3.33,
	}
	
	key := fromToken + "-" + toToken
	rate, ok := rates[key]
	if !ok {
		// Default rate (simplified)
		rate = 1.0
	}
	
	toAmount := amount * rate
	
	quote := &DEXQuote{
		FromToken: fromToken,
		ToToken:   toToken,
		FromAmount: amount,
		ToAmount:   toAmount,
		Rate:      rate,
		Slippage:  0.5, // 0.5% estimated slippage
		GasEstimate: 0.01, // ~$10 in ETH terms
		Expiry:    time.Now().Add(30 * time.Second),
	}
	
	return quote, nil
}

// ExecuteSwap executes a swap
func (s *DEXSwapService) ExecuteSwap(userID, fromToken, toToken string, amount float64) (*SwapOrder, error) {
	// Get quote
	quote, err := s.GetSwapQuote(fromToken, toToken, amount)
	if err != nil {
		return nil, err
	}
	
	// Validate amount
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}
	
	// Create order
	order := &SwapOrder{
		ID:           uuid.New().String(),
		UserID:       userID,
		FromToken:    fromToken,
		ToToken:      toToken,
		FromAmount:   amount,
		ToAmount:     quote.ToAmount,
		Status:       "processing",
		FromAddress:  "", // Would be user's wallet
		ToAddress:    "", // Would be output address
		CreatedAt:    time.Now(),
	}
	
	// In production, would:
	// 1. Validate user balance
	// 2. Call DEX API to execute swap
	// 3. Monitor transaction
	// 4. Credit user account
	
	// Simulate successful swap
	order.Status = "completed"
	now := time.Now()
	order.CompletedAt = &now
	
	// Generate mock tx hash
	order.TxHash = "0x" + fmt.Sprintf("%x", uuid.New().ID())
	
	return order, nil
}

// GetSwapHistory returns user's swap history
func (s *DEXSwapService) GetSwapHistory(userID string, limit int) ([]SwapOrder, error) {
	// In production, would query database
	// Return mock data
	orders := []SwapOrder{}
	
	return orders, nil
}

// ============ 1inch Integration (Production) ============

// OneInchQuote represents 1inch API quote
type OneInchQuote struct {
	FromToken  string  `json:"fromToken"`
	ToToken    string  `json:"toToken"`
	FromTokenDecimals int  `json:"fromTokenDecimals"`
	ToTokenDecimals   int  `json:"toTokenDecimals"`
	FromTokenSymbol   string  `json:"fromTokenSymbol"`
	ToTokenSymbol     string  `json:"toTokenSymbol"`
	ToAmount         string  `json:"toAmount"`
	FromAmount       string  `json:"fromAmount"`
	Slippage         float64 `json:"slippage"`
}

// OneInchSwap represents 1inch swap transaction
type OneInchSwap struct {
	FromToken  string `json:"fromToken"`
	ToToken    string `json:"toToken"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Data      string `json:"data"`
}

// GetOneInchQuote gets quote from 1inch API
func (s *DEXSwapService) GetOneInchQuote(chainID int, fromToken, toToken string, amount string) (*OneInchQuote, error) {
	// In production, would use actual 1inch API
	apiURL := fmt.Sprintf("https://api.1inch.io/v5.0/%d/quote?fromTokenAddress=%s&toTokenAddress=%s&amount=%s",
		chainID, fromToken, toToken, amount)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var quote OneInchQuote
	if err := json.Unmarshal(body, &quote); err != nil {
		return nil, err
	}
	
	return &quote, nil
}

// BuildOneInchSwap builds swap transaction from 1inch
func (s *DEXSwapService) BuildOneInchSwap(chainID int, fromToken, toToken, fromAddress, amount, slippage string) (*OneInchSwap, error) {
	apiURL := fmt.Sprintf("https://api.1inch.io/v5.0/%d/swap?fromTokenAddress=%s&toTokenAddress=%s&fromAddress=%s&amount=%s&slippage=%s&disableEstimate=true",
		chainID, fromToken, toToken, fromAddress, amount, slippage)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	tx, ok := result["tx"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no transaction returned")
	}
	
	swap := &OneInchSwap{
		FromToken: fromToken,
		ToToken:   toToken,
		From:      tx["from"].(string),
		To:        tx["to"].(string),
		Value:     tx["value"].(string),
		Data:      tx["data"].(string),
	}
	
	return swap, nil
}

// ============ Uniswap Integration (Production) ============

// UniswapQuote represents Uniswap quote
type UniswapQuote struct {
	AmountIn        string `json:"amountIn"`
	AmountOut       string `json:"amountOut"`
	AmountOutMin    string `json:"amountOutMin"`
	Path            []string `json:"path"`
	GasEstimate     string `json:"gasEstimate"`
}

// GetUniswapQuote gets quote from Uniswap
func (s *DEXSwapService) GetUniswapQuote(tokenIn, tokenOut string, amountIn *big.Int) (*UniswapQuote, error) {
	// In production, would use Uniswap SDK or API
	// Simplified mock implementation
	quote := &UniswapQuote{
		AmountIn:    amountIn.String(),
		AmountOut:   new(big.Int).Div(amountIn, big.NewInt(2500)).String(),
		AmountOutMin: new(big.Int).Div(amountIn, big.NewInt(2600)).String(),
		Path:       []string{tokenIn, tokenOut},
		GasEstimate: "210000",
	}
	
	return quote, nil
}

// BuildUniswapRoute builds swap route for Uniswap
func (s *DEXSwapService) BuildUniswapRoute(tokenIn, tokenOut string, amountIn *big.Int, to string) ([]byte, error) {
	// Would build smart contract swap data
	// Simplified
	return []byte{}, nil
}

// ============ Price Aggregation ============

// Price represents aggregated price
type Price struct {
	Token   string  `json:"token"`
	USDPrice float64 `json:"usd_price"`
	Change24h float64 `json:"change_24h"`
	Source   string  `json:"source"`
}

// GetPrices gets current prices for all tokens
func (s *DEXSwapService) GetPrices() ([]Price, error) {
	prices := []Price{
		{"BTC", 62500.0, 2.5, "aggregated"},
		{"ETH", 2500.0, 1.8, "aggregated"},
		{"USDT", 1.0, 0.01, "pegged"},
		{"USDC", 1.0, 0.0, "pegged"},
		{"BNB", 550.0, 3.2, "aggregated"},
		{"TRX", 0.12, -0.5, "aggregated"},
		{"SOL", 145.0, 5.2, "aggregated"},
		{"MATIC", 0.85, 1.2, "aggregated"},
	}
	
	return prices, nil
}

// GetHistoricalPrices gets historical prices
func (s *DEXSwapService) GetHistoricalPrices(token string, interval string, limit int) ([]map[string]interface{}, error) {
	// Would fetch from price API (CoinGecko, CoinMarketCap)
	// Simplified mock
	prices := []map[string]interface{}{}
	
	for i := 0; i < limit; i++ {
		prices = append(prices, map[string]interface{}{
			"timestamp": time.Now().Add(-time.Duration(i) * 24 * time.Hour).Unix(),
			"price":    100.0 + float64(i),
		})
	}
	
	return prices, nil
}

// ValidateSwapAmount validates swap amount against limits
func (s *DEXSwapService) ValidateSwapAmount(amount float64, token string) error {
	limits := map[string]map[string]float64{
		"min": {
			"BTC":  0.001,
			"ETH":  0.01,
			"USDT": 10,
			"USDC": 10,
			"BNB":  0.1,
		},
		"max": {
			"BTC":  10.0,
			"ETH":  100.0,
			"USDT": 100000,
			"USDC": 100000,
			"BNB":  500.0,
		},
	}
	
	if min, ok := limits["min"][token]; ok {
		if amount < min {
			return fmt.Errorf("minimum swap amount is %f %s", min, token)
		}
	}
	
	if max, ok := limits["max"][token]; ok {
		if amount > max {
			return fmt.Errorf("maximum swap amount is %f %s", max, token)
		}
	}
	
	return nil
}

// EstimateSwapTime estimates time for swap completion
func (s *DEXSwapService) EstimateSwapTime(fromToken, toToken string) string {
	// Different chains have different confirmation times
	times := map[string]map[string]string{
		"bitcoin":   {"avg": "30 minutes", "max": "2 hours"},
		"ethereum":   {"avg": "5 minutes", "max": "30 minutes"},
		"bsc":       {"avg": "3 minutes", "max": "15 minutes"},
		"polygon":   {"avg": "2 minutes", "max": "10 minutes"},
		"tron":      {"avg": "5 minutes", "max": "15 minutes"},
		"solana":    {"avg": "30 seconds", "max": "2 minutes"},
	}
	
	network := s.getNetwork(fromToken)
	if times, ok := times[network]; ok {
		return times["avg"]
	}
	
	return "5 minutes"
}

func (s *DEXSwapService) getNetwork(token string) string {
	networks := map[string]string{
		"BTC":  "bitcoin",
		"ETH":  "ethereum",
		"USDT": "ethereum",
		"USDC": "ethereum",
		"BNB":  "bsc",
		"TRX":  "tron",
		"SOL":  "solana",
		"MATIC": "polygon",
	}
	
	if net, ok := networks[token]; ok {
		return net
	}
	
	return "ethereum"
}

// ContextWithTimeout creates context with timeout
func ContextWithTimeout(seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds) * time.Second)
}

// ParseTokenAddress parses token symbol to address
func ParseTokenAddress(symbol string) (string, string, error) {
	symbol = strings.ToUpper(symbol)
	
	addresses := map[string][2]string{
		"ETH":   {"0x0000000000000000000000000000000000000000", "ethereum"},
		"USDT":  {"0xdAC17F958D2ee523a2206206994597C13D831ec7", "ethereum"},
		"USDC":  {"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "ethereum"},
		"BNB":   {"0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bd095", "bsc"},
		"MATIC": {"0x0000000000000000000000000000000000000000", "polygon"},
	}
	
	if info, ok := addresses[symbol]; ok {
		return info[0], info[1], nil
	}
	
	return "", "", fmt.Errorf("unsupported token: %s", symbol)
}
