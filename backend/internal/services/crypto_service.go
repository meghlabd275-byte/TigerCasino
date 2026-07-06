package services

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/tigercasino/backend/internal/models"
	"gorm.io/gorm"
)

type CryptoService struct {
	db *gorm.DB
}

func NewCryptoService(db *gorm.DB) *CryptoService {
	return &CryptoService{db: db}
}

// GetAllNetworks returns all active crypto networks
func (s *CryptoService) GetAllNetworks() ([]models.CryptoNetwork, error) {
	var networks []models.CryptoNetwork
	err := s.db.Where("is_active = ?", true).Find(&networks).Error
	return networks, err
}

// GetNetworkByID returns a network by ID
func (s *CryptoService) GetNetworkByID(id uuid.UUID) (*models.CryptoNetwork, error) {
	var network models.CryptoNetwork
	err := s.db.First(&network, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// GetAllAssets returns all active crypto assets
func (s *CryptoService) GetAllAssets() ([]models.CryptoAsset, error) {
	var assets []models.CryptoAsset
	err := s.db.Where("is_active = ?", true).Find(&assets).Error
	return assets, err
}

// GetAssetByID returns an asset by ID
func (s *CryptoService) GetAssetByID(id uuid.UUID) (*models.CryptoAsset, error) {
	var asset models.CryptoAsset
	err := s.db.First(&asset, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAssetBySymbol returns an asset by symbol
func (s *CryptoService) GetAssetBySymbol(symbol string) (*models.CryptoAsset, error) {
	var asset models.CryptoAsset
	err := s.db.First(&asset, "symbol = ? AND is_active = ?", symbol, true).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAssetsByNetwork returns all assets available on a specific network
func (s *CryptoService) GetAssetsByNetwork(networkID uuid.UUID) ([]models.CryptoAssetNetwork, error) {
	var assetNetworks []models.CryptoAssetNetwork
	err := s.db.Preload("Asset").Preload("Network").
		Where("network_id = ?", networkID).
		Find(&assetNetworks).Error
	return assetNetworks, err
}

// GetNetworksByAsset returns all networks available for a specific asset
func (s *CryptoService) GetNetworksByAsset(assetID uuid.UUID) ([]models.CryptoAssetNetwork, error) {
	var assetNetworks []models.CryptoAssetNetwork
	err := s.db.Preload("Network").
		Where("asset_id = ?", assetID).
		Find(&assetNetworks).Error
	return assetNetworks, err
}

// GetAssetWithNetworks returns an asset with all its network information
func (s *CryptoService) GetAssetWithNetworks(assetID uuid.UUID) (*models.CryptoAsset, []models.CryptoAssetNetwork, error) {
	asset, err := s.GetAssetByID(assetID)
	if err != nil {
		return nil, nil, err
	}

	assetNetworks, err := s.GetNetworksByAsset(assetID)
	if err != nil {
		return nil, nil, err
	}

	return asset, assetNetworks, nil
}

// GetAllAssetsWithNetworks returns all assets with their network information
func (s *CryptoService) GetAllAssetsWithNetworks() ([]models.CryptoAssetWithNetwork, error) {
	var result []models.CryptoAssetWithNetwork

	rows, err := s.db.Raw(`
		SELECT 
			ca.id, ca.name, ca.symbol, ca.decimals, ca.contract_address,
			ca.is_active, ca.min_deposit_amount, ca.min_withdrawal_amount,
			can.id as network_id,
			can.name as network_name, can.chain_id, can.symbol as network_symbol,
			can.explorer_url, can.is_withdrawal_enabled, can.is_deposit_enabled,
			can.min_confirmation_blocks,
			can.deposit_enabled, can.withdrawal_enabled, can.contract_address as asset_contract_address,
			nf.deposit_fee, nf.withdrawal_fee, nf.deposit_fee_percent, nf.withdrawal_fee_percent
		FROM crypto_assets ca
		JOIN crypto_asset_networks can ON ca.id = can.asset_id
		JOIN crypto_networks cn ON can.network_id = cn.id
		LEFT JOIN network_fees nf ON ca.id = nf.asset_id AND cn.id = nf.network_id
		WHERE ca.is_active = true AND cn.is_active = true
		ORDER BY ca.symbol, cn.name
	`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assetMap := make(map[string]models.CryptoAssetWithNetwork)

	for rows.Next() {
		var asset models.CryptoAsset
		var networkID uuid.UUID
		var networkName, chainID, networkSymbol, explorerURL string
		var depositEnabled, withdrawalEnabled bool
		var minConfirmations int
		var contractAddress string
		var depositFee, withdrawalFee, depositFeePercent, withdrawalFeePercent float64

		err := rows.Scan(
			&asset.ID, &asset.Name, &asset.Symbol, &asset.Decimals, &asset.ContractAddress,
			&asset.IsActive, &asset.MinDepositAmount, &asset.MinWithdrawalAmount,
			&networkID, &networkName, &chainID, &networkSymbol, &explorerURL,
			&depositEnabled, &withdrawalEnabled, &minConfirmations,
			&depositEnabled, &withdrawalEnabled, &contractAddress,
			&depositFee, &withdrawalFee, &depositFeePercent, &withdrawalFeePercent,
		)
		if err != nil {
			continue
		}

		key := asset.Symbol
		if existing, ok := assetMap[key]; ok {
			continue
		}

		assetMap[key] = models.CryptoAssetWithNetwork{
			CryptoAsset:               asset,
			NetworkID:                 networkID,
			DepositEnabled:            depositEnabled,
			WithdrawalEnabled:         withdrawalEnabled,
			ContractAddress:           contractAddress,
			DepositFee:                depositFee,
			WithdrawalFee:             withdrawalFee,
			DepositFeePercent:         depositFeePercent,
			WithdrawalFeePercent:      withdrawalFeePercent,
			MinDepositAmount:          asset.MinDepositAmount,
			MinWithdrawalAmount:       asset.MinWithdrawalAmount,
		}
	}

	for _, v := range assetMap {
		result = append(result, v)
	}

	return result, nil
}

// CalculateDepositFee calculates the deposit fee for a given amount
func (s *CryptoService) CalculateDepositFee(assetID, networkID uuid.UUID, amount float64, userLevel int) (float64, error) {
	var networkFee models.NetworkFee
	err := s.db.First(&networkFee, "asset_id = ? AND network_id = ?", assetID, networkID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // No fee
		}
		return 0, err
	}

	var brandLevel models.BrandLevel
	err = s.db.First(&brandLevel, "level = ? AND is_active = ?", userLevel, true).Error
	discountPercent := 0.0
	if err == nil {
		discountPercent = brandLevel.DepositFeeDiscountPercent
	}

	fee := networkFee.DepositFee
	feePercent := networkFee.DepositFeePercent * (1 - discountPercent/100)
	fee += amount * feePercent / 100

	return math.Round(fee*100000000) / 100000000, nil
}

// CalculateWithdrawalFee calculates the withdrawal fee for a given amount
func (s *CryptoService) CalculateWithdrawalFee(assetID, networkID uuid.UUID, amount float64, userLevel int) (float64, error) {
	var networkFee models.NetworkFee
	err := s.db.First(&networkFee, "asset_id = ? AND network_id = ?", assetID, networkID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // No fee
		}
		return 0, err
	}

	var brandLevel models.BrandLevel
	err = s.db.First(&brandLevel, "level = ? AND is_active = ?", userLevel, true).Error
	discountPercent := 0.0
	if err == nil {
		discountPercent = brandLevel.WithdrawalFeeDiscountPercent
	}

	fee := networkFee.WithdrawalFee
	feePercent := networkFee.WithdrawalFeePercent * (1 - discountPercent/100)
	fee += amount * feePercent / 100

	return math.Round(fee*100000000) / 100000000, nil
}

// GetDepositAddress generates or retrieves a deposit address for a user
func (s *CryptoService) GetDepositAddress(userID, assetID, networkID uuid.UUID) (string, error) {
	// Check if asset/network is valid and deposit is enabled
	var assetNetwork models.CryptoAssetNetwork
	err := s.db.First(&assetNetwork, "asset_id = ? AND network_id = ? AND deposit_enabled = ?", assetID, networkID, true).Error
	if err != nil {
		return "", errors.New("deposit not enabled for this asset on this network")
	}

	// Get admin wallet address for this asset/network
	var adminWallet models.AdminWalletAddress
	err = s.db.First(&adminWallet, "asset_id = ? AND network_id = ? AND is_active = ?", assetID, networkID, true).Error
	if err != nil {
		return "", errors.New("no admin wallet configured for this asset/network")
	}

	// In production, you would generate unique deposit addresses per user
	// For now, return the admin wallet address with a prefix
	return adminWallet.Address, nil
}

// CreateDeposit creates a new deposit record
func (s *CryptoService) CreateDeposit(userID, assetID, networkID uuid.UUID, amount float64, address, txHash string) (*models.CryptoDeposit, error) {
	// Get user's brand level
	var user models.User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	// Calculate fee
	fee, err := s.CalculateDepositFee(assetID, networkID, amount, user.VIPLevel)
	if err != nil {
		return nil, err
	}

	netAmount := amount - fee

	deposit := models.CryptoDeposit{
		UserID:       userID,
		AssetID:     assetID,
		NetworkID:   networkID,
		Amount:      amount,
		Fee:         fee,
		NetAmount:   netAmount,
		Address:     address,
		TxHash:      txHash,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	err = s.db.Create(&deposit).Error
	if err != nil {
		return nil, err
	}

	return &deposit, nil
}

// ConfirmDeposit confirms a deposit after sufficient confirmations
func (s *CryptoService) ConfirmDeposit(depositID uuid.UUID) error {
	var deposit models.CryptoDeposit
	err := s.db.First(&deposit, "id = ?", depositID).Error
	if err != nil {
		return err
	}

	if deposit.Status != "pending" {
		return errors.New("deposit already processed")
	}

	now := time.Now()
	deposit.Status = "completed"
	deposit.ProcessedAt = &now

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Update deposit status
		if err := tx.Save(&deposit).Error; err != nil {
			return err
		}

		// Add funds to user balance
		if err := tx.Model(&models.User{}).
			Where("id = ?", deposit.UserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", deposit.NetAmount)).Error; err != nil {
			return err
		}

		// Create transaction record
		transaction := models.Transaction{
			UserID:     deposit.UserID,
			Type:       "deposit",
			Amount:     deposit.NetAmount,
			Currency:   deposit.AssetID.String(),
			Status:     "completed",
			TxHash:     deposit.TxHash,
			Address:    deposit.Address,
			Fee:        deposit.Fee,
			ProcessedAt: &now,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// CreateWithdrawal creates a new withdrawal request
func (s *CryptoService) CreateWithdrawal(userID, assetID, networkID uuid.UUID, amount float64, address string) (*models.CryptoWithdrawal, error) {
	// Get user
	var user models.User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	// Check if asset/network is valid and withdrawal is enabled
	var assetNetwork models.CryptoAssetNetwork
	err = s.db.First(&assetNetwork, "asset_id = ? AND network_id = ? AND withdrawal_enabled = ?", assetID, networkID, true).Error
	if err != nil {
		return nil, errors.New("withdrawal not enabled for this asset on this network")
	}

	// Get asset for min withdrawal amount
	asset, err := s.GetAssetByID(assetID)
	if err != nil {
		return nil, err
	}

	if amount < asset.MinWithdrawalAmount {
		return nil, errors.New("amount below minimum withdrawal limit")
	}

	// Calculate fee
	fee, err := s.CalculateWithdrawalFee(assetID, networkID, amount, user.VIPLevel)
	if err != nil {
		return nil, err
	}

	netAmount := amount - fee

	// Check balance
	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	withdrawal := models.CryptoWithdrawal{
		UserID:     userID,
		AssetID:   assetID,
		NetworkID: networkID,
		Amount:    amount,
		Fee:       fee,
		NetAmount: netAmount,
		Address:   address,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Deduct from balance
		if err := tx.Model(&models.User{}).
			Where("id = ? AND balance >= ?", userID, amount).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// Create withdrawal record
		if err := tx.Create(&withdrawal).Error; err != nil {
			return err
		}

		// Create transaction record
		transaction := models.Transaction{
			UserID:    userID,
			Type:      "withdrawal",
			Amount:    amount,
			Currency:  assetID.String(),
			Status:    "pending",
			Address:   address,
			Fee:       fee,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &withdrawal, nil
}

// ConfirmWithdrawal confirms a withdrawal after processing
func (s *CryptoService) ConfirmWithdrawal(withdrawalID uuid.UUID, txHash string) error {
	var withdrawal models.CryptoWithdrawal
	err := s.db.First(&withdrawal, "id = ?", withdrawalID).Error
	if err != nil {
		return err
	}

	if withdrawal.Status != "pending" {
		return errors.New("withdrawal already processed")
	}

	now := time.Now()
	withdrawal.Status = "completed"
	withdrawal.TxHash = txHash
	withdrawal.ProcessedAt = &now

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&withdrawal).Error; err != nil {
			return err
		}

		// Update transaction record
		if err := tx.Model(&models.Transaction{}).
			Where("user_id = ? AND type = 'withdrawal' AND status = 'pending'", withdrawal.UserID).
			Updates(map[string]interface{}{
				"status":     "completed",
				"tx_hash":    txHash,
				"processed_at": now,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// CancelWithdrawal cancels a pending withdrawal
func (s *CryptoService) CancelWithdrawal(withdrawalID uuid.UUID) error {
	var withdrawal models.CryptoWithdrawal
	err := s.db.First(&withdrawal, "id = ? AND status = ?", withdrawalID, "pending").Error
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Refund balance
		if err := tx.Model(&models.User{}).
			Where("id = ?", withdrawal.UserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", withdrawal.Amount)).Error; err != nil {
			return err
		}

		// Update withdrawal status
		now := time.Now()
		withdrawal.Status = "cancelled"
		withdrawal.ProcessedAt = &now
		if err := tx.Save(&withdrawal).Error; err != nil {
			return err
		}

		// Update transaction record
		if err := tx.Model(&models.Transaction{}).
			Where("user_id = ? AND type = 'withdrawal' AND status = 'pending'", withdrawal.UserID).
			Update("status", "cancelled").Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserDeposits returns all deposits for a user
func (s *CryptoService) GetUserDeposits(userID uuid.UUID) ([]models.CryptoDeposit, error) {
	var deposits []models.CryptoDeposit
	err := s.db.Preload("Asset").Preload("Network").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&deposits).Error
	return deposits, err
}

// GetUserWithdrawals returns all withdrawals for a user
func (s *CryptoService) GetUserWithdrawals(userID uuid.UUID) ([]models.CryptoWithdrawal, error) {
	var withdrawals []models.CryptoWithdrawal
	err := s.db.Preload("Asset").Preload("Network").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&withdrawals).Error
	return withdrawals, err
}

// Admin functions

// CreateNetwork creates a new crypto network
func (s *CryptoService) CreateNetwork(network *models.CryptoNetwork) error {
	return s.db.Create(network).Error
}

// UpdateNetwork updates a network
func (s *CryptoService) UpdateNetwork(network *models.CryptoNetwork) error {
	network.UpdatedAt = time.Now()
	return s.db.Save(network).Error
}

// CreateAsset creates a new crypto asset
func (s *CryptoService) CreateAsset(asset *models.CryptoAsset) error {
	return s.db.Create(asset).Error
}

// UpdateAsset updates a crypto asset
func (s *CryptoService) UpdateAsset(asset *models.CryptoAsset) error {
	asset.UpdatedAt = time.Now()
	return s.db.Save(asset).Error
}

// LinkAssetToNetwork links an asset to a network
func (s *CryptoService) LinkAssetToNetwork(assetNetwork *models.CryptoAssetNetwork) error {
	return s.db.Create(assetNetwork).Error
}

// UpdateAssetNetwork updates asset-network relationship
func (s *CryptoService) UpdateAssetNetwork(assetNetwork *models.CryptoAssetNetwork) error {
	assetNetwork.UpdatedAt = time.Now()
	return s.db.Save(assetNetwork).Error
}

// SetNetworkFee sets the fee for an asset-network combination
func (s *CryptoService) SetNetworkFee(fee *models.NetworkFee) error {
	var existingFee models.NetworkFee
	err := s.db.First(&existingFee, "asset_id = ? AND network_id = ?", fee.AssetID, fee.NetworkID).Error
	
	fee.UpdatedAt = time.Now()
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(fee).Error
	}
	
	existingFee.DepositFee = fee.DepositFee
	existingFee.WithdrawalFee = fee.WithdrawalFee
	existingFee.DepositFeePercent = fee.DepositFeePercent
	existingFee.WithdrawalFeePercent = fee.WithdrawalFeePercent
	existingFee.UpdatedAt = time.Now()
	
	return s.db.Save(&existingFee).Error
}

// CreateBrandLevel creates a new brand level
func (s *CryptoService) CreateBrandLevel(level *models.BrandLevel) error {
	return s.db.Create(level).Error
}

// UpdateBrandLevel updates a brand level
func (s *CryptoService) UpdateBrandLevel(level *models.BrandLevel) error {
	level.UpdatedAt = time.Now()
	return s.db.Save(level).Error
}

// GetAllBrandLevels returns all brand levels
func (s *CryptoService) GetAllBrandLevels() ([]models.BrandLevel, error) {
	var levels []models.BrandLevel
	err := s.db.Where("is_active = ?", true).Order("level").Find(&levels).Error
	return levels, err
}

// SetAdminWallet sets an admin wallet address for an asset-network
func (s *CryptoService) SetAdminWallet(wallet *models.AdminWalletAddress) error {
	var existingWallet models.AdminWalletAddress
	err := s.db.First(&existingWallet, "asset_id = ? AND network_id = ?", wallet.AssetID, wallet.NetworkID).Error

	wallet.UpdatedAt = time.Now()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(wallet).Error
	}

	existingWallet.Address = wallet.Address
	existingWallet.PrivateKeyEncrypted = wallet.PrivateKeyEncrypted
	existingWallet.IsActive = wallet.IsActive
	existingWallet.UpdatedAt = time.Now()

	return s.db.Save(&existingWallet).Error
}

// GetAllDeposits returns all deposits (admin)
func (s *CryptoService) GetAllDeposits(status string, limit, offset int) ([]models.CryptoDeposit, int64, error) {
	var deposits []models.CryptoDeposit
	var total int64

	query := s.db.Model(&models.CryptoDeposit{}).Preload("Asset").Preload("Network").Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&deposits).Error
	return deposits, total, err
}

// GetAllWithdrawals returns all withdrawals (admin)
func (s *CryptoService) GetAllWithdrawals(status string, limit, offset int) ([]models.CryptoWithdrawal, int64, error) {
	var withdrawals []models.CryptoWithdrawal
	var total int64

	query := s.db.Model(&models.CryptoWithdrawal{}).Preload("Asset").Preload("Network").Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&withdrawals).Error
	return withdrawals, total, err
}

// GetDepositByID returns a deposit by ID
func (s *CryptoService) GetDepositByID(id uuid.UUID) (*models.CryptoDeposit, error) {
	var deposit models.CryptoDeposit
	err := s.db.Preload("Asset").Preload("Network").Preload("User").
		First(&deposit, "id = ?", id).Error
	return &deposit, err
}

// GetWithdrawalByID returns a withdrawal by ID
func (s *CryptoService) GetWithdrawalByID(id uuid.UUID) (*models.CryptoWithdrawal, error) {
	var withdrawal models.CryptoWithdrawal
	err := s.db.Preload("Asset").Preload("Network").Preload("User").
		First(&withdrawal, "id = ?", id).Error
	return &withdrawal, err
}

// UpdateUserBrandLevel updates a user's brand/VIP level
func (s *CryptoService) UpdateUserBrandLevel(userID uuid.UUID, level int) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("vip_level", level).Error
}

// GetUserByID returns a user by ID
func (s *CryptoService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
