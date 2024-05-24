package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var store *session.Store

func Session(s *session.Store) {
	store = s
}

func readCSV(filename string) ([]models.Question, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var quizzes []models.Question
	for _, record := range records[1:] { // Skip header
		quizzes = append(quizzes, models.Question{
			Question: record[0],
			Options:  record[1:5],
			Answer:   record[5],
		})
	}
	return quizzes, nil
}

func shuffleOptions(options []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
}

func QuizPage(c fiber.Ctx) error {
	QID := c.Query("QID")
	SID := c.Query("SID")

	// Redirect setup (if needed)
	c.Set("HX-Redirect", "/student/dashboard/quiz?QID="+QID+"&SID="+SID)

	// Retrieve user ID from session
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
	}

	userID, ok := sess.Get("user_id").(uint)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized. Please Login"})
	}

	// Convert SID to integer
	SIDInt, err := strconv.Atoi(SID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Student ID format"})
	}

	if userID != uint(SIDInt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Student ID"})
	}

	if QID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("QID is required")
	}

	// Check if the student has already completed the quiz
	var count int64
	err = database.DB.Table("marks").
		Where("student_id = ? AND quiz_id = ?", SID, QID).
		Count(&count).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	if count > 0 {
		return c.Status(fiber.StatusConflict).SendString("Quiz already completed")
	}

	// Fetch quiz details
	var quiz models.Quiz
	err = database.DB.Table("quizzes").Where("quiz_id = ?", QID).First(&quiz).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch quiz details")
	}

	// Parse quiz duration
	value, err := strconv.Atoi(quiz.Duration)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid duration value")
	}
	dur := time.Duration(value) * time.Minute

	if sess.Get("quiz") == nil {
		quizzes, err := readCSV(quiz.Path)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read quizzes")
		}

		for i := range quizzes {
			shuffleOptions(quizzes[i].Options)
		}

		quizzesJSON, err := json.Marshal(quizzes)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to marshal quizzes")
		}

		sess.Set("quiz", string(quizzesJSON))
		sess.Set("start_time", time.Now())
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
		}
	}

	quizzesJSON := sess.Get("quiz")
	var quizzes []models.Question
	if quizzesJSON != nil {
		quizzesStr, ok := quizzesJSON.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Quiz data in session is not a string")
		}

		if err := json.Unmarshal([]byte(quizzesStr), &quizzes); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to unmarshal quizzes")
		}
	}

	startTimeInterface := sess.Get("start_time")
	startTime, ok := startTimeInterface.(time.Time)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid start time in session")
	}
	remainingTime := dur - time.Since(startTime)
	if remainingTime < 0 {
		remainingTime = 0
	}

	return c.Render("quiz", fiber.Map{
		"QID":      QID,
		"SID":      SID,
		"Quizzes":  quizzes,
		"TimeLeft": int(remainingTime.Seconds()),
	})
}

func SubmitQuiz(c fiber.Ctx) error {
	QuizID := c.FormValue("QID")
	StudentID := c.FormValue("SID")

	if QuizID == "" || StudentID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("QuizID and StudentID are required")
	}

	spl := strings.Split(QuizID, "-")
	if len(spl) < 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid QuizID format")
	}

	var quiz models.Quiz
	if err := database.DB.Table("quizzes").Where("quiz_id = ?", QuizID).First(&quiz).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to find quiz")
	}

	submittedAnswers := make(map[string]string)
	values, err := url.ParseQuery(string(c.Request().Body()))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse submitted answers")
	}

	for key, value := range values {
		if len(value) > 0 {
			submittedAnswers[key] = value[0]
		}
	}

	questions, err := readCSV(quiz.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to load quiz questions")
	}

	totalQuestions := len(questions)
	correctQuestions := 0

	for i, question := range questions {
		if submittedAnswers[fmt.Sprintf("answer[%d]", i)] == question.Answer {
			correctQuestions++
		}
	}

	// Save the results in the database
	marks := models.Marks{
		QuizName:         quiz.Name,
		CollageID:        spl[0],
		MentorID:         spl[1],
		QuizID:           QuizID,
		StudentID:        StudentID,
		TotalQuestions:   strconv.Itoa(totalQuestions),
		CorrectQuestions: strconv.Itoa(correctQuestions),
	}

	if err := database.DB.Table("marks").Create(&marks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save quiz results")
	}

	sess, err := store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	sess.Delete("quiz")
	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Quiz submitted successfully",
		"correct": correctQuestions,
		"total":   totalQuestions,
	})
}
