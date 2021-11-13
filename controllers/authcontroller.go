package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/raLaaaa/gorala/models"
	"github.com/raLaaaa/gorala/services"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

type JwtCustomClaims struct {
	Name string
	ID   uint
	jwt.StandardClaims
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthController) Login(c echo.Context) error {
	userLoginDTO := new(UserLoginDTO)

	if err := c.Bind(userLoginDTO); err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	dbService := services.DatabaseService{}
	user, err := dbService.FindByEmail(userLoginDTO.Username)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Throws unauthorized error
	if userLoginDTO.Username != user.Email || !checkPasswordHash(userLoginDTO.Password, user.Password) {
		return echo.ErrUnauthorized
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user.Email,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"email": user.Email,
		"id":    user.ID,
		"token": t,
	})
}

func (a *AuthController) Register(c echo.Context) error {
	userLoginDTO := new(UserLoginDTO)

	if err := c.Bind(userLoginDTO); err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	hashedPW, err := hashPassword(userLoginDTO.Password)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := models.User{
		Email:    userLoginDTO.Username,
		Password: hashedPW,
		AllTasks: []models.Task{},
	}

	dbService := services.DatabaseService{}
	err = dbService.CreateUser(&user)

	print(err)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (a *AuthController) CheckLogin(c echo.Context) error {
	return c.String(http.StatusOK, "Success")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
