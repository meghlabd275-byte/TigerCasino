#pragma once

#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"
#include <string>
#include <vector>
#include <memory>
#include <array>
#include <map>

namespace TigerCasino {

/**
 * Base class for game show games
 */
class GameShowGame {
public:
    virtual ~GameShowGame() = default;
    virtual std::string getName() const = 0;
    virtual std::string getType() const = 0;
    virtual double getRTP() const = 0;
};

/**
 * Crazy Time - Multi-bonus wheel game show
 */
class CrazyTimeGame : public GameShowGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 5000.00;
    static constexpr double BASE_RTP = 0.96; // 96% RTP

    // Wheel segments
    enum class SegmentType {
        NUMBER_1,    // 1x
        NUMBER_2,    // 2x
        NUMBER_5,   // 5x
        NUMBER_10,   // 10x
        COIN_FLIP,   // Bonus game
        CASH_HUNT,  // Bonus game
        PACHINKO,   // Bonus game
        CRAZY_TIME  // Bonus game
    };

    struct WheelResult {
        SegmentType segment;
        uint8_t segment_index;
        double multiplier;
        std::string bonus_game_result;
        double total_win;
        std::string server_seed;
    };

    struct Bet {
        std::string player_id;
        std::map<SegmentType, double> bets; // Bet on specific segment
        std::chrono::steady_clock::time_point bet_time;
    };

private:
    static constexpr std::array<std::pair<SegmentType, uint8_t>, 64> WHEEL_SEGMENTS = {{
        {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1},
        {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1},
        {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1}, {SegmentType::NUMBER_1, 1},
        {SegmentType::NUMBER_2, 2}, {SegmentType::NUMBER_2, 2}, {SegmentType::NUMBER_2, 2}, {SegmentType::NUMBER_2, 2},
        {SegmentType::NUMBER_2, 2}, {SegmentType::NUMBER_2, 2}, {SegmentType::NUMBER_2, 2}, {SegmentType::NUMBER_2, 2},
        {SegmentType::NUMBER_5, 5}, {SegmentType::NUMBER_5, 5}, {SegmentType::NUMBER_5, 5}, {SegmentType::NUMBER_5, 5},
        {SegmentType::NUMBER_5, 5}, {SegmentType::NUMBER_5, 5}, {SegmentType::NUMBER_5, 5}, {SegmentType::NUMBER_5, 5},
        {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10},
        {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10},
        {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10}, {SegmentType::NUMBER_10, 10},
        {SegmentType::COIN_FLIP, 1}, {SegmentType::COIN_FLIP, 1}, {SegmentType::COIN_FLIP, 1}, {SegmentType::COIN_FLIP, 1},
        {SegmentType::CASH_HUNT, 1}, {SegmentType::CASH_HUNT, 1}, {SegmentType::CASH_HUNT, 1}, {SegmentType::CASH_HUNT, 1},
        {SegmentType::PACHINKO, 1}, {SegmentType::PACHINKO, 1}, {SegmentType::PACHINKO, 1}, {SegmentType::PACHINKO, 1},
        {SegmentType::CRAZY_TIME, 1}, {SegmentType::CRAZY_TIME, 1}, {SegmentType::CRAZY_TIME, 1}, {SegmentType::CRAZY_TIME, 1},
    }};

    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t round_counter_;

public:
    CrazyTimeGame();
    
    std::string getName() const override { return "Crazy Time"; }
    std::string getType() const override { return "Game Show"; }
    double getRTP() const override { return BASE_RTP; }

    WheelResult spin();
    WheelResult playBonusGame(SegmentType bonus_type);
    
    std::string segmentToString(SegmentType type) const;
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);
};

inline CrazyTimeGame::CrazyTimeGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {}

inline std::string CrazyTimeGame::segmentToString(SegmentType type) const {
    switch (type) {
        case SegmentType::NUMBER_1: return "1";
        case SegmentType::NUMBER_2: return "2";
        case SegmentType::NUMBER_5: return "5";
        case SegmentType::NUMBER_10: return "10";
        case SegmentType::COIN_FLIP: return "COIN_FLIP";
        case SegmentType::CASH_HUNT: return "CASH_HUNT";
        case SegmentType::PACHINKO: return "PACHINKO";
        case SegmentType::CRAZY_TIME: return "CRAZY_TIME";
        default: return "UNKNOWN";
    }
}

