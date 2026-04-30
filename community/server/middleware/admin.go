package middleware

import (
	"community-server/DB/mysql"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		var user mysql.User
		if err := mysql.DB.Where("id = ?", userID).First(&user).Error; err != nil {
			response.ErrorWithMsg(c, response.CodeServerBusy, "获取用户信息失败")
			c.Abort()
			return
		}

		if user.AdminType != 1 {
			response.ErrorWithMsg(c, response.CodeForbidden, "需要管理员权限")
			c.Abort()
			return
		}

		c.Set("admin_type", user.AdminType)
		c.Next()
	}
}
