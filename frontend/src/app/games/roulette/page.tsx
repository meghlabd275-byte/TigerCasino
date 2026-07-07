'use client';

import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Roulette configuration
const EUROPEAN_WHEEL = [0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10, 5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26];
const AMERICAN_WHEEL = [0, 28, 9, 26, 30, 11, 7, 20, 32, 17, 5, 22, 34, 15, 3, 24, 36, 13, 1, 37, 27, 10, 25, 29, 12, 8, 19, 31, 18, 6, 21, 33, 16, 4, 23, 35, 14, 2];

// Bet types and their payouts
const BET_PAYOUTS: Record<string, number> = {
  'straight': 35,      // Single number
  'split': 17,         // Two adjacent numbers
  'street': 11,        // Three numbers in a row
  'corner': 8,         // Four numbers in a square
  'line': 5,           // Six numbers in two rows
  'column': 2,        // 12 numbers in column
  'dozen': 2,          // 12 numbers (1-12, 13-24, 25-36)
  'even_odd': 1,       // Even or Odd
  'red_black': 1,      // Red or Black
  'low_high': 1,       // 1-18 or 19-36
};

interface Bet {
  id: string;
  amount: number;
  type: string;
  value: string | number;
  result: number;
  won: boolean;
  profit: number;
  timestamp: Date;
}

