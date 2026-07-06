#pragma once

#include <string>
#include <vector>
#include <map>
#include <array>
#include <memory>
#include <random>
#include <mutex>
#include <atomic>
#include <functional>

namespace TigerCasino {

// Slot game configuration
struct SlotConfig {
    std::string gameId;
    std::string gameName;
    int reels;
    int rows;
    int paylines;
    double minBet;
    double maxBet;
    double rtp;
    std::vector<std::string> symbols;
    std::map<std::string, double> symbolValues;  // Multiplier for each symbol
    std::map<std::string, int> symbolWeights;    // Weights for symbol distribution
};

// Slot payline win
struct PaylineWin {
    int paylineIndex;
    std::string symbol;
    int count;
    double winAmount;
    std::vector<int> positions;
};

// Slot game result
struct SlotResult {
    bool success;
    double betAmount;
    double winAmount;
    double multiplier;
    std::string outcome;
    std::vector<std::vector<std::string>> reelSymbols;
    std::vector<PaylineWin> paylineWins;
    bool bonusTriggered;
    std::string bonusType;
    std::map<std::string, std::string> metadata;
    
    SlotResult() : success(false), betAmount(0), winAmount(0), 
                   multiplier(0), bonusTriggered(false) {}
};

// Slot game server
class SlotGameServer {
private:
    SlotConfig config_;
    std::mt19937_64 rng_;
    std::mutex rngMutex_;
    std::atomic<uint64_t> totalSpins_;
    std::atomic<double> totalPayout_;
    
    // Symbol management
    std::vector<std::string> weightedSymbols_;
    void buildWeightedSymbols();
    
    // Reel generation
    std::vector<std::vector<std::string>> generateReels();
    
    // Win calculation
    std::vector<PaylineWin> calculateWins(const std::vector<std::vector<std::string>>& reels);
    std::vector<int> getPaylinePositions(int paylineIndex) const;
    
    // Bonus games
    bool shouldTriggerBonus();
    int generateFreeSpins();
    
public:
    explicit SlotGameServer(const SlotConfig& config);
    ~SlotGameServer() = default;
    
    // Gameplay
    SlotResult spin(const std::string& playerId, double betAmount);
    
    // Configuration
    SlotConfig getConfig() const;
    void updateConfig(const SlotConfig& config);
    
    // Statistics
    std::map<std::string, double> getStatistics() const;
    void resetStatistics();
};

// Pre-built slot game templates
namespace SlotTemplates {
    
SlotConfig createTigerSlots();
SlotConfig createMegaMoolah();
SlotConfig createStarburst();
SlotConfig createGonzoQuest();
SlotConfig createBookOfDead();
SlotConfig createSweetBonanza();
SlotConfig createMoneyTrain();
SlotConfig createDeadOrAlive();
SlotConfig createBigBassBonanza();
SlotConfig createTheDogHouse();
std::vector<SlotConfig> getAllTemplates();

} // namespace SlotTemplates

} // namespace TigerCasino
