//! Patagonia Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct PatagoniaProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl PatagoniaProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "patagonia_001".to_string(), name: "Wild Andes".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/wild-andes/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_002".to_string(), name: "Wolf Legend".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/wolf-legend/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_003".to_string(), name: "Candy Palace".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/candy-palace/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_004".to_string(), name: "Magician Secrets".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/magician-secrets/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_005".to_string(), name: "Aztec Gold".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/aztec-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_006".to_string(), name: "Mighty Buffalo".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/mighty-buffalo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_007".to_string(), name: "Safari Adventure".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/safari-adventure/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_008".to_string(), name: "Golden Egypt".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/golden-egypt/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_009".to_string(), name: "Pirate Queen".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/pirate-queen/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "patagonia_010".to_string(), name: "Dragon Kingdom".to_string(), provider: "Patagonia".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.patagonia.com/dragon-kingdom/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for PatagoniaProvider {
    fn name(&self) -> &str { "Patagonia" }
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
