'use client';

import React from 'react';

const LIVE_DEALERS = [
  { id: 'bj-1', game: 'Blackjack', table: 'VIP Blackjack', dealers: 'Sarah', minBet: 10, maxBet: 5000, players: 3 },
  { id: 'bj-2', game: 'Blackjack', table: 'Speed Blackjack', dealers: 'Mike', minBet: 5, maxBet: 1000, players: 5 },
  { id: 'rl-1', game: 'Roulette', table: 'European Roulette', dealers: 'Emma', minBet: 1, maxBet: 5000, players: 8 },
  { id: 'bc-1', game: 'Baccarat', table: 'Baccarat', dealers: 'David', minBet: 20, maxBet: 10000, players: 6 },
  { id: 'ps-1', game: 'Poker', table: 'Casino Hold\'em', dealers: 'Lisa', minBet: 5, maxBet: 2000, players: 4 },
];

const GAME_SHOWS = [
  { id: 'crazy-time', game: 'Crazy Time', icon: '🎡', players: 12 },
  { id: 'monopoly', game: 'Monopoly Live', icon: '🎲', players: 8 },
  { id: 'lightning', game: 'Lightning Roulette', icon: '⚡', players: 10 },
];

export default function LiveCasinoPage() {
  return (
    <div className="min-h-screen bg-tiger-dark p-6">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-4xl font-heading text-gradient mb-2">🎰 Live Casino</h1>
        <p className="text-gray-400 mb-8">Real dealers, real-time action</p>
        
        <h2 className="text-2xl font-bold mb-4">Live Tables</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          {LIVE_DEALERS.map(table => (
            <div key={table.id} className="glass rounded-xl p-4">
              <div className="aspect-video bg-tiger-surface rounded-lg mb-3 flex items-center justify-center text-4xl">
                🎰
              </div>
              <h3 className="font-bold">{table.game}</h3>
              <p className="text-gray-400 text-sm">{table.table}</p>
              <div className="flex justify-between mt-2 text-sm text-gray-400">
                <span>Dealer: {table.dealers}</span>
                <span>👥 {table.players}</span>
              </div>
              <div className="flex justify-between mt-2 text-sm">
                <span className="text-green-400">${table.minBet} - ${table.maxBet}</span>
                <button className="bg-primary-500 px-4 py-1 rounded">Join</button>
              </div>
            </div>
          ))}
        </div>
        
        <h2 className="text-2xl font-bold mb-4">Game Shows</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {GAME_SHOWS.map(show => (
            <div key={show.id} className="glass rounded-xl p-4">
              <div className="aspect-video bg-tiger-surface rounded-lg mb-3 flex items-center justify-center text-6xl">
                {show.icon}
              </div>
              <h3 className="font-bold">{show.game}</h3>
              <p className="text-gray-400 text-sm">{show.players} players watching</p>
              <button className="w-full mt-3 bg-primary-500 py-2 rounded">Play Now</button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
