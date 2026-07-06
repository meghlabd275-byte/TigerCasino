#pragma once

#include <string>
#include <vector>
#include <map>
#include <memory>

namespace TigerCasino {

// Extended slot game configuration with 50+ games
struct ExtendedSlotConfig {
    std::string gameId;
    std::string gameName;
    std::string provider;
    int reels;
    int rows;
    int paylines;
    double minBet;
    double maxBet;
    double rtp;
    std::vector<std::string> symbols;
    std::map<std::string, double> symbolValues;
    std::map<std::string, int> symbolWeights;
    bool hasBonus;
    bool hasFreeSpins;
    bool hasJackpot;
    std::string volatility; // low, medium, high
    std::string theme; // adventure, classic, fantasy, etc.
};

// Extended slot game server with 50+ games
class ExtendedSlotGameServer {
private:
    std::map<std::string, ExtendedSlotConfig> gameTemplates_;
    std::map<std::string, std::shared_ptr<SlotGameServer>> activeGames_;
    
public:
    ExtendedSlotGameServer();
    ~ExtendedSlotGameServer() = default;
    
    // Initialize all 50+ slot games
    void initializeAllGames();
    
    // Get game by ID
    ExtendedSlotConfig getGameConfig(const std::string& gameId) const;
    
    // Get games by category
    std::vector<ExtendedSlotConfig> getGamesByProvider(const std::string& provider) const;
    std::vector<ExtendedSlotConfig> getGamesByTheme(const std::string& theme) const;
    std::vector<ExtendedSlotConfig> getGamesByVolatility(const std::string& volatility) const;
    
    // Get all games
    std::vector<ExtendedSlotConfig> getAllGames() const;
    
    // Get game count
    size_t getGameCount() const;
    
    // Get games by provider
    std::map<std::string, int> getProviderGameCounts() const;
};

// Game template creators - 50+ games
namespace SlotGameTemplates {

// Pragmatic Play style games (20 games)
ExtendedSlotConfig createWolfGold();
ExtendedSlotConfig createGreatRhino();
ExtendedSlotConfig createMustangGold();
ExtendedSlotConfig createFruitParty();
ExtendedSlotConfig createStarlightPrincess();
ExtendedSlotConfig createGatesOfOlympus();
ExtendedSlotConfig createTheDogHouse();
ExtendedSlotConfig createSweetBonanza();
ExtendedSlotConfig createBigBassBonanza();
ExtendedSlotConfig createJohnHunter();
ExtendedSlotConfig createMadameDestiny();
ExtendedSlotConfig createWildWestGold();
ExtendedSlotConfig createAztecGems();
ExtendedSlotConfig createPyramidKing();
ExtendedSlotConfig createGoldenOx();
ExtendedSlotConfig createMysterious();
ExtendedSlotConfig createPowerOfThor();
ExtendedSlotConfig createBuffaloKing();
ExtendedSlotConfig createDragonHot();
ExtendedSlotConfig createJokerJewels();

// NetEnt style games (10 games)
ExtendedSlotConfig createStarburst();
ExtendedSlotConfig createGonzoQuest();
ExtendedSlotConfig createDeadOrAlive();
ExtendedSlotConfig createTwinSpin();
ExtendedSlotConfig createHallOfGods();
ExtendedSlotConfig createMegaFortune();
ExtendedSlotConfig createBloodSuckers();
ExtendedSlotConfig createJackAndTheBeanstalk();
ExtendedSlotConfig createSteamTower();
ExtendedSlotConfig createFlowers();

// Play'n GO games (10 games)
ExtendedSlotConfig createBookOfDead();
ExtendedSlotConfig createReactoonz();
ExtendedSlotConfig createLegacyOfDead();
ExtendedSlotConfig createRiseOfMerlin();
ExtendedSlotConfig createMoonPrincess();
ExtendedSlotConfig createPlayInHell();
ExtendedSlotConfig createFireJoker();
ExtendedSlotConfig createRingOfOdysseus();
ExtendedSlotConfig createTomeOfMadness();
ExtendedSlotConfig createHoneyRush();

// Big Time Gaming / Megaways games (10 games)
ExtendedSlotConfig createBonanza();
ExtendedSlotConfig createExtraChilli();
ExtendedSlotConfig createWhiteRabbit();
ExtendedSlotConfig createHolyDiver();
ExtendedSlotConfig createWhoWantsToBeAMillionaire();
ExtendedSlotConfig createBuffaloLightning();
ExtendedSlotConfig createDiamondSun();
ExtendedSlotConfig createExtraChilli();
ExtendedSlotConfig createMonsoon();
ExtendedSlotConfig createDangerHighVoltage();

// Relax Gaming / Push Gaming (10 games)
ExtendedSlotConfig createMoneyTrain();
ExtendedSlotConfig createMoneyTrain2();
ExtendedSlotConfig createSpaceMiners();
ExtendedSlotConfig createTempleTumble();
ExtendedSlotConfig createDropz();
ExtendedSlotConfig createMegaStack();
ExtendedSlotConfig createPushUp();
ExtendedSlotConfig createRazorSharks();
ExtendedSlotConfig createStickyBirds();
ExtendedSlotConfig createWildFlowers();

// Hacksaw Gaming (10 games)
ExtendedSlotConfig createWantedDeadOrAlive();
ExtendedSlotConfig createStackEm();
ExtendedSlotConfig createTheBounty();
ExtendedSlotConfig createDice();
ExtendedSlotConfig createAurora();
ExtendedSlotConfig createChaosCrew();
ExtendedSlotConfig createCabinCrashers();
ExtendedSlotConfig createBeastMode();
ExtendedSlotConfig createOutsmart();
ExtendedSlotConfig createTheEmperor();

// Nolimit City (10 games)
ExtendedSlotConfig createSan Quentin();
ExtendedSlotConfig createBookOfShadows();
ExtendedSlotConfig createMental();
ExtendedSlotConfig createX x投;
ExtendedSlotConfig createDeadPanda();
ExtendedSlotConfig createFireInTheHole();
ExtendedSlotConfig createKiss();
ExtendedSlotConfig createLarryTheLeprechaun();
ExtendedSlotConfig createMillion777();
ExtendedSlotConfig createPontiusPilate();

} // namespace SlotGameTemplates

} // namespace TigerCasino
