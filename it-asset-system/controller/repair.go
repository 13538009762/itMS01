package controller

import (
	"github.com/gin-gonic/gin"
	"it-asset-system/core"
	"it-asset-system/models"
	"it-asset-system/repository"
	"it-asset-system/utils"
)

type SubmitRepairRequest struct {
	AssetID uint   `json:"asset_id" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

func SubmitRepair(c *gin.Context) {
	var req SubmitRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters: "+err.Error())
		return
	}

	userID, _ := c.Get("userID")
	repair := models.SysRepairLog{
		AssetID: req.AssetID,
		UserID:  userID.(uint),
		Reason:  req.Reason,
		Status:  0,
	}

	if err := core.DB.Create(&repair).Error; err != nil {
		utils.Error(c, 500, "Failed to submit repair")
		return
	}

	utils.Success(c, repair)
}

type ConfirmRepairRequest struct {
	RepairID uint `json:"repair_id" binding:"required"`
	AssetID  uint `json:"asset_id" binding:"required"`
}

func ConfirmRepair(c *gin.Context) {
	var req ConfirmRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters")
		return
	}

	if err := repository.ConfirmRepair(req.RepairID, req.AssetID); err != nil {
		utils.Error(c, 500, "Failed to confirm repair: "+err.Error())
		return
	}

	utils.Success(c, "Confirmed repair successfully")
}

func GetRepairList(c *gin.Context) {
	var list []models.SysRepairLog
	err := core.DB.Preload("User").Preload("Asset").Order("status ASC, id DESC").Find(&list).Error
	if err != nil {
		utils.Error(c, 500, "Failed to fetch repair list")
		return
	}
	utils.Success(c, list)
}

func CompleteRepair(c *gin.Context) {
	var req ConfirmRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters: "+err.Error())
		return
	}

	if err := repository.CompleteRepair(req.RepairID, req.AssetID); err != nil {
		utils.Error(c, 500, "Failed to complete repair: "+err.Error())
		return
	}

	utils.Success(c, "Completed repair successfully")
}
