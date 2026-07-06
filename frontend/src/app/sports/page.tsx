'use client';

import React, { useState } from 'react';

const SPORTS = [
  { id: 'football', name: 'Football', icon: '⚽', events: 150 },
  { id: 'basketball', name: 'Basketball', icon: '🏀', events: 80 },
  { id: 'tennis', name: 'Tennis', icon: '🎾', events: 60 },
  { id: 'esports', name: 'Esports', icon: '🎮', events: 120 },
];

const LIVE_EVENTS = [
  { id: 1, sport: 'football', home: 'Arsenal', away: 'Liverpool', homeScore: 2, awayScore: 1, time: '67\'', homeOdds: 1.85, drawOdds: 3.40, awayOdds: 4.20 },
  { id: 2, sport: 'basketball', home: 'Lakers', away: 'Warriors', homeScore: 89, awayScore: 92, time: 'Q3', homeOdds: 1.95, awayOdds: 1.88 },
];

const UPCOMING_EVENTS = [
  { id: 3, sport: 'football', home: 'Real Madrid', away: 'Barcelona', startTime: 'Tomorrow 20:00', homeOdds: 2.10, drawOdds: 3.50, awayOdds: 3.00 },
  { id: 4, sport: 'tennis', home: 'Djokovic', away: 'Alcaraz', startTime: 'Today 22:00', homeOdds: 1.75, awayOdds: 2.10 },
];

export default function SportsPage() {
  const [activeSport, setActiveSport] = useState('all');
  const [showLive, setShowLive] = useState(true);

  return (
    <div className="min-h-screen bg-tiger-dark p-6">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-4xl font-heading text-gradient mb-2">🏆 Sports Betting</h1>
        <p className="text-gray-400 mb-8">Bet on your favorite sports with crypto</p>
        
        {/* Sport Tabs */}
        <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
          <button
            onClick={() => setActiveSport('all')}
            className={`px-4 py-2 rounded-lg whitespace-nowrap ${activeSport === 'all' ? 'bg-primary-500' : 'bg-tiger-surface'}`}
          >
            All Sports
          </button>
          {SPORTS.map(sport => (
            <button
              key={sport.id}
              onClick={() => setActiveSport(sport.id)}
              className={`px-4 py-2 rounded-lg whitespace-nowrap flex items-center gap-2 ${activeSport === sport.id ? 'bg-primary-500' : 'bg-tiger-surface'}`}
            >
              {sport.icon} {sport.name}
            </button>
          ))}
        </div>

        {/* Live Toggle */}
        <div className="flex gap-4 mb-6">
          <button
            onClick={() => setShowLive(true)}
            className={`text-xl font-bold ${showLive ? 'text-green-400 border-b-2 border-green-400' : 'text-gray-400'}`}
          >
            🔴 Live ({LIVE_EVENTS.length})
          </button>
          <button
            onClick={() => setShowLive(false)}
            className={`text-xl font-bold ${!showLive ? 'text-primary-400 border-b-2 border-primary-400' : 'text-gray-400'}`}
          >
            📅 Upcoming ({UPCOMING_EVENTS.length})
          </button>
        </div>

        {/* Events */}
        {showLive ? (
          <div className="space-y-4">
            {LIVE_EVENTS.map(event => (
              <div key={event.id} className="glass rounded-xl p-4">
                <div className="flex justify-between items-center mb-3">
                  <span className="text-red-500 font-bold text-sm">🔴 LIVE - {event.time}</span>
                  <span className="text-gray-400 text-sm">{event.sport}</span>
                </div>
                <div className="flex justify-between items-center text-center">
                  <div className="flex-1">
                    <p className="text-xl font-bold">{event.home}</p>
                    <p className="text-3xl">{event.homeScore}</p>
                  </div>
                  <div className="text-gray-500 px-4">vs</div>
                  <div className="flex-1">
                    <p className="text-xl font-bold">{event.away}</p>
                    <p className="text-3xl">{event.awayScore}</p>
                  </div>
                </div>
                <div className="flex gap-2 mt-4">
                  <button className="flex-1 bg-green-500/20 hover:bg-green-500/30 text-green-400 py-2 rounded">
                    {event.homeOdds}
                  </button>
                  <button className="flex-1 bg-gray-500/20 hover:bg-gray-500/30 text-gray-300 py-2 rounded">
                    {event.drawOdds}
                  </button>
                  <button className="flex-1 bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 py-2 rounded">
                    {event.awayOdds}
                  </button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="space-y-4">
            {UPCOMING_EVENTS.map(event => (
              <div key={event.id} className="glass rounded-xl p-4">
                <div className="flex justify-between items-center mb-3">
                  <span className="text-primary-400 font-bold text-sm">{event.startTime}</span>
                  <span className="text-gray-400 text-sm">{event.sport}</span>
                </div>
                <div className="flex justify-between items-center">
                  <div className="flex-1 text-center">
                    <p className="text-xl font-bold">{event.home}</p>
                    <button className="mt-2 bg-green-500/20 hover:bg-green-500/30 text-green-400 px-4 py-1 rounded">
                      {event.homeOdds}
                    </button>
                  </div>
                  <div className="text-gray-500 px-4">vs</div>
                  <div className="flex-1 text-center">
                    <p className="text-xl font-bold">{event.away}</p>
                    <button className="mt-2 bg-blue-500/20 hover:bg-blue-500/30 text-blue-400 px-4 py-1 rounded">
                      {event.awayOdds}
                    </button>
                  </div>
                </div>
                <div className="text-center mt-3">
                  <button className="bg-gray-500/20 hover:bg-gray-500/30 text-gray-300 px-6 py-1 rounded">
                    Draw: {event.drawOdds}
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
