'use client';

import React, { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

// Lottery/Bingo data
const lotteryGames = [
  // Keno
  { id: 'keno-001', name: 'Classic Keno', provider: 'BGaming', rtp: '96.00%', minBet: 0.10, maxBet: 100, category: 'keno', featured: true },
  { id: 'keno-002', name: 'Instant Keno', provider: 'BGaming', rtp: '96.00%', minBet: 0.10, maxBet: 100, category: 'keno', featured: false },
  { id: 'keno-003', name: 'Multi-Race Keno', provider: 'TigerCasino Originals', rtp: '95.50%', minBet: 0.10, maxBet: 50, category: 'keno', featured: false },
  // Bingo
  { id: 'bingo-001', name: 'Bingo Blast', provider: 'Pragmatic Play', rtp: '95.00%', minBet: 0.10, maxBet: 100, category: 'bingo', featured: true },
  { id: 'bingo-002', name: '90-Ball Bingo', provider: 'TigerCasino Originals', rtp: '95.00%', minBet: 0.10, maxBet: 50, category: 'bingo', featured: false },
  { id: 'bingo-003', name: 'Speed Bingo', provider: 'TigerCasino Originals', rtp: '94.50%', minBet: 0.10, maxBet: 25, category: 'bingo', featured: false },
  // Scratch Cards
  { id: 'scratch-001', name: 'Lucky 7s', provider: 'TigerCasino Originals', rtp: '94.00%', minBet: 0.10, maxBet: 100, category: 'scratch', featured: true },
  { id: 'scratch-002', name: 'Gold Rush', provider: 'TigerCasino Originals', rtp: '94.00%', minBet: 0.10, maxBet: 50, category: 'scratch', featured: false },
  { id: 'scratch-003', name: 'Diamond Dazzle', provider: 'TigerCasino Originals', rtp: '94.50%', minBet: 0.10, maxBet: 100, category: 'scratch', featured: false },
  { id: 'scratch-004', name: 'Cash Coin', provider: 'TigerCasino Originals', rtp: '94.00%', minBet: 0.10, maxBet: 50, category: 'scratch', featured: false },
  { id: 'scratch-005', name: 'Mega Win', provider: 'TigerCasino Originals', rtp: '95.00%', minBet: 0.10, maxBet: 100, category: 'scratch', featured: false },
];

export default function LotteryPage() {
  const { isAuthenticated } = useAuth();
  const router = useRouter();
  const [searchTerm, setSearchTerm] = useState('');
  const [category, setCategory] = useState('all');
  const [filteredGames, setFilteredGames] = useState(lotteryGames);

  useEffect(() => {
    let filtered = lotteryGames;
    
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
    router.push(`/games/lottery/${gameId}`);
  };

  const featuredGames = filteredGames.filter(g => g.featured);
  const categories = ['all', 'keno', 'bingo', 'scratch'];
  const categoryNames: Record<string, string> = {
    'all': 'All Games',
    'keno': '🎱 Keno',
    'bingo': '🎯 Bingo',
    'scratch': '🎫 Scratch Cards',
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800">
      <div className="relative py-16 px-4 bg-gradient-to-r from-orange-900 to-red-900">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">Lottery & Bingo</h1>
          <p className="text-xl text-orange-200 max-w-2xl">
            Try your luck with instant wins, bingo rooms, and exciting scratch cards
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex flex-wrap gap-2 mb-6">
          {categories.map((cat) => (
            <button
              key={cat}
              onClick={() => setCategory(cat)}
              className={`px-4 py-2 rounded-full font-medium transition-all ${
                category === cat
                  ? 'bg-orange-500 text-white'
                  : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
              }`}
            >
              {categoryNames[cat]}
            </button>
          ))}
        </div>

        <input
          type="text"
          placeholder="Search lottery games..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full px-6 py-4 bg-gray-800 border border-gray-700 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-500"
        />
      </div>

      {featuredGames.length > 0 && (
        <div className="max-w-7xl mx-auto px-4 mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Featured Games</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {featuredGames.map((game) => (
              <div
                key={game.id}
                className="group relative bg-gray-800 rounded-2xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-orange-500/30 hover:border-orange-500"
                onClick={() => handlePlayGame(game.id)}
              >
                <div className="aspect-video bg-gradient-to-br from-orange-600 to-red-600 flex items-center justify-center">
                  <span className="text-6xl">
                    {game.category === 'keno' ? '🎱' : 
                     game.category === 'bingo' ? '🎯' : '🎫'}
                  </span>
                </div>
                <div className="p-6">
                  <h3 className="text-xl font-bold text-white mb-2">{game.name}</h3>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-orange-400">{game.provider}</span>
                    <span className="text-green-400">RTP: {game.rtp}</span>
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
        <h2 className="text-2xl font-bold text-white mb-6">All Lottery Games</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {filteredGames.map((game) => (
            <div
              key={game.id}
              className="group bg-gray-800 rounded-xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-gray-700 hover:border-orange-500"
              onClick={() => handlePlayGame(game.id)}
            >
              <div className="aspect-square bg-gradient-to-br from-gray-700 to-gray-600 flex items-center justify-center">
                <span className="text-4xl">
                  {game.category === 'keno' ? '🎱' : 
                   game.category === 'bingo' ? '🎯' : '🎫'}
                </span>
              </div>
              <div className="p-3">
                <h3 className="text-sm font-semibold text-white truncate">{game.name}</h3>
                <p className="text-xs text-orange-400">{game.provider}</p>
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
          <h2 className="text-2xl font-bold text-white mb-4">How to Play</h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div>
              <h3 className="text-lg font-semibold text-orange-400 mb-2">🎱 Keno</h3>
              <p className="text-gray-300">Pick 1-10 numbers from 1-80. Match more numbers to win bigger prizes!</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-orange-400 mb-2">🎯 Bingo</h3>
              <p className="text-gray-300">Get your card and mark off numbers as they're called. Complete patterns to win!</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-orange-400 mb-2">🎫 Scratch Cards</h3>
              <p className="text-gray-300">Click to reveal symbols. Match 3 identical symbols to win instant prizes!</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
