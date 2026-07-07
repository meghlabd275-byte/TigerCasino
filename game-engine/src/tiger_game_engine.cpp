/**
 * TigerCasino Ultra-Low Latency Game Engine
 * Implementation of provably fair casino games
 * 
 * Copyright (c) 2024 TigerCasino
 * Distributed under MIT License
 */

#include "tiger_game_engine.h"
#include <openssl/hmac.h>
#include <openssl/sha.h>
#include <openssl/rand.h>
#include <sstream>
#include <iomanip>
#include <cmath>
#include <algorithm>
#include <random>
#include <chrono>

namespace TigerCasino {

// Slot game symbols with payouts
const std::map<std::string, double> SLOT_SYMBOLS = {
    {"Cherry", 2.0},
    {"Lemon", 2.0},
    {"Orange", 3.0},
    {"Grape", 4.0},
    {"Watermelon", 5.0},
    {"Bell", 10.0},
    {"Diamond", 20.0},
    {"Seven", 50.0},
    {"Jackpot", 100.0}
};

// Hi-Lo card deck
const std::vector<std::string> HILO_DECK = {
    "A♠", "2♠", "3♠", "4♠", "5♠", "6♠", "7♠", "8♠", "9♠", "10♠", "J♠", "Q♠", "K♠",
    "A♥", "2♥", "3♥", "4♥", "5♥", "6♥", "7♥", "8♥", "9♥", "10♥", "J♥", "Q♥", "K♥",
    "A♦", "2♦", "3♦", "4♦", "5♦", "6♦", "7♦", "8♦", "9♦", "10♦", "J♦", "Q♦", "K♦",
    "A♣", "2♣", "3♣", "4♣", "5♣", "6♣", "7♣", "8♣", "9♣", "10♣", "J♣", "Q♣", "K♣"
};

// Card values for hi-lo (A=1, 2-10=face value, J=11, Q=12, K=13)
const std::map<std::string, int> HILO_VALUES = {
    {"A♠", 1}, {"2♠", 2}, {"3♠", 3}, {"4♠", 4}, {"5♠", 5}, {"6♠", 6}, {"7♠", 7}, {"8♠", 8}, {"9♠", 9}, {"10♠", 10}, {"J♠", 11}, {"Q♠", 12}, {"K♠", 13},
    {"A♥", 1}, {"2♥", 2}, {"3♥", 3}, {"4♥", 4}, {"5♥", 5}, {"6♥", 6}, {"7♥", 7}, {"8♥", 8}, {"9♥", 9}, {"10♥", 10}, {"J♥", 11}, {"Q♥", 12}, {"K♥", 13},
    {"A♦", 1}, {"2♦", 2}, {"3♦", 3}, {"4♦", 4}, {"5♦", 5}, {"6♦", 6}, {"7♦", 7}, {"8♦", 8}, {"9♦", 9}, {"10♦", 10}, {"J♦", 11}, {"Q♦", 12}, {"K♦", 13},
    {"A♣", 1}, {"2♣", 2}, {"3♣", 3}, {"4♣", 4}, {"5♣", 5}, {"6♣", 6}, {"7♣", 7}, {"8♣", 8}, {"9♣", 9}, {"10♣", 10}, {"J♣", 11}, {"Q♣", 12}, {"K♣", 13}
};

TigerGameEngine::TigerGameEngine() : round_counter_(0), total_rounds_played_(0), total_winnings_(0) {
    // Initialize entropy pool
    RAND_bytes(entropy_pool_.data(), SEED_SIZE);
}

TigerGameEngine::~TigerGameEngine() = default;

// ===== PROVABLY FAIR SYSTEM =====

Seeds TigerGameEngine::generateSeeds(const std::string& client_seed) {
    Seeds seeds;
    
    // Generate cryptographically secure random server seed
    std::vector<uint8_t> server_seed_bytes(SEED_SIZE);
    RAND_bytes(server_seed_bytes.data(), SEED_SIZE);
    seeds.server_seed = bytesToHex(server_seed_bytes);
    seeds.server_seed_hash = hashSeed(seeds.server_seed);
    seeds.client_seed = client_seed.empty() ? bytesToHex(generateRandomBytes(16)) : client_seed;
    seeds.nonce = 0;
    
    return seeds;
}

std::string TigerGameEngine::hashSeed(const std::string& seed) {
    uint8_t hash[SHA256_DIGEST_LENGTH];
    SHA256_CTX ctx;
    SHA256_Init(&ctx);
    SHA256_Update(&ctx, seed.c_str(), seed.length());
    SHA256_Final(hash, &ctx);
    
    std::vector<uint8_t> hash_vec(hash, hash + SHA256_DIGEST_LENGTH);
    return bytesToHex(hash_vec);
}

bool TigerGameEngine::verifySeed(const std::string& server_seed, const std::string& server_seed_hash) {
    return hashSeed(server_seed) == server_seed_hash;
}

int TigerGameEngine::generateOutcome(const std::string& server_seed,
                                    const std::string& client_seed,
                                    int nonce,
                                    int max) {
    if (max <= 0) return 0;
    
    // Create HMAC-SHA256 of seed + nonce
    std::string data = server_seed + ":" + client_seed + ":" + std::to_string(nonce);
    std::string hmac = hmacSha256(server_seed, data);
    std::vector<uint8_t> hash = hexToBytes(hmac);
    
    // Convert first 8 bytes to uint64
    uint64_t value = 0;
    for (size_t i = 0; i < 8 && i < hash.size(); i++) {
        value = (value << 8) | hash[i];
    }
    
    // Use modulo to get result in range [0, max)
    return static_cast<int>(value % max);
}

double TigerGameEngine::generateFloatOutcome(const std::string& server_seed,
                                            const std::string& client_seed,
                                            int nonce) {
    // Generate a 64-bit random number and normalize to [0, 1)
    std::string data = server_seed + ":" + client_seed + ":" + std::to_string(nonce);
    std::string hmac = hmacSha256(server_seed, data);
    std::vector<uint8_t> hash = hexToBytes(hmac);
    
    uint64_t value = 0;
    for (size_t i = 0; i < 8 && i < hash.size(); i++) {
        value = (value << 8) | hash[i];
    }
    
    // Convert to [0, 1) using division by 2^64
    static const double DIVISOR = 18446744073709551616.0;
    return static_cast<double>(value) / DIVISOR;
}

std::vector<uint8_t> TigerGameEngine::generateRandomBytes(size_t length) {
    std::vector<uint8_t> bytes(length);
    RAND_bytes(bytes.data(), length);
    return bytes;
}

// ===== DICE GAME =====

DiceResult TigerGameEngine::playDice(const std::string& user_id,
                                      const std::string& server_seed,
                                      const std::string& client_seed,
                                      int nonce,
                                      double bet_amount,
                                      double target,
                                      bool over) {
    DiceResult result;
    
    // Generate roll (0-100)
    double roll = generateFloatOutcome(server_seed, client_seed, nonce) * 100.0;
    
    // Calculate multiplier based on target and direction
    double raw_multiplier;
    bool is_win;
    
    if (over) {
        is_win = roll > target;
        raw_multiplier = is_win ? 100.0 / (100.0 - target) : 0.0;
    } else {
        is_win = roll < target;
        raw_multiplier = is_win ? target / (100.0 - target) : 0.0;
    }
    
    // Apply house edge
    double multiplier = raw_multiplier * (1.0 - HOUSE_EDGE_DICE / 10000.0);
    
    result.roll = roll;
    result.target = target;
    result.multiplier = multiplier;
    result.direction_over = over;
    result.is_win = is_win;
    result.win_amount = is_win ? bet_amount * multiplier : 0.0;
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    result.nonce = nonce;
    result.verified = verifySeed(server_seed, hashSeed(server_seed));
    
    // Update statistics
    total_rounds_played_.fetch_add(1);
    if (is_win) {
        total_winnings_.fetch_add(static_cast<uint64_t>(result.win_amount * 100));
    }
    
    return result;
}

// ===== CRASH GAME =====

CrashResult TigerGameEngine::startCrashRound(const std::string& round_id,
                                              const std::string& server_seed,
                                              const std::string& server_seed_hash,
                                              const std::string& client_seed,
                                              int nonce) {
    CrashResult result;
    result.round_id = round_id;
    
    // Generate crash point using exponential distribution for realistic curve
    double f = generateFloatOutcome(server_seed, client_seed, nonce);
    double crash_point;
    
    if (f < 0.70) {
        // 70% of crashes happen between 1.0x and 8.0x
        crash_point = 1.0 + (f * 10.0);
    } else {
        // 30% can go much higher (up to 1000x in rare cases, but typically under 100x)
        crash_point = 8.0 + ((f - 0.70) * 300.0);
    }
    
    // Apply house edge (reduce crash point slightly)
    crash_point = crash_point * (1.0 - HOUSE_EDGE_CRASH / 10000.0);
    
    result.crash_point = crash_point;
    result.did_crash = false;
    result.did_cashout = false;
    result.server_seed = server_seed;
    result.server_seed_hash = server_seed_hash;
    result.client_seed = client_seed;
    result.nonce = nonce;
    
    // Store round
    std::lock_guard<std::mutex> lock(games_mutex_);
    crash_rounds_[round_id] = result;
    
    return result;
}

double TigerGameEngine::getCrashMultiplier(const std::string& round_id, uint64_t elapsed_ms) {
    // Crash multiplier increases exponentially over time
    // Base formula: 1 + (elapsed_ms / 1000)^1.5
    // This creates a smooth curve that accelerates
    
    std::lock_guard<std::mutex> lock(games_mutex_);
    auto it = crash_rounds_.find(round_id);
    if (it == crash_rounds_.end()) {
        return 1.0;
    }
    
    double elapsed_seconds = elapsed_ms / 1000.0;
    double multiplier = 1.0 + std::pow(elapsed_seconds, 1.5);
    
    // Cap at crash point
    if (multiplier > it->second.crash_point) {
        multiplier = it->second.crash_point;
    }
    
    return multiplier;
}

CrashResult TigerGameEngine::crashCashout(const std::string& round_id,
                                          const std::string& user_id,
                                          double bet_amount,
                                          double multiplier) {
    std::lock_guard<std::mutex> lock(games_mutex_);
    auto it = crash_rounds_.find(round_id);
    
    CrashResult result;
    if (it == crash_rounds_.end()) {
        return result; // Round not found
    }
    
    // Check if already crashed
    if (multiplier >= it->second.crash_point) {
        result.did_crash = true;
        result.did_cashout = false;
        return result;
    }
    
    // Calculate winnings
    double win_amount = bet_amount * multiplier;
    
    it->second.did_cashout = true;
    it->second.player_cashout = multiplier;
    it->second.win_amount = win_amount;
    it->second.bet_amount = bet_amount;
    it->second.multiplier = multiplier;
    
    // Update statistics
    total_rounds_played_.fetch_add(1);
    total_winnings_.fetch_add(static_cast<uint64_t>(win_amount * 100));
    
    return it->second;
}

CrashResult TigerGameEngine::crashRound(const std::string& round_id) {
    std::lock_guard<std::mutex> lock(games_mutex_);
    auto it = crash_rounds_.find(round_id);
    
    CrashResult result;
    if (it == crash_rounds_.end()) {
        return result;
    }
    
    it->second.did_crash = true;
    return it->second;
}

// ===== MINES GAME =====

MinesResult TigerGameEngine::startMinesGame(const std::string& game_id,
                                             const std::string& user_id,
                                             const std::string& server_seed,
                                             const std::string& client_seed,
                                             int nonce,
                                             double bet_amount,
                                             int mines_count) {
    MinesResult result;
    result.game_id = game_id;
    result.mines_count = mines_count;
    result.bet_amount = bet_amount;
    result.current_step = 0;
    result.game_over = false;
    result.multiplier = 1.0;
    result.win_amount = bet_amount;
    
    // Generate mine positions
    result.mine_locations.resize(MINES_MAX_GRID, false);
    std::vector<int> mine_positions;
    
    while (mine_positions.size() < static_cast<size_t>(mines_count)) {
        int pos = generateOutcome(server_seed, client_seed, nonce + mine_positions.size(), MINES_MAX_GRID);
        if (!result.mine_locations[pos]) {
            result.mine_locations[pos] = true;
            mine_positions.push_back(pos);
        }
    }
    
    // Store game state
    std::lock_guard<std::mutex> lock(games_mutex_);
    mines_games_[game_id] = std::make_shared<MinesResult>(result);
    
    return result;
}

MinesResult TigerGameEngine::revealMinesTile(const std::string& game_id,
                                              const std::string& user_id,
                                              int tile) {
    std::lock_guard<std::mutex> lock(games_mutex_);
    auto it = mines_games_.find(game_id);
    
    MinesResult result;
    if (it == mines_games_.end()) {
        return result; // Game not found
    }
    
    result = *it->second;
    
    // Check if tile already revealed
    for (int revealed : result.revealed_tiles) {
        if (revealed == tile) {
            return result; // Already revealed
        }
    }
    
    result.revealed_tile = tile;
    result.revealed_tiles.push_back(tile);
    result.current_step++;
    
    // Check if mine
    if (tile >= 0 && tile < MINES_MAX_GRID && result.mine_locations[tile]) {
        result.is_mine = true;
        result.game_over = true;
        result.win_amount = 0;
        result.multiplier = 0;
    } else {
        // Calculate new multiplier
        result.multiplier = calculateMinesMultiplier(
            result.current_step,
            result.mines_count,
            result.bet_amount
        );
        result.win_amount = result.bet_amount * result.multiplier;
    }
    
    // Update game state
    *it->second = result;
    
    return result;
}

double TigerGameEngine::calculateMinesMultiplier(int revealed, int mines, double bet) {
    // Multiplier increases as you reveal more safe tiles
    // Formula: multiplier = 1 + (revealed * 0.1) * (mines / 10)
    // Then apply house edge
    double base_multiplier = 1.0 + (revealed * 0.1) * (mines / 10.0);
    double with_house_edge = base_multiplier * (1.0 - HOUSE_EDGE_MINES / 10000.0);
    return with_house_edge;
}

MinesResult TigerGameEngine::cashoutMines(const std::string& game_id, const std::string& user_id) {
    std::lock_guard<std::mutex> lock(games_mutex_);
    auto it = mines_games_.find(game_id);
    
    MinesResult result;
    if (it == mines_games_.end()) {
        return result;
    }
    
    result = *it->second;
    result.game_over = true;
    
    // Update statistics
    total_rounds_played_.fetch_add(1);
    if (result.win_amount > result.bet_amount) {
        total_winnings_.fetch_add(static_cast<uint64_t>((result.win_amount - result.bet_amount) * 100));
    }
    
    return result;
}

// ===== PLINKO GAME =====

std::vector<double> TigerGameEngine::getPlinkoPayouts(int rows, const std::string& risk) {
    // Payout multipliers for each bucket position
    // Center buckets have higher multipliers in high risk
    std::vector<double> payouts;
    
    if (risk == "low") {
        // Low risk: Center has 1.5x, edges have lower
        payouts = {1.5, 1.2, 1.0, 0.8, 0.5, 0.5, 0.8, 1.0, 1.2, 1.5};
    } else if (risk == "high") {
        // High risk: Center has much higher
        payouts = {10.0, 5.0, 2.0, 1.0, 0.5, 0.5, 1.0, 2.0, 5.0, 10.0};
    } else {
        // Medium (default)
        payouts = {5.0, 2.5, 1.5, 1.0, 0.5, 0.5, 1.0, 1.5, 2.5, 5.0};
    }
    
    // Extend or shrink to match rows
    while (payouts.size() < static_cast<size_t>(rows + 1)) {
        payouts.insert(payouts.begin() + payouts.size() / 2, 1.0);
    }
    
    return payouts;
}

PlinkoResult TigerGameEngine::playPlinko(const std::string& user_id,
                                         const std::string& server_seed,
                                         const std::string& client_seed,
                                         int nonce,
                                         double bet_amount,
                                         int rows,
                                         const std::string& risk) {
    PlinkoResult result;
    result.rows = rows;
    result.risk = risk;
    result.bet_amount = bet_amount;
    result.path.resize(rows + 1);
    
    // Simulate ball path through plinko board
    int position = 0; // Start in middle
    result.path[0] = position;
    
    for (int i = 0; i < rows; i++) {
        // Each step, ball goes left or right
        double rand_val = generateFloatOutcome(server_seed, client_seed, nonce + i);
        if (rand_val > 0.5) {
            position++;
        }
        result.path[i + 1] = position;
    }
    
    // Get final bucket position
    result.final_bucket = position;
    
    // Get payout
    std::vector<double> payouts = getPlinkoPayouts(rows, risk);
    int payout_index = position;
    if (payout_index >= static_cast<int>(payouts.size())) {
        payout_index = payouts.size() - 1;
    }
    
    result.multiplier = payouts[payout_index];
    result.win_amount = bet_amount * result.multiplier * (1.0 - HOUSE_EDGE_PLINKO / 10000.0);
    
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    result.nonce = nonce;
    
    // Update statistics
    total_rounds_played_.fetch_add(1);
    if (result.win_amount > bet_amount) {
        total_winnings_.fetch_add(static_cast<uint64_t>((result.win_amount - bet_amount) * 100));
    }
    
    return result;
}

// ===== SLOTS GAME =====

std::string TigerGameEngine::generateSlotSymbols(const std::string& seed, int reel) {
    // Generate symbols for a single reel using seeded random
    static const char* symbols[] = {
        "Cherry", "Lemon", "Orange", "Grape", "Watermelon",
        "Bell", "Diamond", "Seven", "Jackpot"
    };
    static const int weights[] = {30, 25, 20, 15, 10, 5, 3, 1, 0.5};
    static const int total_weight = 134.5;
    
    int outcome = generateOutcome(seed, seed + std::to_string(reel), reel, 1000);
    double value = outcome / 10.0; // Convert to percentage
    
    double cumulative = 0;
    for (int i = 0; i < 9; i++) {
        cumulative += weights[i] / total_weight * 100;
        if (value <= cumulative) {
            return symbols[i];
        }
    }
    
    return symbols[0];
}

SlotsResult TigerGameEngine::playSlots(const std::string& user_id,
                                        const std::string& server_seed,
                                        const std::string& client_seed,
                                        int nonce,
                                        double bet_amount,
                                        int lines) {
    SlotsResult result;
    result.bet_amount = bet_amount;
    
    // Generate reel outcomes
    for (int i = 0; i < SLOT_REELS; i++) {
        std::string symbols = generateSlotSymbols(server_seed + client_seed, nonce + i);
        
        // Map to index
        if (symbols == "Cherry") result.reels[i] = 0;
        else if (symbols == "Lemon") result.reels[i] = 1;
        else if (symbols == "Orange") result.reels[i] = 2;
        else if (symbols == "Grape") result.reels[i] = 3;
        else if (symbols == "Watermelon") result.reels[i] = 4;
        else if (symbols == "Bell") result.reels[i] = 5;
        else if (symbols == "Diamond") result.reels[i] = 6;
        else if (symbols == "Seven") result.reels[i] = 7;
        else result.reels[i] = 8;
    }
    
    // Check for wins
    bool is_jackpot = (result.reels[0] == 8) && (result.reels[1] == 8) && (result.reels[2] == 8);
    bool is_three_of_kind = (result.reels[0] == result.reels[1]) && (result.reels[1] == result.reels[2]);
    bool is_two_of_kind = (result.reels[0] == result.reels[1]) || (result.reels[1] == result.reels[2]) || (result.reels[0] == result.reels[2]);
    
    result.is_jackpot = is_jackpot;
    
    if (is_jackpot) {
        result.multiplier = 100.0;
        result.win_line = "JACKPOT!";
    } else if (is_three_of_kind) {
        result.multiplier = 10.0;
        result.win_line = "Three of a Kind!";
    } else if (is_two_of_kind) {
        result.multiplier = 2.0;
        result.win_line = "Two of a Kind!";
    } else {
        result.multiplier = 0.0;
        result.win_line = "";
    }
    
    // Apply house edge
    result.multiplier = result.multiplier * (1.0 - HOUSE_EDGE_SLOTS / 10000.0);
    result.win_amount = bet_amount * result.multiplier * lines;
    
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    result.nonce = nonce;
    
    // Update statistics
    total_rounds_played_.fetch_add(1);
    if (result.win_amount > 0) {
        total_winnings_.fetch_add(static_cast<uint64_t>(result.win_amount * 100));
    }
    
    return result;
}

// ===== LIMBO GAME =====

LimboResult TigerGameEngine::playLimbo(const std::string& user_id,
                                       const std::string& server_seed,
                                       const std::string& client_seed,
                                       int nonce,
                                       double bet_amount,
                                       double target_multiplier) {
    LimboResult result;
    result.target_multiplier = target_multiplier;
    result.bet_amount = bet_amount;
    
    // Generate result multiplier using exponential distribution
    double f = generateFloatOutcome(server_seed, client_seed, nonce);
    
    // Use inverse transform sampling for exponential-like distribution
    // Higher multipliers are less likely
    double result_multiplier = 1.0 / (1.0 - f);
    
    // Cap at reasonable maximum
    if (result_multiplier > 1000.0) {
        result_multiplier = 1000.0;
    }
    
    result.result_multiplier = result_multiplier;
    result.is_win = result_multiplier > target_multiplier;
    
    // Calculate multiplier (target / result, then apply house edge)
    if (result.is_win) {
        double raw_multiplier = result_multiplier / target_multiplier;
        result.multiplier = raw_multiplier * (1.0 - HOUSE_EDGE_LIMBO / 10000.0);
        result.win_amount = bet_amount * result.multiplier;
    } else {
        result.multiplier = 0;
        result.win_amount = 0;
    }
    
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    result.nonce = nonce;
    
    // Update statistics
    total_rounds_played_.fetch_add(1);
    if (result.is_win) {
        total_winnings_.fetch_add(static_cast<uint64_t>(result.win_amount * 100));
    }
    
    return result;
}

// ===== HI-LO GAME =====

HiloResult TigerGameEngine::startHiloGame(const std::string& game_id,
                                          const std::string& user_id,
                                          const std::string& server_seed,
                                          const std::string& client_seed,
                                          int nonce,
                                          double bet_amount) {
    HiloResult result;
    
    // Draw first card
    int card_index = generateOutcome(server_seed, client_seed, nonce, 52);
    result.current_card = HILO_DECK[card_index];
    result.bet_amount = bet_amount;
    result.multiplier = 1.0;
    result.win_amount = bet_amount;
    result.streak = 0;
    result.is_correct = true;
    
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    result.nonce = nonce;
    
    return result;
}

HiloResult TigerGameEngine::playHiloChoice(const std::string& game_id,
                                            const std::string& user_id,
                                            const std::string& choice) {
    HiloResult result;
    // This would need game state storage - simplified version
    result.is_correct = false;
    return result;
}

// ===== UTILITY FUNCTIONS =====

double TigerGameEngine::calculateMultiplier(const std::string& game_type,
                                           const std::map<std::string, std::string>& params) {
    if (game_type == "dice") {
        double target = std::stod(params.at("target"));
        bool over = params.at("direction") == "over";
        return over ? 100.0 / (100.0 - target) : target / (100.0 - target);
    } else if (game_type == "limbo") {
        double target = std::stod(params.at("target"));
        return target;
    }
    return 1.0;
}

std::map<std::string, uint64_t> TigerGameEngine::getStatistics() {
    std::lock_guard<std::mutex> lock(stats_mutex_);
    std::map<std::string, uint64_t> stats;
    stats["total_rounds"] = total_rounds_played_.load();
    stats["total_winnings"] = total_winnings_.load();
    stats["active_games"] = mines_games_.size() + crash_rounds_.size();
    return stats;
}

void TigerGameEngine::resetStatistics() {
    total_rounds_played_.store(0);
    total_winnings_.store(0);
}

int TigerGameEngine::getHouseEdge(const std::string& game_type) {
    if (game_type == "dice") return HOUSE_EDGE_DICE;
    if (game_type == "crash") return HOUSE_EDGE_CRASH;
    if (game_type == "mines") return HOUSE_EDGE_MINES;
    if (game_type == "plinko") return HOUSE_EDGE_PLINKO;
    if (game_type == "slots") return HOUSE_EDGE_SLOTS;
    if (game_type == "limbo") return HOUSE_EDGE_LIMBO;
    if (game_type == "hilo") return HOUSE_EDGE_HILO;
    return 500; // Default 5%
}

// ===== PRIVATE HELPERS =====

std::string TigerGameEngine::hmacSha256(const std::string& key, const std::string& data) {
    unsigned char hmac[HMAC_MAX_MD_SIZE];
    unsigned int hmac_len;
    
    HMAC(EVP_sha256(),
         key.c_str(), key.length(),
         reinterpret_cast<const unsigned char*>(data.c_str()), data.length(),
         hmac, &hmac_len);
    
    std::vector<uint8_t> hmac_vec(hmac, hmac + hmac_len);
    return bytesToHex(hmac_vec);
}

std::string TigerGameEngine::bytesToHex(const std::vector<uint8_t>& bytes) {
    std::stringstream ss;
    ss << std::hex << std::setfill('0');
    for (uint8_t byte : bytes) {
        ss << std::setw(2) << static_cast<int>(byte);
    }
    return ss.str();
}

std::vector<uint8_t> TigerGameEngine::hexToBytes(const std::string& hex) {
    std::vector<uint8_t> bytes;
    for (size_t i = 0; i < hex.length(); i += 2) {
        std::string byte_str = hex.substr(i, 2);
        uint8_t byte = static_cast<uint8_t>(std::stoi(byte_str, nullptr, 16));
        bytes.push_back(byte);
    }
    return bytes;
}

} // namespace TigerCasino
