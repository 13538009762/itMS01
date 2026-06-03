package controller

import (
	"github.com/gin-gonic/gin"
	"it-asset-system/core"
	"it-asset-system/models"
	"it-asset-system/utils"
)

type SubmitApplyRequest struct {
	AssetID uint   `json:"asset_id" binding:"required"`
	Type    int    `json:"type" binding:"required"`
	Reason  string `json:"reason"`
}

func SubmitApply(c *gin.Context) {
	var req SubmitApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters: "+err.Error())
		return
	}

	userID, _ := c.Get("userID")

	// If it's a claim application (Type == 1), auto-allocate the next available idle unit in sequential order
	if req.Type == 1 {
		var reqAsset models.SysAsset
		if err := core.DB.First(&reqAsset, req.AssetID).Error; err != nil {
			utils.Error(c, 404, "Asset category not found")
			return
		}

		baseNo := reqAsset.BaseNo
		if baseNo == "" {
			baseNo = reqAsset.AssetNo
		}

		var availableAsset models.SysAsset
		subQuery := core.DB.Model(&models.SysApplyLog{}).Select("asset_id").Where("type = 1 AND status = 0")
		err := core.DB.Where("base_no = ? AND status = 0 AND id NOT IN (?)", baseNo, subQuery).Order("asset_no ASC").First(&availableAsset).Error

		if err != nil {
			utils.Error(c, 400, "该批次资产当前已无闲置，或已被其他申请预约完")
			return
		}

		req.AssetID = availableAsset.ID
	}

	apply := models.SysApplyLog{
		UserID:  userID.(uint),
		AssetID: req.AssetID,
		Type:    req.Type,
		Status:  0,
		Reason:  req.Reason,
	}

	if err := core.DB.Create(&apply).Error; err != nil {
		utils.Error(c, 500, "Failed to submit apply")
		return
	}

	utils.Success(c, apply)
}

type AuditRequest struct {
	ID           uint   `json:"id" binding:"required"`
	Status       int    `json:"status" binding:"required"` // 1:已通过, 2:已驳回
	RejectReason string `json:"reject_reason"`
}

// AuditApply 审批申请（领用或归还）
// 每个 AssetID 对应单台物理设备，审批通过后直接翻转该设备的状态。
func AuditApply(c *gin.Context) {
	var req AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters")
		return
	}

	tx := core.DB.Begin()

	var apply models.SysApplyLog
	if err := tx.First(&apply, req.ID).Error; err != nil {
		tx.Rollback()
		utils.Error(c, 404, "Apply record not found")
		return
	}

	apply.Status = req.Status
	if req.Status == 2 {
		apply.RejectReason = req.RejectReason
	}
	if err := tx.Save(&apply).Error; err != nil {
		tx.Rollback()
		utils.Error(c, 500, "Failed to update apply status")
		return
	}

	// 审批通过 + 领用申请 → 该设备状态改为"使用中"，绑定持有人
	if req.Status == 1 && apply.Type == 1 {
		var asset models.SysAsset
		if err := tx.First(&asset, apply.AssetID).Error; err != nil {
			tx.Rollback()
			utils.Error(c, 404, "Asset not found")
			return
		}
		if asset.Status != 0 {
			tx.Rollback()
			utils.Error(c, 400, "Asset is not available (not idle)")
			return
		}
		asset.Status = 1         // 使用中
		asset.UserID = apply.UserID // 绑定持有人
		if err := tx.Save(&asset).Error; err != nil {
			tx.Rollback()
			utils.Error(c, 500, "Failed to update asset status")
			return
		}
	} else if req.Status == 1 && apply.Type == 2 {
		// 审批通过 + 归还申请 → 该设备状态改回"闲置"，清空持有人
		var asset models.SysAsset
		if err := tx.First(&asset, apply.AssetID).Error; err != nil {
			tx.Rollback()
			utils.Error(c, 404, "Asset not found")
			return
		}
		asset.Status = 0 // 闲置
		asset.UserID = 0 // 清空持有人
		if err := tx.Save(&asset).Error; err != nil {
			tx.Rollback()
			utils.Error(c, 500, "Failed to update asset status")
			return
		}
	}

	tx.Commit()

	utils.Success(c, "Audited successfully")
}

func GetApplyList(c *gin.Context) {
	var list []models.SysApplyLog
	err := core.DB.Preload("User").Preload("Asset").Order("status ASC, id DESC").Find(&list).Error
	if err != nil {
		utils.Error(c, 500, "Failed to fetch apply list")
		return
	}
	utils.Success(c, list)
}

func GetMyApplyList(c *gin.Context) {
	userID, _ := c.Get("userID")
	var list []models.SysApplyLog
	err := core.DB.Preload("Asset").Where("user_id = ?", userID).Order("id DESC").Find(&list).Error
	if err != nil {
		utils.Error(c, 500, "Failed to fetch my apply list")
		return
	}
	utils.Success(c, list)
}
