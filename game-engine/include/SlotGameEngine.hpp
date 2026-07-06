#ifndef SLOT_GAME_ENGINE_HPP
#define SLOT_GAME_ENGINE_HPP

#include <string>
#include <vector>
#include <map>
#include <memory>
#include <random>
#include <chrono>
#include <cmath>
#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

/**
 * Slot Game Engine - Ultra Low Latency
 * Supports multiple providers and game types
 */
class SlotGameEngine {
public:
    // Slot game configuration
    struct SlotConfig {
        std::string gameId;
        std::string provider;
        std::string name;
        uint8_t reels;
        uint8_t rows;
        uint8_t paylines;
        double minBet;
        double maxBet;
        double rtp;           // Return to player percentage
        double volatility;     // Low, Medium, High
        bool hasBonus;
        bool hasJackpot;
        std::vector<std::string> symbols;
        std::map<std::string, std::vector<std::string>> paytable; // symbol -> wins
    };

    // Slot game result
    struct SlotResult {
        std::string gameId;
        uint64_t roundId;
        double betAmount;
        std::vector<std::vector<std::string>> reels;  // 2D grid of symbols
        std::vector<std::string> winningLines;
        double totalWin;
        double multiplier;
        bool isBonus;
        bool isJackpot;
        std::string bonusType;
        ServerSeedInfo serverSeed;
    };

    struct ServerSeedInfo {
        std::string seed;
        std::string hash;
        std::string clientSeed;
        uint64_t nonce;
    };

    // Progressive jackpot
    struct Jackpot {
        std::string name;
        double currentAmount;
        double minWin;
        double maxWin;
        double contributionPercent;
    };

private:
    std::map<std::string, SlotConfig> games_;
    std::map<std::string, Jackpot> jackpots_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t roundCounter_{0};

    // Provider games (pre-configured)
    void initProviderGames();

public:
    SlotGameEngine() : rng_(std::make_unique<RandomNumberGenerator>()) {
        initProviderGames();
    }

    // Get all available games
    std::vector<SlotConfig> getAllGames() const;

    // Get games by provider
    std::vector<SlotConfig> getGamesByProvider(const std::string& provider) const;

    // Get game by ID
    const SlotConfig* getGame(const std::string& gameId) const;

    // Spin the reels
    SlotResult spin(const std::string& gameId, double betAmount, 
                   const std::string& clientSeed = "");

    // Calculate winnings from symbol matches
    double calculateWin(const std::string& gameId, 
                      const std::vector<std::vector<std::string>>& reels,
                      double betAmount);

    // Check for bonus round trigger
    bool checkBonusTrigger(const std::vector<std::vector<std::string>>& reels);

    // Progressive jackpot
    void contributeJackpot(const std::string& jackpotName, double amount);
    double getJackpotAmount(const std::string& jackpotName) const;
    bool triggerJackpot(const std::string& gameId);

    // Game providers
    static std::vector<std::string> getSupportedProviders() {
        return {"PragmaticPlay", "NetEnt", "PlayNGO", "BGaming", "Microgaming", "InHouse"};
    }

private:
    // Generate reel positions
    std::vector<std::vector<std::string>> generateReels(const SlotConfig& config,
                                                         const std::string& serverSeed,
                                                         const std::string& clientSeed,
                                                         uint64_t nonce);

    // Check for winning combinations
    std::vector<std::string> findWins(const SlotConfig& config,
                                      const std::vector<std::vector<std::string>>& reels,
                                      double betAmount);
};

