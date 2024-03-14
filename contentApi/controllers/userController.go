package controllers

import (
	"contentApi/dto"
	"contentApi/models"
	"contentApi/utils/token"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	var usersResponse []dto.UserResponse

	for _, user := range users {
		usersResponse = append(usersResponse, GetUserResponse(user))
	}

	c.JSON(http.StatusOK, usersResponse)
}

func GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	userResponse := GetUserResponse(user)

	c.JSON(http.StatusOK, userResponse)
}

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if exists, _ := UserExists(user.Username); exists == true {
		c.JSON(http.StatusFound, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "message": "User created successfully"})
}

func UserExists(username string) (bool, error) {
	var user models.User

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func UserExistsByID(id uint) (bool, error) {
	var user models.User

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, result.Error
	}

	return true, nil
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var id uint
	if _, err := fmt.Sscan(userID, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	exists, err := UserExistsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	result := models.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetUserByID(id uint) (models.User, error) {
	var user models.User

	exists, err := UserExistsByID(id)

	if err != nil {
		return user, err
	}

	if !exists {
		return user, errors.New("User not found!")
	}

	models.DbMutex.Lock()
	defer models.DbMutex.Unlock()

	models.DB.First(&user, id)

	return user, nil

}

func CurrentUser(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse := GetUserResponse(user)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": userResponse})
}
