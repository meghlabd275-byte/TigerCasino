//! BetSolutions Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BetSolutionsProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BetSolutionsProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            GameInfo { id: "bs_crash_001".to_string(), name: "BetAviator".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/betaviator/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_crash_002".to_string(), name: "CrashX".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/crashx/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_crash_003".to_string(), name: "SpeedCrash".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/speedcrash/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_mines_001".to_string(), name: "Minesweeper".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Arcade, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/minesweeper/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_mines_002".to_string(), name: "Mines Deluxe".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Arcade, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/mines-deluxe/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_dice_001".to_string(), name: "Dice Duels".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Dice, rtp: 99.00, volatility: Volatility::Low, min_bet: 0.01, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/dice-duels/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_dice_002".to_string(), name: "Dice Race".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Dice, rtp: 99.00, volatility: Volatility::Low, min_bet: 0.01, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/dice-race/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_keno_001".to_string(), name: "Virtual Keno".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/keno/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_lottery_001".to_string(), name: "Lucky Lotto".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/lucky-lotto/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bs_plinko_001".to_string(), name: "Plinko Star".to_string(), provider: "BetSolutions".to_string(), category: GameCategory::Crash, rtp: 98.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betsolutions.com/plinko-star/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BetSolutionsProvider {
    fn name(&self) -> &str { "BetSolutions" }
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
