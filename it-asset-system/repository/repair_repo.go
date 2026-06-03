package repository

import (
	"errors"
	"it-asset-system/core"
	"it-asset-system/models"
)

// ConfirmRepair 确认送修核心事务
// 将报修单状态改为"已送修(1)"，将对应单件资产状态改为"维修中(2)"
func ConfirmRepair(repairID uint, assetID uint) error {
	tx := core.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 步骤 1: 报修单状态 -> "1: 已送修"
	if err := tx.Model(&models.SysRepairLog{}).Where("id = ?", repairID).Update("status", 1).Error; err != nil {
		tx.Rollback()
		return errors.New("更新报修单状态失败")
	}

	// 步骤 2: 将该单件资产状态直接置为"维修中(2)"
	if err := tx.Model(&models.SysAsset{}).Where("id = ?", assetID).Update("status", 2).Error; err != nil {
		tx.Rollback()
		return errors.New("更新资产状态为维修中失败")
	}

	return tx.Commit().Error
}

// CompleteRepair 完成维修核心事务
// 将报修单状态改为"已完成(2)"。若设备在送修前有持有人（UserID > 0），则恢复为"使用中(1)"并保留持有人；否则恢复为"闲置(0)"并清空持有人。
func CompleteRepair(repairID uint, assetID uint) error {
	tx := core.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 步骤 1: 报修单状态 -> "2: 已完成"
	if err := tx.Model(&models.SysRepairLog{}).Where("id = ?", repairID).Update("status", 2).Error; err != nil {
		tx.Rollback()
		return errors.New("更新报修单状态失败")
	}

	// 步骤 2: 获取该资产单件信息以判断是否有持有人
	var asset models.SysAsset
	if err := tx.First(&asset, assetID).Error; err != nil {
		tx.Rollback()
		return errors.New("找不到对应的资产")
	}

	targetStatus := 0
	targetUserID := uint(0)
	if asset.UserID > 0 {
		targetStatus = 1
		targetUserID = asset.UserID
	}

	// 步骤 3: 更新资产状态与归属
	if err := tx.Model(&models.SysAsset{}).Where("id = ?", assetID).Updates(map[string]interface{}{
		"status":  targetStatus,
		"user_id": targetUserID,
	}).Error; err != nil {
		tx.Rollback()
		return errors.New("更新资产状态失败")
	}

	return tx.Commit().Error
}
