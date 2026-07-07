'use client';

import React, { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import toast from 'react-hot-toast';

// Comprehensive slot games data from multiple providers
const SLOT_GAMES = [
  // Pragmatic Play
  { id: 'pp_gates_of_olympus', name: 'Gates of Olympus', provider: 'Pragmatic Play', icon: '⚡', rtp: 96.50, volatility: 'high', features: ['Tumble', 'Multiplier', 'Free Spins'], isHot: true, category: 'slots' },
  { id: 'pp_sweet_bonanza', name: 'Sweet Bonanza', provider: 'Pragmatic Play', icon: '🍬', rtp: 96.48, volatility: 'medium', features: ['Tumble', 'Multiplier'], isHot: true, category: 'slots' },
  { id: 'pp_starlight_princess', name: 'Starlight Princess', provider: 'Pragmatic Play', icon: '👸', rtp: 96.50, volatility: 'high', features: ['Tumble', 'Multiplier', 'Free Spins'], isNew: true, category: 'slots' },
  { id: 'pp_big_bass_bonanza', name: 'Big Bass Bonanza', provider: 'Pragmatic Play', icon: '🎣', rtp: 96.71, volatility: 'medium', features: ['Free Spins', 'Respin'], isHot: true, category: 'slots' },
  { id: 'pp_wolf_gold', name: 'Wolf Gold', provider: 'Pragmatic Play', icon: '🐺', rtp: 96.00, volatility: 'medium', features: ['Blazing Reels', 'Money Respin'], isHot: true, category: 'slots' },
  { id: 'pp_fruit_party', name: 'Fruit Party', provider: 'Pragmatic Play', icon: '🍓', rtp: 96.47, volatility: 'high', features: ['Tumble', 'Multiplier'], isHot: true, category: 'slots' },
  { id: 'pp_the_dog_house', name: 'The Dog House', provider: 'Pragmatic Play', icon: '🐕', rtp: 96.51, volatility: 'high', features: ['Sticky Wild', 'Free Spins'], isHot: true, category: 'slots' },
  { id: 'pp_wild_west_gold', name: 'Wild West Gold', provider: 'Pragmatic Play', icon: '🤠', rtp: 96.51, volatility: 'high', features: ['Sticky Wild', 'Free Spins'], category: 'slots' },
  { id: 'pp_book_of_dead', name: 'Book of Dead', provider: 'Play\'n GO', icon: '📖', rtp: 96.21, volatility: 'high', features: ['Expandable Wild', 'Free Spins'], isHot: true, category: 'slots' },
  { id: 'pp_legacy_of_dead', name: 'Legacy of Dead', provider: 'Play\'n GO', icon: '🏛️', rtp: 96.50, volatility: 'high', features: ['Expandable Symbol', 'Free Spins'], isHot: true, category: 'slots' },
  
  // NetEnt
  { id: 'net_starburst', name: 'Starburst', provider: 'NetEnt', icon: '💎', rtp: 96.09, volatility: 'low', features: ['Expanding Wild', 'Re-spins'], isHot: true, category: 'slots' },
  { id: 'net_gonzo_quest', name: 'Gonzo\'s Quest', provider: 'NetEnt', icon: '🌋', rtp: 96.00, volatility: 'high', features: ['Avalanche', 'Free Fall', 'Multiplier'], isHot: true, category: 'slots' },
  { id: 'net_dead_or_alive', name: 'Dead or Alive', provider: 'NetEnt', icon: '🔫', rtp: 96.80, volatility: 'high', features: ['Free Spins', 'Sticky Wild'], category: 'slots' },
  { id: 'net_twin_spin', name: 'Twin Spin', provider: 'NetEnt', icon: '🎰', rtp: 96.60, volatility: 'medium', features: ['Twin Reels', 'Wild'], category: 'slots' },
  { id: 'net_jack_and_beanstalk', name: 'Jack and the Beanstalk', provider: 'NetEnt', icon: '🫘', rtp: 96.30, volatility: 'high', features: ['Walking Wild', 'Free Spins'], category: 'slots' },
  { id: 'net_blood_suckers', name: 'Blood Suckers', provider: 'NetEnt', icon: '🧛', rtp: 98.00, volatility: 'low', features: ['Free Spins', 'Bonus Game'], category: 'slots' },
  
  // Microgaming
  { id: 'mg_mega_moolah', name: 'Mega Moolah', provider: 'Microgaming', icon: '🦁', rtp: 88.12, volatility: 'high', features: ['Progressive Jackpot', 'Free Spins'], isHot: true, category: 'jackpot' },
  { id: 'mg_immortal_romance', name: 'Immortal Romance', provider: 'Microgaming', icon: '🧛', rtp: 96.86, volatility: 'high', features: ['Free Spins', 'Wild'], category: 'slots' },
  { id: 'mg_thunderstruck_2', name: 'Thunderstruck II', provider: 'Microgaming', icon: '⚡', rtp: 96.65, volatility: 'high', features: ['Free Spins', 'Wild', 'Multipliers'], category: 'slots' },
  { id: 'mg_game_of_thrones', name: 'Game of Thrones', provider: 'Microgaming', icon: '👑', rtp: 96.50, volatility: 'high', features: ['Free Spins', 'Stacked Wild'], category: 'slots' },
  
  // BGaming
  { id: 'bg_aztec_magic', name: 'Aztec Magic', provider: 'BGaming', icon: '🗿', rtp: 96.96, volatility: 'medium', features: ['Free Spins', 'Bonus Game'], category: 'slots' },
  { id: 'bg_aztec_magic_megaways', name: 'Aztec Magic Megaways', provider: 'BGaming', icon: '🗿', rtp: 96.70, volatility: 'high', features: ['Megaways', 'Cascade', 'Free Spins'], category: 'megaways' },
  { id: 'bg_egypt_horus', name: 'Elvis Frog in Vegas', provider: 'BGaming', icon: '🐸', rtp: 96.00, volatility: 'high', features: ['Free Spins', 'Multiplier'], category: 'slots' },
  
  // Spribe (Crash Games)
  { id: 'sp_aviator', name: 'Aviator', provider: 'Spribe', icon: '✈️', rtp: 97.00, volatility: 'medium', features: ['Auto Cashout', 'Multiplayer'], isHot: true, category: 'crash' },
  { id: 'sp_mines', name: 'Mines', provider: 'Spribe', icon: '💣', rtp: 97.00, volatility: 'medium', features: ['Customizable Mines', 'Auto'], category: 'mine' },
  { id: 'sp_dice', name: 'Dice', provider: 'Spribe', icon: '🎲', rtp: 97.00, volatility: 'medium', features: ['Slider', 'Auto'], category: 'dice' },
  { id: 'sp_plinko', name: 'Plinko', provider: 'Spribe', icon: '🔮', rtp: 98.00, volatility: 'medium', features: ['Risk Level', 'Auto'], category: 'plinko' },
  
  // Hacksaw Gaming
  { id: 'hack_wanted_dead', name: 'Wanted Dead or a Wild', provider: 'Hacksaw Gaming', icon: '🔫', rtp: 96.24, volatility: 'high', features: ['Duel', 'Free Spins', 'RTG'], category: 'slots' },
  { id: 'hack_sticky_bands', name: 'Sticky Bandits', provider: 'Hacksaw Gaming', icon: '💰', rtp: 96.20, volatility: 'high', features: ['Sticky Wild', 'Free Spins'], category: 'slots' },
  
  // Yggdrasil
  { id: 'yg_valley_of_gods', name: 'Valley of the Gods', provider: 'Yggdrasil', icon: '🏺', rtp: 96.10, volatility: 'high', features: ['Extinctions', 'Respins', 'Win Multiplier'], category: 'slots' },
  { id: 'yg_holmes', name: 'Holmes and the Stolen Stones', provider: 'Yggdrasil', icon: '🕵️', rtp: 96.80, volatility: 'medium', features: ['Jackpot', 'Free Spins', 'Wild'], category: 'jackpot' },
  
  // Big Time Gaming (Megaways)
  { id: 'btg_white_rabbit', name: 'White Rabbit', provider: 'Big Time Gaming', icon: '🐇', rtp: 97.70, volatility: 'high', features: ['Megaways', 'Cascade', 'Free Spins'], category: 'megaways' },
  { id: 'btg_extra_chilli', name: 'Extra Chilli', provider: 'Big Time Gaming', icon: '🌶️', rtp: 96.80, volatility: 'high', features: ['Megaways', 'Free Spins', 'Feature Drop'], category: 'megaways' },
  
  // PG Soft
  { id: 'pg_fortune_ox', name: 'Fortune Ox', provider: 'PG Soft', icon: '牛', rtp: 96.50, volatility: 'medium', features: ['Multiplier', 'Free Spins'], category: 'slots' },
  { id: 'pg_fortune_mouse', name: 'Fortune Mouse', provider: 'PG Soft', icon: '🐭', rtp: 96.50, volatility: 'medium', features: ['Respin', 'Multiplier'], category: 'slots' },
  
  // Betsoft
  { id: 'bs_gold_digger', name: 'Gold Digger', provider: 'Betsoft', icon: '⛏️', rtp: 96.00, volatility: 'high', features: ['Link & Win', 'Free Spins', 'Respin'], category: 'slots' },
  { id: 'bs_atlantis', name: 'Quest to the West', provider: 'Betsoft', icon: '🌊', rtp: 97.50, volatility: 'high', features: ['Walking Wild', 'Multiplier'], category: 'slots' },
  
  // Evoplay
  { id: 'evo_fruit_cocktail', name: 'Fruit Cocktail', provider: 'Evoplay', icon: '🍹', rtp: 96.00, volatility: 'medium', features: ['Bonus', 'Free Spins'], category: 'slots' },
  { id: 'evo_garage', name: 'Garage', provider: 'Evoplay', icon: '🚗', rtp: 96.00, volatility: 'medium', features: ['Bonus', 'Respin'], category: 'slots' },
  { id: 'evo_crazy_monkey', name: 'Crazy Monkey', provider: 'Evoplay', icon: '🐒', rtp: 96.00, volatility: 'medium', features: ['Bonus', 'Gamble'], category: 'slots' },
  
  // More Pragmatic Play
  { id: 'pp_john_hunter', name: 'John Hunter Aztec Treasure', provider: 'Pragmatic Play', icon: '🎯', rtp: 96.00, volatility: 'high', features: ['Tumble', 'Free Spins', 'Multipliers'], category: 'slots' },
  { id: 'pp_great_rhino', name: 'Great Rhino', provider: 'Pragmatic Play', icon: '🦏', rtp: 96.50, volatility: 'medium', features: ['Super Respin', 'Free Spins'], category: 'slots' },
  { id: 'pp_mustang_gold', name: 'Mustang Gold', provider: 'Pragmatic Play', icon: '🐴', rtp: 96.53, volatility: 'high', features: ['Collect', 'Free Spins', 'JACKPOT'], category: 'slots' },
  { id: 'pp_aztec_gems', name: 'Aztec Gems', provider: 'Pragmatic Play', icon: '💎', rtp: 96.52, volatility: 'high', features: ['Multiplier', 'Free Spins'], category: 'slots' },
  { id: 'pp_five_lions_megaways', name: 'Five Lions Megaways', provider: 'Pragmatic Play', icon: '🦁', rtp: 96.50, volatility: 'high', features: ['Megaways', 'Cascading', 'Multiplier'], category: 'megaways' },
  { id: 'pp_power_of_thor', name: 'Power of Thor Megaways', provider: 'Pragmatic Play', icon: '🔨', rtp: 96.50, volatility: 'high', features: ['Megaways', 'Free Spins', 'Multiplier'], category: 'megaways' },
  { id: 'pp_buffalo_king', name: 'Buffalo King', provider: 'Pragmatic Play', icon: '🦬', rtp: 96.50, volatility: 'high', features: ['Free Spins', 'Tumble', 'Multiplier'], category: 'slots' },
  { id: 'pp_fruit_party_2', name: 'Fruit Party 2', provider: 'Pragmatic Play', icon: '🍒', rtp: 96.50, volatility: 'high', features: ['Tumble', 'Multiplier', 'Free Spins'], isNew: true, category: 'slots' },
  { id: 'pp_christmas_carol', name: 'Christmas Carol Megaways', provider: 'Pragmatic Play', icon: '🎄', rtp: 96.50, volatility: 'high', features: ['Megaways', 'Free Spins', 'Multiplier'], category: 'megaways' },
  { id: 'pp_magic_gems', name: 'Magic Gems', provider: 'Pragmatic Play', icon: '💍', rtp: 96.50, volatility: 'high', features: ['Tumble', 'Multiplier', 'Free Spins'], category: 'slots' },
  { id: 'pp_treasure_wild', name: 'Treasure Wild', provider: 'Pragmatic Play', icon: '💰', rtp: 96.50, volatility: 'high', features: ['Tumble', 'Free Spins', 'Wild'], category: 'slots' },
];

const PROVIDERS = [
  { id: 'all', name: 'All Providers', count: SLOT_GAMES.length },
  { id: 'Pragmatic Play', name: 'Pragmatic Play', count: SLOT_GAMES.filter(g => g.provider === 'Pragmatic Play').length },
  { id: 'NetEnt', name: 'NetEnt', count: SLOT_GAMES.filter(g => g.provider === 'NetEnt').length },
  { id: 'Microgaming', name: 'Microgaming', count: SLOT_GAMES.filter(g => g.provider === 'Microgaming').length },
  { id: 'BGaming', name: 'BGaming', count: SLOT_GAMES.filter(g => g.provider === 'BGaming').length },
  { id: 'Spribe', name: 'Spribe', count: SLOT_GAMES.filter(g => g.provider === 'Spribe').length },
  { id: 'Hacksaw Gaming', name: 'Hacksaw Gaming', count: SLOT_GAMES.filter(g => g.provider === 'Hacksaw Gaming').length },
  { id: 'Yggdrasil', name: 'Yggdrasil', count: SLOT_GAMES.filter(g => g.provider === 'Yggdrasil').length },
  { id: 'Big Time Gaming', name: 'Big Time Gaming', count: SLOT_GAMES.filter(g => g.provider === 'Big Time Gaming').length },
  { id: 'PG Soft', name: 'PG Soft', count: SLOT_GAMES.filter(g => g.provider === 'PG Soft').length },
  { id: 'Betsoft', name: 'Betsoft', count: SLOT_GAMES.filter(g => g.provider === 'Betsoft').length },
  { id: 'Evoplay', name: 'Evoplay', count: SLOT_GAMES.filter(g => g.provider === 'Evoplay').length },
];

const CATEGORIES = [
  { id: 'all', name: 'All', icon: '🎰' },
  { id: 'slots', name: 'Slots', icon: '🎰' },
  { id: 'megaways', name: 'Megaways', icon: '🌀' },
  { id: 'jackpot', name: 'Jackpots', icon: '💰' },
  { id: 'crash', name: 'Crash', icon: '✈️' },
  { id: 'mine', name: 'Mines', icon: '💣' },
  { id: 'dice', name: 'Dice', icon: '🎲' },
  { id: 'plinko', name: 'Plinko', icon: '🔮' },
];

export default function SlotsPage() {
  const [activeProvider, setActiveProvider] = useState('all');
  const [activeCategory, setActiveCategory] = useState('all');
  const [search, setSearch] = useState('');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [sortBy, setSortBy] = useState<'name' | 'rtp' | 'provider'>('name');
  const [showFavorites, setShowFavorites] = useState(false);

  const filteredGames = useMemo(() => {
    let games = SLOT_GAMES;
    
    if (activeProvider !== 'all') {
      games = games.filter(g => g.provider === activeProvider);
    }
    
    if (activeCategory !== 'all') {
      games = games.filter(g => g.category === activeCategory);
    }
    
    if (search) {
      games = games.filter(g => 
        g.name.toLowerCase().includes(search.toLowerCase()) ||
        g.provider.toLowerCase().includes(search.toLowerCase())
      );
    }
    
    // Sort
    games = [...games].sort((a, b) => {
      if (sortBy === 'name') return a.name.localeCompare(b.name);
      if (sortBy === 'rtp') return b.rtp - a.rtp;
      if (sortBy === 'provider') return a.provider.localeCompare(b.provider);
      return 0;
    });
    
    // Hot and New first
    games = [...games].sort((a, b) => {
      if (a.isHot && !b.isHot) return -1;
      if (!a.isHot && b.isHot) return 1;
      if (a.isNew && !b.isNew) return -1;
      if (!a.isNew && b.isNew) return 1;
      return 0;
    });
    
    return games;
  }, [activeProvider, activeCategory, search, sortBy, showFavorites]);

  const handleGameClick = (gameId: string) => {
    toast.loading(`Launching ${gameId}...`, { duration: 1000 });

    // Call real API
    fetch(`/api/games/slots/bet`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        amount: 1.0 // Default bet
      })
    })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        toast.error(data.error);
        return;
      }

      if (data.win_amount > 0) {
        toast.success(`JACKPOT! You won $${data.win_amount.toFixed(2)} on ${gameId}!`, {
          icon: '💰',
          duration: 5000
        });
      } else {
        toast.error(`No win this time on ${gameId}. Try again!`);
      }
    })
    .catch(err => {
      toast.error('Failed to connect to game server');
    });
  };

  return (
    <div className="min-h-screen bg-tiger-dark p-6">
      <div className="max-w-7xl mx-auto">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-4xl font-heading text-gradient mb-2">🎰 Slot Games</h1>
            <p className="text-gray-400">{filteredGames.length} games available</p>
          </div>
          
          {/* Search */}
          <div className="relative">
            <input
              type="text"
              placeholder="Search games..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="bg-tiger-surface border border-tiger-border rounded-lg px-4 py-2 w-64 text-white"
            />
          </div>
        </div>

        {/* Category Filter */}
        <div className="flex gap-2 mb-4 overflow-x-auto pb-2">
          {CATEGORIES.map(cat => (
            <button
              key={cat.id}
              onClick={() => setActiveCategory(cat.id)}
              className={`px-4 py-2 rounded-lg whitespace-nowrap ${
                activeCategory === cat.id 
                  ? 'bg-primary-500 text-white' 
                  : 'bg-tiger-surface text-gray-400 hover:bg-tiger-border'
              }`}
            >
              {cat.icon} {cat.name}
            </button>
          ))}
        </div>

        {/* Provider Filter */}
        <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
          {PROVIDERS.map(provider => (
            <button
              key={provider.id}
              onClick={() => setActiveProvider(provider.id)}
              className={`px-3 py-1.5 rounded-full text-sm whitespace-nowrap ${
                activeProvider === provider.id 
                  ? 'bg-primary-500 text-white' 
                  : 'bg-tiger-surface text-gray-400 hover:bg-tiger-border'
              }`}
            >
              {provider.name} ({provider.count})
            </button>
          ))}
        </div>

        {/* Sort and View Controls */}
        <div className="flex justify-between items-center mb-4">
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as any)}
            className="bg-tiger-surface border border-tiger-border rounded-lg px-3 py-2 text-white"
          >
            <option value="name">Sort by Name</option>
            <option value="rtp">Sort by RTP</option>
            <option value="provider">Sort by Provider</option>
          </select>
          
          <div className="flex gap-2">
            <button
              onClick={() => setShowFavorites(!showFavorites)}
              className={`px-4 py-2 rounded-lg ${showFavorites ? 'bg-yellow-500 text-black' : 'bg-tiger-surface text-gray-400'}`}
            >
              ⭐ Favorites
            </button>
            <button
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className="px-4 py-2 bg-tiger-surface rounded-lg text-gray-400"
            >
              {viewMode === 'grid' ? '📋 List' : '🔲 Grid'}
            </button>
          </div>
        </div>

        {/* Games Grid */}
        <div className={viewMode === 'grid' 
          ? "grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4"
          : "flex flex-col gap-2"
        }>
          {filteredGames.map(game => (
            <motion.div 
              key={game.id} 
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => handleGameClick(game.id)}
              className={`glass rounded-xl overflow-hidden cursor-pointer ${
                viewMode === 'list' ? 'flex items-center p-3' : ''
              }`}
            >
              <div className={`relative ${viewMode === 'grid' ? 'aspect-square' : 'w-20 h-20'} bg-tiger-surface flex items-center justify-center text-5xl`}>
                {game.icon}
                {(game.isHot || game.isNew) && (
                  <div className="absolute top-2 right-2">
                    {game.isHot && <span className="bg-red-500 text-white text-xs px-2 py-1 rounded-full">🔥 HOT</span>}
                    {game.isNew && <span className="bg-green-500 text-white text-xs px-2 py-1 rounded-full ml-1">🆕 NEW</span>}
                  </div>
                )}
              </div>
              <div className={`${viewMode === 'grid' ? 'p-3' : 'ml-4 flex-1'}`}>
                <h3 className="font-bold text-white truncate">{game.name}</h3>
                <p className="text-gray-400 text-sm">{game.provider}</p>
                <div className="flex items-center gap-2 mt-1">
                  <span className="text-green-400 text-sm">RTP: {game.rtp}%</span>
                  <span className={`text-xs px-2 py-0.5 rounded ${
                    game.volatility === 'high' ? 'bg-red-500/20 text-red-400' :
                    game.volatility === 'medium' ? 'bg-yellow-500/20 text-yellow-400' :
                    'bg-green-500/20 text-green-400'
                  }`}>
                    {game.volatility}
                  </span>
                </div>
                {viewMode === 'list' && (
                  <p className="text-gray-500 text-sm mt-1 truncate">
                    {game.features.join(' • ')}
                  </p>
                )}
              </div>
            </motion.div>
          ))}
        </div>

        {filteredGames.length === 0 && (
          <div className="text-center py-20">
            <p className="text-6xl mb-4">🎮</p>
            <p className="text-gray-400 text-xl">No games found</p>
            <p className="text-gray-500">Try adjusting your filters</p>
          </div>
        )}
      </div>
    </div>
  );
}
