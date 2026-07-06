package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tigercasino/backend/internal/config"
	"github.com/tigercasino/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// Register creates a new user account
func (h *Handler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		ID:            uuid.New(),
		Email:         input.Email,
		Username:      input.Username,
		PasswordHash:  string(hashedPassword),
		Balance:       100, // Welcome bonus
		BonusBalance:  50,
		VIPLevel:      0,
		KYCStatus:     "pending",
		IsVerified:    false,
		IsAdmin:       false,
		IsBanned:      false,
		Is2FAEnabled:  false,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// Login authenticates a user
func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is banned
	if user.IsBanned {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is banned"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login
	now := time.Now()
	h.db.Model(&user).Update("last_login", &now)

	// Generate JWT token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Create session
	session := models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	h.db.Create(&session)

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Logout logs out the current user
func (h *Handler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		token = token[7:] // Remove "Bearer "
		h.db.Where("token = ?", token).Delete(&models.Session{})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken refreshes the JWT token
func (h *Handler) RefreshToken(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	token, err := h.generateToken(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SendOTP sends OTP to user email/phone
func (h *Handler) SendOTP(c *gin.Context) {
	var input struct {
		Type  string `json:"type" binding:"required,oneof=email phone"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In production, generate and send actual OTP
	// For demo, we simulate success
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"otp":     "123456", // Demo only
	})
}

// VerifyOTP verifies the OTP
func (h *Handler) VerifyOTP(c *gin.Context) {
	var input struct {
		OTP string `json:"otp" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In production, verify actual OTP
	// For demo, accept any 6-digit code
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

// GetCurrentUser returns the current user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		Username string `json:"username"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}

	h.db.Save(&user)
	c.JSON(http.StatusOK, user)
}

// Setup2FA sets up two-factor authentication
func (h *Handler) Setup2FA(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// In production, generate actual TOTP secret
	secret := "JBSWY3DPEHPK3PXP"
	
	user.TwoFASecret = secret
	user.Is2FAEnabled = true
	h.db.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA enabled successfully",
		"secret": secret,
	})
}

// Verify2FA verifies two-factor authentication
func (h *Handler) Verify2FA(c *gin.Context) {
	var input struct {
		Code string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In production, verify actual TOTP code
	c.JSON(http.StatusOK, gin.H{"message": "2FA verified successfully"})
}

// Disable2FA disables two-factor authentication
func (h *Handler) Disable2FA(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.TwoFASecret = ""
	user.Is2FAEnabled = false
	h.db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// GetBalance returns the user's balance
func (h *Handler) GetBalance(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":      user.Balance,
		"bonusBalance": user.BonusBalance,
	})
}

// GetDepositAddress returns the user's deposit address
func (h *Handler) GetDepositAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	currency := c.Query("currency")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// In production, generate actual deposit address
	address := "0x" + uuid.New().String()[:40]

	c.JSON(http.StatusOK, gin.H{
		"currency": currency,
		"address":  address,
	})
}

// Withdraw processes a withdrawal request
func (h *Handler) Withdraw(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		Amount   float64 `json:"amount" binding:"required,gt=0"`
		Address  string  `json:"address" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check balance
	if user.Balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Check KYC status
	if user.KYCStatus != "verified" {
		c.JSON(http.StatusForbidden, gin.H{"error": "KYC verification required for withdrawals"})
		return
	}

	// Create transaction
	tx := models.Transaction{
		ID:       uuid.New(),
		UserID:   userID.(uuid.UUID),
		Type:     "withdrawal",
		Amount:   input.Amount,
		Currency: input.Currency,
		Status:   "pending",
		Address:  input.Address,
		Fee:      input.Amount * 0.001, // 0.1% fee
		CreatedAt: time.Now(),
	}

	h.db.Create(&tx)

	// Deduct balance
	user.Balance -= input.Amount
	h.db.Save(&user)

	c.JSON(http.StatusCreated, tx)
}

// GetTransactions returns the user's transactions
func (h *Handler) GetTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	var transactions []models.Transaction
	query := h.db.Where("user_id = ?", userID).Order("created_at desc")

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"items":     transactions,
		"total":     total,
		"page":      parseInt(page),
		"pageSize":  parseInt(pageSize),
		"totalPages": (total + int64(parseInt(pageSize)) - 1) / int64(parseInt(pageSize)),
	})
}

// GetGames returns all available games
func (h *Handler) GetGames(c *gin.Context) {
	var games []models.Game
	h.db.Where("is_active = ?", true).Find(&games)

	c.JSON(http.StatusOK, games)
}

// GetGame returns a specific game
func (h *Handler) GetGame(c *gin.Context) {
	id := c.Param("id")

	var game models.Game
	if err := h.db.Where("id = ?", id).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// PlaceBet places a bet
func (h *Handler) PlaceBet(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameID := c.Param("id")

	var input struct {
		Amount  float64               `json:"amount" binding:"required,gt=0"`
		GameData map[string]any       `json:"gameData"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var game models.Game
	if err := h.db.Where("id = ?", gameID).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Check balance
	if user.Balance+user.BonusBalance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Create bet
	bet := models.Bet{
		ID:         uuid.New(),
		UserID:     userID.(uuid.UUID),
		GameID:     uuid.MustParse(gameID),
		BetAmount:  input.Amount,
		WinAmount:  0,
		Multiplier: 0,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	h.db.Create(&bet)

	// Deduct balance
	if user.BonusBalance >= input.Amount {
		user.BonusBalance -= input.Amount
	} else {
		diff := input.Amount - user.BonusBalance
		user.BonusBalance = 0
		user.Balance -= diff
	}
	h.db.Save(&user)

	// In production, call C++ game engine for result
	// For demo, simulate random win/loss
	win := time.Now().UnixNano()%2 == 0
	if win {
		multiplier := 2.0
		bet.Status = "won"
		bet.Multiplier = multiplier
		bet.WinAmount = input.Amount * multiplier
		user.Balance += bet.WinAmount
	} else {
		bet.Status = "lost"
	}

	now := time.Now()
	bet.SettledAt = &now
	h.db.Save(&bet)
	h.db.Save(&user)

	c.JSON(http.StatusCreated, bet)
}

// GetBetHistory returns the user's bet history
func (h *Handler) GetBetHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameID := c.Query("gameId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	var bets []models.Bet
	query := h.db.Where("user_id = ?", userID)

	if gameID != "" {
		query = query.Where("game_id = ?", gameID)
	}

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at desc").Find(&bets)

	c.JSON(http.StatusOK, gin.H{
		"items":      bets,
		"total":      total,
		"page":       parseInt(page),
		"pageSize":   parseInt(pageSize),
		"totalPages": (total + int64(parseInt(pageSize)) - 1) / int64(parseInt(pageSize)),
	})
}

// GetAdminStats returns admin dashboard statistics
func (h *Handler) GetAdminStats(c *gin.Context) {
	var totalUsers int64
	var activeUsers int64
	var totalRevenue float64
	var totalBets int64
	var pendingWithdrawals int64

	h.db.Model(&models.User{}).Count(&totalUsers)
	h.db.Model(&models.User{}).Where("is_banned = ?", false).Count(&activeUsers)
	h.db.Model(&models.Transaction{}).Where("type = ? AND status = ?", "deposit", "completed").Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
	h.db.Model(&models.Bet{}).Count(&totalBets)
	h.db.Model(&models.Transaction{}).Where("type = ? AND status = ?", "withdrawal", "pending").Count(&pendingWithdrawals)

	c.JSON(http.StatusOK, gin.H{
		"totalUsers":         totalUsers,
		"activeUsers":        activeUsers,
		"totalRevenue":       totalRevenue,
		"totalBets":          totalBets,
		"pendingWithdrawals":  pendingWithdrawals,
		"systemHealth":       "healthy",
	})
}

// GetUsers returns all users (admin)
func (h *Handler) GetUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	search := c.Query("search")

	var users []models.User
	query := h.db.Model(&models.User{})

	if search != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at desc").Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"items":      users,
		"total":      total,
		"page":       parseInt(page),
		"pageSize":   parseInt(pageSize),
		"totalPages": (total + int64(parseInt(pageSize)) - 1) / int64(parseInt(pageSize)),
	})
}

// UpdateUser updates a user (admin)
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		IsBanned    *bool   `json:"isBanned"`
		KYCStatus   string  `json:"kycStatus"`
		Balance     float64 `json:"balance"`
		VIPLevel    int     `json:"vipLevel"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.IsBanned != nil {
		user.IsBanned = *input.IsBanned
	}
	if input.KYCStatus != "" {
		user.KYCStatus = input.KYCStatus
	}
	if input.Balance > 0 {
		user.Balance = input.Balance
	}
	if input.VIPLevel > 0 {
		user.VIPLevel = input.VIPLevel
	}

	h.db.Save(&user)

	// Create audit log
	audit := models.AuditLog{
		ID:        uuid.New(),
		UserID:    &user.ID,
		Action:    "user_update",
		Details:   "Admin updated user",
		CreatedAt: time.Now(),
	}
	h.db.Create(&audit)

	c.JSON(http.StatusOK, user)
}

// GetAllTransactions returns all transactions (admin)
func (h *Handler) GetAllTransactions(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	status := c.Query("status")

	var transactions []models.Transaction
	query := h.db.Model(&models.Transaction{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at desc").Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"items":      transactions,
		"total":      total,
		"page":       parseInt(page),
		"pageSize":   parseInt(pageSize),
		"totalPages": (total + int64(parseInt(pageSize)) - 1) / int64(parseInt(pageSize)),
	})
}

// ApproveTransaction approves a transaction
func (h *Handler) ApproveTransaction(c *gin.Context) {
	id := c.Param("id")

	var tx models.Transaction
	if err := h.db.First(&tx, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	tx.Status = "completed"
	now := time.Now()
	tx.ProcessedAt = &now
	h.db.Save(&tx)

	c.JSON(http.StatusOK, tx)
}

// RejectTransaction rejects a transaction
func (h *Handler) RejectTransaction(c *gin.Context) {
	id := c.Param("id")

	var tx models.Transaction
	if err := h.db.First(&tx, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// If withdrawal, refund the user
	if tx.Type == "withdrawal" {
		var user models.User
		if err := h.db.First(&user, tx.UserID).Error; err == nil {
			user.Balance += tx.Amount
			h.db.Save(&user)
		}
	}

	tx.Status = "rejected"
	now := time.Now()
	tx.ProcessedAt = &now
	h.db.Save(&tx)

	c.JSON(http.StatusOK, tx)
}

// GetAdminGames returns all games (admin)
func (h *Handler) GetAdminGames(c *gin.Context) {
	var games []models.Game
	h.db.Find(&games)

	c.JSON(http.StatusOK, games)
}

// UpdateGame updates a game (admin)
func (h *Handler) UpdateGame(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		IsActive *bool   `json:"isActive"`
		RTP      float64 `json:"rtp"`
		MinBet   float64 `json:"minBet"`
		MaxBet   float64 `json:"maxBet"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var game models.Game
	if err := h.db.First(&game, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if input.IsActive != nil {
		game.IsActive = *input.IsActive
	}
	if input.RTP > 0 {
		game.RTP = input.RTP
	}
	if input.MinBet > 0 {
		game.MinBet = input.MinBet
	}
	if input.MaxBet > 0 {
		game.MaxBet = input.MaxBet
	}

	h.db.Save(&game)

	c.JSON(http.StatusOK, game)
}

// GetAuditLogs returns audit logs (admin)
func (h *Handler) GetAuditLogs(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "50")

	var logs []models.AuditLog
	query := h.db.Model(&models.AuditLog{})

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at desc").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"items":      logs,
		"total":      total,
		"page":       parseInt(page),
		"pageSize":   parseInt(pageSize),
		"totalPages": (total + int64(parseInt(pageSize)) - 1) / int64(parseInt(pageSize)),
	})
}

// Helper function to generate JWT token
func (h *Handler) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// Helper function to parse int
func parseInt(s string) int {
	var i int
	for _, c := range s {
		i = i*10 + int(c-'0')
	}
	return i
}
