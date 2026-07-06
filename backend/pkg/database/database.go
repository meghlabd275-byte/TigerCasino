package database

import (
	"fmt"

	"github.com/tigercasino/backend/internal/config"
	"github.com/tigercasino/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	// First connect to postgres to create the database if it doesn't exist
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
	)

	defaultDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Create database if not exists
	var count int64
	defaultDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", cfg.DBName).Scan(&count)
	if count == 0 {
		defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	}

	// Close default connection
	sqlDB, _ := defaultDB.DB()
	sqlDB.Close()

	// Connect to the actual database
	dsn = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// SeedGames seeds the default games - 200+ games
func SeedGames(db *gorm.DB) {
	games := []models.Game{
		// SLOTS - Top Providers (50+ games)
		{Name: "Tiger King", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Mega Moolah", Type: "slots", Provider: "Microgaming", RTP: 88.12, MinBet: 0.25, MaxBet: 6.25, IsActive: true},
		{Name: "Book of Dead", Type: "slots", Provider: "Play n GO", RTP: 96.21, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "Starburst", Type: "slots", Provider: "NetEnt", RTP: 96.09, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "Gonzo's Quest", Type: "slots", Provider: "NetEnt", RTP: 95.97, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Wolf Gold", Type: "slots", Provider: "Pragmatic Play", RTP: 96.01, MinBet: 0.25, MaxBet: 125, IsActive: true},
		{Name: "Sweet Bonanza", Type: "slots", Provider: "Pragmatic Play", RTP: 96.48, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Gates of Olympus", Type: "slots", Provider: "Pragmatic Play", RTP: 95.51, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Money Train", Type: "slots", Provider: "Relax Gaming", RTP: 96.15, MinBet: 0.1, MaxBet: 20, IsActive: true},
		{Name: "Big Bass Bonanza", Type: "slots", Provider: "Pragmatic Play", RTP: 96.71, MinBet: 0.1, MaxBet: 125, IsActive: true},
		{Name: "Divine Fortune", Type: "slots", Provider: "NetEnt", RTP: 96.59, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Fruit Party", Type: "slots", Provider: "Pragmatic Play", RTP: 96.47, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "The Dog House", Type: "slots", Provider: "Pragmatic Play", RTP: 96.51, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Wild West Gold", Type: "slots", Provider: "Pragmatic Play", RTP: 96.51, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Joker's Jewels", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Great Rhino", Type: "slots", Provider: "Pragmatic Play", RTP: 96.65, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Caishen's Gold", Type: "slots", Provider: "Pragmatic Play", RTP: 96.08, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Fire 88", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Chilli Heat", Type: "slots", Provider: "Pragmatic Play", RTP: 96.52, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Panda Fortune", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Ancient Egypt", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Hercules Son of Zeus", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Madame Destiny", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Pyramid King", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Dragon Hot", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Aztec Gems", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Egyptian Fortunes", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Gold Rush", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Super Joker", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Master Joker", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Diamond Strike", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "5 Lions", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "7 Piggies", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Aladdin's Treasure", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Asian Gaming", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Aztec King", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Bingo Tennis", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Caishens Gold", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Dancing King", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Dragon Tiger", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Emperor's Gate", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Fortune 888", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Fu Fu Fu", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Golden Ox", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Happy Golden Fish", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Hot to Burn", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Jewel Box", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Jungle Gorilla", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Kingdom of the Sun", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Lady of the Moon", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Lucky Dragon", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Magic Crystals", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Magic Journey", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Mighty Kong", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Money Dome", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Mystic Sea", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Nian Nian You Yu", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Orion", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Panda 88", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Phoenix Forge", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Piggy Bankers", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Piggy Gold", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Power of Thor", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Queen of Gods", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Riche Jungle", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Santa's Wonderland", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Secret of the Temple", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Seven Seven Seven", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Speed Winner", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Triple Diamond", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Ultra Hold and Spin", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Vampire's Charm", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Wild Beach Life", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Wild Gladiators", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Wonderland", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Yummy Bonanza", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Zhao Cai Jin Bao", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Zombie Carnival", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Book of Vikings", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Mysterious Egypt", Type: "slots", Provider: "Pragmatic Play", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Secret of the Stones", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Twin Spin", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Jack and the Beanstalk", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Dead or Alive", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Dead or Alive 2", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Blood Suckers", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Joker Pro", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Wild Wild West", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Steam Tower", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Planet of the Apes", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Narcos", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Street Fighter", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Finn and the Swirly Spin", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Berry Burst", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Turn Your Fortune", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Tiki Fruits", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Reel Rush", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Eggomatic", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Creature from the Black Lagoon", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Magic Mirror", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Jackpot 6000", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Mega Fortune", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Hall of Gods", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Mega Fortune Dreams", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		{Name: "Arabian Nights", Type: "slots", Provider: "NetEnt", RTP: 96.5, MinBet: 0.2, MaxBet: 100, IsActive: true},
		// Dice Games (15+ games)
		{Name: "Classic Dice", Type: "dice", Provider: "TigerCasino", RTP: 99.0, MinBet: 0.01, MaxBet: 1000, IsActive: true},
		{Name: "Plinko", Type: "dice", Provider: "TigerCasino", RTP: 98.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Mines", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "HiLo", Type: "dice", Provider: "TigerCasino", RTP: 96.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Keno", Type: "dice", Provider: "TigerCasino", RTP: 95.0, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "Draft", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Limbo", Type: "dice", Provider: "TigerCasino", RTP: 97.5, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Minefield", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Wheel", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Goal", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Rocket", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Tower", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Dice Duel", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Fast Dice", Type: "dice", Provider: "TigerCasino", RTP: 98.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		{Name: "Lucky Dice", Type: "dice", Provider: "TigerCasino", RTP: 97.0, MinBet: 0.01, MaxBet: 100, IsActive: true},
		// Roulette (15+ games)
		{Name: "European Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "American Roulette", Type: "roulette", Provider: "Evolution", RTP: 94.74, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "French Roulette", Type: "roulette", Provider: "Evolution", RTP: 98.65, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Speed Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Immersive Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Lightning Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 0.1, MaxBet: 2000, IsActive: true},
		{Name: "Auto Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Double Ball Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Gold Vault Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Instant Roulette", Type: "roulette", Provider: "Evolution", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Quantum Roulette", Type: "roulette", Provider: "Playtech", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Age of Gods Roulette", Type: "roulette", Provider: "Playtech", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Diamond Roulette", Type: "roulette", Provider: "Playtech", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Pro Roulette", Type: "roulette", Provider: "Playtech", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Mini Roulette", Type: "roulette", Provider: "Playtech", RTP: 97.3, MinBet: 1, MaxBet: 5000, IsActive: true},
		// Blackjack (15+ games)
		{Name: "Classic Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "VIP Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.5, MinBet: 50, MaxBet: 10000, IsActive: true},
		{Name: "Party Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.4, MinBet: 10, MaxBet: 2500, IsActive: true},
		{Name: "Infinite Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.47, MinBet: 1, MaxBet: 2500, IsActive: true},
		{Name: "Power Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.5, MinBet: 25, MaxBet: 10000, IsActive: true},
		{Name: "Free Bet Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Speed Blackjack", Type: "blackjack", Provider: "Evolution", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Blackjack Party", Type: "blackjack", Provider: "Evolution", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Quantum Blackjack", Type: "blackjack", Provider: "Playtech", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Progressive Blackjack", Type: "blackjack", Provider: "Playtech", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Blackjack Switch", Type: "blackjack", Provider: "Playtech", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Super 21", Type: "blackjack", Provider: "Playtech", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Face Up 21", Type: "blackjack", Provider: "Playtech", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Match Play 21", Type: "blackjack", Provider: "Playtech", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "European Blackjack", Type: "blackjack", Provider: "Betsoft", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		{Name: "Perfect Pairs Blackjack", Type: "blackjack", Provider: "Betsoft", RTP: 99.4, MinBet: 10, MaxBet: 5000, IsActive: true},
		// Baccarat (15+ games)
		{Name: "Classic Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Speed Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 5, MaxBet: 5000, IsActive: true},
		{Name: "Baccarat Squeeze", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "No Commission Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.76, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Golden Wealth Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Lightning Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Peek Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Baccarat Control", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "First Person Baccarat", Type: "baccarat", Provider: "Evolution", RTP: 98.94, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Progressive Baccarat", Type: "baccarat", Provider: "Playtech", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Lucky Bonus Baccarat", Type: "baccarat", Provider: "Playtech", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Baccarat Pro", Type: "baccarat", Provider: "Playtech", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Mini Baccarat", Type: "baccarat", Provider: "Playtech", RTP: 98.94, MinBet: 5, MaxBet: 5000, IsActive: true},
		{Name: "Baccarat Deluxe", Type: "baccarat", Provider: "Betsoft", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		{Name: "Punto Banco", Type: "baccarat", Provider: "Betsoft", RTP: 98.94, MinBet: 10, MaxBet: 10000, IsActive: true},
		// Poker (15+ games)
		{Name: "Texas Hold'em", Type: "poker", Provider: "Evolution", RTP: 97.8, MinBet: 5, MaxBet: 5000, IsActive: true},
		{Name: "Caribbean Stud", Type: "poker", Provider: "Evolution", RTP: 96.3, MinBet: 10, MaxBet: 500, IsActive: true},
		{Name: "Three Card Poker", Type: "poker", Provider: "Evolution", RTP: 96.63, MinBet: 10, MaxBet: 1000, IsActive: true},
		{Name: "Ultimate Texas Hold'em", Type: "poker", Provider: "Evolution", RTP: 97.8, MinBet: 10, MaxBet: 2500, IsActive: true},
		{Name: "Casino Hold'em", Type: "poker", Provider: "Evolution", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Side Bet City", Type: "poker", Provider: "Evolution", RTP: 95.5, MinBet: 1, MaxBet: 1000, IsActive: true},
		{Name: "2 Hand Casino Hold'em", Type: "poker", Provider: "Evolution", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Teen Patti", Type: "poker", Provider: "Evolution", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Aussie Rules", Type: "poker", Provider: "Evolution", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Joker Poker", Type: "poker", Provider: "Betsoft", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Deuces Wild", Type: "poker", Provider: "Betsoft", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Jacks or Better", Type: "poker", Provider: "Betsoft", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Bonus Poker", Type: "poker", Provider: "Betsoft", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Double Double Bonus Poker", Type: "poker", Provider: "Betsoft", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		{Name: "Aces and Eights", Type: "poker", Provider: "Betsoft", RTP: 97.8, MinBet: 5, MaxBet: 500, IsActive: true},
		// Live Shows (15+ games)
		{Name: "Dream Catcher", Type: "show", Provider: "Evolution", RTP: 96.58, MinBet: 0.1, MaxBet: 1000, IsActive: true},
		{Name: "Monopoly Live", Type: "show", Provider: "Evolution", RTP: 96.23, MinBet: 0.1, MaxBet: 1000, IsActive: true},
		{Name: "Crazy Time", Type: "show", Provider: "Evolution", RTP: 95.5, MinBet: 0.1, MaxBet: 10000, IsActive: true},
		{Name: "Lightning Roulette", Type: "show", Provider: "Evolution", RTP: 97.3, MinBet: 0.1, MaxBet: 2000, IsActive: true},
		{Name: "Fan Tan", Type: "show", Provider: "Evolution", RTP: 97.5, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Super Sic Bo", Type: "show", Provider: "Evolution", RTP: 97.22, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Craps Live", Type: "show", Provider: "Evolution", RTP: 97.22, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Dragon Tiger", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Andar Bahar", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Teen Patti Live", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Super Color Dragon", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "War of Bets", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Top Card", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Lucky 7", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		{Name: "Sports Studio", Type: "show", Provider: "Evolution", RTP: 97.0, MinBet: 1, MaxBet: 5000, IsActive: true},
		// Other Games
		{Name: "Aviator", Type: "crash", Provider: "Spribe", RTP: 97.0, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "Spaceman", Type: "crash", Provider: "Pragmatic Play", RTP: 97.0, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "JetX", Type: "crash", Provider: "SmartSoft", RTP: 97.0, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "Balloon", Type: "crash", Provider: "SmartSoft", RTP: 97.0, MinBet: 0.1, MaxBet: 100, IsActive: true},
		{Name: "Zeppelin", Type: "crash", Provider: "Betsoft", RTP: 97.0, MinBet: 0.1, MaxBet: 100, IsActive: true},
	}

	for _, game := range games {
		var existing models.Game
		if err := db.Where("name = ?", game.Name).First(&existing).Error; err != nil {
			db.Create(&game)
		}
	}
}
