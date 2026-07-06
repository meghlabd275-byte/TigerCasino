//! Plinko Game Implementation
//! 
//! A popular game where balls fall through a pegged board, landing in multiplier pockets at the bottom.

use serde::{Deserialize, Serialize};

/// Plinko game configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlinkoConfig {
    pub rows: u8,           // Number of rows (8, 10, 12, 16)
    pub risk: PlinkoRisk,   // Low, Medium, High
}

impl Default for PlinkoConfig {
    fn default() -> Self {
        Self {
            rows: 8,
            risk: PlinkoRisk::Medium,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum PlinkoRisk {
    Low,
    Medium,
    High,
}

/// Plinko game state
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlinkoGame {
    pub game_id: String,
    pub config: PlinkoConfig,
    pub balls: Vec<PlinkoBall>,
    pub status: PlinkoStatus,
    pub house_edge: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum PlinkoStatus {
    Waiting,
    Dropping,
    Complete,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlinkoBall {
    pub ball_id: String,
    pub path: Vec<u8>,      // Path through rows
    pub final_position: u8, // Final pocket
    pub multiplier: f64,
    pub payout: f64,
}

impl PlinkoGame {
    pub fn new(rows: u8, risk: PlinkoRisk) -> Self {
        Self {
            game_id: uuid::Uuid::new_v4().to_string(),
            config: PlinkoConfig {
                rows: rows.max(8).min(16),
                risk,
            },
            balls: Vec::new(),
            status: PlinkoStatus::Waiting,
            house_edge: 0.04,
        }
    }
    
    /// Get payout table based on rows and risk
    pub fn get_payout_table(&self) -> Vec<f64> {
        match self.config.rows {
            8 => match self.config.risk {
                PlinkoRisk::Low => vec![1.5, 1.2, 0.8, 0.5, 0.5, 0.8, 1.2, 1.5],
                PlinkoRisk::Medium => vec![5.0, 2.0, 1.0, 0.5, 0.5, 1.0, 2.0, 5.0],
                PlinkoRisk::High => vec![10.0, 5.0, 2.0, 0.5, 0.5, 2.0, 5.0, 10.0],
            },
            10 => match self.config.risk {
                PlinkoRisk::Low => vec![1.2, 1.0, 0.8, 0.6, 0.4, 0.4, 0.6, 0.8, 1.0, 1.2],
                PlinkoRisk::Medium => vec![5.0, 2.5, 1.5, 0.8, 0.4, 0.4, 0.8, 1.5, 2.5, 5.0],
                PlinkoRisk::High => vec![15.0, 8.0, 3.0, 1.5, 0.5, 0.5, 1.5, 3.0, 8.0, 15.0],
            },
            12 => match self.config.risk {
                PlinkoRisk::Low => vec![1.1, 0.9, 0.7, 0.5, 0.3, 0.3, 0.3, 0.3, 0.5, 0.7, 0.9, 1.1],
                PlinkoRisk::Medium => vec![4.0, 2.0, 1.2, 0.7, 0.3, 0.3, 0.3, 0.3, 0.7, 1.2, 2.0, 4.0],
                PlinkoRisk::High => vec![15.0, 10.0, 4.0, 2.0, 0.5, 0.3, 0.3, 0.5, 2.0, 4.0, 10.0, 15.0],
            },
            16 => match self.config.risk {
                PlinkoRisk::Low => vec![0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 1.0, 1.1, 1.1, 1.0, 1.0, 0.9, 0.8, 0.7, 0.6, 0.5],
                PlinkoRisk::Medium => vec![1.0, 1.5, 2.0, 3.0, 4.0, 5.0, 5.0, 6.0, 6.0, 5.0, 5.0, 4.0, 3.0, 2.0, 1.5, 1.0],
                PlinkoRisk::High => vec![10.0, 15.0, 20.0, 30.0, 40.0, 50.0, 60.0, 100.0, 100.0, 60.0, 50.0, 40.0, 30.0, 20.0, 15.0, 10.0],
            },
            _ => vec![1.0; 8],
        }
    }
    
    /// Simulate ball drop using provably fair seeds
    pub fn drop_ball(&mut self, server_seed: &str, client_seed: &str, nonce: u64, bet_amount: f64) -> PlinkoBall {
        use sha2::{Sha256, Digest};
        
        let mut path = Vec::new();
        let mut current_position: u8 = 0; // Start at center
        let rows = self.config.rows as usize;
        
        // Simulate ball path through each row
        for row in 0..rows {
            let combined = format!("{}:{}:{}_{}", server_seed, client_seed, nonce, row);
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
            
            // Determine direction: 0 = left, 1 = right
            let direction = (seed_value % 2) as u8;
            
            // Calculate new position
            // At row 0, positions are 0-1
            // At row 1, positions are 0-2
            // etc.
            let max_pos = (row + 2) as u8;
            
            if direction == 0 && current_position > 0 {
                current_position -= 1;
            } else if direction == 0 && current_position == 0 {
                // Can't go left, go right
                current_position += 1;
            }
            // If direction is 1, stay or move (would need adjustment based on position)
            
            // Keep position in bounds
            current_position = current_position.min(max_pos - 1);
            
            path.push(current_position);
        }
        
        // Calculate final multiplier based on position
        let payout_table = self.get_payout_table();
        let final_position = current_position.min((payout_table.len() - 1) as u8);
        let multiplier = payout_table[final_position as usize];
        
        let payout = bet_amount * multiplier * (1.0 - self.house_edge);
        
        let ball = PlinkoBall {
            ball_id: uuid::Uuid::new_v4().to_string(),
            path,
            final_position,
            multiplier,
            payout,
        };
        
        self.balls.push(ball.clone());
        ball
    }
    
    /// Get active balls in game
    pub fn get_balls(&self) -> &[PlinkoBall] {
        &self.balls
    }
}

/// Plinko bet
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlinkoBet {
    pub bet_id: String,
    pub player_id: String,
    pub bet_amount: f64,
    pub rows: u8,
    pub risk: PlinkoRisk,
    pub balls_count: u8,
    pub total_payout: f64,
    pub created_at: i64,
}

impl PlinkoBet {
    pub fn new(player_id: &str, bet_amount: f64, rows: u8, risk: PlinkoRisk, balls_count: u8) -> Self {
        Self {
            bet_id: uuid::Uuid::new_v4().to_string(),
            player_id: player_id.to_string(),
            bet_amount,
            rows,
            risk,
            balls_count,
            total_payout: 0.0,
            created_at: chrono::Utc::now().timestamp(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_payout_table() {
        let game = PlinkoGame::new(8, PlinkoRisk::Medium);
        let table = game.get_payout_table();
        assert_eq!(table.len(), 8);
    }
}
