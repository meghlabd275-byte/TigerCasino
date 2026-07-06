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
 * Aviator Game - High-speed crash game with airplane multiplier
 * Ultra-low latency implementation for real-time multiplayer gaming
 */
class AviatorGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 10000.00;
    static constexpr double MIN_MULTIPLIER = 1.00;
    static constexpr double MAX_MULTIPLIER = 1000.00;
    static constexpr double BASE_RTP = 0.97; // 97% RTP

    struct GameState {
        uint64_t round_id;
        double current_multiplier;
        bool crashed;
        uint64_t crash_point_crypto; // SHA256 hash for provably fair
        std::chrono::steady_clock::time_point start_time;
        bool round_active;
    };

    struct Bet {
        std::string player_id;
        double amount;
        double cashout_multiplier; // 0 if not cashed out yet
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
        std::string error_message;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    GameState current_state_;
    std::vector<Bet> active_bets_;
    uint64_t round_counter_;

public:
    AviatorGame();
    ~AviatorGame() = default;

    // Game lifecycle
    GameResult startRound();
    GameResult placeBet(const std::string& player_id, double amount);
    GameResult cashOut(const std::string& player_id, double target_multiplier);
    GameResult crash();
    
    // Query current state
    const GameState& getCurrentState() const;
    double getCurrentMultiplier() const;
    bool isRoundActive() const;
    std::vector<Bet> getActiveBets() const;
    std::vector<Bet> getPlayerBets(const std::string& player_id) const;
    
    // Provably fair
    std::string getServerSeed() const;
    std::string getClientSeed() const;
    
    // Configuration
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    double calculateCrashPoint(const std::string& hash) const;
    double calculateExponentialMultiplier(double time_ms) const;
    bool validateBet(double amount) const;
    GameResult createErrorResult(const std::string& error) const;
};

// Inline implementation for ultra-low latency
inline AviatorGame::AviatorGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
    current_state_.crashed = false;
    current_state_.round_active = false;
    current_state_.current_multiplier = 1.0;
}

inline double AviatorGame::calculateExponentialMultiplier(double time_ms) const {
    // Multiplier grows exponentially: 1x at 0ms, increases over time
    // Formula: multiplier = exp(time_ms / 1000.0 * 0.0006)
    // Average crash around 15-25 seconds (multiplier ~10x-15x)
    return std::exp(time_ms / 1000.0 * 0.0006);
}

inline double AviatorGame::calculateCrashPoint(const std::string& hash) const {
    // Convert hash to number and normalize to [0, 1]
    // Then apply crash distribution curve
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash[i]);
    }
    
    double normalized = static_cast<double>(hash_value) / static_cast<double>(UINT64_MAX);
    
    // Exponential distribution for realistic crash points
    // P(crash < x) = 1 - e^(-lambda * x)
    // lambda = 0.04 gives average crash around 25x
    const double lambda = 0.04;
    double crash_multiplier = -std::log(1.0 - normalized) / lambda + 1.0;
    
    // Cap at MAX_MULTIPLIER
    return std::min(crash_multiplier, MAX_MULTIPLIER);
}

inline bool AviatorGame::validateBet(double amount) const {
    return amount >= MIN_BET && amount <= MAX_BET;
}

inline AviatorGame::GameResult AviatorGame::createErrorResult(const std::string& error) const {
    GameResult result;
    result.success = false;
    result.error_message = error;
    result.round_id = current_state_.round_id;
    return result;
}

inline AviatorGame::GameResult AviatorGame::startRound() {
    GameResult result;
    
    // Generate new round
    round_counter_++;
    current_state_.round_id = round_counter_;
    current_state_.crashed = false;
    current_state_.round_active = true;
    current_state_.current_multiplier = 1.0;
    current_state_.start_time = std::chrono::steady_clock::now();
    
    // Generate provably fair crash point
    auto hash_result = provably_fair_->generateRandomHash();
    current_state_.crash_point_crypto = hash_result.hash;
    
    // Calculate crash point
    double crash_point = calculateCrashPoint(hash_result.hash);
    
    result.success = true;
    result.round_id = current_state_.round_id;
    result.server_seed = provably_fair_->getServerSeed();
    result.client_seed = provably_fair_->getClientSeed();
    result.multiplier = crash_point;
    result.outcome = "ROUND_STARTED";
    
    active_bets_.clear();
    
    return result;
}

inline AviatorGame::GameResult AviatorGame::placeBet(const std::string& player_id, double amount) {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    if (!validateBet(amount)) {
        return createErrorResult("Invalid bet amount");
    }
    
    // Check if player already placed bet
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

inline AviatorGame::GameResult AviatorGame::cashOut(const std::string& player_id, double target_multiplier) {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    // Find player's bet
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
            
            return result;
        }
    }
    
    return createErrorResult("No active bet found");
}

inline AviatorGame::GameResult AviatorGame::crash() {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    current_state_.crashed = true;
    current_state_.round_active = false;
    
    double crash_point = calculateCrashPoint(current_state_.crash_point_crypto);
    current_state_.current_multiplier = crash_point;
    
    // Process uncashed bets (they lose)
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
    result.server_seed = provably_fair_->getServerSeed();
    result.client_seed = provably_fair_->getClientSeed();
    
    return result;
}

inline const AviatorGame::GameState& AviatorGame::getCurrentState() const {
    return current_state_;
}

inline double AviatorGame::getCurrentMultiplier() const {
    if (!current_state_.round_active) {
        return current_state_.current_multiplier;
    }
    
    auto now = std::chrono::steady_clock::now();
    auto elapsed = std::chrono::duration_cast<std::chrono::milliseconds>(now - current_state_.start_time).count();
    
    return calculateExponentialMultiplier(static_cast<double>(elapsed));
}

inline bool AviatorGame::isRoundActive() const {
    return current_state_.round_active;
}

inline std::vector<AviatorGame::Bet> AviatorGame::getActiveBets() const {
    return active_bets_;
}

inline std::vector<AviatorGame::Bet> AviatorGame::getPlayerBets(const std::string& player_id) const {
    std::vector<Bet> player_bets;
    for (const auto& bet : active_bets_) {
        if (bet.player_id == player_id) {
            player_bets.push_back(bet);
        }
    }
    return player_bets;
}

inline std::string AviatorGame::getServerSeed() const {
    return provably_fair_->getServerSeed();
}

inline std::string AviatorGame::getClientSeed() const {
    return provably_fair_->getClientSeed();
}

inline void AviatorGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void AviatorGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
