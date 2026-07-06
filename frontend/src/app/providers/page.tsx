'use client';

import React from 'react';
import { useRouter } from 'next/navigation';

// Provider data
const providers = [
  {
    id: 'evolution',
    name: 'Evolution Gaming',
    logo: '🎰',
    description: 'Leading live casino provider with game shows and table games',
    games: 150,
    categories: ['Live Casino', 'Game Shows', 'Blackjack', 'Roulette', 'Baccarat'],
    featured: true,
  },
  {
    id: 'pragmatic',
    name: 'Pragmatic Play',
    logo: '🎲',
    description: 'Multi-product content provider with slots, live casino, bingo',
    games: 300,
    categories: ['Slots', 'Live Casino', 'Game Shows', 'Bingo', 'Virtual Sports'],
    featured: true,
  },
  {
    id: 'playngo',
    name: "Play'n GO",
    logo: '🎮',
    description: 'Innovative slot games with unique mechanics and features',
    games: 200,
    categories: ['Slots'],
    featured: true,
  },
  {
    id: 'netent',
    name: 'NetEnt',
    logo: '⭐',
    description: 'Pioneer in online casino gaming with classic titles',
    games: 200,
    categories: ['Slots', 'Table Games', 'Jackpots'],
    featured: true,
  },
  {
    id: 'bgaming',
    name: 'BGaming',
    logo: '💎',
    description: 'Crypto-friendly games with unique mechanics',
    games: 100,
    categories: ['Slots', 'Crash', 'Dice', 'Plinko', 'Mines'],
    featured: false,
  },
  {
    id: 'spribe',
    name: 'Spribe',
    logo: '🚀',
    description: 'Innovative crash games and lottery-style entertainment',
    games: 50,
    categories: ['Crash', 'Arcade', 'Dice', 'Keno'],
    featured: false,
  },
  {
    id: 'hacksaw',
    name: 'Hacksaw Gaming',
    logo: '🔪',
    description: 'Award-winning scratch cards and innovative slots',
    games: 80,
    categories: ['Slots', 'Scratch Cards'],
    featured: false,
  },
  {
    id: 'nolimit',
    name: 'Nolimit City',
    logo: '🔥',
    description: 'High-volatility slots with innovative bonus features',
    games: 50,
    categories: ['Slots'],
    featured: false,
  },
  {
    id: 'relax',
    name: 'Relax Gaming',
    logo: '🧩',
    description: 'B2B gaming supplier with slots and table games',
    games: 100,
    categories: ['Slots', 'Table Games', 'Jackpots'],
    featured: false,
  },
  {
    id: 'pushgaming',
    name: 'Push Gaming',
    logo: '🎯',
    description: 'Mobile-first gaming with innovative slot mechanics',
    games: 40,
    categories: ['Slots'],
    featured: false,
  },
  {
    id: 'quickspin',
    name: 'Quickspin',
    logo: '⚡',
    description: 'Swedish slot provider with high-quality graphics',
    games: 60,
    categories: ['Slots'],
    featured: false,
  },
  {
    id: 'yggdrasil',
    name: 'Yggdrasil',
    logo: '🌳',
    description: 'Innovative slots with unique bonus mechanics',
    games: 80,
    categories: ['Slots', 'Jackpots'],
    featured: false,
  },
];

export default function ProvidersPage() {
  const router = useRouter();

  const featuredProviders = providers.filter(p => p.featured);

  const handleProviderClick = (providerId: string) => {
    router.push(`/providers/${providerId}`);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800">
      <div className="relative py-16 px-4 bg-gradient-to-r from-indigo-900 to-purple-900">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">Game Providers</h1>
          <p className="text-xl text-indigo-200 max-w-2xl">
            Partnering with the world's best game providers to bring you 10,000+ games
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-12">
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-8">
          {providers.map((provider) => (
            <div
              key={provider.id}
              onClick={() => handleProviderClick(provider.id)}
              className="bg-gray-800 rounded-xl p-4 cursor-pointer hover:bg-gray-700 transition-colors border border-gray-700 hover:border-purple-500"
            >
              <div className="text-4xl mb-2">{provider.logo}</div>
              <h3 className="font-semibold text-white">{provider.name}</h3>
              <p className="text-sm text-gray-400">{provider.games}+ games</p>
            </div>
          ))}
        </div>
      </div>

      {featuredProviders.length > 0 && (
        <div className="max-w-7xl mx-auto px-4 mb-12">
          <h2 className="text-2xl font-bold text-white mb-6">Featured Providers</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {featuredProviders.map((provider) => (
              <div
                key={provider.id}
                className="group bg-gray-800 rounded-2xl overflow-hidden hover:transform hover:scale-105 transition-all duration-300 cursor-pointer border border-purple-500/30 hover:border-purple-500"
                onClick={() => handleProviderClick(provider.id)}
              >
                <div className="p-8 bg-gradient-to-br from-purple-600 to-indigo-600">
                  <div className="text-6xl mb-4">{provider.logo}</div>
                  <h3 className="text-2xl font-bold text-white">{provider.name}</h3>
                  <p className="text-purple-200 mt-2">{provider.games}+ Games</p>
                </div>
                <div className="p-6">
                  <p className="text-gray-300 mb-4">{provider.description}</p>
                  <div className="flex flex-wrap gap-2">
                    {provider.categories.slice(0, 3).map((cat) => (
                      <span key={cat} className="px-3 py-1 bg-purple-900/50 text-purple-300 text-xs rounded-full">
                        {cat}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="max-w-7xl mx-auto px-4 pb-16">
        <h2 className="text-2xl font-bold text-white mb-6">All Providers</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {providers.map((provider) => (
            <div
              key={provider.id}
              className="bg-gray-800 rounded-xl p-6 cursor-pointer hover:bg-gray-700 transition-colors border border-gray-700 hover:border-purple-500"
              onClick={() => handleProviderClick(provider.id)}
            >
              <div className="flex items-start gap-4">
                <div className="text-4xl">{provider.logo}</div>
                <div className="flex-1">
                  <h3 className="text-lg font-bold text-white">{provider.name}</h3>
                  <p className="text-sm text-gray-400 mt-1">{provider.games}+ games</p>
                  <div className="flex flex-wrap gap-1 mt-3">
                    {provider.categories.slice(0, 2).map((cat) => (
                      <span key={cat} className="px-2 py-0.5 bg-gray-700 text-gray-300 text-xs rounded">
                        {cat}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 pb-16">
        <div className="bg-gray-800/50 rounded-2xl p-8 border border-gray-700">
          <h2 className="text-2xl font-bold text-white mb-4">Why Our Providers?</h2>
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-lg font-semibold text-purple-400 mb-2">10,000+ Games</h3>
              <p className="text-gray-300">Massive selection of slots, table games, and live casino</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-purple-400 mb-2">Fair Play</h3>
              <p className="text-gray-300">All games certified fair by independent auditors</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-purple-400 mb-2">Top RTP</h3>
              <p className="text-gray-300">Competitive return-to-player percentages</p>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-purple-400 mb-2">Global Brands</h3>
              <p className="text-gray-300">World's leading game developers</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
