package services

import (
	"fmt"
)

// ============ Extended Game Providers ============

type GameProviderService struct {
	providers map[string]*GameProvider
}

type GameProvider struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Logo        string   `json:"logo"`
	GameCount   int      `json:"game_count"`
	IsLive      bool     `json:"is_live"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
	APIEndpoint string   `json:"api_endpoint"`
	RTP         float64  `json:"rtp"`
}

type SlotGame struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	Category    string   `json:"category"`
	Reels       int      `json:"reels"`
	Rows        int      `json:"rows"`
	Paylines    int      `json:"paylines"`
	RTP         float64  `json:"rtp"`
	Volatility  string   `json:"volatility"`
	MinBet      float64  `json:"min_bet"`
	MaxBet      float64  `json:"max_bet"`
	MaxWin      float64  `json:"max_win"`
	Features    []string `json:"features"`
	Thumbnail   string   `json:"thumbnail"`
	IsNew       bool     `json:"is_new"`
	IsHot       bool     `json:"is_hot"`
}

func NewGameProviderService() *GameProviderService {
	svc := &GameProviderService{
		providers: make(map[string]*GameProvider),
	}
	svc.initializeProviders()
	return svc
}

func (s *GameProviderService) initializeProviders() {
	providers := []*GameProvider{
		// Tier 1 - Top Providers
		{
			ID: "pragmatic_play", Name: "Pragmatic Play", Logo: "/providers/pragmatic.png",
			GameCount: 300, IsLive: true, Categories: []string{"slots", "live_casino", "table_games", " bingo"},
			Description: "Leading game provider with 300+ games including popular slots and live casino",
			APIEndpoint: "https://api.pragmaticplay.com", RTP: 96.5,
		},
		{
			ID: "evolution", Name: "Evolution Gaming", Logo: "/providers/evolution.png",
			GameCount: 150, IsLive: true, Categories: []string{"live_casino", "game_shows"},
			Description: "Premier live casino provider with innovative game shows",
			APIEndpoint: "https://api.evolutiongaming.com", RTP: 97.0,
		},
		{
			ID: "microgaming", Name: "Microgaming", Logo: "/providers/microgaming.png",
			GameCount: 800, IsLive: true, Categories: []string{"slots", "table_games", "poker"},
			Description: "Industry pioneer with largest jackpot network",
			APIEndpoint: "https://api.microgaming.com", RTP: 96.0,
		},
		// Tier 2 - Major Providers
		{
			ID: "netent", Name: "NetEnt", Logo: "/providers/netent.png",
			GameCount: 200, IsLive: true, Categories: []string{"slots", "table_games", "jackpots"},
			Description: "Award-winning slots and progressive jackpots",
			APIEndpoint: "https://api.netent.com", RTP: 96.5,
		},
		{
			ID: "playngo", Name: "Play'n GO", Logo: "/providers/playngo.png",
			GameCount: 300, IsLive: true, Categories: []string{"slots", "table_games"},
			Description: "Mobile-first gaming entertainment",
			APIEndpoint: "https://api.playngo.com", RTP: 96.5,
		},
		{
			ID: "yggdrasil", Name: "Yggdrasil", Logo: "/providers/yggdrasil.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "jackpots"},
			Description: "Innovative slots with unique game mechanics",
			APIEndpoint: "https://api.yggdrasil.com", RTP: 96.5,
		},
		{
			ID: "quickspin", Name: "Quickspin", Logo: "/providers/quickspin.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots"},
			Description: "Story-driven slot experiences",
			APIEndpoint: "https://api.quickspin.com", RTP: 97.0,
		},
		// Tier 3 - Growing Providers
		{
			ID: "bgaming", Name: "BGaming", Logo: "/providers/bgaming.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "table_games", "bitcoin"},
			Description: "Crypto-friendly gaming provider",
			APIEndpoint: "https://api.bgaming.com", RTP: 96.5,
		},
		{
			ID: "spribe", Name: "Spribe", Logo: "/providers/spribe.png",
			GameCount: 20, IsLive: true, Categories: []string{"crash", "instant_win", "dice"},
			Description: "Innovative crash games and instant wins",
			APIEndpoint: "https://api.spribe.com", RTP: 97.0,
		},
		{
			ID: "hacksaw", Name: "Hacksaw Gaming", Logo: "/providers/hacksaw.png",
			GameCount: 50, IsLive: true, Categories: []string{"slots", "scratch", "instant_win"},
			Description: "高频 slots and scratch cards",
			APIEndpoint: "https://api.hacksawgaming.com", RTP: 96.5,
		},
		// Tier 4 - Specialty Providers
		{
			ID: "red_tiger", Name: "Red Tiger", Logo: "/providers/redtiger.png",
			GameCount: 200, IsLive: true, Categories: []string{"slots", "jackpots", "daily_drops"},
			Description: "Daily jackpot network and innovative slots",
			APIEndpoint: "https://api.redtigergaming.com", RTP: 96.0,
		},
		{
			ID: "big_time_gaming", Name: "Big Time Gaming", Logo: "/providers/btg.png",
			GameCount: 50, IsLive: true, Categories: []string{"slots", "megaways"},
			Description: "Inventors of Megaways and Megaclusters",
			APIEndpoint: "https://api.bigtimergaming.com", RTP: 97.0,
		},
		{
			ID: "nolimit_city", Name: "Nolimit City", Logo: "/providers/nolimit.png",
			GameCount: 50, IsLive: true, Categories: []string{"slots"},
			Description: "Bold slots with unique features",
			APIEndpoint: "https://api.nolimitcity.com", RTP: 96.5,
		},
		{
			ID: "relax_gaming", Name: "Relax Gaming", Logo: "/providers/relax.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "table_games"},
			Description: "Award-winning aggregation platform",
			APIEndpoint: "https://api.relaxgaming.com", RTP: 96.5,
		},
		{
			ID: "push_gaming", Name: "Push Gaming", Logo: "/providers/push.png",
			GameCount: 30, IsLive: true, Categories: []string{"slots"},
			Description: "Mobile-first HTML5 slots",
			APIEndpoint: "https://api.pushgaming.com", RTP: 97.0,
		},
		// Asian Providers
		{
			ID: "pg_soft", Name: "PG Soft", Logo: "/providers/pgsoft.png",
			GameCount: 200, IsLive: true, Categories: []string{"slots", "table_games"},
			Description: "Asian market specialist with mobile games",
			APIEndpoint: "https://api.pgsoft.com", RTP: 96.5,
		},
		{
			ID: "ka_gaming", Name: "KA Gaming", Logo: "/providers/kagaming.png",
			GameCount: 200, IsLive: true, Categories: []string{"slots", "fishing", "arcade"},
			Description: "Multi-category gaming provider",
			APIEndpoint: "https://api.kagaming.com", RTP: 96.0,
		},
		{
			ID: "spade_gaming", Name: "Spade Gaming", Logo: "/providers/spadegaming.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots", "table_games", "arcade"},
			Description: "Asian-themed premium games",
			APIEndpoint: "https://api.spadegaming.com", RTP: 96.5,
		},
		// Additional Providers
		{
			ID: "betsoft", Name: "Betsoft", Logo: "/providers/betsoft.png",
			GameCount: 150, IsLive: true, Categories: []string{"slots", "table_games", "video_poker"},
			Description: "3D slots and classic casino games",
			APIEndpoint: "https://api.betsoft.com", RTP: 96.0,
		},
		{
			ID: "vivo_gaming", Name: "Vivo Gaming", Logo: "/providers/vivo.png",
			GameCount: 50, IsLive: true, Categories: []string{"live_casino", "baccarat", "blackjack"},
			Description: "Live dealer casino specialist",
			APIEndpoint: "https://api.vivogaming.com", RTP: 97.0,
		},
		{
			ID: "ezugi", Name: "Ezugi", Logo: "/providers/ezugi.png",
			GameCount: 50, IsLive: true, Categories: []string{"live_casino", "lottery"},
			Description: "Live casino with unique variants",
			APIEndpoint: "https://api.ezugi.com", RTP: 97.0,
		},
		{
			ID: " Authentic Gaming", Name: "Authentic Gaming", Logo: "/providers/authentic.png",
			GameCount: 20, IsLive: true, Categories: []string{"live_casino", "roulette"},
			Description: "Premium live roulette from real casinos",
			APIEndpoint: "https://api.authenticgaming.com", RTP: 97.5,
		},
		{
			ID: "endorphina", Name: "Endorphina", Logo: "/providers/endorphina.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "bitcoin"},
			Description: "Crypto-friendly slots provider",
			APIEndpoint: "https://api.endorphina.com", RTP: 96.0,
		},
		{
			ID: "platipus", Name: "Platipus", Logo: "/providers/platipus.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots", "jackpots"},
			Description: "Progressive jackpot slots",
			APIEndpoint: "https://api.platipus.com", RTP: 96.5,
		},
		{
			ID: "belatra", Name: "Belatra", Logo: "/providers/belatra.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "table_games"},
			Description: "Classic casino games specialist",
			APIEndpoint: "https://api.belatra.com", RTP: 96.0,
		},
		{
			ID: "booming_games", Name: "Booming Games", Logo: "/providers/booming.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots"},
			Description: "Premium slot content",
			APIEndpoint: "https://api.boominggames.com", RTP: 96.5,
		},
		{
			ID: "spinomenal", Name: "Spinomenal", Logo: "/providers/spinomenal.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots", "table_games"},
			Description: "Multi-category game developer",
			APIEndpoint: "https://api.spinomenal.com", RTP: 96.0,
		},
		{
			ID: "evoplay", Name: "Evoplay", Logo: "/providers/evoplay.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "instant_win", "3d_slots"},
			Description: "3D and innovative slot experiences",
			APIEndpoint: "https://api.evoplay.com", RTP: 96.5,
		},
		{
			ID: "fugaso", Name: "Fugaso", Logo: "/providers/fugaso.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots", "jackpots"},
			Description: "Jackpot and slot specialist",
			APIEndpoint: "https://api.fugaso.com", RTP: 96.5,
		},
		{
			ID: "caleta", Name: "Caleta", Logo: "/providers/caleta.png",
			GameCount: 100, IsLive: true, Categories: []string{"slots", "lottery", "instant_win"},
			Description: "Lottery and instant win games",
			APIEndpoint: "https://api.caleta.com", RTP: 96.0,
		},
		{
			ID: "mascot_gaming", Name: "Mascot Gaming", Logo: "/providers/mascot.png",
			GameCount: 80, IsLive: true, Categories: []string{"slots", "table_games"},
			Description: "Risk(R)Qi certified games",
			APIEndpoint: "https://api.mascotgaming.com", RTP: 96.5,
		},
	}

	for _, p := range providers {
		s.providers[p.ID] = p
	}
}

func (s *GameProviderService) GetAllProviders() []*GameProvider {
	providers := make([]*GameProvider, 0, len(s.providers))
	for _, p := range s.providers {
		providers = append(providers, p)
	}
	return providers
}

func (s *GameProviderService) GetProvider(id string) (*GameProvider, error) {
	p, ok := s.providers[id]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", id)
	}
	return p, nil
}

func (s *GameProviderService) GetProvidersByCategory(category string) []*GameProvider {
	var result []*GameProvider
	for _, p := range s.providers {
		for _, c := range p.Categories {
			if c == category {
				result = append(result, p)
				break
			}
		}
	}
	return result
}

// GetAllSlotGames returns all available slot games
func (s *GameProviderService) GetAllSlotGames() []*SlotGame {
	games := make([]*SlotGame, 0)

	// Add games from various providers
	games = append(games, s.getPragmaticPlaySlots()...)
	games = append(games, s.getMicrogamingSlots()...)
	games = append(games, s.getNetEntSlots()...)
	games = append(games, s.getPlayNGOSlots()...)
	games = append(games, s.getBGamingSlots()...)
	games = append(games, s.getSpribeGames()...)
	games = append(games, s.getHacksawSlots()...)
	games = append(games, s.getYggdrasilSlots()...)
	games = append(games, s.getQuickspinSlots()...)
	games = append(games, s.getRedTigerSlots()...)
	games = append(games, s.getBigTimeGamingSlots()...)
	games = append(games, s.getNolimitCitySlots()...)
	games = append(games, s.getPGSoftSlots()...)
	games = append(games, s.getBetsoftSlots()...)
	games = append(games, s.getEvoplaySlots()...)

	return games
}

func (s *GameProviderService) getPragmaticPlaySlots() []*SlotGame {
	return []*SlotGame{
		{ID: "pp_gates_of_olympus", Name: "Gates of Olympus", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 5, Paylines: 20, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Multiplier", "Free Spins"}, IsHot: true},
		{ID: "pp_sweet_bonanza", Name: "Sweet Bonanza", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 5, Paylines: 20, RTP: 96.48, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Multiplier", "Free Spins"}, IsHot: true},
		{ID: "pp_starlight_princess", Name: "Starlight Princess", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 5, Paylines: 20, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Multiplier", "Free Spins"}, IsNew: true},
		{ID: "pp_book_of_dead", Name: "Book of Dead", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.21, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Expandable Wild", "Free Spins"}, IsHot: true},
		{ID: "pp_big_bass_bonanza", Name: "Big Bass Bonanza", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.71, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 2000, Features: []string{"Free Spins", "Respin"}, IsHot: true},
		{ID: "pp_wolf_gold", Name: "Wolf Gold", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.0, Volatility: "medium", MinBet: 0.25, MaxBet: 125, MaxWin: 2500, Features: []string{"Blazing Reels", "Money Respin", "Free Spins"}, IsHot: true},
		{ID: "pp_fruit_party", Name: "Fruit Party", Provider: "pragmatic_play", Category: "slots", Reels: 7, Rows: 7, Paylines: 0, RTP: 96.47, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Multiplier", "Free Spins"}, IsHot: true},
		{ID: "pp_the_dog_house", Name: "The Dog House", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.51, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 6750, Features: []string{"Sticky Wild", "Free Spins"}, IsHot: true},
		{ID: "pp_john_hunter", Name: "John Hunter and the Aztec Treasure", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 5, Paylines: 0, RTP: 96.0, Volatility: "high", MinBet: 0.25, MaxBet: 250, MaxWin: 9000, Features: []string{"Tumble", "Free Spins", "Multipliers"}},
		{ID: "pp_great_rhino", Name: "Great Rhino", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Super Respin", "Free Spins"}},
		{ID: "pp_mustang_gold", Name: "Mustang Gold", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.53, Volatility: "high", MinBet: 0.25, MaxBet: 125, MaxWin: 4500, Features: []string{"Collect", "Free Spins", "JACKPOT"}},
		{ID: "pp_aztec_gems", Name: "Aztec Gems", Provider: "pragmatic_play", Category: "slots", Reels: 3, Rows: 3, Paylines: 5, RTP: 96.52, Volatility: "high", MinBet: 0.05, MaxBet: 25, MaxWin: 1000, Features: []string{"Multiplier", "Free Spins"}},
		{ID: "pp_five_lions_megaways", Name: "Five Lions Megaways", Provider: "pragmatic_play", Category: "megaways", Reels: 6, Rows: 7, Paylines: 200000, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Megaways", "Cascading", "Multiplier"}},
		{ID: "pp_power_of_thor", Name: "Power of Thor Megaways", Provider: "pragmatic_play", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Megaways", "Free Spins", "Multiplier"}},
		{ID: "pp_buffalo_king", Name: "Buffalo King", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 4, Paylines: 4096, RTP: 96.5, Volatility: "high", MinBet: 0.4, MaxBet: 200, MaxWin: 5000, Features: []string{"Free Spins", "Tumble", "Multiplier"}},
		{ID: "pp_fruit_party_2", Name: "Fruit Party 2", Provider: "pragmatic_play", Category: "slots", Reels: 7, Rows: 7, Paylines: 0, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Multiplier", "Free Spins"}, IsNew: true},
		{ID: "pp_christmas_carol", Name: "Christmas Carol Megaways", Provider: "pragmatic_play", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Megaways", "Free Spins", "Multiplier"}},
		{ID: "pp_magic_gems", Name: "Magic Gems", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 5, Paylines: 0, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Multiplier", "Free Spins"}},
		{ID: "pp_wild_west_gold", Name: "Wild West Gold", Provider: "pragmatic_play", Category: "slots", Reels: 5, Rows: 3, Paylines: 40, RTP: 96.51, Volatility: "high", MinBet: 0.25, MaxBet: 125, MaxWin: 10000, Features: []string{"Sticky Wild", "Free Spins"}},
		{ID: "pp_treasure_wild", Name: "Treasure Wild", Provider: "pragmatic_play", Category: "slots", Reels: 6, Rows: 5, Paylines: 20, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Tumble", "Free Spins", "Wild"}},
	}
}

func (s *GameProviderService) getMicrogamingSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "mg_mega_moolah", Name: "Mega Moolah", Provider: "microgaming", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 25, RTP: 88.12, Volatility: "high", MinBet: 0.25, MaxBet: 6.25, MaxWin: 1000000, Features: []string{"Progressive Jackpot", "Free Spins", "Multiplier"}},
		{ID: "mg_immortal_romance", Name: "Immortal Romance", Provider: "microgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.86, Volatility: "high", MinBet: 0.3, MaxBet: 30, MaxWin: 12000, Features: []string{"Free Spins", "Wild", "Multipliers"}},
		{ID: "mg_thunderstruck_2", Name: "Thunderstruck II", Provider: "microgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.65, Volatility: "high", MinBet: 0.3, MaxBet: 30, MaxWin: 1000, Features: []string{"Free Spins", "Wild", "Multipliers"}},
		{ID: "mg_game_of_thrones", Name: "Game of Thrones", Provider: "microgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.5, Volatility: "high", MinBet: 0.3, MaxBet: 30, MaxWin: 5000, Features: []string{"Free Spins", "Stacked Wild"}},
		{ID: "mg_breakaway", Name: "Break Away", Provider: "microgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.42, Volatility: "high", MinBet: 0.5, MaxBet: 50, MaxWin: 2000, Features: []string{"Stacked Wild", "Free Spins"}},
		{ID: "mg_terminator_2", Name: "Terminator 2", Provider: "microgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.5, Volatility: "high", MinBet: 0.3, MaxBet: 30, MaxWin: 1000, Features: []string{"Free Spins", "T-800 Vision"}},
		{ID: "mg_playboy", Name: "Playboy", Provider: "microgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.5, Volatility: "medium", MinBet: 0.3, MaxBet: 60, MaxWin: 5000, Features: []string{"Free Spins", "Wild"}},
		{ID: "mg_jurassic_world", Name: "Jurassic World", Provider: "microgaming", Category: "slots", Reels: 6, Rows: 3, Paylines: 243, RTP: 95.5, Volatility: "high", MinBet: 0.3, MaxBet: 37.5, MaxWin: 6000, Features: []string{"Free Spins", "Walking Wild"}},
		{ID: "mg_mega_moolah_isis", Name: "Mega Moolah Isis", Provider: "microgaming", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 25, RTP: 88.0, Volatility: "high", MinBet: 0.25, MaxBet: 6.25, MaxWin: 1000000, Features: []string{"Progressive Jackpot", "Free Spins"}},
		{ID: "mg_major_millions", Name: "Major Millions", Provider: "microgaming", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 15, RTP: 89.9, Volatility: "high", MinBet: 0.3, MaxBet: 15, MaxWin: 250000, Features: []string{"Progressive Jackpot", "Multipliers"}},
	}
}

func (s *GameProviderService) getNetEntSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "net_starburst", Name: "Starburst", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.1, Volatility: "low", MinBet: 0.1, MaxBet: 100, MaxWin: 500, Features: []string{"Expanding Wild", "Re-spins"}},
		{ID: "net_gonzo_quest", Name: "Gonzo's Quest", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.0, Volatility: "high", MinBet: 0.2, MaxBet: 50, MaxWin: 1875, Features: []string{"Avalanche", "Free Fall", "Multiplier"}},
		{ID: "net_dead_or_alive", Name: "Dead or Alive", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 12, RTP: 96.8, Volatility: "high", MinBet: 0.09, MaxBet: 18, MaxWin: 3000, Features: []string{"Free Spins", "Sticky Wild"}},
		{ID: "net_twin_spin", Name: "Twin Spin", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 243, RTP: 96.6, Volatility: "medium", MinBet: 0.25, MaxBet: 125, MaxWin: 1080, Features: []string{"Twin Reels", "Wild"}},
		{ID: "net_jack_and_beanstalk", Name: "Jack and the Beanstalk", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.3, Volatility: "high", MinBet: 0.2, MaxBet: 50, MaxWin: 1000, Features: []string{"Walking Wild", "Free Spins", "Keys"}},
		{ID: "net_aliens", Name: "Aliens", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.4, Volatility: "high", MinBet: 0.15, MaxBet: 150, MaxWin: 9000, Features: []string{"Multiplier", "Free Spins", "Cluster"}},
		{ID: "net_drive_multiplier", Name: "Drive: Multiplier Mayhem", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 200, MaxWin: 5000, Features: []string{"Free Spins", "Nitro", "Multiplier"}},
		{ID: "net_joker_pro", Name: "Joker Pro", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.8, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 1000, Features: []string{"Sticky Wild", "Free Spins"}},
		{ID: "net_secrets_of_atlantis", Name: "Secrets of Atlantis", Provider: "netent", Category: "slots", Reels: 5, Rows: 4, Paylines: 40, RTP: 97.0, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 800, Features: []string{"Nudge", "Highlight", "Free Spins"}},
		{ID: "net_blood_suckers", Name: "Blood Suckers", Provider: "netent", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 98.0, Volatility: "low", MinBet: 0.25, MaxBet: 50, MaxWin: 1000, Features: []string{"Free Spins", "Bonus Game"}},
	}
}

func (s *GameProviderService) getPlayNGOSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "png_book_of_dead", Name: "Book of Dead", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.21, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Expandable Symbol", "Free Spins"}},
		{ID: "png_rich_wildeat", Name: "Rich Wilde and the Tome of Insanity", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.5, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Expanding Symbol", "Free Spins"}},
		{ID: "png_legacy_of_dead", Name: "Legacy of Dead", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.5, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Expandable Symbol", "Free Spins"}, IsHot: true},
		{ID: "png_gems_gem", Name: "Gemix", Provider: "playngo", Category: "slots", Reels: 0, Rows: 0, Paylines: 0, RTP: 96.5, Volatility: "medium", MinBet: 0.5, MaxBet: 100, MaxWin: 5000, Features: []string{"Cascade", "Cruz", "Crystal"}},
		{ID: "png_fox_funnel", Name: "The Fox Fur", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.2, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 2000, Features: []string{"Expanding Wild", "Free Spins"}},
		{ID: "png_panda_fortune", Name: "Panda Fortune", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.5, Volatility: "medium", MinBet: 0.25, MaxBet: 100, MaxWin: 1000, Features: []string{"Free Spins", "Sticky"}},
		{ID: "png_cash_punk", Name: "Cash Punk", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Expanding Wild", "Free Spins"}},
		{ID: "png_eyes_of_horus", Name: "Eyes of Horus", Provider: "playngo", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.2, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 20000, Features: []string{"Expandable Symbol", "Free Spins", "Gamble"}},
		{ID: "png_mystery_joker", Name: "Mystery Joker", Provider: "playngo", Category: "slots", Reels: 3, Rows: 3, Paylines: 5, RTP: 96.5, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Joker", "Multiplier"}},
		{ID: "png_rainbow_gold", Name: "Rainbow Gold", Provider: "playngo", Category: "slots", Reels: 6, Rows: 4, Paylines: 20, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Free Spins", "Pick"}},
	}
}

func (s *GameProviderService) getBGamingSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "bg_aztec_magic", Name: "Aztec Magic", Provider: "bgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.96, Volatility: "medium", MinBet: 0.1, MaxBet: 25, MaxWin: 5000, Features: []string{"Free Spins", "Bonus Game"}},
		{ID: "bg_aztec_magic_megaways", Name: "Aztec Magic Megaways", Provider: "bgaming", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.7, Volatility: "high", MinBet: 0.2, MaxBet: 20, MaxWin: 5000, Features: []string{"Megaways", "Cascade", "Free Spins"}},
		{ID: "bg_egypt_horus", Name: "Elvis Frog in Vegas", Provider: "bgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.0, Volatility: "high", MinBet: 0.2, MaxBet: 25, MaxWin: 5000, Features: []string{"Free Spins", "Multiplier"}},
		{ID: "bg_brave_ninja", Name: "Brave Ninja", Provider: "bgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.95, Volatility: "medium", MinBet: 0.2, MaxBet: 25, MaxWin: 4000, Features: []string{"Free Spins", "Sticky"}},
		{ID: "bg_diamond_flash", Name: "Diamond Flash", Provider: "bgaming", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 25, MaxWin: 5000, Features: []string{"Free Spins", "Multiplier"}},
		{ID: "bg_lucky_bunny", Name: "Lucky Lady Clover", Provider: "bgaming", Category: "slots", Reels: 6, Rows: 5, Paylines: 50, RTP: 96.0, Volatility: "high", MinBet: 0.1, MaxBet: 25, MaxWin: 5000, Features: []string{"Free Spins", "Bonus"}},
	}
}

func (s *GameProviderService) getSpribeGames() []*SlotGame {
	return []*SlotGame{
		{ID: "sp_aviator", Name: "Aviator", Provider: "spribe", Category: "crash", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Auto Cashout", "Multiplayer"}, IsHot: true},
		{ID: "sp_mines", Name: "Mines", Provider: "spribe", Category: "mine", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Customizable Mines", "Auto"}},
		{ID: "sp_dice", Name: "Dice", Provider: "spribe", Category: "dice", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Slider", "Auto"}},
		{ID: "sp_plinko", Name: "Plinko", Provider: "spribe", Category: "plinko", Reels: 0, Rows: 0, Paylines: 0, RTP: 98.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Risk Level", "Auto"}},
		{ID: "sp_hilo", Name: "Hi-Lo", Provider: "spribe", Category: "hi-lo", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Card Prediction"}},
		{ID: "sp_keno", Name: "Keno", Provider: "spribe", Category: "keno", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Multi-ball"}},
		{ID: "sp_limbo", Name: "Limbo", Provider: "spribe", Category: "limbo", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Target Multiplier"}},
		{ID: "sp_stairs", Name: "Stairs", Provider: "spribe", Category: "arcade", Reels: 0, Rows: 0, Paylines: 0, RTP: 97.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Climb"}},
	}
}

func (s *GameProviderService) getHacksawSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "hack_wanted_dead", Name: "Wanted Dead or a Wild", Provider: "hacksaw", Category: "slots", Reels: 5, Rows: 4, Paylines: 14, RTP: 96.24, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 12500, Features: []string{"Duel", "Free Spins", "RTG"}},
		{ID: "hack_sticky_bands", Name: "Sticky Bandits", Provider: "hacksaw", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.2, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Sticky Wild", "Free Spins"}},
		{ID: "hack_jet_lucky", Name: "Jet Lucky", Provider: "hacksaw", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.3, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Bonus", "Multiplier"}},
		{ID: "hack_cash_quest", Name: "Cash Quest", Provider: "hacksaw", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.5, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 10000, Features: []string{"Free Spins", "Multiplier"}},
		{ID: "hack_punch_bob", Name: "Punch Bob", Provider: "hacksaw", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.5, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Free Spins", "Bonus"}},
	}
}

func (s *GameProviderService) getYggdrasilSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "yg_valley_of_gods", Name: "Valley of the Gods", Provider: "yggdrasil", Category: "slots", Reels: 6, Rows: 5, Paylines: 3125, RTP: 96.1, Volatility: "high", MinBet: 0.1, MaxBet: 50, MaxWin: 5000, Features: []string{"Extinctions", "Respins", "Win Multiplier"}},
		{ID: "yg_holmes", Name: "Holmes and the Stolen Stones", Provider: "yggdrasil", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.8, Volatility: "medium", MinBet: 0.2, MaxBet: 20, MaxWin: 5000, Features: []string{"Jackpot", "Free Spins", "Wild"}},
		{ID: "yg_tutan_khan", Name: "Tutankhamun", Provider: "yggdrasil", Category: "slots", Reels: 6, Rows: 4, Paylines: 50, RTP: 96.0, Volatility: "high", MinBet: 0.2, MaxBet: 40, MaxWin: 2000, Features: []string{"Expanding Symbol", "Free Spins"}},
		{ID: "yg_joker_millions", Name: "Joker Millions", Provider: "yggdrasil", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 20, RTP: 97.0, Volatility: "high", MinBet: 0.25, MaxBet: 25, MaxWin: 50000, Features: []string{"Progressive Jackpot", "Respin"}},
		{ID: "yg_egypt_sky", Name: "Egyption Heroes", Provider: "yggdrasil", Category: "slots", Reels: 5, Rows: 3, Paylines: 30, RTP: 96.5, Volatility: "medium", MinBet: 0.3, MaxBet: 60, MaxWin: 2000, Features: []string{"Free Spins", "Stacked"}},
	}
}

func (s *GameProviderService) getQuickspinSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "qs_big_bang", Name: "Big Bang", Provider: "quickspin", Category: "slots", Reels: 3, Rows: 3, Paylines: 5, RTP: 96.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 1000, Features: []string{"Multiplier", "Progressive"}},
		{ID: "qs_silence", Name: "Silence", Provider: "quickspin", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 1000, Features: []string{"Free Spins", "Sticky"}},
		{ID: "qs_tickets", Name: "Tickets of Fortune", Provider: "quickspin", Category: "slots", Reels: 5, Rows: 3, Paylines: 30, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 500, Features: []string{"Collect", "Multipliers"}},
		{ID: "qs_second_strike", Name: "Second Strike", Provider: "quickspin", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.5, Volatility: "high", MinBet: 0.1, MaxBet: 50, MaxWin: 1260, Features: []string{"Second Strike", "Jackpot"}},
		{ID: "qs_golden_lotus", Name: "Golden Lotus", Provider: "quickspin", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.5, Volatility: "medium", MinBet: 0.25, MaxBet: 100, MaxWin: 1000, Features: []string{"Free Spins", "Sticky"}},
	}
}

func (s *GameProviderService) getRedTigerSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "rt_daily_jackpot", Name: "Daily Jackpot", Provider: "red_tiger", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.0, Volatility: "medium", MinBet: 0.2, MaxBet: 10, MaxWin: 500, Features: []string{"Daily Jackpot", "Must Drop"}},
		{ID: "rt_gonzos", Name: "Gonzo's Quest Megaways", Provider: "red_tiger", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.0, Volatility: "high", MinBet: 0.2, MaxBet: 20, MaxWin: 21000, Features: []string{"Megaways", "Avalanche", "Free Spins"}},
		{ID: "rt_rainbow_jackpots", Name: "Rainbow Jackpots", Provider: "red_tiger", Category: "jackpot", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.0, Volatility: "medium", MinBet: 0.2, MaxBet: 50, MaxWin: 1000, Features: []string{"Jackpot", "Free Spins", "Symbol Transformation"}},
		{ID: "rt_wild_west", Name: "Wild West", Provider: "red_tiger", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.2, Volatility: "medium", MinBet: 0.1, MaxBet: 50, MaxWin: 500, Features: []string{"Multiplier", "Respin"}},
		{ID: "rt_lucky_angels", Name: "Lucky Angels", Provider: "red_tiger", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.0, Volatility: "medium", MinBet: 0.2, MaxBet: 50, MaxWin: 500, Features: []string{"Stacked", "Respin"}},
	}
}

func (s *GameProviderService) getBigTimeGamingSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "btg_white_rabbit", Name: "White Rabbit", Provider: "big_time_gaming", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 97.7, Volatility: "high", MinBet: 0.1, MaxBet: 20, MaxWin: 50000, Features: []string{"Megaways", "Cascade", "Free Spins"}},
		{ID: "btg_extra_chilli", Name: "Extra Chilli", Provider: "big_time_gaming", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.8, Volatility: "high", MinBet: 0.2, MaxBet: 20, MaxWin: 5000, Features: []string{"Megaways", "Free Spins", "Feature Drop"}},
		{ID: "btg_book_of_relics", Name: "Book of Relics", Provider: "big_time_gaming", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 20, MaxWin: 5000, Features: []string{"Megaways", "Free Spins"}},
		{ID: "btg_monkey_king", Name: "The Monkey King", Provider: "big_time_gaming", Category: "megaways", Reels: 6, Rows: 7, Paylines: 117649, RTP: 96.5, Volatility: "high", MinBet: 0.2, MaxBet: 20, MaxWin: 5000, Features: []string{"Megaways", "Reel Clone"}},
		{ID: "btg_star_cluster", Name: "Star Cluster", Provider: "big_time_gaming", Category: "megaways", Reels: 8, Rows: 8, Paylines: 0, RTP: 96.25, Volatility: "high", MinBet: 0.2, MaxBet: 20, MaxWin: 5000, Features: []string{"Cluster Pays", "Avalanche"}},
	}
}

func (s *GameProviderService) getNolimitCitySlots() []*SlotGame {
	return []*SlotGame{
		{ID: "nl_mind_of_god", Name: "Mental", Provider: "nolimit_city", Category: "slots", Reels: 6, Rows: 3, Paylines: 0, RTP: 96.0, Volatility: "high", MinBet: 0.1, MaxBet: 70, MaxWin: 5000, Features: []string{"Nudge", "Joker", "Free Spins"}},
		{ID: "nl_xxtreme", Name: "xxxtreme", Provider: "nolimit_city", Category: "slots", Reels: 6, Rows: 3, Paylines: 0, RTP: 96.0, Volatility: "high", MinBet: 0.1, MaxBet: 70, MaxWin: 5000, Features: []string{"Spins", "Multiplier"}},
		{ID: "nl_el_danger", Name: "El Patron", Provider: "nolimit_city", Category: "slots", Reels: 5, Rows: 3, Paylines: 0, RTP: 96.0, Volatility: "high", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Nudge", "Bonus"}},
		{ID: "nl_tomb", Name: "Tomb of Dead", Provider: "nolimit_city", Category: "slots", Reels: 6, Rows: 3, Paylines: 0, RTP: 96.0, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Expand", "Free Spins"}},
		{ID: "nl_jpush", Name: "JPush", Provider: "nolimit_city", Category: "slots", Reels: 6, Rows: 3, Paylines: 0, RTP: 96.0, Volatility: "high", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Joker", "Multiplier"}},
	}
}

func (s *GameProviderService) getPGSoftSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "pg_fortune_ox", Name: "Fortune Ox", Provider: "pg_soft", Category: "slots", Reels: 5, Rows: 3, Paylines: 10, RTP: 96.5, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Multiplier", "Free Spins"}},
		{ID: "pg_fortune_mouse", Name: "Fortune Mouse", Provider: "pg_soft", Category: "slots", Reels: 3, Rows: 3, Paylines: 9, RTP: 96.5, Volatility: "medium", MinBet: 0.25, MaxBet: 250, MaxWin: 1000, Features: []string{"Respin", "Multiplier"}},
		{ID: "pg_santas_gift", Name: "Santa's Gift", Provider: "pg_soft", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.5, Volatility: "medium", MinBet: 0.25, MaxBet: 250, MaxWin: 2000, Features: []string{"Free Spins", "Multiplier"}},
		{ID: "pg_guardian_of_atlantis", Name: "Guardian of Atlantis", Provider: "pg_soft", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 2000, Features: []string{"Free Spins", "Respin"}},
		{ID: "pg_mayan_calendars", Name: "Mayan Calendars", Provider: "pg_soft", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 2000, Features: []string{"Cascade", "Jackpot"}},
	}
}

func (s *GameProviderService) getBetsoftSlots() []*SlotGame {
	return []*SlotGame{
		{ID: "bs_gold_digger", Name: "Gold Digger", Provider: "betsoft", Category: "slots", Reels: 6, Rows: 5, Paylines: 30, RTP: 96.0, Volatility: "high", MinBet: 0.2, MaxBet: 30, MaxWin: 7500, Features: []string{"Link & Win", "Free Spins", "Respin"}},
		{ID: "bs_atlantis", Name: "Quest to the West", Provider: "betsoft", Category: "slots", Reels: 5, Rows: 3, Paylines: 25, RTP: 97.5, Volatility: "high", MinBet: 0.25, MaxBet: 25, MaxWin: 10000, Features: []string{"Walking Wild", "Multiplier"}},
		{ID: "bs_tiger_claw", Name: "Tiger Claw", Provider: "betsoft", Category: "slots", Reels: 6, Rows: 3, Paylines: 729, RTP: 97.0, Volatility: "high", MinBet: 0.3, MaxBet: 45, MaxWin: 5000, Features: []string{"Frogin", "Multiplier"}},
		{ID: "bs_greek_gods", Name: "Greek Gods", Provider: "betsoft", Category: "slots", Reels: 6, Rows: 5, Paylines: 25, RTP: 96.0, Volatility: "medium", MinBet: 0.25, MaxBet: 125, MaxWin: 2500, Features: []string{"Free Spins", "Jackpot"}},
		{ID: "bs_windy_bee", Name: "Windy Bee", Provider: "betsoft", Category: "slots", Reels: 6, Rows: 5, Paylines: 25, RTP: 96.5, Volatility: "medium", MinBet: 0.2, MaxBet: 100, MaxWin: 5000, Features: []string{"Cluster", "Respin"}},
	}
}

func (s *GameProviderService) getEvoplaySlots() []*SlotGame {
	return []*SlotGame{
		{ID: "evo_fruit_cocktail", Name: "Fruit Cocktail", Provider: "evoplay", Category: "slots", Reels: 5, Rows: 3, Paylines: 9, RTP: 96.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Bonus", "Free Spins"}},
		{ID: "evo_garage", Name: "Garage", Provider: "evoplay", Category: "slots", Reels: 5, Rows: 3, Paylines: 9, RTP: 96.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Bonus", "Respin"}},
		{ID: "evo_naughty_girls", Name: "Naughty Girls", Provider: "evoplay", Category: "slots", Reels: 5, Rows: 3, Paylines: 20, RTP: 96.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Free Spins", "Multiplier"}},
		{ID: "evo_crazy_monkey", Name: "Crazy Monkey", Provider: "evoplay", Category: "slots", Reels: 5, Rows: 3, Paylines: 9, RTP: 96.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Bonus", "Gamble"}},
		{ID: "evo_lucky_drink", Name: "Lucky Drink", Provider: "evoplay", Category: "slots", Reels: 5, Rows: 3, Paylines: 15, RTP: 96.0, Volatility: "medium", MinBet: 0.1, MaxBet: 100, MaxWin: 5000, Features: []string{"Bonus", "Free Spins"}},
	}
}
