#include "SlotsBatch3.hpp"

namespace TigerCasino {

Batch3Manager::Batch3Manager() {
    using namespace SlotsBatch3;
    
    // Betsoft (20)
    for (int i = 0; i < 20; i++) {
        std::string id = "betsoft_" + std::to_string(i);
        slots_[id] = createBetsoftSlots();
        slots_[id].gameId = id;
    }
    
    // Booming Games (20)
    for (int i = 0; i < 20; i++) {
        std::string id = "booming_" + std::to_string(i);
        slots_[id] = createBoomingSlots();
        slots_[id].gameId = id;
    }
    
    // Spinomenal (20)
    for (int i = 0; i < 20; i++) {
        std::string id = "spinomenal_" + std::to_string(i);
        slots_[id] = createSpinomenalSlots();
        slots_[id].gameId = id;
    }
    
    // Wazdan (20)
    for (int i = 0; i < 20; i++) {
        std::string id = "wazdan_" + std::to_string(i);
        slots_[id] = createWazdanSlots();
        slots_[id].gameId = id;
    }
    
    // Endorphina (20)
    for (int i = 0; i < 20; i++) {
        std::string id = "endorphina_" + std::to_string(i);
        slots_[id] = createEndorphinaSlots();
        slots_[id].gameId = id;
    }
}

std::vector<MoreSlotConfig> Batch3Manager::getAll() const {
    std::vector<MoreSlotConfig> result;
    for (const auto& p : slots_) result.push_back(p.second);
    return result;
}

size_t Batch3Manager::count() const { return slots_.size(); }

namespace SlotsBatch3 {

MoreSlotConfig createBetsoftSlots() {
    return {"betsoft_1", "Betsoft Slot", "Betsoft", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "various"};
}

MoreSlotConfig createBoomingSlots() {
    return {"booming_1", "Booming Slot", "BoomingGames", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "various"};
}

MoreSlotConfig createSpinomenalSlots() {
    return {"spinomenal_1", "Spinomenal Slot", "Spinomenal", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "various"};
}

MoreSlotConfig createWazdanSlots() {
    return {"wazdan_1", "Wazdan Slot", "Wazdan", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "various"};
}

MoreSlotConfig createEndorphinaSlots() {
    return {"endorphina_1", "Endorphina Slot", "Endorphina", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "various"};
}

} // namespace SlotsBatch3
} // namespace TigerCasino
