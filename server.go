package main

import (
	"github.com/raLaaaa/gorala/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	a := controllers.AuthController{}
	u := controllers.UserController{}
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

	// Routes
	e.POST("/login", a.Login)
	authAPIGroup.GET("/checklogin", a.CheckLogin)

	authAPIGroup.GET("/tasks", t.GetAllTasks)
	authAPIGroup.GET("/tasks/:date", t.GetTasksForDate)
	authAPIGroup.GET("/tasks/today", t.GetTasksForToday)
	authAPIGroup.GET("/me", u.GetUser)
	authAPIGroup.POST("/tasks", t.CreateTask)
	authAPIGroup.PUT("/tasks/:id", t.UpdateTask)
	authAPIGroup.DELETE("/tasks/:id", t.DeleteTask)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
