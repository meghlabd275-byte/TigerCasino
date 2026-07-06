cmake_minimum_required(VERSION 3.15)
project(TigerCasinoGameEngine VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find OpenSSL for cryptographic functions
find_package(OpenSSL REQUIRED)

# Create the library
add_library(tigercasino_game_engine STATIC
    src/GameEngine.cpp
    src/RandomNumberGenerator.cpp
    src/SlotMachine.cpp
    src/DiceGame.cpp
    src/Roulette.cpp
    src/Blackjack.cpp
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
