package middleware

import (
	"github.com/U-to-E/dashboard/config"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

func Protected(c fiber.Ctx) error {
	token := c.Cookies("jwt")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized. Please Login"})
	}

	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config("SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized. Please Login"})
	}

	c.Locals("claims", claims)

	return c.Next()
}
