'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Badge } from '@/components/ui';
import { useAuth } from '@/contexts/AuthContext';
import styles from './dashboard.module.css';

interface RecentTransaction {
  id: string;
  type: 'deposit' | 'withdrawal' | 'bet' | 'win';
  amount: number;
  status: 'completed' | 'pending' | 'rejected';
  date: string;
  currency: string;
}

interface RecentGame {
  id: string;
  name: string;
  result: 'won' | 'lost';
  amount: number;
  multiplier: number;
  date: string;
}

const recentTransactions: RecentTransaction[] = [
  { id: '1', type: 'deposit', amount: 0.5, status: 'completed', date: '2024-01-15 14:30', currency: 'BTC' },
  { id: '2', type: 'win', amount: 0.125, status: 'completed', date: '2024-01-15 13:45', currency: 'BTC' },
  { id: '3', type: 'bet', amount: 0.01, status: 'completed', date: '2024-01-15 12:20', currency: 'BTC' },
  { id: '4', type: 'withdrawal', amount: 0.25, status: 'pending', date: '2024-01-15 11:00', currency: 'BTC' },
  { id: '5', type: 'win', amount: 0.05, status: 'completed', date: '2024-01-14 22:30', currency: 'BTC' },
];

const recentGames: RecentGame[] = [
  { id: '1', name: 'Tiger King Slots', result: 'won', amount: 0.05, multiplier: 2.5, date: '2024-01-15 14:30' },
  { id: '2', name: 'Classic Dice', result: 'lost', amount: 0.01, multiplier: 0, date: '2024-01-15 13:45' },
  { id: '3', name: 'European Roulette', result: 'won', amount: 0.1, multiplier: 2, date: '2024-01-15 12:20' },
  { id: '4', name: 'Blackjack VIP', result: 'won', amount: 0.2, multiplier: 1.5, date: '2024-01-14 22:30' },
];

const quickPlayGames = [
  { id: 'slots-tiger-king', name: 'Tiger King', icon: '🐯', type: 'slots' },
  { id: 'dice-classic', name: 'Classic Dice', icon: '🎲', type: 'dice' },
  { id: 'roulette-european', name: 'European Roulette', icon: '🎡', type: 'roulette' },
  { id: 'blackjack-classic', name: 'Blackjack', icon: '🃏', type: 'blackjack' },
  { id: 'baccarat-classic', name: 'Baccarat', icon: '🪙', type: 'baccarat' },
  { id: 'show-crazy-time', name: 'Crazy Time', icon: '🎪', type: 'show' },
];

