package services

import (
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

// ProviderService handles game provider integrations
type ProviderService struct {
	db           *gorm.DB
	httpClient   *http.Client
	providers    map[string]ProviderClient
	providerLock sync.RWMutex
}

// ProviderClient interface for game providers
type ProviderClient interface {
	GetName() string
	GetGames() ([]models.Game, error)
	LaunchGame(userID, gameID, currency string) (string, error)
}

// NewProviderService creates a new provider service
func NewProviderService(db *gorm.DB) *ProviderService {
	s := &ProviderService{
		db:         db,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		providers:  make(map[string]ProviderClient),
	}
	
	// Initialize default providers
	s.initializeProviders()
	
	return s
}

func (s *ProviderService) initializeProviders() {
	// These would be initialized with actual API credentials in production
	s.providers["pragmatic_play"] = &PragmaticPlayClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.pragmaticplay.com",
	}
	
	s.providers["evolution"] = &EvolutionClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.evolutiongaming.com",
	}
	
	s.providers["netent"] = &NetEntClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.netent.com",
	}
	
	s.providers["bgaming"] = &BGamingClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.bgaming.com",
	}
	
	s.providers["spribe"] = &SpribeClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.spribe.co",
	}
	
	s.providers["playngo"] = &PlayNGOClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.playngo.com",
	}
	
	s.providers["hacksaw"] = &HacksawClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.hacksawgaming.com",
	}
	
	s.providers["nolimit"] = &NolimitClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.nolimitcity.com",
	}
	
	s.providers["relax"] = &RelaxClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.relaxgaming.com",
	}
	
	s.providers["push"] = &PushClient{
		httpClient: s.httpClient,
		apiKey:     "",
		baseURL:    "https://api.pushgaming.com",
	}
}

// GetProviders returns all available providers
func (s *ProviderService) GetProviders() []map[string]interface{} {
	s.providerLock.RLock()
	defer s.providerLock.RUnlock()
	
	var result []map[string]interface{}
	for name, provider := range s.providers {
		result = append(result, map[string]interface{}{
			"name":     provider.GetName(),
			"id":       name,
			"games":    0, // Would fetch from provider
			"status":   "active",
		})
	}
	
	return result
}

// GetProviderGames returns games from a specific provider
func (s *ProviderService) GetProviderGames(providerID string) ([]models.Game, error) {
	s.providerLock.RLock()
	provider, ok := s.providers[providerID]
	s.providerLock.RUnlock()
	
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", providerID)
	}
	
	return provider.GetGames()
}

// GetAllGames returns all games from all providers
func (s *ProviderService) GetAllGames() ([]models.Game, error) {
	var allGames []models.Game
	
	s.providerLock.RLock()
	defer s.providerLock.RUnlock()
	
	for _, provider := range s.providers {
		games, err := provider.GetGames()
		if err != nil {
			continue // Skip providers that fail
		}
		allGames = append(allGames, games...)
	}
	
	return allGames, nil
}

// LaunchGame launches a game from a specific provider
func (s *ProviderService) LaunchGame(providerID, userID, gameID, currency string) (string, error) {
	s.providerLock.RLock()
	provider, ok := s.providers[providerID]
	s.providerLock.RUnlock()
	
	if !ok {
		return "", fmt.Errorf("provider not found: %s", providerID)
	}
	
	return provider.LaunchGame(userID, gameID, currency)
}

// SyncProviderGames syncs games from a provider to the database
func (s *ProviderService) SyncProviderGames(providerID string) (int, error) {
	games, err := s.GetProviderGames(providerID)
	if err != nil {
		return 0, err
	}
	
	count := 0
	for _, game := range games {
		var existing models.Game
		result := s.db.Where("provider = ? AND external_id = ?", providerID, game.ExternalID).First(&existing)
		
		if result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&game).Error; err == nil {
				count++
			}
		} else if result.Error == nil {
			// Update existing game
			if err := s.db.Model(&existing).Updates(game).Error; err == nil {
				count++
			}
		}
	}
	
	return count, nil
}

// Provider Clients

type PragmaticPlayClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *PragmaticPlayClient) GetName() string { return "Pragmatic Play" }

