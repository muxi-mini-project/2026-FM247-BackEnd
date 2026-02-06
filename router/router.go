package router

import (
	handler "2026-FM247-BackEnd/handlers"
	middleware "2026-FM247-BackEnd/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	authhandler *handler.AuthHandler,
) {
	publicGroup := r.Group("/api")
	{
		// 用户相关
		publicGroup.POST("/auth/register", authhandler.RegisterUserHandler)
		publicGroup.POST("/auth/login", authhandler.LoginHandler)
	}

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware(authhandler.Tokenservice))
	{
		// 用户相关
		authGroup.POST("/auth/logout", authhandler.LogoutHandler)
		authGroup.POST("/auth/cancel", authhandler.CancelHandler)
		authGroup.POST("/user/update_info", authhandler.UpdateUserInfoHandler)
		authGroup.POST("/user/update_email", authhandler.UpdateEmailHandler)
		authGroup.POST("/user/update_password", authhandler.UpdatePasswordHandler)
		authGroup.GET("/user/info", authhandler.GetUserInfoHandler)
	}
}
