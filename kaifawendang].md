一、 选题背景与意义
在小型企业日常办公中，IT设备与固定资产频繁流转，传统表格管理易导致账实不符、审批混乱与资产流失。本系统旨在构建高效的资产后台，实现设备入库、多角色审批及报修闭环；通过引入 Redis 缓存与 GORM 事务，保障高频查询性能与核心资产数据一致性，具备极强的企业实际应用价值。

二、 系统核心功能模块
系统严格划分三个用户角色（普通员工、资产管理员、系统管理员），实现完整的资产流转闭环。

资产与设备管理 (核心 CRUD)：

支持资产入库（SN 码、分类、状态录入）。

资产状态流转机（闲置 -> 使用中 -> 维修中 -> 报废）。

领用与归还审批流：

员工发起领用/归还申请，填写用途与时间。

管理员进行审批（同意/驳回），驱动资产状态变更。

极简报修流（高分事务专属模块）：

员工针对名下设备提交“一句话报修”。

管理员“确认送修”，后端利用事务同时变更报修单与设备状态。

RBAC 权限与系统监控：

基于 JWT + Casbin 实现多角色接口级别鉴权。

基于 Zap 实现操作审计日志与全局异常拦截。

三、 技术栈选型与环境配置
基础环境：Go 1.20+、MySQL 8.0、Redis 7.0

Web 框架：Gin (路由与控制器)

配置管理：Viper (解析 config.yaml)

持久层：GORM (开启 MySQL 连接池)

权限安全：JWT-Go (鉴权) + Casbin (RBAC 角色隔离)

数据校验：Validator (请求入参强校验)

日志系统：Zap + Lumberjack (按天/按大小分级切割日志)

四、 项目目录结构 (分层架构)
严格遵循高内聚低耦合规范，满足课程对工程化的考核标准：

Plaintext


it-asset-system/
├── conf/             # Viper 配置目录 (config.yaml)
├── core/             # 核心组件初始化 (MySQL, Redis, ZapLogger, CasbinEnforcer)
├── middleware/       # 自定义中间件 (JWTAuth, CasbinRBAC, ErrorHandle, ZapLog)
├── models/           # GORM 结构体定义 (包含 Validator 校验标签)
├── repository/       # 数据访问层 (纯数据库 CRUD 与 GORM 事务)
├── service/          # 业务逻辑层 (处理 Redis 缓存策略、数据组装)
├── controller/       # 接口控制层 (接收入参, 校验, 返回标准 JSON)
├── routers/          # Gin 路由注册与权限分组
├── utils/            # 工具类 (密码加密、Token生成、统一返回结构)
└── main.go           # 项目启动入口
五、 核心数据库表设计 (GORM 模型)
通过 GORM AutoMigrate 自动建表，包含完整的关联与状态字段：

用户表 (sys_users)

字段：id, username (唯一索引), password, real_name, role_id (1:员工, 2:IT管理员, 3:系统管理员)。

资产分类表 (sys_categories)

字段：id, category_name (如: 笔记本、显示器)。

资产明细表 (sys_assets)

字段：id, asset_no (唯一编号), name, category_id, status (0:闲置, 1:使用中, 2:维修中, 3:报废), user_id (当前持有人)。

优化：建立 category_id 与 status 的联合索引。

领用记录表 (sys_apply_logs)

字段：id, user_id, asset_id, type (1:领用, 2:归还), status (0:待审批, 1:已通过, 2:已驳回)。

极简报修单表 (sys_repair_logs)

字段：id, asset_id, user_id, reason (报修原因), status (0:待处理, 1:已送修)。

六、 核心业务逻辑与高分技术落地
1. Redis 缓存策略设计
首页“可用资产列表”读多写少，是完美的缓存场景。

执行逻辑：请求进入 service 层，先查询 Redis Key: asset:available:list。若命中直接返回；若未命中查询 MySQL，将结果写入 Redis 并设置 1 小时 TTL。

一致性保障：管理员入库新资产或审批通过领用后，触发 Redis.Del 删除对应缓存。

2. GORM 事务控制 (极简报修核心)
这是答辩时展示数据库安全性的核心代码。在 repository/repair_repo.go 中实现：

Go


package repository

import (
	"errors"
	"gorm.io/gorm"
	"it-asset-system/core"
	"it-asset-system/models"
)

// ConfirmRepair 确认送修核心事务
func ConfirmRepair(repairID uint, assetID uint) error {
	tx := core.DB.Begin() // 开启事务
	if tx.Error != nil {
		return tx.Error
	}

	// 捕获恐慌自动回滚
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

	// 步骤 2: 资产状态 -> "2: 维修中"
	if err := tx.Model(&models.SysAsset{}).Where("id = ?", assetID).Update("status", 2).Error; err != nil {
		tx.Rollback()
		return errors.New("更新资产状态失败")
	}

	return tx.Commit().Error // 提交事务
}
3. Casbin 多角色路由控制与日志切割
权限分配：在 middleware.CasbinRBAC() 中拦截请求。员工无法访问 /api/v1/admin/* 路径下的任何路由。

日志规范：接入 Lumberjack，配置 Zap 将常规日志与错误日志分离（info.log 与 error.log），按 50MB 自动切割，防止服务器磁盘撑爆。

七、 接口路由规划 (RESTful API)
开发人员需在 routers/router.go 中严格按以下分组注册路由：

Go


// 1. 公共路由 (无需 Token)
publicGroup := r.Group("/api/v1")
{
    publicGroup.POST("/login", controller.UserLogin) 
}

// 2. 员工权限分组 (需 JWT 认证 + Casbin 校验 role_1 权限)
employeeGroup := r.Group("/api/v1/employee").Use(middleware.JWTAuth(), middleware.CasbinRBAC())
{
    employeeGroup.GET("/assets/available", controller.GetAvailableAssets) // 走 Redis 缓存
    employeeGroup.POST("/apply", controller.SubmitApply)                  // 发起领用
    employeeGroup.POST("/repair", controller.SubmitRepair)                // 提交极简报修
}

// 3. IT管理员权限分组 (需 JWT 认证 + Casbin 校验 role_2 权限)
adminGroup := r.Group("/api/v1/admin").Use(middleware.JWTAuth(), middleware.CasbinRBAC())
{
    adminGroup.POST("/asset", controller.CreateAsset)                     // 资产入库
    adminGroup.PUT("/apply/audit", controller.AuditApply)                 // 领用审批
    adminGroup.PUT("/repair/confirm", controller.ConfirmRepair)           // 确认送修 (触发核心事务)
}