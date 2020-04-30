package interfaces

import (
	"auth/application"
	"auth/domain/entity"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
)

//Users struct defines the dependencies that will be used
type Users struct {
	us application.UserAppInterface
}

//Users constructor
func NewUsers(us application.UserAppInterface) *Users {
	return &Users{
		us: us,
	}
}

func (s *Users) SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	newUser, err := s.us.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetUsers(c *gin.Context) {
	users := entity.Users{} //customize user
	var err error
	//us, err = application.UserApp.GetUsers()
	users, err = s.us.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.us.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}

func (s *Users) AskForResetPassword(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}

	u, err := s.us.GetUserByEmail(&user)
	if err != nil {
		go SendResetPasswordEmail(u.Email)
		c.JSON(http.StatusOK, gin.H{})
	}
	c.JSON(http.StatusBadRequest, err)
}

func SendResetPasswordEmail(email string) {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "882dbe391f324d", "601d4054d42720", "smtp.mailtrap.io")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"qohatpp@gmail.com"}
	msg := []byte("To: qohatpp@gmail.com\r\n" +
		"Subject: Why are you not using Mailtrap yet?\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	err := smtp.SendMail("smtp.mailtrap.io:25", auth, "piotr@mailtrap.io", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
