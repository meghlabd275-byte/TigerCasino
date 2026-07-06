// Additional Table Games Implementation
#include "TableGames.hpp"
#include <iostream>
#include <random>
#include <sstream>
#include <algorithm>

namespace TigerCasino {

// Caribbean Stud Poker
CaribbeanStudGame::CaribbeanStudGame() 
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<std::string> CaribbeanStudGame::dealPlayerCards(int num_cards) {
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

std::string CaribbeanStudGame::evaluateHand(const std::vector<std::string>& hand) {
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

double CaribbeanStudGame::play(double anteBet) {
    std::vector<std::string> playerHand = dealPlayerCards(5);
    std::vector<std::string> dealerHand = dealPlayerCards(5);
    
    std::string playerRank = evaluateHand(playerHand);
    std::string dealerRank = evaluateHand(dealerHand);
    
    // Compare hands
    int playerValue = getHandValue(playerRank);
    int dealerValue = getHandValue(dealerRank);
    
    if (playerValue > dealerValue) {
        // Player wins
        return anteBet * 2; // Even money for ante
    } else if (playerValue < dealerValue) {
        return 0; // Player loses
    } else {
        // Tie - push
        return anteBet;
    }
}

// Three Card Poker
ThreeCardPokerGame::ThreeCardPokerGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double ThreeCardPokerGame::play(double bet) {
    std::vector<std::string> playerHand = dealPlayerCards(3);
    std::vector<std::string> dealerHand = dealPlayerCards(3);
    
    std::string playerRank = evaluateHand(playerHand);
    std::string dealerRank = evaluateHand(dealerHand);
    
    int playerValue = getHandValue(playerRank);
    int dealerValue = getHandValue(dealerRank);
    
    if (playerValue > dealerValue) {
        return bet * 2;
    } else if (playerValue < dealerValue) {
        return 0;
    }
    return bet; // Tie
}

// Let It Ride
LetItRideGame::LetItRideGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double LetItRideGame::play(double bet) {
    std::vector<std::string> hand = dealPlayerCards(5);
    std::string rank = evaluateHand(hand);
    int value = getHandValue(rank);
    
    // Payout table
    double multipliers[] = {0, 1, 1, 1, 2, 3, 4, 5, 10, 25, 50, 100, 200, 250, 1000};
    int idx = std::min(value, 14);
    
    return bet * multipliers[idx];
}

// Casino War
CasinoWarGame::CasinoWarGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double CasinoWarGame::play(double bet) {
    std::string playerCard = dealPlayerCards(1)[0];
    std::string dealerCard = dealPlayerCards(1)[0];
    
    int playerValue = getCardValue(playerCard);
    int dealerValue = getCardValue(dealerCard);
    
    if (playerValue > dealerValue) {
        return bet * 2;
    } else if (playerValue < dealerValue) {
        return 0;
    }
    return bet; // War - return bet
}

// Mississippi Stud
MississippiStudGame::MississippiStudGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double MississippiStudGame::play(double baseBet) {
    std::vector<std::string> playerHand = dealPlayerCards(5);
    std::string rank = evaluateHand(playerHand);
    int value = getHandValue(rank);
    
    // 3rd street, 4th street, 5th street bets
    double totalBet = baseBet * 4; // Ante + 3x call bets
    
    double multipliers[] = {0, 1, 1, 1, 2, 3, 4, 5, 10, 25, 50, 100, 200, 250, 1000};
    int idx = std::min(value, 14);
    
    return baseBet + (totalBet - baseBet) * multipliers[idx];
}

// Pai Gow Poker
PaiGowGame::PaiGowGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double PaiGowGame::play(double bet) {
    std::vector<std::string> allCards = dealPlayerCards(7);
    std::vector<std::string> highHand(allCards.begin(), allCards.begin() + 5);
    std::vector<std::string> lowHand(allCards.begin() + 5, allCards.end());
    
    std::string highRank = evaluateHand(highHand);
    std::string lowRank = evaluateHand(lowHand);
    
    int highValue = getHandValue(highRank);
    int lowValue = getHandValue(lowRank);
    
    // Both win
    if (highValue > 1 && lowValue > 1) {
        return bet * 2;
    }
    // Both lose
    else if (highValue == 1 && lowValue == 1) {
        return 0;
    }
    // Push
    return bet;
}

// Red Dog
RedDogGame::RedDogGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double RedDogGame::play(double bet) {
    std::vector<std::string> cards = dealPlayerCards(3);
    int card1 = getCardValue(cards[0]);
    int card2 = getCardValue(cards[1]);
    int card3 = getCardValue(cards[2]);
    
    // Sort first two cards
    int low = std::min(card1, card2);
    int high = std::max(card1, card2);
    
    // Check if third card is between
    if (card3 > low && card3 < high) {
        return bet * 2; // Win
    } else if (card3 == low || card3 == high) {
        return bet; // Push
    }
    return 0; // Lose
}

// Sic Bo
SicBoGame::SicBoGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

std::vector<int> SicBoGame::rollThreeDice() {
    std::string seed = provably_fair_->getServerSeed() + provably_fair_->getClientSeed();
    uint64_t hash = provably_fair_->generateHash(seed);
    uint64_t h = hash;
    
    std::vector<int> dice(3);
    dice[0] = (h % 6) + 1;
    h = provably_fair_->generateHash(std::to_string(h) + "_d1");
    dice[1] = (h % 6) + 1;
    h = provably_fair_->generateHash(std::to_string(h) + "_d2");
    dice[2] = (h % 6) + 1;
    
    return dice;
}

double SicBoGame::play(const std::string& betType, double bet, const std::vector<int>& dice) {
    int total = dice[0] + dice[1] + dice[2];
    double multiplier = 0;
    
    if (betType == "small") { // 4-10
        multiplier = (total >= 4 && total <= 10) ? 1 : 0;
    } else if (betType == "big") { // 11-17
        multiplier = (total >= 11 && total <= 17) ? 1 : 0;
    } else if (betType == "triple") {
        multiplier = (dice[0] == dice[1] && dice[1] == dice[2]) ? 30 : 0;
    } else if (betType == "any_triple") {
        multiplier = (dice[0] == dice[1] || dice[1] == dice[2] || dice[0] == dice[2]) ? 8 : 0;
    } else {
        // Number bet
        int num = std::stoi(betType);
        int count = 0;
        for (int d : dice) if (d == num) count++;
        if (count == 1) multiplier = 1;
        else if (count == 2) multiplier = 2;
        else if (count == 3) multiplier = 3;
    }
    
    return bet * (1 + multiplier);
}

// War
WarGame::WarGame()
    : provably_fair_(std::make_unique<ProvablyFair>()) {
}

double WarGame::play(double bet) {
    std::string playerCard = dealPlayerCards(1)[0];
    std::string dealerCard = dealPlayerCards(1)[0];
    
    int playerValue = getCardValue(playerCard);
    int dealerValue = getCardValue(dealerCard);
    
    if (playerValue > dealerValue) return bet * 2;
    else if (playerValue < dealerValue) return 0;
    else {
        // Go to war - simplified
        return bet;
    }
}

} // namespace TigerCasino
