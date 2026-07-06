package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tigercasino/backend/internal/config"
	"tigercasino/backend/internal/models"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.VIPLevel{},
		&models.Referral{},
		&models.Transaction{},
		&models.Game{},
		&models.Bet{},
		&models.CrashGameRound{},
		&models.SportsEvent{},
		&models.SportsBet{},
		&models.LeaderboardEntry{},
		&models.FraudAlert{},
		&models.AuditLog{},
		&models.Session{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Seed default VIP levels
	seedVIPLevels(db)

	log.Println("Database migrations completed")
	return nil
}

func seedVIPLevels(db *gorm.DB) {
	var count int64
	db.Model(&models.VIPLevel{}).Count(&count)

	if count == 0 {
		levels := []models.VIPLevel{
			{Level: 0, Name: "Bronze", RakebackPercent: 0, RequiredPoints: 0},
			{Level: 1, Name: "Silver", RakebackPercent: 5, RequiredPoints: 1000},
			{Level: 2, Name: "Gold", RakebackPercent: 10, RequiredPoints: 10000},
			{Level: 3, Name: "Platinum", RakebackPercent: 15, RequiredPoints: 50000},
			{Level: 4, Name: "Diamond", RakebackPercent: 20, RequiredPoints: 200000},
			{Level: 5, Name: "VIP", RakebackPercent: 25, RequiredPoints: 1000000},
		}

		for i := range levels {
			db.Create(&levels[i])
		}
		log.Println("Seeded default VIP levels")
	}
}
