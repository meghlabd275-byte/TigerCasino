#pragma once

#include <string>
#include <vector>
#include <map>
#include <array>
#include <memory>
#include <random>
#include <mutex>
#include <functional>

namespace TigerCasino {

// Card representation
struct Card {
    std::string suit;  // hearts, diamonds, clubs, spades
    std::string rank;  // 2-10, J, Q, K, A
    int value;
    
    std::string toString() const;
    bool isFaceCard() const;
    bool isAce() const;
};

// Hand representation
struct Hand {
    std::vector<Card> cards;
    bool isSplit = false;
    bool isDoubled = false;
    bool isSurrendered = false;
    
    int getValue() const;
    bool isBust() const;
    bool isBlackjack() const;
    bool canSplit() const;
    void clear();
};

// Table game types
enum class TableGameType {
    BLACKJACK,
    BACCARAT,
    POKER,
    ROULETTE,
    CRAPS
};

// Bet types
enum class BetType {
    STANDARD,
    SIDE_BET,
    PROP_BET
};

// Bet information
struct TableBet {
    std::string betId;
    std::string playerId;
    BetType type;
    double amount;
    std::string selection;
    std::map<std::string, std::string> metadata;
};

// Game result
struct TableGameResult {
    bool success;
    std::string outcome;
    double winAmount;
    double multiplier;
    std::vector<Hand> playerHands;
    Hand dealerHand;
    std::map<std::string, std::string> metadata;
    
    TableGameResult() : success(false), winAmount(0), multiplier(0) {}
};

// Base table game server
class TableGameServer {
protected:
    std::string gameId_;
    std::string gameName_;
    TableGameType gameType_;
    double minBet_;
    double maxBet_;
    double rtp_;
    std::mt19937_64 rng_;
    std::mutex rngMutex_;
    
    // Deck management
    std::vector<Card> shoe_;
    int currentCardIndex_;
    int decksInShoe_;
    
    void createShoe(int numDecks);
    Card drawCard();
    void shuffleShoe();
    
public:
    TableGameServer(const std::string& gameId, const std::string& gameName,
                    TableGameType type, double minBet, double maxBet, double rtp);
    virtual ~TableGameServer() = default;
    
    virtual TableGameResult play(const std::string& playerId, 
                                 const TableBet& bet,
                                 const std::map<std::string, std::string>& actions) = 0;
    
    std::string getGameId() const { return gameId_; }
    std::string getGameName() const { return gameName_; }
    TableGameType getGameType() const { return gameType_; }
    double getMinBet() const { return minBet_; }
    double getMaxBet() const { return maxBet_; }
    double getRTP() const { return rtp_; }
    
    virtual std::vector<std::string> getValidActions() const = 0;
};

// Blackjack server
class BlackjackServer : public TableGameServer {
private:
    bool surrenderAllowed_;
    bool doubleAfterSplitAllowed_;
    bool hitSplitAces_;
    int blackjackPayout_;
    
public:
    BlackjackServer();
    ~BlackjackServer() = default;
    
    TableGameResult play(const std::string& playerId, 
                        const TableBet& bet,
                        const std::map<std::string, std::string>& actions) override;
    
    std::vector<std::string> getValidActions() const override;
    
    // Special blackjack actions
    TableGameResult splitHand(const Hand& hand);
    TableGameResult doubleDown(const Hand& hand, const Card& dealerUpCard);
    TableGameResult surrender(const Hand& hand);
};

// Baccarat server
class BaccaratServer : public TableGameServer {
private:
    bool tieBetAllowed_;
    bool pairBetAllowed_;
    
public:
    BaccaratServer();
    ~BaccaratServer() = default;
    
    TableGameResult play(const std::string& playerId, 
                        const TableGameResult& bet,
                        const std::map<std::string, std::string>& actions) override;
    
    std::vector<std::string> getValidActions() const override;
    
    // Baccarat-specific
    int getHandValue(const Hand& hand) const;
    std::string determineWinner(const Hand& playerHand, const Hand& bankerHand) const;
};

// Roulette server
class RouletteServer : public TableGameServer {
private:
    std::vector<int> wheelNumbers_;  // European wheel: 0, 32, 15, ...
    std::map<std::string, std::vector<int>> betSelections_;
    
public:
    RouletteServer();
    ~RouletteServer() = default;
    
    TableGameResult play(const std::string& playerId, 
                        const TableBet& bet,
                        const std::map<std::string, std::string>& actions) override;
    
    std::vector<std::string> getValidActions() const override;
    
    // Roulette-specific
    int spinWheel();
    std::string getColor(int number) const;
    bool isRed(int number) const;
    bool isBlack(int number) const;
    bool isGreen(int number) const;
};

// Poker server (Texas Hold'em style)
class PokerServer : public TableGameServer {
private:
    std::map<std::string, int> handRankings_;
    
public:
    PokerServer();
    ~PokerServer() = default;
    
    TableGameResult play(const std::string& playerId, 
                        const TableBet& bet,
                        const std::map<std::string, std::string>& actions) override;
    
    std::vector<std::string> getValidActions() const override;
    
    // Poker-specific
    int evaluateHand(const Hand& hand, const std::vector<Card>& community) const;
    std::string getHandRankName(int rank) const;
};

} // namespace TigerCasino
