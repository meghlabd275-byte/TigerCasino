//! Crash Game Implementation
//! 
//! A fast-paced betting game where players watch a multiplier rise and try to cash out before it crashes.

use serde::{Deserialize, Serialize};
use super::super::provably_fair::{GameRound, GameType};

/// Crash game state
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CrashGame {
    pub game_id: String,
    pub current_multiplier: f64,
    pub crash_point: f64,
    pub status: CrashStatus,
    pub history: Vec<CrashHistoryItem>,
    pub house_edge: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CrashHistoryItem {
    pub crash_point: f64,
    pub timestamp: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum CrashStatus {
    Waiting,
    Rising,
    Crashed,
}

impl CrashGame {
    /// Create a new crash game instance
    pub fn new() -> Self {
        Self {
            game_id: uuid::Uuid::new_v4().to_string(),
            current_multiplier: 1.0,
            crash_point: 0.0,
            status: CrashStatus::Waiting,
            history: Vec::new(),
            house_edge: 0.03, // 3% house edge
        }
    }
    
    /// Generate crash point using provably fair system
    pub fn generate_crash_point(server_seed: &str, client_seed: &str, nonce: u64) -> f64 {
        use sha2::{Sha256, Digest};
        
        let combined = format!("{}:{}:{}", server_seed, client_seed, nonce);
        let hash = {
            let mut hasher = Sha256::new();
            hasher.update(combined.as_bytes());
            hex::encode(hasher.finalize())
        };
        
        let hash_bytes = &hex::decode(&hash).unwrap_or_default()[..8];
        let seed_value = u64::from_le_bytes([
            hash_bytes[0], hash_bytes[1], hash_bytes[2], hash_bytes[3],
            hash_bytes[4], hash_bytes[5], hash_bytes[6], hash_bytes[7],
        ]);
        
        // Generate crash point using exponential distribution
        // Most crashes happen early, but rare big multipliers possible
        let hash_mod = seed_valu  % 100000;
        
        if hash_mod < 30000 {
            // 30% chance of instant crash (below 1.1x)
            1.0 + (hash_mod as f64 / 30000.0) * 0.1
        } else {
            // 70% chance of higher crash point
            let x = (hash_mod - 30000) as f64 / 70000.0;
            // Exponential distribution with average around 2.5x
            1.1 + ((-x.ln()) * 2.5).min(100.0)
        }
    }
    
    /// Start a new crash round
    pub fn start_round(&mut self) {
        self.status = CrashStatus::Rising;
        self.current_multiplier = 1.0;
        // Crash point will be set when round ends
    }
    
    /// Calculate payout for a cashed out bet
    pub fn calculate_payout(&self, bet_amount: f64, multiplier: f64) -> f64 {
        bet_amount * multiplier * (1.0 - self.house_edge)
    }
    
    /// Add to crash history
    pub fn add_to_history(&mut self, crash_point: f64) {
        self.history.insert(0, CrashHistoryItem {
            crash_point,
            timestamp: chrono::Utc::now().timestamp(),
        });
        
        // Keep only last 100 items
        if self.history.len() > 100 {
            self.history.pop();
        }
    }
    
    /// Get recent crash history
    pub fn get_history(&self, count: usize) -> Vec<f64> {
        self.history
            .iter()
            .take(count)
            .map(|item| item.crash_point)
            .collect()
    }
}

impl Default for CrashGame {
    fn default() -> Self {
        Self::new()
    }
}

/// Player bet in crash game
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CrashBet {
    pub bet_id: String,
    pub player_id: String,
    pub bet_amount: f64,
    pub auto_cashout_at: Option<f64>,
    pub cashed_out: bool,
    pub cashout_multiplier: f64,
    pub payout: f64,
    pub created_at: i64,
}

impl CrashBet {
    pub fn new(player_id: &str, bet_amount: f64, auto_cashout_at: Option<f64>) -> Self {
        Self {
            bet_id: uuid::Uuid::new_v4().to_string(),
            player_id: player_id.to_string(),
            bet_amount,
            auto_cashout_at,
            cashed_out: false,
            cashout_multiplier: 0.0,
            payout: 0.0,
            created_at: chrono::Utc::now().timestamp(),
        }
    }
    
    /// Try to auto cashout
    pub fn try_auto_cashout(&mut self, current_multiplier: f64) -> bool {
        if !self.cashed_out {
            if let Some(threshold) = self.auto_cashout_at {
                if current_multiplier >= threshold {
                    self.cashed_out = true;
                    self.cashout_multiplier = current_multiplier;
                    return true;
                }
            }
        }
        false
    }
}

/// Crash game API response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CrashGameState {
    pub game_id: String,
    pub status: String,
    pub current_multiplier: f64,
    pub crash_point: Option<f64>,
    pub history: Vec<f64>,
    pub bets: Vec<CrashBetInfo>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CrashBetInfo {
    pub player_id: String,
    pub bet_amount: f64,
    pub auto_cashout_at: Option<f64>,
    pub cashed_out: bool,
    pub cashout_multiplier: Option<f64>,
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_crash_point_generation() {
        let crash = CrashGame::generate_crash_point(
            "server_seed_123",
            "client_seed_456",
            1
        );
        assert!(crash >= 1.0);
    }
    
    #[test]
    fn test_payout_calculation() {
        let game = CrashGame::new();
        let payout = game.calculate_payout(100.0, 2.5);
        // 100 * 2.5 * 0.97 = 242.5
        assert!((payout - 242.5).abs() < 0.01);
    }
}
