#include "LiveDealerServer.hpp"
#include <sstream>
#include <chrono>

namespace TigerCasino {

LiveDealerServer::LiveDealerServer() : running_(false) {
}

LiveDealerServer::~LiveDealerServer() {
    stop();
}

void LiveDealerServer::start() {
    running_ = true;
}

void LiveDealerServer::stop() {
    running_ = false;
}

bool LiveDealerServer::isRunning() const {
    return running_;
}

void LiveDealerServer::addTable(const LiveTable& table) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    tables_[table.tableId] = table;
}

void LiveDealerServer::removeTable(const std::string& tableId) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    tables_.erase(tableId);
    activeRounds_.erase(tableId);
}

LiveTable LiveDealerServer::getTable(const std::string& tableId) const {
    std::lock_guard<std::mutex> lock(const_cast<std::mutex&>(tablesMutex_));
    auto it = tables_.find(tableId);
    if (it != tables_.end()) {
        return it->second;
    }
    return LiveTable();
}

std::vector<LiveTable> LiveDealerServer::getTablesByGame(LiveGameType type) const {
    std::lock_guard<std::mutex> lock(const_cast<std::mutex&>(tablesMutex_));
    std::vector<LiveTable> result;
    
    for (const auto& pair : tables_) {
        if (pair.second.gameType == type) {
            result.push_back(pair.second);
        }
    }
    
    return result;
}

std::vector<LiveTable> LiveDealerServer::getAllTables() const {
    std::lock_guard<std::mutex> lock(const_cast<std::mutex&>(tablesMutex_));
    std::vector<LiveTable> result;
    
    for (const auto& pair : tables_) {
        result.push_back(pair.second);
    }
    
    return result;
}

bool LiveDealerServer::addPlayerToTable(const std::string& tableId, const LivePlayer& player) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    
    auto it = tables_.find(tableId);
    if (it == tables_.end()) return false;
    
    // Check if player already exists
    for (auto& p : it->second.players) {
        if (p.playerId == player.playerId) {
            p = player;
            return true;
        }
    }
    
    it->second.players.push_back(player);
    return true;
}

bool LiveDealerServer::removePlayerFromTable(const std::string& tableId, const std::string& playerId) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    
    auto it = tables_.find(tableId);
    if (it == tables_.end()) return false;
    
    auto& players = it->second.players;
    for (auto it2 = players.begin(); it2 != players.end(); ++it2) {
        if (it2->playerId == playerId) {
            players.erase(it2);
            return true;
        }
    }
    
    return false;
}

bool LiveDealerServer::placeBet(const std::string& tableId, const std::string& playerId,
                                 const std::string& betType, double amount) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    
    auto it = tables_.find(tableId);
    if (it == tables_.end()) return false;
    
    if (it->second.status != TableStatus::WAITING && 
        it->second.status != TableStatus::BETTING_CLOSED) {
        return false;
    }
    
    // Find player
    for (auto& player : it->second.players) {
        if (player.playerId == playerId) {
            if (player.balance < amount) return false;
            if (amount < it->second.minBet || amount > it->second.maxBet) return false;
            
            player.bets[betType] += amount;
            player.balance -= amount;
            return true;
        }
    }
    
    return false;
}

LiveRound LiveDealerServer::startRound(const std::string& tableId) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    
    LiveRound round;
    round.tableId = tableId;
    round.startTime = std::chrono::steady_clock::now();
    
    // Generate round ID
    std::stringstream ss;
    ss << tableId << "_" << std::chrono::steady_clock::now().time_since_epoch().count();
    round.roundId = ss.str();
    
    auto it = tables_.find(tableId);
    if (it != tables_.end()) {
        round.dealerName = it->second.dealerName;
        it->second.status = TableStatus::DEALING;
    }
    
    activeRounds_[tableId] = round;
    return round;
}

void LiveDealerServer::endRound(const std::string& roundId, const std::vector<std::string>& result) {
    std::lock_guard<std::mutex> lock(tablesMutex_);
    
    for (auto& pair : activeRounds_) {
        if (pair.second.roundId == roundId) {
            pair.second.result = result;
            pair.second.endTime = std::chrono::steady_clock::now();
            
            // Update table status
            auto it = tables_.find(pair.first);
            if (it != tables_.end()) {
                it->second.status = TableStatus::RESULT;
            }
            
            // Calculate wins
            if (it != tables_.end()) {
                for (auto& player : it->second.players) {
                    for (const auto& bet : player.bets) {
                        // Simplified win calculation
                        double win = bet.second * 2;  // Placeholder
                        player.wins[bet.first] = win;
                        player.balance += win;
                    }
                    player.bets.clear();
                }
            }
            
            if (onRoundComplete_) {
                onRoundComplete_(pair.second);
            }
            
            break;
        }
    }
}

LiveRound LiveDealerServer::getActiveRound(const std::string& tableId) const {
    auto it = activeRounds_.find(tableId);
    if (it != activeRounds_.end()) {
        return it->second;
    }
    return LiveRound();
}

void LiveDealerServer::registerProvider(const std::string& providerName,
                                        std::function<std::string(const std::string&)> handler) {
    providers_[providerName] = handler;
}

void LiveDealerServer::setTableUpdateCallback(std::function<void(const LiveTable&)> callback) {
    onTableUpdate_ = callback;
}

