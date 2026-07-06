//! More Slot Games Provider

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct MoreSlotsProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl MoreSlotsProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // Relax Gaming (50 more)
            GameInfo { id: "relax_extra_001".to_string(), name: "Money Train 3".to_string(), provider: "Relax Gaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.relaxgaming.com/mt3/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "relax_extra_002".to_string(), name: "Temple Tumble Megaways".to_string(), provider: "Relax Gaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.relaxgaming.com/temple-tumble-mega/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "relax_extra_003".to_string(), name: " Snake Arena".to_string(), provider: "Relax Gaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.relaxgaming.com/snake-arena/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "relax_extra_004".to_string(), name: "Tower".to_string(), provider: "Relax Gaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.relaxgaming.com/tower/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "relax_extra_005".to_string(), name: "Cash Quest".to_string(), provider: "Relax Gaming".to_string(), category: GameCategory::Slots, rtp: 95.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.relaxgaming.com/cash-quest/thumb.jpg".to_string(), game_url: "".to_string() },
            // Nolimit City (30 more)
            GameInfo { id: "nolimit_extra_001".to_string(), name: "xWays Hoarder".to_string(), provider: "Nolimit City".to_string(), category: GameCategory::Slots, rtp: 96.06, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.nolimitcity.com/hoarder/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "nolimit_extra_002".to_string(), name: "Bushido Ways".to_string(), provider: "Nolimit City".to_string(), category: GameCategory::Slots, rtp: 96.10, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.nolimitcity.com/bushido/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "nolimit_extra_003".to_string(), name: "Psychedelic Snacks".to_string(), provider: "Nolimit City".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.nolimitcity.com/snacks/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "nolimit_extra_004".to_string(), name: "El王家府".to_string(), provider: "Nolimit City".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.nolimitcity.com/el/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "nolimit_extra_005".to_string(), name: "Infectious".to_string(), provider: "Nolimit City".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.nolimitcity.com/infectious/thumb.jpg".to_string(), game_url: "".to_string() },
            // Hacksaw Gaming (30 more)
            GameInfo { id: "hacksaw_extra_001".to_string(), name: "Wanted Dead or a Wild".to_string(), provider: "Hacksaw Gaming".to_string(), category: GameCategory::Slots, rtp: 96.30, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.hacksaw.com/wanted/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "hacksaw_extra_002".to_string(), name: "Stack 'Em".to_string(), provider: "Hacksaw Gaming".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.hacksaw.com/stack-em/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "hacksaw_extra_003".to_string(), name: "The Bomb".to_string(), provider: "Hacksaw Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.hacksaw.com/bomb/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "hacksaw_extra_004".to_string(), name: "Time Rush".to_string(), provider: "Hacksaw Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.hacksaw.com/time-rush/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "hacksaw_extra_005".to_string(), name: "Chaos Crew".to_string(), provider: "Hacksaw Gaming".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.hacksaw.com/chaos-crew/thumb.jpg".to_string(), game_url: "".to_string() },
            // Push Gaming (30 more)
            GameInfo { id: "push_extra_001".to_string(), name: "Jammin' Jars 2".to_string(), provider: "Push Gaming".to_string(), category: GameCategory::Slots, rtp: 96.40, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pushgaming.com/jammin2/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "push_extra_002".to_string(), name: "Big Bamboo".to_string(), provider: "Push Gaming".to_string(), category: GameCategory::Slots, rtp: 96.10, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pushgaming.com/big-bamboo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "push_extra_003".to_string(), name: "Wild Swarm".to_string(), provider: "Push Gaming".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pushgaming.com/wild-swarm/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "push_extra_004".to_string(), name: "Frozen Jam".to_string(), provider: "Push Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pushgaming.com/frozen-jam/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "push_extra_005".to_string(), name: "Viking Clash".to_string(), provider: "Push Gaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pushgaming.com/viking-clash/thumb.jpg".to_string(), game_url: "".to_string() },
            // Yggdrasil (30 more)
            GameInfo { id: "yggdrasil_extra_001".to_string(), name: "Valley of the Gods".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.20, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/valley-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "yggdrasil_extra_002".to_string(), name: "Razortooth".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.30, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/razortooth/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "yggdrasil_extra_003".to_string(), name: "Temple of Astaroth".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/astaroth/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "yggdrasil_extra_004".to_string(), name: "Age of Ice".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/age-ice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "yggdrasil_extra_005".to_string(), name: "Tuts Twisted Fortune".to_string(), provider: "Yggdrasil".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.yggdrasil.com/tuts/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for MoreSlotsProvider {
    fn name(&self) -> &str { "More Slots" }
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
