package route

import (
	"peramalan-stok-be/src/delivery/api/handler"

	"github.com/labstack/echo/v4"
)

// PingRoute ...
func (t *NewRoute) HistoryRoute(g *echo.Group) {
	h := handler.HistoryHandler{
		Response: t.Response,
		DB:       t.DB,
	}

	g.GET("/sales-order", h.GetHistorySalesOrder)
	g.GET("/sales-order-monthly", h.GetHistorySalesOrderMonthly)
}
