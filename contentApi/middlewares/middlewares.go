package middlewares

import (
	"contentApi/models"
	"errors"
	"github.com/jinzhu/gorm"
	"net/http"

	"contentApi/utils/token"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)

		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userID, err := token.ExtractTokenID(c)

		if err != nil {
			c.String(http.StatusUnauthorized, "Cant get user")
			c.Abort()
			return
		}

		var user models.User

		models.DbMutex.Lock()
		result := models.DB.First(&user, userID)
		models.DbMutex.Unlock()

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		c.Set("userRole", user.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("userRole")

		if userRole != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admin can access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
