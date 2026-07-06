#pragma once

#include <string>
#include <vector>
#include <map>

namespace TigerCasino {

// Additional slot games - batch 2 (another 100 games)
struct AdditionalSlotConfig {
    std::string gameId;
    std::string gameName;
    std::string provider;
    int reels;
    int rows;
    int paylines;
    double minBet;
    double maxBet;
    double rtp;
    std::string volatility;
    std::string theme;
};

// Batch 2 slot games - Quick hits, Lightning links, etc.
namespace AdditionalSlotsBatch2 {

AdditionalSlotConfig createQuickHits();
AdditionalSlotConfig createLightningLink();
AdditionalSlotConfig create50Dragons();
AdditionalSlotConfig create40SuperHot();
AdditionalSlotConfig create100Rabbits();
AdditionalSlotConfig createHotFortune();
AdditionalSlotConfig createBurningPearl();
AdditionalSlotConfig createOrchidGarden();
AdditionalSlotConfig createAztecRiches();
AdditionalSlotConfig createDiamondEmpire();
AdditionalSlotConfig createFortuneKong();
AdditionalSlotConfig createRacingInferno();
AdditionalSlotConfig createLegendaryRed();
AdditionalSlotConfig createJackpotInferno();
AdditionalSlotConfig createGoldenTree();
AdditionalSlotConfig createRoyalKingdom();
AdditionalSlotConfig createCosmicCrystal();
AdditionalSlotConfig createDragonPhoenix();
AdditionalSlotConfig createLuckyBalls();
AdditionalSlotConfig createMonkeyKing();

// Yggdrasil games (20)
AdditionalSlotConfig createValleyOfTheGods();
AdditionalSlotConfig createRagnarok();
AdditionalSlotConfig createHolmesAndTheStolenStones();
AdditionalSlotConfig createTutsTwist();
AdditionalSlotConfig createGoldenFishTank();
AdditionalSlotConfig createEasterIsland();
AdditionalSlotConfig createPenguinCity();
AdditionalSlotConfig createWickedCircus();
AdditionalSlotConfig createJackpotExpress();
AdditionalSlotConfig createTheDrop();
AdditionalSlotConfig createCaishenRiches();
AdditionalSlotConfig createHades();
AdditionalSlotConfig createGiantGrizzly();
AdditionalSlotConfig createMysticWheel();
AdditionalSlotConfig createPirateQueens();
AdditionalSlotConfig createDragonReborn();
AdditionalSlotConfig createLotusWarrior();
AdditionalSlotConfig createGemsOfBuddha();
AdditionalSlotConfig createMedallion();
AdditionalSlotConfig createCrystalGeminis();

// Quickspin games (20)
AdditionalSlotConfig createSakuraWind();
AdditionalSlotConfig createDawnOfEgypt();
AdditionalSlotConfig createBlueOcean();
AdditionalSlotConfig createHiddenCity();
AdditionalSlotConfig createGrandSpinnSuperChip();
AdditionalSlotConfig createMightyArthur();
AdditionalSlotConfig createJokerStoker();
AdditionalSlotConfig createRiskyRabbit();
AdditionalSlotConfig createGoldilocks();
AdditionalSlotConfig createSpinionsBeachParty();
AdditionalSlotConfig createTheEpicQuest();
AdditionalSlotConfig createAliBaba();
AdditionalSlotConfig createTheLastKingOfAsgard();
AdditionalSlotConfig createFlyingDutchman();
AdditionalSlotConfig createVampireire();
AdditionalSlotConfig createKingColossus();
AdditionalSlotConfig createEl迭();
AdditionalSlotConfig createPersianWonders();
AdditionalSlotConfig createDiamondStrike();
AdditionalSlotConfig createFlipFlip();

// Playtech games (20)
AdditionalSlotConfig createAgeOfGods();
AdditionalSlotConfig createGladiator();
AdditionalSlotConfig createJackpotGiant();
AdditionalSlotConfig createMegaJackpots();
AdditionalSlotConfig createBeachLife();
AdditionalSlotConfig createLifeOfTheParty();
AdditionalSlotConfig createGoldRush();
AdditionalSlotConfig createSuperHeroes();
AdditionalSlotConfig createKingsTreasure();
AdditionalSlotConfig createPrinceOfOlympus();
AdditionalSlotConfig createThaiParadise();
AdditionalSlotConfig createKingOfCards();
AdditionalSlotConfig createWhiteOrchid();
AdditionalSlotConfig createSpartans();
AdditionalSlotConfig createStacksOfCash();
AdditionalSlotConfig createWildGems();
AdditionalSlotConfig createSuperChip();
AdditionalSlotConfig createPantherMoon();
AdditionalSlotConfig createJewelThief();
AdditionalSlotConfig createRocky();

// IGT games (20)
AdditionalSlotConfig createCleopatra();
AdditionalSlotConfig createDaVinciDiamonds();
AdditionalSlotConfig createWheelOfFortune();
AdditionalSlotConfig createStarTrek();
AdditionalSlotConfig createMonopoly();
AdditionalSlotConfig createKittyGlitter();
AdditionalSlotConfig createWolfRun();
AdditionalSlotConfig createSiberianStorm();
AdditionalSlotConfig createEgyptianRiches();
AdditionalSlotConfig createPixiesOfTheForest();
AdditionalSlotConfig createCrystalForest();
AdditionalSlotConfig createGongXiFaCai();
AdditionalSlotConfig createLuckyLarrys();
AdditionalSlotConfig createStarburstXXXtreme();
AdditionalSlotConfig createDoubleDiamond();
AdditionalSlotConfig createTripleDiamond();
AdditionalSlotConfig createRedWhiteAndBlue();
AdditionalSlotConfig createDoubleCasino();
AdditionalSlotConfig createMegaCrown();
AdditionalSlotConfig createGoldWins();

// Ainsworth games (20)
AdditionalSlotConfig createThunderCash();
AdditionalSlotConfig createJungleSpirit();
AdditionalSlotConfig createSunStrike();
AdditionalSlotConfig createPureFixedOdds();
AdditionalSlotConfig createSuperRedPhoenix();
AdditionalSlotConfig createWinStorm();
AdditionalSlotConfig createJungleJim();
AdditionalSlotConfig createKingmaker();
AdditionalSlotConfig createFlyingHigh();
AdditionalSlotConfig createMagicMonkey();
AdditionalSlotConfig createDiamondDuke();
AdditionalSlotConfig createRoaringForties();
AdditionalSlotConfig createJokerJackpots();
AdditionalSlotConfig createGrandGiant();
AdditionalSlotConfig createGoldenWolf();
AdditionalSlotConfig createApolloRising();
AdditionalSlotConfig createThunderCashLightning();
AdditionalSlotConfig createChinaRiver();
AdditionalSlotConfig createWildFrog();
AdditionalSlotConfig createJadePower();

// All additional slots manager
class AdditionalSlotsManager {
private:
    std::map<std::string, AdditionalSlotConfig> slots_;
    
public:
    AdditionalSlotsManager();
    std::vector<AdditionalSlotConfig> getAll() const;
    size_t count() const;
    AdditionalSlotConfig get(const std::string& id) const;
};

} // namespace AdditionalSlotsBatch2

} // namespace TigerCasino
