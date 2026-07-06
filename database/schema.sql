-- TigerCasino Database Schema
-- PostgreSQL 15+

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    two_fa_secret VARCHAR(255),
    is_2fa_enabled BOOLEAN DEFAULT false,
    email_verified BOOLEAN DEFAULT false,
    phone_verified BOOLEAN DEFAULT false,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    currency VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    tx_hash VARCHAR(100),
    address VARCHAR(100),
    fee DECIMAL(20, 8) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);

-- Games table
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    provider VARCHAR(50),
    rtp DECIMAL(5, 2),
    min_bet DECIMAL(20, 8),
    max_bet DECIMAL(20, 8),
    is_active BOOLEAN DEFAULT true,
    thumbnail_url TEXT,
    game_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bets table
CREATE TABLE IF NOT EXISTS bets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID REFERENCES games(id),
    bet_amount DECIMAL(20, 8) NOT NULL,
    win_amount DECIMAL(20, 8) DEFAULT 0,
    multiplier DECIMAL(10, 2) DEFAULT 0,
    game_data JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    settled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crypto Networks table (e.g., Ethereum, BSC, Tron, etc.)
CREATE TABLE IF NOT EXISTS crypto_networks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    chain_id VARCHAR(20),
    symbol VARCHAR(20) NOT NULL,
    explorer_url VARCHAR(255),
    rpc_url TEXT,
    is_withdrawal_enabled BOOLEAN DEFAULT true,
    is_deposit_enabled BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    min_confirmation_blocks INTEGER DEFAULT 6,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, symbol)
);

-- Crypto Assets table (e.g., USDT, BTC, ETH)
CREATE TABLE IF NOT EXISTS crypto_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    decimals INTEGER DEFAULT 18,
    contract_address VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    min_deposit_amount DECIMAL(20, 8) DEFAULT 0,
    min_withdrawal_amount DECIMAL(20, 8) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol)
);

-- Crypto Asset Networks table (many-to-many: asset supports multiple networks)
CREATE TABLE IF NOT EXISTS crypto_asset_networks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES crypto_assets(id) ON DELETE CASCADE,
    network_id UUID REFERENCES crypto_networks(id) ON DELETE CASCADE,
    deposit_enabled BOOLEAN DEFAULT true,
    withdrawal_enabled BOOLEAN DEFAULT true,
    contract_address VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(asset_id, network_id)
);

-- Network Fees table (deposit/withdrawal fees per network per asset)
CREATE TABLE IF NOT EXISTS network_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES crypto_assets(id) ON DELETE CASCADE,
    network_id UUID REFERENCES crypto_networks(id) ON DELETE CASCADE,
    deposit_fee DECIMAL(20, 8) DEFAULT 0,
    withdrawal_fee DECIMAL(20, 8) DEFAULT 0,
    deposit_fee_percent DECIMAL(5, 4) DEFAULT 0,
    withdrawal_fee_percent DECIMAL(5, 4) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(asset_id, network_id)
);

-- Brand/Level table (for fee discounts)
CREATE TABLE IF NOT EXISTS brand_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    level INTEGER NOT NULL,
    deposit_fee_discount_percent DECIMAL(5, 4) DEFAULT 0,
    withdrawal_fee_discount_percent DECIMAL(5, 4) DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(level)
);

-- Crypto Deposits table
CREATE TABLE IF NOT EXISTS crypto_deposits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID REFERENCES crypto_assets(id) ON DELETE CASCADE,
    network_id UUID REFERENCES crypto_networks(id) ON DELETE CASCADE,
    amount DECIMAL(20, 8) NOT NULL,
    fee DECIMAL(20, 8) DEFAULT 0,
    net_amount DECIMAL(20, 8) NOT NULL,
    address VARCHAR(100) NOT NULL,
    tx_hash VARCHAR(100),
    confirmations INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);

-- Crypto Withdrawals table
CREATE TABLE IF NOT EXISTS crypto_withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID REFERENCES crypto_assets(id) ON DELETE CASCADE,
    network_id UUID REFERENCES crypto_networks(id) ON DELETE CASCADE,
    amount DECIMAL(20, 8) NOT NULL,
    fee DECIMAL(20, 8) DEFAULT 0,
    net_amount DECIMAL(20, 8) NOT NULL,
    address VARCHAR(100) NOT NULL,
    tx_hash VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);

-- Admin User Wallet Addresses (for each asset/network combination)
CREATE TABLE IF NOT EXISTS admin_wallet_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES crypto_assets(id) ON DELETE CASCADE,
    network_id UUID REFERENCES crypto_networks(id) ON DELETE CASCADE,
    address VARCHAR(100) NOT NULL,
    private_key_encrypted TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(asset_id, network_id)
);

