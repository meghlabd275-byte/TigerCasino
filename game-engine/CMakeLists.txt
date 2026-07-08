cmake_minimum_required(VERSION 3.15)
project(TigerCasinoGameEngine VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

# Optimization flags
set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} -O3 -march=native -DNDEBUG")
set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -g -O0")

# Find OpenSSL for cryptographic functions
find_package(OpenSSL REQUIRED)

# Include directories
include_directories(${CMAKE_SOURCE_DIR}/include)
include_directories(${OPENSSL_INCLUDE_DIR})

# Source files for the main game engine library
set(GAME_ENGINE_SOURCES
    src/tiger_game_engine.cpp
    src/GameEngine.cpp
    src/DiceGame.cpp
    src/CrashGameServer.cpp
    src/SlotGameServer.cpp
    src/TableGameServer.cpp
    src/LiveDealerServer.cpp
    src/RandomNumberGenerator.cpp
)

# Create the main game engine shared library
add_library(tigercasino_game_engine SHARED ${GAME_ENGINE_SOURCES})

# Link OpenSSL
target_link_libraries(tigercasino_game_engine ${OPENSSL_LIBRARIES})

# Set output directory
set_target_properties(tigercasino_game_engine PROPERTIES
    LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/build
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/build
)

# Build directory
file(MAKE_DIRECTORY ${CMAKE_BINARY_DIR}/build)

# Simple C-compatible library for Go CGo bindings
add_library(tigercasino_simple SHARED src/lib.cpp)
target_link_libraries(tigercasino_simple OpenSSL::SSL OpenSSL::Crypto)
set_target_properties(tigercasino_simple PROPERTIES
    LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/build
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/build
)
