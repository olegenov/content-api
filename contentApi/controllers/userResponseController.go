package controllers

import "contentApi/models"

func GetUserResponse(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
