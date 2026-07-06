#pragma once

#include <string>
#include <vector>
#include <map>

namespace TigerCasino {

// More slot games - batch 3 (100 more games)
struct MoreSlotConfig {
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

namespace SlotsBatch3 {

// Betsoft games (20)
MoreSlotConfig createBetsoftSlots();
// Inspired/Booming Games (20)
MoreSlotConfig createBoomingSlots();
// Spinomenal (20)
MoreSlotConfig createSpinomenalSlots();
// Wazdan (20)
MoreSlotConfig createWazdanSlots();
// Endorphina (20)
MoreSlotConfig createEndorphinaSlots();

class Batch3Manager {
private:
    std::map<std::string, MoreSlotConfig> slots_;
public:
    Batch3Manager();
    std::vector<MoreSlotConfig> getAll() const;
    size_t count() const;
};

} // namespace SlotsBatch3
} // namespace TigerCasino
