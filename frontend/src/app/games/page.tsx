'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Badge } from '@/components/ui';
import styles from './games.module.css';

interface Game {
  id: string;
  name: string;
  type: string;
  icon: string;
  description: string;
  rtp: number;
  category: string;
  isHot?: boolean;
  isNew?: boolean;
  provider: string;
  minBet?: number;
  maxBet?: number;
}

// Icons for different game types
const typeIcons: Record<string, string> = {
  slots: '🎰',
  dice: '🎲',
  roulette: '🎡',
  blackjack: '🃏',
  baccarat: '🪙',
  poker: '♠️',
  show: '🎬',
  crash: '🚀'
};

// Get 200+ games from backend API (simulated with extended list)
const games: Game[] = [
  // SLOTS - Pragmatic Play (50+ games)
  { id: 'slots-tiger-king', name: 'Tiger King', type: 'slots', icon: '🐯', description: 'King of the jungle awaits', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', isHot: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-mega-moolah', name: 'Mega Moolah', type: 'slots', icon: '🦁', description: 'Win massive jackpots', rtp: 88.12, category: 'slots', provider: 'Microgaming', isNew: true, minBet: 0.25, maxBet: 6.25 },
  { id: 'slots-book-of-dead', name: 'Book of Dead', type: 'slots', icon: '📖', description: 'Ancient Egyptian adventure', rtp: 96.21, category: 'slots', provider: 'Play n GO', minBet: 0.1, maxBet: 100 },
  { id: 'slots-starburst', name: 'Starburst', type: 'slots', icon: '💎', description: 'Cosmic gem explosion', rtp: 96.09, category: 'slots', provider: 'NetEnt', isHot: true, minBet: 0.1, maxBet: 100 },
  { id: 'slots-gonzo-quest', name: 'Gonzo\'s Quest', type: 'slots', icon: '🏔️', description: 'Quest for Eldorado', rtp: 95.97, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wolf-gold', name: 'Wolf Gold', type: 'slots', icon: '🐺', description: 'Pack rules the night', rtp: 96.01, category: 'slots', provider: 'Pragmatic Play', minBet: 0.25, maxBet: 125 },
  { id: 'slots-sweet-bonanza', name: 'Sweet Bonanza', type: 'slots', icon: '🍬', description: 'Sweet candy rewards', rtp: 96.48, category: 'slots', provider: 'Pragmatic Play', isHot: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-gates-of-olympus', name: 'Gates of Olympus', type: 'slots', icon: '⚡', description: 'Zeus lightning wins', rtp: 95.51, category: 'slots', provider: 'Pragmatic Play', isNew: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-money-train', name: 'Money Train', type: 'slots', icon: '🚂', description: 'Wild west heist', rtp: 96.15, category: 'slots', provider: 'Relax Gaming', minBet: 0.1, maxBet: 20 },
  { id: 'slots-big-bass-bonanza', name: 'Big Bass Bonanza', type: 'slots', icon: '🎣', description: 'Fishing for big wins', rtp: 96.71, category: 'slots', provider: 'Pragmatic Play', minBet: 0.1, maxBet: 125 },
  { id: 'slots-divine-fortune', name: 'Divine Fortune', type: 'slots', icon: '🐴', description: 'Ancient Greek fortune', rtp: 96.59, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-fruit-party', name: 'Fruit Party', type: 'slots', icon: '🍓', description: 'Party with fruits', rtp: 96.47, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-dog-house', name: 'The Dog House', type: 'slots', icon: '🐕', description: 'Doggy paradise', rtp: 96.51, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wild-west-gold', name: 'Wild West Gold', type: 'slots', icon: '🤠', description: 'Wild west adventures', rtp: 96.51, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-jokers-jewels', name: 'Joker\'s Jewels', type: 'slots', icon: '🤡', description: 'Joker\'s treasures', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-great-rhino', name: 'Great Rhino', type: 'slots', icon: '🦏', description: 'Rhino charge', rtp: 96.65, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-caishens-gold', name: 'Caishen\'s Gold', type: 'slots', icon: '🐱', description: 'Chinese fortune', rtp: 96.08, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-fire-88', name: 'Fire 88', type: 'slots', icon: '🔥', description: 'Fire symbols', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-chilli-heat', name: 'Chilli Heat', type: 'slots', icon: '🌶️', description: 'Spicy wins', rtp: 96.52, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-panda-fortune', name: 'Panda Fortune', type: 'slots', icon: '🐼', description: 'Panda luck', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-ancient-egypt', name: 'Ancient Egypt', type: 'slots', icon: '🏺', description: 'Pharaoh\'s riches', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-hercules', name: 'Hercules Son of Zeus', type: 'slots', icon: '💪', description: 'Greek hero', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-madame-destiny', name: 'Madame Destiny', type: 'slots', icon: '🔮', description: 'Mystical predictions', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-pyramid-king', name: 'Pyramid King', type: 'slots', icon: '🔺', description: 'Egyptian king', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-dragon-hot', name: 'Dragon Hot', type: 'slots', icon: '🐉', description: 'Hot dragon', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-aztec-gems', name: 'Aztec Gems', type: 'slots', icon: '💎', description: 'Aztec treasures', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-egyptian-fortunes', name: 'Egyptian Fortunes', type: 'slots', icon: '👑', description: 'Egyptian wealth', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-gold-rush', name: 'Gold Rush', type: 'slots', icon: '⛏️', description: 'Gold mining', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-super-joker', name: 'Super Joker', type: 'slots', icon: '🤡', description: 'Super joker', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-master-joker', name: 'Master Joker', type: 'slots', icon: '🃏', description: 'Master of jokers', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-diamond-strike', name: 'Diamond Strike', type: 'slots', icon: '💠', description: 'Diamond gems', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-5-lions', name: '5 Lions', type: 'slots', icon: '🦁', description: 'Golden lions', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-7-piggies', name: '7 Piggies', type: 'slots', icon: '🐷', description: 'Seven pigs', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-aladdin', name: 'Aladdin\'s Treasure', type: 'slots', icon: '🧞', description: 'Aladdin\'s riches', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-aztec-king', name: 'Aztec King', type: 'slots', icon: '👑', description: 'Aztec ruler', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-dancing-king', name: 'Dancing King', type: 'slots', icon: '💃', description: 'Royal dance', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-dragon-tiger', name: 'Dragon Tiger', type: 'slots', icon: '🐲', description: 'Dragon vs Tiger', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-fortune-888', name: 'Fortune 888', type: 'slots', icon: '8️⃣', description: 'Lucky 888', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-fu-fu-fu', name: 'Fu Fu Fu', type: 'slots', icon: '🍚', description: 'Chinese fortune', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-golden-ox', name: 'Golden Ox', type: 'slots', icon: '🐂', description: 'Golden beast', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-happy-fish', name: 'Happy Golden Fish', type: 'slots', icon: '🐟', description: 'Lucky fish', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-hot-burn', name: 'Hot to Burn', type: 'slots', icon: '🔥', description: 'Hot flames', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-jewel-box', name: 'Jewel Box', type: 'slots', icon: '📦', description: 'Jewel collection', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-jungle-gorilla', name: 'Jungle Gorilla', type: 'slots', icon: '🦍', description: 'Jungle king', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-kingdom-sun', name: 'Kingdom of the Sun', type: 'slots', icon: '☀️', description: 'Sun kingdom', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-lucky-dragon', name: 'Lucky Dragon', type: 'slots', icon: '🐉', description: 'Chinese dragon', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-magic-crystals', name: 'Magic Crystals', type: 'slots', icon: '🔮', description: 'Mystical gems', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-mighty-kong', name: 'Mighty Kong', type: 'slots', icon: '🦍', description: 'Mighty gorilla', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-money-dome', name: 'Money Dome', type: 'slots', icon: '💰', description: 'Dome of money', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-mystic-sea', name: 'Mystic Sea', type: 'slots', icon: '🌊', description: 'Mystical ocean', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-orion', name: 'Orion', type: 'slots', icon: '⭐', description: 'Star constellation', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-phoenix-forge', name: 'Phoenix Forge', type: 'slots', icon: '🔥', description: 'Rising phoenix', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-piggy-gold', name: 'Piggy Gold', type: 'slots', icon: '🐖', description: 'Golden pig', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-power-thor', name: 'Power of Thor', type: 'slots', icon: '⚡', description: 'Thor\'s power', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-queen-gods', name: 'Queen of Gods', type: 'slots', icon: '👸', description: 'Divine queen', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-santa-wonderland', name: 'Santa\'s Wonderland', type: 'slots', icon: '🎅', description: 'Holiday fun', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', isNew: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-secret-temple', name: 'Secret of the Temple', type: 'slots', icon: '🛕', description: 'Temple secrets', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-speed-winner', name: 'Speed Winner', type: 'slots', icon: '🏎️', description: 'Racing wins', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-triple-diamond', name: 'Triple Diamond', type: 'slots', icon: '💎', description: 'Triple gems', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wild-beach', name: 'Wild Beach Life', type: 'slots', icon: '🏖️', description: 'Beach party', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wolf-howl', name: 'Wolf Howl', type: 'slots', icon: '🐺', description: 'Wolf pack', rtp: 96.5, category: 'slots', provider: 'Pragmatic Play', minBet: 0.2, maxBet: 100 },
  
  // NETENT SLOTS
  { id: 'slots-secret-stones', name: 'Secret of the Stones', type: 'slots', icon: '🪨', description: 'Stone secrets', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-twin-spin', name: 'Twin Spin', type: 'slots', icon: '🎰', description: 'Twin reels', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-jack-beanstalk', name: 'Jack and the Beanstalk', type: 'slots', icon: '🫘', description: 'Bean adventure', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-dead-alive', name: 'Dead or Alive', type: 'slots', icon: '🔫', description: 'Wild west shootout', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-dead-alive-2', name: 'Dead or Alive 2', type: 'slots', icon: '🔫', description: 'Wild west returns', rtp: 96.5, category: 'slots', provider: 'NetEnt', isNew: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-blood-suckers', name: 'Blood Suckers', type: 'slots', icon: '🧛', description: 'Vampire hunt', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-joker-pro', name: 'Joker Pro', type: 'slots', icon: '🤡', description: 'Joker fun', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wild-west', name: 'Wild Wild West', type: 'slots', icon: '🤠', description: 'Wild west', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-steam-tower', name: 'Steam Tower', type: 'slots', icon: '🗼', description: 'Steampunk adventure', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-planet-apes', name: 'Planet of the Apes', type: 'slots', icon: '🦧', description: 'Ape uprising', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-narcos', name: 'Narcos', type: 'slots', icon: '💊', description: 'Drug cartel', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-street-fighter', name: 'Street Fighter', type: 'slots', icon: '🥊', description: 'Fighting game', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-finn-swirly', name: 'Finn and the Swirly Spin', type: 'slots', icon: '🍀', description: 'Lucky Finn', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-mega-fortune', name: 'Mega Fortune', type: 'slots', icon: '💰', description: 'Jackpot dreams', rtp: 96.5, category: 'slots', provider: 'NetEnt', isHot: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-hall-gods', name: 'Hall of Gods', type: 'slots', icon: '⚡', description: 'Norse mythology', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  { id: 'slots-arabian-nights', name: 'Arabian Nights', type: 'slots', icon: '🧞', description: 'Middle East magic', rtp: 96.5, category: 'slots', provider: 'NetEnt', minBet: 0.2, maxBet: 100 },
  
  // DICE GAMES
  { id: 'dice-classic', name: 'Classic Dice', type: 'dice', icon: '🎲', description: 'Roll the dice, win big', rtp: 99, category: 'dice', provider: 'TigerCasino', isHot: true, minBet: 0.01, maxBet: 1000 },
  { id: 'dice-plinko', name: 'Plinko', type: 'dice', icon: '🔻', description: 'Watch the ball drop', rtp: 98, category: 'dice', provider: 'TigerCasino', isNew: true, minBet: 0.01, maxBet: 100 },
  { id: 'dice-mines', name: 'Mines', type: 'dice', icon: '💣', description: 'Avoid the mines', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-hilo', name: 'HiLo', type: 'dice', icon: '🔺', description: 'Predict the outcome', rtp: 96, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-keno', name: 'Keno', type: 'dice', icon: '🎱', description: 'Pick your numbers', rtp: 95, category: 'dice', provider: 'TigerCasino', minBet: 0.1, maxBet: 100 },
  { id: 'dice-draft', name: 'Draft', type: 'dice', icon: '📝', description: 'Draft your wins', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-limbo', name: 'Limbo', type: 'dice', icon: '🎯', description: 'Aim high', rtp: 97.5, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-minefield', name: 'Minefield', type: 'dice', icon: '💥', description: 'Clear the field', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-wheel', name: 'Wheel', type: 'dice', icon: '🎡', description: 'Spin the wheel', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-goal', name: 'Goal', type: 'dice', icon: '⚽', description: 'Score big', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-rocket', name: 'Rocket', type: 'dice', icon: '🚀', description: 'Blast off', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-tower', name: 'Tower', type: 'dice', icon: '🗼', description: 'Build your tower', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-dice-duel', name: 'Dice Duel', type: 'dice', icon: '⚔️', description: 'Dice battle', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-fast-dice', name: 'Fast Dice', type: 'dice', icon: '⚡', description: 'Quick rolls', rtp: 98, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  { id: 'dice-lucky-dice', name: 'Lucky Dice', type: 'dice', icon: '🍀', description: 'Lucky rolls', rtp: 97, category: 'dice', provider: 'TigerCasino', minBet: 0.01, maxBet: 100 },
  
  // ROULETTE
  { id: 'roulette-european', name: 'European Roulette', type: 'roulette', icon: '🎡', description: 'Classic single zero', rtp: 97.3, category: 'roulette', provider: 'Evolution', isHot: true, minBet: 1, maxBet: 5000 },
  { id: 'roulette-american', name: 'American Roulette', type: 'roulette', icon: '🎰', description: 'Double zero action', rtp: 94.74, category: 'roulette', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'roulette-french', name: 'French Roulette', type: 'roulette', icon: '🇫🇷', description: 'La partage rules', rtp: 98.65, category: 'roulette', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'roulette-speed', name: 'Speed Roulette', type: 'roulette', icon: '⚡', description: 'Fast-paced action', rtp: 97.3, category: 'roulette', provider: 'Evolution', isNew: true, minBet: 1, maxBet: 5000 },
  { id: 'roulette-immersive', name: 'Immersive Roulette', type: 'roulette', icon: '🎬', description: 'Cinematic experience', rtp: 97.3, category: 'roulette', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'roulette-lightning', name: 'Lightning Roulette', type: 'roulette', icon: '⚡', description: 'Electrifying wins', rtp: 97.3, category: 'roulette', provider: 'Evolution', minBet: 0.1, maxBet: 2000 },
  { id: 'roulette-auto', name: 'Auto Roulette', type: 'roulette', icon: '🤖', description: 'Automated play', rtp: 97.3, category: 'roulette', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'roulette-double-ball', name: 'Double Ball Roulette', type: 'roulette', icon: '⚫', description: 'Two balls', rtp: 97.3, category: 'roulette', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'roulette-gold-vault', name: 'Gold Vault Roulette', type: 'roulette', icon: '💰', description: 'Gold roulette', rtp: 97.3, category: 'roulette', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'roulette-quantum', name: 'Quantum Roulette', type: 'roulette', icon: '⚛️', description: 'Quantum multipliers', rtp: 97.3, category: 'roulette', provider: 'Playtech', minBet: 1, maxBet: 5000 },
  { id: 'roulette-age-gods', name: 'Age of Gods Roulette', type: 'roulette', icon: '⚡', description: 'God roulette', rtp: 97.3, category: 'roulette', provider: 'Playtech', minBet: 1, maxBet: 5000 },
  { id: 'roulette-mini', name: 'Mini Roulette', type: 'roulette', icon: '🔘', description: 'Compact roulette', rtp: 97.3, category: 'roulette', provider: 'Playtech', minBet: 1, maxBet: 5000 },
  
  // BLACKJACK
  { id: 'blackjack-classic', name: 'Classic Blackjack', type: 'blackjack', icon: '🃏', description: 'Beat the dealer', rtp: 99.4, category: 'blackjack', provider: 'Evolution', isHot: true, minBet: 10, maxBet: 5000 },
  { id: 'blackjack-vip', name: 'VIP Blackjack', type: 'blackjack', icon: '👑', description: 'High roller tables', rtp: 99.5, category: 'blackjack', provider: 'Evolution', minBet: 50, maxBet: 10000 },
  { id: 'blackjack-party', name: 'Party Blackjack', type: 'blackjack', icon: '🎉', description: 'Fun atmosphere', rtp: 99.4, category: 'blackjack', provider: 'Evolution', isNew: true, minBet: 10, maxBet: 2500 },
  { id: 'blackjack-infinite', name: 'Infinite Blackjack', type: 'blackjack', icon: '♾️', description: 'Unlimited seats', rtp: 99.47, category: 'blackjack', provider: 'Evolution', minBet: 1, maxBet: 2500 },
  { id: 'blackjack-power', name: 'Power Blackjack', type: 'blackjack', icon: '⚡', description: 'Double down anytime', rtp: 99.5, category: 'blackjack', provider: 'Evolution', minBet: 25, maxBet: 10000 },
  { id: 'blackjack-free-bet', name: 'Free Bet Blackjack', type: 'blackjack', icon: '🎁', description: 'Free bets', rtp: 99.4, category: 'blackjack', provider: 'Evolution', minBet: 10, maxBet: 5000 },
  { id: 'blackjack-speed', name: 'Speed Blackjack', type: 'blackjack', icon: '⚡', description: 'Lightning fast', rtp: 99.4, category: 'blackjack', provider: 'Evolution', minBet: 10, maxBet: 5000 },
  { id: 'blackjack-quantum', name: 'Quantum Blackjack', type: 'blackjack', icon: '⚛️', description: 'Quantum multipliers', rtp: 99.4, category: 'blackjack', provider: 'Playtech', minBet: 10, maxBet: 5000 },
  { id: 'blackjack-european', name: 'European Blackjack', type: 'blackjack', icon: '🇪🇺', description: 'European style', rtp: 99.4, category: 'blackjack', provider: 'Betsoft', minBet: 10, maxBet: 5000 },
  { id: 'blackjack-perfect-pairs', name: 'Perfect Pairs', type: 'blackjack', icon: '♠️', description: 'Pair bets', rtp: 99.4, category: 'blackjack', provider: 'Betsoft', minBet: 10, maxBet: 5000 },
  
  // BACCARAT
  { id: 'baccarat-classic', name: 'Classic Baccarat', type: 'baccarat', icon: '🪙', description: 'Punto Banco', rtp: 98.94, category: 'baccarat', provider: 'Evolution', isHot: true, minBet: 10, maxBet: 10000 },
  { id: 'baccarat-speed', name: 'Speed Baccarat', type: 'baccarat', icon: '⚡', description: 'Lightning fast', rtp: 98.94, category: 'baccarat', provider: 'Evolution', minBet: 5, maxBet: 5000 },
  { id: 'baccarat-squeeze', name: 'Baccarat Squeeze', type: 'baccarat', icon: '🤏', description: 'Dramatic reveals', rtp: 98.94, category: 'baccarat', provider: 'Evolution', minBet: 10, maxBet: 10000 },
  { id: 'baccarat-no-comm', name: 'No Commission Baccarat', type: 'baccarat', icon: '💰', description: 'No banker fee', rtp: 98.76, category: 'baccarat', provider: 'Evolution', isNew: true, minBet: 10, maxBet: 10000 },
  { id: 'baccarat-golden', name: 'Golden Wealth Baccarat', type: 'baccarat', icon: '✨', description: 'Golden multipliers', rtp: 98.94, category: 'baccarat', provider: 'Evolution', minBet: 10, maxBet: 10000 },
  { id: 'baccarat-lightning', name: 'Lightning Baccarat', type: 'baccarat', icon: '⚡', description: 'Lightning baccarat', rtp: 98.94, category: 'baccarat', provider: 'Evolution', minBet: 10, maxBet: 10000 },
  { id: 'baccarat-first-person', name: 'First Person Baccarat', type: 'baccarat', icon: '🎮', description: '3D baccarat', rtp: 98.94, category: 'baccarat', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'baccarat-punto-banco', name: 'Punto Banco', type: 'baccarat', icon: '🎯', description: 'Classic punto banco', rtp: 98.94, category: 'baccarat', provider: 'Betsoft', minBet: 10, maxBet: 10000 },
  
  // POKER
  { id: 'poker-texas', name: 'Texas Hold\'em', type: 'poker', icon: '🃏', description: 'The classic game', rtp: 97.8, category: 'poker', provider: 'Evolution', isHot: true, minBet: 5, maxBet: 5000 },
  { id: 'poker-caribbean', name: 'Caribbean Stud', type: 'poker', icon: '🏝️', description: 'Island style poker', rtp: 96.3, category: 'poker', provider: 'Evolution', minBet: 10, maxBet: 500 },
  { id: 'poker-three-card', name: 'Three Card Poker', type: 'poker', icon: '3️⃣', description: 'Quick poker action', rtp: 96.63, category: 'poker', provider: 'Evolution', minBet: 10, maxBet: 1000 },
  { id: 'poker-ultimate', name: 'Ultimate Texas Hold\'em', type: 'poker', icon: '🏆', description: 'Ultimate poker', rtp: 97.8, category: 'poker', provider: 'Evolution', minBet: 10, maxBet: 2500 },
  { id: 'poker-casino-hold', name: 'Casino Hold\'em', type: 'poker', icon: '🎯', description: 'Side bet action', rtp: 97.8, category: 'poker', provider: 'Evolution', minBet: 5, maxBet: 500 },
  { id: 'poker-side-city', name: 'Side Bet City', type: 'poker', icon: '🏙️', description: 'City betting', rtp: 95.5, category: 'poker', provider: 'Evolution', minBet: 1, maxBet: 1000 },
  { id: 'poker-2-hand', name: '2 Hand Casino Hold\'em', type: 'poker', icon: '✌️', description: 'Two hands', rtp: 97.8, category: 'poker', provider: 'Evolution', minBet: 5, maxBet: 500 },
  { id: 'poker-teen-patti', name: 'Teen Patti', type: 'poker', icon: '🃏', description: 'Indian poker', rtp: 97.8, category: 'poker', provider: 'Evolution', minBet: 5, maxBet: 500 },
  { id: 'poker-joker', name: 'Joker Poker', type: 'poker', icon: '🤡', description: 'Joker wilds', rtp: 97.8, category: 'poker', provider: 'Betsoft', minBet: 5, maxBet: 500 },
  { id: 'poker-deuces-wild', name: 'Deuces Wild', type: 'poker', icon: '🃏', description: 'Wild deuces', rtp: 97.8, category: 'poker', provider: 'Betsoft', minBet: 5, maxBet: 500 },
  { id: 'poker-jacks-better', name: 'Jacks or Better', type: 'poker', icon: '🃏', description: 'Classic video poker', rtp: 97.8, category: 'poker', provider: 'Betsoft', minBet: 5, maxBet: 500 },
  
  // LIVE SHOWS
  { id: 'show-dream-catcher', name: 'Dream Catcher', type: 'show', icon: '🎯', description: 'Spin the wheel', rtp: 96.58, category: 'live', provider: 'Evolution', isHot: true, minBet: 0.1, maxBet: 1000 },
  { id: 'show-monopoly', name: 'Monopoly Live', type: 'show', icon: '🏠', description: 'Board game magic', rtp: 96.23, category: 'live', provider: 'Evolution', isNew: true, minBet: 0.1, maxBet: 1000 },
  { id: 'show-crazy-time', name: 'Crazy Time', type: 'show', icon: '🎪', description: 'Insane multipliers', rtp: 95.5, category: 'live', provider: 'Evolution', isHot: true, minBet: 0.1, maxBet: 10000 },
  { id: 'show-fan-tan', name: 'Fan Tan', type: 'show', icon: '🔴', description: 'Ancient Chinese', rtp: 97.5, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-super-sic-bo', name: 'Super Sic Bo', type: 'show', icon: '🎲', description: 'Dice game', rtp: 97.22, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-craps', name: 'Craps Live', type: 'show', icon: '🎲', description: 'Dice casino', rtp: 97.22, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-dragon-tiger', name: 'Dragon Tiger Live', type: 'show', icon: '🐉', description: 'East vs West', rtp: 97.0, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-andar-bahar', name: 'Andar Bahar', type: 'show', icon: '🃏', description: 'Indian classic', rtp: 97.0, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-teen-patti-live', name: 'Teen Patti Live', type: 'show', icon: '🃏', description: 'Live Indian poker', rtp: 97.0, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-color-dragon', name: 'Super Color Dragon', type: 'show', icon: '🌈', description: 'Color betting', rtp: 97.0, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-war-bets', name: 'War of Bets', type: 'show', icon: '⚔️', description: 'Battle betting', rtp: 97.0, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  { id: 'show-top-card', name: 'Top Card', type: 'show', icon: '🃏', description: 'High card wins', rtp: 97.0, category: 'live', provider: 'Evolution', minBet: 1, maxBet: 5000 },
  
  // CRASH GAMES
  { id: 'crash-aviator', name: 'Aviator', type: 'crash', icon: '✈️', description: 'Fly high', rtp: 97.0, category: 'crash', provider: 'Spribe', isHot: true, minBet: 0.1, maxBet: 100 },
  { id: 'crash-spaceman', name: 'Spaceman', type: 'crash', icon: '🚀', description: 'Space adventure', rtp: 97.0, category: 'crash', provider: 'Pragmatic Play', minBet: 0.1, maxBet: 100 },
  { id: 'crash-jetx', name: 'JetX', type: 'crash', icon: '✈️', description: 'Jet flying', rtp: 97.0, category: 'crash', provider: 'SmartSoft', minBet: 0.1, maxBet: 100 },
  { id: 'crash-balloon', name: 'Balloon', type: 'crash', icon: '🎈', description: 'Balloon ride', rtp: 97.0, category: 'crash', provider: 'SmartSoft', minBet: 0.1, maxBet: 100 },
  { id: 'crash-zeppelin', name: 'Zeppelin', type: 'crash', icon: '🎈', description: 'Airship ride', rtp: 97.0, category: 'crash', provider: 'Betsoft', minBet: 0.1, maxBet: 100 },

  // MORE SLOTS - Additional Top Games
  { id: 'slots-mega-sphere', name: 'Mega Sphere', type: 'slots', icon: '🌐', description: 'Global jackpot', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-age-gods', name: 'Age of Gods', type: 'slots', icon: '⚡', description: 'Greek mythology', rtp: 96.0, category: 'slots', provider: 'Playtech', isHot: true, minBet: 0.2, maxBet: 100 },
  { id: 'slots-gladiator', name: 'Gladiator', type: 'slots', icon: '⚔️', description: 'Roman arena', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-great-blue', name: 'Great Blue', type: 'slots', icon: '🌊', description: 'Ocean adventure', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-lottery', name: 'Lottery', type: 'slots', icon: '🎰', description: 'Lottery fun', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-pink-panther', name: 'Pink Panther', type: 'slots', icon: '🎀', description: 'Comedy detective', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-captain-america', name: 'Captain America', type: 'slots', icon: '🛡️', description: 'Superhero slot', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-avengers', name: 'Avengers', type: 'slots', icon: '🦸', description: 'Marvel heroes', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-iron-man', name: 'Iron Man', type: 'slots', icon: '🤖', description: 'Tech hero', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-x-men', name: 'X-Men', type: 'slots', icon: '🧬', description: 'Mutant heroes', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-thor', name: 'Thor', type: 'slots', icon: '🔨', description: 'Thunder god', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-hulk', name: 'Hulk', type: 'slots', icon: '💪', description: 'Green giant', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-spider-man', name: 'Spider-Man', type: 'slots', icon: '🕷️', description: 'Web slinger', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wonder-woman', name: 'Wonder Woman', type: 'slots', icon: '👸', description: 'DC hero', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-batman', name: 'Batman', type: 'slots', icon: '🦇', description: 'Dark knight', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-superman', name: 'Superman', type: 'slots', icon: '💥', description: 'Man of steel', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-jackpot-giant', name: 'Jackpot Giant', type: 'slots', icon: '🏔️', description: 'Giant jackpot', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-friendly-froger', name: 'Frog Prince', type: 'slots', icon: '🐸', description: 'Fairy tale', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-white-king', name: 'White King', type: 'slots', icon: '👑', description: 'Royal slot', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-silent-samurai', name: 'Silent Samurai', type: 'slots', icon: '🥷', description: 'Ninja slot', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-kong', name: 'Kong', type: 'slots', icon: '🦍', description: 'King Kong', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-pirates', name: 'Pirates', type: 'slots', icon: '🏴‍☠️', description: 'Treasure hunt', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-halloween', name: 'Halloween', type: 'slots', icon: '🎃', description: 'Spooky fun', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-cook-man', name: 'Cook Man', type: 'slots', icon: '👨‍🍳', description: 'Chef slot', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-sushi', name: 'Sushi', type: 'slots', icon: '🍣', description: 'Japanese food', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-fortune-keeper', name: 'Fortune Keeper', type: 'slots', icon: '🔮', description: 'Mystic fortune', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-robin-hood', name: 'Robin Hood', type: 'slots', icon: '🏹', description: 'Sherwood forest', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-lost-temple', name: 'Lost Temple', type: 'slots', icon: '🛕', description: 'Temple adventure', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-treasure-hunt', name: 'Treasure Hunt', type: 'slots', icon: '💎', description: 'Find treasure', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-jungle-joy', name: 'Jungle Joy', type: 'slots', icon: '🦁', description: 'Jungle adventure', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-wild-spirit', name: 'Wild Spirit', type: 'slots', icon: '🐺', description: 'Wild wolves', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-golden-ark', name: 'Golden Ark', type: 'slots', icon: '方舟', description: 'Noah\'s ark', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-graceful-china', name: 'Zhao Cai Tong Zi', type: 'slots', icon: '🐉', description: 'Chinese fortune', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-aztec-temple', name: 'Aztec Temple', type: 'slots', icon: '🗿', description: 'Aztec adventure', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-cat-queen', name: 'Cat Queen', type: 'slots', icon: '🐱', description: 'Cat goddess', rtp: 96.0, category: 'slots', provider: 'Playtech', minBet: 0.2, maxBet: 100 },
  { id: 'slots-sticky-diamonds', name: 'Sticky Diamonds', type: 'slots', icon: '💎', description: 'Diamond rush', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-tipsy-tourist', name: 'Tipsy Tourist', type: 'slots', icon: '🍹', description: 'Beach fun', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-tiger-claw', name: 'Tiger Claw', type: 'slots', icon: '🐯', description: 'Wild tiger', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-under-water', name: 'Under Water', type: 'slots', icon: '🐠', description: 'Ocean deep', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-birds', name: 'Birds', type: 'slots', icon: '🐦', description: 'Bird slot', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-tekken', name: 'Tekken', type: 'slots', icon: '🥊', description: 'Fighting game', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-rough-rider', name: 'Rough Rider', type: 'slots', icon: '🏇', description: 'Wild west', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-mad-pinocchio', name: 'Mad Pinocchio', type: 'slots', icon: '🤥', description: 'Wooden boy', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-gold-diggers', name: 'Gold Diggers', type: 'slots', icon: '⛏️', description: 'Gold mining', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-reel-of-fortune', name: 'Reel of Fortune', type: 'slots', icon: '🎰', description: 'Lucky reels', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
  { id: 'slots-weekend-in-vegas', name: 'Weekend in Vegas', type: 'slots', icon: '🎰', description: 'Vegas fun', rtp: 96.0, category: 'slots', provider: 'Betsoft', minBet: 0.2, maxBet: 100 },
];

const categories = [
  { id: 'all', name: 'All Games', icon: '🎮' },
  { id: 'slots', name: 'Slots', icon: '🎰' },
  { id: 'dice', name: 'Dice', icon: '🎲' },
  { id: 'roulette', name: 'Roulette', icon: '🎡' },
  { id: 'blackjack', name: 'Blackjack', icon: '🃏' },
  { id: 'baccarat', name: 'Baccarat', icon: '🪙' },
  { id: 'poker', name: 'Poker', icon: '♠️' },
  { id: 'live', name: 'Live Shows', icon: '🎬' },
];

export default function GamesPage() {
  const [activeCategory, setActiveCategory] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'rtp' | 'provider'>('name');

  const filteredGames = games
    .filter(game => 
      (activeCategory === 'all' || game.category === activeCategory) &&
      (game.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
       game.provider.toLowerCase().includes(searchQuery.toLowerCase()))
    )
    .sort((a, b) => {
      if (sortBy === 'name') return a.name.localeCompare(b.name);
      if (sortBy === 'rtp') return b.rtp - a.rtp;
      return a.provider.localeCompare(b.provider);
    });

  return (
    <>
      <Header />
      <main className={styles.main}>
        <section className={styles.hero}>
          <div className={styles.heroContent}>
            <h1 className={styles.title}>
              <span className={styles.highlight}>200+</span> Casino Games
            </h1>
            <p className={styles.subtitle}>
              Experience the thrill of real casino games with instant crypto payouts
            </p>
          </div>
        </section>

        <section className={styles.controls}>
          <div className={styles.searchBar}>
            <span className={styles.searchIcon}>🔍</span>
            <input
              type="text"
              placeholder="Search games or providers..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className={styles.searchInput}
            />
          </div>

          <div className={styles.sortWrapper}>
            <label>Sort by:</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as 'name' | 'rtp' | 'provider')}
              className={styles.sortSelect}
            >
              <option value="name">Name</option>
              <option value="rtp">RTP</option>
              <option value="provider">Provider</option>
            </select>
          </div>
        </section>

        <section className={styles.categories}>
          {categories.map(cat => (
            <button
              key={cat.id}
              className={`${styles.categoryBtn} ${activeCategory === cat.id ? styles.active : ''}`}
              onClick={() => setActiveCategory(cat.id)}
            >
              <span className={styles.categoryIcon}>{cat.icon}</span>
              <span className={styles.categoryName}>{cat.name}</span>
            </button>
          ))}
        </section>

        <section className={styles.gamesGrid}>
          {filteredGames.map((game, index) => (
            <Link href={`/games/${game.type}/${game.id}`} key={game.id}>
              <Card variant="glow" padding="none" className={`${styles.gameCard} stagger-${(index % 5) + 1}`}>
                <div className={styles.gameIcon}>
                  {game.icon}
                  {game.isHot && <span className={styles.hotBadge}>🔥 HOT</span>}
                  {game.isNew && <span className={styles.newBadge}>🆕 NEW</span>}
                </div>
                <div className={styles.gameInfo}>
                  <h3 className={styles.gameName}>{game.name}</h3>
                  <p className={styles.gameDescription}>{game.description}</p>
                  <div className={styles.gameMeta}>
                    <span className={styles.provider}>{game.provider}</span>
                    <span className={styles.rtp}>RTP: {game.rtp}%</span>
                  </div>
                </div>
                <div className={styles.playBtn}>
                  <Button variant="primary" size="sm">Play Now</Button>
                </div>
              </Card>
            </Link>
          ))}
        </section>

        {filteredGames.length === 0 && (
          <div className={styles.noResults}>
            <span className={styles.noResultsIcon}>🎮</span>
            <h3>No games found</h3>
            <p>Try adjusting your search or filters</p>
          </div>
        )}
      </main>
      <Footer />
    </>
  );
}
