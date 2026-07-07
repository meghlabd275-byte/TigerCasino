'use client';

import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Keno configuration
const TOTAL_NUMBERS = 80;
const MIN_PICKS = 1;
const MAX_PICKS = 10;
const DRAW_COUNT = 20;

// Payout table for Keno (multiplier based on hits)
const PAYOUT_TABLE: Record<number, number[]> = {
  1: [3], // 1 pick: 3x for 1 hit
  2: [0, 9], // 2 picks: 9x for 2 hits
  3: [0, 2, 27], // 3 picks: 27x for 3 hits
  4: [0, 1, 4, 80], // 4 picks
  5: [0, 0, 2, 10, 400], // 5 picks
  6: [0, 0, 1, 5, 50, 1000], // 6 picks
  7: [0, 0, 0, 2, 20, 100, 3000], // 7 picks
  8: [0, 0, 0, 1, 10, 50, 500, 10000], // 8 picks
  9: [0, 0, 0, 0, 5, 25, 200, 2000, 20000], // 9 picks
  10: [0, 0, 0, 0, 2, 10, 100, 1000, 10000, 50000], // 10 picks
};

interface Bet {
  id: string;
  amount: number;
  picks: number[];
  hits: number[];
  totalHits: number;
  multiplier: number;
  won: boolean;
  profit: number;
  timestamp: Date;
}

