package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"it-asset-system/core"
	"it-asset-system/models"
	"it-asset-system/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters")
		return
	}

	var user models.SysUser
	if err := core.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		utils.Error(c, 401, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		utils.Error(c, 500, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"role_id":   user.RoleID,
		},
	})
}

func GetUsers(c *gin.Context) {
	var users []models.SysUser
	if err := core.DB.Find(&users).Error; err != nil {
		utils.Error(c, 500, "Failed to get users")
		return
	}
	utils.Success(c, users)
}

func CreateUser(c *gin.Context) {
	var user models.SysUser
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Error(c, 400, "Invalid parameters")
		return
	}

	var count int64
	core.DB.Model(&models.SysUser{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	if err := core.DB.Create(&user).Error; err != nil {
		utils.Error(c, 500, "Failed to create user")
		return
	}

	utils.Success(c, user)
}

func UpdateUser(c *gin.Context) {
	var req struct {
		ID       uint   `json:"id" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password"`
		RealName string `json:"real_name"`
		RoleID   int    `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters: "+err.Error())
		return
	}

	var user models.SysUser
	if err := core.DB.First(&user, req.ID).Error; err != nil {
		utils.Error(c, 404, "User not found")
		return
	}

	var count int64
	core.DB.Model(&models.SysUser{}).Where("username = ? AND id != ?", req.Username, req.ID).Count(&count)
	if count > 0 {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	user.Username = req.Username
	user.RealName = req.RealName
	user.RoleID = req.RoleID
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := core.DB.Save(&user).Error; err != nil {
		utils.Error(c, 500, "Failed to update user")
		return
	}

	utils.Success(c, user)
}

// DeleteUser 删除用户
// DELETE /admin/users
func DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Error(c, 400, "Missing user ID")
		return
	}

	userIDVal, _ := c.Get("userID")
	if fmt.Sprintf("%v", userIDVal) == id {
		utils.Error(c, 400, "不能删除当前登录用户")
		return
	}
	if id == "1" {
		utils.Error(c, 400, "不能删除超级管理员")
		return
	}

	var user models.SysUser
	if err := core.DB.First(&user, id).Error; err != nil {
		utils.Error(c, 404, "User not found")
		return
	}

	if err := core.DB.Delete(&user).Error; err != nil {
		utils.Error(c, 500, "Failed to delete user")
		return
	}

	utils.Success(c, "Deleted successfully")
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func UpdateMyProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.Error(c, 401, "User ID not found in token")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		utils.Error(c, 401, "Invalid User ID type in token")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters")
		return
	}

	var user models.SysUser
	if err := core.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, 404, "User not found")
		return
	}

	// Check if username is already taken by another user
	var count int64
	core.DB.Model(&models.SysUser{}).Where("username = ? AND id != ?", req.Username, userID).Count(&count)
	if count > 0 {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	user.Username = req.Username
	user.RealName = req.RealName

	if err := core.DB.Save(&user).Error; err != nil {
		utils.Error(c, 500, "Failed to update profile")
		return
	}

	utils.Success(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"role_id":   user.RoleID,
	})
}

func UpdateMyPassword(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.Error(c, 401, "User ID not found in token")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		utils.Error(c, 401, "Invalid User ID type in token")
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid parameters")
		return
	}

	var user models.SysUser
	if err := core.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, 404, "User not found")
		return
	}

	// Check if old password is correct
	if user.Password != req.OldPassword {
		utils.Error(c, 400, "原密码错误")
		return
	}

	user.Password = req.NewPassword
	if err := core.DB.Save(&user).Error; err != nil {
		utils.Error(c, 500, "Failed to update password")
		return
	}

	utils.Success(c, "密码修改成功")
}

