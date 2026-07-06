//! Advanced Fraud Detection Module for TigerCasino
//!
//! Machine learning enhanced fraud detection with real-time pattern recognition

use serde::{Deserialize, Serialize};
use std::collections::{HashMap, VecDeque};
use std::sync::Arc;
use tokio::sync::RwLock;
use chrono::{DateTime, Utc, Duration};

// ============== Data Structures ==============

/// User behavioral profile
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BehavioralProfile {
    pub user_id: String,
    pub created_at: DateTime<Utc>,
    pub last_updated: DateTime<Utc>,
    
    // Betting patterns
    pub avg_bet_size: f64,
    pub bet_size_std_dev: f64,
    pub avg_session_duration: i64,  // seconds
    pub games_preferred: Vec<String>,
    pub betting_times: Vec<i32>,  // hour of day (0-23)
    
    // Statistical features
    pub win_rate: f64,
    pub total_bets: u64,
    pub total_wins: u64,
    pub biggest_win: f64,
    pub biggest_loss: f64,
    
    // Anomaly scores
    pub anomaly_score: f64,
    pub risk_score: f64,
}

/// Real-time betting session
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BettingSession {
    pub session_id: String,
    pub user_id: String,
    pub start_time: DateTime<Utc>,
    pub end_time: Option<DateTime<Utc>>,
    pub bets: Vec<BetRecord>,
    pub total_wagered: f64,
    pub total_won: f64,
    pub ip_address: String,
    pub user_agent: String,
    pub device_fingerprint: String,
}

/// Individual bet record
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BetRecord {
    pub bet_id: String,
    pub timestamp: DateTime<Utc>,
    pub game_id: String,
    pub bet_amount: f64,
    pub win_amount: f64,
    pub outcome: String,
    pub ip_address: String,
    pub session_id: String,
}

