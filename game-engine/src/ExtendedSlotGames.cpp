#include "ExtendedSlotGames.hpp"
#include <algorithm>

namespace TigerCasino {

ExtendedSlotGameServer::ExtendedSlotGameServer() {
    initializeAllGames();
}

void ExtendedSlotGameServer::initializeAllGames() {
    using namespace SlotGameTemplates;
    
    // Pragmatic Play style (20 games)
    gameTemplates_["wolf_gold"] = createWolfGold();
    gameTemplates_["great_rhino"] = createGreatRhino();
    gameTemplates_["mustang_gold"] = createMustangGold();
    gameTemplates_["fruit_party"] = createFruitParty();
    gameTemplates_["starlight_princess"] = createStarlightPrincess();
    gameTemplates_["gates_of_olympus"] = createGatesOfOlympus();
    gameTemplates_["the_dog_house"] = createTheDogHouse();
    gameTemplates_["sweet_bonanza"] = createSweetBonanza();
    gameTemplates_["big_bass_bonanza"] = createBigBassBonanza();
    gameTemplates_["john_hunter"] = createJohnHunter();
    gameTemplates_["madame_destin"] = createMadameDestiny();
    gameTemplates_["wild_west_gold"] = createWildWestGold();
    gameTemplates_["aztec_gems"] = createAztecGems();
    gameTemplates_["pyramid_king"] = createPyramidKing();
    gameTemplates_["golden_ox"] = createGoldenOx();
    gameTemplates_["mysterious"] = createMysterious();
    gameTemplates_["power_of_thor"] = createPowerOfThor();
    gameTemplates_["buffalo_king"] = createBuffaloKing();
    gameTemplates_["dragon_hot"] = createDragonHot();
    gameTemplates_["joker_jewels"] = createJokerJewels();
    
    // NetEnt style (10 games)
    gameTemplates_["starburst"] = createStarburst();
    gameTemplates_["gonzo_quest"] = createGonzoQuest();
    gameTemplates_["dead_or_alive"] = createDeadOrAlive();
    gameTemplates_["twin_spin"] = createTwinSpin();
    gameTemplates_["hall_of_gods"] = createHallOfGods();
    gameTemplates_["mega_fortune"] = createMegaFortune();
    gameTemplates_["blood_suckers"] = createBloodSuckers();
    gameTemplates_["jack_and_beanstalk"] = createJackAndTheBeanstalk();
    gameTemplates_["steam_tower"] = createSteamTower();
    gameTemplates_["flowers"] = createFlowers();
    
    // Play'n GO (10 games)
    gameTemplates_["book_of_dead"] = createBookOfDead();
    gameTemplates_["reactoonz"] = createReactoonz();
    gameTemplates_["legacy_of_dead"] = createLegacyOfDead();
    gameTemplates_["rise_of_merlin"] = createRiseOfMerlin();
    gameTemplates_["moon_princess"] = createMoonPrincess();
    gameTemplates_["play_in_hell"] = createPlayInHell();
    gameTemplates_["fire_joker"] = createFireJoker();
    gameTemplates_["ring_of_odysseus"] = createRingOfOdysseus();
    gameTemplates_["tome_of_madness"] = createTomeOfMadness();
    gameTemplates_["honey_rush"] = createHoneyRush();
    
    // Megaways (10 games)
    gameTemplates_["bonanza"] = createBonanza();
    gameTemplates_["extra_chilli"] = createExtraChilli();
    gameTemplates_["white_rabbit"] = createWhiteRabbit();
    gameTemplates_["holy_diver"] = createHolyDiver();
    gameTemplates_["millionaire"] = createWhoWantsToBeAMillionaire();
    gameTemplates_["buffalo_lightning"] = createBuffaloLightning();
    gameTemplates_["diamond_sun"] = createDiamondSun();
    gameTemplates_["monsoon"] = createMonsoon();
    gameTemplates_["danger_voltage"] = createDangerHighVoltage();
    gameTemplates_["extra_chilli_2"] = createExtraChilli();
    
    // Relax/Push Gaming (10 games)
    gameTemplates_["money_train"] = createMoneyTrain();
    gameTemplates_["money_train_2"] = createMoneyTrain2();
    gameTemplates_["space_miners"] = createSpaceMiners();
    gameTemplates_["temple_tumble"] = createTempleTumble();
    gameTemplates_["dropz"] = createDropz();
    gameTemplates_["mega_stack"] = createMegaStack();
    gameTemplates_["push_up"] = createPushUp();
    gameTemplates_["razor_sharks"] = createRazorSharks();
    gameTemplates_["sticky_birds"] = createStickyBirds();
    gameTemplates_["wild_flowers"] = createWildFlowers();
    
    // Hacksaw Gaming (10 games)
    gameTemplates_["wanted_dead"] = createWantedDeadOrAlive();
    gameTemplates_["stack_em"] = createStackEm();
    gameTemplates_["the_bounty"] = createTheBounty();
    gameTemplates_["dice"] = createDice();
    gameTemplates_["aurora"] = createAurora();
    gameTemplates_["chaos_crew"] = createChaosCrew();
    gameTemplates_["cabin_crashers"] = createCabinCrashers();
    gameTemplates_["beast_mode"] = createBeastMode();
    gameTemplates_["outsmart"] = createOutsmart();
    gameTemplates_["the_emperor"] = createTheEmperor();
    
    // Nolimit City (10 games)
    gameTemplates_["san_quentin"] = createSanQuentin();
    gameTemplates_["book_of_shadows"] = createBookOfShadows();
    gameTemplates_["mental"] = createMental();
    gameTemplates_["xxxtreme"] = createXxxtreme();
    gameTemplates_["dead_panda"] = createDeadPanda();
    gameTemplates_["fire_in_hole"] = createFireInTheHole();
    gameTemplates_["kiss"] = createKiss();
    gameTemplates_["larry_leprechaun"] = createLarryTheLeprechaun();
    gameTemplates_["million_777"] = createMillion777();
    gameTemplates_["pontius_pilate"] = createPontiusPilate();
}

ExtendedSlotConfig ExtendedSlotGameServer::getGameConfig(const std::string& gameId) const {
    auto it = gameTemplates_.find(gameId);
    if (it != gameTemplates_.end()) {
        return it->second;
    }
    return ExtendedSlotConfig{};
}

std::vector<ExtendedSlotConfig> ExtendedSlotGameServer::getGamesByProvider(const std::string& provider) const {
    std::vector<ExtendedSlotConfig> result;
    for (const auto& pair : gameTemplates_) {
        if (pair.second.provider == provider) {
            result.push_back(pair.second);
        }
    }
    return result;
}

std::vector<ExtendedSlotConfig> ExtendedSlotGameServer::getGamesByTheme(const std::string& theme) const {
    std::vector<ExtendedSlotConfig> result;
    for (const auto& pair : gameTemplates_) {
        if (pair.second.theme == theme) {
            result.push_back(pair.second);
        }
    }
    return result;
}

std::vector<ExtendedSlotConfig> ExtendedSlotGameServer::getGamesByVolatility(const std::string& volatility) const {
    std::vector<ExtendedSlotConfig> result;
    for (const auto& pair : gameTemplates_) {
        if (pair.second.volatility == volatility) {
            result.push_back(pair.second);
        }
    }
    return result;
}

std::vector<ExtendedSlotConfig> ExtendedSlotGameServer::getAllGames() const {
    std::vector<ExtendedSlotConfig> result;
    for (const auto& pair : gameTemplates_) {
        result.push_back(pair.second);
    }
    return result;
}

size_t ExtendedSlotGameServer::getGameCount() const {
    return gameTemplates_.size();
}

std::map<std::string, int> ExtendedSlotGameServer::getProviderGameCounts() const {
    std::map<std::string, int> counts;
    for (const auto& pair : gameTemplates_) {
        counts[pair.second.provider]++;
    }
    return counts;
}

// Template implementations
namespace SlotGameTemplates {

ExtendedSlotConfig createWolfGold() {
    ExtendedSlotConfig config;
    config.gameId = "wolf_gold";
    config.gameName = "Wolf Gold";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.25;
    config.maxBet = 125.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = true;
    config.volatility = "medium";
    config.theme = "nature";
    config.symbols = {"WOLF", "BISON", "EAGLE", "HORSE", "A", "K", "Q", "J"};
    config.symbolValues = {
        {"WOLF", 20.0}, {"BISON", 10.0}, {"EAGLE", 5.0}, {"HORSE", 2.5},
        {"A", 1.5}, {"K", 1.2}, {"Q", 1.0}, {"J", 0.8}
    };
    config.symbolWeights = {
        {"WOLF", 2}, {"BISON", 3}, {"EAGLE", 4}, {"HORSE", 5},
        {"A", 10}, {"K", 12}, {"Q", 14}, {"J", 16}
    };
    return config;
}

ExtendedSlotConfig createGreatRhino() {
    ExtendedSlotConfig config;
    config.gameId = "great_rhino";
    config.gameName = "Great Rhino";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.965;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = true;
    config.volatility = "high";
    config.theme = "nature";
    config.symbols = {"RHINO", "ELEPHANT", "LION", "ZEBRA", "A", "K", "Q", "J"};
    config.symbolValues = {
        {"RHINO", 20.0}, {"ELEPHANT", 10.0}, {"LION", 5.0}, {"ZEBRA", 2.5},
        {"A", 1.5}, {"K", 1.2}, {"Q", 1.0}, {"J", 0.8}
    };
    config.symbolWeights = {{"RHINO", 2}, {"ELEPHANT", 3}, {"LION", 4}, {"ZEBRA", 5}, {"A", 10}, {"K", 12}, {"Q", 14}, {"J", 16}};
    return config;
}

ExtendedSlotConfig createMustangGold() {
    ExtendedSlotConfig config;
    config.gameId = "mustang_gold";
    config.gameName = "Mustang Gold";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.25;
    config.maxBet = 150.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = true;
    config.volatility = "high";
    config.theme = "western";
    config.symbols = {"HORSE", "SHERIFF", "BOOT", "HAT", "A", "K", "Q", "J"};
    config.symbolValues = {
        {"HORSE", 20.0}, {"SHERIFF", 10.0}, {"BOOT", 5.0}, {"HAT", 2.5},
        {"A", 1.5}, {"K", 1.2}, {"Q", 1.0}, {"J", 0.8}
    };
    config.symbolWeights = {{"HORSE", 2}, {"SHERIFF", 3}, {"BOOT", 4}, {"HAT", 5}, {"A", 10}, {"K", 12}, {"Q", 14}, {"J", 16}};
    return config;
}

ExtendedSlotConfig createFruitParty() {
    ExtendedSlotConfig config;
    config.gameId = "fruit_party";
    config.gameName = "Fruit Party";
    config.provider = "PragmaticPlay";
    config.reels = 7;
    config.rows = 7;
    config.paylines = 0; // Cluster pays
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.965;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "fruit";
    config.symbols = {"FRUIT", "GRAPE", "ORANGE", "APPLE", "PLUM", "CHERRY"};
    config.symbolValues = {
        {"FRUIT", 20.0}, {"GRAPE", 10.0}, {"ORANGE", 5.0}, {"APPLE", 3.0},
        {"PLUM", 2.0}, {"CHERRY", 1.5}
    };
    config.symbolWeights = {{"FRUIT", 2}, {"GRAPE", 4}, {"ORANGE", 6}, {"APPLE", 8}, {"PLUM", 10}, {"CHERRY", 12}};
    return config;
}

ExtendedSlotConfig createStarlightPrincess() {
    ExtendedSlotConfig config;
    config.gameId = "starlight_princess";
    config.gameName = "Starlight Princess";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 200.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "fantasy";
    config.symbols = {"PRINCESS", "STAR", "MOON", "HEART", "DIAMOND"};
    config.symbolValues = {
        {"PRINCESS", 50.0}, {"STAR", 20.0}, {"MOON", 10.0}, {"HEART", 5.0}, {"DIAMOND", 3.0}
    };
    config.symbolWeights = {{"PRINCESS", 1}, {"STAR", 2}, {"MOON", 4}, {"HEART", 6}, {"DIAMOND", 10}};
    return config;
}

ExtendedSlotConfig createGatesOfOlympus() {
    ExtendedSlotConfig config;
    config.gameId = "gates_of_olympus";
    config.gameName = "Gates of Olympus";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 200.0;
    config.rtp = 0.965;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "mythology";
    config.symbols = {"ZEUS", "HERMES", "POSEIDON", "HERA", "TRIDENT", "RING"};
    config.symbolValues = {
        {"ZEUS", 50.0}, {"HERMES", 20.0}, {"POSEIDON", 10.0}, {"HERA", 5.0}, {"TRIDENT", 3.0}, {"RING", 2.0}
    };
    config.symbolWeights = {{"ZEUS", 1}, {"HERMES", 2}, {"POSEIDON", 3}, {"HERA", 5}, {"TRIDENT", 8}, {"RING", 10}};
    return config;
}

ExtendedSlotConfig createTheDogHouse() {
    ExtendedSlotConfig config;
    config.gameId = "the_dog_house";
    config.gameName = "The Dog House";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 20;
    config.minBet = 0.20;
    config.maxBet = 200.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "animals";
    config.symbols = {"DOG", "BONE", "COLLAR", "HOUSE", "A", "K", "Q", "J"};
    config.symbolValues = {
        {"DOG", 15.0}, {"BONE", 5.0}, {"COLLAR", 3.0}, {"HOUSE", 2.5},
        {"A", 1.5}, {"K", 1.2}, {"Q", 1.0}, {"J", 0.8}
    };
    config.symbolWeights = {{"DOG", 2}, {"BONE", 4}, {"COLLAR", 6}, {"HOUSE", 8}, {"A", 10}, {"K", 12}, {"Q", 14}, {"J", 16}};
    return config;
}

ExtendedSlotConfig createSweetBonanza() {
    ExtendedSlotConfig config;
    config.gameId = "sweet_bonanza";
    config.gameName = "Sweet Bonanza";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 200.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "candy";
    config.symbols = {"CANDY", "BANANA", "GRAPE", "APPLE", "BERRY", "HEART"};
    config.symbolValues = {
        {"CANDY", 20.0}, {"BANANA", 10.0}, {"GRAPE", 5.0}, {"APPLE", 3.0}, {"BERRY", 2.0}, {"HEART", 1.5}
    };
    config.symbolWeights = {{"CANDY", 2}, {"BANANA", 4}, {"GRAPE", 6}, {"APPLE", 8}, {"BERRY", 10}, {"HEART", 12}};
    return config;
}

ExtendedSlotConfig createBigBassBonanza() {
    ExtendedSlotConfig config;
    config.gameId = "big_bass_bonanza";
    config.gameName = "Big Bass Bonanza";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 250.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "fishing";
    config.symbols = {"FISHER", "FISH", "DRAGONFLY", "HOOK", "A", "K"};
    config.symbolValues = {
        {"FISHER", 15.0}, {"FISH", 10.0}, {"DRAGONFLY", 5.0}, {"HOOK", 2.5}, {"A", 1.5}, {"K", 1.2}
    };
    config.symbolWeights = {{"FISHER", 2}, {"FISH", 4}, {"DRAGONFLY", 6}, {"HOOK", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createJohnHunter() {
    ExtendedSlotConfig config;
    config.gameId = "john_hunter";
    config.gameName = "John Hunter and the Tomb of Scarab";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.10;
    config.maxBet = 187.5;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "adventure";
    config.symbols = {"JOHN", "SCARAB", "MAP", "COMPASS", "ANKH"};
    config.symbolValues = {
        {"JOHN", 20.0}, {"SCARAB", 10.0}, {"MAP", 5.0}, {"COMPASS", 3.0}, {"ANKH", 2.0}
    };
    config.symbolWeights = {{"JOHN", 2}, {"SCARAB", 4}, {"MAP", 6}, {"COMPASS", 8}, {"ANKH", 10}};
    return config;
}

// Additional template implementations (abbreviated for brevity)
ExtendedSlotConfig createMadameDestiny() {
    ExtendedSlotConfig config;
    config.gameId = "madame_destin";
    config.gameName = "Madame Destiny";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.20;
    config.maxBet = 125.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "mystical";
    config.symbols = {"MADAME", "CRYSTAL", "OWL", "CANDLE", "A", "K"};
    config.symbolValues = {{"MADAME", 20.0}, {"CRYSTAL", 10.0}, {"OWL", 5.0}, {"CANDLE", 2.5}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"MADAME", 2}, {"CRYSTAL", 4}, {"OWL", 6}, {"CANDLE", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createWildWestGold() {
    ExtendedSlotConfig config;
    config.gameId = "wild_west_gold";
    config.gameName = "Wild West Gold";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.25;
    config.maxBet = 125.0;
    config.rtp = 0.965;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "western";
    config.symbols = {"SHERIFF", "COWBOY", "HORSE", "GUN", "A", "K"};
    config.symbolValues = {{"SHERIFF", 20.0}, {"COWBOY", 10.0}, {"HORSE", 5.0}, {"GUN", 2.5}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"SHERIFF", 2}, {"COWBOY", 4}, {"HORSE", 6}, {"GUN", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createAztecGems() {
    ExtendedSlotConfig config;
    config.gameId = "aztec_gems";
    config.gameName = "Aztec Gems";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 5;
    config.minBet = 0.05;
    config.maxBet = 25.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = true;
    config.volatility = "high";
    config.theme = "aztec";
    config.symbols = {"GEM", "AZTEC", "SKULL", "JAGUAR"};
    config.symbolValues = {{"GEM", 20.0}, {"AZTEC", 10.0}, {"SKULL", 5.0}, {"JAGUAR", 2.5}};
    config.symbolWeights = {{"GEM", 2}, {"AZTEC", 4}, {"SKULL", 6}, {"JAGUAR", 10}};
    return config;
}

ExtendedSlotConfig createPyramidKing() {
    ExtendedSlotConfig config;
    config.gameId = "pyramid_king";
    config.gameName = "Pyramid King";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.50;
    config.maxBet = 250.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "egypt";
    config.symbols = {"KING", "QUEEN", "ANUBIS", "SCARAB", "A", "K"};
    config.symbolValues = {{"KING", 20.0}, {"QUEEN", 10.0}, {"ANUBIS", 5.0}, {"SCARAB", 2.5}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"KING", 2}, {"QUEEN", 4}, {"ANUBIS", 6}, {"SCARAB", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createGoldenOx() {
    ExtendedSlotConfig config;
    config.gameId = "golden_ox";
    config.gameName = "Golden Ox";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.60;
    config.maxBet = 300.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "asian";
    config.symbols = {"OX", "LAMP", "COIN", "BELL", "A", "K"};
    config.symbolValues = {{"OX", 25.0}, {"LAMP", 12.0}, {"COIN", 6.0}, {"BELL", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"OX", 2}, {"LAMP", 4}, {"COIN", 6}, {"BELL", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createMysterious() {
    ExtendedSlotConfig config;
    config.gameId = "mysterious";
    config.gameName = "Mysterious";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 240.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "mystery";
    config.symbols = {"DETECTIVE", "CLUE", "GUN", "MAGNIFY", "A", "K"};
    config.symbolValues = {{"DETECTIVE", 25.0}, {"CLUE", 12.0}, {"GUN", 6.0}, {"MAGNIFY", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"DETECTIVE", 2}, {"CLUE", 4}, {"GUN", 6}, {"MAGNIFY", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createPowerOfThor() {
    ExtendedSlotConfig config;
    config.gameId = "power_of_thor";
    config.gameName = "Power of Thor";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 200.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "mythology";
    config.symbols = {"THOR", "HAMMER", "LIGHTNING", "SHIELD", "A", "K"};
    config.symbolValues = {{"THOR", 25.0}, {"HAMMER", 12.0}, {"LIGHTNING", 6.0}, {"SHIELD", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"THOR", 2}, {"HAMMER", 4}, {"LIGHTNING", 6}, {"SHIELD", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createBuffaloKing() {
    ExtendedSlotConfig config;
    config.gameId = "buffalo_king";
    config.gameName = "Buffalo King";
    config.provider = "PragmaticPlay";
    config.reels = 6;
    config.rows = 4;
    config.paylines = 0;
    config.minBet = 0.40;
    config.maxBet = 200.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "nature";
    config.symbols = {"BUFFALO", "WOLF", "EAGLE", "LION", "A", "K"};
    config.symbolValues = {{"BUFFALO", 25.0}, {"WOLF", 12.0}, {"EAGLE", 6.0}, {"LION", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"BUFFALO", 2}, {"WOLF", 4}, {"EAGLE", 6}, {"LION", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createDragonHot() {
    ExtendedSlotConfig config;
    config.gameId = "dragon_hot";
    config.gameName = "Dragon Hot";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 5;
    config.minBet = 0.05;
    config.maxBet = 25.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "dragon";
    config.symbols = {"DRAGON", "FIRE", "COIN", "PEARL"};
    config.symbolValues = {{"DRAGON", 20.0}, {"FIRE", 10.0}, {"COIN", 5.0}, {"PEARL", 2.5}};
    config.symbolWeights = {{"DRAGON", 2}, {"FIRE", 4}, {"COIN", 8}, {"PEARL", 12}};
    return config;
}

ExtendedSlotConfig createJokerJewels() {
    ExtendedSlotConfig config;
    config.gameId = "joker_jewels";
    config.gameName = "Joker Jewels";
    config.provider = "PragmaticPlay";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.05;
    config.maxBet = 250.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "classic";
    config.symbols = {"JOKER", "CROWN", "GEM", "BELL", "A", "K"};
    config.symbolValues = {{"JOKER", 25.0}, {"CROWN", 12.0}, {"GEM", 6.0}, {"BELL", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"JOKER", 2}, {"CROWN", 4}, {"GEM", 6}, {"BELL", 8}, {"A", 10}, {"K", 12}};
    return config;
}

// NetEnt style games (abbreviated implementations)
ExtendedSlotConfig createStarburst() {
    ExtendedSlotConfig config;
    config.gameId = "starburst";
    config.gameName = "Starburst";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 1.0;
    config.maxBet = 1000.0;
    config.rtp = 0.96;
    config.hasBonus = false;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "low";
    config.theme = "space";
    config.symbols = {"STAR", "BAR", "7", "BELL", "CHERRY", "GEM"};
    config.symbolValues = {{"STAR", 10.0}, {"BAR", 5.0}, {"7", 4.0}, {"BELL", 3.0}, {"CHERRY", 2.0}, {"GEM", 1.5}};
    config.symbolWeights = {{"STAR", 2}, {"BAR", 4}, {"7", 4}, {"BELL", 6}, {"CHERRY", 10}, {"GEM", 12}};
    return config;
}

ExtendedSlotConfig createGonzoQuest() {
    ExtendedSlotConfig config;
    config.gameId = "gonzo_quest";
    config.gameName = "Gonzo's Quest";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 20;
    config.minBet = 0.20;
    config.maxBet = 50.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "adventure";
    config.symbols = {"GONZO", "JAGUAR", "FROG", "FISH", "A", "K"};
    config.symbolValues = {{"GONZO", 15.0}, {"JAGUAR", 8.0}, {"FROG", 5.0}, {"FISH", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"GONZO", 2}, {"JAGUAR", 4}, {"FROG", 6}, {"FISH", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createDeadOrAlive() {
    ExtendedSlotConfig config;
    config.gameId = "dead_or_alive";
    config.gameName = "Dead or Alive";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 12;
    config.minBet = 0.09;
    config.maxBet = 18.0;
    config.rtp = 0.97;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "western";
    config.symbols = {"WANTED", "GUN", "BOOT", "HAT", "A", "K"};
    config.symbolValues = {{"WANTED", 15.0}, {"GUN", 8.0}, {"BOOT", 5.0}, {"HAT", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"WANTED", 2}, {"GUN", 4}, {"BOOT", 6}, {"HAT", 8}, {"A", 10}, {"K", 12}};
    return config;
}

// Placeholder implementations for remaining games
ExtendedSlotConfig createTwinSpin() {
    ExtendedSlotConfig config;
    config.gameId = "twin_spin";
    config.gameName = "Twin Spin";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 243;
    config.minBet = 0.25;
    config.maxBet = 250.0;
    config.rtp = 0.96;
    config.hasBonus = false;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "classic";
    config.symbols = {"DIAMOND", "7", "BAR", "CHERRY", "BELL", "A"};
    config.symbolValues = {{"DIAMOND", 15.0}, {"7", 10.0}, {"BAR", 5.0}, {"CHERRY", 3.0}, {"BELL", 2.0}, {"A", 1.5}};
    config.symbolWeights = {{"DIAMOND", 2}, {"7", 4}, {"BAR", 6}, {"CHERRY", 8}, {"BELL", 10}, {"A", 12}};
    return config;
}

ExtendedSlotConfig createHallOfGods() {
    ExtendedSlotConfig config;
    config.gameId = "hall_of_gods";
    config.gameName = "Hall of Gods";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 20;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.95;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = true;
    config.volatility = "high";
    config.theme = "mythology";
    config.symbols = {"ODIN", "THOR", "LOKI", "VALKYRIE", "A", "K"};
    config.symbolValues = {{"ODIN", 20.0}, {"THOR", 12.0}, {"LOKI", 8.0}, {"VALKYRIE", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"ODIN", 1}, {"THOR", 2}, {"LOKI", 3}, {"VALKYRIE", 5}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createMegaFortune() {
    ExtendedSlotConfig config;
    config.gameId = "mega_fortune";
    config.gameName = "Mega Fortune";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.25;
    config.maxBet = 125.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = true;
    config.volatility = "high";
    config.theme = "luxury";
    config.symbols = {"YACHT", "RING", "CHAMPAGNE", "AUTO", "A", "K"};
    config.symbolValues = {{"YACHT", 20.0}, {"RING", 12.0}, {"CHAMPAGNE", 8.0}, {"AUTO", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"YACHT", 1}, {"RING", 2}, {"CHAMPAGNE", 3}, {"AUTO", 5}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createBloodSuckers() {
    ExtendedSlotConfig config;
    config.gameId = "blood_suckers";
    config.gameName = "Blood Suckers";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 25;
    config.minBet = 0.25;
    config.maxBet = 62.5;
    config.rtp = 0.98;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "low";
    config.theme = "horror";
    config.symbols = {"VAMPIRE", "COFFIN", "BAT", "GARLIC", "A", "K"};
    config.symbolValues = {{"VAMPIRE", 15.0}, {"COFFIN", 10.0}, {"BAT", 6.0}, {"GARLIC", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"VAMPIRE", 2}, {"COFFIN", 4}, {"BAT", 6}, {"GARLIC", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createJackAndTheBeanstalk() {
    ExtendedSlotConfig config;
    config.gameId = "jack_and_beanstalk";
    config.gameName = "Jack and the Beanstalk";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 20;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "fantasy";
    config.symbols = {"JACK", "BEANSTALK", "GOLDEN", "COW", "A", "K"};
    config.symbolValues = {{"JACK", 15.0}, {"BEANSTALK", 10.0}, {"GOLDEN", 6.0}, {"COW", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"JACK", 2}, {"BEANSTALK", 4}, {"GOLDEN", 6}, {"COW", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createSteamTower() {
    ExtendedSlotConfig config;
    config.gameId = "steam_tower";
    config.gameName = "Steam Tower";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 15;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.97;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "steampunk";
    config.symbols = {"GIRL", "MECH", "GEAR", "WATCH", "A", "K"};
    config.symbolValues = {{"GIRL", 15.0}, {"MECH", 10.0}, {"GEAR", 6.0}, {"WATCH", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"GIRL", 2}, {"MECH", 4}, {"GEAR", 6}, {"WATCH", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createFlowers() {
    ExtendedSlotConfig config;
    config.gameId = "flowers";
    config.gameName = "Flowers";
    config.provider = "NetEnt";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 30;
    config.minBet = 0.20;
    config.maxBet = 75.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "nature";
    config.symbols = {"FLOWER", "SUN", "TULIP", "ROSE", "A", "K"};
    config.symbolValues = {{"FLOWER", 15.0}, {"SUN", 10.0}, {"TULIP", 6.0}, {"ROSE", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"FLOWER", 2}, {"SUN", 4}, {"TULIP", 6}, {"ROSE", 8}, {"A", 10}, {"K", 12}};
    return config;
}

// Play'n GO games
ExtendedSlotConfig createBookOfDead() {
    ExtendedSlotConfig config;
    config.gameId = "book_of_dead";
    config.gameName = "Book of Dead";
    config.provider = "PlaynGO";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "egypt";
    config.symbols = {"BOOK", "PHARAOH", "ANKH", "SCARAB", "A", "K"};
    config.symbolValues = {{"BOOK", 20.0}, {"PHARAOH", 12.0}, {"ANKH", 8.0}, {"SCARAB", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"BOOK", 2}, {"PHARAOH", 3}, {"ANKH", 4}, {"SCARAB", 6}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createReactoonz() {
    ExtendedSlotConfig config;
    config.gameId = "reactoonz";
    config.gameName = "Reactoonz";
    config.provider = "PlaynGO";
    config.reels = 7;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "space";
    config.symbols = {"ALIEN", "ONE", "TWO", "THREE", "FOUR"};
    config.symbolValues = {{"ALIEN", 20.0}, {"ONE", 10.0}, {"TWO", 6.0}, {"THREE", 4.0}, {"FOUR", 2.0}};
    config.symbolWeights = {{"ALIEN", 2}, {"ONE", 4}, {"TWO", 6}, {"THREE", 8}, {"FOUR", 12}};
    return config;
}

ExtendedSlotConfig createLegacyOfDead() {
    ExtendedSlotConfig config;
    config.gameId = "legacy_of_dead";
    config.gameName = "Legacy of Dead";
    config.provider = "PlaynGO";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "egypt";
    config.symbols = {"BOOK", "PHARAOH", "MUMMY", "ANKH", "A", "K"};
    config.symbolValues = {{"BOOK", 20.0}, {"PHARAOH", 12.0}, {"MUMMY", 8.0}, {"ANKH", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"BOOK", 2}, {"PHARAOH", 3}, {"MUMMY", 4}, {"ANKH", 6}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createRiseOfMerlin() {
    ExtendedSlotConfig config;
    config.gameId = "rise_of_merlin";
    config.gameName = "Rise of Merlin";
    config.provider = "PlaynGO";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "magic";
    config.symbols = {"MERLIN", "OWL", "BOOK", "CRYSTAL", "A", "K"};
    config.symbolValues = {{"MERLIN", 20.0}, {"OWL", 10.0}, {"BOOK", 6.0}, {"CRYSTAL", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"MERLIN", 2}, {"OWL", 4}, {"BOOK", 6}, {"CRYSTAL", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createMoonPrincess() {
    ExtendedSlotConfig config;
    config.gameId = "moon_princess";
    config.gameName = "Moon Princess";
    config.provider = "PlaynGO";
    config.reels = 5;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "anime";
    config.symbols = {"PRINCESS", "MOON", "STAR", "HEART", "A", "K"};
    config.symbolValues = {{"PRINCESS", 20.0}, {"MOON", 10.0}, {"STAR", 6.0}, {"HEART", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"PRINCESS", 2}, {"MOON", 4}, {"STAR", 6}, {"HEART", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createPlayInHell() {
    ExtendedSlotConfig config;
    config.gameId = "play_in_hell";
    config.gameName = "Play in Hell";
    config.provider = "PlaynGO";
    config.reels = 5;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "horror";
    config.symbols = {"DEVIL", "GOD", "TRIDENT", "FIRE", "A", "K"};
    config.symbolValues = {{"DEVIL", 25.0}, {"GOD", 12.0}, {"TRIDENT", 8.0}, {"FIRE", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"DEVIL", 1}, {"GOD", 3}, {"TRIDENT", 5}, {"FIRE", 7}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createFireJoker() {
    ExtendedSlotConfig config;
    config.gameId = "fire_joker";
    config.gameName = "Fire Joker";
    config.provider = "PlaynGO";
    config.reels = 3;
    config.rows = 3;
    config.paylines = 5;
    config.minBet = 0.05;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "classic";
    config.symbols = {"JOKER", "7", "BAR", "BELL", "FRUIT"};
    config.symbolValues = {{"JOKER", 20.0}, {"7", 10.0}, {"BAR", 5.0}, {"BELL", 3.0}, {"FRUIT", 2.0}};
    config.symbolWeights = {{"JOKER", 2}, {"7", 4}, {"BAR", 6}, {"BELL", 10}, {"FRUIT", 14}};
    return config;
}

ExtendedSlotConfig createRingOfOdysseus() {
    ExtendedSlotConfig config;
    config.gameId = "ring_of_odysseus";
    config.gameName = "Ring of Odysseus";
    config.provider = "PlaynGO";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "mythology";
    config.symbols = {"ODYSSEUS", "SHIP", "SWORD", "SHIELD", "A", "K"};
    config.symbolValues = {{"ODYSSEUS", 20.0}, {"SHIP", 10.0}, {"SWORD", 6.0}, {"SHIELD", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"ODYSSEUS", 2}, {"SHIP", 4}, {"SWORD", 6}, {"SHIELD", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createTomeOfMadness() {
    ExtendedSlotConfig config;
    config.gameId = "tome_of_madness";
    config.gameName = "Tome of Madness";
    config.provider = "PlaynGO";
    config.reels = 5;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "horror";
    config.symbols = {"CULTIST", "BOOK", "SKULL", "EYE", "A", "K"};
    config.symbolValues = {{"CULTIST", 20.0}, {"BOOK", 10.0}, {"SKULL", 6.0}, {"EYE", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"CULTIST", 2}, {"BOOK", 4}, {"SKULL", 6}, {"EYE", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createHoneyRush() {
    ExtendedSlotConfig config;
    config.gameId = "honey_rush";
    config.gameName = "Honey Rush";
    config.provider = "PlaynGO";
    config.reels = 7;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "nature";
    config.symbols = {"BEE", "HONEY", "FLOWER", "HIVE", "A", "K"};
    config.symbolValues = {{"BEE", 25.0}, {"HONEY", 12.0}, {"FLOWER", 8.0}, {"HIVE", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"BEE", 1}, {"HONEY", 3}, {"FLOWER", 5}, {"HIVE", 7}, {"A", 10}, {"K", 12}};
    return config;
}

// Megaways implementations (abbreviated)
ExtendedSlotConfig createBonanza() {
    ExtendedSlotConfig config;
    config.gameId = "bonanza";
    config.gameName = "Bonanza";
    config.provider = "BigTimeGaming";
    config.reels = 6;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 500.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "western";
    config.symbols = {"MINER", "GOLD", "LETTER", "TNT", "A", "K"};
    config.symbolValues = {{"MINER", 20.0}, {"GOLD", 10.0}, {"LETTER", 5.0}, {"TNT", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"MINER", 2}, {"GOLD", 4}, {"LETTER", 6}, {"TNT", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createExtraChilli() {
    ExtendedSlotConfig config;
    config.gameId = "extra_chilli";
    config.gameName = "Extra Chilli";
    config.provider = "BigTimeGaming";
    config.reels = 6;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 500.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "mexican";
    config.symbols = {"CHILLI", "MEAT", "HAT", "PEPPER", "A", "K"};
    config.symbolValues = {{"CHILLI", 20.0}, {"MEAT", 10.0}, {"HAT", 5.0}, {"PEPPER", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"CHILLI", 2}, {"MEAT", 4}, {"HAT", 6}, {"PEPPER", 8}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createWhiteRabbit() {
    ExtendedSlotConfig config;
    config.gameId = "white_rabbit";
    config.gameName = "White Rabbit";
    config.provider = "BigTimeGaming";
    config.reels = 6;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 800.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "fantasy";
    config.symbols = {"RABBIT", "CLOCK", "CAKE", "KEY", "A", "K"};
    config.symbolValues = {{"RABBIT", 25.0}, {"CLOCK", 12.0}, {"CAKE", 8.0}, {"KEY", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"RABBIT", 1}, {"CLOCK", 3}, {"CAKE", 5}, {"KEY", 7}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createHolyDiver() {
    ExtendedSlotConfig config;
    config.gameId = "holy_diver";
    config.gameName = "Holy Diver";
    config.provider = "BigTimeGaming";
    config.reels = 6;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.30;
    config.maxBet = 300.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "adventure";
    config.symbols = {"KNIGHT", "SWORD", "DRAGON", "SHIELD", "A", "K"};
    config.symbolValues = {{"KNIGHT", 25.0}, {"SWORD", 12.0}, {"DRAGON", 8.0}, {"SHIELD", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"KNIGHT", 1}, {"SWORD", 3}, {"DRAGON", 5}, {"SHIELD", 7}, {"A", 10}, {"K", 12}};
    return config;
}

ExtendedSlotConfig createWhoWantsToBeAMillionaire() {
    ExtendedSlotConfig config;
    config.gameId = "millionaire";
    config.gameName = "Who Wants to Be a Millionaire";
    config.provider = "BigTimeGaming";
    config.reels = 6;
    config.rows = 7;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 500.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "quiz";
    config.symbols = {"TROPHY", "QUESTION", "LADDER", "PHONE", "A", "K"};
    config.symbolValues = {{"TROPHY", 50.0}, {"QUESTION", 20.0}, {"LADDER", 10.0}, {"PHONE", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"TROPHY", 1}, {"QUESTION", 2}, {"LADDER", 4}, {"PHONE", 6}, {"A", 10}, {"K", 12}};
    return config;
}

// Placeholder functions for remaining games
ExtendedSlotConfig createBuffaloLightning() { return createBuffaloKing(); }
ExtendedSlotConfig createDiamondSun() { return createWolfGold(); }
ExtendedSlotConfig createMonsoon() { return createGreatRhino(); }
ExtendedSlotConfig createDangerHighVoltage() { return createExtraChilli(); }
ExtendedSlotConfig createMoneyTrain() {
    ExtendedSlotConfig config;
    config.gameId = "money_train";
    config.gameName = "Money Train";
    config.provider = "RelaxGaming";
    config.reels = 5;
    config.rows = 4;
    config.paylines = 40;
    config.minBet = 0.40;
    config.maxBet = 40.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "western";
    config.symbols = {"TRAIN", "GUN", "BAG", "MASK", "A", "K"};
    config.symbolValues = {{"TRAIN", 20.0}, {"GUN", 10.0}, {"BAG", 5.0}, {"MASK", 3.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"TRAIN", 2}, {"GUN", 4}, {"BAG", 6}, {"MASK", 8}, {"A", 10}, {"K", 12}};
    return config;
}
ExtendedSlotConfig createMoneyTrain2() { return createMoneyTrain(); }
ExtendedSlotConfig createSpaceMiners() {
    ExtendedSlotConfig config;
    config.gameId = "space_miners";
    config.gameName = "Space Miners";
    config.provider = "RelaxGaming";
    config.reels = 6;
    config.rows = 5;
    config.paylines = 0;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "space";
    config.symbols = {"ALIEN", "ROCKET", "PLANET", "ASTRONAUT", "A", "K"};
    config.symbolValues = {{"ALIEN", 20.0}, {"ROCKET", 10.0}, {"PLANET", 6.0}, {"ASTRONAUT", 4.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"ALIEN", 2}, {"ROCKET", 4}, {"PLANET", 6}, {"ASTRONAUT", 8}, {"A", 10}, {"K", 12}};
    return config;
}
ExtendedSlotConfig createTempleTumble() { return createSpaceMiners(); }
ExtendedSlotConfig createDropz() { return createSpaceMiners(); }
ExtendedSlotConfig createMegaStack() { return createSpaceMiners(); }
ExtendedSlotConfig createPushUp() { return createSpaceMiners(); }
ExtendedSlotConfig createRazorSharks() { return createSpaceMiners(); }
ExtendedSlotConfig createStickyBirds() { return createSpaceMiners(); }
ExtendedSlotConfig createWildFlowers() { return createSpaceMiners(); }

// Hacksaw Gaming implementations
ExtendedSlotConfig createWantedDeadOrAlive() {
    ExtendedSlotConfig config;
    config.gameId = "wanted_dead";
    config.gameName = "Wanted Dead or Alive";
    config.provider = "HacksawGaming";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 15;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "western";
    config.symbols = {"WANTED", "SKULL", "GUN", "BOOT", "A", "K"};
    config.symbolValues = {{"WANTED", 25.0}, {"SKULL", 12.0}, {"GUN", 8.0}, {"BOOT", 5.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"WANTED", 1}, {"SKULL", 3}, {"GUN", 5}, {"BOOT", 7}, {"A", 10}, {"K", 12}};
    return config;
}
ExtendedSlotConfig createStackEm() { return createWantedDeadOrAlive(); }
ExtendedSlotConfig createTheBounty() { return createWantedDeadOrAlive(); }
ExtendedSlotConfig createDice() {
    ExtendedSlotConfig config;
    config.gameId = "dice";
    config.gameName = "Dice";
    config.provider = "HacksawGaming";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 10;
    config.minBet = 0.10;
    config.maxBet = 100.0;
    config.rtp = 0.96;
    config.hasBonus = false;
    config.hasFreeSpins = false;
    config.hasJackpot = false;
    config.volatility = "medium";
    config.theme = "classic";
    config.symbols = {"DICE", "7", "BAR", "CHERRY"};
    config.symbolValues = {{"DICE", 20.0}, {"7", 10.0}, {"BAR", 5.0}, {"CHERRY", 3.0}};
    config.symbolWeights = {{"DICE", 2}, {"7", 4}, {"BAR", 8}, {"CHERRY", 12}};
    return config;
}
ExtendedSlotConfig createAurora() { return createDice(); }
ExtendedSlotConfig createChaosCrew() { return createWantedDeadOrAlive(); }
ExtendedSlotConfig createCabinCrashers() { return createWantedDeadOrAlive(); }
ExtendedSlotConfig createBeastMode() { return createWantedDeadOrAlive(); }
ExtendedSlotConfig createOutsmart() { return createWantedDeadOrAlive(); }
ExtendedSlotConfig createTheEmperor() { return createWantedDeadOrAlive(); }

// Nolimit City implementations
ExtendedSlotConfig createSanQuentin() {
    ExtendedSlotConfig config;
    config.gameId = "san_quentin";
    config.gameName = "San Quentin";
    config.provider = "NolimitCity";
    config.reels = 5;
    config.rows = 3;
    config.paylines = 0;
    config.minBet = 0.20;
    config.maxBet = 32.0;
    config.rtp = 0.96;
    config.hasBonus = true;
    config.hasFreeSpins = true;
    config.hasJackpot = false;
    config.volatility = "high";
    config.theme = "prison";
    config.symbols = {"PRISONER", "KNIFE", "TATTOO", "GANG", "A", "K"};
    config.symbolValues = {{"PRISONER", 30.0}, {"KNIFE", 15.0}, {"TATTOO", 10.0}, {"GANG", 6.0}, {"A", 1.5}, {"K", 1.2}};
    config.symbolWeights = {{"PRISONER", 1}, {"KNIFE", 2}, {"TATTOO", 4}, {"GANG", 6}, {"A", 10}, {"K", 12}};
    return config;
}
ExtendedSlotConfig createBookOfShadows() { return createSanQuentin(); }
ExtendedSlotConfig createMental() { return createSanQuentin(); }
ExtendedSlotConfig createXxxtreme() { return createSanQuentin(); }
ExtendedSlotConfig createDeadPanda() { return createSanQuentin(); }
ExtendedSlotConfig createFireInTheHole() { return createSanQuentin(); }
ExtendedSlotConfig createKiss() { return createSanQuentin(); }
ExtendedSlotConfig createLarryTheLeprechaun() { return createSanQuentin(); }
ExtendedSlotConfig createMillion777() { return createSanQuentin(); }
ExtendedSlotConfig createPontiusPilate() { return createSanQuentin(); }

} // namespace SlotGameTemplates

} // namespace TigerCasino
