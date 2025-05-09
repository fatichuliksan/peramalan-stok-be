package route

import (
	"peramalan-stok-be/src/delivery/api/handler"

	"github.com/labstack/echo/v4"
)

// PingRoute ...
func (t *NewRoute) MainRoute(g *echo.Group) {
	h := handler.MainHandler{
		Response: t.Response,
		DB:       t.DB,
	}
	g.GET("/warehouses", h.GetWarehouses)
	g.GET("/items", h.GetItems)
}
