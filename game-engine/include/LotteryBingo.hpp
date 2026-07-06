#pragma once

#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"
#include <string>
#include <vector>
#include <memory>
#include <array>
#include <set>
#include <map>
#include <algorithm>

namespace TigerCasino {

/**
 * Lottery, Bingo, and Scratch Cards games
 */

/**
 * Keno - Classic lottery-style game
 */
class KenoGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 1000.00;
    static constexpr double BASE_RTP = 0.95;
    static constexpr size_t NUM_SPOTS = 10;
    static constexpr size_t MAX_SPOTS = 10;
    static constexpr size_t DRAW_NUMBERS = 20;

    struct GameResult {
        std::string player_id;
        std::set<uint8_t> selected_numbers;
        std::set<uint8_t> drawn_numbers;
        std::set<uint8_t> matched_numbers;
        size_t hits;
        double bet_amount;
        double payout;
        double multiplier;
        std::string server_seed;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t game_counter_;

    // Paytable: hits -> multiplier for 10 spots
    static constexpr std::array<double, 11> PAYTABLE_10 = {{
        0.0, 0.0, 0.0, 1.0, 2.0, 10.0, 50.0, 100.0, 500.0, 1000.0, 10000.0
    }};

    // Paytable for different spot counts
    static constexpr std::array<double, 11> PAYTABLE_9 = {{
        0.0, 0.0, 0.0, 1.0, 2.0, 5.0, 20.0, 80.0, 300.0, 600.0, 5000.0
    }};

    static constexpr std::array<double, 11> PAYTABLE_8 = {{
        0.0, 0.0, 0.5, 1.0, 2.0, 5.0, 15.0, 40.0, 100.0, 300.0, 1000.0
    }};

    static constexpr std::array<double, 11> PAYTABLE_7 = {{
        0.0, 0.0, 0.5, 1.0, 2.0, 5.0, 15.0, 30.0, 80.0, 200.0, 500.0
    }};

    static constexpr std::array<double, 11> PAYTABLE_6 = {{
        0.0, 0.5, 0.5, 1.0, 2.0, 5.0, 10.0, 25.0, 50.0, 100.0, 300.0
    }};

    static constexpr std::array<double, 11> PAYTABLE_5 = {{
        0.5, 0.5, 1.0, 1.0, 2.0, 5.0, 10.0, 15.0, 30.0, 50.0, 100.0
    }};

public:
    KenoGame();
    
    std::string getName() const { return "Keno"; }
    std::string getType() const { return "Lottery"; }
    double getRTP() const { return BASE_RTP; }

    GameResult play(const std::string& player_id, 
                   const std::set<uint8_t>& selected_numbers, 
                   double bet_amount);
    
    std::vector<uint8_t> drawNumbers();
    std::set<uint8_t> getRandomNumbers(size_t count);
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    const std::array<double, 11>& getPaytable(size_t num_spots) const;
};

inline KenoGame::KenoGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , game_counter_(0) {}

inline const std::array<double, 11>& KenoGame::getPaytable(size_t num_spots) const {
    switch (num_spots) {
        case 10: return PAYTABLE_10;
        case 9: return PAYTABLE_9;
        case 8: return PAYTABLE_8;
        case 7: return PAYTABLE_7;
        case 6: return PAYTABLE_6;
        default: return PAYTABLE_5;
    }
}

inline std::set<uint8_t> KenoGame::getRandomNumbers(size_t count) {
    std::set<uint8_t> numbers;
    while (numbers.size() < count && numbers.size() < 80) {
        uint8_t num = static_cast<uint8_t>(rng_->generateInt(1, 80));
        numbers.insert(num);
    }
    return numbers;
}

inline std::vector<uint8_t> KenoGame::drawNumbers() {
    auto numbers = getRandomNumbers(DRAW_NUMBERS);
    return std::vector<uint8_t>(numbers.begin(), numbers.end());
}

inline KenoGame::GameResult KenoGame::play(const std::string& player_id,
                                           const std::set<uint8_t>& selected_numbers,
                                           double bet_amount) {
    GameResult result;
    game_counter_++;
    
    result.player_id = player_id;
    result.selected_numbers = selected_numbers;
    result.bet_amount = bet_amount;
    result.server_seed = provably_fair_->getServerSeed();
    
    // Validate selection
    if (selected_numbers.size() < 1 || selected_numbers.size() > MAX_SPOTS) {
        result.hits = 0;
        result.payout = 0;
        return result;
    }
    
    // Draw numbers
    result.drawn_numbers = getRandomNumbers(DRAW_NUMBERS);
    
    // Find matches
    std::set<uint8_t> intersection;
    std::set_intersection(selected_numbers.begin(), selected_numbers.end(),
                         result.drawn_numbers.begin(), result.drawn_numbers.end(),
                         std::inserter(intersection, intersection.begin()));
    
    result.matched_numbers = intersection;
    result.hits = intersection.size();
    
    // Calculate payout
    const auto& paytable = getPaytable(selected_numbers.size());
    if (result.hits < paytable.size()) {
        result.multiplier = paytable[result.hits];
    } else {
        result.multiplier = 0;
    }
    
    result.payout = bet_amount * result.multiplier;
    
    return result;
}

inline void KenoGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void KenoGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Bingo - Classic 75-ball bingo
 */
class BingoGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 100.00;
    static constexpr double BASE_RTP = 0.95;
    static constexpr size_t GRID_SIZE = 5; // 5x5 grid
    static constexpr uint8_t FREE_SPACE = 0;

