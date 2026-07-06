#pragma once

#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"
#include <string>
#include <vector>
#include <memory>
#include <array>
#include <map>
#include <chrono>
#include <random>

namespace TigerCasino {

/**
 * Virtual Sports Games - Simulated sports events with realistic outcomes
 */

/**
 * Virtual Football - Simulated soccer matches
 */
class VirtualFootballGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 10000.00;
    static constexpr double BASE_RTP = 0.95;

    struct Team {
        std::string name;
        double attack_rating;
        double defense_rating;
        double form; // 0-1 recent form
    };

    struct MatchResult {
        std::string home_team;
        std::string away_team;
        uint8_t home_score;
        uint8_t away_score;
        std::string outcome; // home_win, draw, away_win
        double home_win_odds;
        double draw_odds;
        double away_win_odds;
        std::vector<std::pair<std::string, uint8_t>> goal_scorers;
        std::string server_seed;
    };

    struct Bet {
        std::string player_id;
        std::string bet_type; // home, draw, away, over_under, exact_score
        double amount;
        double odds;
        double potential_win;
    };

private:
    static constexpr std::array<std::string, 20> TEAM_NAMES = {{
        "Tigers", "Eagles", "Lions", "Wolves", "Panthers",
        "Sharks", "Bears", "Hawks", "Falcons", "Dragons",
        "Bulls", "Stallions", "Raptors", "Cobras", "Vipers",
        "Storm", "Thunder", "Lightning", "Warriors", "Knights"
    }};

    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t match_counter_;
    std::vector<Team> teams_;

public:
    VirtualFootballGame();
    
    std::string getName() const { return "Virtual Football"; }
    std::string getType() const { return "Virtual Sports"; }
    double getRTP() const { return BASE_RTP; }

    MatchResult playMatch(const std::string& home_team, const std::string& away_team);
    MatchResult playRandomMatch();
    std::vector<std::string> getAvailableTeams() const;
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);

private:
    Team generateTeam(const std::string& name);
    std::vector<std::pair<std::string, uint8_t>> generateGoalScorers(uint8_t goals, const std::string& team);
    double calculateOdds(double probability);
};

inline VirtualFootballGame::VirtualFootballGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , match_counter_(0) {
    // Generate teams
    for (const auto& name : TEAM_NAMES) {
        teams_.push_back(generateTeam(name));
    }
}

inline VirtualFootballGame::Team VirtualFootballGame::generateTeam(const std::string& name) {
    Team team;
    team.name = name;
    team.attack_rating = rng_->generateDouble(60.0, 95.0);
    team.defense_rating = rng_->generateDouble(60.0, 95.0);
    team.form = rng_->generateDouble(0.3, 1.0);
    return team;
}

inline double VirtualFootballGame::calculateOdds(double probability) {
    if (probability <= 0) return 0;
    // Apply margin for house edge
    return (1.0 / probability) * 0.92;
}

inline std::vector<std::string> VirtualFootballGame::getAvailableTeams() const {
    std::vector<std::string> names;
    for (const auto& team : teams_) {
        names.push_back(team.name);
    }
    return names;
}

inline std::vector<std::pair<std::string, uint8_t>> 
VirtualFootballGame::generateGoalScorers(uint8_t goals, const std::string& team) {
    std::vector<std::pair<std::string, uint8_t>> scorers;
    static const std::array<std::string, 10> PLAYER_NAMES = {{
        "Smith", "Johnson", "Williams", "Brown", "Jones",
        "Garcia", "Miller", "Davis", "Rodriguez", "Martinez"
    }};
    
    std::uniform_int_distribution<size_t> dist(0, PLAYER_NAMES.size() - 1);
    
    for (uint8_t i = 0; i < goals; ++i) {
        std::string player = PLAYER_NAMES[dist(rng_->getEngine())];
        // Check if player already scored
        bool found = false;
        for (auto& scorer : scorers) {
            if (scorer.first == player) {
                scorer.second++;
                found = true;
                break;
            }
        }
        if (!found) {
            scorers.push_back({player, 1});
        }
    }
    
    return scorers;
}

