package controllers

import (
	"contentApi/controllers/responseControllers"
	"contentApi/dto/requests"
	"contentApi/models"
	"contentApi/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		Preload("Assign").
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

func EditPost(c *gin.Context) {
	postID := c.Param("post-id")

	var post models.Post

	models.DbMutex.Lock()
	if err := models.DB.First(&post, postID).Error; err != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var request requests.EditPostRequest
	if err := c.BindJSON(&request); err != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.Title != "" {
		post.Title = request.Title
	}
	if request.Content != "" {
		post.Content = request.Content
	}

	post.Deadline = request.Deadline
	post.PublishDate = request.PublishDate

	var assign models.User

	if err := models.DB.Where("username = ?", request.Assign).First(&assign).Error; err != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "assign user not found"})
		return
	}

	post.Assign = assign

	if err := models.DB.Save(&post).Error; err != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	models.DbMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
}

func CreatePost(c *gin.Context) {
	var request requests.EditPostRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var assign models.User

	models.DbMutex.Lock()

	if err := models.DB.Where("username = ?", request.Assign).First(&assign).Error; err != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "assign user not found"})
		return
	}

	projectID := c.Param("id")
	id, _ := strconv.Atoi(projectID)

	post := models.Post{
		Title:       request.Title,
		Assign:      assign,
		PublishDate: request.PublishDate,
		Deadline:    request.Deadline,
		ProjectID:   uint(id),
		Content:     request.Content,
	}

	if err := models.DB.Create(&post).Error; err != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Creation error"})
		return
	}

	models.DbMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "post created successfully"})
}
