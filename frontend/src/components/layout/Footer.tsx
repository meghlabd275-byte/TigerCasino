'use client';

import React from 'react';
import Link from 'next/link';
import styles from './Footer.module.css';

export default function Footer() {
  return (
    <footer className={styles.footer}>
      <div className={styles.container}>
        <div className={styles.grid}>
          <div className={styles.brand}>
            <div className={styles.logo}>
              <span className={styles.logoIcon}>🐯</span>
              <span className={styles.logoText}>TigerCasino</span>
            </div>
            <p className={styles.description}>
              The ultimate cryptocurrency casino experience. Play your favorite games with instant crypto transactions.
            </p>
          </div>

          <div className={styles.links}>
            <h4 className={styles.title}>Quick Links</h4>
            <Link href="/games">Games</Link>
            <Link href="/dashboard">Dashboard</Link>
            <Link href="/wallet">Wallet</Link>
            <Link href="/auth/register">Get Started</Link>
          </div>

          <div className={styles.links}>
            <h4 className={styles.title}>Support</h4>
            <Link href="/faq">FAQ</Link>
            <Link href="/contact">Contact Us</Link>
            <Link href="/terms">Terms of Service</Link>
            <Link href="/privacy">Privacy Policy</Link>
          </div>

          <div className={styles.links}>
            <h4 className={styles.title}>Games</h4>
            <Link href="/games/slots">Slots</Link>
            <Link href="/games/dice">Dice</Link>
            <Link href="/games/roulette">Roulette</Link>
            <Link href="/games/blackjack">Blackjack</Link>
          </div>
        </div>

        <div className={styles.bottom}>
          <p className={styles.copyright}>
            © 2024 TigerCasino. All rights reserved.
          </p>
          <div className={styles.crypto}>
            <span>Accepts:</span>
            <span className={styles.icons}>₿ Ξ USDT</span>
          </div>
        </div>
      </div>
    </footer>
  );
}