-- Insert default brand levels
INSERT INTO brand_levels (name, level, deposit_fee_discount_percent, withdrawal_fee_discount_percent) VALUES
    ('Bronze', 1, 0, 0),
    ('Silver', 2, 5.00, 5.00),
    ('Gold', 3, 10.00, 10.00),
    ('Platinum', 4, 15.00, 15.00),
    ('Diamond', 5, 20.00, 20.00),
    ('White', 6, 20.00, 20.00) ON CONFLICT (level) DO NOTHING;

-- Insert default crypto networks
INSERT INTO crypto_networks (name, chain_id, symbol, explorer_url, min_confirmation_blocks) VALUES
    ('Ethereum', '1', 'ETH', 'https://etherscan.io', 12),
    ('Binance Smart Chain', '56', 'BNB', 'https://bscscan.com', 15),
    ('Tron', '0x', 'TRX', 'https://tronscan.org', 19),
    ('Polygon', '137', 'MATIC', 'https://polygonscan.com', 15),
    ('Avalanche', '43114', 'AVAX', 'https://snowtrace.io', 12),
    ('Solana', 'sol', 'SOL', 'https://solscan.io', 15),
    ('Arbitrum', '42161', 'ETH', 'https://arbiscan.io', 15),
    ('Optimism', '10', 'ETH', 'https://optimistic.etherscan.io', 15)
ON CONFLICT (name, symbol) DO NOTHING;

-- Insert top 50 crypto assets
INSERT INTO crypto_assets (name, symbol, decimals, min_deposit_amount, min_withdrawal_amount) VALUES
    ('Tether USD', 'USDT', 6, 10, 10),
    ('Bitcoin', 'BTC', 8, 0.0001, 0.0002),
    ('Ethereum', 'ETH', 18, 0.001, 0.002),
    ('Binance Coin', 'BNB', 18, 0.001, 0.002),
    ('USD Coin', 'USDC', 6, 10, 10),
    ('Tron', 'TRX', 6, 10, 10),
    ('Polygon', 'MATIC', 18, 1, 1),
    ('Avalanche', 'AVAX', 18, 0.1, 0.1),
    ('Solana', 'SOL', 9, 0.1, 0.1),
    ('Litecoin', 'LTC', 8, 0.01, 0.02),
    ('Ripple', 'XRP', 6, 10, 10),
    ('Dogecoin', 'DOGE', 8, 50, 50),
    ('Chainlink', 'LINK', 18, 0.5, 0.5),
    ('Polkadot', 'DOT', 18, 0.1, 0.1),
    ('Uniswap', 'UNI', 18, 0.1, 0.1),
    ('Cosmos', 'ATOM', 18, 0.1, 0.1),
    ('Stellar', 'XLM', 7, 10, 10),
    ('Monero', 'XMR', 12, 0.001, 0.002),
    ('Ethereum Classic', 'ETC', 18, 0.01, 0.02),
    ('Aave', 'AAVE', 18, 0.01, 0.01),
    ('Maker', 'MKR', 18, 0.001, 0.001),
    ('Compound', 'COMP', 18, 0.01, 0.01),
    ('Synthetix', 'SNX', 18, 0.1, 0.1),
    ('Curve DAO', 'CRV', 18, 1, 1),
    ('SushiSwap', 'SUSHI', 18, 0.1, 0.1),
    ('1inch', '1INCH', 18, 1, 1),
    ('Axie Infinity', 'AXS', 18, 0.1, 0.1),
    ('The Sandbox', 'SAND', 18, 10, 10),
    ('Decentraland', 'MANA', 18, 10, 10),
    ('Enjin Coin', 'ENJ', 18, 1, 1),
    ('Flow', 'FLOW', 8, 1, 1),
    ('Hedera', 'HBAR', 18, 10, 10),
    ('Algorand', 'ALGO', 6, 1, 1),
    ('VeChain', 'VET', 18, 100, 100),
    ('IOTA', 'MIOTA', 6, 1, 1),
    ('Quant', 'QNT', 18, 0.01, 0.01),
    ('Fantom', 'FTM', 18, 10, 10),
    ('Near Protocol', 'NEAR', 24, 0.1, 0.1),
    ('Aptos', 'APT', 8, 0.1, 0.1),
    ('Arweave', 'AR', 12, 0.1, 0.1),
    ('Filecoin', 'FIL', 18, 0.1, 0.1),
    ('THORChain', 'RUNE', 8, 0.1, 0.1),
    ('Osmosis', 'OSMO', 6, 0.1, 0.1),
    ('Kava', 'KAVA', 18, 0.1, 0.1),
    ('Celestia', 'TIA', 6, 0.1, 0.1),
    ('Render', 'RNDR', 18, 0.1, 0.1),
    ('Injective', 'INJ', 18, 0.1, 0.1),
    ('Sei', 'SEI', 6, 0.1, 0.1),
    ('Sui', 'SUI', 9, 0.1, 0.1),
    ('Pepe', 'PEPE', 18, 1000000, 1000000)
