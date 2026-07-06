'use client';

import { useState, useEffect } from 'react';
import { getAssets, createWithdrawal, calculateFees, getUserWithdrawals } from '@/lib/crypto';
import { CryptoAssetWithNetwork, CryptoWithdrawal } from '@/types/crypto';

interface CryptoWithdrawalProps {
  balance: number;
  onSuccess?: () => void;
}

export default function CryptoWithdrawal({ balance, onSuccess }: CryptoWithdrawalProps) {
  const [assets, setAssets] = useState<CryptoAssetWithNetwork[]>([]);
  const [selectedAsset, setSelectedAsset] = useState<CryptoAssetWithNetwork | null>(null);
  const [amount, setAmount] = useState<string>('');
  const [address, setAddress] = useState<string>('');
  const [withdrawals, setWithdrawals] = useState<CryptoWithdrawal[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [fees, setFees] = useState<{ deposit_fee: number; withdrawal_fee: number }>({
    deposit_fee: 0,
    withdrawal_fee: 0,
  });

  useEffect(() => {
    loadAssets();
    loadWithdrawals();
  }, []);

  useEffect(() => {
    if (selectedAsset && amount) {
      calculateWithdrawalFees();
    }
  }, [selectedAsset, amount]);

  const loadAssets = async () => {
    try {
      const data = await getAssets();
      setAssets(data);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const loadWithdrawals = async () => {
    try {
      const data = await getUserWithdrawals();
      setWithdrawals(data);
    } catch (err: any) {
      console.error('Failed to load withdrawals:', err);
    }
  };

  const calculateWithdrawalFees = async () => {
    if (!selectedAsset || !amount) return;
    
    try {
      const data = await calculateFees(
        selectedAsset.id,
        selectedAsset.network_id,
        parseFloat(amount)
      );
      setFees(data);
    } catch (err: any) {
      console.error('Failed to calculate fees:', err);
    }
  };

  const handleAssetSelect = (asset: CryptoAssetWithNetwork) => {
    if (!asset.withdrawal_enabled) {
      setError('Withdrawal not enabled for this asset');
      return;
    }
    setSelectedAsset(asset);
    setError('');
    setSuccess('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (!selectedAsset) {
      setError('Please select an asset');
      return;
    }

    const amountNum = parseFloat(amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      setError('Please enter a valid amount');
      return;
    }

    if (amountNum < selectedAsset.min_withdrawal_amount) {
      setError(`Minimum withdrawal amount is ${selectedAsset.min_withdrawal_amount} ${selectedAsset.symbol}`);
      return;
    }

    if (amountNum > balance) {
      setError('Insufficient balance');
      return;
    }

    const totalWithFees = amountNum + fees.withdrawal_fee;
    if (totalWithFees > balance) {
      setError('Insufficient balance (including fees)');
      return;
    }

    if (!address) {
      setError('Please enter a withdrawal address');
      return;
    }

    try {
      setSubmitting(true);
      await createWithdrawal({
        asset_id: selectedAsset.id,
        network_id: selectedAsset.network_id,
        amount: amountNum,
        address,
      });
      setSuccess('Withdrawal request submitted successfully');
      setAmount('');
      setAddress('');
      loadWithdrawals();
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  const netAmount = amount ? parseFloat(amount) - fees.withdrawal_fee : 0;

  return (
    <div className="crypto-withdrawal">
      <h2>Withdraw Crypto</h2>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <form onSubmit={handleSubmit} className="withdrawal-form">
        <div className="balance-display">
          Available Balance: <strong>{balance.toFixed(2)}</strong>
        </div>

        <div className="asset-grid">
          {assets.filter(a => a.withdrawal_enabled).map((asset) => (
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
            </div>
          ))}
        </div>

        {selectedAsset && (
          <>
            <div className="form-group">
              <label>Amount</label>
              <input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder={`Min: ${selectedAsset.min_withdrawal_amount}`}
                step="any"
                min="0"
              />
              <span className="min-note">
                Min: {selectedAsset.min_withdrawal_amount} {selectedAsset.symbol}
              </span>
            </div>

            <div className="form-group">
              <label>Withdrawal Address</label>
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                placeholder="Enter your wallet address"
              />
            </div>

            {amount && parseFloat(amount) > 0 && (
              <div className="fee-summary">
                <div className="fee-row">
                  <span>Withdrawal Fee:</span>
                  <span>{fees.withdrawal_fee} {selectedAsset.symbol}</span>
                </div>
                <div className="fee-row total">
                  <span>You will receive:</span>
                  <span>{netAmount.toFixed(8)} {selectedAsset.symbol}</span>
                </div>
              </div>
            )}
          </>
        )}

        <button
          type="submit"
          className="submit-btn"
          disabled={submitting || !selectedAsset || !amount || !address}
        >
          {submitting ? 'Processing...' : 'Withdraw'}
        </button>
      </form>

      {withdrawals.length > 0 && (
        <div className="withdrawal-history">
          <h3>Recent Withdrawals</h3>
          <div className="history-list">
            {withdrawals.slice(0, 5).map((w) => (
              <div key={w.id} className="history-item">
                <div className="history-asset">
                  <span className="symbol">{w.asset?.symbol || 'Crypto'}</span>
                  <span className="network">{w.network?.name || 'Unknown'}</span>
                </div>
                <div className="history-amount">
                  <span className="amount">-{w.amount}</span>
                  <span className={`status ${w.status}`}>{w.status}</span>
                </div>
                <div className="history-date">
                  {new Date(w.created_at).toLocaleDateString()}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <style jsx>{`
        .crypto-withdrawal {
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

        .success-message {
          background: #efe;
          color: #080;
          padding: 10px;
          border-radius: 4px;
          margin-bottom: 20px;
        }

        .balance-display {
          background: #e7f1ff;
          padding: 15px;
          border-radius: 8px;
          margin-bottom: 20px;
          font-size: 16px;
        }

        .asset-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
          gap: 15px;
          margin-bottom: 20px;
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
        }

        .asset-card.selected {
          border-color: #007bff;
          background: #e7f1ff;
        }

        .asset-icon {
          width: 36px;
          height: 36px;
          background: linear-gradient(135deg, #007bff, #00d4ff);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-weight: bold;
          margin: 0 auto 8px;
        }

        .asset-info {
          display: flex;
          flex-direction: column;
        }

        .asset-name {
          font-weight: 600;
          font-size: 13px;
        }

        .asset-symbol {
          color: #666;
          font-size: 11px;
        }

        .form-group {
          margin-bottom: 20px;
        }

        .form-group label {
          display: block;
          margin-bottom: 8px;
          font-weight: 600;
          color: #333;
        }

        .form-group input {
          width: 100%;
          padding: 12px;
          border: 1px solid #ddd;
          border-radius: 4px;
          font-size: 16px;
        }

        .min-note {
          display: block;
          margin-top: 5px;
          color: #666;
          font-size: 12px;
        }

        .fee-summary {
          background: #f8f9fa;
          padding: 15px;
          border-radius: 8px;
          margin-bottom: 20px;
        }

        .fee-row {
          display: flex;
          justify-content: space-between;
          padding: 8px 0;
        }

        .fee-row.total {
          border-top: 1px solid #ddd;
          font-weight: 600;
          margin-top: 8px;
          padding-top: 16px;
        }

        .submit-btn {
          width: 100%;
          padding: 15px;
          background: #28a745;
          color: white;
          border: none;
          border-radius: 4px;
          font-size: 16px;
          font-weight: 600;
          cursor: pointer;
        }

        .submit-btn:disabled {
          background: #ccc;
          cursor: not-allowed;
        }

        .submit-btn:hover:not(:disabled) {
          background: #218838;
        }

        .withdrawal-history {
          margin-top: 30px;
        }

        .history-list {
          background: #f8f9fa;
          border-radius: 8px;
          overflow: hidden;
        }

        .history-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 15px;
          border-bottom: 1px solid #e9ecef;
        }

        .history-item:last-child {
          border-bottom: none;
        }

        .history-asset {
          display: flex;
          flex-direction: column;
        }

        .history-asset .symbol {
          font-weight: 600;
        }

        .history-asset .network {
          font-size: 12px;
          color: #666;
        }

        .history-amount {
          text-align: right;
        }

        .history-amount .amount {
          font-weight: 600;
          color: #dc3545;
        }

        .history-amount .status {
          display: block;
          font-size: 12px;
          padding: 2px 8px;
          border-radius: 4px;
          margin-top: 4px;
        }

        .history-amount .status.pending {
          background: #fff3cd;
          color: #856404;
        }

        .history-amount .status.completed {
          background: #d4edda;
          color: #155724;
        }

        .history-amount .status.failed {
          background: #f8d7da;
          color: #721c24;
        }

        .history-date {
          color: #666;
          font-size: 12px;
        }
      `}</style>
    </div>
  );
}
