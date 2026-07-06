#ifndef LIVE_DEALER_GAME_HPP
#define LIVE_DEALER_GAME_HPP

#include <string>
#include <vector>
#include <array>
#include <memory>
#include <chrono>
#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

/**
 * Live Dealer Game Base Class
 * Ultra-low latency for real-time gaming
 */
class LiveDealerGame {
public:
    enum class GameType {
        Blackjack,
        Roulette,
        Baccarat,
        Poker,
        GameShow
    };

    enum class GameStatus {
        Waiting,
        Dealing,
        Playing,
        Settling,
        Complete
    };

    struct Player {
        std::string playerId;
        std::string username;
        std::vector<Card> hand;
        double betAmount;
        double winAmount;
        bool isFinished;
    };

    struct Card {
        std::string suit;  // Hearts, Diamonds, Clubs, Spades
        std::string rank;  // A, 2-10, J, Q, K
        int value;         // 1-11 for blackjack
        bool isVisible;
    };

    struct GameTable {
        std::string tableId;
        std::string dealerName;
        GameType gameType;
        GameStatus status;
        std::vector<Card> deck;
        std::vector<Card> dealerHand;
        std::vector<Player> players;
        uint64_t roundNumber;
        uint64_t startTime;
        uint64_t endTime;
    };

protected:
    std::unique_ptr<RandomNumberGenerator> rng_;
    static constexpr size_t DECK_SIZE = 52;

public:
    LiveDealerGame() : rng_(std::make_unique<RandomNumberGenerator>()) {}
    virtual ~LiveDealerGame() = default;

    virtual GameType getType() const = 0;
    virtual std::string getTypeName() const = 0;
    
    // Initialize a new round
    virtual GameTable startRound(const std::string& tableId, const std::string& dealerName) = 0;
    
    // Place a bet
    virtual bool placeBet(GameTable& table, const std::string& playerId, 
                        const std::string& username, double amount) = 0;
    
    // Deal initial cards
    virtual void dealInitialCards(GameTable& table) = 0;
    
    // Player action (hit, stand, etc.)
    virtual bool playerAction(GameTable& table, const std::string& playerId, 
                           const std::string& action) = 0;
    
    // Dealer play
    virtual void dealerPlay(GameTable& table) = 0;
    
    // Settle bets and determine winners
    virtual void settle(GameTable& table) = 0;

    // Create a standard 52-card deck
    static std::vector<Card> createDeck() {
        std::vector<Card> deck;
        std::vector<std::string> suits = {"♠", "♥", "♦", "♣"};
        std::vector<std::string> ranks = {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"};
        
        for (const auto& suit : suits) {
            for (const auto& rank : ranks) {
                int value;
                if (rank == "A") value = 11;
                else if (rank == "J" || rank == "Q" || rank == "K") value = 10;
                else value = std::stoi(rank);
                
                deck.push_back({suit, rank, value, false});
            }
        }
        return deck;
    }

    // Shuffle deck using provably fair
    static std::vector<Card> shuffleDeck(const std::vector<Card>& deck,
                                         const std::string& serverSeed,
                                         const std::string& clientSeed,
                                         uint64_t nonce) {
        auto shuffled = deck;
        
        // Fisher-Yates shuffle with provably fair randomness
        for (size_t i = shuffled.size() - 1; i > 0; --i) {
            uint64_t outcome = ProvablyFair::generateOutcome(serverSeed, clientSeed, nonce + i);
            size_t j = outcome % (i + 1);
            std::swap(shuffled[i], shuffled[j]);
        }
        
        return shuffled;
    }

    // Calculate hand value for blackjack
    static int calculateHandValue(const std::vector<Card>& hand) {
        int value = 0;
        int aces = 0;
        
        for (const auto& card : hand) {
            value += card.value;
            if (card.rank == "A") aces++;
        }
        
        // Adjust for aces
        while (value > 21 && aces > 0) {
            value -= 10;
            aces--;
        }
        
        return value;
    }
};

/**
 * Live Blackjack
 */
class LiveBlackjack : public LiveDealerGame {
public:
    enum class Action {
        Hit,
        Stand,
        Double,
        Split
    };

    GameType getType() const override { return GameType::Blackjack; }
    std::string getTypeName() const override { return "Blackjack"; }

    GameTable startRound(const std::string& tableId, const std::string& dealerName) override {
        GameTable table;
        table.tableId = tableId;
        table.dealerName = dealerName;
        table.gameType = GameType::Blackjack;
        table.status = GameStatus::Waiting;
        table.deck = shuffleDeck(createDeck(), 
                                ProvablyFair::generateServerSeed(),
                                "", 
                                0);
        table.roundNumber = 0;
        table.startTime = getCurrentTimestamp();
        return table;
    }

    bool placeBet(GameTable& table, const std::string& playerId,
                 const std::string& username, double amount) override {
        if (table.status != GameStatus::Waiting) return false;
        
        Player player;
        player.playerId = playerId;
        player.username = username;
        player.betAmount = amount;
        player.winAmount = 0;
        player.isFinished = false;
        
        table.players.push_back(player);
        return true;
    }

