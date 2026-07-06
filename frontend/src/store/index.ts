import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id: string;
  email: string;
  username: string;
  balance: number;
  bonusBalance: number;
  vipLevel: number;
  vipPoints: number;
  isAdmin: boolean;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, token: string) => void;
  updateBalance: (balance: number, bonusBalance?: number) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (user, token) => set({ user, token, isAuthenticated: true }),
      updateBalance: (balance, bonusBalance) => set((state) => ({
        user: state.user ? { ...state.user, balance, bonusBalance: bonusBalance ?? state.user.bonusBalance } : null
      })),
      logout: () => set({ user: null, token: null, isAuthenticated: false }),
    }),
    {
      name: 'tigercasino-auth',
    }
  )
);

// Game state store
interface GameState {
  currentGame: string | null;
  currentMultiplier: number;
  betAmount: number;
  autoCashout: number | null;
  isPlaying: boolean;
  setCurrentGame: (game: string | null) => void;
  setMultiplier: (multiplier: number) => void;
  setBetAmount: (amount: number) => void;
  setAutoCashout: (value: number | null) => void;
  setIsPlaying: (playing: boolean) => void;
}

export const useGameStore = create<GameState>()((set) => ({
  currentGame: null,
  currentMultiplier: 1.0,
  betAmount: 1.0,
  autoCashout: null,
  isPlaying: false,
  setCurrentGame: (game) => set({ currentGame: game }),
  setMultiplier: (multiplier) => set({ currentMultiplier: multiplier }),
  setBetAmount: (amount) => set({ betAmount: amount }),
  setAutoCashout: (value) => set({ autoCashout: value }),
  setIsPlaying: (playing) => set({ isPlaying: playing }),
}));

// UI state store
interface UIState {
  sidebarOpen: boolean;
  theme: 'dark' | 'light';
  toggleSidebar: () => void;
  setTheme: (theme: 'dark' | 'light') => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      sidebarOpen: false,
      theme: 'dark',
      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      setTheme: (theme) => set({ theme }),
    }),
    {
      name: 'tigercasino-ui',
    }
  )
);
