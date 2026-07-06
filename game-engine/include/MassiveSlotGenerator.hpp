#pragma once

#include <string>
#include <vector>
#include <map>

namespace TigerCasino {

// Massive slot generator - 3000+ games
struct MassiveSlotConfig {
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
    std::string category;
};

namespace MassiveSlots {

class MassiveSlotGenerator {
private:
    std::vector<std::string> providers;
    std::vector<std::string> themes;
    std::vector<std::string> volatilities;
    std::vector<std::string> slotPrefixes;
    std::vector<std::string> slotSuffixes;
    
public:
    MassiveSlotGenerator();
    std::vector<MassiveSlotConfig> generateAll();
    std::map<std::string, std::vector<MassiveSlotConfig>> groupByProvider();
    size_t totalCount() const { return 3300; }
};

} // namespace MassiveSlots
} // namespace TigerCasino
