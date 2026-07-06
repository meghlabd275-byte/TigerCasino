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
}

export const api = new ApiClient();
