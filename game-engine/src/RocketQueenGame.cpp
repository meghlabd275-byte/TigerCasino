#include "RocketQueenGame.hpp"
#include <iostream>

namespace TigerCasino {

RocketQueenGame::RocketQueenGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
    current_state_.crashed = false;
    current_state_.round_active = false;
    current_state_.current_multiplier = 1.0;
    current_state_.distance_km = 0.0;
}

double RocketQueenGame::calculateExponentialMultiplier(double time_ms) const {
    // Exponential growth with different curve than Aviator
    // Grows faster initially, then slows
    return 1.0 + std::pow(time_ms / 1000.0, 1.5) * 0.5;
}

double RocketQueenGame::calculateCrashPoint(const std::string& hash) const {
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash[i]);
    }
    
    double normalized = static_cast<double>(hash_value) / static_cast<double>(UINT64_MAX);
    
    // Different distribution - more volatile
    const double lambda = 0.035;
    double crash_multiplier = -std::log(1.0 - normalized) / lambda + 1.0;
    
    return std::min(crash_multiplier, MAX_MULTIPLIER);
}

bool RocketQueenGame::validateBet(double amount) const {
    return amount >= MIN_BET && amount <= MAX_BET;
}

RocketQueenGame::GameResult RocketQueenGame::createErrorResult(const std::string& error) const {
    GameResult result;
    result.success = false;
    result.error_message = error;
    result.round_id = current_state_.round_id;
    return result;
}

RocketQueenGame::GameResult RocketQueenGame::startRound() {
    GameResult result;
    
    round_counter_++;
    current_state_.round_id = round_counter_;
    current_state_.crashed = false;
    current_state_.round_active = true;
    current_state_.current_multiplier = 1.0;
    current_state_.distance_km = 0.0;
    current_state_.start_time = std::chrono::steady_clock::now();
    
    auto hash_result = provably_fair_->generateRandomHash();
    
    result.success = true;
    result.round_id = current_state_.round_id;
    result.server_seed = provably_fair_->getServerSeed();
    result.client_seed = provably_fair_->getClientSeed();
    result.multiplier = 1.0;
    result.outcome = "ROUND_STARTED";
    result.distance_km = 0.0;
    
    active_bets_.clear();
    
    return result;
}

RocketQueenGame::GameResult RocketQueenGame::placeBet(const std::string& player_id, double amount) {
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

RocketQueenGame::GameResult RocketQueenGame::cashOut(const std::string& player_id, double target_multiplier) {
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
            result.distance_km = current_state_.distance_km;
            
            return result;
        }
    }
    
    return createErrorResult("No active bet found");
}

RocketQueenGame::GameResult RocketQueenGame::crash() {
    GameResult result;
    
    if (!current_state_.round_active) {
        return createErrorResult("No active round");
    }
    
    current_state_.crashed = true;
    current_state_.round_active = false;
    
    double crash_point = calculateCrashPoint(provably_fair_->generateRandomHash().hash);
    current_state_.current_multiplier = crash_point;
    current_state_.distance_km = crash_point * 50.0; // 50km per multiplier
    
    result.success = true;
    result.win_amount = 0.0;
    result.multiplier = crash_point;
    result.round_id = current_state_.round_id;
    result.outcome = "CRASHED";
    result.distance_km = current_state_.distance_km;
    
    return result;
}

const RocketQueenGame::GameState& RocketQueenGame::getCurrentState() const {
    return current_state_;
}

double RocketQueenGame::getCurrentMultiplier() const {
    if (!current_state_.round_active) {
        return current_state_.current_multiplier;
    }
    
    auto now = std::chrono::steady_clock::now();
    auto elapsed = std::chrono::duration_cast<std::chrono::milliseconds>(now - current_state_.start_time).count();
    
    double mult = calculateExponentialMultiplier(static_cast<double>(elapsed));
    current_state_.current_multiplier = mult;
    current_state_.distance_km = mult * 50.0;
    
    return mult;
}

bool RocketQueenGame::isRoundActive() const {
    return current_state_.round_active;
}

std::vector<RocketQueenGame::Bet> RocketQueenGame::getActiveBets() const {
    return active_bets_;
}

std::string RocketQueenGame::getServerSeed() const {
    return provably_fair_->getServerSeed();
}

std::string RocketQueenGame::getClientSeed() const {
    return provably_fair_->getClientSeed();
}

void RocketQueenGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

void RocketQueenGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