    struct Card {
        uint8_t grid[GRID_SIZE][GRID_SIZE];
        std::string card_id;
    };

    struct GameResult {
        std::string player_id;
        std::string card_id;
        Card bingo_card;
        std::vector<uint8_t> called_numbers;
        std::vector<std::pair<size_t, size_t>> winning_pattern;
        std::string pattern_name;
        size_t calls_to_win;
        double bet_amount;
        double payout;
        std::string server_seed;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t game_counter_;
    std::vector<uint8_t> called_numbers_;
    std::set<uint8_t> remaining_numbers_;

    // Patterns: vector of (row, col) positions
    std::map<std::string, std::vector<std::pair<size_t, size_t>>> patterns_;

public:
    BingoGame();
    
    std::string getName() const { return "Bingo"; }
    std::string getType() const { return "Bingo"; }
    double getRTP() const { return BASE_RTP; }

    Card generateCard(const std::string& card_id);
    uint8_t callNumber();
    GameResult checkWin(const std::string& player_id, const Card& card, double bet_amount);
    
    void startNewGame();
    const std::vector<uint8_t>& getCalledNumbers() const;
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    bool checkPattern(const Card& card, const std::vector<std::pair<size_t, size_t>>& pattern) const;
};

inline BingoGame::BingoGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , game_counter_(0) {
    // Initialize remaining numbers (1-75)
    for (uint8_t i = 1; i <= 75; ++i) {
        remaining_numbers_.insert(i);
    }
    
    // Define winning patterns
    // Horizontal lines
    for (size_t row = 0; row < GRID_SIZE; ++row) {
        std::vector<std::pair<size_t, size_t>> pattern;
        for (size_t col = 0; col < GRID_SIZE; ++col) {
            pattern.push_back({row, col});
        }
        patterns_["row_" + std::to_string(row + 1)] = pattern;
    }
    
    // Vertical lines
    for (size_t col = 0; col < GRID_SIZE; ++col) {
        std::vector<std::pair<size_t, size_t>> pattern;
        for (size_t row = 0; row < GRID_SIZE; ++row) {
            pattern.push_back({row, col});
        }
        patterns_["col_" + std::to_string(col + 1)] = pattern;
    }
    
    // Diagonals
    std::vector<std::pair<size_t, size_t>> diag1, diag2;
    for (size_t i = 0; i < GRID_SIZE; ++i) {
        diag1.push_back({i, i});
        diag2.push_back({i, GRID_SIZE - 1 - i});
    }
    patterns_["diagonal_1"] = diag1;
    patterns_["diagonal_2"] = diag2;
}

inline BingoGame::Card BingoGame::generateCard(const std::string& card_id) {
    Card card;
    card.card_id = card_id;
    
    // Each column (B-I-N-G-O) has numbers 1-15, 16-30, 31-45, 46-60, 61-75
    for (size_t col = 0; col < GRID_SIZE; ++col) {
        std::set<uint8_t> column_numbers;
        uint8_t min_num = static_cast<uint8_t>(col * 15 + 1);
        uint8_t max_num = static_cast<uint8_t>((col + 1) * 15);
        
        while (column_numbers.size() < GRID_SIZE) {
            uint8_t num = static_cast<uint8_t>(rng_->generateInt(min_num, max_num));
            column_numbers.insert(num);
        }
        
        size_t row = 0;
        for (uint8_t num : column_numbers) {
            card.grid[row][col] = num;
            row++;
        }
    }
    
    // Set free space in center
    card.grid[2][2] = FREE_SPACE;
    
    return card;
}

inline uint8_t BingoGame::callNumber() {
    if (remaining_numbers_.empty()) {
        return 0;
    }
    
    auto it = remaining_numbers_.begin();
    std::advance(it, rng_->generateInt(0, remaining_numbers_.size() - 1));
    
    uint8_t number = *it;
    remaining_numbers_.erase(it);
    called_numbers_.push_back(number);
    
    return number;
}

inline bool BingoGame::checkPattern(const Card& card, 
                                    const std::vector<std::pair<size_t, size_t>>& pattern) const {
    for (const auto& [row, col] : pattern) {
        if (card.grid[row][col] != FREE_SPACE) {
            // Check if this number has been called
            bool called = false;
            for (uint8_t num : called_numbers_) {
                if (num == card.grid[row][col]) {
                    called = true;
                    break;
                }
            }
            if (!called) return false;
        }
    }
    return true;
}

inline BingoGame::GameResult BingoGame::checkWin(const std::string& player_id, 
                                                const Card& card, 
                                                double bet_amount) {
    GameResult result;
    result.player_id = player_id;
    result.card_id = card.card_id;
    result.bingo_card = card;
    result.called_numbers = called_numbers_;
    result.bet_amount = bet_amount;
    result.server_seed = provably_fair_->getServerSeed();
    result.calls_to_win = called_numbers_.size();
    
    // Check all patterns
    for (const auto& [pattern_name, pattern] : patterns_) {
        if (checkPattern(card, pattern)) {
            result.pattern_name = pattern_name;
            result.winning_pattern = pattern;
            
            // Calculate payout based on how early they won
            double base_payout = 100.0;
            double call_bonus = std::max(0.0, 50.0 - static_cast<double>(called_numbers_.size()));
            result.payout = bet_amount * (base_payout + call_bonus) / 10.0;
            
            return result;
        }
    }
    
    result.pattern_name = "none";
    result.payout = 0;
    
    return result;
}

inline void BingoGame::startNewGame() {
    game_counter_++;
    called_numbers_.clear();
    
    // Reset remaining numbers
    remaining_numbers_.clear();
    for (uint8_t i = 1; i <= 75; ++i) {
        remaining_numbers_.insert(i);
    }
}

inline const std::vector<uint8_t>& BingoGame::getCalledNumbers() const {
    return called_numbers_;
}

inline void BingoGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void BingoGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Scratch Cards - Instant win scratch-off games
 */
class ScratchCardGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 100.00;
    static constexpr double BASE_RTP = 0.94;

