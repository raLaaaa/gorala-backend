package main

import (
	"github.com/raLaaaa/gorala/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	a := controllers.AuthController{}
	t := controllers.TaskController{}

	authAPIGroup := e.Group("/api/v1")
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	config := middleware.JWTConfig{
		Claims:     &controllers.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	authAPIGroup.Use(middleware.JWTWithConfig(config))

	e.POST("/login", a.Login)
	e.POST("/register", a.Register)

	authAPIGroup.GET("/checklogin", a.CheckLogin)

	authAPIGroup.GET("/tasks", t.GetAllTasks)
	authAPIGroup.GET("/tasks/:date", t.GetTasksForDate)
	authAPIGroup.GET("/tasks/:date/:range", t.GetTasksForDateInRange)
	authAPIGroup.GET("/tasks/today", t.GetTasksForToday)
	authAPIGroup.POST("/tasks/add", t.CreateTask)
	authAPIGroup.PUT("/tasks/edit/:id", t.UpdateTask)
	authAPIGroup.DELETE("/tasks/delete/:id", t.DeleteTask)

	e.Logger.Fatal(e.Start(":1323"))
}
