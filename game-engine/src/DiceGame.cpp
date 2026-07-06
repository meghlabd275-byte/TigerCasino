#include "DiceGame.hpp"
#include <sstream>

namespace TigerCasino {

DiceGame::DiceGame() = default;

std::string DiceGame::getName() const {
    return "Classic Dice";
}

std::string DiceGame::getType() const {
    return "dice";
}

double DiceGame::getRTP() const {
    return rtp_;
}

GameResult DiceGame::play(const BetInfo& bet) {
    GameResult result;
    
    if (bet.amount <= 0) {
        result.success = false;
        result.outcome = "Invalid bet amount";
        return result;
    }
    
    // Get roll target from params (default 50)
    double target = 50.0;
    auto it = bet.gameParams.find("target");
    if (it != bet.gameParams.end()) {
        target = it->second;
    }
    
    // Roll the dice (0-100)
    double roll = rng_.nextDouble() * 100.0;
    
    // Determine win/loss
    bool win = roll > target;
    double multiplier = win ? (100.0 / (100.0 - target)) : 0.0;
    double payout = win ? bet.amount * multiplier : 0.0;
    
    result.success = true;
    result.winAmount = payout;
    result.multiplier = multiplier;
    
    std::ostringstream oss;
    oss << "Roll: " << roll << ", Target: " << target;
    result.outcome = oss.str();
    
    result.metadata["roll"] = std::to_string(roll);
    result.metadata["target"] = std::to_string(target);
    result.metadata["win"] = win ? "true" : "false";
    
    return result;
}

} // namespace TigerCasino
