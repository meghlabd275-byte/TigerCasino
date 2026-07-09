package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
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
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// GameAggregatorService handles connections to game aggregators (Via, Slotegrator, etc.)
type GameAggregatorService struct {
	db           *gorm.DB
	redis        *redis.Client
	providers    map[string]GameAggregator
	config       *AggregatorConfig
	providerMu   sync.RWMutex
}

type AggregatorConfig struct {
	ViaAPIURL        string
	ViaAPIKey        string
	ViaSecretKey     string
	SlotegratorURL  string
	SlotegratorAPIKey string
	SoftGamingsURL  string
	SoftGamingsKey  string
	GamePlayURL     string
	GamePlayKey     string
	BetConstructURL string
	BetConstructKey string
	Timeout         time.Duration
	RetryCount      int
}

// GameAggregator interface for different aggregator implementations
type GameAggregator interface {
	GetName() string
	GetGames(ctx context.Context, category string, page, limit int) (*AggregatorGameList, error)
	GetGameDetails(ctx context.Context, gameID string) (*AggregatorGame, error)
	LaunchGame(ctx context.Context, userID uuid.UUID, gameID string, mode string) (*GameLaunchInfo, error)
	ProcessTransaction(ctx context.Context, txn *TransactionRequest) (*TransactionResult, error)
	GetBalance(ctx context.Context, userID uuid.UUID, gameID string) (*BalanceInfo, error)
	GetJackpots(ctx context.Context) (map[string]float64, error)
	IsAvailable() bool
}

// AggregatorGameList represents a list of games from aggregator
type AggregatorGameList struct {
	Games      []AggregatorGame `json:"games"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Categories []string         `json:"categories"`
}

// AggregatorGame represents a game from aggregator
type AggregatorGame struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Provider    string                 `json:"provider"`
	Category    string                 `json:"category"`
	SubCategory string                 `json:"sub_category"`
	ImageURL    string                 `json:"image_url"`
	IconURL     string                 `json:"icon_url"`
	RTP         float64                `json:"rtp"`
	Volatility  string                 `json:"volatility"`
	MinBet      float64                `json:"min_bet"`
	MaxBet      float64                `json:"max_bet"`
	MaxWin      float64                `json:"max_win"`
	Features    []string               `json:"features"`
	Tags        []string               `json:"tags"`
	IsMobile    bool                   `json:"is_mobile"`
	IsLive      bool                   `json:"is_live"`
	IsNew       bool                   `json:"is_new"`
	IsPopular   bool                   `json:"is_popular"`
	HasDemo     bool                   `json:"has_demo"`
}

// GameLaunchInfo contains information to launch a game
type GameLaunchInfo struct {
	GameID       string            `json:"game_id"`
	URL          string            `json:"url"`
	Token        string            `json:"token"`
	Mode         string            `json:"mode"`
	Language     string            `json:"language"`
	Currency     string            `json:"currency"`
	Balance      float64           `json:"balance"`
	Expiry       time.Time        `json:"expiry"`
	ExtraData    map[string]string `json:"extra_data"`
}

// TransactionRequest represents a game transaction
type TransactionRequest struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	UserID        uuid.UUID `json:"user_id"`
	GameID        string    `json:"game_id"`
	Type          string    `json:"type"` // bet, win, bonus, refund
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	RoundID       string    `json:"round_id"`
	RefID         string    `json:"ref_id"`
	Timestamp     time.Time `json:"timestamp"`
}

// TransactionResult represents transaction processing result
type TransactionResult struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"` // success, failed, pending
	NewBalance    float64   `json:"new_balance"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
}

// BalanceInfo represents user balance for a game
type BalanceInfo struct {
	UserID      uuid.UUID `json:"user_id"`
	GameID      string    `json:"game_id"`
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	BonusBalance float64  `json:"bonus_balance"`
}

// ViaAggregator implements Via game aggregator
type ViaAggregator struct {
	config    *AggregatorConfig
	httpClient *http.Client
	enabled   bool
}

func NewViaAggregator(config *AggregatorConfig) *ViaAggregator {
	return &ViaAggregator{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		enabled: config.ViaAPIURL != "" && config.ViaAPIKey != "",
	}
}

func (v *ViaAggregator) GetName() string {
	return "Via"
}

func (v *ViaAggregator) IsAvailable() bool {
	return v.enabled
}

func (v *ViaAggregator) buildAuthHeaders() map[string]string {
	timestamp := time.Now().Unix()
	signature := v.generateSignature(fmt.Sprintf("%d%s", timestamp, v.config.ViaSecretKey))
	
	return map[string]string{
		"X-API-Key":    v.config.ViaAPIKey,
		"X-Timestamp":  fmt.Sprintf("%d", timestamp),
		"X-Signature":  signature,
		"Content-Type": "application/json",
	}
}

func (v *ViaAggregator) generateSignature(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func (v *ViaAggregator) GetGames(ctx context.Context, category string, page, limit int) (*AggregatorGameList, error) {
	if !v.enabled {
		return nil, fmt.Errorf("Via aggregator not configured")
	}

	url := fmt.Sprintf("%s/api/v1/games?category=%s&page=%d&limit=%d", 
		v.config.ViaAPIURL, category, page, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range v.buildAuthHeaders() {
		req.Header.Set(k, vv)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Via API error: %d", resp.StatusCode)
	}

	var result struct {
		Success bool              `json:"success"`
		Data    AggregatorGameList `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("Via API returned error")
	}

	return &result.Data, nil
}

