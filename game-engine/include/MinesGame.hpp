#ifndef MINES_GAME_HPP
#define MINES_GAME_HPP

#include <string>
#include <vector>
#include <cstdint>
#include <set>
#include <memory>
#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

/**
 * Mines Game - Ultra Low Latency Implementation
 * Players reveal tiles avoiding mines
 */
class MinesGame {
public:
    enum class Status {
        Playing,
        Won,
        Lost
    };
    
    struct Config {
        uint8_t minesCount;  // 1-24
        uint8_t gridSize;    // Default 25 (5x5)
        
        Config() : minesCount(3), gridSize(25) {}
    };
    
    struct Bet {
        std::string betId;
        std::string playerId;
        double betAmount;
        uint8_t minesCount;
        std::vector<uint8_t> revealedTiles;
        uint8_t mineHit;
        bool cashedOut;
        double payout;
        uint64_t timestamp;
    };
    
    struct GameState {
        std::string gameId;
        Status status;
        Config config;
        std::vector<uint8_t> mines;
        std::vector<Bet> bets;
    };

private:
    static constexpr double HOUSE_EDGE = 0.05; // 5%
    static constexpr uint8_t TOTAL_TILES = 25;
    
    Config config_;
    GameState currentState_;
    std::unique_ptr<RandomNumberGenerator> rng_;

public:
    MinesGame() : rng_(std::make_unique<RandomNumberGenerator>()) {}
    
    void setConfig(uint8_t minesCount) {
        config_.minesCount = std::min(std::max(minesCount, (uint8_t)1), (uint8_t)24);
        config_.gridSize = TOTAL_TILES;
    }
    
    /**
     * Generate mine positions using provably fair seeds
     */
    std::vector<uint8_t> generateMines(const std::string& serverSeed,
                                        const std::string& clientSeed,
                                        uint64_t nonce) {
        std::set<uint8_t> mineSet;
        std::string seed = serverSeed;
        
        while (mineSet.size() < config_.minesCount) {
            uint64_t outcome = ProvablyFair::generateOutcome(seed, clientSeed, nonce + mineSet.size());
            uint8_t position = outcome % TOTAL_TILES;
            
            if (mineSet.find(position) == mineSet.end()) {
                mineSet.insert(position);
            }
            
            // Update seed for next iteration
            seed = ProvablyFair::hashSeed(seed);
        }
        
        return std::vector<uint8_t>(mineSet.begin(), mineSet.end());
    }
    
    /**
     * Start a new mines game
     */
    GameState startGame(const std::string& serverSeed,
                       const std::string& clientSeed,
                       uint64_t nonce,
                       double betAmount) {
        currentState_.gameId = generateGameId();
        currentState_.status = Status::Playing;
        currentState_.config = config_;
        currentState_.mines = generateMines(serverSeed, clientSeed, nonce);
        currentState_.bets.clear();
        
        return currentState_;
    }
    
    /**
     * Place a bet
     */
    Bet placeBet(const std::string& playerId, double betAmount) {
        Bet bet;
        bet.betId = generateBetId();
        bet.playerId = playerId;
        bet.betAmount = betAmount;
        bet.minesCount = config_.minesCount;
        bet.cashedOut = false;
        bet.payout = 0.0;
        bet.mineHit = 255;
        bet.timestamp = getCurrentTimestamp();
        
        currentState_.bets.push_back(bet);
        return bet;
    }
    
    /**
     * Reveal a tile - returns result
     */
    struct RevealResult {
        enum class Type {
            Success,
            MineHit,
            AlreadyRevealed,
            Invalid,
            GameOver
        } type;
        
        uint8_t tile;
        double multiplier;
        double potentialWin;
        double payout;
        std::vector<uint8_t> allMines;
    };
    
    RevealResult revealTile(const std::string& betId, uint8_t tile) {
        RevealResult result;
        result.type = RevealResult::Type::Invalid;
        result.tile = tile;
        result.multiplier = 1.0;
        result.potentialWin = 0.0;
        
        if (tile >= TOTAL_TILES) {
            result.type = RevealResult::Type::Invalid;
            return result;
        }
        
        // Find the bet
        for (auto& bet : currentState_.bets) {
            if (bet.betId == betId) {
                // Check if already revealed
                for (uint8_t revealed : bet.revealedTiles) {
                    if (revealed == tile) {
                        result.type = RevealResult::Type::AlreadyRevealed;
                        return result;
                    }
                }
                
                // Check if hit mine
                for (uint8_t mine : currentState_.mines) {
                    if (mine == tile) {
                        result.type = RevealResult::Type::MineHit;
                        result.allMines = currentState_.mines;
                        result.payout = 0.0;
                        bet.mineHit = tile;
                        currentState_.status = Status::Lost;
                        
                        // Mark all mines revealed
                        bet.revealedTiles = currentState_.mines;
                        return result;
                    }
                }
                
                // Success - tile is safe
                bet.revealedTiles.push_back(tile);
                
                // Calculate multiplier
                uint8_t revealedCount = bet.revealedTiles.size();
                uint8_t safeTiles = TOTAL_TILES - config_.minesCount;
                
                result.type = RevealResult::Type::Success;
                result.tile = tile;
                result.multiplier = calculateMultiplier(revealedCount, safeTiles);
                result.potentialWin = bet.betAmount * result.multiplier * (1.0 - HOUSE_EDGE);
                
                // Check if all safe tiles revealed
                if (revealedCount >= safeTiles) {
                    result.type = RevealResult::Type::GameOver;
                    result.payout = bet.betAmount * result.multiplier * (1.0 - HOUSE_EDGE);
                    bet.payout = result.payout;
                    bet.cashedOut = true;
                    currentState_.status = Status::Won;
                }
                
                return result;
            }
        }
        
        return result;
    }
    
    /**
     * Cash out current winnings
     */
    double cashout(const std::string& betId) {
        for (auto& bet : currentState_.bets) {
            if (bet.betId == betId && !bet.cashedOut) {
                uint8_t revealedCount = bet.revealedTiles.size();
                uint8_t safeTiles = TOTAL_TILES - config_.minesCount;
                
                double multiplier = calculateMultiplier(revealedCount, safeTiles);
                bet.payout = bet.betAmount * multiplier * (1.0 - HOUSE_EDGE);
                bet.cashedOut = true;
                
                return bet.payout;
            }
        }
        return 0.0;
    }
    
    /**
     * Calculate current multiplier
     */
    double calculateMultiplier(uint8_t revealedCount, uint8_t safeTiles) {
        if (revealedCount == 0) return 1.0;
        
        // Progressive multiplier based on revealed tiles
        double base = 1.0;
        double increment = (config_.minesCount * 0.3);
        increment = std::max(increment, 0.5);
        
        return base + (revealedCount * increment);
    }
    
    const GameState& getState() const { return currentState_; }

private:
    uint64_t getCurrentTimestamp() {
        auto now = std::chrono::system_clock::now();
        return std::chrono::duration_cast<std::chrono::milliseconds>(now.time_since_epoch()).count();
    }
    
    std::string generateGameId() {
        uint64_t id = rng_->generateSeed();
        return "MINES-" + std::to_string(id);
    }
    
    std::string generateBetId() {
        uint64_t id = rng_->generateSeed();
        return "BET-" + std::to_string(id);
    }
};

} // namespace TigerCasino

#endif // MINES_GAME_HPP
