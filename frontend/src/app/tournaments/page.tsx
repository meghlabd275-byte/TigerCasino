'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';

// Tournament configurations
const TOURNAMENTS = [
  {
    id: 'daily_slots',
    name: 'Daily Slot Challenge',
    type: 'slots',
    status: 'active',
    prizePool: 5000,
    currency: 'USDT',
    participants: 234,
    maxParticipants: 500,
    endTime: '2024-01-15T23:59:59Z',
    minBet: 0.10,
    game: 'Any Slot Game',
    description: 'Spin your way to the top!',
    rewards: [
      { position: 1, amount: 1500 },
      { position: 2, amount: 1000 },
      { position: 3, amount: 750 },
      { position: 4, amount: 500 },
      { position: 5, amount: 250 },
    ]
  },
  {
    id: 'weekly_crash',
    name: 'Weekly Crash Championship',
    type: 'crash',
    status: 'upcoming',
    prizePool: 10000,
    currency: 'USDT',
    participants: 0,
    maxParticipants: 1000,
    startTime: '2024-01-20T12:00:00Z',
    endTime: '2024-01-27T12:00:00Z',
    minBet: 1.00,
    game: 'Crash',
    description: '7 days of crash action!',
  },
  {
    id: 'vip_monthly',
    name: 'VIP Monthly Marathon',
    type: 'vip',
    status: 'upcoming',
    prizePool: 50000,
    currency: 'USDT',
    participants: 0,
    maxParticipants: 200,
    startTime: '2024-02-01T00:00:00Z',
    endTime: '2024-02-28T23:59:59Z',
    minBet: 10.00,
    game: 'Any Game',
    description: 'Exclusive VIP tournament!',
  },
  {
    id: 'live_casino',
    name: 'Live Casino Weekend',
    type: 'live',
    status: 'active',
    prizePool: 2500,
    currency: 'USDT',
    participants: 156,
    maxParticipants: 300,
    endTime: '2024-01-14T23:59:59Z',
    minBet: 5.00,
    game: 'Any Live Game',
    description: 'Play live dealer games!',
  },
];

export default function TournamentsPage() {
  const [selectedTournament, setSelectedTournament] = useState(TOURNAMENTS[0]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'bg-green-500';
      case 'upcoming': return 'bg-blue-500';
      default: return 'bg-gray-500';
    }
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl md:text-4xl font-heading text-gradient">🏆 Tournaments</h1>
          <p className="text-gray-400">Compete for massive prizes!</p>
        </div>

        {/* Featured */}
        <div className="glass rounded-2xl p-6 mb-8">
          <div className="flex flex-col md:flex-row gap-6">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <span className={`px-3 py-1 rounded-full text-xs text-white ${getStatusColor(selectedTournament.status)}`}>
                  {selectedTournament.status.toUpperCase()}
                </span>
              </div>
              <h2 className="text-3xl font-bold mb-2">{selectedTournament.name}</h2>
              <p className="text-gray-400 mb-4">{selectedTournament.description}</p>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                <div className="bg-tiger-surface/50 rounded-lg p-3 text-center">
                  <p className="text-gray-400 text-xs">Prize Pool</p>
                  <p className="font-mono text-yellow-400 text-xl">{selectedTournament.prizePool} {selectedTournament.currency}</p>
                </div>
                <div className="bg-tiger-surface/50 rounded-lg p-3 text-center">
                  <p className="text-gray-400 text-xs">Participants</p>
                  <p className="font-mono text-primary-400">{selectedTournament.participants}/{selectedTournament.maxParticipants}</p>
                </div>
                <div className="bg-tiger-surface/50 rounded-lg p-3 text-center">
                  <p className="text-gray-400 text-xs">Min Bet</p>
                  <p className="font-mono text-green-400">${selectedTournament.minBet}</p>
                </div>
                <div className="bg-tiger-surface/50 rounded-lg p-3 text-center">
                  <p className="text-gray-400 text-xs">Game</p>
                  <p className="text-white">{selectedTournament.game}</p>
                </div>
              </div>
              <button className="bg-primary-500 hover:bg-primary-600 text-white px-8 py-3 rounded-xl font-bold transition">
                Join Tournament →
              </button>
            </div>
            <div className="w-full md:w-72 bg-tiger-surface/30 rounded-xl p-4">
              <h3 className="font-bold mb-4">🏅 Prizes</h3>
              <div className="space-y-2">
                {selectedTournament.rewards?.map((reward, i) => (
                  <div key={i} className="flex justify-between items-center">
                    <span>#{reward.position}</span>
                    <span className="font-mono text-green-400">{reward.amount} {selectedTournament.currency}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* All Tournaments */}
        <h2 className="text-2xl font-bold mb-4">All Tournaments</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {TOURNAMENTS.map(tournament => (
            <motion.div
              key={tournament.id}
              whileHover={{ scale: 1.01 }}
              className="glass rounded-xl p-4 cursor-pointer"
              onClick={() => setSelectedTournament(tournament)}
            >
              <div className="flex justify-between items-start">
                <div>
                  <span className={`px-2 py-1 rounded-full text-xs text-white ${getStatusColor(tournament.status)}`}>
                    {tournament.status.toUpperCase()}
                  </span>
                  <h3 className="font-bold mt-2">{tournament.name}</h3>
                </div>
                <div className="text-right">
                  <p className="text-yellow-400 font-bold">{tournament.prizePool}</p>
                  <p className="text-xs text-gray-500">{tournament.currency}</p>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
}
