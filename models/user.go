package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type SignUpInput struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" `
	PasswordConfirm string    `json:"password_confirm" bson:"password_confirm,omitempty"`
	PhoneNumber     string    `json:"phone_number" bson:"phone_number" binding:"required"`
	ImageUrl        string    `json:"image_url" bson:"image_url" `
	Role            string    `json:"role" bson:"role"`
	Gender          Gender    `json:"gender" bson:"gender"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type DBResponse struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Role        string             `json:"role" bson:"role"`
	Verified    bool               `json:"verified" bson:"verified"`
	Gender      string             `json:"gender" bson:"gender"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	ImageUrl    string             `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Role        string             `json:"role,omitempty" bson:"role,omitempty"`
	ImageUrl    string             `json:"image_url" bson:"image_url"`
	Gender      string             `json:"gender,omitempty" bson:"gender,omitempty"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	//CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type PasswordResponse struct {
	CurrentPassword    string `json:"current_password,omitempty" bson:"current_password,omitempty"`
	NewPassword        string `json:"new_password,omitempty" bson:"new_password,omitempty"`
	ConfirmNewPassword string `json:"confirm_new_password,omitempty" bson:"confirm_new_password,omitempty"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		Gender:      user.Gender,
		ImageUrl:    user.ImageUrl,
		PhoneNumber: user.PhoneNumber,
		//CreatedAt:   user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
