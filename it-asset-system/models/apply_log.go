package models

import (
	"gorm.io/gorm"
)

type SysApplyLog struct {
	gorm.Model
	UserID       uint     `gorm:"index;not null" json:"user_id"`
	User         SysUser  `gorm:"foreignKey:UserID" json:"user"`
	AssetID      uint     `gorm:"index;not null" json:"asset_id" binding:"required"`
	Asset        SysAsset `gorm:"foreignKey:AssetID" json:"asset"`
	Type         int      `gorm:"type:int;not null;comment:1:领用, 2:归还" json:"type" binding:"required"`
	Status       int      `gorm:"type:int;not null;default:0;comment:0:待审批, 1:已通过, 2:已驳回" json:"status"`
	Reason       string   `gorm:"type:varchar(255)" json:"reason"`
	RejectReason string   `gorm:"type:varchar(255)" json:"reject_reason"`
}
