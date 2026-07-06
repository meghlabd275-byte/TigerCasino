'use client';

import React, { useState } from 'react';
import { Header, Footer } from '@/components/layout';
import { Card, Button, Input } from '@/components/ui';
import styles from './wallet.module.css';

interface WalletAsset {
  chain: string;
  name: string;
  symbol: string;
  balance: number;
  usdValue: number;
  icon: string;
  depositEnabled: boolean;
  withdrawEnabled: boolean;
  minDeposit: number;
  minWithdraw: number;
}

const walletAssets: WalletAsset[] = [
  { chain: 'BTC', name: 'Bitcoin', symbol: 'BTC', balance: 1.2456, usdValue: 49824, depositEnabled: true, withdrawEnabled: true, minDeposit: 0.0001, minWithdraw: 0.0002 },
  { chain: 'ETH', name: 'Ethereum', symbol: 'ETH', balance: 5.8923, usdValue: 14730, depositEnabled: true, withdrawEnabled: true, minDeposit: 0.001, minWithdraw: 0.005 },
  { chain: 'USDT', name: 'Tether USD', symbol: 'USDT', balance: 25000, usdValue: 25000, depositEnabled: true, withdrawEnabled: true, minDeposit: 10, minWithdraw: 20 },
  { chain: 'BNB', name: 'BNB Chain', symbol: 'BNB', balance: 12.567, usdValue: 3760, depositEnabled: true, withdrawEnabled: true, minDeposit: 0.01, minWithdraw: 0.05 },
  { chain: 'SOL', name: 'Solana', symbol: 'SOL', balance: 156.78, usdValue: 12542, depositEnabled: true, withdrawEnabled: true, minDeposit: 0.01, minWithdraw: 0.1 },
  { chain: 'TRX', name: 'Tron', symbol: 'TRX', balance: 5000, usdValue: 450, depositEnabled: true, withdrawEnabled: true, minDeposit: 10, minWithdraw: 50 },
  { chain: 'USDC', name: 'USD Coin', symbol: 'USDC', balance: 10000, usdValue: 10000, depositEnabled: true, withdrawEnabled: true, minDeposit: 10, minWithdraw: 20 },
  { chain: 'MATIC', name: 'Polygon', symbol: 'MATIC', balance: 5000, usdValue: 4000, depositEnabled: true, withdrawEnabled: true, minDeposit: 1, minWithdraw: 10 },
];

