'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Badge, Input } from '@/components/ui';
import { useAuth } from '@/contexts/AuthContext';
import styles from './admin.module.css';

type Tab = 'overview' | 'users' | 'transactions' | 'games' | 'security' | 'settings';

interface User {
  id: string;
  username: string;
  email: string;
  balance: number;
  status: 'active' | 'banned' | 'pending';
  kyc: 'verified' | 'pending' | 'rejected';
  joinedAt: string;
  lastLogin: string;
}

interface Transaction {
  id: string;
  user: string;
  type: 'deposit' | 'withdrawal';
  amount: number;
  currency: string;
  status: 'pending' | 'completed' | 'rejected';
  date: string;
}

interface Game {
  id: string;
  name: string;
  type: string;
  provider: string;
  rtp: number;
  status: 'active' | 'inactive';
  minBet: number;
  maxBet: number;
}

const mockUsers: User[] = [
  { id: '1', username: 'john_doe', email: 'john@example.com', balance: 1250.50, status: 'active', kyc: 'verified', joinedAt: '2024-01-01', lastLogin: '2024-01-15' },
  { id: '2', username: 'jane_smith', email: 'jane@example.com', balance: 567.25, status: 'active', kyc: 'pending', joinedAt: '2024-01-05', lastLogin: '2024-01-14' },
  { id: '3', username: 'bob_wilson', email: 'bob@example.com', balance: 0, status: 'banned', kyc: 'rejected', joinedAt: '2023-12-20', lastLogin: '2024-01-10' },
  { id: '4', username: 'alice_brown', email: 'alice@example.com', balance: 2345.00, status: 'active', kyc: 'verified', joinedAt: '2024-01-08', lastLogin: '2024-01-15' },
  { id: '5', username: 'charlie_davis', email: 'charlie@example.com', balance: 890.75, status: 'active', kyc: 'verified', joinedAt: '2024-01-10', lastLogin: '2024-01-15' },
];

const mockTransactions: Transaction[] = [
  { id: '1', user: 'john_doe', type: 'deposit', amount: 1.5, currency: 'BTC', status: 'completed', date: '2024-01-15 14:30' },
  { id: '2', user: 'alice_brown', type: 'withdrawal', amount: 500, currency: 'USD', status: 'pending', date: '2024-01-15 13:45' },
  { id: '3', user: 'jane_smith', type: 'deposit', amount: 0.5, currency: 'ETH', status: 'completed', date: '2024-01-15 12:20' },
  { id: '4', user: 'charlie_davis', type: 'withdrawal', amount: 250, currency: 'USD', status: 'rejected', date: '2024-01-15 11:00' },
  { id: '5', user: 'john_doe', type: 'deposit', amount: 1000, currency: 'USDT', status: 'completed', date: '2024-01-14 22:30' },
];

const mockGames: Game[] = [
  { id: '1', name: 'Tiger King Slots', type: 'slots', provider: 'Pragmatic Play', rtp: 96.5, status: 'active', minBet: 0.2, maxBet: 100 },
  { id: '2', name: 'Classic Dice', type: 'dice', provider: 'TigerCasino', rtp: 99, status: 'active', minBet: 0.01, maxBet: 1000 },
  { id: '3', name: 'European Roulette', type: 'roulette', provider: 'Evolution', rtp: 97.3, status: 'active', minBet: 1, maxBet: 5000 },
  { id: '4', name: 'Blackjack VIP', type: 'blackjack', provider: 'Evolution', rtp: 99.5, status: 'active', minBet: 10, maxBet: 10000 },
  { id: '5', name: 'Mega Moolah', type: 'slots', provider: 'Microgaming', rtp: 88.12, status: 'inactive', minBet: 0.25, maxBet: 6.25 },
];

