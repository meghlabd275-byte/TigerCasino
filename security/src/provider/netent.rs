//! NetEnt Provider Integration
//! 
//! Integration with NetEnt's game API.
//! NetEnt is a pioneer in online casino gaming with many classic titles.

use super::*;
use serde_json::json;

/// NetEnt API client
pub struct NetEntProvider {
    config: ProviderConfig,
    base_url: String,
}

impl NetEntProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // NetEnt portfolio - classic and innovative games
        Ok(vec![
            // Starburst series
            GameInfo {
                id: "netent_starburst".to_string(),
                name: "Starburst".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.09,
                volatility: Volatility::Low,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.netent.com/starburst.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "netent_starburst_xxx".to_string(),
                name: "Starburst XXXtreme".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.45,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/starburst-xxtreme.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Gonzo's Quest
            GameInfo {
                id: "netent_gonzo".to_string(),
                name: "Gonzo's Quest".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 95.97,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 50.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/gonzos-quest.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "netent_gonzo_quest".to_string(),
                name: "Gonzo's Quest Megaways".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/gonzos-quest-megaways.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Dead or Alive
            GameInfo {
                id: "netent_dead_alive".to_string(),
                name: "Dead or Alive".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.82,
                volatility: Volatility::High,
                min_bet: 0.09,
                max_bet: 18.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/dead-or-alive.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "netent_dead_alive2".to_string(),
                name: "Dead or Alive 2".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.8,
                volatility: Volatility::High,
                min_bet: 0.09,
                max_bet: 18.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/dead-or-alive-2.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Twin Spin
            GameInfo {
                id: "netent_twinspin".to_string(),
                name: "Twin Spin".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.6,
                volatility: Volatility::Medium,
                min_bet: 0.25,
                max_bet: 125.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.netent.com/twin-spin.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Jack and the Beanstalk
            GameInfo {
                id: "netent_beanstalk".to_string(),
                name: "Jack and the Beanstalk".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.7,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 37.5,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/jack-beanstalk.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Mega Fortune
            GameInfo {
                id: "netent_mega_fortune".to_string(),
                name: "Mega Fortune".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.0,
                volatility: Volatility::High,
                min_bet: 0.25,
                max_bet: 62.5,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/mega-fortune.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Divine Fortune
            GameInfo {
                id: "netent_divine".to_string(),
                name: "Divine Fortune".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.59,
                volatility: Volatility::Medium,
                min_bet: 0.20,
                max_bet: 100.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/divine-fortune.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Fruit Shop
            GameInfo {
                id: "netent_fruit_shop".to_string(),
                name: "Fruit Shop".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.7,
                volatility: Volatility::Medium,
                min_bet: 0.01,
                max_bet: 40.0,
                has_free_spins: true,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.netent.com/fruit-shop.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Jimi Hendrix
            GameInfo {
                id: "netent_jimi".to_string(),
                name: "Jimi Hendrix".to_string(),
                provider: "NetEnt".to_string(),
                category: GameCategory::Slots,
                rtp: 96.9,
                volatility: Volatility::Low,
                min_bet: 0.20,
                max_bet: 40.0,
                has_free_spins: true,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.netent.com/jimi-hendrix.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for NetEntProvider {
    fn name(&self) -> &str {
        "NetEnt"
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
