package models

import (
	"gorm.io/gorm"
)

// SysAsset 代表一台独立的物理资产单件。
// 每台设备都有唯一的 AssetNo（如 AST20260601001-002）。
// BaseNo 是管理员输入的基础编号（如 AST20260601001），同一批次的多台设备共享同一 BaseNo。
type SysAsset struct {
	gorm.Model
	AssetNo    string `gorm:"type:varchar(60);uniqueIndex;not null" json:"asset_no"`
	BaseNo     string `gorm:"type:varchar(50);index;not null" json:"base_no"`
	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	CategoryID uint   `gorm:"index;not null" json:"category_id"`
	// Status: 0=闲置, 1=使用中, 2=维修中, 3=已报废
	Status        int    `gorm:"type:int;not null;default:0;comment:0:闲置,1:使用中,2:维修中,3:已报废" json:"status"`
	UserID        uint   `gorm:"index;comment:当前持有人ID" json:"user_id"`
	HoldingStatus string `gorm:"-" json:"holding_status"` // 瞬态字段，不存库
	// UserName 是运行时填充的持有人姓名，不存库（-）
	UserName string `gorm:"-" json:"user_name"`
}
