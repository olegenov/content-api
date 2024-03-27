package controllers

import (
	"contentApi/controllers/responseControllers"
	"contentApi/dto/requests"
	"contentApi/dto/responses"
	"contentApi/models"
	"contentApi/utils/token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetTeams(c *gin.Context) {
	var teams []models.Team

	models.DbMutex.Lock()

	result := models.DB.Preload("Creator").Find(&teams)

	models.DbMutex.Unlock()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams"})
		return
	}

	var teamResponse []responses.TeamResponse

	for _, team := range teams {
		teamResponse = append(teamResponse, responseControllers.GetTeamResponse(team))
	}

	c.JSON(http.StatusOK, teamResponse)
}

func GetTeam(c *gin.Context) {
	teamID := c.Param("id")

	var team models.Team

	models.DbMutex.Lock()

	if err := models.DB.Preload("Creator").First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	models.DbMutex.Unlock()

	userID, _ := token.ExtractTokenID(c)

	if userID != team.CreatorID && team.Creator.Role != "ADMIN" {
		c.JSON(http.StatusNotFound, gin.H{"error": "It is not your team"})
		return
	}

	teamResponse := responseControllers.GetTeamResponse(team)
	c.JSON(http.StatusOK, teamResponse)
}

func CreateTeam(c *gin.Context) {
	var teamRequest requests.TeamRequest

	if err := c.ShouldBindJSON(&teamRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := token.ExtractTokenID(c)

	team := models.Team{
		Name:      teamRequest.Name,
		CreatorID: userID,
	}

	models.DbMutex.Lock()

	var creatorUser models.User

	userResult := models.DB.First(&creatorUser, userID)

	if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Creator user not found"})
		return
	} else if userResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve creator user"})
		return
	}

	team.Users = append(team.Users, creatorUser)

	result := models.DB.Create(&team)

	models.DbMutex.Unlock()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": team.ID, "message": "Team created successfully"})
}

func GetMyTeams(c *gin.Context) {
	var teams []models.Team
	userID, _ := token.ExtractTokenID(c)

	models.DbMutex.Lock()

	result := models.DB.
		Preload("Creator").
		Preload("Users").
		Joins("JOIN team_users ON teams.id = team_users.team_id").
		Where("team_users.user_id = ?", userID).
		Find(&teams)

	models.DbMutex.Unlock()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams"})
		return
	}

	if len(teams) == 0 {
		c.JSON(http.StatusOK, gin.H{"teams": []responses.TeamResponse{}})
		return
	}

	var teamResponse []responses.TeamResponse

	for _, team := range teams {
		teamResponse = append(teamResponse, responseControllers.GetTeamResponse(team))
	}

	c.JSON(http.StatusOK, gin.H{"teams": teamResponse})
}
