import { ApiResponse, PaginatedResponse, User, Transaction, Game, Bet, AdminStats, AuditLog, DepositAddress, LoginForm, RegisterForm, WithdrawForm } from '@/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  private token: string | null = null;

  setToken(token: string | null) {
    this.token = token;
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('token', token);
      } else {
        localStorage.removeItem('token');
      }
    }
  }

  getToken(): string | null {
    if (this.token) return this.token;
    if (typeof window !== 'undefined') {
      return localStorage.getItem('token');
    }
    return null;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const token = this.getToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        return {
          success: false,
          error: data.error || 'An error occurred',
        };
      }

      return {
        success: true,
        data,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Network error',
      };
    }
  }

  // Auth endpoints
  async login(credentials: LoginForm): Promise<ApiResponse<{ user: User; token: string }>> {
    return this.request('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async register(data: RegisterForm): Promise<ApiResponse<{ user: User; token: string }>> {
    return this.request('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async logout(): Promise<ApiResponse<void>> {
    const result = await this.request<void>('/api/auth/logout', {
      method: 'POST',
    });
    this.setToken(null);
    return result;
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    return this.request('/api/users/me');
  }

  // Wallet endpoints
  async getBalance(): Promise<ApiResponse<{ balance: number; bonusBalance: number }>> {
    return this.request('/api/wallet/balance');
  }

  async getDepositAddress(currency: string): Promise<ApiResponse<DepositAddress>> {
    return this.request(`/api/wallet/deposit/address?currency=${currency}`);
  }

  async withdraw(data: WithdrawForm): Promise<ApiResponse<Transaction>> {
    return this.request('/api/wallet/withdraw', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getTransactions(page = 1, pageSize = 20): Promise<ApiResponse<PaginatedResponse<Transaction>>> {
    return this.request(`/api/wallet/transactions?page=${page}&pageSize=${pageSize}`);
  }

  // Game endpoints
  async getGames(): Promise<ApiResponse<Game[]>> {
    return this.request('/api/games');
  }

  async getGame(id: string): Promise<ApiResponse<Game>> {
    return this.request(`/api/games/${id}`);
  }

  async placeBet(gameId: string, betAmount: number, betData?: Record<string, unknown>): Promise<ApiResponse<Bet>> {
    return this.request(`/api/games/${gameId}/bet`, {
      method: 'POST',
      body: JSON.stringify({ amount: betAmount, ...betData }),
    });
  }

  async getBetHistory(gameId?: string, page = 1, pageSize = 20): Promise<ApiResponse<PaginatedResponse<Bet>>> {
    const query = gameId ? `&gameId=${gameId}` : '';
    return this.request(`/api/games/history?page=${page}&pageSize=${pageSize}${query}`);
  }

  // Admin endpoints
  async getAdminStats(): Promise<ApiResponse<AdminStats>> {
    return this.request('/api/admin/dashboard');
  }

  async getUsers(page = 1, pageSize = 20, search?: string): Promise<ApiResponse<PaginatedResponse<User>>> {
    const query = search ? `&search=${search}` : '';
    return this.request(`/api/admin/users?page=${page}&pageSize=${pageSize}${query}`);
  }

  async updateUser(userId: string, data: Partial<User>): Promise<ApiResponse<User>> {
    return this.request(`/api/admin/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async getAllTransactions(page = 1, pageSize = 20, status?: string): Promise<ApiResponse<PaginatedResponse<Transaction>>> {
    const query = status ? `&status=${status}` : '';
    return this.request(`/api/admin/transactions?page=${page}&pageSize=${pageSize}${query}`);
  }

  async approveTransaction(transactionId: string): Promise<ApiResponse<Transaction>> {
    return this.request(`/api/admin/transactions/${transactionId}/approve`, {
      method: 'POST',
    });
  }

  async rejectTransaction(transactionId: string): Promise<ApiResponse<Transaction>> {
    return this.request(`/api/admin/transactions/${transactionId}/reject`, {
      method: 'POST',
    });
  }

  async getAuditLogs(page = 1, pageSize = 50): Promise<ApiResponse<PaginatedResponse<AuditLog>>> {
    return this.request(`/api/admin/audit-logs?page=${page}&pageSize=${pageSize}`);
  }

  async getAdminGames(): Promise<ApiResponse<Game[]>> {
    return this.request('/api/admin/games');
  }

  async updateGame(gameId: string, data: Partial<Game>): Promise<ApiResponse<Game>> {
    return this.request(`/api/admin/games/${gameId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // ============ VIP & LOYALTY ENDPOINTS ============

  async getVIPStatus(): Promise<ApiResponse<VIPStatus>> {
    return this.request('/api/vip/status');
  }

  async claimRakeback(): Promise<ApiResponse<{ amount: number }>> {
    return this.request('/api/vip/rakeback/claim', {
      method: 'POST',
    });
  }

  async redeemPoints(points: number): Promise<ApiResponse<{ amount: number }>> {
    return this.request('/api/vip/points/redeem', {
      method: 'POST',
      body: JSON.stringify({ points }),
    });
  }

  async claimWelcomeBonus(): Promise<ApiResponse<{ bonusId: string; amount: number; wagerReq: number }>> {
    return this.request('/api/vip/bonus/welcome/claim', {
      method: 'POST',
    });
  }

  async claimDepositBonus(depositAmount: number): Promise<ApiResponse<{ bonusId: string; amount: number; wagerReq: number }>> {
    return this.request('/api/vip/bonus/deposit/claim', {
      method: 'POST',
      body: JSON.stringify({ amount: depositAmount }),
    });
  }

  async getLeaderboard(period: string = 'weekly', limit: number = 100): Promise<ApiResponse<LeaderboardEntry[]>> {
    return this.request(`/api/vip/leaderboard?period=${period}&limit=${limit}`);
  }

  async getPromotions(): Promise<ApiResponse<Promotion[]>> {
    return this.request('/api/vip/promotions');
  }

  // ============ TOURNAMENT ENDPOINTS ============

  async getTournaments(): Promise<ApiResponse<Tournament[]>> {
    return this.request('/api/tournaments');
  }

  async getTournamentDetails(tournamentId: string): Promise<ApiResponse<TournamentDetails>> {
    return this.request(`/api/tournaments/${tournamentId}`);
  }

  async registerTournament(tournamentId: string): Promise<ApiResponse<{ success: boolean }>> {
    return this.request(`/api/tournaments/${tournamentId}/register`, {
      method: 'POST',
    });
  }

  async getTournamentLeaderboard(tournamentId: string): Promise<ApiResponse<LeaderboardEntry[]>> {
    return this.request(`/api/tournaments/${tournamentId}/leaderboard`);
  }

  async getMyTournaments(): Promise<ApiResponse<Tournament[]>> {
    return this.request('/api/tournaments/my');
  }

  async getTournamentResults(): Promise<ApiResponse<TournamentResult[]>> {
    return this.request('/api/tournaments/results');
  }

  // ============ SPORTSBOOK ENDPOINTS ============

  async getSports(): Promise<ApiResponse<Sport[]>> {
    return this.request('/api/sportsbook/sports');
  }

  async getLeagues(sportId: string): Promise<ApiResponse<League[]>> {
    return this.request(`/api/sportsbook/leagues?sportId=${sportId}`);
  }

  async getEvents(sportId: string, leagueId: string, date?: string): Promise<ApiResponse<SportsEvent[]>> {
    let query = `?sportId=${sportId}&leagueId=${leagueId}`;
    if (date) query += `&date=${date}`;
    return this.request(`/api/sportsbook/events${query}`);
  }

  async getLiveEvents(sportId?: string): Promise<ApiResponse<SportsEvent[]>> {
    const query = sportId ? `?sportId=${sportId}` : '';
    return this.request(`/api/sportsbook/live${query}`);
  }

  async getEventDetails(eventId: string): Promise<ApiResponse<EventDetails>> {
    return this.request(`/api/sportsbook/events/${eventId}`);
  }

  async placeSportsBet(eventId: string, marketId: string, selection: string, stake: number): Promise<ApiResponse<{ betId: string; potentialWin: number }>> {
    return this.request('/api/sportsbook/bet', {
      method: 'POST',
      body: JSON.stringify({ eventId, marketId, selection, stake }),
    });
  }

  async getMySportsBets(status?: string, page = 1): Promise<ApiResponse<PaginatedResponse<SportsBet>>> {
    const query = status ? `?status=${status}&page=${page}` : `?page=${page}`;
    return this.request(`/api/sportsbook/bets${query}`);
  }

  async getSportsBettingStats(): Promise<ApiResponse<BettingStats>> {
    return this.request('/api/sportsbook/stats');
  }

  // ============ GAME AGGREGATOR ENDPOINTS ============

  async getGameProviders(): Promise<ApiResponse<GameProvider[]>> {
    return this.request('/api/games/providers');
  }

  async launchGame(gameId: string, mode: string = 'real'): Promise<ApiResponse<{ gameUrl: string; token: string }>> {
    return this.request('/api/games/launch', {
      method: 'POST',
      body: JSON.stringify({ gameId, mode }),
    });
  }

  async getJackpots(): Promise<ApiResponse<JackpotInfo>> {
    return this.request('/api/games/jackpots');
  }
}

// Types for new endpoints
export interface VIPStatus {
  level: number;
  levelName: string;
  totalWagered: number;
  points: number;
  rakebackPercent: number;
  rakebackBalance: number;
  benefits: VIPBenefits;
  nextLevel?: {
    level: number;
    name: string;
    minWagered: number;
  };
  progressToNext: number;
}

export interface VIPBenefits {
  maxBet: number;
  maxWin: number;
  withdrawalLimit: number;
  withdrawalFee: number;
  cashbackPercent: number;
  pointsMultiplier: number;
  prioritySupport: boolean;
  instantWithdraw: boolean;
  personalHost: boolean;
}

export interface LeaderboardEntry {
  rank: number;
  userId: string;
  username: string;
  score: number;
  wagered: number;
  wins: number;
}

export interface Promotion {
  id: string;
  name: string;
  description: string;
  type: string;
  bonusAmount: number;
  wagerReq: number;
  startDate: string;
  endDate: string;
}

export interface Tournament {
  id: string;
  name: string;
  description: string;
  type: string;
  status: string;
  prizePool: number;
  currency: string;
  startTime: string;
  endTime: string;
  participantCount?: number;
}

export interface TournamentDetails extends Tournament {
  minBet: number;
  gameFilter: string[];
  prizeDistribution: PrizeBreakdown[];
  leaderboard: LeaderboardEntry[];
}

export interface PrizeBreakdown {
  position: number;
  percent: number;
  amount: number;
}

export interface TournamentResult {
  tournamentId: string;
  tournamentName: string;
  position: number;
  prizeAmount: number;
  currency: string;
  date: string;
}

export interface Sport {
  id: string;
  name: string;
  shortName: string;
  icon: string;
}

export interface League {
  id: string;
  sportId: string;
  name: string;
  country: string;
  logo: string;
}

export interface SportsEvent {
  id: string;
  sportId: string;
  leagueId: string;
  homeTeam: string;
  awayTeam: string;
  startTime: string;
  status: string;
  homeScore?: number;
  awayScore?: number;
  period?: string;
  timeRemaining?: string;
}

export interface EventDetails {
  event: SportsEvent;
  markets: Market[];
}

export interface Market {
  id: string;
  name: string;
  marketType: string;
  outcomes: Outcome[];
  suspended: boolean;
}

export interface Outcome {
  id: string;
  name: string;
  odds: number;
  line?: number;
}

export interface SportsBet {
  id: string;
  eventId: string;
  eventName: string;
  selection: string;
  odds: number;
  stake: number;
  potentialWin: number;
  actualWin?: number;
  status: string;
  createdAt: string;
}

export interface BettingStats {
  totalBets: number;
  pendingBets: number;
  wonBets: number;
  lostBets: number;
  totalStaked: number;
  totalWon: number;
  profit: number;
  winRate: number;
}

export interface GameProvider {
  id: string;
  name: string;
  code: string;
  logo: string;
  gameCount: number;
  isAggregator: boolean;
  status: string;
}

export interface JackpotInfo {
  mini: number;
  major: number;
  grand: number;
}

export const api = new ApiClient();
