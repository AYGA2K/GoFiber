package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type User struct {
	// This is not the model, more like a serializer
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}
type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(token_type string, user User, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(expiry).Unix(),
	})
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if token_type == "access" {
		// Generate encoded token and send it as response.
		access := os.Getenv("ACCESS_KEY")

		tokenString, err := token.SignedString([]byte(access))

		return tokenString, err
	}
	if token_type == "refresh" {
		// Generate encoded token and send it as response.
		refresh := os.Getenv("REFRESH_KEY")

		tokenString, err := token.SignedString([]byte(refresh))

		return tokenString, err
	}
	return "", nil
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
