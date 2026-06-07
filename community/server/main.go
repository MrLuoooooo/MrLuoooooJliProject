// @title           Community Server API
// @version         1.0
// @description     社区论坛后端服务 — 帖子、评论、用户、通知、AI 搜索
// @host            localhost:1807
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

package main

import (
	"community-server/internal/di"
	"community-server/internal/logger"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		di.Module,
		fx.Invoke(logger.InitLogger),
	)

	app.Run()
}
