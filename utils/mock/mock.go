package mock

import (
	"auth/domain/entity"
	"auth/infrastructure/auth"
	"mime/multipart"
	"net/http"
)

//UserAppInterface is a mock user app interface
type UserAppInterface struct {
	SaveUserFn                  func(*entity.User) (*entity.User, map[string]string)
	GetUsersFn                  func() ([]entity.User, error)
	GetUserFn                   func(uint64) (*entity.User, error)
	GetUserByEmailAndPasswordFn func(*entity.User) (*entity.User, map[string]string)
	GetUserByEmailFn            func(*entity.User) (*entity.User, map[string]string)
	GetUserByExternalIdFn       func(*entity.User) (*entity.User, map[string]string)
	AskForResetPasswordFn       func(*entity.User) (*entity.User, map[string]string)
}

//SaveUser calls the SaveUserFn
func (u *UserAppInterface) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.SaveUserFn(user)
}

//GetUsersFn calls the GetUsers
func (u *UserAppInterface) GetUsers() ([]entity.User, error) {
	return u.GetUsersFn()
}

//GetUserFn calls the GetUser
func (u *UserAppInterface) GetUser(userId uint64) (*entity.User, error) {
	return u.GetUserFn(userId)
}

//GetUserByEmailAndPasswordFn calls the GetUserByEmailAndPassword
func (u *UserAppInterface) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.GetUserByEmailAndPasswordFn(user)
}

func (u *UserAppInterface) GetUserByEmail(user *entity.User) (*entity.User, map[string]string) {
	return u.GetUserByEmailFn(user)
}

func (u *UserAppInterface) GetUserByExternalId(user *entity.User) (*entity.User, map[string]string) {
	return u.GetUserByExternalIdFn(user)
}

func (u *UserAppInterface) AskForResetPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.AskForResetPasswordFn(user)
}

//AuthInterface is a mock auth interface
type AuthInterface struct {
	CreateAuthFn    func(uint64, *auth.TokenDetails) error
	FetchAuthFn     func(string) (uint64, error)
	DeleteRefreshFn func(string) error
	DeleteTokensFn  func(*auth.AccessDetails) error
}

func (f *AuthInterface) DeleteRefresh(refreshUuid string) error {
	return f.DeleteRefreshFn(refreshUuid)
}
func (f *AuthInterface) DeleteTokens(authD *auth.AccessDetails) error {
	return f.DeleteTokensFn(authD)
}
func (f *AuthInterface) FetchAuth(uuid string) (uint64, error) {
	return f.FetchAuthFn(uuid)
}
func (f *AuthInterface) CreateAuth(userId uint64, authD *auth.TokenDetails) error {
	return f.CreateAuthFn(userId, authD)
}

//TokenInterface is a mock token interface
type TokenInterface struct {
	CreateTokenFn          func(userId uint64) (*auth.TokenDetails, error)
	ExtractTokenMetadataFn func(*http.Request) (*auth.AccessDetails, error)
}

func (f *TokenInterface) CreateToken(userid uint64) (*auth.TokenDetails, error) {
	return f.CreateTokenFn(userid)
}
func (f *TokenInterface) ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	return f.ExtractTokenMetadataFn(r)
}

type UploadFileInterface struct {
	UploadFileFn func(file *multipart.FileHeader) (string, error)
}

func (up *UploadFileInterface) UploadFile(file *multipart.FileHeader) (string, error) {
	return up.UploadFileFn(file)
}
