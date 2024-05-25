package controller

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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
	var quizs []models.Quiz

	database.DB.Table("logins").Where("email = ?", claims.Issuer).First(&login)
	database.DB.Table(login.CollageID).Where("email = ?", login.Email).First(&student)
	database.DB.Table("mentors").Where("id = ?", student.MentorID).First(&mentor)

	files, err := ioutil.ReadDir("./materials/" + student.CollageID + "-" + student.MentorID)
	if err != nil {
		return err
	}

	quizes, err := ioutil.ReadDir("./quiz/" + student.CollageID + "-" + student.MentorID)
	if err != nil {
		return err
	}

	for _, file := range files {

		sp := strings.Split(file.Name(), ".")
		mat := models.Material{
			Name:     sp[0],
			FilePath: "/" + student.CollageID + "-" + student.MentorID + "/" + file.Name(),
		}
		materials = append(materials, mat)
	}

	for _, quiz := range quizes {
		splt := strings.Split(quiz.Name(), "|")

		var count int64
		err := database.DB.
			Table("marks").
			Where("student_id = ? AND collage_id = ? AND quiz_id = ?", strconv.Itoa(int(student.ID)), student.CollageID, splt[2]).
			Count(&count).Error

		if err != nil {
			log.Printf("Error checking record existence: %v", err)
			continue
		}

		// If no records were found, create a new Quiz entry and append it to quizs slice
		if count == 0 {
			qz := models.Quiz{
				Name:     splt[0],
				Duration: splt[1],
				QuizID:   splt[2],
				Path:     "./quiz/" + student.CollageID + "-" + student.MentorID + "/" + quiz.Name(),
			}
			quizs = append(quizs, qz)
		}
	}

	return c.Render("dashboard", fiber.Map{
		"User":      student,
		"SID":       student.ID,
		"Level":     student.Level,
		"Marks":     student.Marks,
		"Mentor":    mentor,
		"Materials": materials,
		"Quiz":      quizs,
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
	sess, err := store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to destroy session")
	}
	err = sess.Destroy()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to destroy session")
	}

	return c.Redirect().To("/")
}
