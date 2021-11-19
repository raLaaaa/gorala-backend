package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/raLaaaa/gorala/models"
	"github.com/raLaaaa/gorala/services"
	"github.com/raLaaaa/gorala/utilities"
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

type PasswordRequestDTO struct {
	Email string `json:"username"`
}

type PasswordResetDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
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

	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	// Throws unauthorized error
	if userLoginDTO.Username != user.Email || !checkPasswordHash(userLoginDTO.Password, user.Password) {
		return echo.ErrUnauthorized
	}

	if !user.Accepted {
		return echo.NewHTTPError(http.StatusUnauthorized, "Your account has not been confirmed yet")
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
		Accepted: false,
		AllTasks: []models.Task{},
	}

	dbService := services.DatabaseService{}
	err = dbService.CreateUser(&user)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User exists already")
	}

	confirmToken, err := dbService.CreateConfirmationToken(&user)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	e := utilities.EmailUtil{}
	e.SendRegistryConfirmation(*confirmToken, user)

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func (a *AuthController) ConfirmRegistration(c echo.Context) error {
	token := c.Param("token")

	dbService := services.DatabaseService{}
	success, err := dbService.ResolveConfirmationToken(token)

	if success && err == nil {
		return c.String(http.StatusOK, "Success! You can now login")
	} else {
		return c.String(http.StatusBadRequest, "Invalid Token")
	}
}

func (a *AuthController) ShowResetPasswordPage(c echo.Context) error {
	token := c.Param("token")

	dbService := services.DatabaseService{}

	resetToken, err := dbService.FindResetToken(token)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid token")
	}

	if resetToken.Activated {
		return echo.NewHTTPError(http.StatusBadRequest, "Token already used")
	}

	now := time.Now()
	diff := now.Sub(resetToken.CreatedAt)

	if diff.Hours() > 24 {
		resetToken.Activated = true
		dbService.UpdateResetToken(resetToken)
		return echo.NewHTTPError(http.StatusBadRequest, "Token already used")
	}

	return c.Render(http.StatusOK, "reset_password", resetToken.Token)
}

func (a *AuthController) ShowPrivacyPage(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy", "")
}

func (a *AuthController) DoPasswordReset(c echo.Context) error {
	token := c.Param("token")

	passwordResetDTO := new(PasswordResetDTO)

	if err := c.Bind(passwordResetDTO); err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	dbService := services.DatabaseService{}
	resetToken, err := dbService.FindResetToken(token)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid token")
	}

	if resetToken.Activated {
		return c.String(http.StatusBadRequest, "Token already used")
	}

	user, err := dbService.FindUserByID(resetToken.UserID)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !checkPasswordHash(passwordResetDTO.OldPassword, user.Password) {
		return c.String(http.StatusBadRequest, "Your old password was incorrect")
	} else {
		hashedPW, err := hashPassword(passwordResetDTO.NewPassword)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		user.Password = hashedPW
		dbService.UpdateUser(user)
		dbService.ResolveResetToken(token)
		return c.String(http.StatusOK, "Success")
	}
}

func (a *AuthController) RequestPasswordReset(c echo.Context) error {
	passwordRequestDTO := new(PasswordRequestDTO)

	if err := c.Bind(passwordRequestDTO); err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	if len(passwordRequestDTO.Email) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user information")
	}

	dbService := services.DatabaseService{}
	user, err := dbService.FindByEmail(passwordRequestDTO.Email)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	alreadyExistingTokens, err := dbService.FindAllResetTokensByUserID(user.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, v := range alreadyExistingTokens {
		v.Activated = true
		dbService.UpdateResetToken(&v)
	}

	token, err := dbService.CreateResetToken(user)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	e := utilities.EmailUtil{}
	e.SendResetPassword(*token, *user)

	return c.String(http.StatusOK, "Success")
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
