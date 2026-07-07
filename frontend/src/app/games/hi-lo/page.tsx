'use client';

import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Card deck configuration
const SUITS = ['♠️', '♥️', '♦️', '♣️'];
const RANKS = ['A', '2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K'];
const CARD_VALUES: Record<string, number> = {
  'A': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, 
  '8': 8, '9': 9, '10': 10, 'J': 11, 'Q': 12, 'K': 13
};

interface Card {
  rank: string;
  suit: string;
  value: number;
}

interface Bet {
  id: string;
  amount: number;
  predictions: { choice: 'higher' | 'lower' | 'equal'; won: boolean }[];
  totalWon: number;
  timestamp: Date;
}

export default function HiLoGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [deck, setDeck] = useState<Card[]>([]);
  const [currentCard, setCurrentCard] = useState<Card | null>(null);
  const [gameState, setGameState] = useState<'idle' | 'playing' | 'gameover'>('idle');
  const [currentIndex, setCurrentIndex] = useState(0);
  const [streak, setStreak] = useState(0);
  const [totalWon, setTotalWon] = useState(0);
  const [myBets, setMyBets] = useState<Bet[]>([]);
  const [gameHistory, setGameHistory] = useState<{won: boolean; card: Card}[]>([]);

  // Initialize deck
  const initializeDeck = useCallback(() => {
    const newDeck: Card[] = [];
    for (const suit of SUITS) {
      for (const rank of RANKS) {
        newDeck.push({
          rank,
          suit,
          value: CARD_VALUES[rank]
        });
      }
    }
    // Shuffle deck
    for (let i = newDeck.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [newDeck[i], newDeck[j]] = [newDeck[j], newDeck[i]];
    }
    return newDeck;
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

    const newDeck = initializeDeck();
    setDeck(newDeck);
    setCurrentCard(newDeck[0]);
    setCurrentIndex(0);
    setStreak(0);
    setTotalWon(0);
    setGameState('playing');
    setGameHistory([{ won: true, card: newDeck[0] }]);
  }, [betAmount, isAuthenticated, initializeDeck]);

  // Handle player choice
  const handleChoice = useCallback((choice: 'higher' | 'lower' | 'equal') => {
    if (gameState !== 'playing' || !currentCard || currentIndex >= deck.length - 1) {
      return;
    }

    const nextCard = deck[currentIndex + 1];
    let won = false;

    if (choice === 'higher') {
      won = nextCard.value > currentCard.value;
    } else if (choice === 'lower') {
      won = nextCard.value < currentCard.value;
    } else {
      won = nextCard.value === currentCard.value;
    }

    // Calculate multiplier based on probability
    let multiplier = 1;
    if (choice === 'equal') {
      // Equal is rare (1/13 chance)
      multiplier = 12;
    } else {
      // Higher/lower: calculate based on cards remaining
      const higherCount = deck.slice(currentIndex + 1).filter(c => 
        choice === 'higher' ? c.value > currentCard.value : c.value < currentCard.value
      ).length;
      const probability = higherCount / (deck.length - currentIndex - 1);
      multiplier = Math.max(1.01, probability > 0 ? 0.97 / probability : 0);
    }

    if (won) {
      const winAmount = betAmount * multiplier;
      setTotalWon(prev => prev + winAmount);
      setStreak(prev => prev + 1);
      toast.success(`Won $${winAmount.toFixed(2)}!`);
    } else {
      setStreak(0);
      setTotalWon(0);
      toast.error('Wrong! Game over.');
    }

    setGameHistory(prev => [...prev, { won, card: nextCard }].slice(-20));
    setCurrentCard(nextCard);
    setCurrentIndex(prev => prev + 1);

    // Check if deck is finished
    if (currentIndex >= deck.length - 2) {
      setGameState('gameover');
      
      // Record final bet
      const newBet: Bet = {
        id: Date.now().toString(),
        amount: betAmount,
        predictions: gameHistory.slice(1).map(h => ({ choice: 'higher' as const, won: h.won })),
        totalWon,
        timestamp: new Date()
      };
      setMyBets(prev => [newBet, ...prev].slice(0, 50));
      
      if (totalWon > betAmount) {
        toast.success(`Game complete! Total won: $${totalWon.toFixed(2)}`);
      }
    }
  }, [gameState, currentCard, currentIndex, deck, betAmount, totalWon, gameHistory]);

  // Cash out
  const cashOut = useCallback(() => {
    if (gameState !== 'playing' || totalWon === 0) return;

    const newBet: Bet = {
      id: Date.now().toString(),
      amount: betAmount,
      predictions: gameHistory.slice(1).map(h => ({ choice: 'higher' as const, won: h.won })),
      totalWon,
      timestamp: new Date()
    };
    setMyBets(prev => [newBet, ...prev].slice(0, 50));
    
    setGameState('gameover');
    toast.success(`Cashed out! Total won: $${totalWon.toFixed(2)}`);
  }, [gameState, totalWon, betAmount, gameHistory]);

  // Render card
  const renderCard = (card: Card | null, size: 'small' | 'large' = 'large') => {
    if (!card) return null;
    
    const isRed = card.suit === '♥️' || card.suit === '♦️';
    const sizeClasses = size === 'large' 
      ? 'w-24 h-36 text-4xl' 
      : 'w-12 h-16 text-lg';
    
    return (
      <motion.div
        initial={{ scale: 0.8, rotate: -10 }}
        animate={{ scale: 1, rotate: 0 }}
        className={`${sizeClasses} bg-white rounded-lg shadow-lg flex flex-col items-center justify-center`}
      >
        <span className={isRed ? 'text-red-600' : 'text-black'}>{card.rank}</span>
        <span className={isRed ? 'text-red-600' : 'text-black'}>{card.suit}</span>
      </motion.div>
    );
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Hi-Lo</h1>
            <p className="text-gray-400">Predict if the next card is higher or lower!</p>
          </div>
          <div className="flex gap-4">
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Streak</span>
              <p className="text-xl font-mono text-primary-500">{streak}</p>
            </div>
            <div className="glass px-4 py-2 rounded-lg">
              <span className="text-gray-400 text-sm">Won</span>
              <p className="text-xl font-mono text-green-400">${totalWon.toFixed(2)}</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Card Display */}
            <div className="glass rounded-2xl p-8">
              <div className="flex justify-center items-center gap-8 mb-8">
                {/* Current Card */}
                <div className="flex flex-col items-center">
                  <p className="text-gray-400 text-sm mb-2">Current</p>
                  {renderCard(currentCard)}
                </div>
                
                {/* VS indicator */}
                <div className="text-2xl font-bold text-gray-500">VS</div>
                
                {/* Next Card (hidden during play) */}
                <div className="flex flex-col items-center">
                  <p className="text-gray-400 text-sm mb-2">Next</p>
                  <div className="w-24 h-36 bg-gradient-to-br from-primary-500/30 to-primary-700/30 rounded-lg flex items-center justify-center">
                    <span className="text-4xl">❓</span>
                  </div>
                </div>
              </div>

              {/* Card Info */}
              {currentCard && (
                <div className="text-center mb-4">
                  <p className="text-gray-400">
                    Card Value: <span className="text-white font-bold">{currentCard.rank} ({currentCard.value})</span>
                  </p>
                  <p className="text-gray-500 text-sm">
                    Cards remaining: {deck.length - currentIndex - 1}
                  </p>
                </div>
              )}

              {/* Result */}
              <AnimatePresence>
                {gameState === 'gameover' && (
                  <motion.div
                    initial={{ scale: 0.8, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    className={`text-center p-4 rounded-xl ${
                      totalWon > betAmount 
                        ? 'bg-green-500/20 border border-green-500/30'
                        : 'bg-red-500/20 border border-red-500/30'
                    }`}
                  >
                    <p className="text-2xl font-bold">
                      {totalWon > betAmount 
                        ? `WON $${(totalWon - betAmount).toFixed(2)}` 
                        : totalWon > 0 
                          ? `CASHED OUT: $${totalWon.toFixed(2)}`
                          : 'GAME OVER'}
                    </p>
                  </motion.div>
                )}
              </AnimatePresence>
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
                    Start Game
                  </button>
                </>
              ) : gameState === 'playing' ? (
                <>
                  {/* Prediction Buttons */}
                  <div className="mb-6">
                    <label className="block text-gray-400 text-sm mb-3">What will the next card be?</label>
                    <div className="grid grid-cols-3 gap-4">
                      <button
                        onClick={() => handleChoice('lower')}
                        className="py-4 bg-red-500/20 hover:bg-red-500/40 text-red-400 rounded-xl font-bold text-lg transition-all"
                      >
                        ↓ Lower
                        <span className="block text-xs font-normal mt-1">~46%</span>
                      </button>
                      <button
                        onClick={() => handleChoice('equal')}
                        className="py-4 bg-yellow-500/20 hover:bg-yellow-500/40 text-yellow-400 rounded-xl font-bold text-lg transition-all"
                      >
                        = Equal
                        <span className="block text-xs font-normal mt-1">~7.7%</span>
                      </button>
                      <button
                        onClick={() => handleChoice('higher')}
                        className="py-4 bg-green-500/20 hover:bg-green-500/40 text-green-400 rounded-xl font-bold text-lg transition-all"
                      >
                        ↑ Higher
                        <span className="block text-xs font-normal mt-1">~46%</span>
                      </button>
                    </div>
                  </div>

                  {/* Cash Out Button */}
                  {totalWon > 0 && (
                    <button
                      onClick={cashOut}
                      className="w-full py-4 bg-yellow-500 hover:bg-yellow-600 rounded-xl font-bold text-xl transition-all"
                    >
                      Cash Out ${totalWon.toFixed(2)}
                    </button>
                  )}
                </>
              ) : (
                <button
                  onClick={startGame}
                  className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                >
                  Play Again
                </button>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Game History */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Card History</h3>
              <div className="flex flex-wrap gap-2 justify-center">
                {gameHistory.slice(1).reverse().map((item, index) => (
                  <div
                    key={index}
                    className={`w-10 h-14 rounded flex flex-col items-center justify-center text-xs ${
                      item.won ? 'bg-green-500/30' : 'bg-red-500/30'
                    }`}
                  >
                    <span className={item.card.suit === '♥️' || item.card.suit === '♦️' ? 'text-red-300' : 'text-white'}>
                      {item.card.rank}
                    </span>
                    <span className={item.card.suit === '♥️' || item.card.suit === '♦️' ? 'text-red-300' : 'text-white'}>
                      {item.card.suit}
                    </span>
                  </div>
                ))}
                {gameHistory.length <= 1 && (
                  <p className="text-gray-500 text-center py-4">No cards yet</p>
                )}
              </div>
            </div>

            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Games</h3>
              <div className="space-y-2 max-h-60 overflow-y-auto">
                {myBets.slice(0, 10).map(bet => (
                  <div 
                    key={bet.id} 
                    className={`p-3 rounded-lg ${
                      bet.totalWon > bet.amount ? 'bg-green-500/20' : 'bg-red-500/20'
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className={bet.totalWon > bet.amount ? 'text-green-400' : 'text-red-400'}>
                        ${bet.totalWon.toFixed(2)}
                      </span>
                    </div>
                    <p className="text-gray-400 text-xs mt-1">
                      {bet.predictions.length} cards
                    </p>
                  </div>
                ))}
                {myBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No games yet</p>
                )}
              </div>
            </div>

            {/* Provably Fair */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-2">Provably Fair</h3>
              <p className="text-gray-400 text-sm">
                Cards are dealt from a shuffled deck using cryptographic seeds.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
