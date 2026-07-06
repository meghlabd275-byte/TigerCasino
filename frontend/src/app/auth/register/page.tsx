'use client';

import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Card } from '@/components/ui';
import styles from '../auth.module.css';

type InputMode = 'email' | 'phone';

type RegisterStep = 'identity' | 'otp' | 'password';

export default function RegisterPage() {
  const router = useRouter();
  const { register, isAuthenticated, isLoading: authLoading } = useAuth();
  
  const [step, setStep] = useState<RegisterStep>('identity');
  const [inputValue, setInputValue] = useState('');
  const [inputMode, setInputMode] = useState<InputMode>('email');
  const [otp, setOtp] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [referralCode, setReferralCode] = useState('');
  const [acceptTerms, setAcceptTerms] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [passwordStrength, setPasswordStrength] = useState<'weak' | 'medium' | 'strong'>('weak');

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

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    
    // Calculate password strength
    let strength: 'weak' | 'medium' | 'strong' = 'weak';
    if (value.length >= 8) {
      const hasUpper = /[A-Z]/.test(value);
      const hasLower = /[a-z]/.test(value);
      const hasNumber = /[0-9]/.test(value);
      const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(value);
      const score = [hasUpper, hasLower, hasNumber, hasSpecial].filter(Boolean).length;
      
      if (score >= 3 && value.length >= 12) strength = 'strong';
      else if (score >= 2) strength = 'medium';
    }
    setPasswordStrength(strength);
    setError('');
  };

  const handleContinue = async () => {
    if (!inputValue.trim()) {
      setError('Please enter your email or phone number');
      return;
    }

    setIsLoading(true);
    setError('');
    
    // Simulate API call to check account
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    setStep('otp');
    setSuccess('OTP sent successfully!');
    setIsLoading(false);
  };

  const handleOtpVerify = async () => {
    if (otp.length !== 6) {
      setError('Please enter the 6-digit OTP');
      return;
    }

    setIsLoading(true);
    setError('');
    
    // Simulate OTP verification
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    setStep('password');
    setIsLoading(false);
  };

  const handleRegister = async () => {
    if (!password) {
      setError('Please enter a password');
      return;
    }

    if (password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (!acceptTerms) {
      setError('Please accept the Terms & Conditions');
      return;
    }

    setIsLoading(true);
    setError('');

    const result = await register({
      email: inputValue,
      username: inputValue.split('@')[0],
      password,
      confirmPassword,
    });

    if (result.success) {
      router.push('/dashboard');
    } else {
      setError(result.error || 'Registration failed');
    }

    setIsLoading(false);
  };

  const handleResendOtp = async () => {
    setIsLoading(true);
    await new Promise(resolve => setTimeout(resolve, 1000));
    setSuccess('OTP resent successfully!');
    setIsLoading(false);
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
          <h1 className={styles.title}>Create Account</h1>
          <p className={styles.subtitle}>Join the ultimate crypto casino experience</p>
        </div>

        {step === 'identity' && (
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
                  placeholder={inputMode === 'email' ? 'Enter your email' : 'Enter phone number'}
                  className={styles.input}
                  autoComplete="off"
                />
              </div>
            </div>

            {error && <div className={styles.error}>{error}</div>}
            {success && <div className={styles.success}>{success}</div>}

            <Button
              variant="primary"
              size="lg"
              fullWidth
              onClick={handleContinue}
              isLoading={isLoading}
            >
              Continue
            </Button>
          </>
        )}

        {step === 'otp' && (
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
              {success && <div className={styles.success}>{success}</div>}

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
                onClick={handleResendOtp}
                disabled={isLoading}
              >
                Resend Code
              </button>
            </div>
          </>
        )}

        {step === 'password' && (
          <>
            <div className={styles.formGroup}>
              <label className={styles.label}>Password</label>
              <div className={styles.passwordWrapper}>
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={handlePasswordChange}
                  placeholder="Create a password"
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
              <div className={styles.passwordStrength}>
                <div className={styles.strengthBar}>
                  <div className={`${styles.strengthFill} ${styles[`strength${passwordStrength.charAt(0).toUpperCase() + passwordStrength.slice(1)}`]}`}></div>
                </div>
                <span className={styles.strengthLabel}>
                  Password strength: {passwordStrength}
                </span>
              </div>
            </div>

            <div className={styles.formGroup}>
              <label className={styles.label}>Confirm Password</label>
              <div className={styles.passwordWrapper}>
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="Confirm your password"
                  className={styles.input}
                />
              </div>
            </div>

            <div className={styles.formGroup}>
              <label className={styles.label}>Referral Code (Optional)</label>
              <div className={styles.inputWrapper}>
                <input
                  type="text"
                  value={referralCode}
                  onChange={(e) => setReferralCode(e.target.value)}
                  placeholder="Enter referral code"
                  className={styles.input}
                />
              </div>
            </div>

            <div className={styles.options}>
              <label className={styles.checkbox}>
                <input
                  type="checkbox"
                  checked={acceptTerms}
                  onChange={(e) => setAcceptTerms(e.target.checked)}
                />
                <span>I accept the <Link href="/terms" className={styles.link}>Terms & Conditions</Link></span>
              </label>
            </div>

            {error && <div className={styles.error}>{error}</div>}

            <Button
              variant="primary"
              size="lg"
              fullWidth
              onClick={handleRegister}
              isLoading={isLoading}
            >
              Create Account
            </Button>
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
          Already have an account?{' '}
          <Link href="/auth/login" className={styles.link}>
            Sign in
          </Link>
        </p>

        <Link href="/" className={styles.backLink}>
          ← Back to Home
        </Link>
      </Card>
    </div>
  );
}
