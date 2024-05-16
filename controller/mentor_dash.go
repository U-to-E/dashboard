package controller

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/U-to-E/dashboard/config"
	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

func RenderMentorDash(c fiber.Ctx) error {
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
	var mentor models.Mentor

	database.DB.Table("logins").Where("email = ?", claims.Issuer).First(&login)
	database.DB.Table("mentors").Where("email = ?", login.Email).First(&mentor)
	files, err := ioutil.ReadDir("./materials/" + mentor.AssignedCollage + "-" + strconv.FormatUint(uint64(mentor.ID), 10))
	if err != nil {
		return err
	}
	var materials []models.Material
	for _, file := range files {
		mat := models.Material{
			Name:     file.Name(),
			FilePath: "./materials/" + mentor.AssignedCollage + "-" + strconv.FormatUint(uint64(mentor.ID), 10),
		}
		materials = append(materials, mat)
	}

	var students []models.Student

	if err := database.DB.Table(mentor.AssignedCollage).
		Where("mentor_id = ?", strconv.FormatUint(uint64(mentor.ID), 10)).
		Find(&students).Error; err != nil {
		return c.SendString("Issue with DB")
	}

	return c.Render("mentordash", fiber.Map{
		"User":     mentor,
		"Students": students,
		"Material": materials,
	})
}

func PostMaterial(c fiber.Ctx) error {
	file, err := c.FormFile("material")
	mentorID := c.FormValue("mentorId")

	var mentor models.Mentor

	database.DB.Table("mentors").Where("id = ?", mentorID).First(&mentor)

	if err != nil {
		return err
	}
	if err := c.SaveFile(file, "./materials/"+mentor.AssignedCollage+"-"+mentorID+"/"+file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving the file")
	}

	return c.SendString("Done")

}

func DeleteMaterial(c fiber.Ctx) error {

	mentorID := c.FormValue("mentorId")
	file := c.FormValue("file")
	var mentor models.Mentor

	database.DB.Table("mentors").Where("id = ?", mentorID).First(&mentor)

	if err := os.Remove("./materials/" + mentor.AssignedCollage + "-" + mentorID + "/" + file); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting the file")
	}

	return c.SendString("Done")

}
