// Tiger Original Games Implementation
#include "TigerOriginals.hpp"
#include <iostream>
#include <random>
#include <sstream>
#include <cmath>

namespace TigerCasino {

// TigerCrash Implementation
TigerCrashGame::TigerCrashGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , current_multiplier_(1.0)
    , game_state_(GameState::WAITING)
    , crash_point_(0.0) {
}

void TigerCrashGame::startRound() {
    // Generate crash point using provably fair
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    // Crash point calculation (exponential distribution)
    double r = (double)(hash % 10000) / 10000.0;
    crash_point_ = 1.0 / (1.0 - r);
    
    // Cap at 1000x
    if (crash_point_ > 1000.0) crash_point_ = 1000.0;
    if (crash_point_ < 1.0) crash_point_ = 1.0;
    
    game_state_ = GameState::RUNNING;
    current_multiplier_ = 1.0;
    start_time_ = std::time(nullptr);
}

bool TigerCrashGame::cashOut(const std::string& bet_id) {
    if (game_state_ != GameState::RUNNING) return false;
    
    auto it = bets_.find(bet_id);
    if (it == bets_.end()) return false;
    
    double payout = it->second.bet_amount * current_multiplier_;
    it->second.payout = payout;
    it->second.cashed_out = true;
    it->second.cash_out_multiplier = current_multiplier_;
    
    return true;
}

double TigerCrashGame::getCurrentMultiplier() const {
    if (game_state_ != GameState::RUNNING) return 0.0;
    
    time_t now = std::time(nullptr);
    double elapsed = difftime(now, start_time_);
    
    // Multiplier grows exponentially after 1 second
    double multiplier = 1.0;
    if (elapsed > 0) {
        multiplier = 1.0 + (elapsed * 0.1); // 10% per second
        if (multiplier > crash_point_) multiplier = crash_point_;
    }
    
    return multiplier;
}

void TigerCrashGame::endRound() {
    game_state_ = GameState::CRASHED;
    // Settle all remaining bets
    for (auto& bet : bets_) {
        if (!bet.second.cashed_out) {
            bet.second.payout = 0.0;
        }
    }
}

std::string TigerCrashGame::placeBet(const std::string& user_id, double amount) {
    std::string bet_id = "CRASH_" + std::to_string(bets_.size() + 1);
    bets_[bet_id] = {
        user_id: user_id,
        bet_amount: amount,
        payout: 0.0,
        cashed_out: false,
        cash_out_multiplier: 0.0
    };
    return bet_id;
}

// TigerMines Implementation
TigerMinesGame::TigerMinesGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , num_mines_(3)
    , current_step_(0)
    , game_over_(false) {
}

void TigerMinesGame::startRound(int num_mines) {
    num_mines_ = num_mines;
    current_step_ = 0;
    game_over_ = false;
    revealed_tiles_.clear();
    mine_positions_.clear();
    
    // Generate mine positions using provably fair
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    // Place mines randomly
    std::set<int> mines;
    uint64_t h = hash;
    while ((int)mines.size() < num_mines_) {
        int pos = h % 25;
        mines.insert(pos);
        h = provably_fair_->generateHash(std::to_string(h) + "_mine");
    }
    
    for (int m : mines) {
        mine_positions_.insert(m);
    }
}

bool TigerMinesGame::revealTile(int position) {
    if (game_over_) return false;
    if (position < 0 || position >= 25) return false;
    if (revealed_tiles_.count(position)) return false;
    
    revealed_tiles_.insert(position);
    current_step_++;
    
    if (mine_positions_.count(position)) {
        game_over_ = true;
        return false; // Hit a mine
    }
    
    return true; // Safe tile
}

bool TigerMinesGame::cashOut() {
    if (game_over_) return false;
    if (revealed_tiles_.empty()) return false;
    
    // Calculate payout based on safe tiles revealed
    double multiplier = 1.0;
    int safe_tiles = revealed_tiles_.size();
    for (int i = 0; i < safe_tiles; i++) {
        multiplier *= (25.0 - num_mines_) / (25.0 - i);
    }
    
    current_payout_ = current_bet_ * multiplier;
    game_over_ = true;
    
    return true;
}

int TigerMinesGame::getTileValue(int position) const {
    if (mine_positions_.count(position)) return -1; // Mine
    
    // Calculate diamond count for adjacent tiles
    int diamonds = 0;
    int row = position / 5;
    int col = position % 5;
    
    for (int dr = -1; dr <= 1; dr++) {
        for (int dc = -1; dc <= 1; dc++) {
            if (dr == 0 && dc == 0) continue;
            int nr = row + dr;
            int nc = col + dc;
            if (nr >= 0 && nr < 5 && nc >= 0 && nc < 5) {
                int np = nr * 5 + nc;
                if (mine_positions_.count(np)) diamonds++;
            }
        }
    }
    
    return diamonds;
}

// TigerPlinko Implementation
TigerPlinkoGame::TigerPlinkoGame()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rows_(8) {
}

void TigerPlinkoGame::startRound(int num_rows) {
    rows_ = num_rows;
    drop_position_ = -1;
    final_bucket_ = -1;
    path_.clear();
    game_over_ = false;
}

std::vector<int> TigerPlinkoGame::simulateDrop(int start_position) {
    std::vector<int> path;
    int position = start_position;
    path.push_back(position);
    
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    for (int row = 0; row < rows_; row++) {
        // Each step, ball goes left or right
        bool go_right = (h % 2 == 1);
        h = provably_fair_->generateHash(std::to_string(h) + "_drop");
        
        if (go_right) {
            position = position + (row + 1);
        }
        path.push_back(position);
    }
    
    drop_position_ = start_position;
    final_bucket_ = position;
    path_ = path;
    game_over_ = true;
    
    return path;
}

double TigerPlinkoGame::getMultiplier(int bucket) const {
    // Multipliers based on bucket (center = highest)
    int center = rows_ / 2;
    int distance = std::abs(bucket - center);
    
    // Higher multipliers at edges, lower in center
    static const double multipliers[] = {100.0, 50.0, 25.0, 10.0, 5.0, 2.0, 1.0, 0.5, 0.2};
    
    int index = std::min(distance, (int)(sizeof(multipliers) / sizeof(multipliers[0]) - 1));
    return multipliers[index];
}

// TigerDice Implementation
TigerDiceGame::TigerDiceGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double TigerDiceGame::roll(uint64_t target, double bet) {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    // Roll 0-9999 for 4 decimal precision
    uint64_t roll = hash % 10000;
    double roll_value = roll / 100.0; // 0.00 to 99.99
    
    bool win = (roll_value < target);
    double payout = win ? bet * (100.0 / target) : 0.0;
    
    return roll_value;
}

// TigerSlot Implementation
TigerSlotGame::TigerSlotGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<std::vector<std::string>> TigerSlotGame::spin(int bet_level) {
    std::vector<std::vector<std::string>> reels(5, std::vector<std::string>(3));
    
    std::string symbols[] = {"🍒", "🍋", "🍇", "💎", "⭐", "🔔", "7️⃣", "WILD"};
    int num_symbols = 8;
    
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    for (int reel = 0; reel < 5; reel++) {
        for (int row = 0; row < 3; h = provably_firm_->generateHash(std::to_string(h) + "_symbol"), row++) {
            int symbol_idx = h % num_symbols;
            reels[reel][row] = symbols[symbol_idx];
        }
    }
    
    return reels;
}

double TigerSlotGame::calculateWin(const std::vector<std::vector<std::string>>& reels, int bet_level) {
    double total_win = 0.0;
    
    // Check paylines (simplified)
    // In real implementation, check all paylines
    
    return total_win * bet_level;
}

// TigerGameFactory Implementation
std::unique_ptr<BaseGame> TigerGameFactory::createGame(const std::string& game_type) {
    if (game_type == "crash") {
        return std::make_unique<TigerCrashGame>();
    } else if (game_type == "mines") {
        return std::make_unique<TigerMinesGame>();
    } else if (game_type == "plinko") {
        return std::make_unique<TigerPlinkoGame>();
    } else if (game_type == "dice") {
        return std::make_unique<TigerDiceGame>();
    } else if (game_type == "slot") {
        return std::make_unique<TigerSlotGame>();
    }
    
    return nullptr;
}

} // namespace TigerCasino
