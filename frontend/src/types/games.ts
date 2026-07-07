// Game-specific types for all casino games

// Dice Game
export interface DiceGameState {
  betAmount: number;
  targetNumber: number;
  rollOver: boolean; // true = over, false = under
  result?: number;
  winAmount: number;
  multiplier: number;
}

// Crash Game
export interface CrashGameState {
  betAmount: number;
  autoCashout?: number;
  currentMultiplier: number;
  isPlaying: boolean;
  hasCashedOut: boolean;
  history: number[];
  crashPoint: number;
}

// Mines Game
export interface MinesGameState {
  betAmount: number;
  minesCount: number;
  revealedCells: number[];
  currentWin: number;
  isGameOver: boolean;
  isWin: boolean;
  safeCells: number;
}

// Plinko Game
export interface PlinkoGameState {
  betAmount: number;
  rows: number;
  risk: 'low' | 'medium' | 'high';
  balls: number;
  currentBall?: number;
  results: number[];
  totalWin: number;
}

// Slots
export interface SlotSpinResult {
  reels: string[][];
  paylines: PaylineResult[];
  totalWin: number;
  freeSpins: number;
  bonusTriggered: boolean;
}

export interface PaylineResult {
  payline: number[];
  symbol: string;
  count: number;
  win: number;
}

// Roulette
export interface RouletteGameState {
  bets: RouletteBet[];
  wheelType: 'european' | 'american';
  currentNumber: number | null;
  isSpinning: boolean;
  recentHistory: number[];
}

export interface RouletteBet {
  type: RouletteBetType;
  numbers: number[];
  amount: number;
  payout: number;
}

export type RouletteBetType = 
  | 'straight'
  | 'split'
  | 'street'
  | 'corner'
  | 'line'
  | 'column'
  | 'dozen'
  | 'red'
  | 'black'
  | 'even'
  | 'odd'
  | '1-18'
  | '19-36';

// Blackjack
export interface BlackjackGameState {
  deck: string[];
  playerHand: string[];
  dealerHand: string[];
  playerScore: number;
  dealerScore: number;
  betAmount: number;
  gamePhase: 'betting' | 'playing' | 'dealerTurn' | 'finished';
  result?: 'win' | 'lose' | 'push' | 'blackjack' | 'bust';
  insurance?: boolean;
  canSplit: boolean;
  canDouble: boolean;
}

// Baccarat
export interface BaccaratGameState {
  playerHand: string[];
  bankerHand: string[];
  playerScore: number;
  bankerScore: number;
  betAmount: number;
  betType: 'player' | 'banker' | 'tie' | 'pair';
  result?: 'player' | 'banker' | 'tie' | 'playerPair' | 'bankerPair';
  isPlayerDone: boolean;
  isBankerDone: boolean;
}

// Poker
export interface PokerGameState {
  hand: string[];
  betAmount: number;
  heldCards: boolean[];
  drawCount: number;
  finalHand: PokerHandType;
  winAmount: number;
  payTable: Record<PokerHandType, number>;
}

export type PokerHandType = 
  | 'royalFlush'
  | 'straightFlush'
  | 'fourOfAKind'
  | 'fullHouse'
  | 'flush'
  | 'straight'
  | 'threeOfAKind'
  | 'twoPair'
  | 'jacksOrBetter'
  | 'nothing';

// Sports Betting
export interface SportsBet {
  id: string;
  userId: string;
  selections: BetSelection[];
  stake: number;
  potentialWin: number;
  status: 'pending' | 'won' | 'lost' | 'cancelled';
  settledAt?: string;
  createdAt: string;
}

export interface BetSelection {
  eventId: string;
  market: string;
  odds: number;
  selection: string;
  result?: 'won' | 'lost' | 'pending';
}

export interface SportsEvent {
  id: string;
  sport: string;
  league: string;
  homeTeam: string;
  awayTeam: string;
  startTime: string;
  markets: BetMarket[];
  status: 'upcoming' | 'live' | 'finished';
}

export interface BetMarket {
  id: string;
  name: string;
  selections: MarketSelection[];
}

export interface MarketSelection {
  id: string;
  name: string;
  odds: number;
}

// Live Casino
export interface LiveCasinoTable {
  id: string;
  provider: string;
  game: string;
  name: string;
  minBet: number;
  maxBet: number;
  players: number;
  isLive: boolean;
  thumbnail: string;
}

// Tournament
export interface Tournament {
  id: string;
  name: string;
  game: string;
  startTime: string;
  endTime: string;
  prizePool: number;
  entryFee: number;
  maxParticipants: number;
  currentParticipants: number;
  status: 'upcoming' | 'active' | 'completed';
  leaderboard: TournamentLeaderboardEntry[];
}

export interface TournamentLeaderboardEntry {
  rank: number;
  userId: string;
  username: string;
  score: number;
  prize: number;
}

// VIP System
export interface VIPLevel {
  level: number;
  name: string;
  rakeback: number;
  depositBonus: number;
  withdrawalLimit: number;
  exclusiveGames: boolean;
  vipManager: boolean;
}

// Leaderboard
export interface LeaderboardEntry {
  rank: number;
  userId: string;
  username: string;
  avatar?: string;
  amount: number;
  wins: number;
  biggestWin: number;
}

// Provider
export interface GameProvider {
  id: string;
  name: string;
  logo: string;
  gameCount: number;
  isLive: boolean;
  categories: string[];
}

// Jackpot
export interface Jackpot {
  id: string;
  name: string;
  game: string;
  currentAmount: number;
  minBet: number;
  lastWon: string;
  winner?: string;
}
