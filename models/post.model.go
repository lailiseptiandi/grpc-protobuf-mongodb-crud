package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePostRequest struct {
	Title     string    `json:"title" bson:"title" binding:"required"`
	Content   string    `json:"content" bson:"content" binding:"required"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty"`
	User      string    `json:"user" bson:"user" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DBPost struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" binding:"required"`
	Content   string             `json:"content" bson:"content" binding:"required"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty"`
	User      string             `json:"user" bson:"user" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdatePost struct {
	Title     string    `json:"title" bson:"title" binding:"required"`
	Content   string    `json:"content" bson:"content" binding:"required"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty"`
	User      string    `json:"user" bson:"user" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
