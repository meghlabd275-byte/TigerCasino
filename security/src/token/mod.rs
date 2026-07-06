//! Native Token ($TIGER) Module for TigerCasino
//! 
//! Tokenomics, staking, rewards, and governance

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Token configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenConfig {
    pub name: String,
    pub symbol: String,
    pub decimals: u8,
    pub total_supply: u64,
    pub initial_price: f64,
}

/// Token holder
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenHolder {
    pub address: String,
    pub balance: u64,
    pub staked_balance: u64,
    pub rewards_earned: f64,
    pub vip_level: u8,
}

/// Staking position
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StakingPosition {
    pub id: String,
    pub owner: String,
    pub amount: u64,
    pub start_time: u64,
    pub lock_period: u64,
    pub rewards_claimed: f64,
    pub multiplier: f64,
}

/// Token distribution
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenDistribution {
    pub category: String,
    pub percentage: f64,
    pub amount: u64,
    pub cliff_period: u64,
    pub total_duration: u64,
}

/// Rewards pool
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RewardsPool {
    pub total_allocated: u64,
    pub distributed: u64,
    pub reward_rate: f64,
}

/// Native token manager
pub struct TokenManager {
    config: TokenConfig,
    holders: HashMap<String, TokenHolder>,
    staking_positions: HashMap<String, StakingPosition>,
    distributions: Vec<TokenDistribution>,
    rewards_pool: RewardsPool,
    exchange_rate: f64,
}

impl TokenManager {
    pub fn new() -> Self {
        Self {
            config: TokenConfig {
                name: "TigerCasino Token".to_string(),
                symbol: "TIGER".to_string(),
                decimals: 18,
                total_supply: 1_000_000_000,
                initial_price: 0.10,
            },
            holders: HashMap::new(),
            staking_positions: HashMap::new(),
            distributions: Vec::new(),
            rewards_pool: RewardsPool {
                total_allocated: 0,
                distributed: 0,
                reward_rate: 0.05,
            },
            exchange_rate: 0.10,
        }
    }
    
    /// Initialize token distribution
    pub fn init_distribution(&mut self) {
        self.distributions.push(TokenDistribution {
            category: "airdrop".to_string(),
            percentage: 20.0,
            amount: 200_000_000,
            cliff_period: 0,
            total_duration: 365 * 24 * 60 * 60,
        });
        
        self.distributions.push(TokenDistribution {
            category: "staking".to_string(),
            percentage: 30.0,
            amount: 300_000_000,
            cliff_period: 90 * 24 * 60 * 60,
            total_duration: 730 * 24 * 60 * 60,
        });
        
        self.distributions.push(TokenDistribution {
            category: "team".to_string(),
            percentage: 15.0,
            amount: 150_000_000,
            cliff_period: 365 * 24 * 60 * 60,
            total_duration: 1095 * 24 * 60 * 60,
        });
    }
    
    /// Stake tokens
    pub fn stake(&mut self, address: &str, amount: u64, lock_period: u64) -> String {
        let position_id = format!("STAKE-{}", uuid::Uuid::new_v4());
        
        let multiplier = match lock_period {
            2592000 => 1.2,
            7776000 => 1.5,
            15552000 => 2.0,
            31536000 => 3.0,
            _ => 1.0,
        };
        
        let position = StakingPosition {
            id: position_id.clone(),
            owner: address.to_string(),
            amount,
            start_time: 1700000000,
            lock_period,
            rewards_claimed: 0.0,
            multiplier,
        };
        
        self.staking_positions.insert(position_id.clone(), position);
        
        let holder = self.holders.entry(address.to_string()).or_insert(TokenHolder {
            address: address.to_string(),
            balance: 0,
            staked_balance: 0,
            rewards_earned: 0.0,
            vip_level: 0,
        });
        holder.staked_balance += amount;
        
        position_id
    }
    
    /// Calculate staking rewards
    pub fn calculate_rewards(&self, position_id: &str) -> f64 {
        if let Some(position) = self.staking_positions.get(position_id) {
            let days = 30.0;
            let apy = self.rewards_pool.reward_rate * position.multiplier;
            return position.amount as f64 * apy * days / 365.0;
        }
        0.0
    }
    
    /// Get token price
    pub fn get_price(&self) -> f64 {
        self.exchange_rate
    }
    
    /// Convert TIGER to USD
    pub fn to_usd(&self, tiger_amount: u64) -> f64 {
        tiger_amount as f64 * self.exchange_rate
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_staking() {
        let mut manager = TokenManager::new();
        let pos_id = manager.stake("user1", 10000, 7776000);
        assert!(!pos_id.is_empty());
    }
}
