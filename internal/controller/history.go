package controller

import (
	"avito_project/internal/service"
	"github.com/labstack/echo/v4"
)

type historyRoutes struct {
	historyService service.History
}

func newHistoryRoutes(g *echo.Group, service service.History) {
	r := &historyRoutes{historyService: service}

	g.GET("/export", r.export)
}

func (r *historyRoutes) export(c echo.Context) error {
	return nil
}
