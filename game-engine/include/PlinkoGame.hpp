#ifndef PLINKO_GAME_HPP
#define PLINKO_GAME_HPP

#include <string>
#include <vector>
#include <cstdint>
#include <memory>
#include <map>
#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

/**
 * Plinko Game - Ultra Low Latency Implementation
 * Balls fall through pegged board to multiplier pockets
 */
class PlinkoGame {
public:
    enum class Risk {
        Low,
        Medium,
        High
    };
    
    enum class Status {
        Waiting,
        Dropping,
        Complete
    };
    
    struct Config {
        uint8_t rows;
        Risk risk;
        
        Config() : rows(8), risk(Risk::Medium) {}
    };
    
    struct Ball {
        std::string ballId;
        std::vector<uint8_t> path;
        uint8_t finalPosition;
        double multiplier;
        double payout;
    };
    
    struct Bet {
        std::string betId;
        std::string playerId;
        double betAmount;
        Config config;
        std::vector<Ball> balls;
        double totalPayout;
        uint64_t timestamp;
    };
    
    struct GameState {
        std::string gameId;
        Status status;
        Config config;
        std::vector<Bet> bets;
    };

private:
    static constexpr double HOUSE_EDGE = 0.04; // 4%
    
    Config config_;
    GameState currentState_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    std::map<uint8_t, std::vector<double>> payoutTables_;

public:
    PlinkoGame() : rng_(std::make_unique<RandomNumberGenerator>()) {
        initPayoutTables();
    }
    
    void setConfig(uint8_t rows, Risk risk) {
        config_.rows = std::min(std::max(rows, (uint8_t)8), (uint8_t)16);
        config_.risk = risk;
    }
    
    /**
     * Initialize payout tables for all configurations
     */
    void initPayoutTables() {
        // 8 rows
        payoutTables_[8] = getPayoutTable(8, Risk::Medium);
        // 10 rows
        payoutTables_[10] = getPayoutTable(10, Risk::Medium);
        // 12 rows
        payoutTables_[12] = getPayoutTable(12, Risk::Medium);
        // 16 rows
        payoutTables_[16] = getPayoutTable(16, Risk::Medium);
    }
    
    std::vector<double> getPayoutTable(uint8_t rows, Risk risk) {
        switch(rows) {
            case 8:
                return {5.0, 2.0, 1.0, 0.5, 0.5, 1.0, 2.0, 5.0};
            case 10:
                return {5.0, 2.5, 1.5, 0.8, 0.4, 0.4, 0.8, 1.5, 2.5, 5.0};
            case 12:
                return {4.0, 2.0, 1.2, 0.7, 0.3, 0.3, 0.3, 0.3, 0.7, 1.2, 2.0, 4.0};
            case 16:
                return {1.0, 1.5, 2.0, 3.0, 4.0, 5.0, 5.0, 6.0, 
                        6.0, 5.0, 5.0, 4.0, 3.0, 2.0, 1.5, 1.0};
            default:
                return {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0};
        }
    }
    
    /**
     * Drop ball using provably fair seeds
     */
    Ball dropBall(const std::string& serverSeed,
                  const std::string& clientSeed,
                  uint64_t nonce) {
        Ball ball;
        ball.ballId = generateBallId();
        
        uint8_t rows = config_.rows;
        uint8_t position = rows / 2; // Start at center
        
        // Simulate ball path through each row
        for (uint8_t row = 0; row < rows; ++row) {
            uint64_t outcome = ProvablyFair::generateOutcome(
                serverSeed, 
                clientSeed, 
                nonce + row
            );
            
            // 50/50 chance left or right
            bool goRight = (outcome % 2) == 1;
            
            // Can only go right if not at edge
            uint8_t maxPos = row + 1;
            if (goRight && position < maxPos) {
                position++;
            }
            // If goLeft and position > 0, would go left
            
            ball.path.push_back(position);
        }
        
        // Get payout table for final multiplier
        auto& table = payoutTables_[rows];
        ball.finalPosition = std::min(position, (uint8_t)(table.size() - 1));
        ball.multiplier = table[ball.finalPosition];
        ball.payout = 0.0; // Will be set when bet is processed
        
        return ball;
    }
    
    /**
     * Place bet and process ball drop
     */
    Bet placeBet(const std::string& playerId,
                 double betAmount,
                 const std::string& serverSeed,
                 const std::string& clientSeed,
                 uint64_t nonce,
                 uint8_t ballsCount = 1) {
        Bet bet;
        bet.betId = generateBetId();
        bet.playerId = playerId;
        bet.betAmount = betAmount;
        bet.config = config_;
        bet.totalPayout = 0.0;
        bet.timestamp = getCurrentTimestamp();
        
        // Drop requested number of balls
        for (uint8_t i = 0; i < ballsCount; ++i) {
            Ball ball = dropBall(serverSeed, clientSeed, nonce + i);
            ball.payout = betAmount * ball.multiplier * (1.0 - HOUSE_EDGE);
            bet.totalPayout += ball.payout;
            bet.balls.push_back(ball);
        }
        
        currentState_.bets.push_back(bet);
        return bet;
    }
    
    /**
     * Get payout table for current config
     */
    const std::vector<double>& getCurrentPayoutTable() {
        return payoutTables_[config_.rows];
    }
    
    const GameState& getState() const { return currentState_; }

private:
    uint64_t getCurrentTimestamp() {
        auto now = std::chrono::system_clock::now();
        return std::chrono::duration_cast<std::chrono::milliseconds>(now.time_since_epoch()).count();
    }
    
    std::string generateBallId() {
        uint64_t id = rng_->generateSeed();
        return "BALL-" + std::to_string(id);
    }
    
    std::string generateBetId() {
        uint64_t id = rng_->generateSeed();
        return "BET-" + std::to_string(id);
    }
};

} // namespace TigerCasino

#endif // PLINKO_GAME_HPP
