package routes

import (
	"time"

	"example.com/api/database"
	"example.com/api/models"
	"github.com/gofiber/fiber/v2"
)

type Token struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Value     string `json:"value"`
}

func refreshToken(c *fiber.Ctx) error {
	var token models.Token
	res := database.Database.Db.Create(&token)
	if res.Error != nil {
		return c.Status(500).JSON(res.Error.Error())

	} else {

		return c.Status(200).JSON(token.Value)
	}

}
