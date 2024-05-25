package controller

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

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
	quizes, err := ioutil.ReadDir("./quiz/" + mentor.AssignedCollage + "-" + strconv.FormatUint(uint64(mentor.ID), 10))
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
	var quizzzz []models.Quiz
	for _, quiz := range quizes {
		qw := models.Quiz{
			Name: quiz.Name(),
			Path: "./quiz/" + mentor.AssignedCollage + "-" + strconv.FormatUint(uint64(mentor.ID), 10),
		}
		quizzzz = append(quizzzz, qw)
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
		"Quiz":     quizzzz,
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

func generateQuizID(mentor models.Mentor) string {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(900) + 100
	quizID := fmt.Sprintf("%s-%d-%03d", mentor.AssignedCollage, mentor.ID, randNum)
	return quizID
}

func CreateQuiz(c fiber.Ctx) error {
	Qduration := c.FormValue("Qduration")
	Qname := c.FormValue("Qname")
	mentorID := c.FormValue("mentorId")
	file, err := c.FormFile("quizfile")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error uploading file",
		})
	}

	var mentor models.Mentor

	database.DB.Table("mentors").Where("id = ?", mentorID).First(&mentor)

	if err := c.SaveFile(file, "./quiz/"+mentor.AssignedCollage+"-"+mentorID+"/"+file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving the quiz file")
	}

	quizID := generateQuizID(mentor)

	if err := os.Rename("./quiz/"+mentor.AssignedCollage+"-"+mentorID+"/"+file.Filename, "./quiz/"+mentor.AssignedCollage+"-"+mentorID+"/"+Qname+"|"+Qduration+"|"+quizID+"|.csv"); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving quiz")
	}

	quiz := models.Quiz{
		Name:     Qname,
		QuizID:   quizID,
		Duration: Qduration,
		Path:     "./quiz/" + mentor.AssignedCollage + "-" + mentorID + "/" + Qname + "|" + Qduration + "|" + quizID + "|.csv",
	}

	database.DB.Table("quizzes").Create(&quiz)

	return c.SendString("Done")

}

func GetStudentPage(c fiber.Ctx) error {

	studentID := c.Params("id")

	spl := strings.Split(studentID, "-")

	collageID := spl[0]
	stuID := spl[1]

	studID, err := strconv.Atoi(stuID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error Converting number")
	}

	var student models.Student

	if err := database.DB.Table(collageID).
		Where("id = ?", studID).
		Find(&student).Error; err != nil {
		return c.SendString("Issue with DB")
	}

	sess, err := store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get Session")
	}

	mentorID, ok := sess.Get("mentor_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access")
	}
	if student.MentorID != strconv.FormatUint(uint64(mentorID), 10) {
		return c.Status(fiber.StatusForbidden).SendString("Not Your Student")
	}

	var marks []models.Marks
	err = database.DB.Table("marks").Where("student_id = ? AND collage_id = ?", stuID, collageID).Find(&marks).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch marks")
	}

	certFilePath := fmt.Sprintf("./certificates/%s-%s.pdf", collageID, stuID)
	certExists := false
	var filename string
	if file, err := os.Stat(certFilePath); err == nil {
		certExists = true
		filename = file.Name()
	}

	return c.Render("student", fiber.Map{
		"Student":    student,
		"Marks":      marks,
		"CertExists": certExists,
		"CertPath":   filename,
	})
}

func SetLevel(c fiber.Ctx) error {
	studentID := c.Params("id")
	level := c.FormValue("level")

	spl := strings.Split(studentID, "-")
	if len(spl) != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid student ID format")
	}

	collageID := spl[0]
	stuID := spl[1]

	studID, err := strconv.Atoi(stuID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error Converting number")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get Session")
	}

	mentorID, ok := sess.Get("mentor_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access")
	}

	var student models.Student

	if err := database.DB.Table(collageID).
		Where("id = ?", studID).
		Find(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Issue with DB")
	}

	fmt.Println(student.MentorID, mentorID)
	if student.MentorID != strconv.FormatUint(uint64(mentorID), 10) {
		return c.Status(fiber.StatusForbidden).SendString("Not Your Student")
	}

	student.Level, err = strconv.Atoi(level)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error Conversion")
	}
	if err := database.DB.Table(collageID).Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update student level")
	}

	return c.SendString("Set level to " + level)
}

func UploadCert(c fiber.Ctx) error {
	studentID := c.Params("id")
	file, err := c.FormFile("cert")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to retrieve file")
	}

	spl := strings.Split(studentID, "-")
	if len(spl) != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid student ID format")
	}

	collageID := spl[0]
	stuID := spl[1]

	_, err = strconv.Atoi(stuID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error converting student ID")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
	}

	mentorID, ok := sess.Get("mentor_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access")
	}

	var student models.Student
	if err := database.DB.Table(collageID).
		Where("id = ?", stuID).
		Find(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database issue: unable to find student")
	}

	if student.MentorID != strconv.FormatUint(uint64(mentorID), 10) {
		return c.Status(fiber.StatusForbidden).SendString("Not your student")
	}

	filename := fmt.Sprintf("%s-%s.pdf", collageID, stuID)
	filePath := fmt.Sprintf("./certificates/%s", filename)

	if _, err := os.Stat("./certificates"); os.IsNotExist(err) {
		err := os.MkdirAll("./certificates", os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create directory")
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete existing certificate")
		}
	}

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
	}

	return c.SendString("Certificate uploaded successfully")
}

func DeleteCert(c fiber.Ctx) error {
	studentID := c.Params("id")

	spl := strings.Split(studentID, "-")
	if len(spl) != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid student ID format")
	}

	collageID := spl[0]
	stuID := spl[1]

	studID, err := strconv.Atoi(stuID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error converting student ID")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
	}

	mentorID, ok := sess.Get("mentor_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access")
	}

	var student models.Student
	if err := database.DB.Table(collageID).
		Where("id = ?", studID).
		Find(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database issue: unable to find student")
	}

	if student.MentorID != strconv.FormatUint(uint64(mentorID), 10) {
		return c.Status(fiber.StatusForbidden).SendString("Not your student")
	}

	filePath := fmt.Sprintf("./certificates/%s-%s.pdf", collageID, stuID)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).SendString("Certificate not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete certificate")
	}

	return c.SendString("Certificate deleted successfully")
}
