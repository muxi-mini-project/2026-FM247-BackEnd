package main

import (
	"2026-FM247-BackEnd/config"
	handler "2026-FM247-BackEnd/handlers"
	repository "2026-FM247-BackEnd/repositories"
	"2026-FM247-BackEnd/router"
	"2026-FM247-BackEnd/service"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 优先从 go.env 加载环境变量（开发环境）
	_ = godotenv.Overload("go.env")

	fmt.Println("正在初始化配置")
	config.LoadConfig()

	db, err := config.ConnectDatabase()
	if err != nil {
		fmt.Printf("无法连接到mysql数据库: %v\n", err)
		return
	}

	redisClient, err := config.ConnectRedis()
	if err != nil {
		fmt.Printf("无法连接到Redis: %v\n", err)
		return
	}

	fmt.Println("数据库连接成功")

	//dao层初始化
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenBlacklistRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	studyDataRepo := repository.NewStudyDataRepository(db, redisClient)

	//service层初始化
	userService := service.NewUserService(userRepo, tokenRepo)
	tokenService := service.NewTokenBlacklistService(tokenRepo)
	todoService := service.NewTodoService(todoRepo)
	studyDataService := service.NewStudyDataService(studyDataRepo)

	//handler层初始化
	authhandler := handler.NewAuthHandler(tokenService, userService)
	todohandler := handler.NewTodoHandler(todoService)
	studydatahandler := handler.NewStudyDataHandler(studyDataService)

	// 启动服务器
	r := gin.Default()

	// 打印请求头的中间件，需放在其他中间件和路由注册前
	r.Use(func(c *gin.Context) {
		fmt.Println("请求头：")
		for k, v := range c.Request.Header {
			fmt.Printf("%s: %v\n", k, v)
		}
		c.Next()
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // 允许的请求源
			"http://localhost:5173", // 前端vite的默认启动地址
			"http://localhost:3000", // 前端自己定义的启动地址
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                             // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"}, // 允许的请求头
		AllowCredentials: true,
	}))

	router.RegisterRoutes(r, authhandler, todohandler, studydatahandler)
	port := ":" + config.AppConfig.ServerPort
	fmt.Printf("服务器正在运行，监听端口 %s\n", port)
	if err := r.Run(port); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}

	r.Use(func(c *gin.Context) {
		fmt.Println("请求头：")
		for k, v := range c.Request.Header {
			fmt.Printf("%s: %v\n", k, v)
		}
		c.Next()
	})

	fmt.Println("正在与数据库断开连接")
	err = config.CloseDatabase(db)
	if err != nil {
		fmt.Printf("关闭数据库连接时出错: %v\n", err)
		return
	}
	fmt.Println("数据库连接已关闭")
}
