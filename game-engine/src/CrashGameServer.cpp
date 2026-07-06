#include "CrashGameServer.hpp"
#include <sstream>
#include <iomanip>
#include <thread>
#include <chrono>

namespace TigerCasino {

// Crash point calculation using hash
double CrashGameServer::calculateCrashPoint(uint64_t seed) {
    // Convert seed to bytes
    uint8_t bytes[8];
    for (int i = 0; i < 8; i++) {
        bytes[i] = (seed >> (i * 8)) & 0xFF;
    }
    
    // Use OpenSSL for secure hash
    unsigned char hash[EVP_MAX_MD_SIZE];
    unsigned int hashLen;
    EVP_MD_CTX* ctx = EVP_MD_CTX_new();
    
    if (ctx) {
        EVP_DigestInit_ex(ctx, EVP_sha256(), nullptr);
        EVP_DigestUpdate(ctx, bytes, 8);
        EVP_DigestFinal_ex(ctx, hash, &hashLen);
        EVP_MD_CTX_free(ctx);
    }
    
    // Convert hash to number
    uint64_t hashNum = 0;
    for (unsigned int i = 0; i < std::min(hashLen, (unsigned int)8); i++) {
        hashNum = (hashNum << 8) | hash[i];
    }
    
    // Apply exponential distribution for crash point
    double result = BASE_MULTIPLIER + (static_cast<double>(hashNum % 1000000) / 1000000.0);
    result = std::pow(result, 0.5);  // Square root for better distribution
    
    // Scale to desired range
    result = MIN_MULTIPLIER + (result * (MAX_MULTIPLIER - MIN_MULTIPLIER));
    
    // Cap at max multiplier
    return std::min(result, MAX_MULTIPLIER);
}

uint64_t CrashGameServer::generateSecureSeed() {
    std::lock_guard<std::mutex> lock(rngMutex_);
    uint8_t seedBytes[8];
    RAND_bytes(seedBytes, 8);
    
    uint64_t seed = 0;
    for (int i = 0; i < 8; i++) {
        seed = (seed << 8) | seedBytes[i];
    }
    return seed;
}

CrashGameServer::CrashGameServer() 
    : running_(false)
    , currentMultiplier_(1.0) {
    // Initialize RNG with secure seed
    std::random_device rd;
    rng_.seed(rd());
}

CrashGameServer::~CrashGameServer() {
    stop();
}

void CrashGameServer::start() {
    if (running_) return;
    
    running_ = true;
    startNewRound();
}

void CrashGameServer::stop() {
    running_ = false;
}

bool CrashGameServer::isRunning() const {
    return running_;
}

std::string CrashGameServer::startNewRound() {
    currentSeed_ = generateSecureSeed();
    
    // Generate round ID
    std::stringstream ss;
    ss << "crash_" << std::chrono::steady_clock::now().time_since_epoch().count();
    currentRoundId_ = ss.str();
    
    // Reset round
    currentRound_.roundId = currentRoundId_;
    currentRound_.seed = currentSeed_;
    currentRound_.currentMultiplier = MIN_MULTIPLIER;
    currentRound_.state = CrashState::WAITING;
    currentRound_.startTime = std::chrono::steady_clock::now();
    currentMultiplier_ = MIN_MULTIPLIER;
    
    // Calculate crash point in advance (provably fair)
    double crashPoint = calculateCrashPoint(currentSeed_);
    currentRound_.crashTime = currentRound_.startTime + 
        std::chrono::milliseconds(static_cast<int>(crashPoint * 1000));
    
    return currentRoundId_;
}

bool CrashGameServer::cashout(const std::string& playerId, double multiplier) {
    if (currentRound_.state != CrashState::FLYING) {
        return false;
    }
    
    for (auto& player : currentRound_.players) {
        if (player.playerId == playerId && !player.hasCashedOut) {
            player.cashoutMultiplier = multiplier;
            player.hasCashedOut = true;
            player.winAmount = player.betAmount * multiplier;
            return true;
        }
    }
    return false;
}

CrashRound CrashGameServer::getCurrentRound() const {
    return currentRound_;
}

bool CrashGameServer::addPlayer(const std::string& playerId, double betAmount) {
    for (const auto& player : currentRound_.players) {
        if (player.playerId == playerId) {
            return false;  // Player already in game
        }
    }
    
    CrashPlayer player;
    player.playerId = playerId;
    player.betAmount = betAmount;
    player.cashoutMultiplier = 0;
    player.hasCashedOut = false;
    player.winAmount = 0;
    
    currentRound_.players.push_back(player);
    return true;
}

bool CrashGameServer::removePlayer(const std::string& playerId) {
    for (auto it = currentRound_.players.begin(); it != currentRound_.players.end(); ++it) {
        if (it->playerId == playerId) {
            currentRound_.players.erase(it);
            return true;
        }
    }
    return false;
}

void CrashGameServer::updateMultiplier(double multiplier) {
    currentMultiplier_ = multiplier;
    currentRound_.currentMultiplier = multiplier;
    
    if (multiplier >= MAX_MULTIPLIER) {
        currentRound_.state = CrashState::CRASHED;
    }
}

void CrashGameServer::setGameUpdateCallback(std::function<void(const std::string&)> callback) {
    onGameUpdate_ = callback;
}

void CrashGameServer::setRoundEndCallback(std::function<void(const CrashRound&)> callback) {
    onRoundEnd_ = callback;
}

std::map<std::string, double> CrashGameServer::getGameStats() const {
    std::map<std::string, double> stats;
    
    stats["currentMultiplier"] = currentMultiplier_.load();
    stats["playerCount"] = currentRound_.players.size();
    stats["totalBets"] = 0;
    
    for (const auto& player : currentRound_.players) {
        stats["totalBets"] += player.betAmount;
    }
    
    return stats;
}

// SHA256 hash function
std::string hashSHA256(const std::string& input) {
    unsigned char hash[EVP_MAX_MD_SIZE];
    unsigned int hashLen;
    
    EVP_MD_CTX* ctx = EVP_MD_CTX_new();
    if (!ctx) return "";
    
    EVP_DigestInit_ex(ctx, EVP_sha256(), nullptr);
    EVP_DigestUpdate(ctx, input.c_str(), input.length());
    EVP_DigestFinal_ex(ctx, hash, &hashLen);
    EVP_MD_CTX_free(ctx);
    
    std::stringstream ss;
    for (unsigned int i = 0; i < hashLen; i++) {
        ss << std::hex << std::setw(2) << std::setfill('0') << (int)hash[i];
    }
    return ss.str();
}

// Generate server seed
std::string generateServerSeed() {
    uint8_t seedBytes[32];
    RAND_bytes(seedBytes, 32);
    
    std::stringstream ss;
    for (int i = 0; i < 32; i++) {
        ss << std::hex << std::setw(2) << std::setfill('0') << (int)seedBytes[i];
    }
    return ss.str();
}

} // namespace TigerCasino
