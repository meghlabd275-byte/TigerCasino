//! NFT Module for TigerCasino
//! 
//! NFT-based rewards, NFT-gated access, and NFT collections

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// NFT Collection
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NFTCollection {
    pub id: String,
    pub name: String,
    pub description: String,
    pub image_url: String,
    pub total_supply: u32,
    pub floor_price: f64,
    pub traits: Vec<NFTTrait>,
}

/// NFT Trait (for rarity)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NFTTrait {
    pub trait_type: String,
    pub value: String,
    pub rarity: f64,
}

/// Individual NFT
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NFT {
    pub id: String,
    pub collection_id: String,
    pub owner: String,
    pub token_id: u64,
    pub metadata_url: String,
    pub image_url: String,
    pub traits: Vec<NFTTrait>,
    pub rarity_score: f64,
    pub is_staked: bool,
    pub staked_at: Option<u64>,
}

/// NFT Reward Tier
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NFRewardTier {
    pub tier: String,
    pub nft_collection: String,
    pub bonus_percent: f64,
    pub rakeback_boost: f64,
    pub daily_free_spins: u32,
    pub vip_points_multiplier: f64,
}

/// NFT Staking Pool
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NFTStakingPool {
    pub collection_id: String,
    pub staked_nfts: HashMap<String, Vec<StakedNFT>>,
    pub rewards_per_day: f64,
    pub total_staked: u32,
    pub lock_period: u64,
}

/// Staked NFT
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StakedNFT {
    pub nft_id: String,
    pub owner: String,
    pub staked_at: u64,
    pub pending_rewards: f64,
    pub multiplier: f64,
}

/// NFT Manager
pub struct NFTManager {
    collections: HashMap<String, NFTCollection>,
    nfts: HashMap<String, NFT>,
    reward_tiers: HashMap<String, NFRewardTier>,
    staking_pools: HashMap<String, NFTStakingPool>,
}

impl NFTManager {
    pub fn new() -> Self {
        let mut manager = Self {
            collections: HashMap::new(),
            nfts: HashMap::new(),
            reward_tiers: HashMap::new(),
            staking_pools: HashMap::new(),
        };
        
        manager.init_default_collections();
        manager.init_reward_tiers();
        
        manager
    }
    
    fn init_default_collections(&mut self) {
        self.collections.insert("tiger-royal".to_string(), NFTCollection {
            id: "tiger-royal".to_string(),
            name: "Tiger Royal Collection".to_string(),
            description: "Exclusive NFT collection for TigerCasino VIP".to_string(),
            image_url: "https://tigercasino.com/nft/tiger-royal".to_string(),
            total_supply: 1000,
            floor_price: 0.5,
            traits: vec![
                NFTTrait { trait_type: "Background".to_string(), value: "Gold".to_string(), rarity: 0.1 },
                NFTTrait { trait_type: "Fur".to_string(), value: "Orange".to_string(), rarity: 0.3 },
            ],
        });
        
        self.collections.insert("lucky-cat".to_string(), NFTCollection {
            id: "lucky-cat".to_string(),
            name: "Lucky Cat NFTs".to_string(),
            description: "Lucky cat NFTs with bonus features".to_string(),
            image_url: "https://tigercasino.com/nft/lucky-cat".to_string(),
            total_supply: 5000,
            floor_price: 0.1,
            traits: vec![],
        });
    }
    
    fn init_reward_tiers(&mut self) {
        self.reward_tiers.insert("bronze".to_string(), NFRewardTier {
            tier: "bronze".to_string(),
            nft_collection: "lucky-cat".to_string(),
            bonus_percent: 5.0,
            rakeback_boost: 1.1,
            daily_free_spins: 10,
            vip_points_multiplier: 1.2,
        });
        
        self.reward_tiers.insert("gold".to_string(), NFRewardTier {
            tier: "gold".to_string(),
            nft_collection: "tiger-royal".to_string(),
            bonus_percent: 20.0,
            rakeback_boost: 1.5,
            daily_free_spins: 50,
            vip_points_multiplier: 2.0,
        });
        
        self.reward_tiers.insert("diamond".to_string(), NFRewardTier {
            tier: "diamond".to_string(),
            nft_collection: "tiger-royal".to_string(),
            bonus_percent: 50.0,
            rakeback_boost: 3.0,
            daily_free_spins: 0,
            vip_points_multiplier: 5.0,
        });
    }
    
    pub fn get_collections(&self) -> Vec<&NFTCollection> {
        self.collections.values().collect()
    }
    
    pub fn get_reward_tier(&self, tier: &str) -> Option<&NFRewardTier> {
        self.reward_tiers.get(tier)
    }
    
    pub fn calculate_nft_bonus(&self, _user_id: &str, bet_amount: f64) -> f64 {
        bet_amount * 0.10 // 10% base bonus
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_collections() {
        let manager = NFTManager::new();
        let collections = manager.get_collections();
        assert!(collections.len() > 0);
    }
}