ON CONFLICT (symbol) DO NOTHING;

-- Link USDT to all networks
INSERT INTO crypto_asset_networks (asset_id, network_id, deposit_enabled, withdrawal_enabled, contract_address)
SELECT ca.id, cn.id, true, true,
    CASE 
        WHEN cn.name = 'Ethereum' THEN '0xdAC17F958D2ee523a2206206994597C13D831ec7'
        WHEN cn.name = 'Binance Smart Chain' THEN '0x55d398326f99059fF775485246999027B3197955'
        WHEN cn.name = 'Tron' THEN 'TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t'
        WHEN cn.name = 'Polygon' THEN '0xc2132D05D31c914a87C6611C10748AEb04B58e8F'
        WHEN cn.name = 'Avalanche' THEN '0x970979025d0e5eA7D4Bc5a4a4F7B2e5d7eF6F8c9'
        WHEN cn.name = 'Solana' THEN 'Es9vMFrzaCERmfrfcEmMsqJCxEs1bR7cDHEEW1BAxhmG'
        WHEN cn.name = 'Arbitrum' THEN '0xFd086bC7CD5D481aC85e1a55c405B4Cc5B1eB123'
        WHEN cn.name = 'Optimism' THEN '0x94b008aA00579c1307B0EF2c499aD98a8ce58e58'
    END
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol = 'USDT';

-- Link BTC to all networks
INSERT INTO crypto_asset_networks (asset_id, network_id, deposit_enabled, withdrawal_enabled, contract_address)
SELECT ca.id, cn.id, true, true,
    CASE 
        WHEN cn.name = 'Bitcoin' THEN ''
        WHEN cn.name = 'Ethereum' THEN '0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599'
        WHEN cn.name = 'Binance Smart Chain' THEN '0x7130d2A12B9BCbDAe0978D17ce2708b7Fc1d3F1E'
    END
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol = 'BTC' AND cn.name IN ('Bitcoin', 'Ethereum', 'Binance Smart Chain');

-- Link ETH to Ethereum-compatible networks
INSERT INTO crypto_asset_networks (asset_id, network_id, deposit_enabled, withdrawal_enabled, contract_address)
SELECT ca.id, cn.id, true, true, ''
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol = 'ETH' AND cn.name IN ('Ethereum', 'Binance Smart Chain', 'Polygon', 'Avalanche', 'Arbitrum', 'Optimism');

-- Link other major assets to relevant networks
INSERT INTO crypto_asset_networks (asset_id, network_id, deposit_enabled, withdrawal_enabled, contract_address)
SELECT ca.id, cn.id, true, true, ''
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol IN ('BNB', 'USDC', 'MATIC', 'AVAX', 'SOL', 'LINK', 'DOT', 'UNI')
AND (
    (ca.symbol = 'BNB' AND cn.name = 'Binance Smart Chain')
    OR (ca.symbol = 'USDC' AND cn.name IN ('Ethereum', 'Binance Smart Chain', 'Polygon', 'Avalanche', 'Arbitrum', 'Optimism'))
    OR (ca.symbol = 'MATIC' AND cn.name = 'Polygon')
    OR (ca.symbol = 'AVAX' AND cn.name = 'Avalanche')
    OR (ca.symbol = 'SOL' AND cn.name = 'Solana')
    OR (ca.symbol IN ('LINK', 'DOT', 'UNI') AND cn.name IN ('Ethereum', 'Binance Smart Chain'))
);

-- Insert network fees for USDT (default 1 USDT fee for withdrawals)
INSERT INTO network_fees (asset_id, network_id, withdrawal_fee)
SELECT ca.id, cn.id, 1
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol = 'USDT';

-- Insert network fees for BTC (default 0.0005 BTC fee for withdrawals)
INSERT INTO network_fees (asset_id, network_id, withdrawal_fee)
SELECT ca.id, cn.id, 0.0005
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol = 'BTC' AND cn.name = 'Bitcoin';

-- Insert network fees for ETH (default 0.002 ETH fee for withdrawals)
INSERT INTO network_fees (asset_id, network_id, withdrawal_fee)
SELECT ca.id, cn.id, 0.002
FROM crypto_assets ca
CROSS JOIN crypto_networks cn
WHERE ca.symbol = 'ETH';
