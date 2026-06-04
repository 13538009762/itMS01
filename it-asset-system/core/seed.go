package core

import (
	"fmt"
	"it-asset-system/models"
)

func SeedDatabase() {
	Logger.Info("Checking if database seeding is required...")

	// 1. Seed Categories if empty
	var categoryCount int64
	DB.Model(&models.SysCategory{}).Count(&categoryCount)
	if categoryCount == 0 {
		Logger.Info("Seeding default categories...")
		categories := []models.SysCategory{
			{CategoryName: "笔记本电脑"},
			{CategoryName: "台式电脑"},
			{CategoryName: "显示器"},
			{CategoryName: "服务器及网络设备"},
			{CategoryName: "办公外设"},
		}
		for _, cat := range categories {
			if err := DB.Create(&cat).Error; err != nil {
				Logger.Error(fmt.Sprintf("Failed to seed category %s: %v", cat.CategoryName, err))
			}
		}
	}

	// 2. Seed Users if empty
	var userCount int64
	DB.Model(&models.SysUser{}).Count(&userCount)
	if userCount == 0 {
		Logger.Info("Seeding default users...")
		users := []models.SysUser{
			{
				Username: "admin",
				Password: "123456",
				RealName: "系统管理员",
				RoleID:   3,
			},
			{
				Username: "it_admin",
				Password: "123456",
				RealName: "IT管理员",
				RoleID:   2,
			},
			{
				Username: "employee",
				Password: "123456",
				RealName: "普通员工",
				RoleID:   1,
			},
		}
		for _, user := range users {
			if err := DB.Create(&user).Error; err != nil {
				Logger.Error(fmt.Sprintf("Failed to seed user %s: %v", user.Username, err))
			}
		}
	}

	// 3. Seed Assets if empty — 每台物理设备独立一行，BaseNo 相同的设备属同一批次
	var assetCount int64
	DB.Model(&models.SysAsset{}).Count(&assetCount)
	if assetCount == 0 {
		Logger.Info("Seeding default individual asset units...")

		var laptopCat models.SysCategory
		DB.Where("category_name = ?", "笔记本电脑").First(&laptopCat)

		var monitorCat models.SysCategory
		DB.Where("category_name = ?", "显示器").First(&monitorCat)

		// 定义批次：baseNo, name, categoryID, quantity
		type batch struct {
			baseNo     string
			name       string
			categoryID uint
			quantity   int
		}
		batches := []batch{
			{"AST20260601001", "MacBook Pro 16 M3 Max", laptopCat.ID, 2},
			{"AST20260601002", "ThinkPad T14 Gen 4", laptopCat.ID, 1},
			{"AST20260601003", "Dell U2723QE 4K Monitor", monitorCat.ID, 3},
		}

		for _, b := range batches {
			for i := 1; i <= b.quantity; i++ {
				unit := models.SysAsset{
					BaseNo:     b.baseNo,
					AssetNo:    fmt.Sprintf("%s-%03d", b.baseNo, i),
					Name:       b.name,
					CategoryID: b.categoryID,
					Status:     0,
				}
				if err := DB.Create(&unit).Error; err != nil {
					Logger.Error(fmt.Sprintf("Failed to seed asset unit %s: %v", unit.AssetNo, err))
				}
			}
		}
	}

	// 4. Seed Casbin Policies
	Logger.Info("Seeding Casbin policies...")

	policies := [][]string{
		// 员工权限
		{"role_1", "/api/v1/employee/assets/available", "GET"},
		{"role_1", "/api/v1/employee/apply", "POST"},
		{"role_1", "/api/v1/employee/repair", "POST"},
		{"role_1", "/api/v1/employee/apply/my", "GET"},
		{"role_1", "/api/v1/employee/assets/held", "GET"},
		{"role_1", "/api/v1/employee/profile", "PUT"},
		{"role_1", "/api/v1/employee/profile/password", "PUT"},
		// IT管理员权限
		{"role_2", "/api/v1/admin/asset", "POST"},
		{"role_2", "/api/v1/admin/asset", "PUT"},
		{"role_2", "/api/v1/admin/apply/audit", "PUT"},
		{"role_2", "/api/v1/admin/repair/confirm", "PUT"},
		{"role_2", "/api/v1/admin/repair/complete", "PUT"},
		{"role_2", "/api/v1/admin/apply/list", "GET"},
		{"role_2", "/api/v1/admin/repair/list", "GET"},
		{"role_2", "/api/v1/admin/categories", "GET"},
		{"role_2", "/api/v1/admin/asset", "DELETE"},
		// 系统管理员专属权限 (不被IT管理员继承)
		{"role_3", "/api/v1/admin/users", "GET"},
		{"role_3", "/api/v1/admin/users", "POST"},
		{"role_3", "/api/v1/admin/users", "PUT"},
		{"role_3", "/api/v1/admin/users", "DELETE"},
	}

	// 清理旧的 role_2 用户管理权限 (以防数据库中残留)
	legacyUserPolicies := [][]string{
		{"role_2", "/api/v1/admin/users", "GET"},
		{"role_2", "/api/v1/admin/users", "POST"},
		{"role_2", "/api/v1/admin/users", "PUT"},
		{"role_2", "/api/v1/admin/users", "DELETE"},
	}
	for _, p := range legacyUserPolicies {
		has, _ := Enforcer.HasPolicy(p[0], p[1], p[2])
		if has {
			_, err := Enforcer.RemovePolicy(p[0], p[1], p[2])
			if err != nil {
				Logger.Error(fmt.Sprintf("Failed to remove legacy policy %v: %v", p, err))
			}
		}
	}

	for _, policy := range policies {
		has, err := Enforcer.HasPolicy(policy[0], policy[1], policy[2])
		if err != nil {
			Logger.Error(fmt.Sprintf("Failed to check Casbin policy: %v", err))
			continue
		}
		if !has {
			_, err = Enforcer.AddPolicy(policy[0], policy[1], policy[2])
			if err != nil {
				Logger.Error(fmt.Sprintf("Failed to add Casbin policy %v: %v", policy, err))
			}
		}
	}

	// 角色继承：role_3(系统管理员) > role_2(IT管理员) > role_1(员工)
	roleInheritance := [][]string{
		{"role_2", "role_1"},
		{"role_3", "role_2"},
	}

	for _, rule := range roleInheritance {
		has, err := Enforcer.HasGroupingPolicy(rule[0], rule[1])
		if err != nil {
			Logger.Error(fmt.Sprintf("Failed to check Casbin grouping policy: %v", err))
			continue
		}
		if !has {
			_, err = Enforcer.AddGroupingPolicy(rule[0], rule[1])
			if err != nil {
				Logger.Error(fmt.Sprintf("Failed to add Casbin grouping policy %v: %v", rule, err))
			}
		}
	}

	if err := Enforcer.SavePolicy(); err != nil {
		Logger.Error(fmt.Sprintf("Failed to save Casbin policy changes: %v", err))
	} else {
		Logger.Info("Casbin policies checked and synced successfully")
	}

	Logger.Info("Database seeding verification completed.")
}
