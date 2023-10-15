package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleUser string

const (
	UserRole  RoleUser = "user"
	AdminRole RoleUser = "admin"
)

type RegiserUser struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"password_confirm" bson:"password_confirm" binding:"required"`
	Role            string    `json:"role" bson:"role"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type LoginUser struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type DBResponseUser struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *DBResponseUser) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type FormatterLoginRegisterUser struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Email        string             `json:"email" bson:"email"`
	Role         string             `json:"role" bson:"role"`
	Verified     bool               `json:"verified" bson:"verified"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}

func FormatterLoginRegister(user *DBResponseUser, token string, refreshToken string) FormatterLoginRegisterUser {
	formatter := FormatterLoginRegisterUser{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		Verified:     user.Verified,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	return formatter
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
