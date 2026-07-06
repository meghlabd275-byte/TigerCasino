package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdmin creates the default admin user
func SeedAdmin(db *gorm.DB, email, password string) {
	var admin User
	if err := db.Where("email = ?", email).First(&admin).Error; err != nil {
		// Admin doesn't exist, create it
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		admin = User{
			ID:           uuid.New(),
			Email:        email,
			Username:     "admin",
			PasswordHash: string(hashedPassword),
			Balance:      0,
			BonusBalance: 0,
			VIPLevel:     0,
			KYCStatus:    "verified",
			IsVerified:   true,
			IsAdmin:      true,
			IsBanned:     false,
			Is2FAEnabled: false,
			EmailVerified: true,
			PhoneVerified: false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		db.Create(&admin)
	}
}

// SeedGames seeds the default games (moved to database package)
func SeedGames(db *gorm.DB) {
	// Games are seeded in database package
}
