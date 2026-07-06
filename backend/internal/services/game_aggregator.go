package services

import (
	"sync"

	"tigercasino/backend/internal/models"
)

// GameAggregator aggregates games from all providers
type GameAggregator struct {
	providers map[string]ProviderGames
	mu        sync.RWMutex
}

// ProviderGames holds games from a specific provider
type ProviderGames struct {
	Provider string
	Games    []models.Game
}

// NewGameAggregator creates a new game aggregator
func NewGameAggregator() *GameAggregator {
	return &GameAggregator{
		providers: make(map[string]ProviderGames),
	}
}

// RegisterProvider registers a provider with its games
func (ga *GameAggregator) RegisterProvider(provider string, games []models.Game) {
	ga.mu.Lock()
	defer ga.mu.Unlock()
	ga.providers[provider] = ProviderGames{
		Provider: provider,
		Games:    games,
	}
}

// GetAllGames returns all games from all providers
func (ga *GameAggregator) GetAllGames() []models.Game {
	ga.mu.RLock()
	defer ga.mu.RUnlock()

	var allGames []models.Game
	for _, pg := range ga.providers {
		allGames = append(allGames, pg.Games...)
	}
	return allGames
}

// GetGamesByProvider returns games from a specific provider
func (ga *GameAggregator) GetGamesByProvider(provider string) []models.Game {
	ga.mu.RLock()
	defer ga.mu.RUnlock()

	if pg, ok := ga.providers[provider]; ok {
		return pg.Games
	}
	return nil
}

// GetGamesByCategory returns games filtered by category
func (ga *GameAggregator) GetGamesByCategory(category string) []models.Game {
	ga.mu.RLock()
	defer ga.mu.RUnlock()

	var games []models.Game
	for _, pg := range ga.providers {
		for _, game := range pg.Games {
			if game.Type == category {
				games = append(games, game)
			}
		}
	}
	return games
}

// SearchGames searches games by name
func (ga *GameAggregator) SearchGames(query string) []models.Game {
	ga.mu.RLock()
	defer ga.mu.RUnlock()

	var games []models.Game
	lowerQuery := toLower(query)
	for _, pg := range ga.providers {
		for _, game := range pg.Games {
			if contains(toLower(game.Name), lowerQuery) {
				games = append(games, game)
			}
		}
	}
	return games
}

