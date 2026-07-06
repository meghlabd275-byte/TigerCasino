#pragma once

#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"
#include <string>
#include <vector>
#include <memory>
#include <cmath>
#include <chrono>

namespace TigerCasino {

/**
 * Spaceman Game - Astronaut flying game with rising multiplier
 * Ultra-low latency implementation with smooth animation
 */
class SpacemanGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 10000.00;
    static constexpr double MIN_MULTIPLIER = 1.00;
    static constexpr double MAX_MULTIPLIER = 500.00;
    static constexpr double BASE_RTP = 0.965; // 96.5% RTP

    struct GameState {
        uint64_t round_id;
        double current_multiplier;
        bool crashed;
        uint64_t crash_point_crypto;
        std::chrono::steady_clock::time_point start_time;
        bool round_active;
        double altitude_km; // Virtual altitude for visual
    };

    struct Bet {
        std::string player_id;
        double amount;
        double cashout_multiplier;
        bool cashed_out;
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
        double altitude_km;
        std::string error_message;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    GameState current_state_;
    std::vector<Bet> active_bets_;
    uint64_t round_counter_;

public:
    SpacemanGame();
    ~SpacemanGame() = default;

    GameResult startRound();
    GameResult placeBet(const std::string& player_id, double amount);
    GameResult cashOut(const std::string& player_id, double target_multiplier);
    GameResult crash();
    
    const GameState& getCurrentState() const;
    double getCurrentMultiplier() const;
    bool isRoundActive() const;
    std::vector<Bet> getActiveBets() const;
    
    std::string getServerSeed() const;
    std::string getClientSeed() const;
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    double calculateCrashPoint(const std::string& hash) const;
    double calculateLinearMultiplier(double time_ms) const;
    bool validateBet(double amount) const;
    GameResult createErrorResult(const std::string& error) const;
};

inline SpacemanGame::SpacemanGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
    current_state_.crashed = false;
    current_state_.round_active = false;
    current_state_.current_multiplier = 1.0;
    current_state_.altitude_km = 0.0;
}

inline double SpacemanGame::calculateLinearMultiplier(double time_ms) const {
    // Linear growth: 1x at 0s, increases at 0.1x per second
    // More predictable than exponential
    return 1.0 + (time_ms / 1000.0) * 0.15;
}

inline double SpacemanGame::calculateCrashPoint(const std::string& hash) const {
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash[i]);
    }
    
    double normalized = static_cast<double>(hash_value) / static_cast<double>(UINT64_MAX);
    
    // Linear-ish distribution with higher chance of early crashes
    const double lambda = 0.05;
    double crash_multiplier = -std::log(1.0 - normalized) / lambda + 1.0;
    
    return std::min(crash_multiplier, MAX_MULTIPLIER);
}

inline bool SpacemanGame::validateBet(double amount) const {
    return amount >= MIN_BET && amount <= MAX_BET;
}

inline SpacemanGame::GameResult SpacemanGame::createErrorResult(const std::string& error) const {
    GameResult result;
    result.success = false;
    result.error_message = error;
    result.round_id = current_state_.round_id;
    return result;
}

inline SpacemanGame::GameResult SpacemanGame::startRound() {
    GameResult result;
    
    round_counter_++;
    current_state_.round_id = round_counter_;
    current_state_.crashed = false;
    current_state_.round_active = true;
    current_state_.current_multiplier = 1.0;
    current_state_.altitude_km = 0.0;
    current_state_.start_time = std::chrono::steady_clock::now();
    
    auto hash_result = provably_fair_->generateRandomHash();
    
    result.success = true;
    result.round_id = current_state_.round_id;
    result.server_seed = provably_fair_->getServerSeed();
    result.client_seed = provably_fair_->getClientSeed();
    result.multiplier = 1.0;
    result.outcome = "ROUND_STARTED";
    result.altitude_km = 0.0;
    
    active_bets_.clear();
    
    return result;
}

inline SpacemanGame::GameResult SpacemanGame::placeBet(const std::string& player_id, double amount) {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    if (!validateBet(amount)) {
        return createErrorResult("Invalid bet amount");
    }
    
    for (const auto& bet : active_bets_) {
        if (bet.player_id == player_id) {
            return createErrorResult("Player already placed a bet");
        }
    }
    
    Bet new_bet;
    new_bet.player_id = player_id;
    new_bet.amount = amount;
    new_bet.cashed_out = false;
    new_bet.cashout_multiplier = 0.0;
    new_bet.bet_time = std::chrono::steady_clock::now();
    
    active_bets_.push_back(new_bet);
    
    result.success = true;
    result.win_amount = amount;
    result.round_id = current_state_.round_id;
    result.outcome = "BET_PLACED";
    
    return result;
}

inline SpacemanGame::GameResult SpacemanGame::cashOut(const std::string& player_id, double target_multiplier) {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    for (auto& bet : active_bets_) {
        if (bet.player_id == player_id && !bet.cashed_out) {
            double current_mult = getCurrentMultiplier();
            
            if (target_multiplier > current_mult) {
                return createErrorResult("Target multiplier too high");
            }
            
            bet.cashed_out = true;
            bet.cashout_multiplier = current_mult;
            
            double win = bet.amount * current_mult;
            
            result.success = true;
            result.win_amount = win;
            result.multiplier = current_mult;
            result.round_id = current_state_.round_id;
            result.outcome = "CASHED_OUT";
            result.altitude_km = current_state_.altitude_km;
            
            return result;
        }
    }
    
    return createErrorResult("No active bet found");
}

inline SpacemanGame::GameResult SpacemanGame::crash() {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    current_state_.crashed = true;
    current_state_.round_active = false;
    
    double crash_point = calculateCrashPoint(provably_fair_->generateRandomHash().hash);
    current_state_.current_multiplier = crash_point;
    current_state_.altitude_km = crash_point * 10.0; // 10km per multiplier
    
    result.success = true;
    result.win_amount = 0.0;
    result.multiplier = crash_point;
    result.round_id = current_state_.round_id;
    result.outcome = "CRASHED";
    result.altitude_km = current_state_.altitude_km;
    
    return result;
}

inline const SpacemanGame::GameState& SpacemanGame::getCurrentState() const {
    return current_state_;
}

inline double SpacemanGame::getCurrentMultiplier() const {
    if (!current_state_.round_active) {
        return current_state_.current_multiplier;
    }
    
    auto now = std::chrono::steady_clock::now();
    auto elapsed = std::chrono::duration_cast<std::chrono::milliseconds>(now - current_state_.start_time).count();
    
    double mult = calculateLinearMultiplier(static_cast<double>(elapsed));
    current_state_.current_multiplier = mult;
    current_state_.altitude_km = mult * 10.0;
    
    return mult;
}

inline bool SpacemanGame::isRoundActive() const {
    return current_state_.round_active;
}

inline std::vector<SpacemanGame::Bet> SpacemanGame::getActiveBets() const {
    return active_bets_;
}

inline std::string SpacemanGame::getServerSeed() const {
    return provably_fair_->getServerSeed();
}

inline std::string SpacemanGame::getClientSeed() const {
    return provably_fair_->getClientSeed();
}

inline void SpacemanGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void SpacemanGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