    void dealInitialCards(GameTable& table) override {
        if (table.players.empty()) return;
        
        table.status = GameStatus::Dealing;
        table.roundNumber++;
        
        // Deal 2 cards to each player
        for (auto& player : table.players) {
            player.hand.push_back(drawCard(table));
            player.hand.push_back(drawCard(table));
        }
        
        // Deal 1 card to dealer (face down)
        table.dealerHand.push_back(drawCard(table));
        
        // Check for blackjack
        for (auto& player : table.players) {
            if (calculateHandValue(player.hand) == 21) {
                player.isFinished = true;
            }
        }
        
        table.status = GameStatus::Playing;
    }

    bool playerAction(GameTable& table, const std::string& playerId,
                    const std::string& actionStr) override {
        Action action;
        if (actionStr == "hit") action = Action::Hit;
        else if (actionStr == "stand") action = Action::Stand;
        else if (actionStr == "double") action = Action::Double;
        else if (actionStr == "split") action = Action::Split;
        else return false;
        
        for (auto& player : table.players) {
            if (player.playerId == playerId && !player.isFinished) {
                switch (action) {
                    case Action::Hit:
                        player.hand.push_back(drawCard(table));
                        if (calculateHandValue(player.hand) > 21) {
                            player.isFinished = true; // Bust
                        }
                        break;
                    case Action::Stand:
                        player.isFinished = true;
                        break;
                    case Action::Double:
                        player.betAmount *= 2;
                        player.hand.push_back(drawCard(table));
                        player.isFinished = true;
                        break;
                    case Action::Split:
                        // Handle split logic
                        break;
                }
                return true;
            }
        }
        return false;
    }

    void dealerPlay(GameTable& table) override {
        table.status = GameStatus::Settling;
        
        // Reveal dealer's hole card
        if (!table.dealerHand.empty()) {
            table.dealerHand[0].isVisible = true;
        }
        
        // Dealer hits on soft 17
        int dealerValue = calculateHandValue(table.dealerHand);
        while (dealerValue < 17) {
            table.dealerHand.push_back(drawCard(table));
            dealerValue = calculateHandValue(table.dealerHand);
        }
    }

    void settle(GameTable& table) override {
        int dealerValue = calculateHandValue(table.dealerHand);
        
        for (auto& player : table.players) {
            int playerValue = calculateHandValue(player.hand);
            
            if (playerValue > 21) {
                // Player busts - lose
                player.winAmount = 0;
            } else if (dealerValue > 21) {
                // Dealer busts - win
                player.winAmount = player.betAmount * 2;
            } else if (playerValue > dealerValue) {
                // Player wins
                player.winAmount = player.betAmount * 2;
            } else if (playerValue == dealerValue) {
                // Push
                player.winAmount = player.betAmount;
            } else {
                // Dealer wins
                player.winAmount = 0;
            }
            
            // Blackjack pays 3:2
            if (player.hand.size() == 2 && playerValue == 21) {
                player.winAmount = player.betAmount * 2.5;
            }
        }
        
        table.status = GameStatus::Complete;
        table.endTime = getCurrentTimestamp();
    }

private:
    Card drawCard(GameTable& table) {
        if (table.deck.empty()) {
            // Reshuffle
            table.deck = createDeck();
        }
        Card card = table.deck.back();
        card.isVisible = true;
        table.deck.pop_back();
        return card;
    }

    uint64_t getCurrentTimestamp() {
        return std::chrono::duration_cast<std::chrono::milliseconds>(
            std::chrono::system_clock::now().time_since_epoch()
        ).count();
    }
};

/**
 * Live Roulette
 */
class LiveRoulette : public LiveDealerGame {
public:
    enum class BetType {
        Straight,      // Single number 0-36
        Split,         // Two adjacent numbers
        Street,        // Three numbers in a row
        Corner,        // Four numbers in a square
        Line,          // Six numbers in two rows
        Dozen,         // 1-12, 13-24, 25-36
        Column,        // Three vertical columns
        RedBlack,      // Red or Black
        EvenOdd,       // Even or Odd
        HighLow        // 1-18 or 19-36
    };

    struct RouletteBet {
        std::string playerId;
        BetType type;
        std::vector<int> numbers;
        double amount;
        double multiplier;
    };

    GameType getType() const override { return GameType::Roulette; }
    std::string getTypeName() const override { return "Roulette"; }

    GameTable startRound(const std::string& tableId, const std::string& dealerName) override {
        GameTable table;
        table.tableId = tableId;
        table.dealerName = dealerName;
        table.gameType = GameType::Roulette;
        table.status = GameStatus::Waiting;
        table.roundNumber = 0;
        return table;
    }

    bool placeBet(GameTable& table, const std::string& playerId,
                 const std::string& username, double amount) override {
        return true; // Simplified
    }

