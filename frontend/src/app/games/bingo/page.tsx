'use client';

import React, { useState, useCallback, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Bingo configuration
const GRID_SIZE = 5; // 5x5 grid
const FREE_SPACE = 12; // Center is free
const TOTAL_CALLS = 75;

// Payouts based on pattern
const PAYOUTS = {
  '1_line': 5,     // Any line (horizontal, vertical, diagonal)
  '2_lines': 20,    // Any 2 lines
  '3_lines': 100,   // Any 3 lines
  '4_lines': 500,   // Any 4 lines  
  'bingo': 10000,   // Full card
};

interface BingoCard {
  numbers: (number | null)[][];
  marked: boolean[][];
}

interface Bet {
  id: string;
  amount: number;
  lines: number;
  won: boolean;
  pattern: string;
  profit: number;
  timestamp: Date;
}

export default function BingoGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [card, setCard] = useState<BingoCard>({ numbers: [], marked: [] });
  const [calledNumbers, setCalledNumbers] = useState<number[]>([]);
  const [currentCall, setCurrentCall] = useState<number | null>(null);
  const [gameState, setGameState] = useState<'idle' | 'playing' | 'completed'>('idle');
  const [myBets, setMyBets] = useState<Bet[]>([]);
  const [lines, setLines] = useState(0);

  // Generate random bingo card
  const generateCard = useCallback(() => {
    // Generate numbers for each column (B-I-N-G-O)
    const columnRanges = [
      [1, 15],   // B: 1-15
      [16, 30],  // I: 16-30
      [31, 45],  // N: 31-45
      [46, 60],  // G: 46-60
      [61, 75],  // O: 61-75
    ];
    
    const numbers: (number | null)[][] = [];
    const marked: boolean[][] = [];
    
    for (let col = 0; col < GRID_SIZE; col++) {
      const colNumbers: (number | null)[] = [];
      const colMarked: boolean[] = [];
      const range = columnRanges[col];
      const available = Array.from({ length: range[1] - range[0] + 1 }, (_, i) => range[0] + i);
      
      // Shuffle and take 5
      for (let i = available.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [available[i], available[j]] = [available[j], available[i]];
      }
      
      for (let row = 0; row < GRID_SIZE; row++) {
        // Center is free space
        if (row === 2 && col === 2) {
          colNumbers.push(null);
          colMarked.push(true);
        } else {
          colNumbers.push(available.pop() || null);
          colMarked.push(false);
        }
      }
      
      numbers.push(colNumbers);
      marked.push(colMarked);
    }
    
    // Transpose to get row-major order
    const transposedNumbers: (number | null)[][] = [];
    const transposedMarked: boolean[][] = [];
    for (let row = 0; row < GRID_SIZE; row++) {
      transposedNumbers.push([]);
      transposedMarked.push([]);
      for (let col = 0; col < GRID_SIZE; col++) {
        transposedNumbers[row].push(numbers[col][row]);
        transposedMarked[row].push(marked[col][row]);
      }
    }
    
    return { numbers: transposedNumbers, marked: transposedMarked };
  }, []);

  // Check for lines
  const checkLines = useCallback((cardState: BingoCard): number => {
    let lineCount = 0;
    const { marked, numbers } = cardState;
    
    // Check rows
    for (let row = 0; row < GRID_SIZE; row++) {
      if (marked[row].every(m => m)) lineCount++;
    }
    
    // Check columns
    for (let col = 0; col < GRID_SIZE; col++) {
      let colComplete = true;
      for (let row = 0; row < GRID_SIZE; row++) {
        if (!marked[row][col]) {
          colComplete = false;
          break;
        }
      }
      if (colComplete) lineCount++;
    }
    
    // Check diagonals
    let diag1Complete = true;
    let diag2Complete = true;
    for (let i = 0; i < GRID_SIZE; i++) {
      if (!marked[i][i]) diag1Complete = false;
      if (!marked[i][GRID_SIZE - 1 - i]) diag2Complete = false;
    }
    if (diag1Complete) lineCount++;
    if (diag2Complete) lineCount++;
    
    return lineCount;
  }, []);

  // Start new game
  const startGame = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    const newCard = generateCard();
    setCard(newCard);
    setCalledNumbers([]);
    setCurrentCall(null);
    setLines(0);
    setGameState('playing');
  }, [betAmount, isAuthenticated, generateCard]);

  // Mark number on card
  const markNumber = useCallback((num: number) => {
    setCard(prev => {
      const newMarked = prev.marked.map(row => [...row]);
      const newNumbers = prev.numbers;
      
      for (let row = 0; row < GRID_SIZE; row++) {
        for (let col = 0; col < GRID_SIZE; col++) {
          if (newNumbers[row][col] === num) {
            newMarked[row][col] = true;
          }
        }
      }
      
      return { ...prev, marked: newMarked };
    });
  }, []);

  // Draw next number
  const drawNextNumber = useCallback(() => {
    if (gameState !== 'playing') return;
    
    // Find available numbers
    const available = [];
    for (let i = 1; i <= TOTAL_CALLS; i++) {
      if (!calledNumbers.includes(i)) {
        available.push(i);
      }
    }
    
    if (available.length === 0) {
      toast.error('All numbers called!');
      setGameState('completed');
      return;
    }
    
    // Pick random number
    const randomIndex = Math.floor(Math.random() * available.length);
    const newNumber = available[randomIndex];
    
    setCurrentCall(newNumber);
    setCalledNumbers(prev => [...prev, newNumber]);
    
    // Mark on card
    markNumber(newNumber);
    
    // Check lines
    const newLines = checkLines(card);
    if (newLines > lines) {
      setLines(newLines);
      
      // Determine payout
      let won = false;
      let pattern = '';
      let multiplier = 0;
      
      if (newLines >= 5) {
        won = true;
        pattern = 'BINGO!';
        multiplier = PAYOUTS['bingo'];
      } else if (newLines >= 4) {
        won = true;
        pattern = '4 Lines';
        multiplier = PAYOUTS['4_lines'];
      } else if (newLines >= 3) {
        won = true;
        pattern = '3 Lines';
        multiplier = PAYOUTS['3_lines'];
      } else if (newLines >= 2) {
        won = true;
        pattern = '2 Lines';
        multiplier = PAYOUTS['2_lines'];
      } else if (newLines >= 1) {
        won = true;
        pattern = '1 Line';
        multiplier = PAYOUTS['1_line'];
      }
      
      if (won) {
        const profit = (betAmount * multiplier) - betAmount;
        
        const newBet: Bet = {
          id: Date.now().toString(),
          amount: betAmount,
          lines: newLines,
          won: true,
          pattern,
          profit,
          timestamp: new Date()
        };
        
        setMyBets(prev => [newBet, ...prev].slice(0, 50));
        toast.success(`${pattern}! Won $${profit.toFixed(2)}!`);
      }
    }
    
    // Check if card is complete
    if (newLines >= 5) {
      setGameState('completed');
    }
  }, [gameState, calledNumbers, card, lines, betAmount, checkLines, markNumber]);

  // Initialize card on mount
  useEffect(() => {
    setCard(generateCard());
  }, [generateCard]);

  // Get letter for number
  const getLetter = (num: number) => {
    if (num <= 15) return 'B';
    if (num <= 30) return 'I';
    if (num <= 45) return 'N';
    if (num <= 60) return 'G';
    return 'O';
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Bingo</h1>
            <p className="text-gray-400">Get 5 in a row to win!</p>
          </div>
          <div className="flex gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Lines</span>
              <p className="text-xl font-mono text-primary-500">{lines}</p>
            </div>
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Called</span>
              <p className="text-xl font-mono text-green-400">{calledNumbers.length}/{TOTAL_CALLS}</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Bingo Card */}
            <div className="glass rounded-2xl p-4">
              {/* B-I-N-G-O Header */}
              <div className="grid grid-cols-5 gap-2 mb-4">
                {['B', 'I', 'N', 'G', 'O'].map((letter, i) => (
                  <div key={letter} className="text-center py-2 bg-primary-500 rounded-lg font-bold text-xl">
                    {letter}
                  </div>
                ))}
              </div>
              
              {/* Card Grid */}
              <div className="grid grid-cols-5 gap-2">
                {card.numbers.map((row, rowIndex) => (
                  row.map((num, colIndex) => {
                    const isMarked = card.marked[rowIndex][colIndex];
                    const isCurrentCall = currentCall === num;
                    
                    return (
                      <motion.div
                        key={`${rowIndex}-${colIndex}`}
                        animate={isCurrentCall ? { scale: [1, 1.2, 1] } : {}}
                        className={`
                          aspect-square rounded-lg flex items-center justify-center text-lg font-bold transition-all
                          ${num === null 
                            ? 'bg-yellow-500/30 text-yellow-300' 
                            : isMarked 
                              ? 'bg-primary-500 text-white' 
                              : isCurrentCall 
                                ? 'bg-green-500 text-white'
                                : 'bg-tiger-surface text-gray-400'
                          }
                        `}
                      >
                        {num || '★'}
                      </motion.div>
                    );
                  })
                ))}
              </div>
              
              {/* Current Call */}
              {currentCall && (
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className="mt-4 text-center"
                >
                  <p className="text-gray-400 text-sm">Latest Call:</p>
                  <div className="inline-flex items-center gap-2 bg-primary-500/20 px-6 py-3 rounded-lg">
                    <span className="text-2xl font-bold text-primary-400">{getLetter(currentCall)}</span>
                    <span className="text-4xl font-bold text-white">{currentCall}</span>
                  </div>
                </motion.div>
              )}
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {gameState === 'idle' ? (
                <>
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

                  <button
                    onClick={startGame}
                    className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                  >
                    New Card (Generate Random)
                  </button>
                </>
              ) : gameState === 'playing' ? (
                <>
                  <button
                    onClick={drawNextNumber}
                    className="w-full py-4 bg-green-500 hover:bg-green-600 rounded-xl font-bold text-xl transition-all"
                  >
                    🎲 Call Next Number
                  </button>
                  <p className="text-center text-gray-400 text-sm mt-2">
                    Call all numbers to complete the card!
                  </p>
                </>
              ) : (
                <>
                  <div className="text-center mb-4">
                    <p className="text-2xl font-bold text-green-400">BINGO!</p>
                    <p className="text-gray-400">You completed the card!</p>
                  </div>
                  <button
                    onClick={startGame}
                    className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Play Again
                  </button>
                </>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Called Numbers */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Called Numbers</h3>
              <div className="grid grid-cols-5 gap-2">
                {calledNumbers.map(num => (
                  <div
                    key={num}
                    className="w-8 h-8 rounded flex items-center justify-center text-xs font-bold bg-primary-500/30 text-primary-300"
                  >
                    {num}
                  </div>
                ))}
                {calledNumbers.length === 0 && (
                  <p className="text-gray-500 text-center w-full col-span-5 py-4">No numbers called yet</p>
                )}
              </div>
            </div>

            {/* Payouts */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Payouts</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-400">1 Line:</span>
                  <span className="text-green-400">{PAYOUTS['1_line']}x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">2 Lines:</span>
                  <span className="text-green-400">{PAYOUTS['2_lines']}x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">3 Lines:</span>
                  <span className="text-green-400">{PAYOUTS['3_lines']}x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">4 Lines:</span>
                  <span className="text-green-400">{PAYOUTS['4_lines']}x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">BINGO:</span>
                  <span className="text-purple-400">{PAYOUTS['bingo']}x</span>
                </div>
              </div>
            </div>

            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Wins</h3>
              <div className="space-y-2 max-h-60 overflow-y-auto">
                {myBets.slice(0, 10).map(bet => (
                  <div 
                    key={bet.id} 
                    className="p-3 rounded-lg bg-green-500/20"
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className="text-green-400 font-bold">{bet.pattern}</span>
                    </div>
                    <p className="text-green-300 text-sm">+${bet.profit.toFixed(2)}</p>
                  </div>
                ))}
                {myBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No wins yet</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
