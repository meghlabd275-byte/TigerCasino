//! Endorphina Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct EndorphinaProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl EndorphinaProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "endorphina_001".to_string(), name: "Diamond Wild".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.05, max_bet: 45.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/diamond-wild/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_002".to_string(), name: "Safari".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 45.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/safari/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_003".to_string(), name: "Inferno Star".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/inferno-star/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_004".to_string(), name: "Football".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 45.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/football/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_005".to_string(), name: "Charming Lady".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 45.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/charming-lady/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_006".to_string(), name: "Blazing Star".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/blazing-star/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_007".to_string(), name: "Mines".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.endorphina.com/mines/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_008".to_string(), name: "Plinko".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Crash, rtp: 98.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.endorphina.com/plinko/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_009".to_string(), name: "Dice".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Dice, rtp: 99.00, volatility: Volatility::Low, min_bet: 0.01, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.endorphina.com/dice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_010".to_string(), name: "Keno".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.endorphina.com/keno/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for EndorphinaProvider {
    fn name(&self) -> &str { "Endorphina" }
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
