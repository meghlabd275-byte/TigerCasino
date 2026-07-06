#pragma once

#include "GameEngine.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

class DiceGame : public Game {
private:
    RandomNumberGenerator rng_;
    double rtp_ = 0.99; // 99% RTP

public:
    DiceGame();
    virtual ~DiceGame() = default;
    
    std::string getName() const override;
    std::string getType() const override;
    double getRTP() const override;
    GameResult play(const BetInfo& bet) override;
};

} // namespace TigerCasino