export default function AdminPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const [activeTab, setActiveTab] = useState<Tab>('overview');
  const [searchQuery, setSearchQuery] = useState('');

  const stats = {
    totalUsers: 12543,
    activeUsers: 8456,
    totalRevenue: 2456789,
    totalBets: 5678901,
    pendingWithdrawals: 23,
    systemHealth: 'healthy' as const,
  };

  const tabs = [
    { id: 'overview', name: 'Overview', icon: '📊' },
    { id: 'users', name: 'Users', icon: '👥' },
    { id: 'transactions', name: 'Transactions', icon: '💰' },
    { id: 'games', name: 'Games', icon: '🎮' },
    { id: 'security', name: 'Security', icon: '🔒' },
    { id: 'settings', name: 'Settings', icon: '⚙️' },
  ];

  const renderContent = () => {
    switch (activeTab) {
      case 'overview':
        return (
          <div className={styles.overviewGrid}>
            <Card variant="bordered" padding="lg">
              <div className={styles.statCard}>
                <span className={styles.statIcon}>👥</span>
                <div className={styles.statContent}>
                  <span className={styles.statValue}>{stats.totalUsers.toLocaleString()}</span>
                  <span className={styles.statLabel}>Total Users</span>
                </div>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statCard}>
                <span className={styles.statIcon}>🔥</span>
                <div className={styles.statContent}>
                  <span className={styles.statValue}>{stats.activeUsers.toLocaleString()}</span>
                  <span className={styles.statLabel}>Active Users</span>
                </div>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statCard}>
                <span className={styles.statIcon}>💵</span>
                <div className={styles.statContent}>
                  <span className={styles.statValue}>${(stats.totalRevenue / 1000000).toFixed(2)}M</span>
                  <span className={styles.statLabel}>Total Revenue</span>
                </div>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statCard}>
                <span className={styles.statIcon}>🎯</span>
                <div className={styles.statContent}>
                  <span className={styles.statValue}>{(stats.totalBets / 1000000).toFixed(1)}M</span>
                  <span className={styles.statLabel}>Total Bets</span>
                </div>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statCard}>
                <span className={styles.statIcon}>⏳</span>
                <div className={styles.statContent}>
                  <span className={styles.statValue}>{stats.pendingWithdrawals}</span>
                  <span className={styles.statLabel}>Pending Withdrawals</span>
                </div>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <div className={styles.statCard}>
                <span className={styles.statIcon}>💚</span>
                <div className={styles.statContent}>
                  <span className={styles.statValue}>99.9%</span>
                  <span className={styles.statLabel}>System Health</span>
                </div>
              </div>
            </Card>
          </div>
        );

      case 'users':
        return (
          <div className={styles.tableContainer}>
            <div className={styles.tableHeader}>
              <Input
                placeholder="Search users..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Username</th>
                  <th>Email</th>
                  <th>Balance</th>
                  <th>Status</th>
                  <th>KYC</th>
                  <th>Joined</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {mockUsers.filter(u => u.username.includes(searchQuery) || u.email.includes(searchQuery)).map(u => (
                  <tr key={u.id}>
                    <td>{u.id}</td>
                    <td>{u.username}</td>
                    <td>{u.email}</td>
                    <td>${u.balance.toFixed(2)}</td>
                    <td>
                      <Badge variant={u.status === 'active' ? 'success' : u.status === 'banned' ? 'error' : 'warning'} size="sm">
                        {u.status}
                      </Badge>
                    </td>
                    <td>
                      <Badge variant={u.kyc === 'verified' ? 'success' : u.kyc === 'pending' ? 'warning' : 'error'} size="sm">
                        {u.kyc}
                      </Badge>
                    </td>
                    <td>{u.joinedAt}</td>
                    <td>
                      <div className={styles.actionBtns}>
                        <Button variant="outline" size="sm">Edit</Button>
                        <Button variant={u.status === 'banned' ? 'primary' : 'danger'} size="sm">
                          {u.status === 'banned' ? 'Unban' : 'Ban'}
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        );

      case 'transactions':
        return (
          <div className={styles.tableContainer}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>User</th>
                  <th>Type</th>
                  <th>Amount</th>
                  <th>Currency</th>
                  <th>Status</th>
                  <th>Date</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {mockTransactions.map(tx => (
                  <tr key={tx.id}>
                    <td>{tx.id}</td>
                    <td>{tx.user}</td>
                    <td>
                      <Badge variant={tx.type === 'deposit' ? 'info' : 'warning'} size="sm">
                        {tx.type}
                      </Badge>
                    </td>
                    <td>{tx.amount}</td>
                    <td>{tx.currency}</td>
                    <td>
                      <Badge variant={tx.status === 'completed' ? 'success' : tx.status === 'pending' ? 'warning' : 'error'} size="sm">
                        {tx.status}
                      </Badge>
                    </td>
                    <td>{tx.date}</td>
                    <td>
                      <div className={styles.actionBtns}>
                        {tx.status === 'pending' && (
                          <>
                            <Button variant="primary" size="sm">Approve</Button>
                            <Button variant="danger" size="sm">Reject</Button>
                          </>
                        )}
                        <Button variant="outline" size="sm">View</Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        );

      case 'games':
        return (
          <div className={styles.tableContainer}>
            <div className={styles.tableHeader}>
              <Button variant="primary">Add New Game</Button>
            </div>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Provider</th>
                  <th>RTP</th>
                  <th>Min Bet</th>
                  <th>Max Bet</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {mockGames.map(game => (
                  <tr key={game.id}>
                    <td>{game.id}</td>
                    <td>{game.name}</td>
                    <td>{game.type}</td>
                    <td>{game.provider}</td>
                    <td>{game.rtp}%</td>
                    <td>${game.minBet}</td>
                    <td>${game.maxBet}</td>
                    <td>
                      <Badge variant={game.status === 'active' ? 'success' : 'error'} size="sm">
                        {game.status}
                      </Badge>
                    </td>
                    <td>
                      <div className={styles.actionBtns}>
                        <Button variant="outline" size="sm">Edit</Button>
                        <Button variant={game.status === 'active' ? 'danger' : 'primary'} size="sm">
                          {game.status === 'active' ? 'Disable' : 'Enable'}
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        );

      case 'security':
        return (
          <div className={styles.securityGrid}>
            <Card variant="bordered" padding="lg">
              <h3>Recent Security Events</h3>
              <div className={styles.securityList}>
                <div className={styles.securityItem}>
                  <span className={styles.securityIcon}>🔔</span>
                  <div className={styles.securityContent}>
                    <span className={styles.securityTitle}>Failed Login Attempt</span>
                    <span className={styles.securityDesc}>User: john_doe, IP: 192.168.1.1, Attempts: 5</span>
                    <span className={styles.securityTime}>2 minutes ago</span>
                  </div>
                </div>
                <div className={styles.securityItem}>
                  <span className={styles.securityIcon}>⚠️</span>
                  <div className={styles.securityContent}>
                    <span className={styles.securityTitle}>Suspicious Withdrawal</span>
                    <span className={styles.securityDesc}>User: bob_wilson, Amount: $10,000</span>
                    <span className={styles.securityTime}>1 hour ago</span>
                  </div>
                </div>
                <div className={styles.securityItem}>
                  <span className={styles.securityIcon}>✅</span>
                  <div className={styles.securityContent}>
                    <span className={styles.securityTitle}>KYC Approved</span>
                    <span className={styles.securityDesc}>User: alice_brown verified</span>
                    <span className={styles.securityTime}>3 hours ago</span>
                  </div>
                </div>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <h3>Security Settings</h3>
              <div className={styles.settingsList}>
                <div className={styles.settingItem}>
                  <span>Two-Factor Authentication</span>
                  <Button variant="primary" size="sm">Required</Button>
                </div>
                <div className={styles.settingItem}>
                  <span>KYC Verification</span>
                  <Button variant="primary" size="sm">Required</Button>
                </div>
                <div className={styles.settingItem}>
                  <span>Withdrawal Lock</span>
                  <Button variant="secondary" size="sm">48 hours</Button>
                </div>
                <div className={styles.settingItem}>
                  <span>IP Whitelist</span>
                  <Button variant="outline" size="sm">Configure</Button>
                </div>
              </div>
            </Card>
          </div>
        );

      case 'settings':
        return (
          <div className={styles.settingsGrid}>
            <Card variant="bordered" padding="lg">
              <h3>General Settings</h3>
              <div className={styles.settingsForm}>
                <Input label="Site Name" defaultValue="TigerCasino" />
                <Input label="Support Email" defaultValue="support@tigercasino.com" />
                <Input label="Minimum Withdrawal" defaultValue="10" />
                <Input label="Maximum Withdrawal" defaultValue="100000" />
                <Button variant="primary">Save Changes</Button>
              </div>
            </Card>
            <Card variant="bordered" padding="lg">
              <h3>API Keys</h3>
              <div className={styles.settingsForm}>
                <Input label="API Key" defaultValue="sk_live_xxxxx" type="password" />
                <Input label="Webhook URL" defaultValue="https://tigercasino.com/webhook" />
                <Button variant="primary">Generate New Key</Button>
              </div>
            </Card>
          </div>
        );

      default:
        return null;
    }
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
      </div>
    );
  }

  if (!isAuthenticated || !user?.isAdmin) {
    return (
      <div className={styles.container}>
        <Header />
        <main className={styles.main}>
          <div className={styles.notAuthorized}>
            <h2>Access Denied</h2>
            <p>You don't have permission to access this page.</p>
            <Link href="/">
              <Button variant="primary">Go Home</Button>
            </Link>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <>
      <Header />
      <main className={styles.main}>
        <div className={styles.adminHeader}>
          <h1>Admin Dashboard</h1>
          <p>Manage your casino platform</p>
        </div>

        <div className={styles.adminLayout}>
          <aside className={styles.sidebar}>
            <nav className={styles.nav}>
              {tabs.map(tab => (
                <button
                  key={tab.id}
                  className={`${styles.navItem} ${activeTab === tab.id ? styles.active : ''}`}
                  onClick={() => setActiveTab(tab.id as Tab)}
                >
                  <span className={styles.navIcon}>{tab.icon}</span>
                  <span className={styles.navName}>{tab.name}</span>
                </button>
              ))}
            </nav>
          </aside>

          <div className={styles.content}>
            {renderContent()}
          </div>
        </div>
      </main>
    </>
  );
}
