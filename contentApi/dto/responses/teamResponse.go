package responses

type TeamResponse struct {
	ID      uint           `json:"id"`
	Name    string         `json:"name"`
	Creator UserResponse   `json:"creator"`
	Users   []UserResponse `json:"users"`
}

type ProjectTeamResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatorID uint   `json:"creator"`
}
