package controller

import (
	"fmt"
	"net/http"
	"src/pkg/entities"
	"src/pkg/helpers"
	"src/pkg/repositories"

	"github.com/gin-gonic/gin"
)

// @Router /getUsers [get]
// @Security ApiKeyAuth
func GetAllUsers(c *gin.Context) {

	// Get the JWT token from the request
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Validate the JWT token
	_, err := helpers.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Get the user role from the token
	role, err := helpers.GetRole(tokenString)
	if err != nil {
		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
	}

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not exist"})
		return
	}

	if user.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := helpers.GenerateJWT(entities.UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": entities.UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}})
}

// @Router /createUser [post]
// @Accept json
// @Param user body entities.User true "User"
// @Security ApiKeyAuth
func CreateUser(c *gin.Context) {
	var user entities.User
	var token string = c.GetHeader("Authorization")
	var role, _ = helpers.GetRole(token)
	fmt.Println(role)
	//fix role implementation

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repositories.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": user})
}
