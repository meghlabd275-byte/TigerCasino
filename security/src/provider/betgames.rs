//! BetGames Provider Integration

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct BetGamesProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl BetGamesProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Lottery/Game Shows
            GameInfo { id: "betgames_001".to_string(), name: "BetGames Lucky 5".to_string(), provider: "BetGames".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/lucky5/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_002".to_string(), name: "BetGames Lucky 6".to_string(), provider: "BetGames".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/lucky6/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_003".to_string(), name: "BetGames Lucky 7".to_string(), provider: "BetGames".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/lucky7/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_004".to_string(), name: "BetGames Wheel".to_string(), provider: "BetGames".to_string(), category: GameCategory::GameShows, rtp: 95.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.betgames.com/wheel/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_005".to_string(), name: "BetGames Dice".to_string(), provider: "BetGames".to_string(), category: GameCategory::Dice, rtp: 98.00, volatility: Volatility::Medium, min_bet: 0.50, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/dice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_006".to_string(), name: "BetGames War".to_string(), provider: "BetGames".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Low, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/war/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_007".to_string(), name: "BetGames Bet on Poker".to_string(), provider: "BetGames".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/poker/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_008".to_string(), name: "BetGames Teen Patti".to_string(), provider: "BetGames".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/teen-patti/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_009".to_string(), name: "BetGames Andar Bahar".to_string(), provider: "BetGames".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/andar-bahar/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "betgames_010".to_string(), name: "BetGames 32 Cards".to_string(), provider: "BetGames".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.betgames.com/32cards/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for BetGamesProvider {
    fn name(&self) -> &str { "BetGames" }
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