/// Risk rule evaluation result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RiskEvaluation {
    pub user_id: String,
    pub session_id: String,
    pub timestamp: DateTime<Utc>,
    pub risk_score: f64,
    pub risk_level: RiskLevel,
    pub triggered_rules: Vec<TriggeredRule>,
    pub recommended_action: RecommendedAction,
    pub details: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum RiskLevel {
    VeryLow,
    Low,
    Medium,
    High,
    Critical,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TriggeredRule {
    pub rule_id: String,
    pub rule_name: String,
    pub severity: ActivitySeverity,
    pub score_contribution: f64,
    pub description: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum RecommendedAction {
    Allow,
    Monitor,
    RequireVerification,
    Block,
    FreezeAccount,
}

// ============== Advanced Fraud Detector ==============

/// Advanced fraud detection engine with ML-like features
pub struct AdvancedFraudDetector {
    config: FraudConfig,
    user_profiles: Arc<RwLock<HashMap<String, BehavioralProfile>>>,
    sessions: Arc<RwLock<HashMap<String, BettingSession>>>,
    recent_bets: Arc<RwLock<VecDeque<BetRecord>>>,
    suspicious_activities: Arc<RwLock<VecDeque<SuspiciousActivity>>>,
    
    // ML-like thresholds
    bet_amount_zscore_threshold: f64,
    win_rate_zscore_threshold: f64,
    session_count_threshold: u32,
    rapid_bet_threshold: u32,  // bets per minute
}

impl AdvancedFraudDetector {
    pub fn new(config: FraudConfig) -> Self {
        Self {
            config,
            user_profiles: Arc::new(RwLock::new(HashMap::new())),
            sessions: Arc::new(RwLock::new(HashMap::new())),
            recent_bets: Arc::new(RwLock::new(VecDeque::new())),
            suspicious_activities: Arc::new(RwLock::new(VecDeque::new())),
            bet_amount_zscore_threshold: 3.0,
            win_rate_zscore_threshold: 2.5,
            session_count_threshold: 5,
            rapid_bet_threshold: 10,
        }
    }

    /// Initialize or update user behavioral profile
    pub async fn update_profile(&self, user_id: &str, bet: &BetRecord) {
        let mut profiles = self.user_profiles.write().await;
        
        let profile = profiles.entry(user_id.to_string())
            .or_insert_with(|| BehavioralProfile {
                user_id: user_id.to_string(),
                created_at: Utc::now(),
                last_updated: Utc::now(),
                avg_bet_size: 0.0,
                bet_size_std_dev: 0.0,
                avg_session_duration: 0,
                games_preferred: Vec::new(),
                betting_times: Vec::new(),
                win_rate: 0.0,
                total_bets: 0,
                total_wins: 0,
                biggest_win: 0.0,
                biggest_loss: 0.0,
                anomaly_score: 0.0,
                risk_score: 0.0,
            });
        
        // Update statistics
        let total_bets = profile.total_bets + 1;
        let total_wins = if bet.win_amount > bet.bet_amount { profile.total_wins + 1 } else { profile.total_wins };
        
        // Update average bet size
        profile.avg_bet_size = ((profile.avg_bet_size * profile.total_bets as f64) + bet.bet_amount) / total_bets as f64;
        
        // Update win rate
        profile.win_rate = total_wins as f64 / total_bets as f64;
        
        // Update biggest win/loss
        if bet.win_amount > profile.biggest_win {
            profile.biggest_win = bet.win_amount;
        }
        if bet.bet_amount > profile.biggest_loss {
            profile.biggest_loss = bet.bet_amount;
        }
        
        profile.total_bets = total_bets;
        profile.total_wins = total_wins;
        profile.last_updated = Utc::now();
        
        // Update games preferred
        if !profile.games_preferred.contains(&bet.game_id) {
            profile.games_preferred.push(bet.game_id.clone());
        }
        
        // Update betting times
        let hour = bet.timestamp.format("%H").to_string().parse::<i32>().unwrap_or(0);
        if !profile.betting_times.contains(&hour) {
            profile.betting_times.push(hour);
        }
    }

    /// Evaluate risk for a new bet
    pub async fn evaluate_risk(&self, user_id: &str, session_id: &str, bet: &BetRecord) -> RiskEvaluation {
        let mut triggered_rules = Vec::new();
        let mut total_score = 0.0;
        
        let profiles = self.user_profiles.read().await;
        let profile = profiles.get(user_id);
        
        // Rule 1: Bet amount anomaly (Z-score based)
        if let Some(p) = profile {
            if p.total_bets > 10 {
                let z_score = (bet.bet_amount - p.avg_bet_size) / p.bet_size_std_dev.max(1.0);
                if z_score.abs() > self.bet_amount_zscore_threshold {
                    triggered_rules.push(TriggeredRule {
                        rule_id: "BET_ANOMALY".to_string(),
                        rule_name: "Unusual Bet Amount".to_string(),
                        severity: ActivitySeverity::High,
                        score_contribution: 0.3,
                        description: format!("Bet amount {} deviates significantly from user's average", bet.bet_amount),
                    });
                    total_score += 0.3;
                }
            }
        }
        
        // Rule 2: Win rate anomaly
        if let Some(p) = profile {
            if p.total_bets > 20 {
                let expected_win_rate = 0.5; // Expected for fair games
                let z_score = (p.win_rate - expected_win_rate) / 0.2; // Approximate std dev
                if z_score.abs() > self.win_rate_zscore_threshold {
                    triggered_rules.push(TriggeredRule {
                        rule_id: "WIN_RATE_ANOMALY".to_string(),
                        rule_name: "Suspicious Win Rate".to_string(),
                        severity: ActivitySeverity::High,
                        score_contribution: 0.4,
                        description: format!("Win rate {} is statistically anomalous", p.win_rate),
                    });
                    total_score += 0.4;
                }
            }
        }
        
        // Rule 3: Rapid betting detection
        let recent_bets = self.recent_bets.read().await;
        let recent_user_bets: Vec<_> = recent_bets.iter()
            .filter(|b| b.user_id == user_id)
            .filter(|b| b.timestamp > Utc::now() - Duration::seconds(60))
            .collect();
        
        if recent_user_bets.len() > self.rapid_bet_threshold {
            triggered_rules.push(TriggeredRule {
                rule_id: "RAPID_BETTING".to_string(),
                rule_name: "Rapid Betting Detected".to_string(),
                severity: ActivitySeverity::Medium,
                score_contribution: 0.25,
                description: format!("{} bets in the last minute", recent_user_bets.len()),
            });
            total_score += 0.25;
        }
        
        // Rule 4: Large bet relative to balance (requires balance check)
        // This would need balance info from the caller
        
        // Determine risk level
        let risk_level = match total_score {
            s if s >= 0.8 => RiskLevel::Critical,
            s if s >= 0.6 => RiskLevel::High,
            s if s >= 0.4 => RiskLevel::Medium,
            s if s >= 0.2 => RiskLevel::Low,
            _ => RiskLevel::VeryLow,
        };
        
        // Determine recommended action
        let recommended_action = match risk_level {
            RiskLevel::Critical => RecommendedAction::Block,
            RiskLevel::High => RecommendedAction::RequireVerification,
            RiskLevel::Medium => RecommendedAction::Monitor,
            _ => RecommendedAction::Allow,
        };
        
        RiskEvaluation {
            user_id: user_id.to_string(),
            session_id: session_id.to_string(),
            timestamp: Utc::now(),
            risk_score: total_score,
            risk_level: risk_level.clone(),
            triggered_rules,
            recommended_action,
            details: HashMap::new(),
        }
    }

    /// Record a bet for analysis
    pub async fn record_bet(&self, bet: BetRecord) {
        // Update profile
        self.update_profile(&bet.user_id, &bet).await;
        
        // Add to recent bets
        let mut recent_bets = self.recent_bets.write().await;
        recent_bets.push_back(bet);
        
        // Keep only last 10000 bets
        if recent_bets.len() > 10000 {
            recent_bets.pop_front();
        }
    }

    /// Start a new betting session
    pub async fn start_session(&self, session: BettingSession) {
        let mut sessions = self.sessions.write().await;
        sessions.insert(session.session_id.clone(), session);
    }

    /// End a betting session
    pub async fn end_session(&self, session_id: &str) {
        let mut sessions = self.sessions.write().await;
        if let Some(session) = sessions.get_mut(session_id) {
            session.end_time = Some(Utc::now());
        }
    }

    /// Get user profile
    pub async fn get_profile(&self, user_id: &str) -> Option<BehavioralProfile> {
        let profiles = self.user_profiles.read().await;
        profiles.get(user_id).cloned()
    }

    /// Get suspicious activities
    pub async fn get_suspicious_activities(&self, limit: usize) -> Vec<SuspiciousActivity> {
        let activities = self.suspicious_activities.read().await;
        activities.iter().rev().take(limit).cloned().collect()
    }
}

// ============== Rate Limiter ==============

/// Token bucket rate limiter
pub struct RateLimiter {
    capacity: u32,
    refill_rate: f64,  // tokens per second
    tokens: f64,
    last_refill: DateTime<Utc>,
    mutex: tokio::sync::Mutex<()>,
}

impl RateLimiter {
    pub fn new(capacity: u32, refill_rate: f64) -> Self {
        Self {
            capacity,
            refill_rate,
            tokens: capacity as f64,
            last_refill: Utc::now(),
            mutex: tokio::sync::Mutex::new(()),
        }
    }

    /// Try to consume a token
    pub async fn try_acquire(&self) -> bool {
        let _guard = self.mutex.lock().await;
        
        let now = Utc::now();
        let elapsed = (now - self.last_refill).num_milliseconds() as f64 / 1000.0;
        
        // Refill tokens
        self.tokens = (self.tokens + elapsed * self.refill_rate).min(self.capacity as f64);
        self.last_refill = now;
        
        if self.tokens >= 1.0 {
            self.tokens -= 1.0;
            true
        } else {
            false
        }
    }

    /// Get remaining tokens
    pub async fn remaining(&self) -> u32 {
        let _guard = self.mutex.lock().await;
        
        let now = Utc::now();
        let elapsed = (now - self.last_refill).num_milliseconds() as f64 / 1000.0;
        let tokens = (self.tokens + elapsed * self.refill_rate).min(self.capacity as f64);
        
        tokens.floor() as u32
    }
}

// ============== IP Reputation ==============

/// IP reputation data
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IPReputation {
    pub ip_address: String,
    pub score: f64,  // 0-100, lower is worse
    pub is_vpn: bool,
    pub is_proxy: bool,
    pub is_tor: bool,
    pub is_datacenter: bool,
    pub country: String,
    pub isp: String,
    pub last_checked: DateTime<Utc>,
    pub reports: u32,
}

/// IP reputation checker
pub struct IPReputationChecker {
    cache: Arc<RwLock<HashMap<String, IPReputation>>>,
}

impl IPReputationChecker {
    pub fn new() -> Self {
        Self {
            cache: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    /// Check IP reputation
    pub async fn check(&self, ip_address: &str) -> IPReputation {
        // Check cache first
        {
            let cache = self.cache.read().await;
            if let Some(rep) = cache.get(ip_address) {
                // Return cached if less than 1 hour old
                if Utc::now() - rep.last_checked < Duration::hours(1) {
                    return rep.clone();
                }
            }
        }
        
        // In production, this would call an IP reputation API
        // For now, return a default reputation
        let reputation = IPReputation {
            ip_address: ip_address.to_string(),
            score: 80.0,
            is_vpn: false,
            is_proxy: false,
            is_tor: false,
            is_datacenter: false,
            country: "Unknown".to_string(),
            isp: "Unknown".to_string(),
            last_checked: Utc::now(),
            reports: 0,
        };
        
        // Cache the result
        {
            let mut cache = self.cache.write().await;
            cache.insert(ip_address.to_string(), reputation.clone());
        }
        
        reputation
    }

    /// Get cached reputation
    pub async fn get_cached(&self, ip_address: &str) -> Option<IPReputation> {
        let cache = self.cache.read().await;
        cache.get(ip_address).cloned()
    }
}

impl Default for IPReputationChecker {
    fn default() -> Self {
        Self::new()
    }
}
