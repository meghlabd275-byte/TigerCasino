package observability

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tigercasino_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tigercasino_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Game metrics
	gameBetsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tigercasino_game_bets_total",
			Help: "Total number of game bets",
		},
		[]string{"game_id", "provider"},
	)

	gameWinsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tigercasino_game_wins_total",
			Help: "Total amount won by players",
		},
		[]string{"game_id", "provider"},
	)

	activePlayers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "tigercasino_active_players",
			Help: "Number of currently active players",
		},
	)

	activeGames = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tigercasino_active_games",
			Help: "Number of active game sessions",
		},
		[]string{"game_id"},
	)

	// Crypto metrics
	cryptoDepositsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tigercasino_crypto_deposits_total",
			Help: "Total crypto deposits",
		},
		[]string{"currency", "network", "status"},
	)

	cryptoWithdrawalsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tigercasino_crypto_withdrawals_total",
			Help: "Total crypto withdrawals",
		},
		[]string{"currency", "network", "status"},
	)

	// Database metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tigercasino_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "tigercasino_db_connections_active",
			Help: "Number of active database connections",
		},
	)

	// Cache metrics
	cacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "tigercasino_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "tigercasino_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// WebSocket metrics
	wsConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "tigercasino_ws_connections_active",
			Help: "Number of active WebSocket connections",
		},
	)

	wsMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tigercasino_ws_messages_total",
			Help: "Total WebSocket messages",
		},
		[]string{"direction", "message_type"},
	)
)

// RecordHTTPRequest records an HTTP request
func RecordHTTPRequest(method, endpoint string, status int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordGameBet records a game bet
func RecordGameBet(gameID, provider string, amount float64) {
	gameBetsTotal.WithLabelValues(gameID, provider).Inc()
}

// RecordGameWin records a game win
func RecordGameWin(gameID, provider string, amount float64) {
	gameWinsTotal.WithLabelValues(gameID, provider).Add(amount)
}

// SetActivePlayers sets the number of active players
func SetActivePlayers(count int) {
	activePlayers.Set(float64(count))
}

// SetActiveGames sets the number of active games for a game type
func SetActiveGames(gameID string, count int) {
	activeGames.WithLabelValues(gameID).Set(float64(count))
}

// RecordCryptoDeposit records a crypto deposit
func RecordCryptoDeposit(currency, network, status string, amount float64) {
	cryptoDepositsTotal.WithLabelValues(currency, network, status).Add(amount)
}

// RecordCryptoWithdrawal records a crypto withdrawal
func RecordCryptoWithdrawal(currency, network, status string, amount float64) {
	cryptoWithdrawalsTotal.WithLabelValues(currency, network, status).Add(amount)
}

// RecordDBQuery records a database query
func RecordDBQuery(queryType string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(queryType).Observe(duration.Seconds())
}

// SetDBConnections sets the number of active database connections
func SetDBConnections(count int) {
	dbConnectionsActive.Set(float64(count))
}

// RecordCacheHit records a cache hit
func RecordCacheHit() {
	cacheHits.Inc()
}

// RecordCacheMiss records a cache miss
func RecordCacheMiss() {
	cacheMisses.Inc()
}

// SetWSConnections sets the number of active WebSocket connections
func SetWSConnections(count int) {
	wsConnectionsActive.Set(float64(count))
}

// RecordWSMessage records a WebSocket message
func RecordWSMessage(direction, messageType string) {
	wsMessagesTotal.WithLabelValues(direction, messageType).Inc()
}

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request is allowed
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Clean old requests
	requests := r.requests[key]
	var valid []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= r.limit {
		r.requests[key] = valid
		return false
	}

	r.requests[key] = append(valid, now)
	return true
}

// Reset clears rate limit for a key
func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.requests, key)
}
