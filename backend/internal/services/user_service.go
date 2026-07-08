package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tigercasino/backend/internal/config"
	"github.com/tigercasino/backend/internal/models"
)

type AuthService struct {
	db       *gorm.DB
	cfg      *config.Config
	security *SecurityBridge
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:       db,
		cfg:      cfg,
		security: NewSecurityBridge(),
	}
}

func (s *AuthService) Register(email, username, password string, referralCode string) (*models.User, error) {
	hashedPassword := s.security.HashPassword(password)
	if hashedPassword == "" {
		return nil, errors.New("failed to hash password")
	}

	user := models.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		Balance:      1000, // Starting balance for demo
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !s.security.VerifyPassword(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) GetBalance(userID uuid.UUID) (float64, float64, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 0, 0, err
	}
	return user.Balance, user.BonusBalance, nil
}

func (s *UserService) UpdateBalance(userID uuid.UUID, amount float64) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (s *UserService) CheckVIPUpgrade(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Simple logic: upgrade based on total wagered (mocked here as balance for simplicity)
	newLevel := int(user.Balance / 10000)
	if newLevel > user.VIPLevel {
		return s.db.Model(&user).Update("vip_level", newLevel).Error
	}
	return nil
}

func (s *UserService) Setup2FA(userID uuid.UUID) (string, error) {
	return "secret", nil
}

func (s *UserService) Verify2FA(userID uuid.UUID, code string) (bool, error) {
	return true, nil
}

func (s *UserService) Disable2FA(userID uuid.UUID) error {
	return nil
}

func (s *UserService) GetTransactions(userID uuid.UUID, txType string, limit int) ([]models.Transaction, error) {
	var txs []models.Transaction
	return txs, nil
}

func (s *UserService) UpdateWagered(userID uuid.UUID, amount float64) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("total_wagered", gorm.Expr("total_wagered + ?", amount)).Error
}

func (s *UserService) UpdateDeposited(userID uuid.UUID, amount float64) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("total_deposited", gorm.Expr("total_deposited + ?", amount)).Error
}

func (s *UserService) UpdateWithdrawn(userID uuid.UUID, amount float64) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("total_withdrawn", gorm.Expr("total_withdrawn + ?", amount)).Error
}

func (s *UserService) GetVIPLevel(userID uuid.UUID) (*models.VIPLevel, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var vipLevel models.VIPLevel
	err := s.db.Where("level = ?", user.VIPLevel).First(&vipLevel).Error
	if err != nil {
		return nil, err
	}

	return &vipLevel, nil
}

func (s *UserService) CheckAndUpgradeVIP(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Check if user qualifies for a higher VIP level based on total wagered
	var levels []models.VIPLevel
	s.db.Order("level DESC").Find(&levels)

	for _, level := range levels {
		if user.TotalWagered >= level.RequiredPoints && user.VIPLevel < level.Level {
			s.db.Model(&user).Update("vip_level", level.Level)
			break
		}
	}

	return nil
}
