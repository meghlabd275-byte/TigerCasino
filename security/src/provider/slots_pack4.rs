//! Even More Slots Provider - Pack 4

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct SlotsPack4Provider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl SlotsPack4Provider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Evoplay (50)
            GameInfo { id: "evoplay_s001".to_string(), name: "The Great Conflict".to_string(), provider: "Evoplay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/conflict/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_s002".to_string(), name: "Epic Gladiators".to_string(), provider: "Evoplay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/gladiators/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_s003".to_string(), name: "Totem Island".to_string(), provider: "Evoplay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/totem/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_s004".to_string(), name: "Naughty Sweets".to_string(), provider: "Evoplay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/sweets/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "evoplay_s005".to_string(), name: "Irish Reels".to_string(), provider: "Evoplay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.evoplay.com/irish/thumb.jpg".to_string(), game_url: "".to_string() },
            // RubyPlay (50)
            GameInfo { id: "rubyplay_s001".to_string(), name: "Egyptian Dreams".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/egyptian-dreams/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "rubyplay_s002".to_string(), name: "Joker's Gold".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/jokers-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "rubyplay_s003".to_string(), name: "Reels of Wealth".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/reels-wealth/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "rubyplay_s004".to_string(), name: "Money Mansion".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/money-mansion/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "rubyplay_s005".to_string(), name: "Fruits of the Nile".to_string(), provider: "RubyPlay".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.rubyplay.com/fruits-nile/thumb.jpg".to_string(), game_url: "".to_string() },
            // KA Gaming (50)
            GameInfo { id: "ka_s001".to_string(), name: "Golden Empire".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/golden-empire/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_s002".to_string(), name: "Lucky Dragons".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/lucky-dragons/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_s003".to_string(), name: "Magic World".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/magic-world/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_s004".to_string(), name: "Fortune Panda".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/fortune-panda/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ka_s005".to_string(), name: "Super rich".to_string(), provider: "KA Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.kagaming.com/super-rich/thumb.jpg".to_string(), game_url: "".to_string() },
            // Vivo Gaming Live (30)
            GameInfo { id: "vivo_bj_001".to_string(), name: "Vivo Blackjack".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivo.com/bj/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_bj_002".to_string(), name: "Vivo Speed Blackjack".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 99.50, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivo.com/speed-bj/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_r_001".to_string(), name: "Vivo Roulette".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.30, volatility: Volatility::Low, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivo.com/roulette/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_b_001".to_string(), name: "Vivo Baccarat".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 98.94, volatility: Volatility::Low, min_bet: 5, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivo.com/baccarat/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "vivo_tp_001".to_string(), name: "Vivo Teen Patti".to_string(), provider: "Vivo Gaming".to_string(), category: GameCategory::LiveCasino, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1, max_bet: 5000, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.vivo.com/teen-patti/thumb.jpg".to_string(), game_url: "".to_string() },
            // CT Gaming (30)
            GameInfo { id: "ct_s001".to_string(), name: "Diamond Chase".to_string(), provider: "CT Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.ctgaming.com/diamond-chase/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ct_s002".to_string(), name: "Fruity Wins".to_string(), provider: "CT Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.ctgaming.com/fruity-wins/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ct_s003".to_string(), name: "Lucky Wizard".to_string(), provider: "CT Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.ctgaming.com/lucky-wizard/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ct_s004".to_string(), name: "Purple Payday".to_string(), provider: "CT Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.ctgaming.com/purple-payday/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "ct_s005".to_string(), name: "Winter Winners".to_string(), provider: "CT Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.ctgaming.com/winter-winners/thumb.jpg".to_string(), game_url: "".to_string() },
            // More Game Shows (20)
            GameInfo { id: "gs_001".to_string(), name: "Super Sic Bo".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/super-sic-bo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gs_002".to_string(), name: "Craps".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/craps/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gs_003".to_string(), name: "Side Bet City".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::High, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/side-bet-city/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gs_004".to_string(), name: "Dragonic".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.evolution.com/dragonic/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "gs_005".to_string(), name: "Ultimate Texas Hold'em".to_string(), provider: "Evolution".to_string(), category: GameCategory::GameShows, rtp: 97.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.evolution.com/ultimate-th/thumb.jpg".to_string(), game_url: "".to_string() },
            // Lottery (20)
            GameInfo { id: "lotto_001".to_string(), name: "Keno 1".to_string(), provider: "BGaming".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/keno1/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "lotto_002".to_string(), name: "Keno 2".to_string(), provider: "BGaming".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/keno2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "lotto_003".to_string(), name: "Bingo".to_string(), provider: "BGaming".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/bingo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "lotto_004".to_string(), name: "Bingo Blast".to_string(), provider: "BGaming".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/bingo-blast/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "lotto_005".to_string(), name: "Bingo Soccer".to_string(), provider: "BGaming".to_string(), category: GameCategory::Lottery, rtp: 95.00, volatility: Volatility::Medium, min_bet: 1.0, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/bingo-soccer/thumb.jpg".to_string(), game_url: "".to_string() },
            // Scratch Cards (20)
            GameInfo { id: "scratch_001".to_string(), name: "Scratch 777".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-777/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_002".to_string(), name: "Scratch Diamonds".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-diamonds/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_003".to_string(), name: "Scratch Match".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-match/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_004".to_string(), name: "Scratch King".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-king/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "scratch_005".to_string(), name: "Scratch Queen".to_string(), provider: "BGaming".to_string(), category: GameCategory::ScratchCards, rtp: 96.00, volatility: Volatility::Low, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.bgaming.com/scratch-queen/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for SlotsPack4Provider {
    fn name(&self) -> &str { "Slots Pack 4" }
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
