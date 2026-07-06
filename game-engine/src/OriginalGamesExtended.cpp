#include "OriginalGamesExtended.hpp"
#include <algorithm>
#include <numeric>
#include <openssl/rand.h>
#include <sstream>
#include <iomanip>

namespace TigerCasino {

// Base constructor
OriginalGameServer::OriginalGameServer() {
    std::random_device rd;
    rng_.seed(rd());
}

// ============== Plinko Game ==============

PlinkoGame::PlinkoGame() : rows_(16), riskMultiplier_(1.0) {
    // Payout table for 16 rows
    payouts_ = {
        // row 0 to row 16 (center to edges)
        {0.1, 0.2, 0.3, 0.5, 1.0, 2.0, 5.0, 10.0, 10.0, 5.0, 2.0, 1.0, 0.5, 0.3, 0.2, 0.1}
    };
}

OriginalGameConfig PlinkoGame::getConfig() const {
    return {
        "plinko",
        "Plinko",
        OriginalGameType::PLINKO,
        1.0,
        1000.0,
        0.98,
        true,
        "Drop the ball through pegs to win multipliers"
    };
}

OriginalGameResult PlinkoGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    int rows = params.count("rows") ? static_cast<int>(params.at("rows")) : rows_;
    double risk = params.count("risk") ? params.at("risk") : riskMultiplier_;
    
    // Determine landing position using random
    int position = 0;
    for (int i = 0; i < rows; i++) {
        std::lock_guard<std::mutex> lock(rngMutex_);
        bool goRight = std::uniform_int_distribution<int>(0, 1)(rng_) == 1;
        if (goRight) position++;
    }
    
    // Calculate multiplier based on position and risk
    double multiplier = 1.0;
    if (position >= 0 && position < static_cast<int>(payouts_[0].size())) {
        multiplier = payouts_[0][position] * risk;
    }
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = multiplier >= 1.0 ? "WIN" : "LOSE";
    
    return result;
}

// ============== Mines Game ==============

MinesGame::MinesGame() : mines_(3), gridSize_(25) {
}

OriginalGameConfig MinesGame::getConfig() const {
    return {
        "mines",
        "Mines",
        OriginalGameType::MINES,
        1.0,
        500.0,
        0.97,
        false,
        "Find gems, avoid mines"
    };
}

std::vector<std::string> MinesGame::generateMines(int count) {
    std::lock_guard<std::mutex> lock(rngMutex_);
    std::vector<int> indices(gridSize_);
    std::iota(indices.begin(), indices.end(), 0);
    std::shuffle(indices.begin(), indices.end(), rng_);
    
    std::vector<std::string> result;
    for (int i = 0; i < count && i < static_cast<int>(indices.size()); i++) {
        result.push_back("mine_" + std::to_string(indices[i]));
    }
    return result;
}

double MinesGame::calculatePayout(int gemsFound, double bet) {
    // Exponential payout for each gem found
    double multiplier = 1.0;
    for (int i = 0; i < gemsFound; i++) {
        multiplier *= 1.5;
    }
    return bet * multiplier;
}

OriginalGameResult MinesGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    int mines = params.count("mines") ? static_cast<int>(params.at("mines")) : mines_;
    
    result.success = true;
    result.multiplier = 1.0;
    result.winAmount = betAmount;
    result.outcome = "WIN";
    
    // Simulate finding gems
    int gemsFound = 5; // Simulated
    result.winAmount = calculatePayout(gemsFound, betAmount);
    result.multiplier = result.winAmount / betAmount;
    
    return result;
}

// ============== Limbo Game ==============

LimboGame::LimboGame() : targetMultiplier_(2.0) {
}

OriginalGameConfig LimboGame::getConfig() const {
    return {
        "limbo",
        "Limbo",
        OriginalGameType::LIMBO,
        1.0,
        1000.0,
        0.99,
        false,
        "Guess if the next number is higher or lower"
    };
}

OriginalGameResult LimboGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    double target = params.count("target") ? params.at("target") : targetMultiplier_;
    
    // Generate random multiplier (exponential distribution)
    std::lock_guard<std::mutex> lock(rngMutex_);
    double multiplier = 0.0;
    do {
        multiplier = -std::log(std::uniform_real_distribution<double>(0.0001, 1.0)(rng_)) * 10.0;
    } while (multiplier < 0.01 || multiplier > 1000.0);
    
    multiplier = std::min(multiplier, 1000.0);
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = (multiplier >= target) ? betAmount * multiplier : 0.0;
    result.outcome = multiplier >= target ? "WIN" : "LOSE";
    
    return result;
}

