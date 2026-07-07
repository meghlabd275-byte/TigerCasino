package handlers

import (
	"net/http"

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

// Auth Handlers
func (h *Handler) Register(c *gin.Context) {
	var input struct {
		Email        string `json:"email" binding:"required"`
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required"`
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

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
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

// User Handlers
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

// Wallet Handlers
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
	c.JSON(http.StatusOK, gin.H{"address": "0x1234567890abcdef"})
}

func (h *Handler) Withdraw(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal request submitted"})
}

func (h *Handler) GetTransactions(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	txs, err := h.userService.GetTransactions(userID, "", 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txs)
}

// Game Handlers
func (h *Handler) GetGames(c *gin.Context) {
	var games []models.Game
	h.db.Find(&games)
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
		result, err := h.gameService.PlayDice(userID, input.Amount, 50, "over")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	if game.Type == "slots" {
		result, err := h.gameService.PlaySlots(userID, input.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	if game.Type == "roulette" {
		result, err := h.gameService.PlayRoulette(userID, input.Amount, "red", 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	if game.Type == "blackjack" {
		result, err := h.gameService.PlayBlackjack(userID, input.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	if game.Type == "video_poker" {
		result, err := h.gameService.PlayVideoPoker(userID, input.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bet placed", "game": game.Name})
}

func (h *Handler) GetBetHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	bets, err := h.gameService.GetUserBets(userID, "", 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bets)
}

// Admin Handlers
func (h *Handler) GetAdminStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"total_users": 1000, "total_revenue": 50000})
}

func (h *Handler) GetUsers(c *gin.Context) {
	var users []models.User
	h.db.Find(&users)
	c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func (h *Handler) GetAllTransactions(c *gin.Context) {
	var txs []models.Transaction
	h.db.Find(&txs)
	c.JSON(http.StatusOK, txs)
}

func (h *Handler) ApproveTransaction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Transaction approved"})
}

func (h *Handler) RejectTransaction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Transaction rejected"})
}

func (h *Handler) GetAdminGames(c *gin.Context) {
	var games []models.Game
	h.db.Find(&games)
	c.JSON(http.StatusOK, games)
}

func (h *Handler) UpdateGame(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Game updated"})
}

func (h *Handler) GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	h.db.Find(&logs)
	c.JSON(http.StatusOK, logs)
}
