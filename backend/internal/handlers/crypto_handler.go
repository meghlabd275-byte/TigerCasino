package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tigercasino/backend/internal/models"
	"github.com/tigercasino/backend/internal/services"
)

type CryptoHandler struct {
	cryptoService *services.CryptoService
}

func NewCryptoHandler(cryptoService *services.CryptoService) *CryptoHandler {
	return &CryptoHandler{cryptoService: cryptoService}
}

// GetNetworks returns all available crypto networks
// @Summary Get all crypto networks
// @Description Returns all active crypto networks
// @Tags Crypto
// @Accept json
// @Produce json
// @Success 200 {array} models.CryptoNetwork
// @Router /api/crypto/networks [get]
func (h *CryptoHandler) GetNetworks(c *gin.Context) {
	networks, err := h.cryptoService.GetAllNetworks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, networks)
}

// GetAssets returns all available crypto assets with their networks
// @Summary Get all crypto assets
// @Description Returns all active crypto assets with network information
// @Tags Crypto
// @Accept json
// @Produce json
// @Success 200 {array} models.CryptoAssetWithNetwork
// @Router /api/crypto/assets [get]
func (h *CryptoHandler) GetAssets(c *gin.Context) {
	assets, err := h.cryptoService.GetAllAssetsWithNetworks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

// GetAssetNetworks returns all networks for a specific asset
// @Summary Get networks for an asset
// @Description Returns all networks available for a specific crypto asset
// @Tags Crypto
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {array} models.CryptoAssetNetwork
// @Router /api/crypto/assets/{id}/networks [get]
func (h *CryptoHandler) GetAssetNetworks(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	assetNetworks, err := h.cryptoService.GetNetworksByAsset(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assetNetworks)
}

// GetDepositAddress returns a deposit address for the user
// @Summary Get deposit address
// @Description Returns a deposit address for the specified asset and network
// @Tags Crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.DepositRequest true "Deposit request"
// @Success 200 {object} gin.H{"address": string}
// @Router /api/crypto/deposit/address [post]
func (h *CryptoHandler) GetDepositAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.cryptoService.GetDepositAddress(
		userID.(uuid.UUID),
		req.AssetID,
		req.NetworkID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

// CreateDeposit creates a new deposit record
// @Summary Create deposit
// @Description Creates a new deposit record (for manual deposits or webhooks)
// @Tags Crypto
// @Accept json
// @Produce json
// @Param request body models.DepositRequest true "Deposit request"
// @Success 200 {object} models.CryptoDeposit
// @Router /api/crypto/deposits [post]
func (h *CryptoHandler) CreateDeposit(c *gin.Context) {
	var req struct {
		UserID   uuid.UUID `json:"user_id" binding:"required"`
		AssetID  uuid.UUID `json:"asset_id" binding:"required"`
		NetworkID uuid.UUID `json:"network_id" binding:"required"`
		Amount   float64   `json:"amount" binding:"required,gt=0"`
		Address  string    `json:"address" binding:"required"`
		TxHash   string    `json:"tx_hash"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deposit, err := h.cryptoService.CreateDeposit(
		req.UserID,
		req.AssetID,
		req.NetworkID,
		req.Amount,
		req.Address,
		req.TxHash,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deposit)
}

// GetUserDeposits returns all deposits for the current user
// @Summary Get user deposits
// @Description Returns all deposits for the authenticated user
// @Tags Crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CryptoDeposit
// @Router /api/crypto/deposits [get]
func (h *CryptoHandler) GetUserDeposits(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deposits, err := h.cryptoService.GetUserDeposits(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deposits)
}

// CreateWithdrawal creates a new withdrawal request
// @Summary Create withdrawal
// @Description Creates a new withdrawal request
// @Tags Crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.WithdrawalRequest true "Withdrawal request"
// @Success 201 {object} models.CryptoWithdrawal
// @Router /api/crypto/withdrawals [post]
func (h *CryptoHandler) CreateWithdrawal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdrawal, err := h.cryptoService.CreateWithdrawal(
		userID.(uuid.UUID),
		req.AssetID,
		req.NetworkID,
		req.Amount,
		req.Address,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, withdrawal)
}

// GetUserWithdrawals returns all withdrawals for the current user
// @Summary Get user withdrawals
// @Description Returns all withdrawals for the authenticated user
// @Tags Crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CryptoWithdrawal
// @Router /api/crypto/withdrawals [get]
func (h *CryptoHandler) GetUserWithdrawals(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	withdrawals, err := h.cryptoService.GetUserWithdrawals(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

// CalculateFees calculates the deposit and withdrawal fees
// @Summary Calculate fees
// @Description Calculates deposit and withdrawal fees for a given amount
// @Tags Crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body struct{AssetID uuid.UUID `json:"asset_id";NetworkID uuid.UUID `json:"network_id";Amount float64 `json:"amount"`} true "Fee calculation request"
// @Success 200 {object} gin.H{"deposit_fee": float64, "withdrawal_fee": float64}
// @Router /api/crypto/fees [post]
func (h *CryptoHandler) CalculateFees(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		AssetID   uuid.UUID `json:"asset_id" binding:"required"`
		NetworkID uuid.UUID `json:"network_id" binding:"required"`
		Amount    float64   `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user for brand level
	user, err := h.cryptoService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	depositFee, err := h.cryptoService.CalculateDepositFee(req.AssetID, req.NetworkID, req.Amount, user.VIPLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	withdrawalFee, err := h.cryptoService.CalculateWithdrawalFee(req.AssetID, req.NetworkID, req.Amount, user.VIPLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deposit_fee":    depositFee,
		"withdrawal_fee": withdrawalFee,
	})
}

// ADMIN ENDPOINTS

// AdminGetAllNetworks returns all networks (including inactive)
// @Summary Get all networks (admin)
// @Description Returns all crypto networks including inactive ones
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CryptoNetwork
// @Router /api/admin/crypto/networks [get]
func (h *CryptoHandler) AdminGetAllNetworks(c *gin.Context) {
	networks, err := h.cryptoService.GetAllNetworks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, networks)
}

// AdminCreateNetwork creates a new crypto network
// @Summary Create network (admin)
// @Description Creates a new crypto network
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CryptoNetwork true "Network data"
// @Success 201 {object} models.CryptoNetwork
// @Router /api/admin/crypto/networks [post]
func (h *CryptoHandler) AdminCreateNetwork(c *gin.Context) {
	var network models.CryptoNetwork
	if err := c.ShouldBindJSON(&network); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.CreateNetwork(&network); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, network)
}

// AdminUpdateNetwork updates a crypto network
// @Summary Update network (admin)
// @Description Updates a crypto network
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Network ID"
// @Param request body models.CryptoNetwork true "Network data"
// @Success 200 {object} models.CryptoNetwork
// @Router /api/admin/crypto/networks/{id} [put]
func (h *CryptoHandler) AdminUpdateNetwork(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid network ID"})
		return
	}

	var network models.CryptoNetwork
	if err := c.ShouldBindJSON(&network); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network.ID = id
	if err := h.cryptoService.UpdateNetwork(&network); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, network)
}

