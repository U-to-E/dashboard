package controller

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type Table struct {
	TableName string `gorm:"column:table_name"`
}

func RenderAdmin(c fiber.Ctx) error {

	table := CollageIDList()
	if c.IP() != "127.0.0.1" && c.IP() != "::1" {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	}

	return c.Render("admin", fiber.Map{
		"Table": table,
	})
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

	for _, record := range records[1:] {
		if len(record) != 5 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parameters in CSV file",
			})
		}
		var mentor models.Mentor
		if err := database.DB.Table("mentors").Where("id = ? AND collage_id = ?", record[4], record[3]).First(&mentor).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Mentor not assigned to this collage",
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
			MentorID:  record[4],
			Level:     1,
			Marks:     0,
		}

		if !database.DB.Migrator().HasTable(record[3]) {
			err := database.DB.Table(record[3]).Migrator().CreateTable(&models.Student{})
			if err != nil {
				panic("failed to create table")
			}
		}

		if err := database.DB.Table(record[3]).Create(&student).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating student record",
			})
		}

		materialsDir := "./materials/" + record[3] + "-" + record[4]
		quizDir := "./quiz/" + record[3] + "-" + record[4]

		if err := os.MkdirAll(materialsDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating materials directory",
			})
		}

		if err := os.MkdirAll(quizDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating quiz directory",
			})
		}

		if err := database.DB.Table("logins").Create(&login).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating student login",
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

	for _, record := range records[1:] {
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

		if err := database.DB.Table("logins").Create(&login).Error; err != nil {
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

func filterTablesWithPrefix(tables []Table, prefix string) []Table {
	var filtered []Table
	for _, table := range tables {
		if strings.HasPrefix(strings.ToLower(table.TableName), strings.ToLower(prefix)) {
			filtered = append(filtered, table)
		}
	}
	return filtered
}

func CollageIDList() *[]Table {
	var tables []Table
	if err := database.DB.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		panic(err)
	}

	filteredTables := filterTablesWithPrefix(tables, "U")

	return &filteredTables

}

func AddSingleStudent(c fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	pass := c.FormValue("password")
	cID := c.FormValue("collageID")
	mentorID := c.FormValue("mentorID")

	login := models.Login{
		Name:      name,
		Email:     email,
		Password:  []byte(pass),
		CollageID: cID,
		Role:      "Student",
	}

	student := models.Student{
		Name:      name,
		Email:     email,
		CollageID: cID,
		MentorID:  mentorID,
		Level:     1,
		Marks:     0,
	}
	var mentor models.Mentor
	if err := database.DB.Table("mentors").Where("id = ? AND collage_id = ?", mentorID, cID).First(&mentor).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Mentor not assigned to this collage",
		})
	}

	if !database.DB.Migrator().HasTable(cID) {
		err := database.DB.Table(cID).Migrator().CreateTable(&models.Student{})
		if err != nil {
			panic("failed to create table")
		}
	}

	if err := database.DB.Table(cID).Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating student record",
		})
	}

	materialsDir := "./materials/" + cID
	quizDir := "./quiz/" + cID

	if err := os.MkdirAll(materialsDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating materials directory",
		})
	}

	if err := os.MkdirAll(quizDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating quiz directory",
		})
	}

	if err := database.DB.Table("logins").Create(&login).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating student login",
		})
	}

	return c.SendString("Student added successfully")
}
func AddSingleMentor(c fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	pass := c.FormValue("password")
	phone := c.FormValue("phnumber")

	login := models.Login{
		Name:     name,
		Email:    email,
		Password: []byte(pass),
		Role:     "Mentor",
	}

	mentor := models.Mentor{
		Name:            name,
		Email:           email,
		PhoneNumber:     phone,
		AssignedCollage: "0",
	}

	if err := database.DB.Table("logins").Create(&login).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating mentor login",
		})
	}

	if err := database.DB.Table("mentors").Create(&mentor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating mentor record",
		})
	}
	return c.SendString("Mentors added successfully")

}

func PostCID(c fiber.Ctx) error {
	CID := c.FormValue("collageId")
	return c.Redirect().To("/admin/panel/collage/" + CID)
}

func GetStudentList(c fiber.Ctx) error {
	cid := c.Params("id")
	var students []models.Student

	err := database.DB.Table(cid).Find(&students).Error
	if err != nil {
		return c.SendString("Issue with DB")
	}

	return c.Render("studenttable", fiber.Map{
		"Students": students,
	})
}

func MapMentorToCollage(c fiber.Ctx) error {
	mentorID := c.FormValue("mentorID")
	collageID := c.FormValue("collageId")

	if err := database.DB.Table("mentors").Where("id = ?", mentorID).Update("CID", collageID).Error; err != nil {
		return err
	}

	return c.SendString("Updated")
}
