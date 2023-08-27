package controller

import (
	"avito_project/internal/service"
	"github.com/labstack/echo/v4"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(g *echo.Group, userService service.User) {
	r := &userRoutes{userService: userService}

	g.POST("/user", r.NewUser)
	g.DELETE("/user", r.DeleteUser)
	g.GET("/user", r.GetUser)
}

func (r *userRoutes) NewUser(c echo.Context) error {
	return nil
}

func (r *userRoutes) DeleteUser(c echo.Context) error {
	return nil
}

func (r *userRoutes) GetUser(c echo.Context) error {
	return nil
}
