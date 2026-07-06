// Game Shows Implementation
#include "GameShows.hpp"
#include <iostream>
#include <random>
#include <sstream>

namespace TigerCasino {

// Crazy Time Implementation
CrazyTimeGame::CrazyTimeGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
}

std::string CrazyTimeGame::segmentToString(SegmentType type) const {
    switch (type) {
        case SegmentType::NUMBER_1: return "1";
        case SegmentType::NUMBER_2: return "2";
        case SegmentType::NUMBER_5: return "5";
        case SegmentType::NUMBER_10: return "10";
        case SegmentType::COIN_FLIP: return "COIN FLIP";
        case SegmentType::CASH_HUNT: return "CASH HUNT";
        case SegmentType::PACHINKO: return "PACHINKO";
        case SegmentType::CRAZY_TIME: return "CRAZY TIME";
        default: return "UNKNOWN";
    }
}

CrazyTimeGame::WheelResult CrazyTimeGame::spin() {
    WheelResult result;
    
    // Use provably fair random
    std::string combined = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(combined);
    uint64_t index = hash % 64;
    
    auto segment = WHEEL_SEGMENTS[index];
    result.segment = segment.first;
    result.segment_index = static_cast<uint8_t>(index);
    result.multiplier = static_cast<double>(segment.second);
    
    // Generate result
    if (result.multiplier > 0) {
        result.total_win = result.multiplier;
        result.bonus_game_result = "";
    } else {
        // Play bonus game
        result = playBonusGame(result.segment);
    }
    
    result.server_seed = provably_fair_->getServerSeed();
    
    return result;
}

CrazyTimeGame::WheelResult CrazyTimeGame::playBonusGame(SegmentType bonus_type) {
    WheelResult result;
    result.segment = bonus_type;
    result.multiplier = 0;
    
    switch (bonus_type) {
        case SegmentType::COIN_FLIP: {
            // Simple coin flip with 2 multipliers
            double multipliers[] = {10.0, 10.0};
            std::string combined = provably_fair_->getServerSeed() + "COIN";
            uint64_t hash = provably_fair_->generateHash(combined);
            size_t idx = hash % 2;
            result.total_win = multipliers[idx];
            result.bonus_game_result = idx == 0 ? "HEADS" : "TAILS";
            break;
        }
        case SegmentType::CASH_HUNT: {
            // Cash hunt with multiple targets
            double total = 0;
            for (int i = 0; i < 8; i++) {
                std::string combined = provably_fair_->getServerSeed() + "HUNT" + std::to_string(i);
                uint64_t hash = provably_fair_->generateHash(combined);
                double value = (hash % 500 + 50) / 10.0;
                total += value;
            }
            result.total_win = total;
            result.bonus_game_result = "CASH HUNT COMPLETED";
            break;
        }
        case SegmentType::PACHINKO: {
            // Pachinko with falling ball
            double total = 0;
            for (int i = 0; i < 12; i++) {
                std::string combined = provably_fair_->getServerSeed() + "PACHINKO" + std::to_string(i);
                uint64_t hash = provably_fair_->generateHash(combined);
                double value = (hash % 500 + 10) / 10.0;
                total += value;
            }
            result.total_win = total;
            result.bonus_game_result = "PACHINKO COMPLETED";
            break;
        }
        case SegmentType::CRAZY_TIME: {
            // Multiple spins of the big wheel
            double total = 0;
            for (int i = 0; i < 3; i++) {
                std::string combined = provably_fair_->getServerSeed() + "CRAZY" + std::to_string(i);
                uint64_t hash = provably_fair_->generateHash(combined);
                double value = (hash % 1000 + 100) / 10.0;
                total += value;
            }
            result.total_win = total;
            result.bonus_game_result = "CRAZY TIME COMPLETED";
            break;
        }
        default:
            result.total_win = 0;
            result.bonus_game_result = "NO BONUS";
    }
    
    return result;
}

void CrazyTimeGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

void CrazyTimeGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

// Monopoly Live Implementation
MonopolyLiveGame::MonopolyLiveGame()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
}

std::string MonopolyLiveGame::spin() {
    std::string combined = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(combined);
    uint64_t index = hash % 54;
    
    // Map to segments (simplified)
    if (index < 6) return "1";
    else if (index < 12) return "2";
    else if (index < 18) return "5";
    else if (index < 22) return "10";
    else if (index < 24) return "2 ROLLS";
    else if (index < 26) return "4 ROLLS";
    else if (index < 27) return "CHANCE";
    else return "JAIL";
}

MonopolyLiveGame::PropertyResult MonopolyLiveGame::playBonus() {
    PropertyResult result;
    result.total_payout = 0;
    result.rolls = 2;
    
    for (int i = 0; i < result.rolls; i++) {
        std::string combined = provably_fair_->getServerSeed() + "PROPERTY" + std::to_string(i);
        uint64_t hash = provably_fair_->generateHash(combined);
        
        // Simplified property value
        double value = (hash % 500 + 50) / 10.0;
        result.properties.push_back("Property " + std::to_string(i + 1));
        result.payouts.push_back(value);
        result.total_payout += value;
    }
    
    return result;
}

void MonopolyLiveGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

void MonopolyLiveGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

// Dream Catcher Implementation
DreamCatcherGame::DreamCatcherGame()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>()) {
}

std::string DreamCatcherGame::spin() {
    std::string combined = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(combined);
    uint64_t index = hash % 54;
    
    if (index < 20) return "1x";
    else if (index < 32) return "2x";
    else if (index < 40) return "5x";
    else if (index < 44) return "10x";
    else if (index < 47) return "15x";
    else if (index < 49) return "20x";
    else if (index < 51) return "40x";
    else return "2x";
}

void DreamCatcherGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

void DreamCatcherGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

// Lightning Roulette Implementation
LightningRouletteGame::LightningRouletteGame()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>()) {
}

LightningRouletteGame::SpinResult LightningRouletteGame::spin() {
    SpinResult result;
    
    // Generate winning number using provably fair
    std::string combined = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(combined);
    result.winningNumber = static_cast<int>(hash % 37);
    
    // Generate lightning numbers (5 numbers)
    for (int i = 0; i < 5; i++) {
        std::string hash_input = combined + "LIGHTNING" + std::to_string(i);
        uint64_t lhash = provably_fair_->generateHash(hash_input);
        int lightningNum = static_cast<int>(lhash % 37);
        result.lightningNumbers.push_back(lightningNum);
    }
    
    // Calculate multiplier if winning number is a lightning number
    for (int ln : result.lightningNumbers) {
        if (ln == result.winningNumber) {
            std::string mult_input = combined + "MULTIPLIER";
            uint64_t mhash = provably_fair_->generateHash(mult_input);
            result.multiplier = static_cast<int>((mhash % 500) + 50) / 10.0; // 50x to 500x
            break;
        }
    }
    
    if (result.multiplier == 0) result.multiplier = 1;
    
    return result;
}

void LightningRouletteGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

void LightningRouletteGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
