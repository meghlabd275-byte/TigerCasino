#pragma once

#include <string>
#include <vector>
#include <map>

namespace TigerCasino {

// Bulk slot games generator - 500 games
struct BulkSlotConfig {
    std::string gameId;
    std::string gameName;
    std::string provider;
    int reels;
    int rows;
    int paylines;
    double minBet;
    double maxBet;
    double rtp;
    std::string volatility;
    std::string theme;
};

namespace BulkSlots {

class BulkSlotGenerator {
private:
    std::vector<std::string> providers = {
        "PragmaticPlay", "NetEnt", "PlaynGO", "Yggdrasil", "Quickspin",
        "Betsoft", "BGaming", "Spinomenal", "Wazdan", "Evoplay",
        "Endorphina", "BoomingGames", "NolimitCity", "Hacksaw", "BigTimeGaming"
    };
    
    std::vector<std::string> themes = {
        "classic", "adventure", "mythology", "fantasy", "nature", "ocean",
        "animal", "egypt", "asian", "western", "horror", "sci-fi",
        "fruit", "luxury", "sports", "music", "movie", "history"
    };
    
    std::vector<std::string> volatilities = {"low", "medium", "high"};
    
public:
    std::vector<BulkSlotConfig> generateSlots(int count);
    std::map<std::string, std::vector<BulkSlotConfig>> groupByProvider();
    size_t totalCount() const;
};

} // namespace BulkSlots
} // namespace TigerCasino
