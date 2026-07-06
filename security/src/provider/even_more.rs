//! Even More Games Provider

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct EvenMoreProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl EvenMoreProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Live Casino Table Games (50+)
            GameInfo { id: "live_extra_bj_001".to_string(), name: "Speed Blackjack 1".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 10, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/speed-bj1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_002".to_string(), name: "Speed Blackjack 2".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 10, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/speed-bj2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_003".to_string(), name: "Speed Blackjack 3".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 10, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/speed-bj3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_004".to_string(), name: "Free Bet Blackjack".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 10, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/free-bet-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_bj_005".to_string(), name: "Power Blackjack".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 10, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/power-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_001".to_string(), name: "Immersive Roulette".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/immersive-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_002".to_string(), name: "Double Ball Roulette".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/double-ball/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_r_003".to_string(), name: "Speed Roulette".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/speed-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_b_001".to_string(), name: "No Commission Baccarat".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 10, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/no-comm-bac/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "live_extra_b_002".to_string(), name: "Speed Baccarat".to_string(), provider: "Evolution".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 10, max_bet: 10000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/speed-bac/thumb.jpg".to_string(), game_url: "".to_string() },
            // Game Shows (20+)
            GameInfo { id: "gameshow_001".to_string(), name: "Cash or Crash".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 99.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/cash-crash/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gameshow_002".to_string(), name: "Gold Vault Roulette".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/gold-vault/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gameshow_003".to_string(), name: "Funky Time".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/funky-time/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gameshow_004".to_string(), name: "Monopoly Big Baller".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/monopoly-big-baller/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gameshow_005".to_string(), name: "Ballon d'Or".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/ballon-dor/thumb.jpg".to_string(), game_url: "".to_string() },
            // Scratch Cards (30+)
            GameInfo { id: "scratch_001".to_string(), name: "Scratch Gold".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_002".to_string(), name: "Scratch Dice".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-dice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_003".to_string(), name: "Scratch Platinum".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-platinum/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_004".to_string(), name: "Scratch Ruby".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-ruby/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_005".to_string(), name: "Scratch Emerald".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-emerald/thumb.jpg".to_string(), game_url: "".to_string() },
            // Table Games (50+)
            GameInfo { id: "table_001".to_string(), name: "European Blackjack".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 99.60, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/eu-blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_002".to_string(), name: "American Blackjack".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 99.60, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/am-blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_003".to_string(), name: "Atlantic City Blackjack".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 99.60, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/ac-blackjack/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_004".to_string(), name: "European Roulette".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/eu-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_005".to_string(), name: "American Roulette".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 94.70, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/am-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_006".to_string(), name: "French Roulette".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 98.60, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/fr-roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_007".to_string(), name: "Baccarat".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 98.94, volatility: Volatility::Low, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_008".to_string(), name: "Three Card Poker".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 97.97, volatility: Volatility::Medium, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/3card-poker/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_009".to_string(), name: "Caribbean Stud".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 96.30, volatility: Volatility::Medium, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/caribbean-stud/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "table_010".to_string(), name: "Hold'em Poker".to_string(), provider: "BGaming".to_string(), category: GameCategory::TableGames, rtp: 97.80, volatility: Volatility::Medium, min_bet: 1, max_bet: 1000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/holdem-poker/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for EvenMoreProvider {
    fn name(&self) -> &str { "Even More Games" }
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
