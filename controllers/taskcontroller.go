package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/raLaaaa/gorala/models"
	"github.com/raLaaaa/gorala/services"

	jwtv3 "github.com/dgrijalva/jwt-go"
)

type TaskController struct{}

type TaskDTO struct {
	Description   string    `json:"description"`
	ExecutionDate time.Time `json:"executionDate"`
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
		Description:   taskDTO.Description,
		ExecutionDate: taskDTO.ExecutionDate,
		UserID:        claims.ID,
	}

	dbService := services.DatabaseService{}
	dbService.CreateTask(&task)

	return c.JSON(http.StatusCreated, task)
}

func (t *TaskController) UpdateTask(c echo.Context) error {

	i, err := strconv.ParseUint(c.Param("id"), 10, 64)

	dbService := services.DatabaseService{}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Task not found (ID Error)")
	}

	taskFromDB, err := dbService.FindTaskByID(i)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Task not found (Database Error)")
	}

	taskDTO := new(TaskDTO)
	if err := c.Bind(taskDTO); err != nil {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	taskFromDB.Description = taskDTO.Description
	taskFromDB.ExecutionDate = taskDTO.ExecutionDate

	dbService.UpdateTask(taskFromDB)

	return c.JSON(http.StatusCreated, taskFromDB)
}

func (t *TaskController) DeleteTask(c echo.Context) error {

	i, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	dbService := services.DatabaseService{}
	dbService.DeleteTask(i)

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

func (t *TaskController) GetTasksForDate(c echo.Context) error {

	dbService := services.DatabaseService{}

	user := c.Get("user").(*jwtv3.Token)
	claims := user.Claims.(*JwtCustomClaims)

	tasks, err := dbService.FindAllTasksByUserID(claims.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	return c.JSON(http.StatusOK, tasks)
}
