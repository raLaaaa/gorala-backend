package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/raLaaaa/gorala/models"
	"github.com/raLaaaa/gorala/services"

	jwtv3 "github.com/dgrijalva/jwt-go"
)

type TaskController struct{}

type TaskDTO struct {
	Description string `json:"description"`
}

func (t *TaskController) CreateTask(c echo.Context) error {

	user := c.Get("user").(*jwtv3.Token)
	claims := user.Claims.(*JwtCustomClaims)

	taskDTO := new(TaskDTO)
	if err := c.Bind(taskDTO); err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	task := models.Task{
		Description: taskDTO.Description,
		UserID:      claims.ID,
	}

	dbService := services.DatabaseService{}
	dbService.CreateTask(&task)

	return c.JSON(http.StatusCreated, task)
}

func (t *TaskController) UpdateTask(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}

func (t *TaskController) DeleteTask(c echo.Context) error {

	taskID, err := strconv.ParseUint(c.QueryParam("id"), 10, 32)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	dbService := services.DatabaseService{}

	dbService.DeleteTask(taskID)

	return c.NoContent(http.StatusNoContent)
}

func (t *TaskController) GetAllTasks(c echo.Context) error {

	dbService := services.DatabaseService{}

	user := c.Get("user").(*jwtv3.Token)
	claims := user.Claims.(*JwtCustomClaims)

	tasks, err := dbService.FindAllTasksByUserID(claims.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	return c.JSON(http.StatusOK, tasks)
}
