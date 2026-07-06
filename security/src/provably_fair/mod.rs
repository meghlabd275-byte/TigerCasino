//! Provably Fair Gaming System for TigerCasino
//! 
//! This module implements cryptographic verification for all casino games,
//! allowing players to independently verify the fairness of each game outcome.

use sha2::{Sha256, Digest};
use serde::{Deserialize, Serialize};
use uuid::Uuid;
use std::time::{SystemTime, UNIX_EPOCH};

/// Represents a game round with provably fair verification data
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GameRound {
    pub round_id: String,
    pub server_seed: String,
    pub server_seed_hash: String,
    pub client_seed: String,
    pub nonce: u64,
    pub game_type: GameType,
    pub outcome: GameOutcome,
    pub timestamp: u64,
    pub player_id: String,
}

impl GameRound {
    /// Creates a new game round with generated server seed
    pub fn new(game_type: GameType, client_seed: &str, player_id: &str, nonce: u64) -> Self {
        let server_seed = generate_random_seed();
        let server_seed_hash = hash_seed(&server_seed);
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        let outcome = calculate_outcome(&server_seed, client_seed, nonce, &game_type);
        
        Self {
            round_id: Uuid::new_v4().to_string(),
            server_seed,
            server_seed_hash,
            client_seed: client_seed.to_string(),
            nonce,
            game_type,
            outcome,
            timestamp,
            player_id: player_id.to_string(),
        }
    }
    
    /// Verify the fairness of a game round
    pub fn verify(&self, provided_server_seed: &str) -> bool {
        let provided_hash = hash_seed(provided_server_seed);
        if provided_hash != self.server_seed_hash {
            return false;
        }
        
        let recalculated_outcome = calculate_outcome(
            provided_server_seed, 
            &self.client_seed, 
            self.nonce, 
            &self.game_type
        );
        
        recalculated_outcome == self.outcome
    }
}

/// Game types supported by the provably fair system
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum GameType {
    Dice,
    Crash,
    Mines,
    Plinko,
    Blackjack,
    Roulette,
    Slots,
    Baccarat,
    Poker,
    HiLo,
    Keno,
    Lottery,
}

impl GameType {
    pub fn as_str(&self) -> &str {
        match self {
            GameType::Dice => "dice",
            GameType::Crash => "crash",
            GameType::Mines => "mines",
            GameType::Plinko => "plinko",
            GameType::Blackjack => "blackjack",
            GameType::Roulette => "roulette",
            GameType::Slots => "slots",
            GameType::Baccarat => "baccarat",
            GameType::Poker => "poker",
            GameType::HiLo => "hilo",
            GameType::Keno => "keno",
            GameType::Lottery => "lottery",
        }
    }
}

/// Outcome of a provably fair game round
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct GameOutcome {
    pub result: String,
    pub multiplier: f64,
    pub win_amount: f64,
    pub is_win: bool,
}

/// Generate a cryptographically secure random seed
pub fn generate_random_seed() -> String {
    use rand::RngCore;
    let mut bytes = [0u8; 32];
    rand::rngs::OsRng.fill_bytes(&mut bytes);
    hex::encode(bytes)
}

/// Hash a seed using SHA-256
pub fn hash_seed(seed: &str) -> String {
    let mut hasher = Sha256::new();
    hasher.update(seed.as_bytes());
    let result = hasher.finalize();
    hex::encode(result)
}

/// Calculate game outcome from seeds
fn calculate_outcome(server_seed: &str, client_seed: &str, nonce: u64, game_type: &GameType) -> GameOutcome {
    let combined = format!("{}:{}:{}", server_seed, client_seed, nonce);
    let hash = hash_seed(&combined);
    let hash_bytes = hex::decode(&hash).unwrap_or_default();
    
    // Use first 8 bytes as u64 for randomness
    let seed_value = u64::from_le_bytes([
        hash_bytes.get(0).copied().unwrap_or(0),
        hash_bytes.get(1).copied().unwrap_or(0),
        hash_bytes.get(2).copied().unwrap_or(0),
        hash_bytes.get(3).copied().unwrap_or(0),
        hash_bytes.get(4).copied().unwrap_or(0),
        hash_bytes.get(5).copied().unwrap_or(0),
        hash_bytes.get(6).copied().unwrap_or(0),
        hash_bytes.get(7).copied().unwrap_or(0),
    ]);
    
    match game_type {
        GameType::Dice => calculate_dice_outcome(seed_value),
        GameType::Crash => calculate_crash_outcome(seed_value),
        GameType::Mines => calculate_mines_outcome(seed_value),
        GameType::Plinko => calculate_plinko_outcome(seed_value),
        GameType::Slots => calculate_slots_outcome(seed_value),
        GameType::HiLo => calculate_hilo_outcome(seed_value),
        GameType::Keno => calculate_keno_outcome(seed_value),
        GameType::Lottery => calculate_lottery_outcome(seed_value),
        _ => GameOutcome {
            result: "pending".to_string(),
            multiplier: 0.0,
            win_amount: 0.0,
            is_win: false,
        }
    }
}

/// Calculate dice game outcome (0-100)
fn calculate_dice_outcome(seed: u64) -> GameOutcome {
    let result = (seed % 10001) as f64 / 100.0;
    GameOutcome {
        result: format!("{:.2}", result),
        multiplier: 1.0,
        win_amount: 0.0,
        is_win: false,
    }
}

