package router

import (
	handler "2026-FM247-BackEnd/handlers"
	middleware "2026-FM247-BackEnd/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	authhandler *handler.AuthHandler,
	userhandler *handler.UserHandler,
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
		authGroup.POST("/user/update_info", userhandler.UpdateUserInfoHandler)
		authGroup.POST("/user/update_password", userhandler.UpdatePasswordHandler)
		authGroup.GET("/user/info", userhandler.GetUserInfoHandler)
	}
}
