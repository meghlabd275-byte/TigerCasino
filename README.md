# TigerCasino - Crypto Deposit & Withdrawal System

A full-stack online casino platform with comprehensive cryptocurrency support for deposits and withdrawals across multiple blockchain networks.

## 🚀 Features

### Multi-Blockchain Crypto Support
- **50+ Popular Cryptocurrencies**: USDT, BTC, ETH, BNB, USDC, TRX, MATIC, AVAX, SOL, and many more
- **Multiple Networks per Asset**: Each cryptocurrency supports multiple blockchain networks (e.g., USDT on Ethereum, BSC, Tron, Polygon, Avalanche, Solana, Arbitrum, Optimism)
- **Flexible Network Management**: Admin can enable/disable deposit and withdrawal for any asset on any network

### Admin Features
- **Add New Crypto Assets**: Easily add any new cryptocurrency to the platform
- **Add New Blockchain Networks**: Support any new blockchain network
- **Link Assets to Networks**: Connect any asset to any network
- **Fee Management**: Configure deposit and withdrawal fees for each asset-network combination
- **Brand/Level System**: Manage VIP levels with customizable fee discounts
- **White Level Special Discount**: Configure up to 20% fee discount for White level brands
- **Wallet Management**: Configure admin wallet addresses for each asset-network
- **Transaction Management**: View, confirm, or cancel deposits and withdrawals

### User Features
- **Multi-chain Deposits**: Deposit crypto using any supported network
- **Multi-chain Withdrawals**: Withdraw crypto to any supported network address
- **Real-time Fee Calculation**: See fees before confirming transactions
- **Brand Level Benefits**: Higher levels get fee discounts on deposits and withdrawals
- **Transaction History**: View complete deposit and withdrawal history

## 📋 Supported Cryptocurrencies

### Top 50 Crypto Assets
1. USDT (Tether)
2. BTC (Bitcoin)
3. ETH (Ethereum)
4. BNB (Binance Coin)
5. USDC (USD Coin)
6. TRX (Tron)
7. MATIC (Polygon)
8. AVAX (Avalanche)
9. SOL (Solana)
10. LTC (Litecoin)
11. XRP (Ripple)
12. DOGE (Dogecoin)
13. LINK (Chainlink)
14. DOT (Polkadot)
15. UNI (Uniswap)
16. ATOM (Cosmos)
17. XLM (Stellar)
18. XMR (Monero)
19. ETC (Ethereum Classic)
20. AAVE (Aave)
21. MKR (Maker)
22. COMP (Compound)
23. SNX (Synthetix)
24. CRV (Curve DAO)
25. SUSHI (SushiSwap)
26. 1INCH (1inch)
27. AXS (Axie Infinity)
28. SAND (The Sandbox)
29. MANA (Decentraland)
30. ENJ (Enjin Coin)
31. FLOW (Flow)
32. HBAR (Hedera)
33. ALGO (Algorand)
34. VET (VeChain)
35. MIOTA (IOTA)
36. QNT (Quant)
37. FTM (Fantom)
38. NEAR (Near Protocol)
39. APT (Aptos)
40. AR (Arweave)
41. FIL (Filecoin)
42. RUNE (THORChain)
43. OSMO (Osmosis)
44. KAVA (Kava)
45. TIA (Celestia)
46. RNDR (Render)
47. INJ (Injective)
48. SEI (Sei)
49. SUI (Sui)
50. PEPE (Pepe)

### Supported Blockchain Networks
1. **Ethereum** (ETH)
2. **Binance Smart Chain** (BNB)
3. **Tron** (TRX)
4. **Polygon** (MATIC)
5. **Avalanche** (AVAX)
6. **Solana** (SOL)
7. **Arbitrum** (ETH)
8. **Optimism** (ETH)

## 💰 Brand Levels & Fee Discounts

| Level | Name | Deposit Discount | Withdrawal Discount |
|-------|------|------------------|---------------------|
| 1 | Bronze | 0% | 0% |
| 2 | Silver | 5% | 5% |
| 3 | Gold | 10% | 10% |
| 4 | Platinum | 15% | 15% |
| 5 | Diamond | 20% | 20% |
| 6 | **White** | **20%** | **20%** |

Admin can update any brand level discount percentages. The White level is pre-configured with 20% fee discount as requested.

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 15+
- **ORM**: GORM

### Frontend
- **Framework**: Next.js
- **Language**: TypeScript

## 📦 Installation

### Prerequisites

- **Go** 1.21 or later
- **Node.js** 18 or later
- **PostgreSQL** 15 or later
- **Docker** & Docker Compose (optional)

### Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE tigercasino;
```

2. Run the schema file:
```bash
psql -U postgres -d tigercasino -f database/schema.sql
```

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create configuration file:
```bash
cp config.yaml.example config.yaml
```

4. Update `config.yaml` with your settings:
```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  name: tigercasino

app:
  host: 0.0.0.0
  port: 8080
  jwt_secret: your_jwt_secret_key
