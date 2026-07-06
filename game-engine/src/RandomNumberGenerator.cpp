#include "RandomNumberGenerator.hpp"
#include <openssl/rand.h>
#include <algorithm>

namespace TigerCasino {

RandomNumberGenerator::RandomNumberGenerator() 
    : gen_(rd_()) {
}

int64_t RandomNumberGenerator::nextInt(int64_t min, int64_t max) {
    std::uniform_int_distribution<int64_t> dist(min, max);
    return dist(gen_);
}

double RandomNumberGenerator::nextDouble() {
    std::uniform_real_distribution<double> dist(0.0, 1.0);
    return dist(gen_);
}

bool RandomNumberGenerator::nextBool(double probability) {
    return nextDouble() < probability;
}

std::vector<uint8_t> RandomNumberGenerator::nextBytes(size_t length) {
    std::vector<uint8_t> bytes(length);
    RAND_bytes(bytes.data(), static_cast<int>(length));
    return bytes;
}

uint64_t RandomNumberGenerator::generateSeed() {
    std::vector<uint8_t> seed = nextBytes(8);
    uint64_t result = 0;
    for (size_t i = 0; i < 8; ++i) {
        result = (result << 8) | seed[i];
    }
    return result;
}

} // namespace TigerCasino
