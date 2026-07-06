#ifndef ORIGINAL_GAMES_HPP
#define ORIGINAL_GAMES_HPP

#include <string>
#include <vector>
#include <memory>
#include <cmath>
#include <random>
#include "ProvablyFair.hpp"
#include "RandomNumberGenerator.hpp"

namespace TigerCasino {

/**
 * Original Games - Proprietary Games
 * These are unique to TigerCasino
 */

// ============ KENO ============
class KenoGame {
public:
    struct KenoResult {
        std::string roundId;
        std::vector<int> selectedNumbers;
        std::vector<int> drawnNumbers;
        int matches;
        double betAmount;
        double winAmount;
        double multiplier;
        bool isBonus;
    };

    static KenoResult play(const std::string& serverSeed,
                         const std::string& clientSeed,
                         uint64_t nonce,
                         const std::vector<int>& selections,
                         double betAmount) {
        KenoResult result;
        result.roundId = "KENO-" + std::to_string(nonce);
        result.selectedNumbers = selections;
        result.betAmount = betAmount;
        
        // Draw 20 numbers using provably fair
        std::vector<int> allNumbers;
        for (int i = 1; i <= 80; i++) allNumbers.push_back(i);
        
        std::vector<int> drawn;
        for (int i = 0; i < 20; i++) {
            uint64_t outcome = ProvablyFair::generateOutcome(
                serverSeed, clientSeed, nonce + i
            );
            int idx = outcome % allNumbers.size();
            drawn.push_back(allNumbers[idx]);
            allNumbers.erase(allNumbers.begin() + idx);
        }
        
        result.drawnNumbers = drawn;
        
        // Count matches
        int matches = 0;
        for (auto sel : selections) {
            for (auto dr : drawn) {
                if (sel == dr) matches++;
            }
        }
        result.matches = matches;
        
        // Calculate payout based on matches
        static double payouts[] = {0, 0, 0, 1, 2, 5, 10, 25, 100, 500, 2000};
        int idx = std::min(matches, 10);
        result.multiplier = payouts[idx];
        result.winAmount = betAmount * result.multiplier;
        result.isBonus = matches >= 7;
        
        return result;
    }
};

// ============ VIDEO POKER ============
class VideoPokerGame {
public:
    enum class HandType {
        Nothing,
        JacksOrBetter,
        TwoPair,
        ThreeOfAKind,
        Straight,
        Flush,
        FullHouse,
        FourOfAKind,
        StraightFlush,
        RoyalFlush
    };

    struct VideoPokerResult {
        std::string roundId;
        std::vector<std::string> hand;        // Current hand
        std::vector<std::string> finalHand;   // After draw
        std::vector<bool> held;               // Which cards held
        HandType handRank;
        double betAmount;
        double winAmount;
        double multiplier;
    };

    static HandType evaluateHand(const std::vector<std::string>& hand) {
        // Simplified evaluation
        // In production, would check actual poker hands
        
        // Count ranks
        std::map<std::string, int> rankCount;
        for (const auto& card : hand) {
            std::string rank = card.substr(0, card.length() - 1);
            rankCount[rank]++;
        }
        
        // Check for pairs, trips, quads
        for (const auto& rc : rankCount) {
            if (rc.second == 4) return HandType::FourOfAKind;
            if (rc.second == 3) return HandType::ThreeOfAKind;
            if (rc.second == 2) {
                // Check for two pair or jacks or better
                int pairCount = 0;
                for (const auto& rc2 : rankCount) {
                    if (rc2.second == 2) pairCount++;
                }
                if (pairCount == 2) return HandType::TwoPair;
                if (rc.first == "J" || rc.first == "Q" || rc.first == "K" || rc.first == "A") {
                    return HandType::JacksOrBetter;
                }
            }
        }
        
        return HandType::Nothing;
    }

    static VideoPokerResult play(const std::string& serverSeed,
                                const std::string& clientSeed,
                                uint64_t nonce,
                                const std::vector<std::string>& initialHand,
                                const std::vector<bool>& held,
                                double betAmount) {
        VideoPokerResult result;
        result.roundId = "VP-" + std::to_string(nonce);
        result.hand = initialHand;
        result.held = held;
        result.betAmount = betAmount;
        
        // Replace non-held cards
        std::vector<std::string> newHand = initialHand;
        std::vector<std::string> deck = createPokerDeck();
        
        for (size_t i = 0; i < held.size(); i++) {
            if (!held[i] && !deck.empty()) {
                uint64_t outcome = ProvablyFair::generateOutcome(
                    serverSeed, clientSeed, nonce + i
                );
                size_t idx = outcome % deck.size();
                newHand[i] = deck[idx];
                deck.erase(deck.begin() + idx);
            }
        }
        
        result.finalHand = newHand;
        result.handRank = evaluateHand(newHand);
        
        // Calculate payout
        static double payouts[] = {0, 1, 2, 3, 4, 6, 9, 25, 50, 250};
        int idx = static_cast<int>(result.handRank);
        result.multiplier = payouts[std::min(idx, 9)];
        result.winAmount = betAmount * result.multiplier;
        
        return result;
    }

private:
    static std::vector<std::string> createPokerDeck() {
        std::vector<std::string> deck;
        std::vector<std::string> suits = {"♠", "♥", "♦", "♣"};
        std::vector<std::string> ranks = {"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"};
        
        for (const auto& suit : suits) {
            for (const auto& rank : ranks) {
                deck.push_back(rank + suit);
            }
        }
        return deck;
    }
};

// ============ HILO ============
class HiLoGame {
public:
    struct HiLoResult {
        std::string roundId;
        std::string currentCard;
        std::string nextCard;
        std::string prediction;  // "higher", "lower", "equal"
        bool won;
        double betAmount;
        double winAmount;
        double multiplier;
        uint8_t streak;
    };

