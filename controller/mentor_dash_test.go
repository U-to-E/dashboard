package controller

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestRenderMentorDash(t *testing.T) {
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
			if err := RenderMentorDash(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("RenderMentorDash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostMaterial(t *testing.T) {
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
			if err := PostMaterial(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("PostMaterial() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteMaterial(t *testing.T) {
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
			if err := DeleteMaterial(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMaterial() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
