package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"it-asset-system/core"
	"it-asset-system/utils"
)

func CasbinRBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			utils.Error(c, 401, "Role not found in token")
			c.Abort()
			return
		}

		sub := fmt.Sprintf("role_%d", roleID)
		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := core.Enforcer.Enforce(sub, obj, act)
		if err != nil {
			utils.Error(c, 500, "Error occurred when authorizing")
			c.Abort()
			return
		}

		if !ok {
			utils.Error(c, 403, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
