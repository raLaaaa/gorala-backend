package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

type UserController struct{}

type UserDTO struct {
	Description string `json:"description"`
}

func (u *UserController) CreateUser(c echo.Context) error {

	return c.JSON(http.StatusCreated, nil)
}

func (u *UserController) GetUser(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}

func (u *UserController) UpdateUser(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}

func (u *UserController) DeleteUser(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func (u *UserController) GetAllUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
