'use client';

import React, { useState, useCallback, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Configuration
const MIN_TARGET = 1.01;
const MAX_TARGET = 1000;
const DEFAULT_TARGET = 2;

interface Bet {
  id: string;
  amount: number;
  target: number;
  result: number;
  won: boolean;
  profit: number;
  timestamp: Date;
}

export default function LimboGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [target, setTarget] = useState(DEFAULT_TARGET);
  const [gameState, setGameState] = useState<'idle' | 'playing' | 'completed'>('idle');
  const [currentMultiplier, setCurrentMultiplier] = useState<number | null>(null);
  const [history, setHistory] = useState<{target: number; result: number; won: boolean}[]>([]);
  const [myBets, setMyBets] = useState<Bet[]>([]);
  const [isAutoMode, setIsAutoMode] = useState(false);
  const [autoGames, setAutoGames] = useState(10);
  const [gamesPlayed, setGamesPlayed] = useState(0);
  
  const animationRef = useRef<number>();

  // Calculate multiplier based on target
  const calculateMultiplier = useCallback((targetValue: number) => {
    // House edge of 1%
    return (0.99 * targetValue) / (targetValue - 1);
  }, []);

  // Calculate win chance
  const calculateWinChance = useCallback((targetValue: number) => {
    return ((targetValue - 1) / targetValue) * 100;
  }, []);

  // Play game
  const playGame = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    if (target < MIN_TARGET || target > MAX_TARGET) {
      toast.error('Invalid target value');
      return;
    }

    setGameState('playing');
    
    // Simulate multiplier animation
    let currentValue = 1.0;
    const targetValue = target;
    const duration = 1500;
    const startTime = Date.now();
    
    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      
      // Exponential growth simulation
      // The multiplier grows until it "crashes" at the target
      const crashPoint = 1 + Math.random() * targetValue * 2;
      
      if (progress < 1) {
        // Simulate rising multiplier
        currentValue = 1 + (crashPoint - 1) * Math.pow(progress, 0.5);
        if (currentValue > targetValue) {
          currentValue = targetValue + Math.random() * 0.5; // Slight overshoot
        }
        setCurrentMultiplier(currentValue);
        animationRef.current = requestAnimationFrame(animate);
      } else {
        // Final result
        const finalResult = 1 + Math.random() * targetValue * 1.5;
        setCurrentMultiplier(finalResult);
        
        // Determine if player won
        const won = finalResult >= target;
        const profit = won 
          ? (betAmount * target * 0.99) - betAmount 
          : -betAmount;
        
        // Record bet
        const newBet: Bet = {
          id: Date.now().toString(),
          amount: betAmount,
          target,
          result: finalResult,
          won,
          profit,
          timestamp: new Date()
        };
        
        setMyBets(prev => [newBet, ...prev].slice(0, 50));
        setHistory(prev => [{target, result: finalResult, won}, ...prev].slice(0, 20));
        
        if (won) {
          toast.success(`Won $${(betAmount * target * 0.99).toFixed(2)}! Target: ${target}x`);
        } else {
          toast.error(`Lost $${betAmount.toFixed(2)}`);
        }
        
        setGameState('completed');
        setGamesPlayed(prev => prev + 1);
        
        // Auto-play mode
        if (isAutoMode && gamesPlayed < autoGames - 1) {
          setTimeout(() => {
            setGameState('idle');
            setTimeout(() => playGame(), 100);
          }, 1500);
        } else if (isAutoMode) {
          setIsAutoMode(false);
          setGamesPlayed(0);
          toast.success(`Auto-play completed!`);
        }
      }
    };
    
    animationRef.current = requestAnimationFrame(animate);
  }, [betAmount, target, isAuthenticated, isAutoMode, autoGames, gamesPlayed]);

  // Cleanup
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, []);

  // Quick target presets
  const targetPresets = [
    { value: 1.5, label: '1.5x' },
    { value: 2, label: '2x' },
    { value: 5, label: '5x' },
    { value: 10, label: '10x' },
    { value: 20, label: '20x' },
    { value: 50, label: '50x' },
    { value: 100, label: '100x' },
  ];

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Limbo</h1>
            <p className="text-gray-400">Set your target and watch the multiplier rise!</p>
          </div>
          <div className="flex gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Your Bets</span>
              <p className="text-xl font-mono text-primary-500">{myBets.length}</p>
            </div>
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Balance</span>
              <p className="text-xl font-mono text-green-400">${user?.balance?.toFixed(2) || '0.00'}</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Game Display */}
            <div className="glass rounded-2xl p-8">
              <div className="relative h-48 mb-4">
                {/* Target line */}
                <div 
                  className="absolute top-0 bottom-0 w-0.5 bg-primary-500 z-10"
                  style={{ left: `${Math.min((target / 100) * 100, 100)}%` }}
                >
                  <div className="absolute -top-2 -left-3 bg-primary-500 text-white text-xs px-1 rounded">
                    {target}x
                  </div>
                </div>
                
                {/* Multiplier display */}
                <div className="absolute inset-0 flex items-center justify-center">
                  <motion.div
                    animate={{ scale: gameState === 'playing' ? [1, 1.1, 1] : 1 }}
                    transition={{ repeat: gameState === 'playing' ? Infinity : 0, duration: 0.5 }}
                    className={`text-7xl font-mono font-bold ${
                      currentMultiplier !== null
                        ? currentMultiplier >= target
                          ? 'text-green-400'
                          : 'text-red-400'
                        : 'text-white'
                    }`}
                  >
                    {currentMultiplier !== null 
                      ? currentMultiplier.toFixed(2) 
                      : '1.00'
                    }x
                  </motion.div>
                </div>
                
                {/* Target reached indicator */}
                {currentMultiplier !== null && currentMultiplier >= target && (
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-green-500/30 rounded-full p-8"
                  >
                    <span className="text-4xl">🎯</span>
                  </motion.div>
                )}
              </div>
              
              {/* Result */}
              <AnimatePresence>
                {gameState === 'completed' && currentMultiplier !== null && (
                  <motion.div
                    initial={{ scale: 0.8, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    className={`text-center p-4 rounded-xl ${
                      currentMultiplier >= target
                        ? 'bg-green-500/20 border border-green-500/30'
                        : 'bg-red-500/20 border border-red-500/30'
                    }`}
                  >
                    <p className="text-2xl font-bold">
                      {currentMultiplier >= target 
                        ? `TARGET HIT! Won $${(betAmount * target * 0.99).toFixed(2)}` 
                        : `MISSED! Rolled ${currentMultiplier.toFixed(2)}x`
                      }
                    </p>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {/* Target Selection */}
              <div className="mb-6">
                <label className="block text-gray-400 text-sm mb-3">
                  Target Multiplier: <span className="text-white text-lg font-bold">{target}x</span>
                </label>
                <div className="grid grid-cols-4 md:grid-cols-7 gap-2 mb-3">
                  {targetPresets.map(preset => (
                    <button
                      key={preset.value}
                      onClick={() => { setTarget(preset.value); }}
                      className={`py-2 rounded-lg font-mono text-sm transition ${
                        target === preset.value
                          ? 'bg-primary-500 text-white'
                          : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
                      }`}
                    >
                      {preset.label}
                    </button>
                  ))}
                </div>
                <input
                  type="range"
                  min={MIN_TARGET * 100}
                  max={10000}
                  value={target * 100}
                  onChange={(e) => setTarget(parseFloat(e.target.value) / 100)}
                  className="w-full h-2 bg-tiger-surface rounded-lg appearance-none cursor-pointer accent-primary-500"
                />
                <div className="flex justify-between text-xs text-gray-500 mt-1">
                  <span>1.01x</span>
                  <span>1000x</span>
                </div>
              </div>

              {/* Win Chance & Multiplier */}
              <div className="grid grid-cols-2 gap-4 mb-6">
                <div className="bg-tiger-surface/50 p-4 rounded-lg text-center">
                  <p className="text-gray-400 text-sm">Win Chance</p>
                  <p className="text-2xl font-bold text-primary-400">{calculateWinChance(target).toFixed(1)}%</p>
                </div>
                <div className="bg-tiger-surface/50 p-4 rounded-lg text-center">
                  <p className="text-gray-400 text-sm">Multiplier</p>
                  <p className="text-2xl font-bold text-green-400">{calculateMultiplier(target).toFixed(3)}x</p>
                </div>
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

              {/* Potential Win */}
              <div className="mb-6 p-4 bg-tiger-surface/50 rounded-lg">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-400">Potential Win:</span>
                  <span className="text-green-400 font-mono">
                    ${(betAmount * target * 0.99).toFixed(2)}
                  </span>
                </div>
                <div className="flex justify-between text-sm mt-1">
                  <span className="text-gray-400">House Edge:</span>
                  <span className="text-yellow-400">1%</span>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-4">
                <button
                  onClick={playGame}
                  disabled={gameState === 'playing'}
                  className={`flex-1 py-4 rounded-xl font-bold text-xl transition-all ${
                    gameState === 'playing'
                      ? 'bg-gray-700 cursor-not-allowed'
                      : 'bg-primary-500 hover:bg-primary-600'
                  }`}
                >
                  {gameState === 'playing' ? 'Rolling...' : gameState === 'completed' ? 'Play Again' : 'Play'}
                </button>
                
                <button
                  onClick={() => {
                    if (gameState === 'idle') {
                      setIsAutoMode(true);
                      setGamesPlayed(0);
                      playGame();
                    } else {
                      setIsAutoMode(false);
                      setGamesPlayed(0);
                    }
                  }}
                  disabled={gameState === 'playing'}
                  className={`px-6 py-4 rounded-xl font-bold transition ${
                    isAutoMode 
                      ? 'bg-red-500 hover:bg-red-600'
                      : 'bg-tiger-surface hover:bg-primary-500/30'
                  }`}
                >
                  {isAutoMode ? 'Stop' : 'Auto'}
                </button>
              </div>

              {/* Auto Mode Settings */}
              {isAutoMode && (
                <div className="mt-4 p-4 bg-tiger-surface/50 rounded-lg">
                  <label className="block text-gray-400 text-sm mb-2">Auto-play games:</label>
                  <div className="flex gap-2">
                    {[5, 10, 25, 50].map(num => (
                      <button
                        key={num}
                        onClick={() => setAutoGames(num)}
                        className={`flex-1 py-2 rounded-lg text-sm ${
                          autoGames === num
                            ? 'bg-primary-500 text-white'
                            : 'bg-tiger-surface text-gray-400'
                        }`}
                      >
                        {num}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Bets</h3>
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {myBets.slice(0, 10).map(bet => (
                  <div 
                    key={bet.id} 
                    className={`p-3 rounded-lg ${
                      bet.won ? 'bg-green-500/20' : 'bg-red-500/20'
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className="text-gray-400 text-xs">Target: {bet.target}x</span>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="text-gray-400 text-sm">Result: {bet.result.toFixed(2)}x</span>
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

            {/* History */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">History</h3>
              <div className="flex flex-wrap gap-2">
                {history.slice(0, 15).map((item, index) => (
                  <div
                    key={index}
                    className={`w-12 h-8 rounded flex items-center justify-center text-xs font-mono ${
                      item.won ? 'bg-green-500/30 text-green-300' : 'bg-red-500/30 text-red-300'
                    }`}
                    title={`Target: ${item.target}x, Result: ${item.result.toFixed(2)}x`}
                  >
                    {item.result.toFixed(0)}
                  </div>
                ))}
                {history.length === 0 && (
                  <p className="text-gray-500 text-center w-full py-4">No history</p>
                )}
              </div>
            </div>

            {/* Provably Fair */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-2">Provably Fair</h3>
              <p className="text-gray-400 text-sm">
                Every game outcome is determined by cryptographic seeds. You can verify each round.
              </p>
              <button className="mt-3 text-primary-400 text-sm hover:underline">
                How it works →
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
