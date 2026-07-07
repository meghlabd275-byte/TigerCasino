'use client';

import React, { useState, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';

// Sports data with comprehensive leagues
const SPORTS = [
  { id: 'football', name: 'Football', icon: '⚽', events: 450, leagues: ['Premier League', 'La Liga', 'Bundesliga', 'Serie A', 'Ligue 1', 'MLS', 'Champions League'] },
  { id: 'basketball', name: 'Basketball', icon: '🏀', events: 180, leagues: ['NBA', 'EuroLeague', 'WNBA', 'CBA', 'NBL'] },
  { id: 'tennis', name: 'Tennis', icon: '🎾', events: 120, leagues: ['ATP', 'WTA', 'Grand Slams', 'Masters 1000'] },
  { id: 'esports', name: 'eSports', icon: '🎮', events: 280, leagues: ['CS2', 'Dota 2', 'League of Legends', 'Valorant', 'Overwatch 2'] },
  { id: 'baseball', name: 'Baseball', icon: '⚾', events: 90, leagues: ['MLB', 'NPB', 'KBO', 'CPBL'] },
  { id: 'hockey', name: 'Ice Hockey', icon: '🏒', events: 75, leagues: ['NHL', 'KHL', 'SHL', 'World Championship'] },
  { id: 'mma', name: 'MMA/UFC', icon: '🥊', events: 45, leagues: ['UFC', 'Bellator', 'ONE Championship'] },
  { id: 'boxing', name: 'Boxing', icon: '🥊', events: 30, leagues: ['Heavyweight', 'Welterweight', 'Lightweight'] },
  { id: 'cricket', name: 'Cricket', icon: '🏏', events: 150, leagues: ['IPL', 'World Cup', 'Ashes', 'Big Bash'] },
  { id: 'rugby', name: 'Rugby', icon: '🏉', events: 60, leagues: ['Super Rugby', 'Premiership', 'Six Nations'] },
  { id: 'volleyball', name: 'Volleyball', icon: '🏐', events: 50, leagues: ['Superliga', 'Serie A1', 'CEV Champions'] },
  { id: 'american_football', name: 'American Football', icon: '🏈', events: 65, leagues: ['NFL', 'NCAAF', 'CFL'] },
];

// Sample events data
const generateEvents = () => {
  const events = [];
  
  // Football
  const footballMatches = [
    { home: 'Arsenal', away: 'Liverpool', league: 'Premier League', time: 'Today 19:00' },
    { home: 'Real Madrid', away: 'Barcelona', league: 'La Liga', time: 'Today 20:00' },
    { home: 'Bayern Munich', away: 'Dortmund', league: 'Bundesliga', time: 'Tomorrow 17:30' },
    { home: 'PSG', away: 'Marseille', league: 'Ligue 1', time: 'Tomorrow 20:00' },
    { home: 'Man City', away: 'Man United', league: 'Premier League', time: 'Sat 12:30' },
    { home: 'Inter', away: 'Juventus', league: 'Serie A', time: 'Sun 19:00' },
    { home: 'Tottenham', away: 'Chelsea', league: 'Premier League', time: 'Sat 17:00' },
    { home: 'Atletico', away: 'Real Betis', league: 'La Liga', time: 'Sun 21:00' },
  ];
  
  footballMatches.forEach((m, i) => {
    events.push({
      id: `fb_${i}`,
      sport: 'football',
      sportIcon: '⚽',
      ...m,
      homeOdds: (1.5 + Math.random() * 2).toFixed(2),
      drawOdds: (2.5 + Math.random() * 1.5).toFixed(2),
      awayOdds: (1.5 + Math.random() * 2.5).toFixed(2),
      totalOver25: (1.75 + Math.random() * 0.3).toFixed(2),
      totalUnder25: (1.75 + Math.random() * 0.3).toFixed(2),
      bttsYes: (1.65 + Math.random() * 0.3).toFixed(2),
      bttsNo: (1.85 + Math.random() * 0.3).toFixed(2),
    });
  });
  
  // Basketball
  const basketballMatches = [
    { home: 'Lakers', away: 'Warriors', league: 'NBA', time: 'Today 02:00' },
    { home: 'Celtics', away: 'Heat', league: 'NBA', time: 'Today 00:30' },
    { home: 'Nets', away: '76ers', league: 'NBA', time: 'Tomorrow 00:00' },
    { home: 'Bucks', away: 'Celtics', league: 'NBA', time: 'Fri 01:30' },
  ];
  
  basketballMatches.forEach((m, i) => {
    events.push({
      id: `bb_${i}`,
      sport: 'basketball',
      sportIcon: '🏀',
      ...m,
      homeOdds: (1.7 + Math.random() * 0.6).toFixed(2),
      awayOdds: (1.7 + Math.random() * 0.6).toFixed(2),
      spreadHome: (-3.5 + Math.random() * 7).toFixed(1),
      spreadAway: (-3.5 + Math.random() * 7).toFixed(1),
      totalOver: (210 + Math.random() * 20).toFixed(1),
      totalUnder: (210 + Math.random() * 20).toFixed(1),
    });
  });
  
  // Tennis
  const tennisMatches = [
    { home: 'Djokovic', away: 'Alcaraz', league: 'ATP', time: 'Today 22:00' },
    { home: 'Sinner', away: 'Medvedev', league: 'ATP', time: 'Tomorrow 18:00' },
    { home: 'Swiatek', away: 'Gauff', league: 'WTA', time: 'Today 20:00' },
    { home: 'Fritz', away: 'Zverev', league: 'ATP', time: 'Fri 16:00' },
  ];
  
  tennisMatches.forEach((m, i) => {
    events.push({
      id: `tn_${i}`,
      sport: 'tennis',
      sportIcon: '🎾',
      ...m,
      homeOdds: (1.6 + Math.random() * 1.2).toFixed(2),
      awayOdds: (1.6 + Math.random() * 1.2).toFixed(2),
      totalSetsOver: (2.5 + Math.random() * 0.5).toFixed(1),
      totalSetsUnder: (2.5 + Math.random() * 0.5).toFixed(1),
    });
  });
  
  // eSports
  const esportsMatches = [
    { home: 'Team Liquid', away: 'Cloud9', league: 'CS2', time: 'Today 23:00' },
    { home: 'G2 Esports', away: 'Fnatic', league: 'League of Legends', time: 'Tomorrow 18:00' },
    { home: 'FaZe Clan', away: 'Natus Vincere', league: 'CS2', time: 'Today 21:00' },
    { home: 'T1', away: 'Gen.G', league: 'League of Legends', time: 'Fri 04:00' },
    { home: 'Sentinels', away: 'LOUD', league: 'Valorant', time: 'Tomorrow 03:00' },
  ];
  
  esportsMatches.forEach((m, i) => {
    events.push({
      id: `esports_${i}`,
      sport: 'esports',
      sportIcon: '🎮',
      ...m,
      homeOdds: (1.7 + Math.random() * 0.8).toFixed(2),
      awayOdds: (1.7 + Math.random() * 0.8).toFixed(2),
      mapHandicapHome: (-1.5 + Math.random() * 3).toFixed(1),
      mapHandicapAway: (-1.5 + Math.random() * 3).toFixed(1),
    });
  });
  
  return events;
};

const EVENTS_DATA = generateEvents();

const LIVE_EVENTS = [
  { id: 'live1', sport: 'football', sportIcon: '⚽', home: 'Arsenal', away: 'Liverpool', homeScore: 2, awayScore: 1, time: "67'", homeOdds: 2.15, drawOdds: 3.50, awayOdds: 3.25 },
  { id: 'live2', sport: 'basketball', sportIcon: '🏀', home: 'Lakers', away: 'Warriors', homeScore: 89, awayScore: 92, time: "Q4 8:23", homeOdds: 1.95, awayOdds: 1.88 },
  { id: 'live3', sport: 'esports', sportIcon: '🎮', home: 'FaZe', away: 'NAVI', homeScore: 1, awayScore: 0, time: "MAP 2", homeOdds: 2.10, awayOdds: 1.75 },
];

export default function SportsPage() {
  const [activeSport, setActiveSport] = useState('all');
  const [activeLeague, setActiveLeague] = useState('all');
  const [showLive, setShowLive] = useState(true);
  const [showParlay, setShowParlay] = useState(false);
  const [selectedOdds, setSelectedOdds] = useState<Record<string, { eventId: string; type: string; odds: number; name: string }>>({});
  const [parlayStake, setParlayStake] = useState(10);

  const filteredEvents = useMemo(() => {
    let events = activeSport === 'all' ? EVENTS_DATA : EVENTS_DATA.filter(e => e.sport === activeSport);
    if (activeLeague !== 'all') {
      events = events.filter(e => e.league === activeLeague);
    }
    return events;
  }, [activeSport, activeLeague]);

  const liveEvents = useMemo(() => {
    return activeSport === 'all' ? LIVE_EVENTS : LIVE_EVENTS.filter(e => e.sport === activeSport);
  }, [activeSport]);

  const availableLeagues = useMemo(() => {
    const leagues = new Set<string>();
    filteredEvents.forEach(e => leagues.add(e.league));
    return Array.from(leagues);
  }, [filteredEvents]);

  const handleSelectOdds = (eventId: string, type: string, odds: number, name: string) => {
    const key = `${eventId}_${type}`;
    if (selectedOdds[key]) {
      const newSelected = { ...selectedOdds };
      delete newSelected[key];
      setSelectedOdds(newSelected);
    } else {
      setSelectedOdds({ ...selectedOdds, [key]: { eventId, type, odds, name } });
    }
  };

  const parlayOdds = useMemo(() => {
    let total = 1.0;
    Object.values(selectedOdds).forEach(sel => {
      total *= sel.odds;
    });
    return total;
  }, [selectedOdds]);

  const parlayWin = parlayStake * parlayOdds;

  const placeBet = (eventId: string, oddsType: string, odds: number) => {
    toast.success(`Bet placed! Odds: ${odds.toFixed(2)}`);
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 lg:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl lg:text-4xl font-heading text-gradient mb-1">🏆 Sports Betting</h1>
            <p className="text-gray-400 text-sm">50,000+ markets across all sports</p>
          </div>
          
          {/* Live/In-Play Toggle */}
          <div className="flex bg-tiger-surface rounded-lg p-1">
            <button
              onClick={() => setShowLive(true)}
              className={`px-4 py-2 rounded-md transition ${showLive ? 'bg-primary-500 text-white' : 'text-gray-400'}`}
            >
              🔴 Live ({liveEvents.length})
            </button>
            <button
              onClick={() => setShowLive(false)}
              className={`px-4 py-2 rounded-md transition ${!showLive ? 'bg-primary-500 text-white' : 'text-gray-400'}`}
            >
              📅 Upcoming
            </button>
          </div>
        </div>

        {/* Live Events Banner */}
        <AnimatePresence>
          {showLive && liveEvents.length > 0 && (
            <motion.div 
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="mb-6 bg-gradient-to-r from-red-900/50 to-primary-900/50 rounded-xl p-4 border border-red-500/30"
            >
              <h2 className="text-lg font-bold mb-3 flex items-center gap-2">
                <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
                🔴 Live Now
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                {liveEvents.map(event => (
                  <div key={event.id} className="bg-tiger-dark rounded-lg p-3">
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-xs text-gray-400">{event.league}</span>
                      <span className="text-xs text-red-400 font-mono">{event.time}</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <div className="flex-1">
                        <div className="font-semibold text-sm">{event.home}</div>
                        <div className="font-semibold text-sm">{event.away}</div>
                      </div>
                      <div className="text-right">
                        <div className="text-xl font-bold text-green-400">{event.homeScore} - {event.awayScore}</div>
                      </div>
                    </div>
                    <div className="flex gap-1 mt-2">
                      <button 
                        onClick={() => placeBet(event.id, 'ml_home', parseFloat(event.homeOdds))}
                        className="flex-1 bg-primary-600 hover:bg-primary-500 text-white text-xs py-1 rounded"
                      >
                        {event.homeOdds}
                      </button>
                      {event.drawOdds && (
                        <button 
                          onClick={() => placeBet(event.id, 'draw', parseFloat(event.drawOdds))}
                          className="flex-1 bg-primary-600 hover:bg-primary-500 text-white text-xs py-1 rounded"
                        >
                          {event.drawOdds}
                        </button>
                      )}
                      <button 
                        onClick={() => placeBet(event.id, 'ml_away', parseFloat(event.awayOdds))}
                        className="flex-1 bg-primary-600 hover:bg-primary-500 text-white text-xs py-1 rounded"
                      >
                        {event.awayOdds}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Sport Categories */}
        <div className="flex gap-2 mb-4 overflow-x-auto pb-2">
          <button
            onClick={() => setActiveSport('all')}
            className={`px-4 py-2 rounded-lg whitespace-nowrap ${activeSport === 'all' ? 'bg-primary-500' : 'bg-tiger-surface'}`}
          >
            All ({EVENTS_DATA.length})
          </button>
          {SPORTS.map(sport => (
            <button
              key={sport.id}
              onClick={() => setActiveSport(sport.id)}
              className={`px-4 py-2 rounded-lg whitespace-nowrap flex items-center gap-2 ${activeSport === sport.id ? 'bg-primary-500' : 'bg-tiger-surface'}`}
            >
              {sport.icon} {sport.name}
              <span className="text-xs opacity-60">({sport.events})</span>
            </button>
          ))}
        </div>

        {/* League Filter */}
        {availableLeagues.length > 1 && (
          <div className="flex gap-2 mb-4 overflow-x-auto pb-2">
            <button
              onClick={() => setActiveLeague('all')}
              className={`px-3 py-1 rounded-full text-sm whitespace-nowrap ${activeLeague === 'all' ? 'bg-primary-500' : 'bg-tiger-surface'}`}
            >
              All Leagues
            </button>
            {availableLeagues.map(league => (
              <button
                key={league}
                onClick={() => setActiveLeague(league)}
                className={`px-3 py-1 rounded-full text-sm whitespace-nowrap ${activeLeague === league ? 'bg-primary-500' : 'bg-tiger-surface'}`}
              >
                {league}
              </button>
            ))}
          </div>
        )}

        <div className="flex gap-4">
          {/* Events List */}
          <div className="flex-1">
            {!showLive && (
              <div className="space-y-3">
                {filteredEvents.map(event => (
                  <motion.div
                    key={event.id}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="bg-tiger-surface rounded-xl p-4"
                  >
                    <div className="flex justify-between items-center mb-3">
                      <div className="flex items-center gap-2">
                        <span className="text-xl">{event.sportIcon}</span>
                        <div>
                          <div className="text-xs text-gray-400">{event.league}</div>
                          <div className="text-sm text-gray-500">{event.time}</div>
                        </div>
                      </div>
                      <span className="text-xs bg-primary-500/20 text-primary-400 px-2 py-1 rounded">
                        Upcoming
                      </span>
                    </div>
                    
                    <div className="flex justify-between items-center mb-4">
                      <div className="flex-1">
                        <div className="font-semibold">{event.home}</div>
                        <div className="font-semibold">vs</div>
                        <div className="font-semibold">{event.away}</div>
                      </div>
                    </div>

                    {/* Betting Markets */}
                    <div className="grid grid-cols-2 lg:grid-cols-4 gap-2">
                      {/* Match Result */}
                      <div>
                        <div className="text-xs text-gray-500 mb-1">Match Result</div>
                        <div className="flex flex-col gap-1">
                          <button
                            onClick={() => handleSelectOdds(event.id, 'ml_home', parseFloat(event.homeOdds), event.home)}
                            className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_ml_home`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                          >
                            {event.home} @ {event.homeOdds}
                          </button>
                          {event.drawOdds && (
                            <button
                              onClick={() => handleSelectOdds(event.id, 'draw', parseFloat(event.drawOdds), 'Draw')}
                              className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_draw`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                            >
                              Draw @ {event.drawOdds}
                            </button>
                          )}
                          <button
                            onClick={() => handleSelectOdds(event.id, 'ml_away', parseFloat(event.awayOdds), event.away)}
                            className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_ml_away`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                          >
                            {event.away} @ {event.awayOdds}
                          </button>
                        </div>
                      </div>

                      {/* Total Goals */}
                      {event.totalOver25 && (
                        <div>
                          <div className="text-xs text-gray-500 mb-1">Total Goals</div>
                          <div className="flex flex-col gap-1">
                            <button
                              onClick={() => handleSelectOdds(event.id, 'over25', parseFloat(event.totalOver25), 'Over 2.5')}
                              className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_over25`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                            >
                              Over 2.5 @ {event.totalOver25}
                            </button>
                            <button
                              onClick={() => handleSelectOdds(event.id, 'under25', parseFloat(event.totalUnder25), 'Under 2.5')}
                              className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_under25`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                            >
                              Under 2.5 @ {event.totalUnder25}
                            </button>
                          </div>
                        </div>
                      )}

                      {/* Both Teams to Score */}
                      {event.bttsYes && (
                        <div>
                          <div className="text-xs text-gray-500 mb-1">BTTS</div>
                          <div className="flex flex-col gap-1">
                            <button
                              onClick={() => handleSelectOdds(event.id, 'btts_yes', parseFloat(event.bttsYes), 'Yes')}
                              className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_btts_yes`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                            >
                              Yes @ {event.bttsYes}
                            </button>
                            <button
                              onClick={() => handleSelectOdds(event.id, 'btts_no', parseFloat(event.bttsNo), 'No')}
                              className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_btts_no`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                            >
                              No @ {event.bttsNo}
                            </button>
                          </div>
                        </div>
                      )}

                      {/* Half Time / Full Time */}
                      <div>
                        <div className="text-xs text-gray-500 mb-1">Double Chance</div>
                        <div className="flex flex-col gap-1">
                          <button
                            onClick={() => handleSelectOdds(event.id, 'dc_home', parseFloat((parseFloat(event.homeOdds) * 0.7).toFixed(2)), `${event.home} or Draw`)}
                            className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_dc_home`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                          >
                            1X @ {(parseFloat(event.homeOdds) * 0.7).toFixed(2)}
                          </button>
                          <button
                            onClick={() => handleSelectOdds(event.id, 'dc_away', parseFloat((parseFloat(event.awayOdds) * 0.7).toFixed(2)), `${event.away} or Draw`)}
                            className={`text-xs py-1.5 rounded ${selectedOdds[`${event.id}_dc_away`] ? 'bg-green-500 text-white' : 'bg-tiger-dark hover:bg-primary-600'}`}
                          >
                            X2 @ {(parseFloat(event.awayOdds) * 0.7).toFixed(2)}
                          </button>
                        </div>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            )}
          </div>

          {/* Bet Slip */}
          {(Object.keys(selectedOdds).length > 0 || showParlay) && (
            <div className="w-80 hidden lg:block">
              <div className="bg-tiger-surface rounded-xl p-4 sticky top-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-bold text-lg">🎫 Bet Slip</h3>
                  <button
                    onClick={() => setShowParlay(!showParlay)}
                    className={`text-xs px-2 py-1 rounded ${showParlay ? 'bg-primary-500' : 'bg-tiger-dark'}`}
                  >
                    Parlay
                  </button>
                </div>

                {showParlay ? (
                  <>
                    <div className="space-y-2 mb-4 max-h-60 overflow-y-auto">
                      {Object.values(selectedOdds).map((sel, idx) => (
                        <div key={idx} className="flex justify-between items-center text-sm bg-tiger-dark p-2 rounded">
                          <div>
                            <div className="font-medium">{sel.name}</div>
                            <div className="text-gray-500 text-xs">@ {sel.odds.toFixed(2)}</div>
                          </div>
                          <button
                            onClick={() => {
                              const newSelected = { ...selectedOdds };
                              delete newSelected[`${sel.eventId}_${sel.type}`];
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
                        <span className="text-gray-400">Parlay Odds:</span>
                        <span className="font-bold text-primary-400">{parlayOdds.toFixed(2)}</span>
                      </div>
                      
                      <div className="mb-4">
                        <label className="text-sm text-gray-400 block mb-1">Stake</label>
                        <input
                          type="number"
                          value={parlayStake}
                          onChange={(e) => setParlayStake(parseFloat(e.target.value) || 0)}
                          className="w-full bg-tiger-dark border border-tiger-border rounded-lg px-3 py-2 text-white"
                        />
                      </div>
                      
                      <div className="flex justify-between items-center mb-4">
                        <span className="text-gray-400">Potential Win:</span>
                        <span className="font-bold text-green-400 text-xl">{parlayWin.toFixed(2)}</span>
                      </div>
                      
                      <button
                        onClick={() => {
                          toast.success(`Parlay bet placed! Potential win: $${parlayWin.toFixed(2)}`);
                          setSelectedOdds({});
                          setParlayStake(10);
                        }}
                        className="w-full bg-green-500 hover:bg-green-600 text-white font-bold py-3 rounded-lg"
                      >
                        Place Parlay Bet
                      </button>
                    </div>
                  </>
                ) : (
                  <div className="space-y-2">
                    <p className="text-gray-400 text-sm">Switch to Parlay mode to combine multiple selections</p>
                    <button
                      onClick={() => setShowParlay(true)}
                      className="w-full bg-primary-500 hover:bg-primary-600 text-white py-2 rounded-lg text-sm"
                    >
                      Enable Parlay ({Object.keys(selectedOdds).length} selections)
                    </button>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Quick Stats */}
        <div className="mt-8 grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-primary-400">50,000+</div>
            <div className="text-gray-400 text-sm">Betting Markets</div>
          </div>
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-green-400">12+</div>
            <div className="text-gray-400 text-sm">Sports</div>
          </div>
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-yellow-400">100+</div>
            <div className="text-gray-400 text-sm">Leagues</div>
          </div>
          <div className="bg-tiger-surface rounded-xl p-4 text-center">
            <div className="text-3xl font-bold text-red-400">24/7</div>
            <div className="text-gray-400 text-sm">Live Betting</div>
          </div>
        </div>
      </div>
    </div>
  );
}
