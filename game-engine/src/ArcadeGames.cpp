// Additional Crash & Arcade Games
#include "ArcadeGames.hpp"
#include <iostream>
#include <random>
#include <sstream>
#include <cmath>

namespace TigerCasino {

// SpaceXY Crash Game
SpaceXYGame::SpaceXYGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double SpaceXYGame::calculateCrashPoint(uint64_t hash) {
    // Similar to crash but with space theme
    double r = (double)(hash % 10000) / 10000.0;
    double point = 1.0 / (1.0 - r * 0.95);
    if (point > 1000.0) point = 1000.0;
    if (point < 1.0) point = 1.0;
    return point;
}

double SpaceXYGame::play(double bet, double cashoutAt) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    double crashPoint = calculateCrashPoint(hash);
    
    if (cashoutAt <= crashPoint) {
        return bet * cashoutAt;
    }
    return 0;
}

// Zeppelin Crash Game
ZeppelinGame::ZeppelinGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double ZeppelinGame::calculateCrashPoint(uint64_t hash) {
    double r = (double)(hash % 10000) / 10000.0;
    // Different curve for zeppelin
    double point = 1.0 + (exp(r * 3.5) - 1) * 10;
    if (point > 1000.0) point = 1000.0;
    return point;
}

double ZeppelinGame::play(double bet, double target) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    double crashPoint = calculateCrashPoint(hash);
    
    if (target <= crashPoint) {
        return bet * target;
    }
    return 0;
}

// Balloon Game
BalloonGame::BalloonGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

bool BalloonGame::inflate(double probability) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    uint64_t check = hash % 100;
    return (double)check < (probability * 100);
}

double BalloonGame::play(double bet, int pumps) {
    double multiplier = 1.0;
    for (int i = 0; i < pumps; i++) {
        multiplier += 0.1; // 10% per pump
        
        if (!inflate(0.95 - (i * 0.02))) {
            return 0; // Popped
        }
    }
    return bet * multiplier;
}

// Penalty Shootout
PenaltyShootoutGame::PenaltyShootoutGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double PenaltyShootoutGame::play(double bet, const std::string& prediction) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    // 70% chance to score
    bool scored = (hash % 100) < 70;
    
    if (prediction == "goal" && scored) return bet * 1.8;
    if (prediction == "miss" && !scored) return bet * 3.0;
    return 0;
}

// Dino Run
DinoRunGame::DinoRunGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double DinoRunGame::play(double bet, double distance) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    uint64_t random = hash % 100;
    double multiplier = 1.0 + (distance * 0.01);
    
    if (random < 90) { // 90% success rate
        return bet * multiplier;
    }
    return 0;
}

// Minefield
MinefieldGame::MinefieldGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

bool MinefieldGame::revealTile(int tile, int numMines) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed + std::to_string(tile));
    
    int totalTiles = 25;
    int safeTiles = totalTiles - numMines;
    
    // Probability decreases as more tiles revealed
    uint64_t check = hash % 100;
    double probability = (double)safeTiles / (double)totalTiles * 100;
    
    return (double)check < probability;
}

double MinefieldGame::play(double bet, int numMines) {
    double multiplier = 1.0;
    int wins = 0;
    
    for (int i = 0; i < 25 - numMines; i++) {
        if (revealTile(i, numMines)) {
            wins++;
            multiplier *= 1.1;
        } else {
            return 0; // Hit mine
        }
    }
    
    return bet * multiplier;
}

// Tower Crash
TowerCrashGame::TowerCrashGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double TowerCrashGame::play(double bet, int floorsCleared, int totalFloors) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    // Each floor has 90% success rate
    for (int i = 0; i < floorsCleared; i++) {
        uint64_t check = (hash + i) % 100;
        if (check < 10) return 0; // Failed
    }
    
    double multiplier = 1.0 + (floorsCleared * 0.25);
    return bet * multiplier;
}

// Fruit Ninja
FruitNinjaGame::FruitNinjaGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double FruitNinjaGame::play(double bet, int fruits) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    double multiplier = 1.0;
    int bombs = 0;
    
    for (int i = 0; i < fruits; i++) {
        uint64_t check = (hash + i) % 100;
        
        if (check < 5) { // 5% bomb chance
            bombs++;
        } else {
            multiplier += 0.2;
        }
    }
    
    if (bombs > 0) return 0;
    return bet * multiplier;
}

// Racing
RacingGame::RacingGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

int RacingGame::getWinner() {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    return hash % 6; // 6 horses/racers
}

double RacingGame::play(double bet, int selection) {
    int winner = getWinner();
    
    if (selection == winner) {
        return bet * 5.0; // 5x payout
    }
    return 0;
}

} // namespace TigerCasino
