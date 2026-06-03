package main

import (
	"fmt"
	"log"

	"it-asset-system/core"
	"it-asset-system/models"
	"it-asset-system/routers"
)

func main() {
	// Initialize configurations and core components
	core.InitConfig()
	core.InitLogger()
	core.InitDB()
	core.InitRedis()
	core.InitCasbin()

	// Auto migrate database models
	err := core.DB.AutoMigrate(
		&models.SysUser{},
		&models.SysCategory{},
		&models.SysAsset{},
		&models.SysApplyLog{},
		&models.SysRepairLog{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Seed database data and Casbin policies
	core.SeedDatabase()

	// Setup router
	r := routers.SetupRouter()

	// Run server
	host := core.AppConfig.Server.Host
	port := core.AppConfig.Server.Port
	if host == "" {
		host = "127.0.0.1"
	}
	core.Logger.Info(fmt.Sprintf("Server is running on http://%s:%d", host, port))
	if err := r.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
