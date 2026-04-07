package controller

import (
	"net/http"
	"src/pkg/entities"
	"src/pkg/helpers"
	"src/pkg/repositories"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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
	}, user.Role)
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

// @Router /createOrder [post]
// @Security ApiKeyAuth
func CreateOrder(c *gin.Context) {
	var order entities.Order
	token := c.GetHeader("Authorization")

	// Validação básica de token já existente no seu helpers
	_, err := helpers.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repositories.CreateOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

// @Router /getOrders [get]
// @Security ApiKeyAuth
func GetMyOrders(c *gin.Context) {
	token := c.GetHeader("Authorization")

	// Extraindo o ID do usuário do Token (baseado no seu JWTgenerator.go)
	t, _ := helpers.ValidateJWT(token)
	claims := t.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64))

	orders, err := repositories.GetOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Router /deleteOrder/{id} [delete]
// @Security ApiKeyAuth
func DeleteOrder(c *gin.Context) {
	token := c.GetHeader("Authorization")

	// Validação de permissão: Apenas admin pode deletar
	role, err := helpers.GetRole(token)
	if err != nil || role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso negado: apenas administradores podem apagar pedidos"})
		return
	}

	// Captura o ID da URL (ex: /deleteOrder/5)
	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do pedido inválido"})
		return
	}

	// Chama o repositório para deletar
	if err := repositories.DeleteOrder(orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar pedido: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pedido " + idParam + " removido com sucesso!"})
}