// GetProviderStats returns statistics about providers
func (ga *GameAggregator) GetProviderStats() []map[string]interface{} {
	ga.mu.RLock()
	defer ga.mu.RUnlock()

	var stats []map[string]interface{}
	for provider, pg := range ga.providers {
		stats = append(stats, map[string]interface{}{
			"provider":     provider,
			"game_count":   len(pg.Games),
			"categories":   getCategories(pg.Games),
		})
	}
	return stats
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func getCategories(games []models.Game) []string {
	catMap := make(map[string]bool)
	for _, game := range games {
		catMap[game.Type] = true
	}
	categories := make([]string, 0, len(catMap))
	for cat := range catMap {
		categories = append(categories, cat)
	}
	return categories
}

// InitializeDefaultGames initializes the aggregator with all providers
func InitializeDefaultGames() *GameAggregator {
	ga := NewGameAggregator()

	// Evolution Gaming
	ga.RegisterProvider("evolution", []models.Game{
		{ExternalID: "evo_bj_001", Name: "Blackjack", Type: "live_blackjack", Provider: "evolution", RTP: 99.50, MinBet: 10, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_bj_002", Name: "Speed Blackjack", Type: "live_blackjack", Provider: "evolution", RTP: 99.50, MinBet: 10, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_bj_003", Name: "Infinite Blackjack", Type: "live_blackjack", Provider: "evolution", RTP: 99.47, MinBet: 1, MaxBet: 2500, IsActive: true},
		{ExternalID: "evo_r_001", Name: "Live Roulette", Type: "live_roulette", Provider: "evolution", RTP: 97.30, MinBet: 1, MaxBet: 10000, IsActive: true},
		{ExternalID: "evo_r_002", Name: "Lightning Roulette", Type: "live_roulette", Provider: "evolution", RTP: 97.30, MinBet: 1, MaxBet: 5000, IsActive: true},
		{ExternalID: "evo_b_001", Name: "Baccarat", Type: "live_baccarat", Provider: "evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{ExternalID: "evo_gs_001", Name: "Crazy Time", Type: "game_show", Provider: "evolution", RTP: 95.50, MinBet: 0.10, MaxBet: 10000, IsActive: true},
		{ExternalID: "evo_gs_002", Name: "Monopoly Live", Type: "game_show", Provider: "evolution", RTP: 96.23, MinBet: 0.10, MaxBet: 1000, IsActive: true},
		{ExternalID: "evo_gs_003", Name: "Dream Catcher", Type: "game_show", Provider: "evolution", RTP: 96.58, MinBet: 0.10, MaxBet: 1000, IsActive: true},
		{ExternalID: "evo_gs_004", Name: "Lightning Dice", Type: "game_show", Provider: "evolution", RTP: 96.03, MinBet: 0.10, MaxBet: 5000, IsActive: true},
	})

	// Pragmatic Play
	ga.RegisterProvider("pragmatic_play", []models.Game{
		{ExternalID: "pp_001", Name: "Sweet Bonanza", Type: "slots", Provider: "pragmatic_play", RTP: 96.48, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_002", Name: "Wolf Gold", Type: "slots", Provider: "pragmatic_play", RTP: 96.01, MinBet: 0.25, MaxBet: 125, IsActive: true},
		{ExternalID: "pp_003", Name: "The Dog House", Type: "slots", Provider: "pragmatic_play", RTP: 96.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_004", Name: "Gates of Olympus", Type: "slots", Provider: "pragmatic_play", RTP: 95.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pp_005", Name: "Big Bass Bonanza", Type: "slots", Provider: "pragmatic_play", RTP: 96.71, MinBet: 0.10, MaxBet: 125, IsActive: true},
	})

	// NetEnt
	ga.RegisterProvider("netent", []models.Game{
		{ExternalID: "netent_001", Name: "Starburst", Type: "slots", Provider: "netent", RTP: 96.09, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "netent_002", Name: "Gonzo's Quest", Type: "slots", Provider: "netent", RTP: 95.97, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "netent_003", Name: "Dead or Alive", Type: "slots", Provider: "netent", RTP: 96.80, MinBet: 0.09, MaxBet: 18, IsActive: true},
		{ExternalID: "netent_004", Name: "Twin Spin", Type: "slots", Provider: "netent", RTP: 96.60, MinBet: 0.25, MaxBet: 125, IsActive: true},
	})

	// BGaming
	ga.RegisterProvider("bgaming", []models.Game{
		{ExternalID: "bgaming_001", Name: "Avalon", Type: "slots", Provider: "bgaming", RTP: 95.00, MinBet: 0.20, MaxBet: 50, IsActive: true},
		{ExternalID: "bgaming_002", Name: "Lucky Lady's Clover", Type: "slots", Provider: "bgaming", RTP: 96.00, MinBet: 0.20, MaxBet: 50, IsActive: true},
		{ExternalID: "bgaming_003", Name: "Plinko", Type: "crash", Provider: "bgaming", RTP: 98.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "bgaming_004", Name: "Mines", Type: "mines", Provider: "bgaming", RTP: 97.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
	})

	// Spribe
	ga.RegisterProvider("spribe", []models.Game{
		{ExternalID: "spribe_001", Name: "Aviator", Type: "crash", Provider: "spribe", RTP: 97.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_002", Name: "Mines", Type: "mines", Provider: "spribe", RTP: 97.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "spribe_003", Name: "Plinko", Type: "crash", Provider: "spribe", RTP: 98.00, MinBet: 0.10, MaxBet: 100, IsActive: true},
	})

	// Play'n GO
	ga.RegisterProvider("playngo", []models.Game{
		{ExternalID: "pgo_001", Name: "Book of Dead", Type: "slots", Provider: "playngo", RTP: 96.21, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_002", Name: "Reactoonz", Type: "slots", Provider: "playngo", RTP: 96.51, MinBet: 0.20, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_003", Name: "Legacy of Dead", Type: "slots", Provider: "playngo", RTP: 96.53, MinBet: 0.10, MaxBet: 100, IsActive: true},
		{ExternalID: "pgo_004", Name: "Rise of Olympus", Type: "slots", Provider: "playngo", RTP: 96.50, MinBet: 0.20, MaxBet: 100, IsActive: true},
	})

	return ga
}
