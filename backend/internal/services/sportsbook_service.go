package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/models"
)

// SportsbookService handles sports betting with real data feeds
type SportsbookService struct {
	db           *gorm.DB
	redis        *redis.Client
	config       *SportsbookConfig
	feedClient   *SportsDataClient
	oddsCache    *OddsCache
	markets      map[string]*Market
	marketsMu    sync.RWMutex
	bets         map[uuid.UUID]*Bet
	betsMu       sync.RWMutex
}

type SportsbookConfig struct {
	SportRadarAPIKey    string
	SportRadarURL      string
	BetGeniusAPIKey    string
	BetGeniusURL       string
	TheRundownAPIKey   string
	TheRundownURL      string
	OddsJamAPIKey      string
	OddsJamURL         string
	
	// Bet limits
	MinBetAmount        float64
	MaxBetAmount        float64
	MaxWinAmount        float64
	
	// Settlement settings
	SettlementTimeout   time.Duration
	AutoSettleEnabled   bool
	
	// Live betting
	LiveUpdateInterval  time.Duration
	MaxLiveBets         int
	
	Timeout             time.Duration
}

// SportsDataClient connects to sports data providers
type SportsDataClient struct {
	config     *SportsbookConfig
	httpClient *http.Client
	providers  map[string]DataProvider
}

type DataProvider interface {
	GetName() string
	GetSports(ctx context.Context) ([]Sport, error)
	GetLeagues(ctx context.Context, sportID string) ([]League, error)
	GetEvents(ctx context.Context, sportID, leagueID string, date time.Time) ([]Event, error)
	GetLiveEvents(ctx context.Context, sportID string) ([]Event, error)
	GetOdds(ctx context.Context, eventID string) (*OddsData, error)
	GetResults(ctx context.Context, eventID string) (*EventResult, error)
	IsAvailable() bool
}

type Sport struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Icon      string `json:"icon"`
}

