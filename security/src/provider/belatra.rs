//! Belatra Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BelatraProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BelatraProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "belatra_001".to_string(), name: "Book of Zoo".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/book-zoo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_002".to_string(), name: "The Explorers".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/explorers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_003".to_string(), name: "Neon Tam".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/neon-tam/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_004".to_string(), name: "King of Slots".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/king-slots/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_005".to_string(), name: "Joker Explosion".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/joker-explosion/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_006".to_string(), name: "Catch the Gold".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/catch-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_007".to_string(), name: "Mighty Mustang".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/mighty-mustang/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_008".to_string(), name: "Rhino".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/rhino/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_009".to_string(), name: "Fruit Cocktail".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/fruit-cocktail/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "belatra_010".to_string(), name: "Frog War".to_string(), provider: "Belatra".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.belatra.com/frog-war/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BelatraProvider {
    fn name(&self) -> &str { "Belatra" }
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
