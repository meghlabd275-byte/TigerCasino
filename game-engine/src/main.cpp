#include <iostream>
#include "GameEngine.hpp"

using namespace TigerCasino;

int main() {
    std::cout << "TigerCasino Game Engine" << std::endl;
    std::cout << "=======================" << std::endl << std::endl;
    
    // Create game engine
    GameEngine engine;
    
    // List available games
    std::cout << "Available Games:" << std::endl;
    auto games = engine.getAvailableGames();
    for (const auto& game : games) {
        std::cout << "  - " << game << std::endl;
    }
    std::cout << std::endl;
    
    // Test slot machine
    std::cout << "Testing Slot Machine:" << std::endl;
    BetInfo slotBet;
    slotBet.amount = 1.0;
    slotBet.balance = 100.0;
    
    auto slotResult = engine.play("slots", slotBet);
    std::cout << "  Outcome: " << slotResult.outcome << std::endl;
    std::cout << "  Win Amount: $" << slotResult.winAmount << std::endl;
    std::cout << "  Multiplier: " << slotResult.multiplier << "x" << std::endl << std::endl;
    
    // Test dice game
    std::cout << "Testing Dice Game:" << std::endl;
    BetInfo diceBet;
    diceBet.amount = 1.0;
    diceBet.balance = 100.0;
    diceBet.gameParams["target"] = 50.0;
    
    auto diceResult = engine.play("dice", diceBet);
    std::cout << "  Outcome: " << diceResult.outcome << std::endl;
    std::cout << "  Win Amount: $" << diceResult.winAmount << std::endl;
    std::cout << "  Multiplier: " << diceResult.multiplier << "x" << std::endl << std::endl;
    
    // Test roulette
    std::cout << "Testing Roulette:" << std::endl;
    BetInfo rouletteBet;
    rouletteBet.amount = 10.0;
    rouletteBet.balance = 1000.0;
    
    auto rouletteResult = engine.play("roulette", rouletteBet);
    std::cout << "  Outcome: " << rouletteResult.outcome << std::endl;
    std::cout << "  Win Amount: $" << rouletteResult.winAmount << std::endl;
    std::cout << "  Multiplier: " << rouletteResult.multiplier << "x" << std::endl << std::endl;
    
    // Show RTP values
    std::cout << "RTP Values:" << std::endl;
    std::cout << "  Slots: " << engine.getGameRTP("slots") * 100 << "%" << std::endl;
    std::cout << "  Dice: " << engine.getGameRTP("dice") * 100 << "%" << std::endl;
    std::cout << "  Roulette: " << engine.getGameRTP("roulette") * 100 << "%" << std::endl;
    std::cout << "  Blackjack: " << engine.getGameRTP("blackjack") * 100 << "%" << std::endl;
    
    std::cout << std::endl << "Game engine ready for production!" << std::endl;
    
    return 0;
}
