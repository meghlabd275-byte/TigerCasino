//! BGaming Provider Integration
//! 
//! Integration with BGaming's game API.
//! BGaming specializes in crypto-friendly games with unique mechanics.

use super::*;
use serde_json::json;

/// BGaming API client
pub struct BGamingProvider {
    config: ProviderConfig,
    base_url: String,
}

impl BGamingProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // BGaming portfolio - known for crypto games and unique mechanics
        Ok(vec![
            // Popular slots
            GameInfo {
                id: "bg ElvisFrog".to_string(),
                name: "Elvis Frog in Vegas".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 25.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.bgaming.com/elvis-frog.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "bg_aztec".to_string(),
                name: "Aztec Magic".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Slots,
                rtp: 95.95,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 25.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.bgaming.com/aztec-magic.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "bg_west".to_string(),
                name: "West Town".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 40.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.bgaming.com/west-town.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "bg_bunny".to_string(),
                name: "Bunny Pop".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Slots,
                rtp: 96.19,
                volatility: Volatility::Medium,
                min_bet: 0.30,
                max_bet: 48.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.bgaming.com/bunny-pop.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Crash games - BGaming's specialty
            GameInfo {
                id: "bg_crash_aviator".to_string(),
                name: "Aviator".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Crash,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/aviator.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Dice games
            GameInfo {
                id: "bg_dice".to_string(),
                name: "Dice Dice".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Arcade,
                rtp: 99.0,
                volatility: Volatility::Low,
                min_bet: 0.10,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/dice-dice.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Plinko
            GameInfo {
                id: "bg_plinko".to_string(),
                name: "Plinko".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Arcade,
                rtp: 99.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/plinko.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Mines
            GameInfo {
                id: "bg_mines".to_string(),
                name: "Mines".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Arcade,
                rtp: 97.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/mines.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Keno
            GameInfo {
                id: "bg_keno".to_string(),
                name: "Instant Keno".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Lottery,
                rtp: 96.0,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/instant-keno.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Table games
            GameInfo {
                id: "bg_blackjack".to_string(),
                name: "Blackjack".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Blackjack,
                rtp: 99.5,
                volatility: Volatility::Low,
                min_bet: 1.0,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/blackjack.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "bg_roulette".to_string(),
                name: "European Roulette".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Roulette,
                rtp: 97.3,
                volatility: Volatility::Medium,
                min_bet: 1.0,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/european-roulette.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Video poker
            GameInfo {
                id: "bg_jacks".to_string(),
                name: "Jacks or Better".to_string(),
                provider: "BGaming".to_string(),
                category: GameCategory::Poker,
                rtp: 99.56,
                volatility: Volatility::Low,
                min_bet: 0.05,
                max_bet: 45.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.bgaming.com/jacks-or-better.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for BGamingProvider {
    fn name(&self) -> &str {
        "BGaming"
    }
    
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        self.fetch_games()
    }
    
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        Ok(LaunchGameResponse {
            game_url: format!("{}/game/{}", self.base_url, request.game_id),
            session_id: uuid::Uuid::new_v4().to_string(),
            token: "session_token".to_string(),
            expires_at: chrono::Utc::now().timestamp() + 3600,
        })
    }
    
    fn process_transaction(&self, request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        Ok(TransactionResult {
            transaction_id: uuid::Uuid::new_v4().to_string(),
                status: TransactionStatus::Completed,
            amount: request.amount,
            balance_after: 0.0,
            game_round_id: request.round_id,
            timestamp: chrono::Utc::now().timestamp(),
        })
    }
    
    fn get_game_info(&self, game_id: &str) -> Result<GameInfo, ProviderError> {
        let games = self.get_games()?;
        games.into_iter()
            .find(|g| g.id == game_id)
            .ok_or_else(|| ProviderError::GameNotFound(game_id.to_string()))
    }
    
    fn is_available(&self) -> bool {
        self.config.enabled
    }
}
