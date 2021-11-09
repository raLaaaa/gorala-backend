package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/raLaaaa/gorala/services"
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
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	dbService := services.DatabaseService{}
	user, err := dbService.FindByEmail(userLoginDTO.Username)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Throws unauthorized error
	if userLoginDTO.Username != user.Email || userLoginDTO.Password != user.Password {
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

func (a *AuthController) CheckLogin(c echo.Context) error {
	return c.String(http.StatusOK, "Success")
}
