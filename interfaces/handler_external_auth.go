package interfaces

import (
	"auth/application"
	"auth/domain/entity"
	"auth/infrastructure/auth"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ExtAuth struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

type ExtAuthResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
	Id      string `json:"id"`
}

var fbOauthConf = &oauth2.Config{
	/*ClientID:     os.Getenv("FB_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("FB_OAUTH_SECRET"),
	RedirectURL:  os.Getenv("FB_OAUTH_REDIRECT_URI"),*/
	ClientID:     "2724842771078378",
	ClientSecret: "e3347c0c4059333e3550bd3a07165ce3",
	RedirectURL:  "http://localhost:9191/auth/facebook/callback",
	Scopes:       []string{"public_profile", "email"},
	Endpoint:     facebook.Endpoint,
}

var googleOauthConfig = &oauth2.Config{
	/*ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),*/
	ClientID:     "127999312447-9rnltevu3v1n1fmij6alat0r08b224af.apps.googleusercontent.com",
	ClientSecret: "-wIVprfkEjCfw9e5njbmYe8j",
	RedirectURL:  "http://localhost:9191/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

var oauthLocalRedirect = "http://dev.app.dochealth.co/auth/verification?key="

//NewExtAuth constructor
func NewExtAuth(uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *ExtAuth {
	return &ExtAuth{
		us: uApp,
		rd: rd,
		tk: tk,
	}
}

func (extAuth *ExtAuth) HandleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(generateStateOauthCookie(c.Writer))
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func (extAuth *ExtAuth) HandleFacebookLogin(c *gin.Context) {
	url := fbOauthConf.AuthCodeURL(generateStateOauthCookie(c.Writer))
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func (extAuth *ExtAuth) HandleGoogleCallback(c *gin.Context) {
	// Read oauthState from Cookie
	state := c.Request.FormValue("state")
	oauthState, _ := c.Request.Cookie("oauthstate")

	if state != oauthState.Value {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthState, state)
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}

	code := c.Request.FormValue("code")
	data, err := getUserDataFromGoogle(code)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}

	extAuthData := toExtAuthResponse(data)

	newUser := &entity.User{
		FirstName:  extAuthData.Name,
		LastName:   "",
		Email:      extAuthData.Email,
		ExternalId: extAuthData.Id,
		Password:   "",
	}
	// GetOrCreate User in your db.
	u, userErr := extAuth.us.GetUserByEmail(newUser)

	if userErr != nil && userErr["no_user"] == "user not found" {
		newUser, err := extAuth.us.SaveUser(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		http.Redirect(c.Writer, c.Request, oauthLocalRedirect+newUser.ExternalId, http.StatusTemporaryRedirect)
	} else if userErr != nil && userErr["no_user"] != "user not found" {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	}

	if u != nil {
		log.Printf("user exist")
		u.FirstName = extAuthData.Name
		u.ExternalId = extAuthData.Id
		extAuth.us.UpdateUser(u)
		http.Redirect(c.Writer, c.Request, oauthLocalRedirect+u.ExternalId, http.StatusTemporaryRedirect)
	}
}

func (extAuth *ExtAuth) HandleFacebookCallback(c *gin.Context) {
	state := c.Request.FormValue("state")
	oauthState, _ := c.Request.Cookie("oauthstate")
	if c.Request.FormValue("state") != oauthState.Value {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthState, state)
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}

	code := c.Request.FormValue("code")

	token, err := fbOauthConf.Exchange(context.TODO(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken) + "&fields=id,name,email")
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("parseResponseBody: %s\n", string(data))

	extAuthData := toExtAuthResponse(data)

	newUser := &entity.User{
		FirstName:  extAuthData.Name,
		LastName:   "",
		Email:      extAuthData.Email,
		ExternalId: extAuthData.Id,
		Password:   "",
	}
	// GetOrCreate User in your db.
	u, userErr := extAuth.us.GetUserByEmail(newUser)

	if userErr != nil && userErr["no_user"] == "user not found" {
		newUser, err := extAuth.us.SaveUser(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		http.Redirect(c.Writer, c.Request, oauthLocalRedirect+newUser.ExternalId, http.StatusTemporaryRedirect)
	} else if userErr != nil && userErr["no_user"] != "user not found" {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	}

	if u != nil {
		log.Printf("user exist")
		u.FirstName = extAuthData.Name
		u.ExternalId = extAuthData.Id
		extAuth.us.UpdateUser(u)
		http.Redirect(c.Writer, c.Request, oauthLocalRedirect+u.ExternalId, http.StatusTemporaryRedirect)
	}
}

func (extAuth *ExtAuth) HanldeExternalClientId(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	usr, err := extAuth.us.GetUserByExternalId(&user)
	if usr == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "token_expired",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
		return
	}

	var tokenErr = map[string]string{}
	ts, tErr := extAuth.tk.CreateToken(usr.ID)
	if tErr != nil {
		tokenErr["token_error"] = tErr.Error()
		c.JSON(http.StatusUnprocessableEntity, tErr.Error())
		return
	}
	saveErr := extAuth.rd.CreateAuth(usr.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr.Error())
		return
	}
	userData := make(map[string]interface{})
	userData["accessToken"] = ts.AccessToken
	userData["refreshToken"] = ts.RefreshToken
	userData["id"] = usr.ID
	userData["firstName"] = usr.FirstName
	userData["lastName"] = usr.LastName

	c.JSON(http.StatusOK, userData)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func toExtAuthResponse(data []byte) ExtAuthResponse {
	var extAuthResponse ExtAuthResponse
	json.Unmarshal(data, &extAuthResponse)
	return extAuthResponse
}
