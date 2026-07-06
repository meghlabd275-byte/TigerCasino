#pragma once

#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"
#include <string>
#include <vector>
#include <memory>
#include <cmath>
#include <chrono>
#include <random>

namespace TigerCasino {

/**
 * Rocket Queen Game - Female astronaut themed crash game
 * Dual-bet capability with unique social features
 */
class RocketQueenGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 10000.00;
    static constexpr double MIN_MULTIPLIER = 1.00;
    static constexpr double MAX_MULTIPLIER = 200.00;
    static constexpr double BASE_RTP = 0.96; // 96% RTP
    static constexpr size_t MAX_CONCURRENT_BETS = 2;

    struct GameState {
        uint64_t round_id;
        double current_multiplier;
        bool crashed;
        std::string crash_hash;
        std::chrono::steady_clock::time_point start_time;
        bool round_active;
        double distance_km;
        double fuel_level;
    };

    struct Bet {
        std::string player_id;
        double amount;
        double cashout_multiplier;
        bool cashed_out;
        size_t bet_index; // 0 or 1 for dual bets
        std::chrono::steady_clock::time_point bet_time;
    };

    struct GameResult {
        bool success;
        double win_amount;
        double multiplier;
        std::string outcome;
        std::string server_seed;
        std::string client_seed;
        uint64_t round_id;
        double distance_km;
        double fuel_level;
        size_t bet_index;
        std::string error_message;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    GameState current_state_;
    std::vector<Bet> active_bets_;
    uint64_t round_counter_;

public:
    RocketQueenGame();
    ~RocketQueenGame() = default;

    GameResult startRound();
    GameResult placeBet(const std::string& player_id, double amount, size_t bet_index = 0);
    GameResult cashOut(const std::string& player_id, size_t bet_index = 0);
    GameResult crash();
    
    const GameState& getCurrentState() const;
    double getCurrentMultiplier() const;
    bool isRoundActive() const;
    std::vector<Bet> getActiveBets() const;
    std::vector<Bet> getPlayerBets(const std::string& player_id) const;
    
    std::string getServerSeed() const;
    std::string getClientSeed() const;
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    double calculateCrashPoint(const std::string& hash) const;
    double calculateMultiplier(double time_ms) const;
    bool validateBet(double amount) const;
    bool validateBetIndex(size_t index) const;
    GameResult createErrorResult(const std::string& error) const;
};

inline RocketQueenGame::RocketQueenGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
    current_state_.crashed = false;
    current_state_.round_active = false;
    current_state_.current_multiplier = 1.0;
    current_state_.distance_km = 0.0;
    current_state_.fuel_level = 100.0;
}

inline double RocketQueenGame::calculateMultiplier(double time_ms) const {
    // Rocket-style acceleration curve
    double seconds = time_ms / 1000.0;
    return 1.0 + std::pow(seconds, 1.5) * 0.12;
}

inline double RocketQueenGame::calculateCrashPoint(const std::string& hash) const {
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash[i]);
    }
    
    double normalized = static_cast<double>(hash_value) / static_cast<double>(UINT64_MAX);
    
    // Higher volatility - more frequent low crashes
    const double lambda = 0.06;
    double crash_multiplier = -std::log(1.0 - normalized) / lambda + 1.0;
    
    return std::min(crash_multiplier, MAX_MULTIPLIER);
}

inline bool RocketQueenGame::validateBet(double amount) const {
    return amount >= MIN_BET && amount <= MAX_BET;
}

inline bool RocketQueenGame::validateBetIndex(size_t index) const {
    return index < MAX_CONCURRENT_BETS;
}

inline RocketQueenGame::GameResult RocketQueenGame::createErrorResult(const std::string& error) const {
    GameResult result;
    result.success = false;
    result.error_message = error;
    result.round_id = current_state_.round_id;
    return result;
}

inline RocketQueenGame::GameResult RocketQueenGame::startRound() {
    GameResult result;
    
    round_counter_++;
    current_state_.round_id = round_counter_;
    current_state_.crashed = false;
    current_state_.round_active = true;
    current_state_.current_multiplier = 1.0;
    current_state_.distance_km = 0.0;
    current_state_.fuel_level = 100.0;
    current_state_.start_time = std::chrono::steady_clock::now();
    
    auto hash_result = provably_fair_->generateRandomHash();
    current_state_.crash_hash = hash_result.hash;
    
    result.success = true;
    result.round_id = current_state_.round_id;
    result.server_seed = provably_fair_->getServerSeed();
    result.client_seed = provably_fair_->getClientSeed();
    result.multiplier = 1.0;
    result.outcome = "ROUND_STARTED";
    result.distance_km = 0.0;
    result.fuel_level = 100.0;
    result.bet_index = 0;
    
    active_bets_.clear();
    
    return result;
}

