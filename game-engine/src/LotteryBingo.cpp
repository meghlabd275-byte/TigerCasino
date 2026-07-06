#include "LotteryBingo.hpp"
#include <algorithm>
#include <sstream>

namespace TigerCasino {

LotteryGame::LotteryGame()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>()) {
    initializeLotteries();
}

void LotteryGame::initializeLotteries() {
    // Powerball style
    lotteries_["powerball"] = {
        id = "powerball",
        name = "Tiger Powerball",
        main_numbers = 69,
        powerball_numbers = 26,
        ticket_price = 2.0,
        min_bet = 2.0,
        max_bet = 100.0,
        rtp = 0.70
    };
    
    // Mega Millions style
    lotteries_["megamillions"] = {
        id = "megamillions",
        name = "Tiger Mega Millions",
        main_numbers = 70,
        powerball_numbers = 25,
        ticket_price = 2.0,
        min_bet = 2.0,
        max_bet = 100.0,
        rtp = 0.72
    };
    
    // Daily lottery
    lotteries_["daily"] = {
        id = "daily",
        name = "Daily Draw",
        main_numbers = 50,
        powerball_numbers = 10,
        ticket_price = 1.0,
        min_bet = 1.0,
        max_bet = 50.0,
        rtp = 0.80
    };
    
    // Instant win
    lotteries_["instant"] = {
        id = "instant",
        name = "Instant Win",
        main_numbers = 30,
        powerball_numbers = 1,
        ticket_price = 0.5,
        min_bet = 0.5,
        max_bet = 20.0,
        rtp = 0.85
    };
}

LotteryGame::Ticket LotteryGame::buyTicket(
    const std::string& lottery_id,
    const std::vector<int>& main_numbers,
    int powerball,
    double bet_amount) {
    
    Ticket ticket;
    ticket.id = generateTicketId();
    ticket.lottery_id = lottery_id;
    ticket.main_numbers = main_numbers;
    ticket.powerball = powerball;
    ticket.bet_amount = bet_amount;
    ticket.purchase_time = std::chrono::steady_clock::now();
    ticket.status = TicketStatus::PENDING;
    
    return ticket;
}

LotteryGame::DrawResult LotteryGame::drawLottery(
    const std::string& lottery_id,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    DrawResult result;
    result.lottery_id = lottery_id;
    result.draw_time = std::chrono::steady_clock::now();
    result.server_seed = server_seed;
    result.client_seed = client_seed;
    
    auto it = lotteries_.find(lottery_id);
    if (it == lotteries_.end()) {
        result.is_valid = false;
        result.error_message = "Lottery not found";
        return result;
    }
    
    const auto& lottery = it->second;
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, 
        lottery.main_numbers + lottery.powerball_numbers);
    
    // Generate main numbers (unique)
    std::set<int> selected_main;
    int idx = 0;
    while (selected_main.size() < 5) {
        int num = static_cast<int>(randoms[idx++] * lottery.main_numbers) + 1;
        if (selected_main.find(num) == selected_main.end()) {
            selected_main.insert(num);
        }
    }
    
    for (int num : selected_main) {
        result.winning_numbers.push_back(num);
    }
    
    // Generate powerball
    result.powerball = static_cast<int>(randoms[lottery.main_numbers] * lottery.powerball_numbers) + 1;
    result.is_valid = true;
    
    return result;
}

LotteryGame::PayoutResult LotteryGame::calculatePayout(
    const Ticket& ticket,
    const DrawResult& draw) {
    
    PayoutResult result;
    result.ticket_id = ticket.id;
    
    if (!draw.is_valid || ticket.status == TicketStatus::CANCELLED) {
        result.payout = 0;
        result.winning_matches = 0;
        result.powerball_match = false;
        return result;
    }
    
    // Count matching main numbers
    std::set<int> ticket_main(ticket.main_numbers.begin(), ticket.main_numbers.end());
    std::set<int> draw_main(draw.winning_numbers.begin(), draw.winning_numbers.end());
    
    std::vector<int> intersection;
    std::set_intersection(ticket_main.begin(), ticket_main.end(),
                         draw_main.begin(), draw_main.end(),
                         std::back_inserter(intersection));
    
    result.winning_matches = intersection.size();
    result.powerball_match = (ticket.powerball == draw.powerball);
    
    // Calculate payout based on matches
    double multiplier = 0;
    
    if (result.powerball_match) {
        switch (result.winning_matches) {
            case 0: multiplier = 4; break;
            case 1: multiplier = 4; break;
            case 2: multiplier = 7; break;
            case 3: multiplier = 100; break;
            case 4: multiplier = 10000; break;
            case 5: multiplier = 100000000; break; // Jackpot!
        }
    } else {
        switch (result.winning_matches) {
            case 3: multiplier = 7; break;
            case 4: multiplier = 500; break;
            case 5: multiplier = 1000000; break;
            default: multiplier = 0;
        }
    }
    
    result.payout = ticket.bet_amount * multiplier;
    
    if (result.payout > 0) {
        result.tier = getPayoutTier(result.winning_matches, result.powerball_match);
    }
    
    return result;
}