inline VirtualFootballGame::MatchResult 
VirtualFootballGame::playMatch(const std::string& home_team_name, const std::string& away_team_name) {
    MatchResult result;
    match_counter_++;
    
    // Find teams
    const Team* home_team = nullptr;
    const Team* away_team = nullptr;
    
    for (const auto& team : teams_) {
        if (team.name == home_team_name) home_team = &team;
        if (team.name == away_team_name) away_team = &team;
    }
    
    if (!home_team || !away_team) {
        // Use random teams
        size_t home_idx = rng_->generateInt(0, teams_.size() - 1);
        size_t away_idx = (home_idx + rng_->generateInt(1, teams_.size() - 1)) % teams_.size();
        home_team = &teams_[home_idx];
        away_team = &teams_[away_idx];
    }
    
    result.home_team = home_team->name;
    result.away_team = away_team->name;
    result.server_seed = provably_fair_->getServerSeed();
    
    // Calculate expected goals using Poisson-like distribution
    double home_xg = (home_team->attack_rating * home_team->form) / 50.0;
    double away_xg = (away_team->attack_rating * away_team->form) / 50.0;
    
    // Adjust for defense
    home_xg *= (away_team->defense_rating / 100.0);
    away_xg *= (home_team->defense_rating / 100.0);
    
    // Generate goals using Poisson-like distribution
    result.home_score = 0;
    result.away_score = 0;
    
    // Home goals
    for (int i = 0; i < 10; ++i) { // Simulate 10 scoring chances
        if (rng_->generateDouble(0.0, 1.0) < home_xg / 10.0) {
            result.home_score++;
        }
    }
    
    // Away goals
    for (int i = 0; i < 10; ++i) {
        if (rng_->generateDouble(0.0, 1.0) < away_xg / 10.0) {
            result.away_score++;
        }
    }
    
    // Determine outcome
    if (result.home_score > result.away_score) {
        result.outcome = "home_win";
    } else if (result.home_score < result.away_score) {
        result.outcome = "away_win";
    } else {
        result.outcome = "draw";
    }
    
    // Calculate odds
    double home_prob = static_cast<double>(result.home_score + 1) / 
                        (result.home_score + result.away_score + 2);
    double away_prob = static_cast<double>(result.away_score + 1) / 
                        (result.home_score + result.away_score + 2);
    double draw_prob = 1.0 - home_prob - away_prob;
    
    // Recalculate based on team ratings
    double total_rating = home_team->attack_rating + away_team->attack_rating;
    home_prob = (home_team->attack_rating / total_rating) * 0.5 + 0.25;
    away_prob = (away_team->attack_rating / total_rating) * 0.5 + 0.25;
    draw_prob = 0.5 - (home_prob - away_prob) * 0.3;
    
    home_prob = std::max(0.1, std::min(0.8, home_prob));
    away_prob = std::max(0.1, std::min(0.8, away_prob));
    draw_prob = std::max(0.1, std::min(0.5, draw_prob));
    
    result.home_win_odds = calculateOdds(home_prob);
    result.draw_odds = calculateOdds(draw_prob);
    result.away_win_odds = calculateOdds(away_prob);
    
    // Generate goal scorers
    if (result.home_score > 0) {
        auto home_scorers = generateGoalScorers(result.home_score, result.home_team);
        result.goal_scorers.insert(result.goal_scorers.end(), home_scorers.begin(), home_scorers.end());
    }
    if (result.away_score > 0) {
        auto away_scorers = generateGoalScorers(result.away_score, result.away_team);
        result.goal_scorers.insert(result.goal_scorers.end(), away_scorers.begin(), away_scorers.end());
    }
    
    return result;
}

