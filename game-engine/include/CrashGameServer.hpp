#pragma once

#include <string>
#include <vector>
#include <map>
#include <memory>
#include <functional>
#include <mutex>
#include <atomic>
#include <chrono>
#include <cmath>
#include <random>
#include <openssl/rand.h>
#include <openssl/evp.h>

namespace TigerCasino {

// Crash game states
enum class CrashState {
    WAITING,
    FLYING,
    CRASHED
};

// Player bet in crash game
struct CrashPlayer {
    std::string playerId;
    double betAmount;
    double cashoutMultiplier;
    bool hasCashedOut;
    double winAmount;
};

// Crash game round
struct CrashRound {
    std::string roundId;
    uint64_t seed;
    double currentMultiplier;
    CrashState state;
    std::vector<CrashPlayer> players;
    std::chrono::steady_clock::time_point startTime;
    std::chrono::steady_clock::time_point crashTime;
    
    CrashRound() : currentMultiplier(1.0), state(CrashState::WAITING) {}
};

// Crash game server with provably fair outcomes
class CrashGameServer {
private:
    static constexpr double MIN_MULTIPLIER = 1.0;
    static constexpr double MAX_MULTIPLIER = 100.0;
    static constexpr double BASE_MULTIPLIER = 0.01;
    static constexpr int WAITING_TIME_MS = 5000;
    static constexpr int FLYING_INTERVAL_MS = 100;
    
    std::mt19937_64 rng_;
    std::mutex rngMutex_;
    
    std::string currentRoundId_;
    uint64_t currentSeed_;
    CrashRound currentRound_;
    std::atomic<bool> running_;
    std::atomic<double> currentMultiplier_;
    
    std::function<void(const std::string&)> onGameUpdate_;
    std::function<void(const CrashRound&)> onRoundEnd_;
    
    // Provably fair calculation
    double calculateCrashPoint(uint64_t seed);
    uint64_t generateSecureSeed();
    
public:
    CrashGameServer();
    ~CrashGameServer();
    
    // Game control
    void start();
    void stop();
    bool isRunning() const;
    
    // Round management
    std::string startNewRound();
    bool cashout(const std::string& playerId, double multiplier);
    CrashRound getCurrentRound() const;
    
    // Player management
    bool addPlayer(const std::string& playerId, double betAmount);
    bool removePlayer(const std::string& playerId);
    std::vector<CrashPlayer> getPlayers() const;
    
    // Multiplier updates
    void updateMultiplier(double multiplier);
    
    // Callbacks
    void setGameUpdateCallback(std::function<void(const std::string&)> callback);
    void setRoundEndCallback(std::function<void(const CrashRound&)> callback);
    
    // Statistics
    std::map<std::string, double> getGameStats() const;
};

// Hash for provably fair verification
std::string hashSHA256(const std::string& input);

// Generate server seed for provably fair games
std::string generateServerSeed();

} // namespace TigerCasino
