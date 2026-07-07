'use client';

import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Card configuration
const SUITS = ['♠', '♥', '♦', '♣'];
const RANKS = ['A', '2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K'];

interface Card {
  rank: string;
  suit: string;
  value: number;
}

interface Bet {
  id: string;
  amount: number;
  playerScore: number;
  dealerScore: number;
  result: 'win' | 'lose' | 'push' | 'blackjack';
  profit: number;
  timestamp: Date;
}

// Calculate hand value (handle Ace as 1 or 11)
const calculateScore = (cards: Card[]): number => {
  let score = 0;
  let aces = 0;
  
  cards.forEach(card => {
    if (card.rank === 'A') {
      aces++;
      score += 11;
    } else if (['K', 'Q', 'J'].includes(card.rank)) {
      score += 10;
    } else {
      score += card.value;
    }
  });
  
  // Adjust for aces
  while (score > 21 && aces > 0) {
    score -= 10;
    aces--;
  }
  
  return score;
};

// Check for blackjack (Ace + 10-value card)
const isBlackjack = (cards: Card[]): boolean => {
  return cards.length === 2 && calculateScore(cards) === 21;
};

export default function BlackjackGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [deck, setDeck] = useState<Card[]>([]);
  const [playerHand, setPlayerHand] = useState<Card[]>([]);
  const [dealerHand, setDealerHand] = useState<Card[]>([]);
  const [gameState, setGameState] = useState<'betting' | 'playing' | 'dealerTurn' | 'completed'>('betting');
  const [myBets, setMyBets] = useState<Bet[]>([]);

  // Create and shuffle deck
  const createDeck = useCallback((): Card[] => {
    const newDeck: Card[] = [];
    for (const suit of SUITS) {
      for (const rank of RANKS) {
        let value = parseInt(rank);
        if (isNaN(value)) {
          if (rank === 'A') value = 11;
          else value = 10;
        }
        newDeck.push({ rank, suit, value });
      }
    }
    
    // Shuffle
    for (let i = newDeck.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [newDeck[i], newDeck[j]] = [newDeck[j], newDeck[i]];
    }
    
    return newDeck;
  }, []);

  // Deal initial hands
  const dealGame = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    const newDeck = createDeck();
    
    // Deal cards
    const p1 = newDeck.pop()!;
    const d1 = newDeck.pop()!;
    const p2 = newDeck.pop()!;
    const d2 = newDeck.pop()!;
    
    setDeck(newDeck);
    setPlayerHand([p1, p2]);
    setDealerHand([d1, d2]);
    setGameState('playing');
    
    // Check for blackjacks
    const playerBJ = isBlackjack([p1, p2]);
    const dealerBJ = isBlackjack([d1, d2]);
    
    if (playerBJ || dealerBJ) {
      setGameState('completed');
      
      let result: 'win' | 'lose' | 'push' | 'blackjack' = 'push';
      let profit = 0;
      
      if (playerBJ && dealerBJ) {
        result = 'push';
        profit = 0;
      } else if (playerBJ) {
        result = 'blackjack';
        profit = betAmount * 1.5;
      } else if (dealerBJ) {
        result = 'lose';
        profit = -betAmount;
      }
      
      const newBet: Bet = {
        id: Date.now().toString(),
        amount: betAmount,
        playerScore: calculateScore([p1, p2]),
        dealerScore: calculateScore([d1, d2]),
        result,
        profit,
        timestamp: new Date()
      };
      
      setMyBets(prev => [newBet, ...prev].slice(0, 50));
      
      if (result === 'blackjack') {
        toast.success(`Blackjack! Won $${profit.toFixed(2)}!`);
      } else if (result === 'lose') {
        toast.error('Dealer has Blackjack!');
      } else {
        toast('Push - no winner');
      }
    }
  }, [betAmount, isAuthenticated, createDeck]);

  // Player hit
  const hit = useCallback(() => {
    if (gameState !== 'playing') return;
    
    const newCard = deck.pop()!;
    const newHand = [...playerHand, newCard];
    setPlayerHand(newHand);
    setDeck([...deck]);
    
    // Check for bust
    if (calculateScore(newHand) > 21) {
      setGameState('completed');
      
      const profit = -betAmount;
      const newBet: Bet = {
        id: Date.now().toString(),
        amount: betAmount,
        playerScore: calculateScore(newHand),
        dealerScore: calculateScore(dealerHand),
        result: 'lose',
        profit,
        timestamp: new Date()
      };
      
      setMyBets(prev => [newBet, ...prev].slice(0, 50));
      toast.error('Bust! You lost.');
    }
  }, [gameState, deck, playerHand, dealerHand, betAmount]);

  // Player stand
  const stand = useCallback(() => {
    if (gameState !== 'playing') return;
    
    setGameState('dealerTurn');
    
    // Dealer plays
    let currentDeck = [...deck];
    let currentDealerHand = [...dealerHand];
    
    const dealerPlay = () => {
      const score = calculateScore(currentDealerHand);
      
      if (score < 17) {
        const newCard = currentDeck.pop()!;
        currentDealerHand.push(newCard);
        setDealerHand([...currentDealerHand]);
        setTimeout(dealerPlay, 500);
      } else {
        // Compare scores
        const playerScore = calculateScore(playerHand);
        const dealerScore = calculateScore(currentDealerHand);
        
        let result: 'win' | 'lose' | 'push' = 'push';
        let profit = 0;
        
        if (dealerScore > 21) {
          result = 'win';
          profit = betAmount;
        } else if (playerScore > dealerScore) {
          result = 'win';
          profit = betAmount;
        } else if (dealerScore > playerScore) {
          result = 'lose';
          profit = -betAmount;
        } else {
          result = 'push';
          profit = 0;
        }
        
        const newBet: Bet = {
          id: Date.now().toString(),
          amount: betAmount,
          playerScore,
          dealerScore,
          result,
          profit,
          timestamp: new Date()
        };
        
        setMyBets(prev => [newBet, ...prev].slice(0, 50));
        setGameState('completed');
        
        if (result === 'win') {
          toast.success(`You won $${profit.toFixed(2)}!`);
        } else if (result === 'lose') {
          toast.error(`Dealer wins with ${dealerScore}`);
        } else {
          toast('Push - no winner');
        }
      }
    };
    
    setTimeout(dealerPlay, 500);
  }, [gameState, deck, playerHand, dealerHand, betAmount]);

  // Double down
  const doubleDown = useCallback(() => {
    if (gameState !== 'playing') return;
    if (playerHand.length !== 2) {
      toast.error('Can only double down on first two cards');
      return;
    }
    
    // Double the bet
    const newCard = deck.pop()!;
    const newHand = [...playerHand, newCard];
    setPlayerHand(newHand);
    setDeck([...deck]);
    
    const score = calculateScore(newHand);
    
    if (score > 21) {
      // Bust
      setGameState('completed');
      const profit = -betAmount * 2;
      
      const newBet: Bet = {
        id: Date.now().toString(),
        amount: betAmount * 2,
        playerScore: score,
        dealerScore: calculateScore(dealerHand),
        result: 'lose',
        profit,
        timestamp: new Date()
      };
      
      setMyBets(prev => [newBet, ...prev].slice(0, 50));
      toast.error('Bust! You lost.');
    } else {
      // Auto-stand after double
      stand();
    }
  }, [gameState, deck, playerHand, dealerHand, betAmount, stand]);

  // Render card
  const renderCard = (card: Card | null, hidden: boolean = false) => {
    if (hidden) {
      return (
        <div className="w-16 h-24 bg-gradient-to-br from-red-700 to-red-900 rounded-lg border-2 border-white/30 flex items-center justify-center">
          <span className="text-2xl">?</span>
        </div>
      );
    }
    
    if (!card) return null;
    
    const isRed = card.suit === '♥' || card.suit === '♦';
    
    return (
      <motion.div
        initial={{ rotateY: 0 }}
        animate={{ rotateY: 0 }}
        className="w-16 h-24 bg-white rounded-lg shadow-lg flex flex-col items-center justify-between p-1"
      >
        <span className={`text-sm font-bold ${isRed ? 'text-red-600' : 'text-black'}`}>
          {card.rank}
        </span>
        <span className={`text-2xl ${isRed ? 'text-red-600' : 'text-black'}`}>
          {card.suit}
        </span>
        <span className={`text-sm font-bold ${isRed ? 'text-red-600' : 'text-black'}`}>
          {card.rank}
        </span>
      </motion.div>
    );
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Blackjack</h1>
            <p className="text-gray-400">Beat the dealer without going over 21!</p>
          </div>
          <div className="glass px-4 py-2 rounded-lg">
            <span className="text-gray-400 text-sm">Balance</span>
            <p className="text-xl font-mono text-green-400">${user?.balance?.toFixed(2) || '0.00'}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Game Table */}
            <div className="glass rounded-2xl p-6">
              {/* Dealer's Hand */}
              <div className="mb-8">
                <div className="flex justify-between items-center mb-2">
                  <h3 className="text-lg font-bold">Dealer's Hand</h3>
                  {gameState !== 'betting' && (
                    <span className="text-xl font-mono text-white">
                      {gameState === 'dealerTurn' || gameState === 'completed' 
                        ? calculateScore(dealerHand) 
                        : '?'
                      }
                    </span>
                  )}
                </div>
                <div className="flex gap-2 justify-center flex-wrap">
                  {dealerHand.map((card, index) => (
                    <motion.div
                      key={`dealer-${index}`}
                      initial={{ scale: 0, rotate: -10 }}
                      animate={{ scale: 1, rotate: 0 }}
                      transition={{ delay: index * 0.1 }}
                    >
                      {renderCard(
                        gameState === 'betting' || (gameState === 'playing' && index === 0) 
                          ? (index === 0 ? null : card) 
                          : card,
                        gameState === 'playing' && index === 0
                      )}
                    </motion.div>
                  ))}
                </div>
              </div>

              {/* Player's Hand */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <h3 className="text-lg font-bold">Your Hand</h3>
                  {gameState !== 'betting' && (
                    <span className="text-xl font-mono text-white">
                      {calculateScore(playerHand)}
                      {isBlackjack(playerHand) && ' - BLACKJACK!'}
                    </span>
                  )}
                </div>
                <div className="flex gap-2 justify-center flex-wrap">
                  {playerHand.map((card, index) => (
                    <motion.div
                      key={`player-${index}`}
                      initial={{ scale: 0, rotate: 10 }}
                      animate={{ scale: 1, rotate: 0 }}
                      transition={{ delay: index * 0.1 }}
                    >
                      {renderCard(card)}
                    </motion.div>
                  ))}
                </div>
              </div>

              {/* Result */}
              {gameState === 'completed' && (
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className="mt-6 text-center p-4 rounded-xl bg-tiger-surface/50"
                >
                  {isBlackjack(playerHand) ? (
                    <p className="text-2xl font-bold text-yellow-400">BLACKJACK!</p>
                  ) : calculateScore(playerHand) > 21 ? (
                    <p className="text-2xl font-bold text-red-400">BUST!</p>
                  ) : calculateScore(dealerHand) > 21 ? (
                    <p className="text-2xl font-bold text-green-400">DEALER BUSTS - YOU WIN!</p>
                  ) : calculateScore(playerHand) > calculateScore(dealerHand) ? (
                    <p className="text-2xl font-bold text-green-400">YOU WIN!</p>
                  ) : calculateScore(playerHand) < calculateScore(dealerHand) ? (
                    <p className="text-2xl font-bold text-red-400">DEALER WINS</p>
                  ) : (
                    <p className="text-2xl font-bold text-gray-400">PUSH</p>
                  )}
                </motion.div>
              )}
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {gameState === 'betting' ? (
                <>
                  {/* Bet Amount */}
                  <div className="mb-6">
                    <label className="block text-gray-400 text-sm mb-2">Bet Amount</label>
                    <div className="flex gap-2 mb-3">
                      {[1, 10, 25, 100, 500].map(amount => (
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
                      step={1}
                      min={1}
                    />
                  </div>

                  <button
                    onClick={dealGame}
                    className="w-full py-4 bg-green-500 hover:bg-green-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Deal
                  </button>
                </>
              ) : gameState === 'completed' ? (
                <button
                  onClick={() => { setPlayerHand([]); setDealerHand([]); setGameState('betting'); }}
                  className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                >
                  New Hand
                </button>
              ) : (
                <div className="flex gap-4">
                  <button
                    onClick={hit}
                    className="flex-1 py-4 bg-blue-500 hover:bg-blue-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Hit
                  </button>
                  <button
                    onClick={stand}
                    className="flex-1 py-4 bg-red-500 hover:bg-red-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Stand
                  </button>
                  {playerHand.length === 2 && (
                    <button
                      onClick={doubleDown}
                      className="flex-1 py-4 bg-yellow-500 hover:bg-yellow-600 rounded-xl font-bold text-xl transition-all"
                    >
                      Double
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Rules */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Rules</h3>
              <ul className="text-gray-400 text-sm space-y-2">
                <li>• Get closer to 21 than dealer</li>
                <li>• Face cards = 10, Ace = 1 or 11</li>
                <li>• Blackjack pays 3:2</li>
                <li>• Dealer stands on 17</li>
                <li>• Bust = automatic loss</li>
              </ul>
            </div>

            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Hands</h3>
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {myBets.slice(0, 10).map(bet => (
                  <div 
                    key={bet.id} 
                    className={`p-3 rounded-lg ${
                      bet.result === 'win' || bet.result === 'blackjack' 
                        ? 'bg-green-500/20' 
                        : bet.result === 'push'
                          ? 'bg-gray-500/20'
                          : 'bg-red-500/20'
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className={`font-bold ${
                        bet.result === 'win' || bet.result === 'blackjack' 
                          ? 'text-green-400' 
                          : bet.result === 'push'
                            ? 'text-gray-400'
                            : 'text-red-400'
                      }`}>
                        {bet.result.toUpperCase()}
                      </span>
                    </div>
                    <div className="flex justify-between text-xs text-gray-400 mt-1">
                      <span>You: {bet.playerScore}</span>
                      <span>Dealer: {bet.dealerScore}</span>
                    </div>
                    {bet.profit !== 0 && (
                      <p className={`text-sm ${bet.profit > 0 ? 'text-green-300' : 'text-red-300'}`}>
                        {bet.profit > 0 ? '+' : ''}${bet.profit.toFixed(2)}
                      </p>
                    )}
                  </div>
                ))}
                {myBets.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No hands played</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
