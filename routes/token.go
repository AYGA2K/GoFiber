package routes

import (
	"errors"
	"example.com/api/database"
	"example.com/api/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Token struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Value  string `json:"value"`
	UserID uint   `json:"userId"`
}

func CreateAccessToken(c *fiber.Ctx) error {
	var token *models.Token
	var user *models.User
	if err := c.BodyParser(&token); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	res := database.Database.Db.Find(&token)

	if res.Error != nil {
		return c.Status(500).JSON(res.Error.Error())

	} else {
		Usererr := findUser(int(token.UserID), user)
		if Usererr != nil {
			access_token, access_err := GenerateJWT("access", CreateResponseUser(*user), time.Minute*15)
			if access_err == nil {
				return c.Status(200).JSON(fiber.Map{"access_token": access_token})
			}
		}
	}
	return nil
}
func findToken(id int, token *models.Token) error {

	database.Database.Db.Find(&token, "id = ?", id)
	if token.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}
func DeleteRefreshToken(c *fiber.Ctx) error {
	var token models.Token
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findToken(id, &token)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = database.Database.Db.Delete(&token).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Token")

}
