#include "BulkSlotGenerator.hpp"
#include <random>
#include <sstream>

namespace TigerCasino {

std::vector<BulkSlotConfig> BulkSlots::BulkSlotGenerator::generateSlots(int count) {
    std::vector<BulkSlotConfig> slots;
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> reelDist(3, 6);
    std::uniform_int_distribution<> rowDist(3, 5);
    std::uniform_int_distribution<> paylineDist(5, 100);
    std::uniform_real_distribution<> rtpDist(0.92, 0.98);
    std::uniform_real_distribution<> betDist(0.01, 1.0);
    std::uniform_real_distribution<> maxBetDist(100.0, 1000.0);
    
    std::vector<std::string> slotNames = {
        "Golden", "Lucky", "Royal", "Mega", "Super", "Extreme", "Pro", "Master",
        "Royal", "Grand", "Epic", "Ultimate", "Premium", "Deluxe", "Platinum", "Diamond",
        "Ruby", "Emerald", "Sapphire", "Crystal", "Thunder", "Storm", "Fire", "Ice",
        "Dragon", "Phoenix", "Tiger", "Lion", "Eagle", "Wolf", "Bear", "Shark",
        "Magic", "Wizard", "Wizard", "Spell", "Enchanted", "Mystic", "Secret", "Hidden",
        "Treasure", "Gold", "Cash", "Money", "Fortune", "Wealth", "Rich", "Jackpot"
    };
    
    std::vector<std::string> slotSuffixes = {
        "King", "Queen", "Prince", "Princess", "Lord", "Lady", "Hero", "Legend",
        "Gems", "Wins", "Fortune", "Riches", "Treasure", "Jackpot", "Millions", "Cash",
        "Wheel", "Quest", "Adventure", "Journey", "Quest", "Quest", "Quest", "Quest",
        "Spins", "Reels", "Hot", "Cool", "Wild", "Free", "Bonus", "Deluxe"
    };
    
    for (int i = 0; i < count; i++) {
        BulkSlotConfig slot;
        slot.gameId = "slot_" + std::to_string(i + 1);
        
        // Generate random name
        std::uniform_int_distribution<> nameDist(0, (int)slotNames.size() - 1);
        std::uniform_int_distribution<> suffixDist(0, (int)slotSuffixes.size() - 1);
        slot.gameName = slotNames[nameDist(gen)] + " " + slotSuffixes[suffixDist(gen)];
        
        // Random provider
        std::uniform_int_distribution<> provDist(0, (int)providers.size() - 1);
        slot.provider = providers[provDist(gen)];
        
        // Random specs
        slot.reels = reelDist(gen);
        slot.rows = rowDist(gen);
        slot.paylines = paylineDist(gen);
        slot.minBet = betDist(gen);
        slot.maxBet = maxBetDist(gen);
        slot.rtp = rtpDist(gen);
        
        std::uniform_int_distribution<> volDist(0, (int)volatilities.size() - 1);
        slot.volatility = volatilities[volDist(gen)];
        
        std::uniform_int_distribution<> themeDist(0, (int)themes.size() - 1);
        slot.theme = themes[themeDist(gen)];
        
        slots.push_back(slot);
    }
    
    return slots;
}

std::map<std::string, std::vector<BulkSlotConfig>> BulkSlots::BulkSlotGenerator::groupByProvider() {
    auto allSlots = generateSlots(500);
    std::map<std::string, std::vector<BulkSlotConfig>> grouped;
    
    for (const auto& slot : allSlots) {
        grouped[slot.provider].push_back(slot);
    }
    
    return grouped;
}

size_t BulkSlots::BulkSlotGenerator::totalCount() const {
    return 500;
}

} // namespace TigerCasino
