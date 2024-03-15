package responseControllers

import (
	"contentApi/dto/responses"
	"contentApi/models"
)

func GetUserResponse(user models.User) responses.UserResponse {
	return responses.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
