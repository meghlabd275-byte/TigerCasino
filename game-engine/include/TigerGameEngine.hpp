/**
 * TigerCasino Ultra-Low Latency Game Engine
 * C++20 High-Performance Gaming Library
 * 
 * This header defines the core game engine interface for all casino games.
 * Optimized for sub-microsecond processing with lock-free data structures.
 */

#ifndef TIGER_GAME_ENGINE_HPP
#define TIGER_GAME_ENGINE_HPP

#include <atomic>
#include <chrono>
#include <cstdint>
#include <memory>
#include <optional>
#include <string>
#include <string_view>
#include <tuple>
#include <variant>
#include <vector>

namespace tiger {
namespace engine {

// ============ Core Types ============

using GameID = uint64_t;
using UserID = uint64_t;
using BetAmount = double;
using Multiplier = double;
using Timestamp = std::chrono::steady_clock::time_point;

// Card representation
enum class Suit : uint8_t { Hearts, Diamonds, Clubs, Spades };
enum class Rank : uint8_t { 
    Two = 2, Three = 3, Four = 4, Five = 5, Six = 6, Seven = 7,
    Eight = 8, Nine = 9, Ten = 10, Jack = 11, Queen = 12, 
    King = 13, Ace = 14 
};

struct Card {
    Suit suit;
    Rank rank;
    uint8_t value() const { 
        return (rank >= Rank::Ten) ? 10 : static_cast<uint8_t>(rank); 
    }
};

// ============ Game Results ============

struct SpinResult {
    std::vector<std::vector<std::string>> reels;
    std::vector<std::tuple<std::vector<uint8_t>, std::string, uint8_t, double>> paylines;
    double total_win;
    int free_spins;
    bool bonus_triggered;
};

struct BlackjackResult {
    std::vector<Card> player_hand;
    std::vector<Card> dealer_hand;
    int player_score;
    int dealer_score;
    std::string result; // "win", "lose", "push", "blackjack", "bust"
    double payout;
};

struct RouletteResult {
    int winning_number;
    std::string color; // "red", "black", "green"
    bool is_even;
    std::string zone; // "1-18", "19-36", "dozen", "column"
    double payout;
};

struct BaccaratResult {
    std::vector<Card> player_hand;
    std::vector<Card> banker_hand;
    int player_score;
    int banker_score;
    std::string result; // "player", "banker", "tie"
    double payout;
};

struct PokerResult {
    std::vector<Card> hand;
    std::string hand_type;
    double payout;
    std::vector<bool> held_cards;
};

struct CrashResult {
    double crash_point;
    double cashout_multiplier;
    bool crashed;
    bool won;
};

struct DiceResult {
    int roll;
    double multiplier;
    bool won;
};

struct MinesResult {
    std::vector<int> revealed_positions;
    std::vector<int> mine_positions;
    int safe_cells_remaining;
    double current_multiplier;
    bool hit_mine;
    bool won;
};

struct PlinkoResult {
    int ball_position;
    double multiplier;
    std::vector<int> path;
};

// ============ Game State ============

struct GameState {
    GameID game_id;
    UserID user_id;
    BetAmount bet_amount;
    Timestamp started_at;
    std::string game_type;
};

// ============ RNG Engine (Hardware Accelerated) ============

class RNGEngine {
public:
    static RNGEngine& get_instance();
    
    // Generate random 64-bit integer
    uint64_t next_uint64();
    
    // Generate random in range [min, max]
    uint64_t next_range(uint64_t min, uint64_t max);
    
    // Generate random float [0.0, 1.0]
    double next_double();
    
    // Generate random bool with probability
    bool next_bool(double probability = 0.5);
    
    // Shuffle vector in-place
    template<typename T>
    void shuffle(std::vector<T>& vec) {
        for (size_t i = vec.size() - 1; i > 0; --i) {
            std::swap(vec[i], vec[next_uint64() % (i + 1)]);
        }
    }
    
private:
    RNGEngine() = default;
    // Hardware RNG would be initialized here
};

// ============ Slot Machine Engine ============

class SlotEngine {
public:
    struct SlotConfig {
        int reels;
        int rows;
        int paylines;
        double rtp;
        std::string volatility;
        double min_bet;
        double max_bet;
        double max_win_multiplier;
    };
    
