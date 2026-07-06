'use client';

import React from 'react';
import Link from 'next/link';
import { Header, Footer } from '@/components/layout';
import { Button, Card } from '@/components/ui';
import ThemeToggle from '@/components/ui/ThemeToggle';
import styles from './page.module.css';

export default function HomePage() {
  const games = [
    { id: 'slots', name: 'Slot Machines', icon: '🎰', description: 'Spin to win with our collection of themed slot machines' },
    { id: 'dice', name: 'Dice', icon: '🎲', description: 'Roll the dice with customizable odds' },
    { id: 'roulette', name: 'Roulette', icon: '🎡', description: 'Classic European and American roulette' },
    { id: 'blackjack', name: 'Blackjack', icon: '🃏', description: 'Beat the dealer in this classic card game' },
  ];

  return (
    <>
      <Header />
      <ThemeToggle />
      <main className={styles.main}>
        {/* Hero Section */}
        <section className={styles.hero}>
          <div className={styles.heroBackground}>
            <div className={styles.heroGlow}></div>
            <div className={styles.heroPattern}></div>
          </div>
          <div className={styles.heroContent}>
            <div className={styles.heroBadge}>
              <span className={styles.badgeIcon}>⚡</span>
              <span>Instant Crypto Transactions</span>
            </div>
            <h1 className={styles.heroTitle}>
              Experience the <span className={styles.highlight}>Ultimate</span> Crypto Casino
            </h1>
            <p className={styles.heroSubtitle}>
              Play your favorite casino games with cryptocurrency. Lightning-fast payouts, provably fair gaming, and unmatched excitement await.
            </p>
            <div className={styles.heroActions}>
              <Link href="/auth/register">
                <Button variant="primary" size="lg">
                  Get Started
                </Button>
              </Link>
              <Link href="/games">
                <Button variant="outline" size="lg">
                  Browse Games
                </Button>
              </Link>
            </div>
            <div className={styles.heroStats}>
              <div className={styles.stat}>
                <span className={styles.statValue}>$10M+</span>
                <span className={styles.statLabel}>Paid Out</span>
              </div>
              <div className={styles.stat}>
                <span className={styles.statValue}>50K+</span>
                <span className={styles.statLabel}>Players</span>
              </div>
              <div className={styles.stat}>
                <span className={styles.statValue}>99.9%</span>
                <span className={styles.statLabel}>Uptime</span>
              </div>
            </div>
          </div>
          <div className={styles.heroVisual}>
            <div className={styles.tigerMascot}>
              <span className={styles.tigerEmoji}>🐯</span>
              <div className={styles.tigerGlow}></div>
            </div>
          </div>
        </section>

        {/* Games Section */}
        <section className={styles.games}>
          <div className={styles.container}>
            <div className={styles.sectionHeader}>
              <h2 className={styles.sectionTitle}>Featured Games</h2>
              <p className={styles.sectionSubtitle}>
                Explore our collection of provably fair casino games
              </p>
            </div>
            <div className={styles.gamesGrid}>
              {games.map((game, index) => (
                <Link href={`/games/${game.id}`} key={game.id}>
                  <Card variant="glow" padding="lg" className={`${styles.gameCard} stagger-${index + 1}`}>
                    <div className={styles.gameIcon}>{game.icon}</div>
                    <h3 className={styles.gameName}>{game.name}</h3>
                    <p className={styles.gameDescription}>{game.description}</p>
                  </Card>
                </Link>
              ))}
            </div>
            <div className={styles.viewAll}>
              <Link href="/games">
                <Button variant="secondary" size="lg">
                  View All Games
                </Button>
              </Link>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className={styles.features}>
          <div className={styles.container}>
            <div className={styles.sectionHeader}>
              <h2 className={styles.sectionTitle}>Why Choose TigerCasino?</h2>
              <p className={styles.sectionSubtitle}>
                Experience the future of online gambling
              </p>
            </div>
            <div className={styles.featuresGrid}>
              <div className={styles.feature}>
                <div className={styles.featureIcon}>🔒</div>
                <h3>Secure & Fair</h3>
                <p>Provably fair algorithms ensure every game outcome is random and verifiable</p>
              </div>
              <div className={styles.feature}>
                <div className={styles.featureIcon}>⚡</div>
                <h3>Instant Transactions</h3>
                <p>Deposit and withdraw with cryptocurrency in seconds, not days</p>
              </div>
              <div className={styles.feature}>
                <div className={styles.featureIcon}>🎯</div>
                <h3>Low House Edge</h3>
                <p>Enjoy some of the best odds in the industry with our low house edge games</p>
              </div>
              <div className={styles.feature}>
                <div className={styles.featureIcon}>🏆</div>
                <h3>VIP Program</h3>
                <p>Earn rewards and climb the VIP tiers for exclusive benefits</p>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className={styles.cta}>
          <div className={styles.container}>
            <div className={styles.ctaContent}>
              <h2>Ready to Start Winning?</h2>
              <p>Join thousands of players and experience the thrill of crypto gambling</p>
              <Link href="/auth/register">
                <Button variant="primary" size="lg">
                  Create Account
                </Button>
              </Link>
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </>
  );
}
