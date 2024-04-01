package controllers

import (
	"contentApi/controllers/responseControllers"
	"contentApi/models"
	"contentApi/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPost(c *gin.Context) {
	projectID := c.Param("id")
	postID := c.Param("post-id")

	var project models.Project

	models.DbMutex.Lock()
	if err := models.DB.
		First(&project, projectID).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	models.DbMutex.Unlock()

	userID, _ := token.ExtractTokenID(c)

	var team models.Team

	models.DbMutex.Lock()
	if err := models.DB.
		Preload("Users").
		First(&team, project.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	models.DbMutex.Unlock()

	isMember := false

	for _, user := range team.Users {
		if user.ID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a member of the team"})
		return
	}

	var post models.Post

	models.DbMutex.Lock()
	if err := models.DB.
		Preload("Project").
		First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	models.DbMutex.Unlock()

	if project.ID != post.ProjectID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post is not from the project"})
		return
	}

	postResponse := responseControllers.GetSinglePostResponse(post)
	c.JSON(http.StatusOK, postResponse)
}
