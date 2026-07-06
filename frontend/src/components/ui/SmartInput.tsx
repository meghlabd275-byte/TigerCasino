'use client';

import React, { useState, useEffect, useRef, useCallback } from 'react';
import styles from './SmartInput.module.css';

interface SmartInputProps {
  value: string;
  onChange: (value: string, type: 'email' | 'phone') => void;
  placeholder?: string;
  autoFocus?: boolean;
  disabled?: boolean;
  error?: string;
}

// Country data for phone input
const countries = [
  { code: 'US', name: 'United States', dial: '+1', flag: '🇺🇸' },
  { code: 'GB', name: 'United Kingdom', dial: '+44', flag: '🇬🇧' },
  { code: 'CA', name: 'Canada', dial: '+1', flag: '🇨🇦' },
  { code: 'AU', name: 'Australia', dial: '+61', flag: '🇦🇺' },
  { code: 'DE', name: 'Germany', dial: '+49', flag: '🇩🇪' },
  { code: 'FR', name: 'France', dial: '+33', flag: '🇫🇷' },
  { code: 'JP', name: 'Japan', dial: '+81', flag: '🇯🇵' },
  { code: 'CN', name: 'China', dial: '+86', flag: '🇨🇳' },
  { code: 'IN', name: 'India', dial: '+91', flag: '🇮🇳' },
  { code: 'BR', name: 'Brazil', dial: '+55', flag: '🇧🇷' },
  { code: 'MX', name: 'Mexico', dial: '+52', flag: '🇲🇽' },
  { code: 'ES', name: 'Spain', dial: '+34', flag: '🇪🇸' },
  { code: 'IT', name: 'Italy', dial: '+39', flag: '🇮🇹' },
  { code: 'KR', name: 'South Korea', dial: '+82', flag: '🇰🇷' },
  { code: 'ID', name: 'Indonesia', dial: '+62', flag: '🇮🇩' },
  { code: 'TR', name: 'Turkey', dial: '+90', flag: '🇹🇷' },
  { code: 'SA', name: 'Saudi Arabia', dial: '+966', flag: '🇸🇦' },
  { code: 'AE', name: 'UAE', dial: '+971', flag: '🇦🇪' },
  { code: 'TH', name: 'Thailand', dial: '+66', flag: '🇹🇭' },
  { code: 'VN', name: 'Vietnam', dial: '+84', flag: '🇻🇳' },
  { code: 'PH', name: 'Philippines', dial: '+63', flag: '🇵🇭' },
  { code: 'MY', name: 'Malaysia', dial: '+60', flag: '🇲🇾' },
  { code: 'SG', name: 'Singapore', dial: '+65', flag: '🇸🇬' },
  { code: 'PK', name: 'Pakistan', dial: '+92', flag: '🇵🇰' },
  { code: 'BD', name: 'Bangladesh', dial: '+880', flag: '🇧🇩' },
  { code: 'NG', name: 'Nigeria', dial: '+234', flag: '🇳🇬' },
  { code: 'EG', name: 'Egypt', dial: '+20', flag: '🇪🇬' },
  { code: 'ZA', name: 'South Africa', dial: '+27', flag: '🇿🇦' },
  { code: 'KE', name: 'Kenya', dial: '+254', flag: '🇰🇪' },
  { code: 'RU', name: 'Russia', dial: '+7', flag: '🇷🇺' },
  { code: 'UA', name: 'Ukraine', dial: '+380', flag: '🇺🇦' },
  { code: 'PL', name: 'Poland', dial: '+48', flag: '🇵🇱' },
  { code: 'NL', name: 'Netherlands', dial: '+31', flag: '🇳🇱' },
  { code: 'BE', name: 'Belgium', dial: '+32', flag: '🇧🇪' },
  { code: 'SE', name: 'Sweden', dial: '+46', flag: '🇸🇪' },
  { code: 'NO', name: 'Norway', dial: '+47', flag: '🇳🇴' },
  { code: 'DK', name: 'Denmark', dial: '+45', flag: '🇩🇰' },
  { code: 'FI', name: 'Finland', dial: '+358', flag: '🇫🇮' },
  { code: 'CH', name: 'Switzerland', dial: '+41', flag: '🇨🇭' },
  { code: 'AT', name: 'Austria', dial: '+43', flag: '🇦🇹' },
  { code: 'PT', name: 'Portugal', dial: '+351', flag: '🇵🇹' },
  { code: 'GR', name: 'Greece', dial: '+30', flag: '🇬🇷' },
  { code: 'IE', name: 'Ireland', dial: '+353', flag: '🇮🇪' },
  { code: 'NZ', name: 'New Zealand', dial: '+64', flag: '🇳🇿' },
  { code: 'AR', name: 'Argentina', dial: '+54', flag: '🇦🇷' },
  { code: 'CO', name: 'Colombia', dial: '+57', flag: '🇨🇴' },
  { code: 'CL', name: 'Chile', dial: '+56', flag: '🇨🇱' },
  { code: 'PE', name: 'Peru', dial: '+51', flag: '🇵🇪' },
  { code: 'VE', name: 'Venezuela', dial: '+58', flag: '🇻🇪' },
];

