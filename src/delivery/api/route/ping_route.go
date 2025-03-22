package route

import (
	"peramalan-stok-be/src/delivery/api/handler"

	"github.com/labstack/echo/v4"
)

// PingRoute ...
func (t *NewRoute) PingRoute(g *echo.Group) {
	h := handler.PingHandler{
		Response: t.Response,
	}
	g.GET("", h.Get)
}
