package services

import "grcp-api-client-mongo/models"

type UserService interface {
	FindUserById(string) (*models.DBResponseUser, error)
	FindUserByEmail(string) (*models.DBResponseUser, error)
}
