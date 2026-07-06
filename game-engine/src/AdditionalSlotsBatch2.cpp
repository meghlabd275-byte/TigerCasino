#include "AdditionalSlotsBatch2.hpp"

namespace TigerCasino {

AdditionalSlotsManager::AdditionalSlotsManager() {
    using namespace AdditionalSlotsBatch2;
    
    // Quick hits style games (20)
    slots_["quick_hits"] = createQuickHits();
    slots_["lightning_link"] = createLightningLink();
    slots_["50_dragons"] = create50Dragons();
    slots_["40_super_hot"] = create40SuperHot();
    slots_["100_rabbits"] = create100Rabbits();
    slots_["hot_fortune"] = createHotFortune();
    slots_["burning_pearl"] = createBurningPearl();
    slots_["orchid_garden"] = createOrchidGarden();
    slots_["aztec_riches"] = createAztecRiches();
    slots_["diamond_empire"] = createDiamondEmpire();
    slots_["fortune_kong"] = createFortuneKong();
    slots_["racing_inferno"] = createRacingInferno();
    slots_["legendary_red"] = createLegendaryRed();
    slots_["jackpot_inferno"] = createJackpotInferno();
    slots_["golden_tree"] = createGoldenTree();
    slots_["royal_kingdom"] = createRoyalKingdom();
    slots_["cosmic_crystal"] = createCosmicCrystal();
    slots_["dragon_phoenix"] = createDragonPhoenix();
    slots_["lucky_balls"] = createLuckyBalls();
    slots_["monkey_king"] = createMonkeyKing();
    
    // Yggdrasil (20)
    slots_["valley_of_gods"] = createValleyOfTheGods();
    slots_["ragnarok"] = createRagnarok();
    slots_["holmes"] = createHolmesAndTheStolenStones();
    slots_["tuts_twist"] = createTutsTwist();
    slots_["golden_fish_tank"] = createGoldenFishTank();
    slots_["easter_island"] = createEasterIsland();
    slots_["penguin_city"] = createPenguinCity();
    slots_["wicked_circus"] = createWickedCircus();
    slots_["jackpot_express"] = createJackpotExpress();
    slots_["the_drop"] = createTheDrop();
    slots_["caishen_riches"] = createCaishenRiches();
    slots_["hades"] = createHades();
    slots_["giant_grizzly"] = createGiantGrizzly();
    slots_["mystic_wheel"] = createMysticWheel();
    slots_["pirate_queens"] = createPirateQueens();
    slots_["dragon_reborn"] = createDragonReborn();
    slots_["lotus_warrior"] = createLotusWarrior();
    slots_["gems_of_buddha"] = createGemsOfBuddha();
    slots_["medallion"] = createMedallion();
    slots_["crystal_geminis"] = createCrystalGeminis();
    
    // Quickspin (20)
    slots_["sakura_wind"] = createSakuraWind();
    slots_["dawn_of_egypt"] = createDawnOfEgypt();
    slots_["blue_ocean"] = createBlueOcean();
    slots_["hidden_city"] = createHiddenCity();
    slots_["grand_spinn"] = createGrandSpinnSuperChip();
    slots_["mighty_arthur"] = createMightyArthur();
    slots_["joker_stoker"] = createJokerStoker();
    slots_["risky_rabbit"] = createRiskyRabbit();
    slots_["goldilocks"] = createGoldilocks();
    slots_["spinions"] = createSpinionsBeachParty();
    slots_["epic_quest"] = createTheEpicQuest();
    slots_["ali_baba"] = createAliBaba();
    slots_["last_king"] = createTheLastKingOfAsgard();
    slots_["flying_dutchman"] = createFlyingDutchman();
    slots_["vampireire"] = createVampireire();
    slots_["king_colossus"] = createKingColossus();
    slots_["el_dorado"] = createEl迭();
    slots_["persian_wonders"] = createPersianWonders();
    slots_["diamond_strike"] = createDiamondStrike();
    slots_["flip_flip"] = createFlipFlip();
    
    // Playtech (20)
    slots_["age_of_gods"] = createAgeOfGods();
    slots_["gladiator"] = createGladiator();
    slots_["jackpot_giant"] = createJackpotGiant();
    slots_["mega_jackpots"] = createMegaJackpots();
    slots_["beach_life"] = createBeachLife();
    slots_["life_of_party"] = createLifeOfTheParty();
    slots_["gold_rush"] = createGoldRush();
    slots_["super_heroes"] = createSuperHeroes();
    slots_["kings_treasure"] = createKingsTreasure();
    slots_["prince_olympus"] = createPrinceOfOlympus();
    slots_["thai_paradise"] = createThaiParadise();
    slots_["king_of_cards"] = createKingOfCards();
    slots_["white_orchid"] = createWhiteOrchid();
    slots_["spartans"] = createSpartans();
    slots_["stacks_of_cash"] = createStacksOfCash();
    slots_["wild_gems"] = createWildGems();
    slots_["super_chip"] = createSuperChip();
    slots_["panther_moon"] = createPantherMoon();
    slots_["jewel_thief"] = createJewelThief();
    slots_["rocky"] = createRocky();
    
    // IGT (20)
    slots_["cleopatra"] = createCleopatra();
    slots_["davinci_diamonds"] = createDaVinciDiamonds();
    slots_["wheel_of_fortune"] = createWheelOfFortune();
    slots_["star_trek"] = createStarTrek();
    slots_["monopoly"] = createMonopoly();
    slots_["kitty_glitter"] = createKittyGlitter();
    slots_["wolf_run"] = createWolfRun();
    slots_["siberian_storm"] = createSiberianStorm();
    slots_["egyptian_riches"] = createEgyptianRiches();
    slots_["pixies_forest"] = createPixiesOfTheForest();
    slots_["crystal_forest"] = createCrystalForest();
    slots_["gong_xi"] = createGongXiFaCai();
    slots_["lucky_larrys"] = createLuckyLarrys();
    slots_["starburst_xtreme"] = createStarburstXXXtreme();
    slots_["double_diamond"] = createDoubleDiamond();
    slots_["triple_diamond"] = createTripleDiamond();
    slots_["red_white_blue"] = createRedWhiteAndBlue();
    slots_["double_casino"] = createDoubleCasino();
    slots_["mega_crown"] = createMegaCrown();
    slots_["gold_wins"] = createGoldWins();
    
    // Ainsworth (20)
    slots_["thunder_cash"] = createThunderCash();
    slots_["jungle_spirit"] = createJungleSpirit();
    slots_["sun_strike"] = createSunStrike();
    slots_["pure_fixed"] = createPureFixedOdds();
    slots_["super_red_phoenix"] = createSuperRedPhoenix();
    slots_["win_storm"] = createWinStorm();
    slots_["jungle_jim"] = createJungleJim();
    slots_["kingmaker"] = createKingmaker();
    slots_["flying_high"] = createFlyingHigh();
    slots_["magic_monkey"] = createMagicMonkey();
    slots_["diamond_duke"] = createDiamondDuke();
    slots_["roaring_forties"] = createRoaringForties();
    slots_["joker_jackpots"] = createJokerJackpots();
    slots_["grand_giant"] = createGrandGiant();
    slots_["golden_wolf"] = createGoldenWolf();
    slots_["apollo_rising"] = createApolloRising();
    slots_["thunder_lightning"] = createThunderCashLightning();
    slots_["china_river"] = createChinaRiver();
    slots_["wild_frog"] = createWildFrog();
    slots_["jade_power"] = createJadePower();
}

std::vector<AdditionalSlotConfig> AdditionalSlotsManager::getAll() const {
    std::vector<AdditionalSlotConfig> result;
    for (const auto& p : slots_) result.push_back(p.second);
    return result;
}

size_t AdditionalSlotsManager::count() const { return slots_.size(); }

AdditionalSlotConfig AdditionalSlotsManager::get(const std::string& id) const {
    auto it = slots_.find(id);
    return it != slots_.end() ? it->second : AdditionalSlotConfig{};
}

// Template implementations
namespace AdditionalSlotsBatch2 {

AdditionalSlotConfig createQuickHits() {
    return {"quick_hits", "Quick Hits", "Ainsworth", 5, 3, 25, 0.05, 100.0, 0.94, "medium", "classic"};
}
AdditionalSlotConfig createLightningLink() {
    return {"lightning_link", "Lightning Link", "Ainsworth", 5, 3, 25, 0.05, 100.0, 0.95, "medium", "progressive"};
}
AdditionalSlotConfig create50Dragons() {
    return {"50_dragons", "50 Dragons", "Ainsworth", 5, 3, 50, 0.50, 500.0, 0.95, "medium", "asian"};
}
AdditionalSlotConfig create40SuperHot() {
    return {"40_super_hot", "40 Super Hot", "Ainsworth", 5, 3, 40, 0.40, 400.0, 0.95, "low", "classic"};
}
AdditionalSlotConfig create100Rabbits() {
    return {"100_rabbits", "100 Rabbits", "Ainsworth", 5, 3, 100, 0.20, 200.0, 0.95, "medium", "nature"};
}
AdditionalSlotConfig createHotFortune() {
    return {"hot_fortune", "Hot Fortune", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "asian"};
}
AdditionalSlotConfig createBurningPearl() {
    return {"burning_pearl", "Burning Pearl", "Ainsworth", 5, 3, 50, 0.40, 400.0, 0.95, "high", "ocean"};
}
AdditionalSlotConfig createOrchidGarden() {
    return {"orchid_garden", "Orchid Garden", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "nature"};
}
AdditionalSlotConfig createAztecRiches() {
    return {"aztec_riches", "Aztec Riches", "Ainsworth", 5, 3, 25, 0.30, 300.0, 0.95, "medium", "aztec"};
}
AdditionalSlotConfig createDiamondEmpire() {
    return {"diamond_empire", "Diamond Empire", "Ainsworth", 6, 5, 0, 0.20, 200.0, 0.96, "high", "luxury"};
}
AdditionalSlotConfig createFortuneKong() {
    return {"fortune_kong", "Fortune Kong", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "animal"};
}
AdditionalSlotConfig createRacingInferno() {
    return {"racing_inferno", "Racing Inferno", "Ainsworth", 5, 3, 25, 0.30, 300.0, 0.95, "high", "racing"};
}
AdditionalSlotConfig createLegendaryRed() {
    return {"legendary_red", "Legendary Red", "Ainsworth", 5, 3, 25, 0.40, 400.0, 0.95, "high", "dragon"};
}
AdditionalSlotConfig createJackpotInferno() {
    return {"jackpot_inferno", "Jackpot Inferno", "Ainsworth", 5, 3, 25, 0.50, 500.0, 0.94, "high", "jackpot"};
}
AdditionalSlotConfig createGoldenTree() {
    return {"golden_tree", "Golden Tree", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "fantasy"};
}
AdditionalSlotConfig createRoyalKingdom() {
    return {"royal_kingdom", "Royal Kingdom", "Ainsworth", 5, 3, 25, 0.30, 300.0, 0.95, "medium", "kingdom"};
}
AdditionalSlotConfig createCosmicCrystal() {
    return {"cosmic_crystal", "Cosmic Crystal", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.96, "high", "space"};
}
AdditionalSlotConfig createDragonPhoenix() {
    return {"dragon_phoenix", "Dragon Phoenix", "Ainsworth", 5, 3, 50, 0.50, 500.0, 0.95, "high", "dragon"};
}
AdditionalSlotConfig createLuckyBalls() {
    return {"lucky_balls", "Lucky Balls", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "lottery"};
}
AdditionalSlotConfig createMonkeyKing() {
    return {"monkey_king", "Monkey King", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "mythology"};
}

// Yggdrasil
AdditionalSlotConfig createValleyOfTheGods() { return {"valley_gods", "Valley of the Gods", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "egypt"}; }
AdditionalSlotConfig createRagnarok() { return {"ragnarok", "Ragnarok", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "mythology"}; }
AdditionalSlotConfig createHolmesAndTheStolenStones() { return {"holmes", "Holmes", "Yggdrasil", 5, 3, 20, 0.20, 100.0, 0.96, "medium", "adventure"}; }
AdditionalSlotConfig createTutsTwist() { return {"tuts_twist", "Tut's Twist", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "egypt"}; }
AdditionalSlotConfig createGoldenFishTank() { return {"gold_fish", "Golden Fish Tank", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "ocean"}; }
AdditionalSlotConfig createEasterIsland() { return {"easter", "Easter Island", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "island"}; }
AdditionalSlotConfig createPenguinCity() { return {"penguin", "Penguin City", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "animal"}; }
AdditionalSlotConfig createWickedCircus() { return {"wicked", "Wicked Circus", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "circus"}; }
AdditionalSlotConfig createJackpotExpress() { return {"jackpot_express", "Jackpot Express", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "jackpot"}; }
AdditionalSlotConfig createTheDrop() { return {"the_drop", "The Drop", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "adventure"}; }
AdditionalSlotConfig createCaishenRiches() { return {"caishen", "Caishen Riches", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "asian"}; }
AdditionalSlotConfig createHades() { return {"hades", "Hades", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "mythology"}; }
AdditionalSlotConfig createGiantGrizzly() { return {"grizzly", "Giant Grizzly", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "animal"}; }
AdditionalSlotConfig createMysticWheel() { return {"mystic", "Mystic Wheel", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "mystery"}; }
AdditionalSlotConfig createPirateQueens() { return {"pirate", "Pirate Queens", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "pirate"}; }
AdditionalSlotConfig createDragonReborn() { return {"dragon_reborn", "Dragon Reborn", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "dragon"}; }
AdditionalSlotConfig createLotusWarrior() { return {"lotus", "Lotus Warrior", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "high", "asian"}; }
AdditionalSlotConfig createGemsOfBuddha() { return {"buddha", "Gems of Buddha", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "buddha"}; }
AdditionalSlotConfig createMedallion() { return {"medallion", "Medallion", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "treasure"}; }
AdditionalSlotConfig createCrystalGeminis() { return {"geminis", "Crystal Geminis", "Yggdrasil", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "crystal"}; }

// Quickspin
AdditionalSlotConfig createSakuraWind() { return {"sakura", "Sakura Wind", "Quickspin", 5, 3, 50, 0.10, 100.0, 0.97, "medium", "asian"}; }
AdditionalSlotConfig createDawnOfEgypt() { return {"dawn_egypt", "Dawn of Egypt", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "egypt"}; }
AdditionalSlotConfig createBlueOcean() { return {"blue_ocean", "Blue Ocean", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "low", "ocean"}; }
AdditionalSlotConfig createHiddenCity() { return {"hidden_city", "Hidden City", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "adventure"}; }
AdditionalSlotConfig createGrandSpinnSuperChip() { return {"grand_spinn", "Grand Spinn Super Chip", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.97, "high", "classic"}; }
AdditionalSlotConfig createMightyArthur() { return {"arthur", "Mighty Arthur", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "fantasy"}; }
AdditionalSlotConfig createJokerStoker() { return {"joker_stoker", "Joker Stoker", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "horror"}; }
AdditionalSlotConfig createRiskyRabbit() { return {"risky_rabbit", "Risky Rabbit", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "animal"}; }
AdditionalSlotConfig createGoldilocks() { return {"goldilocks", "Goldilocks", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "fairytale"}; }
AdditionalSlotConfig createSpinionsBeachParty() { return {"spinions", "Spinions Beach Party", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "party"}; }
AdditionalSlotConfig createTheEpicQuest() { return {"epic_quest", "The Epic Quest", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "adventure"}; }
AdditionalSlotConfig createAliBaba() { return {"ali_baba", "Ali Baba", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "adventure"}; }
AdditionalSlotConfig createTheLastKingOfAsgard() { return {"asgard_king", "The Last King of Asgard", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "mythology"}; }
AdditionalSlotConfig createFlyingDutchman() { return {"flying_dutchman", "Flying Dutchman", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "pirate"}; }
AdditionalSlotConfig createVampireire() { return {"vampire", "Vampireire", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "horror"}; }
AdditionalSlotConfig createKingColossus() { return {"king_colossus", "King Colossus", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "fantasy"}; }
AdditionalSlotConfig createEl迭() { return {"el_dorado", "El Dorado", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "adventure"}; }
AdditionalSlotConfig createPersianWonders() { return {"persian", "Persian Wonders", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "persia"}; }
AdditionalSlotConfig createDiamondStrike() { return {"diamond_strike", "Diamond Strike", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "high", "gem"}; }
AdditionalSlotConfig createFlipFlip() { return {"flip_flip", "Flip Flip", "Quickspin", 5, 3, 25, 0.10, 100.0, 0.96, "medium", "fruit"}; }

// Playtech
AdditionalSlotConfig createAgeOfGods() { return {"age_of_gods", "Age of Gods", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "mythology"}; }
AdditionalSlotConfig createGladiator() { return {"gladiator", "Gladiator", "Playtech", 5, 3, 25, 0.25, 250.0, 0.95, "high", "historical"}; }
AdditionalSlotConfig createJackpotGiant() { return {"jackpot_giant", "Jackpot Giant", "Playtech", 5, 3, 50, 0.50, 500.0, 0.94, "high", "jackpot"}; }
AdditionalSlotConfig createMegaJackpots() { return {"mega_jackpots", "Mega Jackpots", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "jackpot"}; }
AdditionalSlotConfig createBeachLife() { return {"beach_life", "Beach Life", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "beach"}; }
AdditionalSlotConfig createLifeOfTheParty() { return {"life_party", "Life of the Party", "Playtech", 5, 3, 25, 0.20, 200.0, 0.96, "high", "party"}; }
AdditionalSlotConfig createGoldRush() { return {"gold_rush", "Gold Rush", "Playtech", 5, 3, 25, 0.25, 250.0, 0.95, "high", "western"}; }
AdditionalSlotConfig createSuperHeroes() { return {"super_heroes", "Super Heroes", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "superhero"}; }
AdditionalSlotConfig createKingsTreasure() { return {"kings_treasure", "King's Treasure", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "kingdom"}; }
AdditionalSlotConfig createPrinceOfOlympus() { return {"prince_olympus", "Prince of Olympus", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "mythology"}; }
AdditionalSlotConfig createThaiParadise() { return {"thai_paradise", "Thai Paradise", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "thailand"}; }
AdditionalSlotConfig createKingOfCards() { return {"king_cards", "King of Cards", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "cards"}; }
AdditionalSlotConfig createWhiteOrchid() { return {"white_orchid", "White Orchid", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "nature"}; }
AdditionalSlotConfig createSpartans() { return {"spartans", "Spartans", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "historical"}; }
AdditionalSlotConfig createStacksOfCash() { return {"stacks_cash", "Stacks of Cash", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "money"}; }
AdditionalSlotConfig createWildGems() { return {"wild_gems", "Wild Gems", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "gem"}; }
AdditionalSlotConfig createSuperChip() { return {"super_chip", "Super Chip", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "tech"}; }
AdditionalSlotConfig createPantherMoon() { return {"panther_moon", "Panther Moon", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "animal"}; }
AdditionalSlotConfig createJewelThief() { return {"jewel_thief", "Jewel Thief", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "thief"}; }
AdditionalSlotConfig createRocky() { return {"rocky", "Rocky", "Playtech", 5, 3, 25, 0.20, 200.0, 0.95, "high", "movie"}; }

// IGT
AdditionalSlotConfig createCleopatra() { return {"cleopatra", "Cleopatra", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "egypt"}; }
AdditionalSlotConfig createDaVinciDiamonds() { return {"davinci", "Da Vinci Diamonds", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "art"}; }
AdditionalSlotConfig createWheelOfFortune() { return {"wheel_fortune", "Wheel of Fortune", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "classic"}; }
AdditionalSlotConfig createStarTrek() { return {"star_trek", "Star Trek", "IGT", 5, 3, 30, 0.40, 400.0, 0.95, "high", "scifi"}; }
AdditionalSlotConfig createMonopoly() { return {"monopoly", "Monopoly", "IGT", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "boardgame"}; }
AdditionalSlotConfig createKittyGlitter() { return {"kitty_glitter", "Kitty Glitter", "IGT", 5, 3, 30, 0.20, 200.0, 0.95, "medium", "animal"}; }
AdditionalSlotConfig createWolfRun() { return {"wolf_run", "Wolf Run", "IGT", 5, 3, 40, 0.20, 200.0, 0.95, "medium", "animal"}; }
AdditionalSlotConfig createSiberianStorm() { return {"siberian_storm", "Siberian Storm", "IGT", 5, 3, 50, 0.20, 200.0, 0.95, "high", "animal"}; }
AdditionalSlotConfig createEgyptianRiches() { return {"egyptian_riches", "Egyptian Riches", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "egypt"}; }
AdditionalSlotConfig createPixiesOfTheForest() { return {"pixies", "Pixies of the Forest", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "fantasy"}; }
AdditionalSlotConfig createCrystalForest() { return {"crystal_forest", "Crystal Forest", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "nature"}; }
AdditionalSlotConfig createGongXiFaCai() { return {"gong_xi", "Gong Xi Fa Cai", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "asian"}; }
AdditionalSlotConfig createLuckyLarrys() { return {"lucky_larrys", "Lucky Larry's", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "animal"}; }
AdditionalSlotConfig createStarburstXXXtreme() { return {"starburst_xtreme", "Starburst XXXtreme", "IGT", 5, 3, 20, 0.10, 100.0, 0.96, "high", "classic"}; }
AdditionalSlotConfig createDoubleDiamond() { return {"double_diamond", "Double Diamond", "IGT", 3, 3, 5, 0.05, 50.0, 0.95, "low", "classic"}; }
AdditionalSlotConfig createTripleDiamond() { return {"triple_diamond", "Triple Diamond", "IGT", 3, 3, 5, 0.05, 50.0, 0.95, "low", "classic"}; }
AdditionalSlotConfig createRedWhiteAndBlue() { return {"red_white_blue", "Red White and Blue", "IGT", 3, 3, 5, 0.05, 50.0, 0.95, "low", "classic"}; }
AdditionalSlotConfig createDoubleCasino() { return {"double_casino", "Double Casino", "IGT", 3, 3, 5, 0.05, 50.0, 0.95, "low", "classic"}; }
AdditionalSlotConfig createMegaCrown() { return {"mega_crown", "Mega Crown", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "high", "jackpot"}; }
AdditionalSlotConfig createGoldWins() { return {"gold_wins", "Gold Wins", "IGT", 5, 3, 20, 0.20, 200.0, 0.95, "medium", "gold"}; }

// Ainsworth
AdditionalSlotConfig createThunderCash() { return {"thunder_cash", "Thunder Cash", "Ainsworth", 5, 3, 25, 0.30, 300.0, 0.95, "high", "progressive"}; }
AdditionalSlotConfig createJungleSpirit() { return {"jungle_spirit", "Jungle Spirit", "Ainsworth", 5, 3, 25, 0.40, 400.0, 0.95, "high", "nature"}; }
AdditionalSlotConfig createSunStrike() { return {"sun_strike", "Sun Strike", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "high", "jackpot"}; }
AdditionalSlotConfig createPureFixedOdds() { return {"pure_fixed", "Pure Fixed Odds", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.96, "medium", "classic"}; }
AdditionalSlotConfig createSuperRedPhoenix() { return {"super_red", "Super Red Phoenix", "Ainsworth", 5, 3, 25, 0.40, 400.0, 0.95, "high", "dragon"}; }
AdditionalSlotConfig createWinStorm() { return {"win_storm", "Win Storm", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "high", "storm"}; }
AdditionalSlotConfig createJungleJim() { return {"jungle_jim", "Jungle Jim", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "adventure"}; }
AdditionalSlotConfig createKingmaker() { return {"kingmaker", "Kingmaker", "Ainsworth", 5, 3, 50, 0.40, 400.0, 0.95, "high", "kingdom"}; }
AdditionalSlotConfig createFlyingHigh() { return {"flying_high", "Flying High", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "bird"}; }
AdditionalSlotConfig createMagicMonkey() { return {"magic_monkey", "Magic Monkey", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "animal"}; }
AdditionalSlotConfig createDiamondDuke() { return {"diamond_duke", "Diamond Duke", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "high", "gem"}; }
AdditionalSlotConfig createRoaringForties() { return {"roaring_forties", "Roaring Forties", "Ainsworth", 5, 3, 40, 0.40, 400.0, 0.95, "low", "classic"}; }
AdditionalSlotConfig createJokerJackpots() { return {"joker_jackpots", "Joker Jackpots", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.95, "high", "jackpot"}; }
AdditionalSlotConfig createGrandGiant() { return {"grand_giant", "Grand Giant", "Ainsworth", 5, 3, 25, 0.40, 400.0, 0.95, "high", "giant"}; }
AdditionalSlotConfig createGoldenWolf() { return {"golden_wolf", "Golden Wolf", "Ainsworth", 5, 3, 25, 0.30, 300.0, 0.95, "high", "animal"}; }
AdditionalSlotConfig createApolloRising() { return {"apollo_rising", "Apollo Rising", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.95, "high", "mythology"}; }
AdditionalSlotConfig createThunderCashLightning() { return {"thunder_lightning", "Thunder Cash Lightning", "Ainsworth", 5, 3, 25, 0.30, 300.0, 0.95, "high", "storm"}; }
AdditionalSlotConfig createChinaRiver() { return {"china_river", "China River", "Ainsworth", 5, 3, 25, 0.25, 250.0, 0.95, "medium", "asian"}; }
AdditionalSlotConfig createWildFrog() { return {"wild_frog", "Wild Frog", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "animal"}; }
AdditionalSlotConfig createJadePower() { return {"jade_power", "Jade Power", "Ainsworth", 5, 3, 25, 0.20, 200.0, 0.95, "medium", "asian"}; }

} // namespace AdditionalSlotsBatch2
} // namespace TigerCasino
