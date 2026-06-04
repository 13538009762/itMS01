package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"it-asset-system/core"
	"it-asset-system/models"
	"it-asset-system/utils"
)

// GetAvailableAssets 获取资产列表
// 所有角色均返回全部资产，以便前端正确统计总数量和闲置数量
func GetAvailableAssets(c *gin.Context) {
	var assets []models.SysAsset

	// 统一返回全部资产，按状态升序（闲置优先）再按ID降序
	if err := core.DB.Order("status ASC, id DESC").Find(&assets).Error; err != nil {
		utils.Error(c, 500, "Failed to get assets")
		return
	}

	// 批量填充持有人姓名：收集所有非零的 user_id，一次性查询对应用户
	userIDSet := map[uint]bool{}
	for _, a := range assets {
		if a.UserID > 0 {
			userIDSet[a.UserID] = true
		}
	}
	if len(userIDSet) > 0 {
		var userIDs []uint
		for id := range userIDSet {
			userIDs = append(userIDs, id)
		}
		var users []models.SysUser
		core.DB.Where("id IN ?", userIDs).Find(&users)
		userMap := make(map[uint]models.SysUser)
		for _, u := range users {
			userMap[u.ID] = u
		}
		for i := range assets {
			if u, ok := userMap[assets[i].UserID]; ok {
				if u.RealName != "" {
					assets[i].UserName = u.RealName
				} else {
					assets[i].UserName = u.Username
				}
			}
		}
	}

	utils.Success(c, assets)
}

// CreateAssetRequest 添加资产请求体
type CreateAssetRequest struct {
	BaseNo     string `json:"base_no" binding:"required"`    // 基础编号，如 AST20260601001
	Name       string `json:"name" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Quantity   int    `json:"quantity"` // 数量，默认 1，将生成 Quantity 条独立记录
}

// CreateAsset 批量创建资产单件
// POST /admin/asset
// 若 Quantity=3，则自动生成：BaseNo-001, BaseNo-002, BaseNo-003
func CreateAsset(c *gin.Context) {
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters: "+err.Error())
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	tx := core.DB.Begin()

	var created []models.SysAsset
	for i := 1; i <= req.Quantity; i++ {
		unit := models.SysAsset{
			BaseNo:     req.BaseNo,
			AssetNo:    fmt.Sprintf("%s-%03d", req.BaseNo, i),
			Name:       req.Name,
			CategoryID: req.CategoryID,
			Status:     0,
		}
		if err := tx.Create(&unit).Error; err != nil {
			tx.Rollback()
			utils.Error(c, 500, fmt.Sprintf("Failed to create unit %d: %v", i, err))
			return
		}
		created = append(created, unit)
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, 500, "Transaction commit failed")
		return
	}

	utils.Success(c, created)
}

// UpdateAsset 更新单件资产信息及其同批次数量（名称、分类、数量）
// PUT /admin/asset
func UpdateAsset(c *gin.Context) {
	var req struct {
		ID         uint   `json:"id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
		Quantity   int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters: "+err.Error())
		return
	}

	var asset models.SysAsset
	if err := core.DB.First(&asset, req.ID).Error; err != nil {
		utils.Error(c, 404, "Asset not found")
		return
	}

	baseNo := asset.BaseNo
	if baseNo == "" {
		baseNo = asset.AssetNo
		asset.BaseNo = baseNo
		core.DB.Save(&asset)
	}

	var batch []models.SysAsset
	core.DB.Where("base_no = ?", baseNo).Find(&batch)

	currentTotal := len(batch)
	borrowedCount := 0
	var idleAssets []models.SysAsset
	for _, a := range batch {
		if a.Status != 0 {
			borrowedCount++
		} else {
			idleAssets = append(idleAssets, a)
		}
	}

	if req.Quantity > 0 && req.Quantity != currentTotal {
		if req.Quantity < borrowedCount {
			utils.Error(c, 400, fmt.Sprintf("数量不能少于被借走/维修中的数量（当前有 %d 个非闲置状态）", borrowedCount))
			return
		}

		if req.Quantity > currentTotal {
			// create new units
			maxSuffix := 0
			for _, a := range batch {
				var suffix int
				fmt.Sscanf(a.AssetNo, baseNo+"-%03d", &suffix)
				if suffix > maxSuffix {
					maxSuffix = suffix
				}
			}
			if maxSuffix == 0 {
				maxSuffix = currentTotal
			}

			toCreate := req.Quantity - currentTotal
			for i := 1; i <= toCreate; i++ {
				maxSuffix++
				newUnit := models.SysAsset{
					BaseNo:     baseNo,
					AssetNo:    fmt.Sprintf("%s-%03d", baseNo, maxSuffix),
					Name:       req.Name,
					CategoryID: req.CategoryID,
					Status:     0,
				}
				core.DB.Create(&newUnit)
			}
		} else {
			// delete idle units
			toDelete := currentTotal - req.Quantity
			for i := 0; i < toDelete; i++ {
				core.DB.Delete(&idleAssets[i])
			}
		}
	}

	// Update Name and Category for all units in batch
	if err := core.DB.Model(&models.SysAsset{}).Where("base_no = ?", baseNo).Updates(map[string]interface{}{
		"name":        req.Name,
		"category_id": req.CategoryID,
	}).Error; err != nil {
		utils.Error(c, 500, "Failed to update asset batch")
		return
	}

	utils.Success(c, asset)
}

