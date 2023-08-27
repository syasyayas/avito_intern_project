package controller

import (
	"avito_project/internal/service"
	"github.com/labstack/echo/v4"
)

type featureRoutes struct {
	featureService service.Feature
}

func newFeatureRoutes(g *echo.Group, featureService service.Feature) {
	r := &featureRoutes{featureService: featureService}

	g.POST("/feature", r.NewFeature)
	g.DELETE("/feature", r.DeleteFeature)
	g.POST("/features", r.AddFeaturesToUser)
}

func (r *featureRoutes) NewFeature(c echo.Context) error {
	return nil
}

func (r *featureRoutes) DeleteFeature(c echo.Context) error {
	return nil
}

func (r *featureRoutes) AddFeaturesToUser(c echo.Context) error {
	return nil
}
