#pragma once

#include "GameEngine.hpp"
#include "RandomNumberGenerator.hpp"
#include <vector>

namespace TigerCasino {

class Roulette : public Game {
private:
    RandomNumberGenerator rng_;
    double rtp_ = 0.973; // 97.3% RTP for European Roulette
    bool isAmerican_ = false;

    // European wheel (0-36)
    static constexpr int WHEEL[] = {
        0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10,
        5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26
    };

public:
    Roulette(bool american = false);
    virtual ~Roulette() = default;
    
    std::string getName() const override;
    std::string getType() const override;
    double getRTP() const override;
    GameResult play(const BetInfo& bet) override;
};

} // namespace TigerCasino