inline CrazyTimeGame::WheelResult CrazyTimeGame::spin() {
    WheelResult result;
    
    round_counter_++;
    
    // Use provably fair to determine spin result
    auto hash = provably_fair_->generateRandomHash();
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash.hash[i]);
    }
    
    // Map to wheel segment
    size_t segment_index = hash_value % WHEEL_SEGMENTS.size();
    auto [segment_type, _] = WHEEL_SEGMENTS[segment_index];
    
    result.segment = segment_type;
    result.segment_index = static_cast<uint8_t>(segment_index);
    result.server_seed = provably_fair_->getServerSeed();
    
    // Calculate multiplier based on segment
    switch (segment_type) {
        case SegmentType::NUMBER_1: result.multiplier = 1.0; break;
        case SegmentType::NUMBER_2: result.multiplier = 2.0; break;
        case SegmentType::NUMBER_5: result.multiplier = 5.0; break;
        case SegmentType::NUMBER_10: result.multiplier = 10.0; break;
        default: result.multiplier = 1.0; break;
    }
    
    // If bonus game, play bonus
    if (segment_type == SegmentType::COIN_FLIP || 
        segment_type == SegmentType::CASH_HUNT ||
        segment_type == SegmentType::PACHINKO ||
        segment_type == SegmentType::CRAZY_TIME) {
        auto bonus_result = playBonusGame(segment_type);
        result.bonus_game_result = bonus_result.bonus_game_result;
        result.total_win = bonus_result.total_win;
    }
    
    return result;
}

inline CrazyTimeGame::WheelResult CrazyTimeGame::playBonusGame(SegmentType bonus_type) {
    WheelResult result;
    result.segment = bonus_type;
    
    switch (bonus_type) {
        case SegmentType::COIN_FLIP: {
            // Simple coin flip with multipliers
            auto flip_result = provably_fair_->generateRandomHash();
            bool heads = flip_result.hash[0] % 2 == 0;
            double multiplier = heads ? rng_->generateDouble(1.0, 50.0) : rng_->generateDouble(1.0, 50.0);
            result.bonus_game_result = heads ? "HEADS" : "TAILS";
            result.total_win = multiplier;
            result.multiplier = multiplier;
            break;
        }
        case SegmentType::CASH_HUNT: {
            // Shooting gallery with hidden multipliers
            result.bonus_game_result = "CASH_HUNT";
            result.total_win = rng_->generateDouble(1.0, 100.0);
            result.multiplier = result.total_win;
            break;
        }
        case SegmentType::PACHINKO: {
            // Pachinko-style pinball
            result.bonus_game_result = "PACHINKO";
            result.total_win = rng_->generateDouble(1.0, 100.0);
            result.multiplier = result.total_win;
            break;
        }
        case SegmentType::CRAZY_TIME: {
            // Big wheel with huge multipliers
            result.bonus_game_result = "CRAZY_TIME";
            result.total_win = rng_->generateDouble(1.0, 500.0);
            result.multiplier = result.total_win;
            break;
        }
        default:
            break;
    }
    
    return result;
}

inline void CrazyTimeGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void CrazyTimeGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Monopoly Live - Monopoly-themed game show
 */
class MonopolyLiveGame : public GameShowGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 5000.00;
    static constexpr double BASE_RTP = 0.96;

    struct WheelResult {
        uint8_t segment_index;
        std::string segment_name;
        double multiplier;
        std::string chance_result;
        double total_win;
        std::string server_seed;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t round_counter_;

    static constexpr std::array<std::pair<std::string, double>, 54> WHEEL = {{
        {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0},
        {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0},
        {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0},
        {"10", 10.0}, {"10", 10.0}, {"10", 10.0}, {"10", 10.0},
        {"CHANCE", 0.0}, {"CHANCE", 0.0}, {"CHANCE", 0.0}, {"CHANCE", 0.0},
        {"2 ROLLS", 0.0}, {"2 ROLLS", 0.0}, {"2 ROLLS", 0.0}, {"2 ROLLS", 0.0},
        {"4 ROLLS", 0.0}, {"4 ROLLS", 0.0}, {"4 ROLLS", 0.0}, {"4 ROLLS", 0.0},
    }};

public:
    MonopolyLiveGame();
    
    std::string getName() const override { return "Monopoly Live"; }
    std::string getType() const override { return "Game Show"; }
    double getRTP() const override { return BASE_RTP; }

    WheelResult spin();
    WheelResult playBonus(uint8_t rolls);
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);
};

inline MonopolyLiveGame::MonopolyLiveGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {}

inline MonopolyLiveGame::WheelResult MonopolyLiveGame::spin() {
    WheelResult result;
    round_counter_++;
    
    auto hash = provably_fair_->generateRandomHash();
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash.hash[i]);
    }
    
    size_t segment_index = hash_value % WHEEL.size();
    auto [name, multiplier] = WHEEL[segment_index];
    
    result.segment_index = static_cast<uint8_t>(segment_index);
    result.segment_name = name;
    result.multiplier = multiplier;
    result.server_seed = provably_fair_->getServerSeed();
    
    if (name == "CHANCE") {
        result.chance_result = std::to_string(rng_->generateDouble(1.0, 10.0));
        result.total_win = std::stod(result.chance_result);
    } else if (name == "2 ROLLS" || name == "4 ROLLS") {
        uint8_t rolls = (name == "2 ROLLS") ? 2 : 4;
        auto bonus_result = playBonus(rolls);
        result.total_win = bonus_result.total_win;
        result.chance_result = name;
    }
    
    return result;
}