    struct Symbol {
        std::string name;
        std::string display;
        uint8_t weight;
        bool is_wild;
        bool is_scatter;
    };
    
    SlotEngine(const SlotConfig& config, const std::vector<Symbol>& symbols);
    
    // Execute a spin - returns result in < 10 microseconds
    SpinResult spin(BetAmount bet, UserID user_id);
    
    // Check for free spins
    int check_free_spins(const std::vector<std::vector<std::string>>& reels);
    
    // Calculate payline wins
    double calculate_payline_wins(
        const std::vector<std::vector<std::string>>& reels,
        const std::vector<std::tuple<std::vector<uint8_t>, std::string, uint8_t, double>>& paylines
    );
    
    SlotConfig config() const { return config_; }
    
private:
    SlotConfig config_;
    std::vector<Symbol> symbols_;
    std::vector<std::string> reel_strips_[10]; // Support up to 10 reels
    
    std::vector<std::string> generate_reel(int reel_index);
    std::vector<std::vector<std::string>> generate_all_reels();
};

// ============ Table Games Engine ============

class TableGamesEngine {
public:
    // ============ Blackjack ============
    BlackjackResult play_blackjack(
        BetAmount bet,
        UserID user_id,
        const std::string& action // "hit", "stand", "double", "split"
    );
    
    std::vector<Card> create_deck(int num_decks = 6);
    int calculate_score(const std::vector<Card>& hand);
    bool is_blackjack(const std::vector<Card>& hand);
    
    // ============ Roulette ============
    RouletteResult spin_roulette(
        BetAmount bet,
        const std::vector<std::tuple<std::string, std::vector<int>, double>>& bets,
        bool american = false
    );
    
    // ============ Baccarat ============
    BaccaratResult play_baccarat(
        BetAmount bet,
        const std::string& bet_on, // "player", "banker", "tie"
        UserID user_id
    );
    
    // ============ Video Poker ============
    PokerResult play_video_poker(
        BetAmount bet,
        const std::vector<bool>& hold,
        const std::string& variant // "jacks_or_better", "deuces_wild", "joker_poker"
    );
    
    std::string evaluate_poker_hand(const std::vector<Card>& hand);
    double get_payout(const std::string& hand_type, double bet);
    
private:
    std::vector<Card> deck_;
    std::vector<std::vector<Card>> player_hands_;
    std::vector<Card> dealer_hand_;
    
    void shuffle_deck();
    Card draw_card();
};

// ============ Crash/Provably Fair Games ============

class CrashGameEngine {
public:
    struct CrashConfig {
        double min_bet;
        double max_bet;
        double auto_cashout_default;
        bool provably_fair_enabled;
    };
    
    CrashGameEngine(const CrashConfig& config);
    
    // Start a new crash round
    CrashResult start_round(UserID user_id, BetAmount bet);
    
    // User cashes out
    CrashResult cashout(UserID user_id, double target_multiplier);
    
    // Force crash at point (for provably fair)
    void force_crash(double point);
    
    // Get current multiplier (called every ~100ms)
    double get_current_multiplier() const;
    
    // Check if crashed
    bool is_crashed() const;
    
    CrashConfig config() const { return config_; }
    
private:
    CrashConfig config_;
    double current_multiplier_{1.0};
    double crash_point_;
    std::atomic<bool> crashed_{false};
    std::chrono::steady_clock::time_point round_start_;
    
    double calculate_crash_point();
};

// ============ Mines Game ============

class MinesGameEngine {
public:
    MinesGameEngine(int mines_count = 3, int grid_size = 5);
    
    // Start new game
    MinesResult start_game(BetAmount bet, UserID user_id);
    
    // Reveal cell
    MinesResult reveal_cell(int position);
    
    // Cashout
    MinesResult cashout();
    
