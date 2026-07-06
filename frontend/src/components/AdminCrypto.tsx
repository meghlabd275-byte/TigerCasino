'use client';

import { useState, useEffect } from 'react';
import {
  adminGetAllNetworks,
  adminCreateNetwork,
  adminUpdateNetwork,
  adminGetAllAssets,
  adminCreateAsset,
  adminUpdateAsset,
  adminLinkAssetToNetwork,
  adminUpdateAssetNetwork,
  adminSetNetworkFee,
  adminGetBrandLevels,
  adminUpdateBrandLevel,
  adminSetAdminWallet,
  adminGetAllDeposits,
  adminGetAllWithdrawals,
  adminConfirmDeposit,
  adminConfirmWithdrawal,
  adminCancelWithdrawal,
} from '@/lib/crypto';
import {
  CryptoNetwork,
  CryptoAssetWithNetwork,
  BrandLevel,
  NetworkFee,
  AdminWalletAddress,
  CryptoDeposit,
  CryptoWithdrawal,
} from '@/types/crypto';

type TabType = 'networks' | 'assets' | 'fees' | 'brands' | 'wallets' | 'deposits' | 'withdrawals';

export default function AdminCrypto() {
  const [activeTab, setActiveTab] = useState<TabType>('networks');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Data states
  const [networks, setNetworks] = useState<CryptoNetwork[]>([]);
  const [assets, setAssets] = useState<CryptoAssetWithNetwork[]>([]);
  const [brandLevels, setBrandLevels] = useState<BrandLevel[]>([]);
  const [deposits, setDeposits] = useState<CryptoDeposit[]>([]);
  const [withdrawals, setWithdrawals] = useState<CryptoWithdrawal[]>([]);

  // Form states
  const [editingNetwork, setEditingNetwork] = useState<CryptoNetwork | null>(null);
  const [editingAsset, setEditingAsset] = useState<CryptoAssetWithNetwork | null>(null);
  const [editingBrandLevel, setEditingBrandLevel] = useState<BrandLevel | null>(null);

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async () => {
    setLoading(true);
    setError('');
    try {
      switch (activeTab) {
        case 'networks':
          const nets = await adminGetAllNetworks();
          setNetworks(nets);
          break;
        case 'assets':
          const assts = await adminGetAllAssets();
          setAssets(assts);
          break;
        case 'brands':
          const levels = await adminGetBrandLevels();
          setBrandLevels(levels);
          break;
        case 'deposits':
          const deps = await adminGetAllDeposits();
          setDeposits(deps.deposits);
          break;
        case 'withdrawals':
          const wds = await adminGetAllWithdrawals();
          setWithdrawals(withdrawals);
          break;
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateNetwork = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    
    const formData = new FormData(e.currentTarget);
    const network = {
      name: formData.get('name') as string,
      chain_id: formData.get('chain_id') as string,
      symbol: formData.get('symbol') as string,
      explorer_url: formData.get('explorer_url') as string,
      min_confirmation_blocks: parseInt(formData.get('min_confirmation_blocks') as string) || 6,
    };

    try {
      await adminCreateNetwork(network);
      setSuccess('Network created successfully');
      loadData();
      (e.target as HTMLFormElement).reset();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleUpdateNetwork = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!editingNetwork) return;
    
    setError('');
    setSuccess('');
    
    const formData = new FormData(e.currentTarget);
    const network = {
      name: formData.get('name') as string,
      chain_id: formData.get('chain_id') as string,
      symbol: formData.get('symbol') as string,
      explorer_url: formData.get('explorer_url') as string,
      is_active: formData.get('is_active') === 'true',
      is_deposit_enabled: formData.get('is_deposit_enabled') === 'true',
      is_withdrawal_enabled: formData.get('is_withdrawal_enabled') === 'true',
      min_confirmation_blocks: parseInt(formData.get('min_confirmation_blocks') as string) || 6,
    };

    try {
      await adminUpdateNetwork(editingNetwork.id, network);
      setSuccess('Network updated successfully');
      setEditingNetwork(null);
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleUpdateBrandLevel = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!editingBrandLevel) return;
    
    setError('');
    setSuccess('');
    
    const formData = new FormData(e.currentTarget);
    const level = {
      name: formData.get('name') as string,
      level: parseInt(formData.get('level') as string),
      deposit_fee_discount_percent: parseFloat(formData.get('deposit_fee_discount_percent') as string) || 0,
      withdrawal_fee_discount_percent: parseFloat(formData.get('withdrawal_fee_discount_percent') as string) || 0,
      is_active: formData.get('is_active') === 'true',
    };

    try {
      await adminUpdateBrandLevel(editingBrandLevel.id, level);
      setSuccess('Brand level updated successfully');
      setEditingBrandLevel(null);
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleConfirmDeposit = async (id: string) => {
    try {
      await adminConfirmDeposit(id);
      setSuccess('Deposit confirmed');
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleConfirmWithdrawal = async (id: string) => {
    const txHash = prompt('Enter transaction hash:');
    if (!txHash) return;
    
    try {
      await adminConfirmWithdrawal(id, txHash);
      setSuccess('Withdrawal confirmed');
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleCancelWithdrawal = async (id: string) => {
    if (!confirm('Are you sure you want to cancel this withdrawal?')) return;
    
    try {
      await adminCancelWithdrawal(id);
      setSuccess('Withdrawal cancelled');
      loadData();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const renderNetworksTab = () => (
    <div className="tab-content">
      <h3>Manage Networks</h3>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      {editingNetwork ? (
        <form onSubmit={handleUpdateNetwork} className="edit-form">
          <h4>Edit Network</h4>
          <input type="hidden" name="id" value={editingNetwork.id} />
          
          <div className="form-group">
            <label>Name</label>
            <input type="text" name="name" defaultValue={editingNetwork.name} required />
          </div>
          
          <div className="form-group">
            <label>Chain ID</label>
            <input type="text" name="chain_id" defaultValue={editingNetwork.chain_id} />
          </div>
          
          <div className="form-group">
            <label>Symbol</label>
            <input type="text" name="symbol" defaultValue={editingNetwork.symbol} required />
          </div>
          
          <div className="form-group">
            <label>Explorer URL</label>
            <input type="text" name="explorer_url" defaultValue={editingNetwork.explorer_url} />
          </div>
          
          <div className="form-group">
            <label>Min Confirmations</label>
            <input type="number" name="min_confirmation_blocks" defaultValue={editingNetwork.min_confirmation_blocks} />
          </div>
          
          <div className="form-group">
            <label>Active</label>
            <select name="is_active" defaultValue={editingNetwork.is_active.toString()}>
              <option value="true">Yes</option>
              <option value="false">No</option>
            </select>
          </div>
          
          <div className="form-group">
            <label>Deposit Enabled</label>
            <select name="is_deposit_enabled" defaultValue={editingNetwork.is_deposit_enabled.toString()}>
              <option value="true">Yes</option>
              <option value="false">No</option>
            </select>
          </div>
          
          <div className="form-group">
            <label>Withdrawal Enabled</label>
            <select name="is_withdrawal_enabled" defaultValue={editingNetwork.is_withdrawal_enabled.toString()}>
              <option value="true">Yes</option>
              <option value="false">No</option>
            </select>
          </div>
          
          <div className="form-actions">
            <button type="submit">Update</button>
            <button type="button" onClick={() => setEditingNetwork(null)}>Cancel</button>
          </div>
        </form>
      ) : (
        <form onSubmit={handleCreateNetwork} className="create-form">
          <h4>Add New Network</h4>
          
          <div className="form-row">
            <div className="form-group">
              <label>Name</label>
              <input type="text" name="name" placeholder="e.g., Ethereum" required />
            </div>
            
            <div className="form-group">
              <label>Chain ID</label>
              <input type="text" name="chain_id" placeholder="e.g., 1" />
            </div>
          </div>
          
          <div className="form-row">
            <div className="form-group">
              <label>Symbol</label>
              <input type="text" name="symbol" placeholder="e.g., ETH" required />
            </div>
            
            <div className="form-group">
              <label>Min Confirmations</label>
              <input type="number" name="min_confirmation_blocks" defaultValue="12" />
            </div>
          </div>
          
          <div className="form-group">
            <label>Explorer URL</label>
            <input type="text" name="explorer_url" placeholder="https://etherscan.io" />
          </div>
          
          <button type="submit">Create Network</button>
        </form>
      )}

      <div className="data-table">
        <h4>Existing Networks</h4>
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Chain ID</th>
              <th>Symbol</th>
              <th>Confirmations</th>
              <th>Deposit</th>
              <th>Withdrawal</th>
              <th>Active</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {networks.map((network) => (
              <tr key={network.id}>
                <td>{network.name}</td>
                <td>{network.chain_id}</td>
                <td>{network.symbol}</td>
                <td>{network.min_confirmation_blocks}</td>
                <td>{network.is_deposit_enabled ? '✓' : '✗'}</td>
                <td>{network.is_withdrawal_enabled ? '✓' : '✗'}</td>
                <td>{network.is_active ? '✓' : '✗'}</td>
                <td>
                  <button onClick={() => setEditingNetwork(network)}>Edit</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );

  const renderBrandLevelsTab = () => (
    <div className="tab-content">
      <h3>Manage Brand Levels (Fee Discounts)</h3>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      {editingBrandLevel ? (
        <form onSubmit={handleUpdateBrandLevel} className="edit-form">
          <h4>Edit Brand Level: {editingBrandLevel.name}</h4>
          
          <div className="form-group">
            <label>Name</label>
            <input type="text" name="name" defaultValue={editingBrandLevel.name} required />
          </div>
          
          <div className="form-group">
            <label>Level</label>
            <input type="number" name="level" defaultValue={editingBrandLevel.level} required />
          </div>
          
          <div className="form-group">
            <label>Deposit Fee Discount (%)</label>
            <input 
              type="number" 
              name="deposit_fee_discount_percent" 
              defaultValue={editingBrandLevel.deposit_fee_discount_percent} 
              step="0.01"
            />
            <small>Current: {editingBrandLevel.deposit_fee_discount_percent}% discount</small>
          </div>
          
          <div className="form-group">
            <label>Withdrawal Fee Discount (%)</label>
            <input 
              type="number" 
              name="withdrawal_fee_discount_percent" 
              defaultValue={editingBrandLevel.withdrawal_fee_discount_percent} 
              step="0.01"
            />
            <small>Current: {editingBrandLevel.withdrawal_fee_discount_percent}% discount</small>
          </div>
          
          <div className="form-group">
            <label>Active</label>
            <select name="is_active" defaultValue={editingBrandLevel.is_active.toString()}>
              <option value="true">Yes</option>
              <option value="false">No</option>
            </select>
          </div>
          
          <div className="form-actions">
            <button type="submit">Update</button>
            <button type="button" onClick={() => setEditingBrandLevel(null)}>Cancel</button>
          </div>
        </form>
      ) : (
        <div className="info-box">
          <h4>Default Brand Levels</h4>
          <p>The system includes 6 default brand levels. You can edit the fee discounts for each level below. The "White" level (Level 6) has a default 20% fee discount as requested.</p>
        </div>
      )}

      <div className="brand-levels-grid">
        {brandLevels.map((level) => (
          <div key={level.id} className="brand-level-card">
            <h4>{level.name}</h4>
            <p className="level-badge">Level {level.level}</p>
            <div className="discounts">
              <p>Deposit Discount: {level.deposit_fee_discount_percent}%</p>
              <p>Withdrawal Discount: {level.withdrawal_fee_discount_percent}%</p>
            </div>
            <button onClick={() => setEditingBrandLevel(level)}>Edit</button>
          </div>
        ))}
      </div>
    </div>
  );

  const renderDepositsTab = () => (
    <div className="tab-content">
      <h3>Manage Deposits</h3>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <div className="data-table">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>User</th>
              <th>Asset</th>
              <th>Amount</th>
              <th>Fee</th>
              <th>Net</th>
              <th>Status</th>
              <th>Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {deposits.map((deposit) => (
              <tr key={deposit.id}>
                <td>{deposit.id.slice(0, 8)}...</td>
                <td>{deposit.user_id?.slice(0, 8)}...</td>
                <td>{deposit.asset?.symbol}</td>
                <td>{deposit.amount}</td>
                <td>{deposit.fee}</td>
                <td>{deposit.net_amount}</td>
                <td>
                  <span className={`status-badge ${deposit.status}`}>{deposit.status}</span>
                </td>
                <td>{new Date(deposit.created_at).toLocaleDateString()}</td>
                <td>
                  {deposit.status === 'pending' && (
                    <button onClick={() => handleConfirmDeposit(deposit.id)}>Confirm</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );

  const renderWithdrawalsTab = () => (
    <div className="tab-content">
      <h3>Manage Withdrawals</h3>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <div className="data-table">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>User</th>
              <th>Asset</th>
              <th>Amount</th>
              <th>Fee</th>
              <th>Net</th>
              <th>Address</th>
              <th>Status</th>
              <th>Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {withdrawals.map((withdrawal) => (
              <tr key={withdrawal.id}>
                <td>{withdrawal.id.slice(0, 8)}...</td>
                <td>{withdrawal.user_id?.slice(0, 8)}...</td>
                <td>{withdrawal.asset?.symbol}</td>
                <td>{withdrawal.amount}</td>
                <td>{withdrawal.fee}</td>
                <td>{withdrawal.net_amount}</td>
                <td className="address-cell">{withdrawal.address?.slice(0, 10)}...</td>
                <td>
                  <span className={`status-badge ${withdrawal.status}`}>{withdrawal.status}</span>
                </td>
                <td>{new Date(withdrawal.created_at).toLocaleDateString()}</td>
                <td>
                  {withdrawal.status === 'pending' && (
                    <div className="action-buttons">
                      <button onClick={() => handleConfirmWithdrawal(withdrawal.id)}>Confirm</button>
                      <button onClick={() => handleCancelWithdrawal(withdrawal.id)} className="cancel">Cancel</button>
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );

  const renderAssetsTab = () => (
    <div className="tab-content">
      <h3>Manage Crypto Assets</h3>
      <p>Add new crypto assets or edit existing ones. Each asset can be linked to multiple blockchain networks.</p>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <div className="data-table">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Symbol</th>
              <th>Decimals</th>
              <th>Min Deposit</th>
              <th>Min Withdrawal</th>
              <th>Active</th>
            </tr>
          </thead>
          <tbody>
            {assets.map((asset) => (
              <tr key={asset.id}>
                <td>{asset.name}</td>
                <td>{asset.symbol}</td>
                <td>{asset.decimals}</td>
                <td>{asset.min_deposit_amount}</td>
                <td>{asset.min_withdrawal_amount}</td>
                <td>{asset.is_active ? '✓' : '✗'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );

  const renderFeesTab = () => (
    <div className="tab-content">
      <h3>Manage Network Fees</h3>
      <p>Set deposit and withdrawal fees for each asset-network combination.</p>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}
      
      <div className="info-box">
        <h4>Fee Management</h4>
        <p>To update fees for a specific asset-network combination, use the API or contact development team to add fee management UI.</p>
      </div>
    </div>
  );

  const renderWalletsTab = () => (
    <div className="tab-content">
      <h3>Manage Admin Wallets</h3>
      <p>Configure wallet addresses for receiving deposits and processing withdrawals.</p>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}
      
      <div className="info-box">
        <h4>Wallet Management</h4>
        <p>To add admin wallets for each asset-network, use the API or contact development team to add wallet management UI.</p>
      </div>
    </div>
  );

  return (
    <div className="admin-crypto">
      <h2>Crypto Management</h2>

      <div className="tabs">
        <button 
          className={activeTab === 'networks' ? 'active' : ''} 
          onClick={() => setActiveTab('networks')}
        >
          Networks
        </button>
        <button 
          className={activeTab === 'assets' ? 'active' : ''} 
          onClick={() => setActiveTab('assets')}
        >
          Assets
        </button>
        <button 
          className={activeTab === 'fees' ? 'active' : ''} 
          onClick={() => setActiveTab('fees')}
        >
          Fees
        </button>
        <button 
          className={activeTab === 'brands' ? 'active' : ''} 
          onClick={() => setActiveTab('brands')}
        >
          Brand Levels
        </button>
        <button 
          className={activeTab === 'wallets' ? 'active' : ''} 
          onClick={() => setActiveTab('wallets')}
        >
          Wallets
        </button>
        <button 
          className={activeTab === 'deposits' ? 'active' : ''} 
          onClick={() => setActiveTab('deposits')}
        >
          Deposits
        </button>
        <button 
          className={activeTab === 'withdrawals' ? 'active' : ''} 
          onClick={() => setActiveTab('withdrawals')}
        >
          Withdrawals
        </button>
      </div>

      {loading ? (
        <div className="loading">Loading...</div>
      ) : (
        <>
          {activeTab === 'networks' && renderNetworksTab()}
          {activeTab === 'assets' && renderAssetsTab()}
          {activeTab === 'fees' && renderFeesTab()}
          {activeTab === 'brands' && renderBrandLevelsTab()}
          {activeTab === 'wallets' && renderWalletsTab()}
          {activeTab === 'deposits' && renderDepositsTab()}
          {activeTab === 'withdrawals' && renderWithdrawalsTab()}
        </>
      )}

      <style jsx>{`
        .admin-crypto {
          padding: 20px;
        }

        h2 {
          margin-bottom: 20px;
        }

        .tabs {
          display: flex;
          gap: 10px;
          margin-bottom: 20px;
          border-bottom: 2px solid #e9ecef;
          padding-bottom: 10px;
          overflow-x: auto;
        }

        .tabs button {
          padding: 10px 20px;
          border: none;
          background: none;
          cursor: pointer;
          font-size: 14px;
          border-radius: 4px;
          white-space: nowrap;
        }

        .tabs button:hover {
          background: #f8f9fa;
        }

        .tabs button.active {
          background: #007bff;
          color: white;
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

        .loading {
          text-align: center;
          padding: 40px;
          color: #666;
        }

        .tab-content h3 {
          margin-bottom: 20px;
        }

        .tab-content h4 {
          margin: 20px 0 10px;
        }

        .create-form, .edit-form {
          background: #f8f9fa;
          padding: 20px;
          border-radius: 8px;
          margin-bottom: 30px;
        }

        .form-row {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 15px;
        }

        .form-group {
          margin-bottom: 15px;
        }

        .form-group label {
          display: block;
          margin-bottom: 5px;
          font-weight: 600;
        }

        .form-group input,
        .form-group select {
          width: 100%;
          padding: 8px;
          border: 1px solid #ddd;
          border-radius: 4px;
        }

        .form-group small {
          display: block;
          margin-top: 5px;
          color: #666;
        }

        .form-actions {
          display: flex;
          gap: 10px;
        }

        .form-actions button {
          padding: 10px 20px;
          border: none;
          border-radius: 4px;
          cursor: pointer;
        }

        .form-actions button[type="submit"] {
          background: #007bff;
          color: white;
        }

        .form-actions button[type="button"] {
          background: #6c757d;
          color: white;
        }

        button[type="submit"] {
          padding: 10px 20px;
          background: #007bff;
          color: white;
          border: none;
          border-radius: 4px;
          cursor: pointer;
        }

        .data-table {
          overflow-x: auto;
        }

        .data-table table {
          width: 100%;
          border-collapse: collapse;
        }

        .data-table th,
        .data-table td {
          padding: 12px;
          text-align: left;
          border-bottom: 1px solid #e9ecef;
        }

        .data-table th {
          background: #f8f9fa;
          font-weight: 600;
        }

        .data-table td button {
          padding: 5px 10px;
          background: #007bff;
          color: white;
          border: none;
          border-radius: 4px;
          cursor: pointer;
          font-size: 12px;
        }

        .data-table td button.cancel {
          background: #dc3545;
          margin-left: 5px;
        }

        .action-buttons {
          display: flex;
          gap: 5px;
        }

        .status-badge {
          padding: 4px 8px;
          border-radius: 4px;
          font-size: 12px;
        }

        .status-badge.pending {
          background: #fff3cd;
          color: #856404;
        }

        .status-badge.completed {
          background: #d4edda;
          color: #155724;
        }

        .status-badge.failed, .status-badge.cancelled {
          background: #f8d7da;
          color: #721c24;
        }

        .info-box {
          background: #e7f1ff;
          padding: 20px;
          border-radius: 8px;
          margin-bottom: 20px;
        }

        .brand-levels-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
          gap: 20px;
        }

        .brand-level-card {
          background: #f8f9fa;
          padding: 20px;
          border-radius: 8px;
          text-align: center;
        }

        .brand-level-card h4 {
          margin: 0 0 10px;
        }

        .level-badge {
          background: #007bff;
          color: white;
          padding: 4px 12px;
          border-radius: 12px;
          font-size: 12px;
          display: inline-block;
          margin-bottom: 15px;
        }

        .discounts {
          margin-bottom: 15px;
        }

        .discounts p {
          margin: 5px 0;
          font-size: 14px;
        }

        .brand-level-card button {
          padding: 8px 16px;
          background: #28a745;
          color: white;
          border: none;
          border-radius: 4px;
          cursor: pointer;
        }

        .address-cell {
          font-family: monospace;
          font-size: 12px;
        }
      `}</style>
    </div>
  );
}
