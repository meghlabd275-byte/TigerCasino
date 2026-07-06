#ifndef CRASH_GAME_HPP
#define CRASH_GAME_HPP

#include <string>
#include <vector>
#include <cstdint>
#include <memory>
#include <atomic>
#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

/**
 * Crash Game - Ultra Low Latency Implementation
 * Players watch multiplier rise and cash out before crash
 */
class CrashGame {
public:
    enum class Status {
        Waiting,
        Rising,
        Crashed
    };
    
    struct Bet {
        std::string betId;
        std::string playerId;
        double betAmount;
        double autoCashoutAt;
        bool cashedOut;
        double cashoutMultiplier;
        double payout;
        uint64_t timestamp;
    };
    
    struct GameState {
        std::string gameId;
        Status status;
        double currentMultiplier;
        double crashPoint;
        uint64_t startTime;
        std::vector<Bet> bets;
    };

private:
    static constexpr double HOUSE_EDGE = 0.03; // 3%
    static constexpr double MIN_MULTIPLIER = 1.0;
    static constexpr double MAX_MULTIPLIER = 100.0;
    
    std::atomic<uint64_t> currentGameId_{0};
    GameState currentState_;
    std::vector<double> history_;
    
public:
    CrashGame() : rng_(std::make_unique<RandomNumberGenerator>()) {
        history_.reserve(100);
    }
    
    /**
     * Generate crash point using provably fair seeds
     * Ultra-fast calculation for real-time gaming
     */
    double generateCrashPoint(const std::string& serverSeed,
                              const std::string& clientSeed,
                              uint64_t nonce) {
        uint64_t outcome = ProvablyFair::generateOutcome(serverSeed, clientSeed, nonce);
        
        // Exponential distribution for realistic crash points
        // Most crashes happen early, rare big multipliers
        uint64_t mod = outcome % 100000;
        
        if (mod < 35000) {
            // 35% - crash below 1.1x (instant crash)
            return 1.0 + (mod % 1000) / 10000.0;
        } else {
            // 65% - higher crash point
            double x = static_cast<double>(mod - 35000) / 65000.0;
            // Logarithmic distribution for natural feel
            double crash = 1.1 + (std::log(1.0 / (1.0 - x * 0.99)) * 2.0);
            return std::min(crash, MAX_MULTIPLIER);
        }
    }
    
    /**
     * Start a new crash round
     */
    GameState startRound(const std::string& serverSeed,
                         const std::string& clientSeed) {
        currentGameId_++;
        
        currentState_.gameId = std::to_string(currentGameId_.load());
        currentState_.status = Status::Rising;
        currentState_.currentMultiplier = 1.0;
        currentState_.crashPoint = generateCrashPoint(serverSeed, clientSeed, currentGameId_.load());
        currentState_.startTime = getCurrentTimestamp();
        currentState_.bets.clear();
        
        return currentState_;
    }
    
    /**
     * Place a bet in current round
     */
    Bet placeBet(const std::string& playerId,
                  double betAmount,
                  double autoCashoutAt = 0.0) {
        Bet bet;
        bet.betId = generateBetId();
        bet.playerId = playerId;
        bet.betAmount = betAmount;
        bet.autoCashoutAt = autoCashoutAt;
        bet.cashedOut = false;
        bet.cashoutMultiplier = 0.0;
        bet.payout = 0.0;
        bet.timestamp = getCurrentTimestamp();
        
        currentState_.bets.push_back(bet);
        return bet;
    }
    
    /**
     * Process cashout for a player
     */
    double processCashout(const std::string& betId, double currentMultiplier) {
        for (auto& bet : currentState_.bets) {
            if (bet.betId == betId && !bet.cashedOut) {
                bet.cashedOut = true;
                bet.cashoutMultiplier = currentMultiplier;
                bet.payout = bet.betAmount * currentMultiplier * (1.0 - HOUSE_EDGE);
                return bet.payout;
            }
        }
        return 0.0;
    }
    
    /**
     * Check and process auto-cashouts
     */
    std::vector<Bet> checkAutoCashouts(double currentMultiplier) {
        std::vector<Bet> cashedOut;
        
        for (auto& bet : currentState_.bets) {
            if (!bet.cashedOut && bet.autoCashoutAt > 0.0 && 
                currentMultiplier >= bet.autoCashoutAt) {
                bet.cashedOut = true;
                bet.cashoutMultiplier = currentMultiplier;
                bet.payout = bet.betAmount * currentMultiplier * (1.0 - HOUSE_EDGE);
                cashedOut.push_back(bet);
            }
        }
        
        return cashedOut;
    }
    
    /**
     * End the current round (crash)
     */
    void endRound() {
        currentState_.status = Status::Crashed;
        
        // Record to history
        history_.insert(history_.begin(), currentState_.crashPoint);
        if (history_.size() > 100) {
            history_.pop_back();
        }
    }
    
    /**
     * Get current game state
     */
    const GameState& getState() const { return currentState_; }
    
    /**
     * Get crash history
     */
    const std::vector<double>& getHistory() const { return history_; }
    
    /**
     * Get current multiplier with auto-increment for rising state
     */
    double getCurrentMultiplier(uint64_t elapsedMs) const {
        if (currentState_.status != Status::Rising) {
            return currentState_.crashPoint;
        }
        
        // Multiplier grows exponentially: 1x at start, +0.03% per ms
        double multiplier = 1.0 + (elapsedMs * 0.0003);
        return std::min(multiplier, currentState_.crashPoint);
    }

private:
    std::unique_ptr<RandomNumberGenerator> rng_;
    
    uint64_t getCurrentTimestamp() {
        auto now = std::chrono::system_clock::now();
        auto duration = now.time_since_epoch();
        return std::chrono::duration_cast<std::chrono::milliseconds>(duration).count();
    }
    
    std::string generateBetId() {
        uint64_t id = rng_->generateSeed();
        return "BET-" + std::to_string(id);
    }
};

} // namespace TigerCasino

#endif // CRASH_GAME_HPP
