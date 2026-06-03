package models

import (
	"gorm.io/gorm"
)

type SysUser struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" binding:"required"`
	Password string `gorm:"type:varchar(100);not null" json:"password" binding:"required"`
	RealName string `gorm:"type:varchar(50);not null" json:"real_name"`
	RoleID   int    `gorm:"type:int;not null;comment:1:员工, 2:IT管理员, 3:系统管理员" json:"role_id"`
}
