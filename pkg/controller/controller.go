package controller

import (
	"src/pkg/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"src/pkg/entities"
)

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

func login(c *gin.Context){
	var request entities.LoginRequest
}