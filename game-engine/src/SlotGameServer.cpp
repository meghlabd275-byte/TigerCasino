#include "SlotGameServer.hpp"
#include <algorithm>
#include <numeric>
#include <random>
#include <openssl/rand.h>

namespace TigerCasino {

// Build weighted symbol list for random selection
void SlotGameServer::buildWeightedSymbols() {
    weightedSymbols_.clear();
    for (const auto& symbol : config_.symbols) {
        int weight = config_.symbolWeights[symbol];
        for (int i = 0; i < weight; i++) {
            weightedSymbols_.push_back(symbol);
        }
    }
}

// Generate random reel symbols
std::vector<std::vector<std::string>> SlotGameServer::generateReels() {
    std::lock_guard<std::mutex> lock(rngMutex_);
    
    std::vector<std::vector<std::string>> reels;
    
    for (int r = 0; r < config_.reels; r++) {
        std::vector<std::string> reel;
        for (int pos = 0; pos < config_.rows; pos++) {
            std::uniform_int_distribution<size_t> dist(0, weightedSymbols_.size() - 1);
            reel.push_back(weightedSymbols_[dist(rng_)]);
        }
        reels.push_back(reel);
    }
    
    return reels;
}

// Get payline positions
std::vector<int> SlotGameServer::getPaylinePositions(int paylineIndex) const {
    std::vector<int> positions;
    int symbolsPerReel = config_.rows;
    
    // Different payline patterns
    switch (paylineIndex % 10) {
        case 0: // Top row
            for (int r = 0; r < config_.reels; r++) positions.push_back(r * symbolsPerReel);
            break;
        case 1: // Middle row
            for (int r = 0; r < config_.reels; r++) positions.push_back(r * symbolsPerReel + 1);
            break;
        case 2: // Bottom row
            for (int r = 0; r < config_.reels; r++) positions.push_back(r * symbolsPerReel + 2);
            break;
        case 3: // V shape
            for (int r = 0; r < config_.reels; r++) {
                positions.push_back(r * symbolsPerReel + (r % 2 == 0 ? 0 : 2));
            }
            break;
        case 4: // Inverted V
            for (int r = 0; r < config_.reels; r++) {
                positions.push_back(r * symbolsPerReel + (r % 2 == 0 ? 2 : 0));
            }
            break;
        default: // Zigzag
            for (int r = 0; r < config_.reels; r++) {
                positions.push_back(r * symbolsPerReel + ((r + paylineIndex) % 3));
            }
            break;
    }
    
    return positions;
}

// Calculate wins across paylines
std::vector<PaylineWin> SlotGameServer::calculateWins(
    const std::vector<std::vector<std::string>>& reels) {
    
    std::vector<PaylineWin> wins;
    
    // Check each payline
    for (int p = 0; p < config_.paylines; p++) {
        std::vector<int> positions = getPaylinePositions(p);
        
        // Get symbols on this payline
        std::string firstSymbol;
        std::vector<std::string> paylineSymbols;
        int matchCount = 0;
        
        for (size_t i = 0; i < positions.size(); i++) {
            int reelIdx = i;
            int posIdx = positions[i];
            
            if (reelIdx < reels.size() && posIdx < (int)reels[reelIdx].size()) {
                std::string symbol = reels[reelIdx][posIdx];
                paylineSymbols.push_back(symbol);
                
                if (i == 0) {
                    firstSymbol = symbol;
                    matchCount = 1;
                } else if (symbol == firstSymbol) {
                    matchCount++;
                } else {
                    break;  // Sequence broken
                }
            }
        }
        
        // Calculate win if we have matching symbols
        if (matchCount >= 3 && !firstSymbol.empty()) {
            double symbolMultiplier = config_.symbolValues[firstSymbol];
            double winAmount = config_.minBet * symbolMultiplier * matchCount;
            
            PaylineWin win;
            win.paylineIndex = p;
            win.symbol = firstSymbol;
            win.count = matchCount;
            win.winAmount = winAmount;
            win.positions = positions;
            wins.push_back(win);
        }
    }
    
    return wins;
}

// Determine if bonus should trigger
bool SlotGameServer::shouldTriggerBonus() {
    std::lock_guard<std::mutex> lock(rngMutex_);
    std::uniform_real_distribution<double> dist(0.0, 1.0);
    return dist(rng_) < 0.05;  // 5% bonus trigger rate
}

// Generate free spins count
int SlotGameServer::generateFreeSpins() {
    std::lock_guard<std::mutex> lock(rngMutex_);
    std::uniform_int_distribution<int> dist(5, 20);
    return dist(rng_);
}

// Constructor
SlotGameServer::SlotGameServer(const SlotConfig& config) 
    : config_(config)
    , totalSpins_(0)
    , totalPayout_(0.0) {
    buildWeightedSymbols();
}

// Spin the slots
SlotResult SlotGameServer::spin(const std::string& playerId, double betAmount) {
    SlotResult result;
    result.success = true;
    result.betAmount = betAmount;
    
    // Validate bet
    if (betAmount < config_.minBet || betAmount > config_.maxBet) {
        result.success = false;
        result.outcome = "Invalid bet amount";
        return result;
    }
    
    // Generate reels
    result.reelSymbols = generateReels();
    
    // Calculate wins
    result.paylineWins = calculateWins(result.reelSymbols);
    
    // Sum up wins
    double totalWin = 0;
    for (const auto& win : result.paylineWins) {
        totalWin += win.winAmount;
    }
    
    result.winAmount = totalWin;
    result.multiplier = totalWin / betAmount;
    
    // Check for bonus
    result.bonusTriggered = shouldTriggerBonus();
    if (result.bonusTriggered) {
        result.bonusType = "free_spins";
        result.metadata["freeSpins"] = std::to_string(generateFreeSpins());
    }
    
    // Update statistics
    totalSpins_++;
    totalPayout_ += totalWin;
    
    // Set outcome
    if (result.winAmount > 0) {
        result.outcome = "WIN";
    } else {
        result.outcome = "LOSE";
    }
    
    return result;
}

// Get configuration
SlotConfig SlotGameServer::getConfig() const {
    return config_;
}

// Update configuration
void SlotGameServer::updateConfig(const SlotConfig& config) {
    config_ = config;
    buildWeightedSymbols();
}

// Get statistics
std::map<std::string, double> SlotGameServer::getStatistics() const {
    std::map<std::string, double> stats;
    stats["totalSpins"] = totalSpins_.load();
    stats["totalPayout"] = totalPayout_.load();
    stats["rtp"] = totalSpins_.load() > 0 ? 
        (totalPayout_.load() / (totalSpins_.load() * config_.minBet)) * 100 : 0;
    return stats;
}

// Reset statistics
void SlotGameServer::resetStatistics() {
    totalSpins_ = 0;
    totalPayout_ = 0.0;
}

// Slot templates
namespace SlotTemplates {

SlotConfig createTigerSlots() {
    SlotConfig config;
    config.gameId = "tiger_slots";
    config.gameName = "Tiger Slots";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.01;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    
    config.symbols = {"TIGER", "WILD", "SCATTER", "A", "K", "Q", "J", "10", "9"};
    config.symbolValues = {
        {"TIGER", 10.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"A", 2.0}, {"K", 1.5}, {"Q", 1.2}, {"J", 1.0}, {"10", 0.8}, {"9", 0.5}
    };
    config.symbolWeights = {
        {"TIGER", 2}, {"WILD", 3}, {"SCATTER", 2},
        {"A", 8}, {"K", 10}, {"Q", 12}, {"J", 14}, {"10", 16}, {"9", 18}
    };
    
    return config;
}

SlotConfig createMegaMoolah() {
    SlotConfig config;
    config.gameId = "mega_moolah";
    config.gameName = "Mega Moolah";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.25;
    config.maxBet = 6.25;
    config.rtp = 0.88;  // Lower RTP due to progressive jackpot
    
    config.symbols = {"LION", "WILD", "SCATTER", "A", "K", "Q", "J", "10"};
    config.symbolValues = {
        {"LION", 15.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"A", 2.0}, {"K", 1.5}, {"Q", 1.2}, {"J", 1.0}, {"10", 0.8}
    };
    config.symbolWeights = {
        {"LION", 1}, {"WILD", 3}, {"SCATTER", 2},
        {"A", 10}, {"K", 12}, {"Q", 14}, {"J", 16}, {"10", 18}
    };
    
    return config;
}

SlotConfig createStarburst() {
    SlotConfig config;
    config.gameId = "starburst";
    config.gameName = "Starburst";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    
    config.symbols = {"STAR", "WILD", "SEVEN", "BAR", "BELL", "CHERRY"};
    config.symbolValues = {
        {"STAR", 10.0}, {"WILD", 5.0}, {"SEVEN", 4.0},
        {"BAR", 3.0}, {"BELL", 2.0}, {"CHERRY", 1.0}
    };
    config.symbolWeights = {
        {"STAR", 2}, {"WILD", 4}, {"SEVEN", 3},
        {"BAR", 8}, {"BELL", 12}, {"CHERRY", 20}
    };
    
    return config;
}

SlotConfig createGonzoQuest() {
    SlotConfig config;
    config.gameId = "gonzo_quest";
    config.gameName = "Gonzo's Quest";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 20;
    config.minBet = 0.20;
    config.maxBet = 50.0;
    config.rtp = 0.96;
    
    config.symbols = {"GONZO", "WILD", "SCATTER", "JAGUAR", "FROG", "FISH"};
    config.symbolValues = {
        {"GONZO", 12.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"JAGUAR", 2.5}, {"FROG", 1.5}, {"FISH", 1.0}
    };
    config.symbolWeights = {
        {"GONZO", 1}, {"WILD", 3}, {"SCATTER", 2},
        {"JAGUAR", 6}, {"FROG", 10}, {"FISH", 15}
    };
    
    return config;
}

SlotConfig createBookOfDead() {
    SlotConfig config;
    config.gameId = "book_of_dead";
    config.gameName = "Book of Dead";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    
    config.symbols = {"DEAD", "WILD", "SCATTER", "PHARAOH", "ANKH", "HORUS"};
    config.symbolValues = {
        {"DEAD", 15.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"PHARAOH", 2.5}, {"ANKH", 1.5}, {"HORUS", 1.0}
    };
    config.symbolWeights = {
        {"DEAD", 1}, {"WILD", 3}, {"SCATTER": 2},
        {"PHARAOH", 5}, {"ANKH", 10}, {"HORUS", 14}
    };
    
    return config;
}

SlotConfig createSweetBonanza() {
    SlotConfig config;
    config.gameId = "sweet_bonanza";
    config.gameName = "Sweet Bonanza";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;  // Cluster pays
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    
    config.symbols = {"CANDY", "WILD", "SCATTER", "BANANA", "APPLE", "GRAPE", "BERRY"};
    config.symbolValues = {
        {"CANDY", 20.0}, {"WILD", 0}, {"SCATTER", 5.0},
        {"BANANA", 3.0}, {"APPLE", 2.0}, {"GRAPE", 1.5}, {"BERRY", 1.0}
    };
    config.symbolWeights = {
        {"CANDY", 2}, {"WILD", 0}, {"SCATTER", 3},
        {"BANANA", 8}, {"APPLE", 10}, {"GRAPE", 14}, {"BERRY", 18}
    };
    
    return config;
}

SlotConfig createMoneyTrain() {
    SlotConfig config;
    config.gameId = "money_train";
    config.gameName = "Money Train";
    config.reels = 5;
    config.rows = 4;
    config.paylines = 40;
    config.minBet = 0.40;
    config.maxBet = 20.0;
    config.rtp = 0.96;
    
    config.symbols = {"TRAIN", "WILD", "SCATTER", "SHERIFF", "BANDIT", "COIN"};
    config.symbolValues = {
        {"TRAIN", 15.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"SHERIFF", 3.0}, {"BANDIT", 2.0}, {"COIN", 1.5}
    };
    config.symbolWeights = {
        {"TRAIN", 1}, {"WILD", 4}, {"SCATTER", 2},
        {"SHERIFF", 5}, {"BANDIT", 10}, {"COIN", 15}
    };
    
    return config;
}

SlotConfig createDeadOrAlive() {
    SlotConfig config;
    config.gameId = "dead_or_alive";
    config.gameName = "Dead or Alive";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 12;
    config.minBet = 0.09;
    config.maxBet = 9.00;
    config.rtp = 0.97;
    
    config.symbols = {"WANTED", "WILD", "SCATTER", "PISTOL", "BOOT", "HAT"};
    config.symbolValues = {
        {"WANTED", 12.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"PISTOL", 2.0}, {"BOOT", 1.5}, {"HAT", 1.0}
    };
    config.symbolWeights = {
        {"WANTED": 2}, {"WILD", 3}, {"SCATTER", 2},
        {"PISTOL", 8}, {"BOOT", 12}, {"HAT", 16}
    };
    
    return config;
}

SlotConfig createBigBassBonanza() {
    SlotConfig config;
    config.gameId = "big_bass_bonanza";
    config.gameName = "Big Bass Bonanza";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 250.0;
    config.rtp = 0.96;
    
    config.symbols = {"FISHER", "WILD", "SCATTER", "BASS", "DRAGON", "TROPHY"};
    config.symbolValues = {
        {"FISHER", 12.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"BASS", 2.5}, {"DRAGON", 1.8}, {"TROPHY", 1.2}
    };
    config.symbolWeights = {
        {"FISHER", 1}, {"WILD", 3}, {"SCATTER", 2},
        {"BASS", 6}, {"DRAGON", 10}, {"TROPHY", 14}
    };
    
    return config;
}

SlotConfig createTheDogHouse() {
    SlotConfig config;
    config.gameId = "the_dog_house";
    config.gameName = "The Dog House";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 20;
    config.minBet = 0.20;
    config.maxBet = 200.0;
    config.rtp = 0.96;
    
    config.symbols = {"DOG", "WILD", "SCATTER", "BONE", "COLLAR", "HOUSE"};
    config.symbolValues = {
        {"DOG", 12.0}, {"WILD", 5.0}, {"SCATTER", 3.0},
        {"BONE", 2.0}, {"COLLAR", 1.5}, {"HOUSE", 1.0}
    };
    config.symbolWeights = {
        {"DOG", 2}, {"WILD", 3}, {"SCATTER", 2},
        {"BONE", 8}, {"COLLAR", 12}, {"HOUSE", 16}
    };
    
    return config;
}

std::vector<SlotConfig> getAllTemplates() {
    return {
        createTigerSlots(),
        createMegaMoolah(),
        createStarburst(),
        createGonzoQuest(),
        createBookOfDead(),
        createSweetBonanza(),
        createMoneyTrain(),
        createDeadOrAlive(),
        createBigBassBonanza(),
        createTheDogHouse()
    };
}

} // namespace SlotTemplates

} // namespace TigerCasino
