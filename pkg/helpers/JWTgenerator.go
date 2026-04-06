package helpers

import (
	"fmt"
	"src/pkg/entities"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user entities.UserDTO, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key"))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GetRole(tokenString string) (string, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	role, ok := token.Claims.(jwt.MapClaims)["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found in token")
	}
	return role, nil
}
