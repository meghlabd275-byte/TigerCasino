//! Play'n GO Provider Integration
//! 
//! Integration with Play'n GO's game API.
//! Play'n GO is a leading slot game provider with innovative titles.

use super::*;
use serde_json::json;

/// Play'n GO API client
pub struct PlaynGoProvider {
    config: ProviderConfig,
    base_url: String,
}

impl PlaynGoProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // Play'n GO portfolio - known for innovative slots
        Ok(vec![
            // Book of Dead series
            GameInfo {
                id: "png_book_of_dead".to_string(),
                name: "Book of Dead".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.21,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/book-of-dead.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Reactoonz series
            GameInfo {
                id: "png_reactoonz".to_string(),
                name: "Reactoonz".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.51,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/reactoonz.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "png_reactoonz2".to_string(),
                name: "Reactoonz 2".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.20,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/reactoonz2.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Rich Wilde series
            GameInfo {
                id: "png_rich_wilde".to_string(),
                name: "Rich Wilde and the Shield of Athena".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.25,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/shield-of-athena.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Legacy series
            GameInfo {
                id: "pngLegacyOfDead".to_string(),
                name: "Legacy of Dead".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.58,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/legacy-of-dead.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Egyptian adventures
            GameInfo {
                id: "png_scrollOfDead".to_string(),
                name: "Scroll of Dead".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.28,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/scroll-of-dead.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Viking series
            GameInfo {
                id: "png_vikings".to_string(),
                name: "Vikings".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/vikings.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Moon Princess series
            GameInfo {
                id: "png_moon_princess".to_string(),
                name: "Moon Princess".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/moon-princess.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "png_moon_princess_100".to_string(),
                name: "Moon Princess 100".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.50,
                volatility: Volatility::High,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/moon-princess-100.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Cat Wilde series
            GameInfo {
                id: "png_cat_wilde".to_string(),
                name: "Cat Wilde and the Doom of Death".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.29,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/cat-wilde-doom.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Newer releases
            GameInfo {
                id: "png_bee_hive".to_string(),
                name: "Bee Hive".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.20,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.playngo.com/bee-hive.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "png_fire_joker".to_string(),
                name: "Fire Joker".to_string(),
                provider: "Play'n GO".to_string(),
                category: GameCategory::Slots,
                rtp: 96.15,
                volatility: Volatility::Medium,
                min_bet: 0.05,
                max_bet: 100.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.playngo.com/fire-joker.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for PlaynGoProvider {
    fn name(&self) -> &str {
        "Play'n GO"
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
