'use client';

import React, { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import toast from 'react-hot-toast';

// Live Dealer Tables Data
const LIVE_TABLES = [
  // Blackjack
  { id: 'bj_vip', game: 'Blackjack', table: 'VIP Blackjack', provider: 'Evolution', dealer: 'Sarah', minBet: 100, maxBet: 50000, players: 3, status: 'live', image: '🃏' },
  { id: 'bj_speed', game: 'Blackjack', table: 'Speed Blackjack', provider: 'Evolution', dealer: 'Mike', minBet: 5, maxBet: 1000, players: 5, status: 'live', image: '🃏' },
  { id: 'bj_infinite', game: 'Blackjack', table: 'Infinite Blackjack', provider: 'Evolution', dealer: 'Emma', minBet: 1, maxBet: 5000, players: 12, status: 'live', image: '🃏' },
  { id: 'bj_free', game: 'Blackjack', table: 'Free Bet Blackjack', provider: 'Evolution', dealer: 'David', minBet: 10, maxBet: 2500, players: 4, status: 'live', image: '🃏' },
  { id: 'bj_power', game: 'Blackjack', table: 'Power Blackjack', provider: 'Evolution', dealer: 'Lisa', minBet: 25, maxBet: 10000, players: 6, status: 'live', image: '🃏' },
  { id: 'bj_quantum', game: 'Blackjack', table: 'Quantum Blackjack', provider: 'Evolution', dealer: 'James', minBet: 50, maxBet: 25000, players: 3, status: 'live', image: '🃏' },
  
  // Roulette
  { id: 'rl_eu', game: 'Roulette', table: 'European Roulette', provider: 'Evolution', dealer: 'Emma', minBet: 1, maxBet: 5000, players: 8, status: 'live', image: '🎡' },
  { id: 'rl_lightning', game: 'Roulette', table: 'Lightning Roulette', provider: 'Evolution', dealer: 'Sophie', minBet: 0.5, maxBet: 5000, players: 15, status: 'live', image: '⚡' },
  { id: 'rl_immersive', game: 'Roulette', table: 'Immersive Roulette', provider: 'Evolution', dealer: 'Marco', minBet: 1, maxBet: 10000, players: 10, status: 'live', image: '🎬' },
  { id: 'rl_auto', game: 'Roulette', table: 'Auto Roulette', provider: 'Evolution', dealer: 'Auto', minBet: 0.1, maxBet: 1000, players: 20, status: 'live', image: '🤖' },
  { id: 'rl_french', game: 'Roulette', table: 'French Roulette', provider: 'Evolution', dealer: 'Pierre', minBet: 10, maxBet: 5000, players: 6, status: 'live', image: '🇫🇷' },
  { id: 'rl_speed', game: 'Roulette', table: 'Speed Auto Roulette', provider: 'Evolution', dealer: 'Auto', minBet: 0.1, maxBet: 500, players: 25, status: 'live', image: '💨' },
  
  // Baccarat
  { id: 'bc_classic', game: 'Baccarat', table: 'Classic Baccarat', provider: 'Evolution', dealer: 'David', minBet: 20, maxBet: 10000, players: 6, status: 'live', image: '🎴' },
  { id: 'bc_speed', game: 'Baccarat', table: 'Speed Baccarat', provider: 'Evolution', dealer: 'Jennifer', minBet: 10, maxBet: 5000, players: 8, status: 'live', image: '💨' },
  { id: 'bc_squeeze', game: 'Baccarat', table: 'Baccarat Squeeze', provider: 'Evolution', dealer: 'Chen', minBet: 50, maxBet: 25000, players: 4, status: 'live', image: '🖐️' },
  { id: 'bc_golden', game: 'Baccarat', table: 'Golden Wealth Baccarat', provider: 'Evolution', dealer: 'Linda', minBet: 1, maxBet: 2000, players: 10, status: 'live', image: '✨' },
  { id: 'bc_no_comm', game: 'Baccarat', table: 'No Commission Baccarat', provider: 'Evolution', dealer: 'Alex', minBet: 5, maxBet: 5000, players: 7, status: 'live', image: '💰' },
  { id: 'bc_control', game: 'Baccarat', table: 'Baccarat Control Squeeze', provider: 'Evolution', dealer: 'Kevin', minBet: 100, maxBet: 50000, players: 3, status: 'live', image: '🎯' },
  
  // Poker
  { id: 'pk_casino', game: 'Poker', table: 'Casino Hold\'em', provider: 'Evolution', dealer: 'Lisa', minBet: 5, maxBet: 2000, players: 4, status: 'live', image: '🃏' },
  { id: 'pk_texas', game: 'Poker', table: 'Ultimate Texas Hold\'em', provider: 'Evolution', dealer: 'Mark', minBet: 10, maxBet: 5000, players: 3, status: 'live', image: '🤠' },
  { id: 'pk_caribbean', game: 'Poker', table: 'Caribbean Stud Poker', provider: 'Evolution', dealer: 'Anna', minBet: 10, maxBet: 5000, players: 5, status: 'live', image: '🏝️' },
  { id: 'pk_three', game: 'Poker', table: 'Three Card Poker', provider: 'Evolution', dealer: 'Tom', minBet: 5, maxBet: 2500, players: 6, status: 'live', image: '3️⃣' },
  { id: 'pk_side', game: 'Poker', table: 'Side Bet City', provider: 'Evolution', dealer: 'Rachel', minBet: 1, maxBet: 1000, players: 8, status: 'live', image: '🎰' },
];

// Game Shows
const GAME_SHOWS = [
  { id: 'crazy_time', game: 'Crazy Time', provider: 'Evolution', minBet: 0.1, maxBet: 5000, players: 20, status: 'live', image: '🎡', icon: '🎡' },
  { id: 'monopoly', game: 'Monopoly Live', provider: 'Evolution', minBet: 0.1, maxBet: 5000, players: 18, status: 'live', image: '🎲', icon: '🎲' },
  { id: 'lightning_dice', game: 'Lightning Dice', provider: 'Evolution', minBet: 0.2, maxBet: 2500, players: 15, status: 'live', image: '🎲', icon: '⚡' },
  { id: 'dream_catcher', game: 'Dream Catcher', provider: 'Evolution', minBet: 0.1, maxBet: 5000, players: 12, status: 'live', image: '🎯', icon: '🎯' },
  { id: 'megawheel', game: 'Mega Wheel', provider: 'Pragmatic Play', minBet: 0.1, maxBet: 5000, players: 16, status: 'live', image: '🎰', icon: '🎰' },
  { id: 'sweet_bonanza', game: 'Sweet Bonanza CandyLand', provider: 'Pragmatic Play', minBet: 0.1, maxBet: 2500, players: 14, status: 'live', image: '🍬', icon: '🍬' },
  { id: 'gonzo_treasure', game: 'Gonzo\'s Treasure Hunt', provider: 'Evolution', minBet: 0.1, maxBet: 1250, players: 10, status: 'live', image: '🌋', icon: '🌋' },
  { id: 'deal_no_deal', game: 'Deal or No Deal Live', provider: 'Evolution', minBet: 0.1, maxBet: 5000, players: 8, status: 'live', image: '💼', icon: '💼' },
  { id: 'football_studio', game: 'Football Studio', provider: 'Evolution', minBet: 1, maxBet: 5000, players: 12, status: 'live', image: '⚽', icon: '⚽' },
  { id: 'moneywheel', game: 'The Money Wheel', provider: 'Evolution', minBet: 0.1, maxBet: 2500, players: 15, status: 'live', image: '💵', icon: '💵' },
];

// Providers
const PROVIDERS = [
  { id: 'all', name: 'All Providers', count: LIVE_TABLES.length + GAME_SHOWS.length },
  { id: 'Evolution', name: 'Evolution Gaming', count: LIVE_TABLES.filter(t => t.provider === 'Evolution').length },
  { id: 'Pragmatic Play', name: 'Pragmatic Play', count: GAME_SHOWS.filter(s => s.provider === 'Pragmatic Play').length },
];

// Game Categories
const CATEGORIES = [
  { id: 'all', name: 'All Games', icon: '🎰' },
  { id: 'Blackjack', name: 'Blackjack', icon: '🃏' },
  { id: 'Roulette', name: 'Roulette', icon: '🎡' },
  { id: 'Baccarat', name: 'Baccarat', icon: '🎴' },
  { id: 'Poker', name: 'Poker', icon: '🃏' },
  { id: 'Game Shows', name: 'Game Shows', icon: '🎬' },
];

export default function LiveCasinoPage() {
  const [activeCategory, setActiveCategory] = useState('all');
  const [activeProvider, setActiveProvider] = useState('all');
  const [showLiveOnly, setShowLiveOnly] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  const filteredTables = useMemo(() => {
    let tables = activeCategory === 'all' || activeCategory === 'Game Shows' ? [] : LIVE_TABLES.filter(t => t.game === activeCategory);
    
    if (activeCategory === 'Game Shows') {
      tables = [];
    }
    
    if (activeProvider !== 'all') {
      tables = tables.filter(t => t.provider === activeProvider);
    }
    
    if (searchQuery) {
      tables = tables.filter(t => 
        t.table.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.dealer.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    return tables;
  }, [activeCategory, activeProvider, searchQuery]);

  const filteredShows = useMemo(() => {
    let shows = activeCategory === 'all' || activeCategory !== 'Game Shows' ? GAME_SHOWS : [];
    
    if (activeCategory === 'Game Shows') {
      shows = GAME_SHOWS;
    }
    
    if (activeProvider !== 'all') {
      shows = shows.filter(s => s.provider === activeProvider);
    }
    
    if (searchQuery) {
      shows = shows.filter(s => 
        s.game.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    return shows;
  }, [activeCategory, activeProvider, searchQuery]);

  const joinTable = (tableId: string, tableName: string) => {
    toast.success(`Joining ${tableName}...`);
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 lg:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl lg:text-4xl font-heading text-gradient mb-1">🎰 Live Casino</h1>
            <p className="text-gray-400">Real dealers, real-time action • 500+ tables</p>
          </div>
          
          {/* Search */}
          <div className="relative">
            <input
              type="text"
              placeholder="Search tables..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="bg-tiger-surface border border-tiger-border rounded-lg px-4 py-2 w-64 text-white"
            />
          </div>
        </div>

        {/* Live Banner */}
        <div className="bg-gradient-to-r from-red-900/50 to-primary-900/50 rounded-xl p-4 mb-6 border border-red-500/30">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <span className="w-3 h-3 bg-red-500 rounded-full animate-pulse"></span>
              <span className="text-lg font-semibold">🔴 {LIVE_TABLES.filter(t => t.status === 'live').length + GAME_SHOWS.filter(s => s.status === 'live').length} Live Tables Now</span>
            </div>
            <div className="text-gray-400 text-sm">
              Average wait time: <span className="text-green-400">Under 30 seconds</span>
            </div>
          </div>
        </div>

        {/* Category Filter */}
        <div className="flex gap-2 mb-4 overflow-x-auto pb-2">
          {CATEGORIES.map(cat => (
            <button
              key={cat.id}
              onClick={() => setActiveCategory(cat.id)}
              className={`px-4 py-2 rounded-lg whitespace-nowrap flex items-center gap-2 ${
                activeCategory === cat.id 
                  ? 'bg-primary-500 text-white' 
                  : 'bg-tiger-surface text-gray-400 hover:bg-tiger-border'
              }`}
            >
              {cat.icon} {cat.name}
            </button>
          ))}
        </div>

        {/* Provider Filter */}
        <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
          {PROVIDERS.map(provider => (
            <button
              key={provider.id}
              onClick={() => setActiveProvider(provider.id)}
              className={`px-3 py-1.5 rounded-full text-sm whitespace-nowrap ${
                activeProvider === provider.id 
                  ? 'bg-primary-500 text-white' 
                  : 'bg-tiger-surface text-gray-400 hover:bg-tiger-border'
              }`}
            >
              {provider.name} ({provider.count})
            </button>
          ))}
        </div>

        {/* Live Tables */}
        {(activeCategory === 'all' || activeCategory !== 'Game Shows') && filteredTables.length > 0 && (
          <div className="mb-8">
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-2">
              🃏 Live Tables
              <span className="text-sm font-normal text-gray-400">({filteredTables.length})</span>
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {filteredTables.map(table => (
                <motion.div
                  key={table.id}
                  whileHover={{ scale: 1.02 }}
                  className="glass rounded-xl overflow-hidden cursor-pointer"
                >
                  <div className="relative">
                    <div className="aspect-video bg-gradient-to-br from-tiger-surface to-tiger-dark flex items-center justify-center text-6xl">
                      {table.image}
                    </div>
                    <div className="absolute top-2 right-2 flex items-center gap-1 bg-black/60 px-2 py-1 rounded-full">
                      <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                      <span className="text-xs text-white">{table.players} playing</span>
                    </div>
                    {table.status === 'live' && (
                      <div className="absolute top-2 left-2 bg-red-600 text-white text-xs px-2 py-1 rounded">
                        LIVE
                      </div>
                    )}
                  </div>
                  
                  <div className="p-3">
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <h3 className="font-bold text-white">{table.table}</h3>
                        <p className="text-gray-400 text-xs">{table.game} • {table.provider}</p>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2 mb-2 text-sm text-gray-400">
                      <span>👤 {table.dealer}</span>
                    </div>
                    
                    <div className="flex justify-between items-center">
                      <span className="text-green-400 text-sm font-medium">
                        ${table.minBet} - ${table.maxBet.toLocaleString()}
                      </span>
                      <button
                        onClick={() => joinTable(table.id, table.table)}
                        className="bg-primary-500 hover:bg-primary-600 text-white px-4 py-1.5 rounded-lg text-sm font-medium"
                      >
                        Join
                      </button>
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        )}

        {/* Game Shows */}
        {(activeCategory === 'all' || activeCategory === 'Game Shows') && filteredShows.length > 0 && (
          <div className="mb-8">
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-2">
              🎬 Game Shows
              <span className="text-sm font-normal text-gray-400">({filteredShows.length})</span>
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {filteredShows.map(show => (
                <motion.div
                  key={show.id}
                  whileHover={{ scale: 1.02 }}
                  className="glass rounded-xl overflow-hidden cursor-pointer"
                >
                  <div className="relative">
                    <div className="aspect-video bg-gradient-to-br from-purple-900/50 to-tiger-dark flex items-center justify-center text-6xl">
                      {show.icon}
                    </div>
                    <div className="absolute top-2 right-2 flex items-center gap-1 bg-black/60 px-2 py-1 rounded-full">
                      <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                      <span className="text-xs text-white">{show.players} watching</span>
                    </div>
                    {show.status === 'live' && (
                      <div className="absolute top-2 left-2 bg-red-600 text-white text-xs px-2 py-1 rounded">
                        LIVE
                      </div>
                    )}
                  </div>
                  
                  <div className="p-3">
                    <h3 className="font-bold text-white">{show.game}</h3>
                    <p className="text-gray-400 text-xs mb-2">{show.provider}</p>
                    
                    <div className="flex justify-between items-center">
                      <span className="text-green-400 text-sm font-medium">
                        ${show.minBet} - ${show.maxBet.toLocaleString()}
                      </span>
                      <button
                        onClick={() => joinTable(show.id, show.game)}
                        className="bg-purple-500 hover:bg-purple-600 text-white px-4 py-1.5 rounded-lg text-sm font-medium"
                      >
                        Play
                      </button>
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        )}

        {/* Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8">
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-primary-400">500+</div>
            <div className="text-gray-400 text-sm">Live Tables</div>
          </div>
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-green-400">24/7</div>
            <div className="text-gray-400 text-sm">Live Gaming</div>
          </div>
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-yellow-400">2s</div>
            <div className="text-gray-400 text-sm">Avg Response</div>
          </div>
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-purple-400">HD</div>
            <div className="text-gray-400 text-sm">Stream Quality</div>
          </div>
        </div>
      </div>
    </div>
  );
}
