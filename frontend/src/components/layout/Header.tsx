'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui';
import styles from './Header.module.css';

export default function Header() {
  const { user, isAuthenticated, logout } = useAuth();
  const [scrolled, setScrolled] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <header className={`${styles.header} ${scrolled ? styles.scrolled : ''}`}>
      <div className={styles.container}>
        <Link href="/" className={styles.logo}>
          <span className={styles.logoIcon}>🐯</span>
          <span className={styles.logoText}>TigerCasino</span>
        </Link>

        <nav className={`${styles.nav} ${mobileMenuOpen ? styles.open : ''}`}>
          <Link href="/games" className={styles.navLink}>Games</Link>
          <Link href="/dashboard" className={styles.navLink}>Dashboard</Link>
          <Link href="/wallet" className={styles.navLink}>Wallet</Link>
          {isAuthenticated && user?.isAdmin && (
            <Link href="/admin" className={styles.navLink}>Admin</Link>
          )}
        </nav>

        <div className={styles.actions}>
          {isAuthenticated ? (
            <>
              <div className={styles.balance}>
                <span className={styles.balanceLabel}>Balance:</span>
                <span className={styles.balanceValue}>
                  ${user?.balance.toFixed(2)}
                </span>
              </div>
              <div className={styles.userMenu}>
                <span className={styles.username}>{user?.username}</span>
                <Button variant="outline" size="sm" onClick={logout}>
                  Logout
                </Button>
              </div>
            </>
          ) : (
            <>
              <Link href="/auth/login">
                <Button variant="ghost" size="sm">Login</Button>
              </Link>
              <Link href="/auth/register">
                <Button variant="primary" size="sm">Register</Button>
              </Link>
            </>
          )}
        </div>

        <button 
          className={styles.mobileToggle}
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>
    </header>
  );
}
