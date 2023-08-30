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

// @Summary Create new feature
// @Description Creates new feature in the database with provided slug.
// @Description If "percent" value was provided, also binds new feature to given % of users
// @Tags feature
// @Accept json
// @Param feature body model.NewFeatureRequest true "feature"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.Error
// @Router /v1/feature [post]
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

// @Summary Delete feature
// @Description Deletes feature with provided slug from database.
// @Description All connecte user to feature bindings will also be deleted, but saved to the history to retrieve later.
// @Tags feature
// @Accept json
// @Param feature body model.DeleteFeatureRequest true "feature"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.Error
// @Router /v1/feature [delete]
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

// @Summary Add features to user
// @Description Adds provided features to user.
// @Description The execution will fail if one of the features doesn't exist.
// @Description You may also provide expiration date for each feature individuall in format "2023-08-29T23:01:00Z".
// @Description Additions will be recorded to history.
// @Tags feature
// @Accept json
// @Param feature body model.AddFeaturesToUserRequest true "User with features to add"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.Error
// @Router /v1/feature/features [post]
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

// @Summary Delete features from user
// @Description Deletes feature with provided slug from user, but doesn't delete the feature itself.
// @Description Deletions will be recorded in history.
// @Tags feature
// @Accept json
// @Param feature body model.DeleteFeaturesFromUser true "User with features to delete"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.Error
// @Router /v1/feature/features [delete]
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
