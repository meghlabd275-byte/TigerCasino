//! Mines Game Implementation
//! 
//! A strategic game where players click tiles to uncover rewards while avoiding hidden mines.

use serde::{Deserialize, Serialize};
use super::super::provably_fair::{GameRound, GameType};

/// Mines game configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MinesConfig {
    pub mines_count: u8,  // 1-24 mines
    pub grid_size: u8,     // Typically 5x5 = 25 tiles
}

impl Default for MinesConfig {
    fn default() -> Self {
        Self {
            mines_count: 3,
            grid_size: 25,
        }
    }
}

/// Mines game state
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MinesGame {
    pub game_id: String,
    pub config: MinesConfig,
    pub revealed_tiles: Vec<u8>,
    pub mines_positions: Vec<u8>,
    pub status: MinesStatus,
    pub current_multiplier: f64,
    pub bet_amount: f64,
    pub potential_win: f64,
    pub house_edge: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum MinesStatus {
    Playing,
    Won,
    Lost,
}

impl MinesGame {
    pub fn new(mines_count: u8) -> Self {
        Self {
            game_id: uuid::Uuid::new_v4().to_string(),
            config: MinesConfig {
                mines_count: mines_count.min(24).max(1),
                grid_size: 25,
            },
            revealed_tiles: Vec::new(),
            mines_positions: Vec::new(),
            status: MinesStatus::Playing,
            current_multiplier: 1.0,
            bet_amount: 0.0,
            potential_win: 0.0,
            house_edge: 0.05,
        }
    }
    
    /// Generate mine positions using provably fair seeds
    pub fn generate_mines(server_seed: &str, client_seed: &str, nonce: u64, mines_count: u8) -> Vec<u8> {
        use sha2::{Sha256, Digest};
        
        let mut mines = Vec::new();
        let mut hash = {
            let combined = format!("{}:{}:{}", server_seed, client_seed, nonce);
            let mut hasher = Sha256::new();
            hasher.update(combined.as_bytes());
            hex::encode(hasher.finalize())
        };
        
        while mines.len() < mines_count as usize {
            let hash_bytes = &hex::decode(&hash).unwrap_or_default()[..8];
            let seed_value = u64::from_le_bytes([
                hash_bytes[0], hash_bytes[1], hash_bytes[2], hash_bytes[3],
                hash_bytes[4], hash_bytes[5], hash_bytes[6], hash_bytes[7],
            ]);
            
            let position = (seed_value % 25) as u8;
            if !mines.contains(&position) {
                mines.push(position);
            }
            
            // Generate next hash for next mine
            let mut hasher = Sha256::new();
            hasher.update(hash.as_bytes());
            hash = hex::encode(hasher.finalize());
        }
        
        mines.sort();
        mines
    }
    
    /// Initialize game with bet
    pub fn start_game(&mut self, bet_amount: f64) {
        self.bet_amount = bet_amount;
        self.status = MinesStatus::Playing;
        self.revealed_tiles.clear();
        self.current_multiplier = 1.0;
    }
    
    /// Calculate current multiplier based on revealed tiles
    pub fn calculate_multiplier(&self) -> f64 {
        let revealed = self.revealed_tiles.len() as u8;
        let total_tiles = 25u8;
        let safe_tiles = total_tiles - self.config.mines_count;
        
        if revealed == 0 {
            return 1.0;
        }
        
        // Progressive multiplier calculation
        let base = 1.0;
        let increment = (self.config.mines_count as f64 * 0.3).max(0.5);
        
        base + (revealed as f64 * increment)
    }
    
    /// Reveal a tile
    pub fn reveal_tile(&mut self, tile_index: u8) -> RevealResult {
        if tile_index >= 25 {
            return RevealResult::Invalid;
        }
        
        if self.revealed_tiles.contains(&tile_index) {
            return RevealResult::AlreadyRevealed;
        }
        
        if self.mines_positions.is_empty() {
            return RevealResult::GameNotStarted;
        }
        
        if self.mines_positions.contains(&tile_index) {
            self.status = MinesStatus::Lost;
            return RevealResult::MineHit {
                mines: self.mines_positions.clone(),
                payout: 0.0,
            };
        }
        
        self.revealed_tiles.push(tile_index);
        self.current_multiplier = self.calculate_multiplier();
        self.potential_win = self.bet_amount * self.current_multiplier * (1.0 - self.house_edge);
        
        // Check if all safe tiles revealed
        let safe_tiles = 25 - self.config.mines_count;
        if self.revealed_tiles.len() as u8 >= safe_tiles {
            self.status = MinesStatus::Won;
            RevealResult::GameWon {
                payout: self.potential_win,
            }
        } else {
            RevealResult::TileRevealed {
                tile: tile_index,
                multiplier: self.current_multiplier,
                potential_win: self.potential_win,
            }
        }
    }
    
    /// Cash out current winnings
    pub fn cashout(&self) -> f64 {
        if self.status == MinesStatus::Won {
            self.potential_win
        } else {
            0.0
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum RevealResult {
    Invalid,
    AlreadyRevealed,
    GameNotStarted,
    TileRevealed { tile: u8, multiplier: f64, potential_win: f64 },
    MineHit { mines: Vec<u8>, payout: f64 },
    GameWon { payout: f64 },
}

/// Player bet in mines
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MinesBet {
    pub bet_id: String,
    pub player_id: String,
    pub bet_amount: f64,
    pub mines_count: u8,
    pub revealed_count: u8,
    pub cashed_out: bool,
    pub payout: f64,
    pub created_at: i64,
}

impl MinesBet {
    pub fn new(player_id: &str, bet_amount: f64, mines_count: u8) -> Self {
        Self {
            bet_id: uuid::Uuid::new_v4().to_string(),
            player_id: player_id.to_string(),
            bet_amount,
            mines_count,
            revealed_count: 0,
            cashed_out: false,
            payout: 0.0,
            created_at: chrono::Utc::now().timestamp(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_mine_generation() {
        let mines = MinesGame::generate_mines("server", "client", 1, 3);
        assert_eq!(mines.len(), 3);
    }
    
    #[test]
    fn test_multiplier_calculation() {
        let mut game = MinesGame::new(3);
        game.start_game(100.0);
        assert_eq!(game.current_multiplier, 1.0);
    }
}