func (c *PragmaticPlayClient) GetGames() ([]models.Game, error) {
	// Simulated game list
	games := []models.Game{
		{ExternalID: "pp_001", Name: "Sweet Bonanza", Type: "slots", Provider: "pragmatic_play", RTP: 96.48, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_002", Name: "Wolf Gold", Type: "slots", Provider: "pragmatic_play", RTP: 96.01, MinBet: 0.25, MaxBet: 125, IsActive: true},
		{ExternalID: "pp_003", Name: "The Dog House", Type: "slots", Provider: "pragmatic_play", RTP: 96.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_004", Name: "Gates of Olympus", Type: "slots", Provider: "pragmatic_play", RTP: 95.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_005", Name: "Big Bass Bonanza", Type: "slots", Provider: "pragmatic_play", RTP: 96.71, MinBet: 0.10, MaxBet: 125, IsActive: true},
		{ExternalID: "pp_006", Name: "Fruit Party", Type: "slots", Provider: "pragmatic_play", RTP: 96.47, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_007", Name: "Starlight Princess", Type: "slots", Provider: "pragmatic_play", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_008", Name: "Power of Thor", Type: "slots", Provider: "pragmatic_play", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_009", Name: "Joker's Jewels", Type: "slots", Provider: "pragmatic_play", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_010", Name: "Wild West Gold", Type: "slots", Provider: "pragmatic_play", RTP: 96.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
	}
	return games, nil
}

func (c *PragmaticPlayClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.pragmaticplay.com/%s?user=%s&currency=%s", gameID, userID, currency), nil
}

type EvolutionClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *EvolutionClient) GetName() string { return "Evolution Gaming" }

func (c *EvolutionClient) GetGames() ([]models.Game, error) {
	games := []models.Game{
		// Live Blackjack
		{ExternalID: "evo_bj_001", Name: "Blackjack", Type: "live_blackjack", Provider: "evolution", RTP: 99.50, MinBet: 10, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_bj_002", Name: "Speed Blackjack", Type: "live_blackjack", Provider: "evolution", RTP: 99.50, MinBet: 10, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_bj_003", Name: "Infinite Blackjack", Type: "live_blackjack", Provider: "evolution", RTP: 99.47, MinBet: 1, MaxBet: 2500, IsActive: true},
		// Live Roulette
		{ExternalID: "evo_r_001", Name: "Live Roulette", Type: "live_roulette", Provider: "evolution", RTP: 97.30, MinBet: 1, MaxBet: 10000, IsActive: true},
		{ExternalID: "evo_r_002", Name: "Speed Roulette", Type: "live_roulette", Provider: "evolution", RTP: 97.30, MinBet: 1, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_r_003", Name: "Lightning Roulette", Type: "live_roulette", Provider: "evolution", RTP: 97.30, MinBet: 1, MaxBet: 5000, IsActive: true},
		// Live Baccarat
		{ExternalID: "evo_b_001", Name: "Baccarat", Type: "live_baccarat", Provider: "evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{ExternalID: "evo_b_002", Name: "Speed Baccarat", Type: "live_baccarat", Provider: "evolution", RTP: 98.94, MinBet: 5, MaxBet: 5000, IsActive: true},
		// Game Shows
		{ExternalID: "evo_gs_001", Name: "Crazy Time", Type: "game_show", Provider: "evolution", RTP: 95.50, MinBet: 0.10, MaxBet: 10000, IsActive: true},
		{ExternalID: "evo_gs_002", Name: "Monopoly Live", Type: "game_show", Provider: "evolution", RTP: 96.23, MinBet: 0.10, MaxBet: 1000, IsActive: true},
		{ExternalID: "evo_gs_003", Name: "Dream Catcher", Type: "game_show", Provider: "evolution", RTP: 96.58, MinBet: 0.10, MaxBet: 1000, IsActive: true},
		{ExternalID: "evo_gs_004", Name: "Lightning Dice", Type: "game_show", Provider: "evolution", RTP: 96.03, MinBet: 0.10, MaxBet: 5000, IsActive: true},
		// Poker
		{ExternalID: "evo_p_001", Name: "Texas Hold'em", Type: "live_poker", Provider: "evolution", RTP: 97.80, MinBet: 5, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_p_002", Name: "Three Card Poker", Type: "live_poker", Provider: "evolution", RTP: 96.63, MinBet: 10, MaxBet: 1000, IsActive: true},
		{ExternalID: "evo_p_003", Name: "Caribbean Stud", Type: "live_poker", Provider: "evolution", RTP: 96.30, MinBet: 10, MaxBet: 500, IsActive: true},
	}
	return games, nil
}

func (c *EvolutionClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.evolutiongaming.com/%s?user=%s&currency=%s", gameID, userID, currency), nil
}

type NetEntClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *NetEntClient) GetName() string { return "NetEnt" }

func (c *NetEntClient) GetGames() ([]models.Game, error) {
	games := []models.Game{
		{ExternalID: "netent_001", Name: "Starburst", Type: "slots", Provider: "netent", RTP: 96.09, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "netent_002", Name: "Gonzo's Quest", Type: "slots", Provider: "netent", RTP: 95.97, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "netent_003", Name: "Dead or Alive", Type: "slots", Provider: "netent", RTP: 96.80, MinBet: 0.09, MaxBet: 18, IsActive: true},
		{ExternalID: "netent_004", Name: "Twin Spin", Type: "slots", Provider: "netent", RTP: 96.60, MinBet: 0.25, MaxBet: 125, IsActive: true},
		{ExternalID: "netent_005", Name: "Mega Fortune", Type: "slots", Provider: "netent", RTP: 96.60, MinBet: 0.25, MaxBet: 62.50, IsActive: true},
		{ExternalID: "netent_006", Name: "Hall of Gods", Type: "slots", Provider: "netent", RTP: 95.70, MinBet: 0.20, MaxBet: 50, IsActive: true},
		{ExternalID: "netent_007", Name: "Blood Suckers", Type: "slots", Provider: "netent", RTP: 98.00, MinBet: 0.01, MaxBet: 50, IsActive: true},
		{ExternalID: "netent_008", Name: "Jack and the Beanstalk", Type: "slots", Provider: "netent", RTP: 96.30, MinBet: 0.20, MaxBet: 100, IsActive: true},
	}
	return games, nil
}

func (c *NetEntClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.netent.com/%s?user=%s&currency=%s", gameID, userID, currency), nil
}

type BGamingClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *BGamingClient) GetName() string { return "BGaming" }

func (c *BGamingClient) GetGames() ([]models.Game, error) {
	games := []models.Game{
		{ExternalID: "bgaming_001", Name: "Avalon", Type: "slots", Provider: "bgaming", RTP: 95.00, MinBet: 0.20, MaxBet: 50, IsActive: true},
		{ExternalID: "bgaming_002", Name: "Lucky Lady's Clover", Type: "slots", Provider: "bgaming", RTP: 96.00, MinBet: 0.20, MaxBet: 50, IsActive: true},
		{ExternalID: "bgaming_003", Name: "Plinko", Type: "crash", Provider: "bgaming", RTP: 98.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "bgaming_004", Name: "Mines", Type: "mines", Provider: "bgaming", RTP: 97.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "bgaming_005", Name: "Dice", Type: "dice", Provider: "bgaming", RTP: 99.00, MinBet: 0.01, MaxBet: 1000, IsActive: true},
		{ExternalID: "bgaming_006", Name: "Keno", Type: "lottery", Provider: "bgaming", RTP: 95.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
	}
	return games, nil
}

func (c *BGamingClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.bgaming.com/%s?user=%s&currency=%s", gameID, userID, currency), nil
}

type SpribeClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *SpribeClient) GetName() string { return "Spribe" }

func (c *SpribeClient) GetGames() ([]models.Game, error) {
	games := []models.Game{
		{ExternalID: "spribe_001", Name: "Aviator", Type: "crash", Provider: "spribe", RTP: 97.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_002", Name: "Mines", Type: "mines", Provider: "spribe", RTP: 97.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_003", Name: "Plinko", Type: "crash", Provider: "spribe", RTP: 98.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_004", Name: "Dice", Type: "dice", Provider: "spribe", RTP: 99.00, MinBet: 0.01, MaxBet: 1000, IsActive: true},
		{ExternalID: "spribe_005", Name: "Hilo", Type: "hilo", Provider: "spribe", RTP: 96.50, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_006", Name: "Keno", Type: "lottery", Provider: "spribe", RTP: 95.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_007", Name: "Goal", Type: "arcade", Provider: "spribe", RTP: 95.50, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_008", Name: "Mini Games", Type: "arcade", Provider: "spribe", RTP: 96.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
	}
	return games, nil
}

func (c *SpribeClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.spribe.co/%s?user=%s&currency=%s", gameID, userID, currency), nil
}

type PlayNGOClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *PlayNGOClient) GetName() string { return "Play'n GO" }

func (c *PlayNGOClient) GetGames() ([]models.Game, error) {
	games := []models.Game{
		{ExternalID: "pgo_001", Name: "Book of Dead", Type: "slots", Provider: "playngo", RTP: 96.21, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_002", Name: "Reactoonz", Type: "slots", Provider: "playngo", RTP: 96.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_003", Name: "Legacy of Dead", Type: "slots", Provider: "playngo", RTP: 96.53, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_004", Name: "Rise of Olympus", Type: "slots", Provider: "playngo", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_005", Name: "Moon Princess", Type: "slots", Provider: "playngo", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_006", Name: "Fire Joker", Type: "slots", Provider: "playngo", RTP: 96.15, MinBet: 0.05, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_007", Name: "Rich Wilde and the Aztec Idols", Type: "slots", Provider: "playngo", RTP: 96.60, MinBet: 0.10, MaxBet: 50, IsActive: true},
		{ExternalID: "pgo_008", Name: "Gems Pay 300", Type: "slots", Provider: "playngo", RTP: 96.30, MinBet: 0.20, MaxBet: 100, IsActive: true},
	}
	return games, nil
}

func (c *PlayNGOClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.playngo.com/%s?user=%s&currency=%s", gameID, userID, currency), nil
}

// Simplified clients for other providers
type HacksawClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *HacksawClient) GetName() string { return "Hacksaw Gaming" }
func (c *HacksawClient) GetGames() ([]models.Game, error) {
	return []models.Game{
		{ExternalID: "hack_001", Name: "Stick 'Em", Type: "slots", Provider: "hacksaw", RTP: 96.30, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "hack_002", Name: "Chaos Crew", Type: "slots", Provider: "hacksaw", RTP: 96.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
	}, nil
}
func (c *HacksawClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.hacksawgaming.com/%s?user=%s", gameID, userID), nil
}

type NolimitClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *NolimitClient) GetName() string { return "Nolimit City" }
func (c *NolimitClient) GetGames() ([]models.Game, error) {
	return []models.Game{
		{ExternalID: "nlc_001", Name: "San Quentin", Type: "slots", Provider: "nolimit", RTP: 96.00, MinBet: 0.20, MaxBet: 70, IsActive: true},
		{ExternalID: "nlc_002", Name: "Book of Shadows", Type: "slots", Provider: "nolimit", RTP: 96.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "nlc_003", Name: "Mental", Type: "slots", Provider: "nolimit", RTP: 96.00, MinBet: 0.10, MaxBet: 40, IsActive: true},
	}, nil
}
func (c *NolimitClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.nolimitcity.com/%s?user=%s", gameID, userID), nil
}

type RelaxClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *RelaxClient) GetName() string { return "Relax Gaming" }
func (c *RelaxClient) GetGames() ([]models.Game, error) {
	return []models.Game{
		{ExternalID: "relax_001", Name: "Money Train", Type: "slots", Provider: "relax", RTP: 96.15, MinBet: 0.10, MaxBet: 20, IsActive: true},
		{ExternalID: "relax_002", Name: "Temple Tumble", Type: "slots", Provider: "relax", RTP: 96.25, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "relax_003", Name: "Money Train 2", Type: "slots", Provider: "relax", RTP: 96.20, MinBet: 0.10, MaxBet: 20, IsActive: true},
	}, nil
}
func (c *RelaxClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.relaxgaming.com/%s?user=%s", gameID, userID), nil
}

type PushClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func (c *PushClient) GetName() string { return "Push Gaming" }
func (c *PushClient) GetGames() ([]models.Game, error) {
	return []models.Game{
		{ExternalID: "push_001", Name: "Jammin' Jars", Type: "slots", Provider: "push", RTP: 96.83, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "push_002", Name: "Wild Swarm", Type: "slots", Provider: "push", RTP: 97.00, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "push_003", Name: "Big Bamboo", Type: "slots", Provider: "push", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
	}, nil
}
func (c *PushClient) LaunchGame(userID, gameID, currency string) (string, error) {
	return fmt.Sprintf("https://game.pushgaming.com/%s?user=%s", gameID, userID), nil
}

// Add more provider clients as needed...
