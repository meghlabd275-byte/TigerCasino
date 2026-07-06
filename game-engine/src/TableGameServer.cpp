#include "TableGameServer.hpp"
#include <algorithm>
#include <random>

namespace TigerCasino {

// Card methods
std::string Card::toString() const {
    return rank + " of " + suit;
}

bool Card::isFaceCard() const {
    return rank == "J" || rank == "Q" || rank == "K";
}

bool Card::isAce() const {
    return rank == "A";
}

// Hand methods
int Hand::getValue() const {
    int value = 0;
    int aces = 0;
    
    for (const auto& card : cards) {
        if (card.isAce()) {
            aces++;
            value += 11;
        } else if (card.isFaceCard()) {
            value += 10;
        } else {
            value += card.value;
        }
    }
    
    // Handle aces
    while (value > 21 && aces > 0) {
        value -= 10;
        aces--;
    }
    
    return value;
}

bool Hand::isBust() const {
    return getValue() > 21;
}

bool Hand::isBlackjack() const {
    return cards.size() == 2 && getValue() == 21;
}

bool Hand::canSplit() const {
    return cards.size() == 2 && cards[0].value == cards[1].value;
}

void Hand::clear() {
    cards.clear();
    isSplit = false;
    isDoubled = false;
    isSurrendered = false;
}

// Base table game server
TableGameServer::TableGameServer(const std::string& gameId, const std::string& gameName,
                                  TableGameType type, double minBet, double maxBet, double rtp)
    : gameId_(gameId)
    , gameName_(gameName)
    , gameType_(type)
    , minBet_(minBet)
    , maxBet_(maxBet)
    , rtp_(rtp)
    , currentCardIndex_(0)
    , decksInShoe_(8) {
    
    std::random_device rd;
    rng_.seed(rd());
    createShoe(decksInShoe_);
}

void TableGameServer::createShoe(int numDecks) {
    shoe_.clear();
    
    std::vector<std::string> suits = {"Hearts", "Diamonds", "Clubs", "Spades"};
    std::vector<std::string> ranks = {"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"};
    
    for (int d = 0; d < numDecks; d++) {
        for (const auto& suit : suits) {
            for (size_t i = 0; i < ranks.size(); i++) {
                Card card;
                card.suit = suit;
                card.rank = ranks[i];
                
                if (ranks[i] == "A") {
                    card.value = 11;
                } else if (ranks[i] == "J" || ranks[i] == "Q" || ranks[i] == "K") {
                    card.value = 10;
                } else {
                    card.value = std::stoi(ranks[i]);
                }
                
                shoe_.push_back(card);
            }
        }
    }
    
    shuffleShoe();
}

Card TableGameServer::drawCard() {
    if (currentCardIndex_ >= shoe_.size() - (decksInShoe_ * 10)) {
        createShoe(decksInShoe_);
        currentCardIndex_ = 0;
    }
    return shoe_[currentCardIndex_++];
}

void TableGameServer::shuffleShoe() {
    std::lock_guard<std::mutex> lock(rngMutex_);
    std::shuffle(shoe_.begin(), shoe_.end(), rng_);
    currentCardIndex_ = 0;
}

// Blackjack server implementation
BlackjackServer::BlackjackServer()
    : TableGameServer("blackjack", "Blackjack", TableGameType::BLACKJACK, 1.0, 10000.0, 0.995)
    , surrenderAllowed_(true)
    , doubleAfterSplitAllowed_(true)
    , hitSplitAces_(false)
    , blackjackPayout_(3) {
}

std::vector<std::string> BlackjackServer::getValidActions() const {
    return {"hit", "stand", "double", "split", "surrender"};
}

TableGameResult BlackjackServer::play(const std::string& playerId, 
                                      const TableBet& bet,
                                      const std::map<std::string, std::string>& actions) {
    TableGameResult result;
    
    // Validate bet
    if (bet.amount < minBet_ || bet.amount > maxBet_) {
        result.success = false;
        result.outcome = "Invalid bet amount";
        return result;
    }
    
    // Initial deal
    Hand playerHand;
    Hand dealerHand;
    
    playerHand.cards.push_back(drawCard());
    dealerHand.cards.push_back(drawCard());
    playerHand.cards.push_back(drawCard());
    dealerHand.cards.push_back(drawCard());
    
    result.playerHands.push_back(playerHand);
    result.dealerHand = dealerHand;
    
    // Check for blackjack
    if (playerHand.isBlackjack()) {
        if (dealerHand.isBlackjack()) {
            result.outcome = "push";
            result.winAmount = 0;
        } else {
            result.outcome = "blackjack";
            result.winAmount = bet.amount * (blackjackPayout_ / 2.0);
            result.multiplier = blackjackPayout_ / 2.0;
        }
        result.success = true;
        return result;
    }
    
    // Player actions
    std::string action = actions.count("action") ? actions.at("action") : "stand";
    
    if (action == "surrender" && surrenderAllowed_) {
        result.outcome = "surrender";
        result.winAmount = bet.amount * 0.5;
        result.multiplier = 0.5;
        result.success = true;
        return result;
    }
    
    // Hit loop
    while (action == "hit" && playerHand.getValue() < 21) {
        playerHand.cards.push_back(drawCard());
        if (playerHand.isBust()) {
            break;
        }
    }
    
    // Check bust
    if (playerHand.isBust()) {
        result.outcome = "bust";
        result.winAmount = 0;
        result.multiplier = 0;
        result.playerHands[0] = playerHand;
        result.success = true;
        return result;
    }
    
    // Dealer plays
    while (result.dealerHand.getValue() < 17) {
        result.dealerHand.cards.push_back(drawCard());
    }
    
    // Determine winner
    int playerValue = playerHand.getValue();
    int dealerValue = result.dealerHand.getValue();
    
    if (dealerValue > 21 || playerValue > dealerValue) {
        result.outcome = "win";
        result.winAmount = bet.amount * 2;
        result.multiplier = 2.0;
    } else if (playerValue < dealerValue) {
        result.outcome = "lose";
        result.winAmount = 0;
        result.multiplier = 0;
    } else {
        result.outcome = "push";
        result.winAmount = bet.amount;
        result.multiplier = 1.0;
    }
    
    result.success = true;
    return result;
}

// Baccarat server implementation
BaccaratServer::BaccaratServer()
    : TableGameServer("baccarat", "Baccarat", TableGameType::BACCARAT, 1.0, 10000.0, 0.98)
    , tieBetAllowed_(true)
    , pairBetAllowed_(true) {
}

std::vector<std::string> BaccaratServer::getValidActions() const {
    return {"player", "banker", "tie", "player_pair", "banker_pair"};
}

int BaccaratServer::getHandValue(const Hand& hand) const {
    int value = 0;
    for (const auto& card : hand.cards) {
        value += card.value;
    }
    return value % 10;
}

std::string BaccaratServer::determineWinner(const Hand& playerHand, const Hand& bankerHand) const {
    int playerValue = getHandValue(playerHand);
    int bankerValue = getHandValue(bankerHand);
    
    if (playerValue > bankerValue) return "player";
    if (bankerValue > playerValue) return "banker";
    return "tie";
}

TableGameResult BaccaratServer::play(const std::string& playerId, 
                                      const TableBet& bet,
                                      const std::map<std::string, std::string>& actions) {
    TableGameResult result;
    
    // Deal initial cards
    Hand playerHand;
    Hand bankerHand;
    
    playerHand.cards.push_back(drawCard());
    bankerHand.cards.push_back(drawCard());
    playerHand.cards.push_back(drawCard());
    bankerHand.cards.push_back(drawCard());
    
    // Check for natural
    int playerValue = getHandValue(playerHand);
    int bankerValue = getHandValue(bankerHand);
    
    if (playerValue >= 8 || bankerValue >= 8) {
        // Natural - no more cards
    } else {
        // Third card rules
        if (playerValue <= 5) {
            playerHand.cards.push_back(drawCard());
            playerValue = getHandValue(playerHand);
        }
        
        // Banker third card rule
        bool bankerDraws = false;
        if (playerHand.cards.size() == 3) {
            int thirdCardValue = playerHand.cards[2].value;
            if (thirdCardValue >= 0 && thirdCardValue <= 9) {
                if (bankerValue <= 2) {
                    bankerDraws = true;
                } else if (bankerValue == 3 && thirdCardValue != 8) {
                    bankerDraws = true;
                } else if (bankerValue == 4 && thirdCardValue >= 2 && thirdCardValue <= 7) {
                    bankerDraws = true;
                } else if (bankerValue == 5 && thirdCardValue >= 4 && thirdCardValue <= 7) {
                    bankerDraws = true;
                } else if (bankerValue == 6 && thirdCardValue >= 6 && thirdCardValue <= 7) {
                    bankerDraws = true;
                }
            }
        }
        
        if (bankerDraws) {
            bankerHand.cards.push_back(drawCard());
        }
    }
    
    result.playerHands.push_back(playerHand);
    result.dealerHand = bankerHand;
    
    // Determine winner
    std::string winner = determineWinner(playerHand, bankerHand);
    std::string selection = bet.selection;
    
    if (selection == winner) {
        double payout = 1.0;
        if (selection == "banker") {
            payout = 0.95;  // Banker has 5% commission
        }
        result.outcome = "win";
        result.winAmount = bet.amount * (1 + payout);
        result.multiplier = 1 + payout;
    } else {
        result.outcome = "lose";
        result.winAmount = 0;
        result.multiplier = 0;
    }
    
    result.success = true;
    return result;
}

// Roulette server implementation
RouletteServer::RouletteServer()
    : TableGameServer("roulette", "Roulette", TableGameType::ROULETTE, 1.0, 10000.0, 0.973) {
    // European wheel numbers
    wheelNumbers_ = {0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10, 5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26};
    
    // Bet selections
    betSelections_["red"] = {1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36};
    betSelections_["black"] = {2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35};
    betSelections_["even"] = {2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36};
    betSelections_["odd"] = {1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35};
    betSelections_["1st_dozen"] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12};
    betSelections_["2nd_dozen"] = {13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24};
    betSelections_["3rd_dozen"] = {25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36};
}

