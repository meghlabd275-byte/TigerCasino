#ifndef TIGER_GAME_ENGINE_H
#define TIGER_GAME_ENGINE_H

#include <string>
#include <vector>
#include <map>
#include <memory>
#include <functional>
#include <mutex>
#include <atomic>
#include <array>

namespace TigerCasino {

// Constants
constexpr int MAX_DICE_VALUE = 100;
constexpr int MIN_DICE_VALUE = 0;
constexpr int SLOT_REELS = 5;
constexpr int SLOT_SYMBOLS_PER_REEL = 3;
constexpr int PLINKO_MIN_ROWS = 8;
constexpr int PLINKO_MAX_ROWS = 16;
constexpr int MINES_MAX_GRID = 25;
constexpr int MINES_MIN_COUNT = 1;
constexpr int MINES_MAX_COUNT = 24;

// Seed size for cryptographic operations
constexpr int SEED_SIZE = 32;
constexpr int HASH_SIZE = 64;

// House edge percentages (multiplied by 10000 for precision)
constexpr int HOUSE_EDGE_DICE = 100;      // 1%
constexpr int HOUSE_EDGE_CRASH = 300;     // 3%
constexpr int HOUSE_EDGE_MINES = 200;     // 2%
constexpr int HOUSE_EDGE_PLINKO = 400     // 4%
constexpr int HOUSE_EDGE_SLOTS = 500;     // 5%
constexpr int HOUSE_EDGE_LIMBO = 200;     // 2%
constexpr int HOUSE_EDGE_HILO = 150;      // 1.5%

// Game result structures
struct DiceResult {
    double roll;
    double target;
    double multiplier;
    double win_amount;
    bool direction_over;
    bool is_win;
    std::string server_seed;
    std::string client_seed;
    int nonce;
    bool verified;
};

struct CrashResult {
    std::string round_id;
    double crash_point;
    double player_cashout;
    bool did_crash;
    bool did_cashout;
    double bet_amount;
    double multiplier;
    double win_amount;
    std::string server_seed;
    std::string client_seed;
    int nonce;
};

struct MinesResult {
    std::string game_id;
    int mines_count;
    int current_step;
    int revealed_tile;
    bool is_mine;
    bool game_over;
    double multiplier;
    double bet_amount;
    double win_amount;
    std::vector<int> revealed_tiles;
    std::vector<bool> mine_locations;
};

struct PlinkoResult {
    std::string game_id;
    int rows;
    std::string risk;  // low, medium, high
    std::vector<int> path;  // Ball path through rows
    int final_bucket;
    double multiplier;
    double bet_amount;
    double win_amount;
    std::string server_seed;
    std::string client_seed;
    int nonce;
};

struct SlotsResult {
    std::array<int, SLOT_REELS> reels;
    std::vector<std::string> symbols;
    int scatter_count;
    int wild_count;
    double multiplier;
    double bet_amount;
    double win_amount;
    bool is_jackpot;
    bool is_free_spin;
    std::string win_line;
    std::string server_seed;
    std::string client_seed;
    int nonce;
};

struct LimboResult {
    double target_multiplier;
    double result_multiplier;
    double bet_amount;
    double win_amount;
    bool is_win;
    std::string server_seed;
    std::string client_seed;
    int nonce;
};

struct HiloResult {
    std::string current_card;
    std::string next_card;
    std::string choice;  // higher, lower, equal
    bool is_correct;
    int streak;
    double multiplier;
    double bet_amount;
    double win_amount;
    std::string server_seed;
    std::string client_seed;
    int nonce;
};

// Seed management
struct Seeds {
    std::string server_seed;
    std::string server_seed_hash;
    std::string client_seed;
    int nonce;
    
    Seeds() : nonce(0) {}
};

// Game state
enum class GameStatus {
    PENDING,
    RUNNING,
    WON,
    LOST,
    CRASHED,
    CANCELLED
};

// Player bet info
struct BetInfo {
    std::string bet_id;
    std::string user_id;
    std::string game_type;
    double bet_amount;
    double win_amount;
    GameStatus status;
    std::string game_data;
    uint64_t timestamp;
};

// Game engine class
class TigerGameEngine {
private:
    // Cryptographic random number generator state
    std::array<uint8_t, SEED_SIZE> entropy_pool_;
    std::mutex rng_mutex_;
    std::atomic<uint64_t> round_counter_;
    
