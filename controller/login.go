package controller

import (
	"strconv"
	"time"

	"github.com/U-to-E/dashboard/config"
	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func RenderLogin(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func Register(c fiber.Ctx) error {
	email := c.FormValue("email")
	passwd := c.FormValue("password")

	password, _ := bcrypt.GenerateFromPassword([]byte(passwd), 14)

	user := models.Student{
		Username: email,
		Email:    email,
		Password: password,
	}

	database.DB.Create(&user)

	return c.SendString("Now login in")
}

func Handlelogin(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if isAdmin(email, password) {
		return c.Redirect().To("/adminpanel")
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
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Redirect().To("/dashboard")

}

func isAdmin(email, password string) bool {
	return email == config.Config("ADMIN_EMAIL") && password == config.Config("ADMIN_PASS")
}

func authenticateUser(email, password string) (*models.Student, error) {
	var user models.Student
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func generateJWT(userID uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(userID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	return claims.SignedString([]byte(config.Config("SECRET")))
}
