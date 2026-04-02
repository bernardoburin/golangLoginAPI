package controller

import (
	"net/http"
	"src/pkg/entities"
	"src/pkg/repositories"

	"github.com/gin-gonic/gin"
)

// @Router /getUsers [get]
func GetAllUsers(c *gin.Context) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users found"})
		return
	}
	c.JSON(http.StatusOK, users)

}

// @Router /login [post]
// @Accept json
// @Param loginRequest body entities.LoginRequest true "Login Request"
func Login(c *gin.Context) {
	var request entities.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repositories.GetUser(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if user.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
