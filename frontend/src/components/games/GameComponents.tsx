'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/contexts/AuthContext';

// Types
interface GameState {
  balance: number;
  currentBet: number;
  gameData: any;
}

interface DiceGameProps {
  initialBalance?: number;
}

interface CrashGameProps {
  initialBalance?: number;
}

interface MinesGameProps {
  initialBalance?: number;
  minesCount?: number;
}

interface PlinkoGameProps {
  initialBalance?: number;
  rows?: number;
  risk?: 'low' | 'medium' | 'high';
}

// ============ DICE GAME ============
export function DiceGame({ initialBalance = 1000 }: DiceGameProps) {
  const [balance, setBalance] = useState(initialBalance);
  const [betAmount, setBetAmount] = useState(1);
  const [target, setTarget] = useState(50);
  const [direction, setDirection] = useState<'over' | 'under'>('over');
  const [isPlaying, setIsPlaying] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [serverSeed, setServerSeed] = useState('');
  const [clientSeed, setClientSeed] = useState('');

  const playDice = async () => {
    if (balance < betAmount) {
      alert('Insufficient balance!');
      return;
    }

    setIsPlaying(true);

    try {
      const roll = Math.random() * 100;
      const isWin = direction === 'over' ? roll > target : roll < target;
      
      let multiplier = 0;
      if (isWin) {
        multiplier = direction === 'over' 
          ? 100 / (100 - target)
          : target / (100 - target);
        multiplier = multiplier * 0.99;
      }

      const winAmount = isWin ? betAmount * multiplier : 0;
      
      setResult({
        roll,
        target,
        direction,
        multiplier,
        winAmount,
        isWin,
        serverSeed: serverSeed || 'generated-server-seed',
        clientSeed: clientSeed || 'client-seed'
      });

      setBalance(prev => prev - betAmount + winAmount);
    } catch (error) {
      console.error('Error playing dice:', error);
    } finally {
      setIsPlaying(false);
    }
  };

  const maxWin = (direction === 'over' 
    ? (100 - target) / target 
    : target / (100 - target)) * betAmount * 0.99;

  return (
    <div className="game-container">
      <div className="game-header">
        <h2>🎲 Dice</h2>
        <div className="balance-display">
          Balance: <span className="balance-amount">${balance.toFixed(2)}</span>
        </div>
      </div>

      <div className="dice-game">
        <div className="dice-display">
          <div className={`dice-result ${result?.isWin ? 'win' : ''}`}>
            {result?.roll?.toFixed(2) || '🎲'}
          </div>
          {result && (
            <div className={`result-text ${result.isWin ? 'win' : 'lose'}`}>
              {result.isWin ? `WIN $${result.winAmount.toFixed(2)}` : 'LOSE'}
            </div>
          )}
        </div>

        <div className="betting-panel">
          <div className="bet-amount">
            <label>Bet Amount: ${betAmount.toFixed(2)}</label>
            <input
              type="range"
              min="0.01"
              max={balance}
              step="0.01"
              value={betAmount}
              onChange={(e) => setBetAmount(parseFloat(e.target.value))}
              disabled={isPlaying}
            />
            <div className="quick-bets">
              {[0.1, 1, 10, 100].map(amt => (
                <button
                  key={amt}
                  onClick={() => setBetAmount(Math.min(amt, balance))}
                  disabled={isPlaying || amt > balance}
                >
                  ${amt}
                </button>
              ))}
            </div>
          </div>

          <div className="target-selector">
            <label>Target: {target}</label>
            <input
              type="range"
              min="2"
              max="98"
              value={target}
              onChange={(e) => setTarget(parseInt(e.target.value))}
              disabled={isPlaying}
            />
            <div className="target-info">
              <span>Win Chance: {direction === 'over' ? (100 - target).toFixed(1) : target.toFixed(1)}%</span>
              <span>Max Win: ${maxWin.toFixed(2)}</span>
            </div>
          </div>

          <div className="direction-selector">
            <button
              className={`direction-btn ${direction === 'over' ? 'active' : ''}`}
              onClick={() => setDirection('over')}
              disabled={isPlaying}
            >
              OVER {target}
            </button>
            <button
              className={`direction-btn ${direction === 'under' ? 'active' : ''}`}
              onClick={() => setDirection('under')}
              disabled={isPlaying}
            >
              UNDER {target}
            </button>
          </div>

          <button 
            className="play-btn"
            onClick={playDice}
            disabled={isPlaying || balance < betAmount}
          >
            {isPlaying ? 'Rolling...' : `ROLL DICE - $${betAmount.toFixed(2)}`}
          </button>
        </div>

        {result && (
          <div className="game-info">
            <h4>Provably Fair Verification</h4>
            <div className="seed-info">
              <p>Server Seed: <code>{result.serverSeed}</code></p>
              <p>Client Seed: <code>{result.clientSeed}</code></p>
              <p>Roll: {result.roll.toFixed(6)}</p>
            </div>
          </div>
        )}
      </div>

      <style jsx>{`
        .game-container {
          background: linear-gradient(145deg, #1a1a2e 0%, #16213e 100%);
          border-radius: 16px;
          padding: 24px;
          color: white;
          max-width: 600px;
          margin: 0 auto;
        }

        .game-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }

        .balance-amount {
          color: #00D26A;
          font-weight: bold;
          font-size: 1.2em;
        }

        .dice-display {
          text-align: center;
          padding: 40px;
          background: #0f0f1a;
          border-radius: 16px;
          margin-bottom: 24px;
        }

        .dice-result {
          font-size: 64px;
          font-weight: bold;
          color: #FFD700;
        }

        .result-text {
          font-size: 24px;
          margin-top: 16px;
        }

        .result-text.win { color: #00D26A; }
        .result-text.lose { color: #FF4757; }

        .betting-panel {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .bet-amount, .target-selector {
          background: rgba(255,255,255,0.05);
          padding: 16px;
          border-radius: 8px;
        }

        label {
          display: block;
          margin-bottom: 8px;
          color: #B0B0B0;
        }

        input[type="range"] {
          width: 100%;
          accent-color: #FF6B35;
        }

        .quick-bets {
          display: flex;
          gap: 8px;
          margin-top: 12px;
        }

        .quick-bets button {
          flex: 1;
          padding: 8px;
          background: #FF6B35;
          border: none;
          border-radius: 4px;
          color: white;
          cursor: pointer;
        }

        .target-info {
          display: flex;
          justify-content: space-between;
          margin-top: 8px;
          font-size: 0.9em;
          color: #B0B0B0;
        }

        .direction-selector {
          display: flex;
          gap: 12px;
        }

        .direction-btn {
          flex: 1;
          padding: 16px;
          background: #2a2a4a;
          border: 2px solid transparent;
          border-radius: 8px;
          color: white;
          font-size: 16px;
          cursor: pointer;
          transition: all 0.3s;
        }

        .direction-btn.active {
          border-color: #FF6B35;
          background: rgba(255,107,53,0.2);
        }

        .play-btn {
          padding: 20px;
          background: linear-gradient(135deg, #FF6B35 0%, #FF8F5A 100%);
          border: none;
          border-radius: 12px;
          color: white;
          font-size: 18px;
          font-weight: bold;
          cursor: pointer;
        }

        .play-btn:disabled {
          opacity: 0.5;
        }

        .game-info {
          margin-top: 24px;
          padding: 16px;
          background: rgba(0,0,0,0.3);
          border-radius: 8px;
        }

        .seed-info code {
          display: block;
          background: #0f0f1a;
          padding: 4px 8px;
          border-radius: 4px;
          font-size: 0.8em;
          word-break: break-all;
        }
      `}</style>
    </div>
  );
}

