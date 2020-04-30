package repository

import (
	"auth/domain/entity"
)

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	UpdateUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
	GetUserByEmail(*entity.User) (*entity.User, map[string]string)
	GetUserByExternalId(*entity.User) (*entity.User, map[string]string)
}
