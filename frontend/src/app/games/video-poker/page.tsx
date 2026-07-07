'use client';

import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import toast from 'react-hot-toast';
import { useAuthStore, useGameStore } from '@/store';

// Card configuration
const SUITS = ['♠', '♥', '♦', '♣'];
const RANKS = ['A', '2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K'];

// Video Poker paytable (Jacks or Better)
const PAYTABLE: Record<string, number> = {
  'royal_flush': 800,
  'straight_flush': 50,
  'four_kind': 25,
  'full_house': 9,
  'flush': 6,
  'straight': 4,
  'three_kind': 3,
  'two_pair': 2,
  'jacks_or_better': 1,
};

interface Card {
  rank: string;
  suit: string;
  value: number;
  held: boolean;
}

interface Bet {
  id: string;
  amount: number;
  hand: string;
  won: boolean;
  profit: number;
  timestamp: Date;
}

// Create and shuffle deck
const createDeck = (): Card[] => {
  const deck: Card[] = [];
  for (const suit of SUITS) {
    for (const rank of RANKS) {
      let value: number;
      if (rank === 'A') value = 14;
      else if (rank === 'K') value = 13;
      else if (rank === 'Q') value = 12;
      else if (rank === 'J') value = 11;
      else value = parseInt(rank);
      
      deck.push({ rank, suit, value, held: false });
    }
  }
  
  // Shuffle
  for (let i = deck.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [deck[i], deck[j]] = [deck[j], deck[i]];
  }
  
  return deck;
};

// Evaluate hand
const evaluateHand = (cards: Card[]): string => {
  const ranks = cards.map(c => c.value).sort((a, b) => a - b);
  const suits = cards.map(c => c.suit);
  const rankCounts: Record<number, number> = {};
  
  ranks.forEach(r => {
    rankCounts[r] = (rankCounts[r] || 0) + 1;
  });
  
  const isFlush = suits.every(s => s === suits[0]);
  
  // Check for straight
  let isStraight = true;
  for (let i = 0; i < ranks.length - 1; i++) {
    if (ranks[i + 1] !== ranks[i] + 1) {
      isStraight = false;
      break;
    }
  }
  // Also check Ace-low straight (A-2-3-4-5)
  if (!isStraight && ranks[0] === 2 && ranks[4] === 14) {
    isStraight = true;
  }
  
  // Count pairs, trips, quads
  const counts = Object.values(rankCounts);
  
  // Royal flush (10-A, same suit)
  if (isFlush && isStraight && ranks[0] === 10 && ranks[4] === 14) {
    return 'royal_flush';
  }
  
  // Straight flush
  if (isFlush && isStraight) {
    return 'straight_flush';
  }
  
  // Four of a kind
  if (counts.includes(4)) {
    return 'four_kind';
  }
  
  // Full house
  if (counts.includes(3) && counts.includes(2)) {
    return 'full_house';
  }
  
  // Flush
  if (isFlush) {
    return 'flush';
  }
  
  // Straight
  if (isStraight) {
    return 'straight';
  }
  
  // Three of a kind
  if (counts.includes(3)) {
    return 'three_kind';
  }
  
  // Two pair
  if (counts.filter(c => c === 2).length === 2) {
    return 'two_pair';
  }
  
  // Jacks or better
  const hasJackOrBetter = ranks.some(r => (r === 11 || r === 12 || r === 13 || r === 14) && rankCounts[r] >= 1);
  if (hasJackOrBetter) {
    return 'jacks_or_better';
  }
  
  return 'none';
};