    struct Card {
        std::string game_id;
        std::vector<uint8_t> symbols;
        std::vector<uint8_t> prize_symbols;
        bool revealed;
    };

    struct GameResult {
        std::string player_id;
        std::string game_id;
        Card scratch_card;
        double bet_amount;
        double payout;
        bool is_winner;
        std::string prize_name;
        std::string server_seed;
    };

private:
    enum class ScratchGameType {
        LUCKY_7S,
        GOLD_RUSH,
        DIAMOND_DAZZLE,
        CASH_COIN,
        MEGA_WIN
    };

    struct GameConfig {
        std::string name;
        std::vector<std::pair<std::string, uint8_t>> prize_table; // prize name, symbol
        std::vector<size_t> match_counts; // how many matches needed to win
    };

    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t game_counter_;
    std::map<ScratchGameType, GameConfig> game_configs_;

public:
    ScratchCardGame();
    
    std::string getName() const { return "Scratch Cards"; }
    std::string getType() const { return "Instant Win"; }
    double getRTP() const { return BASE_RTP; }

    Card createCard(ScratchGameType game_type);
    GameResult scratch(const std::string& player_id, 
                       ScratchGameType game_type,
                       double bet_amount);
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    const GameConfig& getConfig(ScratchGameType type) const;
    double calculatePayout(ScratchGameType type, size_t matches, double bet_amount) const;
};

inline ScratchCardGame::ScratchCardGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , game_counter_(0) {
    // Configure Lucky 7s
    GameConfig lucky7s;
    lucky7s.name = "Lucky 7s";
    lucky7s.prize_table = {
        {"JACKPOT", 1}, {"GRAND", 2}, {"MAJOR", 3}, {"MINOR", 4}, {"MINI", 5}
    };
    lucky7s.match_counts = {3};
    game_configs_[ScratchGameType::LUCKY_7S] = lucky7s;
    
    // Configure Gold Rush
    GameConfig gold_rush;
    gold_rush.name = "Gold Rush";
    gold_rush.prize_table = {
        {"JACKPOT", 7}, {"GRAND", 6}, {"MAJOR", 5}, {"MINOR", 4}, {"MINI", 3}
    };
    gold_rush.match_counts = {3};
    game_configs_[ScratchGameType::GOLD_RUSH] = gold_rush;
    
    // Configure Diamond Dazzle
    GameConfig diamond;
    diamond.name = "Diamond Dazzle";
    diamond.prize_table = {
        {"JACKPOT", 10}, {"GRAND", 9}, {"MAJOR", 8}, {"MINOR", 7}, {"MINI", 6}
    };
    diamond.match_counts = {3};
    game_configs_[ScratchGameType::DIAMOND_DAZZLE] = diamond;
    
    // Configure Cash Coin
    GameConfig cash_coin;
    cash_coin.name = "Cash Coin";
    cash_coin.prize_table = {
        {"$1,000,000", 1}, {"$100,000", 2}, {"$10,000", 3}, {"$1,000", 4}, {"$100", 5}
    };
    cash_coin.match_counts = {3};
    game_configs_[ScratchGameType::CASH_COIN] = cash_coin;
    
    // Configure Mega Win
    GameConfig mega_win;
    mega_win.name = "Mega Win";
    mega_win.prize_table = {
        {"MEGA JACKPOT", 1}, {"SUPER PRIZE", 2}, {"BIG WIN", 3}, {"WIN", 4}, {"SMALL WIN", 5}
    };
    mega_win.match_counts = {3};
    game_configs_[ScratchGameType::MEGA_WIN] = mega_win;
}

