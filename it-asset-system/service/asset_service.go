package service

import (
	"context"
	"encoding/json"
	"time"

	"it-asset-system/core"
	"it-asset-system/models"
)

const AvailableAssetsCacheKey = "asset:available:list"

func GetAvailableAssets() ([]models.SysAsset, error) {
	ctx := context.Background()

	// 1. Try to get from Redis
	val, err := core.RedisClient.Get(ctx, AvailableAssetsCacheKey).Result()
	if err == nil && val != "" {
		var assets []models.SysAsset
		if err := json.Unmarshal([]byte(val), &assets); err == nil {
			return assets, nil
		}
	}

	// 2. Not in cache, query DB
	var assets []models.SysAsset
	if err := core.DB.Where("status = ?", 0).Find(&assets).Error; err != nil {
		return nil, err
	}

	// 3. Set cache
	if bytes, err := json.Marshal(assets); err == nil {
		core.RedisClient.Set(ctx, AvailableAssetsCacheKey, string(bytes), time.Hour)
	}

	return assets, nil
}

func ClearAvailableAssetsCache() {
	core.RedisClient.Del(context.Background(), AvailableAssetsCacheKey)
}
