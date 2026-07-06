//! Evolution Gaming Provider Integration
//! 
//! Integration with Evolution Gaming's live casino API.
//! Evolution is the leading provider of live dealer games and game shows.

use super::*;
use reqwest::Client;
use serde_json::json;

/// Evolution Gaming API client
pub struct EvolutionProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl EvolutionProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self {
            client: Client::new(),
            config,
            base_url: config.api_url.clone(),
        }
    }
    
    /// Generate HMAC signature for Evolution API
    fn generate_signature(&self, method: &str, path: &str, body: &str, timestamp: i64) -> String {
        use hmac::{Hmac, Mac};
        use sha2::Sha256;
        
        type HmacSha256 = Hmac<Sha256>;
        
        let message = format!("{}\n{}\n{}\n{}", method, path, body, timestamp);
        
        let mut mac = HmacSha256::new_from_slice(self.config.secret_key.as_bytes())
            .expect("HMAC can take key of any size");
        mac.update(message.as_bytes());
        
        hex::encode(mac.finalize().into_bytes())
    }
    
    /// Get available tables from Evolution
    pub fn fetch_tables(&self) -> Result<Vec<GameInfo>, ProviderError> {
        // Evolution offers live dealer games and game shows
        // This is a representative sample of their portfolio
        Ok(vec![
            // Blackjack tables
            GameInfo {
                id: "evo_bj_001".to_string(),
                name: "Infinite Blackjack".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 99.47,
                volatility: Volatility::Low,
                min_bet: 1.0,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/infinite-blackjack.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_bj_002".to_string(),
                name: "Speed Blackjack".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 99.42,
                volatility: Volatility::Low,
                min_bet: 1.0,
                max_bet: 10000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/speed-blackjack.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_bj_003".to_string(),
                name: "Power Blackjack".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 99.47,
                volatility: Volatility::Low,
                min_bet: 25.0,
                max_bet: 10000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/power-blackjack.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Roulette tables
            GameInfo {
                id: "evo_roulette_001".to_string(),
                name: "Lightning Roulette".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 97.30,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/lightning-roulette.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_roulette_002".to_string(),
                name: "Immersive Roulette".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 97.30,
                volatility: Volatility::Medium,
                min_bet: 1.0,
                max_bet: 10000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/immersive-roulette.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_roulette_003".to_string(),
                name: "Speed Roulette".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 97.30,
                volatility: Volatility::Medium,
                min_bet: 1.0,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/speed-roulette.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Baccarat tables
            GameInfo {
                id: "evo_baccarat_001".to_string(),
                name: "Speed Baccarat".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 98.64,
                volatility: Volatility::Low,
                min_bet: 1.0,
                max_bet: 10000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/speed-baccarat.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_baccarat_002".to_string(),
                name: "Lightning Baccarat".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::LiveCasino,
                rtp: 98.64,
                volatility: Volatility::High,
                min_bet: 1.0,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/lightning-baccarat.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Game Shows
            GameInfo {
                id: "evo_show_001".to_string(),
                name: "Crazy Time".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 96.08,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/crazy-time.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_002".to_string(),
                name: "Monopoly Live".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 96.23,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/monopoly-live.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_003".to_string(),
                name: "Dream Catcher".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 96.58,
                volatility: Volatility::Medium,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/dream-catcher.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_004".to_string(),
                name: "Deal or No Deal".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 95.42,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/deal-or-no-deal.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_005".to_string(),
                name: "Mega Ball".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 95.50,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/mega-ball.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_006".to_string(),
                name: "Funky Time".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 95.99,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/funky-time.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_007".to_string(),
                name: "Cash or Crash".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 96.10,
                volatility: Volatility::High,
                min_bet: 0.10,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/cash-or-crash.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_show_008".to_string(),
                name: "Gold Vault Roulette".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::GameShows,
                rtp: 97.30,
                volatility: Volatility::Medium,
                min_bet: 1.0,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/gold-vault-roulette.jpg".to_string(),
                game_url: "".to_string(),
            },
            // Poker
            GameInfo {
                id: "evo_poker_001".to_string(),
                name: "Casino Hold'em".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::Poker,
                rtp: 97.84,
                volatility: Volatility::Medium,
                min_bet: 1.0,
                max_bet: 1000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/casino-holdem.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_poker_002".to_string(),
                name: "Texas Hold'em Bonus".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::Poker,
                rtp: 97.64,
                volatility: Volatility::Medium,
                min_bet: 5.0,
                max_bet: 5000.0,
                has_free_spins: false,
                has_bonus_game: true,
                thumbnail_url: "https://cdn.evolutiongaming.com/texas-holdem.jpg".to_string(),
                game_url: "".to_string(),
            },
            GameInfo {
                id: "evo_poker_003".to_string(),
                name: "Three Card Poker".to_string(),
                provider: "Evolution Gaming".to_string(),
                category: GameCategory::Poker,
                rtp: 98.31,
                volatility: Volatility::Low,
                min_bet: 1.0,
                max_bet: 2500.0,
                has_free_spins: false,
                has_bonus_game: false,
                thumbnail_url: "https://cdn.evolutiongaming.com/three-card-poker.jpg".to_string(),
                game_url: "".to_string(),
            },
        ])
    }
}

impl GameProvider for EvolutionProvider {
    fn name(&self) -> &str {
        "Evolution Gaming"
    }
    
    fn get_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        self.fetch_tables()
    }
    
    fn launch_game(&self, request: LaunchGameRequest) -> Result<LaunchGameResponse, ProviderError> {
        let timestamp = chrono::Utc::now().timestamp();
        let path = format!("/v1/tables/{}/session", request.game_id);
        
        let body = json!({
            "user_id": request.user_id,
            "currency": request.currency,
            "language": request.language,
            "stake": request.stake,
            "token": request.session_token,
        }).to_string();
        
        let signature = self.generate_signature("POST", &path, &body, timestamp);
        
        // In production, make actual API call
        Ok(LaunchGameResponse {
            game_url: format!("{}/{}", self.base_url, request.game_id),
            session_id: uuid::Uuid::new_v4().to_string(),
            token: signature,
            expires_at: timestamp + 7200,
        })
    }
    
    fn process_transaction(&self, request: TransactionRequest) -> Result<TransactionResult, ProviderError> {
        let timestamp = chrono::Utc::now().timestamp();
        let path = "/v1/transactions";
        
        let body = json!({
            "session_id": request.session_id,
            "game_id": request.game_id,
            "user_id": request.user_id,
            "type": request.transaction_type,
            "amount": request.amount,
            "round_id": request.round_id,
            "timestamp": timestamp,
        }).to_string();
        
        let _signature = self.generate_signature("POST", path, &body, timestamp);
        
        Ok(TransactionResult {
            transaction_id: uuid::Uuid::new_v4().to_string(),
            status: TransactionStatus::Completed,
            amount: request.amount,
            balance_after: 0.0,
            game_round_id: request.round_id,
            timestamp,
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
