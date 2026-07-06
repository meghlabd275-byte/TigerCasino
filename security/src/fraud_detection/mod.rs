//! Fraud Detection Module for TigerCasino
//! 
//! Real-time fraud detection and prevention system

pub mod advanced;

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::time::{SystemTime, UNIX_EPOCH};

pub use advanced::*;

/// Types of suspicious activities
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum SuspiciousActivityType {
    MultipleAccountRegistration,
    RapidBetting,
    UnusualWinPattern,
    BonusAbuse,
    Collusion,
    AccountSharing,
    VPNProxyUsage,
    OddArbitrage,
    WithdrawalAnomaly,
    ChipDumping,
}

/// A suspicious activity record
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SuspiciousActivity {
    pub id: String,
    pub activity_type: SuspiciousActivityType,
    pub user_id: String,
    pub severity: ActivitySeverity,
    pub description: String,
    pub evidence: HashMap<String, serde_json::Value>,
    pub timestamp: u64,
    pub resolved: bool,
    pub resolution_note: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum ActivitySeverity {
    Low,
    Medium,
    High,
    Critical,
}

/// Fraud detection configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FraudConfig {
    pub max_accounts_per_ip: u32,
    pub max_bets_per_minute: u32,
    pub max_deposits_per_hour: u32,
    pub suspicious_win_rate_threshold: f64,
    pub min_bet_to_win_ratio: f64,
    pub enable_geo_check: bool,
    pub enable_vpn_check: bool,
    pub enable_device_fingerprint: bool,
}

impl Default for FraudConfig {
    fn default() -> Self {
        Self {
            max_accounts_per_ip: 3,
            max_bets_per_minute: 60,
            max_deposits_per_hour: 10,
            suspicious_win_rate_threshold: 0.75,
            min_bet_to_win_ratio: 0.1,
            enable_geo_check: true,
            enable_vpn_check: true,
            enable_device_fingerprint: true,
        }
    }
}

