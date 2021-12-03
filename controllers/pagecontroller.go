package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

type PageController struct{}

func (p *PageController) ShowMainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "main", "")
}

func (a *PageController) ShowPrivacyPage(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy", "")
}
