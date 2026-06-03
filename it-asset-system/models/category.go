package models

import (
	"gorm.io/gorm"
)

type SysCategory struct {
	gorm.Model
	CategoryName string `gorm:"type:varchar(100);not null" json:"category_name" binding:"required"`
}