void SlotGameEngine::initProviderGames() {
    // Pragmatic Play Games
    games_["pp-sweet-bonanza"] = {
        "pp-sweet-bonanza", "PragmaticPlay", "Sweet Bonanza",
        6, 5, 0, 0.20, 100.00, 96.48, 4.0, true, true,
        {"🍬", "🍇", "🍊", "🍋", "🍌", "💎", "❤️"},
        {}
    };

    games_["pp-gates-olympus"] = {
        "pp-gates-olympus", "PragmaticPlay", "Gates of Olympus",
        6, 5, 0, 0.20, 100.00, 95.51, 5.0, true, false,
        {"⚡", "👑", "💍", "💎", "🔮", "💰", "⭐"},
        {}
    };

    games_["pp-starlight"] = {
        "pp-starlight", "PragmaticPlay", "Starlight Princess",
        6, 5, 0, 0.20, 100.00, 96.50, 5.0, true, true,
        {"👸", "⭐", "🌙", "💎", "💍", "🔮", "❤️"},
        {}
    };

    games_["pp-tiger-fortune"] = {
        "pp-tiger-fortune", "PragmaticPlay", "Tiger Fortune",
        5, 3, 25, 0.25, 125.00, 96.50, 3.5, true, true,
        {"🐯", "💎", "🔔", "⭐", "🎴", "🎋", "🏮"},
        {}
    };

    // NetEnt Games
    games_["net-starburst"] = {
        "net-starburst", "NetEnt", "Starburst",
        10, 3, 10, 0.10, 100.00, 96.09, 3.0, false, false,
        {"💎", "⭐", "🔷", "🔴", "7️⃣", "🎰", "🍇"},
        {}
    };

    games_["net-gonzo-quest"] = {
        "net-gonzo-quest", "NetEnt", "Gonzo's Quest",
        5, 3, 20, 0.20, 100.00, 95.97, 4.0, true, false,
        {"🏔️", "👤", "📜", "🗿", "🐒", "💎", "❓"},
        {}
    };

    // BGaming Games (Crypto-focused)
    games_["bg-plinko"] = {
        "bg-plinko", "BGaming", "Plinko",
        8, 8, 0, 0.10, 100.00, 98.00, 2.0, true, false,
        {"🟢", "🔵", "🔴", "🟡", "⚪"},
        {}
    };

    games_["bg-mines"] = {
        "bg-mines", "BGaming", "Mines",
        5, 5, 0, 0.10, 100.00, 97.00, 3.0, true, false,
        {"💎", "💣", "⭐", "🎯", "🏆"},
        {}
    };

    // In-House Originals
    games_["original-tiger-slots"] = {
        "original-tiger-slots", "InHouse", "Tiger Slots",
        5, 3, 25, 0.10, 100.00, 97.50, 3.0, true, true,
        {"🐯", "💰", "💎", "⭐", "🔔", "🎴", "7️⃣"},
        {}
    };

    games_["original-dragon-race"] = {
        "original-dragon-race", "InHouse", "Dragon Race",
        6, 4, 50, 0.50, 500.00, 96.80, 4.5, true, true,
        {"🐉", "🦁", "🐯", "🦅", "🐺", "⭐", "💎"},
        {}
    };

    // Initialize jackpots
    jackpots_["mini"] = {"Mini Jackpot", 1000, 1000, 50000, 0.5};
    jackpots_["major"] = {"Major Jackpot", 10000, 10000, 500000, 1.0};
    jackpots_["grand"] = {"Grand Jackpot", 100000, 100000, 5000000, 2.0};
}

std::vector<SlotGameEngine::SlotConfig> SlotGameEngine::getAllGames() const {
    std::vector<SlotConfig> result;
    for (const auto& game : games_) {
        result.push_back(game.second);
    }
    return result;
}

std::vector<SlotGameEngine::SlotConfig> SlotGameEngine::getGamesByProvider(
    const std::string& provider) const {
    std::vector<SlotConfig> result;
    for (const auto& game : games_) {
        if (game.second.provider == provider) {
            result.push_back(game.second);
        }
    }
    return result;
}

const SlotGameEngine::SlotConfig* SlotGameEngine::getGame(
    const std::string& gameId) const {
    auto it = games_.find(gameId);
    if (it != games_.end()) {
        return &it->second;
    }
    return nullptr;
}

std::vector<std::vector<std::string>> SlotGameEngine::generateReels(
    const SlotConfig& config,
    const std::string& serverSeed,
    const std::string& clientSeed,
    uint64_t nonce) {
    
    std::vector<std::vector<std::string>> result(config.reels);
    
    for (size_t reel = 0; reel < config.reels; ++reel) {
        result[reel].resize(config.rows);
        
        for (size_t row = 0; row < config.rows; ++row) {
            uint64_t outcome = ProvablyFair::generateOutcome(
                serverSeed,
                clientSeed,
                nonce + (reel * 100) + row
            );
            
            size_t symbolIndex = outcome % config.symbols.size();
            result[reel][row] = config.symbols[symbolIndex];
        }
    }
    
    return result;
}