    // Game-specific state
    std::map<std::string, std::shared_ptr<Seeds>> active_seeds_;
    std::map<std::string, std::shared_ptr<MinesResult>> mines_games_;
    std::map<std::string, CrashResult> crash_rounds_;
    std::mutex games_mutex_;
    
    // Statistics
    std::atomic<uint64_t> total_rounds_played_;
    std::atomic<uint64_t> total_winnings_;
    std::mutex stats_mutex_;
    
public:
    TigerGameEngine();
    ~TigerGameEngine();
    
    // ===== PROVABLY FAIR SYSTEM =====
    
    /**
     * Generate cryptographically secure seeds
     * @param client_seed Optional client-provided seed for additional randomness
     * @return Generated seeds
     */
    Seeds generateSeeds(const std::string& client_seed = "");
    
    /**
     * Hash a seed using SHA-256
     * @param seed The seed to hash
     * @return Hex-encoded hash
     */
    std::string hashSeed(const std::string& seed);
    
    /**
     * Verify that a server seed matches its hash
     * @param server_seed The original server seed
     * @param server_seed_hash The claimed hash
     * @return true if valid
     */
    bool verifySeed(const std::string& server_seed, const std::string& server_seed_hash);
    
    /**
     * Generate a random outcome using HMAC-based deterministic RNG
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce for this specific outcome
     * @param max Maximum value (exclusive)
     * @return Random integer in [0, max)
     */
    int generateOutcome(const std::string& server_seed, 
                       const std::string& client_seed, 
                       int nonce, 
                       int max);
    
    /**
     * Generate a float outcome in [0, 1)
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @return Random float in [0, 1)
     */
    double generateFloatOutcome(const std::string& server_seed,
                                const std::string& client_seed,
                                int nonce);
    
    /**
     * Generate cryptographic random bytes
     * @param length Number of bytes to generate
     * @return Vector of random bytes
     */
    std::vector<uint8_t> generateRandomBytes(size_t length);
    
    // ===== DICE GAME =====
    
    /**
     * Play a dice game
     * @param user_id User identifier
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Current nonce
     * @param bet_amount Bet amount
     * @param target Target number (0-100)
     * @param over Direction: true = over, false = under
     * @return Dice result
     */
    DiceResult playDice(const std::string& user_id,
                        const std::string& server_seed,
                        const std::string& client_seed,
                        int nonce,
                        double bet_amount,
                        double target,
                        bool over);
    
    // ===== CRASH GAME =====
    
    /**
     * Start a new crash round
     * @param round_id Round identifier
     * @param server_seed Server seed
     * @param server_seed_hash Hash of server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @return Crash result
     */
    CrashResult startCrashRound(const std::string& round_id,
                                const std::string& server_seed,
                                const std::string& server_seed_hash,
                                const std::string& client_seed,
                                int nonce);
    
    /**
     * Cash out from a crash round
     * @param round_id Round identifier
     * @param user_id User identifier
     * @param bet_amount Original bet amount
     * @param multiplier Multiplier at cashout
     * @return Crash result with winnings
     */
    CrashResult crashCashout(const std::string& round_id,
                             const std::string& user_id,
                             double bet_amount,
                             double multiplier);
    
    /**
     * Get current crash multiplier for a round (for animation)
     * @param round_id Round identifier
     * @param elapsed_ms Milliseconds since round start
     * @return Current multiplier
     */
    double getCrashMultiplier(const std::string& round_id, uint64_t elapsed_ms);
    
    /**
     * Crash a round (called when rocket flies away)
     * @param round_id Round identifier
     * @return Final crash result
     */
    CrashResult crashRound(const std::string& round_id);
    
    // ===== MINES GAME =====
    
    /**
     * Start a new mines game
     * @param game_id Game identifier
     * @param user_id User identifier
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @param bet_amount Bet amount
     * @param mines_count Number of mines (1-24)
     * @return Mines result
     */
    MinesResult startMinesGame(const std::string& game_id,
                              const std::string& user_id,
                              const std::string& server_seed,
                              const std::string& client_seed,
                              int nonce,
                              double bet_amount,
                              int mines_count);
    
