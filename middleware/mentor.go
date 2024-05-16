package middleware

import (
	"github.com/U-to-E/dashboard/config"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

func MentorProtected(c fiber.Ctx) error {
	token := c.Cookies("jwt")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config("SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if claims.Subject != "2" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden access",
		})
	}
	return c.Next()
}
