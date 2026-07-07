package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tigercasino/backend/internal/models"
)

// ProviderIntegrationService handles game provider API integrations
type ProviderIntegrationService struct {
	db          *gorm.DB
	mu          sync.RWMutex
	providers   map[string]*ProviderConfig
	httpClient  *http.Client
}

type ProviderConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	APIKey      string            `json:"api_key"`
	APISecret   string            `json:"api_secret"`
	BaseURL     string            `json:"base_url"`
	WebhookURL  string            `json:"webhook_url"`
	IsActive    bool              `json:"is_active"`
	Currencies []string          `json:"currencies"`
	GameTypes   []string         `json:"game_types"`
	Categories []string         `json:"categories"`
	Config      map[string]string `json:"config"`
}

type ProviderGame struct {
	ExternalID   string            `json:"external_id"`
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	Type         string            `json:"type"`
	Category     string            `json:"category"`
	RTP          float64           `json:"rtp"`
	MinBet       float64           `json:"min_bet"`
	MaxBet       float64           `json:"max_bet"`
	Volatility   string            `json:"volatility"`
	ThumbnailURL string            `json:"thumbnail_url"`
	MobileReady bool              `json:"mobile_ready"`
	Tags         []string         `json:"tags"`
}

type GameSession struct {
	SessionID    string            `json:"session_id"`
	UserID       string            `json:"user_id"`
	Provider     string            `json:"provider"`
	GameID       string            `json:"game_id"`
	GameName     string            `json:"game_name"`
	Currency     string            `json:"currency"`
	Balance      float64           `json:"balance"`
	Stake        float64           `json:"stake"`
	Win          float64            `json:"win"`
	Status       string            `json:"status"` // active, completed, cancelled
	StartedAt    time.Time         `json:"started_at"`
	EndedAt      *time.Time       `json:"ended_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type Transaction struct {
	TransactionID string    `json:"transaction_id"`
	SessionID    string    `json:"session_id"`
	UserID       string    `json:"user_id"`
	Provider     string    `json:"provider"`
	Type         string    `json:"type"` // bet, win, refund
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"` // pending, completed, failed
	ExternalRef  string    `json:"external_ref"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewProviderIntegrationService(db *gorm.DB) *ProviderIntegrationService {
	s := &ProviderIntegrationService{
		db:         db,
		providers:  make(map[string]*ProviderConfig),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Initialize default provider configurations
	s.initializeProviders()

	return s
}

func (s *ProviderIntegrationService) initializeProviders() {
	// Pragmatic Play
	s.providers["pragmatic_play"] = &ProviderConfig{
		ID:          "pragmatic_play",
		Name:        "Pragmatic Play",
		BaseURL:     "https://api.pragmaticplay.com",
		IsActive:    false, // Requires API key
		Currencies:  []string{"USD", "EUR", "BTC", "ETH", "USDT"},
		GameTypes:   []string{"slot", "live_casino", "virtual_sports", " Bingo", "lottery"},
		Categories: []string{"slots", "table_games", "game_shows"},
	}

	// Evolution Gaming
	s.providers["evolution"] = &ProviderConfig{
		ID:          "evolution",
		Name:        "Evolution Gaming",
		BaseURL:     "https://api.evolutiongaming.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "GBP", "BTC", "ETH"},
		GameTypes:   []string{"live_casino", "game_shows"},
		Categories:  []string{"blackjack", "roulette", "baccarat", "poker", "game_shows"},
	}

	// NetEnt
	s.providers["netent"] = &ProviderConfig{
		ID:          "netent",
		Name:        "NetEnt",
		BaseURL:     "https://api.netent.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "SEK", "GBP"},
		GameTypes:   []string{"slot", "table_games"},
		Categories:  []string{"slots", "jackpots"},
	}

	// Play'n GO
	s.providers["playngo"] = &ProviderConfig{
		ID:          "playngo",
		Name:        "Play'n GO",
		BaseURL:     "https://api.playngo.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "GBP", "NOK", "SEK"},
		GameTypes:   []string{"slot", "table_games", "scratch"},
		Categories:  []string{"slots", "table_games"},
	}

	// Microgaming
	s.providers["microgaming"] = &ProviderConfig{
		ID:          "microgaming",
		Name:        "Microgaming",
		BaseURL:     "https://api.microgaming.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "GBP", "CAD", "AUD"},
		GameTypes:   []string{"slot", "table_games", "poker", "bingo"},
		Categories:  []string{"slots", "jackpots", "table_games"},
	}

	// Red Tiger
	s.providers["redtiger"] = &ProviderConfig{
		ID:          "redtiger",
		Name:        "Red Tiger Gaming",
		BaseURL:     "https://api.redtigergaming.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "GBP", "SEK", "NOK"},
		GameTypes:   []string{"slot", "jackpot"},
		Categories:  []string{"slots", "daily_jackpots"},
	}

	// Yggdrasil
	s.providers["yggdrasil"] = &ProviderConfig{
		ID:          "yggdrasil",
		Name:        "Yggdrasil Gaming",
		BaseURL:     "https://api.yggdrasilgaming.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "GBP", "NOK", "SEK", "BTC"},
		GameTypes:   []string{"slot", "jackpot"},
		Categories:  []string{"slots", "jackpots"},
	}

	// BGaming
	s.providers["bgaming"] = &ProviderConfig{
		ID:          "bgaming",
		Name:        "BGaming",
		BaseURL:     "https://api.bgaming.com",
		IsActive:    true, // Demo mode
		Currencies:  []string{"USD", "EUR", "BTC", "ETH", "USDT", "DOGE"},
		GameTypes:   []string{"slot", "table_games", "crash"},
		Categories:  []string{"slots", "table_games", "crash_games"},
	}

	// Spribe
	s.providers["spribe"] = &ProviderConfig{
		ID:          "spribe",
		Name:        "Spribe",
		BaseURL:     "https://api.spribe.co",
		IsActive:    true, // Demo mode
		Currencies:  []string{"USD", "EUR", "USDT", "BTC", "ETH"},
		GameTypes:   []string{"crash", "mines", "plinko", "hilo", "dice", "keno"},
		Categories:  []string{"crash_games", "instant_win"},
	}

	// Hacksaw Gaming
	s.providers["hacksaw"] = &ProviderConfig{
		ID:          "hacksaw",
		Name:        "Hacksaw Gaming",
		BaseURL:     "https://api.hacksawgaming.com",
		IsActive:    true, // Demo mode
		Currencies:  []string{"USD", "EUR", "GBP", "NOK", "SEK", "BTC"},
		GameTypes:   []string{"slot", "scratch", "instant"},
		Categories:  []string{"slots", "scratch_cards", "instant_win"},
	}

	// Nolimit City
	s.providers["nolimitcity"] = &ProviderConfig{
		ID:          "nolimitcity",
		Name:        "Nolimit City",
		BaseURL:     "https://api.nolimitcity.com",
		IsActive:    false,
		Currencies:  []string{"USD", "EUR", "GBP", "BTC", "ETH"},
		GameTypes:   []string{"slot"},
		Categories:  []string{"slots"},
	}
}

// ConfigureProvider configures a provider with API credentials
func (s *ProviderIntegrationService) ConfigureProvider(providerID, apiKey, apiSecret, baseURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providers[providerID]
	if !ok {
		return fmt.Errorf("provider not found: %s", providerID)
	}

	provider.APIKey = apiKey
	provider.APISecret = apiSecret
	provider.BaseURL = baseURL

	if apiKey != "" && apiSecret != "" {
		provider.IsActive = true
	}

	return nil
}

// GetProvider returns provider configuration
func (s *ProviderIntegrationService) GetProvider(providerID string) (*ProviderConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, fmt.Errorf("provider not found")
	}
	return provider, nil
}

// GetProviders returns all providers
func (s *ProviderIntegrationService) GetProviders() []*ProviderConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*ProviderConfig
	for _, provider := range s.providers {
		result = append(result, provider)
	}
	return result
}

// GetActiveProviders returns only active providers
func (s *ProviderIntegrationService) GetActiveProviders() []*ProviderConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*ProviderConfig
	for _, provider := range s.providers {
		if provider.IsActive {
			result = append(result, provider)
		}
	}
	return result
}

// GetProviderGames fetches games from a provider
func (s *ProviderIntegrationService) GetProviderGames(providerID string) ([]ProviderGame, error) {
	provider, err := s.GetProvider(providerID)
	if err != nil {
		return nil, err
	}

	// If not active, return demo games
	if !provider.IsActive {
		return s.getDemoGames(providerID), nil
	}

	// Fetch from API (in production, make actual API call)
	// For now, return demo games
	return s.getDemoGames(providerID), nil
}

func (s *ProviderIntegrationService) getDemoGames(providerID string) []ProviderGame {
	switch providerID {
	case "pragmatic_play":
		return []ProviderGame{
			{ExternalID: "pp_slots_001", Name: "Sweet Bonanza", Provider: "pragmatic_play", Type: "slot", Category: "slots", RTP: 96.48, MinBet: 0.20, MaxBet: 100, MobileReady: true, Tags: []string{"candy", "cascade", "free_spins"}},
			{ExternalID: "pp_slots_002", Name: "Wolf Gold", Provider: "pragmatic_play", Type: "slot", Category: "slots", RTP: 96.01, MinBet: 0.25, MaxBet: 125, MobileReady: true, Tags: []string{"wilds", "jackpot"}},
			{ExternalID: "pp_slots_003", Name: "The Dog House", Provider: "pragmatic_play", Type: "slot", Category: "slots", RTP: 96.51, MinBet: 0.20, MaxBet: 100, MobileReady: true, Tags: []string{"dogs", "sticky_wilds"}},
			{ExternalID: "pp_slots_004", Name: "Gates of Olympus", Provider: "pragmatic_play", Type: "slot", Category: "slots", RTP: 95.51, MinBet: 0.20, MaxBet: 100, MobileReady: true, Tags: []string{"greek", "multiplier", "free_spins"}},
			{ExternalID: "pp_slots_005", Name: "Big Bass Bonanza", Provider: "pragmatic_play", Type: "slot", Category: "slots", RTP: 96.71, MinBet: 0.10, MaxBet: 125, MobileReady: true, Tags: []string{"fishing", "free_spins"}},
		}
	case "evolution":
		return []ProviderGame{
			{ExternalID: "evo_blk_001", Name: "Lightning Blackjack", Provider: "evolution", Type: "live_blackjack", Category: "blackjack", RTP: 99.56, MinBet: 10, MaxBet: 5000, MobileReady: true},
			{ExternalID: "evo_rou_001", Name: "Lightning Roulette", Provider: "evolution", Type: "live_roulette", Category: "roulette", RTP: 97.10, MinBet: 1, MaxBet: 10000, MobileReady: true},
			{ExternalID: "evo_bac_001", Name: "Speed Baccarat", Provider: "evolution", Type: "live_baccarat", Category: "baccarat", RTP: 98.94, MinBet: 1, MaxBet: 15000, MobileReady: true},
			{ExternalID: "evo_pok_001", Name: "Casino Hold'em", Provider: "evolution", Type: "live_poker", Category: "poker", RTP: 97.84, MinBet: 5, MaxBet: 1000, MobileReady: true},
		}
	case "bgaming":
		return []ProviderGame{
			{ExternalID: "bg_pli_001", Name: "Plinko", Provider: "bgaming", Type: "crash", Category: "crash_games", RTP: 98.0, MinBet: 0.10, MaxBet: 100, MobileReady: true},
			{ExternalID: "bg_min_001", Name: "Mines", Provider: "bgaming", Type: "mines", Category: "crash_games", RTP: 97.0, MinBet: 0.10, MaxBet: 100, MobileReady: true},
			{ExternalID: "bg_dic_001", Name: "Dice", Provider: "bgaming", Type: "dice", Category: "crash_games", RTP: 99.0, MinBet: 0.01, MaxBet: 1000, MobileReady: true},
			{ExternalID: "bg_ken_001", Name: "Keno", Provider: "bgaming", Type: "lottery", Category: "lottery", RTP: 95.0, MinBet: 0.10, MaxBet: 100, MobileReady: true},
		}
	case "spribe":
		return []ProviderGame{
			{ExternalID: "sp_avi_001", Name: "Aviator", Provider: "spribe", Type: "crash", Category: "crash_games", RTP: 97.0, MinBet: 0.10, MaxBet: 100, MobileReady: true, Tags: []string{"crash", "multiplier"}},
			{ExternalID: "sp_min_001", Name: "Mines", Provider: "spribe", Type: "mines", Category: "crash_games", RTP: 97.0, MinBet: 0.10, MaxBet: 100, MobileReady: true, Tags: []string{"mines", "grid"}},
			{ExternalID: "sp_pli_001", Name: "Plinko", Provider: "spribe", Type: "plinko", Category: "crash_games", RTP: 98.0, MinBet: 0.10, MaxBet: 100, MobileReady: true, Tags: []string{"plinko", "balls"}},
			{ExternalID: "sp_hil_001", Name: "Hi Lo", Provider: "spribe", Type: "hilo", Category: "card_games", RTP: 96.50, MinBet: 0.10, MaxBet: 100, MobileReady: true, Tags: []string{"hilo", "cards"}},
		}
	case "hacksaw":
		return []ProviderGame{
			{ExternalID: "hack_stick_001", Name: "Stick 'Em", Provider: "hacksaw", Type: "slot", Category: "slots", RTP: 96.30, MinBet: 0.20, MaxBet: 100, MobileReady: true},
			{ExternalID: "hack_chaos_001", Name: "Chaos Crew", Provider: "hacksaw", Type: "slot", Category: "slots", RTP: 96.0, MinBet: 0.10, MaxBet: 100, MobileReady: true},
			{ExternalID: "hack_den_001", Name: "Denali", Provider: "hacksaw", Type: "slot", Category: "slots", RTP: 96.0, MinBet: 0.20, MaxBet: 100, MobileReady: true},
		}
	default:
		return []ProviderGame{}
	}
}

// GetAllGames returns games from all active providers
func (s *ProviderIntegrationService) GetAllGames() ([]ProviderGame, error) {
	activeProviders := s.GetActiveProviders()
	var allGames []ProviderGame

	for _, provider := range activeProviders {
		games, err := s.GetProviderGames(provider.ID)
		if err != nil {
			continue
		}
		allGames = append(allGames, games...)
	}

	return allGames, nil
}

// CreateSession creates a game session for a user
func (s *ProviderIntegrationService) CreateSession(userID, providerID, gameID, currency string) (*GameSession, error) {
	provider, err := s.GetProvider(providerID)
	if err != nil {
		return nil, err
	}

	// Get user's balance (in production, fetch from wallet)
	balance := 1000.0 // Demo balance

	session := &GameSession{
		SessionID:  uuid.New().String(),
		UserID:     userID,
		Provider:   providerID,
		GameID:     gameID,
		Currency:   currency,
		Balance:    balance,
		Status:     "active",
		StartedAt:  time.Now(),
		Metadata:   make(map[string]interface{}),
	}

	// In production, this would call provider's API to get game URL
	// For demo, generate a placeholder URL
	session.Metadata["game_url"] = fmt.Sprintf("https://games.%s.com/%s?session=%s", providerID, gameID, session.SessionID)

	return session, nil
}

// RecordTransaction records a game transaction
func (s *ProviderIntegrationService) RecordTransaction(tx *Transaction) error {
	// In production, save to database
	return nil
}

// GenerateSignature generates HMAC signature for provider API
func (s *ProviderIntegrationService) GenerateSignature(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidateSignature validates HMAC signature
func (s *ProviderIntegrationService) ValidateSignature(secret, message, signature string) bool {
	expected := s.GenerateSignature(secret, message)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// SyncGames syncs games from provider to database
func (s *ProviderIntegrationService) SyncGames(providerID string) (int, error) {
	games, err := s.GetProviderGames(providerID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, game := range games {
		// Check if game exists
		var existing models.Game
		result := s.db.Where("provider = ? AND external_id = ?", providerID, game.ExternalID).First(&existing)

		gameModel := models.Game{
			Name:          game.Name,
			Type:          game.Type,
			Provider:      game.Provider,
			RTP:           game.RTP,
			MinBet:        game.MinBet,
			MaxBet:        game.MaxBet,
			ThumbnailURL:   game.ThumbnailURL,
			IsActive:       true,
		}

		if result.Error == gorm.ErrRecordNotFound {
			gameModel.ExternalID = game.ExternalID
			s.db.Create(&gameModel)
			count++
		} else if result.Error == nil {
			// Update existing
			s.db.Model(&existing).Updates(gameModel)
			count++
		}
	}

	return count, nil
}

// GetGamesByCategory returns games filtered by category
func (s *ProviderIntegrationService) GetGamesByCategory(category string) ([]ProviderGame, error) {
	allGames, err := s.GetAllGames()
	if err != nil {
		return nil, err
	}

	var filtered []ProviderGame
	for _, game := range allGames {
		if game.Category == category || contains(game.Tags, category) {
			filtered = append(filtered, game)
		}
	}

	return filtered, nil
}

// GetGamesByType returns games filtered by type
func (s *ProviderIntegrationService) GetGamesByType(gameType string) ([]ProviderGame, error) {
	allGames, err := s.GetAllGames()
	if err != nil {
		return nil, err
	}

	var filtered []ProviderGame
	for _, game := range allGames {
		if game.Type == gameType {
			filtered = append(filtered, game)
		}
	}

	return filtered, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SearchGames searches games by name or tags
func (s *ProviderIntegrationService) SearchGames(query string) ([]ProviderGame, error) {
	allGames, err := s.GetAllGames()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var filtered []ProviderGame
	for _, game := range allGames {
		if strings.Contains(strings.ToLower(game.Name), query) {
			filtered = append(filtered, game)
			continue
		}
		for _, tag := range game.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				filtered = append(filtered, game)
				break
			}
		}
	}

	return filtered, nil
}

// MakeAPIRequest makes an authenticated API request to provider
func (s *ProviderIntegrationService) MakeAPIRequest(providerID, endpoint, method string, body []byte) ([]byte, error) {
	provider, err := s.GetProvider(providerID)
	if err != nil {
		return nil, err
	}

	if !provider.IsActive {
		return nil, fmt.Errorf("provider is not active")
	}

	url := provider.BaseURL + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add authentication headers
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", provider.APIKey)
	req.Header.Set("X-Timestamp", timestamp)

	// Sign request
	signature := s.GenerateSignature(provider.APISecret, timestamp+string(body))
	req.Header.Set("X-Signature", signature)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetProviderStats returns statistics for a provider
func (s *ProviderIntegrationService) GetProviderStats(providerID string) (map[string]interface{}, error) {
	provider, err := s.GetProvider(providerID)
	if err != nil {
		return nil, err
	}

	games, err := s.GetProviderGames(providerID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"provider":     provider.Name,
		"is_active":   provider.IsActive,
		"games_count": len(games),
		"currencies":   len(provider.Currencies),
		"game_types":  len(provider.GameTypes),
	}, nil
}
