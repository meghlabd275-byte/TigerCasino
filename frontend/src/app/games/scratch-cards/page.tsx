'use client';

import React, { useState, useCallback, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Scratch card themes and their configurations
const SCRATCH_CARDS = [
  {
    id: 'golden_7s',
    name: 'Golden 7s',
    icon: '💰',
    theme: 'gold',
    symbols: ['🍒', '🍋', '🍇', '💎', '7️⃣', '⭐'],
    match3: 10,
    match4: 50,
    match5: 200,
    match6: 1000,
  },
  {
    id: 'diamond_mine',
    name: 'Diamond Mine',
    icon: '💎',
    theme: 'blue',
    symbols: ['💎', '🔷', '⛏️', '💰', '📦', '🧨'],
    match3: 15,
    match4: 75,
    match5: 300,
    match6: 2000,
  },
  {
    id: 'lucky_clover',
    name: 'Lucky Clover',
    icon: '🍀',
    theme: 'green',
    symbols: ['🍀', '🍀', '🍀', '🍀', '🌸', '🌟'],
    match3: 20,
    match4: 100,
    match5: 500,
    match6: 5000,
  },
  {
    id: 'fruity_gems',
    name: 'Fruity Gems',
    icon: '🍓',
    theme: 'red',
    symbols: ['🍒', '🍓', '🍋', '🍇', '🍊', '🥝'],
    match3: 10,
    match4: 50,
    match5: 250,
    match6: 1500,
  },
];

interface Bet {
  id: string;
  cardId: string;
  cardName: string;
  amount: number;
  matches: number;
  won: boolean;
  profit: number;
  timestamp: Date;
}

export default function ScratchCardsGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [selectedCard, setSelectedCard] = useState(SCRATCH_CARDS[0]);
  const [revealedSymbols, setRevealedSymbols] = useState<Set<number>>(new Set());
  const [allSymbols, setAllSymbols] = useState<string[]>([]);
  const [gameState, setGameState] = useState<'selecting' | 'scratching' | 'completed'>('selecting');
  const [myBets, setMyBets] = useState<Bet[]>([]);

  // Generate random symbols
  const generateSymbols = useCallback(() => {
    const symbols = [];
    for (let i = 0; i < 9; i++) {
      const randomSymbol = selectedCard.symbols[Math.floor(Math.random() * selectedCard.symbols.length)];
      symbols.push(randomSymbol);
    }
    return symbols;
  }, [selectedCard]);

  // Initialize card
  useEffect(() => {
    setAllSymbols(generateSymbols());
  }, [selectedCard, generateSymbols]);

  // Select card
  const selectCard = useCallback((card: typeof SCRATCH_CARDS[0]) => {
    setSelectedCard(card);
    setRevealedSymbols(new Set());
    setAllSymbols(generateSymbols());
    setGameState('selecting');
  }, [generateSymbols]);

  // Reveal symbol (click to scratch)
  const revealSymbol = useCallback((index: number) => {
    if (gameState !== 'scratching') return;
    if (revealedSymbols.has(index)) return;

    const newRevealed = new Set(revealedSymbols);
    newRevealed.add(index);
    setRevealedSymbols(newRevealed);

    // Check if all revealed
    if (newRevealed.size === 9) {
      // Calculate matches
      const symbolCounts: Record<string, number> = {};
      allSymbols.forEach((symbol, i) => {
        if (newRevealed.has(i)) {
          symbolCounts[symbol] = (symbolCounts[symbol] || 0) + 1;
        }
      });

      // Find max match
      const matches = Math.max(...Object.values(symbolCounts), 0);

      // Calculate winnings
      let multiplier = 0;
      if (matches === 6) multiplier = selectedCard.match6;
      else if (matches === 5) multiplier = selectedCard.match5;
      else if (matches === 4) multiplier = selectedCard.match4;
      else if (matches === 3) multiplier = selectedCard.match3;

      const won = multiplier > 0;
      const profit = won ? (betAmount * multiplier) - betAmount : -betAmount;

      // Record bet
      const newBet: Bet = {
        id: Date.now().toString(),
        cardId: selectedCard.id,
        cardName: selectedCard.name,
        amount: betAmount,
        matches,
        won,
        profit,
        timestamp: new Date()
      };

      setMyBets(prev => [newBet, ...prev].slice(0, 50));

      if (won) {
        toast.success(`Won $${(betAmount * multiplier).toFixed(2)} with ${matches} matches!`);
      } else {
        toast.error('No matches. Try again!');
      }

      setGameState('completed');
    }
  }, [gameState, revealedSymbols, allSymbols, selectedCard, betAmount]);

  // Start scratching
  const startScratching = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    setGameState('scratching');
  }, [isAuthenticated, betAmount]);

  // Reset game
  const resetGame = useCallback(() => {
    setRevealedSymbols(new Set());
    setAllSymbols(generateSymbols());
    setGameState('selecting');
  }, [generateSymbols]);

  // Auto-reveal all
  const autoReveal = useCallback(() => {
    if (gameState !== 'scratching') return;

    for (let i = 0; i < 9; i++) {
      if (!revealedSymbols.has(i)) {
        revealSymbol(i);
      }
    }
  }, [gameState, revealedSymbols, revealSymbol]);

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Scratch Cards</h1>
            <p className="text-gray-400">Scratch to reveal your fortune!</p>
          </div>
          <div className="glass px-4 py-2 rounded-lg">
            <span className="text-gray-400 text-sm">Balance</span>
            <p className="text-xl font-mono text-green-400">${user?.balance?.toFixed(2) || '0.00'}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Card Selection */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Select Card</h3>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {SCRATCH_CARDS.map(card => (
                  <button
                    key={card.id}
                    onClick={() => selectCard(card)}
                    className={`p-4 rounded-xl transition-all ${
                      selectedCard.id === card.id
                        ? 'bg-primary-500/30 border-2 border-primary-500'
                        : 'bg-tiger-surface hover:bg-primary-500/20'
                    }`}
                  >
                    <div className="text-3xl mb-2">{card.icon}</div>
                    <p className="font-bold text-sm">{card.name}</p>
                  </button>
                ))}
              </div>
            </div>

            {/* Scratch Card Display */}
            <div className="glass rounded-2xl p-6">
              <div className="flex justify-center mb-6">
                <div 
                  className="relative w-80 h-80"
                  style={{
                    background: selectedCard.theme === 'gold' 
                      ? 'linear-gradient(135deg, #FFD700 0%, #FFA500 100%)'
                      : selectedCard.theme === 'blue'
                        ? 'linear-gradient(135deg, #1E90FF 0%, #4169E1 100%)'
                        : selectedCard.theme === 'green'
                          ? 'linear-gradient(135deg, #32CD32 0%, #228B22 100%)'
                          : 'linear-gradient(135deg, #FF6347 0%, #DC143C 100%)'
                  }}
                  style={{ borderRadius: '1rem' }}
                >
                  {/* Hidden layer (silver scratch-off) */}
                  {gameState !== 'completed' && (
                    <div 
                      className="absolute inset-0 cursor-crosshair"
                      style={{ 
                        background: 'linear-gradient(135deg, #C0C0C0 0%, #A9A9A9 100%)',
                        borderRadius: '1rem'
                      }}
                    >
                      {/* Scratch overlay text */}
                      <div className="absolute inset-0 flex items-center justify-center">
                        <span className="text-4xl font-bold text-gray-500">💨</span>
                      </div>
                    </div>
                  )}
                  
                  {/* Revealed symbols */}
                  <div 
                    className="absolute inset-0 p-4"
                    style={{ 
                      display: 'grid', 
                      gridTemplateColumns: 'repeat(3, 1fr)',
                      gap: '0.5rem'
                    }}
                  >
                    {allSymbols.map((symbol, index) => (
                      <motion.button
                        key={index}
                        onClick={() => revealSymbol(index)}
                        disabled={gameState !== 'scratching'}
                        initial={{ scale: 0 }}
                        animate={{ 
                          scale: revealedSymbols.has(index) ? 1 : 0,
                          opacity: revealedSymbols.has(index) ? 1 : 0
                        }}
                        className={`
                          aspect-square rounded-lg flex items-center justify-center text-3xl
                          ${revealedSymbols.has(index) 
                            ? 'bg-white/90' 
                            : 'bg-transparent'
                          }
                          ${gameState !== 'scratching' ? 'cursor-default' : 'cursor-pointer'}
                        `}
                      >
                        {symbol}
                      </motion.button>
                    ))}
                  </div>

                  {/* Scratch to reveal overlay */}
                  {gameState === 'scratching' && (
                    <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                      <div className="bg-black/50 px-4 py-2 rounded-full">
                        <span className="text-white font-bold">Click to Scratch!</span>
                      </div>
                    </div>
                  )}
                </div>
              </div>

              {/* Instructions */}
              {gameState === 'selecting' && (
                <p className="text-center text-gray-400">
                  Select a card and place your bet to start scratching!
                </p>
              )}

              {/* Auto-scratch button */}
              {gameState === 'scratching' && revealedSymbols.size < 9 && revealedSymbols.size > 0 && (
                <button
                  onClick={autoReveal}
                  className="w-full py-3 bg-primary-500/20 hover:bg-primary-500/40 text-primary-400 rounded-xl font-bold transition"
                >
                  Reveal All
                </button>
              )}
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {/* Bet Amount */}
              <div className="mb-6">
                <label className="block text-gray-400 text-sm mb-2">Bet Amount</label>
                <div className="flex gap-2 mb-3">
                  {[0.10, 1, 10, 50, 100].map(amount => (
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

              {/* Payouts */}
              <div className="mb-6 p-4 bg-tiger-surface/50 rounded-lg">
                <p className="text-gray-400 text-sm mb-2">Payouts:</p>
                <div className="grid grid-cols-4 gap-2 text-center text-sm">
                  <div>
                    <span className="text-gray-500">3 Match</span>
                    <p className="text-green-400">{selectedCard.match3}x</p>
                  </div>
                  <div>
                    <span className="text-gray-500">4 Match</span>
                    <p className="text-green-400">{selectedCard.match4}x</p>
                  </div>
                  <div>
                    <span className="text-gray-500">5 Match</span>
                    <p className="text-green-400">{selectedCard.match5}x</p>
                  </div>
                  <div>
                    <span className="text-gray-500">6 Match</span>
                    <p className="text-yellow-400">{selectedCard.match6}x</p>
                  </div>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-4">
                {gameState === 'selecting' ? (
                  <button
                    onClick={startScratching}
                    className="flex-1 py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Start Scratching
                  </button>
                ) : gameState === 'completed' ? (
                  <button
                    onClick={resetGame}
                    className="flex-1 py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Play Again
                  </button>
                ) : (
                  <button
                    onClick={resetGame}
                    className="flex-1 py-4 bg-tiger-surface hover:bg-red-500/30 rounded-xl font-bold text-xl transition-all"
                  >
                    Reset
                  </button>
                )}
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Current Card Info */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">{selectedCard.name}</h3>
              <div className="flex justify-center mb-4">
                <div className="text-6xl">{selectedCard.icon}</div>
              </div>
              <p className="text-gray-400 text-sm text-center">
                Match 3 or more identical symbols to win!
              </p>
            </div>

            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Cards</h3>
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
                      <span className="text-gray-400 text-xs">{bet.cardName}</span>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="text-sm">{bet.matches} matches</span>
                      <span className={bet.won ? 'text-green-400' : 'text-red-400'}>
                        {bet.won ? `+$${bet.profit.toFixed(2)}` : `-$${bet.amount.toFixed(2)}`}
                      </span>
                    </div>
                  </div>
                ))}
                {myBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No cards scratched yet</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
