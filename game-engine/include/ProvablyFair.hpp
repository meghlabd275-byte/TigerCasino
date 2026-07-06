#ifndef PROVABLY_FAIR_HPP
#define PROVABLY_FAIR_HPP

#include <string>
#include <vector>
#include <cstdint>
#include <sstream>
#include <iomanip>
#include <openssl/sha.h>
#include <openssl/rand.h>

namespace TigerCasino {

/**
 * Provably Fair System for Ultra-Low Latency Game Outcomes
 * Uses SHA-256 for cryptographic verification
 */
class ProvablyFair {
public:
    struct GameSeeds {
        std::string serverSeed;
        std::string serverSeedHash;
        std::string clientSeed;
        uint64_t nonce;
    };
    
    struct GameResult {
        std::string result;
        double multiplier;
        double winAmount;
        bool isWin;
    };

    /**
     * Generate a cryptographically secure random server seed
     */
    static std::string generateServerSeed() {
        uint8_t buffer[32];
        RAND_bytes(buffer, 32);
        
        std::stringstream ss;
        for (int i = 0; i < 32; ++i) {
            ss << std::hex << std::setw(2) << std::setfill('0') << (int)buffer[i];
        }
        return ss.str();
    }
    
    /**
     * Hash a seed using SHA-256
     */
    static std::string hashSeed(const std::string& seed) {
        unsigned char hash[SHA256_DIGEST_LENGTH];
        SHA256((const unsigned char*)seed.c_str(), seed.length(), hash);
        
        std::stringstream ss;
        for (int i = 0; i < SHA256_DIGEST_LENGTH; ++i) {
            ss << std::hex << std::setw(2) << std::setfill('0') << (int)hash[i];
        }
        return ss.str();
    }
    
    /**
     * Generate a deterministic outcome from seeds
     */
    static uint64_t generateOutcome(const std::string& serverSeed, 
                                     const std::string& clientSeed, 
                                     uint64_t nonce) {
        std::stringstream ss;
        ss << serverSeed << ":" << clientSeed << ":" << nonce;
        std::string combined = ss.str();
        
        unsigned char hash[SHA256_DIGEST_LENGTH];
        SHA256((const unsigned char*)combined.c_str(), combined.length(), hash);
        
        // Use first 8 bytes as outcome
        uint64_t result = 0;
        for (int i = 0; i < 8; ++i) {
            result = (result << 8) | hash[i];
        }
        return result;
    }
    
    /**
     * Verify a game round was fair
     */
    static bool verify(const std::string& providedServerSeed,
                       const std::string& storedSeedHash,
                       const std::string& clientSeed,
                       uint64_t nonce) {
        return hashSeed(providedServerSeed) == storedSeedHash;
    }
};

} // namespace TigerCasino

#endif // PROVABLY_FAIR_HPP
