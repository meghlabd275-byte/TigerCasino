'use client';

import React, { useState } from 'react';
import styles from './Sportsbook.module.css';

interface Odds {
  homeWin: number;
  draw: number;
  awayWin: number;
}

interface Match {
  id: string;
  sport: string;
  league: string;
  homeTeam: string;
  awayTeam: string;
  startTime: Date;
  status: 'upcoming' | 'live' | 'finished';
  homeScore?: number;
  awayScore?: number;
  odds: Odds;
}

interface BetSlip {
  id: string;
  matchId: string;
  selection: string;
  odds: number;
  stake: number;
}

interface SportsbookProps {
  matches: Match[];
  onPlaceBet: (bets: BetSlip[]) => Promise<void>;
  onCashout?: (betId: string) => Promise<void>;
}

export const Sportsbook: React.FC<SportsbookProps> = ({
  matches,
  onPlaceBet,
  onCashout,
}) => {
  const [betSlip, setBetSlip] = useState<BetSlip[]>([]);
  const [activeTab, setActiveTab] = useState<'live' | 'upcoming'>('upcoming');
  const [selectedSport, setSelectedSport] = useState<string>('all');
  const [stake, setStake] = useState(10);
  const [isPlacingBet, setIsPlacingBet] = useState(false);

  const addToSlip = (match: Match, selection: string, odds: number) => {
    const existing = betSlip.find(b => b.matchId === match.id && b.selection === selection);
    if (existing) return;

    const newBet: BetSlip = {
      id: `${match.id}-${selection}`,
      matchId: match.id,
      selection,
      odds,
      stake,
    };
    setBetSlip([...betSlip, newBet]);
  };

  const removeFromSlip = (betId: string) => {
    setBetSlip(betSlip.filter(b => b.id !== betId));
  };

  const updateStake = (newStake: number) => {
    setStake(newStake);
    setBetSlip(betSlip.map(b => ({ ...b, stake: newStake })));
  };

  const calculatePotentialWin = () => {
    return betSlip.reduce((acc, bet) => acc + (bet.stake * bet.odds), 0);
  };

  const handlePlaceBet = async () => {
    if (betSlip.length === 0) return;
    
    setIsPlacingBet(true);
    try {
      await onPlaceBet(betSlip);
      setBetSlip([]);
    } catch (err) {
      console.error('Failed to place bet:', err);
    } finally {
      setIsPlacingBet(false);
    }
  };

  const filteredMatches = matches.filter(m => {
    if (activeTab === 'live' && m.status !== 'live') return false;
    if (activeTab === 'upcoming' && m.status !== 'upcoming') return false;
    if (selectedSport !== 'all' && m.sport !== selectedSport) return false;
    return true;
  });

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className={styles.container}>
      <div className={styles.main}>
        <div className={styles.header}>
          <div className={styles.tabs}>
            <button
              className={`${styles.tab} ${activeTab === 'live' ? styles.active : ''}`}
              onClick={() => setActiveTab('live')}
            >
              🔴 Live
            </button>
            <button
              className={`${styles.tab} ${activeTab === 'upcoming' ? styles.active : ''}`}
              onClick={() => setActiveTab('upcoming')}
            >
              Upcoming
            </button>
          </div>
          
          <select
            className={styles.sportFilter}
            value={selectedSport}
            onChange={(e) => setSelectedSport(e.target.value)}
          >
            <option value="all">All Sports</option>
            <option value="football">Football</option>
            <option value="basketball">Basketball</option>
            <option value="tennis">Tennis</option>
            <option value="hockey">Hockey</option>
          </select>
        </div>

        <div className={styles.matches}>
          {filteredMatches.map(match => (
            <div key={match.id} className={styles.matchCard}>
              <div className={styles.matchHeader}>
                <span className={styles.league}>{match.league}</span>
                {match.status === 'live' && (
                  <span className={styles.liveIndicator}>● LIVE</span>
                )}
                <span className={styles.time}>
                  {match.status === 'live' 
                    ? `${match.homeScore} - ${match.awayScore}`
                    : formatTime(match.startTime)}
                </span>
              </div>
              
              <div className={styles.teams}>
                <span className={styles.team}>{match.homeTeam}</span>
                <span className={styles.vs}>vs</span>
                <span className={styles.team}>{match.awayTeam}</span>
              </div>
              
              <div className={styles.odds}>
                <button
                  className={`${styles.oddBtn} ${betSlip.some(b => b.matchId === match.id && b.selection === 'home') ? styles.selected : ''}`}
                  onClick={() => addToSlip(match, 'home', match.odds.homeWin)}
                >
                  <span className={styles.oddLabel}>1</span>
                  <span className={styles.oddValue}>{match.odds.homeWin.toFixed(2)}</span>
                </button>
                
                {match.odds.draw > 0 && (
                  <button
                    className={`${styles.oddBtn} ${betSlip.some(b => b.matchId === match.id && b.selection === 'draw') ? styles.selected : ''}`}
                    onClick={() => addToSlip(match, 'draw', match.odds.draw)}
                  >
                    <span className={styles.oddLabel}>X</span>
                    <span className={styles.oddValue}>{match.odds.draw.toFixed(2)}</span>
                  </button>
                )}
                
                <button
                  className={`${styles.oddBtn} ${betSlip.some(b => b.matchId === match.id && b.selection === 'away') ? styles.selected : ''}`}
                  onClick={() => addToSlip(match, 'away', match.odds.awayWin)}
                >
                  <span className={styles.oddLabel}>2</span>
                  <span className={styles.oddValue}>{match.odds.awayWin.toFixed(2)}</span>
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className={styles.slip}>
        <h3 className={styles.slipTitle}>Bet Slip</h3>
        
        {betSlip.length === 0 ? (
          <p className={styles.emptySlip}>Add selections to your bet slip</p>
        ) : (
          <>
            <div className={styles.slipItems}>
              {betSlip.map(bet => {
                const match = matches.find(m => m.id === bet.matchId);
                return (
                  <div key={bet.id} className={styles.slipItem}>
                    <div className={styles.slipItemInfo}>
                      <span className={styles.slipMatch}>
                        {match?.homeTeam} vs {match?.awayTeam}
                      </span>
                      <span className={styles.slipSelection}>
                        {bet.selection === 'home' ? '1' : bet.selection === 'draw' ? 'X' : '2'} @ {bet.odds.toFixed(2)}
                      </span>
                    </div>
                    <button
                      className={styles.removeBtn}
                      onClick={() => removeFromSlip(bet.id)}
                    >
                      ×
                    </button>
                  </div>
                );
              })}
            </div>
            
            <div className={styles.stakeInput}>
              <label>Stake</label>
              <input
                type="number"
                value={stake}
                onChange={(e) => updateStake(Number(e.target.value))}
                min={1}
              />
            </div>
            
            <div className={styles.potentialWin}>
              <span>Potential Win</span>
              <span className={styles.winAmount}>${calculatePotentialWin().toFixed(2)}</span>
            </div>
            
            <button
              className={styles.placeBetBtn}
              onClick={handlePlaceBet}
              disabled={isPlacingBet}
            >
              {isPlacingBet ? 'Placing Bet...' : `Place Bet - $${stake}`}
            </button>
          </>
        )}
      </div>
    </div>
  );
};

// Cashout component
interface CashoutProps {
  betId: string;
  originalStake: number;
  originalOdds: number;
  currentOdds: number;
  onCashout: (betId: string) => Promise<void>;
}

export const CashoutWidget: React.FC<CashoutProps> = ({
  betId,
  originalStake,
  originalOdds,
  currentOdds,
  onCashout,
}) => {
  const [isProcessing, setIsProcessing] = useState(false);

  const cashoutAmount = (currentOdds / originalOdds) * originalStake * 0.95; // 5% margin

  const handleCashout = async () => {
    setIsProcessing(true);
    try {
      await onCashout(betId);
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className={styles.cashoutWidget}>
      <div className={styles.cashoutInfo}>
        <span className={styles.cashoutLabel}>Cash Out</span>
        <span className={styles.cashoutAmount}>${cashoutAmount.toFixed(2)}</span>
      </div>
      <button
        className={styles.cashoutBtn}
        onClick={handleCashout}
        disabled={isProcessing}
      >
        {isProcessing ? 'Processing...' : 'Cash Out'}
      </button>
    </div>
  );
};

export default Sportsbook;
