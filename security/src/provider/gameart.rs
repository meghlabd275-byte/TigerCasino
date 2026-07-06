//! GameArt Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct GameArtProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl GameArtProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "ga_001".to_string(), name: "Legend of the Nile".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/nile/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_002".to_string(), name: "Joker Poker".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/joker-poker/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_003".to_string(), name: "Lady Fire".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/lady-fire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_004".to_string(), name: "Wild Dolphin".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/dolphin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_005".to_string(), name: "Book of Sheba".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/sheba/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_006".to_string(), name: "Golden Tiger".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/golden-tiger/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_007".to_string(), name: "Mighty Dragon".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/mighty-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_008".to_string(), name: "Battle for Atlantis".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/atlantis/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_009".to_string(), name: "Wolf Moon".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/wolf-moon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ga_010".to_string(), name: "Maid Marian".to_string(), provider: "GameArt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.gameart.com/maid-marian/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for GameArtProvider {
    fn name(&self) -> &str { "GameArt" }
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