inline VirtualFootballGame::MatchResult VirtualFootballGame::playRandomMatch() {
    size_t home_idx = rng_->generateInt(0, teams_.size() - 1);
    size_t away_idx = (home_idx + rng_->generateInt(1, teams_.size() - 1)) % teams_.size();
    return playMatch(teams_[home_idx].name, teams_[away_idx].name);
}

inline void VirtualFootballGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void VirtualFootballGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Virtual Basketball - Simulated basketball games
 */
class VirtualBasketballGame {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 10000.00;
    static constexpr double BASE_RTP = 0.95;

    struct GameResult {
        std::string home_team;
        std::string away_team;
        uint16_t home_score;
        uint16_t away_score;
        uint8_t home_quarter_scores[4];
        uint8_t away_quarter_scores[4];
        std::string outcome;
        double home_win_odds;
        double away_win_odds;
        std::string server_seed;
    };

private:
    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t game_counter_;

public:
    VirtualBasketballGame();
    
    std::string getName() const { return "Virtual Basketball"; }
    std::string getType() const { return "Virtual Sports"; }
    double getRTP() const { return BASE_RTP; }

    GameResult playGame(const std::string& home_team, const std::string& away_team);
    GameResult playRandomGame();
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);
};

inline VirtualBasketballGame::VirtualBasketballGame() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , game_counter_(0) {}

inline VirtualBasketballGame::GameResult 
VirtualBasketballGame::playGame(const std::string& home_team, const std::string& away_team) {
    GameResult result;
    game_counter_++;
    
    result.home_team = home_team;
    result.away_team = away_team;
    result.server_seed = provably_fair_->getServerSeed();
    
    // Generate quarter scores
    result.home_score = 0;
    result.away_score = 0;
    
    for (int q = 0; q < 4; ++q) {
        result.home_quarter_scores[q] = rng_->generateInt(15, 35);
        result.away_quarter_scores[q] = rng_->generateInt(15, 35);
        
        result.home_score += result.home_quarter_scores[q];
        result.away_score += result.away_quarter_scores[q];
    }
    
    if (result.home_score > result.away_score) {
        result.outcome = "home_win";
        result.home_win_odds = 1.85 + rng_->generateDouble(-0.1, 0.3);
        result.away_win_odds = 2.0 + rng_->generateDouble(-0.2, 0.2);
    } else {
        result.outcome = "away_win";
        result.away_win_odds = 1.85 + rng_->generateDouble(-0.1, 0.3);
        result.home_win_odds = 2.0 + rng_->generateDouble(-0.2, 0.2);
    }
    
    return result;
}

inline VirtualBasketballGame::GameResult VirtualBasketballGame::playRandomGame() {
    static const std::array<std::string, 20> TEAMS = {{
        "Warriors", "Lakers", "Celtics", "Bulls", "Heat",
        "Knicks", "Nets", "76ers", "Raptors", "Bucks",
        "Suns", "Clippers", "Nuggets", "Mavericks", "Rockets",
        "Hornets", "Pacers", "Cavaliers", "Magic", "Hawks"
    }};
    
    size_t home_idx = rng_->generateInt(0, TEAMS.size() - 1);
    size_t away_idx = (home_idx + rng_->generateInt(1, TEAMS.size() - 1)) % TEAMS.size();
    
    return playGame(TEAMS[home_idx], TEAMS[away_idx]);
}

inline void VirtualBasketballGame::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void VirtualBasketballGame::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

/**
 * Virtual Horse Racing
 */
class VirtualHorseRacing {
public:
    static constexpr double MIN_BET = 0.10;
    static constexpr double MAX_BET = 10000.00;
    static constexpr double BASE_RTP = 0.92;
    static constexpr size_t NUM_HORSES = 8;

    struct Horse {
        std::string name;
        double rating;
        double form;
        uint8_t position;
        double odds;
    };

