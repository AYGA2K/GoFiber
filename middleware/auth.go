package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the token from the authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		c.Status(401).Send([]byte("Authorization header is empty"))
		return nil
	}
	godotenv.Load()

	tokenString := authHeader[7:] // Remove "Bearer " from the header

	secret := os.Getenv("ACCESS_KEY")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	// Check the signing method
	if err != nil {

		c.Status(401).Send([]byte(err.Error()))
	} else if token.Valid {

		return c.Next()
	}

	return nil
}
