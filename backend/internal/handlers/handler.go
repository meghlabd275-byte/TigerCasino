package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/config"
	"github.com/tigercasino/backend/internal/models"
	"github.com/tigercasino/backend/internal/services"
)

type Handler struct {
	db          *gorm.DB
	cfg         *config.Config
	authService *services.AuthService
	userService *services.UserService
	gameService *services.GameService
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:          db,
		cfg:         cfg,
		authService: services.NewAuthService(db, cfg),
		userService: services.NewUserService(db),
		gameService: services.NewGameService(db),
	}
}

// ============ Auth Handlers ============

func (h *Handler) Register(c *gin.Context) {
	var input struct {
		Email        string `json:"email" binding:"required,email"`
		Username     string `json:"username" binding:"required,min=3,max=20"`
		Password     string `json:"password" binding:"required,min=6"`
		ReferralCode string `json:"referral_code"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(input.Email, input.Username, input.Password, input.ReferralCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Initialize provably fair seeds for user
	h.gameService.InitializeUserSeeds(user.ID)

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "new_token"})
}

func (h *Handler) SendOTP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
}

// ============ User Handlers ============

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove protected fields
	delete(updates, "is_admin")
	delete(updates, "balance")
	delete(updates, "password_hash")

	if err := h.userService.UpdateProfile(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func (h *Handler) Setup2FA(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	secret, err := h.userService.Setup2FA(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": secret})
}

func (h *Handler) Verify2FA(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := h.userService.Verify2FA(userID, input.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA verified"})
}

func (h *Handler) Disable2FA(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.userService.Disable2FA(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}

// ============ Wallet Handlers ============

func (h *Handler) GetBalance(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	balance, bonus, err := h.userService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance, "bonus_balance": bonus})
}

func (h *Handler) GetDepositAddress(c *gin.Context) {
	// In production, this would generate a real deposit address
	c.JSON(http.StatusOK, gin.H{"address": "0x1234567890abcdef"})
}

func (h *Handler) Withdraw(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Address string  `json:"address" binding:"required"`
		Amount  float64 `json:"amount" binding:"required,gt=0"`
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create withdrawal request
	tx := models.Transaction{
		ID:      uuid.New(),
		UserID:  userID,
		Type:    "withdrawal",
		Amount:  input.Amount,
		Currency: input.Currency,
		Address: input.Address,
		Status:  "pending",
	}

	if err := h.db.Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal request submitted", "transaction_id": tx.ID})
}

func (h *Handler) GetTransactions(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	txType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	txs, err := h.userService.GetTransactions(userID, txType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txs)
}

// ============ Game Handlers ============

func (h *Handler) GetGames(c *gin.Context) {
	var games []models.Game
	query := h.db.Where("is_active = ?", true)

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if gameType := c.Query("type"); gameType != "" {
		query = query.Where("type = ?", gameType)
	}
	if provider := c.Query("provider"); provider != "" {
		query = query.Where("provider = ?", provider)
	}

	query.Find(&games)
	c.JSON(http.StatusOK, games)
}

func (h *Handler) GetGame(c *gin.Context) {
	id := c.Param("id")
	var game models.Game
	if err := h.db.First(&game, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	c.JSON(http.StatusOK, game)
}

func (h *Handler) PlaceBet(c *gin.Context) {
	gameID := c.Param("id")
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var game models.Game
	h.db.First(&game, "id = ?", gameID)

	if game.Type == "dice" {
		result, err := h.gameService.PlayDice(userID, input.Amount, 50, "over", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	if game.Type == "slots" {
		result, err := h.gameService.PlaySlots(userID, input.Amount, 20)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bet placed", "game": game.Name})
}

// ============ Dice Game ============

func (h *Handler) PlayDice(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount    float64 `json:"amount" binding:"required"`
		Target    float64 `json:"target" binding:"required"`
		Direction string  `json:"direction" binding:"required,oneof=over under"`
		ClientSeed string `json:"client_seed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.PlayDice(userID, input.Amount, input.Target, input.Direction, input.ClientSeed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Slots Game ============

func (h *Handler) PlaySlots(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
		Lines  int     `json:"lines"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Lines <= 0 {
		input.Lines = 20
	}

	result, err := h.gameService.PlaySlots(userID, input.Amount, input.Lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Roulette Game ============

func (h *Handler) PlayRoulette(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount   float64       `json:"amount" binding:"required"`
		BetType string        `json:"bet_type" binding:"required"`
		BetValue interface{}  `json:"bet_value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.PlayRoulette(userID, input.Amount, input.BetType, input.BetValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Blackjack Game ============

func (h *Handler) PlayBlackjack(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.PlayBlackjack(userID, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Crash Game ============

func (h *Handler) GetCrashState(c *gin.Context) {
	state := h.gameService.GetCurrentCrashState()
	c.JSON(http.StatusOK, state)
}

func (h *Handler) GetCrashHistory(c *gin.Context) {
	history := h.gameService.GetCrashHistory()
	c.JSON(http.StatusOK, gin.H{"history": history})
}

func (h *Handler) PlaceCrashBet(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		RoundID    string  `json:"round_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		AutoCashout float64 `json:"auto_cashout"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	h.db.First(&user, userID)

	err := h.gameService.PlaceCrashBet(userID, user.Username, input.RoundID, input.Amount, input.AutoCashout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bet placed"})
}

func (h *Handler) CashoutCrash(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		RoundID   string  `json:"round_id" binding:"required"`
		Multiplier float64 `json:"multiplier" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.CashoutCrash(userID, input.RoundID, input.Multiplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Mines Game ============

func (h *Handler) StartMines(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount      float64 `json:"amount" binding:"required"`
		MinesCount int     `json:"mines_count" binding:"required,min=1,max=24"`
		GridSize   int     `json:"grid_size"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.GridSize <= 0 {
		input.GridSize = 5
	}

	result, err := h.gameService.StartMinesGame(userID, input.Amount, input.MinesCount, input.GridSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) RevealMinesTile(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		GameID string `json:"game_id" binding:"required"`
		Tile   int    `json:"tile" binding:"required,min=0,max=24"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameID, _ := uuid.Parse(input.GameID)
	result, err := h.gameService.RevealMinesTile(userID, gameID.String(), input.Tile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CashoutMines(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		GameID string  `json:"game_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.CashoutMines(userID, input.GameID, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Plinko Game ============

func (h *Handler) PlayPlinko(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
		Rows   int     `json:"rows"`
		Risk   string  `json:"risk"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Rows <= 0 {
		input.Rows = 8
	}
	if input.Risk == "" {
		input.Risk = "medium"
	}

	result, err := h.gameService.PlayPlinko(userID, input.Amount, input.Rows, input.Risk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Limbo Game ============

func (h *Handler) PlayLimbo(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount            float64 `json:"amount" binding:"required"`
		TargetMultiplier float64 `json:"target_multiplier" binding:"required,gt=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.PlayLimbo(userID, input.Amount, input.TargetMultiplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Hi-Lo Game ============

func (h *Handler) StartHilo(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.StartHiloGame(userID, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) PlayHiloChoice(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Choice string `json:"choice" binding:"required,oneof=higher lower equal"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.PlayHiloChoice(userID, input.Choice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CashoutHilo(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	result, err := h.gameService.CashoutHilo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============ Provably Fair ============

func (h *Handler) GetSeeds(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	seeds, err := h.gameService.GetUserSeeds(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No seeds found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"server_seed_hash": seeds.ServerSeedHash,
		"client_seed":      seeds.ClientSeed,
		"nonce":            seeds.NextRevealNonce,
	})
}

func (h *Handler) RegenerateSeeds(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		ClientSeed string `json:"client_seed"`
	}

	c.ShouldBindJSON(&input)

	err := h.gameService.RegenerateSeeds(userID, input.ClientSeed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seeds regenerated"})
}

// ============ Bet History ============

func (h *Handler) GetBetHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	gameType := c.Query("game_type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	bets, err := h.gameService.GetUserBets(userID, gameType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bets)
}

// ============ Leaderboard ============

func (h *Handler) GetLeaderboard(c *gin.Context) {
	gameType := c.Query("game_type")
	timeframe := c.DefaultQuery("timeframe", "daily")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	entries, err := h.gameService.GetLeaderboard(gameType, timeframe, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// ============ User Stats ============

func (h *Handler) GetUserStats(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	stats, err := h.gameService.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ============ Admin Handlers ============

func (h *Handler) GetAdminStats(c *gin.Context) {
	var totalUsers int64
	var totalBets int64
	var totalRevenue float64

	h.db.Model(&models.User{}).Count(&totalUsers)
	h.db.Model(&models.Bet{}).Count(&totalBets)

	h.db.Model(&models.Transaction{}).Where("type = ?", "deposit").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_users":   totalUsers,
		"total_bets":    totalBets,
		"total_revenue": totalRevenue,
	})
}

func (h *Handler) GetUsers(c *gin.Context) {
	var users []models.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	query := h.db.Model(&models.User{})
	if search != "" {
		query = query.Where("email LIKE ? OR username LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Offset((page - 1) * limit).Limit(limit).Find(&users)
	c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove protected fields
	delete(updates, "id")
	delete(updates, "password_hash")

	h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func (h *Handler) GetAllTransactions(c *gin.Context) {
	var txs []models.Transaction
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	query := h.db.Model(&models.Transaction{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&txs)
	c.JSON(http.StatusOK, txs)
}

func (h *Handler) ApproveTransaction(c *gin.Context) {
	txID := c.Param("id")

	var tx models.Transaction
	if err := h.db.First(&tx, "id = ?", txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	tx.Status = "confirmed"
	h.db.Save(&tx)

	// Update user balance for deposits
	if tx.Type == "deposit" {
		h.userService.UpdateBalance(tx.UserID, tx.Amount)
		h.userService.UpdateDeposited(tx.UserID, tx.Amount)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction approved"})
}

func (h *Handler) RejectTransaction(c *gin.Context) {
	txID := c.Param("id")

	var tx models.Transaction
	if err := h.db.First(&tx, "id = ?", txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	tx.Status = "rejected"
	h.db.Save(&tx)

	// Refund user for withdrawals
	if tx.Type == "withdrawal" {
		h.userService.UpdateBalance(tx.UserID, tx.Amount)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction rejected"})
}

func (h *Handler) GetAdminGames(c *gin.Context) {
	var games []models.Game
	h.db.Find(&games)
	c.JSON(http.StatusOK, games)
}

func (h *Handler) UpdateGame(c *gin.Context) {
	gameID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(updates, "id")

	h.db.Model(&models.Game{}).Where("id = ?", gameID).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "Game updated"})
}

func (h *Handler) GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	h.db.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&logs)
	c.JSON(http.StatusOK, logs)
}

// Admin: Create a new game
func (h *Handler) CreateGame(c *gin.Context) {
	var game models.Game
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game.ID = uuid.New()
	if err := h.db.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// Admin: Update VIP level for user
func (h *Handler) UpdateUserVIP(c *gin.Context) {
	userID := c.Param("id")

	var input struct {
		VIPLevel int `json:"vip_level"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Model(&models.User{}).Where("id = ?", userID).Update("vip_level", input.VIPLevel)
	c.JSON(http.StatusOK, gin.H{"message": "VIP level updated"})
}

// ============ Additional Game Handlers ============

func (h *Handler) PlayVideoPoker(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.gameService.PlayVideoPoker(userID, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Helper function to get user ID from token (simplified)
func getUserIDFromToken(c *gin.Context) (uuid.UUID, error) {
	// This would normally extract the user ID from the JWT token
	// For now, return a placeholder
	return uuid.New(), nil
}
