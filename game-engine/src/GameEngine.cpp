#include "GameEngine.hpp"
#include "SlotMachine.hpp"
#include "DiceGame.hpp"
#include "Roulette.hpp"
#include "Blackjack.hpp"

namespace TigerCasino {

GameEngine::GameEngine() {
    // Register all available games
    registerGame(std::make_shared<SlotMachine>());
    registerGame(std::make_shared<DiceGame>());
    registerGame(std::make_shared<Roulette>());
    registerGame(std::make_shared<Blackjack>());
}

void GameEngine::registerGame(std::shared_ptr<Game> game) {
    games_[game->getType()] = game;
}

GameResult GameEngine::play(const std::string& gameType, const BetInfo& bet) {
    auto it = games_.find(gameType);
    if (it == games_.end()) {
        GameResult result;
        result.success = false;
        result.outcome = "Game not found";
        return result;
    }
    
    return it->second->play(bet);
}

std::vector<std::string> GameEngine::getAvailableGames() const {
    std::vector<std::string> games;
    for (const auto& [type, game] : games_) {
        games.push_back(game->getName());
    }
    return games;
}

double GameEngine::getGameRTP(const std::string& gameType) const {
    auto it = games_.find(gameType);
    if (it == games_.end()) {
        return 0.0;
    }
    return it->second->getRTP();
}

} // namespace TigerCasino
