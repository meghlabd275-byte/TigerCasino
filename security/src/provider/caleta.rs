//! Caleta Provider Integration (Lottery, Scratch Cards, Bingo)

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct CaletaProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl CaletaProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Scratch Cards
            GameInfo { id: "caleta_scratch_001".to_string(), name: "Scratch King".to_string(), provider: "Caleta".to_string(), category: GameCategory::ScratchCards, rtp: 95.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/scratch-king/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_scratch_002".to_string(), name: "Lucky Scratch".to_string(), provider: "Caleta".to_string(), category: GameCategory::ScratchCards, rtp: 95.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/lucky-scratch/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_scratch_003".to_string(), name: "Diamond Scratch".to_string(), provider: "Caleta".to_string(), category: GameCategory::ScratchCards, rtp: 95.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/diamond-scratch/thumb.jpg".to_string(), game_url: "".to_string() },
            // Lottery
            GameInfo { id: "caleta_lottery_001".to_string(), name: "Keno 10".to_string(), provider: "Caleta".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/keno10/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_lottery_002".to_string(), name: "Power Keno".to_string(), provider: "Caleta".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/power-keno/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_lottery_003".to_string(), name: "Bingo 75".to_string(), provider: "Caleta".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/bingo75/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_lottery_004".to_string(), name: "Bingo 90".to_string(), provider: "Caleta".to_string(), category: GameCategory::Bingo, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/bingo90/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_lottery_005".to_string(), name: "Pick 3".to_string(), provider: "Caleta".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/pick3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_lottery_006".to_string(), name: "Pick 4".to_string(), provider: "Caleta".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/pick4/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "caleta_lottery_007".to_string(), name: "Instant Lottery".to_string(), provider: "Caleta".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.caleta.com/instant-lottery/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for CaletaProvider {
    fn name(&self) -> &str { "Caleta" }
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
