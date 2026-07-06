#include "SlotMachine.hpp"
#include <sstream>

namespace TigerCasino {

SlotMachine::SlotMachine() = default;

std::string SlotMachine::getName() const {
    return "Tiger King Slots";
}

std::string SlotMachine::getType() const {
    return "slots";
}

double SlotMachine::getRTP() const {
    return rtp_;
}

std::array<int, 3> SlotMachine::spin() {
    std::array<int, 3> result;
    for (int i = 0; i < 3; ++i) {
        result[i] = rng_.nextInt(0, 7); // 8 symbols
    }
    return result;
}

double SlotMachine::calculatePayout(const std::array<int, 3>& symbols, double betAmount) {
    // Check for three matching symbols
    if (symbols[0] == symbols[1] && symbols[1] == symbols[2]) {
        return betAmount * PAYOUTS[symbols[0]];
    }
    
    // Check for two matching symbols
    if (symbols[0] == symbols[1] || symbols[1] == symbols[2] || symbols[0] == symbols[2]) {
        return betAmount * 0.5;
    }
    
    return 0.0;
}

GameResult SlotMachine::play(const BetInfo& bet) {
    GameResult result;
    
    if (bet.amount <= 0) {
        result.success = false;
        result.outcome = "Invalid bet amount";
        return result;
    }
    
    // Spin the reels
    auto symbols = spin();
    
    // Calculate payout
    double payout = calculatePayout(symbols, bet.amount);
    
    result.success = true;
    result.winAmount = payout;
    result.multiplier = payout / bet.amount;
    
    std::ostringstream oss;
    oss << SYMBOLS[symbols[0]] << " " << SYMBOLS[symbols[1]] << " " << SYMBOLS[symbols[2]];
    result.outcome = oss.str();
    
    result.metadata["symbol1"] = SYMBOLS[symbols[0]];
    result.metadata["symbol2"] = SYMBOLS[symbols[1]];
    result.metadata["symbol3"] = SYMBOLS[symbols[2]];
    result.metadata["isWin"] = (payout > 0) ? "true" : "false";
    
    return result;
}

} // namespace TigerCasino
