package controller

import (
	"avito_project/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *echo.Echo, services service.Services) {
	handler.Use(middleware.Recover())

	handler.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(200)
	})

	// TODO swagger

	v1 := handler.Group("/v1")
	{
		newFeatureRoutes(v1.Group("/feature"), services.Feature)
		newUserRoutes(v1.Group("/user"), services.User)
		newHistoryRoutes(v1.Group("/history"), services.History)
	}
}