std::vector<std::string> RouletteServer::getValidActions() const {
    return {"red", "black", "even", "odd", "1st_dozen", "2nd_dozen", "3rd_dozen"};
}

int RouletteServer::spinWheel() {
    std::lock_guard<std::mutex> lock(rngMutex_);
    std::uniform_int_distribution<int> dist(0, 36);
    return wheelNumbers_[dist(rng_)];
}

std::string RouletteServer::getColor(int number) const {
    if (isGreen(number)) return "green";
    if (isRed(number)) return "red";
    return "black";
}

bool RouletteServer::isRed(int number) const {
    return std::find(betSelections_.at("red").begin(), 
                     betSelections_.at("red").end(), number) != betSelections_.at("red").end();
}

bool RouletteServer::isBlack(int number) const {
    return std::find(betSelections_.at("black").begin(), 
                     betSelections_.at("black").end(), number) != betSelections_.at("black").end();
}

bool RouletteServer::isGreen(int number) const {
    return number == 0;
}

TableGameResult RouletteServer::play(const std::string& playerId, 
                                      const TableBet& bet,
                                      const std::map<std::string, std::string>& actions) {
    TableGameResult result;
    
    int winningNumber = spinWheel();
    std::string winningColor = getColor(winningNumber);
    
    result.metadata["winningNumber"] = std::to_string(winningNumber);
    result.metadata["winningColor"] = winningColor;
    
    // Check if bet wins
    bool wins = false;
    double multiplier = 0;
    
    std::string selection = bet.selection;
    
    // Number bet
    if (selection == std::to_string(winningNumber)) {
        wins = true;
        multiplier = 35.0;  // 35:1 payout
    }
    // Color bets
    else if (selection == winningColor) {
        wins = true;
        multiplier = 1.0;  // 1:1 payout
    }
    // Dozen bets
    else if ((selection == "1st_dozen" && winningNumber >= 1 && winningNumber <= 12) ||
             (selection == "2nd_dozen" && winningNumber >= 13 && winningNumber <= 24) ||
             (selection == "3rd_dozen" && winningNumber >= 25 && winningNumber <= 36)) {
        wins = true;
        multiplier = 2.0;  // 2:1 payout
    }
    
    if (wins) {
        result.outcome = "win";
        result.winAmount = bet.amount * (1 + multiplier);
        result.multiplier = 1 + multiplier;
    } else {
        result.outcome = "lose";
        result.winAmount = 0;
        result.multiplier = 0;
    }
    
    result.success = true;
    return result;
}

