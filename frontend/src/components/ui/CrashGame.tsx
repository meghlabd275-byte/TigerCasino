'use client';

import React, { useState, useEffect, useRef } from 'react';
import styles from './CrashGame.module.css';

interface CrashGameProps {
  gameId: string;
  minBet: number;
  maxBet: number;
  currency: string;
  onPlaceBet: (amount: number, autoCashout?: number) => Promise<void>;
  onCashout: () => Promise<void>;
}

export const CrashGame: React.FC<CrashGameProps> = ({
  gameId,
  minBet,
  maxBet,
  currency,
  onPlaceBet,
  onCashout,
}) => {
  const [multiplier, setMultiplier] = useState(1.0);
  const [gameState, setGameState] = useState<'waiting' | 'flying' | 'crashed'>('waiting');
  const [betAmount, setBetAmount] = useState(minBet);
  const [hasBet, setHasBet] = useState(false);
  const [hasCashedOut, setHasCashedOut] = useState(false);
  const [cashoutMultiplier, setCashoutMultiplier] = useState(0);
  const [betHistory, setBetHistory] = useState<number[]>([]);
  const [error, setError] = useState<string | null>(null);
  
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationRef = useRef<number>();

  // Game loop
  useEffect(() => {
    if (gameState === 'flying') {
      const startTime = Date.now();
      const startMultiplier = multiplier;
      
      const animate = () => {
        const elapsed = (Date.now() - startTime) / 1000;
        // Exponential growth: e^(0.3*t)
        const newMultiplier = startMultiplier * Math.exp(0.3 * elapsed);
        
        setMultiplier(newMultiplier);
        
        // Auto-crash at random point (provably fair in production)
        if (newMultiplier > 100 || Math.random() < 0.001) {
          setGameState('crashed');
          if (!hasCashedOut && hasBet) {
            setBetHistory(prev => [...prev, 0]);
          }
          setTimeout(() => {
            setGameState('waiting');
            setMultiplier(1.0);
            setHasBet(false);
            setHasCashedOut(false);
          }, 2000);
          return;
        }
        
        animationRef.current = requestAnimationFrame(animate);
      };
      
      animationRef.current = requestAnimationFrame(animate);
    }

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [gameState, hasBet, hasCashedOut]);

  // Canvas drawing
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Draw background
    ctx.fillStyle = '#0F0F1A';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Draw grid
    ctx.strokeStyle = 'rgba(255, 107, 53, 0.1)';
    ctx.lineWidth = 1;
    
    for (let i = 0; i < canvas.width; i += 40) {
      ctx.beginPath();
      ctx.moveTo(i, 0);
      ctx.lineTo(i, canvas.height);
      ctx.stroke();
    }
    
    for (let i = 0; i < canvas.height; i += 40) {
      ctx.beginPath();
      ctx.moveTo(0, i);
      ctx.lineTo(canvas.width, i);
      ctx.stroke();
    }

    // Draw multiplier line
    if (gameState !== 'waiting') {
      const x = 50;
      const maxHeight = canvas.height - 50;
      const height = Math.min((multiplier / 100) * maxHeight, maxHeight);
      const y = canvas.height - 50 - height;

      // Line
      ctx.strokeStyle = '#FF6B35';
      ctx.lineWidth = 3;
      ctx.beginPath();
      ctx.moveTo(x, canvas.height - 50);
      ctx.lineTo(x, y);
      ctx.stroke();

      // Glow
      ctx.shadowColor = '#FF6B35';
      ctx.shadowBlur = 20;
      ctx.stroke();
      ctx.shadowBlur = 0;

      // Point
      ctx.fillStyle = '#FF6B35';
      ctx.beginPath();
      ctx.arc(x, y, 8, 0, Math.PI * 2);
      ctx.fill();

      // Multiplier text
      ctx.fillStyle = '#FFFFFF';
      ctx.font = 'bold 24px Orbitron';
      ctx.fillText(`${multiplier.toFixed(2)}x`, x + 20, y + 8);
    }

    // Draw "waiting" state
    if (gameState === 'waiting') {
      ctx.fillStyle = '#FFFFFF';
      ctx.font = 'bold 20px Rajdhani';
      ctx.textAlign = 'center';
      ctx.fillText('Waiting for next round...', canvas.width / 2, canvas.height / 2);
    }

    // Draw crash point
    if (gameState === 'crashed') {
      ctx.fillStyle = '#FF4757';
      ctx.font = 'bold 32px Orbitron';
      ctx.textAlign = 'center';
      ctx.fillText(`CRASHED @ ${multiplier.toFixed(2)}x`, canvas.width / 2, canvas.height / 2);
    }
  }, [multiplier, gameState]);

  const handlePlaceBet = async () => {
    if (betAmount < minBet || betAmount > maxBet) {
      setError(`Bet must be between ${minBet} and ${maxBet}`);
      return;
    }

    setError(null);
    try {
      await onPlaceBet(betAmount);
      setHasBet(true);
      setGameState('flying');
    } catch (err) {
      setError('Failed to place bet');
    }
  };

  const handleCashout = async () => {
    if (!hasBet || hasCashedOut) return;

    try {
      await onCashout();
      setHasCashedOut(true);
      setCashoutMultiplier(multiplier);
      setBetHistory(prev => [...prev, multiplier]);
    } catch (err) {
      setError('Failed to cash out');
    }
  };

  const quickBetAmounts = [minBet, minBet * 10, minBet * 100, maxBet];

  return (
    <div className={styles.container}>
      <div className={styles.gameArea}>
        <canvas 
          ref={canvasRef} 
          width={600} 
          height={400}
          className={styles.canvas}
        />
        
        <div className={styles.controls}>
          {gameState === 'waiting' && !hasBet && (
            <div className={styles.betControls}>
              <label className={styles.label}>Bet Amount</label>
              <div className={styles.quickBets}>
                {quickBetAmounts.map((amount) => (
                  <button
                    key={amount}
                    className={`${styles.quickBetBtn} ${betAmount === amount ? styles.active : ''}`}
                    onClick={() => setBetAmount(amount)}
                  >
                    {amount}
                  </button>
                ))}
              </div>
              <input
                type="number"
                value={betAmount}
                onChange={(e) => setBetAmount(Number(e.target.value))}
                min={minBet}
                max={maxBet}
                className={styles.betInput}
              />
              <button className={styles.placeBetBtn} onClick={handlePlaceBet}>
                Place Bet
              </button>
            </div>
          )}

          {hasBet && gameState === 'flying' && !hasCashedOut && (
            <button className={styles.cashoutBtn} onClick={handleCashout}>
              Cash Out @ {multiplier.toFixed(2)}x
            </button>
          )}

          {hasCashedOut && (
            <div className={styles.cashedOut}>
              Cashed out at {cashoutMultiplier.toFixed(2)}x
            </div>
          )}

          {error && <div className={styles.error}>{error}</div>}
        </div>
      </div>

      <div className={styles.history}>
        <h4>History</h4>
        <div className={styles.historyList}>
          {betHistory.slice(-10).reverse().map((multiplier, i) => (
            <span 
              key={i} 
              className={`${styles.historyItem} ${multiplier > 0 ? styles.win : styles.loss}`}
            >
              {multiplier > 0 ? `${multiplier.toFixed(2)}x` : 'CRASH'}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
};

export default CrashGame;
