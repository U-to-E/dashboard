package controller

import (
	"strconv"

	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
)

func RenderMarks(c fiber.Ctx) error {

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
	}
	userID, ok := sess.Get("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized. Please Login"})
	}

	struser := strconv.FormatUint(uint64(userID), 10)

	collageID := sess.Get("collage_id")

	var marks []models.Marks
	err = database.DB.Table("marks").Where("student_id = ? AND collage_id = ?", struser, collageID).Find(&marks).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch marks")
	}
	return c.Render("marks", fiber.Map{
		"Marks": marks,
	})
}
