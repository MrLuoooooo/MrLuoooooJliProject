package di

import (
	"context"

	"community-server/DB/mysql"
	"community-server/internal/ai"
	"community-server/internal/config"
	"community-server/internal/handler"
	"community-server/internal/logger"
	"community-server/internal/router"
	"community-server/internal/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	ConfigModule,
	LoggerModule,
	DatabaseModule,
	AIModule,
	ServiceModule,
	HandlerModule,
	RouterModule,
)

var ConfigModule = fx.Options(
	fx.Provide(config.New),
)

var LoggerModule = fx.Options(
	fx.Provide(logger.NewLogger),
)

var DatabaseModule = fx.Options(
	fx.Provide(mysql.NewMySQL),
	fx.Invoke(mysql.AutoMigrate),
)

var AIModule = fx.Options(
	fx.Provide(ai.NewEngine),
)

var ServiceModule = fx.Options(
	fx.Provide(
		service.NewUserService,
		service.NewPostService,
		service.NewCommentService,
		service.NewTagService,
		service.NewSearchService,
		service.NewAdminService,
		service.NewFollowService,
		service.NewMessageService,
	),
)

var HandlerModule = fx.Options(
	fx.Provide(
		handler.NewUserHandler,
		handler.NewPostHandler,
		handler.NewCommentHandler,
		handler.NewTagHandler,
		handler.NewAIHandler,
		handler.NewSearchHandler,
		handler.NewAISearchHandler,
		handler.NewAdminHandler,
		handler.NewFollowHandler,
		handler.NewMessageHandler,
	),
)

var RouterModule = fx.Options(
	fx.Provide(router.NewRouter),
	fx.Invoke(router.RegisterRoutes),
)

func OnStartHook(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Application started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Application stopping")
			mysql.Close()
			return nil
		},
	})
}
