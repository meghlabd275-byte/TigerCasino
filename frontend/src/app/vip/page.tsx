'use client';

import React, { useState, useEffect } from 'react';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Badge } from '@/components/ui';
import { api, VIPStatus, LeaderboardEntry, Promotion } from '@/lib/api';
import styles from './vip.module.css';

export default function VIPPage() {
  const [vipStatus, setVipStatus] = useState<VIPStatus | null>(null);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [promotions, setPromotions] = useState<Promotion[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [statusRes, leaderboardRes, promosRes] = await Promise.all([
        api.getVIPStatus(),
        api.getLeaderboard('weekly', 10),
        api.getPromotions()
      ]);

      if (statusRes.success && statusRes.data) {
        setVipStatus(statusRes.data);
      }
      if (leaderboardRes.success && leaderboardRes.data) {
        setLeaderboard(leaderboardRes.data);
      }
      if (promosRes.success && promosRes.data) {
        setPromotions(promosRes.data);
      }
    } catch (error) {
      console.error('Failed to load VIP data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleClaimRakeback = async () => {
    try {
      const result = await api.claimRakeback();
      if (result.success) {
        alert(`Successfully claimed $${result.data?.amount.toFixed(2)} rakeback!`);
        loadData();
      }
    } catch (error) {
      console.error('Failed to claim rakeback:', error);
    }
  };

  const handleClaimWelcomeBonus = async () => {
    try {
      const result = await api.claimWelcomeBonus();
      if (result.success) {
        alert(`Welcome bonus claimed! $${result.data?.amount} with ${result.data?.wagerReq}x wager requirement`);
        loadData();
      }
    } catch (error: any) {
      alert(error?.response?.data?.error || 'Failed to claim bonus');
    }
  };

  if (loading) {
    return (
      <div className={styles.loading}>
        <div className={styles.spinner}></div>
        <p>Loading VIP Status...</p>
      </div>
    );
  }

  const levelColors = ['#CD7F32', '#C0C0C0', '#FFD700', '#E5E4E2', '#B9F2FF', '#FF0000'];
  const levelColor = vipStatus ? levelColors[vipStatus.level] : '#CD7F32';

  return (
    <>
      <Header />
      <main className={styles.main}>
        <div className={styles.container}>
          <h1 className={styles.title}>VIP Club</h1>

          {/* VIP Status Card */}
          {vipStatus && (
            <div className={styles.statusCard} style={{ borderColor: levelColor }}>
              <div className={styles.levelBadge} style={{ backgroundColor: levelColor }}>
                {vipStatus.levelName}
              </div>
              
              <div className={styles.statsGrid}>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Total Wagered</span>
                  <span className={styles.statValue}>${vipStatus.totalWagered.toLocaleString()}</span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Points</span>
                  <span className={styles.statValue}>{vipStatus.points.toLocaleString()}</span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Rakeback</span>
                  <span className={styles.statValue}>{vipStatus.rakebackPercent}%</span>
                </div>
                <div className={styles.stat}>
                  <span className={styles.statLabel}>Available Rakeback</span>
                  <span className={styles.statValue}>${vipStatus.rakebackBalance.toFixed(2)}</span>
                </div>
              </div>

              {vipStatus.nextLevel && (
                <div className={styles.progressSection}>
                  <div className={styles.progressLabel}>
                    <span>Progress to {vipStatus.nextLevel.name}</span>
                    <span>{vipStatus.progressToNext.toFixed(1)}%</span>
                  </div>
                  <div className={styles.progressBar}>
                    <div 
                      className={styles.progressFill} 
                      style={{ width: `${vipStatus.progressToNext}%`, backgroundColor: levelColor }}
                    ></div>
                  </div>
                </div>
              )}

              <div className={styles.actions}>
                <Button 
                  onClick={handleClaimRakeback}
                  disabled={vipStatus.rakebackBalance <= 0}
                >
                  Claim Rakeback (${vipStatus.rakebackBalance.toFixed(2)})
                </Button>
              </div>
            </div>
          )}

          {/* Benefits Section */}
          {vipStatus?.benefits && (
            <Card className={styles.section}>
              <h2 className={styles.sectionTitle}>Your Benefits</h2>
              <div className={styles.benefitsGrid}>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Max Bet</span>
                  <span className={styles.benefitValue}>${vipStatus.benefits.maxBet.toLocaleString()}</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Max Win</span>
                  <span className={styles.benefitValue}>${vipStatus.benefits.maxWin.toLocaleString()}</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Withdrawal Limit</span>
                  <span className={styles.benefitValue}>${vipStatus.benefits.withdrawalLimit.toLocaleString()}</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Withdrawal Fee</span>
                  <span className={styles.benefitValue}>${vipStatus.benefits.withdrawalFee}</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Cashback</span>
                  <span className={styles.benefitValue}>{vipStatus.benefits.cashbackPercent}%</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Points Multiplier</span>
                  <span className={styles.benefitValue}>{vipStatus.benefits.pointsMultiplier}x</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Priority Support</span>
                  <span className={styles.benefitValue}>{vipStatus.benefits.prioritySupport ? '✓' : '✗'}</span>
                </div>
                <div className={styles.benefit}>
                  <span className={styles.benefitLabel}>Personal Host</span>
                  <span className={styles.benefitValue}>{vipStatus.benefits.personalHost ? '✓' : '✗'}</span>
                </div>
              </div>
            </Card>
          )}

          {/* Promotions */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>Active Promotions</h2>
            <div className={styles.promotions}>
              {promotions.length === 0 ? (
                <p className={styles.noData}>No active promotions</p>
              ) : (
                promotions.map((promo) => (
                  <div key={promo.id} className={styles.promotion}>
                    <h3>{promo.name}</h3>
                    <p>{promo.description}</p>
                    <div className={styles.promoDetails}>
                      <Badge variant="success">Bonus: ${promo.bonusAmount}</Badge>
                      <Badge variant="info">Wager: {promo.wagerReq}x</Badge>
                    </div>
                  </div>
                ))
              )}
              
              {/* Welcome Bonus */}
              <div className={styles.promotion}>
                <h3>🎁 Welcome Bonus</h3>
                <p>Claim your $100 welcome bonus!</p>
                <div className={styles.promoDetails}>
                  <Badge variant="success">$100 Bonus</Badge>
                  <Badge variant="info">10x Wager</Badge>
                </div>
                <Button onClick={handleClaimWelcomeBonus}>Claim Now</Button>
              </div>
            </div>
          </Card>

          {/* Leaderboard */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>🏆 Weekly Leaderboard</h2>
            <div className={styles.leaderboard}>
              {leaderboard.map((entry, index) => (
                <div key={entry.userId} className={styles.leaderboardEntry}>
                  <span className={styles.rank}>
                    {index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : `#${entry.rank}`}
                  </span>
                  <span className={styles.username}>{entry.username}</span>
                  <span className={styles.score}>{entry.score.toLocaleString()} pts</span>
                </div>
              ))}
            </div>
          </Card>

          {/* VIP Levels */}
          <Card className={styles.section}>
            <h2 className={styles.sectionTitle}>VIP Levels</h2>
            <div className={styles.levelsGrid}>
              {['Bronze', 'Silver', 'Gold', 'Platinum', 'Diamond', 'VIP'].map((level, index) => (
                <div key={level} className={styles.levelCard} style={{ borderColor: levelColors[index] }}>
                  <div className={styles.levelName} style={{ color: levelColors[index] }}>{level}</div>
                  <div className={styles.levelRakeback}>{3.5 + index * 2.5}% Rakeback</div>
                  <div className={styles.levelCashback}>{index}% Cashback</div>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </main>
      <Footer />
    </>
  );
}
