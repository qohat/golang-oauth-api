package application

import (
	"auth/domain/entity"
	"auth/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

//UserApp implements the UserAppInterface
var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	UpdateUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
	GetUserByEmail(*entity.User) (*entity.User, map[string]string)
	GetUserByExternalId(*entity.User) (*entity.User, map[string]string)
}

func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) UpdateUser(user *entity.User) (*entity.User, map[string]string) {
	return u.us.UpdateUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entity.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers() ([]entity.User, error) {
	return u.us.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}

func (u *userApp) GetUserByEmail(user *entity.User) (*entity.User, map[string]string) {
	return u.us.GetUserByEmail(user)
}

func (u *userApp) GetUserByExternalId(user *entity.User) (*entity.User, map[string]string) {
	return u.us.GetUserByExternalId(user)
}
