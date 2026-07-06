'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Badge, Input } from '@/components/ui';
import styles from './whitelabel.module.css';

interface WhiteLabel {
  id: string;
  name: string;
  domain: string;
  status: 'pending' | 'active' | 'suspended' | 'rejected';
  commissionRate: number;
  totalRevenue: number;
  totalFees: number;
  createdAt: string;
  owner: string;
}

const mockWhiteLabels: WhiteLabel[] = [
  { id: '1', name: 'LuckyStar Casino', domain: 'luckystar.com', status: 'active', commissionRate: 80, totalRevenue: 1250000, totalFees: 250000, createdAt: '2024-01-15', owner: 'John Doe' },
  { id: '2', name: 'RoyalBet Platform', domain: 'royalbet.io', status: 'pending', commissionRate: 80, totalRevenue: 0, totalFees: 0, createdAt: '2024-02-20', owner: 'Jane Smith' },
  { id: '3', name: 'CryptoGaming Pro', domain: 'cryptogaming.pro', status: 'active', commissionRate: 80, totalRevenue: 890000, totalFees: 178000, createdAt: '2024-01-28', owner: 'Mike Johnson' },
  { id: '4', name: 'BitWager Hub', domain: 'bitwagerhub.com', status: 'suspended', commissionRate: 80, totalRevenue: 450000, totalFees: 90000, createdAt: '2024-02-01', owner: 'Sarah Wilson' },
];

