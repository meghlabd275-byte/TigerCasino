//! More Slots Provider - Pack 2

use super::*;
use reqwest::Client;
use chrono::Utc;

pub struct SlotsPack2Provider {
    client: Client,
    config: ProviderConfig,
    base_url: String,
}

impl SlotsPack2Provider {
    pub fn new(config: ProviderConfig) -> Self {
        Self { client: Client::new(), config, base_url: config.api_url.clone() }
    }
    
    pub fn fetch_games(&self) -> Result<Vec<GameInfo>, ProviderError> {
        Ok(vec![
            // BGaming Slots (50)
            GameInfo { id: "bgaming_s001".to_string(), name: "Elf Princess".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/elf-princess/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s002".to_string(), name: "Wild Melbourne".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/wild-melbourne/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s003".to_string(), name: "Lady Wolf Moon".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/lady-wolf-moon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s004".to_string(), name: "Mighty Buffalo".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/mighty-buffalo/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s005".to_string(), name: "Aztec Magic".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/aztec-magic/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s006".to_string(), name: "Bonanza Billion".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/bonanza-billion/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s007".to_string(), name: "Miss Cherry Fruits".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/miss-cherry/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s008".to_string(), name: "Dragon's Gold".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/dragons-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s009".to_string(), name: "Johnny Cash".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/johnny-cash/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s010".to_string(), name: "West Candy".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/west-candy/thumb.jpg".to_string(), game_url: "".to_string() },
            // More BGaming
            GameInfo { id: "bgaming_s011".to_string(), name: "Royal Coin".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/royal-coin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s012".to_string(), name: "Sweets".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/sweets/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s013".to_string(), name: "Cherry Bomb".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/cherry-bomb/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s014".to_string(), name: "Wolf Power".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/wolf-power/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "bgaming_s015".to_string(), name: "Platinum Lightning".to_string(), provider: "BGaming".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.bgaming.com/platinum-lightning/thumb.jpg".to_string(), game_url: "".to_string() },
            // Spribe (30)
            GameInfo { id: "spribe_s001".to_string(), name: "Aviator".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/aviator/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s002".to_string(), name: "Spaceman".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/spaceman/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s003".to_string(), name: "JetX".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/jetx/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s004".to_string(), name: "Lucky Jet".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/lucky-jet/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s005".to_string(), name: "Rocket Queen".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/rocket-queen/thumb.jpg".to_string(), game_url: "".to_string() },
            // More Spribe
            GameInfo { id: "spribe_s006".to_string(), name: "Space XY".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/space-xy/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s007".to_string(), name: "Zeppelin".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/zeppelin/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s008".to_string(), name: "Balloon".to_string(), provider: "Spribe".to_string(), category: GameCategory::Crash, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/balloon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s009".to_string(), name: "Dice".to_string(), provider: "Spribe".to_string(), category: GameCategory::Dice, rtp: 97.00, volatility: Volatility::Medium, min_bet: 0.10, max_bet: 1000.0, has_free_spins: false, has_bonus_game: false, thumbnail_url: "https://static.spribe.com/dice/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "spribe_s010".to_string(), name: "Mines".to_string(), provider: "Spribe".to_string(), category: GameCategory::Mines, rtp: 97.00, volatility: Volatility::High, min_bet: 0.10, max_bet: 100.0, has_free_spins: false, has_bonus_game: true, thumbnail_url: "https://static.spribe.com/mines/thumb.jpg".to_string(), game_url: "".to_string() },
            // Booongo (30)
            GameInfo { id: "booongo_s001".to_string(), name: "Wolf Night".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/wolf-night/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_s002".to_string(), name: "Lady Wolf".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/lady-wolf/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_s003".to_string(), name: "Magic Apple".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/magic-apple/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_s004".to_string(), name: "Book of Sun".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/book-sun/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "booongo_s005".to_string(), name: "Egyptian Gods".to_string(), provider: "Booongo".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.booongo.com/egyptian-gods/thumb.jpg".to_string(), game_url: "".to_string() },
            // Endorphina (20)
            GameInfo { id: "endorphina_s001".to_string(), name: "Slotomon".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/slotomon/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_s002".to_string(), name: "Cyber Hunter".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/cyber-hunter/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_s003".to_string(), name: "Lunapark".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/lunapark/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_s004".to_string(), name: "Mongol".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/mongol/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "endorphina_s005".to_string(), name: "2027 Hit".to_string(), provider: "Endorphina".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.endorphina.com/2027-hit/thumb.jpg".to_string(), game_url: "".to_string() },
            // Playson (30)
            GameInfo { id: "playson_s001".to_string(), name: "Solar Queen".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playson.com/solar-queen/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playson_s002".to_string(), name: "Book of Gold".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playson.com/book-gold/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playson_s003".to_string(), name: "Royal Coins".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playson.com/royal-coins/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playson_s004".to_string(), name: "Legend of Cleopatra".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::High, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playson.com/cleopatra/thumb.jpg".to_string(), game_url: "".to_string() },
            GameInfo { id: "playson_s005".to_string(), name: "Fortune House".to_string(), provider: "Playson".to_string(), category: GameCategory::Slots, rtp: 96.00, volatility: Volatility::Medium, min_bet: 0.20, max_bet: 100.0, has_free_spins: true, has_bonus_game: true, thumbnail_url: "https://static.playson.com/fortune-house/thumb.jpg".to_string(), game_url: "".to_string() },
        ])
    }
}

impl GameProvider for SlotsPack2Provider {
    fn name(&self) -> &str { "Slots Pack 2" }
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
