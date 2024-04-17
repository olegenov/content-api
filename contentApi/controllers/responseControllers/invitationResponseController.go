package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetInvitationResponse(invitation models.Invitation) responses.InvitationResponse {
	invitationResponse := responses.InvitationResponse{
		ID:     invitation.ID,
		Name:   invitation.Team.Name,
		Sender: invitation.SenderUsername,
	}

	return invitationResponse
}
