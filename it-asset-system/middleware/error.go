package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"it-asset-system/core"
	"it-asset-system/utils"
)

func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				core.Logger.Error(fmt.Sprintf("Panic recovered: %v", err))
				utils.Error(c, 500, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
