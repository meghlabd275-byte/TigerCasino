#include <iostream>
#include <vector>
#include <string>
#include <random>

extern "C" {
    double cpp_generate_dice_roll() {
        static std::random_device rd;
        static std::mt19937 gen(rd());
        std::uniform_real_distribution<> dis(0, 100);
        return dis(gen);
    }

    int cpp_calculate_slots(int bet) {
        static std::random_device rd;
        static std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(0, 100);
        int chance = dis(gen);
        if (chance > 95) return bet * 10;
        if (chance > 80) return bet * 2;
        return 0;
    }

    int cpp_spin_roulette() {
        static std::random_device rd;
        static std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(0, 36);
        return dis(gen);
    }

    int cpp_play_blackjack(int bet) {
        static std::random_device rd;
        static std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(1, 100);
        int chance = dis(gen);
        if (chance > 55) return bet * 2; // Player win
        if (chance > 45) return bet;     // Push
        return 0;                        // Dealer win
    }

    int cpp_play_video_poker(int bet) {
        static std::random_device rd;
        static std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(1, 100);
        int chance = dis(gen);
        if (chance > 98) return bet * 25; // Four of a kind
        if (chance > 90) return bet * 5;  // Flush
        if (chance > 70) return bet * 2;  // Two pair
        if (chance > 50) return bet;      // Jacks or better
        return 0;
    }
}
