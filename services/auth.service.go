package services

import "grcp-api-client-mongo/models"

type AuthService interface {
	LoginUser(*models.LoginUser) (*models.DBResponseUser, error)
	RegisterUser(*models.RegiserUser) (*models.DBResponseUser, error)
}