export default function WalletPage() {
  const [activeTab, setActiveTab] = useState<'assets' | 'deposit' | 'withdraw'>('assets');
  const [selectedChain, setSelectedChain] = useState<string>('BTC');
  const [depositAmount, setDepositAmount] = useState('');
  const [withdrawAmount, setWithdrawAmount] = useState('');
  const [walletAddress, setWalletAddress] = useState('');

  const selectedAsset = walletAssets.find(a => a.chain === selectedChain);
  const totalUSD = walletAssets.reduce((sum, a) => sum + a.usdValue, 0);

  return (
    <>
      <Header />
      <main className={styles.main}>
        <div className={styles.container}>
          <div className={styles.header}>
            <h1 className={styles.title}>Wallet</h1>
            <p className={styles.subtitle}>
              Multi-chain deposit and withdrawal. 11 blockchains supported.
            </p>
          </div>

          {/* Balance Overview */}
          <Card className={styles.balanceCard}>
            <div className={styles.balanceLabel}>Total Balance</div>
            <div className={styles.balanceValue}>${totalUSD.toLocaleString()}</div>
            <div className={styles.balanceActions}>
              <Button onClick={() => setActiveTab('deposit')}>Deposit</Button>
              <Button onClick={() => setActiveTab('withdraw')} variant="secondary">Withdraw</Button>
            </div>
          </Card>

          {/* Tab Navigation */}
          <div className={styles.tabs}>
            <button
              className={`${styles.tab} ${activeTab === 'assets' ? styles.active : ''}`}
              onClick={() => setActiveTab('assets')}
            >
              Assets
            </button>
            <button
              className={`${styles.tab} ${activeTab === 'deposit' ? styles.active : ''}`}
              onClick={() => setActiveTab('deposit')}
            >
              Deposit
            </button>
            <button
              className={`${styles.tab} ${activeTab === 'withdraw' ? styles.active : ''}`}
              onClick={() => setActiveTab('withdraw')}
            >
              Withdraw
            </button>
          </div>

          {/* Assets Tab */}
          {activeTab === 'assets' && (
            <div className={styles.assetGrid}>
              {walletAssets.map(asset => (
                <Card key={asset.chain} className={styles.assetCard}>
                  <div className={styles.assetHeader}>
                    <span className={styles.assetIcon}>{asset.icon}</span>
                    <div className={styles.assetInfo}>
                      <div className={styles.assetName}>{asset.name}</div>
                      <div className={styles.assetSymbol}>{asset.symbol}</div>
                    </div>
                  </div>
                  <div className={styles.assetBalance}>
                    <div className={styles.assetAmount}>{asset.balance.toLocaleString()} {asset.symbol}</div>
                    <div className={styles.assetUSD}>${asset.usdValue.toLocaleString()}</div>
                  </div>
                  <div className={styles.assetActions}>
                    <Button 
                      size="sm" 
                      onClick={() => { setSelectedChain(asset.chain); setActiveTab('deposit'); }}
                      disabled={!asset.depositEnabled}
                    >
                      Deposit
                    </Button>
                    <Button 
                      size="sm" 
                      variant="secondary"
                      onClick={() => { setSelectedChain(asset.chain); setActiveTab('withdraw'); }}
                      disabled={!asset.withdrawEnabled}
                    >
                      Withdraw
                    </Button>
                  </div>
                </Card>
              ))}
            </div>
          )}

          {/* Deposit Tab */}
          {activeTab === 'deposit' && (
            <Card className={styles.actionCard}>
              <h2 className={styles.cardTitle}>Deposit Crypto</h2>
              
              <div className={styles.chainSelector}>
                <label className={styles.label}>Select Blockchain</label>
                <div className={styles.chainGrid}>
                  {walletAssets.map(asset => (
                    <button
                      key={asset.chain}
                      className={`${styles.chainBtn} ${selectedChain === asset.chain ? styles.selected : ''}`}
                      onClick={() => setSelectedChain(asset.chain)}
                    >
                      <span className={styles.chainIcon}>{asset.icon}</span>
                      <span className={styles.chainName}>{asset.chain}</span>
                    </button>
                  ))}
                </div>
              </div>

              {selectedAsset && (
                <>
                  <div className={styles.depositInfo}>
                    <div className={styles.infoRow}>
                      <span>Network</span>
                      <span>{selectedAsset.name}</span>
                    </div>
                    <div className={styles.infoRow}>
                      <span>Minimum Deposit</span>
                      <span>{selectedAsset.minDeposit} {selectedAsset.symbol}</span>
                    </div>
                  </div>

                  <div className={styles.addressBox}>
                    <label className={styles.label}>Deposit Address</label>
                    <div className={styles.addressDisplay}>
                      <code className={styles.address}>
                        0x7a250d5630b4cf539739df2c5bdac4a3a0a2b4
                      </code>
                      <Button 
                        size="sm" 
                        onClick={() => navigator.clipboard.writeText('0x7a250d5630b4cf539739df2c5bdac4a3a0a2b4')}
                      >
                        Copy
                      </Button>
                    </div>
                    <p className={styles.addressHint}>
                      Only send {selectedAsset.symbol} to this address. Other tokens may be lost.
                    </p>
                  </div>

                  <div className={styles.qrCode}>
                    <div className={styles.qrPlaceholder}>QR Code</div>
                  </div>
                </>
              )}
            </Card>
          )}

          {/* Withdraw Tab */}
          {activeTab === 'withdraw' && (
            <Card className={styles.actionCard}>
              <h2 className={styles.cardTitle}>Withdraw Crypto</h2>

              <div className={styles.chainSelector}>
                <label className={styles.label}>Select Blockchain</label>
                <div className={styles.chainGrid}>
                  {walletAssets.map(asset => (
                    <button
                      key={asset.chain}
                      className={`${styles.chainBtn} ${selectedChain === asset.chain ? styles.selected : ''}`}
                      onClick={() => setSelectedChain(asset.chain)}
                    >
                      <span className={styles.chainIcon}>{asset.icon}</span>
                      <span className={styles.chainName}>{asset.chain}</span>
                    </button>
                  ))}
                </div>
              </div>

              {selectedAsset && (
                <>
                  <div className={styles.inputGroup}>
                    <label className={styles.label}>Recipient Address</label>
                    <Input
                      type="text"
                      value={walletAddress}
                      onChange={(e) => setWalletAddress(e.target.value)}
                      placeholder={`Enter ${selectedAsset.symbol} wallet address`}
                    />
                  </div>

                  <div className={styles.inputGroup}>
                    <label className={styles.label}>Amount</label>
                    <Input
                      type="number"
                      value={withdrawAmount}
                      onChange={(e) => setWithdrawAmount(e.target.value)}
                      placeholder={`Min: ${selectedAsset.minWithdraw} ${selectedAsset.symbol}`}
                    />
                    <div className={styles.available}>
                      Available: {selectedAsset.balance} {selectedAsset.symbol}
                    </div>
                  </div>

                  <div className={styles.feeInfo}>
                    <div className={styles.feeRow}>
                      <span>Network Fee</span>
                      <span>0.0001 {selectedAsset.symbol}</span>
                    </div>
                    <div className={styles.feeRow}>
                      <span>Platform Fee (0%)</span>
                      <span>0 {selectedAsset.symbol}</span>
                    </div>
                  </div>

                  <Button className={styles.submitBtn} disabled={!walletAddress || !withdrawAmount}>
                    Withdraw {withdrawAmount || 0} {selectedAsset.symbol}
                  </Button>
                </>
              )}
            </Card>
          )}
        </div>
      </main>
      <Footer />
    </>
  );
}
