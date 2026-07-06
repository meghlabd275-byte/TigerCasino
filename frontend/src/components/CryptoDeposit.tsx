'use client';

import { useState, useEffect } from 'react';
import { getAssets, getDepositAddress, calculateFees } from '@/lib/crypto';
import { CryptoAssetWithNetwork } from '@/types/crypto';

interface CryptoDepositProps {
  onSuccess?: () => void;
}

export default function CryptoDeposit({ onSuccess }: CryptoDepositProps) {
  const [assets, setAssets] = useState<CryptoAssetWithNetwork[]>([]);
  const [selectedAsset, setSelectedAsset] = useState<CryptoAssetWithNetwork | null>(null);
  const [depositAddress, setDepositAddress] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    loadAssets();
  }, []);

  const loadAssets = async () => {
    try {
      const data = await getAssets();
      setAssets(data);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleAssetSelect = async (asset: CryptoAssetWithNetwork) => {
    setSelectedAsset(asset);
    setDepositAddress('');
    setError('');
    
    if (asset.deposit_enabled) {
      try {
        setLoading(true);
        const { address } = await getDepositAddress(asset.id, asset.network_id);
        setDepositAddress(address);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
  };

  const copyToClipboard = async () => {
    if (depositAddress) {
      await navigator.clipboard.writeText(depositAddress);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="crypto-deposit">
      <h2>Deposit Crypto</h2>
      
      {error && <div className="error-message">{error}</div>}

      <div className="asset-grid">
        {assets.filter(a => a.deposit_enabled).map((asset) => (
          <div
            key={asset.id}
            className={`asset-card ${selectedAsset?.id === asset.id ? 'selected' : ''}`}
            onClick={() => handleAssetSelect(asset)}
          >
            <div className="asset-icon">{asset.symbol.charAt(0)}</div>
            <div className="asset-info">
              <span className="asset-name">{asset.name}</span>
              <span className="asset-symbol">{asset.symbol}</span>
            </div>
            <div className="asset-network">{asset.network_id ? 'Multi-chain' : 'Available'}</div>
          </div>
        ))}
      </div>

      {selectedAsset && (
        <div className="deposit-details">
          <h3>Deposit {selectedAsset.symbol}</h3>
          <p className="min-amount">
            Minimum deposit: {selectedAsset.min_deposit_amount} {selectedAsset.symbol}
          </p>
          
          {loading ? (
            <div className="loading">Loading deposit address...</div>
          ) : depositAddress ? (
            <div className="deposit-address-container">
              <label>Deposit Address</label>
              <div className="address-box">
                <input
                  type="text"
                  value={depositAddress}
                  readOnly
                  className="address-input"
                />
                <button onClick={copyToClipboard} className="copy-btn">
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
              <p className="warning">
                ⚠️ Only send {selectedAsset.symbol} to this address. Sending other tokens may result in permanent loss.
              </p>
            </div>
          ) : (
            <p className="no-address">No deposit address available</p>
          )}
        </div>
      )}

      <style jsx>{`
        .crypto-deposit {
          padding: 20px;
        }
        
        h2 {
          margin-bottom: 20px;
          color: #333;
        }

        .error-message {
          background: #fee;
          color: #c00;
          padding: 10px;
          border-radius: 4px;
          margin-bottom: 20px;
        }

        .asset-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
          gap: 15px;
          margin-bottom: 30px;
        }

        .asset-card {
          background: #f8f9fa;
          border: 2px solid #e9ecef;
          border-radius: 8px;
          padding: 15px;
          cursor: pointer;
          transition: all 0.2s;
          text-align: center;
        }

        .asset-card:hover {
          border-color: #007bff;
          transform: translateY(-2px);
        }

        .asset-card.selected {
          border-color: #007bff;
          background: #e7f1ff;
        }

        .asset-icon {
          width: 40px;
          height: 40px;
          background: linear-gradient(135deg, #007bff, #00d4ff);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-weight: bold;
          margin: 0 auto 10px;
        }

        .asset-info {
          display: flex;
          flex-direction: column;
        }

        .asset-name {
          font-weight: 600;
          font-size: 14px;
        }

        .asset-symbol {
          color: #666;
          font-size: 12px;
        }

        .asset-network {
          margin-top: 8px;
          font-size: 11px;
          color: #28a745;
        }

        .deposit-details {
          background: #f8f9fa;
          border-radius: 8px;
          padding: 20px;
        }

        .min-amount {
          color: #666;
          margin-bottom: 20px;
        }

        .loading {
          text-align: center;
          padding: 20px;
          color: #666;
        }

        .deposit-address-container {
          margin-top: 20px;
        }

        .address-box {
          display: flex;
          gap: 10px;
          margin-top: 10px;
        }

        .address-input {
          flex: 1;
          padding: 12px;
          border: 1px solid #ddd;
          border-radius: 4px;
          font-family: monospace;
          font-size: 14px;
        }

        .copy-btn {
          padding: 12px 20px;
          background: #007bff;
          color: white;
          border: none;
          border-radius: 4px;
          cursor: pointer;
        }

        .copy-btn:hover {
          background: #0056b3;
        }

        .warning {
          margin-top: 15px;
          padding: 10px;
          background: #fff3cd;
          border: 1px solid #ffc107;
          border-radius: 4px;
          color: #856404;
          font-size: 14px;
        }

        .no-address {
          color: #dc3545;
          text-align: center;
          padding: 20px;
        }
      `}</style>
    </div>
  );
}
