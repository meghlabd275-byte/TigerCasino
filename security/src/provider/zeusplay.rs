//! ZeusPlay Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct ZeusPlayProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl ZeusPlayProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "zeusplay_001".to_string(), name: "Frog Treasure".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/frog-treasure/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_002".to_string(), name: "Joker Jackpot".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/joker-jackpot/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_003".to_string(), name: "Zeus".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/zeus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_004".to_string(), name: "Viking Gold".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/viking-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_005".to_string(), name: "Egyptian Dreams".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/egyptian-dreams/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_006".to_string(), name: "Wolf Land".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/wolf-land/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_007".to_string(), name: "Dragon's Realm".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/dragons-realm/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_008".to_string(), name: "Neon Fruits".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/neon-fruits/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_009".to_string(), name: "Pirates".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/pirates/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "zeusplay_010".to_string(), name: "Wizard".to_string(), provider: "ZeusPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.zeusplay.com/wizard/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for ZeusPlayProvider {
    fn name(&self) -> &str { "ZeusPlay" }
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> { self.fetch_games() }
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        Ok(LaunchGameResponse { game_url: format!("{}/game/{}", self.base_url, request.game_id), session_id: uuid::Uuid::new_v4().to_string(), token: "session_token".to_string(), expires_at: Utc::now().timestamp() + 3600 })
    }
    fn process_transaction(&self, _request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        Ok(TransactionResult { transaction_id: uuid::Uuid::new_v4().to_string(), status: TransactionStatus::Completed, amount: 0.0, balance_after: 0.0, game_round_id: "round_id".to_string(), timestamp: Utc::now().timestamp() })
    }
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError> {
        let games = self.fetch_games()?;
        games.into_iter().find(|g| g.id == game_id).ok_or_else(|| ProviderError::GameNotFound(game_id.to_string()))
    }
    fn is_available(&self) -> bool { self.config.enabled }
}