// DeleteAsset 删除单个资产单件
// DELETE /admin/asset
func DeleteAsset(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Error(c, 400, "Missing asset ID")
		return
	}

	var asset models.SysAsset
	if err := core.DB.First(&asset, id).Error; err != nil {
		utils.Error(c, 404, "Asset not found")
		return
	}

	if asset.Status != 0 {
		utils.Error(c, 400, "不能删除非闲置状态的资产")
		return
	}

	if err := core.DB.Delete(&asset).Error; err != nil {
		utils.Error(c, 500, "Failed to delete asset")
		return
	}

	utils.Success(c, "Deleted successfully")
}

func GetCategories(c *gin.Context) {
	var categories []models.SysCategory
	if err := core.DB.Find(&categories).Error; err != nil {
		utils.Error(c, 500, "Failed to get categories")
		return
	}
	utils.Success(c, categories)
}

// GetHeldAssets 获取当前登录用户持有的资产列表
// 通过 status=1 且 user_id=当前用户 直接查询，无需日志聚合
func GetHeldAssets(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	var assets []models.SysAsset
	if err := core.DB.Where("user_id = ? AND status IN (1, 2)", userID).Order("id DESC").Find(&assets).Error; err != nil {
		utils.Error(c, 500, "Failed to fetch held assets")
		return
	}

	if len(assets) == 0 {
		utils.Success(c, assets)
		return
	}

	// 收集资产ID列表用于查询日志
	var assetIDs []uint
	for _, a := range assets {
		assetIDs = append(assetIDs, a.ID)
	}

	// 查询待归还的申请（type=2, status=0 待审批）
	var pendingReturns []models.SysApplyLog
	core.DB.Where("user_id = ? AND type = 2 AND status = 0 AND asset_id IN ?", userID, assetIDs).Find(&pendingReturns)
	pendingReturnMap := make(map[uint]bool)
	for _, r := range pendingReturns {
		pendingReturnMap[r.AssetID] = true
	}

	// 查询待报修 / 维修中的报修单（status 0=待确认, 1=已送修）
	var activeRepairs []models.SysRepairLog
	core.DB.Where("user_id = ? AND status IN (0, 1) AND asset_id IN ?", userID, assetIDs).Find(&activeRepairs)
	repairStatusMap := make(map[uint]int)
	for _, r := range activeRepairs {
		repairStatusMap[r.AssetID] = r.Status
	}

	// 填充 HoldingStatus 瞬态字段
	for i := range assets {
		assetID := assets[i].ID
		if assets[i].Status == 2 {
			assets[i].HoldingStatus = "维修中"
		} else if pendingReturnMap[assetID] {
			assets[i].HoldingStatus = "待归还"
		} else if status, ok := repairStatusMap[assetID]; ok {
			if status == 0 {
				assets[i].HoldingStatus = "待报修"
			} else {
				assets[i].HoldingStatus = "维修中"
			}
		} else {
			assets[i].HoldingStatus = "使用中"
		}
	}

	utils.Success(c, assets)
}