void LiveDealerServer::setRoundCompleteCallback(std::function<void(const LiveRound&)> callback) {
    onRoundComplete_ = callback;
}

// Live provider table templates
namespace LiveProviders {

std::vector<LiveTable> createEvolutionTables() {
    std::vector<LiveTable> tables;
    
    // Blackjack tables
    LiveTable blackjack1;
    blackjack1.tableId = "evo_bj_01";
    blackjack1.tableName = "Blackjack A";
    blackjack1.provider = "evolution";
    blackjack1.gameType = LiveGameType::BLACKJACK;
    blackjack1.minBet = 10;
    blackjack1.maxBet = 5000;
    blackjack1.dealerName = "Alex";
    blackjack1.streamUrl = "https://cdn.evolutiongaming.com/blackjack-a";
    tables.push_back(blackjack1);
    
    // Roulette
    LiveTable roulette1;
    roulette1.tableId = "evo_rou_01";
    roulette1.tableName = "Speed Roulette";
    roulette1.provider = "evolution";
    roulette1.gameType = LiveGameType::ROULETTE;
    roulette1.minBet = 1;
    roulette1.maxBet = 10000;
    roulette1.dealerName = "Maria";
    roulette1.streamUrl = "https://cdn.evolutiongaming.com/speed-roulette";
    tables.push_back(roulette1);
    
    // Baccarat
    LiveTable baccarat1;
    baccarat1.tableId = "evo_bac_01";
    baccarat1.tableName = "Baccarat Squeeze";
    baccarat1.provider = "evolution";
    baccarat1.gameType = LiveGameType::BACCARAT;
    baccarat1.minBet = = 15;
    baccarat1.maxBet = 10000;
    baccarat1.dealerName = "Chen";
    baccarat1.streamUrl = "https://cdn.evolutiongaming.com/baccarat-squeeze";
    tables.push_back(baccarat1);
    
    // Game shows
    LiveTable dreamCatcher;
    dreamCatcher.tableId = "evo_dc_01";
    dreamCatcher.tableName = "Dream Catcher";
    dreamCatcher.provider = "evolution";
    dreamCatcher.gameType = LiveGameType::GAME_SHOW;
    dreamCatcher.minBet = 1;
    dreamCatcher.maxBet = 1000;
    dreamCatcher.dealerName = "Jessica";
    dreamCatcher.streamUrl = "https://cdn.evolutiongaming.com/dream-catcher";
    tables.push_back(dreamCatcher);
    
    return tables;
}

std::vector<LiveTable> createPragmaticPlayTables() {
    std::vector<LiveTable> tables;
    
    LiveTable megaWheel;
    megaWheel.tableId = "pp_mw_01";
    megaWheel.tableName = "Mega Wheel";
    megaWheel.provider = "pragmatic_play";
    megaWheel.gameType = LiveGameType::GAME_SHOW;
    megaWheel.minBet = 1;
    megaWheel.maxBet = 500;
    megaWheel.dealerName = "Emma";
    megaWheel.streamUrl = "https://cdn.pragmaticplay.com/mega-wheel";
    tables.push_back(megaWheel);
    
    LiveTable sweetBonanzaCandy;
    sweetBonanzaCandy.tableId = "pp_sb_01";
    sweetBonanzaCandy.tableName = "Sweet Bonanza CandyLand";
    sweetBonanzaCandy.provider = "pragmatic_play";
    sweetBonanzaCandy.gameType = LiveGameType::GAME_SHOW;
    sweetBonanzaCandy.minBet = 1;
    sweetBonanzaCandy.maxBet = 500;
    sweetBonanzaCandy.dealerName = "Sophie";
    sweetBonanzaCandy.streamUrl = "https://cdn.pragmaticplay.com/sweet-bonanza-candyland";
    tables.push_back(sweetBonanzaCandy);
    
    return tables;
}

std::vector<LiveTable> createEzugiTables() {
    std::vector<LiveTable> tables;
    
    LiveTable ktRoulette;
    ktRoulette.tableId = "ez_rou_01";
    ktRoulette.tableName = "KT Roulette";
    ktRoulette.provider = "ezugi";
    ktRoulette.gameType = LiveGameType::ROULETTE;
    ktRoulette.minBet = 1;
    ktRoulette.maxBet = 5000;
    ktRoulette.dealerName = "Katie";
    ktRoulette.streamUrl = "https://cdn.ezugi.com/kt-roulette";
    tables.push_back(ktRoulette);
    
    return tables;
}

std::vector<LiveTable> createAuthenticGamingTables() {
    std::vector<LiveTable> tables;
    
    LiveTable authenticRoulette;
    authenticRoulette.tableId = "ag_rou_01";
    authenticRoulette.tableName = "Authentic Roulette";
    authenticRoulette.provider = "authentic_gaming";
    authenticRoulette.gameType = LiveGameType::ROULETTE;
    authenticRoulette.minBet = 1;
    authenticRoulette.maxBet = 10000;
    authenticRoulette.dealerName = "Live Dealer";
    authenticRoulette.streamUrl = "https://cdn.authenticgaming.com/authentic-roulette";
    tables.push_back(authenticRoulette);
    
    return tables;
}

} // namespace LiveProviders

} // namespace TigerCasino