// ============== Keno Game ==============

KenoGame::KenoGame() : numbersPick_(10), numbersDrawn_(20) {
}

OriginalGameConfig KenoGame::getConfig() const {
    return {
        "keno",
        "Keno",
        OriginalGameType::KENO,
        1.0,
        100.0,
        0.95,
        false,
        "Pick numbers and match drawn numbers"
    };
}

OriginalGameResult KenoGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    int pick = params.count("pick") ? static_cast<int>(params.at("pick")) : numbersPick_;
    
    // Generate 20 random numbers from 1-80
    std::lock_guard<std::mutex> lock(rngMutex_);
    std::vector<int> drawn(80);
    std::iota(drawn.begin(), drawn.end(), 1);
    std::shuffle(drawn.begin(), drawn.end(), rng_);
    drawn.resize(numbersDrawn_);
    
    // Simulate matching (simplified)
    int matches = pick / 2; // Simulated
    
    // Payout table
    double multiplier = 0.0;
    if (matches >= 10) multiplier = 10.0;
    else if (matches >= 9) multiplier = 5.0;
    else if (matches >= 8) multiplier = 2.0;
    else if (matches >= 7) multiplier = 1.0;
    else if (matches >= 5) multiplier = 0.5;
    else multiplier = 0.0;
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = multiplier > 0 ? "WIN" : "LOSE";
    
    return result;
}

// ============== Bingo Game ==============

BingoGame::BingoGame() : ballCount_(75), cardSize_(5) {
}

OriginalGameConfig BingoGame::getConfig() const {
    return {
        "bingo",
        "Bingo",
        OriginalGameType::BINGO,
        1.0,
        100.0,
        0.95,
        false,
        "Complete patterns on your bingo card"
    };
}

OriginalGameResult BingoGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    // Simplified bingo result
    result.success = true;
    result.multiplier = 1.0;
    result.winAmount = betAmount;
    result.outcome = "WIN";
    
    return result;
}

// ============== Scratch Game ==============

ScratchGame::ScratchGame() : symbolsToReveal_(3), winningSymbols_(3) {
}

OriginalGameConfig ScratchGame::getConfig() const {
    return {
        "scratch",
        "Scratch Card",
        OriginalGameType::SCRATCH,
        1.0,
        100.0,
        0.95,
        false,
        "Scratch to reveal winning symbols"
    };
}

OriginalGameResult ScratchGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    // Simulate scratch
    std::lock_guard<std::mutex> lock(rngMutex_);
    bool isWinner = std::uniform_int_distribution<int>(0, 4)(rng_) == 0;
    
    double multiplier = isWinner ? 10.0 : 0.0;
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = isWinner ? "WIN" : "LOSE";
    
    return result;
}

// ============== Diamond Game ==============

DiamondGame::DiamondGame() : rows_(8), cols_(6) {
}

OriginalGameConfig DiamondGame::getConfig() const {
    return {
        "diamond",
        "Diamond Breaker",
        OriginalGameType::DIAMOND,
        1.0,
        500.0,
        0.96,
        false,
        "Break diamonds to multiply your win"
    };
}

OriginalGameResult DiamondGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    // Simulate diamond breaking
    std::lock_guard<std::mutex> lock(rngMutex_);
    int diamonds = std::uniform_int_distribution<int>(3, 15)(rng_);
    
    double multiplier = diamonds * 0.1 + 1.0;
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = multiplier > 1.0 ? "WIN" : "LOSE";
    
    return result;
}

// ============== Hash Dice Game ==============

HashDiceGame::HashDiceGame() : minWin_(0.01), maxWin_(1000.0) {
}

OriginalGameConfig HashDiceGame::getConfig() const {
    return {
        "hash_dice",
        "Hash Dice",
        OriginalGameType::HASH,
        1.0,
        10000.0,
        0.99,
        true,
        "Provably fair dice game"
    };
}

