//! Fantasma Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct FantasmaProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl FantasmaProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "fantasma_001".to_string(), name: "Flower Fortunes".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/flower-fortunes/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_002".to_string(), name: "Fortune Girl".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/fortune-girl/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_003".to_string(), name: "Royal Dragon".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/royal-dragon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_004".to_string(), name: "Mighty Stallion".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/mighty-stallion/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_005".to_string(), name: "Leprechaun".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/leprechaun/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_006".to_string(), name: "Golden Empire".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/golden-empire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_007".to_string(), name: "Wild Jack".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/wild-jack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_008".to_string(), name: "Cleopatra".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/cleopatra/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_009".to_string(), name: "Phoenix".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/phoenix/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "fantasma_010".to_string(), name: "Pirate's Treasure".to_string(), provider: "Fantasma".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.fantasma.com/pirates-treasure/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for FantasmaProvider {
    fn name(&self) -> &str { "Fantasma" }
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
