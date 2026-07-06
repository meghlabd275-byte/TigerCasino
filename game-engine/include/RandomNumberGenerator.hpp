#pragma once

#include <cstdint>
#include <random>

namespace TigerCasino {

// Cryptographically secure random number generator
class RandomNumberGenerator {
private:
    std::random_device rd_;
    std::mt19937_64 gen_;
    
public:
    RandomNumberGenerator();
    ~RandomNumberGenerator() = default;
    
    // Generate random integer in range [min, max]
    int64_t nextInt(int64_t min, int64_t max);
    
    // Generate random double in range [0, 1)
    double nextDouble();
    
    // Generate random boolean with probability
    bool nextBool(double probability = 0.5);
    
    // Generate random bytes
    std::vector<uint8_t> nextBytes(size_t length);
    
    // Generate seed from system entropy
    uint64_t generateSeed();
};

} // namespace TigerCasino