OriginalGameResult HashDiceGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    double chance = params.count("chance") ? params.at("chance") : 50.0;
    double target = 100.0 - chance;
    
    // Generate random roll
    std::lock_guard<std::mutex> lock(rngMutex_);
    double roll = std::uniform_real_distribution<double>(0, 100)(rng_);
    
    bool win = roll <= target;
    double multiplier = win ? (100.0 / chance) : 0.0;
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = win ? betAmount * multiplier : 0.0;
    result.outcome = win ? "WIN" : "LOSE";
    
    return result;
}

// ============== Coin Flip Game ==============

CoinFlipGame::CoinFlipGame() {
}

OriginalGameConfig CoinFlipGame::getConfig() const {
    return {
        "coinflip",
        "Coin Flip",
        OriginalGameType::COINFLIP,
        1.0,
        10000.0,
        0.99,
        false,
        "Double your winnings or lose"
    };
}

OriginalGameResult CoinFlipGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    std::lock_guard<std::mutex> lock(rngMutex_);
    bool heads = std::uniform_int_distribution<int>(0, 1)(rng_) == 0;
    
    std::string choice = params.count("choice") ? 
        (params.at("choice") > 0.5 ? "heads" : "tails") : "heads";
    
    bool win = (heads && choice == "heads") || (!heads && choice == "tails");
    
    result.success = true;
    result.multiplier = win ? 2.0 : 0.0;
    result.winAmount = win ? betAmount * 2.0 : 0.0;
    result.outcome = win ? "WIN" : "LOSE";
    
    return result;
}

// ============== Rock Paper Scissors ==============

RockPaperScissorsGame::RockPaperScissorsGame() {
}

OriginalGameConfig RockPaperScissorsGame::getConfig() const {
    return {
        "rps",
        "Rock Paper Scissors",
        OriginalGameType::ROCK_PAPER_SCISSORS,
        1.0,
        1000.0,
        0.98,
        false,
        "Classic rock paper scissors game"
    };
}

OriginalGameResult RockPaperScissorsGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    std::lock_guard<std::mutex> lock(rngMutex_);
    int playerChoice = params.count("choice") ? static_cast<int>(params.at("choice")) : 0;
    int cpuChoice = std::uniform_int_distribution<int>(0, 2)(rng_);
    
    // 0=rock, 1=paper, 2=scissors
    int diff = (playerChoice - cpuChoice + 3) % 3;
    
    double multiplier = 0.0;
    if (diff == 0) {
        multiplier = 1.0; // Tie
    } else if (diff == 1) {
        multiplier = 2.0; // Win
    } else {
        multiplier = 0.0; // Lose
    }
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = diff == 0 ? "TIE" : (diff == 1 ? "WIN" : "LOSE");
    
    return result;
}

// ============== Fortune Wheel ==============

FortuneWheelGame::FortuneWheelGame() {
    segments_ = {
        {"1", 1.0},
        {"2", 2.0},
        {"5", 5.0},
        {"10", 10.0},
        {"25", 25.0},
        {"50", 50.0},
        {"100", 100.0},
        {"500", 500.0}
    };
}

OriginalGameConfig FortuneWheelGame::getConfig() const {
    return {
        "fortune_wheel",
        "Fortune Wheel",
        OriginalGameType::WHEEL,
        1.0,
        100.0,
        0.95,
        false,
        "Spin the wheel to win multipliers"
    };
}

OriginalGameResult FortuneWheelGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    std::lock_guard<std::mutex> lock(rngMutex_);
    int segment = std::uniform_int_distribution<int>(0, static_cast<int>(segments_.size()) - 1)(rng_);
    
    double multiplier = segments_[segment].second;
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = multiplier >= 1.0 ? "WIN" : "LOSE";
    
    return result;
}

// ============== Tomb Game ==============

TombGame::TombGame() : levels_(5), levelMultiplier_(1.5) {
}

OriginalGameConfig TombGame::getConfig() const {
    return {
        "tomb",
        "Tomb",
        OriginalGameType::TOMB,
        1.0,
        500.0,
        0.96,
        false,
        "Escape the tomb by picking safe stones"
    };
}

OriginalGameResult TombGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    // Simulate tomb escape
    std::lock_guard<std::mutex> lock(rngMutex_);
    int reachedLevel = std::uniform_int_distribution<int>(1, levels_)(rng_);
    
    double multiplier = 1.0;
    for (int i = 0; i < reachedLevel; i++) {
        multiplier *= levelMultiplier_;
    }
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = "WIN";
    
    return result;
}

