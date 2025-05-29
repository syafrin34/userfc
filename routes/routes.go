package routes

import (
	"userfc/cmd/user/handler"
	middleware "userfc/middlaware"

	"github.com/gin-gonic/gin"
	//"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRoutes(router *gin.Engine, userHandler handler.UserHandler, jwtSecret string) {
	// public api
	router.Use(middleware.RequestLogger())
	//router.Use(otelgin.middlaware("userfc"))
	router.Use(otelgin.Middleware("userfc"))
	router.GET("/ping", userHandler.Ping)
	router.POST("/v1/register", userHandler.Register)
	router.POST("v1/login", userHandler.Login)

	// Private api
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	private := router.Group("/api")
	private.Use(authMiddleware)
	private.GET("/v1/user_info", userHandler.GetUserInfo)
}