export default function VideoPokerGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [hand, setHand] = useState<Card[]>([]);
  const [gameState, setGameState] = useState<'initial' | 'drawing' | 'completed'>('initial');
  const [myBets, setMyBets] = useState<Bet[]>([]);

  // Deal initial hand
  const dealHand = useCallback(() => {
    if (!isAuthenticated) {
      toast.error('Please login to play');
      return;
    }
    
    if (betAmount <= 0) {
      toast.error('Please enter a valid bet amount');
      return;
    }

    const deck = createDeck();
    const newHand = deck.slice(0, 5);
    
    setHand(newHand);
    setGameState('drawing');
  }, [betAmount, isAuthenticated]);

  // Toggle hold on card
  const toggleHold = useCallback((index: number) => {
    if (gameState !== 'drawing') return;
    
    setHand(prev => prev.map((card, i) => 
      i === index ? { ...card, held: !card.held } : card
    ));
  }, [gameState]);

  // Draw new cards
  const drawCards = useCallback(() => {
    if (gameState !== 'drawing') return;
    
    const heldCards = hand.filter(c => c.held);
    const currentDeck = createDeck();
    
    // Remove held cards from consideration
    const availableCards = currentDeck.filter(c => 
      !heldCards.some(hc => hc.rank === c.rank && hc.suit === c.suit)
    );
    
    // Draw new cards to replace non-held cards
    const newCardsNeeded = 5 - heldCards.length;
    const newCards = availableCards.slice(0, newCardsNeeded);
    
    // Combine held and new cards
    const newHand: Card[] = [];
    hand.forEach((card, index) => {
      if (card.held) {
        newHand.push(card);
      } else {
        newHand.push(newCards.shift()!);
      }
    });
    
    // Call real API
    fetch(`/api/games/video_poker/bet`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        amount: betAmount
      })
    })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        toast.error(data.error);
        return;
      }

      setHand(newHand);
      setGameState('completed');

      const won = data.win_amount > 0;
      const profit = data.win_amount - betAmount;
      const handResult = evaluateHand(newHand);

      const newBet: Bet = {
        id: Date.now().toString(),
        amount: betAmount,
        hand: handResult.replace('_', ' '),
        won,
        profit,
        timestamp: new Date()
      };

      setMyBets(prev => [newBet, ...prev].slice(0, 50));

      if (won) {
        toast.success(`Won $${data.win_amount.toFixed(2)}!`);
      } else {
        toast.error('No winning hand');
      }
    })
    .catch(err => {
      toast.error('Failed to connect to game server');
    });
  }, [gameState, hand, betAmount]);

  // Render card
  const renderCard = (card: Card, index: number) => {
    const isRed = card.suit === '♥' || card.suit === '♦';
    
    return (
      <motion.div
        key={`${card.rank}-${card.suit}-${index}`}
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={() => toggleHold(index)}
        className={`relative cursor-pointer transition-all ${
          card.held ? 'ring-4 ring-primary-500' : ''
        }`}
      >
        <div className={`
          w-20 h-28 rounded-lg shadow-lg flex flex-col items-center justify-between p-2
          ${gameState === 'initial' ? 'bg-white' : card.held ? 'bg-primary-500/50' : 'bg-white/50'}
        `}>
          <span className={`text-sm font-bold ${isRed ? 'text-red-600' : 'text-black'}`}>
            {card.rank}
          </span>
          <span className={`text-3xl ${isRed ? 'text-red-600' : 'text-black'}`}>
            {card.suit}
          </span>
          <span className={`text-sm font-bold ${isRed ? 'text-red-600' : 'text-black'}`}>
            {card.rank}
          </span>
        </div>
        
        {/* Hold indicator */}
        {card.held && gameState === 'drawing' && (
          <div className="absolute -bottom-8 left-1/2 -translate-x-1/2 bg-primary-500 text-white text-xs px-2 py-1 rounded">
            HOLD
          </div>
        )}
      </motion.div>
    );
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Video Poker</h1>
            <p className="text-gray-400">Jacks or Better - Hold cards to make the best hand!</p>
          </div>
          <div className="glass px-4 py-2 rounded-lg">
            <span className="text-gray-400 text-sm">Balance</span>
            <p className="text-xl font-mono text-green-400">${user?.balance?.toFixed(2) || '0.00'}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Game Area */}
          <div className="lg:col-span-2 space-y-6">
            {/* Poker Table */}
            <div className="glass rounded-2xl p-8">
              <div className="flex flex-col items-center">
                {/* Hand */}
                <div className="flex gap-4 justify-center flex-wrap mb-8">
                  {hand.length > 0 ? (
                    hand.map((card, index) => (
                      <div key={index} className="mb-8">
                        {renderCard(card, index)}
                      </div>
                    ))
                  ) : (
                    // Empty card placeholders
                    Array(5).fill(0).map((_, i) => (
                      <div key={i} className="w-20 h-28 rounded-lg bg-white/10 flex items-center justify-center">
                        <span className="text-3xl text-white/30">{i + 1}</span>
                      </div>
                    ))
                  )}
                </div>
                
                {/* Result */}
                {gameState === 'completed' && hand.length > 0 && (
                  <motion.div
                    initial={{ scale: 0.8, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    className="text-center"
                  >
                    <p className="text-2xl font-bold text-yellow-400">
                      {evaluateHand(hand).replace('_', ' ').toUpperCase()}
                    </p>
                  </motion.div>
                )}
              </div>
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {gameState === 'initial' ? (
                <>
                  {/* Bet Amount */}
                  <div className="mb-6">
                    <label className="block text-gray-400 text-sm mb-2">Bet Amount (1-5 coins)</label>
                    <div className="flex gap-2 mb-3">
                      {[1, 2, 3, 4, 5].map(coins => (
                        <button
                          key={coins}
                          onClick={() => setBetAmount(coins)}
                          className={`flex-1 py-3 rounded-lg font-bold transition ${
                            betAmount === coins
                              ? 'bg-primary-500 text-white'
                              : 'bg-tiger-surface text-gray-400 hover:bg-primary-500/30'
                          }`}
                        >
                          {coins}x
                        </button>
                      ))}
                    </div>
                  </div>

                  <button
                    onClick={dealHand}
                    className="w-full py-4 bg-green-500 hover:bg-green-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Deal
                  </button>
                </>
              ) : gameState === 'drawing' ? (
                <>
                  <p className="text-center text-gray-400 mb-4">
                    Click cards to hold them, then click Draw
                  </p>
                  <button
                    onClick={drawCards}
                    className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Draw
                  </button>
                </>
              ) : (
                <button
                  onClick={() => { setHand([]); setGameState('initial'); }}
                  className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                >
                  Deal New Hand
                </button>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Paytable */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Paytable</h3>
              <div className="space-y-1 text-sm">
                <div className="flex justify-between text-yellow-400">
                  <span>Royal Flush</span>
                  <span>800x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Straight Flush</span>
                  <span className="text-white">50x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Four of a Kind</span>
                  <span className="text-white">25x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Full House</span>
                  <span className="text-white">9x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Flush</span>
                  <span className="text-white">6x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Straight</span>
                  <span className="text-white">4x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Three of a Kind</span>
                  <span className="text-white">3x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Two Pair</span>
                  <span className="text-white">2x</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Jacks or Better</span>
                  <span className="text-white">1x</span>
                </div>
              </div>
            </div>

            {/* My Bets */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">My Hands</h3>
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {myBets.slice(0, 10).map(bet => (
                  <div 
                    key={bet.id} 
                    className={`p-3 rounded-lg ${
                      bet.won ? 'bg-green-500/20' : 'bg-red-500/20'
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${(bet.amount * PAYTABLE[bet.hand.replace(' ', '_') as string] || 0).toFixed(2)}</span>
                      <span className="text-gray-400 text-xs capitalize">{bet.hand}</span>
                    </div>
                    <p className={bet.profit >= 0 ? 'text-green-400 text-sm' : 'text-red-400 text-sm'}>
                      {bet.profit >= 0 ? '+' : ''}${bet.profit.toFixed(2)}
                    </p>
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