export default function DashboardPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted || isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className={styles.container}>
        <Header />
        <main className={styles.main}>
          <div className={styles.notAuthenticated}>
            <h2>Please Login to Access Dashboard</h2>
            <Link href="/auth/login">
              <Button variant="primary" size="lg">Login</Button>
            </Link>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <>
      <Header />
      <main className={styles.main}>
        <section className={styles.welcomeSection}>
          <div className={styles.welcomeContent}>
            <h1>Welcome back, <span className={styles.username}>{user?.username || 'Player'}!</span></h1>
            <p>Ready to win big today?</p>
          </div>
          <div className={styles.balanceCard}>
            <div className={styles.balanceItem}>
              <span className={styles.balanceLabel}>Main Balance</span>
              <span className={styles.balanceValue}>${(user?.balance || 0).toFixed(2)}</span>
            </div>
            <div className={styles.balanceItem}>
              <span className={styles.balanceLabel}>Bonus Balance</span>
              <span className={styles.bonusValue}>${(user?.bonusBalance || 0).toFixed(2)}</span>
            </div>
            <div className={styles.vipBadge}>
              <span className={styles.vipIcon}>👑</span>
              <span>VIP Level {user?.vipLevel || 0}</span>
            </div>
          </div>
        </section>

        <section className={styles.quickActions}>
          <Link href="/wallet/deposit">
            <Button variant="primary" size="lg">💰 Deposit</Button>
          </Link>
          <Link href="/wallet/withdraw">
            <Button variant="secondary" size="lg">💸 Withdraw</Button>
          </Link>
          <Link href="/games">
            <Button variant="outline" size="lg">🎮 Play Now</Button>
          </Link>
        </section>

        <section className={styles.quickPlay}>
          <h2 className={styles.sectionTitle}>Quick Play</h2>
          <div className={styles.quickPlayGrid}>
            {quickPlayGames.map(game => (
              <Link href={`/games/${game.type}/${game.id}`} key={game.id}>
                <Card variant="glow" padding="md" className={styles.quickPlayCard}>
                  <span className={styles.gameIcon}>{game.icon}</span>
                  <span className={styles.gameName}>{game.name}</span>
                </Card>
              </Link>
            ))}
          </div>
        </section>

        <section className={styles.stats}>
          <div className={styles.statsGrid}>
            <Card variant="bordered" padding="lg">
              <div className={styles.statItem}>
                <span className={styles.statIcon}>🎯</span>
                <span className={styles.statValue}>1,234</span>
                <span className={styles.statLabel}>Total Bets</span>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statItem}>
                <span className={styles.statIcon}>🏆</span>
                <span className={styles.statValue}>$5,678</span>
                <span className={styles.statLabel}>Total Wins</span>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statItem}>
                <span className={styles.statIcon}>🔥</span>
                <span className={styles.statValue}>$123.45</span>
                <span className={styles.statLabel}>Biggest Win</span>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statItem}>
                <span className={styles.statIcon}>📈</span>
                <span className={styles.statValue}>98.5%</span>
                <span className={styles.statLabel}>Win Rate</span>
              </div>
            </Card>
          </div>
        </section>

        <section className={styles.activity}>
          <div className={styles.activityGrid}>
            <Card variant="bordered" padding="lg">
              <h3 className={styles.cardTitle}>Recent Transactions</h3>
              <div className={styles.transactionList}>
                {recentTransactions.map(tx => (
                  <div key={tx.id} className={styles.transactionItem}>
                    <div className={styles.txInfo}>
                      <span className={`${styles.txType} ${styles[tx.type]}`}>
                        {tx.type === 'deposit' && '💰'}
                        {tx.type === 'withdrawal' && '💸'}
                        {tx.type === 'bet' && '🎯'}
                        {tx.type === 'win' && '🏆'}
                      </span>
                      <span className={styles.txAmount}>
                        {tx.type === 'bet' || tx.type === 'withdrawal' ? '-' : '+'}{tx.amount} {tx.currency}
                      </span>
                    </div>
                    <div className={styles.txStatus}>
                      <Badge variant={tx.status === 'completed' ? 'success' : tx.status === 'pending' ? 'warning' : 'error'} size="sm">
                        {tx.status}
                      </Badge>
                      <span className={styles.txDate}>{tx.date}</span>
                    </div>
                  </div>
                ))}
              </div>
              <Link href="/wallet/transactions" className={styles.viewAllLink}>
                View All Transactions →
              </Link>
            </Card>

            <Card variant="bordered" padding="lg">
              <h3 className={styles.cardTitle}>Recent Games</h3>
              <div className={styles.gameList}>
                {recentGames.map(game => (
                  <div key={game.id} className={styles.gameItem}>
                    <div className={styles.gameInfo}>
                      <span className={styles.gameName}>{game.name}</span>
                      <span className={styles.gameDate}>{game.date}</span>
                    </div>
                    <div className={styles.gameResult}>
                      <Badge variant={game.result === 'won' ? 'success' : 'error'} size="sm">
                        {game.result === 'won' ? `+${game.multiplier}x` : 'Lost'}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
              <Link href="/games/history" className={styles.viewAllLink}>
                View Game History →
              </Link>
            </Card>
          </div>
        </section>
      </main>
      <Footer />
    </>
  );
}
