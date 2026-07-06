#include "VirtualSports.hpp"
#include <iostream>
#include <algorithm>
#include <random>

namespace TigerCasino {

VirtualSportsEngine::VirtualSportsEngine()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>()) {
    initializeLeagues();
}

void VirtualSportsEngine::initializeLeagues() {
    // Football leagues
    leagues_["football_premier"] = {
        id = "football_premier",
        name = "Virtual Premier League",
        sport = "football",
        teams = {
            {"TIGER", "Tigers FC"}, {"LION", "Lions FC"}, {"EAGLE", "Eagles FC"},
            {"WOLF", "Wolves FC"}, {"BEAR", "Bears FC"}, {"HAWK", "Hawks FC"},
            {"SHARK", "Sharks FC"}, {"PANTHER", "Panthers FC"}
        }
    };
    
    leagues_["football_la_liga"] = {
        id = "football_la_liga",
        name = "Virtual La Liga",
        sport = "football",
        teams = {
            {"DRAGON", "Dragons FC"}, {"TORO", "Toros FC"}, {"REY", "Rey FC"},
            {"LEON", "Leones FC"}, {"GALGO", "Galgos FC"}
        }
    };
    
    // Basketball leagues
    leagues_["basketball_nba"] = {
        id = "basketball_nba",
        name = "Virtual NBA",
        sport = "basketball",
        teams = {
            {"THUNDER", "Thunder"}, {"STORM", "Storm"}, {"BLAZE", "Blaze"},
            {"ICE", "Ice"}, {"STARS", "Stars"}, {"COMETS", "Comets"}
        }
    };
    
    // Tennis
    leagues_["tennis_atp"] = {
        id = "tennis_atp",
        name = "Virtual ATP Tour",
        sport = "tennis",
        teams = {}
    };
    
    // Horse racing
    leagues_["horse_racing"] = {
        id = "horse_racing",
        name = "Virtual Horse Racing",
        sport = "horse_racing",
        teams = {}
    };
}

std::vector<VirtualSportsEngine::Team> VirtualSportsEngine::getTeams(const std::string& league_id) {
    auto it = leagues_.find(league_id);
    if (it != leagues_.end()) {
        return it->second.teams;
    }
    return {};
}

VirtualSportsEngine::MatchResult VirtualSportsEngine::simulateFootballMatch(
    const std::string& league_id,
    const std::string& home_team,
    const std::string& away_team,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    MatchResult result;
    result.league_id = league_id;
    result.home_team = home_team;
    result.away_team = away_team;
    result.sport = "football";
    
    // Generate deterministic random values from seeds
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, 10);
    
    // Home team strength (0.5 - 1.5)
    double home_strength = 0.8 + randoms[0] * 0.7;
    // Away team strength (0.5 - 1.5)
    double away_strength = 0.8 + randoms[1] * 0.7;
    
    // Home advantage
    home_strength *= 1.1;
    
    // Generate goals based on strength
    double total_strength = home_strength + away_strength;
    double home_ratio = home_strength / total_strength;
    
    // Expected goals (1-4 range with some outliers)
    double expected_home_goals = 1.5 * home_ratio + randoms[2];
    double expected_away_goals = 1.5 * (1 - home_ratio) + randoms[3];
    
    // Convert to actual goals using Poisson-like distribution
    result.home_score = poissonRandom(expected_home_goals, randoms[4]);
    result.away_score = poissonRandom(expected_away_goals, randoms[5]);
    
    // Calculate winner
    if (result.home_score > result.away_score) {
        result.winner = "home";
        result.odds = 1.8 + randoms[6] * 1.5;
    } else if (result.away_score > result.home_score) {
        result.winner = "away";
        result.odds = 1.8 + randoms[7] * 1.5;
    } else {
        result.winner = "draw";
        result.odds = 2.5 + randoms[8] * 1.5;
    }
    
    result.is_over = true;
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    
    return result;
}

VirtualSportsEngine::MatchResult VirtualSportsEngine::simulateBasketballMatch(
    const std::string& league_id,
    const std::string& home_team,
    const std::string& away_team,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    MatchResult result;
    result.league_id = league_id;
    result.home_team = home_team;
    result.away_team = away_team;
    result.sport = "basketball";
    
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, 10);
    
    // Basketball scores are higher
    double home_strength = 80 + randoms[0] * 40;
    double away_strength = 80 + randoms[1] * 40;
    
    // Home advantage
    home_strength += 5;
    
    result.home_score = static_cast<int>(home_strength + randoms[2] * 20);
    result.away_score = static_cast<int>(away_strength + randoms[3] * 20);
    
    if (result.home_score > result.away_score) {
        result.winner = "home";
        result.odds = 1.7 + randoms[4] * 1.3;
    } else {
        result.winner = "away";
        result.odds = 1.7 + randoms[5] * 1.3;
    }
    
    result.is_over = true;
    return result;
}

