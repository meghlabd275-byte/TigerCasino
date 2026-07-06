// User types
export interface User {
  id: string;
  email: string;
  username: string;
  walletAddress?: string;
  walletType?: string;
  balance: number;
  bonusBalance: number;
  vipLevel: number;
  kycStatus: 'pending' | 'verified' | 'rejected';
  isVerified: boolean;
  isAdmin: boolean;
  isBanned: boolean;
  banReason?: string;
  createdAt: string;
  updatedAt: string;
  lastLogin?: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  token: string | null;
}

// Transaction types
export interface Transaction {
  id: string;
  userId: string;
  type: 'deposit' | 'withdrawal' | 'bet' | 'win' | 'bonus';
  amount: number;
  currency: string;
  status: 'pending' | 'completed' | 'rejected';
  txHash?: string;
  address?: string;
  fee?: number;
  createdAt: string;
  processedAt?: string;
}

// Game types
export interface Game {
  id: string;
  name: string;
  type: 'slots' | 'dice' | 'roulette' | 'blackjack' | 'baccarat' | 'sports';
  provider?: string;
  rtp: number;
  minBet: number;
  maxBet: number;
  isActive: boolean;
  thumbnailUrl?: string;
}

export interface Bet {
  id: string;
  userId: string;
  gameId: string;
  gameName?: string;
  betAmount: number;
  winAmount: number;
  multiplier: number;
  gameData?: Record<string, unknown>;
  status: 'pending' | 'won' | 'lost';
  settledAt?: string;
  createdAt: string;
}

// Admin types
export interface AdminStats {
  totalUsers: number;
  activeUsers: number;
  totalRevenue: number;
  totalBets: number;
  pendingWithdrawals: number;
  systemHealth: 'healthy' | 'warning' | 'critical';
}

export interface AuditLog {
  id: string;
  userId?: string;
  userName?: string;
  action: string;
  details?: Record<string, unknown>;
  ipAddress?: string;
  createdAt: string;
}

// API Response types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Wallet types
export interface Wallet {
  address: string;
  currency: string;
  balance: number;
  pendingBalance: number;
}

export interface DepositAddress {
  currency: string;
  address: string;
  qrCode?: string;
}

// Game state types
export interface GameState {
  gameId: string;
  betAmount: number;
  potentialWin: number;
  isPlaying: boolean;
  result?: GameResult;
}

export interface GameResult {
  outcome: string;
  multiplier: number;
  winAmount: number;
}

// Form types
export interface LoginForm {
  email: string;
  password: string;
}

export interface RegisterForm {
  email: string;
  username: string;
  password: string;
  confirmPassword: string;
}

export interface WithdrawForm {
  amount: number;
  address: string;
  currency: string;
}
