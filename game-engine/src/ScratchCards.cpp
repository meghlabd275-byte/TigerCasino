// Scratch Cards Game Implementation
#include "ScratchCards.hpp"
#include <iostream>
#include <random>
#include <sstream>

namespace TigerCasino {

// ScratchCard Game Implementation
ScratchCardGame::ScratchCardGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , round_counter_(0) {
}

ScratchCardGame::ScratchResult ScratchCardGame::play(const std::string& cardType, double bet) {
    ScratchResult result;
    result.card_type = cardType;
    result.bet = bet;
    result.card_id = generateCardId();
    
    // Use provably fair for card generation
    std::string combined = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(combined);
    
    // Generate symbols based on card type
    if (cardType == "classic") {
        result = generateClassicCard(hash, bet);
    } else if (cardType == " themed") {
        result = generateThemedCard(hash, bet);
    } else if (cardType == "instant_win") {
        result = generateInstantWinCard(hash, bet);
    } else {
        result = generateClassicCard(hash, bet);
    }
    
    result.server_seed = provably_fair_->getServerSeed();
    result.timestamp = time(nullptr);
    
    return result;
}

std::string ScratchCardGame::generateCardId() {
    std::stringstream ss;
    ss << "SCR_" << round_counter_ << "_" << provably_fair_->getServerSeed();
    return ss.str();
}

ScratchCardGame::ScratchResult ScratchCardGame::generateClassicCard(uint64_t hash, double bet) {
    ScratchResult result;
    result.win = false;
    result.prize = 0;
    result.multiplier = 0;
    
    // Classic scratch card - match 3 symbols
    std::string symbols[] = {"🍒", "🍋", "🍇", "💎", "⭐", "🔔"};
    
    // Generate 3x3 grid
    for (int i = 0; i < 9; i++) {
        uint64_t symbol_hash = provably_fair_->generateHash(
            std::to_string(hash) + "_symbol_" + std::to_string(i)
        );
        int symbol_idx = symbol_hash % 6;
        result.symbols.push_back(symbols[symbol_idx]);
    }
    
    // Check for win (simplified - middle row all same)
    if (result.symbols[3] == result.symbols[4] && result.symbols[4] == result.symbols[5]) {
        result.win = true;
        result.prize = bet * 10;
        result.multiplier = 10;
    } else if (result.symbols[0] == result.symbols[4] && result.symbols[4] == result.symbols[8]) {
        result.win = true;
        result.prize = bet * 25;
        result.multiplier = 25;
    } else if (result.symbols[2] == result.symbols[4] && result.symbols[4] == result.symbols[6]) {
        result.win = true;
        result.prize = bet * 50;
        result.multiplier = 50;
    }
    
    return result;
}

ScratchCardGame::ScratchResult ScratchCardGame::generateThemedCard(uint64_t hash, double bet) {
    ScratchResult result;
    result.win = false;
    result.prize = 0;
    result.multiplier = 0;
    
    // Themed scratch card - themed symbols
    std::string symbols[] = {"🐯", "🦁", "🐻", "🐼", "🐨", "🦊"};
    
    for (int i = 0; i < 9; i++) {
        uint64_t symbol_hash = provably_fair_->generateHash(
            std::to_string(hash) + "_themed_" + std::to_string(i)
        );
        int symbol_idx = symbol_hash % 6;
        result.symbols.push_back(symbols[symbol_idx]);
    }
    
    // Check for win
    if (result.symbols[3] == result.symbols[4] && result.symbols[4] == result.symbols[5]) {
        result.win = true;
        result.prize = bet * 15;
        result.multiplier = 15;
    }
    
    return result;
}

ScratchCardGame::ScratchResult ScratchCardGame::generateInstantWinCard(uint64_t hash, double bet) {
    ScratchResult result;
    
    // Instant win - reveal one symbol
    uint64_t prize_hash = provably_fair_->generateHash(std::to_string(hash) + "_prize");
    uint64_t win_hash = provably_fair_->generateHash(std::to_string(hash) + "_win");
    
    // 30% win rate
    result.win = (win_hash % 100) < 30;
    
    if (result.win) {
        uint64_t prize_val = prize_hash % 100;
        if (prize_val < 50) {
            result.prize = bet * 2;
            result.multiplier = 2;
        } else if (prize_val < 80) {
            result.prize = bet * 5;
            result.multiplier = 5;
        } else if (prize_val < 95) {
            result.prize = bet * 10;
            result.multiplier = 10;
        } else {
            result.prize = bet * 50;
            result.multiplier = 50;
        }
    } else {
        result.prize = 0;
        result.multiplier = 0;
    }
    
    result.symbols.push_back(result.win ? "🎉" : "❌");
    
    return result;
}

void ScratchCardGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

void ScratchCardGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

// Scratch Card Manager
ScratchCardManager::ScratchCardManager() {
    // Initialize available scratch cards
    scratch_cards_["classic"] = ScratchCardConfig{
        id: "classic",
        name: "Classic Scratch",
        min_bet: 0.10,
        max_bet: 100.0,
        rtp: 95.0,
        description: "Match 3 symbols to win!"
    };
    
    scratch_cards_["animal"] = ScratchCardConfig{
        id: "animal",
        name: "Animal Kingdom",
        min_bet: 0.10,
        max_bet: 100.0,
        rtp: 95.5,
        description: "Scratch to find animals!"
    };
    
    scratch_cards_["lucky"] = ScratchCardConfig{
        id: "lucky",
        name: "Lucky Scratch",
        min_bet: 0.10,
        max_bet: 100.0,
        rtp: 96.0,
        description: "Instant win scratch card!"
    };
    
    scratch_cards_["treasure"] = ScratchCardConfig{
        id: "treasure",
        name: "Treasure Hunt",
        min_bet: 0.50,
        max_bet: 500.0,
        rtp: 94.0,
        description: "Find the treasure!"
    };
    
    scratch_cards_["diamond"] = ScratchCardConfig{
        id: "diamond",
        name: "Diamond Riches",
        min_bet: 1.0,
        max_bet: 1000.0,
        rtp: 93.0,
        description: "Win diamond prizes!"
    };
}

std::vector<ScratchCardConfig> ScratchCardManager::getAvailableCards() {
    std::vector<ScratchCardConfig> cards;
    for (const auto& pair : scratch_cards_) {
        cards.push_back(pair.second);
    }
    return cards;
}

ScratchCardConfig ScratchCardManager::getCardConfig(const std::string& cardId) {
    if (scratch_cards_.find(cardId) != scratch_cards_.end()) {
        return scratch_cards_[cardId];
    }
    return scratch_cards_["classic"];
}

ScratchCardGame::ScratchResult ScratchCardManager::scratch(const std::string& cardId, double bet) {
    ScratchCardGame game;
    return game.play(cardId, bet);
}

} // namespace TigerCasino
