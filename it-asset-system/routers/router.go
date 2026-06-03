package routers

import (
	"github.com/gin-gonic/gin"
	"it-asset-system/controller"
	"it-asset-system/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.ZapLog())
	r.Use(middleware.ErrorHandle())
	r.Use(gin.Recovery())

	// 1. 公共路由 (无需 Token)
	publicGroup := r.Group("/api/v1")
	{
		publicGroup.POST("/login", controller.UserLogin)
	}

	// 2. 员工权限分组 (需 JWT 认证 + Casbin 校验)
	employeeGroup := r.Group("/api/v1/employee").Use(middleware.JWTAuth(), middleware.CasbinRBAC())
	{
		employeeGroup.GET("/assets/available", controller.GetAvailableAssets) // 走 Redis 缓存
		employeeGroup.POST("/apply", controller.SubmitApply)                  // 发起领用
		employeeGroup.POST("/repair", controller.SubmitRepair)                // 提交极简报修
		employeeGroup.GET("/apply/my", controller.GetMyApplyList)             // 获取我本人的申请记录
		employeeGroup.GET("/assets/held", controller.GetHeldAssets)           // 获取当前持有的资产
		employeeGroup.PUT("/profile", controller.UpdateMyProfile)
		employeeGroup.PUT("/profile/password", controller.UpdateMyPassword)
	}

	// 3. IT管理员权限分组 (需 JWT 认证 + Casbin 校验)
	adminGroup := r.Group("/api/v1/admin").Use(middleware.JWTAuth(), middleware.CasbinRBAC())
	{
		adminGroup.POST("/asset", controller.CreateAsset)           // 资产入库
		adminGroup.PUT("/asset", controller.UpdateAsset)            // 修改资产 (包括数量、名称等)
		adminGroup.PUT("/apply/audit", controller.AuditApply)       // 领用审批
		adminGroup.PUT("/repair/confirm", controller.ConfirmRepair) // 确认送修 (触发核心事务)
		adminGroup.PUT("/repair/complete", controller.CompleteRepair) // 完成维修
		adminGroup.GET("/apply/list", controller.GetApplyList)      // 获取领用申请列表
		adminGroup.GET("/repair/list", controller.GetRepairList)    // 获取报修申请列表
		adminGroup.GET("/categories", controller.GetCategories)      // 获取分类列表
		adminGroup.GET("/users", controller.GetUsers)               // 获取用户列表
		adminGroup.POST("/users", controller.CreateUser)             // 添加用户
		adminGroup.PUT("/users", controller.UpdateUser)             // 修改用户 (包括角色、密码等)
		adminGroup.DELETE("/asset", controller.DeleteAsset)         // 删除资产单件
		adminGroup.DELETE("/users", controller.DeleteUser)          // 删除用户
	}

	return r
}