/// User risk profile
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RiskProfile {
    pub user_id: String,
    pub risk_score: u32, // 0-100
    pub is_flagged: bool,
    pub is_blocked: bool,
    pub recent_activities: Vec<String>,
    pub betting_pattern: BettingPattern,
    pub last_updated: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BettingPattern {
    pub total_bets: u64,
    pub total_wagered: f64,
    pub total_won: f64,
    pub win_rate: f64,
    pub avg_bet_size: f64,
    pub largest_bet: f64,
    pub bet_frequency_per_hour: f64,
    pub is_unusual: bool,
}

/// Main fraud detector
pub struct FraudDetector {
    config: FraudConfig,
    ip_tracking: HashMap<String, u32>,
    user_bet_counts: HashMap<String, u32>,
    recent_bets: Vec<BetRecord>,
    suspicious_activities: Vec<SuspiciousActivity>,
}

#[derive(Debug, Clone)]
pub struct BetRecord {
    pub user_id: String,
    pub amount: f64,
    pub win_amount: f64,
    pub timestamp: u64,
}

impl FraudDetector {
    pub fn new(config: FraudConfig) -> Self {
        Self {
            config,
            ip_tracking: HashMap::new(),
            user_bet_counts: HashMap::new(),
            recent_bets: Vec::new(),
            suspicious_activities: Vec::new(),
        }
    }
    
    /// Check if a new account registration is suspicious
    pub fn check_registration(&mut self, ip_address: &str) -> bool {
        let count = self.ip_tracking.entry(ip_address.to_string()).or_insert(0);
        *count += 1;
        
        if *count > self.config.max_accounts_per_ip {
            self.record_activity(
                SuspiciousActivityType::MultipleAccountRegistration,
                "unknown".to_string(),
                ActivitySeverity::High,
                format!("IP {} has {} accounts", ip_address, count),
            );
            return true;
        }
        
        false
    }
    
    /// Check if a bet is suspicious
    pub fn check_bet(&mut self, user_id: &str, amount: f64, ip_address: &str) -> Option<ActivitySeverity> {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        // Track bet count
        let count = self.user_bet_counts.entry(user_id.to_string()).or_insert(0);
        *count += 1;
        
        // Check for rapid betting
        if *count > self.config.max_bets_per_minute as u32 {
            self.record_activity(
                SuspiciousActivityType::RapidBetting,
                user_id.to_string(),
                ActivitySeverity::Medium,
                format!("User {} placed {} bets in a minute", user_id, count),
            );
            return Some(ActivitySeverity::Medium);
        }
        
        // Record bet for pattern analysis
        self.recent_bets.push(BetRecord {
            user_id: user_id.to_string(),
            amount,
            win_amount: 0.0,
            timestamp,
        });
        
        // Clean old bets (older than 1 hour)
        self.recent_bets.retain(|b| timestamp - b.timestamp < 3600);
        
        None
    }
    
    /// Analyze betting patterns for a user
    pub fn analyze_betting_pattern(&self, user_id: &str) -> BettingPattern {
        let user_bets: Vec<&BetRecord> = self.recent_bets
            .iter()
            .filter(|b| b.user_id == user_id)
            .collect();
        
        if user_bets.is_empty() {
            return BettingPattern {
                total_bets: 0,
                total_wagered: 0.0,
                total_won: 0.0,
                win_rate: 0.0,
                avg_bet_size: 0.0,
                largest_bet: 0.0,
                bet_frequency_per_hour: 0.0,
                is_unusual: false,
            };
        }
        
        let total_bets = user_bets.len() as u64;
        let total_wagered: f64 = user_bets.iter().map(|b| b.amount).sum();
        let total_won: f64 = user_bets.iter().map(|b| b.win_amount).sum();
        let win_rate = if total_wagered > 0.0 {
            total_won / total_wagered
        } else {
            0.0
        };
        
        let avg_bet_size = total_wagered / total_bets as f64;
        let largest_bet = user_bets.iter().map(|b| b.amount).fold(0.0, f64::max);
        
        // Simple bet frequency calculation
        let timestamps: Vec<u64> = user_bets.iter().map(|b| b.timestamp).collect();
        let time_span = timestamps.iter().max().unwrap_or(&timestamp) - 
                        timestamps.iter().min().unwrap_or(&timestamp);
        let bet_frequency = if time_span > 0 {
            (total_bets as f64) / (time_span as f64 / 3600.0)
        } else {
            0.0
        };
        
        // Check if pattern is unusual
        let is_unusual = win_rate > self.config.suspicious_win_rate_threshold ||
                        largest_bet > avg_bet_size * 10.0;
        
        BettingPattern {
            total_bets,
            total_wagered,
            total_won,
            win_rate,
            avg_bet_size,
            largest_bet,
            bet_frequency_per_hour: bet_frequency,
            is_unusual,
        }
    }
    
    /// Record a suspicious activity
    fn record_activity(
        &mut self,
        activity_type: SuspiciousActivityType,
        user_id: String,
        severity: ActivitySeverity,
        description: String,
    ) {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        let activity = SuspiciousActivity {
            id: uuid::Uuid::new_v4().to_string(),
            activity_type,
            user_id,
            severity,
            description,
            evidence: HashMap::new(),
            timestamp,
            resolved: false,
            resolution_note: None,
        };
        
        self.suspicious_activities.push(activity);
    }
    
    /// Get all suspicious activities
    pub fn get_activities(&self, limit: usize) -> Vec<&SuspiciousActivity> {
        self.suspicious_activities
            .iter()
            .rev()
            .take(limit)
            .collect()
    }
    
    /// Calculate user risk score
    pub fn calculate_risk_score(&self, user_id: &str, pattern: &BettingPattern) -> u32 {
        let mut score = 0u32;
        
        // Win rate factor
        if pattern.win_rate > 0.8 {
            score += 40;
        } else if pattern.win_rate > 0.6 {
            score += 20;
        }
        
        // Bet size factor
        if pattern.largest_bet > pattern.avg_bet_size * 20.0 {
            score += 30;
        } else if pattern.largest_bet > pattern.avg_bet_size * 10.0 {
            score += 15;
        }
        
        // Frequency factor
        if pattern.bet_frequency_per_hour > 120.0 {
            score += 20;
        } else if pattern.bet_frequency_per_hour > 60.0 {
            score += 10;
        }
        
        score.min(100)
    }
}

/// Risk assessment result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RiskAssessment {
    pub user_id: String,
    pub risk_score: u32,
    pub risk_level: String,
    pub recommendations: Vec<String>,
    pub should_block: bool,
    pub should_review: bool,
}

impl FraudDetector {
    pub fn assess_risk(&self, user_id: &str) -> RiskAssessment {
        let pattern = self.analyze_betting_pattern(user_id);
        let risk_score = self.calculate_risk_score(user_id, &pattern);
        
        let (risk_level, should_block, should_review) = if risk_score >= 80 {
            ("Critical".to_string(), true, true)
        } else if risk_score >= 60 {
            ("High".to_string(), false, true)
        } else if risk_score >= 40 {
            ("Medium".to_string(), false, true)
        } else {
            ("Low".to_string(), false, false)
        };
        
        let mut recommendations = Vec::new();
        if risk_score >= 60 {
            recommendations.push("Enable enhanced monitoring".to_string());
        }
        if pattern.is_unusual {
            recommendations.push("Review betting patterns".to_string());
        }
        if pattern.win_rate > 0.7 {
            recommendations.push("Check for bonus abuse".to_string());
        }
        
        RiskAssessment {
            user_id: user_id.to_string(),
            risk_score,
            risk_level,
            recommendations,
            should_block,
            should_review,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_ip_tracking() {
        let mut detector = FraudDetector::new(FraudConfig::default());
        
        assert!(!detector.check_registration("192.168.1.1"));
        assert!(!detector.check_registration("192.168.1.1"));
        assert!(detector.check_registration("192.168.1.1")); // 4th registration
    }
    
    #[test]
    fn test_bet_check() {
        let mut detector = FraudDetector::new(FraudConfig::default());
        
        let result = detector.check_bet("user1", 100.0, "192.168.1.1");
        assert!(result.is_none());
    }
}