export default function AdminWhiteLabelPage() {
  const [whiteLabels, setWhiteLabels] = useState<WhiteLabel[]>(mockWhiteLabels);
  const [selectedStatus, setSelectedStatus] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');

  const filteredWhiteLabels = whiteLabels.filter(wl => {
    const matchesStatus = selectedStatus === 'all' || wl.status === selectedStatus;
    const matchesSearch = wl.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                         wl.domain.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesStatus && matchesSearch;
  });

  const handleApprove = (id: string) => {
    setWhiteLabels(prev => prev.map(wl => 
      wl.id === id ? { ...wl, status: 'active' as const } : wl
    ));
  };

  const handleReject = (id: string) => {
    setWhiteLabels(prev => prev.map(wl => 
      wl.id === id ? { ...wl, status: 'rejected' as const } : wl
    ));
  };

  const handleSuspend = (id: string) => {
    setWhiteLabels(prev => prev.map(wl => 
      wl.id === id ? { ...wl, status: 'suspended' as const } : wl
    ));
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'success';
      case 'pending': return 'warning';
      case 'suspended': return 'error';
      case 'rejected': return 'default';
      default: return 'default';
    }
  };

  const totalRevenue = whiteLabels.reduce((sum, wl) => sum + wl.totalRevenue, 0);
  const totalFees = whiteLabels.reduce((sum, wl) => sum + wl.totalFees, 0);
  const activeBrands = whiteLabels.filter(wl => wl.status === 'active').length;

  return (
    <>
      <Header />
      <main className={styles.main}>
        <div className={styles.container}>
          <div className={styles.header}>
            <h1 className={styles.title}>White Label Management</h1>
            <p className={styles.subtitle}>
              Manage white label brands. 20% platform fee auto-deducted from all earnings.
            </p>
          </div>

          <div className={styles.statsGrid}>
            <Card className={styles.statCard}>
              <div className={styles.statIcon}>💰</div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>${totalRevenue.toLocaleString()}</div>
                <div className={styles.statLabel}>Total Brand Revenue</div>
              </div>
            </Card>
            <Card className={styles.statCard}>
              <div className={styles.statIcon}>🏛️</div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>${totalFees.toLocaleString()}</div>
                <div className={styles.statLabel}>Platform Fees (20%)</div>
              </div>
            </Card>
            <Card className={styles.statCard}>
              <div className={styles.statIcon}>🌐</div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>{activeBrands}</div>
                <div className={styles.statLabel}>Active Brands</div>
              </div>
            </Card>
            <Card className={styles.statCard}>
              <div className={styles.statIcon}>⏳</div>
              <div className={styles.statContent}>
                <div className={styles.statValue}>
                  {whiteLabels.filter(wl => wl.status === 'pending').length}
                </div>
                <div className={styles.statLabel}>Pending Approval</div>
              </div>
            </Card>
          </div>

          <div className={styles.filters}>
            <Input
              type="text"
              placeholder="Search brands..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className={styles.searchInput}
            />
            <div className={styles.statusFilters}>
              <button
                className={`${styles.filterBtn} ${selectedStatus === 'all' ? styles.active : ''}`}
                onClick={() => setSelectedStatus('all')}
              >
                All
              </button>
              <button
                className={`${styles.filterBtn} ${selectedStatus === 'pending' ? styles.active : ''}`}
                onClick={() => setSelectedStatus('pending')}
              >
                Pending
              </button>
              <button
                className={`${styles.filterBtn} ${selectedStatus === 'active' ? styles.active : ''}`}
                onClick={() => setSelectedStatus('active')}
              >
                Active
              </button>
              <button
                className={`${styles.filterBtn} ${selectedStatus === 'suspended' ? styles.active : ''}`}
                onClick={() => setSelectedStatus('suspended')}
              >
                Suspended
              </button>
            </div>
          </div>

          <div className={styles.whiteLabelList}>
            {filteredWhiteLabels.map(wl => (
              <Card key={wl.id} className={styles.whiteLabelCard}>
                <div className={styles.wlHeader}>
                  <div className={styles.wlInfo}>
                    <h3 className={styles.wlName}>{wl.name}</h3>
                    <p className={styles.wlDomain}>{wl.domain}</p>
                  </div>
                  <Badge variant={getStatusColor(wl.status)}>
                    {wl.status.toUpperCase()}
                  </Badge>
                </div>

                <div className={styles.wlStats}>
                  <div className={styles.wlStat}>
                    <span className={styles.wlStatLabel}>Owner</span>
                    <span className={styles.wlStatValue}>{wl.owner}</span>
                  </div>
                  <div className={styles.wlStat}>
                    <span className={styles.wlStatLabel}>Revenue</span>
                    <span className={styles.wlStatValue}>${wl.totalRevenue.toLocaleString()}</span>
                  </div>
                  <div className={styles.wlStat}>
                    <span className={styles.wlStatLabel}>Platform Fee (20%)</span>
                    <span className={styles.wlStatValue}>${wl.totalFees.toLocaleString()}</span>
                  </div>
                  <div className={styles.wlStat}>
                    <span className={styles.wlStatLabel}>Commission</span>
                    <span className={styles.wlStatValue}>{wl.commissionRate}%</span>
                  </div>
                  <div className={styles.wlStat}>
                    <span className={styles.wlStatLabel}>Created</span>
                    <span className={styles.wlStatValue}>{wl.createdAt}</span>
                  </div>
                </div>

                <div className={styles.wlActions}>
                  {wl.status === 'pending' && (
                    <>
                      <Button onClick={() => handleApprove(wl.id)} className={styles.approveBtn}>
                        Approve
                      </Button>
                      <Button onClick={() => handleReject(wl.id)} variant="secondary">
                        Reject
                      </Button>
                    </>
                  )}
                  {wl.status === 'active' && (
                    <Button onClick={() => handleSuspend(wl.id)} variant="secondary">
                      Suspend
                    </Button>
                  )}
                  {wl.status === 'suspended' && (
                    <Button onClick={() => handleApprove(wl.id)}>
                      Reactivate
                    </Button>
                  )}
                  <Link href={`/admin/whitelabel/${wl.id}`} className={styles.viewLink}>
                    View Details →
                  </Link>
                </div>
              </Card>
            ))}
          </div>

          {filteredWhiteLabels.length === 0 && (
            <div className={styles.empty}>
              <p>No white label brands found.</p>
            </div>
          )}
        </div>
      </main>
      <Footer />
    </>
  );
}
