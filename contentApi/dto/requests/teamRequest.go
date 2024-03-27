package requests

type TeamRequest struct {
	Name string `json:"name" binding:"required"`
}
