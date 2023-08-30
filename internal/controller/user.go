package controller

import (
	"avito_project/internal/model"
	"avito_project/internal/service"
	"errors"
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

// @Summary Create new user
// @Description Creates new user in the database with provided id.
// @Tags user
// @Accept json
// @Param feature body model.UserRequest true "user"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.Error
// @Router /v1/user [post]
func (r *userRoutes) NewUser(c echo.Context) error {
	var u model.User

	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	if len(u.ID) == 0 {
		return c.JSON(400, ErrorJson(errors.New("Empty user id")))
	}

	err = r.userService.AddUser(c.Request().Context(), &u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}

	return c.NoContent(200)
}

// @Summary Delete user
// @Description Deletes user from database.
// @Description All user-feature relations will also be deleted, but saved to history.
// @Tags user
// @Accept json
// @Param feature body model.UserRequest true "user"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.Error
// @Router /v1/user [delete]
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

// @Summary Get user
// @Description Gets user from database with all of his currently active features.
// @Tags user
// @Accept json
// @Param feature body model.UserRequest true "user"
// @Produce json
// @Success 200 {object} model.UserWithFeaturesResponse
// @Failure 400 {object} controller.Error
// @Router /v1/user [get]
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
