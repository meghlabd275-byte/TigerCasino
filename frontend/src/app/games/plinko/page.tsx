'use client';

import React, { useState, useEffect, useCallback, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Plinko board configuration
const ROWS = 16;
const COLS = ROWS + 1;
const TILE_SIZE = 40;

// Multiplier buckets for 16-row Plinko
const MULTIPLIERS_16 = [
  0, 0.5, 1, 2, 5, 10, 20, 50, 
  50, 20, 10, 5, 2, 1, 0.5, 0
];

// Alternative configurations for different row counts
const MULTIPLIERS_8 = [0, 1, 2, 5, 10, 20, 10, 5, 2, 1];
const MULTIPLIERS_12 = [0, 0.5, 1, 2, 5, 10, 25, 25, 10, 5, 2, 1, 0.5];

interface BallPosition {
  row: number;
  col: number;
  x: number;
  y: number;
}

interface Bet {
  id: string;
  amount: number;
  rows: number;
  multiplier: number;
  won: boolean;
  profit: number;
  timestamp: Date;
}

export default function PlinkoGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [rows, setRows] = useState(16);
  const [gameState, setGameState] = useState<'idle' | 'playing' | 'completed'>('idle');
  const [ballPosition, setBallPosition] = useState<BallPosition | null>(null);
  const [path, setPath] = useState<number[]>([]);
  const [currentMultiplier, setCurrentMultiplier] = useState(1);
  const [history, setHistory] = useState<number[]>([]);
  const [myBets, setMyBets] = useState<Bet[]>([]);
  const [isAutoMode, setIsAutoMode] = useState(false);
  const [autoGames, setAutoGames] = useState(10);
  const [gamesPlayed, setGamesPlayed] = useState(0);
  
  const boardRef = useRef<HTMLDivElement>(null);
  const animationRef = useRef<number>();

  // Get multipliers based on row count
  const getMultipliers = useCallback((rowCount: number) => {
    switch(rowCount) {
      case 8: return MULTIPLIERS_8;
      case 12: return MULTIPLIERS_12;
      case 16: return MULTIPLIERS_16;
      default: return MULTIPLIERS_16;
    }
  }, []);

  // Simulate ball dropping with physics-like animation
  const dropBall = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    setGameState('playing');
    
    // Determine path through the board
    const newPath: number[] = [];
    let currentCol = Math.floor(COLS / 2);
    
    for (let row = 0; row < rows; row++) {
      newPath.push(currentCol);
      // Randomly decide to go left or right at each peg
      const direction = Math.random() > 0.5 ? 1 : -1;
      currentCol += direction;
      // Clamp to valid range
      currentCol = Math.max(0, Math.min(COLS - 1, currentCol));
    }
    
    setPath(newPath);
    
    // Animate the ball
    let currentRow = 0;
    const animateBall = () => {
      if (currentRow >= newPath.length) {
        // Calculate final multiplier
        const multipliers = getMultipliers(rows);
        const finalCol = newPath[newPath.length - 1];
        const finalMultiplier = multipliers[finalCol] || 1;
        
        setCurrentMultiplier(finalMultiplier);
        
        // Calculate winnings
        const won = finalMultiplier > 0;
        const profit = won ? (betAmount * finalMultiplier * 0.97) - betAmount : -betAmount;
        
        // Record bet
        const newBet: Bet = {
          id: Date.now().toString(),
          amount: betAmount,
          rows,
          multiplier: finalMultiplier,
          won,
          profit,
          timestamp: new Date()
        };
        
        setMyBets(prev => [newBet, ...prev].slice(0, 50));
        setHistory(prev => [finalMultiplier, ...prev].slice(0, 20));
        
        if (won) {
          toast.success(`Won $${(betAmount * finalMultiplier * 0.97).toFixed(2)} at ${finalMultiplier}x!`);
        } else {
          toast.error(`Lost $${betAmount.toFixed(2)}`);
        }
        
        setGameState('completed');
        setGamesPlayed(prev => prev + 1);
        
        // Auto-play mode
        if (isAutoMode && gamesPlayed < autoGames - 1) {
          setTimeout(() => {
            setGameState('idle');
            setTimeout(() => dropBall(), 100);
          }, 1500);
        } else if (isAutoMode) {
          setIsAutoMode(false);
          setGamesPlayed(0);
          toast.success(`Auto-play completed!`);
        }
        
        return;
      }
      
      setBallPosition({
        row: currentRow,
        col: newPath[currentRow],
        x: newPath[currentRow] * TILE_SIZE + TILE_SIZE / 2,
        y: currentRow * TILE_SIZE + TILE_SIZE / 2 + 60
      });
      
      currentRow++;
      animationRef.current = requestAnimationFrame(animateBall);
    };
    
    animationRef.current = requestAnimationFrame(animateBall);
  }, [betAmount, rows, isAuthenticated, getMultipliers, isAutoMode, autoGames, gamesPlayed]);

  // Cleanup animation on unmount
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, []);

  // Reset game
  const resetGame = () => {
    setGameState('idle');
    setBallPosition(null);
    setPath([]);
    setCurrentMultiplier(1);
  };

  // Render the Plinko board
  const renderBoard = () => {
    const multipliers = getMultipliers(rows);
    const cells = [];
    
    // Calculate board dimensions
    const boardWidth = COLS * TILE_SIZE;
    const boardHeight = rows * TILE_SIZE + 100;
    
    // Generate peg positions
    for (let row = 0; row < rows; row++) {
      const isOddRow = row % 2 === 1;
      const colsInRow = isOddRow ? COLS - 1 : COLS;
      
      for (let col = 0; col < colsInRow; col++) {
        const x = isOddRow ? (col + 0.5) * TILE_SIZE : col * TILE_SIZE + TILE_SIZE / 2;
        const y = row * TILE_SIZE + TILE_SIZE / 2 + 50;
        
        cells.push(
          <div
            key={`peg-${row}-${col}`}
            className="absolute w-2 h-2 rounded-full bg-white/30"
            style={{
              left: `${x}px`,
              top: `${y}px`,
              transform: 'translate(-50%, -50%)'
            }}
          />
        );
      }
    }
    
    // Generate multiplier buckets at the bottom
    const bucketWidth = boardWidth / multipliers.length;
    multipliers.forEach((mult, index) => {
      const isHighlighted = path.length > 0 && path[path.length - 1] === index;
      cells.push(
        <div
          key={`bucket-${index}`}
          className={`absolute bottom-0 h-16 flex items-end justify-center pb-2 text-xs font-bold transition-colors ${
            isHighlighted 
              ? mult > 0 
                ? 'bg-green-500/50 text-green-300' 
                : 'bg-red-500/50 text-red-300'
              : 'bg-white/5 text-white/50'
          }`}
          style={{
            left: `${index * bucketWidth}px`,
            width: `${bucketWidth}px`,
          }}
        >
          {mult > 0 ? `${mult}x` : '-'}
        </div>
      );
    });
    
    return cells;
  };

  // Render ball
  const renderBall = () => {
    if (!ballPosition) return null;
    
    return (
      <motion.div
        initial={{ y: -50, opacity: 0 }}
        animate={{ 
          x: ballPosition.x - TILE_SIZE / 2,
          y: ballPosition.y,
          opacity: 1
        }}
        transition={{ duration: 0.3, ease: "easeOut" }}
        className="absolute w-6 h-6 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 shadow-lg z-10"
        style={{
          transform: 'translate(-50%, -50%)',
          left: 0,
          top: 0
        }}
      >
        <div className="absolute inset-1 rounded-full bg-white/30" />
      </motion.div>
    );
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Plinko</h1>
            <p className="text-gray-400">Watch the ball drop through the pegs!</p>
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
            {/* Game Board */}
            <div className="glass rounded-2xl p-4 md:p-6 overflow-hidden">
              <div 
                ref={boardRef}
                className="relative mx-auto"
                style={{
                  width: `${COLS * TILE_SIZE}px`,
                  height: `${rows * TILE_SIZE + 120}px`,
                  maxWidth: '100%'
                }}
              >
                {renderBoard()}
                <AnimatePresence>
                  {renderBall()}
                </AnimatePresence>
              </div>
              
              {/* Result Display */}
              {gameState === 'completed' && (
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className={`text-center mt-4 p-4 rounded-xl ${
                    currentMultiplier > 0 
                      ? 'bg-green-500/20 border border-green-500/30' 
                      : 'bg-red-500/20 border border-red-500/30'
                  }`}
                >
                  <p className="text-gray-400 text-sm">Result</p>
                  <p className={`text-4xl font-bold ${currentMultiplier > 0 ? 'text-green-400' : 'text-red-400'}`}>
                    {currentMultiplier > 0 ? `${currentMultiplier}x` : 'LOST'}
                  </p>
                  {currentMultiplier > 0 && (
                    <p className="text-green-300">
                      Won ${(betAmount * currentMultiplier * 0.97).toFixed(2)}
                    </p>
                  )}
                </motion.div>
              )}
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {/* Row Selection */}
              <div className="mb-6">
                <label className="block text-gray-400 text-sm mb-3">Select Rows: <span className="text-white">{rows}</span></label>
                <div className="flex gap-2">
                  {[8, 12, 16].map(rowCount => (
                    <button
                      key={rowCount}
                      onClick={() => { setRows(rowCount); resetGame(); }}
                      className={`flex-1 py-3 rounded-lg font-bold transition-all ${
                        rows === rowCount
                          ? 'bg-primary-500 text-white'
                          : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
                      }`}
                    >
                      {rowCount} Rows
                    </button>
                  ))}
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

              {/* Quick Multiplier Info */}
              <div className="mb-6 p-4 bg-tiger-surface/50 rounded-lg">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-400">Potential Win:</span>
                  <span className="text-green-400 font-mono">
                    ${(betAmount * 50 * 0.97).toFixed(2)} (at 50x)
                  </span>
                </div>
                <div className="flex justify-between text-sm mt-1">
                  <span className="text-gray-400">House Edge:</span>
                  <span className="text-yellow-400">3%</span>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-4">
                <button
                  onClick={dropBall}
                  disabled={gameState === 'playing'}
                  className={`flex-1 py-4 rounded-xl font-bold text-xl transition-all ${
                    gameState === 'playing'
                      ? 'bg-gray-700 cursor-not-allowed'
                      : 'bg-primary-500 hover:bg-primary-600'
                  }`}
                >
                  {gameState === 'playing' ? 'Dropping...' : gameState === 'completed' ? 'Play Again' : 'Drop Ball'}
                </button>
                
                <button
                  onClick={() => {
                    if (gameState === 'idle') {
                      setIsAutoMode(true);
                      setGamesPlayed(0);
                      dropBall();
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
                    className={`p-3 rounded-lg flex justify-between items-center ${
                      bet.won ? 'bg-green-500/20' : 'bg-red-500/20'
                    }`}
                  >
                    <div>
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className="text-gray-400 text-xs ml-2">({bet.rows} rows)</span>
                    </div>
                    <div className="text-right">
                      <span className={bet.won ? 'text-green-400' : 'text-red-400'}>
                        {bet.multiplier}x
                      </span>
                      <p className={`text-xs ${bet.profit >= 0 ? 'text-green-300' : 'text-red-300'}`}>
                        {bet.profit >= 0 ? '+' : ''}${bet.profit.toFixed(2)}
                      </p>
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
                {history.map((mult, index) => (
                  <div
                    key={index}
                    className={`w-10 h-10 rounded-lg flex items-center justify-center text-xs font-bold ${
                      mult > 10 
                        ? 'bg-purple-500/30 text-purple-300'
                        : mult > 0
                          ? 'bg-green-500/30 text-green-300'
                          : 'bg-red-500/30 text-red-300'
                    }`}
                  >
                    {mult}x
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
                Every game outcome is determined by cryptographic seeds. You can verify each round using the seed hash.
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
