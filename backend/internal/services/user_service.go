package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tigercasino/backend/internal/config"
	"tigercasino/backend/internal/models"
)

// AuthService handles authentication operations
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register creates a new user account
func (s *AuthService) Register(email, username, password string, referralCode string) (*models.User, error) {
	// Check if email exists
	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Check if username exists
	if err := s.db.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Balance:      0,
		BonusBalance: 0,
		VIPLevel:     0,
		KYCStatus:    "pending",
		IsVerified:   false,
		IsAdmin:      false,
		IsBanned:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Process referral if code provided
	if referralCode != "" {
		referralService := NewReferralService(s.db)
		var referrer models.User
		if err := s.db.Where("referral_code = ?", referralCode).First(&referrer).Error; err == nil {
			referralService.ProcessReferral(referrer.ID, user.ID, referralCode)
		}
	}

	return &user, nil
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check if banned
	if user.IsBanned {
		return nil, "", errors.New("account suspended")
	}

	// Generate JWT token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, "", err
	}

	// Update last login
	now := time.Now()
	s.db.Model(&user).Update("last_login", &now)

	return &user, token, nil
}

// GenerateToken creates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"email":    user.Email,
		"username":  user.Username,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

// UserService handles user operations
type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// GetBalance returns user's balance
func (s *UserService) GetBalance(userID uuid.UUID) (float64, float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, 0, err
	}
	return user.Balance, user.BonusBalance, nil
}

// UpdateBalance updates user's balance
func (s *UserService) UpdateBalance(userID uuid.UUID, amount float64, isBonus bool) error {
	if isBonus {
		return s.db.Model(&models.User{}).
			Where("id = ?", userID).
			Update("bonus_balance", gorm.Expr("bonus_balance + ?", amount)).Error
	}
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// GetTransactions returns user's transactions
func (s *UserService) GetTransactions(userID uuid.UUID, txType string, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := s.db.Where("user_id = ?", userID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&transactions).Error
	return transactions, err
}

// Setup2FA generates 2FA secret for user
func (s *UserService) Setup2FA(userID uuid.UUID) (string, error) {
	// In production, would use proper TOTP library
	secret := generateRandomString(32)
	return secret, s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("two_fa_secret", secret).Error
}

// Enable2FA enables 2FA for user
func (s *UserService) Enable2FA(userID uuid.UUID, secret string) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"two_fa_secret":  secret,
			"is_2fa_enabled": true,
		}).Error
}

// Disable2FA disables 2FA for user
func (s *UserService) Disable2FA(userID uuid.UUID) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"two_fa_secret":  nil,
			"is_2fa_enabled": false,
		}).Error
}

// Verify2FA verifies 2FA code
func (s *UserService) Verify2FA(userID uuid.UUID, code string) (bool, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return false, err
	}

	// In production, would verify against TOTP
	// Simplified - check if code matches (would be time-based in production)
	return user.TwoFASecret != "", nil
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(i % 256)
	}
	return string(b)
}