// AdminGetAllAssets returns all assets (including inactive)
// @Summary Get all assets (admin)
// @Description Returns all crypto assets including inactive ones
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CryptoAsset
// @Router /api/admin/crypto/assets [get]
func (h *CryptoHandler) AdminGetAllAssets(c *gin.Context) {
	assets, err := h.cryptoService.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

// AdminCreateAsset creates a new crypto asset
// @Summary Create asset (admin)
// @Description Creates a new crypto asset
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CryptoAsset true "Asset data"
// @Success 201 {object} models.CryptoAsset
// @Router /api/admin/crypto/assets [post]
func (h *CryptoHandler) AdminCreateAsset(c *gin.Context) {
	var asset models.CryptoAsset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.CreateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

// AdminUpdateAsset updates a crypto asset
// @Summary Update asset (admin)
// @Description Updates a crypto asset
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Param request body models.CryptoAsset true "Asset data"
// @Success 200 {object} models.CryptoAsset
// @Router /api/admin/crypto/assets/{id} [put]
func (h *CryptoHandler) AdminUpdateAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset ID"})
		return
	}

	var asset models.CryptoAsset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.ID = id
	if err := h.cryptoService.UpdateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// AdminLinkAssetToNetwork links an asset to a network
// @Summary Link asset to network (admin)
// @Description Links a crypto asset to a blockchain network
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CryptoAssetNetwork true "Asset-Network data"
// @Success 201 {object} models.CryptoAssetNetwork
// @Router /api/admin/crypto/asset-networks [post]
func (h *CryptoHandler) AdminLinkAssetToNetwork(c *gin.Context) {
	var assetNetwork models.CryptoAssetNetwork
	if err := c.ShouldBindJSON(&assetNetwork); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.LinkAssetToNetwork(&assetNetwork); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, assetNetwork)
}

// AdminUpdateAssetNetwork updates asset-network relationship
// @Summary Update asset-network (admin)
// @Description Updates deposit/withdrawal enabled status for an asset on a network
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset-Network ID"
// @Param request body models.CryptoAssetNetwork true "Asset-Network data"
// @Success 200 {object} models.CryptoAssetNetwork
// @Router /api/admin/crypto/asset-networks/{id} [put]
func (h *CryptoHandler) AdminUpdateAssetNetwork(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset-network ID"})
		return
	}

	var assetNetwork models.CryptoAssetNetwork
	if err := c.ShouldBindJSON(&assetNetwork); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assetNetwork.ID = id
	if err := h.cryptoService.UpdateAssetNetwork(&assetNetwork); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assetNetwork)
}

// AdminSetNetworkFee sets the fee for an asset-network combination
// @Summary Set network fee (admin)
// @Description Sets deposit/withdrawal fees for an asset on a network
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.NetworkFee true "Fee data"
// @Success 200 {object} models.NetworkFee
// @Router /api/admin/crypto/fees [post]
func (h *CryptoHandler) AdminSetNetworkFee(c *gin.Context) {
	var fee models.NetworkFee
	if err := c.ShouldBindJSON(&fee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.SetNetworkFee(&fee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fee)
}

// AdminGetBrandLevels returns all brand levels
// @Summary Get brand levels (admin)
// @Description Returns all brand/VIP levels
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.BrandLevel
// @Router /api/admin/crypto/brand-levels [get]
func (h *CryptoHandler) AdminGetBrandLevels(c *gin.Context) {
	levels, err := h.cryptoService.GetAllBrandLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, levels)
}