// ============== Avalanche Game ==============

AvalancheGame::AvalancheGame() : gridSize_(8), cascadeMultiplier_(1.2) {
}

OriginalGameConfig AvalancheGame::getConfig() const {
    return {
        "avalanche",
        "Avalanche",
        OriginalGameType::AVALANCHE,
        1.0,
        500.0,
        0.96,
        false,
        "Cascade wins with multiplier"
    };
}

OriginalGameResult AvalancheGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    // Simulate avalanche
    std::lock_guard<std::mutex> lock(rngMutex_);
    int cascades = std::uniform_int_distribution<int>(2, 8)(rng_);
    
    double multiplier = 1.0;
    for (int i = 0; i < cascades; i++) {
        multiplier *= cascadeMultiplier_;
    }
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = multiplier > 1.0 ? "WIN" : "LOSE";
    
    return result;
}

// ============== Tower Game ==============

TowerGame::TowerGame() : floors_(10), floorMultiplier_(1.5) {
}

OriginalGameConfig TowerGame::getConfig() const {
    return {
        "tower",
        "Tower",
        OriginalGameType::TOWER,
        1.0,
        500.0,
        0.96,
        false,
        "Climb the tower to win"
    };
}

OriginalGameResult TowerGame::play(double betAmount, const std::map<std::string, double>& params) {
    OriginalGameResult result;
    result.betAmount = betAmount;
    
    // Simulate climbing tower
    std::lock_guard<std::mutex> lock(rngMutex_);
    int floorsClimbed = std::uniform_int_distribution<int>(1, floors_)(rng_);
    
    double multiplier = 1.0;
    for (int i = 0; i < floorsClimbed; i++) {
        multiplier *= floorMultiplier_;
    }
    
    result.success = true;
    result.multiplier = multiplier;
    result.winAmount = betAmount * multiplier;
    result.outcome = "WIN";
    
    return result;
}

// ============== Original Games Manager ==============

OriginalGamesManager::OriginalGamesManager() {
    // Register all games
    registerGame("plinko", std::make_shared<PlinkoGame>());
    registerGame("mines", std::make_shared<MinesGame>());
    registerGame("limbo", std::make_shared<LimboGame>());
    registerGame("keno", std::make_shared<KenoGame>());
    registerGame("bingo", std::make_shared<BingoGame>());
    registerGame("scratch", std::make_shared<ScratchGame>());
    registerGame("diamond", std::make_shared<DiamondGame>());
    registerGame("hash_dice", std::make_shared<HashDiceGame>());
    registerGame("coinflip", std::make_shared<CoinFlipGame>());
    registerGame("rps", std::make_shared<RockPaperScissorsGame>());
    registerGame("fortune_wheel", std::make_shared<FortuneWheelGame>());
    registerGame("tomb", std::make_shared<TombGame>());
    registerGame("avalanche", std::make_shared<AvalancheGame>());
    registerGame("tower", std::make_shared<TowerGame>());
}

void OriginalGamesManager::registerGame(const std::string& gameId, 
                                       std::shared_ptr<OriginalGameServer> game) {
    games_[gameId] = game;
}

OriginalGameResult OriginalGamesManager::play(const std::string& gameId, double betAmount,
                                            const std::map<std::string, double>& params) {
    auto it = games_.find(gameId);
    if (it != games_.end()) {
        return it->second->play(betAmount, params);
    }
    return OriginalGameResult{false, betAmount, 0, 0, "Game not found", "", {}};
}

OriginalGameConfig OriginalGamesManager::getConfig(const std::string& gameId) const {
    auto it = games_.find(gameId);
    if (it != games_.end()) {
        return it->second->getConfig();
    }
    return OriginalGameConfig{"", "", OriginalGameType::CRASH, 0, 0, 0, false, ""};
}

std::vector<OriginalGameConfig> OriginalGamesManager::getAllConfigs() const {
    std::vector<OriginalGameConfig> configs;
    for (const auto& pair : games_) {
        configs.push_back(pair.second->getConfig());
    }
    return configs;
}

std::vector<std::string> OriginalGamesManager::getGameIds() const {
    std::vector<std::string> ids;
    for (const auto& pair : games_) {
        ids.push_back(pair.first);
    }
    return ids;
}

size_t OriginalGamesManager::getGameCount() const {
    return games_.size();
}

} // namespace TigerCasino
