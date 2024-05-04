package controller

import (
	"encoding/csv"

	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func RenderAdmin(c fiber.Ctx) error {

	if c.IP() != "127.0.0.1" && c.IP() != "::1" {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	}

	return c.Render("admin", fiber.Map{})
}

func AddStudent(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error uploading file",
		})
	}

	csvFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error opening file",
		})
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error reading CSV file",
		})
	}

	for _, record := range records {
		if len(record) != 4 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameters in CSV file",
			})
		}

		login := models.Login{
			Name:      record[0],
			Email:     record[1],
			Password:  []byte(record[2]),
			CollageID: record[3],
			Role:      "Student",
		}

		student := models.Student{
			Name:      record[0],
			Email:     record[1],
			CollageID: record[3],
			MentorID:  "0",
			Level:     1,
			Marks:     1,
		}

		if err := database.DB.Table("logins").Create(&login).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating student login",
			})
		}

		if !database.DB.Migrator().HasTable(record[3]) {
			err := database.DB.Migrator().CreateTable(&models.Student{})
			if err != nil {
				panic("failed to create table")
			}
		}

		if err := database.DB.Table(record[3]).Create(&student).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating student record",
			})
		}
	}

	return c.SendString("Students added successfully")
}

func AddMentor(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error uploading file",
		})
	}

	csvFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error opening file",
		})
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error reading CSV file",
		})
	}

	for _, record := range records {
		if len(record) != 3 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameters in CSV file",
			})
		}

		login := models.Login{
			Name:     record[0],
			Email:    record[1],
			Password: []byte(record[2]),
			Role:     "Mentor",
		}

		mentor := models.Mentor{
			Name:  record[0],
			Email: record[1],
		}

		if err := database.DB.Table("students").Create(&login).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating mentor login",
			})
		}

		if err := database.DB.Table("mentors").Create(&mentor).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating mentor record",
			})
		}
	}

	return c.SendString("Mentors added successfully")
}

func PasswordGenrator(c fiber.Ctx) error {
	pass := c.FormValue("password")
	password, _ := bcrypt.GenerateFromPassword([]byte(pass), 14)

	return c.Send(password)
}