// AdminCreateBrandLevel creates a new brand level
// @Summary Create brand level (admin)
// @Description Creates a new brand/VIP level with fee discounts
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BrandLevel true "Brand level data"
// @Success 201 {object} models.BrandLevel
// @Router /api/admin/crypto/brand-levels [post]
func (h *CryptoHandler) AdminCreateBrandLevel(c *gin.Context) {
	var level models.BrandLevel
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.CreateBrandLevel(&level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, level)
}

// AdminUpdateBrandLevel updates a brand level
// @Summary Update brand level (admin)
// @Description Updates a brand/VIP level (including 20% discount for White level)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand level ID"
// @Param request body models.BrandLevel true "Brand level data"
// @Success 200 {object} models.BrandLevel
// @Router /api/admin/crypto/brand-levels/{id} [put]
func (h *CryptoHandler) AdminUpdateBrandLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand level ID"})
		return
	}

	var level models.BrandLevel
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	level.ID = id
	if err := h.cryptoService.UpdateBrandLevel(&level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, level)
}

// AdminSetAdminWallet sets an admin wallet address
// @Summary Set admin wallet (admin)
// @Description Sets an admin wallet address for an asset-network
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AdminWalletAddress true "Wallet data"
// @Success 200 {object} models.AdminWalletAddress
// @Router /api/admin/crypto/wallets [post]
func (h *CryptoHandler) AdminSetAdminWallet(c *gin.Context) {
	var wallet models.AdminWalletAddress
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.SetAdminWallet(&wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// AdminGetAllDeposits returns all deposits
// @Summary Get all deposits (admin)
// @Description Returns all deposits with optional status filter
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} gin.H{"deposits": []models.CryptoDeposit, "total": int64, "page": int, "limit": int}
// @Router /api/admin/crypto/deposits [get]
func (h *CryptoHandler) AdminGetAllDeposits(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	deposits, total, err := h.cryptoService.GetAllDeposits(status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deposits": deposits,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// AdminGetAllWithdrawals returns all withdrawals
// @Summary Get all withdrawals (admin)
// @Description Returns all withdrawals with optional status filter
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} gin.H{"withdrawals": []models.CryptoWithdrawal, "total": int64, "page": int, "limit": int}
// @Router /api/admin/crypto/withdrawals [get]
func (h *CryptoHandler) AdminGetAllWithdrawals(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	withdrawals, total, err := h.cryptoService.GetAllWithdrawals(status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"withdrawals": withdrawals,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

// AdminConfirmDeposit confirms a pending deposit
// @Summary Confirm deposit (admin)
// @Description Confirms a pending deposit and credits user balance
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Deposit ID"
// @Success 200 {object} gin.H{"message": string}
// @Router /api/admin/crypto/deposits/{id}/confirm [post]
func (h *CryptoHandler) AdminConfirmDeposit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deposit ID"})
		return
	}

	if err := h.cryptoService.ConfirmDeposit(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deposit confirmed"})
}

// AdminConfirmWithdrawal confirms a pending withdrawal
// @Summary Confirm withdrawal (admin)
// @Description Confirms a pending withdrawal after processing on blockchain
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Withdrawal ID"
// @Param request body struct{TxHash string `json:"tx_hash"`} true "Transaction hash"
// @Success 200 {object} gin.H{"message": string}
// @Router /api/admin/crypto/withdrawals/{id}/confirm [post]
func (h *CryptoHandler) AdminConfirmWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid withdrawal ID"})
		return
	}

	var req struct {
		TxHash string `json:"tx_hash" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.ConfirmWithdrawal(id, req.TxHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "withdrawal confirmed"})
}

// AdminCancelWithdrawal cancels a pending withdrawal
// @Summary Cancel withdrawal (admin)
// @Description Cancels a pending withdrawal and refunds user balance
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Withdrawal ID"
// @Success 200 {object} gin.H{"message": string}
// @Router /api/admin/crypto/withdrawals/{id}/cancel [post]
func (h *CryptoHandler) AdminCancelWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid withdrawal ID"})
		return
	}

	if err := h.cryptoService.CancelWithdrawal(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "withdrawal cancelled"})
}

// AdminUpdateUserBrandLevel updates a user's brand level
// @Summary Update user brand level (admin)
// @Description Updates a user's brand/VIP level for fee discounts
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body struct{Level int `json:"level"`} true "Brand level"
// @Success 200 {object} gin.H{"message": string}
// @Router /api/admin/users/{id}/brand-level [put]
func (h *CryptoHandler) AdminUpdateUserBrandLevel(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Level int `json:"level" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cryptoService.UpdateUserBrandLevel(userID, req.Level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "brand level updated"})
}
