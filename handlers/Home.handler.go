package handlers

import (
	"0xKowalski1/cluster-web-viewer/templates"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (handler *HomeHandler) GetHome(c echo.Context) error {
	return Render(c, 200, templates.HomePage())
}
