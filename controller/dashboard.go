package controller

import (
	"io/ioutil"
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

	var login models.Login
	var student models.Student
	var mentor models.Mentor
	var materials []models.Material

	database.DB.Table("logins").Where("email = ?", claims.Issuer).First(&login)
	database.DB.Table(login.CollageID).Where("email = ?", login.Email).First(&student)
	database.DB.Table("mentors").Where("id = ?", student.MentorID).First(&mentor)

	files, err := ioutil.ReadDir("./materials/" + student.CollageID + "-" + student.MentorID)
	if err != nil {
		return err
	}

	for _, file := range files {
		mat := models.Material{
			Name:     file.Name(),
			FilePath: "./materials/" + student.CollageID + "-" + student.MentorID + "/" + file.Name(),
		}
		materials = append(materials, mat)
	}

	return c.Render("dashboard", fiber.Map{
		"User":      student.Name,
		"Level":     student.Level,
		"Marks":     student.Marks,
		"Mentor":    mentor,
		"Materials": materials,
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

	return c.Redirect().To("/")
}
