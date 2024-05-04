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

func Register(c fiber.Ctx) error {
	email := c.FormValue("email")
	passwd := c.FormValue("password")
	collageID := c.FormValue("CollageID")
	name := c.FormValue("name")

	password, _ := bcrypt.GenerateFromPassword([]byte(passwd), 14)

	user := models.Login{
		Email:     email,
		Password:  password,
		CollageID: collageID,
	}

	student := models.Student{
		Name:      name,
		CollageID: collageID,
		Level:     0,
		Marks:     0,
	}

	database.DB.Create(&user)
	database.DB.Table(collageID).Create(student)

	return c.SendString("Now login in")
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
		return c.Redirect().To("/admin/panel")
	}

	user, err := authenticateUser(email, password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	c.Locals("user", user)

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Redirect().To("/student/dashboard")

}

func isAdmin(email, password string) bool {
	return email == config.Config("ADMIN_EMAIL") && password == config.Config("ADMIN_PASS")
}

func authenticateUser(email, password string) (*models.Student, error) {
	var user models.Login
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("incorrect password for user %s", email)
	}

	var student models.Student

	if err := database.DB.Table(user.CollageID).Where("email = ?", email).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	return &student, nil
}

func generateJWT(userID uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(userID)),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	})

	return claims.SignedString([]byte(config.Config("SECRET")))
}