std::string LotteryGame::getPayoutTier(int matches, bool powerball_match) const {
    if (powerball_match) {
        switch (matches) {
            case 0:
            case 1: return "Match 0+ Powerball";
            case 2: return "Match 2 + Powerball";
            case 3: return "Match 3 + Powerball";
            case 4: return "Match 4 + Powerball";
            case 5: return "JACKPOT";
        }
    } else {
        switch (matches) {
            case 3: return "Match 3";
            case 4: return "Match 4";
            case 5: return "Match 5";
        }
    }
    return "No Win";
}

std::vector<double> LotteryGame::generateDeterministicRandoms(
    const std::string& server_seed,
    const std::string& client_seed,
    size_t count) {
    
    std::vector<double> result;
    std::string combined = server_seed + client_seed;
    
    for (size_t i = 0; i < count; i++) {
        std::string hash_input = combined + "LOTTERY" + std::to_string(i);
        uint64_t hash = provably_fair_->generateHash(hash_input);
        double normalized = static_cast<double>(hash) / static_cast<double>(UINT64_MAX);
        result.push_back(normalized);
    }
    
    return result;
}

std::string LotteryGame::generateTicketId() {
    auto now = std::chrono::steady_clock::now().time_since_epoch().count();
    std::ostringstream oss;
    oss << "TKT" << now << (rng_->next() % 10000);
    return oss.str();
}

std::vector<LotteryGame::LotteryInfo> LotteryGame::getAvailableLotteries() const {
    std::vector<LotteryInfo> result;
    for (const auto& pair : lotteries_) {
        result.push_back(pair.second);
    }
    return result;
}

// Bingo Implementation

BingoGame::BingoGame()
    : provably_fair_(std::make_unique<ProvablyFair>())
    , rng_(std::make_unique<RandomNumberGenerator>())
    , current_game_id_(0) {
    initializeRooms();
}

void BingoGame::initializeRooms() {
    rooms_["bingo75"] = {
        id = "bingo75",
        name = "75-Ball Bingo",
        ball_count = 75,
        card_size = 24,
        min_players = 1,
        max_players = 100,
        ticket_price = 1.0,
        min_bet = 1.0,
        max_bet = 10.0,
        rtp = 0.80
    };
    
    rooms_["bingo90"] = {
        id = "bingo90",
        name = "90-Ball Bingo",
        ball_count = 90,
        card_size = 15,
        min_players = 1,
        max_players = 100,
        ticket_price = 1.0,
        min_bet = 1.0,
        max_bet = 10.0,
        rtp = 0.82
    };
    
    rooms_["speed_bingo"] = {
        id = "speed_bingo",
        name = "Speed Bingo",
        ball_count = 30,
        card_size = 9,
        min_players = 1,
        max_players = 50,
        ticket_price = 0.5,
        min_bet = 0.5,
        max_bet = 5.0,
        rtp = 0.85
    };
}

BingoGame::Card BingoGame::generateCard(
    const std::string& room_id,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    Card card;
    card.id = generateCardId();
    card.room_id = room_id;
    
    auto it = rooms_.find(room_id);
    if (it == rooms_.end()) {
        return card;
    }
    
    const auto& room = it->second;
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, room.card_size + 10);
    
    if (room.id == "bingo75") {
        // 5x5 grid with free space in center
        std::vector<std::vector<int>> grid(5, std::vector<int>(5, 0));
        std::set<int> used;
        
        for (int col = 0; col < 5; col++) {
            int min_val = col * 15 + 1;
            int max_val = (col + 1) * 15;
            
            for (int row = 0; row < 5; row++) {
                if (row == 2 && col == 2) {
                    grid[row][col] = 0; // Free space
                    continue;
                }
                
                int num;
                do {
                    num = min_val + static_cast<int>(randoms[row * 5 + col] * (max_val - min_val + 1));
                } while (used.find(num) != used.end());
                
                used.insert(num);
                grid[row][col] = num;
            }
        }
        
        card.numbers = grid;
    }
    
    card.is_valid = true;
    return card;
}