inline const ScratchCardGame::GameConfig& ScratchCardGame::getConfig(ScratchGameType type) const {
    return game_configs_.at(type);
}

inline double ScratchCardGame::calculatePayout(ScratchGameType type, size_t matches, double bet_amount) const {
    const auto& config = getConfig(type);
    
    if (matches < 3) return 0;
    
    // Higher tier prizes have lower probability
    double multiplier = 0;
    switch (matches) {
        case 3: multiplier = 2.0; break;
        case 4: multiplier = 10.0; break;
        case 5: multiplier = 50.0; break;
        default: multiplier = 100.0; break;
    }
    
    return bet_amount * multiplier;
}

inline ScratchCardGame::Card ScratchCardGame::createCard(ScratchGameType game_type) {
    Card card;
    card.game_id = "SCRATCH_" + std::to_string(++game_counter_);
    card.revealed = false;
    
    const auto& config = getConfig(game_type);
    
    // Create symbol pool with weighted distribution
    std::vector<uint8_t> symbol_pool;
    
    // Add winning symbols (lower frequency)
    for (const auto& [name, symbol] : config.prize_table) {
        // Higher prizes have fewer symbols
        size_t count = 1;
        if (symbol <= 2) count = 1;
        else if (symbol <= 4) count = 2;
        else count = 3;
        
        for (size_t i = 0; i < count; ++i) {
            symbol_pool.push_back(symbol);
        }
    }
    
    // Add filler symbols (higher frequency)
    for (uint8_t i = 10; i <= 20; ++i) {
        for (size_t j = 0; j < 3; ++j) {
            symbol_pool.push_back(i);
        }
    }
    
    // Generate 9 scratch areas (3x3 grid)
    for (size_t i = 0; i < 9; ++i) {
        size_t idx = rng_->generateInt(0, symbol_pool.size() - 1);
        card.symbols.push_back(symbol_pool[idx]);
    }
    
    // Determine winning symbols (3 matching symbols)
    size_t winner_idx = rng_->generateInt(0, config.prize_table.size() - 1);
    uint8_t win_symbol = config.prize_table[winner_idx].second;
    
    // Place winning symbol in 3 positions
    std::vector<size_t> positions = {0, 1, 2, 3, 4, 5, 6, 7, 8};
    std::shuffle(positions.begin(), positions.end(), rng_->getEngine());
    
    for (size_t i = 0; i < 3; ++i) {
        card.prize_symbols.push_back(win_symbol);
    }
    
    return card;
}

inline ScratchCardGame::GameResult ScratchCardGame::scratch(const std::string& player_id,
                                                           ScratchGameType game_type,
                                                           double bet_amount) {
    GameResult result;
    result.player_id = player_id;
    result.bet_amount = bet_amount;
    result.server_seed = provably_fair_->getServerSeed();
    
    const auto& config = getConfig(game_type);
    result.scratch_card = createCard(game_type);
    result.game_id = result.scratch_card.game_id;
    
    // Check for winning combination
    // Count matching symbols in prize spots (first 3)
    std::map<uint8_t, size_t> symbol_counts;
    for (size_t i = 0; i < 3; ++i) {
        symbol_counts[result.scratch_card.symbols[i]]++;
    }
    
    size_t max_matches = 0;
    for (const auto& [symbol, count] : symbol_counts) {
        if (count > max_matches) {
            max_matches = count;
        }
    }
    
    result.is_winner = (max_matches >= 3);
    result.payout = calculatePayout(game_type, max_matches, bet_amount);
    
    if (result.is_winner) {
        for (const auto& [name, symbol] : config.prize_table) {
            if (symbol_counts.count(symbol) && symbol_counts[symbol] >= 3) {
                result.prize_name = name;
                break;
            }
        }
    } else {
        result.prize_name = "NO WIN";
    }
    
    return result;
}

inline void ScratchCardGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void ScratchCardGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
