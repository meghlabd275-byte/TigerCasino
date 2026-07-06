'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import toast from 'react-hot-toast';

const SLOT_GAMES = [
  { id: 'pp-tiger-fortune', name: 'Tiger Fortune', provider: 'PragmaticPlay', icon: '🐯', rtp: 96.50 },
  { id: 'pp-sweet-bonanza', name: 'Sweet Bonanza', provider: 'PragmaticPlay', icon: '🍬', rtp: 96.48 },
  { id: 'net-starburst', name: 'Starburst', provider: 'NetEnt', icon: '💎', rtp: 96.09 },
];

export default function SlotsPage() {
  const [activeProvider, setActiveProvider] = useState('all');
  const [search, setSearch] = useState('');

  const games = SLOT_GAMES.filter(g => 
    (activeProvider === 'all' || g.provider === activeProvider) &&
    g.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-tiger-dark p-6">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-4xl font-heading text-gradient mb-2">🎰 Slot Games</h1>
        <p className="text-gray-400 mb-8">{games.length} games available</p>
        
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {games.map(game => (
            <motion.div key={game.id} whileHover={{ scale: 1.05 }} className="glass rounded-xl overflow-hidden cursor-pointer">
              <div className="aspect-square bg-tiger-surface flex items-center justify-center text-6xl">
                {game.icon}
              </div>
              <div className="p-3">
                <h3 className="font-bold">{game.name}</h3>
                <p className="text-gray-400 text-sm">{game.provider}</p>
                <p className="text-green-400 text-sm">RTP: {game.rtp}%</p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
}
