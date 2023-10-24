package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpInput struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	PhoneNumber     string    `json:"phoneNumber" bson:"phoneNumber" binding:"required"`
	ImageUrl        string    `json:"imageUrl" bson:"imageUrl" `
	Role            string    `json:"role" bson:"role"`
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
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	ImageUrl    string             `json:"imageUrl" bson:"imageUrl"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Role        string             `json:"role,omitempty" bson:"role,omitempty"`
	ImageUrl    string             `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	PhoneNumber string             `json:"phoneNumber,omitempty" bson:"PhoneNumber,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserRole struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Role string             `json:"role,omitempty" bson:"role"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		ImageUrl:    user.ImageUrl,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
