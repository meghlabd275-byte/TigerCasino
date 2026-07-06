//! Slingo Provider Integration (Bingo/Slot hybrid)

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct SlingoProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl SlingoProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "slingo_001".to_string(), name: "Slingo Original".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/original/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_002".to_string(), name: "Slingo Rainbow Riches".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/rainbow-riches/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_003".to_string(), name: "Slingo Deal or No Deal".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/deal-or-no-deal/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_004".to_string(), name: "Slingo Who Wants to be a Millionaire".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/millionaire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_005".to_string(), name: "Slingo Monopoly".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/monopoly/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_006".to_string(), name: "Slingo Fluffy Favourites".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/fluffy-fav/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_007".to_string(), name: "Slingo Berry Blast".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/berry-blast/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_008".to_string(), name: "Slingo Super Spin".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/super-spin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_009".to_string(), name: "Slingo Candy".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/candy/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "slingo_010".to_string(), name: "Slingo Adventure".to_string(), provider: "Slingo".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.slingo.com/adventure/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for SlingoProvider {
    fn name(&self) -> &str { "Slingo" }
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
