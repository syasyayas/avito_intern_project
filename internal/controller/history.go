package controller

import (
	"avito_project/internal/service"
	"github.com/labstack/echo/v4"
	"time"
)

var magic_date = time.Date(2050, time.December, 31, 0, 0, 0, 0, time.Local)

type historyRoutes struct {
	historyService service.History
}

func newHistoryRoutes(g *echo.Group, service service.History) {
	r := &historyRoutes{historyService: service}

	g.GET("/export", r.export)
}

// @Summary Get history
// @Description Retrieves history between to provided dates and returns csv file uri.
// @Tags history
// @Accept json
// @Param feature body model.HistoryRequest true "dates"
// @Produce json
// @Success 200 {object} model.HistoryResponse
// @Failure 400 {object} controller.Error
// @Router /v1/history/export [get]
func (r *historyRoutes) export(c echo.Context) error {
	var req historyRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}
	if req.Before.IsZero() {
		req.Before = magic_date
	}
	url, err := r.historyService.Export(c.Request().Context(), req.After, req.Before)
	if err != nil {
		return c.JSON(400, ErrorJson(err))
	}

	var res = struct {
		Url string `json:"url"`
	}{
		url,
	}
	return c.JSON(200, res)
}

type historyRequest struct {
	After  time.Time `json:"after"`
	Before time.Time `json:"before"`
}