// ============ CRASH GAME ============
export function CrashGame({ initialBalance = 1000 }: CrashGameProps) {
  const [balance, setBalance] = useState(initialBalance);
  const [betAmount, setBetAmount] = useState(1);
  const [autoCashout, setAutoCashout] = useState(2);
  const [isInGame, setIsInGame] = useState(false);
  const [multiplier, setMultiplier] = useState(1);
  const [hasCashedOut, setHasCashedOut] = useState(false);
  const [gameStatus, setGameStatus] = useState<'waiting' | 'flying' | 'crashed'>('waiting');
  const [crashPoint, setCrashPoint] = useState(0);
  const [myBet, setMyBet] = useState(0);
  const animationRef = React.useRef<number>();

  useEffect(() => {
    if (gameStatus === 'flying') {
      const startTime = Date.now();
      const crashAt = crashPoint * 1000;

      const animate = () => {
        const elapsed = Date.now() - startTime;
        const currentMultiplier = Math.pow(1 + elapsed / 10000, elapsed / 1000);
        
        if (currentMultiplier >= crashPoint) {
          setMultiplier(crashPoint);
          setGameStatus('crashed');
        } else {
          setMultiplier(currentMultiplier);
          
          if (!hasCashedOut && currentMultiplier >= autoCashout) {
            handleCashout();
          }
          
          animationRef.current = requestAnimationFrame(animate);
        }
      };

      animationRef.current = requestAnimationFrame(animate);
    }

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [gameStatus, crashPoint]);

  const placeBet = () => {
    if (balance < betAmount) {
      alert('Insufficient balance!');
      return;
    }

    const rand = Math.random();
    let crash;
    if (rand < 0.70) {
      crash = 1 + rand * 7;
    } else {
      crash = 8 + (rand - 0.7) * 300;
    }
    crash = Math.floor(crash * 100) / 100;

    setCrashPoint(crash);
    setBalance(prev => prev - betAmount);
    setMyBet(betAmount);
    setIsInGame(true);
    setMultiplier(1);
    setHasCashedOut(false);
    setGameStatus('flying');
  };

  const handleCashout = () => {
    if (!isInGame || hasCashedOut || gameStatus === 'crashed') return;
    
    const winAmount = betAmount * multiplier * 0.97;
    setBalance(prev => prev + winAmount);
    setHasCashedOut(true);
  };

  return (
    <div className="game-container">
      <div className="game-header">
        <h2>🚀 Crash</h2>
        <div className="balance-display">
          Balance: <span className="balance-amount">${balance.toFixed(2)}</span>
        </div>
      </div>

      <div className="crash-game">
        <div className="crash-display">
          <div className={`multiplier ${gameStatus === 'flying' ? 'flying' : ''} ${gameStatus === 'crashed' ? 'crashed' : ''}`}>
            {multiplier.toFixed(2)}x
          </div>
          {gameStatus === 'crashed' && (
            <div className="crash-indicator">💥 CRASHED AT {crashPoint.toFixed(2)}x</div>
          )}
          {hasCashedOut && (
            <div className="cashed-out-indicator">
              ✅ CASHED OUT AT {multiplier.toFixed(2)}x
            </div>
          )}
        </div>

        <div className="betting-panel">
          {!isInGame ? (
            <>
              <div className="bet-amount">
                <label>Bet Amount: ${betAmount.toFixed(2)}</label>
                <input
                  type="range"
                  min="0.01"
                  max={balance}
                  step="0.01"
                  value={betAmount}
                  onChange={(e) => setBetAmount(parseFloat(e.target.value))}
                />
              </div>

              <button className="play-btn" onClick={placeBet}>
                PLACE BET - ${betAmount.toFixed(2)}
              </button>
            </>
          ) : (
            <>
              <div className="auto-cashout">
                <label>Auto Cashout: {autoCashout}x</label>
                <input
                  type="range"
                  min="1.01"
                  max="10"
                  step="0.01"
                  value={autoCashout}
                  onChange={(e) => setAutoCashout(parseFloat(e.target.value))}
                />
              </div>

              <button 
                className="cashout-btn" 
                onClick={handleCashout}
                disabled={hasCashedOut || gameStatus === 'crashed'}
              >
                CASH OUT NOW - $${(betAmount * multiplier * 0.97).toFixed(2)}
              </button>
            </>
          )}
        </div>
      </div>

      <style jsx>{`
        .game-container {
          background: linear-gradient(145deg, #1a1a2e 0%, #16213e 100%);
          border-radius: 16px;
          padding: 24px;
          color: white;
          max-width: 600px;
          margin: 0 auto;
        }

        .game-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }

        .crash-display {
          text-align: center;
          padding: 60px;
          background: #0f0f1a;
          border-radius: 16px;
          margin-bottom: 24px;
        }

        .multiplier {
          font-size: 72px;
          font-weight: bold;
          color: #00D26A;
        }

        .multiplier.flying {
          animation: pulse 0.5s infinite;
        }

        .multiplier.crashed {
          color: #FF4757;
        }

        .crash-indicator, .cashed-out-indicator {
          margin-top: 16px;
          font-size: 18px;
        }

        .crash-indicator { color: #FF4757; }
        .cashed-out-indicator { color: #00D26A; }

        .betting-panel {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .bet-amount, .auto-cashout {
          background: rgba(255,255,255,0.05);
          padding: 16px;
          border-radius: 8px;
        }

        label {
          display: block;
          margin-bottom: 8px;
          color: #B0B0B0;
        }

        input[type="range"] {
          width: 100%;
          accent-color: #FF6B35;
        }

        .play-btn, .cashout-btn {
          padding: 20px;
          border: none;
          border-radius: 12px;
          color: white;
          font-size: 18px;
          font-weight: bold;
          cursor: pointer;
        }

        .play-btn { background: linear-gradient(135deg, #FF6B35 0%, #FF8F5A 100%); }
        .cashout-btn { background: linear-gradient(135deg, #00D26A 0%, #00E676 100%); }

        .cashout-btn:disabled { opacity: 0.5; }

        @keyframes pulse {
          0%, 100% { transform: scale(1); }
          50% { transform: scale(1.05); }
        }
      `}</style>
    </div>
  );
}

// ============ MINES GAME ============
export function MinesGame({ initialBalance = 1000, minesCount = 3 }: MinesGameProps) {
  const [balance, setBalance] = useState(initialBalance);
  const [betAmount, setBetAmount] = useState(1);
  const [mines, setMines] = useState(minesCount);
  const [mineLocations, setMineLocations] = useState<number[]>([]);
  const [revealed, setRevealed] = useState<number[]>([]);
  const [gameOver, setGameOver] = useState(false);
  const [gameWon, setGameWon] = useState(false);
  const [currentMultiplier, setCurrentMultiplier] = useState(1);
  const [totalBet, setTotalBet] = useState(0);

  const startGame = () => {
    if (balance < betAmount) {
      alert('Insufficient balance!');
      return;
    }

    const minesLoc: number[] = [];
    while (minesLoc.length < mines) {
      const pos = Math.floor(Math.random() * 25);
      if (!minesLoc.includes(pos)) {
        minesLoc.push(pos);
      }
    }

    setMineLocations(minesLoc);
    setBalance(prev => prev - betAmount);
    setTotalBet(betAmount);
    setRevealed([]);
    setGameOver(false);
    setGameWon(false);
    setCurrentMultiplier(1);
  };

  const revealTile = (index: number) => {
    if (gameOver || revealed.includes(index)) return;

    const isMine = mineLocations.includes(index);

    if (isMine) {
      setGameOver(true);
      setGameWon(false);
      setRevealed([...revealed, ...mineLocations.filter(m => !revealed.includes(m))]);
    } else {
      const newRevealed = [...revealed, index];
      setRevealed(newRevealed);
      
      const multiplier = 1 + (newRevealed.length * 0.1) * (mines / 10);
      setCurrentMultiplier(multiplier);

      if (newRevealed.length === 25 - mines) {
        setGameOver(true);
        setGameWon(true);
        const winAmount = betAmount * multiplier * 0.98;
        setBalance(prev => prev + winAmount);
      }
    }
  };

  const cashOut = () => {
    if (revealed.length === 0) return;
    
    const winAmount = betAmount * currentMultiplier * 0.98;
    setBalance(prev => prev + winAmount);
    setGameOver(true);
    setGameWon(true);
  };

  return (
    <div className="game-container">
      <div className="game-header">
        <h2>💣 Mines</h2>
        <div className="balance-display">
          Balance: <span className="balance-amount">${balance.toFixed(2)}</span>
        </div>
      </div>

      <div className="mines-game">
        <div className="game-stats">
          <div className="stat">
            <span className="stat-label">Bet</span>
            <span className="stat-value">${totalBet.toFixed(2)}</span>
          </div>
          <div className="stat">
            <span className="stat-label">Multiplier</span>
            <span className="stat-value">{currentMultiplier.toFixed(2)}x</span>
          </div>
          <div className="stat">
            <span className="stat-label">Win</span>
            <span className="stat-value">${(totalBet * currentMultiplier * 0.98).toFixed(2)}</span>
          </div>
        </div>

        <div className="mines-grid">
          {Array(25).fill(null).map((_, index) => {
            const isRevealed = revealed.includes(index) || (gameOver && mineLocations.includes(index));
            const isMine = mineLocations.includes(index);
            const showMine = gameOver && isMine;

            return (
              <button
                key={index}
                className={`mine-tile ${isRevealed ? 'revealed' : ''} ${showMine ? 'mine' : ''}`}
                onClick={() => !gameOver && revealTile(index)}
                disabled={gameOver || revealed.includes(index)}
              >
                {isRevealed && !isMine ? '💎' : ''}
                {showMine ? '💣' : ''}
              </button>
            );
          })}
        </div>

        {gameOver && (
          <div className={`game-result ${gameWon ? 'win' : 'lose'}`}>
            {gameWon ? `🎉 YOU WON! $${(totalBet * currentMultiplier * 0.98).toFixed(2)}` : '💥 GAME OVER'}
          </div>
        )}

        <div className="betting-panel">
          {!totalBet || gameOver ? (
            <>
              <div className="bet-amount">
                <label>Bet Amount: ${betAmount.toFixed(2)}</label>
                <input
                  type="range"
                  min="0.01"
                  max={balance}
                  step="0.01"
                  value={betAmount}
                  onChange={(e) => setBetAmount(parseFloat(e.target.value))}
                />
              </div>
              <div className="mines-count">
                <label>Mines: {mines}</label>
                <input
                  type="range"
                  min="1"
                  max="24"
                  value={mines}
                  onChange={(e) => setMines(parseInt(e.target.value))}
                />
              </div>
              <button className="play-btn" onClick={startGame}>
                START GAME - ${betAmount.toFixed(2)}
              </button>
            </>
          ) : (
            <button className="cashout-btn" onClick={cashOut}>
              CASH OUT - $${(totalBet * currentMultiplier * 0.98).toFixed(2)}
            </button>
          )}
        </div>
      </div>

      <style jsx>{`
        .game-container {
          background: linear-gradient(145deg, #1a1a2e 0%, #16213e 100%);
          border-radius: 16px;
          padding: 24px;
          color: white;
          max-width: 500px;
          margin: 0 auto;
        }

        .game-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }

        .game-stats {
          display: flex;
          justify-content: space-around;
          margin-bottom: 24px;
          padding: 16px;
          background: rgba(0,0,0,0.3);
          border-radius: 8px;
        }

        .stat {
          text-align: center;
        }

        .stat-label {
          display: block;
          color: #B0B0B0;
          font-size: 0.9em;
        }

        .stat-value {
          display: block;
          color: #00D26A;
          font-size: 1.2em;
          font-weight: bold;
        }

        .mines-grid {
          display: grid;
          grid-template-columns: repeat(5, 1fr);
          gap: 8px;
          margin-bottom: 24px;
        }

        .mine-tile {
          aspect-ratio: 1;
          background: #2a2a4a;
          border: 2px solid #3a3a5a;
          border-radius: 8px;
          font-size: 24px;
          cursor: pointer;
        }

        .mine-tile.revealed {
          background: rgba(0, 210, 106, 0.2);
        }

        .mine-tile.mine {
          background: rgba(255, 71, 87, 0.3);
          border-color: #FF4757;
        }

        .game-result {
          text-align: center;
          font-size: 24px;
          font-weight: bold;
          padding: 16px;
          margin-bottom: 24px;
          border-radius: 8px;
        }

        .game-result.win {
          background: rgba(0, 210, 106, 0.2);
          color: #00D26A;
        }

        .game-result.lose {
          background: rgba(255, 71, 87, 0.2);
          color: #FF4757;
        }

        .betting-panel {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .bet-amount, .mines-count {
          background: rgba(255,255,255,0.05);
          padding: 16px;
          border-radius: 8px;
        }

        label {
          display: block;
          margin-bottom: 8px;
          color: #B0B0B0;
        }

        input[type="range"] {
          width: 100%;
          accent-color: #FF6B35;
        }

        .play-btn, .cashout-btn {
          padding: 20px;
          border: none;
          border-radius: 12px;
          color: white;
          font-size: 18px;
          font-weight: bold;
        }

        .play-btn { background: linear-gradient(135deg, #FF6B35 0%, #FF8F5A 100%); }
        .cashout-btn { background: linear-gradient(135deg, #00D26A 0%, #00E676 100%); }
      `}</style>
    </div>
  );
}

// ============ PLINKO GAME ============
export function PlinkoGame({ initialBalance = 1000, rows = 8, risk = 'medium' }: PlinkoGameProps) {
  const [balance, setBalance] = useState(initialBalance);
  const [betAmount, setBetAmount] = useState(1);
  const [numRows, setNumRows] = useState(rows);
  const [riskLevel, setRiskLevel] = useState(risk);
  const [isPlaying, setIsPlaying] = useState(false);
  const [result, setResult] = useState<any>(null);

  const playPlinko = async () => {
    if (balance < betAmount) {
      alert('Insufficient balance!');
      return;
    }

    setIsPlaying(true);
    setBalance(prev => prev - betAmount);

    const position = Math.floor(numRows / 2) + (Math.random() > 0.5 ? 1 : 0);
    
    const payouts = getPayouts(numRows, riskLevel);
    const multiplier = payouts[Math.min(position, payouts.length - 1)];
    const winAmount = betAmount * multiplier * 0.96;

    await new Promise(resolve => setTimeout(resolve, 1000));

    setResult({ multiplier, winAmount, position });
    setBalance(prev => prev + winAmount);
    setIsPlaying(false);
  };

  const getPayouts = (r: number, risk: string): number[] => {
    if (risk === 'low') return [1.5, 1.2, 1.0, 0.5, 0.5, 1.0, 1.2, 1.5];
    if (risk === 'high') return [10, 5, 2, 0.5, 0.5, 2, 5, 10];
    return [5, 2.5, 1.5, 0.5, 0.5, 1.5, 2.5, 5];
  };

  return (
    <div className="game-container">
      <div className="game-header">
        <h2>🎯 Plinko</h2>
        <div className="balance-display">
          Balance: <span className="balance-amount">${balance.toFixed(2)}</span>
        </div>
      </div>

      <div className="plinko-game">
        <div className="plinko-board">
          {result && (
            <div className="result-display">
              <div className={`multiplier ${result.winAmount > 0 ? 'win' : ''}`}>
                {result.multiplier.toFixed(2)}x
              </div>
              {result.winAmount > 0 && (
                <div className="win-amount">WIN ${result.winAmount.toFixed(2)}</div>
              )}
            </div>
          )}
        </div>

        <div className="betting-panel">
          <div className="bet-amount">
            <label>Bet Amount: ${betAmount.toFixed(2)}</label>
            <input
              type="range"
              min="0.01"
              max={balance}
              step="0.01"
              value={betAmount}
              onChange={(e) => setBetAmount(parseFloat(e.target.value))}
            />
          </div>

          <div className="rows-selector">
            <label>Rows: {numRows}</label>
            <input
              type="range"
              min="8"
              max="16"
              value={numRows}
              onChange={(e) => setNumRows(parseInt(e.target.value))}
            />
          </div>

          <div className="risk-selector">
            <label>Risk Level</label>
            <div className="risk-buttons">
              {(['low', 'medium', 'high'] as const).map(r => (
                <button
                  key={r}
                  className={`risk-btn ${riskLevel === r ? 'active' : ''}`}
                  onClick={() => setRiskLevel(r)}
                >
                  {r.toUpperCase()}
                </button>
              ))}
            </div>
          </div>

          <button 
            className="play-btn" 
            onClick={playPlinko}
            disabled={isPlaying || balance < betAmount}
          >
            {isPlaying ? 'DROPPING...' : `DROP BALL - $${betAmount.toFixed(2)}`}
          </button>
        </div>
      </div>

      <style jsx>{`
        .game-container {
          background: linear-gradient(145deg, #1a1a2e 0%, #16213e 100%);
          border-radius: 16px;
          padding: 24px;
          color: white;
          max-width: 500px;
          margin: 0 auto;
        }

        .game-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 24px;
        }

        .plinko-board {
          background: #0f0f1a;
          border-radius: 16px;
          padding: 24px;
          margin-bottom: 24px;
          min-height: 200px;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .result-display {
          text-align: center;
        }

        .multiplier {
          font-size: 48px;
          font-weight: bold;
          color: #FFD700;
        }

        .multiplier.win { animation: pulse 0.5s; }

        .win-amount {
          font-size: 24px;
          color: #00D26A;
          margin-top: 8px;
        }

        .betting-panel {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .bet-amount, .rows-selector, .risk-selector {
          background: rgba(255,255,255,0.05);
          padding: 16px;
          border-radius: 8px;
        }

        label {
          display: block;
          margin-bottom: 8px;
          color: #B0B0B0;
        }

        input[type="range"] {
          width: 100%;
          accent-color: #FF6B35;
        }

        .risk-buttons {
          display: flex;
          gap: 8px;
        }

        .risk-btn {
          flex: 1;
          padding: 12px;
          background: #2a2a4a;
          border: none;
          border-radius: 8px;
          color: white;
          cursor: pointer;
        }

        .risk-btn.active { background: #FF6B35; }

        .play-btn {
          padding: 20px;
          background: linear-gradient(135deg, #FF6B35 0%, #FF8F5A 100%);
          border: none;
          border-radius: 12px;
          color: white;
          font-size: 18px;
          font-weight: bold;
          cursor: pointer;
        }

        .play-btn:disabled { opacity: 0.5; }

        @keyframes pulse {
          0%, 100% { transform: scale(1); }
          50% { transform: scale(1.1); }
        }
      `}</style>
    </div>
  );
}
