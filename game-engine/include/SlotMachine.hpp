#pragma once

#include "GameEngine.hpp"
#include "RandomNumberGenerator.hpp"
#include <array>

namespace TigerCasino {

// Slot machine game implementation
class SlotMachine : public Game {
private:
    RandomNumberGenerator rng_;
    double rtp_ = 0.965; // 96.5% RTP
    
    // Symbol definitions
    static constexpr const char* SYMBOLS[] = {
        "🔔", "💎", "🍒", "🍋", "🍇", "⭐", "🐯", "💰"
    };
    
    // Payout multipliers for 3 matching symbols
    static constexpr double PAYOUTS[] = {
        10.0,  // Bell
        25.0,  // Diamond
        5.0,   // Cherry
        3.0,   // Lemon
        15.0,  // Grape
        20.0,  // Star
        50.0,  // Tiger (jackpot)
        30.0   // Money
    };

public:
    SlotMachine();
    virtual ~SlotMachine() = default;
    
    std::string getName() const override;
    std::string getType() const override;
    double getRTP() const override;
    GameResult play(const BetInfo& bet) override;
    
private:
    std::array<int, 3> spin();
    double calculatePayout(const std::array<int, 3>& symbols, double betAmount);
};

} // namespace TigerCasino
