cmake_minimum_required(VERSION 3.15)
project(TigerCasinoGameEngine VERSION 2.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find OpenSSL for cryptographic functions
find_package(OpenSSL REQUIRED)

# Enable optimization flags for ultra-low latency
if(CMAKE_BUILD_TYPE STREQUAL "Release")
    add_compile_options(-O3 -march=native -mtune=native)
    add_compile_definitions(NDEBUG)
endif()

# Create the library
add_library(tigercasino_game_engine STATIC
    src/GameEngine.cpp
    src/RandomNumberGenerator.cpp
    src/SlotMachine.cpp
    src/SlotGameServer.cpp
    src/ExtendedSlotGames.cpp
    src/DiceGame.cpp
    src/Roulette.cpp
    src/Blackjack.cpp
    src/CrashGameServer.cpp
    src/TableGameServer.cpp
    src/LiveDealerServer.cpp
    src/OriginalGamesExtended.cpp
    src/ArcadeGames.cpp
    src/GameShows.cpp
    src/LotteryBingo.cpp
    src/ScratchCards.cpp
    src/TigerOriginals.cpp
    src/TigerOriginals2.cpp
    src/VirtualSports.cpp
    src/ZeppelinGame.cpp
    src/RocketQueenGame.cpp
)

target_include_directories(tigercasino_game_engine PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_link_libraries(tigercasino_game_engine PUBLIC
    OpenSSL::Crypto
)

# Create CLI executable
add_executable(game_engine_cli src/main.cpp)
target_link_libraries(game_engine_cli PRIVATE
    tigercasino_game_engine
)

# Install headers for external use
install(DIRECTORY include/ DESTINATION include/tigercasino
    FILES_MATCHING PATTERN "*.hpp")
