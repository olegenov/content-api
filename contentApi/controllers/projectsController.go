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
	result := models.DB.Preload("Creator").Find(&projects)
	models.DbMutex.Unlock()

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
	if err := models.DB.Preload("Creator").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	models.DbMutex.Unlock()

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

	var team models.Team

	models.DbMutex.Lock()
	teamResult := models.DB.First(&team, projectRequest.TeamID)
	models.DbMutex.Unlock()

	if teamResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team"})
		return
	}

	var isMember bool
	for _, user := range team.Users {
		if user.ID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User is not a member of the team"})
		return
	}

	project := models.Project{
		Name:      projectRequest.Name,
		CreatorID: userID,
		TeamID:    projectRequest.TeamID,
	}

	models.DbMutex.Lock()
	result := models.DB.Create(&project)
	models.DbMutex.Unlock()

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
	result := models.DB.
		Preload("Creator").
		Preload("Team").
		Joins("JOIN team_users ON team_users.team_id = projects.team_id").
		Where("projects.creator_id = ? OR team_users.user_id = ?", userID, userID).
		Find(&projects)
	models.DbMutex.Unlock()

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
