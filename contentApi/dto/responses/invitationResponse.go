package responses

type InvitationResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Sender string `json:"sender"`
}