    int mines_count() const { return mines_count_; }
    int grid_size() const { return grid_size_; }
    
private:
    int mines_count_;
    int grid_size_;
    BetAmount current_bet_;
    std::vector<bool> mines_;
    std::vector<bool> revealed_;
    int safe_cells_revealed_;
    double current_multiplier_;
    bool game_over_;
};

// ============ Plinko Game ============

class PlinkoEngine {
public:
    PlinkoEngine(int rows = 8, const std::string& risk = "medium");
    
    // Drop ball
    PlinkoResult drop_ball(BetAmount bet, int pins = -1);
    
    int rows() const { return rows_; }
    std::string risk() const { return risk_; }
    
private:
    int rows_;
    std::string risk_;
    std::vector<double> multipliers_;
    
    void initialize_multipliers();
    int simulate_ball_path();
};

// ============ Dice Game ============

class DiceEngine {
public:
    DiceEngine();
    
    // Roll dice with target
    DiceResult roll(BetAmount bet, double target, bool roll_over, UserID user_id);
    
    // Quick roll (random)
    DiceResult quick_roll(UserID user_id);
    
private:
    double calculate_payout(double target, bool roll_over, int roll);
};

// ============ Game Server (WebSocket Interface) ============

class GameServer {
public:
    static GameServer& get_instance();
    
    // Register a game engine
    void register_engine(const std::string& game_type, std::shared_ptr<void> engine);
    
    // Process game request
    template<typename T>
    std::optional<T> process_request(
        const std::string& game_type,
        const std::string& action,
        const std::string& params
    );
    
    // Broadcast game state to all connected clients
    void broadcast_state(GameID game_id, const std::string& state);
    
    // Get active games count
    size_t active_games() const;
    
private:
    GameServer() = default;
    
    struct GameSession {
        GameID id;
        UserID user_id;
        std::string game_type;
        Timestamp started;
        std::atomic<bool> active;
    };
    
    std::vector<GameSession> active_sessions_;
    std::unordered_map<std::string, std::shared_ptr<void>> engines_;
};

// ============ Performance Metrics ============

struct GameMetrics {
    uint64_t total_spins;
    uint64_t total_wins;
    double average_latency_us;
    double max_latency_us;
    uint64_t active_players;
    
    double win_rate() const {
        return total_spins > 0 ? static_cast<double>(total_wins) / total_spins : 0.0;
    }
};

class MetricsCollector {
public:
    static MetricsCollector& get_instance();
    
    void record_spin(GameID game_id, uint64_t latency_us, bool won);
    void record_bet(GameID game_id, BetAmount amount);
    void record_payout(GameID game_id, BetAmount amount);
    
    GameMetrics get_metrics() const;
    void reset();
    
private:
    std::atomic<uint64_t> total_spins_{0};
    std::atomic<uint64_t> total_wins_{0};
    std::atomic<uint64_t> total_bets_{0};
    std::atomic<uint64_t> total_payouts_{0};
    std::atomic<uint64_t> total_latency_{0};
    std::atomic<uint64_t> max_latency_{0};
};

// ============ Inline Implementations ============

inline RNGEngine& RNGEngine::get_instance() {
    static RNGEngine instance;
    return instance;
}

inline uint64_t RNGEngine::next_uint64() {
    // Hardware RNG fallback would go here
    // For now using splitmix64
    static uint64_t x = 0x123456789ABCDEF0ULL;
    x += 0x9E3779B97F4A7C15ULL;
    uint64_t z = x;
    z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9ULL;
    z = (z ^ (z >> 27)) * 0x94D049BB133111EBULL;
    return z ^ (z >> 31);
}

inline uint64_t RNGEngine::next_range(uint64_t min, uint64_t max) {
    return min + (next_uint64() % (max - min + 1));
}

inline double RNGEngine::next_double() {
    return static_cast<double>(next_uint64()) / static_cast<double>(UINT64_MAX);
}

inline bool RNGEngine::next_bool(double probability) {
    return next_double() < probability;
}

} // namespace engine
} // namespace tiger

#endif // TIGER_GAME_ENGINE_HPP
