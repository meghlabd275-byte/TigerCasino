#pragma once

#include <string>
#include <vector>
#include <map>
#include <memory>
#include <functional>
#include <mutex>
#include <atomic>

namespace TigerCasino {

// Live dealer game types
enum class LiveGameType {
    BLACKJACK,
    ROULETTE,
    BACCARAT,
    POKER,
    GAME_SHOW,
    DRAGON_TIGER,
    SIC_BO,
    DREAM_CATCHER
};

// Live dealer table status
enum class TableStatus {
    WAITING,
    DEALING,
    BETTING_CLOSED,
    RESULT,
    SETTLING
};

// Player at a live table
struct LivePlayer {
    std::string playerId;
    double balance;
    std::map<std::string, double> bets;  // betType -> amount
    std::map<std::string, double> wins;
    bool isConnected;
};

// Live dealer table
struct LiveTable {
    std::string tableId;
    std::string tableName;
    std::string provider;  // evolution, pragmatiplay, etc.
    LiveGameType gameType;
    TableStatus status;
    std::vector<LivePlayer> players;
    int minBet;
    int maxBet;
    std::string dealerName;
    std::string streamUrl;
    std::map<std::string, std::string> metadata;
    
    LiveTable() : status(TableStatus::WAITING), minBet(1), maxBet(10000) {}
};

// Game round data
struct LiveRound {
    std::string roundId;
    std::string tableId;
    std::string dealerName;
    std::vector<std::string> result;
    std::chrono::steady_clock::time_point startTime;
    std::chrono::steady_clock::time_point endTime;
    std::map<std::string, std::string> gameData;
    
    LiveRound() {}
};

// Live dealer server - interfaces with providers
class LiveDealerServer {
private:
    std::map<std::string, LiveTable> tables_;
    std::map<std::string, LiveRound> activeRounds_;
    std::mutex tablesMutex_;
    std::atomic<bool> running_;
    
    // Provider integration
    std::map<std::string, std::function<std::string(const std::string&)>> providers_;
    
    // Callbacks
    std::function<void(const LiveTable&)> onTableUpdate_;
    std::function<void(const LiveRound&)> onRoundComplete_;
    
public:
    LiveDealerServer();
    ~LiveDealerServer();
    
    // Table management
    void addTable(const LiveTable& table);
    void removeTable(const std::string& tableId);
    LiveTable getTable(const std::string& tableId) const;
    std::vector<LiveTable> getTablesByGame(LiveGameType type) const;
    std::vector<LiveTable> getAllTables() const;
    
    // Player management
    bool addPlayerToTable(const std::string& tableId, const LivePlayer& player);
    bool removePlayerFromTable(const std::string& tableId, const std::string& playerId);
    bool placeBet(const std::string& tableId, const std::string& playerId,
                  const std::string& betType, double amount);
    
    // Round management
    LiveRound startRound(const std::string& tableId);
    void endRound(const std::string& roundId, const std::vector<std::string>& result);
    LiveRound getActiveRound(const std::string& tableId) const;
    
    // Provider integration
    void registerProvider(const std::string& providerName,
                         std::function<std::string(const std::string&)> handler);
    
    // Callbacks
    void setTableUpdateCallback(std::function<void(const LiveTable&)> callback);
    void setRoundCompleteCallback(std::function<void(const LiveRound&)> callback);
    
    // Control
    void start();
    void stop();
    bool isRunning() const;
};

// Live dealer provider integration
namespace LiveProviders {

struct ProviderConfig {
    std::string name;
    std::string apiEndpoint;
    std::string apiKey;
    bool isActive;
    std::map<std::string, std::string> settings;
};

// Create tables for major providers
std::vector<LiveTable> createEvolutionTables();
std::vector<LiveTable> createPragmaticPlayTables();
std::vector<LiveTable> createEzugiTables();
std::vector<LiveTable> createAuthenticGamingTables();

} // namespace LiveProviders

} // namespace TigerCasino
