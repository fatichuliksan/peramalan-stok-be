package route

import (
	"peramalan-stok-be/src/delivery/api/handler"

	"github.com/labstack/echo/v4"
)

// PingRoute ...
func (t *NewRoute) ForcastingRoute(g *echo.Group) {
	h := handler.ForecastingHandler{
		Response: t.Response,
		DB:       t.DB,
	}
	g.POST("/generate", h.PostGenerate)
	g.POST("/history", h.PostHistory)
	g.GET("/history", h.GetHistory)
	g.DELETE("/history", h.DeleteHistory)
}
