'use client';

import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react';
import { AuthState, User, LoginForm, RegisterForm } from '@/types';
import { api } from '@/lib/api';

type AuthAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'LOGIN_SUCCESS'; payload: { user: User; token: string } }
  | { type: 'LOGOUT' }
  | { type: 'UPDATE_USER'; payload: User }
  | { type: 'SET_ERROR'; payload: string };

const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
  token: null,
};

function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };
    case 'LOGIN_SUCCESS':
      return {
        ...state,
        user: action.payload.user,
        token: action.payload.token,
        isAuthenticated: true,
        isLoading: false,
      };
    case 'LOGOUT':
      return {
        ...state,
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
      };
    case 'UPDATE_USER':
      return { ...state, user: action.payload };
    case 'SET_ERROR':
      return { ...state, isLoading: false };
    default:
      return state;
  }
}

interface AuthContextType extends AuthState {
  login: (credentials: LoginForm) => Promise<{ success: boolean; error?: string }>;
  register: (data: RegisterForm) => Promise<{ success: boolean; error?: string }>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(authReducer, initialState);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    const token = api.getToken();
    if (token) {
      const response = await api.getCurrentUser();
      if (response.success && response.data) {
        dispatch({ type: 'LOGIN_SUCCESS', payload: { user: response.data, token } });
      } else {
        dispatch({ type: 'LOGOUT' });
      }
    } else {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const login = async (credentials: LoginForm) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    const response = await api.login(credentials);
    
    if (response.success && response.data) {
      api.setToken(response.data.token);
      dispatch({ type: 'LOGIN_SUCCESS', payload: response.data });
      return { success: true };
    }
    
    dispatch({ type: 'SET_ERROR', payload: response.error || 'Login failed' });
    return { success: false, error: response.error };
  };

  const register = async (data: RegisterForm) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    const response = await api.register(data);
    
    if (response.success && response.data) {
      api.setToken(response.data.token);
      dispatch({ type: 'LOGIN_SUCCESS', payload: response.data });
      return { success: true };
    }
    
    dispatch({ type: 'SET_ERROR', payload: response.error || 'Registration failed' });
    return { success: false, error: response.error };
  };

  const logout = async () => {
    await api.logout();
    dispatch({ type: 'LOGOUT' });
  };

  const refreshUser = async () => {
    const response = await api.getCurrentUser();
    if (response.success && response.data) {
      dispatch({ type: 'UPDATE_USER', payload: response.data });
    }
  };

  return (
    <AuthContext.Provider value={{ ...state, login, register, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
