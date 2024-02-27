package routes

import (
	"contentApi/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	userGroup := r.Group("/user")
	{
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/:id", controllers.GetUser)
	}

	err := r.Run()

	if err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}
