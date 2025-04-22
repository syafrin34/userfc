package main

import (
	"fmt"
	"userfc/cmd/user/handler"
	"userfc/cmd/user/repository"
	"userfc/cmd/user/resource"
	"userfc/cmd/user/service"
	"userfc/cmd/user/usecase"
	"userfc/config"
	"userfc/infrastructure/logger"
	"userfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)
	logger.SetupLogger()

	userRepository := repository.NewUserRepository(db, redis)
	userService := service.NewUserService(*userRepository)
	userUseCase := usecase.NewUserUseCase(*userService, cfg.Secret.JWTSecret)
	userHandler := handler.NewUserHandler(*userUseCase)

	port := cfg.App.Port // baca ke config yang sudah kita load diawal

	router := gin.Default()
	routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)
	router.Run(":" + port)
	fmt.Println("sekarang server berjalan di port : ", port)
	logger.Logger.Printf("Server running on port : %s", port)

	// router := gin.Default()
	// router.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "ping pong",
	// 	})
	// })
	// router.Run()

}
