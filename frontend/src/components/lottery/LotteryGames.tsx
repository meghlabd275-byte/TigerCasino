'use client';

import React, { useState } from 'react';

export function KenoGame({ initialBalance = 1000 }) {
  const [balance, setBalance] = useState(initialBalance);
  const [betAmount, setBetAmount] = useState(1);
  const [selected, setSelected] = useState<number[]>([]);
  const [spots, setSpots] = useState(10);
  const [isPlaying, setIsPlaying] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [drawn, setDrawn] = useState<number[]>([]);

  const numbers = Array.from({ length: 80 }, (_, i) => i + 1);

  const toggleNumber = (num: number) => {
    if (isPlaying) return;
    if (selected.includes(num)) {
      setSelected(selected.filter(n => n !== num));
    } else if (selected.length < spots) {
      setSelected([...selected, num]);
    }
  };

  const playKeno = async () => {
    if (selected.length !== spots) {
      alert('Please select exactly ' + spots + ' numbers');
      return;
    }
    if (balance < betAmount) {
      alert('Insufficient balance!');
      return;
    }

    setIsPlaying(true);
    setBalance(prev => prev - betAmount);

    const drawnNumbers: number[] = [];
    const used = new Set<number>();
    
    while (drawnNumbers.length < 20) {
      const num = Math.floor(Math.random() * 80) + 1;
      if (!used.has(num)) {
        used.add(num);
        drawnNumbers.push(num);
        setDrawn([...drawnNumbers]);
        await new Promise(r => setTimeout(r, 200));
      }
    }

    let matches = 0;
    selected.forEach(num => {
      if (drawnNumbers.includes(num)) matches++;
    });

    const paytable: any = {
      1: {1: 3.0}, 2: {1: 1.0, 2: 9.0}, 3: {1: 1.0, 2: 2.0, 3: 27.0},
      4: {1: 0.5, 2: 1.0, 3: 4.0, 4: 80.0}, 5: {1: 0.5, 2: 2.0, 3: 5.0, 4: 10.0, 5: 400.0},
      6: {1: 0.5, 2: 1.5, 3: 3.0, 4: 5.0, 5: 20.0, 6: 1000.0},
      7: {1: 0.5, 2: 1.0, 3: 2.0, 4: 5.0, 5: 20.0, 6: 100.0, 7: 3000.0},
      8: {1: 0.5, 2: 1.0, 3: 2.0, 4: 4.0, 5: 10.0, 6: 50.0, 7: 500.0, 8: 10000.0},
    };

    const multiplier = paytable[spots]?.[matches] || 0;
    const winAmount = betAmount * multiplier;

    setResult({ matches, multiplier, winAmount });
    setBalance(prev => prev + winAmount);
    setIsPlaying(false);
  };

  return (
    <div className="keno-game">
      <h2>🎱 Keno</h2>
      <div className="balance">Balance: ${balance.toFixed(2)}</div>
      <div className="spots-selector">
        <label>Spots: {spots}</label>
        <input type="range" min="1" max="10" value={spots} onChange={(e) => { setSpots(parseInt(e.target.value)); setSelected([]); }} disabled={isPlaying} />
      </div>
      <div className="keno-board">
        <div className="numbers-grid">
          {numbers.map(num => (
            <button key={num} className={'keno-number ' + (selected.includes(num) ? 'selected' : '') + (drawn.includes(num) ? (selected.includes(num) ? 'match' : 'drawn') : '')} onClick={() => toggleNumber(num)} disabled={isPlaying}>{num}</button>
          ))}
        </div>
      </div>
      <div className="bet-controls">
        <div className="bet-amount">
          <label>Bet: ${betAmount.toFixed(2)}</label>
          <input type="range" min="0.1" max={balance} step="0.1" value={betAmount} onChange={(e) => setBetAmount(parseFloat(e.target.value))} disabled={isPlaying} />
        </div>
        <button className="play-btn" onClick={playKeno} disabled={isPlaying || selected.length !== spots}>
          {isPlaying ? 'Drawing...' : 'Play ' + spots + ' Spots - $' + betAmount.toFixed(2)}
        </button>
      </div>
      {result && <div className="result-panel"><h3>Matches: {result.matches} / {spots}</h3><p className={result.winAmount > 0 ? 'win' : 'lose'}>{result.winAmount > 0 ? 'WIN $' + result.winAmount.toFixed(2) : 'No win'}</p></div>}
    </div>
  );
}
