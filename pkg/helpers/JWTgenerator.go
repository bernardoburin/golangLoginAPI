package helpers

import (
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
