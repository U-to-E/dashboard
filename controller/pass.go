package controller

import (
	"errors"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/U-to-E/dashboard/config"
	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/models"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ChangePassword(c fiber.Ctx) error {
	oldpass := c.FormValue("old_password")
	newpass := c.FormValue("new_password")
	confirmpass := c.FormValue("confirm_password")

	if newpass == "" || confirmpass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password and confirmation password cannot be empty",
		})
	}

	if newpass != confirmpass {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New passwords do not match",
		})
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get session")
	}

	email := sess.Get("email").(string)

	var login models.Login
	if err := database.DB.Table("logins").Where("email = ?", email).First(&login).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(oldpass)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Password comparison error",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash new password",
		})
	}

	if err := database.DB.Table("logins").Where("email = ?", email).Update("password", hashedPassword).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

func ChangePassPage(c fiber.Ctx) error {
	return c.Render("changepasswd", fiber.Map{})
}

func generateResetToken() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func sendResetEmail(email, token string) error {
	from := config.Config("EMAIL_ADDR")
	password := config.Config("EMAIL_PASS")
	to := email
	smtpHost := config.Config("EMAIL_SMTP")
	smtpPort := config.Config("EMAIL_PORT")

	msg := "Subject: Password Reset Request\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<html><head><style>" +
		"body { font-family: Arial, sans-serif; }" +
		".container { max-width: 600px; margin: auto; }" +
		".logo { max-width: 150px; }" +
		".reset-link { background-color: #4CAF50; border: none; color: white; padding: 15px 32px; text-align: center; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 8px;}" +
		"</style></head><body>" +
		"<div class='container'>" +
		"<img src='https://localhost:300/asserts/UtoElogo.png' alt='Company Logo' class='logo'>" +
		"<h2>Password Reset Request</h2>" +
		"<p>To reset your password, please click the link below:</p>" +
		"<p><a href='http://localhost:3000/reset/password?token=" + token + "' class='reset-link'>Reset Password</a></p>" +
		"<p>This link will expire in 1 hour.</p>" +
		"<p>Please do not reply to this email.</p>" +
		"</div>" +
		"</body></html>"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}

func ForgotPassword(c fiber.Ctx) error {
	email := c.FormValue("email")

	var user models.Login
	if err := database.DB.Table("logins").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	token := generateResetToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	resetToken := models.PasswordResetToken{
		Email:     user.Email,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Table("password_reset_tokens").Create(&resetToken).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store reset token",
		})
	}

	if err := sendResetEmail(user.Email, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset email sent",
	})
}

func ResetPassword(c fiber.Ctx) error {
	token := c.FormValue("token")
	newpass := c.FormValue("new_password")
	confirmpass := c.FormValue("confirm_password")

	if newpass == "" || confirmpass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password and confirmation password cannot be empty",
		})
	}

	if newpass != confirmpass {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New passwords do not match",
		})
	}

	var resetToken models.PasswordResetToken
	if err := database.DB.Where("token = ?", token).First(&resetToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token has expired",
		})
	}

	var user models.Login
	if err := database.DB.Table("logins").Where("email = ?", resetToken.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash new password",
		})
	}

	user.Password = hashedPassword
	if err := database.DB.Table("logins").Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	if err := database.DB.Table("password_reset_tokens").Delete(&resetToken).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete reset token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully. Now please login",
	})

}

func RenderForgotPass(c fiber.Ctx) error {
	return c.Render("forgotpass", fiber.Map{})
}

func RenderResetPass(c fiber.Ctx) error {
	token := c.Query("token")
	return c.Render("resetpass", fiber.Map{
		"TOKEN": token,
	})
}