std::vector<std::string> SlotGameEngine::findWins(
    const SlotConfig& config,
    const std::vector<std::vector<std::string>>& reels,
    double betAmount) {
    
    std::vector<std::string> wins;
    
    // Simplified win detection - check for matching symbols in paylines
    // In production, would check actual paylines
    
    // Check for scatter wins
    int scatterCount = 0;
    for (size_t r = 0; r < reels.size(); ++r) {
        for (size_t c = 0; c < reels[r].size(); ++c) {
            if (reels[r][c] == "⭐" || reels[r][c] == "💎") {
                scatterCount++;
            }
        }
    }
    
    if (scatterCount >= 3) {
        wins.push_back("Scatter: " + std::to_string(scatterCount) + " free spins!");
    }
    
    // Check for jackpot (rare)
    if (scatterCount >= 6) {
        wins.push_back("JACKPOT TRIGGERED!");
    }
    
    return wins;
}

SlotGameEngine::SlotResult SlotGameEngine::spin(
    const std::string& gameId,
    double betAmount,
    const std::string& clientSeed) {
    
    roundCounter_++;
    
    SlotResult result;
    result.gameId = gameId;
    result.roundId = roundCounter_;
    result.betAmount = betAmount;
    
    auto it = games_.find(gameId);
    if (it == games_.end()) {
        result.totalWin = 0;
        return result;
    }
    
    const SlotConfig& config = it->second;
    
    // Generate server seed
    result.serverSeed.seed = ProvablyFair::generateServerSeed();
    result.serverSeed.hash = ProvablyFair::hashSeed(result.serverSeed.seed);
    result.serverSeed.clientSeed = clientSeed.empty() ? 
        std::to_string(roundCounter_) : clientSeed;
    result.serverSeed.nonce = roundCounter_;
    
    // Generate reels
    result.reels = generateReels(config, result.serverSeed.seed,
                                result.serverSeed.clientSeed,
                                result.serverSeed.nonce);
    
    // Calculate winnings
    result.totalWin = calculateWin(gameId, result.reels, betAmount);
    
    // Check for bonus
    result.isBonus = checkBonusTrigger(result.reels);
    result.isJackpot = result.totalWin > betAmount * 100;
    
    if (result.isBonus) {
        result.bonusType = "Free Spins";
    }
    
    result.multiplier = result.totalWin / betAmount;
    
    // Find winning lines
    result.winningLines = findWins(config, result.reels, betAmount);
    
    return result;
}

double SlotGameEngine::calculateWin(
    const std::string& gameId,
    const std::vector<std::vector<std::string>>& reels,
    double betAmount) {
    
    auto it = games_.find(gameId);
    if (it == games_.end()) return 0.0;
    
    const SlotConfig& config = it->second;
    
    // Simplified win calculation
    // Count matching symbols
    double totalWin = 0.0;
    
    // Check scatter count
    int scatterCount = 0;
    for (size_t r = 0; r < reels.size(); ++r) {
        for (size_t c = 0; c < reels[r].size(); ++c) {
            if (reels[r][c] == "⭐") {
                scatterCount++;
            }
        }
    }
    
    // Scatter pays
    if (scatterCount >= 3) {
        totalWin += betAmount * scatterCount * 2;
    }
    
    // Check for specific symbol wins (simplified)
    // In production, would check actual paylines
    
    return totalWin;
}

bool SlotGameEngine::checkBonusTrigger(
    const std::vector<std::vector<std::string>>& reels) {
    
    int scatterCount = 0;
    for (size_t r = 0; r < reels.size(); ++r) {
        for (size_t c = 0; c < reels[r].size(); ++c) {
            if (reels[r][c] == "⭐") {
                scatterCount++;
            }
        }
    }
    
    return scatterCount >= 3;
}

void SlotGameEngine::contributeJackpot(const std::string& jackpotName, double amount) {
    auto it = jackpots_.find(jackpotName);
    if (it != jackpots_.end()) {
        it->second.currentAmount += amount * (it->second.contributionPercent / 100.0);
    }
}

double SlotGameEngine::getJackpotAmount(const std::string& jackpotName) const {
    auto it = jackpots_.find(jackpotName);
    if (it != jackpots_.end()) {
        return it->second.currentAmount;
    }
    return 0.0;
}

bool SlotGameEngine::triggerJackpot(const std::string& gameId) {
    // Simplified - random chance based on game config
    auto it = games_.find(gameId);
    if (it == games_.end() || !it->second.hasJackpot) {
        return false;
    }
    
    // 1 in 10,000 chance
    uint64_t seed = rng_->generateSeed();
    return seed % 10000 == 0;
}

} // namespace TigerCasino

#endif // SLOT_GAME_ENGINE_HPP
