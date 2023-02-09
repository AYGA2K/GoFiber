package routes

import (
	"errors"

	"log"
	"os"
	"time"

	"example.com/api/database"
	"example.com/api/models"
	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// This is not the model, more like a serializer
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

func CreateResponseUser(user models.User) User {
	return User{ID: user.ID, Name: user.Name, Email: user.Email}
}

func SignUp(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)

	res := database.Database.Db.Create(&user)
	if res.Error != nil {
		return c.Status(500).JSON(res.Error.Error())

	} else {
		responseUser := CreateResponseUser(user)
		return c.Status(200).JSON(responseUser)
	}

}

func GetUsers(c *fiber.Ctx) error {

	users := []models.User{}
	database.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {

	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
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
func Login(c *fiber.Ctx) error {
	var user *models.User
	// Get the email and password from the request body
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Check the email and password if they are valid

	database.Database.Db.Find(&user, "email = ?  ", email)

	if user.Email == "" {
		// if user not found send an 404 status
		return c.Status(404).JSON("User Not found")
	}
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if bcryptErr != nil {

		return c.Status(401).JSON(fiber.Map{"error": "Wrong Password"})
	}

	refresh_token, refresh_err := GenerateJWT("refresh", CreateResponseUser(*user), time.Hour*24)

	access_token, access_err := GenerateJWT("access", CreateResponseUser(*user), time.Minute*15)
	if refresh_err == nil && access_err == nil {
		return c.Status(200).JSON(fiber.Map{"access_token": access_token, "refresh_token": refresh_token})
	}
	return nil
}
func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findUser(id, &user)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.Name = updateData.Name
	user.Password = updateData.Password

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)

}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findUser(id, &user)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted User")
}
