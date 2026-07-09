package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/services"
)

// VIPHandler handles VIP and loyalty endpoints
type VIPHandler struct {
	db        *gorm.DB
	vipService *services.VIPService
}

func NewVIPHandler(db *gorm.DB, vipService *services.VIPService) *VIPHandler {
	return &VIPHandler{
		db:         db,
		vipService: vipService,
	}
}

// GetVIPStatus returns current user's VIP status
func (h *VIPHandler) GetVIPStatus(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status, err := h.vipService.GetUserVIPStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// ClaimRakeback allows user to claim their rakeback
func (h *VIPHandler) ClaimRakeback(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	amount, err := h.vipService.ClaimRakeback(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"amount": amount})
}

// RedeemPoints allows user to redeem loyalty points
func (h *VIPHandler) RedeemPoints(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input struct {
		Points int64 `json:"points" binding:"required,min=100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := h.vipService.RedeemPoints(c.Request.Context(), userID, input.Points)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"amount": amount})
}

// ClaimWelcomeBonus claims welcome bonus for new users
func (h *VIPHandler) ClaimWelcomeBonus(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.vipService.ClaimWelcomeBonus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ClaimDepositBonus claims deposit bonus
func (h *VIPHandler) ClaimDepositBonus(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input struct {
		Amount float64 `json:"amount" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.vipService.ClaimDepositBonus(c.Request.Context(), userID, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLeaderboard returns the leaderboard
func (h *VIPHandler) GetLeaderboard(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	entries, err := h.vipService.GetLeaderboard(c.Request.Context(), period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// GetPromotions returns active promotions
func (h *VIPHandler) GetPromotions(c *gin.Context) {
	promotions, err := h.vipService.GetActivePromotions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, promotions)
}

// GetAllVIPLevels returns all VIP levels
func (h *VIPHandler) GetAllVIPLevels(c *gin.Context) {
	levels := h.vipService.GetAllVIPLevels(c.Request.Context())
	c.JSON(http.StatusOK, levels)
}

// TournamentHandler handles tournament endpoints
type TournamentHandler struct {
	db              *gorm.DB
	tournamentService *services.TournamentService
}

func NewTournamentHandler(db *gorm.DB, ts *services.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		db:               db,
		tournamentService: ts,
	}
}

// GetTournaments returns all active tournaments
func (h *TournamentHandler) GetTournaments(c *gin.Context) {
	tournaments, err := h.tournamentService.GetActiveTournaments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// GetUpcomingTournaments returns upcoming tournaments
func (h *TournamentHandler) GetUpcomingTournaments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	tournaments, err := h.tournamentService.GetUpcomingTournaments(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// GetTournamentDetails returns tournament details
func (h *TournamentHandler) GetTournamentDetails(c *gin.Context) {
	tournamentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tournament ID"})
		return
	}

	tournament, leaderboard, err := h.tournamentService.GetTournamentDetails(c.Request.Context(), tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tournament": tournament,
		"leaderboard": leaderboard,
	})
}

// RegisterTournament registers user for a tournament
func (h *TournamentHandler) RegisterTournament(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tournamentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tournament ID"})
		return
	}

	if err := h.tournamentService.RegisterUser(c.Request.Context(), tournamentID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetTournamentLeaderboard returns tournament leaderboard
func (h *TournamentHandler) GetTournamentLeaderboard(c *gin.Context) {
	tournamentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tournament ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	leaderboard, err := h.tournamentService.GetLeaderboard(c.Request.Context(), tournamentID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetUserTournaments returns tournaments user is participating in
func (h *TournamentHandler) GetUserTournaments(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tournaments, err := h.tournamentService.GetUserTournaments(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// GetTournamentResults returns user's tournament results
func (h *TournamentHandler) GetTournamentResults(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	results, err := h.tournamentService.GetUserTournamentResults(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// SportsbookHandler handles sportsbook endpoints
type SportsbookHandler struct {
	db              *gorm.DB
	sportsbookService *services.SportsbookService
}

func NewSportsbookHandler(db *gorm.DB, ss *services.SportsbookService) *SportsbookHandler {
	return &SportsbookHandler{
		db:               db,
		sportsbookService: ss,
	}
}

// GetSports returns all available sports
func (h *SportsbookHandler) GetSports(c *gin.Context) {
	sports, err := h.sportsbookService.GetSports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sports)
}

// GetLeagues returns leagues for a sport
func (h *SportsbookHandler) GetLeagues(c *gin.Context) {
	sportID := c.Query("sportId")
	if sportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sportId required"})
		return
	}

	leagues, err := h.sportsbookService.GetLeagues(c.Request.Context(), sportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leagues)
}

// GetEvents returns events for a sport/league
func (h *SportsbookHandler) GetEvents(c *gin.Context) {
	sportID := c.Query("sportId")
	leagueID := c.Query("leagueId")
	dateStr := c.Query("date")

	var date interface{}
	if dateStr != "" {
		// Parse date
		date = dateStr
	}

	events, err := h.sportsbookService.GetEvents(c.Request.Context(), sportID, leagueID, date.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetLiveEvents returns live events
func (h *SportsbookHandler) GetLiveEvents(c *gin.Context) {
	sportID := c.Query("sportId")

	events, err := h.sportsbookService.GetLiveEvents(c.Request.Context(), sportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEventDetails returns event details with markets
func (h *SportsbookHandler) GetEventDetails(c *gin.Context) {
	eventID := c.Param("id")

	event, markets, err := h.sportsbookService.GetEventDetails(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": event,
		"markets": markets,
	})
}

// PlaceBet places a sports bet
func (h *SportsbookHandler) PlaceBet(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input struct {
		EventID   string  `json:"event_id" binding:"required"`
		MarketID  string  `json:"market_id" binding:"required"`
		Selection string  `json:"selection" binding:"required"`
		Stake     float64 `json:"stake" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bet, err := h.sportsbookService.PlaceBet(c.Request.Context(), userID, input.EventID, input.MarketID, input.Selection, input.Stake)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"betId":        bet.ID,
		"potentialWin":  bet.PotentialWin,
	})
}

// GetUserBets returns user's sports bets
func (h *SportsbookHandler) GetUserBets(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	bets, total, err := h.sportsbookService.GetUserBets(c.Request.Context(), userID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bets": bets,
		"total": total,
		"page":  page,
	})
}

// GetBettingStats returns user's betting statistics
func (h *SportsbookHandler) GetBettingStats(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.sportsbookService.GetBettingStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Helper function to get user ID from context
func (h *Handler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, fmt.Errorf("user not authenticated")
	}

	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID")
	}

	return userID, nil
}

// Import fmt for error handling
var fmt *struct{} = func() *struct{} { return nil }()
