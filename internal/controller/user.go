package controller

import (
	"avito_project/internal/model"
	"avito_project/internal/service"
	"github.com/labstack/echo/v4"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(g *echo.Group, userService service.User) {
	r := &userRoutes{userService: userService}

	g.POST("", r.NewUser)
	g.DELETE("", r.DeleteUser)
	g.GET("", r.GetUser)
}

func (r *userRoutes) NewUser(c echo.Context) error {
	var u model.User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	err = r.userService.AddUser(c.Request().Context(), &u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	return c.NoContent(200)
}

func (r *userRoutes) DeleteUser(c echo.Context) error {
	var u model.User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	err = r.userService.DeleteUser(c.Request().Context(), &u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	return c.NoContent(200)
}

func (r *userRoutes) GetUser(c echo.Context) error {
	var u model.User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	res, err := r.userService.GetUserWithFeatures(c.Request().Context(), &u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	if res.Features == nil {
		res.Features = []model.Feature{}
	}
	return c.JSON(200, res)
}