std::vector<int> BingoGame::callNumbers(
    const std::string& room_id,
    const std::string& server_seed,
    const std::string& client_seed) {
    
    std::vector<int> called;
    
    auto it = rooms_.find(room_id);
    if (it == rooms_.end()) {
        return called;
    }
    
    int ball_count = it->second.ball_count;
    std::vector<double> randoms = generateDeterministicRandoms(server_seed, client_seed, ball_count);
    
    std::set<int> used;
    for (int i = 0; i < ball_count; i++) {
        int num;
        do {
            num = static_cast<int>(randoms[i] * ball_count) + 1;
        } while (used.find(num) != used.end());
        
        used.insert(num);
        called.push_back(num);
    }
    
    return called;
}

BingoGame::Pattern BingoGame::getPattern(const std::string& pattern_name) {
    // Predefined winning patterns
    if (pattern_name == "full_house") {
        return Pattern{"Full House", "cover_all"};
    } else if (pattern_name == "four_corners") {
        return Pattern{"Four Corners", "corners"};
    } else if (pattern_name == "diagonal") {
        return Pattern{"Diagonal", "diagonal"};
    } else if (pattern_name == "lines") {
        return Pattern{"Any Line", "line"};
    } else if (pattern_name == "x_pattern") {
        return Pattern{"X Pattern", "x"};
    }
    
    return Pattern{"Custom", "custom"};
}

bool BingoGame::checkWin(const Card& card, const std::vector<int>& called, const Pattern& pattern) {
    std::set<int> called_set(called.begin(), called.end());
    
    if (pattern.pattern_type == "cover_all") {
        // Check all numbers
        for (const auto& row : card.numbers) {
            for (int num : row) {
                if (num != 0 && called_set.find(num) == called_set.end()) {
                    return false;
                }
            }
        }
        return true;
    } else if (pattern.pattern_type == "corners") {
        // Check four corners
        if (card.numbers[0][0] != 0 && called_set.find(card.numbers[0][0]) != called_set.end() &&
            card.numbers[0][4] != 0 && called_set.find(card.numbers[0][4]) != called_set.end() &&
            card.numbers[4][0] != 0 && called_set.find(card.numbers[4][0]) != called_set.end() &&
            card.numbers[4][4] != 0 && called_set.find(card.numbers[4][4]) != called_set.end()) {
            return true;
        }
    } else if (pattern.pattern_type == "line") {
        // Check any row, column, or diagonal
        // Check rows
        for (int row = 0; row < 5; row++) {
            bool win = true;
            for (int col = 0; col < 5; col++) {
                if (card.numbers[row][col] != 0 && 
                    called_set.find(card.numbers[row][col]) == called_set.end()) {
                    win = false;
                    break;
                }
            }
            if (win) return true;
        }
    }
    
    return false;
}

std::string BingoGame::generateCardId() {
    auto now = std::chrono::steady_clock::now().time_since_epoch().count();
    std::ostringstream oss;
    oss << "CRD" << now << (rng_->next() % 10000);
    return oss.str();
}

std::vector<double> BingoGame::generateDeterministicRandoms(
    const std::string& server_seed,
    const std::string& client_seed,
    size_t count) {
    
    std::vector<double> result;
    std::string combined = server_seed + client_seed;
    
    for (size_t i = 0; i < count; i++) {
        std::string hash_input = combined + "BINGO" + std::to_string(i);
        uint64_t hash = provably_fair_->generateHash(hash_input);
        double normalized = static_cast<double>(hash) / static_cast<double>(UINT64_MAX);
        result.push_back(normalized);
    }
    
    return result;
}

std::vector<BingoGame::RoomInfo> BingoGame::getAvailableRooms() const {
    std::vector<RoomInfo> result;
    for (const auto& pair : rooms_) {
        result.push_back(pair.second);
    }
    return result;
}

} // namespace TigerCasino
