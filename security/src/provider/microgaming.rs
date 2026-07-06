//! Microgaming Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct MicrogamingProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl MicrogamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "microgaming_001".to_string(), name: "Mega Moolah".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 88.12, volatility: Volatility::High, min_bet: 0.25, max_bet: 6.25, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/mega-moolah/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_002".to_string(), name: "Immortal Romance".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 96.86, volatility: Volatility::High, min_bet: 0.30, max_bet: 30.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/immortal-romance/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_003".to_string(), name: "Thunderstruck II".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 96.65, volatility: Volatility::High, min_bet: 0.30, max_bet: 15.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/thunderstruck2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_004".to_string(), name: "Game of Thrones".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 95.01, volatility: Volatility::Medium, min_bet: 0.30, max_bet: 30.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/game-of-thrones/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_005".to_string(), name: "Avalon".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/avalon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_006".to_string(), name: "Break da Bank".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 95.75, volatility: Volatility::High, min_bet: 0.25, max_bet: 25.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/break-da-bank/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_007".to_string(), name: "Terminator 2".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.30, max_bet: 30.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/terminator2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_008".to_string(), name: "Jurassic Park".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 96.67, volatility: Volatility::Medium, min_bet: 0.30, max_bet: 30.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/jurassic-park/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_009".to_string(), name: "Playboy".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.30, max_bet: 45.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/playboy/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "microgaming_010".to_string(), name: "Lucky Leprechaun".to_string(), provider: "Microgaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 50.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.microgaming.com/lucky-leprechaun/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for MicrogamingProvider {
    fn name(&self) -> &str { "Microgaming" }
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
