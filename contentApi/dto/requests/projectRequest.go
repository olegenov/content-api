package requests

type ProjectRequest struct {
	Name string `json:"name" binding:"required"`
}
