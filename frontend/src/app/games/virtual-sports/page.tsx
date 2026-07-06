'use client';

import React, { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

// Virtual sports data
const virtualSports = [
  // Football
  { id: 'vfoot-001', name: 'Virtual Football', provider: 'Pragmatic Play Virtual', rtp: '95.00%', minBet: 0.50, maxBet: 1000, category: 'football', featured: true },
  { id: 'vfoot-002', name: 'Virtual Football League', provider: 'TigerCasino Originals', rtp: '95.50%', minBet: 1.0, maxBet: 500, category: 'football', featured: true },
  { id: 'vfoot-003', name: 'World Cup Virtual', provider: 'TigerCasino Originals', rtp: '96.00%', minBet: 1.0, maxBet: 1000, category: 'football', featured: false },
  // Basketball
  { id: 'vbasket-001', name: 'Virtual Basketball', provider: 'TigerCasino Originals', rtp: '95.00%', minBet: 1.0, maxBet: 500, category: 'basketball', featured: false },
  { id: 'vbasket-002', name: 'Virtual NBA', provider: 'TigerCasino Originals', rtp: '95.50%', minBet: 1.0, maxBet: 1000, category: 'basketball', featured: false },
  // Tennis
  { id: 'vtennis-001', name: 'Virtual Tennis Open', provider: 'TigerCasino Originals', rtp: '95.00%', minBet: 0.50, maxBet: 500, category: 'tennis', featured: false },
  // Horse Racing
  { id: 'vhorse-001', name: 'Virtual Horse Racing', provider: 'TigerCasino Originals', rtp: '92.00%', minBet: 1.0, maxBet: 500, category: 'horse-racing', featured: true },
  { id: 'vhorse-002', name: 'Virtual Greyhound Racing', provider: 'TigerCasino Originals', rtp: '92.00%', minBet: 1.0, maxBet: 500, category: 'horse-racing', featured: false },
  // Motor Racing
  { id: 'vmotor-001', name: 'Virtual Formula Racing', provider: 'TigerCasino Originals', rtp: '95.00%', minBet: 1.0, maxBet: 500, category: 'motor-racing', featured: false },
];

export default function VirtualSportsPage() {
  const { isAuthenticated } = useAuth();
  const router = useRouter();
  const [searchTerm, setSearchTerm] = useState('');
  const [category, setCategory] = useState('all');
  const [filteredGames, setFilteredGames] = useState(virtualSports);

  useEffect(() => {
    let filtered = virtualSports;
    
    if (category !== 'all') {
      filtered = filtered.filter(game => game.category === category);
    }
    
    if (searchTerm) {
      filtered = filtered.filter(game =>
        game.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        game.provider.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    setFilteredGames(filtered);
  }, [searchTerm, category]);

  const handlePlayGame = (gameId: string) => {
    if (!isAuthenticated) {
      router.push('/auth/login');
      return;
    }
    router.push(`/games/virtual-sports/${gameId}`);
  };

  const featuredGames = filteredGames.filter(g => g.featured);
  const categories = ['all', 'football', 'basketball', 'tennis', 'horse-racing', 'motor-racing'];
  const categoryNames: Record<string, string> = {
    'all': 'All Sports',
    'football': '⚽ Football',
    'basketball': '🏀 Basketball',
    'tennis': '🎾 Tennis',
    'horse-racing': '🏇 Horse Racing',
    'motor-racing': '🏎️ Motor Racing',
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800">
      <div className="relative py-16 px-4 bg-gradient-to-r from-green-900 to-emerald-900">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">Virtual Sports</h1>
          <p className="text-xl text-green-200 max-w-2xl">
            24/7 sports action with realistic simulations and instant results
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Category Filter */}
        <div className="flex flex-wrap gap-2 mb-6">
          {categories.map((cat) => (
            <button
              key={cat}
              onClick={() => setCategory(cat)}
              className={`px-4 py-2 rounded-full font-medium transition-all ${
                category === cat
                  ? 'bg-green-500 text-white'
                  : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
              }`}
            >
              {categoryNames[cat]}
            </button>
          ))}
        </div>

        {/* Search */}
        <input
          type="text"
          placeholder="Search virtual sports..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full px-6 py-4 bg-gray-800 border border-gray-700 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500"
        />
      </div>

      {featuredGames.length > 0 && (
        <div className="max-w-7xl mx-auto px-4 mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Featured Virtual Sports</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {featuredGames.map((game) => (
              <div
                key={game.id}
                className="group relative bg-gray-800 rounded-2xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-green-500/30 hover:border-green-500"
                onClick={() => handlePlayGame(game.id)}
              >
                <div className="aspect-video bg-gradient-to-br from-green-600 to-emerald-600 flex items-center justify-center">
                  <span className="text-6xl">
                    {game.category === 'football' ? '⚽' : 
                     game.category === 'horse-racing' ? '🏇' : '🏆'}
                  </span>
                </div>
                <div className="p-6">
                  <h3 className="text-xl font-bold text-white mb-2">{game.name}</h3>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-green-400">{game.provider}</span>
                    <span className="text-yellow-400">RTP: {game.rtp}</span>
                  </div>
                  <div className="mt-4 flex items-center justify-between text-sm text-gray-400">
                    <span>Min: ${game.minBet}</span>
                    <span>Max: ${game.maxBet}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="max-w-7xl mx-auto px-4 pb-16">
        <h2 className="text-2xl font-bold text-white mb-6">All Virtual Sports</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {filteredGames.map((game) => (
            <div
              key={game.id}
              className="group bg-gray-800 rounded-xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-gray-700 hover:border-green-500"
              onClick={() => handlePlayGame(game.id)}
            >
              <div className="aspect-video bg-gradient-to-br from-gray-700 to-gray-600 flex items-center justify-center">
                <span className="text-4xl">
                  {game.category === 'football' ? '⚽' : 
                   game.category === 'basketball' ? '🏀' :
                   game.category === 'tennis' ? '🎾' :
                   game.category === 'horse-racing' ? '🏇' : '🏎️'}
                </span>
              </div>
              <div className="p-3">
                <h3 className="text-sm font-semibold text-white truncate">{game.name}</h3>
                <p className="text-xs text-green-400">{game.provider}</p>
                <div className="mt-2 flex items-center justify-between text-xs text-gray-400">
                  <span>RTP: {game.rtp}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 pb-16">
        <div className="bg-gray-800/50 rounded-2xl p-8 border border-gray-700">
          <h2 className="text-2xl font-bold text-white mb-4">Why Play Virtual Sports?</h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div>
              <h3 className="text-lg font-semibold text-green-400 mb-2">⏱️ 24/7 Action</h3>
              <p className="text-gray-300">Matches run every few minutes, so there's always a game to bet on.</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-green-400 mb-2">🎯 Instant Results</h3>
              <p className="text-gray-300">No waiting for real matches - get results in minutes.</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-green-400 mb-2">📊 Consistent Odds</h3>
              <p className="text-gray-300">Stable betting markets with competitive odds every time.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
