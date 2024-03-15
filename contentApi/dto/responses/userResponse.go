package responses

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}
