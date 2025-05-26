package main

import (
	"net"
	"userfc/cmd/user/handler"
	"userfc/cmd/user/repository"
	"userfc/cmd/user/resource"
	"userfc/cmd/user/service"
	"userfc/cmd/user/usecase"
	"userfc/config"
	grpcUser "userfc/grpc"
	"userfc/infrastructure/logger"
	"userfc/proto/userpb"
	"userfc/routes"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	go func() {
		port := cfg.App.Port // baca ke config yang sudah kita load diawal
		router := gin.Default()
		routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)
		logger.Logger.Printf("Server running on port : %s", port)
		//router.Run(":" + port)
		if err := router.Run(":" + port); err != nil {
			logger.Logger.Fatalf("Failed to start http server: %v", err.Error())
		}
	}()

	// init grpc server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Logger.Fatal("Failed init GRPC Server: ", err.Error())
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &grpcUser.GRPCServer{UserUsecase: *userUseCase})
	reflection.Register(grpcServer)
	logger.Logger.Printf("GRPC Server running on port : %s", ":50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Logger.Fatalf("Failed to serve GRPC:  %V", err)
	}

	// router := gin.Default()
	// router.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "ping pong",
	// 	})
	// })
	// router.Run()

}