func (v *ViaAggregator) GetGameDetails(ctx context.Context, gameID string) (*AggregatorGame, error) {
	if !v.enabled {
		return nil, fmt.Errorf("Via aggregator not configured")
	}

	url := fmt.Sprintf("%s/api/v1/games/%s", v.config.ViaAPIURL, gameID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range v.buildAuthHeaders() {
		req.Header.Set(k, vv)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Via API error: %d", resp.StatusCode)
	}

	var result struct {
		Success bool            `json:"success"`
		Data    AggregatorGame `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (v *ViaAggregator) LaunchGame(ctx context.Context, userID uuid.UUID, gameID string, mode string) (*GameLaunchInfo, error) {
	if !v.enabled {
		return nil, fmt.Errorf("Via aggregator not configured")
	}

	url := fmt.Sprintf("%s/api/v1/session/launch", v.config.ViaAPIURL)

	payload := map[string]interface{}{
		"user_id":    userID.String(),
		"game_id":    gameID,
		"mode":       mode, // real, demo
		"language":   "en",
		"currency":   "USD",
		"timestamp":  time.Now().Unix(),
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	for k, vv := range v.buildAuthHeaders() {
		req.Header.Set(k, vv)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Via launch error: %d", resp.StatusCode)
	}

	var result struct {
		Success bool           `json:"success"`
		Data    GameLaunchInfo `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	gameLaunchInfo := result.Data
	gameLaunchInfo.Expiry = time.Now().Add(10 * time.Minute)

	return &gameLaunchInfo, nil
}

func (v *ViaAggregator) ProcessTransaction(ctx context.Context, txn *TransactionRequest) (*TransactionResult, error) {
	if !v.enabled {
		return nil, fmt.Errorf("Via aggregator not configured")
	}

	url := fmt.Sprintf("%s/api/v1/transaction", v.config.ViaAPIURL)

	body, _ := json.Marshal(txn)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	for k, vv := range v.buildAuthHeaders() {
		req.Header.Set(k, vv)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool              `json:"success"`
		Data    TransactionResult `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (v *ViaAggregator) GetBalance(ctx context.Context, userID uuid.UUID, gameID string) (*BalanceInfo, error) {
	if !v.enabled {
		return nil, fmt.Errorf("Via aggregator not configured")
	}

	url := fmt.Sprintf("%s/api/v1/balance/%s/%s", v.config.ViaAPIURL, userID.String(), gameID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range v.buildAuthHeaders() {
		req.Header.Set(k, vv)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool        `json:"success"`
		Data    BalanceInfo `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (v *ViaAggregator) GetJackpots(ctx context.Context) (map[string]float64, error) {
	if !v.enabled {
		return nil, fmt.Errorf("Via aggregator not configured")
	}

	url := fmt.Sprintf("%s/api/v1/jackpots", v.config.ViaAPIURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range v.buildAuthHeaders() {
		req.Header.Set(k, vv)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool                `json:"success"`
		Data    map[string]float64  `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// SlotegratorAggregator implements Slotegrator aggregator
type SlotegratorAggregator struct {
	config     *AggregatorConfig
	httpClient *http.Client
	enabled    bool
}

func NewSlotegratorAggregator(config *AggregatorConfig) *SlotegratorAggregator {
	return &SlotegratorAggregator{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		enabled: config.SlotegratorURL != "" && config.SlotegratorAPIKey != "",
	}
}

func (s *SlotegratorAggregator) GetName() string {
	return "Slotegrator"
}

func (s *SlotegratorAggregator) IsAvailable() bool {
	return s.enabled
}

func (s *SlotegratorAggregator) buildAuthHeaders() map[string]string {
	timestamp := time.Now().Format(time.RFC3339)
	signature := s.generateHMACSignature(fmt.Sprintf("%s:%s", timestamp, s.config.SlotegratorAPIKey))

	return map[string]string{
		"X-API-Key":    s.config.SlotegratorAPIKey,
		"X-Timestamp":  timestamp,
		"X-Signature":  signature,
		"Content-Type": "application/json",
	}
}

func (s *SlotegratorAggregator) generateHMACSignature(data string) string {
	h := hmac.New(sha512.New, []byte(s.config.SlotegratorAPIKey))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *SlotegratorAggregator) GetGames(ctx context.Context, category string, page, limit int) (*AggregatorGameList, error) {
	if !s.enabled {
		return nil, fmt.Errorf("Slotegrator not configured")
	}

	url := fmt.Sprintf("%s/api/integration/v2/games?type=%s&page=%d&per_page=%d",
		s.config.SlotegratorURL, category, page, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range s.buildAuthHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Games      []AggregatorGame `json:"data"`
			Total      int              `json:"total"`
			Categories []string         `json:"categories"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &AggregatorGameList{
		Games:      result.Data.Games,
		Total:      result.Data.Total,
		Page:       page,
		Limit:      limit,
		Categories: result.Data.Categories,
	}, nil
}

func (s *SlotegratorAggregator) GetGameDetails(ctx context.Context, gameID string) (*AggregatorGame, error) {
	if !s.enabled {
		return nil, fmt.Errorf("Slotegrator not configured")
	}

	url := fmt.Sprintf("%s/api/integration/v2/games/%s", s.config.SlotegratorURL, gameID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range s.buildAuthHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data AggregatorGame `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (s *SlotegratorAggregator) LaunchGame(ctx context.Context, userID uuid.UUID, gameID string, mode string) (*GameLaunchInfo, error) {
	if !s.enabled {
		return nil, fmt.Errorf("Slotegrator not configured")
	}

	url := fmt.Sprintf("%s/api/integration/v2/game/launch", s.config.SlotegratorURL)

	payload := map[string]interface{}{
		"user_id":   userID.String(),
		"game_id":   gameID,
		"mode":      mode,
		"back_url":  fmt.Sprintf("%s/games", getBaseURL()),
		"lang":      "en",
		"currency":  "USD",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	for k, v := range s.buildAuthHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Slotegrator launch error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Data struct {
			GameID    string            `json:"game_id"`
			URL       string            `json:"game_url"`
			Token     string            `json:"session_token"`
			ExpiresAt int64             `json:"expires_at"`
			ExtraData map[string]string `json:"extra_data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &GameLaunchInfo{
		GameID:    result.Data.GameID,
		URL:       result.Data.URL,
		Token:     result.Data.Token,
		Mode:      mode,
		Language:  "en",
		Currency:  "USD",
		Expiry:    time.Unix(result.Data.ExpiresAt, 0),
		ExtraData: result.Data.ExtraData,
	}, nil
}

func (s *SlotegratorAggregator) ProcessTransaction(ctx context.Context, txn *TransactionRequest) (*TransactionResult, error) {
	if !s.enabled {
		return nil, fmt.Errorf("Slotegrator not configured")
	}

	url := fmt.Sprintf("%s/api/integration/v2/transaction", s.config.SlotegratorURL)

	payload := map[string]interface{}{
		"external_transaction_id": txn.TransactionID.String(),
		"user_id":                txn.UserID.String(),
		"game_id":                txn.GameID,
		"type":                   txn.Type,
		"amount":                 txn.Amount,
		"currency":               txn.Currency,
		"round_id":               txn.RoundID,
		"reference_id":           txn.RefID,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	for k, v := range s.buildAuthHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data TransactionResult `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (s *SlotegratorAggregator) GetBalance(ctx context.Context, userID uuid.UUID, gameID string) (*BalanceInfo, error) {
	return &BalanceInfo{
		UserID:   userID,
		GameID:   gameID,
		Balance:  0,
		Currency: "USD",
	}, nil
}

func (s *SlotegratorAggregator) GetJackpots(ctx context.Context) (map[string]float64, error) {
	return make(map[string]float64), nil
}

// NewGameAggregatorService creates a new game aggregator service
func NewGameAggregatorService(db *gorm.DB, redisClient *redis.Client, config *AggregatorConfig) *GameAggregatorService {
	service := &GameAggregatorService{
		db:        db,
		redis:     redisClient,
		config:    config,
		providers: make(map[string]GameAggregator),
	}

	// Initialize available aggregators
	if config != nil {
		via := NewViaAggregator(config)
		if via.IsAvailable() {
			service.providers["via"] = via
		}

		slotegrator := NewSlotegratorAggregator(config)
		if slotegrator.IsAvailable() {
			service.providers["slotegrator"] = slotegrator
		}
	}

	// Always add fallback internal games
	service.providers["internal"] = NewInternalGameProvider()

	return service
}

// GetAvailableProviders returns list of available game providers
func (s *GameAggregatorService) GetAvailableProviders() []string {
	s.providerMu.RLock()
	defer s.providerMu.RUnlock()

	providers := make([]string, 0, len(s.providers))
	for name, provider := range s.providers {
		if provider.IsAvailable() {
			providers = append(providers, name)
		}
	}
	return providers
}

// GetProvider returns a specific provider
func (s *GameAggregatorService) GetProvider(name string) GameAggregator {
	s.providerMu.RLock()
	defer s.providerMu.RUnlock()

	if provider, ok := s.providers[name]; ok && provider.IsAvailable() {
		return provider
	}
	return nil
}

// SyncGamesFromAggregator syncs games from all configured aggregators
func (s *GameAggregatorService) SyncGamesFromAggregator(ctx context.Context) error {
	// This would sync games to local database for caching
	// In production, this would be called periodically via cron

	for name, provider := range s.providers {
		if !provider.IsAvailable() {
			continue
		}

		categories := []string{"slots", "live_casino", "table_games", "instant_games", "virtual_sports"}

		for _, category := range categories {
			page := 1
			limit := 100

			for {
				games, err := provider.GetGames(ctx, category, page, limit)
				if err != nil {
					break
				}

				if len(games.Games) == 0 {
					break
				}

				// Save games to database
				for _, game := range games.Games {
					s.saveGameToDB(&game, name)
				}

				if len(games.Games) < limit {
					break
				}
				page++
			}
		}
	}

	return nil
}

func (s *GameAggregatorService) saveGameToDB(game *AggregatorGame, provider string) {
	var existingGame models.Game
	result := s.db.Where("external_id = ? AND provider = ?", game.ID, provider).First(&existingGame)

	if result.Error == gorm.ErrRecordNotFound {
		newGame := models.Game{
			ID:           uuid.New(),
			ExternalID:   game.ID,
			Name:         game.Name,
			Provider:     game.Provider,
			Category:     game.Category,
			SubCategory:  game.SubCategory,
			ImageURL:     game.ImageURL,
			IconURL:      game.IconURL,
			RTP:          game.RTP,
			Volatility:   game.Volatility,
			MinBet:       game.MinBet,
			MaxBet:       game.MaxBet,
			MaxWin:       game.MaxWin,
			Features:     strings.Join(game.Features, ","),
			Tags:         strings.Join(game.Tags, ","),
			IsMobile:     game.IsMobile,
			IsLive:       game.IsLive,
			IsNew:        game.IsNew,
			IsPopular:    game.IsPopular,
			HasDemo:      game.HasDemo,
			Status:       "active",
			Source:       provider,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		s.db.Create(&newGame)
	} else if result.Error == nil {
		// Update existing game
		existingGame.Name = game.Name
		existingGame.Provider = game.Provider
		existingGame.Category = game.Category
		existingGame.RTP = game.RTP
		existingGame.ImageURL = game.ImageURL
		existingGame.UpdatedAt = time.Now()
		s.db.Save(&existingGame)
	}
}

// GetGames returns games from cache or database
func (s *GameAggregatorService) GetGames(ctx context.Context, category string, provider string, page, limit int) ([]models.Game, int64, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("games:%s:%s:%d:%d", category, provider, page, limit)
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var games []models.Game
			if json.Unmarshal([]byte(cached), &games) == nil {
				return games, int64(len(games)), nil
			}
		}
	}

	// Query database
	var games []models.Game
	query := s.db.Model(&models.Game{}).Where("status = ?", "active")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if provider != "" {
		query = query.Where("provider = ? OR source = ?", provider, provider)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("is_popular DESC, created_at DESC").Find(&games).Error; err != nil {
		return nil, 0, err
	}

	// Cache results
	if s.redis != nil && len(games) > 0 {
		if data, err := json.Marshal(games); err == nil {
			s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}

	return games, total, nil
}

// LaunchGame launches a game for a user
func (s *GameAggregatorService) LaunchGame(ctx context.Context, userID uuid.UUID, gameID string, mode string) (*GameLaunchInfo, error) {
	// First check if it's an internal game
	if internalProvider := s.GetProvider("internal"); internalProvider != nil {
		if launchInfo, err := internalProvider.LaunchGame(ctx, userID, gameID, mode); err == nil {
			return launchInfo, nil
		}
	}

	// Get game info from database to determine provider
	var game models.Game
	if err := s.db.Where("external_id = ? OR id = ?", gameID, gameID).First(&game).Error; err != nil {
		return nil, fmt.Errorf("game not found: %v", err)
	}

	// Try to launch with the game's source provider
	provider := s.GetProvider(game.Source)
	if provider == nil {
		// Fallback to any available provider
		providers := s.GetAvailableProviders()
		if len(providers) == 0 {
			return nil, fmt.Errorf("no game providers available")
		}
		provider = s.GetProvider(providers[0])
	}

	if provider == nil {
		return nil, fmt.Errorf("cannot launch game")
	}

	return provider.LaunchGame(ctx, userID, gameID, mode)
}

// InternalGameProvider provides internally developed games
type InternalGameProvider struct {
	games map[string]*InternalGame
}

type InternalGame struct {
	ID          string
	Name        string
	Category    string
	Provider    string
	RTP         float64
	MinBet      float64
	MaxBet      float64
	MaxWin      float64
	ImageURL    string
	IsMobile    bool
}

func NewInternalGameProvider() *InternalGameProvider {
	return &InternalGameProvider{
		games: getInternalGames(),
	}
}

func (i *InternalGameProvider) GetName() string {
	return "TigerCasino Originals"
}

func (i *InternalGameProvider) IsAvailable() bool {
	return true
}

func (i *InternalGameProvider) GetGames(ctx context.Context, category string, page, limit int) (*AggregatorGameList, error) {
	var games []AggregatorGame
	for _, g := range i.games {
		if category == "" || g.Category == category {
			games = append(games, AggregatorGame{
				ID:         g.ID,
				Name:       g.Name,
				Provider:  g.Provider,
				Category:   g.Category,
				ImageURL:   g.ImageURL,
				RTP:        g.RTP,
				MinBet:     g.MinBet,
				MaxBet:     g.MaxBet,
				MaxWin:     g.MaxWin,
				IsMobile:   g.IsMobile,
				IsPopular:  true,
				HasDemo:    true,
			})
		}
	}

	start := (page - 1) * limit
	end := start + limit
	if start >= len(games) {
		games = []AggregatorGame{}
	} else {
		if end > len(games) {
			end = len(games)
		}
		games = games[start:end]
	}

	return &AggregatorGameList{
		Games:      games,
		Total:      len(i.games),
		Page:       page,
		Limit:      limit,
		Categories: []string{"slots", "crash", "table_games", "instant_games"},
	}, nil
}

func (i *InternalGameProvider) GetGameDetails(ctx context.Context, gameID string) (*AggregatorGame, error) {
	game, ok := i.games[gameID]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}

	return &AggregatorGame{
		ID:         game.ID,
		Name:       game.Name,
		Provider:   game.Provider,
		Category:   game.Category,
		ImageURL:   game.ImageURL,
		RTP:        game.RTP,
		MinBet:     game.MinBet,
		MaxBet:     game.MaxBet,
		MaxWin:     game.MaxWin,
		IsMobile:   game.IsMobile,
		HasDemo:    true,
	}, nil
}

func (i *InternalGameProvider) LaunchGame(ctx context.Context, userID uuid.UUID, gameID string, mode string) (*GameLaunchInfo, error) {
	game, ok := i.games[gameID]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}

	token := fmt.Sprintf("%s_%d", userID.String(), time.Now().Unix())

	return &GameLaunchInfo{
		GameID:   gameID,
		URL:      fmt.Sprintf("/games/%s?token=%s&mode=%s", gameID, token, mode),
		Token:    token,
		Mode:     mode,
		Language: "en",
		Currency: "USD",
		Expiry:   time.Now().Add(30 * time.Minute),
	}, nil
}

func (i *InternalGameProvider) ProcessTransaction(ctx context.Context, txn *TransactionRequest) (*TransactionResult, error) {
	return &TransactionResult{
		TransactionID: txn.TransactionID,
		Status:        "success",
		NewBalance:    txn.Amount,
		Message:       "Transaction processed",
		Timestamp:    time.Now(),
	}, nil
}

func (i *InternalGameProvider) GetBalance(ctx context.Context, userID uuid.UUID, gameID string) (*BalanceInfo, error) {
	return &BalanceInfo{
		UserID:   userID,
		GameID:   gameID,
		Balance:  0,
		Currency: "USD",
	}, nil
}

func (i *InternalGameProvider) GetJackpots(ctx context.Context) (map[string]float64, error) {
	return map[string]float64{
		"mini":    1000,
		"major":   10000,
		"grand":   100000,
	}, nil
}

func getInternalGames() map[string]*InternalGame {
	return map[string]*InternalGame{
		"tiger-crash": {
			ID:        "tiger-crash",
			Name:      "Tiger Crash",
			Category:  "crash",
			Provider:  "TigerCasino",
			RTP:       97.5,
			MinBet:    0.1,
			MaxBet:    1000,
			MaxWin:    100000,
			ImageURL:  "/images/games/tiger-crash.png",
			IsMobile:  true,
		},
		"tiger-mines": {
			ID:        "tiger-mines",
			Name:      "Tiger Mines",
			Category:  "instant_games",
			Provider:  "TigerCasino",
			RTP:       97.0,
			MinBet:    0.1,
			MaxBet:    500,
			MaxWin:    50000,
			ImageURL:  "/images/games/tiger-mines.png",
			IsMobile:  true,
		},
		"tiger-plinko": {
			ID:        "tiger-plinko",
			Name:      "Tiger Plinko",
			Category:  "instant_games",
			Provider:  "TigerCasino",
			RTP:       98.0,
			MinBet:    0.1,
			MaxBet:    500,
			MaxWin:    100000,
			ImageURL:  "/images/games/tiger-plinko.png",
			IsMobile:  true,
		},
		"tiger-dice": {
			ID:        "tiger-dice",
			Name:      "Tiger Dice",
			Category:  "instant_games",
			Provider:  "TigerCasino",
			RTP:       99.0,
			MinBet:    0.01,
			MaxBet:    5000,
			MaxWin:    50000,
			ImageURL:  "/images/games/tiger-dice.png",
			IsMobile:  true,
		},
		"tiger-limbo": {
			ID:        "tiger-limbo",
			Name:      "Tiger Limbo",
			Category:  "instant_games",
			Provider:  "TigerCasino",
			RTP:       96.5,
			MinBet:    0.1,
			MaxBet:    1000,
			MaxWin:    100000,
			ImageURL:  "/images/games/tiger-limbo.png",
			IsMobile:  true,
		},
		"tiger-keno": {
			ID:        "tiger-keno",
			Name:      "Tiger Keno",
			Category:  "lottery",
			Provider:  "TigerCasino",
			RTP:       95.0,
			MinBet:    0.1,
			MaxBet:    100,
			MaxWin:    10000,
			ImageURL:  "/images/games/tiger-keno.png",
			IsMobile:  true,
		},
		"tiger-hilo": {
			ID:        "tiger-hilo",
			Name:      "Tiger Hi-Lo",
			Category:  "table_games",
			Provider:  "TigerCasino",
			RTP:       98.5,
			MinBet:    1,
			MaxBet:    10000,
			MaxWin:    20000,
			ImageURL:  "/images/games/tiger-hilo.png",
			IsMobile:  true,
		},
		"tiger-video-poker": {
			ID:        "tiger-video-poker",
			Name:      "Tiger Video Poker",
			Category:  "table_games",
			Provider:  "TigerCasino",
			RTP:       99.5,
			MinBet:    0.1,
			MaxBet:    500,
			MaxWin:    5000,
			ImageURL:  "/images/games/tiger-video-poker.png",
			IsMobile:  true,
		},
		"tiger-scratch": {
			ID:        "tiger-scratch",
			Name:      "Tiger Scratch",
			Category:  "instant_games",
			Provider:  "TigerCasino",
			RTP:       94.5,
			MinBet:    0.5,
			MaxBet:    100,
			MaxWin:    10000,
			ImageURL:  "/images/games/tiger-scratch.png",
			IsMobile:  true,
		},
		"tiger-roulette": {
			ID:        "tiger-roulette",
			Name:      "Tiger Roulette",
			Category:  "table_games",
			Provider:  "TigerCasino",
			RTP:       97.3,
			MinBet:    1,
			MaxBet:    10000,
			MaxWin:    35000,
			ImageURL:  "/images/games/tiger-roulette.png",
			IsMobile:  true,
		},
		"tiger-blackjack": {
			ID:        "tiger-blackjack",
			Name:      "Tiger Blackjack",
			Category:  "table_games",
			Provider:  "TigerCasino",
			RTP:       99.5,
			MinBet:    5,
			MaxBet:    10000,
			MaxWin:    15000,
			ImageURL:  "/images/games/tiger-blackjack.png",
			IsMobile:  true,
		},
		"tiger-baccarat": {
			ID:        "tiger-baccarat",
			Name:      "Tiger Baccarat",
			Category:  "table_games",
			Provider:  "TigerCasino",
			RTP:       98.9,
			MinBet:    5,
			MaxBet:    10000,
			MaxWin:    15000,
			ImageURL:  "/images/games/tiger-baccarat.png",
			IsMobile:  true,
		},
	}
}

func getBaseURL() string {
	return "https://tigercasino.com"
}