```

5. Run the backend:
```bash
go run cmd/server/main.go
```

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create environment file:
```bash
cp .env.example .env.local
```

4. Update `.env.local`:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

5. Run the development server:
```bash
npm run dev
```

### Docker Setup (Alternative)

1. Navigate to the project root:
```bash
cd deployments
```

2. Start all services:
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database
- Backend API server
- Frontend Next.js application
- Nginx reverse proxy

## 🔧 API Endpoints

### User Endpoints

| Method | Endpoint | Description |
|--------|----------|--------------|
| GET | `/api/crypto/networks` | Get all available networks |
| GET | `/api/crypto/assets` | Get all assets with network info |
| POST | `/api/crypto/deposit/address` | Get deposit address |
| GET | `/api/crypto/deposits` | Get user deposits |
| POST | `/api/crypto/withdrawals` | Create withdrawal |
| GET | `/api/crypto/withdrawals` | Get user withdrawals |
| POST | `/api/crypto/fees` | Calculate fees |

### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|--------------|
| GET/POST | `/api/admin/crypto/networks` | List/Create networks |
| PUT | `/api/admin/crypto/networks/:id` | Update network |
| GET/POST | `/api/admin/crypto/assets` | List/Create assets |
| PUT | `/api/admin/crypto/assets/:id` | Update asset |
| POST | `/api/admin/crypto/asset-networks` | Link asset to network |
| PUT | `/api/admin/crypto/asset-networks/:id` | Update asset-network |
| POST | `/api/admin/crypto/fees` | Set network fees |
| GET/POST | `/api/admin/crypto/brand-levels` | List/Create brand levels |
| PUT | `/api/admin/crypto/brand-levels/:id` | Update brand level |
| POST | `/api/admin/crypto/wallets` | Set admin wallet |
| GET | `/api/admin/crypto/deposits` | List all deposits |
| POST | `/api/admin/crypto/deposits/:id/confirm` | Confirm deposit |
| GET | `/api/admin/crypto/withdrawals` | List all withdrawals |
| POST | `/api/admin/crypto/withdrawals/:id/confirm` | Confirm withdrawal |
| POST | `/api/admin/crypto/withdrawals/:id/cancel` | Cancel withdrawal |
| PUT | `/api/admin/users/:id/brand-level` | Update user brand level |

## 📝 Configuration Examples

### Adding a New Crypto Asset

```bash
curl -X POST http://localhost:8080/api/admin/crypto/assets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "Bitcoin",
    "symbol": "BTC",
    "decimals": 8,
    "min_deposit_amount": 0.0001,
    "min_withdrawal_amount": 0.0002
  }'
```

### Adding a New Network

```bash
curl -X POST http://localhost:8080/api/admin/crypto/networks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "Base",
    "chain_id": "8453",
    "symbol": "ETH",
    "explorer_url": "https://basescan.org",
    "min_confirmation_blocks": 12
  }'
```

### Linking Asset to Network

```bash
curl -X POST http://localhost:8080/api/admin/crypto/asset-networks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "asset_id": "ASSET_UUID",
    "network_id": "NETWORK_UUID",
    "deposit_enabled": true,
    "withdrawal_enabled": true,
    "contract_address": "0x..."
  }'
```

### Setting Network Fees

```bash
curl -X POST http://localhost:8080/api/admin/crypto/fees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "asset_id": "ASSET_UUID",
    "network_id": "NETWORK_UUID",
    "withdrawal_fee": 1.0,
    "withdrawal_fee_percent": 0.5
  }'
```

### Updating White Level to 20% Discount

```bash
curl -X PUT http://localhost:8080/api/admin/crypto/brand-levels/WHITE_LEVEL_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "White",
    "level": 6,
    "deposit_fee_discount_percent": 20.0,
    "withdrawal_fee_discount_percent": 20.0,
    "is_active": true
  }'
```

### Setting Admin Wallet

```bash
curl -X POST http://localhost:8080/api/admin/crypto/wallets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "asset_id": "USDT_ASSET_UUID",
    "network_id": "ETHEREUM_NETWORK_UUID",
    "address": "0xYourWalletAddress",
    "is_active": true
  }'
```

## 🔒 Security Considerations

1. **Private Key Storage**: Admin wallet private keys should be encrypted at rest
2. **Web3 Integration**: In production, use dedicated wallet services or MPC wallets
3. **Confirmation Requirements**: Configure minimum confirmation blocks per network
4. **Rate Limiting**: Implement rate limiting on withdrawal endpoints
5. **KYC Integration**: Consider integrating KYC for large transactions
6. **Audit Logging**: All admin actions are logged in the audit_logs table

## 📊 Database Schema

### Key Tables

- `crypto_networks` - Blockchain network definitions
- `crypto_assets` - Cryptocurrency definitions
- `crypto_asset_networks` - Asset-network relationships
- `network_fees` - Deposit/withdrawal fees per asset-network
- `brand_levels` - VIP/brand levels with fee discounts
- `crypto_deposits` - User deposit records
- `crypto_withdrawals` - User withdrawal records
- `admin_wallet_addresses` - Admin wallet configurations

## 🚀 Deployment

### Production Checklist

- [ ] Configure SSL/TLS certificates
- [ ] Set up environment variables for secrets
- [ ] Configure database connection pooling
- [ ] Set up monitoring and logging
- [ ] Configure backup strategy
- [ ] Set up Redis for caching (optional)
- [ ] Configure load balancer

### Environment Variables

```bash
# Backend
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=tigercasino
JWT_SECRET=your_jwt_secret
PORT=8080

# Frontend
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

## 📄 License

This project is proprietary software. All rights reserved.