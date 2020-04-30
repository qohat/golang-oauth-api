package interfaces

import "auth/utils/mock"

var (
	userApp    mock.UserAppInterface
	fakeUpload mock.UploadFileInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface

	s  = NewUsers(&userApp)                               //We use all mocked data here
	au = NewAuthenticate(&userApp, &fakeAuth, &fakeToken) //We use all mocked data here

)