VirtualSportsEngine::MatchResult VirtualSportsEngine::simulateTennisMatch(
    const std::string& league_id,
    const std::string& player1,
    const std::string& player2,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    MatchResult result;
    result.league_id = league_id;
    result.home_team = player1;
    result.away_team = player2;
    result.sport = "tennis";
    
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, 10);
    
    // Simulate sets (best of 3)
    int player1_sets = 0;
    int player2_sets = 0;
    
    for (int set = 0; set < 3 && player1_sets < 2 && player2_sets < 2; set++) {
        double p1_strength = 0.5 + randoms[set % 5];
        double p2_strength = 0.5 + randoms[(set + 2) % 5];
        
        if (p1_strength > p2_strength) {
            player1_sets++;
        } else {
            player2_sets++;
        }
    }
    
    result.home_score = player1_sets;
    result.away_score = player2_sets;
    
    if (player1_sets > player2_sets) {
        result.winner = "home";
        result.odds = 1.8 + randoms[6] * 1.4;
    } else {
        result.winner = "away";
        result.odds = 1.8 + randoms[7] * 1.4;
    }
    
    result.is_over = true;
    return result;
}

VirtualSportsEngine::RaceResult VirtualSportsEngine::simulateHorseRace(
    const std::string& race_id,
    const std::vector<std::string>& horses,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    RaceResult result;
    result.race_id = race_id;
    result.horses = horses;
    
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, horses.size() * 2);
    
    // Each horse has a finish time with some randomness
    std::vector<std::pair<std::string, double>> times;
    for (size_t i = 0; i < horses.size(); i++) {
        double base_time = 120.0; // 2 minutes base
        double random_factor = randoms[i * 2] * 10; // Up to 10 seconds variation
        double skill_factor = randoms[i * 2 + 1] * 15; // Up to 15 seconds based on "skill"
        times.push_back({horses[i], base_time + random_factor - skill_factor});
    }
    
    // Sort by time (lower is better)
    std::sort(times.begin(), times.end(), 
        [](const auto& a, const auto& b) { return a.second < b.second; });
    
    // Assign positions
    for (size_t i = 0; i < times.size(); i++) {
        result.positions.push_back(times[i].first);
        result.times[times[i].first] = times[i].second;
    }
    
    result.winner = times[0].first;
    result.is_over = true;
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    
    return result;
}

std::vector<double> VirtualSportsEngine::generateDeterministicRandoms(
    const std::string& server_seed,
    const std::string& client_seed,
    size_t count) {
    
    std::vector<double> result;
    std::string combined = server_seed + client_seed;
    
    for (size_t i = 0; i < count; i++) {
        std::string hash_input = combined + std::to_string(i);
        uint64_t hash = provably_fair_->generateHash(hash_input);
        double normalized = static_cast<double>(hash) / static_cast<double>(UINT64_MAX);
        result.push_back(normalized);
    }
    
    return result;
}

int VirtualSportsEngine::poissonRandom(double lambda, double random) {
    // Simple Poisson random number generation
    double L = std::exp(-lambda);
    double p = 1.0;
    int k = 0;
    
    do {
        k++;
        p *= random;
    } while (p > L);
    
    return k - 1;
}

std::vector<VirtualSportsEngine::League> VirtualSportsEngine::getAvailableLeagues() const {
    std::vector<League> result;
    for (const auto& pair : leagues_) {
        result.push_back(pair.second);
    }
    return result;
}

std::vector<VirtualSportsEngine::UpcomingMatch> VirtualSportsEngine::getUpcomingMatches(
    const std::string& league_id, int count) {
    
    std::vector<UpcomingMatch> matches;
    auto it = leagues_.find(league_id);
    if (it == leagues_.end() || it->second.teams.size() < 2) {
        return matches;
    }
    
    const auto& teams = it->second.teams;
    for (int i = 0; i < count && i * 2 + 1 < static_cast<int>(teams.size()); i++) {
        UpcomingMatch match;
        match.league_id = league_id;
        match.home_team = teams[i * 2].id;
        match.away_team = teams[i * 2 + 1].id;
        match.start_time = std::chrono::steady_clock::now() + std::chrono::seconds(60 * (i + 1));
        
        // Generate placeholder odds
        match.home_odds = 1.8 + (i % 5) * 0.1;
        match.draw_odds = 3.0 + (i % 3) * 0.2;
        match.away_odds = 1.8 + ((i + 1) % 5) * 0.1;
        
        matches.push_back(match);
    }
    
    return matches;
}

} // namespace TigerCasino
