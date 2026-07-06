'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Badge } from '@/components/ui';
import styles from './games.module.css';

interface Game {
  id: string;
  name: string;
  type: string;
  icon: string;
  description: string;
  rtp: number;
  category: string;
  isHot?: boolean;
  isNew?: boolean;
  provider: string;
}

const games: Game[] = [
  // Slot Games
  { id: 'slots-tiger-king', name: 'Tiger King', type: 'slots', icon: '🐯', description: 'King of the jungle awaits', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', isHot: true },
  { id: 'slots-mega-moolah', name: 'Mega Moolah', type: 'slots', icon: '🦁', description: 'Win massive jackpots', rtp: 88.12, category: 'slots', provider: 'Microgaming', isNew: true },
  { id: 'slots-book-of-dead', name: 'Book of Dead', type: 'slots', icon: '📖', description: 'Ancient Egyptian adventure', rtp: 96.21, category: 'slots', provider: 'Play n GO' },
  { id: 'slots-starburst', name: 'Starburst', type: 'slots', icon: '💎', description: 'Cosmic gem explosion', rtp: 96.09, category: 'slots', provider: 'NetEnt', isHot: true },
  { id: 'slots-gonzo-quest', name: 'Gonzo\'s Quest', type: 'slots', icon: '🏔️', description: 'Quest for Eldorado', rtp: 95.97, category: 'slots', provider: 'NetEnt' },
  { id: 'slots-wolf-gold', name: 'Wolf Gold', type: 'slots', icon: '🐺', description: 'Pack rules the night', rtp: 96.01, category: 'slots', provider: 'Pragmatic Play' },
  { id: 'slots-sweet-bonanza', name: 'Sweet Bonanza', type: 'slots', icon: '🍬', description: 'Sweet candy rewards', rtp: 96.48, category: 'slots', provider: 'Pragmatic Play', isHot: true },
  { id: 'slots-gates-of-olympus', name: 'Gates of Olympus', type: 'slots', icon: '⚡', description: 'Zeus lightning wins', rtp: 95.51, category: 'slots', provider: 'Pragmatic Play', isNew: true },
  { id: 'slots-money-train', name: 'Money Train', type: 'slots', icon: '🚂', description: 'Wild west heist', rtp: 96.15, category: 'slots', provider: 'Relax Gaming' },
  { id: 'slots-big-bass-bonanza', name: 'Big Bass Bonanza', type: 'slots', icon: '🎣', description: 'Fishing for big wins', rtp: 96.71, category: 'slots', provider: 'Pragmatic Play' },
  
  // Table Games
  { id: 'dice-classic', name: 'Classic Dice', type: 'dice', icon: '🎲', description: 'Roll the dice, win big', rtp: 99, category: 'dice', provider: 'TigerCasino', isHot: true },
  { id: 'dice-plinko', name: 'Plinko', type: 'dice', icon: '🔻', description: 'Watch the ball drop', rtp: 98, category: 'dice', provider: 'TigerCasino', isNew: true },
  { id: 'dice-mines', name: 'Mines', type: 'dice', icon: '💣', description: 'Avoid the mines', rtp: 97, category: 'dice', provider: 'TigerCasino' },
  { id: 'dice-hilo', name: 'HiLo', type: 'dice', icon: '🔺', description: 'Predict the outcome', rtp: 96, category: 'dice', provider: 'TigerCasino' },
  { id: 'dice-keno', name: 'Keno', type: 'dice', icon: '🎱', description: 'Pick your numbers', rtp: 95, category: 'dice', provider: 'TigerCasino' },
  
  // Roulette
  { id: 'roulette-european', name: 'European Roulette', type: 'roulette', icon: '🎡', description: 'Classic single zero', rtp: 97.3, category: 'roulette', provider: 'Evolution', isHot: true },
  { id: 'roulette-american', name: 'American Roulette', type: 'roulette', icon: '🎰', description: 'Double zero action', rtp: 94.74, category: 'roulette', provider: 'Evolution' },
  { id: 'roulette-french', name: 'French Roulette', type: 'roulette', icon: '🇫🇷', description: 'La partage rules', rtp: 98.65, category: 'roulette', provider: 'Evolution' },
  { id: 'roulette-speed', name: 'Speed Roulette', type: 'roulette', icon: '⚡', description: 'Fast-paced action', rtp: 97.3, category: 'roulette', provider: 'Evolution', isNew: true },
  { id: 'roulette-immersive', name: 'Immersive Roulette', type: 'roulette', icon: '🎬', description: 'Cinematic experience', rtp: 97.3, category: 'roulette', provider: 'Evolution' },
  
  // Blackjack
  { id: 'blackjack-classic', name: 'Classic Blackjack', type: 'blackjack', icon: '🃏', description: 'Beat the dealer', rtp: 99.4, category: 'blackjack', provider: 'Evolution', isHot: true },
  { id: 'blackjack-vip', name: 'VIP Blackjack', type: 'blackjack', icon: '👑', description: 'High roller tables', rtp: 99.5, category: 'blackjack', provider: 'Evolution' },
  { id: 'blackjack-party', name: 'Party Blackjack', type: 'blackjack', icon: '🎉', description: 'Fun atmosphere', rtp: 99.4, category: 'blackjack', provider: 'Evolution', isNew: true },
  { id: 'blackjack-infinite', name: 'Infinite Blackjack', type: 'blackjack', icon: '♾️', description: 'Unlimited seats', rtp: 99.47, category: 'blackjack', provider: 'Evolution' },
  { id: 'blackjack-power', name: 'Power Blackjack', type: 'blackjack', icon: '⚡', description: 'Double down anytime', rtp: 99.5, category: 'blackjack', provider: 'Evolution' },
  
  // Baccarat
  { id: 'baccarat-classic', name: 'Classic Baccarat', type: 'baccarat', icon: '🪙', description: 'Punto Banco', rtp: 98.94, category: 'baccarat', provider: 'Evolution', isHot: true },
  { id: 'baccarat-speed', name: 'Speed Baccarat', type: 'baccarat', icon: '⚡', description: 'Lightning fast', rtp: 98.94, category: 'baccarat', provider: 'Evolution' },
  { id: 'baccarat-squeeze', name: 'Baccarat Squeeze', type: 'baccarat', icon: '🤏', description: 'Dramatic reveals', rtp: 98.94, category: 'baccarat', provider: 'Evolution' },
  { id: 'baccarat-no-comm', name: 'No Commission Baccarat', type: 'baccarat', icon: '💰', description: 'No banker fee', rtp: 98.76, category: 'baccarat', provider: 'Evolution', isNew: true },
  { id: 'baccarat-golden', name: 'Golden Wealth Baccarat', type: 'baccarat', icon: '✨', description: 'Golden multipliers', rtp: 98.94, category: 'baccarat', provider: 'Evolution' },
  
  // Poker
  { id: 'poker-texas-hold-em', name: 'Texas Hold\'em', type: 'poker', icon: '🃏', description: 'The classic game', rtp: 97.8, category: 'poker', provider: 'Evolution', isHot: true },
  { id: 'poker-caribbean', name: 'Caribbean Stud', type: 'poker', icon: '🏝️', description: 'Island style poker', rtp: 96.3, category: 'poker', provider: 'Evolution' },
  { id: 'poker-three-card', name: 'Three Card Poker', type: 'poker', icon: '3️⃣', description: 'Quick poker action', rtp: 96.63, category: 'poker', provider: 'Evolution' },
  { id: 'poker-ultimate', name: 'Ultimate Texas Hold\'em', type: 'poker', icon: '🏆', description: 'Ultimate poker', rtp: 97.8, category: 'poker', provider: 'Evolution' },
  { id: 'poker-side-bet', name: 'Casino Hold\'em', type: 'poker', icon: '🎯', description: 'Side bet action', rtp: 97.8, category: 'poker', provider: 'Evolution' },
  
  // Live Shows
  { id: 'show-dream-catcher', name: 'Dream Catcher', type: 'show', icon: '🎯', description: 'Spin the wheel', rtp: 96.58, category: 'live', provider: 'Evolution', isHot: true },
  { id: 'show-monopoly', name: 'Monopoly Live', type: 'show', icon: '🏠', description: 'Board game magic', rtp: 96.23, category: 'live', provider: 'Evolution', isNew: true },
  { id: 'show-crazy-time', name: 'Crazy Time', type: 'show', icon: '🎪', description: 'Insane multipliers', rtp: 95.5, category: 'live', provider: 'Evolution', isHot: true },
  { id: 'show-lightning-roulette', name: 'Lightning Roulette', type: 'show', icon: '⚡', description: 'Electrifying wins', rtp: 97.3, category: 'live', provider: 'Evolution' },
  { id: 'show-fan-tan', name: 'Fan Tan', type: 'show', icon: '🔴', description: 'Ancient Chinese', rtp: 97.5, category: 'live', provider: 'Evolution' },
];

const categories = [
  { id: 'all', name: 'All Games', icon: '🎮' },
  { id: 'slots', name: 'Slots', icon: '🎰' },
  { id: 'dice', name: 'Dice', icon: '🎲' },
  { id: 'roulette', name: 'Roulette', icon: '🎡' },
  { id: 'blackjack', name: 'Blackjack', icon: '🃏' },
  { id: 'baccarat', name: 'Baccarat', icon: '🪙' },
  { id: 'poker', name: 'Poker', icon: '♠️' },
  { id: 'live', name: 'Live Shows', icon: '🎬' },
];

export default function GamesPage() {
  const [activeCategory, setActiveCategory] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'rtp' | 'provider'>('name');

  const filteredGames = games
    .filter(game => 
      (activeCategory === 'all' || game.category === activeCategory) &&
      (game.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
       game.provider.toLowerCase().includes(searchQuery.toLowerCase()))
    )
    .sort((a, b) => {
      if (sortBy === 'name') return a.name.localeCompare(b.name);
      if (sortBy === 'rtp') return b.rtp - a.rtp;
      return a.provider.localeCompare(b.provider);
    });

  return (
    <>
      <Header />
      <main className={styles.main}>
        <section className={styles.hero}>
          <div className={styles.heroContent}>
            <h1 className={styles.title}>
              <span className={styles.highlight}>500+</span> Casino Games
            </h1>
            <p className={styles.subtitle}>
              Experience the thrill of real casino games with instant crypto payouts
            </p>
          </div>
        </section>

        <section className={styles.controls}>
          <div className={styles.searchBar}>
            <span className={styles.searchIcon}>🔍</span>
            <input
              type="text"
              placeholder="Search games or providers..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className={styles.searchInput}
            />
          </div>

          <div className={styles.sortWrapper}>
            <label>Sort by:</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as 'name' | 'rtp' | 'provider')}
              className={styles.sortSelect}
            >
              <option value="name">Name</option>
              <option value="rtp">RTP</option>
              <option value="provider">Provider</option>
            </select>
          </div>
        </section>

        <section className={styles.categories}>
          {categories.map(cat => (
            <button
              key={cat.id}
              className={`${styles.categoryBtn} ${activeCategory === cat.id ? styles.active : ''}`}
              onClick={() => setActiveCategory(cat.id)}
            >
              <span className={styles.categoryIcon}>{cat.icon}</span>
              <span className={styles.categoryName}>{cat.name}</span>
            </button>
          ))}
        </section>

        <section className={styles.gamesGrid}>
          {filteredGames.map((game, index) => (
            <Link href={`/games/${game.type}/${game.id}`} key={game.id}>
              <Card variant="glow" padding="none" className={`${styles.gameCard} stagger-${(index % 5) + 1}`}>
                <div className={styles.gameIcon}>
                  {game.icon}
                  {game.isHot && <span className={styles.hotBadge}>🔥 HOT</span>}
                  {game.isNew && <span className={styles.newBadge}>🆕 NEW</span>}
                </div>
                <div className={styles.gameInfo}>
                  <h3 className={styles.gameName}>{game.name}</h3>
                  <p className={styles.gameDescription}>{game.description}</p>
                  <div className={styles.gameMeta}>
                    <span className={styles.provider}>{game.provider}</span>
                    <span className={styles.rtp}>RTP: {game.rtp}%</span>
                  </div>
                </div>
                <div className={styles.playBtn}>
                  <Button variant="primary" size="sm">Play Now</Button>
                </div>
              </Card>
            </Link>
          ))}
        </section>

        {filteredGames.length === 0 && (
          <div className={styles.noResults}>
            <span className={styles.noResultsIcon}>🎮</span>
            <h3>No games found</h3>
            <p>Try adjusting your search or filters</p>
          </div>
        )}
      </main>
      <Footer />
    </>
  );
}
