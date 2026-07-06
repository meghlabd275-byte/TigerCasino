// Additional Tiger Original Games
#include "TigerOriginals2.hpp"
#include <iostream>
#include <random>
#include <sstream>
#include <algorithm>

namespace TigerCasino {

// TigerPoker Implementation
TigerPokerGame::TigerPokerGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<std::string> TigerPokerGame::dealCards(int num_cards) {
    std::vector<std::string> cards;
    std::string suits[] = {"♠️", "♥️", "♦️", "♣️"};
    std::string ranks[] = {"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"};
    
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    for (int i = 0; i < num_cards; i++) {
        int suit_idx = h % 4;
        h = provably_fair_->generateHash(std::to_string(h) + "_suit");
        int rank_idx = h % 13;
        h = provably_fair_->generateHash(std::to_string(h) + "_rank");
        
        cards.push_back(ranks[rank_idx] + suits[suit_idx]);
    }
    
    return cards;
}

std::string TigerPokerGame::evaluateHand(const std::vector<std::string>& hand) {
    // Simplified hand evaluation
    if (isRoyalFlush(hand)) return "Royal Flush";
    if (isStraightFlush(hand)) return "Straight Flush";
    if (isFourOfAKind(hand)) return "Four of a Kind";
    if (isFullHouse(hand)) return "Full House";
    if (isFlush(hand)) return "Flush";
    if (isStraight(hand)) return "Straight";
    if (isThreeOfAKind(hand)) return "Three of a Kind";
    if (isTwoPair(hand)) return "Two Pair";
    if (isPair(hand)) return "Pair";
    return "High Card";
}

bool TigerPokerGame::isRoyalFlush(const std::vector<std::string>& hand) {
    // Simplified - real implementation would check properly
    return isFlush(hand) && isStraight(hand);
}

// TigerBingo Implementation
TigerBingoGame::TigerBingoGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<std::vector<int>> TigerBingoGame::generateCard() {
    std::vector<std::vector<int>> card(5, std::vector<int>(5));
    
    // Each column has numbers in specific ranges
    int ranges[][2] = {{1,15}, {16,30}, {31,45}, {46,60}, {61,75}};
    
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    for (int col = 0; col < 5; col++) {
        std::set<int> numbers;
        while ((int)numbers.size() < 5) {
            int num = (h % (ranges[col][1] - ranges[col][0] + 1)) + ranges[col][0];
            numbers.insert(num);
            h = provably_fair_->generateHash(std::to_string(h) + "_num");
        }
        
        int row = 0;
        for (int num : numbers) {
            card[row][col] = num;
            row++;
        }
    }
    
    // Free space in center
    card[2][2] = 0;
    
    return card;
}

std::vector<int> TigerBingoGame::drawBall() {
    std::vector<int> drawn;
    
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    std::set<int> balls;
    while ((int)balls.size() < 75) {
        int ball = (h % 75) + 1;
        balls.insert(ball);
        h = provably_fair_->generateHash(std::to_string(h) + "_ball");
        drawn.push_back(ball);
    }
    
    return drawn;
}

bool TigerBingoGame::checkWin(const std::vector<std::vector<int>>& card, const std::vector<int>& drawn) {
    std::set<int> drawn_set(drawn.begin(), drawn.end());
    
    // Check rows
    for (int row = 0; row < 5; row++) {
        bool win = true;
        for (int col = 0; col < 5; col++) {
            if (card[row][col] != 0 && drawn_set.count(card[row][col]) == 0) {
                win = false;
                break;
            }
        }
        if (win) return true;
    }
    
    // Check columns
    for (int col = 0; col < 5; col++) {
        bool win = true;
        for (int row = 0; row < 5; row++) {
            if (card[row][col] != 0 && drawn_set.count(card[row][col]) == 0) {
                win = false;
                break;
            }
        }
        if (win) return true;
    }
    
    // Check diagonals
    bool win = true;
    for (int i = 0; i < 5; i++) {
        if (card[i][i] != 0 && drawn_set.count(card[i][i]) == 0) {
            win = false;
            break;
        }
    }
    if (win) return true;
    
    win = true;
    for (int i = 0; i < 5; i++) {
        if (card[i][4-i] != 0 && drawn_set.count(card[i][4-i]) == 0) {
            win = false;
            break;
        }
    }
    
    return win;
}

// TigerLottery Implementation
TigerLotteryGame::TigerLotteryGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<int> TigerLotteryGame::drawNumbers(int num_numbers, int max_number) {
    std::vector<int> numbers;
    std::set<int> drawn;
    
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    while ((int)drawn.size() < num_numbers) {
        int num = (h % max_number) + 1;
        if (drawn.count(num) == 0) {
            drawn.insert(num);
            numbers.push_back(num);
        }
        h = provably_fair_->generateHash(std::to_string(h) + "_num");
    }
    
    std::sort(numbers.begin(), numbers.end());
    return numbers;
}

double TigerLotteryGame::calculatePrize(const std::vector<int>& user_numbers, const std::vector<int>& winning_numbers, double bet) {
    int matches = 0;
    std::set<int> winning_set(winning_numbers.begin(), winning_numbers.end());
    
    for (int num : user_numbers) {
        if (winning_set.count(num)) {
            matches++;
        }
    }
    
    // Payout table
    double multipliers[] = {0, 2, 10, 100, 1000, 10000};
    int idx = std::min(matches, 5);
    
    return bet * multipliers[idx];
}

// TigerKeno Implementation
TigerKenoGame::TigerKenoGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<int> TigerKenoGame::selectNumbers(int num_select) {
    std::vector<int> selected;
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    std::set<int> numbers;
    while ((int)numbers.size() < num_select) {
        int num = (h % 80) + 1;
        numbers.insert(num);
        h = provably_fair_->generateHash(std::to_string(h) + "_num");
    }
    
    for (int num : numbers) {
        selected.push_back(num);
    }
    
    std::sort(selected.begin(), selected.end());
    return selected;
}

std::vector<int> TigerKenoGame::drawNumbers(int num_draw) {
    std::vector<int> drawn;
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    std::set<int> numbers;
    while ((int)numbers.size() < num_draw) {
        int num = (h % 80) + 1;
        numbers.insert(num);
        h = provably_fair_->generateHash(std::to_string(h) + "_num");
    }
    
    for (int num : numbers) {
        drawn.push_back(num);
    }
    
    return drawn;
}

double TigerKenoGame::calculatePrizes(const std::vector<int>& selected, const std::vector<int>& drawn, double bet) {
    int matches = 0;
    std::set<int> drawn_set(drawn.begin(), drawn.end());
    
    for (int num : selected) {
        if (drawn_set.count(num)) {
            matches++;
        }
    }
    
    // Simplified payout
    int num_selected = selected.size();
    double payout = 0;
    
    // Payout table varies by number of spots selected
    // This is simplified
    if (num_selected == 10) {
        double payouts[] = {0, 0, 0, 0, 1, 5, 50, 500, 5000, 50000, 100000};
        payout = bet * payouts[std::min(matches, 10)];
    } else if (num_selected == 5) {
        double payouts[] = {0, 0, 1, 5, 50, 400};
        payout = bet * payouts[std::min(matches, 5)];
    } else {
        // Default: catch all = win
        if (matches > 0) {
            payout = bet * matches;
        }
    }
    
    return payout;
}

// TigerBlackjack Implementation
TigerBlackjackGame::TigerBlackjackGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

int TigerBlackjackGame::getCardValue(const std::string& card) {
    // Extract rank from card string
    if (card.find("J") != std::string::npos || 
        card.find("Q") != std::string::npos || 
        card.find("K") != std::string::npos) {
        return 10;
    }
    if (card.find("A") != std::string::npos) {
        return 11; // Could be 1 or 11
    }
    return std::stoi(card);
}

int TigerBlackjackGame::calculateScore(const std::vector<std::string>& hand) {
    int score = 0;
    int aces = 0;
    
    for (const auto& card : hand) {
        if (card.find("A") != std::string::npos) {
            aces++;
            score += 11;
        } else {
            score += getCardValue(card);
        }
    }
    
    // Convert aces from 11 to 1 if over 21
    while (score > 21 && aces > 0) {
        score -= 10;
        aces--;
    }
    
    return score;
}

bool TigerBlackjackGame::isBlackjack(const std::vector<std::string>& hand) {
    return hand.size() == 2 && calculateScore(hand) == 21;
}

// TigerRoulette Implementation
TigerRouletteGame::TigerRouletteGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

int TigerRouletteGame::spinWheel() {
    std::string seed = provably_fair_->getServerSeed() + provactly_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    
    return hash % 37; // 0-36
}

std::string TigerRouletteGame::getColor(int number) {
    if (number == 0) return "green";
    
    std::set<int> red = {1,3,5,7,9,12,14,16,18,19,21,23,25,27,30,32,34,36};
    if (red.count(number)) return "red";
    return "black";
}

double TigerRouletteGame::calculatePayout(int bet_type, int number, double bet) {
    double multiplier = 0;
    
    switch (bet_type) {
        case 0: // Straight up
            multiplier = 35;
            break;
        case 1: // Split
            multiplier = 17;
            break;
        case 2: // Street
            multiplier = 11;
            break;
        case 3: // Corner
            multiplier = 8;
            break;
        case 4: // Line
            multiplier = 5;
            break;
        case 5: // Column/Dozen
            multiplier = 2;
            break;
        case 6: // Even/Odd
        case 7: // Red/Black
        case 8: // High/Low
            multiplier = 1;
            break;
    }
    
    return bet * multiplier;
}

} // namespace TigerCasino
