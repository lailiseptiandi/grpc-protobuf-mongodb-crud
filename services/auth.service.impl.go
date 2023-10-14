package services

import (
	"context"
	"errors"
	"grcp-api-client-mongo/models"
	"grcp-api-client-mongo/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthService(collection *mongo.Collection, ctx context.Context) *AuthServiceImpl {
	return &AuthServiceImpl{collection, ctx}
}

func (uc *AuthServiceImpl) RegisterUser(user *models.RegiserUser) (*models.DBResponseUser, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = string(models.UserRole)

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	res, err := uc.collection.InsertOne(uc.ctx, &user)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}
	if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	var newUser *models.DBResponseUser
	query := bson.M{"_id": res.InsertedID}

	err = uc.collection.FindOne(uc.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) LoginUser(*models.LoginUser) (*models.DBResponseUser, error) {
	return nil, nil
}
