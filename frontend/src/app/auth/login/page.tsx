'use client';

import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Card } from '@/components/ui';
import styles from '../auth.module.css';

type InputMode = 'email' | 'phone';

export default function LoginPage() {
  const router = useRouter();
  const { login, isAuthenticated, isLoading: authLoading } = useAuth();
  
  const [inputValue, setInputValue] = useState('');
  const [inputMode, setInputMode] = useState<InputMode>('email');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [otpSent, setOtpSent] = useState(false);
  const [otp, setOtp] = useState('');
  const [failedAttempts, setFailedAttempts] = useState(0);
  const [isLocked, setIsLocked] = useState(false);
  const [lockEndTime, setLockEndTime] = useState<Date | null>(null);

  const detectInputMode = useCallback((value: string) => {
    const trimmed = value.trim().toLowerCase();
    if (trimmed.includes('@') && trimmed.includes('.')) return 'email';
    const phonePattern = /^[\d\s\-\+\(\)]+$/;
    if (trimmed.startsWith('+') || phonePattern.test(trimmed.replace(/\s/g, ''))) {
      if (!trimmed.includes('@') && trimmed.length >= 7) return 'phone';
    }
    return 'email';
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setInputValue(value);
    const detectedMode = detectInputMode(value);
    setInputMode(detectedMode);
    setError('');
  };

  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace' && inputValue.length <= 1) {
      const newMode = detectInputMode(inputValue.slice(0, -1));
      setInputMode(newMode);
    }
  };

  const handleContinue = async () => {
    if (!inputValue.trim()) {
      setError('Please enter your email or phone number');
      return;
    }
    if (isLocked && lockEndTime && new Date() < lockEndTime) {
      const remaining = Math.ceil((lockEndTime.getTime() - Date.now()) / 1000 / 60);
      setError(`Account locked. Try again in ${remaining} minutes`);
      return;
    }
    setIsLoading(true);
    setError('');
    await new Promise(resolve => setTimeout(resolve, 1000));
    setOtpSent(true);
    setIsLoading(false);
  };

  const handleLogin = async () => {
    if (!password) {
      setError('Please enter your password');
      return;
    }
    if (isLocked && lockEndTime && new Date() < lockEndTime) {
      const remaining = Math.ceil((lockEndTime.getTime() - Date.now()) / 1000 / 60);
      setError(`Account locked. Try again in ${remaining} minutes`);
      return;
    }
    setIsLoading(true);
    setError('');
    const result = await login({ email: inputValue, password });
    if (result.success) {
      router.push('/dashboard');
    } else {
      setFailedAttempts(prev => prev + 1);
      if (failedAttempts + 1 >= 5) {
        const lockTime = new Date(Date.now() + 48 * 60 * 60 * 1000);
        setLockEndTime(lockTime);
        setIsLocked(true);
        setError('Too many failed attempts. Account locked for 48 hours.');
      } else {
        setError(result.error || 'Invalid credentials');
      }
    }
    setIsLoading(false);
  };

  const handleOtpVerify = async () => {
    if (otp.length !== 6) {
      setError('Please enter the 6-digit OTP');
      return;
    }
    setIsLoading(true);
    setError('');
    await new Promise(resolve => setTimeout(resolve, 1000));
    setIsLoading(false);
    setOtpSent(false);
  };

  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      router.push('/dashboard');
    }
  }, [authLoading, isAuthenticated, router]);

  if (authLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.background}>
        <div className={styles.glow}></div>
      </div>
      
      <Card variant="glow" padding="lg" className={styles.card}>
        <div className={styles.header}>
          <Link href="/" className={styles.logo}>
            <span className={styles.logoIcon}>🐯</span>
            <span className={styles.logoText}>TigerCasino</span>
          </Link>
          <h1 className={styles.title}>Welcome Back</h1>
          <p className={styles.subtitle}>Sign in to continue to your account</p>
        </div>

        {!otpSent ? (
          <>
            <div className={styles.formGroup}>
              <label className={styles.label}>
                {inputMode === 'phone' ? 'Phone Number' : 'Email Address'}
              </label>
              <div className={styles.inputWrapper}>
                {inputMode === 'phone' && (
                  <div className={styles.countrySelector}>
                    <select className={styles.countrySelect}>
                      <option value="+1">🇺🇸 +1</option>
                      <option value="+44">🇬🇧 +44</option>
                      <option value="+91">🇮🇳 +91</option>
                      <option value="+86">🇨🇳 +86</option>
                      <option value="+81">🇯🇵 +81</option>
                      <option value="+49">🇩🇪 +49</option>
                      <option value="+33">🇫🇷 +33</option>
                      <option value="+61">🇦🇺 +61</option>
                      <option value="+55">🇧🇷 +55</option>
                      <option value="+7">🇷🇺 +7</option>
                    </select>
                  </div>
                )}
                <input
                  type={inputMode === 'email' ? 'email' : 'tel'}
                  value={inputValue}
                  onChange={handleInputChange}
                  onKeyDown={handleInputKeyDown}
                  placeholder={inputMode === 'email' ? 'Enter your email' : 'Enter phone number'}
                  className={styles.input}
                  autoComplete="off"
                />
              </div>
            </div>

            <div className={styles.formGroup}>
              <label className={styles.label}>Password</label>
              <div className={styles.passwordWrapper}>
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  className={styles.input}
                />
                <button
                  type="button"
                  className={styles.togglePassword}
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? '👁️' : '👁️‍🗨️'}
                </button>
              </div>
            </div>

            <div className={styles.options}>
              <label className={styles.checkbox}>
                <input
                  type="checkbox"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                />
                <span>Remember me</span>
              </label>
              <Link href="/auth/forgot-password" className={styles.link}>
                Forgot Password?
              </Link>
            </div>

            {error && <div className={styles.error}>{error}</div>}

            <Button
              variant="primary"
              size="lg"
              fullWidth
              onClick={handleLogin}
              isLoading={isLoading}
            >
              Sign In
            </Button>
          </>
        ) : (
          <>
            <div className={styles.otpSection}>
              <p className={styles.otpInfo}>
                Enter the 6-digit code sent to your {inputMode === 'phone' ? 'phone' : 'email'}
              </p>
              
              <div className={styles.otpInputs}>
                {[0, 1, 2, 3, 4, 5].map((i) => (
                  <input
                    key={i}
                    type="text"
                    maxLength={1}
                    className={styles.otpInput}
                    value={otp[i] || ''}
                    onChange={(e) => {
                      const val = e.target.value;
                      if (val.match(/^\d$/)) {
                        const newOtp = otp.split('');
                        newOtp[i] = val;
                        setOtp(newOtp.join(''));
                        if (i < 5) {
                          const inputs = document.querySelectorAll('.otp-input');
                          (inputs[i + 1] as HTMLInputElement)?.focus();
                        }
                      }
                    }}
                    onKeyDown={(e) => {
                      if (e.key === 'Backspace' && !otp[i] && i > 0) {
                        const inputs = document.querySelectorAll('.otp-input');
                        (inputs[i - 1] as HTMLInputElement)?.focus();
                      }
                    }}
                  />
                ))}
              </div>

              {error && <div className={styles.error}>{error}</div>}

              <Button
                variant="primary"
                size="lg"
                fullWidth
                onClick={handleOtpVerify}
                isLoading={isLoading}
              >
                Verify OTP
              </Button>

              <button
                className={styles.resendBtn}
                onClick={handleContinue}
                disabled={isLoading}
              >
                Resend Code
              </button>
            </div>
          </>
        )}

        <div className={styles.divider}>
          <span>or continue with</span>
        </div>

        <div className={styles.socialButtons}>
          <button className={styles.socialBtn} type="button">
            <span>🔵</span> Google
          </button>
          <button className={styles.socialBtn} type="button">
            <span>🍎</span> Apple
          </button>
          <button className={styles.socialBtn} type="button">
            <span>✈️</span> Telegram
          </button>
        </div>

        <div className={styles.divider}></div>

        <p className={styles.footer}>
          Don't have an account?{' '}
          <Link href="/auth/register" className={styles.link}>
            Sign up
          </Link>
        </p>

        <Link href="/" className={styles.backLink}>
          ← Back to Home
        </Link>
      </Card>
    </div>
  );
}
