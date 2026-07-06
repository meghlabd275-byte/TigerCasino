'use client';

import React, { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

// Game show data
const gameShows = [
  { id: 'crazy-time', name: 'Crazy Time', provider: 'Evolution Gaming', rtp: '96.08%', minBet: 0.10, maxBet: 5000, image: '/games/crazy-time.jpg', featured: true },
  { id: 'monopoly-live', name: 'Monopoly Live', provider: 'Evolution Gaming', rtp: '96.23%', minBet: 0.10, maxBet: 5000, image: '/games/monopoly-live.jpg', featured: true },
  { id: 'lightning-roulette', name: 'Lightning Roulette', provider: 'Evolution Gaming', rtp: '97.30%', minBet: 0.10, maxBet: 5000, image: '/games/lightning-roulette.jpg', featured: true },
  { id: 'dream-catcher', name: 'Dream Catcher', provider: 'Evolution Gaming', rtp: '96.58%', minBet: 0.10, maxBet: 5000, image: '/games/dream-catcher.jpg', featured: false },
  { id: 'deal-or-no-deal', name: 'Deal or No Deal', provider: 'Evolution Gaming', rtp: '95.42%', minBet: 0.10, maxBet: 5000, image: '/games/deal-or-no-deal.jpg', featured: false },
  { id: 'mega-ball', name: 'Mega Ball', provider: 'Evolution Gaming', rtp: '95.50%', minBet: 0.10, maxBet: 1000, image: '/games/mega-ball.jpg', featured: false },
  { id: 'funky-time', name: 'Funky Time', provider: 'Evolution Gaming', rtp: '95.99%', minBet: 0.10, maxBet: 5000, image: '/games/funky-time.jpg', featured: false },
  { id: 'cash-or-crash', name: 'Cash or Crash', provider: 'Evolution Gaming', rtp: '96.10%', minBet: 0.10, maxBet: 5000, image: '/games/cash-or-crash.jpg', featured: false },
  { id: 'mega-wheel', name: 'Mega Wheel', provider: 'Pragmatic Play', rtp: '96.00%', minBet: 0.10, maxBet: 5000, image: '/games/mega-wheel.jpg', featured: false },
  { id: 'gold-vault', name: 'Gold Vault Roulette', provider: 'Evolution Gaming', rtp: '97.30%', minBet: 1.0, maxBet: 5000, image: '/games/gold-vault.jpg', featured: false },
];

export default function GameShowsPage() {
  const { user, isAuthenticated } = useAuth();
  const router = useRouter();
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredGames, setFilteredGames] = useState(gameShows);

  useEffect(() => {
    if (searchTerm) {
      setFilteredGames(
        gameShows.filter(game =>
          game.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          game.provider.toLowerCase().includes(searchTerm.toLowerCase())
        )
      );
    } else {
      setFilteredGames(gameShows);
    }
  }, [searchTerm]);

  const handlePlayGame = (gameId: string) => {
    if (!isAuthenticated) {
      router.push('/auth/login');
      return;
    }
    router.push(`/games/game-shows/${gameId}`);
  };

  const featuredGames = filteredGames.filter(g => g.featured);
  const regularGames = filteredGames.filter(g => !g.featured);

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800">
      <div className="relative py-16 px-4 bg-gradient-to-r from-purple-900 to-indigo-900">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">Live Game Shows</h1>
          <p className="text-xl text-purple-200 max-w-2xl">
            Experience the excitement of live game shows with real dealers and interactive bonus rounds
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        <input
          type="text"
          placeholder="Search game shows..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full px-6 py-4 bg-gray-800 border border-gray-700 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500"
        />
      </div>

      {featuredGames.length > 0 && (
        <div className="max-w-7xl mx-auto px-4 mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Featured Game Shows</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {featuredGames.map((game) => (
              <div
                key={game.id}
                className="group relative bg-gray-800 rounded-2xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-purple-500/30 hover:border-purple-500"
                onClick={() => handlePlayGame(game.id)}
              >
                <div className="aspect-video bg-gradient-to-br from-purple-600 to-indigo-600 flex items-center justify-center">
                  <span className="text-6xl">🎡</span>
                </div>
                <div className="p-6">
                  <h3 className="text-xl font-bold text-white mb-2">{game.name}</h3>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-purple-400">{game.provider}</span>
                    <span className="text-green-400">RTP: {game.rtp}</span>
                  </div>
                  <div className="mt-4 flex items-center justify-between text-sm text-gray-400">
                    <span>Min: ${game.minBet}</span>
                    <span>Max: ${game.maxBet.toLocaleString()}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="max-w-7xl mx-auto px-4 pb-16">
        <h2 className="text-2xl font-bold text-white mb-6">All Game Shows</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {regularGames.map((game) => (
            <div
              key={game.id}
              className="group bg-gray-800 rounded-xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-gray-700 hover:border-purple-500"
              onClick={() => handlePlayGame(game.id)}
            >
              <div className="aspect-square bg-gradient-to-br from-gray-700 to-gray-600 flex items-center justify-center">
                <span className="text-4xl">🎲</span>
              </div>
              <div className="p-3">
                <h3 className="text-sm font-semibold text-white truncate">{game.name}</h3>
                <p className="text-xs text-purple-400">{game.provider}</p>
                <div className="mt-2 flex items-center justify-between text-xs text-gray-400">
                  <span>RTP: {game.rtp}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
