package controller

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/U-to-E/dashboard/config"
	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RenderLogin(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func Handlelogin(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if isAdmin(email, password) {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(9899)),
			Subject:   "1",
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // 1 hour
		})

		token, err := claims.SignedString([]byte(config.Config("SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 1),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})
		c.Set("HX-Redirect", "/admin/panel")

		return c.SendStatus(fiber.StatusOK)
	}

	user, mentor, err := authenticateUser(email, password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
	}

	if mentor != nil {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    mentor.Email,
			Id:        strconv.FormatUint(uint64(mentor.ID), 10),
			Subject:   "2",
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // 1 hour
		})
		token, err := claims.SignedString([]byte(config.Config("SECRET")))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		c.Locals("mentor", mentor)

		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
		}
		sess.Set("mentor_id", mentor.ID)
		sess.Set("email", mentor.Email)

		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 1),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})

		c.Set("HX-Redirect", "/mentor/dashboard")

		return c.SendStatus(fiber.StatusOK)
	}

	token, err := generateJWT(user.Email, user.CollageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	c.Locals("user", user)

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
	}
	sess.Set("user_id", user.ID)
	sess.Set("collage_id", user.CollageID)
	sess.Set("email", user.Email)

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
	}

	c.Cookie(&fiber.Cookie{
		Name:        "jwt",
		Value:       token,
		Expires:     time.Now().Add(time.Hour * 1),
		HTTPOnly:    true,
		Secure:      true,
		SameSite:    "Strict",
		SessionOnly: true,
	})

	c.Set("HX-Redirect", "/student/dashboard")

	return c.SendStatus(fiber.StatusOK)

}

func isAdmin(email, password string) bool {
	return email == config.Config("ADMIN_EMAIL") && password == config.Config("ADMIN_PASS")
}

func authenticateUser(email, password string) (*models.Student, *models.Mentor, error) {
	var user models.Login
	if err := database.DB.Table("logins").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, fmt.Errorf("incorrect password for user %s", email)
	}

	if user.Role == "Mentor" {
		var mentor models.Mentor
		if err := database.DB.Table("mentors").Where("email = ?", email).First(&mentor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, fmt.Errorf("user with email %s not found", email)
			}
			return nil, nil, err
		}

		return nil, &mentor, nil
	}
	var student models.Student
	if err := database.DB.Table(user.CollageID).Where("email = ?", email).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, nil, err
	}

	return &student, nil, nil
}

func generateJWT(userEmail string, collageId string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        collageId,
		Issuer:    userEmail,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	})

	return claims.SignedString([]byte(config.Config("SECRET")))
}