inline MonopolyLiveGame::WheelResult MonopolyLiveGame::playBonus(uint8_t rolls) {
    WheelResult result;
    result.total_win = rng_->generateDouble(1.0, 50.0) * rolls;
    result.multiplier = result.total_win;
    result.chance_result = std::to_string(rolls) + " ROLLS";
    return result;
}

inline void MonopolyLiveGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void MonopolyLiveGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Lightning Roulette - Fast-paced roulette with lightning strikes
 */
class LightningRouletteGame : public GameShowGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 5000.00;
    static constexpr double BASE_RTP = 0.973; // 97.3% RTP

    struct RouletteResult {
        uint8_t number; // 0-37 (American roulette)
        std::string color; // red/black
        std::string parity; // even/odd
        std::string range; // 1-18 / 19-36
        std::vector<uint8_t> lightning_numbers;
        std::vector<double> lightning_multipliers;
        double total_win;
        std::string server_seed;
    };

private:
    static constexpr std::array<std::string, 38> NUMBER_COLORS = {{
        "green", // 0
        "red", "black", "red", "black", "red", "black", "red", "black", "red", "black", "red", // 1-12
        "black", "red", "black", "red", "black", "red", "black", "red", "black", "red", "black", "red", // 13-24
        "red", "black", "red", "black", "red", "black", "red", "black", "red", "black", "red", "black", // 25-36
        "green" // 00
    }};

    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t round_counter_;

public:
    LightningRouletteGame();
    
    std::string getName() const override { return "Lightning Roulette"; }
    std::string getType() const override { return "Game Show"; }
    double getRTP() const override { return BASE_RTP; }

    RouletteResult spin();
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);
};

inline LightningRouletteGame::LightningRouletteGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {}

inline LightningRouletteGame::RouletteResult LightningRouletteGame::spin() {
    RouletteResult result;
    round_counter_++;
    
    // Generate lightning numbers (1-5 numbers with 50x-500x multipliers)
    size_t num_lightning = rng_->generateInt(1, 5);
    for (size_t i = 0; i < num_lightning; ++i) {
        result.lightning_numbers.push_back(static_cast<uint8_t>(rng_->generateInt(1, 36)));
        // Higher chance of lower multipliers
        result.lightning_multipliers.push_back(rng_->generateDouble(50.0, 500.0));
    }
    
    // Spin the wheel
    auto hash = provably_fair_->generateRandomHash();
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash.hash[i]);
    }
    
    result.number = static_cast<uint8_t>(hash_value % 38);
    result.color = NUMBER_COLORS[result.number];
    result.parity = (result.number % 2 == 0) ? "even" : "odd";
    result.range = (result.number >= 1 && result.number <= 18) ? "1-18" : "19-36";
    result.server_seed = provably_fair_->getServerSeed();
    
    // Check if number hit lightning
    for (size_t i = 0; i < result.lightning_numbers.size(); ++i) {
        if (result.lightning_numbers[i] == result.number) {
            result.total_win = result.lightning_multipliers[i];
            return result;
        }
    }
    
    result.total_win = 0.0;
    return result;
}

inline void LightningRouletteGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void LightningRouletteGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Dream Catcher - Money wheel game show
 */
class DreamCatcherGame : public GameShowGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 5000.00;
    static constexpr double BASE_RTP = 0.965;

    struct WheelResult {
        uint8_t segment_index;
        std::string segment_value;
        double multiplier;
        std::string server_seed;
    };

private:
    static constexpr std::array<std::pair<std::string, double>, 54> WHEEL = {{
        {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0}, {"1", 1.0},
        {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0}, {"2", 2.0},
        {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0}, {"5", 5.0},
        {"10", 10.0}, {"10", 10.0}, {"10", 10.0}, {"10", 10.0}, {"10", 10.0}, {"10", 10.0},
        {"20", 20.0}, {"20", 20.0}, {"20", 20.0}, {"20", 20.0},
        {"40", 40.0}, {"40", 40.0}, {"40", 40.0}, {"40", 40.0},
    }};

    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t round_counter_;

public:
    DreamCatcherGame();
    
    std::string getName() const override { return "Dream Catcher"; }
    std::string getType() const override { return "Game Show"; }
    double getRTP() const override { return BASE_RTP; }

    WheelResult spin();
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);
};

inline DreamCatcherGame::DreamCatcherGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {}

inline DreamCatcherGame::WheelResult DreamCatcherGame::spin() {
    WheelResult result;
    round_counter_++;
    
    auto hash = provably_fair_->generateRandomHash();
    uint64_t hash_value = 0;
    for (size_t i = 0; i < std::min(hash.hash.size(), sizeof(uint64_t)); ++i) {
        hash_value = (hash_value << 8) | static_cast<uint8_t>(hash.hash[i]);
    }
    
    size_t segment_index = hash_value % WHEEL.size();
    auto [value, multiplier] = WHEEL[segment_index];
    
    result.segment_index = static_cast<uint8_t>(segment_index);
    result.segment_value = value;
    result.multiplier = multiplier;
    result.server_seed = provably_fair_->getServerSeed();
    
    return result;
}

inline void DreamCatcherGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void DreamCatcherGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
