//! Kiron Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct KironProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl KironProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "kiron_vs_001".to_string(), name: "Virtual Football World Cup".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vf-worldcup/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_002".to_string(), name: "Virtual Football Champions Cup".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vf-champions/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_003".to_string(), name: "Virtual Football Pro".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vf-pro/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_004".to_string(), name: "Virtual Basketball Pro".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vb-pro/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_005".to_string(), name: "Virtual Tennis Pro".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vt-pro/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_006".to_string(), name: "Virtual Horse Racing".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vhr/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_007".to_string(), name: "Virtual Greyhound Racing".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vgr/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_008".to_string(), name: "Virtual Speedway".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vspeedway/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_009".to_string(), name: "Virtual Cricket".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vcricket/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "kiron_vs_010".to_string(), name: "Virtual Rugby".to_string(), provider: "Kiron".to_string(), category: GameCategory::VirtualSports, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.kiron.com/vrugby/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for KironProvider {
    fn name(&self) -> &str { "Kiron" }
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
