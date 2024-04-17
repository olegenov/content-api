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

func CreateInvitation(c *gin.Context) {
	var invitationRequest requests.InvitationRequest

	if err := c.ShouldBindJSON(&invitationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := token.ExtractTokenID(c)

	var team models.Team

	models.DbMutex.Lock()
	teamResult := models.DB.First(&team, invitationRequest.TeamID)
	models.DbMutex.Unlock()

	if errors.Is(teamResult.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	} else if teamResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve team"})
		return
	}

	var creatorUser models.User

	models.DbMutex.Lock()
	userResult := models.DB.First(&creatorUser, userID)
	models.DbMutex.Unlock()

	if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Creator user not found"})
		return
	} else if userResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve creator user"})
		return
	}

	if team.CreatorID != creatorUser.ID {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "you are not a creator"})
		return
	}

	var receiverUser models.User

	models.DbMutex.Lock()
	receiverUserResult := models.DB.Where("username = ?", invitationRequest.ReceiverUsername).First(&receiverUser)
	models.DbMutex.Unlock()

	if errors.Is(receiverUserResult.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if userResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve creator user"})
		return
	}

	invitation := models.Invitation{
		TeamID:         team.ID,
		ReceiverID:     receiverUser.ID,
		SenderUsername: creatorUser.Username,
	}

	models.DbMutex.Lock()
	result := models.DB.Create(&invitation)
	models.DbMutex.Unlock()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": team.ID, "message": "invitation created successfully"})
}

func GetMyInvitations(c *gin.Context) {
	var invitations []models.Invitation
	userID, _ := token.ExtractTokenID(c)

	models.DbMutex.Lock()

	result := models.DB.
		Preload("Receiver").
		Preload("Team").
		Where("receiver_id = ?", userID).
		Find(&invitations)

	models.DbMutex.Unlock()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get invitation"})
		return
	}

	if len(invitations) == 0 {
		c.JSON(http.StatusOK, gin.H{"invitations": []responses.InvitationResponse{}, "amount": 0})
		return
	}

	var invitationResponse []responses.InvitationResponse

	for _, invitation := range invitations {
		invitationResponse = append(invitationResponse, responseControllers.GetInvitationResponse(invitation))
	}

	c.JSON(http.StatusOK, gin.H{"invitations": invitationResponse, "amount": len(invitationResponse)})
}

func ReceiveInvitation(c *gin.Context) {
	invitationID := c.Param("id")
	userID, _ := token.ExtractTokenID(c)

	models.DbMutex.Lock()

	var invitation models.Invitation

	invitationResult := models.DB.
		Preload("Receiver").
		Preload("Team").
		First(&invitation, invitationID)

	models.DbMutex.Unlock()

	if errors.Is(invitationResult.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	} else if invitationResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invitation"})
		return
	}

	if userID != invitation.Receiver.ID {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not your invitation"})
		return
	}

	for _, user := range invitation.Team.Users {
		if user.ID == userID {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "already in this team"})
			return
		}
	}

	team := invitation.Team
	team.Users = append(team.Users, invitation.Receiver)

	models.DbMutex.Lock()

	updateResult := models.DB.Save(&team)

	models.DbMutex.Unlock()

	if updateResult.Error != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}

	models.DbMutex.Lock()

	deleteResult := models.DB.Delete(&invitation)

	models.DbMutex.Unlock()

	if deleteResult.Error != nil {
		models.DbMutex.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invitation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted successfully"})
}

func RejectInvitation(c *gin.Context) {
	invitationID := c.Param("id")
	userID, _ := token.ExtractTokenID(c)

	models.DbMutex.Lock()

	var invitation models.Invitation

	invitationResult := models.DB.
		Preload("Receiver").
		Preload("Team").
		First(&invitation, invitationID)

	models.DbMutex.Unlock()

	if errors.Is(invitationResult.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	} else if invitationResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invitation"})
		return
	}

	if userID != invitation.Receiver.ID {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not your invitation"})
		return
	}

	for _, user := range invitation.Team.Users {
		if user.ID == userID {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "already in this team"})
			return
		}
	}

	team := invitation.Team
	for i, user := range team.Users {
		if user.ID == userID {
			team.Users = append(team.Users[:i], team.Users[i+1:]...)
			break
		}
	}

	models.DbMutex.Lock()
	updateResult := models.DB.Save(&team)
	models.DbMutex.Unlock()

	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}
	models.DbMutex.Lock()

	deleteResult := models.DB.Delete(&invitation)

	models.DbMutex.Unlock()

	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invitation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation rejected successfully"})
}