    struct RaceResult {
        uint64_t race_id;
        uint8_t distance_furlongs;
        std::vector<Horse> horses;
        std::string winner;
        double winner_odds;
        std::string server_seed;
    };

private:
    static constexpr std::array<std::string, 16> HORSE_NAMES = {{
        "Thunder Bolt", "Lightning Strike", "Star Dancer", "Golden Arrow",
        "Silver Shadow", "Crystal Wing", "Phoenix Rising", "Dragon Fire",
        "Midnight Runner", "Sunset Champion", "Ocean Wave", "Mountain Peak",
        "Desert Storm", "Forest Ghost", "Arctic Wind", "Tropical Storm"
    }};

    std::unique_ptr<ProvablyFair> provably_fair_;
    std::unique_ptr<RandomNumberGenerator> rng_;
    uint64_t race_counter_;

public:
    VirtualHorseRacing();
    
    std::string getName() const { return "Virtual Horse Racing"; }
    std::string getType() const { return "Virtual Sports"; }
    double getRTP() const { return BASE_RTP; }

    RaceResult runRace(uint8_t distance = 12);
    std::vector<Horse> getFormedHorses();
    
    void setClientSeed(const std::string& seed);
    void setServerSeed(const std::string& seed);
};

inline VirtualHorseRacing::VirtualHorseRacing() 
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , race_counter_(0) {}

inline std::vector<VirtualHorseRacing::Horse> VirtualHorseRacing::getFormedHorses() {
    std::vector<Horse> horses;
    
    for (size_t i = 0; i < NUM_HORSES; ++i) {
        Horse horse;
        horse.name = HORSE_NAMES[i];
        horse.rating = rng_->generateDouble(60.0, 100.0);
        horse.form = rng_->generateDouble(0.5, 1.0);
        horse.position = 0;
        
        // Calculate odds based on rating
        double win_prob = (horse.rating / 100.0) * horse.form;
        horse.odds = win_prob > 0 ? (1.0 / win_prob) * 0.88 : 10.0;
        
        horses.push_back(horse);
    }
    
    return horses;
}

inline VirtualHorseRacing::RaceResult VirtualHorseRacing::runRace(uint8_t distance) {
    RaceResult result;
    race_counter_++;
    result.race_id = race_counter_;
    result.distance_furlongs = distance;
    result.server_seed = provably_fair_->getServerSeed();
    
    auto horses = getFormedHorses();
    
    // Run the race - simulate each horse's performance
    std::vector<double> race_times;
    for (auto& horse : horses) {
        // Base time + random variation + rating adjustment
        double base_time = distance * 12.0; // 12 seconds per furlong base
        double variation = rng_->generateDouble(-2.0, 2.0);
        double rating_adjustment = (100.0 - horse.rating) * 0.1;
        double form_adjustment = (1.0 - horse.form) * 1.0;
        
        double race_time = base_time + variation + rating_adjustment + form_adjustment;
        race_times.push_back(race_time);
    }
    
    // Sort horses by race time (lower is better)
    std::vector<size_t> positions(NUM_HORSES);
    std::iota(positions.begin(), positions.end(), 0);
    std::sort(positions.begin(), positions.end(), 
              [&race_times](size_t a, size_t b) {
                  return race_times[a] < race_times[b];
              });
    
    // Assign positions
    for (size_t i = 0; i < NUM_HORSES; ++i) {
        horses[positions[i]].position = static_cast<uint8_t>(i + 1);
    }
    
    result.horses = horses;
    result.winner = horses[positions[0]].name;
    result.winner_odds = horses[positions[0]].odds;
    
    return result;
}

inline void VirtualHorseRacing::setClientSeed(const std::string& seed) {
    provably_fair_->setClientSeed(seed);
}

inline void VirtualHorseRacing::setServerSeed(const std::string& seed) {
    provably_fair_->setServerSeed(seed);
}

} // namespace TigerCasino