    static HiLoResult play(const std::string& serverSeed,
                          const std::string& clientSeed,
                          uint64_t nonce,
                          const std::string& currentCard,
                          const std::string& prediction,
                          double betAmount) {
        HiLoResult result;
        result.roundId = "HILO-" + std::to_string(nonce);
        result.currentCard = currentCard;
        result.prediction = prediction;
        result.betAmount = betAmount;
        
        // Draw next card
        std::vector<std::string> deck = createPokerDeck();
        
        uint64_t outcome = ProvablyFair::generateOutcome(
            serverSeed, clientSeed, nonce
        );
        size_t idx = outcome % deck.size();
        result.nextCard = deck[idx];
        
        // Get card values
        int currentValue = getCardValue(currentCard);
        int nextValue = getCardValue(result.nextCard);
        
        // Determine winner
        if (prediction == "higher") {
            result.won = nextValue > currentValue;
            result.multiplier = 1.95;
        } else if (prediction == "lower") {
            result.won = nextValue < currentValue;
            result.multiplier = 1.95;
        } else {
            result.won = nextValue == currentValue;
            result.multiplier = 9.0; // Very rare!
        }
        
        if (result.won) {
            result.winAmount = betAmount * result.multiplier;
        } else {
            result.winAmount = 0;
        }
        
        return result;
    }

private:
    static std::vector<std::string> createPokerDeck() {
        std::vector<std::string> deck;
        std::vector<std::string> suits = {"♠", "♥", "♦", "♣"};
        std::vector<std::string> ranks = {"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"};
        
        for (const auto& suit : suits) {
            for (const auto& rank : ranks) {
                deck.push_back(rank + suit);
            }
        }
        return deck;
    }
    
    static int getCardValue(const std::string& card) {
        std::string rank = card.substr(0, card.length() - 1);
        if (rank == "A") return 14;
        if (rank == "K") return 13;
        if (rank == "Q") return 12;
        if (rank == "J") return 11;
        return std::stoi(rank);
    }
};

// ============ DICE (Enhanced) ============
class EnhancedDiceGame {
public:
    struct DiceResult {
        std::string roundId;
        double roll;
        double target;
        std::string direction;  // "over" or "under"
        double multiplier;
        double betAmount;
        double winAmount;
        bool won;
        std::string serverSeed;
        std::string clientSeed;
    };

    static DiceResult play(const std::string& serverSeed,
                         const std::string& clientSeed,
                         uint64_t nonce,
                         double target,
                         std::string direction,
                         double betAmount,
                         double maxWin = 1000.0) {
        DiceResult result;
        result.roundId = "DICE-" + std::to_string(nonce);
        result.target = target;
        result.direction = direction;
        result.betAmount = betAmount;
        result.serverSeed = serverSeed;
        result.clientSeed = clientSeed;
        
        // Generate roll using provably fair
        uint64_t outcome = ProvablyFair::generateOutcome(
            serverSeed, clientSeed, nonce
        );
        result.roll = (outcome % 10001) / 100.0; // 0.00 - 100.00
        
        // Calculate multiplier
        if (direction == "over") {
            result.multiplier = (100.0 - target) / target;
            result.won = result.roll > target;
        } else {
            result.multiplier = target / (100.0 - target);
            result.won = result.roll < target;
        }
        
        // Cap multiplier for house edge
        result.multiplier = result.multiplier * 0.98; // 2% house edge
        result.multiplier = std::min(result.multiplier, maxWin);
        
        result.winAmount = result.won ? betAmount * result.multiplier : 0;
        
        return result;
    }
};

// ============ DIAMOND POKER ============
class DiamondPokerGame {
public:
    // Like plinko but with cards
    struct DiamondResult {
        std::string roundId;
        uint8_t row;
        uint8_t position;
        double multiplier;
        double betAmount;
        double winAmount;
        std::vector<uint8_t> path;
    };

