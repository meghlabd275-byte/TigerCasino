// Game card component for game listings
import React, { useState } from 'react';
import styles from './GameCard.module.css';

interface GameCardProps {
  id: string;
  name: string;
  provider: string;
  thumbnail: string;
  rtp: number;
  categories: string[];
  isLive?: boolean;
  isFavorite?: boolean;
  onPlay: (gameId: string) => void;
  onDemo?: (gameId: string) => void;
  onToggleFavorite?: (gameId: string) => void;
}

export const GameCard: React.FC<GameCardProps> = ({
  id,
  name,
  provider,
  thumbnail,
  rtp,
  categories,
  isLive = false,
  isFavorite = false,
  onPlay,
  onDemo,
  onToggleFavorite,
}) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className={`${styles.card} ${isHovered ? styles.hovered : ''}`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <div className={styles.thumbnail}>
        <img src={thumbnail} alt={name} />
        {isLive && <span className={styles.liveBadge}>LIVE</span>}
        {isFavorite && (
          <button 
            className={styles.favoriteBtn} 
            onClick={(e) => {
              e.stopPropagation();
              onToggleFavorite?.(id);
            }}
          >
            ★
          </button>
        )}
        <div className={styles.overlay}>
          <button className={styles.playBtn} onClick={() => onPlay(id)}>
            ▶ Play Now
          </button>
          {onDemo && (
            <button className={styles.demoBtn} onClick={() => onDemo(id)}>
              Demo
            </button>
          )}
        </div>
      </div>
      <div className={styles.info}>
        <h3 className={styles.name}>{name}</h3>
        <p className={styles.provider}>{provider}</p>
        <div className={styles.meta}>
          <span className={styles.rtp}>RTP: {rtp}%</span>
          <div className={styles.categories}>
            {categories.slice(0, 2).map((cat) => (
              <span key={cat} className={styles.category}>{cat}</span>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

// Game grid component
interface GameGridProps {
  games: GameCardProps[];
  loading?: boolean;
  columns?: number;
}

export const GameGrid: React.FC<GameGridProps> = ({ 
  games, 
  loading = false,
  columns = 4 
}) => {
  if (loading) {
    return (
      <div className={styles.grid} style={{ '--columns': columns } as React.CSSProperties}>
        {[...Array(8)].map((_, i) => (
          <div key={i} className={styles.skeleton} />
        ))}
      </div>
    );
  }

  return (
    <div className={styles.grid} style={{ '--columns': columns } as React.CSSProperties}>
      {games.map((game) => (
        <GameCard key={game.id} {...game} />
      ))}
    </div>
  );
};

// Category filter component
interface Category {
  id: string;
  name: string;
  icon?: string;
  count?: number;
}

interface CategoryFilterProps {
  categories: Category[];
  activeCategory: string;
  onSelectCategory: (categoryId: string) => void;
}

export const CategoryFilter: React.FC<CategoryFilterProps> = ({
  categories,
  activeCategory,
  onSelectCategory,
}) => {
  return (
    <div className={styles.categoryFilter}>
      <button
        className={`${styles.categoryBtn} ${activeCategory === 'all' ? styles.active : ''}`}
        onClick={() => onSelectCategory('all')}
      >
        All Games
      </button>
      {categories.map((cat) => (
        <button
          key={cat.id}
          className={`${styles.categoryBtn} ${activeCategory === cat.id ? styles.active : ''}`}
          onClick={() => onSelectCategory(cat.id)}
        >
          {cat.icon && <span className={styles.catIcon}>{cat.icon}</span>}
          {cat.name}
          {cat.count !== undefined && <span className={styles.count}>({cat.count})</span>}
        </button>
      ))}
    </div>
  );
};

// Search component
interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

export const SearchBar: React.FC<SearchBarProps> = ({
  value,
  onChange,
  placeholder = 'Search games...',
}) => {
  return (
    <div className={styles.searchBar}>
      <span className={styles.searchIcon}>🔍</span>
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className={styles.searchInput}
      />
      {value && (
        <button className={styles.clearBtn} onClick={() => onChange('')}>
          ×
        </button>
      )}
    </div>
  );
};

// Sort component
interface SortOption {
  value: string;
  label: string;
}

interface GameSortProps {
  options: SortOption[];
  value: string;
  onChange: (value: string) => void;
}

export const GameSort: React.FC<GameSortProps> = ({ options, value, onChange }) => {
  return (
    <select 
      className={styles.sortSelect} 
      value={value} 
      onChange={(e) => onChange(e.target.value)}
    >
      {options.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.label}
        </option>
      ))}
    </select>
  );
};