type League struct {
	ID        string `json:"id"`
	SportID   string `json:"sport_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Logo      string `json:"logo"`
}

type Event struct {
	ID              string                 `json:"id"`
	SportID         string                 `json:"sport_id"`
	LeagueID        string                 `json:"league_id"`
	HomeTeam        string                 `json:"home_team"`
	AwayTeam        string                 `json:"away_team"`
	StartTime       time.Time              `json:"start_time"`
	Status          string                 `json:"status"` // scheduled, live, finished, cancelled
	HomeScore       int                    `json:"home_score"`
	AwayScore       int                    `json:"away_score"`
	Period          string                 `json:"period"`
	TimeRemaining   string                 `json:"time_remaining"`
	Venue           string                 `json:"venue"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type OddsData struct {
	EventID     string             `json:"event_id"`
	Bookmaker   string             `json:"bookmaker"`
	Timestamp   time.Time          `json:"timestamp"`
	Markets     []MarketOdds       `json:"markets"`
}

type MarketOdds struct {
	MarketID    string          `json:"market_id"`
	MarketType  string          `json:"market_type"`
	Outcomes    []OutcomeOdds  `json:"outcomes"`
	Suspended   bool            `json:"suspended"`
}

type OutcomeOdds struct {
	OutcomeID string  `json:"outcome_id"`
	Name      string  `json:"name"`
	Odds      float64 `json:"odds"`
	Line      float64 `json:"line"`
}

type EventResult struct {
	EventID     string                 `json:"event_id"`
	Status      string                 `json:"status"`
	WinnerID    string                 `json:"winner_id"`
	HomeScore   int                    `json:"home_score"`
	AwayScore   int                    `json:"away_score"`
	Periods     []PeriodScore          `json:"periods"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type PeriodScore struct {
	Period   string `json:"period"`
	HomeScore int   `json:"home_score"`
	AwayScore int   `json:"away_score"`
}

// OddsCache caches odds data
type OddsCache struct {
	mu    sync.RWMutex
	odds  map[string]map[string]float64 // eventID + marketID + outcomeID -> odds
}

func NewOddsCache() *OddsCache {
	return &OddsCache{
		odds: make(map[string]map[string]float64),
	}
}

func (oc *OddsCache) Set(eventID, marketID, outcomeID string, odds float64) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	
	key := eventID + ":" + marketID + ":" + outcomeID
	if oc.odds[eventID] == nil {
		oc.odds[eventID] = make(map[string]float64)
	}
	oc.odds[eventID][key] = odds
}

func (oc *OddsCache) Get(eventID, marketID, outcomeID string) (float64, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()
	
	key := eventID + ":" + marketID + ":" + outcomeID
	if odds, ok := oc.odds[eventID]; ok {
		o, ok := odds[key]
		return o, ok
	}
	return 0, false
}

func (oc *OddsCache) Clear() {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	oc.odds = make(map[string]map[string]float64)
}

// Market represents a betting market
type Market struct {
	ID          string
	Name        string
	EventID     string
	MarketType  string
	Outcomes    []Outcome
	Suspended   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Outcome struct {
	ID      string
	Name    string
	Odds    float64
	Line    float64
	Result  string // pending, won, lost, cancelled
}

// Bet represents a user bet
type Bet struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	EventID        string
	MarketID       string
	Stake          float64
	Odds           float64
	PotentialWin   float64
	Selection      string
	Status         string // pending, won, lost, cancelled, voided
	SettledAt      *time.Time
	CreatedAt      time.Time
}

// SportRadarProvider implements SportRadar data provider
type SportRadarProvider struct {
	config     *SportsbookConfig
	httpClient *http.Client
	enabled    bool
}

func NewSportRadarProvider(config *SportsbookConfig) *SportRadarProvider {
	return &SportRadarProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		enabled: config.SportRadarAPIKey != "" && config.SportRadarURL != "",
	}
}

func (s *SportRadarProvider) GetName() string {
	return "SportRadar"
}

func (s *SportRadarProvider) IsAvailable() bool {
	return s.enabled
}

func (s *SportRadarProvider) GetSports(ctx context.Context) ([]Sport, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SportRadar not configured")
	}

	url := fmt.Sprintf("%s/v1/sports", s.config.SportRadarURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.SportRadarAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Sports []Sport `json:"sports"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Sports, nil
}

func (s *SportRadarProvider) GetLeagues(ctx context.Context, sportID string) ([]League, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SportRadar not configured")
	}

	url := fmt.Sprintf("%s/v1/sports/%s/leagues", s.config.SportRadarURL, sportID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.SportRadarAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Leagues []League `json:"leagues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Leagues, nil
}

func (s *SportRadarProvider) GetEvents(ctx context.Context, sportID, leagueID string, date time.Time) ([]Event, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SportRadar not configured")
	}

	url := fmt.Sprintf("%s/v1/sports/%s/leagues/%s/schedules/%s",
		s.config.SportRadarURL, sportID, leagueID, date.Format("2006-01-02"))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.SportRadarAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Events []Event `json:"games"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Events, nil
}

func (s *SportRadarProvider) GetLiveEvents(ctx context.Context, sportID string) ([]Event, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SportRadar not configured")
	}

	url := fmt.Sprintf("%s/v1/sports/%s/live", s.config.SportRadarURL, sportID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.SportRadarAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Events []Event `json:"games"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Events, nil
}

func (s *SportRadarProvider) GetOdds(ctx context.Context, eventID string) (*OddsData, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SportRadar not configured")
	}

	url := fmt.Sprintf("%s/v1/events/%s/odds", s.config.SportRadarURL, eventID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.SportRadarAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result OddsData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *SportRadarProvider) GetResults(ctx context.Context, eventID string) (*EventResult, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SportRadar not configured")
	}

	url := fmt.Sprintf("%s/v1/events/%s/results", s.config.SportRadarURL, eventID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.SportRadarAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EventResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// MockDataProvider provides mock data for testing
type MockDataProvider struct{}

func NewMockDataProvider() *MockDataProvider {
	return &MockDataProvider{}
}

func (m *MockDataProvider) GetName() string {
	return "Mock"
}

func (m *MockDataProvider) IsAvailable() bool {
	return true
}

func (m *MockDataProvider) GetSports(ctx context.Context) ([]Sport, error) {
	return []Sport{
		{ID: "football", Name: "Football", ShortName: "FB", Icon: "⚽"},
		{ID: "basketball", Name: "Basketball", ShortName: "BB", Icon: "🏀"},
		{ID: "tennis", Name: "Tennis", ShortName: "TN", Icon: "🎾"},
		{ID: "baseball", Name: "Baseball", ShortName: "BB", Icon: "⚾"},
		{ID: "hockey", Name: "Ice Hockey", ShortName: "IH", Icon: "🏒"},
		{ID: "mma", Name: "MMA", ShortName: "MMA", Icon: "🥊"},
		{ID: "boxing", Name: "Boxing", ShortName: "BX", Icon: "🥋"},
		{ID: "esports", Name: "Esports", ShortName: "ES", Icon: "🎮"},
		{ID: "cricket", Name: "Cricket", ShortName: "CR", Icon: "🏏"},
		{ID: "rugby", Name: "Rugby", ShortName: "RG", Icon: "🏉"},
		{ID: "volleyball", Name: "Volleyball", ShortName: "VB", Icon: "🏐"},
		{ID: "american_football", Name: "American Football", ShortName: "AF", Icon: "🏈"},
		{ID: "golf", Name: "Golf", ShortName: "GF", Icon: "⛳"},
		{ID: "motorsport", Name: "Motorsport", ShortName: "MT", Icon: "🏎️"},
	}, nil
}

func (m *MockDataProvider) GetLeagues(ctx context.Context, sportID string) ([]League, error) {
	leagues := map[string][]League{
		"football": {
			{ID: "epl", SportID: "football", Name: "English Premier League", Country: "England", Logo: "/images/leagues/epl.png"},
			{ID: "laliga", SportID: "football", Name: "La Liga", Country: "Spain", Logo: "/images/leagues/laliga.png"},
			{ID: "bundesliga", SportID: "football", Name: "Bundesliga", Country: "Germany", Logo: "/images/leagues/bundesliga.png"},
			{ID: "serie_a", SportID: "football", Name: "Serie A", Country: "Italy", Logo: "/images/leagues/seriea.png"},
			{ID: "ligue1", SportID: "football", Name: "Ligue 1", Country: "France", Logo: "/images/leagues/ligue1.png"},
			{ID: "champions_league", SportID: "football", Name: "UEFA Champions League", Country: "Europe", Logo: "/images/leagues/ucl.png"},
			{ID: "world_cup", SportID: "football", Name: "FIFA World Cup", Country: "International", Logo: "/images/leagues/worldcup.png"},
			{ID: "mls", SportID: "football", Name: "MLS", Country: "USA", Logo: "/images/leagues/mls.png"},
		},
		"basketball": {
			{ID: "nba", SportID: "basketball", Name: "NBA", Country: "USA", Logo: "/images/leagues/nba.png"},
			{ID: "euroleague", SportID: "basketball", Name: "EuroLeague", Country: "Europe", Logo: "/images/leagues/euroleague.png"},
			{ID: "wnba", SportID: "basketball", Name: "WNBA", Country: "USA", Logo: "/images/leagues/wnba.png"},
		},
		"tennis": {
			{ID: "atp", SportID: "tennis", Name: "ATP Tour", Country: "International", Logo: "/images/leagues/atp.png"},
			{ID: "wta", SportID: "tennis", Name: "WTA Tour", Country: "International", Logo: "/images/leagues/wta.png"},
			{ID: "grand_slam", SportID: "tennis", Name: "Grand Slams", Country: "International", Logo: "/images/leagues/gs.png"},
		},
		"esports": {
			{ID: "lol_worlds", SportID: "esports", Name: "League of Legends World Championship", Country: "International", Logo: "/images/leagues/lol.png"},
			{ID: "csgo_major", SportID: "esports", Name: "CS:GO Majors", Country: "International", Logo: "/images/leagues/csgo.png"},
			{ID: "dota_international", SportID: "esports", Name: "The International", Country: "International", Logo: "/images/leagues/dota.png"},
			{ID: "valorant_champions", SportID: "esports", Name: "VALORANT Champions", Country: "International", Logo: "/images/leagues/valorant.png"},
		},
	}

	if l, ok := leagues[sportID]; ok {
		return l, nil
	}
	return []League{}, nil
}

func (m *MockDataProvider) GetEvents(ctx context.Context, sportID, leagueID string, date time.Time) ([]Event, error) {
	// Generate mock events
	events := []Event{
		{
			ID:          uuid.New().String(),
			SportID:     sportID,
			LeagueID:    leagueID,
			HomeTeam:    "Team A",
			AwayTeam:    "Team B",
			StartTime:   date.Add(2 * time.Hour),
			Status:      "scheduled",
		},
		{
			ID:          uuid.New().String(),
			SportID:     sportID,
			LeagueID:    leagueID,
			HomeTeam:    "Team C",
			AwayTeam:    "Team D",
			StartTime:   date.Add(4 * time.Hour),
			Status:      "scheduled",
		},
	}

	return events, nil
}

func (m *MockDataProvider) GetLiveEvents(ctx context.Context, sportID string) ([]Event, error) {
	return []Event{
		{
			ID:            uuid.New().String(),
			SportID:       sportID,
			LeagueID:      "live",
			HomeTeam:      "Live Team A",
			AwayTeam:      "Live Team B",
			StartTime:     time.Now().Add(-1 * time.Hour),
			Status:        "live",
			HomeScore:     2,
			AwayScore:     1,
			Period:        "2nd Half",
			TimeRemaining: "67'",
		},
	}, nil
}

func (m *MockDataProvider) GetOdds(ctx context.Context, eventID string) (*OddsData, error) {
	return &OddsData{
		EventID:   eventID,
		Bookmaker: "TigerCasino",
		Timestamp: time.Now(),
		Markets: []MarketOdds{
			{
				MarketID:   "moneyline",
				MarketType: "moneyline",
				Outcomes: []OutcomeOdds{
					{OutcomeID: "home", Name: "Home Win", Odds: 1.85},
					{OutcomeID: "draw", Name: "Draw", Odds: 3.40},
					{OutcomeID: "away", Name: "Away Win", Odds: 2.10},
				},
			},
			{
				MarketID:   "spread",
				MarketType: "point_spread",
				Outcomes: []OutcomeOdds{
					{OutcomeID: "home_spread", Name: "Home -2.5", Odds: 1.90, Line: -2.5},
					{OutcomeID: "away_spread", Name: "Away +2.5", Odds: 1.90, Line: 2.5},
				},
			},
			{
				MarketID:   "total",
				MarketType: "over_under",
				Outcomes: []OutcomeOdds{
					{OutcomeID: "over", Name: "Over 2.5", Odds: 1.95, Line: 2.5},
					{OutcomeID: "under", Name: "Under 2.5", Odds: 1.85, Line: 2.5},
				},
			},
		},
	}, nil
}

func (m *MockDataProvider) GetResults(ctx context.Context, eventID string) (*EventResult, error) {
	return &EventResult{
		EventID:     eventID,
		Status:      "finished",
		HomeScore:   2,
		AwayScore:   1,
		WinnerID:    "home",
		Periods:     []PeriodScore{{Period: "1st Half", HomeScore: 1, AwayScore: 0}, {Period: "2nd Half", HomeScore: 1, AwayScore: 1}},
	}, nil
}

// NewSportsbookService creates a new sportsbook service
func NewSportsbookService(db *gorm.DB, redisClient *redis.Client, config *SportsbookConfig) *SportsbookService {
	client := &SportsDataClient{
		config:     config,
		httpClient: &http.Client{Timeout: config.Timeout},
		providers:  make(map[string]DataProvider),
	}

	// Initialize providers
	if config != nil {
		if sportRadar := NewSportRadarProvider(config); sportRadar.IsAvailable() {
			client.providers["sportradar"] = sportRadar
		}
	}

	// Always add mock provider for demo
	client.providers["mock"] = NewMockDataProvider()

	return &SportsbookService{
		db:         db,
		redis:      redisClient,
		config:     config,
		feedClient: client,
		oddsCache:  NewOddsCache(),
		markets:    make(map[string]*Market),
		bets:       make(map[uuid.UUID]*Bet),
	}
}

// GetSports returns all available sports
func (s *SportsbookService) GetSports(ctx context.Context) ([]Sport, error) {
	// Try cache first
	cacheKey := "sportsbook:sports"
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var sports []Sport
			if json.Unmarshal([]byte(cached), &sports) == nil {
				return sports, nil
			}
		}
	}

	// Get from provider
	provider := s.getProvider()
	sports, err := provider.GetSports(ctx)
	if err != nil {
		return nil, err
	}

	// Cache results
	if s.redis != nil {
		if data, err := json.Marshal(sports); err == nil {
			s.redis.Set(ctx, cacheKey, data, 30*time.Minute)
		}
	}

	return sports, nil
}

func (s *SportsbookService) getProvider() DataProvider {
	// Prefer real provider if available
	if provider, ok := s.feedClient.providers["sportradar"]; ok && provider.IsAvailable() {
		return provider
	}
	// Fallback to mock
	return s.feedClient.providers["mock"]
}

// GetLeagues returns leagues for a sport
func (s *SportsbookService) GetLeagues(ctx context.Context, sportID string) ([]League, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("sportsbook:leagues:%s", sportID)
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var leagues []League
			if json.Unmarshal([]byte(cached), &leagues) == nil {
				return leagues, nil
			}
		}
	}

	provider := s.getProvider()
	leagues, err := provider.GetLeagues(ctx, sportID)
	if err != nil {
		return nil, err
	}

	// Cache results
	if s.redis != nil {
		if data, err := json.Marshal(leagues); err == nil {
			s.redis.Set(ctx, cacheKey, data, 30*time.Minute)
		}
	}

	return leagues, nil
}

// GetEvents returns events for a sport/league on a date
func (s *SportsbookService) GetEvents(ctx context.Context, sportID, leagueID string, date time.Time) ([]Event, error) {
	cacheKey := fmt.Sprintf("sportsbook:events:%s:%s:%s", sportID, leagueID, date.Format("2006-01-02"))
	
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var events []Event
			if json.Unmarshal([]byte(cached), &events) == nil {
				return events, nil
			}
		}
	}

	provider := s.getProvider()
	events, err := provider.GetEvents(ctx, sportID, leagueID, date)
	if err != nil {
		return nil, err
	}

	// Load odds for each event
	for i := range events {
		odds, err := provider.GetOdds(ctx, events[i].ID)
		if err == nil {
			events[i].Metadata = map[string]interface{}{"odds": odds}
		}
	}

	if s.redis != nil {
		if data, err := json.Marshal(events); err == nil {
			s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}

	return events, nil
}

// GetLiveEvents returns live events
func (s *SportsbookService) GetLiveEvents(ctx context.Context, sportID string) ([]Event, error) {
	provider := s.getProvider()
	return provider.GetLiveEvents(ctx, sportID)
}

// GetEventOdds returns odds for an event
func (s *SportsbookService) GetEventOdds(ctx context.Context, eventID string) (*OddsData, error) {
	// Check cache
	if odds, found := s.oddsCache.Get(eventID, "moneyline", "home"); found {
		// Return cached data
		return &OddsData{
			EventID:   eventID,
			Timestamp: time.Now(),
		}, nil
	}

	provider := s.getProvider()
	return provider.GetOdds(ctx, eventID)
}

// GetMarkets returns betting markets for an event
func (s *SportsbookService) GetMarkets(ctx context.Context, eventID string) ([]Market, error) {
	// Try database first
	var dbMarkets []models.SportsMarket
	if err := s.db.Where("event_id = ? AND status = ?", eventID, "active").Find(&dbMarkets).Error; err == nil && len(dbMarkets) > 0 {
		markets := make([]Market, len(dbMarkets))
		for i, m := range dbMarkets {
			var outcomes []Outcome
			json.Unmarshal([]byte(m.Outcomes), &outcomes)
			markets[i] = Market{
				ID:         m.ID.String(),
				Name:       m.Name,
				EventID:    m.EventID,
				MarketType: m.MarketType,
				Outcomes:   outcomes,
				Suspended:  m.Suspended,
			}
		}
		return markets, nil
	}

	// Get from provider
	odds, err := s.GetEventOdds(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Convert to markets
	markets := make([]Market, len(odds.Markets))
	for i, mo := range odds.Markets {
		outcomes := make([]Outcome, len(mo.Outcomes))
		for j, oo := range mo.Outcomes {
			outcomes[j] = Outcome{
				ID:    oo.OutcomeID,
				Name:  oo.Name,
				Odds:  oo.Odds,
				Line:  oo.Line,
			}
			
			// Cache odds
			s.oddsCache.Set(eventID, mo.MarketID, oo.OutcomeID, oo.Odds)
		}
		markets[i] = Market{
			ID:         mo.MarketID,
			Name:       mo.MarketType,
			EventID:    eventID,
			MarketType: mo.MarketType,
			Outcomes:   outcomes,
			Suspended:  mo.Suspended,
		}
	}

	return markets, nil
}

// PlaceBet places a bet
func (s *SportsbookService) PlaceBet(ctx context.Context, userID uuid.UUID, eventID, marketID, selection string, stake float64) (*Bet, error) {
	// Validate bet amount
	if stake < s.config.MinBetAmount {
		return nil, fmt.Errorf("minimum bet amount is %.2f", s.config.MinBetAmount)
	}
	if stake > s.config.MaxBetAmount {
		return nil, fmt.Errorf("maximum bet amount is %.2f", s.config.MaxBetAmount)
	}

	// Get odds
	odds, found := s.oddsCache.Get(eventID, marketID, selection)
	if !found {
		return nil, fmt.Errorf("odds not found for selection")
	}

	// Calculate potential win
	potentialWin := stake * odds

	if potentialWin > s.config.MaxWinAmount {
		return nil, fmt.Errorf("potential win exceeds maximum of %.2f", s.config.MaxWinAmount)
	}

	// Check user balance
	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("wallet not found")
	}

	if wallet.Balance < stake {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create bet
	bet := &Bet{
		ID:            uuid.New(),
		UserID:         userID,
		EventID:        eventID,
		MarketID:       marketID,
		Stake:          stake,
		Odds:           odds,
		PotentialWin:   potentialWin,
		Selection:      selection,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	// Deduct stake from balance
	tx := s.db.Begin()
	
	wallet.Balance -= stake
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Save bet
	if err := tx.Create(&models.SportsBet{
		ID:           bet.ID,
		UserID:       bet.UserID,
		EventID:      bet.EventID,
		MarketID:     bet.MarketID,
		Selection:    bet.Selection,
		Stake:        bet.Stake,
		Odds:         bet.Odds,
		PotentialWin: bet.PotentialWin,
		Status:       bet.Status,
		CreatedAt:    bet.CreatedAt,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Store in memory
	s.betsMu.Lock()
	s.bets[bet.ID] = bet
	s.betsMu.Unlock()

	return bet, nil
}

// SettleBet settles a bet
func (s *SportsbookService) SettleBet(ctx context.Context, betID uuid.UUID, result string) error {
	s.betsMu.Lock()
	bet, ok := s.bets[betID]
	s.betsMu.Unlock()

	if !ok {
		return fmt.Errorf("bet not found")
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet already settled")
	}

	// Update bet status
	bet.Status = result

	// If won, credit winnings
	if result == "won" {
		var wallet models.Wallet
		if err := s.db.Where("user_id = ?", bet.UserID).First(&wallet).Error; err == nil {
			wallet.Balance += bet.PotentialWin
			wallet.UpdatedAt = time.Now()
			s.db.Save(&wallet)
		}
	}

	// Update database
	now := time.Now()
	return s.db.Model(&models.SportsBet{}).
		Where("id = ?", betID).
		Updates(map[string]interface{}{
			"status":      result,
			"settled_at": now,
		}).Error
}

// GetUserBets returns user's bets
func (s *SportsbookService) GetUserBets(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]models.SportsBet, int64, error) {
	var bets []models.SportsBet
	query := s.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&bets).Error; err != nil {
		return nil, 0, err
	}

	return bets, total, nil
}

// GetBetsByEvent returns all bets for an event
func (s *SportsbookService) GetBetsByEvent(ctx context.Context, eventID string) ([]*Bet, error) {
	s.betsMu.RLock()
	defer s.betsMu.RUnlock()

	var result []*Bet
	for _, bet := range s.bets {
		if bet.EventID == eventID {
			result = append(result, bet)
		}
	}

	return result, nil
}

// GetUpcomingEvents returns upcoming events across all sports
func (s *SportsbookService) GetUpcomingEvents(ctx context.Context, limit int) ([]Event, error) {
	// Get all sports
	sports, err := s.GetSports(ctx)
	if err != nil {
		return nil, err
	}

	var allEvents []Event
	now := time.Now()

	for _, sport := range sports {
		leagues, err := s.GetLeagues(ctx, sport.ID)
		if err != nil {
			continue
		}

		for _, league := range leagues {
			events, err := s.GetEvents(ctx, sport.ID, league.ID, now)
			if err != nil {
				continue
			}

			for _, event := range events {
				if event.Status == "scheduled" && event.StartTime.After(now) {
					event.Metadata = map[string]interface{}{
						"sport":   sport.Name,
						"league":  league.Name,
						"country": league.Country,
					}
					allEvents = append(allEvents, event)
				}
			}
		}
	}

	// Sort by start time
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].StartTime.Before(allEvents[j].StartTime)
	})

	if limit > 0 && len(allEvents) > limit {
		allEvents = allEvents[:limit]
	}

	return allEvents, nil
}

// GetPopularEvents returns popular events for display
func (s *SportsbookService) GetPopularEvents(ctx context.Context) ([]Event, error) {
	return s.GetUpcomingEvents(ctx, 20)
}

// GetEventDetails returns detailed event information
func (s *SportsbookService) GetEventDetails(ctx context.Context, eventID string) (*Event, []Market, error) {
	// Get from database
	var dbEvent models.SportsEvent
	if err := s.db.Where("id = ?", eventID).First(&dbEvent).Error; err != nil {
		return nil, nil, err
	}

	event := Event{
		ID:         dbEvent.ID.String(),
		SportID:    dbEvent.Sport,
		LeagueID:   dbEvent.League,
		HomeTeam:   dbEvent.HomeTeam,
		AwayTeam:   dbEvent.AwayTeam,
		StartTime:  dbEvent.StartTime,
		Status:     dbEvent.Status,
		HomeScore:  dbEvent.HomeScore,
		AwayScore:  dbEvent.AwayScore,
		Period:     dbEvent.Status,
		TimeRemaining: dbEvent.Status,
	}

	markets, err := s.GetMarkets(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}

	return &event, markets, nil
}

// SyncEvents syncs events from data provider
func (s *SportsbookService) SyncEvents(ctx context.Context) error {
	sports, err := s.GetSports(ctx)
	if err != nil {
		return err
	}

	now := time.Now()

	for _, sport := range sports {
		leagues, err := s.GetLeagues(ctx, sport.ID)
		if err != nil {
			continue
		}

		for _, league := range leagues {
			events, err := s.GetEvents(ctx, sport.ID, league.ID, now)
			if err != nil {
				continue
			}

			for _, event := range events {
				var dbEvent models.SportsEvent
				result := s.db.Where("external_id = ?", event.ID).First(&dbEvent)

				eventData := models.SportsEvent{
					ExternalID:   event.ID,
					Sport:        event.SportID,
					League:       event.LeagueID,
					HomeTeam:     event.HomeTeam,
					AwayTeam:     event.AwayTeam,
					StartTime:    event.StartTime,
					Status:       event.Status,
					HomeScore:    event.HomeScore,
					AwayScore:    event.AwayScore,
					UpdatedAt:    time.Now(),
				}

				if result.Error == gorm.ErrRecordNotFound {
					eventData.CreatedAt = time.Now()
					s.db.Create(&eventData)
				} else if result.Error == nil {
					s.db.Save(&eventData)
				}
			}
		}
	}

	return nil
}

// GetBettingStats returns betting statistics
func (s *SportsbookService) GetBettingStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	var totalBets int64
	var pendingBets int64
	var wonBets int64
	var lostBets int64
	var totalStaked float64
	var totalWon float64

	s.db.Model(&models.SportsBet{}).Where("user_id = ?", userID).Count(&totalBets)
	s.db.Model(&models.SportsBet{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&pendingBets)
	s.db.Model(&models.SportsBet{}).Where("user_id = ? AND status = ?", userID, "won").Count(&wonBets)
	s.db.Model(&models.SportsBet{}).Where("user_id = ? AND status = ?", userID, "lost").Count(&lostBets)
	s.db.Model(&models.SportsBet{}).Where("user_id = ?", userID).Select("COALESCE(SUM(stake), 0)").Scan(&totalStaked)
	s.db.Model(&models.SportsBet{}).Where("user_id = ? AND status = ?", userID, "won").Select("COALESCE(SUM(potential_win - stake), 0)").Scan(&totalWon)

	winRate := 0.0
	if wonBets+lostBets > 0 {
		winRate = float64(wonBets) / float64(wonBets+lostBets) * 100
	}

	return map[string]interface{}{
		"total_bets":   totalBets,
		"pending_bets": pendingBets,
		"won_bets":     wonBets,
		"lost_bets":    lostBets,
		"total_staked": totalStaked,
		"total_won":    totalWon,
		"profit":       totalWon - totalStaked,
		"win_rate":     winRate,
	}, nil
}
