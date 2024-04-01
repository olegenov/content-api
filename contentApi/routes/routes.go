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
		userGroup.GET("/", middlewares.AdminOnly(), controllers.GetUsers)
		userGroup.GET("/:id", middlewares.AdminOnly(), controllers.GetUser)
		userGroup.GET("/me", controllers.CurrentUser)
		userGroup.POST("/", middlewares.AdminOnly(), controllers.CreateUser)
		userGroup.DELETE("/:id", middlewares.AdminOnly(), controllers.DeleteUser)
	}

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}

	projectsGroup := api.Group("/projects")
	projectsGroup.Use(middlewares.JwtAuthMiddleware())
	{
		projectsGroup.GET("/", middlewares.AdminOnly(), controllers.GetProjects)
		projectsGroup.POST("/", controllers.CreateProject)
		projectsGroup.GET("/:id", controllers.GetProject)
		projectsGroup.GET("/my", controllers.GetMyProjects)
	}

	postsGroup := api.Group("/projects/:id/post")
	postsGroup.Use(middlewares.JwtAuthMiddleware())
	{
		postsGroup.GET("/:post-id", controllers.GetPost)
	}

	teamsGroup := api.Group("/teams")
	teamsGroup.Use(middlewares.JwtAuthMiddleware())
	{
		teamsGroup.GET("/", middlewares.AdminOnly(), controllers.GetTeams)
		teamsGroup.POST("/", controllers.CreateTeam)
		teamsGroup.GET("/:id", controllers.GetTeam)
		teamsGroup.GET("/my", controllers.GetMyTeams)
	}

	err := r.Run()

	if err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}
