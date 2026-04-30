package middleware

import (
	"strings"

	"community-server/pkg/jwt"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithMsg(c, response.CodeInvalidParam, "请求头中Authorization为空")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithMsg(c, response.CodeInvalidParam, "请求头格式错误，应为'Bearer {token}'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				response.ErrorWithMsg(c, response.CodeTokenExpired, "令牌已过期")
			} else {
				response.ErrorWithMsg(c, response.CodeTokenInvalid, "令牌无效")
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.UserName)
		c.Next()
	}
}
