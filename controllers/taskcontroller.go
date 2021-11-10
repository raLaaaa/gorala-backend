package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/raLaaaa/gorala/models"
	"github.com/raLaaaa/gorala/services"
	"github.com/raLaaaa/gorala/utilities"

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

	tu := utilities.TimeUtil{}
	roundedDate := tu.RoundDate(taskDTO.ExecutionDate)

	task := models.Task{
		Description:   taskDTO.Description,
		ExecutionDate: roundedDate,
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

	tu := utilities.TimeUtil{}
	roundedDate := tu.RoundDate(taskDTO.ExecutionDate)
	taskFromDB.ExecutionDate = roundedDate

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

func (t *TaskController) GetTasksForToday(c echo.Context) error {

	dbService := services.DatabaseService{}

	user := c.Get("user").(*jwtv3.Token)
	claims := user.Claims.(*JwtCustomClaims)

	today := time.Now().UTC()
	tu := utilities.TimeUtil{}
	today = tu.RoundDate(today)

	tasks, err := dbService.FindAllTasksOfDateByUserID(claims.ID, today)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	return c.JSON(http.StatusOK, tasks)
}

func (t *TaskController) GetTasksForDate(c echo.Context) error {

	dbService := services.DatabaseService{}
	slicedParams := strings.Split(c.Param("date"), "-")

	if len(slicedParams) != 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Date need 3 params (eg. 01-01-2021)")
	}

	user := c.Get("user").(*jwtv3.Token)
	claims := user.Claims.(*JwtCustomClaims)

	year, errYear := strconv.Atoi(slicedParams[2])
	month, errMonth := strconv.Atoi(slicedParams[1])
	day, errDay := strconv.Atoi(slicedParams[0])

	if errDay != nil || errMonth != nil || errYear != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Date")
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().UTC().Location())

	tasks, err := dbService.FindAllTasksOfDateByUserID(claims.ID, date)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	return c.JSON(http.StatusOK, tasks)
}
