// Crypto Types

export interface CryptoNetwork {
  id: string;
  name: string;
  chain_id: string;
  symbol: string;
  explorer_url: string;
  rpc_url?: string;
  is_withdrawal_enabled: boolean;
  is_deposit_enabled: boolean;
  is_active: boolean;
  min_confirmation_blocks: number;
  created_at: string;
  updated_at: string;
}

export interface CryptoAsset {
  id: string;
  name: string;
  symbol: string;
  decimals: number;
  contract_address?: string;
  is_active: boolean;
  min_deposit_amount: number;
  min_withdrawal_amount: number;
  created_at: string;
  updated_at: string;
}

export interface CryptoAssetWithNetwork {
  id: string;
  name: string;
  symbol: string;
  decimals: number;
  contract_address?: string;
  is_active: boolean;
  min_deposit_amount: number;
  min_withdrawal_amount: number;
  network_id: string;
  deposit_enabled: boolean;
  withdrawal_enabled: boolean;
  contract_address: string;
  deposit_fee: number;
  withdrawal_fee: number;
  deposit_fee_percent: number;
  withdrawal_fee_percent: number;
}

export interface CryptoAssetNetwork {
  id: string;
  asset_id: string;
  network_id: string;
  deposit_enabled: boolean;
  withdrawal_enabled: boolean;
  contract_address: string;
  asset?: CryptoAsset;
  network?: CryptoNetwork;
}

export interface NetworkFee {
  id: string;
  asset_id: string;
  network_id: string;
  deposit_fee: number;
  withdrawal_fee: number;
  deposit_fee_percent: number;
  withdrawal_fee_percent: number;
}

export interface BrandLevel {
  id: string;
  name: string;
  level: number;
  deposit_fee_discount_percent: number;
  withdrawal_fee_discount_percent: number;
  is_active: boolean;
}

export interface CryptoDeposit {
  id: string;
  user_id: string;
  asset_id: string;
  network_id: string;
  amount: number;
  fee: number;
  net_amount: number;
  address: string;
  tx_hash?: string;
  confirmations: number;
  status: 'pending' | 'completed' | 'failed';
  created_at: string;
  processed_at?: string;
  asset?: CryptoAsset;
  network?: CryptoNetwork;
}

export interface CryptoWithdrawal {
  id: string;
  user_id: string;
  asset_id: string;
  network_id: string;
  amount: number;
  fee: number;
  net_amount: number;
  address: string;
  tx_hash?: string;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  created_at: string;
  processed_at?: string;
  asset?: CryptoAsset;
  network?: CryptoNetwork;
}

export interface AdminWalletAddress {
  id: string;
  asset_id: string;
  network_id: string;
  address: string;
  private_key_encrypted?: string;
  is_active: boolean;
}

export interface DepositRequest {
  asset_id: string;
  network_id: string;
}

export interface WithdrawalRequest {
  asset_id: string;
  network_id: string;
  amount: number;
  address: string;
}

export interface FeeCalculation {
  deposit_fee: number;
  withdrawal_fee: number;
}
