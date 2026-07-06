'use client';

import React, { useState, useEffect } from 'react';
import styles from './ThemeToggle.module.css';

export default function ThemeToggle() {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    const savedTheme = localStorage.getItem('theme') as 'dark' | 'light';
    if (savedTheme) {
      setTheme(savedTheme);
      document.documentElement.setAttribute('data-theme', savedTheme);
    }
  }, []);

  const toggleTheme = () => {
    const newTheme = theme === 'dark' ? 'light' : 'dark';
    setTheme(newTheme);
    localStorage.setItem('theme', newTheme);
    document.documentElement.setAttribute('data-theme', newTheme);
  };

  if (!mounted) {
    return null;
  }

  return (
    <div className={styles.container}>
      <button
        className={styles.toggle}
        onClick={toggleTheme}
        aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
      >
        <span className={`${styles.icon} ${theme === 'dark' ? styles.sun : styles.moon}`}>
          {theme === 'dark' ? '☀️' : '🌙'}
        </span>
        <span className={styles.label}>
          {theme === 'dark' ? 'Light' : 'Dark'}
        </span>
      </button>
    </div>
  );
}
