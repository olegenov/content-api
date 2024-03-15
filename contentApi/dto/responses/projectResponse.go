package responses

type ProjectResponse struct {
	ID      uint           `json:"id"`
	Name    string         `json:"name"`
	Creator UserResponse   `json:"creator"`
	Posts   []PostResponse `json:"posts"`
}
