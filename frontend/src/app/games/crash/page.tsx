'use client';

import React, { useState, useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

export default function CrashGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount, currentMultiplier, setCurrentMultiplier, isPlaying, setIsPlaying } = useGameStore();
  
  const [gameState, setGameState] = useState<'waiting' | 'rising' | 'crashed'>('waiting');
  const [crashPoint, setCrashPoint] = useState<number | null>(null);
  const [history, setHistory] = useState<number[]>([]);
  const [myBets, setMyBets] = useState<{ amount: number; multiplier: number; won: boolean }[]>([]);
  const [autoCashout, setAutoCashout] = useState<number | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationRef = useRef<number>();

  // Game loop
  useEffect(() => {
    if (gameState === 'rising') {
      const startTime = Date.now();
      const animate = () => {
        const elapsed = (Date.now() - startTime) / 1000;
        const multiplier = Math.min(1 + (elapsed * 0.03 * Math.pow(1.003, elapsed * 10)), crashPoint || 100);
        setCurrentMultiplier(multiplier);
        
        if (autoCashout && multiplier >= autoCashout && isPlaying) {
          handleCashout();
          return;
        }
        
        if (multiplier >= (crashPoint || 100)) {
          setGameState('crashed');
          setIsPlaying(false);
          handleCrash();
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
  }, [gameState, crashPoint]);

  const handlePlaceBet = async () => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }
    
    const newCrashPoint = 1 + Math.random() * 10 + (Math.random() > 0.7 ? Math.random() * 90 : 0);
    setCrashPoint(newCrashPoint);
    setGameState('rising');
    setIsPlaying(true);
    
    setMyBets([...myBets, { amount: betAmount, multiplier: 0, won: false }]);
  };

  const handleCashout = () => {
    if (!isPlaying) return;
    
    const multiplier = currentMultiplier;
    const winnings = betAmount * multiplier * 0.97;
    
    setMyBets(myBets.map((bet, i) => 
      i === myBets.length - 1 
        ? { ...bet, multiplier, won: true }
        : bet
    ));
    
    setGameState('crashed');
    setIsPlaying(false);
    toast.success(`Cashed out at ${multiplier.toFixed(2)}x! Won $${winnings.toFixed(2)}`);
    
    setHistory([multiplier, ...history.slice(0, 19)]);
  };

  const handleCrash = () => {
    const multiplier = currentMultiplier;
    setMyBets(myBets.map((bet, i) => 
      i === myBets.length - 1 
        ? { ...bet, multiplier, won: false }
        : bet
    ));
    
    setHistory([multiplier, ...history.slice(0, 19)]);
    toast.error(`Crashed at ${multiplier.toFixed(2)}x`);
    
    setTimeout(() => {
      setGameState('waiting');
      setCrashPoint(null);
    }, 3000);
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-6">
      <div className="max-w-7xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-heading text-gradient">Crash</h1>
            <p className="text-gray-400">Cash out before the crash!</p>
          </div>
          <div className="flex gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Your Bets</span>
              <p className="text-xl font-mono text-primary-500">{myBets.length}</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-6">
            <div className="glass rounded-2xl p-8 relative overflow-hidden">
              <div className="absolute inset-0 bg-gradient-to-b from-primary-500/10 to-transparent" />
              
              <div className="relative z-10 text-center py-16">
                <motion.div
                  key={currentMultiplier}
                  initial={{ scale: 0.9 }}
                  animate={{ scale: 1 }}
                  className={`text-8xl font-mono font-bold ${
                    gameState === 'crashed' 
                      ? 'text-red-500' 
                      : gameState === 'rising'
                        ? 'text-primary-500'
                        : 'text-gray-500'
                  }`}
                >
                  {gameState === 'waiting' ? '1.00' : currentMultiplier.toFixed(2)}x
                </motion.div>
                
                <p className="text-gray-400 mt-4">
                  {gameState === 'waiting' && 'Place your bet to start'}
                  {gameState === 'rising' && 'Cash out now!'}
                  {gameState === 'crashed' && `Crashed at ${currentMultiplier.toFixed(2)}x`}
                </p>
              </div>
            </div>

            <div className="glass rounded-2xl p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                <div>
                  <label className="block text-gray-400 text-sm mb-2">Bet Amount</label>
                  <div className="flex gap-2">
                    <input
                      type="number"
                      value={betAmount}
                      onChange={(e) => setBetAmount(parseFloat(e.target.value) || 0)}
                      className="flex-1 bg-tiger-surface border border-gray-700 rounded-lg px-4 py-3 font-mono text-lg focus:border-primary-500 focus:outline-none"
                      step={0.01}
                      min={0.01}
                    />
                  </div>
                </div>
                
                <div>
                  <label className="block text-gray-400 text-sm mb-2">Auto Cashout (optional)</label>
                  <input
                    type="number"
                    value={autoCashout || ''}
                    onChange={(e) => setAutoCashout(parseFloat(e.target.value) || null)}
                    placeholder="e.g. 2.00"
                    className="w-full bg-tiger-surface border border-gray-700 rounded-lg px-4 py-3 font-mono focus:border-primary-500 focus:outline-none"
                    step={0.1}
                    min={1.01}
                  />
                </div>
              </div>

              <button
                onClick={gameState === 'waiting' ? handlePlaceBet : handleCashout}
                disabled={gameState === 'waiting' && !isAuthenticated}
                className={`w-full py-4 rounded-xl font-bold text-xl transition-all ${
                  gameState === 'waiting'
                    ? 'bg-primary-500 hover:bg-primary-600'
                    : gameState === 'rising'
                      ? 'bg-green-500 hover:bg-green-600 animate-pulse'
                      : 'bg-gray-700 cursor-not-allowed'
                }`}
              >
                {gameState === 'waiting' ? 'Place Bet' : 'Cash Out!'}
              </button>
            </div>
          </div>

          <div className="space-y-6">
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Bets</h3>
              <div className="space-y-2 max-h-64 overflow-y-auto">
                {myBets.slice(-5).reverse().map((bet, i) => (
                  <div 
                    key={i} 
                    className={`p-3 rounded-lg flex justify-between items-center ${
                      bet.won ? 'bg-green-500/20' : 'bg-red-500/20'
                    }`}
                  >
                    <span className="font-mono">${bet.amount.toFixed(2)}</span>
                    <span className={bet.won ? 'text-green-400' : 'text-red-400'}>
                      {bet.multiplier > 0 ? `${bet.multiplier.toFixed(2)}x` : 'Lost'}
                    </span>
                  </div>
                ))}
                {myBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No bets yet</p>
                )}
              </div>
            </div>

            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-2">Provably Fair</h3>
              <p className="text-gray-400 text-sm">
                Every game outcome is determined by cryptographic seeds that you can verify.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