// Email regex pattern
const EMAIL_REGEX = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
// Phone regex - starts with + or digits only, minimum 7 digits
const PHONE_REGEX = /^[\d\s\-\+\(\)]{7,}$/;

export default function SmartInput({
  value,
  onChange,
  placeholder = 'Enter email or phone number',
  autoFocus = false,
  disabled = false,
  error
}: SmartInputProps) {
  const [inputType, setInputType] = useState<'email' | 'phone' | null>(null);
  const [showCountryDropdown, setShowCountryDropdown] = useState(false);
  const [selectedCountry, setSelectedCountry] = useState(countries[0]);
  const [searchQuery, setSearchQuery] = useState('');
  const [isValid, setIsValid] = useState<boolean | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Detect input type based on value
  const detectType = useCallback((val: string): 'email' | 'phone' | null => {
    if (!val || val.trim() === '') return null;
    
    // If contains @, it's email
    if (val.includes('@')) return 'email';
    
    // Check if it's a phone number (starts with + or all digits)
    const cleanVal = val.replace(/[\s\-\(\)]/g, '');
    if (/^\+?\d+$/.test(cleanVal) && cleanVal.length >= 7) return 'phone';
    
    // Check if it looks like email (has letters before @)
    if (EMAIL_REGEX.test(val)) return 'email';
    
    return null;
  }, []);

  // Validate input
  const validateInput = useCallback((val: string, type: 'email' | 'phone'): boolean => {
    if (type === 'email') {
      return EMAIL_REGEX.test(val);
    } else {
      const cleanVal = val.replace(/[\s\-\(\)]/g, '');
      return /^\+?\d+$/.test(cleanVal) && cleanVal.length >= 7;
    }
  }, []);

  // Handle input change
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    const detectedType = detectType(newValue);
    
    setInputType(detectedType);
    
    if (detectedType) {
      const valid = validateInput(newValue, detectedType);
      setIsValid(valid);
    } else {
      setIsValid(null);
    }
    
    onChange(newValue, detectedType || 'email');
  };

  // Handle country selection for phone
  const handleCountrySelect = (country: typeof countries[0]) => {
    setSelectedCountry(country);
    setShowCountryDropdown(false);
    setSearchQuery('');
    
    // Update the input value with new country code if empty
    if (!value) {
      onChange(country.dial, 'phone');
    }
  };

  // Filter countries based on search
  const filteredCountries = countries.filter(country => 
    country.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    country.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
    country.dial.includes(searchQuery)
  );

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setShowCountryDropdown(false);
      }
    };
    
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Show/hide country selector based on input type
  useEffect(() => {
    if (inputType === 'phone') {
      // Keep country selector visible
    } else if (inputType === 'email') {
      setShowCountryDropdown(false);
    }
  }, [inputType]);

  return (
    <div className={styles.container}>
      <div className={`${styles.inputWrapper} ${error ? styles.error : ''} ${isValid === true ? styles.valid : ''} ${inputType ? styles[`type-${inputType}`] : ''}`}>
        {inputType === 'phone' && (
          <div className={styles.countrySelector} ref={dropdownRef}>
            <button
              type="button"
              className={styles.countryButton}
              onClick={() => setShowCountryDropdown(!showCountryDropdown)}
              disabled={disabled}
            >
              <span className={styles.flag}>{selectedCountry.flag}</span>
              <span className={styles.dialCode}>{selectedCountry.dial}</span>
              <span className={styles.arrow}>▼</span>
            </button>
            
            {showCountryDropdown && (
              <div className={styles.countryDropdown}>
                <input
                  type="text"
                  placeholder="Search country..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className={styles.countrySearch}
                  autoFocus
                />
                <div className={styles.countryList}>
                  {filteredCountries.map(country => (
                    <button
                      key={country.code}
                      type="button"
                      className={`${styles.countryOption} ${country.code === selectedCountry.code ? styles.selected : ''}`}
                      onClick={() => handleCountrySelect(country)}
                    >
                      <span className={styles.flag}>{country.flag}</span>
                      <span className={styles.countryName}>{country.name}</span>
                      <span className={styles.dialCode}>{country.dial}</span>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
        
        <input
          ref={inputRef}
          type={inputType === 'email' ? 'email' : 'tel'}
          value={value}
          onChange={handleChange}
          placeholder={placeholder}
          disabled={disabled}
          autoFocus={autoFocus}
          className={styles.input}
          autoComplete="off"
        />
        
        {inputType && (
          <div className={styles.typeIndicator}>
            {inputType === 'email' ? (
              <span className={styles.emailIcon}>✉️</span>
            ) : (
              <span className={styles.phoneIcon}>📱</span>
            )}
          </div>
        )}
        
        {isValid === true && (
          <div className={styles.validIndicator}>✓</div>
        )}
      </div>
      
      {error && <div className={styles.errorMessage}>{error}</div>}
      
      {!error && isValid === false && value.length > 0 && (
        <div className={styles.hint}>
          {inputType === 'email' ? 'Please enter a valid email address' : 'Please enter a valid phone number'}
        </div>
      )}
      
      {inputType && !error && (
        <div className={styles.typeLabel}>
          {inputType === 'email' ? '📧 Email detected' : '📱 Phone detected'}
        </div>
      )}
    </div>
  );
}
