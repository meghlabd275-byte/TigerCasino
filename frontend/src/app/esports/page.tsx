'use client';

import React, { useState, useMemo, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';

// Esports Games
const ESPORTS_GAMES = [
  { id: 'lol', name: 'League of Legends', icon: '⚔️', color: 'from-yellow-500 to-orange-500' },
  { id: 'cs2', name: 'Counter-Strike 2', icon: '🔫', color: 'from-orange-600 to-red-600' },
  { id: 'dota2', name: 'Dota 2', icon: '🛡️', color: 'from-red-600 to-red-800' },
  { id: 'valorant', name: 'Valorant', icon: '🎯', color: 'from-red-500 to-pink-500' },
  { id: 'overwatch', name: 'Overwatch 2', icon: '🦸', color: 'from-orange-500 to-red-400' },
  { id: 'rl', name: 'Rocket League', icon: '🚗', color: 'from-blue-500 to-cyan-400' },
];

// Major Tournaments
const TOURNAMENTS = [
  { id: 'worlds', game: 'lol', name: 'World Championship', prize: '$2,500,000', teams: 20, status: 'live' },
  { id: 'major_paris', game: 'cs2', name: 'Paris Major', prize: '$1,250,000', teams: 24, status: 'upcoming' },
  { id: 'ti12', game: 'dota2', name: 'The International', prize: '$15,000,000', teams: 18, status: 'upcoming' },
  { id: 'champs', game: 'valorant', name: 'Champions', prize: '$1,000,000', teams: 16, status: 'live' },
  { id: 'rlcs', game: 'rl', name: 'World Championship', prize: '$2,100,000', teams: 16, status: 'upcoming' },
];

// Teams data
const TEAMS: Record<string, { name: string; logo: string; wins: number; losses: number }> = {
  't1': { name: 'T1', logo: '🔴', wins: 45, losses: 12 },
  'geng': { name: 'Gen.G', logo: '🟣', wins: 42, losses: 15 },
  'fnc': { name: 'Fnatic', logo: '🟠', wins: 38, losses: 18 },
  'g2': { name: 'G2 Esports', logo: '⚫', wins: 36, losses: 20 },
  'navi': { name: 'Natus Vincere', logo: '🟡', wins: 48, losses: 10 },
  'faze': { name: 'FaZe Clan', logo: '🔴', wins: 35, losses: 22 },
  'liquid': { name: 'Team Liquid', logo: '🔵', wins: 40, losses: 17 },
  'c9': { name: 'Cloud9', logo: '⚪', wins: 37, losses: 19 },
  'sentinels': { name: 'Sentinels', logo: '🔴', wins: 30, losses: 15 },
  'loud': { name: 'LOUD', logo: '🟢', wins: 42, losses: 13 },
  'drx': { name: 'DRX', logo: '⚪', wins: 35, losses: 18 },
  'eg': { name: 'Evil Geniuses', logo: '⚫', wins: 32, losses: 20 },
};

// Generate live matches
const generateMatches = () => {
  const matches = [];
  const gameList = ['lol', 'cs2', 'dota2', 'valorant', 'rl'];
  const teamKeys = Object.keys(TEAMS);
  
  for (let i = 0; i < 30; i++) {
    const game = gameList[Math.floor(Math.random() * gameList.length)];
    const homeIdx = Math.floor(Math.random() * teamKeys.length);
    let awayIdx = Math.floor(Math.random() * teamKeys.length);
    while (awayIdx === homeIdx) awayIdx = Math.floor(Math.random() * teamKeys.length);
    
    const isLive = Math.random() < 0.3;
    const homeScore = isLive ? Math.floor(Math.random() * 3) : 0;
    const awayScore = isLive ? Math.floor(Math.random() * 3) : 0;
    
    matches.push({
      id: `match_${i}`,
      game,
      homeTeam: teamKeys[homeIdx],
      awayTeam: teamKeys[awayIdx],
      homeScore,
      awayScore,
      homeOdds: (1.5 + Math.random() * 2).toFixed(2),
      awayOdds: (1.5 + Math.random() * 2).toFixed(2),
      map: isLive ? `Map ${Math.floor(Math.random() * 5) + 1}` : 'Best of 3',
      status: isLive ? 'live' : 'upcoming',
      startTime: isLive ? 'LIVE' : `${Math.floor(Math.random() * 12) + 1}:00`,
      tournament: TOURNAMENTS[Math.floor(Math.random() * TOURNAMENTS.length)].name,
    });
  }
  
  return matches;
};

const MATCHES = generateMatches();

export default function EsportsPage() {
  const [activeGame, setActiveGame] = useState('all');
  const [activeTournament, setActiveTournament] = useState('all');
  const [showLiveOnly, setShowLiveOnly] = useState(true);
  const [selectedOdds, setSelectedOdds] = useState<Record<string, { matchId: string; team: string; odds: number }>>({});
  const [stake, setStake] = useState(10);
  const [matches, setMatches] = useState(MATCHES);
  const [animateScore, setAnimateScore] = useState<string | null>(null);

  // Filter matches
  const filteredMatches = useMemo(() => {
    let filtered = activeGame === 'all' ? matches : matches.filter(m => m.game === activeGame);
    if (activeTournament !== 'all') {
      filtered = filtered.filter(m => m.tournament === activeTournament);
    }
    if (showLiveOnly) {
      filtered = filtered.filter(m => m.status === 'live');
    }
    return filtered;
  }, [matches, activeGame, activeTournament, showLiveOnly]);

  // Live matches count
  const liveCount = useMemo(() => matches.filter(m => m.status === 'live').length, [matches]);

  // Calculate parlay odds
  const totalOdds = useMemo(() => {
    let odds = 1;
    Object.values(selectedOdds).forEach(sel => odds *= sel.odds);
    return odds;
  }, [selectedOdds]);

  const potentialWin = stake * totalOdds;

  // Handle odds selection
  const handleSelectOdds = (matchId: string, team: string, odds: number) => {
    const key = `${matchId}_${team}`;
    if (selectedOdds[key]) {
      const newSelected = { ...selectedOdds };
      delete newSelected[key];
      setSelectedOdds(newSelected);
    } else {
      setSelectedOdds({ ...selectedOdds, [key]: { matchId, team, odds } });
    }
  };

  // Place bet
  const placeBet = (matchId: string, team: string, odds: number) => {
    const key = `${matchId}_${team}`;
    if (selectedOdds[key]) {
      toast.success(`Bet placed on ${TEAMS[team].name} @ ${odds}x! Potential win: $${(stake * odds).toFixed(2)}`);
      const newSelected = { ...selectedOdds };
      delete newSelected[key];
      setSelectedOdds(newSelected);
    }
  };

  // Simulate live score updates
  useEffect(() => {
    const interval = setInterval(() => {
      setMatches(prev => prev.map(m => {
        if (m.status === 'live' && Math.random() < 0.1) {
          const scoringTeam = Math.random() < 0.5 ? 'homeScore' : 'awayScore';
          setAnimateScore(m.id);
          setTimeout(() => setAnimateScore(null), 500);
          return { ...m, [scoringTeam]: m[scoringTeam] + 1 };
        }
        return m;
      }));
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
          <div>
            <h1 className="text-4xl font-heading text-gradient flex items-center gap-3">
              🎮 Esports Betting
            </h1>
            <p className="text-gray-400 mt-1">
              Live esports betting on League of Legends, CS2, Dota 2, Valorant & more
            </p>
          </div>
          <div className="flex items-center gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Live Matches</span>
              <p className="text-2xl font-bold text-red-400 animate-pulse">{liveCount}</p>
            </div>
          </div>
        </div>

        {/* Game Filters */}
        <div className="flex flex-wrap gap-3 mb-6">
          <button
            onClick={() => setActiveGame('all')}
            className={`px-4 py-2 rounded-lg font-medium transition ${
              activeGame === 'all' ? 'bg-primary-500 text-white' : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
            }`}
          >
            All Games
          </button>
          {ESPORTS_GAMES.map(game => (
            <button
              key={game.id}
              onClick={() => setActiveGame(game.id)}
              className={`px-4 py-2 rounded-lg font-medium transition flex items-center gap-2 ${
                activeGame === game.id ? 'bg-primary-500 text-white' : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
              }`}
            >
              <span>{game.icon}</span>
              <span>{game.name}</span>
            </button>
          ))}
        </div>

        {/* Live Toggle */}
        <div className="flex items-center gap-4 mb-6">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={showLiveOnly}
              onChange={(e) => setShowLiveOnly(e.target.checked)}
              className="w-5 h-5 rounded bg-tiger-surface border-gray-600 text-primary-500 focus:ring-primary-500"
            />
            <span className="text-gray-300">Show Live Only</span>
          </label>
          {liveCount > 0 && (
            <span className="text-red-400 text-sm flex items-center gap-1">
              <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
              {liveCount} matches live now
            </span>
          )}
        </div>

        {/* Tournaments Banner */}
        <div className="mb-6">
          <h3 className="text-lg font-bold mb-3">🏆 Major Tournaments</h3>
          <div className="flex flex-wrap gap-3">
            {TOURNAMENTS.map(tournament => (
              <div
                key={tournament.id}
                className={`bg-gradient-to-r ${ESPORTS_GAMES.find(g => g.id === tournament.game)?.color || 'from-gray-600 to-gray-700'} rounded-xl px-4 py-3 flex items-center gap-4`}
              >
                <div>
                  <p className="font-bold text-white">{tournament.name}</p>
                  <p className="text-sm text-white/80">{ESPORTS_GAMES.find(g => g.id === tournament.game)?.icon} {tournament.game.toUpperCase()}</p>
                </div>
                <div className="text-right">
                  <p className="font-bold text-yellow-400">{tournament.prize}</p>
                  <p className="text-xs text-white/70">{tournament.teams} teams</p>
                </div>
                <div className={`px-3 py-1 rounded-full text-xs font-bold ${
                  tournament.status === 'live' ? 'bg-red-500 text-white' : 'bg-green-500 text-white'
                }`}>
                  {tournament.status.toUpperCase()}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-4 gap-6">
          {/* Matches */}
          <div className="xl:col-span-3 space-y-4">
            {filteredMatches.map(match => (
              <motion.div
                key={match.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className={`glass rounded-xl p-4 ${match.status === 'live' ? 'border border-red-500/30' : ''}`}
              >
                {/* Match Header */}
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <span className="text-2xl">{ESPORTS_GAMES.find(g => g.id === match.game)?.icon}</span>
                    <span className="text-gray-400 text-sm">{match.tournament}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    {match.status === 'live' ? (
                      <span className="flex items-center gap-1 text-red-400 text-sm font-bold">
                        <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
                        LIVE • {match.map}
                      </span>
                    ) : (
                      <span className="text-gray-400 text-sm">{match.startTime}</span>
                    )}
                  </div>
                </div>

                {/* Teams and Odds */}
                <div className="grid grid-cols-3 items-center gap-4">
                  {/* Home Team */}
                  <div 
                    className={`text-center p-3 rounded-lg cursor-pointer transition ${
                      selectedOdds[`${match.id}_${match.homeTeam}`] ? 'bg-green-500/30 border border-green-500' : 'bg-tiger-surface hover:bg-primary-500/20'
                    }`}
                    onClick={() => handleSelectOdds(match.id, match.homeTeam, parseFloat(match.homeOdds))}
                  >
                    <div className="text-3xl mb-1">{TEAMS[match.homeTeam].logo}</div>
                    <div className="font-bold text-white">{TEAMS[match.homeTeam].name}</div>
                    {match.status === 'live' && (
                      <div className={`text-2xl font-bold mt-1 ${animateScore === match.id ? 'text-green-400 scale-125' : 'text-white'}`}>
                        {match.homeScore}
                      </div>
                    )}
                  </div>

                  {/* VS / Score */}
                  <div className="text-center">
                    {match.status === 'live' ? (
                      <div className="text-gray-400 text-sm">vs</div>
                    ) : (
                      <div className="text-gray-500 text-sm">Best of 3</div>
                    )}
                  </div>

                  {/* Away Team */}
                  <div 
                    className={`text-center p-3 rounded-lg cursor-pointer transition ${
                      selectedOdds[`${match.id}_${match.awayTeam}`] ? 'bg-green-500/30 border border-green-500' : 'bg-tiger-surface hover:bg-primary-500/20'
                    }`}
                    onClick={() => handleSelectOdds(match.id, match.awayTeam, parseFloat(match.awayOdds))}
                  >
                    <div className="text-3xl mb-1">{TEAMS[match.awayTeam].logo}</div>
                    <div className="font-bold text-white">{TEAMS[match.awayTeam].name}</div>
                    {match.status === 'live' && (
                      <div className={`text-2xl font-bold mt-1 ${animateScore === match.id ? 'text-green-400 scale-125' : 'text-white'}`}>
                        {match.awayScore}
                      </div>
                    )}
                  </div>
                </div>

                {/* Odds Row */}
                <div className="flex justify-center gap-4 mt-3 pt-3 border-t border-tiger-border">
                  <button
                    onClick={() => placeBet(match.id, match.homeTeam, parseFloat(match.homeOdds))}
                    className={`px-6 py-2 rounded-lg font-bold transition ${
                      selectedOdds[`${match.id}_${match.homeTeam}`] 
                        ? 'bg-green-500 text-white' 
                        : 'bg-tiger-surface hover:bg-primary-500 text-primary-400'
                    }`}
                  >
                    {match.homeOdds}x
                  </button>
                  <button
                    onClick={() => placeBet(match.id, match.awayTeam, parseFloat(match.awayOdds))}
                    className={`px-6 py-2 rounded-lg font-bold transition ${
                      selectedOdds[`${match.id}_${match.awayTeam}`] 
                        ? 'bg-green-500 text-white' 
                        : 'bg-tiger-surface hover:bg-primary-500 text-primary-400'
                    }`}
                  >
                    {match.awayOdds}x
                  </button>
                </div>
              </motion.div>
            ))}

            {filteredMatches.length === 0 && (
              <div className="text-center py-12">
                <div className="text-4xl mb-4">🎮</div>
                <p className="text-gray-400">No matches available</p>
              </div>
            )}
          </div>

          {/* Bet Slip */}
          <div className="xl:col-span-1">
            <div className="bg-tiger-surface rounded-xl p-4 sticky top-4">
              <h3 className="font-bold text-lg mb-4">🎫 Bet Slip</h3>
              
              {Object.keys(selectedOdds).length > 0 ? (
                <>
                  <div className="space-y-2 mb-4 max-h-60 overflow-y-auto">
                    {Object.entries(selectedOdds).map(([key, sel]) => (
                      <div key={key} className="flex justify-between items-center text-sm bg-tiger-dark p-2 rounded">
                        <div>
                          <div className="font-medium">{TEAMS[sel.team]?.name || sel.team}</div>
                          <div className="text-gray-500 text-xs">@ {sel.odds.toFixed(2)}</div>
                        </div>
                        <button
                          onClick={() => {
                            const newSelected = { ...selectedOdds };
                            delete newSelected[key];
                            setSelectedOdds(newSelected);
                          }}
                          className="text-red-400 text-xs"
                        >
                          ✕
                        </button>
                      </div>
                    ))}
                  </div>
                  
                  <div className="border-t border-tiger-border pt-4">
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-gray-400">Total Odds:</span>
                      <span className="font-bold text-primary-400">{totalOdds.toFixed(2)}</span>
                    </div>
                    
                    <div className="mb-4">
                      <label className="text-sm text-gray-400 block mb-1">Stake ($)</label>
                      <input
                        type="number"
                        value={stake}
                        onChange={(e) => setStake(parseFloat(e.target.value) || 0)}
                        className="w-full bg-tiger-dark border border-tiger-border rounded-lg px-3 py-2 text-white"
                      />
                    </div>
                    
                    <div className="flex justify-between items-center mb-4">
                      <span className="text-gray-400">Potential Win:</span>
                      <span className="font-bold text-green-400 text-xl">${potentialWin.toFixed(2)}</span>
                    </div>
                    
                    <button
                      onClick={() => {
                        toast.success(`Parlay bet placed! Potential win: $${potentialWin.toFixed(2)}`);
                        setSelectedOdds({});
                      }}
                      disabled={Object.keys(selectedOdds).length === 0}
                      className="w-full bg-green-500 hover:bg-green-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold py-3 rounded-lg"
                    >
                      Place {Object.keys(selectedOdds).length} Selection(s)
                    </button>
                  </div>
                </>
              ) : (
                <p className="text-gray-400 text-sm">Select odds to place your bets</p>
              )}
            </div>

            {/* Stats */}
            <div className="mt-4 bg-tiger-surface rounded-xl p-4">
              <h4 className="font-bold mb-3">📊 Esports Stats</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-400">Total Markets</span>
                  <span className="text-primary-400">5,000+</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Live Now</span>
                  <span className="text-red-400">{liveCount}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Tournaments</span>
                  <span className="text-primary-400">{TOURNAMENTS.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Teams</span>
                  <span className="text-primary-400">{Object.keys(TEAMS).length}+</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