/// Calculate crash game outcome
fn calculate_crash_outcome(seed: u64) -> GameOutcome {
    // Crash point typically between 1.00x and 100x
    let hash_mod = seed % 1000000;
    let multiplier = if hash_mod < 500000 {
        1.0 + (hash_mod as f64 / 500000.0) * 9.0 // 1.0 - 10.0
    } else {
        10.0 + ((hash_mod - 500000) as f64 / 500000.0) * 90.0 // 10.0 - 100.0
    };
    
    GameOutcome {
        result: format!("{:.2}x", multiplier),
        multiplier,
        win_amount: 0.0,
        is_win: false,
    }
}

/// Calculate mines game outcome
fn calculate_mines_outcome(seed: u64) -> GameOutcome {
    let mines_hit = (seed % 3) as u8; // 0-2 mines typically hit
    GameOutcome {
        result: format!("{} mines hit", mines_hit),
        multiplier: 0.0,
        win_amount: 0.0,
        is_win: false,
    }
}

/// Calculate plinko game outcome
fn calculate_plinko_outcome(seed: u64) -> GameOutcome {
    // Plinko has multiple payout tiers
    let tier = (seed % 16) as u8;
    let multiplier = match tier {
        0 => 0.5,
        1 | 15 => 1.0,
        2 | 14 => 2.0,
        3 | 13 => 5.0,
        4 | 12 => 10.0,
        5 | 11 => 25.0,
        6 | 10 => 50.0,
        7 | 9 => 100.0,
        8 => 1000.0,
        _ => 1.0,
    };
    
    GameOutcome {
        result: format!("Tier {} - {}x", tier, multiplier),
        multiplier,
        win_amount: 0.0,
        is_win: false,
    }
}

/// Calculate slots outcome
fn calculate_slots_outcome(seed: u64) -> GameOutcome {
    // Simplified slot outcome - 3 reels
    let reel1 = (seed % 10) as u8;
    let reel2 = ((seed / 10) % 10) as u8;
    let reel3 = ((seed / 100) % 10) as u8;
    
    let multiplier = if reel1 == reel2 && reel2 == reel3 {
        match reel1 {
            0 => 10.0, // Jackpot
            1..=3 => 5.0,
            _ => 3.0,
        }
    } else if reel1 == reel2 || reel2 == reel3 || reel1 == reel3 {
        2.0
    } else {
        0.0
    };
    
    GameOutcome {
        result: format!("[{}, {}, {}]", reel1, reel2, reel3),
        multiplier,
        win_amount: 0.0,
        is_win: multiplier > 0.0,
    }
}

/// Calculate Hi-Lo outcome
fn calculate_hilo_outcome(seed: u64) -> GameOutcome {
    let card_value = (seed % 13) as u8 + 1; // 1-13
    let suit = (seed / 13) % 4; // 0-3
    
    let suit_char = match suit {
        0 => '♠',
        1 => '♥',
        2 => '♦',
        3 => '♣',
        _ => '?',
    };
    
    GameOutcome {
        result: format!("{}{}", card_value, suit_char),
        multiplier: 1.0,
        win_amount: 0.0,
        is_win: false,
    }
}

/// Calculate Keno outcome
fn calculate_keno_outcome(seed: u64) -> GameOutcome {
    let hits = (seed % 11) as u8; // 0-10 hits
    GameOutcome {
        result: format!("{} hits", hits),
        multiplier: hits as f64 * 0.5,
        win_amount: 0.0,
        is_win: hits >= 3,
    }
}

/// Calculate Lottery outcome
fn calculate_lottery_outcome(seed: u64) -> GameOutcome {
    let matched = (seed % 7) as u8; // 0-6 matches
    let multiplier = match matched {
        6 => 1000000.0, // Jackpot
        5 => 10000.0,
        4 => 1000.0,
        3 => 100.0,
        2 => 10.0,
        1 => 2.0,
        _ => 0.0,
    };
    
    GameOutcome {
        result: format!("{} matches", matched),
        multiplier,
        win_amount: 0.0,
        is_win: matched >= 1,
    }
}

/// Verification data for players to check fairness
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VerificationData {
    pub round_id: String,
    pub server_seed_hash: String,
    pub client_seed: String,
    pub nonce: u64,
    pub game_type: String,
    pub outcome: String,
    pub verify_url: String,
}

impl GameRound {
    /// Generate verification data for player
    pub fn to_verification_data(&self) -> VerificationData {
        VerificationData {
            round_id: self.round_id.clone(),
            server_seed_hash: self.server_seed_hash.clone(),
            client_seed: self.client_seed.clone(),
            nonce: self.nonce,
            game_type: self.game_type.as_str().to_string(),
            outcome: self.outcome.result.clone(),
            verify_url: format!("/api/fairness/verify/{}", self.round_id),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_seed_generation() {
        let seed = generate_random_seed();
        assert_eq!(seed.len(), 64); // SHA-256 hex = 64 chars
    }
    
    #[test]
    fn test_hash_seed() {
        let hash = hash_seed("test_seed");
        assert_eq!(hash.len(), 64);
    }
    
    #[test]
    fn test_dice_outcome() {
        let outcome = calculate_dice_outcome(123456789);
        let value: f64 = outcome.result.parse().unwrap_or(0.0);
        assert!(value >= 0.0 && value <= 100.0);
    }
    
    #[test]
    fn test_round_verification() {
        let round = GameRound::new(
            GameType::Dice, 
            "player_seed", 
            "player123", 
            1
        );
        
        assert!(round.verify(&round.server_seed));
    }
}