// Poker server implementation
PokerServer::PokerServer()
    : TableGameServer("poker", "Texas Hold'em Poker", TableGameType::POKER, 1.0, 10000.0, 0.97) {
}

std::vector<std::string> PokerServer::getValidActions() const {
    return {"check", "bet", "fold", "call", "raise"};
}

int PokerServer::evaluateHand(const Hand& hand, const std::vector<Card>& community) const {
    // Simplified hand evaluation
    // In production, this would use a proper poker hand evaluator
    
    std::map<std::string, int> rankCounts;
    std::map<std::string, int> suitCounts;
    
    std::vector<Card> allCards = hand.cards;
    allCards.insert(allCards.end(), community.begin(), community.end());
    
    for (const auto& card : allCards) {
        rankCounts[card.rank]++;
        suitCounts[card.suit]++;
    }
    
    // Check for flush
    for (const auto& suit : suitCounts) {
        if (suit.second >= 5) return 8;  // Flush
    }
    
    // Check for pairs, trips, quads
    for (const auto& rank : rankCounts) {
        if (rank.second == 4) return 7;  // Four of a kind
        if (rank.second == 3) return 3;  // Three of a kind
        if (rank.second == 2) return 1;  // Pair
    }
    
    return 0;  // High card
}

std::string PokerServer::getHandRankName(int rank) const {
    switch (rank) {
        case 8: return "flush";
        case 7: return "four_of_a_kind";
        case 3: return "three_of_a_kind";
        case 1: return "pair";
        default: return "high_card";
    }
}