inline RocketQueenGame::GameResult RocketQueenGame::placeBet(const std::string& player_id, double amount, size_t bet_index) {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    if (!validateBet(amount)) {
        return createErrorResult("Invalid bet amount");
    }
    
    if (!validateBetIndex(bet_index)) {
        return createErrorResult("Invalid bet index");
    }
    
    // Check if player already placed bet at this index
    for (const auto& bet : active_bets_) {
        if (bet.player_id == player_id && bet.bet_index == bet_index) {
            return createErrorResult("Bet already placed at this index");
        }
    }
    
    Bet new_bet;
    new_bet.player_id = player_id;
    new_bet.amount = amount;
    new_bet.cashed_out = false;
    new_bet.cashout_multiplier = 0.0;
    new_bet.bet_index = bet_index;
    new_bet.bet_time = std::chrono::steady_clock::now();
    
    active_bets_.push_back(new_bet);
    
    result.success = true;
    result.win_amount = amount;
    result.round_id = current_state_.round_id;
    result.outcome = "BET_PLACED";
    result.bet_index = bet_index;
    
    return result;
}

inline RocketQueenGame::GameResult RocketQueenGame::cashOut(const std::string& player_id, size_t bet_index) {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    if (!validateBetIndex(bet_index)) {
        return createErrorResult("Invalid bet index");
    }
    
    for (auto& bet : active_bets_) {
        if (bet.player_id == player_id && bet.bet_index == bet_index && !bet.cashed_out) {
            double current_mult = getCurrentMultiplier();
            
            bet.cashed_out = true;
            bet.cashout_multiplier = current_mult;
            
            double win = bet.amount * current_mult;
            
            result.success = true;
            result.win_amount = win;
            result.multiplier = current_mult;
            result.round_id = current_state_.round_id;
            result.outcome = "CASHED_OUT";
            result.distance_km = current_state_.distance_km;
            result.fuel_level = current_state_.fuel_level;
            result.bet_index = bet_index;
            
            return result;
        }
    }
    
    return createErrorResult("No active bet found at this index");
}

inline RocketQueenGame::GameResult RocketQueenGame::crash() {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    current_state_.crashed = true;
    current_state_.round_active = false;
    
    double crash_point = calculateCrashPoint(current_state_.crash_hash);
    current_state_.current_multiplier = crash_point;
    current_state_.distance_km = crash_point * 5.0;
    current_state_.fuel_level = 0.0;
    
    // Calculate total payout
    double total_payout = 0.0;
    for (const auto& bet : active_bets_) {
        if (bet.cashed_out) {
            total_payout += bet.amount * bet.cashout_multiplier;
        }
    }
    
    result.success = true;
    result.win_amount = total_payout;
    result.multiplier = crash_point;
    result.round_id = current_state_.round_id;
    result.outcome = "CRASHED";
    result.distance_km = current_state_.distance_km;
    result.fuel_level = 0.0;
    
    return result;
}

inline const RocketQueenGame::GameState& RocketQueenGame::getCurrentState() const {
    return current_state_;
}

inline double RocketQueenGame::getCurrentMultiplier() const {
    if (!current_state_.round_active) {
        return current_state_.current_multiplier;
    }
    
    auto now = std::chrono::steady_clock::now();
    auto elapsed = std::chrono::duration_cast<std::chrono::milliseconds>(now - current_state_.start_time).count();
    
    double mult = calculateMultiplier(static_cast<double>(elapsed));
    current_state_.current_multiplier = mult;
    current_state_.distance_km = mult * 5.0;
    current_state_.fuel_level = std::max(0.0, 100.0 - (mult * 0.5));
    
    return mult;
}

inline bool RocketQueenGame::isRoundActive() const {
    return current_state_.round_active;
}

inline std::vector<RocketQueenGame::Bet> RocketQueenGame::getActiveBets() const {
    return active_bets_;
}

inline std::vector<RocketQueenGame::Bet> RocketQueenGame::getPlayerBets(const std::string& player_id) const {
    std::vector<Bet> player_bets;
    for (const auto& bet : active_bets_) {
        if (bet.player_id == player_id) {
            player_bets.push_back(bet);
        }
    }
    return player_bets;
}

inline std::string RocketQueenGame::getServerSeed() const {
    return provably_fair_->getServerSeed();
}

inline std::string RocketQueenGame::getClientSeed() const {
    return provably_fair_->getClientSeed();
}

inline void RocketQueenGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void RocketQueenGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