export default function KenoGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [selectedNumbers, setSelectedNumbers] = useState<number[]>([]);
  const [drawnNumbers, setDrawnNumbers] = useState<number[]>([]);
  const [gameState, setGameState] = useState<'picking' | 'drawing' | 'completed'>('picking');
  const [myBets, setMyBets] = useState<Bet[]>([]);

  // Toggle number selection
  const toggleNumber = useCallback((num: number) => {
    if (gameState !== 'picking') return;
    
    setSelectedNumbers(prev => {
      if (prev.includes(num)) {
        return prev.filter(n => n !== num);
      }
      if (prev.length >= MAX_PICKS) {
        toast.error(`Maximum ${MAX_PICKS} numbers can be selected`);
        return prev;
      }
      return [...prev, num];
    });
  }, [gameState]);

  // Clear selection
  const clearSelection = useCallback(() => {
    setSelectedNumbers([]);
    setDrawnNumbers([]);
    setGameState('picking');
  }, []);

  // Draw numbers
  const drawNumbers = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (selectedNumbers.length < MIN_PICKS) {
      toast.error(`Select at least ${MIN_PICKS} number(s)`);
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    setGameState('drawing');
    
    // Simulate drawing animation
    let drawIndex = 0;
    const drawInterval = setInterval(() => {
      // Generate random number not already drawn
      let num: number;
      do {
        num = Math.floor(Math.random() * TOTAL_NUMBERS) + 1;
      } while (drawnNumbers.includes(num));
      
      setDrawnNumbers(prev => [...prev, num]);
      drawIndex++;
      
      if (drawIndex >= DRAW_COUNT) {
        clearInterval(drawInterval);
        
        // Calculate hits
        const hits = selectedNumbers.filter(n => drawnNumbers.includes(n));
        const totalHits = hits.length;
        
        // Calculate payout
        const payoutRow = PAYOUT_TABLE[selectedNumbers.length] || [];
        const multiplier = payoutRow[totalHits] || 0;
        const won = multiplier > 0;
        const profit = won ? (betAmount * multiplier) - betAmount : -betAmount;
        
        // Record bet
        const newBet: Bet = {
          id: Date.now().toString(),
          amount: betAmount,
          picks: [...selectedNumbers],
          hits,
          totalHits,
          multiplier,
          won,
          profit,
          timestamp: new Date()
        };
        
        setMyBets(prev => [newBet, ...prev].slice(0, 50));
        
        if (won) {
          toast.success(`Won $${(betAmount * multiplier).toFixed(2)} with ${totalHits} hits!`);
        } else {
          toast.error(`No winning hits. Better luck next time!`);
        }
        
        setGameState('completed');
      }
    }, 150);
  }, [selectedNumbers, betAmount, isAuthenticated, drawnNumbers]);

  // Quick pick - random selection
  const quickPick = useCallback(() => {
    const count = selectedNumbers.length || 5; // Default to 5 picks
    const nums: number[] = [];
    while (nums.length < count) {
      const num = Math.floor(Math.random() * TOTAL_NUMBERS) + 1;
      if (!nums.includes(num)) {
        nums.push(num);
      }
    }
    setSelectedNumbers(nums);
  }, [selectedNumbers.length]);

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Keno</h1>
            <p className="text-gray-400">Pick your lucky numbers and win!</p>
          </div>
          <div className="flex gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Selected</span>
              <p className="text-xl font-mono text-primary-500">{selectedNumbers.length}</p>
            </div>
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Drawn</span>
              <p className="text-xl font-mono text-green-400">{drawnNumbers.length}/{DRAW_COUNT}</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Keno Board */}
            <div className="glass rounded-2xl p-4">
              <div className="grid grid-cols-10 gap-2 mb-6">
                {Array.from({ length: TOTAL_NUMBERS }, (_, i) => i + 1).map(num => {
                  const isSelected = selectedNumbers.includes(num);
                  const isDrawn = drawnNumbers.includes(num);
                  const isHit = isSelected && isDrawn;
                  
                  return (
                    <motion.button
                      key={num}
                      whileHover={{ scale: 1.1 }}
                      whileTap={{ scale: 0.95 }}
                      onClick={() => toggleNumber(num)}
                      disabled={gameState !== 'picking'}
                      className={`
                        aspect-square rounded-lg font-bold text-sm transition-all
                        ${isHit 
                          ? 'bg-green-500 text-white' 
                          : isDrawn 
                            ? 'bg-yellow-500/50 text-white'
                            : isSelected 
                              ? 'bg-primary-500 text-white' 
                              : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
                        }
                        ${gameState !== 'picking' ? 'cursor-not-allowed' : ''}
                      `}
                    >
                      {num}
                    </motion.button>
                  );
                })}
              </div>

              {/* Results */}
              {gameState === 'completed' && (
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className={`text-center p-4 rounded-xl ${
                    selectedNumbers.filter(n => drawnNumbers.includes(n)).length > 0
                      ? 'bg-green-500/20 border border-green-500/30'
                      : 'bg-red-500/20 border border-red-500/30'
                  }`}
                >
                  <p className="text-xl font-bold">
                    {selectedNumbers.filter(n => drawnNumbers.includes(n)).length} hits!
                  </p>
                </motion.div>
              )}
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {/* Pick Count */}
              <div className="mb-6">
                <label className="block text-gray-400 text-sm mb-3">
                  Numbers Picked: <span className="text-white text-lg">{selectedNumbers.length}</span>
                </label>
                <div className="flex gap-2">
                  {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(num => (
                    <button
                      key={num}
                      onClick={() => { clearSelection(); setSelectedNumbers([]); }}
                      className={`flex-1 py-2 rounded-lg font-bold transition ${
                        selectedNumbers.length === num
                          ? 'bg-primary-500 text-white'
                          : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
                      }`}
                    >
                      {num}
                    </button>
                  ))}
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-4 mb-6">
                <button
                  onClick={quickPick}
                  disabled={gameState !== 'picking'}
                  className="flex-1 py-3 bg-tiger-surface hover:bg-primary-500/30 rounded-xl font-bold transition"
                >
                  🎲 Quick Pick
                </button>
                <button
                  onClick={clearSelection}
                  disabled={gameState !== 'picking'}
                  className="flex-1 py-3 bg-tiger-surface hover:bg-red-500/30 rounded-xl font-bold transition"
                >
                  🗑️ Clear
                </button>
              </div>

              {/* Bet Amount */}
              <div className="mb-6">
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

              {/* Payout Info */}
              <div className="mb-6 p-4 bg-tiger-surface/50 rounded-lg">
                <p className="text-gray-400 text-sm mb-2">
                  {selectedNumbers.length} picks - Payouts: {PAYOUT_TABLE[selectedNumbers.length]?.filter(p => p > 0).join(', ') || 'N/A'}x
                </p>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-400">Potential Win:</span>
                  <span className="text-green-400 font-mono">
                    ${(betAmount * (PAYOUT_TABLE[selectedNumbers.length]?.[selectedNumbers.length] || 0)).toFixed(2)}
                  </span>
                </div>
              </div>

              {/* Play Button */}
              <button
                onClick={drawNumbers}
                disabled={gameState !== 'picking' || selectedNumbers.length < MIN_PICKS}
                className={`w-full py-4 rounded-xl font-bold text-xl transition-all ${
                  gameState !== 'picking' || selectedNumbers.length < MIN_PICKS
                    ? 'bg-gray-700 cursor-not-allowed'
                    : 'bg-primary-500 hover:bg-primary-600'
                }`}
              >
                {gameState === 'drawing' ? 'Drawing...' : gameState === 'completed' ? 'Play Again' : 'Draw Numbers'}
              </button>
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* My Numbers */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Your Picks</h3>
              <div className="flex flex-wrap gap-2">
                {selectedNumbers.map(num => (
                  <div
                    key={num}
                    className={`w-10 h-10 rounded-lg flex items-center justify-center font-bold ${
                      drawnNumbers.includes(num)
                        ? 'bg-green-500 text-white'
                        : 'bg-primary-500/50 text-white'
                    }`}
                  >
                    {num}
                  </div>
                ))}
                {selectedNumbers.length === 0 && (
                  <p className="text-gray-500 text-center w-full py-4">Select numbers to play</p>
                )}
              </div>
            </div>

            {/* Drawn Numbers */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Drawn Numbers</h3>
              <div className="flex flex-wrap gap-2">
                {drawnNumbers.map(num => (
                  <div
                    key={num}
                    className={`w-10 h-10 rounded-lg flex items-center justify-center font-bold ${
                      selectedNumbers.includes(num)
                        ? 'bg-green-500 text-white'
                        : 'bg-yellow-500/50 text-white'
                    }`}
                  >
                    {num}
                  </div>
                ))}
                {drawnNumbers.length === 0 && (
                  <p className="text-gray-500 text-center w-full py-4">Numbers will appear here</p>
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
                      <span className="text-gray-400 text-xs">{bet.picks.length} picks</span>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="text-sm">{bet.totalHits} hits</span>
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
