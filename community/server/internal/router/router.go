package router

import (
	"context"

	"community-server/internal/ai"
	"community-server/internal/config"
	"community-server/internal/handler"
	"community-server/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RouterParams struct {
	fx.In

	Config          *config.Config
	Logger          *zap.Logger
	UserHandler     *handler.UserHandler
	PostHandler     *handler.PostHandler
	CommentHandler  *handler.CommentHandler
	TagHandler      *handler.TagHandler
	AIHandler       *handler.AIHandler
	SearchHandler   *handler.SearchHandler
	AISearchHandler *handler.AISearchHandler
	AdminHandler    *handler.AdminHandler
	FollowHandler   *handler.FollowHandler
	MessageHandler  *handler.MessageHandler
	AIEngine        ai.Engine
}

func NewRouter(params RouterParams) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.RateLimitMiddleware(middleware.NewRateLimiter(), 100, 10))

	api := r.Group("/api/v1")
	{
		api.POST("/users/register", params.UserHandler.Register)
		api.POST("/users/login", params.UserHandler.Login)
		api.GET("/posts", params.PostHandler.GetPostList)
		api.GET("/posts/:id", params.PostHandler.GetPost)
		api.GET("/users/:user_id/posts", params.PostHandler.GetUserPosts)
		api.GET("/tags", params.TagHandler.GetTagList)
		api.GET("/tags/:id/posts", params.TagHandler.GetPostsByTag)

		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/users/profile", params.UserHandler.GetProfile)
			protected.PUT("/users/profile", params.UserHandler.UpdateProfile)

			protected.POST("/posts", params.PostHandler.CreatePost)
			protected.PUT("/posts/:id", params.PostHandler.UpdatePost)
			protected.DELETE("/posts/:id", params.PostHandler.DeletePost)
			protected.POST("/posts/:id/like", params.PostHandler.LikePost)
			protected.DELETE("/posts/:id/like", params.PostHandler.UnlikePost)
			protected.POST("/posts/:id/favorite", params.PostHandler.FavoritePost)
			protected.DELETE("/posts/:id/favorite", params.PostHandler.UnfavoritePost)

			protected.POST("/comments", params.CommentHandler.CreateComment)
			protected.PUT("/comments/:id", params.CommentHandler.UpdateComment)
			protected.DELETE("/comments/:id", params.CommentHandler.DeleteComment)
			protected.POST("/comments/:id/like", params.CommentHandler.LikeComment)
			protected.DELETE("/comments/:id/like", params.CommentHandler.UnlikeComment)

			protected.POST("/tags", params.TagHandler.CreateTag)
			protected.PUT("/tags/:id", params.TagHandler.UpdateTag)
			protected.DELETE("/tags/:id", params.TagHandler.DeleteTag)
			protected.POST("/posts/:id/tags", params.TagHandler.AddPostTags)
			protected.DELETE("/posts/:id/tags/:tag_id", params.TagHandler.RemovePostTag)
			protected.GET("/posts/:id/tags", params.TagHandler.GetPostTags)

			protected.POST("/follows", params.FollowHandler.FollowUser)
			protected.DELETE("/follows/:id", params.FollowHandler.UnfollowUser)
			protected.GET("/follows/followers", params.FollowHandler.GetFollowers)
			protected.GET("/follows/following", params.FollowHandler.GetFollowing)
			protected.GET("/follows/:id/status", params.FollowHandler.IsFollowing)
			protected.GET("/follows/:id/counts", params.FollowHandler.GetFollowCounts)

			protected.POST("/messages", params.MessageHandler.SendMessage)
			protected.GET("/messages", params.MessageHandler.GetMessageList)
			protected.GET("/messages/conversations", params.MessageHandler.GetConversationList)
			protected.GET("/messages/unread", params.MessageHandler.GetUnreadCount)
			protected.PUT("/messages/:id/read", params.MessageHandler.MarkAsRead)
		}

		api.GET("/comments", params.CommentHandler.GetCommentList)

		if params.Config.AI.ApiKey != "" && params.Config.AI.Url != "" {
			api.POST("/ai/chat", params.AIHandler.Chat)
			api.POST("/ai/chat/stream", params.AIHandler.ChatSSE)
			api.POST("/ai/search", params.AISearchHandler.AISearch)
			api.POST("/ai/search/stream", params.AISearchHandler.AISearchStream)
		}

		api.GET("/search", params.SearchHandler.SearchPosts)

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth())
		admin.Use(middleware.AdminAuth())
		{
			admin.GET("/users", params.AdminHandler.GetUserList)
			admin.DELETE("/users/:id", params.AdminHandler.DeleteUser)
			admin.PUT("/users/:id/admin_type", params.AdminHandler.UpdateUserAdminType)
			admin.PUT("/users/:id/status", params.AdminHandler.UpdateUserStatus)

			admin.GET("/posts", params.AdminHandler.GetPostList)
			admin.DELETE("/posts/:id", params.AdminHandler.DeletePost)
		}
	}

	return r
}

func RegisterRoutes(lc fx.Lifecycle, r *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := ":" + cfg.Server.Prot
			if addr == ":" {
				addr = ":8080"
			}
			logger.Info("Starting server", zap.String("addr", addr))
			go func() {
				if err := r.Run(addr); err != nil {
					logger.Error("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			return nil
		},
	})
}