TableGameResult PokerServer::play(const std::string& playerId, 
                                   const TableBet& bet,
                                   const std::map<std::string, std::string>& actions) {
    TableGameResult result;
    
    // Simplified poker - just evaluate hole cards vs random community
    Hand playerHand;
    playerHand.cards.push_back(drawCard());
    playerHand.cards.push_back(drawCard());
    
    std::vector<Card> community;
    for (int i = 0; i < 5; i++) {
        community.push_back(drawCard());
    }
    
    result.playerHands.push_back(playerHand);
    
    // Evaluate
    int playerRank = evaluateHand(playerHand, community);
    
    // Simulate dealer rank (simplified)
    Hand dealerHand;
    dealerHand.cards.push_back(drawCard());
    dealerHand.cards.push_back(drawCard());
    int dealerRank = evaluateHand(dealerHand, community);
    
    result.dealerHand = dealerHand;
    
    if (playerRank > dealerRank) {
        result.outcome = "win";
        result.winAmount = bet.amount * 2;
        result.multiplier = 2.0;
    } else if (playerRank < dealerRank) {
        result.outcome = "lose";
        result.winAmount = 0;
        result.multiplier = 0;
    } else {
        result.outcome = "split";
        result.winAmount = bet.amount;
        result.multiplier = 1.0;
    }
    
    result.metadata["playerHand"] = getHandRankName(playerRank);
    result.metadata["dealerHand"] = getHandRankName(dealerRank);
    
    result.success = true;
    return result;
}

} // namespace TigerCasino
