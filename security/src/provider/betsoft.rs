//! Betsoft Provider Integration
//! 
//! Integration with Betsoft's game API.

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BetsoftProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BetsoftProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "betsoft_001".to_string(), name: "Slots Angels".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.02, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/slots-angels/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_002".to_string(), name: "Aztec Treasure".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.30, max_bet: 120.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/aztec-treasure/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_003".to_string(), name: "Greedy Goblins".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::High, min_bet: 0.50, max_bet: 150.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/greedy-goblins/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_004".to_string(), name: "Madder Scientist".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::High, min_bet: 0.02, max_bet: 45.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/madder-scientist/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_005".to_string(), name: "At The Copa".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/at-the-copa/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_006".to_string(), name: "The Slot Father".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.25, max_bet: 150.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/slot-father/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_007".to_string(), name: "It Came From Venus".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.30, max_bet: 150.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/it-came-venus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_008".to_string(), name: "Gypsy Rose".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.02, max_bet: 150.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/gypsy-rose/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_009".to_string(), name: "Paco and the Popping Peppers".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.02, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/paco-peppers/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_010".to_string(), name: "Mighty Zeus".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.01, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/mighty-zeus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_011".to_string(), name: "Event Horizon".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/event-horizon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betsoft_012".to_string(), name: "Viking Voyage".to_string(), provider: "Betsoft".to_string(), category: GameCategory::Slots, rtp: 97.00, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.betsoft.com/viking-voyage/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BetsoftProvider {
    fn name(&self) -> &str { "Betsoft" }
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
