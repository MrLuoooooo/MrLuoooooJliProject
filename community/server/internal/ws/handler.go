package ws

import (
	"net/http"
	"strings"

	"community-server/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler 返回 WebSocket 连接端点，通过 URL 参数中的 token 鉴权
func Handler(m *Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			// 也支持从 Authorization header 拿
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				token = auth[7:]
			}
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 token"})
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 无效或已过期"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			zap.S().Warn("ws升级失败", "error", err)
			return
		}

		zap.S().Info("ws连接建立", "userID", claims.UserID)
		m.Add(claims.UserID, conn)
	}
}
