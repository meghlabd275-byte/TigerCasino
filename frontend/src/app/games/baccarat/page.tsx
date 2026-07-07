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
  betType: 'player' | 'banker' | 'tie';
  playerScore: number;
  bankerScore: number;
  result: 'player' | 'banker' | 'tie';
  profit: number;
  timestamp: Date;
}

// Calculate Baccarat score (only last digit matters)
const calculateScore = (cards: Card[]): number => {
  let score = 0;
  
  cards.forEach(card => {
    if (card.rank === 'A') {
      score += 1;
    } else if (['K', 'Q', 'J', '10'].includes(card.rank)) {
      score += 0;
    } else {
      score += card.value;
    }
  });
  
  return score % 10;
};

export default function BaccaratGame() {
  const { user, isAuthenticated } = useAuthStore();
  const { betAmount, setBetAmount } = useGameStore();
  
  const [playerHand, setPlayerHand] = useState<Card[]>([]);
  const [bankerHand, setBankerHand] = useState<Card[]>([]);
  const [deck, setDeck] = useState<Card[]>([]);
  const [gameState, setGameState] = useState<'betting' | 'dealing' | 'thirdCard' | 'completed'>('betting');
  const [selectedBet, setSelectedBet] = useState<'player' | 'banker' | 'tie'>('player');
  const [myBets, setMyBets] = useState<Bet[]>([]);

  // Create and shuffle deck
  const createDeck = useCallback((): Card[] => {
    const newDeck: Card[] = [];
    for (const suit of SUITS) {
      for (const rank of RANKS) {
        let value = parseInt(rank);
        if (isNaN(value)) {
          value = rank === 'A' ? 1 : 0;
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

  // Deal game
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
    
    // Deal initial 4 cards: Player, Banker, Player, Banker
    const p1 = newDeck.pop()!;
    const b1 = newDeck.pop()!;
    const p2 = newDeck.pop()!;
    const b2 = newDeck.pop()!;
    
    setDeck(newDeck);
    setPlayerHand([p1, p2]);
    setBankerHand([b1, b2]);
    setGameState('dealing');
    
    // Check for natural win
    const playerScore = calculateScore([p1, p2]);
    const bankerScore = calculateScore([b1, b2]);
    
    if (playerScore >= 8 || bankerScore >= 8) {
      // Natural - game over
      setGameState('completed');
      setTimeout(() => determineWinner(playerScore, bankerScore, []), 500);
    } else {
      // Check if player needs third card
      setTimeout(() => {
        if (playerScore <= 5) {
          setGameState('thirdCard');
          // Player draws third card
          const p3 = newDeck.pop()!;
          setPlayerHand(prev => [...prev, p3]);
          
          // Check if banker needs third card
          const newPlayerScore = calculateScore([...playerHand, p3]);
          let bankerShouldDraw = false;
          
          // Baccarat third card rules
          if (newPlayerScore <= 2) {
            bankerShouldDraw = true;
          } else if (newPlayerScore === 3 && calculateScore([b1, b2]) !== 8) {
            bankerShouldDraw = true;
          } else if (newPlayerScore === 4 && calculateScore([b1, b2]) >= 2 && calculateScore([b1, b2]) <= 7) {
            bankerShouldDraw = true;
          } else if (newPlayerScore === 5 && calculateScore([b1, b2]) >= 4 && calculateScore([b1, b2]) <= 7) {
            bankerShouldDraw = true;
          } else if (newPlayerScore === 6 && calculateScore([b1, b2]) >= 6 && calculateScore([b1, b2]) <= 7) {
            bankerShouldDraw = true;
          }
          
          if (bankerShouldDraw) {
            setTimeout(() => {
              const b3 = newDeck.pop()!;
              setBankerHand(prev => [...prev, b3]);
              setDeck([...newDeck]);
              setGameState('completed');
              setTimeout(() => determineWinner(
                calculateScore([p1, p2, p3]),
                calculateScore([b1, b2, b3]),
                [p3]
              ), 500);
            }, 1000);
          } else {
            setDeck([...newDeck]);
            setGameState('completed');
            setTimeout(() => determineWinner(newPlayerScore, calculateScore([b1, b2]), [p3]), 500);
          }
        } else if (bankerScore <= 5) {
          // Banker draws
          setGameState('thirdCard');
          setTimeout(() => {
            const b3 = newDeck.pop()!;
            setBankerHand(prev => [...prev, b3]);
            setDeck([...newDeck]);
            setGameState('completed');
            setTimeout(() => determineWinner(
              playerScore,
              calculateScore([b1, b2, b3]),
              []
            ), 500);
          }, 1000);
        } else {
          setDeck([...newDeck]);
          setGameState('completed');
          setTimeout(() => determineWinner(playerScore, bankerScore, []), 500);
        }
      }, 500);
    }
  }, [betAmount, isAuthenticated, createDeck, playerHand, bankerHand]);

  // Determine winner
  const determineWinner = (playerScore: number, bankerScore: number, thirdCards: Card[]) => {
    let result: 'player' | 'banker' | 'tie' = 'tie';
    let multiplier = 0;
    
    if (playerScore > bankerScore) {
      result = 'player';
      multiplier = selectedBet === 'player' ? 2 : 0;
    } else if (bankerScore > playerScore) {
      result = 'banker';
      // Banker wins pay 1.95 (5% commission)
      multiplier = selectedBet === 'banker' ? 1.95 : 0;
    } else {
      result = 'tie';
      // Tie pays 8:1
      multiplier = selectedBet === 'tie' ? 9 : 0;
    }
    
    const profit = (betAmount * multiplier) - betAmount;
    
    const newBet: Bet = {
      id: Date.now().toString(),
      amount: betAmount,
      betType: selectedBet,
      playerScore,
      bankerScore,
      result,
      profit,
      timestamp: new Date()
    };
    
    setMyBets(prev => [newBet, ...prev].slice(0, 50));
    
    if (profit > 0) {
      toast.success(`Won $${profit.toFixed(2)}!`);
    } else if (profit < 0) {
      toast.error(`Lost $${betAmount.toFixed(2)}`);
    } else {
      toast('Push - bet returned');
    }
  };

  // Render card
  const renderCard = (card: Card | null) => {
    if (!card) return null;
    
    const isRed = card.suit === '♥' || card.suit === '♦';
    
    return (
      <motion.div
        initial={{ scale: 0, rotate: -10 }}
        animate={{ scale: 1, rotate: 0 }}
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
            <h1 className="text-3xl md:text-4xl font-heading text-gradient">Baccarat</h1>
            <p className="text-gray-400">Bet on Player, Banker, or Tie!</p>
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
              {/* Banker's Hand */}
              <div className="mb-8">
                <div className="flex justify-between items-center mb-2">
                  <h3 className="text-lg font-bold text-blue-400">Banker</h3>
                  {gameState !== 'betting' && (
                    <span className="text-2xl font-mono text-white">
                      {calculateScore(bankerHand)}
                    </span>
                  )}
                </div>
                <div className="flex gap-2 justify-center flex-wrap">
                  {bankerHand.map((card, index) => (
                    <motion.div
                      key={`banker-${index}`}
                      initial={{ scale: 0, rotate: -10 }}
                      animate={{ scale: 1, rotate: 0 }}
                      transition={{ delay: index * 0.1 }}
                    >
                      {renderCard(card)}
                    </motion.div>
                  ))}
                </div>
              </div>

              {/* Player's Hand */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <h3 className="text-lg font-bold text-green-400">Player</h3>
                  {gameState !== 'betting' && (
                    <span className="text-2xl font-mono text-white">
                      {calculateScore(playerHand)}
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
              {gameState === 'completed' && playerHand.length > 0 && bankerHand.length > 0 && (
                <motion.div
                  initial={{ scale: 0.8, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className="mt-6 text-center p-4 rounded-xl bg-tiger-surface/50"
                >
                  <p className="text-3xl font-bold">
                    {calculateScore(playerHand) > calculateScore(bankerHand) ? (
                      <span className="text-green-400">PLAYER WINS</span>
                    ) : calculateScore(bankerHand) > calculateScore(playerHand) ? (
                      <span className="text-blue-400">BANKER WINS</span>
                    ) : (
                      <span className="text-yellow-400">TIE</span>
                    )}
                  </p>
                </motion.div>
              )}
            </div>

            {/* Controls */}
            <div className="glass rounded-2xl p-6">
              {gameState === 'betting' ? (
                <>
                  {/* Bet Selection */}
                  <div className="mb-6">
                    <label className="block text-gray-400 text-sm mb-3">Select Bet</label>
                    <div className="grid grid-cols-3 gap-4">
                      <button
                        onClick={() => setSelectedBet('player')}
                        className={`py-4 rounded-xl font-bold text-lg transition-all ${
                          selectedBet === 'player'
                            ? 'bg-green-500 text-white'
                            : 'bg-tiger-surface text-gray-400 hover:bg-green-500/30'
                        }`}
                      >
                        Player
                        <span className="block text-xs font-normal">2x</span>
                      </button>
                      <button
                        onClick={() => setSelectedBet('banker')}
                        className={`py-4 rounded-xl font-bold text-lg transition-all ${
                          selectedBet === 'banker'
                            ? 'bg-blue-500 text-white'
                            : 'bg-tiger-surface text-gray-400 hover:bg-blue-500/30'
                        }`}
                      >
                        Banker
                        <span className="block text-xs font-normal">1.95x</span>
                      </button>
                      <button
                        onClick={() => setSelectedBet('tie')}
                        className={`py-4 rounded-xl font-bold text-lg transition-all ${
                          selectedBet === 'tie'
                            ? 'bg-yellow-500 text-white'
                            : 'bg-tiger-surface text-gray-400 hover:bg-yellow-500/30'
                        }`}
                      >
                        Tie
                        <span className="block text-xs font-normal">9x</span>
                      </button>
                    </div>
                  </div>

                  {/* Bet Amount */}
                  <div className="mb-6">
                    <label className="block text-gray-400 text-sm mb-2">Bet Amount</label>
                    <div className="flex gap-2 mb-3">
                      {[1, 10, 50, 100, 500].map(amount => (
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
                    className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                  >
                    Deal
                  </button>
                </>
              ) : (
                <button
                  onClick={() => { setPlayerHand([]); setBankerHand([]); setGameState('betting'); }}
                  className="w-full py-4 bg-primary-500 hover:bg-primary-600 rounded-xl font-bold text-xl transition-all"
                >
                  New Hand
                </button>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Payouts */}
            <div className="glass rounded-2xl p-4">
              <h3 className="text-lg font-bold mb-4">Payouts</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-green-400">Player:</span>
                  <span className="text-white">2x (1:1)</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-blue-400">Banker:</span>
                  <span className="text-white">1.95x (19:20)</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-yellow-400">Tie:</span>
                  <span className="text-white">9x (8:1)</span>
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
                      bet.profit > 0 
                        ? 'bg-green-500/20' 
                        : bet.profit === 0
                          ? 'bg-gray-500/20'
                          : 'bg-red-500/20'
                    }`}
                  >
                    <div className="flex justify-between items-center">
                      <span className="font-mono text-sm">${bet.amount.toFixed(2)}</span>
                      <span className={`font-bold ${
                        bet.result === 'player' ? 'text-green-400' :
                        bet.result === 'banker' ? 'text-blue-400' :
                        'text-yellow-400'
                      }`}>
                        {bet.result.toUpperCase()}
                      </span>
                    </div>
                    <div className="flex justify-between text-xs text-gray-400 mt-1">
                      <span>P: {bet.playerScore}</span>
                      <span>B: {bet.bankerScore}</span>
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