    /**
     * Reveal a tile in mines game
     * @param game_id Game identifier
     * @param user_id User identifier
     * @param tile Tile index (0-24)
     * @return Mines result
     */
    MinesResult revealMinesTile(const std::string& game_id,
                               const std::string& user_id,
                               int tile);
    
    /**
     * Cash out from mines game
     * @param game_id Game identifier
     * @param user_id User identifier
     * @return Mines result with winnings
     */
    MinesResult cashoutMines(const std::string& game_id, const std::string& user_id);
    
    // ===== PLINKO GAME =====
    
    /**
     * Play plinko game
     * @param user_id User identifier
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @param bet_amount Bet amount
     * @param rows Number of rows (8-16)
     * @param risk Risk level (low, medium, high)
     * @return Plinko result
     */
    PlinkoResult playPlinko(const std::string& user_id,
                            const std::string& server_seed,
                            const std::string& client_seed,
                            int nonce,
                            double bet_amount,
                            int rows,
                            const std::string& risk);
    
    // ===== SLOTS GAME =====
    
    /**
     * Play slots game
     * @param user_id User identifier
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @param bet_amount Bet amount
     * @param lines Number of paylines
     * @return Slots result
     */
    SlotsResult playSlots(const std::string& user_id,
                         const std::string& server_seed,
                         const std::string& client_seed,
                         int nonce,
                         double bet_amount,
                         int lines = 20);
    
    // ===== LIMBO GAME =====
    
    /**
     * Play limbo game
     * @param user_id User identifier
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @param bet_amount Bet amount
     * @param target_multiplier Target multiplier to beat
     * @return Limbo result
     */
    LimboResult playLimbo(const std::string& user_id,
                          const std::string& server_seed,
                          const std::string& client_seed,
                          int nonce,
                          double bet_amount,
                          double target_multiplier);
    
    // ===== HI-LO GAME =====
    
    /**
     * Start a hi-lo game
     * @param game_id Game identifier
     * @param user_id User identifier
     * @param server_seed Server seed
     * @param client_seed Client seed
     * @param nonce Nonce
     * @param bet_amount Bet amount
     * @return Hilo result
     */
    HiloResult startHiloGame(const std::string& game_id,
                             const std::string& user_id,
                             const std::string& server_seed,
                             const std::string& client_seed,
                             int nonce,
                             double bet_amount);
    
    /**
     * Make a choice in hi-lo
     * @param game_id Game identifier
     * @param user_id User identifier
     * @param choice higher, lower, or equal
     * @return Hilo result
     */
    HiloResult playHiloChoice(const std::string& game_id,
                              const std::string& user_id,
                              const std::string& choice);
    
    // ===== UTILITY FUNCTIONS =====
    
    /**
     * Calculate payout multiplier based on game type and parameters
     * @param game_type Type of game
     * @param params Game-specific parameters
     * @return Multiplier (1.0 = even money)
     */
    double calculateMultiplier(const std::string& game_type,
                              const std::map<std::string, std::string>& params);
    
    /**
     * Get game statistics
     * @return Map of statistics
     */
    std::map<std::string, uint64_t> getStatistics();
    
    /**
     * Reset statistics
     */
    void resetStatistics();
    
    /**
     * Get house edge for a game type
     * @param game_type Game type
     * @return House edge as percentage * 10000
     */
    int getHouseEdge(const std::string& game_type);
    
private:
    // Internal helper methods
    std::string hmacSha256(const std::string& key, const std::string& data);
    std::string bytesToHex(const std::vector<uint8_t>& bytes);
    std::vector<uint8_t> hexToBytes(const std::string& hex);
    double calculateMinesMultiplier(int revealed, int mines, double bet);
    std::vector<double> getPlinkoPayouts(int rows, const std::string& risk);
    std::string generateSlotSymbols(const std::string& seed, int reel);
    bool checkSlotWin(const std::array<int, SLOT_REELS>& reels);
};

} // namespace TigerCasino

#endif // TIGER_GAME_ENGINE_H
