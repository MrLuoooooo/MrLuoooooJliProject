package di

import (
	"community-server/internal/ai"
	"community-server/internal/config"
	"community-server/internal/db/mysql"
	"community-server/internal/handler"
	"community-server/internal/im"
	"community-server/internal/logger"
	"community-server/internal/repository"
	"community-server/internal/router"
	"community-server/internal/service"
	"community-server/internal/ws"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	// ============================================
	// 配置 + 日志
	// ============================================
	fx.Provide(config.New),
	fx.Invoke(logger.InitLogger),
	fx.Provide(func() *zap.Logger { return zap.L() }),

	// ============================================
	// 数据库
	// ============================================
	fx.Provide(mysql.NewMySQL),
	fx.Invoke(mysql.AutoMigrate),

	// ============================================
	// 数据访问层（GORM 实现）
	// ============================================
	fx.Provide(
		repository.NewUserRepo,
		repository.NewPostRepo,
		repository.NewPostLikeRepo,
		repository.NewPostFavoriteRepo,
		repository.NewCommentRepo,
		repository.NewCommentLikeRepo,
		repository.NewTagRepo,
		repository.NewPostTagRepo,
		repository.NewCategoryRepo,
		repository.NewFollowRepo,
		repository.NewMessageRepo,
		repository.NewNotificationRepo,
		repository.NewPasswordResetRepo,
	),

	// ============================================
	// AI + IM 引擎 + WebSocket
	// ============================================
	fx.Provide(ai.NewEngine),
	fx.Provide(fx.Annotate(im.NewIMClient, fx.As(new(im.IMClient)))),
	fx.Provide(ws.NewManager),

	// ============================================
	// 业务服务（同时注册具体类型和接口）
	// ============================================
	fx.Provide(service.NewNotificationService),
	fx.Provide(fx.Annotate(service.NewNotificationService, fx.As(new(handler.NotificationService)))),
	fx.Provide(service.NewUserService),
	fx.Provide(fx.Annotate(service.NewUserService, fx.As(new(handler.UserService)))),
	fx.Provide(service.NewPostService),
	fx.Provide(fx.Annotate(service.NewPostService, fx.As(new(handler.PostService)))),
	fx.Provide(service.NewCommentService),
	fx.Provide(fx.Annotate(service.NewCommentService, fx.As(new(handler.CommentService)))),
	fx.Provide(service.NewTagService),
	fx.Provide(fx.Annotate(service.NewTagService, fx.As(new(handler.TagService)))),
	fx.Provide(service.NewCategoryService),
	fx.Provide(fx.Annotate(service.NewCategoryService, fx.As(new(handler.CategoryService)))),
	fx.Provide(service.NewUploadService),
	fx.Provide(fx.Annotate(service.NewUploadService, fx.As(new(handler.UploadService)))),
	fx.Provide(service.NewSearchService),
	fx.Provide(fx.Annotate(service.NewSearchService, fx.As(new(handler.SearchService)))),
	fx.Provide(service.NewAdminService),
	fx.Provide(fx.Annotate(service.NewAdminService, fx.As(new(handler.AdminService)))),
	fx.Provide(service.NewFollowService),
	fx.Provide(fx.Annotate(service.NewFollowService, fx.As(new(handler.FollowService)))),
	fx.Provide(service.NewMessageService),
	fx.Provide(fx.Annotate(service.NewMessageService, fx.As(new(handler.MessageService)))),

	// ============================================
	// HTTP 处理层
	// ============================================
	fx.Provide(handler.NewUserHandler),
	fx.Provide(handler.NewPostHandler),
	fx.Provide(handler.NewCommentHandler),
	fx.Provide(handler.NewTagHandler),
	fx.Provide(handler.NewCategoryHandler),
	fx.Provide(handler.NewNotificationHandler),
	fx.Provide(handler.NewUploadHandler),
	fx.Provide(handler.NewStatusHandler),
	fx.Provide(handler.NewBotWebhookHandler),
	fx.Provide(handler.NewAIHandler),
	fx.Provide(handler.NewSearchHandler),
	fx.Provide(handler.NewAISearchHandler),
	fx.Provide(handler.NewAdminHandler),
	fx.Provide(handler.NewFollowHandler),
	fx.Provide(handler.NewMessageHandler),

	// ============================================
	// 路由
	// ============================================
	fx.Provide(router.NewRouter),
	fx.Invoke(router.RegisterRoutes),
)
