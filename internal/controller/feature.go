package controller

import (
	"avito_project/internal/model"
	"avito_project/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
)

type featureRoutes struct {
	featureService service.Feature
}

func newFeatureRoutes(g *echo.Group, featureService service.Feature) {
	r := &featureRoutes{featureService: featureService}

	g.POST("", r.NewFeature)
	g.DELETE("", r.DeleteFeature)
	g.POST("/features", r.AddFeaturesToUser)
	g.DELETE("/features", r.DeleteFeatures)
}

func (r *featureRoutes) NewFeature(c echo.Context) error {
	var f model.Feature
	err := c.Bind(&f)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	if f.Percent == nil {
		err = r.featureService.AddFeature(c.Request().Context(), &f)

		if err != nil {
			return c.JSON(400, ErrorJson(err))
		}
	} else {
		if *f.Percent < 0 || *f.Percent > 100 {
			return c.JSON(400, ErrorJson(errors.New("invalid percent value")))
		}
		err = r.featureService.AddFeatureWithPercent(c.Request().Context(), f)
		if err != nil {
			return c.JSON(400, ErrorJson(err))
		}
	}

	return c.NoContent(200)
}

func (r *featureRoutes) DeleteFeature(c echo.Context) error {
	var f model.Feature
	err := c.Bind(&f)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	err = r.featureService.DeleteFeature(c.Request().Context(), &f)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	return c.NoContent(200)
}

func (r *featureRoutes) AddFeaturesToUser(c echo.Context) error {
	var u model.User

	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	err = r.featureService.AddFeaturesToUser(c.Request().Context(), &u, u.Features)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}

	return c.NoContent(200)
}
func (r *featureRoutes) DeleteFeatures(c echo.Context) error {
	var u model.User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	err = r.featureService.DeleteFeaturesFromUser(c.Request().Context(), u.Features, &u)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	return c.NoContent(200)
}
