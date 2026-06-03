package models

import (
	"gorm.io/gorm"
)

type SysRepairLog struct {
	gorm.Model
	AssetID uint     `gorm:"index;not null" json:"asset_id" binding:"required"`
	Asset   SysAsset `gorm:"foreignKey:AssetID" json:"asset"`
	UserID  uint     `gorm:"index;not null" json:"user_id"`
	User    SysUser  `gorm:"foreignKey:UserID" json:"user"`
	Reason  string   `gorm:"type:text;not null" json:"reason" binding:"required"`
	Status  int      `gorm:"type:int;not null;default:0;comment:0:待处理, 1:已送修, 2:已完成" json:"status"`
}
