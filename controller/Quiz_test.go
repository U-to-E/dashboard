package controller

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestQuizPage(t *testing.T) {
	type args struct {
		c fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := QuizPage(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("QuizPage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmitQuiz(t *testing.T) {
	type args struct {
		c fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SubmitQuiz(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SubmitQuiz() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
