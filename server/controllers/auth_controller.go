package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"server/database"
	"server/dto"
	"server/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gogithub "github.com/google/go-github/v55/github"
	githuboauth "golang.org/x/oauth2/github"
	"gorm.io/gorm"
)


var jwtSecret = []byte("your_secret_key")

var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/api/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

var githubOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/api/auth/github/callback",
	Scopes:       []string{"user:email"},
	Endpoint:     githuboauth.Endpoint,
}

func Register(c echo.Context) error {
	var data dto.RegisterDTO
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	user := models.User{
		Email:    data.Email,
		Password: data.Password,
		Name:     data.Email,
		OAuth:    false,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "User already exists"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered"})
}

func Login(c echo.Context) error {
	var data dto.LoginDTO
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	var user models.User
	if err := database.DB.Where("email = ? AND password = ?", data.Email, data.Password).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	token := generateJWT(user)
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func GoogleLogin(c echo.Context) error {
	url := googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

func GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing code"})
	}

	tok, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to exchange token"})
	}

	client := googleOAuthConfig.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get user info"})
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Decode failed"})
	}

	var user models.User
	err = database.DB.Where("email = ?", userInfo.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = models.User{
			Email: userInfo.Email,
			Name:  userInfo.Name,
			OAuth: true,
		}
		database.DB.Create(&user)
	}

	token := generateJWT(user)
	return c.Redirect(http.StatusFound, "http://localhost:3000/?token="+token)
}

func GitHubLogin(c echo.Context) error {
	url := githubOAuthConfig.AuthCodeURL("state-github")
	return c.Redirect(http.StatusFound, url)
}

func GitHubCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing code"})
	}

	tok, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Token exchange failed"})
	}

	client := gogithub.NewClient(githubOAuthConfig.Client(context.Background(), tok))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user"})
	}

	email := user.GetEmail()
	if email == "" {
		emails, _, err := client.Users.ListEmails(context.Background(), nil)
		if err == nil && len(emails) > 0 {
			email = emails[0].GetEmail()
		}
	}
	if email == "" {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Email not found"})
	}

	var u models.User
	if err := database.DB.Where("email = ?", email).First(&u).Error; err == gorm.ErrRecordNotFound {
		u = models.User{
			Email: email,
			Name:  user.GetName(),
			OAuth: true,
		}
		database.DB.Create(&u)
	}

	token := generateJWT(u)
	return c.Redirect(http.StatusFound, "http://localhost:3000/?token="+token)
}

func generateJWT(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}
