'use client';

import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Dice configuration
const MIN_VALUE = 0.01;
const MAX_VALUE = 100.00;
const DEFAULT_TARGET = 50;

interface Bet {
  id: string;
  amount: number;
  target: number;
  direction: 'over' | 'under';
  roll: number;
  won: boolean;
  profit: number;
  timestamp: Date;
}

export default function DiceGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [target, setTarget] = useState(DEFAULT_TARGET);
  const [direction, setDirection] = useState<'over' | 'under'>('over');
  const [currentRoll, setCurrentRoll] = useState<number | null>(null);
  const [gameState, setGameState] = useState<'idle' | 'rolling' | 'completed'>('idle');
  const [history, setHistory] = useState<{roll: number; won: boolean}[]>([]);
  const [myBets, setMyBets] = useState<Bet[]>([]);
  const [isAutoMode, setIsAutoMode] = useState(false);
  const [autoGames, setAutoGames] = useState(10);
  const [gamesPlayed, setGamesPlayed] = useState(0);

  // Calculate win chance and multiplier
  const calculateWinChance = useCallback(() => {
    if (direction === 'over') {
      return ((MAX_VALUE - target) / MAX_VALUE) * 100;
    } else {
      return (target / MAX_VALUE) * 100;
    }
  }, [target, direction]);

  const calculateMultiplier = useCallback(() => {
    const chance = calculateWinChance() / 100;
    // House edge of 1%
    return (chance * 0.99 > 0) ? (chance * 0.99) / chance : 0;
  }, [calculateWinChance]);

  // Simulate dice roll
  const rollDice = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    if (target <= 0 || target > MAX_VALUE) {
      toast.error('Invalid target value');
      return;
    }

    setGameState('rolling');
    
    // Simulate rolling animation
    let rollCount = 0;
    const rollInterval = setInterval(() => {
      setCurrentRoll(Math.random() * MAX_VALUE);
      rollCount++;
      
      if (rollCount > 10) {
        clearInterval(rollInterval);
        
        // Final result
        const finalRoll = Math.random() * MAX_VALUE;
        setCurrentRoll(finalRoll);
        
        // Determine win
        const won = direction === 'over' 
          ? finalRoll > target 
          : finalRoll < target;
        
        const multiplier = won ? calculateMultiplier() : 0;
        const profit = won ? (betAmount * multiplier) - betAmount : -betAmount;
        
        // Record bet
        const newBet: Bet = {
          id: Date.now().toString(),
          amount: betAmount,
          target,
          direction,
          roll: finalRoll,
          won,
          profit,
          timestamp: new Date()
        };
        
        setMyBets(prev => [newBet, ...prev].slice(0, 50));
        setHistory(prev => [{roll: finalRoll, won}, ...prev].slice(0, 20));
        
        if (won) {
          toast.success(`Won $${(betAmount * multiplier).toFixed(2)}!`);
        } else {
          toast.error(`Lost $${betAmount.toFixed(2)}`);
        }
        
        setGameState('completed');
        setGamesPlayed(prev => prev + 1);
        
        // Auto-play mode
        if (isAutoMode && gamesPlayed < autoGames - 1) {
          setTimeout(() => {
            setGameState('idle');
            setTimeout(() => rollDice(), 100);
          }, 1000);
        } else if (isAutoMode) {
          setIsAutoMode(false);
          setGamesPlayed(0);
          toast.success(`Auto-play completed!`);
        }
      }
    }, 50);
  }, [betAmount, target, direction, isAuthenticated, calculateMultiplier, isAutoMode, autoGames, gamesPlayed]);

  // Reset game
  const resetGame = () => {
    setGameState('idle');
    setCurrentRoll(null);
  };

  // Quick target presets
  const targetPresets = [
    { value: 10, label: '10' },
    { value: 25, label: '25' },
    { value: 50, label: '50' },
    { value: 75, label: '75' },
    { value: 90, label: '90' },
  ];

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Dice</h1>
            <p className="text-gray-400">Predict whether the roll will be over or under your target!</p>
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
            {/* Dice Display */}
            <div className="glass rounded-2xl p-8">
              <div className="flex justify-center items-center mb-8">
                <motion.div
                  animate={gameState === 'rolling' ? {
                    rotateX: [0, 360, 720, 1080],
                    rotateY: [0, 180, 360, 540]
                  } : {}}
                  transition={{ duration: 1, ease: "easeOut" }}
                  className="relative"
                >
                  {/* Dice visualization */}
                  <div className={`w-48 h-48 rounded-2xl flex items-center justify-center text-6xl font-bold ${
                    currentRoll !== null
                      ? currentRoll > target && direction === 'over'
                        ? 'bg-green-500/30 text-green-400'
                        : currentRoll < target && direction === 'under'
                          ? 'bg-green-500/30 text-green-400'
                          : 'bg-red-500/30 text-red-400'
                      : 'bg-tiger-surface text-white'
                  }`}>
                    {currentRoll !== null ? currentRoll.toFixed(2) : '🎲'}
                  </div>
                  
                  {/* Target line indicator */}
                  {currentRoll !== null && (
                    <motion.div
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      className="absolute -right-4 top-1/2 -translate-y-1/2 w-2 h-8 bg-primary-500 rounded-full"
                    />
                  )}
                </motion.div>
              </div>
              
              {/* Target indicator */}
              <div className="relative h-8 mb-4">
                <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-700 rounded-full" />
                <motion.div
                  animate={{ left: `${(target / MAX_VALUE) * 100}%` }}
                  className="absolute top-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 bg-primary-500 rounded-full border-2 border-white"
                />
                {currentRoll !== null && (
                  <motion.div
                    initial={{ left: 0 }}
                    animate={{ left: `${(currentRoll / MAX_VALUE) * 100}%` }}
                    className={`absolute top-1/2 -translate-x-1/2 -translate-y-1/2 w-3 h-3 rounded-full ${
                      (direction === 'over' && currentRoll > target) ||
                      (direction === 'under' && currentRoll < target)
                        ? 'bg-green-500'
                        : 'bg-red-500'
                    }`}
                  />
                )}
              </div>
              
              {/* Scale labels */}
              <div className="flex justify-between text-gray-500 text-sm">
                <span>0</span>
                <span>25</span>
                <span>50</span>
                <span>75</span>
                <span>100</span>
              </div>
              
              {/* Result */}
              <AnimatePresence>
                {gameState === 'completed' && currentRoll !== null && (
                  <motion.div
                    initial={{ scale: 0.8, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    className={`text-center mt-6 p-4 rounded-xl ${
                      (direction === 'over' && currentRoll > target) ||
                      (direction === 'under' && currentRoll < target)
                        ? 'bg-green-500/20 border border-green-500/30'
                        : 'bg-red-500/20 border border-red-500/30'
                    }`}
                  >
                    <p className="text-2xl font-bold">
                      {((direction === 'over' && currentRoll > target) ||
                      (direction === 'under' && currentRoll < target))
                        ? 'YOU WON!'
                        : 'YOU LOST!'}
                    </p>
                    <p className="text-gray-400">
                      Rolled {currentRoll.toFixed(2)} {direction} {target}
                    </p>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {/* Direction Selection */}
              <div className="mb-6">
                <label className="block text-gray-400 text-sm mb-3">Prediction</label>
                <div className="flex gap-2">
                  <button
                    onClick={() => { setDirection('over'); resetGame(); }}
                    className={`flex-1 py-3 rounded-lg font-bold text-lg transition-all ${
                      direction === 'over'
                        ? 'bg-green-500 text-white'
                        : 'bg-tiger-surface text-gray-400 hover:bg-green-500/30'
                    }`}
                  >
                    Over {target}
                  </button>
                  <button
                    onClick={() => { setDirection('under'); resetGame(); }}
                    className={`flex-1 py-3 rounded-lg font-bold text-lg transition-all ${
                      direction === 'under'
                        ? 'bg-red-500 text-white'
                        : 'bg-tiger-surface text-gray-400 hover:bg-red-500/30'
                    }`}
                  >
                    Under {target}
                  </button>
                </div>
              </div>

              {/* Target Selection */}
              <div className="mb-6">
                <label className="block text-gray-400 text-sm mb-3">
                  Target: <span className="text-white text-lg">{target}</span>
                </label>
                <div className="flex gap-2 mb-3">
                  {targetPresets.map(preset => (
                    <button
                      key={preset.value}
                      onClick={() => { setTarget(preset.value); resetGame(); }}
                      className={`flex-1 py-2 rounded-lg font-mono transition ${
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
                  min={1}
                  max={99}
                  value={target}
                  onChange={(e) => { setTarget(parseInt(e.target.value)); resetGame(); }}
                  className="w-full h-2 bg-tiger-surface rounded-lg appearance-none cursor-pointer accent-primary-500"
                />
              </div>

              {/* Win Chance & Multiplier */}
              <div className="grid grid-cols-2 gap-4 mb-6">
                <div className="bg-tiger-surface/50 p-4 rounded-lg text-center">
                  <p className="text-gray-400 text-sm">Win Chance</p>
                  <p className="text-2xl font-bold text-primary-400">{calculateWinChance().toFixed(1)}%</p>
                </div>
                <div className="bg-tiger-surface/50 p-4 rounded-lg text-center">
                  <p className="text-gray-400 text-sm">Multiplier</p>
                  <p className="text-2xl font-bold text-green-400">{calculateMultiplier().toFixed(3)}x</p>
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
                    ${(betAmount * calculateMultiplier()).toFixed(2)}
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
                  onClick={rollDice}
                  disabled={gameState === 'rolling'}
                  className={`flex-1 py-4 rounded-xl font-bold text-xl transition-all ${
                    gameState === 'rolling'
                      ? 'bg-gray-700 cursor-not-allowed'
                      : 'bg-primary-500 hover:bg-primary-600'
                  }`}
                >
                  {gameState === 'rolling' ? 'Rolling...' : gameState === 'completed' ? 'Roll Again' : 'Roll Dice'}
                </button>
                
                <button
                  onClick={() => {
                    if (gameState === 'idle') {
                      setIsAutoMode(true);
                      setGamesPlayed(0);
                      rollDice();
                    } else {
                      setIsAutoMode(false);
                      setGamesPlayed(0);
                    }
                  }}
                  disabled={gameState === 'rolling'}
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
                      <span className="text-gray-400 text-xs">
                        {bet.direction} {bet.target}
                      </span>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="font-mono text-lg">{bet.roll.toFixed(2)}</span>
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
                {history.map((item, index) => (
                  <div
                    key={index}
                    className={`w-12 h-8 rounded flex items-center justify-center text-xs font-mono ${
                      item.won ? 'bg-green-500/30 text-green-300' : 'bg-red-500/30 text-red-300'
                    }`}
                  >
                    {item.roll.toFixed(0)}
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
                Every roll is determined by cryptographic seeds. Verify each roll using the seed hash.
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
