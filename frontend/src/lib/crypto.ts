import {
  CryptoNetwork,
  CryptoAssetWithNetwork,
  CryptoAssetNetwork,
  CryptoDeposit,
  CryptoWithdrawal,
  BrandLevel,
  NetworkFee,
  AdminWalletAddress,
  DepositRequest,
  WithdrawalRequest,
  FeeCalculation,
} from '@/types/crypto';

const API_BASE = '/api/crypto';

// Helper function for API calls
async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || 'Request failed');
  }

  return response.json();
}

// Crypto Networks API

export async function getNetworks(): Promise<CryptoNetwork[]> {
  return fetchApi<CryptoNetwork[]>('/networks');
}

export async function getAssets(): Promise<CryptoAssetWithNetwork[]> {
  return fetchApi<CryptoAssetWithNetwork[]>('/assets');
}

export async function getAssetNetworks(assetId: string): Promise<CryptoAssetNetwork[]> {
  return fetchApi<CryptoAssetNetwork[]>(`/assets/${assetId}/networks`);
}

// Deposit API

export async function getDepositAddress(
  assetId: string,
  networkId: string
): Promise<{ address: string }> {
  return fetchApi<{ address: string }>('/deposit/address', {
    method: 'POST',
    body: JSON.stringify({ asset_id: assetId, network_id: networkId }),
  });
}

export async function getUserDeposits(): Promise<CryptoDeposit[]> {
  return fetchApi<CryptoDeposit[]>('/deposits');
}

// Withdrawal API

export async function createWithdrawal(
  request: WithdrawalRequest
): Promise<CryptoWithdrawal> {
  return fetchApi<CryptoWithdrawal>('/withdrawals', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}

export async function getUserWithdrawals(): Promise<CryptoWithdrawal[]> {
  return fetchApi<CryptoWithdrawal[]>('/withdrawals');
}

// Fees API

export async function calculateFees(
  assetId: string,
  networkId: string,
  amount: number
): Promise<FeeCalculation> {
  return fetchApi<FeeCalculation>('/fees', {
    method: 'POST',
    body: JSON.stringify({
      asset_id: assetId,
      network_id: networkId,
      amount,
    }),
  });
}

// Admin API

export async function adminGetAllNetworks(): Promise<CryptoNetwork[]> {
  return fetchApi<CryptoNetwork[]>('/admin/crypto/networks');
}

export async function adminCreateNetwork(
  network: Partial<CryptoNetwork>
): Promise<CryptoNetwork> {
  return fetchApi<CryptoNetwork>('/admin/crypto/networks', {
    method: 'POST',
    body: JSON.stringify(network),
  });
}

export async function adminUpdateNetwork(
  id: string,
  network: Partial<CryptoNetwork>
): Promise<CryptoNetwork> {
  return fetchApi<CryptoNetwork>(`/admin/crypto/networks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(network),
  });
}

export async function adminGetAllAssets(): Promise<any[]> {
  return fetchApi<any[]>('/admin/crypto/assets');
}

export async function adminCreateAsset(
  asset: Partial<CryptoAssetWithNetwork>
): Promise<any> {
  return fetchApi<any>('/admin/crypto/assets', {
    method: 'POST',
    body: JSON.stringify(asset),
  });
}

export async function adminUpdateAsset(
  id: string,
  asset: Partial<CryptoAssetWithNetwork>
): Promise<any> {
  return fetchApi<any>(`/admin/crypto/assets/${id}`, {
    method: 'PUT',
    body: JSON.stringify(asset),
  });
}

export async function adminLinkAssetToNetwork(
  assetNetwork: Partial<CryptoAssetNetwork>
): Promise<CryptoAssetNetwork> {
  return fetchApi<CryptoAssetNetwork>('/admin/crypto/asset-networks', {
    method: 'POST',
    body: JSON.stringify(assetNetwork),
  });
}

export async function adminUpdateAssetNetwork(
  id: string,
  assetNetwork: Partial<CryptoAssetNetwork>
): Promise<CryptoAssetNetwork> {
  return fetchApi<CryptoAssetNetwork>(`/admin/crypto/asset-networks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(assetNetwork),
  });
}

export async function adminSetNetworkFee(
  fee: Partial<NetworkFee>
): Promise<NetworkFee> {
  return fetchApi<NetworkFee>('/admin/crypto/fees', {
    method: 'POST',
    body: JSON.stringify(fee),
  });
}

export async function adminGetBrandLevels(): Promise<BrandLevel[]> {
  return fetchApi<BrandLevel[]>('/admin/crypto/brand-levels');
}

export async function adminCreateBrandLevel(
  level: Partial<BrandLevel>
): Promise<BrandLevel> {
  return fetchApi<BrandLevel>('/admin/crypto/brand-levels', {
    method: 'POST',
    body: JSON.stringify(level),
  });
}

export async function adminUpdateBrandLevel(
  id: string,
  level: Partial<BrandLevel>
): Promise<BrandLevel> {
  return fetchApi<BrandLevel>(`/admin/crypto/brand-levels/${id}`, {
    method: 'PUT',
    body: JSON.stringify(level),
  });
}

export async function adminSetAdminWallet(
  wallet: Partial<AdminWalletAddress>
): Promise<AdminWalletAddress> {
  return fetchApi<AdminWalletAddress>('/admin/crypto/wallets', {
    method: 'POST',
    body: JSON.stringify(wallet),
  });
}

export async function adminGetAllDeposits(
  status?: string,
  page = 1,
  limit = 20
): Promise<{ deposits: CryptoDeposit[]; total: number; page: number; limit: number }> {
  const params = new URLSearchParams();
  if (status) params.append('status', status);
  params.append('page', page.toString());
  params.append('limit', limit.toString());

  return fetchApi<{ deposits: CryptoDeposit[]; total: number; page: number; limit: number }>(
    `/admin/crypto/deposits?${params}`
  );
}

export async function adminGetAllWithdrawals(
  status?: string,
  page = 1,
  limit = 20
): Promise<{ withdrawals: CryptoWithdrawal[]; total: number; page: number; limit: number }> {
  const params = new URLSearchParams();
  if (status) params.append('status', status);
  params.append('page', page.toString());
  params.append('limit', limit.toString());

  return fetchApi<{ withdrawals: CryptoWithdrawal[]; total: number; page: number; limit: number }>(
    `/admin/crypto/withdrawals?${params}`
  );
}

export async function adminConfirmDeposit(id: string): Promise<void> {
  return fetchApi<void>(`/admin/crypto/deposits/${id}/confirm`, {
    method: 'POST',
  });
}

export async function adminConfirmWithdrawal(
  id: string,
  txHash: string
): Promise<void> {
  return fetchApi<void>(`/admin/crypto/withdrawals/${id}/confirm`, {
    method: 'POST',
    body: JSON.stringify({ tx_hash: txHash }),
  });
}

export async function adminCancelWithdrawal(id: string): Promise<void> {
  return fetchApi<void>(`/admin/crypto/withdrawals/${id}/cancel`, {
    method: 'POST',
  });
}

export async function adminUpdateUserBrandLevel(
  userId: string,
  level: number
): Promise<void> {
  return fetchApi<void>(`/admin/users/${userId}/brand-level`, {
    method: 'PUT',
    body: JSON.stringify({ level }),
  });
}