    static DiamondResult play(const std::string& serverSeed,
                           const std::string& clientSeed,
                           uint64_t nonce,
                           uint8_t startPosition,
                           double betAmount) {
        DiamondResult result;
        result.roundId = "DIA-" + std::to_string(nonce);
        result.betAmount = betAmount;
        result.position = startPosition;
        
        // Payout table for diamond poker (various multipliers)
        double payouts[] = {0.5, 1.0, 2.0, 5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0};
        
        // Simulate drop through rows
        uint8_t rows = 11;
        result.path.push_back(startPosition);
        
        for (uint8_t r = 0; r < rows; r++) {
            uint64_t outcome = ProvablyFair::generateOutcome(
                serverSeed, clientSeed, nonce + r
            );
            bool goRight = (outcome % 2) == 1;
            
            if (goRight && result.position < r + 1) {
                result.position++;
            }
            result.path.push_back(result.position);
        }
        
        uint8_t finalPos = std::min(result.position, (uint8_t)(sizeof(payouts)/sizeof(payouts[0]) - 1));
        result.multiplier = payouts[finalPos];
        result.winAmount = betAmount * result.multiplier;
        
        return result;
    }
};

// ============ WHEEL OF FORTUNE ============
class WheelGame {
public:
    struct WheelResult {
        std::string roundId;
        double result;
        double multiplier;
        double betAmount;
        double winAmount;
        std::string segment;
    };

    static WheelResult spin(const std::string& serverSeed,
                          const std::string& clientSeed,
                          uint64_t nonce,
                          double betAmount) {
        WheelResult result;
        result.roundId = "WHEEL-" + std::to_string(nonce);
        result.betAmount = betAmount;
        
        // Wheel segments (54 segments typical)
        double segments[] = {
            0, 0, 0, 0, 0,  // 10x (5)
            1, 1, 1, 1,     // 5x (4)
            2, 2, 2, 2, 2, 2, // 2x (6)
            5, 5, 5, 5,      // 1.5x (4)
            10, 10,          // 1x (2)
            20,              // 0.5x (1)
            0,               // Jackpot (1)
        };
        
        uint64_t outcome = ProvablyFair::generateOutcome(
            serverSeed, clientSeed, nonce
        );
        
        size_t idx = outcome % (sizeof(segments)/sizeof(segments[0]));
        result.result = segments[idx];
        
        if (result.result == 0) {
            if (idx == 53) { // Jackpot
                result.multiplier = 1000.0;
                result.segment = "JACKPOT!";
            } else {
                result.multiplier = 0.0;
                result.segment = "0 (Bankrupt)";
            }
        } else if (result.result == 20) {
            result.multiplier = 0.5;
            result.segment = "20x";
        } else {
            result.multiplier = result.result;
            result.segment = std::to_string((int)result.result) + "x";
        }
        
        result.winAmount = betAmount * result.multiplier;
        
        return result;
    }
};

// ============ CLAIM CLAIM (Minefield) ============
class ClaimGame {
public:
    struct ClaimResult {
        std::string roundId;
        std::vector<int> mines;
        int currentStep;
        double multiplier;
        double betAmount;
        double winAmount;
        bool hitMine;
        bool claimed;
    };

    static ClaimResult play(const std::string& serverSeed,
                          const std::string& clientSeed,
                          uint64_t nonce,
                          int mineCount,
                          int steps,
                          double betAmount,
                          bool claim) {
        ClaimResult result;
        result.roundId = "CLAIM-" + std::to_string(nonce);
        result.betAmount = betAmount;
        result.currentStep = steps;
        
        // Generate mine positions
        std::vector<int> allPos;
        for (int i = 0; i < 25; i++) allPos.push_back(i);
        
        std::vector<int> mines;
        for (int i = 0; i < mineCount; i++) {
            uint64_t outcome = ProvablyFair::generateOutcome(
                serverSeed, clientSeed, nonce + i
            );
            int idx = outcome % allPos.size();
            mines.push_back(allPos[idx]);
            allPos.erase(allPos.begin() + idx);
        }
        result.mines = mines;
        
        // Calculate multiplier based on steps taken
        double baseMultiplier = 1.0;
        for (int i = 0; i < steps; i++) {
            baseMultiplier *= 1.2; // 20% increase per step
        }
        result.multiplier = baseMultiplier * 0.95; // 5% house edge
        
        // Check if hit mine
        result.hitMine = false;
        for (int m : mines) {
            if (m == steps - 1) {
                result.hitMine = true;
                break;
            }
        }
        
        if (result.hitMine) {
            result.winAmount = 0;
            result.claimed = false;
        } else if (claim) {
            result.winAmount = betAmount * result.multiplier;
            result.claimed = true;
        } else {
            result.winAmount = 0;
            result.claimed = false;
        }
        
        return result;
    }
};

} // namespace TigerCasino

#endif // ORIGINAL_GAMES_HPP
