# TigerCasino - Crypto Casino Platform Specification

## 1. Project Overview

**Project Name:** TigerCasino  
**Project Type:** Crypto Casino Web Application  
**Core Functionality:** A complete online cryptocurrency casino platform with multiple games, user wallet management, admin dashboard, and real-time betting capabilities  
**Target Users:** Cryptocurrency enthusiasts who want to play casino games with digital assets

---

## 2. Technology Stack

### Frontend
- **Framework:** Next.js 14 with App Router
- **Language:** TypeScript
- **UI Library:** React 18
- **Styling:** CSS Modules with custom properties
- **State Management:** React Context + useReducer
- **HTTP Client:** Native fetch with custom hooks

### Backend
- **Primary Language:** Go (Golang) - Main backend API
- **Database:** PostgreSQL 15
- **Cache:** Redis for session management
- **ORM:** GORM
- **Authentication:** JWT with RSA signatures

### High-Performance Components
- **Rust:** Security module for cryptographic operations, RNG, fraud detection
- **C++:** Ultra-low-latency game engine for real-time game outcomes
- **WebSocket:** For real-time game updates

### Infrastructure
- **Container:** Docker with multi-stage builds
- **Web Server:** Nginx for reverse proxy

---

## 3. UI/UX Specification

### Color Palette
- **Primary:** `#FF6B35` (Tiger Orange)
- **Secondary:** `#1A1A2E` (Deep Navy)
- **Accent:** `#FFD700` (Gold)
- **Background:** `#0F0F1A` (Dark Black)
- **Surface:** `#16213E` (Dark Blue)
- **Text Primary:** `#FFFFFF`
- **Text Secondary:** `#B0B0B0`
- **Success:** `#00D26A`
- **Error:** `#FF4757`
- **Warning:** `#FFA502`

### Typography
- **Heading Font:** 'Orbitron', sans-serif (Futuristic)
- **Body Font:** 'Rajdhani', sans-serif (Tech)
- **Monospace:** 'JetBrains Mono', monospace (Numbers)

### Spacing System
- **Base Unit:** 8px
- **XS:** 4px
- **SM:** 8px
- **MD:** 16px
- **LG:** 24px
- **XL:** 32px
- **XXL:** 48px

### Responsive Breakpoints
- **Mobile:** < 768px
- **Tablet:** 768px - 1024px
- **Desktop:** > 1024px
- **Wide:** > 1440px

### Layout Structure

#### Public Pages
1. **Landing Page**
   - Hero section with animated tiger logo
   - Featured games carousel
   - Trust indicators
   - Call-to-action buttons

2. **Games Page**
   - Grid of available games
   - Categories filter
   - Search functionality
   - Game thumbnails with hover effects

3. **About Page**
   - Company information
   - Security guarantees
   - Fair play explanation

4. **Contact Page**
   - Contact form
   - Support channels

#### User Dashboard
1. **Main Dashboard**
   - Wallet balance display
   - Recent transactions
   - Active bets
   - Quick play buttons

2. **Wallet Page**
   - Deposit addresses (BTC, ETH, USDT, etc.)
   - Withdrawal form
   - Transaction history
   - Balance charts

3. **Games Lobby**
   - All available games
   - Favorite games
   - Recent games

4. **Profile Page**
   - User information
   - Security settings
   - Two-factor authentication
   - KYC verification

5. **Transaction History**
   - All deposits/withdrawals
   - Bet history
   - Win/loss statistics

#### Admin Dashboard
1. **Overview**
   - Total users, revenue, active games
   - Real-time statistics
   - System health

2. **User Management**
   - User list with search/filter
   - User details view
   - Account actions (ban, verify, etc.)

3. **Game Management**
   - Game settings
   - Odds management
   - Game enable/disable

4. **Financial Management**
   - All transactions
   - Pending withdrawals
   - Manual adjustments

5. **Security Center**
   - Suspicious activities
   - Fraud detection alerts
   - Audit logs

6. **System Settings**
   - General settings
   - API keys management
   - Notification templates

### Visual Effects
- **Animations:** Smooth transitions (0.3s ease)
- **Hover Effects:** Scale + glow on interactive elements
- **Loading States:** Skeleton loaders + spinners
- **Notifications:** Toast messages with slide animation

---

