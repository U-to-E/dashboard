package controller

import (
	"time"

	"github.com/U-to-E/dashboard/config"
	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

func RenderDashboard(c fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config("SECRET")), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.Student

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.Render("dashboard", fiber.Map{
		"User": user.Username,
	})
}

func Logout(c fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
