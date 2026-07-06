#pragma once

#include <string>
#include <vector>
#include <map>
#include <memory>

namespace TigerCasino {

// Game result structure
struct GameResult {
    bool success;
    double winAmount;
    double multiplier;
    std::string outcome;
    std::map<std::string, std::string> metadata;
    
    GameResult() : success(false), winAmount(0.0), multiplier(0.0) {}
};

// Bet information
struct BetInfo {
    double amount;
    double balance;
    std::map<std::string, double> gameParams;
};

// Base game class
class Game {
public:
    virtual ~Game() = default;
    virtual std::string getName() const = 0;
    virtual std::string getType() const = 0;
    virtual double getRTP() const = 0;
    virtual GameResult play(const BetInfo& bet) = 0;
};

// Game engine main class
class GameEngine {
private:
    std::map<std::string, std::shared_ptr<Game>> games_;
    
public:
    GameEngine();
    ~GameEngine() = default;
    
    void registerGame(std::shared_ptr<Game> game);
    GameResult play(const std::string& gameType, const BetInfo& bet);
    std::vector<std::string> getAvailableGames() const;
    double getGameRTP(const std::string& gameType) const;
};

} // namespace TigerCasino
