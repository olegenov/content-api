package requests

type InvitationRequest struct {
	TeamID           uint   `json:"team_id" binding:"required"`
	ReceiverUsername string `json:"receiver_username" binding:"required"`
}