export default function RouletteGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [wheelType, setWheelType] = useState<'european' | 'american'>('european');
  const [selectedBets, setSelectedBets] = useState<{type: string; value: string | number; amount: number}[]>([]);
  const [currentNumber, setCurrentNumber] = useState<number | null>(null);
  const [isSpinning, setIsSpinning] = useState(false);
  const [gameHistory, setGameHistory] = useState<number[]>([]);
  const [myBets, setMyBets] = useState<Bet[]>([]);

  // Get wheel numbers
  const getWheel = useCallback(() => {
    return wheelType === 'european' ? EUROPEAN_WHEEL : AMERICAN_WHEEL;
  }, [wheelType]);

  // Check if number is red
  const isRed = (num: number) => {
    const reds = [1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36];
    return reds.includes(num);
  };

  // Check if number is black
  const isBlack = (num: number) => {
    return num !== 0 && !isRed(num);
  };

  // Place a bet
  const placeBet = useCallback((type: string, value: string | number) => {
    const existing = selectedBets.find(b => b.type === type && b.value === value);
    if (existing) {
      setSelectedBets(prev => prev.map(b => 
        (b.type === type && b.value === value)
          ? {...b, amount: b.amount + betAmount}
          : b
      ));
    } else {
      setSelectedBets(prev => [...prev, { type, value, amount: betAmount }]);
    }
  }, [selectedBets, betAmount]);

  // Clear all bets
  const clearBets = useCallback(() => {
    setSelectedBets([]);
  }, []);

  // Spin wheel
  const spinWheel = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (selectedBets.length === 0) {
      toast.error('Place at least one bet');
      return;
    }

    const totalBet = selectedBets.reduce((sum, b) => sum + b.amount, 0);
    if (totalBet <= 0) {
      toast.error('Invalid bet amount');
      return;
    }

    setIsSpinning(true);
    
    // Call real API
    // Note: This is simplified to just use the first bet for demo purposes
    const mainBet = selectedBets[0];

    fetch(`/api/games/roulette/bet`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        amount: totalBet,
        bet_type: mainBet.type,
        selected_number: typeof mainBet.value === 'number' ? mainBet.value : 0
      })
    })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        toast.error(data.error);
        setIsSpinning(false);
        return;
      }

      const finalNumber = data.number;
      setCurrentNumber(finalNumber);

      const won = data.win_amount > 0;
      const profit = data.win_amount - totalBet;

      const newBet: Bet = {
        id: Date.now().toString(),
        amount: totalBet,
        type: mainBet.type,
        value: mainBet.value,
        result: finalNumber,
        won,
        profit,
        timestamp: new Date()
      };

      setMyBets(prev => [newBet, ...prev].slice(0, 50));
      setGameHistory(prev => [finalNumber, ...prev].slice(0, 20));

      if (won) {
        toast.success(`Won $${data.win_amount.toFixed(2)}!`);
      } else {
        toast.error(`Lost $${totalBet.toFixed(2)}`);
      }

      setIsSpinning(false);
      setSelectedBets([]);
    })
    .catch(err => {
      toast.error('Failed to connect to game server');
      setIsSpinning(false);
    });
  }, [selectedBets, isAuthenticated, getWheel, betAmount]);

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Roulette</h1>
            <p className="text-gray-400">Place your bets and spin the wheel!</p>
          </div>
          <div className="flex gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Type</span>
              <p className="text-xl font-mono text-primary-500">{wheelType === 'european' ? 'European' : 'American'}</p>
            </div>
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Total Bet</span>
              <p className="text-xl font-mono text-yellow-400">
                ${selectedBets.reduce((sum, b) => sum + b.amount, 0).toFixed(2)}
              </p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Wheel Display */}
            <div className="glass rounded-2xl p-6">
              <div className="flex justify-center mb-4">
                <motion.div
                  animate={isSpinning ? { rotate: 360 } : {}}
                  transition={{ repeat: Infinity, duration: 0.5, ease: "linear" }}
                  className="w-48 h-48 rounded-full bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center"
                >
                  <div className="w-40 h-40 rounded-full bg-tiger-dark flex items-center justify-center">
                    <span className="text-4xl font-bold text-white">
                      {currentNumber !== null ? (currentNumber === 37 ? '00' : currentNumber) : wheelType === 'european' ? '0' : '00'}
                    </span>
                  </div>
                </motion.div>
              </div>
              
              {/* Result */}
              {currentNumber !== null && !isSpinning && (
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className="text-center"
                >
                  <p className="text-gray-400">
                    {isRed(currentNumber) ? '🔴 Red' : isBlack(currentNumber) ? '⚫ Black' : '🟢 Green'}
                  </p>
                  <p className="text-2xl font-bold">
                    {currentNumber === 0 || currentNumber === 37 ? 'ZERO' : 
                     currentNumber % 2 === 0 ? 'Even' : 'Odd'}
                  </p>
                </motion.div>
              )}
            </div>

            {/* Betting Table */}
            <div className="glass rounded-2xl p-4">
              {/* Wheel Type */}
              <div className="flex gap-2 mb-4">
                <button
                  onClick={() => { setWheelType('european'); clearBets(); }}
                  className={`flex-1 py-2 rounded-lg font-bold transition ${
                    wheelType === 'european'
                      ? 'bg-primary-500 text-white'
                      : 'bg-tiger-surface text-gray-400'
                  }`}
                >
                  European (0)
                </button>
                <button
                  onClick={() => { setWheelType('american'); clearBets(); }}
                  className={`flex-1 py-2 rounded-lg font-bold transition ${
                    wheelType === 'american'
                      ? 'bg-primary-500 text-white'
                      : 'bg-tiger-surface text-gray-400'
                  }`}
                >
                  American (0, 00)
                </button>
              </div>

              {/* Number Grid */}
              <div className="grid grid-cols-7 gap-1 mb-4">
                {/* Row 1: 0, 00 for American */}
                {wheelType === 'american' && (
                  <button
                    onClick={() => placeBet('straight', 37)}
                    className="aspect-square rounded bg-green-600 text-white font-bold text-sm hover:bg-green-500"
                  >
                    00
                  </button>
                )}
                <button
                  onClick={() => placeBet('straight', 0)}
                  className="aspect-square rounded bg-green-600 text-white font-bold hover:bg-green-500"
                >
                  0
                </button>
                
                {/* Regular numbers */}
                {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12].map(num => (
                  <button
                    key={num}
                    onClick={() => placeBet('straight', num)}
                    className={`aspect-square rounded font-bold text-sm hover:opacity-80 ${
                      isRed(num) ? 'bg-red-600 text-white' : 'bg-gray-800 text-white'
                    }`}
                  >
                    {num}
                  </button>
                ))}
              </div>

              {/* Outside Bets */}
              <div className="space-y-2">
                {/* Dozens */}
                <div className="grid grid-cols-3 gap-2">
                  <button onClick={() => placeBet('dozen', 1)} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-sm">1st 12</button>
                  <button onClick={() => placeBet('dozen', 2)} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-sm">2nd 12</button>
                  <button onClick={() => placeBet('dozen', 3)} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-sm">3rd 12</button>
                </div>
                
                {/* Even/Odd, Red/Black, Low/High */}
                <div className="grid grid-cols-6 gap-2">
                  <button onClick={() => placeBet('low_high', 'low')} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-xs">1-18</button>
                  <button onClick={() => placeBet('even_odd', 'even')} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-xs">Even</button>
                  <button onClick={() => placeBet('red_black', 'red')} className="py-2 bg-red-600 rounded hover:bg-red-500 text-xs text-white">Red</button>
                  <button onClick={() => placeBet('red_black', 'black')} className="py-2 bg-gray-800 rounded hover:bg-gray-700 text-xs text-white">Black</button>
                  <button onClick={() => placeBet('even_odd', 'odd')} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-xs">Odd</button>
                  <button onClick={() => placeBet('low_high', 'high')} className="py-2 bg-tiger-surface rounded hover:bg-primary-500/30 text-xs">19-36</button>
                </div>
              </div>
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {/* Bet Amount */}
              <div className="mb-4">
                <label className="block text-gray-400 text-sm mb-2">Bet Amount</label>
                <div className="flex gap-2 mb-3">
                  {[0.10, 1, 10, 100].map(amount => (
                    <button
                      key={amount}
                      onClick={() => setBetAmount(amount)}
                      className={`flex-1 py-2 rounded-lg font-mono transition ${
                        betAmount === amount
                          ? 'bg-primary-500 text-white'
                          : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
                      }`}
                    >
                      ${amount}
                    </button>
                  ))}
                </div>
                <input
                  type="number"
                  value={betAmount}
                  onChange={(e) => setBetAmount(parseFloat(e.target.value) || 0)}
                  className="w-full bg-tiger-surface border border-gray-700 rounded-lg px-4 py-3 font-mono text-lg focus:border-primary-500 focus:outline-none"
                  step={0.01}
                  min={0.01}
                />
              </div>

              <div className="flex gap-4">
                <button
                  onClick={clearBets}
                  className="flex-1 py-3 bg-red-500/20 hover:bg-red-500/40 text-red-400 rounded-xl font-bold transition"
                >
                  Clear Bets
                </button>
                <button
                  onClick={spinWheel}
                  disabled={isSpinning || selectedBets.length === 0}
                  className={`flex-1 py-3 rounded-xl font-bold text-xl transition-all ${
                    isSpinning || selectedBets.length === 0
                      ? 'bg-gray-700 cursor-not-allowed'
                      : 'bg-green-500 hover:bg-green-600'
                  }`}
                >
                  {isSpinning ? 'Spinning...' : 'Spin'}
                </button>
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Current Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Current Bets</h3>
              <div className="space-y-2 max-h-40 overflow-y-auto">
                {selectedBets.map((bet, index) => (
                  <div key={index} className="p-2 bg-tiger-surface/50 rounded flex justify-between text-sm">
                    <span className="text-gray-400">{bet.type}: {bet.value}</span>
                    <span className="text-yellow-400">${bet.amount.toFixed(2)}</span>
                  </div>
                ))}
                {selectedBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No bets placed</p>
                )}
              </div>
            </div>

            {/* History */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">History</h3>
              <div className="flex flex-wrap gap-2">
                {gameHistory.map((num, index) => (
                  <div
                    key={index}
                    className={`w-8 h-8 rounded flex items-center justify-center text-xs font-bold ${
                      isRed(num) ? 'bg-red-600 text-white' : 
                      num === 0 || num === 37 ? 'bg-green-600 text-white' :
                      'bg-gray-800 text-white'
                    }`}
                  >
                    {num}
                  </div>
                ))}
                {gameHistory.length === 0 && (
                  <p className="text-gray-500 text-center w-full py-4">No history</p>
                )}
              </div>
            </div>

            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Bets</h3>
              <div className="space-y-2 max-h-60 overflow-y-auto">
                {myBets.slice(0, 10).map(bet => (
                  <div 
                    key={bet.id} 
                    className={`p-3 rounded-lg ${
                      bet.won ? 'bg-green-500/20' : 'bg-red-500/20'
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className="text-gray-400 text-xs">{bet.type}</span>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="text-sm">Result: {bet.result}</span>
                      <span className={bet.won ? 'text-green-400' : 'text-red-400'}>
                        {bet.won ? `+$${bet.profit.toFixed(2)}` : `-$${bet.amount.toFixed(2)}`}
                      </span>
                    </div>
                  </div>
                ))}
                {myBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No bets yet</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
