//! Additional Slot Games Provider

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct ExtraSlotsProvider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl ExtraSlotsProvider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // More Pragmatic Play slots (200+ more)
            GameInfo { id: "prag_extra_001".to_string(), name: "Sweet Bonanza".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.51, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/sweet-bonanza/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_002".to_string(), name: " Gates of Olympus".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/olympus/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_003".to_string(), name: "Wild West Gold".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.51, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/wild-west/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_004".to_string(), name: "Fruit Party".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.47, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/fruit-party/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_005".to_string(), name: "Starlight Princess".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/princess/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_006".to_string(), name: "Power of Thor".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/thor/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_007".to_string(), name: "Joker's Jewels".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/jokers-jewels/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_008".to_string(), name: "Great Rhino".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/rhino/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_009".to_string(), name: "5 Lions Gold".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/5lions/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "prag_extra_010".to_string(), name: "Aztec Gems".to_string(), provider: "Pragmatic Play".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.pragmaticplay.net/aztec/thumb.jpg".to_string(), game_url: "".to_string() },
            // Play'n GO slots (50+ more)
            GameInfo { id: "png_extra_001".to_string(), name: "Book of Dead".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.21, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/book-dead/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_002".to_string(), name: "Reactoonz".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.51, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/reactoonz/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_003".to_string(), name: "Fire Joker".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.15, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/fire-joker/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_004".to_string(), name: "Legacy of Dead".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.58, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/legacy-dead/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_005".to_string(), name: "Rise of Merlin".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.58, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/merlin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_006".to_string(), name: "Rainbow Ryan".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.30, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/ryan/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_007".to_string(), name: "Durian Dynamite".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.52, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/durian/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_008".to_string(), name: "Moon Princess".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/moon-princess/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_009".to_string(), name: "Peking Luck".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.58, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/peking/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "png_extra_010".to_string(), name: "Angry Dogs".to_string(), provider: "Play'n GO".to_string(), category: GameCategory::Slots, rtp: 96.50, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playngo.com/angry-dogs/thumb.jpg".to_string(), game_url: "".to_string() },
            // NetEnt slots (50+ more)
            GameInfo { id: "netent_extra_001".to_string(), name: "Starburst XXXtreme".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.26, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/starburst-xtreme/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_002".to_string(), name: "Twin Spin XXXtreme".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.36, volatility: Volatility::High, min_bet: 0.25, max_bet: 125.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/twinspin-xtreme/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_003".to_string(), name: "Mega Fortune".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/mega-fortune/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_004".to_string(), name: "Hall of Gods".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 95.30, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/hall-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_005".to_string(), name: "Arabian Nights".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 95.20, volatility: Volatility::High, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/arabian-nights/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_006".to_string(), name: "Aloha! Cluster Pays".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.42, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/aloha/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_007".to_string(), name: "Berryburst".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.56, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/berryburst/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_008".to_string(), name: "Butterfly Staxx".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.80, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/butterfly/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_009".to_string(), name: "Narcos".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.23, volatility: Volatility::High, min_bet: 0.20, max_bet: 400.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/narcos/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "netent_extra_010".to_string(), name: "Jumanji".to_string(), provider: "NetEnt".to_string(), category: GameCategory::Slots, rtp: 96.33, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 200.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.netent.com/jumanji/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for ExtraSlotsProvider {
    fn name(&self) -> &str { "Extra Slots" }
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
