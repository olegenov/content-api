package controllers

import (
	"contentApi/dto"
	"contentApi/models"
)

func GetUserResponse(user models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
