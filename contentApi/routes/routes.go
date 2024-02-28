package routes

import (
	"contentApi/controllers"
	"contentApi/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	api := r.Group("/api")

	userGroup := api.Group("/users")
	userGroup.Use(middlewares.JwtAuthMiddleware())
	{
		userGroup.GET("/", controllers.GetUsers)
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.GET("/me", controllers.CurrentUser)
		userGroup.POST("/", controllers.CreateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}

	err := r.Run()

	if err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}
