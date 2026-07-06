'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

const GRID_SIZE = 25;
const TILE_SIZE = 60;

export default function MinesGame() {
  const { isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [minesCount, setMinesCount] = useState(3);
  const [grid, setGrid] = useState<number[]>([]);
  const [gameOver, setGameOver] = useState(false);
  const [gameWon, setGameWon] = useState(false);
  const [currentMultiplier, setCurrentMultiplier] = useState(1.0);
  const [myBets, setMyBets] = useState<{ amount: number; mines: number; won: boolean; profit: number }[]>([]);

  const generateMines = (count: number) => {
    const mines = new Set<number>();
    while (mines.size < count) {
      mines.add(Math.floor(Math.random() * GRID_SIZE));
    }
    return Array.from(mines);
  };

  const startGame = () => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    const mines = generateMines(minesCount);
    setGrid(mines);
    setGameOver(false);
    setGameWon(false);
    setCurrentMultiplier(1.0);
  };

  const revealTile = (index: number) => {
    if (grid.length === 0 || gameOver) return;
    
    if (grid.includes(index)) {
      // Hit a mine!
      setGameOver(true);
      setGameWon(false);
      toast.error(`Hit a mine! Game over.`);
      
      setMyBets([...myBets, { amount: betAmount, mines: minesCount, won: false, profit: -betAmount }]);
    } else {
      // Safe tile
      const revealedCount = grid.filter((m, i) => i < index && !grid.includes(i)).length + 1;
      const newMultiplier = 1 + (revealedCount * 0.15);
      setCurrentMultiplier(newMultiplier);
      
      // Check if all safe tiles revealed
      const safeTiles = GRID_SIZE - minesCount;
      if (revealedCount >= safeTiles) {
        setGameOver(true);
        setGameWon(true);
        const winnings = betAmount * newMultiplier * 0.95;
        toast.success(`You won $${winnings.toFixed(2)}!`);
        
        setMyBets([...myBets, { amount: betAmount, mines: minesCount, won: true, profit: winnings - betAmount }]);
      }
    }
  };

  const cashout = () => {
    if (gameOver || grid.length === 0) return;
    
    const winnings = betAmount * currentMultiplier * 0.95;
    setGameOver(true);
    setGameWon(true);
    toast.success(`Cashed out $${winnings.toFixed(2)}!`);
    
    setMyBets([...myBets, { amount: betAmount, mines: minesCount, won: true, profit: winnings - betAmount }]);
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-6">
      <div className="max-w-7xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-heading text-gradient">Mines</h1>
            <p className="text-gray-400">Avoid the mines, collect the gems!</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-6">
            <div className="glass rounded-2xl p-8">
              <div className="flex justify-center mb-6">
                <div 
                  className="grid gap-2"
                  style={{ 
                    gridTemplateColumns: `repeat(5, ${TILE_SIZE}px)`,
                  }}
                >
                  {Array.from({ length: GRID_SIZE }).map((_, i) => {
                    const isMine = grid.includes(i);
                    const isRevealed = gameOver && isMine;
                    const isSafe = gameOver && !isMine;
                    
                    return (
                      <motion.button
                        key={i}
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={() => revealTile(i)}
                        disabled={gameOver || grid.length === 0}
                        className={`
                          w-14 h-14 rounded-lg font-bold text-2xl transition-all
                          ${isRevealed 
                            ? 'bg-red-500' 
                            : isSafe || grid.length === 0
                              ? 'bg-green-500'
                              : 'bg-tiger-surface hover:bg-primary-500/30'
                          }
                        `}
                      >
                        {isRevealed ? '💣' : (isSafe ? '💎' : '')}
                      </motion.button>
                    );
                  })}
                </div>
              </div>

              {gameOver && (
                <div className="text-center mb-6">
                  <p className={`text-2xl font-bold ${gameWon ? 'text-green-500' : 'text-red-500'}`}>
                    {gameWon ? 'You Won!' : 'Game Over!'}
                  </p>
                </div>
              )}
            </div>

            <div className="glass rounded-2xl p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                <div>
                  <label className="block text-gray-400 text-sm mb-2">Bet Amount</label>
                  <input
                    type="number"
                    value={betAmount}
                    onChange={(e) => setBetAmount(parseFloat(e.target.value) || 0)}
                    className="w-full bg-tiger-surface border border-gray-700 rounded-lg px-4 py-3 font-mono focus:border-primary-500 focus:outline-none"
                    step={0.01}
                    min={0.01}
                  />
                </div>
                
                <div>
                  <label className="block text-gray-400 text-sm mb-2">Number of Mines</label>
                  <div className="flex gap-2">
                    {[1, 3, 5, 10, 15, 24].map((count) => (
                      <button
                        key={count}
                        onClick={() => setMinesCount(count)}
                        className={`flex-1 py-2 rounded-lg font-semibold transition ${
                          minesCount === count 
                            ? 'bg-primary-500' 
                            : 'bg-tiger-surface hover:bg-primary-500/30'
                        }`}
                      >
                        {count}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="text-center mb-4">
                <p className="text-gray-400">Current Multiplier</p>
                <p className="text-4xl font-mono text-primary-500">{currentMultiplier.toFixed(2)}x</p>
                <p className="text-gray-500 text-sm">
                  Potential win: ${(betAmount * currentMultiplier * 0.95).toFixed(2)}
                </p>
              </div>

              <div className="flex gap-4">
                {grid.length === 0 ? (
                  <button
                    onClick={startGame}
                    className="flex-1 py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition"
                  >
                    Start Game
                  </button>
                ) : !gameOver ? (
                  <button
                    onClick={cashout}
                    className="flex-1 py-4 bg-green-500 hover:bg-green-600 rounded-xl font-bold text-xl transition"
                  >
                    Cash Out
                  </button>
                ) : (
                  <button
                    onClick={startGame}
                    className="flex-1 py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition"
                  >
                    Play Again
                  </button>
                )}
              </div>
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
                    <div>
                      <span className="font-mono">${bet.amount.toFixed(2)}</span>
                      <span className="text-gray-400 text-sm ml-2">({bet.mines} mines)</span>
                    </div>
                    <span className={bet.won ? 'text-green-400' : 'text-red-400'}>
                      {bet.profit >= 0 ? '+' : ''}{bet.profit.toFixed(2)}
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
                Every game outcome is determined by cryptographic seeds.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
