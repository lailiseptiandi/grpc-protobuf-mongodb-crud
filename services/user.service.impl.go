package services

import (
	"context"
	"grcp-api-client-mongo/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserService(collection *mongo.Collection, ctx context.Context) *UserServiceImpl {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(userID string) (*models.DBResponseUser, error) {
	objId, _ := primitive.ObjectIDFromHex(userID)

	var user *models.DBResponseUser

	query := bson.M{"_id": objId}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponseUser{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponseUser, error) {

	var user *models.DBResponseUser

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponseUser{}, err
		}
		return nil, err
	}

	return user, nil
}