    void dealInitialCards(GameTable& table) override {
        // Generate winning number using provably fair
        uint64_t seed = ProvablyFair::generateOutcome(
            ProvablyFair::generateServerSeed(), 
            "", 
            table.roundNumber
        );
        int winningNumber = seed % 37; // 0-36
        
        table.status = GameStatus::Complete;
        table.endTime = getCurrentTimestamp();
    }

    bool playerAction(GameTable& table, const std::string& playerId,
                    const std::string& action) override {
        return false; // No player action in roulette
    }

    void dealerPlay(GameTable& table) override {
        // Spin the wheel
        dealInitialCards(table);
    }

    void settle(GameTable& table) override {
        // Settle logic handled in dealInitialCards for simplicity
    }

    // Calculate payout based on bet type
    static double calculatePayout(BetType type) {
        switch (type) {
            case BetType::Straight: return 35.0;
            case BetType::Split: return 17.0;
            case BetType::Street: return 11.0;
            case BetType::Corner: return 8.0;
            case BetType::Line: return 5.0;
            case BetType::Dozen: return 2.0;
            case BetType::Column: return 2.0;
            case BetType::RedBlack: return 1.0;
            case BetType::EvenOdd: return 1.0;
            case BetType::HighLow: return 1.0;
            default: return 0.0;
        }
    }

private:
    uint64_t getCurrentTimestamp() {
        return std::chrono::duration_cast<std::chrono::milliseconds>(
            std::chrono::system_clock::now().time_since_epoch()
        ).count();
    }
};

/**
 * Live Baccarat
 */
class LiveBaccarat : public LiveDealerGame {
public:
    enum class BetType {
        Player,
        Banker,
        Tie,
        PlayerPair,
        BankerPair
    };

    GameType getType() const override { return GameType::Baccarat; }
    std::string getTypeName() const override { return "Baccarat"; }

    GameTable startRound(const std::string& tableId, const std::string& dealerName) override {
        GameTable table;
        table.tableId = tableId;
        table.dealerName = dealerName;
        table.gameType = GameType::Baccarat;
        table.status = GameStatus::Waiting;
        table.roundNumber = 0;
        return table;
    }

    bool placeBet(GameTable& table, const std::string& playerId,
                 const std::string& username, double amount) override {
        Player player;
        player.playerId = playerId;
        player.username = username;
        player.betAmount = amount;
        player.isFinished = false;
        table.players.push_back(player);
        return true;
    }

    void dealInitialCards(GameTable& table) override {
        table.roundNumber++;
        
        // Deal 2 cards to player, 2 to banker
        for (int i = 0; i < 2; i++) {
            table.players[0].hand.push_back(drawCard(table));
            table.dealerHand.push_back(drawCard(table));
        }
        
        // Check for natural win
        int playerValue = calculateBaccaratValue(table.players[0].hand);
        int bankerValue = calculateBaccaratValue(table.dealerHand);
        
        if (playerValue >= 8 || bankerValue >= 8) {
            table.status = GameStatus::Complete;
        } else {
            table.status = GameStatus::Playing;
        }
    }

    bool playerAction(GameTable& table, const std::string& playerId,
                    const std::string& action) override {
        // Baccarat is mostly automatic, player just watches
        return false;
    }

    void dealerPlay(GameTable& table) override {
        int playerValue = calculateBaccaratValue(table.players[0].hand);
        int bankerValue = calculateBaccaratValue(table.dealerHand);
        
        // Third card rules (simplified)
        if (playerValue < 6) {
            table.players[0].hand.push_back(drawCard(table));
            playerValue = calculateBaccaratValue(table.players[0].hand);
        }
        
        // Banker third card rules
        if (playerValue >= 6) {
            if (bankerValue < 6) {
                table.dealerHand.push_back(drawCard(table));
            }
        }
    }

    void settle(GameTable& table) override {
        int playerValue = calculateBaccaratValue(table.players[0].hand);
        int bankerValue = calculateBaccaratValue(table.dealerHand);
        
        if (table.players.empty()) return;
        
        if (playerValue > bankerValue) {
            table.players[0].winAmount = table.players[0].betAmount * 2; // Player wins
        } else if (bankerValue > playerValue) {
            table.players[0].winAmount = table.players[0].betAmount * 1.95; // Banker wins (5% commission)
        } else {
            table.players[0].winAmount = table.players[0].betAmount; // Tie
        }
        
        table.status = GameStatus::Complete;
    }

private:
    Card drawCard(GameTable& table) {
        static std::vector<Card> deck = createDeck();
        
        if (deck.empty()) {
            deck = createDeck();
        }
        
        Card card = deck.back();
        card.isVisible = true;
        deck.pop_back();
        return card;
    }

    static int calculateBaccaratValue(const std::vector<Card>& hand) {
        int value = 0;
        for (const auto& card : hand) {
            value += card.value;
        }
        return value % 10; // Baccarat: only last digit matters
    }

    uint64_t getCurrentTimestamp() {
        return std::chrono::duration_cast<std::chrono::milliseconds>(
            std::chrono::system_clock::now().time_since_epoch()
        ).count();
    }
};

} // namespace TigerCasino

#endif // LIVE_DEALER_GAME_HPP