## 4. Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    wallet_address VARCHAR(100),
    wallet_type VARCHAR(20),
    balance DECIMAL(20, 8) DEFAULT 0,
    bonus_balance DECIMAL(20, 8) DEFAULT 0,
    vip_level INTEGER DEFAULT 0,
    kyc_status VARCHAR(20) DEFAULT 'pending',
    is_verified BOOLEAN DEFAULT false,
    is_admin BOOLEAN DEFAULT false,
    is_banned BOOLEAN DEFAULT false,
    ban_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    currency VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    tx_hash VARCHAR(100),
    address VARCHAR(100),
    fee DECIMAL(20, 8),
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP
);
```

### Games Table
```sql
CREATE TABLE games (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    provider VARCHAR(50),
    rtp DECIMAL(5, 2),
    min_bet DECIMAL(20, 8),
    max_bet DECIMAL(20, 8),
    is_active BOOLEAN DEFAULT true,
    thumbnail_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Bets Table
```sql
CREATE TABLE bets (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    game_id UUID REFERENCES games(id),
    bet_amount DECIMAL(20, 8) NOT NULL,
    win_amount DECIMAL(20, 8) DEFAULT 0,
    multiplier DECIMAL(10, 2) DEFAULT 0,
    game_data JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    settled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Sessions Table
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    token TEXT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Audit Logs Table
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 5. Functionality Specification

### User Features

#### Authentication
- Email/password registration with validation
- Secure login with JWT tokens
- Two-factor authentication (TOTP)
- Password reset via email
- Session management

#### Wallet Management
- Multi-currency support (BTC, ETH, USDT, etc.)
- Deposit address generation
- Withdrawal requests
- Transaction history
- Balance auto-refresh

#### Casino Games
1. **Slot Machines**
   - Multiple themes
   - Bonus rounds
   - Progressive jackpots

2. **Dice**
   - Customizable odds
   - Auto-bet feature
   - Strategy mode

3. **Roulette**
   - European/American variants
   - Live dealer option
   - Multiple betting options

4. **Blackjack**
   - Classic rules
   - Multi-hand mode
   - Card counting indicator

5. **Baccarat**
   - Punto Banco variant
   - Side bets
   - Speed mode

6. **Sports Betting**
   - Live events
   - Pre-match betting
   - Multiple sports

#### Account Features
- Profile management
- Security settings
- KYC verification
- VIP program
- Referral system

### Admin Features

#### Dashboard
- Real-time statistics
- Revenue charts
- User activity
- System health

#### User Management
- View all users
- Search/filter users
- Edit user details
- Ban/unban users
- Verify KYC

#### Game Management
- Enable/disable games
- Adjust game settings
- Set RTPs
- Manage providers

#### Financial Management
- View all transactions
- Approve/reject withdrawals
- Manual deposits
- Refund processing

#### Security
- View suspicious activities
- Block suspicious IPs
- Audit logs
- Fraud detection

---

## 6. API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh token
- `POST /api/auth/2fa/setup` - Setup 2FA
- `POST /api/auth/2fa/verify` - Verify 2FA

### Users
- `GET /api/users/me` - Get current user
- `PUT /api/users/me` - Update profile
- `GET /api/users/:id` - Get user (admin)
- `GET /api/users` - List users (admin)

### Wallet
- `GET /api/wallet/balance` - Get balance
- `GET /api/wallet/deposit/address` - Get deposit address
- `POST /api/wallet/withdraw` - Request withdrawal
- `GET /api/wallet/transactions` - Transaction history

### Games
- `GET /api/games` - List games
- `GET /api/games/:id` - Get game details
- `POST /api/games/:id/bet` - Place bet
- `GET /api/games/:id/history` - Bet history

### Admin
- `GET /api/admin/dashboard` - Dashboard stats
- `GET /api/admin/users` - User management
- `PUT /api/admin/users/:id` - Update user
- `GET /api/admin/transactions` - All transactions
- `POST /api/admin/transactions/:id/approve` - Approve withdrawal
- `GET /api/admin/audit-logs` - Audit logs

---

## 7. Security Implementation

### Rust Security Module
- Cryptographic random number generation
- Hash password with Argon2
- Signature verification
- Encryption/decryption
- Fraud pattern detection

### C++ Game Engine
- Seeded RNG with cryptographic entropy
- Game outcome calculation
- Collision detection
- Performance-critical computations

### Security Measures
- HTTPS only
- Rate limiting
- Input validation
- SQL injection prevention
- XSS protection
- CSRF tokens
- Account lockout
- Audit logging

---

## 8. Acceptance Criteria

### Functional Requirements
- [x] Users can register and login
- [x] Users can deposit cryptocurrency
- [x] Users can withdraw funds
- [x] Users can play all casino games
- [x] Bet results are calculated correctly
- [x] Admin can manage all aspects
- [x] Real-time balance updates work
- [x] Transaction history is complete

### Performance Requirements
- [x] Page load time < 2 seconds
- [x] API response time < 100ms
- [x] Game outcome calculation < 10ms
- [x] Support 10,000+ concurrent users

### Security Requirements
- [x] All passwords are hashed
- [x] JWT tokens are properly validated
- [x] 2FA is available
- [x] All inputs are validated
- [x] Audit logs are complete

### Visual Checkpoints
- [x] Dark theme with tiger orange accents
- [x] Smooth animations on interactions
- [x] Responsive on all devices
- [x] Loading states are visible
- [x] Error messages are clear
