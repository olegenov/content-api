package controllers

import (
	"contentApi/controllers/responseControllers"
	"contentApi/dto/requests"
	"contentApi/dto/responses"
	"contentApi/models"
	"contentApi/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.Preload("Creator").Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	var projectResponse []responses.ProjectResponse

	for _, project := range projects {
		projectResponse = append(projectResponse, responseControllers.GetProjectResponse(project))
	}

	c.JSON(http.StatusOK, projectResponse)
}

func GetProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	if err := models.DB.Preload("Creator").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	userID, _ := token.ExtractTokenID(c)

	if userID != project.CreatorID && project.Creator.Role != "ADMIN" {
		c.JSON(http.StatusNotFound, gin.H{"error": "It is not your project"})
		return
	}

	projectResponse := responseControllers.GetProjectResponse(project)
	c.JSON(http.StatusOK, projectResponse)
}

func CreateProject(c *gin.Context) {
	var projectRequest requests.ProjectRequest

	if err := c.ShouldBindJSON(&projectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := token.ExtractTokenID(c)
	project := models.Project{
		Name:      projectRequest.Name,
		CreatorID: userID,
	}

	result := models.DB.Create(&project)

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": project.ID, "message": "Project created successfully"})
}

func GetMyProjects(c *gin.Context) {
	var projects []models.Project
	userID, _ := token.ExtractTokenID(c)

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.Preload("Creator").Where("creator_id = ?", userID).Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	if len(projects) == 0 {
		c.JSON(http.StatusOK, gin.H{"projects": []responses.ProjectResponse{}})
		return
	}

	var projectResponse []responses.ProjectResponse

	for _, project := range projects {
		projectResponse = append(projectResponse, responseControllers.GetProjectResponse(project))
	}

	c.JSON(http.StatusOK, gin.H{"projects": projectResponse})
}
