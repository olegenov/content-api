package requests

type ProjectRequest struct {
	Name   string `json:"name" binding:"required"`
	TeamID uint   `json:"team_id" binding:"required"`
}
