#include "MassiveSlotGenerator.hpp"
#include <random>
#include <sstream>

namespace TigerCasino {

MassiveSlotGenerator::MassiveSlotGenerator() {
    providers = {
        "PragmaticPlay", "NetEnt", "PlaynGO", "Yggdrasil", "Quickspin",
        "Betsoft", "BGaming", "Spinomenal", "Wazdan", "Evoplay",
        "Endorphina", "BoomingGames", "NolimitCity", "Hacksaw", "BigTimeGaming",
        "RedTiger", "Blueprint", "IronDog", "Thunderkick", "ELKStudios",
        "RelaxGaming", "PushGaming", "Mancala", "RubyPlay", "Dragoon"
    };
    
    themes = {
        "classic", "adventure", "mythology", "fantasy", "nature", "ocean",
        "animal", "egypt", "asian", "western", "horror", "sci-fi",
        "fruit", "luxury", "sports", "music", "movie", "history", "food",
        "party", "jewels", "money", "fantasy", "vintage", "modern"
    };
    
    volatilities = {"low", "medium", "high"};
    
    slotPrefixes = {
        "Golden", "Lucky", "Royal", "Mega", "Super", "Extreme", "Pro", "Master",
        "Grand", "Epic", "Ultimate", "Premium", "Deluxe", "Platinum", "Diamond",
        "Ruby", "Emerald", "Sapphire", "Crystal", "Thunder", "Storm", "Fire", "Ice",
        "Dragon", "Phoenix", "Tiger", "Lion", "Eagle", "Wolf", "Bear", "Shark",
        "Magic", "Wizard", "Spell", "Enchanted", "Mystic", "Secret", "Hidden",
        "Treasure", "Cash", "Fortune", "Wealth", "Rich", "Jackpot", "Winning",
        "King", "Queen", "Prince", "Princess", "Hero", "Legend", "Warrior"
    };
    
    slotSuffixes = {
        "King", "Queen", "Prince", "Princess", "Lord", "Lady", "Hero", "Legend",
        "Gems", "Wins", "Fortune", "Riches", "Treasure", "Jackpot", "Millions", "Cash",
        "Wheel", "Quest", "Adventure", "Journey", "Spins", "Reels", "Hot", "Cool", "Wild",
        "Free", "Bonus", "Deluxe", "Paradise", "Island", "World", "Land", "Kingdom",
        "Valley", "Mountain", "Ocean", "Forest", "Desert", "City", "Town", "Village"
    };
}

std::vector<MassiveSlotConfig> MassiveSlots::MassiveSlotGenerator::generateAll() {
    std::vector<MassiveSlotConfig> slots;
    std::random_device rd;
    std::mt19937 gen(rd());
    
    std::uniform_int_distribution<> reelDist(3, 6);
    std::uniform_int_distribution<> rowDist(3, 5);
    std::uniform_int_distribution<> paylineDist(5, 200);
    std::uniform_real_distribution<> rtpDist(0.90, 0.99);
    std::uniform_real_distribution<> betDist(0.01, 2.0);
    std::uniform_real_distribution<> maxBetDist(100.0, 2000.0);
    
    for (int i = 0; i < 3300; i++) {
        MassiveSlotConfig slot;
        slot.gameId = "mslot_" + std::to_string(i + 1);
        
        std::uniform_int_distribution<> nameDist(0, (int)slotPrefixes.size() - 1);
        std::uniform_int_distribution<> suffixDist(0, (int)slotSuffixes.size() - 1);
        slot.gameName = slotPrefixes[nameDist(gen)] + " " + slotSuffixes[suffixDist(gen)];
        
        std::uniform_int_distribution<> provDist(0, (int)providers.size() - 1);
        slot.provider = providers[provDist(gen)];
        
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
        
        slot.category = "slots";
        slots.push_back(slot);
    }
    
    return slots;
}

std::map<std::string, std::vector<MassiveSlotConfig>> MassiveSlots::MassiveSlotGenerator::groupByProvider() {
    auto allSlots = generateAll();
    std::map<std::string, std::vector<MassiveSlotConfig>> grouped;
    
    for (const auto& slot : allSlots) {
        grouped[slot.provider].push_back(slot);
    }
    
    return grouped;
}

} // namespace TigerCasino
