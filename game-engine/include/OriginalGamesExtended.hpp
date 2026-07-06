#pragma once

#include <string>
#include <vector>
#include <map>
#include <memory>
#include <random>
#include <mutex>
#include <functional>

namespace TigerCasino {

// Original game types
enum class OriginalGameType {
    CRASH,
    PLINKO,
    MINES,
    LIMBO,
    DICE,
    KENO,
    BINGO,
    SCRATCH,
    ROULETTE,
    DIAMOND,
    HASH,
    COINFLIP,
    ROCK_PAPER_SCISSORS,
    WHEEL,
    TOMB,
    AVALANCHE,
    CAVE,
    TOWER,
    BATTLE,
    CUP,
    RACE
};

// Original game configuration
struct OriginalGameConfig {
    std::string gameId;
    std::string gameName;
    OriginalGameType type;
    double minBet;
    double maxBet;
    double rtp;
    bool hasAutoCashout;
    std::string description;
};

// Game result for original games
struct OriginalGameResult {
    bool success;
    double betAmount;
    double winAmount;
    double multiplier;
    std::string outcome;
    std::string hash;
    std::map<std::string, std::string> metadata;
    
    OriginalGameResult() : success(false), betAmount(0), winAmount(0), multiplier(0) {}
};

// Base original game server
class OriginalGameServer {
protected:
    std::mt19937_64 rng_;
    std::mutex rngMutex_;
    
public:
    OriginalGameServer();
    virtual ~OriginalGameServer() = default;
    
    virtual OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) = 0;
    virtual OriginalGameConfig getConfig() const = 0;
};

// Plinko game
class PlinkoGame : public OriginalGameServer {
private:
    int rows_;
    double riskMultiplier_;
    std::vector<std::vector<double>> payouts_;
    
public:
    PlinkoGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Mines game
class MinesGame : public OriginalGameServer {
private:
    int mines_;
    int gridSize_;
    
public:
    MinesGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
    
    std::vector<std::string> generateMines(int count);
    double calculatePayout(int gemsFound, double bet);
};

// Limbo game
class LimboGame : public OriginalGameServer {
private:
    double targetMultiplier_;
    
public:
    LimboGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Keno game
class KenoGame : public OriginalGameServer {
private:
    int numbersPick_;
    int numbersDrawn_;
    
public:
    KenoGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Bingo game
class BingoGame : public OriginalGameServer {
private:
    int ballCount_;
    int cardSize_;
    
public:
    BingoGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Scratch card game
class ScratchGame : public OriginalGameServer {
private:
    int symbolsToReveal_;
    int winningSymbols_;
    
public:
    ScratchGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Diamond game (diamond breaker)
class DiamondGame : public OriginalGameServer {
private:
    int rows_;
    int cols_;
    
public:
    DiamondGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Hash dice game
class HashDiceGame : public OriginalGameServer {
private:
    double minWin_;
    double maxWin_;
    
public:
    HashDiceGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Coin flip game
class CoinFlipGame : public OriginalGameServer {
public:
    CoinFlipGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Rock Paper Scissors
class RockPaperScissorsGame : public OriginalGameServer {
public:
    RockPaperScissorsGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Fortune wheel game
class FortuneWheelGame : public OriginalGameServer {
private:
    std::vector<std::pair<std::string, double>> segments_;
    
public:
    FortuneWheelGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Tomb game (escape tomb)
class TombGame : public OriginalGameServer {
private:
    int levels_;
    double levelMultiplier_;
    
public:
    TombGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Avalanche game
class AvalancheGame : public OriginalGameServer {
private:
    int gridSize_;
    double cascadeMultiplier_;
    
public:
    AvalancheGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Tower game
class TowerGame : public OriginalGameServer {
private:
    int floors_;
    double floorMultiplier_;
    
public:
    TowerGame();
    OriginalGameResult play(double betAmount, const std::map<std::string, double>& params) override;
    OriginalGameConfig getConfig() const override;
};

// Complete original games manager
class OriginalGamesManager {
private:
    std::map<std::string, std::shared_ptr<OriginalGameServer>> games_;
    
public:
    OriginalGamesManager();
    
    void registerGame(const std::string& gameId, std::shared_ptr<OriginalGameServer> game);
    
    OriginalGameResult play(const std::string& gameId, double betAmount, const std::map<std::string, double>& params);
    
    OriginalGameConfig getConfig(const std::string& gameId) const;
    
    std::vector<OriginalGameConfig> getAllConfigs() const;
    
    std::vector<std::string> getGameIds() const;
    
    size_t getGameCount() const;
};

} // namespace TigerCasino
