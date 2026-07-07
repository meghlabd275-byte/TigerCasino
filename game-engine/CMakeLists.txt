cmake_minimum_required(VERSION 3.15)
project(TigerCasinoGameEngine VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)

add_library(tigercasino_game_engine SHARED src/lib.cpp)
